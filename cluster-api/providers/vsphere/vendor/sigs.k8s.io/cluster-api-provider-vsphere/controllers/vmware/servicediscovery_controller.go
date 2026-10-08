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

	pkgerrors "github.com/pkg/errors"
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
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/controllers/clustercache"
	clusterutilv1 "sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/conditions"
	deprecatedv1beta1conditions "sigs.k8s.io/cluster-api/util/conditions/deprecated/v1beta1"
	capicontrollerutil "sigs.k8s.io/cluster-api/util/controller"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/api/supervisor/v1beta2"
	capvcontext "sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	vmwarecontext "sigs.k8s.io/cluster-api-provider-vsphere/pkg/context/vmware"
	inframanager "sigs.k8s.io/cluster-api-provider-vsphere/pkg/manager"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/util"
)

const (
	clusterNotReadyRequeueTime = time.Minute * 2

	supervisorAPIServerPort = 6443

	supervisorHeadlessSvcNamespace = "default"
	supervisorHeadlessSvcName      = "supervisor"
)

// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=services/status,verbs=get
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=configmaps/status,verbs=get

// AddServiceDiscoveryControllerToManager adds the ServiceDiscovery controller to the provided manager.
func AddServiceDiscoveryControllerToManager(ctx context.Context, controllerManagerCtx *capvcontext.ControllerManagerContext, mgr manager.Manager, clusterCache clustercache.ClusterCache, options controller.Options) error {
	networkProvider, err := inframanager.GetNetworkProvider(ctx, controllerManagerCtx.Client, controllerManagerCtx.NetworkProvider)
	if err != nil {
		return pkgerrors.Wrap(err, "failed to create network provider")
	}

	r := &serviceDiscoveryReconciler{
		Client:          controllerManagerCtx.Client,
		Recorder:        mgr.GetEventRecorderFor("servicediscovery/vspherecluster-controller"),
		NetworkProvider: networkProvider,
		clusterCache:    clusterCache,
	}
	predicateLog := ctrl.LoggerFrom(ctx).WithValues("controller", "servicediscovery/vspherecluster")

	return capicontrollerutil.NewControllerManagedBy(mgr, predicateLog).
		For(&vmwarev1.VSphereCluster{}).
		Named("servicediscovery/vspherecluster").
		WithOptions(options).
		Watches(
			&corev1.Service{},
			handler.EnqueueRequestsFromMapFunc(r.serviceToClusters),
		).
		Watches(
			&corev1.ConfigMap{},
			handler.EnqueueRequestsFromMapFunc(r.configMapToClusters),
		).
		// watch the CAPI cluster
		Watches(
			&clusterv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(clusterToSupervisorVSphereClusterFunc(r.Client)),
		).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(mgr.GetScheme(), predicateLog, controllerManagerCtx.WatchFilterValue)).
		WatchesRawSource(r.clusterCache.GetClusterSource("servicediscovery/vspherecluster", clusterToSupervisorVSphereClusterFunc(r.Client))).
		Complete(ctx, r)
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
	Client          client.Client
	Recorder        record.EventRecorder
	NetworkProvider services.NetworkProvider

	clusterCache clustercache.ClusterCache
}

func (r *serviceDiscoveryReconciler) Reconcile(ctx context.Context, req reconcile.Request) (ret reconcile.Result, reterr error) {
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
		return reconcile.Result{RequeueAfter: clusterNotReadyRequeueTime}, pkgerrors.Wrapf(err, "failed to get Cluster from VSphereCluster")
	}
	log = log.WithValues("Cluster", klog.KObj(cluster))
	ctx = ctrl.LoggerInto(ctx, log)

	// In future we might consider to surface a separate paused condition for this controller.
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
		if err := r.patch(ctx, clusterContext); err != nil {
			reterr = kerrors.NewAggregate([]error{reterr, err})
		}

		if reterr != nil {
			// Note: controller-runtime logs a warning if an error is returned in combination with
			// RequeueAfter / Requeue. Dropping RequeueAfter and Requeue here to avoid this warning
			// (while preserving Priority).
			ret.RequeueAfter = 0
			ret.Requeue = false //nolint:staticcheck // We have to handle Requeue until it is removed
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
		if pkgerrors.Is(err, clustercache.ErrClusterNotConnected) {
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

func (r *serviceDiscoveryReconciler) patch(ctx context.Context, clusterCtx *vmwarecontext.ClusterContext) error {
	// NOTE: this controller only owns the ServiceDiscoveryReady condition on the VSphereCluster object.
	return clusterCtx.PatchHelper.Patch(ctx, clusterCtx.VSphereCluster,
		patch.WithOwnedV1Beta1Conditions{Conditions: []clusterv1.ConditionType{
			vmwarev1.ServiceDiscoveryReadyV1Beta1Condition,
		}},
		patch.WithOwnedConditions{Conditions: []string{
			vmwarev1.VSphereClusterServiceDiscoveryReadyCondition,
		}},
	)
}

func (r *serviceDiscoveryReconciler) reconcileNormal(ctx context.Context, guestClusterCtx *vmwarecontext.GuestClusterContext) error {
	if err := r.reconcileSupervisorHeadlessService(ctx, guestClusterCtx); err != nil {
		deprecatedv1beta1conditions.MarkFalse(guestClusterCtx.VSphereCluster, vmwarev1.ServiceDiscoveryReadyV1Beta1Condition, vmwarev1.SupervisorHeadlessServiceSetupFailedV1Beta1Reason,
			clusterv1.ConditionSeverityWarning, "%v", err)
		conditions.Set(guestClusterCtx.VSphereCluster, metav1.Condition{
			Type:    vmwarev1.VSphereClusterServiceDiscoveryReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  vmwarev1.VSphereClusterServiceDiscoveryNotReadyReason,
			Message: err.Error(),
		})
		return pkgerrors.Wrapf(err, "failed to reconcile supervisor headless Service")
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
			return pkgerrors.Wrapf(err, "failed to check if Service %s already exists", klog.KObj(svc))
		}

		// If Secret doesn't exist, create it
		log.Info("Creating supervisor headless Service")
		if err := guestClusterCtx.GuestClient.Create(ctx, svc); err != nil && !apierrors.IsAlreadyExists(err) {
			return pkgerrors.Wrapf(err, "failed to create supervisor headless Service")
		}
	}

	var supervisorHosts []string
	var err error
	supervisorHosts, err = r.getSupervisorAPIServerAddresses(ctx, guestClusterCtx.Cluster)
	if err != nil {
		// Note: We have watches on the LB Svc (VIP) & the cluster-info configmap (FIP).
		// There is no need to return an error to keep re-trying.
		deprecatedv1beta1conditions.MarkFalse(guestClusterCtx.VSphereCluster, vmwarev1.ServiceDiscoveryReadyV1Beta1Condition, vmwarev1.SupervisorHeadlessServiceSetupFailedV1Beta1Reason,
			clusterv1.ConditionSeverityWarning, "%v", err)
		conditions.Set(guestClusterCtx.VSphereCluster, metav1.Condition{
			Type:    vmwarev1.VSphereClusterServiceDiscoveryReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  vmwarev1.VSphereClusterServiceDiscoveryNotReadyReason,
			Message: err.Error(),
		})
		return nil
	}

	log.V(5).Info("Discovered supervisor API server endpoint", "hosts", supervisorHosts, "port", supervisorPort)
	// CreateOrPatch the newEndpoints with the discovered supervisor api server address
	newEndpoints, err := newSupervisorHeadlessServiceEndpoints(
		supervisorHosts,
		supervisorPort,
	)
	if err != nil {
		// Note: We have watches on the LB Svc (VIP) & the cluster-info configmap (FIP).
		// There is no need to return an error to keep re-trying.
		deprecatedv1beta1conditions.MarkFalse(guestClusterCtx.VSphereCluster, vmwarev1.ServiceDiscoveryReadyV1Beta1Condition, vmwarev1.SupervisorHeadlessServiceSetupFailedV1Beta1Reason,
			clusterv1.ConditionSeverityWarning, "%v", err)
		conditions.Set(guestClusterCtx.VSphereCluster, metav1.Condition{
			Type:    vmwarev1.VSphereClusterServiceDiscoveryReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  vmwarev1.VSphereClusterServiceDiscoveryNotReadyReason,
			Message: err.Error(),
		})
		return nil
	}
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
		return pkgerrors.Wrapf(err, "failed to create or patch service Endpoints")
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

	deprecatedv1beta1conditions.MarkTrue(guestClusterCtx.VSphereCluster, vmwarev1.ServiceDiscoveryReadyV1Beta1Condition)
	conditions.Set(guestClusterCtx.VSphereCluster, metav1.Condition{
		Type:   vmwarev1.VSphereClusterServiceDiscoveryReadyCondition,
		Status: metav1.ConditionTrue,
		Reason: vmwarev1.VSphereClusterServiceDiscoveryReadyReason,
	})
	return nil
}

func (r *serviceDiscoveryReconciler) getSupervisorAPIServerAddresses(ctx context.Context, cluster *clusterv1.Cluster) ([]string, error) {
	if r.NetworkProvider.SupportsIPv6DualStack() {
		// 1. If dual stack IS supported, we determine the intended IP family from the cluster configuration.
		ipFamily, err := util.DetermineClusterIPFamily(cluster)
		if err != nil {
			return nil, err
		}

		// 2. For IPv6 single stack or dual-stack (when using NSX-VPC), only check VIP.
		//    No FIP fallback since NSX-VPC does not support non-load-balanced scenarios.
		if ipFamily != util.IPv4SingleStack {
			vips, err := getSupervisorAPIServerVIPs(ctx, r.Client)
			if err != nil {
				return nil, pkgerrors.Wrap(err, "failed to discover supervisor API server VIPs")
			}

			if len(vips) >= 3 {
				return nil, pkgerrors.Errorf("found too many VIPs: %v", vips)
			}

			var ipv4, ipv6 string
			for _, ip := range vips {
				parsedIP := net.ParseIP(ip)
				if parsedIP == nil {
					return nil, pkgerrors.Errorf("invalid supervisor API server VIP %q: must be an IP address", ip)
				}
				if parsedIP.To4() != nil {
					ipv4 = ip
				} else {
					ipv6 = ip
				}
			}

			switch ipFamily {
			case util.IPv6SingleStack:
				if ipv6 == "" {
					return nil, pkgerrors.Errorf("no supervisor apiserver IPv6 VIP found for IPv6 single stack cluster")
				}
				return []string{ipv6}, nil

			case util.DualStackIPv4Primary:
				var result []string
				if ipv4 != "" {
					result = append(result, ipv4)
				}
				if ipv6 != "" {
					result = append(result, ipv6)
				}
				if len(result) == 0 {
					return nil, pkgerrors.Errorf("no supervisor apiserver VIP found for dual stack cluster")
				}
				return result, nil

			case util.DualStackIPv6Primary:
				var result []string
				if ipv6 != "" {
					result = append(result, ipv6)
				}
				if ipv4 != "" {
					result = append(result, ipv4)
				}
				if len(result) == 0 {
					return nil, pkgerrors.Errorf("no supervisor apiserver VIP found for dual stack cluster")
				}
				return result, nil
			}

			return nil, pkgerrors.Errorf("unknown cluster IPFamily")
		}
	}

	// 3. Handle IPv4 single stack (either because dual stack is not supported by the network provider
	// or because the cluster is configured as IPv4).
	// This path is CIDR-agnostic when dual stack is not supported, ensuring no side effects for
	// existing clusters in older environments. It also supports FIP fallback.
	supervisorHost, vipErr := getSupervisorAPIServerVIP(ctx, r.Client)
	if vipErr != nil {
		var fipErr error
		supervisorHost, fipErr = getSupervisorAPIServerFIP(ctx, r.Client)
		if fipErr != nil {
			return nil, pkgerrors.Wrapf(kerrors.NewAggregate([]error{vipErr, fipErr}), "Failed to discover supervisor API server endpoint")
		}
	}
	return []string{supervisorHost}, nil
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
func newSupervisorHeadlessServiceEndpoints(targetHosts []string, targetPort int) (*corev1.Endpoints, error) {
	var addresses []corev1.EndpointAddress
	for _, targetHost := range targetHosts {
		ip := net.ParseIP(targetHost)
		if ip == nil {
			return nil, pkgerrors.Errorf("invalid supervisor API server endpoint %q: must be an IP address", targetHost)
		}
		addresses = append(addresses, corev1.EndpointAddress{IP: ip.String()})
	}

	return &corev1.Endpoints{
		ObjectMeta: metav1.ObjectMeta{
			Name:      vmwarev1.SupervisorHeadlessSvcName,
			Namespace: vmwarev1.SupervisorHeadlessSvcNamespace,
		},
		Subsets: []corev1.EndpointSubset{
			{
				Addresses: addresses,
				Ports: []corev1.EndpointPort{
					{
						Port: int32(targetPort),
					},
				},
			},
		},
	}, nil
}

// getSupervisorAPIServerVIP finds the load balancer IP of the Supervisor APIServer.
func getSupervisorAPIServerVIP(ctx context.Context, client client.Client) (string, error) {
	svc := &corev1.Service{}
	svcKey := types.NamespacedName{Name: vmwarev1.SupervisorLoadBalancerSvcName, Namespace: vmwarev1.SupervisorLoadBalancerSvcNamespace}
	if err := client.Get(ctx, svcKey, svc); err != nil {
		return "", pkgerrors.Wrapf(err, "unable to get supervisor loadbalancer Service %s", svcKey)
	}
	if len(svc.Status.LoadBalancer.Ingress) > 0 {
		ingress := svc.Status.LoadBalancer.Ingress[0]
		if ipAddr := ingress.IP; ipAddr != "" {
			return ipAddr, nil
		}
	}
	return "", pkgerrors.Errorf("no VIP found in the supervisor loadbalancer Service %s", svcKey)
}

// getSupervisorAPIServerVIPs finds all load balancer IPs of the Supervisor APIServer.
func getSupervisorAPIServerVIPs(ctx context.Context, client client.Client) ([]string, error) {
	svc := &corev1.Service{}
	svcKey := types.NamespacedName{Name: vmwarev1.SupervisorLoadBalancerSvcName, Namespace: vmwarev1.SupervisorLoadBalancerSvcNamespace}
	if err := client.Get(ctx, svcKey, svc); err != nil {
		return nil, pkgerrors.Wrapf(err, "unable to get supervisor loadbalancer Service %s", svcKey)
	}

	var vips []string
	for _, ingress := range svc.Status.LoadBalancer.Ingress {
		if ipAddr := ingress.IP; ipAddr != "" {
			vips = append(vips, ipAddr)
		}
	}
	if len(vips) > 0 {
		return vips, nil
	}
	return nil, pkgerrors.Errorf("no VIP found in the supervisor loadbalancer Service %s", svcKey)
}

// getSupervisorAPIServerFIP finds the floating ip of the Supervisor APIServer.
func getSupervisorAPIServerFIP(ctx context.Context, client client.Client) (string, error) {
	urlString, err := getSupervisorAPIServerURLWithFIP(ctx, client)
	if err != nil {
		return "", pkgerrors.Wrap(err, "unable to get supervisor URL")
	}
	urlVal, err := url.Parse(urlString)
	if err != nil {
		return "", pkgerrors.Wrapf(err, "unable to parse supervisor URL from %s", urlString)
	}
	host := urlVal.Hostname()
	if host == "" {
		return "", pkgerrors.Errorf("unable to get supervisor host from URL %s", urlVal)
	}
	return host, nil
}

func getSupervisorAPIServerURLWithFIP(ctx context.Context, client client.Client) (string, error) {
	cm := &corev1.ConfigMap{}
	cmKey := types.NamespacedName{Name: bootstrapapi.ConfigMapClusterInfo, Namespace: metav1.NamespacePublic}
	if err := client.Get(ctx, cmKey, cm); err != nil {
		return "", pkgerrors.Wrapf(err, "unable to get ConfigMap %s", cmKey)
	}
	kubeconfig, err := tryParseClusterInfoFromConfigMap(cm)
	if err != nil {
		return "", err
	}
	clusterConfig := getClusterFromKubeConfig(kubeconfig)
	if clusterConfig != nil {
		return clusterConfig.Server, nil
	}
	return "", pkgerrors.Errorf("unable to get cluster from kubeconfig in ConfigMap %s", cmKey)
}

// tryParseClusterInfoFromConfigMap tries to parse a kubeconfig file from a ConfigMap key.
func tryParseClusterInfoFromConfigMap(cm *corev1.ConfigMap) (*clientcmdapi.Config, error) {
	kubeConfigString, ok := cm.Data[bootstrapapi.KubeConfigKey]
	if !ok || kubeConfigString == "" {
		return nil, pkgerrors.Errorf("no %s key in ConfigMap %s/%s", bootstrapapi.KubeConfigKey, cm.Namespace, cm.Name)
	}
	parsedKubeConfig, err := clientcmd.Load([]byte(kubeConfigString))
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "couldn't parse the kubeconfig file in the ConfigMap %s/%s", cm.Namespace, cm.Name)
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
func (r *serviceDiscoveryReconciler) configMapToClusters(ctx context.Context, o client.Object) []reconcile.Request {
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
