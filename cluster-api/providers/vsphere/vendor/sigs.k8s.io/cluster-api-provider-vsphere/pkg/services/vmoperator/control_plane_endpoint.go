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

package vmoperator

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"slices"

	pkgerrors "github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	controlplanev1 "sigs.k8s.io/cluster-api/api/controlplane/kubeadm/v1beta2"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/util/conditions"
	deprecatedv1beta1conditions "sigs.k8s.io/cluster-api/util/conditions/deprecated/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/api/supervisor/v1beta2"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context/vmware"
	vmoprvhub "sigs.k8s.io/cluster-api-provider-vsphere/pkg/conversion/api/vmoperator/hub"
	conversionclient "sigs.k8s.io/cluster-api-provider-vsphere/pkg/conversion/client"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/util"
)

const (
	defaultAPIBindPort                   = 6443
	controlPlaneServiceAPIServerPortName = "apiserver"

	// ClusterSelectorKey is a label we use to store the cluster name on VirtualMachine objects.
	ClusterSelectorKey = "capv.vmware.com/cluster.name"
	nodeSelectorKey    = "capv.vmware.com/cluster.role"
	roleNode           = "node"
	roleControlPlane   = "controlplane"

	// TODO(lubronzhan): Deprecated, will be removed in a future release.
	// https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/issues/1483
	// legacyClusterSelectorKey and legacyNodeSelectorKey are added for backward compatibility.
	// These will be removed in the future release.
	// Please refer to the issue above for deprecation process.
	//
	// Deprecated: legacyClusterSelectorKey will be removed in a future release.
	legacyClusterSelectorKey = "capw.vmware.com/cluster.name"

	// Please refer to the issue above for deprecation process.
	//
	// Deprecated: legacyClusterSelectorKey will be removed in a future release.
	legacyNodeSelectorKey = "capw.vmware.com/cluster.role"
)

// CPService represents the ability to reconcile a ControlPlaneEndpoint.
type CPService struct {
	Client client.Client
}

// ReconcileControlPlaneEndpointService manages the lifecycle of a control plane endpoint managed by a vmoperator VirtualMachineService.
func (s *CPService) ReconcileControlPlaneEndpointService(ctx context.Context, clusterCtx *vmware.ClusterContext, netProvider services.NetworkProvider) (*vmwarev1.APIEndpoint, error) {
	// If the NetworkProvider does not support a load balancer, this should be a no-op
	if !netProvider.HasLoadBalancer() {
		return nil, nil
	}

	// 1. Determine the cluster's intended IP ipFamily
	ipFamily, err := util.DetermineClusterIPFamily(clusterCtx.Cluster)
	if err != nil {
		deprecatedv1beta1conditions.MarkFalse(clusterCtx.VSphereCluster, vmwarev1.LoadBalancerReadyV1Beta1Condition, vmwarev1.LoadBalancerCreationFailedV1Beta1Reason, clusterv1.ConditionSeverityWarning, "%v", err)
		conditions.Set(clusterCtx.VSphereCluster, metav1.Condition{
			Type:    vmwarev1.VSphereClusterLoadBalancerReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  vmwarev1.VSphereClusterLoadBalancerNotReadyReason,
			Message: err.Error(),
		})
		return nil, err
	}

	// 2. Capability Check
	ipv6DualStackSupported := netProvider.SupportsIPv6DualStack()

	// 3. Validate topology against capabilities
	if !ipv6DualStackSupported {
		if ipFamily == util.IPv6SingleStack || ipFamily == util.DualStackIPv4Primary || ipFamily == util.DualStackIPv6Primary {
			err := fmt.Errorf("IPv6 and DualStack require the IPv6DualStack feature gate, VM Operator v1alpha6+, and a network provider that supports it")
			deprecatedv1beta1conditions.MarkFalse(clusterCtx.VSphereCluster, vmwarev1.LoadBalancerReadyV1Beta1Condition, vmwarev1.LoadBalancerCreationFailedV1Beta1Reason, clusterv1.ConditionSeverityWarning, "%v", err)
			conditions.Set(clusterCtx.VSphereCluster, metav1.Condition{
				Type:    vmwarev1.VSphereClusterLoadBalancerReadyCondition,
				Status:  metav1.ConditionFalse,
				Reason:  vmwarev1.VSphereClusterLoadBalancerNotReadyReason,
				Message: err.Error(),
			})
			return nil, err
		}
	}

	vmService, err := s.getVMControlPlaneService(ctx, clusterCtx)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			err = pkgerrors.Wrapf(err, "failed to check if VirtualMachineService exists")
			deprecatedv1beta1conditions.MarkFalse(clusterCtx.VSphereCluster, vmwarev1.LoadBalancerReadyV1Beta1Condition, vmwarev1.LoadBalancerCreationFailedV1Beta1Reason, clusterv1.ConditionSeverityWarning, "%v", err)
			conditions.Set(clusterCtx.VSphereCluster, metav1.Condition{
				Type:    vmwarev1.VSphereClusterLoadBalancerReadyCondition,
				Status:  metav1.ConditionFalse,
				Reason:  vmwarev1.VSphereClusterLoadBalancerNotReadyReason,
				Message: err.Error(),
			})
			return nil, err
		}

		// Get the provider annotations for the ControlPlane Service.
		annotations, err := netProvider.GetVMServiceAnnotations(ctx, clusterCtx)
		if err != nil {
			err = pkgerrors.Wrapf(err, "failed to get provider VirtualMachineService annotations")
			deprecatedv1beta1conditions.MarkFalse(clusterCtx.VSphereCluster, vmwarev1.LoadBalancerReadyV1Beta1Condition, vmwarev1.LoadBalancerCreationFailedV1Beta1Reason, clusterv1.ConditionSeverityWarning, "%v", err)
			conditions.Set(clusterCtx.VSphereCluster, metav1.Condition{
				Type:    vmwarev1.VSphereClusterLoadBalancerReadyCondition,
				Status:  metav1.ConditionFalse,
				Reason:  vmwarev1.VSphereClusterLoadBalancerNotReadyReason,
				Message: err.Error(),
			})
			return nil, err
		}

		vmService, err = s.createVMControlPlaneService(ctx, clusterCtx, annotations, ipv6DualStackSupported, ipFamily)
		if err != nil {
			err = pkgerrors.Wrapf(err, "failed to create VirtualMachineService")
			deprecatedv1beta1conditions.MarkFalse(clusterCtx.VSphereCluster, vmwarev1.LoadBalancerReadyV1Beta1Condition, vmwarev1.LoadBalancerCreationFailedV1Beta1Reason, clusterv1.ConditionSeverityWarning, "%v", err)
			conditions.Set(clusterCtx.VSphereCluster, metav1.Condition{
				Type:    vmwarev1.VSphereClusterLoadBalancerReadyCondition,
				Status:  metav1.ConditionFalse,
				Reason:  vmwarev1.VSphereClusterLoadBalancerNotReadyReason,
				Message: err.Error(),
			})
			return nil, err
		}
	}

	// See if the LB has VIP(s) assigned, and delay reconciliation until it does
	primaryVIP, requiredVIPs, err := getAndValidateVIPs(vmService, ipFamily)
	if err != nil {
		deprecatedv1beta1conditions.MarkFalse(clusterCtx.VSphereCluster, vmwarev1.LoadBalancerReadyV1Beta1Condition, vmwarev1.WaitingForLoadBalancerIPV1Beta1Reason, clusterv1.ConditionSeverityInfo, "%v", err)
		conditions.Set(clusterCtx.VSphereCluster, metav1.Condition{
			Type:    vmwarev1.VSphereClusterLoadBalancerReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  vmwarev1.VSphereClusterLoadBalancerWaitingForIPReason,
			Message: err.Error(),
		})
		return nil, err
	}

	// Dual-stack specific KCP race-condition check
	if ipFamily == util.DualStackIPv4Primary || ipFamily == util.DualStackIPv6Primary {
		if err := ensureKCPReadyForControlPlaneEndpoint(ctx, s.Client, clusterCtx, requiredVIPs); err != nil {
			deprecatedv1beta1conditions.MarkFalse(clusterCtx.VSphereCluster, vmwarev1.LoadBalancerReadyV1Beta1Condition, vmwarev1.WaitingForKubeadmControlPlaneSpecReadyV1Beta1Reason, clusterv1.ConditionSeverityInfo, "%v", err)
			conditions.Set(clusterCtx.VSphereCluster, metav1.Condition{
				Type:    vmwarev1.VSphereClusterLoadBalancerReadyCondition,
				Status:  metav1.ConditionFalse,
				Reason:  vmwarev1.VSphereClusterLoadBalancerWaitingForKubeadmControlPlaneSpecReadyReason,
				Message: err.Error(),
			})
			return nil, nil
		}
	}

	cpEndpoint, err := getAPIEndpointFromVIP(vmService, primaryVIP)
	if err != nil {
		err = pkgerrors.Wrapf(err, "VirtualMachineService %s does not have an apiserver endpoint", klog.KObj(vmService))
		deprecatedv1beta1conditions.MarkFalse(clusterCtx.VSphereCluster, vmwarev1.LoadBalancerReadyV1Beta1Condition, vmwarev1.WaitingForLoadBalancerIPV1Beta1Reason, clusterv1.ConditionSeverityWarning, "%v", err)
		conditions.Set(clusterCtx.VSphereCluster, metav1.Condition{
			Type:    vmwarev1.VSphereClusterLoadBalancerReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  vmwarev1.VSphereClusterLoadBalancerWaitingForIPReason,
			Message: err.Error(),
		})
		return nil, err
	}

	deprecatedv1beta1conditions.MarkTrue(clusterCtx.VSphereCluster, vmwarev1.LoadBalancerReadyV1Beta1Condition)
	conditions.Set(clusterCtx.VSphereCluster, metav1.Condition{
		Type:   vmwarev1.VSphereClusterLoadBalancerReadyCondition,
		Status: metav1.ConditionTrue,
		Reason: vmwarev1.VSphereClusterLoadBalancerReadyReason,
	})
	return cpEndpoint, nil
}

func controlPlaneVMServiceName(clusterName string) string {
	return clusterName
}

// legacyControlPlaneVMServiceName was used for creating the ControlPlane VirtualMachineService prior to
// v1.13.0. It resulted in limiting the name of a Cluster to 41 characters (besides other places).
func legacyControlPlaneVMServiceName(clusterName string) string {
	return fmt.Sprintf("%s-control-plane-service", clusterName)
}

// ClusterRoleVMLabels returns labels applied to a VirtualMachine in the cluster. The Control Plane
// VM Service uses these labels to select VMs, as does the Cloud Provider.
// Add the legacyNodeSelectorKey and legacyClusterSelectorKey to machines as well.
func clusterRoleVMLabels(ctx *vmware.ClusterContext, controlPlane bool) map[string]string {
	result := map[string]string{
		ClusterSelectorKey:       ctx.Cluster.Name,
		legacyClusterSelectorKey: ctx.Cluster.Name,
	}
	if controlPlane {
		result[nodeSelectorKey] = roleControlPlane
		result[legacyNodeSelectorKey] = roleControlPlane
	} else {
		result[nodeSelectorKey] = roleNode
		result[legacyNodeSelectorKey] = roleNode
	}
	return result
}

func newVirtualMachineService(ctx *vmware.ClusterContext) *vmoprvhub.VirtualMachineService {
	return &vmoprvhub.VirtualMachineService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      controlPlaneVMServiceName(ctx.Cluster.Name),
			Namespace: ctx.Cluster.Namespace,
		},
	}
}

func (s *CPService) createVMControlPlaneService(ctx context.Context, clusterCtx *vmware.ClusterContext, annotations map[string]string, ipv6DualStackSupported bool, ipFamily util.ClusterIPFamily) (*vmoprvhub.VirtualMachineService, error) {
	// Note that the current implementation will only create a VirtualMachineService for a load balanced endpoint
	serviceType := vmoprvhub.VirtualMachineServiceTypeLoadBalancer

	vmService := newVirtualMachineService(clusterCtx)

	vmServiceExists := true
	if err := s.Client.Get(ctx, client.ObjectKeyFromObject(vmService), vmService); err != nil {
		if !apierrors.IsNotFound(err) {
			return nil, err
		}
		vmServiceExists = false
	}
	originalVMService := vmService.DeepCopy()

	if vmService.Annotations == nil {
		vmService.Annotations = annotations
	} else {
		for k, v := range annotations {
			vmService.Annotations[k] = v
		}
	}
	vmService.Annotations = annotations
	spec := vmoprvhub.VirtualMachineServiceSpec{
		Type: serviceType,
		Ports: []vmoprvhub.VirtualMachineServicePort{
			{
				Name:       controlPlaneServiceAPIServerPortName,
				Protocol:   "TCP",
				Port:       defaultAPIBindPort,
				TargetPort: defaultAPIBindPort,
			},
		},
		Selector: clusterRoleVMLabels(clusterCtx, true),
	}

	if ipv6DualStackSupported {
		spec.IPFamilyPolicy, spec.IPFamilies = getIPFamilyConfig(ipFamily)
	}

	vmService.Spec = spec

	if err := ctrlutil.SetOwnerReference(
		clusterCtx.VSphereCluster,
		vmService,
		s.Client.Scheme(),
	); err != nil {
		return nil, pkgerrors.Wrapf(
			err,
			"error setting %s/%s as owner of %s/%s",
			clusterCtx.VSphereCluster.Namespace,
			clusterCtx.VSphereCluster.Name,
			vmService.Namespace,
			vmService.Name,
		)
	}

	if !vmServiceExists {
		log := ctrl.LoggerFrom(ctx)
		log.Info("Creating VirtualMachineService", "VirtualMachineService", klog.KObj(vmService))
		if err := s.Client.Create(ctx, vmService); err != nil {
			return nil, pkgerrors.Wrapf(err, "failed to create VirtualMachineService %s", klog.KObj(vmService))
		}
	} else if !reflect.DeepEqual(originalVMService, vmService) {
		patch, err := conversionclient.MergeFrom(ctx, s.Client, originalVMService)
		if err != nil {
			return nil, pkgerrors.Wrapf(err, "failed to create patch for VirtualMachineService object")
		}
		if err := s.Client.Patch(ctx, vmService, patch); err != nil {
			return nil, pkgerrors.Wrapf(err, "failed to patch VirtualMachineService object")
		}
	}

	return vmService, nil
}

func (s *CPService) getVMControlPlaneService(ctx context.Context, clusterCtx *vmware.ClusterContext) (*vmoprvhub.VirtualMachineService, error) {
	log := ctrl.LoggerFrom(ctx)

	vmService := &vmoprvhub.VirtualMachineService{}
	vmServiceKey := client.ObjectKey{
		Namespace: clusterCtx.Cluster.Namespace,
		Name:      controlPlaneVMServiceName(clusterCtx.Cluster.Name),
	}
	if err := s.Client.Get(ctx, vmServiceKey, vmService); err != nil {
		if !apierrors.IsNotFound(err) {
			return nil, fmt.Errorf("failed to get VirtualMachineService %s: %v", vmServiceKey.Name, err)
		}

		// In case of not finding the ControlPlane VirtualMachineService: fallback to try the legacy name.
		fallbackVMServiceKey := client.ObjectKey{
			Namespace: vmServiceKey.Namespace,
			Name:      legacyControlPlaneVMServiceName(clusterCtx.Cluster.Name),
		}
		if fallbackErr := s.Client.Get(ctx, fallbackVMServiceKey, vmService); fallbackErr != nil {
			if !apierrors.IsNotFound(fallbackErr) {
				return nil, fmt.Errorf("failed to get VirtualMachineService %s: %v", fallbackVMServiceKey.Name, fallbackErr)
			}

			log.Info("VirtualMachineService was not found", "VirtualMachineService", klog.KRef(vmServiceKey.Namespace, vmServiceKey.Name))
			return nil, err
		}
	}

	// Verify OwnerReference UID to prevent adopting a service from a
	// previous cluster with the same name.
	refs := vmService.GetOwnerReferences()
	for _, ref := range refs {
		if ref.Kind == "VSphereCluster" && ref.UID != clusterCtx.VSphereCluster.UID {
			return nil, fmt.Errorf("VirtualMachineService %s exists but is owned by a different VSphereCluster instance %s", vmServiceKey.Name, ref.UID)
		}
	}

	// If the service is being deleted, it could be an old service with the same name,
	// return error.
	if !vmService.DeletionTimestamp.IsZero() {
		return nil, fmt.Errorf("VirtualMachineService %s exists but is being deleted", vmServiceKey.Name)
	}

	return vmService, nil
}

// getIPFamilyConfig constructs the IPFamily settings based on IP topology.
func getIPFamilyConfig(ipFamily util.ClusterIPFamily) (*corev1.IPFamilyPolicy, []corev1.IPFamily) {
	switch ipFamily {
	case util.IPv4SingleStack:
		return ptr.To(corev1.IPFamilyPolicySingleStack), []corev1.IPFamily{corev1.IPv4Protocol}
	case util.IPv6SingleStack:
		return ptr.To(corev1.IPFamilyPolicySingleStack), []corev1.IPFamily{corev1.IPv6Protocol}
	case util.DualStackIPv4Primary:
		return ptr.To(corev1.IPFamilyPolicyRequireDualStack), []corev1.IPFamily{corev1.IPv4Protocol, corev1.IPv6Protocol}
	case util.DualStackIPv6Primary:
		return ptr.To(corev1.IPFamilyPolicyRequireDualStack), []corev1.IPFamily{corev1.IPv6Protocol, corev1.IPv4Protocol}
	default:
		return nil, nil
	}
}

// getAndValidateVIPs validates status VIPs against topology and returns the primary VIP and required VIP list.
func getAndValidateVIPs(vmService *vmoprvhub.VirtualMachineService, topology util.ClusterIPFamily) (primaryVIP string, requiredVIPs []string, err error) {
	ipv4VIP, ipv6VIP, err := getVMServiceVIPs(vmService)
	if err != nil {
		return "", nil, err
	}

	switch topology {
	case util.IPv4SingleStack:
		if ipv4VIP == "" {
			return "", nil, fmt.Errorf("VirtualMachineService %s does not yet have IPv4 VIP assigned", klog.KObj(vmService))
		}
		return ipv4VIP, []string{ipv4VIP}, nil

	case util.IPv6SingleStack:
		if ipv6VIP == "" {
			return "", nil, fmt.Errorf("VirtualMachineService %s does not yet have IPv6 VIP assigned", klog.KObj(vmService))
		}
		return ipv6VIP, []string{ipv6VIP}, nil

	case util.DualStackIPv4Primary:
		if ipv4VIP == "" || ipv6VIP == "" {
			return "", nil, fmt.Errorf("VirtualMachineService %s must have both IPv4 and IPv6 ingress for dual stack cluster (have IPv4: %v, IPv6: %v)", klog.KObj(vmService), ipv4VIP != "", ipv6VIP != "")
		}
		return ipv4VIP, []string{ipv4VIP, ipv6VIP}, nil

	case util.DualStackIPv6Primary:
		if ipv4VIP == "" || ipv6VIP == "" {
			return "", nil, fmt.Errorf("VirtualMachineService %s must have both IPv4 and IPv6 ingress for dual stack cluster (have IPv4: %v, IPv6: %v)", klog.KObj(vmService), ipv4VIP != "", ipv6VIP != "")
		}
		return ipv6VIP, []string{ipv6VIP, ipv4VIP}, nil

	default:
		return "", nil, fmt.Errorf("unknown cluster IPFamily")
	}
}

// getVMServiceVIPs returns IPv4 and IPv6 from the VirtualMachineService LoadBalancer ingress.
// For dual stack we require both families to be present before setting the control plane endpoint.
func getVMServiceVIPs(vmService *vmoprvhub.VirtualMachineService) (ipv4, ipv6 string, err error) {
	if vmService.Spec.Type != vmoprvhub.VirtualMachineServiceTypeLoadBalancer {
		return "", "", fmt.Errorf("VirtualMachineService %s does not have load balancer", klog.KObj(vmService))
	}

	var ipv4Addr, ipv6Addr string
	for _, ingress := range vmService.Status.LoadBalancer.Ingress {
		if ingress.IP != "" {
			ip := net.ParseIP(ingress.IP)
			if ip == nil {
				return "", "", fmt.Errorf("invalid load balancer ingress IP address %q for VirtualMachineService %s", ingress.IP, klog.KObj(vmService))
			}
			if ip.To4() != nil {
				if ipv4Addr == "" {
					ipv4Addr = ingress.IP
				}
			} else {
				if ipv6Addr == "" {
					ipv6Addr = ingress.IP
				}
			}
		}
	}

	if ipv4Addr == "" && ipv6Addr == "" {
		return "", "", fmt.Errorf("VirtualMachineService %s LoadBalancer does not have any Ingresses", klog.KObj(vmService))
	}
	return ipv4Addr, ipv6Addr, nil
}

// ensureKCPReadyForControlPlaneEndpoint ensures KubeadmControlPlane has (a) requiredIPs in certSANs and
// (b) status.observedGeneration == metadata.generation to avoid a race between CAPV and KCP controllers.
// This prevents the first control plane machine from being created with stale certSANs (e.g. single stack)
// and avoids an extra rollout for dual stack clusters.
func ensureKCPReadyForControlPlaneEndpoint(ctx context.Context, c client.Client, clusterCtx *vmware.ClusterContext, requiredIPs []string) error {
	// Early return for single stack or if cluster is not set.
	if clusterCtx.Cluster == nil || len(requiredIPs) <= 1 {
		return nil
	}

	if !clusterCtx.Cluster.Spec.ControlPlaneRef.IsDefined() {
		return fmt.Errorf("ControlPlaneRef is not defined for Cluster %s", klog.KObj(clusterCtx.Cluster))
	}

	kcp := &controlplanev1.KubeadmControlPlane{}
	kcpKey := client.ObjectKey{
		Namespace: clusterCtx.Cluster.Namespace,
		Name:      clusterCtx.Cluster.Spec.ControlPlaneRef.Name,
	}
	if err := c.Get(ctx, kcpKey, kcp); err != nil {
		return pkgerrors.Wrapf(err, "failed to get KubeadmControlPlane %s", klog.KRef(kcpKey.Namespace, kcpKey.Name))
	}

	if !kcp.GetDeletionTimestamp().IsZero() {
		return nil
	}

	// (b) Avoid race: KCP controller must have observed and successfully reconciled the current spec generation.
	// We check both the top level ObservedGeneration and the ObservedGeneration on individual conditions
	// to avoid a deadlock: KCP might be waiting for CAPV to set the ControlPlaneEndpoint before it can
	// fully reconcile, while CAPV is waiting for KCP to observe the latest spec (e.g. certSANs for dual-stack).
	observed := kcp.Status.ObservedGeneration == kcp.GetGeneration()
	if !observed {
		for _, cond := range kcp.Status.Conditions {
			if cond.ObservedGeneration == kcp.GetGeneration() {
				observed = true
				break
			}
		}
	}
	if !observed {
		return fmt.Errorf("KubeadmControlPlane %s has not yet observed the current generation %d (observed: %d, no condition has observed the current generation)",
			klog.KObj(kcp), kcp.GetGeneration(), kcp.Status.ObservedGeneration)
	}

	// (a) For dual stack, certSANs must contain all required IPs (both primary and secondary LB VIPs) so the first machine gets both in certs.
	// The runtime extension must add both IPs to apiServer.certSANs; adding only the secondary-family IP is insufficient.
	certSANs := kcp.Spec.KubeadmConfigSpec.ClusterConfiguration.APIServer.CertSANs
	for _, ip := range requiredIPs {
		if !slices.Contains(certSANs, ip) {
			return fmt.Errorf("KubeadmControlPlane %s certSANs does not contain %q yet; the runtime extension might not have set the IPs yet",
				klog.KObj(kcp), ip)
		}
	}
	return nil
}

func getAPIEndpointFromVIP(vmService *vmoprvhub.VirtualMachineService, vip string) (*vmwarev1.APIEndpoint, error) {
	name := controlPlaneServiceAPIServerPortName
	servicePort := int32(-1)
	for _, port := range vmService.Spec.Ports {
		if port.Name == name && port.Protocol == "TCP" {
			servicePort = port.Port
			break
		}
	}

	if servicePort == -1 {
		return nil, fmt.Errorf("VirtualMachineService does not have port entry for %q", name)
	}

	return &vmwarev1.APIEndpoint{
		Host: vip,
		Port: servicePort,
	}, nil
}
