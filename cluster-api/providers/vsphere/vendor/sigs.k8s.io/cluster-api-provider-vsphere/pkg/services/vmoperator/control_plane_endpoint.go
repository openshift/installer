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
	"fmt"

	"github.com/pkg/errors"
	vmoprv1 "github.com/vmware-tanzu/vm-operator/api/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/vmware/v1beta1"
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
	legacyClusterSelectorKey = "capw.vmware.com/cluster.name"
	legacyNodeSelectorKey    = "capw.vmware.com/cluster.role"
)

// CPService represents the ability to reconcile a ControlPlaneEndpoint.
type CPService struct{}

// ReconcileControlPlaneEndpointService manages the lifecycle of a control plane endpoint managed by a vmoperator VirtualMachineService.
func (s CPService) ReconcileControlPlaneEndpointService(ctx *vmware.ClusterContext, netProvider services.NetworkProvider) (*clusterv1.APIEndpoint, error) {
	ctx.Logger.V(4).Info("Reconciling control plane VirtualMachineService for cluster")

	// If the NetworkProvider does not support a load balancer, this should be a no-op
	if !netProvider.HasLoadBalancer() {
		return nil, nil
	}

	vmService, err := s.getVMControlPlaneService(ctx)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			err = errors.Wrapf(err, "failed to check if VirtualMachineService exists")
			conditions.MarkFalse(ctx.VSphereCluster, infrav1.LoadBalancerReadyCondition, infrav1.LoadBalancerCreationFailedReason, clusterv1.ConditionSeverityWarning, err.Error())
			return nil, err
		}

		// Get the provider annotations for the ControlPlane Service.
		annotations, err := netProvider.GetVMServiceAnnotations(ctx)
		if err != nil {
			err = errors.Wrapf(err, "failed to get provider VirtualMachineService annotations")
			conditions.MarkFalse(ctx.VSphereCluster, infrav1.LoadBalancerReadyCondition, infrav1.LoadBalancerCreationFailedReason, clusterv1.ConditionSeverityWarning, err.Error())
			return nil, err
		}

		vmService, err = s.createVMControlPlaneService(ctx, annotations)
		if err != nil {
			err = errors.Wrapf(err, "failed to create VirtualMachineService")
			conditions.MarkFalse(ctx.VSphereCluster, infrav1.LoadBalancerReadyCondition, infrav1.LoadBalancerCreationFailedReason, clusterv1.ConditionSeverityWarning, err.Error())
			return nil, err
		}
	}

	// See if the LB has a VIP assigned, and delay reconciliation until it does
	vip, err := getVMServiceVIP(vmService)
	if err != nil {
		err = errors.Wrapf(err, "VirtualMachineService LB does not yet have VIP assigned")
		conditions.MarkFalse(ctx.VSphereCluster, infrav1.LoadBalancerReadyCondition, infrav1.WaitingForLoadBalancerIPReason, clusterv1.ConditionSeverityInfo, err.Error())
		return nil, err
	}

	cpEndpoint, err := getAPIEndpointFromVIP(vmService, vip)
	if err != nil {
		err = errors.Wrapf(err, "VirtualMachineService LB does not have an apiserver endpoint")
		conditions.MarkFalse(ctx.VSphereCluster, infrav1.LoadBalancerReadyCondition, infrav1.WaitingForLoadBalancerIPReason, clusterv1.ConditionSeverityWarning, err.Error())
		return nil, err
	}

	conditions.MarkTrue(ctx.VSphereCluster, infrav1.LoadBalancerReadyCondition)
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

func (s CPService) createVMControlPlaneService(ctx *vmware.ClusterContext, annotations map[string]string) (*vmoprv1.VirtualMachineService, error) {
	// Note that the current implementation will only create a VirtualMachineService for a load balanced endpoint
	serviceType := vmoprv1.VirtualMachineServiceTypeLoadBalancer

	vmService := newVirtualMachineService(ctx)

	_, err := ctrlutil.CreateOrPatch(ctx, ctx.Client, vmService, func() error {
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
			Selector: clusterRoleVMLabels(ctx, true),
		}

		if err := ctrlutil.SetOwnerReference(
			ctx.VSphereCluster,
			vmService,
			ctx.Scheme,
		); err != nil {
			return errors.Wrapf(
				err,
				"error setting %s/%s as owner of %s/%s",
				ctx.VSphereCluster.Namespace,
				ctx.VSphereCluster.Name,
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

func (s CPService) getVMControlPlaneService(ctx *vmware.ClusterContext) (*vmoprv1.VirtualMachineService, error) {
	vmService := &vmoprv1.VirtualMachineService{}
	vmServiceName := client.ObjectKey{
		Namespace: ctx.Cluster.Namespace,
		Name:      controlPlaneVMServiceName(ctx),
	}
	if err := ctx.Client.Get(ctx, vmServiceName, vmService); err != nil {
		if !apierrors.IsNotFound(err) {
			return nil, fmt.Errorf("failed to get VirtualMachineService %s: %v", vmServiceName.Name, err)
		}

		ctx.Logger.V(2).Info("VirtualMachineService was not found", "cluster", ctx.Cluster.Name, "service", vmServiceName.Name)
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
