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

package controllers

import (
	goctx "context"
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	bootstrapapi "k8s.io/cluster-bootstrap/token/api"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/controllers/remote"
	clusterutilv1 "sigs.k8s.io/cluster-api/util"
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
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	vmwarecontext "sigs.k8s.io/cluster-api-provider-vsphere/pkg/context/vmware"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/record"
)

const (
	clusterNotReadyRequeueTime     = time.Minute * 2
	ServiceDiscoveryControllerName = "servicediscovery-controller"

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
func AddServiceDiscoveryControllerToManager(ctx *context.ControllerManagerContext, mgr manager.Manager, tracker *remote.ClusterCacheTracker, options controller.Options) error {
	var (
		controllerNameShort = ServiceDiscoveryControllerName
		controllerNameLong  = fmt.Sprintf("%s/%s/%s", ctx.Namespace, ctx.Name, ServiceDiscoveryControllerName)
	)
	controllerContext := &context.ControllerContext{
		ControllerManagerContext: ctx,
		Name:                     controllerNameShort,
		Recorder:                 record.New(mgr.GetEventRecorderFor(controllerNameLong)),
		Logger:                   ctx.Logger.WithName(controllerNameShort),
	}
	r := serviceDiscoveryReconciler{
		ControllerContext:         controllerContext,
		remoteClusterCacheTracker: tracker,
	}

	configMapCache, err := cache.New(mgr.GetConfig(), cache.Options{
		Scheme: mgr.GetScheme(),
		Mapper: mgr.GetRESTMapper(),
		// TODO: Reintroduce the cache sync period
		// Resync:    ctx.SyncPeriod,
		Namespaces: []string{metav1.NamespacePublic},
	})
	if err != nil {
		return errors.Wrapf(err, "failed to create configmap cache")
	}
	if err := mgr.Add(configMapCache); err != nil {
		return errors.Wrapf(err, "failed to start configmap cache")
	}
	src := source.Kind(configMapCache, &corev1.ConfigMap{})

	return ctrl.NewControllerManagedBy(mgr).For(&vmwarev1.VSphereCluster{}).
		WithOptions(options).
		Watches(
			&corev1.Service{},
			handler.EnqueueRequestsFromMapFunc(r.serviceToClusters),
		).
		WatchesRawSource(
			src,
			handler.EnqueueRequestsFromMapFunc(r.configMapToClusters),
		).
		// watch the CAPI cluster
		Watches(
			&clusterv1.Cluster{},
			handler.EnqueueRequestForOwner(
				mgr.GetScheme(), mgr.GetRESTMapper(),
				&vmwarev1.VSphereCluster{},
				handler.OnlyControllerOwner(),
			)).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(ctrl.LoggerFrom(ctx), ctx.WatchFilterValue)).
		Complete(r)
}

type serviceDiscoveryReconciler struct {
	*context.ControllerContext

	remoteClusterCacheTracker *remote.ClusterCacheTracker
}

func (r serviceDiscoveryReconciler) Reconcile(_ goctx.Context, req reconcile.Request) (_ reconcile.Result, reterr error) {
	logger := r.Logger.WithName(req.Namespace).WithName(req.Name)
	logger.V(4).Info("Starting Reconcile")

	// Get the vspherecluster for this request.
	vsphereCluster := &vmwarev1.VSphereCluster{}
	if err := r.Client.Get(r, req.NamespacedName, vsphereCluster); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("Cluster not found, won't reconcile", "key", req.NamespacedName)
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// Create the patch helper.
	patchHelper, err := patch.NewHelper(vsphereCluster, r.Client)
	if err != nil {
		return reconcile.Result{}, errors.Wrapf(
			err,
			"failed to init patch helper for %s %s/%s",
			vsphereCluster.GroupVersionKind(),
			vsphereCluster.Namespace,
			vsphereCluster.Name)
	}

	// Create the cluster context for this request.
	clusterContext := &vmwarecontext.ClusterContext{
		ControllerContext: r.ControllerContext,
		VSphereCluster:    vsphereCluster,
		Logger:            logger,
		PatchHelper:       patchHelper,
	}

	// Always issue a patch when exiting this function so changes to the
	// resource are patched back to the API server.
	defer func() {
		if err := clusterContext.Patch(); err != nil {
			if reterr == nil {
				reterr = err
			} else {
				clusterContext.Logger.Error(err, "patch failed", "cluster", clusterContext.String())
			}
		}
	}()

	// This type of controller doesn't care about delete events.
	if !vsphereCluster.DeletionTimestamp.IsZero() {
		return reconcile.Result{}, nil
	}

	cluster, err := clusterutilv1.GetClusterFromMetadata(r, r.Client, vsphereCluster.ObjectMeta)
	if err != nil {
		logger.Info("unable to get capi cluster from vsphereCluster", "err", err)
		return reconcile.Result{RequeueAfter: clusterNotReadyRequeueTime}, nil
	}

	// We cannot proceed until we are able to access the target cluster. Until
	// then just return a no-op and wait for the next sync.
	guestClient, err := r.remoteClusterCacheTracker.GetClient(clusterContext, client.ObjectKeyFromObject(cluster))
	if err != nil {
		if errors.Is(err, remote.ErrClusterLocked) {
			r.Logger.V(5).Info("Requeuing because another worker has the lock on the ClusterCacheTracker")
			return ctrl.Result{Requeue: true}, nil
		}
		logger.Info("The control plane is not ready yet", "err", err)
		return reconcile.Result{RequeueAfter: clusterNotReadyRequeueTime}, nil
	}

	// Defer to the Reconciler for reconciling a non-delete event.
	return r.ReconcileNormal(&vmwarecontext.GuestClusterContext{
		ClusterContext: clusterContext,
		GuestClient:    guestClient,
	})
}

func (r serviceDiscoveryReconciler) ReconcileNormal(ctx *vmwarecontext.GuestClusterContext) (reconcile.Result, error) {
	ctx.Logger.V(4).Info("Reconciling Service Discovery", "cluster", ctx.VSphereCluster.Name)
	if err := r.reconcileSupervisorHeadlessService(ctx); err != nil {
		conditions.MarkFalse(ctx.VSphereCluster, vmwarev1.ServiceDiscoveryReadyCondition, vmwarev1.SupervisorHeadlessServiceSetupFailedReason,
			clusterv1.ConditionSeverityWarning, err.Error())
		return reconcile.Result{}, errors.Wrapf(err, "failed to configure supervisor headless service for %v", ctx.VSphereCluster)
	}

	return reconcile.Result{}, nil
}

// Setup a local k8s service in the target cluster that proxies to the Supervisor Cluster API Server. The add-ons are
// dependent on this local service to connect to the Supervisor Cluster.
func (r serviceDiscoveryReconciler) reconcileSupervisorHeadlessService(ctx *vmwarecontext.GuestClusterContext) error {
	// Create the headless service to the supervisor api server on the target cluster.
	supervisorPort := vmwarev1.SupervisorAPIServerPort
	svc := NewSupervisorHeadlessService(vmwarev1.SupervisorHeadlessSvcPort, supervisorPort)
	if err := ctx.GuestClient.Create(ctx, svc); err != nil && !apierrors.IsAlreadyExists(err) {
		return errors.Wrapf(err, "cannot create k8s service %s/%s in ", svc.Namespace, svc.Name)
	}

	supervisorHost, err := GetSupervisorAPIServerAddress(ctx.ClusterContext)
	if err != nil {
		// Note: We have watches on the LB Svc (VIP) & the cluster-info configmap (FIP). There is no need to return an error to keep
		// re-trying.
		conditions.MarkFalse(ctx.VSphereCluster, vmwarev1.ServiceDiscoveryReadyCondition, vmwarev1.SupervisorHeadlessServiceSetupFailedReason,
			clusterv1.ConditionSeverityWarning, err.Error())
		return nil
	}

	ctx.Logger.Info("Discovered supervisor apiserver address", "host", supervisorHost, "port", supervisorPort)
	// CreateOrPatch the newEndpoints with the discovered supervisor api server address
	newEndpoints := NewSupervisorHeadlessServiceEndpoints(
		supervisorHost,
		supervisorPort,
	)
	endpointsKey := types.NamespacedName{
		Namespace: newEndpoints.Namespace,
		Name:      newEndpoints.Name,
	}
	endpoints := &corev1.Endpoints{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: newEndpoints.Namespace,
			Name:      newEndpoints.Name,
		},
	}
	result, err := controllerutil.CreateOrPatch(
		ctx,
		ctx.GuestClient,
		endpoints,
		func() error {
			endpoints.Subsets = newEndpoints.Subsets
			return nil
		})
	if err != nil {
		return errors.Wrapf(
			err,
			"cannot create k8s service endpoints %s",
			endpointsKey,
		)
	}

	endpointsSubsetsStr := fmt.Sprintf("%+v", endpoints.Subsets)

	switch result {
	case controllerutil.OperationResultNone:
		ctx.Logger.Info(
			"no update required for k8s service endpoints",
			"endpointsKey",
			endpointsKey,
			"endpointsSubsets",
			endpointsSubsetsStr,
		)
	case controllerutil.OperationResultCreated:
		ctx.Logger.Info(
			"created k8s service endpoints",
			"endpointsKey",
			endpointsKey,
			"endpointsSubsets",
			endpointsSubsetsStr,
		)
	case controllerutil.OperationResultUpdated:
		ctx.Logger.Info(
			"updated k8s service endpoints",
			"endpointsKey",
			endpointsKey,
			"endpointsSubsets",
			endpointsSubsetsStr,
		)
	default:
		ctx.Logger.Error(
			fmt.Errorf(
				"unexpected result during createOrPatch k8s service endpoints",
			),
			"endpointsKey",
			endpointsKey,
			"endpointsSubsets",
			endpointsSubsetsStr,
			"result",
			result,
		)
	}

	conditions.MarkTrue(ctx.VSphereCluster, vmwarev1.ServiceDiscoveryReadyCondition)
	return nil
}

func GetSupervisorAPIServerAddress(ctx *vmwarecontext.ClusterContext) (string, error) {
	// Discover the supervisor api server address
	// 1. Check if a k8s service "kube-system/kube-apiserver-lb-svc" is available, if so, fetch the loadbalancer IP.
	// 2. If not, get the Supervisor Cluster Management Network Floating IP (FIP) from the cluster-info configmap. This is
	// to support non-NSX-T development usecases only. If we are unable to find the cluster-info configmap for some reason,
	// we log the error.
	supervisorHost, err := GetSupervisorAPIServerVIP(ctx.Client)
	if err != nil {
		ctx.Logger.Info("Unable to discover supervisor apiserver virtual ip, fallback to floating ip", "reason", err.Error())
		supervisorHost, err = GetSupervisorAPIServerFIP(ctx.Client)
		if err != nil {
			ctx.Logger.Error(err, "Unable to discover supervisor apiserver address")
			return "", errors.Wrapf(err, "Unable to discover supervisor apiserver address")
		}
	}
	return supervisorHost, nil
}

func NewSupervisorHeadlessService(port, targetPort int) *corev1.Service {
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

func NewSupervisorHeadlessServiceEndpoints(targetHost string, targetPort int) *corev1.Endpoints {
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

func GetSupervisorAPIServerVIP(client client.Client) (string, error) {
	svc := &corev1.Service{}
	svcKey := types.NamespacedName{Name: vmwarev1.SupervisorLoadBalancerSvcName, Namespace: vmwarev1.SupervisorLoadBalancerSvcNamespace}
	if err := client.Get(goctx.Background(), svcKey, svc); err != nil {
		return "", errors.Wrapf(err, "unable to get supervisor loadbalancer svc %s", svcKey)
	}
	if len(svc.Status.LoadBalancer.Ingress) > 0 {
		ingress := svc.Status.LoadBalancer.Ingress[0]
		if ipAddr := ingress.IP; ipAddr != "" {
			return ipAddr, nil
		}
		return ingress.Hostname, nil
	}
	return "", errors.Errorf("no VIP found in the supervisor loadbalancer svc %s", svcKey)
}

func GetSupervisorAPIServerFIP(client client.Client) (string, error) {
	urlString, err := getSupervisorAPIServerURLWithFIP(client)
	if err != nil {
		return "", errors.Wrap(err, "unable to get supervisor url")
	}
	urlVal, err := url.Parse(urlString)
	if err != nil {
		return "", errors.Wrapf(err, "unable to parse supervisor url from %s", urlString)
	}
	host := urlVal.Hostname()
	if host == "" {
		return "", errors.Errorf("unable to get supervisor host from url %s", urlVal)
	}
	return host, nil
}

func getSupervisorAPIServerURLWithFIP(client client.Client) (string, error) {
	cm := &corev1.ConfigMap{}
	cmKey := types.NamespacedName{Name: bootstrapapi.ConfigMapClusterInfo, Namespace: metav1.NamespacePublic}
	if err := client.Get(goctx.Background(), cmKey, cm); err != nil {
		return "", err
	}
	kubeconfig, err := tryParseClusterInfoFromConfigMap(cm)
	if err != nil {
		return "", err
	}
	clusterConfig := getClusterFromKubeConfig(kubeconfig)
	if clusterConfig != nil {
		return clusterConfig.Server, nil
	}
	return "", errors.Errorf("unable to get cluster from kubeconfig in ConfigMap %s/%s", cm.Namespace, cm.Name)
}

// tryParseClusterInfoFromConfigMap tries to parse a kubeconfig file from a ConfigMap key.
func tryParseClusterInfoFromConfigMap(cm *corev1.ConfigMap) (*clientcmdapi.Config, error) {
	kubeConfigString, ok := cm.Data[bootstrapapi.KubeConfigKey]
	if !ok || len(kubeConfigString) == 0 {
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
func (r serviceDiscoveryReconciler) serviceToClusters(ctx goctx.Context, o client.Object) []reconcile.Request {
	if o.GetNamespace() != vmwarev1.SupervisorLoadBalancerSvcNamespace || o.GetName() != vmwarev1.SupervisorLoadBalancerSvcName {
		return nil
	}
	return allClustersRequests(ctx, r.Client)
}

// configMapToClusters is a mapper function used to enqueue reconcile.Requests
// It watches for cluster-info configmaps for the supervisor api-server.
func (r serviceDiscoveryReconciler) configMapToClusters(ctx goctx.Context, o client.Object) []reconcile.Request {
	if o.GetNamespace() != metav1.NamespacePublic || o.GetName() != bootstrapapi.ConfigMapClusterInfo {
		return nil
	}
	return allClustersRequests(ctx, r.Client)
}

func allClustersRequests(ctx goctx.Context, c client.Client) []reconcile.Request {
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
