/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package vmware

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"reflect"
	"time"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/tools/record"
	bootstrapapi "k8s.io/cluster-bootstrap/token/api"
	"k8s.io/klog/v2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/controllers/clustercache"
	clusterutilv1 "sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/vmware/v1beta1"
	capvcontext "sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	vmwarecontext "sigs.k8s.io/cluster-api-provider-vsphere/pkg/context/vmware"
)

const (
	clusterNotReadyRequeueTime = time.Minute * 2

	supervisorLoadBalancerSvcNamespace = "kube-system"
	supervisorLoadBalancerSvcName      = "kube-apiserver-lb-svc"
	supervisorAPIServerPort            = 6443

	supervisorHeadlessSvcNamespace = "default"
	supervisorHeadlessSvcName      = "supervisor"
)

// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=services/status,verbs=get
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=configmaps/status,verbs=get

// AddServiceDiscoveryControllerToManager adds the ServiceDiscovery controller to the provided manager.
func AddServiceDiscoveryControllerToManager(ctx context.Context, controllerManagerCtx *capvcontext.ControllerManagerContext, mgr manager.Manager, clusterCache clustercache.ClusterCache, options controller.Options) error {
	r := &serviceDiscoveryReconciler{
		Client:       controllerManagerCtx.Client,
		Recorder:     mgr.GetEventRecorderFor("servicediscovery/vspherecluster-controller"),
		clusterCache: clusterCache,
	}
	predicateLog := ctrl.LoggerFrom(ctx).WithValues("controller", "servicediscovery/vspherecluster")

	configMapCache, err := cache.New(mgr.GetConfig(), cache.Options{
		Scheme: mgr.GetScheme(),
		Mapper: mgr.GetRESTMapper(),
		// TODO: Reintroduce the cache sync period
		// Resync:    ctx.SyncPeriod,
		DefaultNamespaces: map[string]cache.Config{metav1.NamespacePublic: {}},
	})
	if err != nil {
		return errors.Wrapf(err, "failed to create ConfigMap cache")
	}
	if err := mgr.Add(configMapCache); err != nil {
		return errors.Wrapf(err, "failed to add ConfigMap cache")
	}
	return ctrl.NewControllerManagedBy(mgr).For(&vmwarev1.VSphereCluster{}).
		Named("servicediscovery/vspherecluster").
		WithOptions(options).
		Watches(
			&corev1.Service{},
			handler.EnqueueRequestsFromMapFunc(r.serviceToClusters),
		).
		WatchesRawSource(
			source.Kind(
				configMapCache,
				&corev1.ConfigMap{},
				handler.TypedEnqueueRequestsFromMapFunc(r.configMapToClusters),
			),
		).
		// watch the CAPI cluster
		Watches(
			&clusterv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(clusterToSupervisorVSphereClusterFunc(r.Client)),
		).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(mgr.GetScheme(), predicateLog, controllerManagerCtx.WatchFilterValue)).
		WatchesRawSource(r.clusterCache.GetClusterSource("servicediscovery/vspherecluster", clusterToSupervisorVSphereClusterFunc(r.Client))).
		Complete(r)
}

func clusterToSupervisorVSphereClusterFunc(ctrlclient client.Client) func(ctx context.Context, obj client.Object) []reconcile.Request {
	return func(ctx context.Context, obj client.Object) []reconcile.Request {
		gvk := vmwarev1.GroupVersion.WithKind(reflect.TypeOf(&vmwarev1.VSphereCluster{}).Elem().Name())
		requests := clusterutilv1.ClusterToInfrastructureMapFunc(ctx, gvk, ctrlclient, &vmwarev1.VSphereCluster{})(ctx, obj)
		if len(requests) == 0 {
			return nil
		}

		log := ctrl.LoggerFrom(ctx, "Cluster", klog.KObj(obj), "VSphereCluster", klog.KRef(requests[0].Namespace, requests[0].Name))
		ctx = ctrl.LoggerInto(ctx, log)

		c := &vmwarev1.VSphereCluster{}
		if err := ctrlclient.Get(ctx, requests[0].NamespacedName, c); err != nil {
			log.V(4).Error(err, "Failed to get VSphereCluster")
			return nil
		}

		if annotations.IsExternallyManaged(c) {
			log.V(6).Info("VSphereCluster is externally managed, will not attempt to map resource")
			return nil
		}
		return requests
	}
}

type serviceDiscoveryReconciler struct {
	Client   client.Client
	Recorder record.EventRecorder

	clusterCache clustercache.ClusterCache
}

func (r *serviceDiscoveryReconciler) Reconcile(ctx context.Context, req reconcile.Request) (_ reconcile.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)

	// Get the vspherecluster for this request.
	vsphereCluster := &vmwarev1.VSphereCluster{}
	// Note: VSphereCluster doesn't have to be added to the logger as controller-runtime
	// already adds the reconciled object (which is VSphereCluster).
	if err := r.Client.Get(ctx, req.NamespacedName, vsphereCluster); err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	cluster, err := clusterutilv1.GetClusterFromMetadata(ctx, r.Client, vsphereCluster.ObjectMeta)
	if err != nil {
		return reconcile.Result{RequeueAfter: clusterNotReadyRequeueTime}, errors.Wrapf(err, "failed to get Cluster from VSphereCluster")
	}
	log = log.WithValues("Cluster", klog.KObj(cluster))
	ctx = ctrl.LoggerInto(ctx, log)

	if annotations.IsPaused(cluster, vsphereCluster) {
		log.Info("Reconciliation is paused for this object")
		return ctrl.Result{}, nil
	}

	// Create the patch helper.
	patchHelper, err := patch.NewHelper(vsphereCluster, r.Client)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Create the cluster context for this request.
	clusterContext := &vmwarecontext.ClusterContext{
		Cluster:        cluster,
		VSphereCluster: vsphereCluster,
		PatchHelper:    patchHelper,
	}

	// Always issue a patch when exiting this function so changes to the
	// resource are patched back to the API server.
	defer func() {
		if err := clusterContext.Patch(ctx); err != nil {
			reterr = kerrors.NewAggregate([]error{reterr, err})
		}
	}()

	// This type of controller doesn't care about delete events.
	if !vsphereCluster.DeletionTimestamp.IsZero() {
		return reconcile.Result{}, nil
	}

	// We cannot proceed until we are able to access the target cluster. Until
	// then just return a no-op and wait for the next sync.
	guestClient, err := r.clusterCache.GetClient(ctx, client.ObjectKeyFromObject(cluster))
	if err != nil {
		if errors.Is(err, clustercache.ErrClusterNotConnected) {
			log.V(5).Info("Requeuing because connection to the workload cluster is down")
			return ctrl.Result{RequeueAfter: time.Minute}, nil
		}
		log.Error(err, "The control plane is not ready yet")
		return reconcile.Result{RequeueAfter: clusterNotReadyRequeueTime}, nil
	}

	// Defer to the Reconciler for reconciling a non-delete event.
	return reconcile.Result{}, r.reconcileNormal(ctx, &vmwarecontext.GuestClusterContext{
		ClusterContext: clusterContext,
		GuestClient:    guestClient,
	})
}

func (r *serviceDiscoveryReconciler) reconcileNormal(ctx context.Context, guestClusterCtx *vmwarecontext.GuestClusterContext) error {
	if err := r.reconcileSupervisorHeadlessService(ctx, guestClusterCtx); err != nil {
		conditions.MarkFalse(guestClusterCtx.VSphereCluster, vmwarev1.ServiceDiscoveryReadyCondition, vmwarev1.SupervisorHeadlessServiceSetupFailedReason,
			clusterv1.ConditionSeverityWarning, err.Error())
		return errors.Wrapf(err, "failed to reconcile supervisor headless Service")
	}

	return nil
}

// reconcileSupervisorHeadlessService sets up a local k8s service in the workload cluster that
// proxies to the Supervisor Cluster API Server. The add-ons are depend on this local service
// to connect to the Supervisor Cluster.
func (r *serviceDiscoveryReconciler) reconcileSupervisorHeadlessService(ctx context.Context, guestClusterCtx *vmwarecontext.GuestClusterContext) error {
	log := ctrl.LoggerFrom(ctx)

	// Create the headless service to the supervisor api server on the target cluster.
	supervisorPort := vmwarev1.SupervisorAPIServerPort
	svc := newSupervisorHeadlessService(vmwarev1.SupervisorHeadlessSvcPort, supervisorPort)

	log = log.WithValues("Service", klog.KObj(svc))
	ctx = ctrl.LoggerInto(ctx, log)

	testObj := svc.DeepCopyObject().(client.Object)
	if err := guestClusterCtx.GuestClient.Get(ctx, client.ObjectKeyFromObject(svc), testObj); err != nil {
		if !apierrors.IsNotFound(err) {
			return errors.Wrapf(err, "failed to check if Service %s already exists", klog.KObj(svc))
		}

		// If Secret doesn't exist, create it
		log.Info("Creating supervisor headless Service")
		if err := guestClusterCtx.GuestClient.Create(ctx, svc); err != nil && !apierrors.IsAlreadyExists(err) {
			return errors.Wrapf(err, "failed to create supervisor headless Service")
		}
	}

	supervisorHost, err := r.getSupervisorAPIServerAddress(ctx)
	if err != nil {
		// Note: We have watches on the LB Svc (VIP) & the cluster-info configmap (FIP).
		// There is no need to return an error to keep re-trying.
		conditions.MarkFalse(guestClusterCtx.VSphereCluster, vmwarev1.ServiceDiscoveryReadyCondition, vmwarev1.SupervisorHeadlessServiceSetupFailedReason,
			clusterv1.ConditionSeverityWarning, err.Error())
		return nil
	}

	log.Info("Discovered supervisor API server endpoint", "host", supervisorHost, "port", supervisorPort)
	// CreateOrPatch the newEndpoints with the discovered supervisor api server address
	newEndpoints := newSupervisorHeadlessServiceEndpoints(
		supervisorHost,
		supervisorPort,
	)
	endpointsKey := types.NamespacedName{
		Namespace: newEndpoints.Namespace,
		Name:      newEndpoints.Name,
	}
	log = log.WithValues("Endpoints", klog.KRef(endpointsKey.Namespace, endpointsKey.Name))
	ctx = ctrl.LoggerInto(ctx, log)

	endpoints := &corev1.Endpoints{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: newEndpoints.Namespace,
			Name:      newEndpoints.Name,
		},
	}

	result, err := controllerutil.CreateOrPatch(
		ctx,
		guestClusterCtx.GuestClient,
		endpoints,
		func() error {
			endpoints.Subsets = newEndpoints.Subsets
			return nil
		})
	if err != nil {
		return errors.Wrapf(err, "failed to create or patch service Endpoints")
	}

	endpointsSubsetsStr := fmt.Sprintf("%+v", endpoints.Subsets)

	switch result {
	case controllerutil.OperationResultNone:
		log.V(3).Info("No update required for service Endpoints", "endpointsSubsets", endpointsSubsetsStr)
	case controllerutil.OperationResultCreated:
		log.Info("Created service Endpoints", "endpointsSubsets", endpointsSubsetsStr)
	case controllerutil.OperationResultUpdated:
		log.Info("Updated service Endpoints", "endpointsSubsets", endpointsSubsetsStr)
	default:
		log.Error(nil, "Unexpected result during createOrPatch service Endpoints", "endpointsSubsets", endpointsSubsetsStr, "operationResult", result)
	}

	conditions.MarkTrue(guestClusterCtx.VSphereCluster, vmwarev1.ServiceDiscoveryReadyCondition)
	return nil
}

func (r *serviceDiscoveryReconciler) getSupervisorAPIServerAddress(ctx context.Context) (string, error) {
	// Discover the supervisor api server address
	// 1. Check if a k8s service "kube-system/kube-apiserver-lb-svc" is available, if so, fetch the loadbalancer IP.
	// 2. If not, get the Supervisor Cluster Management Network Floating IP (FIP) from the cluster-info configmap. This is
	// to support non-NSX-T development use cases only. If we are unable to find the cluster-info configmap for some reason,
	// we log the error.
	supervisorHost, vipErr := getSupervisorAPIServerVIP(ctx, r.Client)
	if vipErr != nil {
		var fipErr error
		supervisorHost, fipErr = getSupervisorAPIServerFIP(ctx, r.Client)
		if fipErr != nil {
			return "", errors.Wrapf(kerrors.NewAggregate([]error{vipErr, fipErr}), "Failed to discover supervisor API server endpoint")
		}
	}
	return supervisorHost, nil
}

// newSupervisorHeadlessService returns a new Supervisor headless service.
func newSupervisorHeadlessService(port, targetPort int) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      vmwarev1.SupervisorHeadlessSvcName,
			Namespace: vmwarev1.SupervisorHeadlessSvcNamespace,
		},
		Spec: corev1.ServiceSpec{
			// Note: This will be a headless service with no selectors. The endpoints will be manually created.
			ClusterIP: corev1.ClusterIPNone,
			Ports: []corev1.ServicePort{
				{
					Port:       int32(port),
					TargetPort: intstr.FromInt(targetPort),
				},
			},
		},
	}
}

// newSupervisorHeadlessServiceEndpoints returns Kubernetes Endpoints for the supervisor apiserver address.
func newSupervisorHeadlessServiceEndpoints(targetHost string, targetPort int) *corev1.Endpoints {
	var endpointAddr corev1.EndpointAddress
	if ip := net.ParseIP(targetHost); ip != nil {
		endpointAddr.IP = ip.String()
	} else {
		endpointAddr.Hostname = targetHost
	}
	return &corev1.Endpoints{
		ObjectMeta: metav1.ObjectMeta{
			Name:      vmwarev1.SupervisorHeadlessSvcName,
			Namespace: vmwarev1.SupervisorHeadlessSvcNamespace,
		},
		Subsets: []corev1.EndpointSubset{
			{
				Addresses: []corev1.EndpointAddress{
					endpointAddr,
				},
				Ports: []corev1.EndpointPort{
					{
						Port: int32(targetPort),
					},
				},
			},
		},
	}
}

// getSupervisorAPIServerVIP finds the load balancer IP of the Supervisor APIServer.
func getSupervisorAPIServerVIP(ctx context.Context, client client.Client) (string, error) {
	svc := &corev1.Service{}
	svcKey := types.NamespacedName{Name: vmwarev1.SupervisorLoadBalancerSvcName, Namespace: vmwarev1.SupervisorLoadBalancerSvcNamespace}
	if err := client.Get(ctx, svcKey, svc); err != nil {
		return "", errors.Wrapf(err, "unable to get supervisor loadbalancer Service %s", svcKey)
	}
	if len(svc.Status.LoadBalancer.Ingress) > 0 {
		ingress := svc.Status.LoadBalancer.Ingress[0]
		if ipAddr := ingress.IP; ipAddr != "" {
			return ipAddr, nil
		}
		return ingress.Hostname, nil
	}
	return "", errors.Errorf("no VIP found in the supervisor loadbalancer Service %s", svcKey)
}

// getSupervisorAPIServerFIP finds the floating ip of the Supervisor APIServer.
func getSupervisorAPIServerFIP(ctx context.Context, client client.Client) (string, error) {
	urlString, err := getSupervisorAPIServerURLWithFIP(ctx, client)
	if err != nil {
		return "", errors.Wrap(err, "unable to get supervisor URL")
	}
	urlVal, err := url.Parse(urlString)
	if err != nil {
		return "", errors.Wrapf(err, "unable to parse supervisor URL from %s", urlString)
	}
	host := urlVal.Hostname()
	if host == "" {
		return "", errors.Errorf("unable to get supervisor host from URL %s", urlVal)
	}
	return host, nil
}

func getSupervisorAPIServerURLWithFIP(ctx context.Context, client client.Client) (string, error) {
	cm := &corev1.ConfigMap{}
	cmKey := types.NamespacedName{Name: bootstrapapi.ConfigMapClusterInfo, Namespace: metav1.NamespacePublic}
	if err := client.Get(ctx, cmKey, cm); err != nil {
		return "", errors.Wrapf(err, "unable to get ConfigMap %s", cmKey)
	}
	kubeconfig, err := tryParseClusterInfoFromConfigMap(cm)
	if err != nil {
		return "", err
	}
	clusterConfig := getClusterFromKubeConfig(kubeconfig)
	if clusterConfig != nil {
		return clusterConfig.Server, nil
	}
	return "", errors.Errorf("unable to get cluster from kubeconfig in ConfigMap %s", cmKey)
}

// tryParseClusterInfoFromConfigMap tries to parse a kubeconfig file from a ConfigMap key.
func tryParseClusterInfoFromConfigMap(cm *corev1.ConfigMap) (*clientcmdapi.Config, error) {
	kubeConfigString, ok := cm.Data[bootstrapapi.KubeConfigKey]
	if !ok || kubeConfigString == "" {
		return nil, errors.Errorf("no %s key in ConfigMap %s/%s", bootstrapapi.KubeConfigKey, cm.Namespace, cm.Name)
	}
	parsedKubeConfig, err := clientcmd.Load([]byte(kubeConfigString))
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't parse the kubeconfig file in the ConfigMap %s/%s", cm.Namespace, cm.Name)
	}
	return parsedKubeConfig, nil
}

// GetClusterFromKubeConfig returns the default Cluster of the specified KubeConfig.
func getClusterFromKubeConfig(config *clientcmdapi.Config) *clientcmdapi.Cluster {
	// If there is an unnamed cluster object, use it
	if config.Clusters[""] != nil {
		return config.Clusters[""]
	}
	if config.Contexts[config.CurrentContext] != nil {
		return config.Clusters[config.Contexts[config.CurrentContext].Cluster]
	}
	return nil
}

// serviceToClusters is a mapper function used to enqueue reconcile.Requests
// It watches for Service objects of type LoadBalancer for the supervisor api-server.
func (r *serviceDiscoveryReconciler) serviceToClusters(ctx context.Context, o client.Object) []reconcile.Request {
	if o.GetNamespace() != vmwarev1.SupervisorLoadBalancerSvcNamespace || o.GetName() != vmwarev1.SupervisorLoadBalancerSvcName {
		return nil
	}
	return allClustersRequests(ctx, r.Client)
}

// configMapToClusters is a mapper function used to enqueue reconcile.Requests
// It watches for cluster-info configmaps for the supervisor api-server.
func (r *serviceDiscoveryReconciler) configMapToClusters(ctx context.Context, o *corev1.ConfigMap) []reconcile.Request {
	if o.GetNamespace() != metav1.NamespacePublic || o.GetName() != bootstrapapi.ConfigMapClusterInfo {
		return nil
	}
	return allClustersRequests(ctx, r.Client)
}

func allClustersRequests(ctx context.Context, c client.Client) []reconcile.Request {
	vsphereClusterList := &vmwarev1.VSphereClusterList{}
	if err := c.List(ctx, vsphereClusterList, &client.ListOptions{}); err != nil {
		return nil
	}

	requests := make([]reconcile.Request, 0, len(vsphereClusterList.Items))
	for _, vSphereCluster := range vsphereClusterList.Items {
		key := client.ObjectKey{
			Namespace: vSphereCluster.GetNamespace(),
			Name:      vSphereCluster.GetName(),
		}
		requests = append(requests, reconcile.Request{NamespacedName: key})
	}
	return requests
}
