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

	"github.com/pkg/errors"
	vmoprv1 "github.com/vmware-tanzu/vm-operator/api/v1alpha2"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/vmware/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context/vmware"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services"
)

const (
	defaultAPIBindPort                   = 6443
	controlPlaneServiceAPIServerPortName = "apiserver"

	clusterSelectorKey = "capv.vmware.com/cluster.name"
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
func (s *CPService) ReconcileControlPlaneEndpointService(ctx context.Context, clusterCtx *vmware.ClusterContext, netProvider services.NetworkProvider) (*clusterv1.APIEndpoint, error) {
	log := ctrl.LoggerFrom(ctx)
	log.V(4).Info("Reconciling control plane VirtualMachineService for cluster")

	// If the NetworkProvider does not support a load balancer, this should be a no-op
	if !netProvider.HasLoadBalancer() {
		return nil, nil
	}

	vmService, err := s.getVMControlPlaneService(ctx, clusterCtx)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			err = errors.Wrapf(err, "failed to check if VirtualMachineService exists")
			conditions.MarkFalse(clusterCtx.VSphereCluster, vmwarev1.LoadBalancerReadyCondition, vmwarev1.LoadBalancerCreationFailedReason, clusterv1.ConditionSeverityWarning, err.Error())
			return nil, err
		}

		// Get the provider annotations for the ControlPlane Service.
		annotations, err := netProvider.GetVMServiceAnnotations(ctx, clusterCtx)
		if err != nil {
			err = errors.Wrapf(err, "failed to get provider VirtualMachineService annotations")
			conditions.MarkFalse(clusterCtx.VSphereCluster, vmwarev1.LoadBalancerReadyCondition, vmwarev1.LoadBalancerCreationFailedReason, clusterv1.ConditionSeverityWarning, err.Error())
			return nil, err
		}

		vmService, err = s.createVMControlPlaneService(ctx, clusterCtx, annotations)
		if err != nil {
			err = errors.Wrapf(err, "failed to create VirtualMachineService")
			conditions.MarkFalse(clusterCtx.VSphereCluster, vmwarev1.LoadBalancerReadyCondition, vmwarev1.LoadBalancerCreationFailedReason, clusterv1.ConditionSeverityWarning, err.Error())
			return nil, err
		}
	}

	// See if the LB has a VIP assigned, and delay reconciliation until it does
	vip, err := getVMServiceVIP(vmService)
	if err != nil {
		err = errors.Wrapf(err, "VirtualMachineService LB does not yet have VIP assigned")
		conditions.MarkFalse(clusterCtx.VSphereCluster, vmwarev1.LoadBalancerReadyCondition, vmwarev1.WaitingForLoadBalancerIPReason, clusterv1.ConditionSeverityInfo, err.Error())
		return nil, err
	}

	cpEndpoint, err := getAPIEndpointFromVIP(vmService, vip)
	if err != nil {
		err = errors.Wrapf(err, "VirtualMachineService LB does not have an apiserver endpoint")
		conditions.MarkFalse(clusterCtx.VSphereCluster, vmwarev1.LoadBalancerReadyCondition, vmwarev1.WaitingForLoadBalancerIPReason, clusterv1.ConditionSeverityWarning, err.Error())
		return nil, err
	}

	conditions.MarkTrue(clusterCtx.VSphereCluster, vmwarev1.LoadBalancerReadyCondition)
	return cpEndpoint, nil
}

func controlPlaneVMServiceName(ctx *vmware.ClusterContext) string {
	return fmt.Sprintf("%s-control-plane-service", ctx.Cluster.Name)
}

// ClusterRoleVMLabels returns labels applied to a VirtualMachine in the cluster. The Control Plane
// VM Service uses these labels to select VMs, as does the Cloud Provider.
// Add the legacyNodeSelectorKey and legacyClusterSelectorKey to machines as well.
func clusterRoleVMLabels(ctx *vmware.ClusterContext, controlPlane bool) map[string]string {
	result := map[string]string{
		clusterSelectorKey:       ctx.Cluster.Name,
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

func newVirtualMachineService(ctx *vmware.ClusterContext) *vmoprv1.VirtualMachineService {
	return &vmoprv1.VirtualMachineService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      controlPlaneVMServiceName(ctx),
			Namespace: ctx.Cluster.Namespace,
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: vmoprv1.SchemeGroupVersion.String(),
			Kind:       "VirtualMachineService",
		},
	}
}

func (s *CPService) createVMControlPlaneService(ctx context.Context, clusterCtx *vmware.ClusterContext, annotations map[string]string) (*vmoprv1.VirtualMachineService, error) {
	// Note that the current implementation will only create a VirtualMachineService for a load balanced endpoint
	serviceType := vmoprv1.VirtualMachineServiceTypeLoadBalancer

	vmService := newVirtualMachineService(clusterCtx)

	_, err := ctrlutil.CreateOrPatch(ctx, s.Client, vmService, func() error {
		if vmService.Annotations == nil {
			vmService.Annotations = annotations
		} else {
			for k, v := range annotations {
				vmService.Annotations[k] = v
			}
		}
		vmService.Annotations = annotations
		vmService.Spec = vmoprv1.VirtualMachineServiceSpec{
			Type: serviceType,
			Ports: []vmoprv1.VirtualMachineServicePort{
				{
					Name:       controlPlaneServiceAPIServerPortName,
					Protocol:   "TCP",
					Port:       defaultAPIBindPort,
					TargetPort: defaultAPIBindPort,
				},
			},
			Selector: clusterRoleVMLabels(clusterCtx, true),
		}

		if err := ctrlutil.SetOwnerReference(
			clusterCtx.VSphereCluster,
			vmService,
			s.Client.Scheme(),
		); err != nil {
			return errors.Wrapf(
				err,
				"error setting %s/%s as owner of %s/%s",
				clusterCtx.VSphereCluster.Namespace,
				clusterCtx.VSphereCluster.Name,
				vmService.Namespace,
				vmService.Name,
			)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return vmService, nil
}

func (s *CPService) getVMControlPlaneService(ctx context.Context, clusterCtx *vmware.ClusterContext) (*vmoprv1.VirtualMachineService, error) {
	log := ctrl.LoggerFrom(ctx)

	vmService := &vmoprv1.VirtualMachineService{}
	vmServiceKey := client.ObjectKey{
		Namespace: clusterCtx.Cluster.Namespace,
		Name:      controlPlaneVMServiceName(clusterCtx),
	}
	if err := s.Client.Get(ctx, vmServiceKey, vmService); err != nil {
		if !apierrors.IsNotFound(err) {
			return nil, fmt.Errorf("failed to get VirtualMachineService %s: %v", vmServiceKey.Name, err)
		}

		log.Info("VirtualMachineService was not found", "VirtualMachineService", klog.KRef(vmServiceKey.Namespace, vmServiceKey.Name))
		return nil, err
	}

	return vmService, nil
}

func getVMServiceVIP(vmService *vmoprv1.VirtualMachineService) (string, error) {
	if vmService.Spec.Type != vmoprv1.VirtualMachineServiceTypeLoadBalancer {
		return "", fmt.Errorf("VirtualMachineService for control plane does not have load balancer")
	}

	for _, ingress := range vmService.Status.LoadBalancer.Ingress {
		if ingress.IP != "" {
			return ingress.IP, nil
		}
		// BMV: Supported?
		// if ingress.Hostname != "" {
		// 	return ingress.Hostname, nil
		// }
	}

	return "", fmt.Errorf("VirtualMachineService LoadBalancer does not have any Ingresses")
}

func getAPIEndpointFromVIP(vmService *vmoprv1.VirtualMachineService, vip string) (*clusterv1.APIEndpoint, error) {
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

	return &clusterv1.APIEndpoint{
		Host: vip,
		Port: servicePort,
	}, nil
}
