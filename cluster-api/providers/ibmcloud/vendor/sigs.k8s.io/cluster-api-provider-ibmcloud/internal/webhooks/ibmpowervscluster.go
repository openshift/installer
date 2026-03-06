/*
Copyright 2022 The Kubernetes Authors.

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

package webhooks

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	regionUtil "github.com/ppc64le-cloud/powervs-utils"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/internal/genutil"
)

//+kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-ibmpowervscluster,mutating=true,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsclusters,verbs=create;update,versions=v1beta2,name=mibmpowervscluster.kb.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
//+kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-ibmpowervscluster,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsclusters,versions=v1beta2,name=vibmpowervscluster.kb.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

func (r *IBMPowerVSCluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&infrav1.IBMPowerVSCluster{}).
		WithValidator(r).
		WithDefaulter(r).
		Complete()
}

// IBMPowerVSCluster implements a validation and defaulting webhook for IBMPowerVSCluster.
type IBMPowerVSCluster struct{}

var _ webhook.CustomDefaulter = &IBMPowerVSCluster{}
var _ webhook.CustomValidator = &IBMPowerVSCluster{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the type.
func (r *IBMPowerVSCluster) Default(_ context.Context, _ runtime.Object) error {
	return nil
}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type.
func (r *IBMPowerVSCluster) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	objValue, ok := obj.(*infrav1.IBMPowerVSCluster)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a IBMPowerVSCluster but got a %T", obj))
	}
	return validateIBMPowerVSCluster(nil, objValue)
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type.
func (r *IBMPowerVSCluster) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (warnings admission.Warnings, err error) {
	oldObjValue, ok := oldObj.(*infrav1.IBMPowerVSCluster)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a IBMPowerVSCluster but got a %T", oldObj))
	}
	newObjValue, ok := newObj.(*infrav1.IBMPowerVSCluster)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a IBMPowerVSCluster but got a %T", newObj))
	}
	return validateIBMPowerVSCluster(oldObjValue, newObjValue)
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type.
func (r *IBMPowerVSCluster) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

func validateIBMPowerVSCluster(oldCluster, newCluster *infrav1.IBMPowerVSCluster) (admission.Warnings, error) {
	var allErrs field.ErrorList
	if err := validateIBMPowerVSClusterNetwork(newCluster); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := validateIBMPowerVSClusterCreateInfraPrereq(newCluster); err != nil {
		allErrs = append(allErrs, err...)
	}
	// Need not validate for create operation
	if oldCluster != nil {
		if err := validateAdditionalListenerSelector(newCluster, oldCluster); err != nil {
			allErrs = append(allErrs, err...)
		}
	}

	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(
		schema.GroupKind{Group: "infrastructure.cluster.x-k8s.io", Kind: "IBMPowerVSCluster"},
		newCluster.Name, allErrs)
}

func validateIBMPowerVSClusterNetwork(cluster *infrav1.IBMPowerVSCluster) *field.Error {
	if res, err := validateIBMPowerVSNetworkReference(cluster.Spec.Network); !res {
		return err
	}
	if (cluster.Spec.Network.Name != nil || cluster.Spec.Network.ID != nil) && (cluster.Spec.DHCPServer != nil && cluster.Spec.DHCPServer.Name != nil) {
		return field.Invalid(field.NewPath("spec.dhcpServer.name"), cluster.Spec.DHCPServer.Name, "either one of network or dhcpServer details can be provided")
	}
	if (cluster.Spec.Network.Name != nil || cluster.Spec.Network.ID != nil) && (cluster.Spec.DHCPServer != nil && cluster.Spec.DHCPServer.ID != nil) {
		return field.Invalid(field.NewPath("spec.dhcpServer.id"), cluster.Spec.DHCPServer.ID, "either one of network or dhcpServer details can be provided")
	}
	return nil
}

func validateIBMPowerVSClusterLoadBalancers(cluster *infrav1.IBMPowerVSCluster) (allErrs field.ErrorList) {
	if err := validateIBMPowerVSClusterLoadBalancerNames(cluster); err != nil {
		allErrs = append(allErrs, err...)
	}

	if len(cluster.Spec.LoadBalancers) == 0 {
		return allErrs
	}

	for _, loadBalancer := range cluster.Spec.LoadBalancers {
		if *loadBalancer.Public {
			return allErrs
		}
	}

	return append(allErrs, field.Invalid(field.NewPath("spec.LoadBalancers"), cluster.Spec.LoadBalancers, "Expect atleast one of the load balancer to be public"))
}

func validateIBMPowerVSClusterLoadBalancerNames(cluster *infrav1.IBMPowerVSCluster) (allErrs field.ErrorList) {
	found := make(map[string]bool)
	for i, loadbalancer := range cluster.Spec.LoadBalancers {
		if loadbalancer.Name == "" {
			continue
		}

		if found[loadbalancer.Name] {
			allErrs = append(allErrs, field.Duplicate(field.NewPath("spec", fmt.Sprintf("loadbalancers[%d]", i)), map[string]interface{}{"Name": loadbalancer.Name}))
			continue
		}
		found[loadbalancer.Name] = true
	}

	return allErrs
}

func validateIBMPowerVSClusterVPCSubnetNames(cluster *infrav1.IBMPowerVSCluster) (allErrs field.ErrorList) {
	found := make(map[string]bool)
	for i, subnet := range cluster.Spec.VPCSubnets {
		if subnet.Name == nil {
			continue
		}
		if found[*subnet.Name] {
			allErrs = append(allErrs, field.Duplicate(field.NewPath("spec", fmt.Sprintf("vpcSubnets[%d]", i)), map[string]interface{}{"Name": *subnet.Name}))
			continue
		}
		found[*subnet.Name] = true
	}

	return allErrs
}

func validateIBMPowerVSClusterTransitGateway(cluster *infrav1.IBMPowerVSCluster) *field.Error {
	if cluster.Spec.Zone == nil && cluster.Spec.VPC == nil {
		return nil
	}
	if cluster.Spec.TransitGateway == nil {
		return nil
	}
	if _, globalRouting, _ := genutil.GetTransitGatewayLocationAndRouting(cluster.Spec.Zone, cluster.Spec.VPC.Region); cluster.Spec.TransitGateway.GlobalRouting != nil && !*cluster.Spec.TransitGateway.GlobalRouting && globalRouting != nil && *globalRouting {
		return field.Invalid(field.NewPath("spec.transitGateway.globalRouting"), cluster.Spec.TransitGateway.GlobalRouting, "global routing is required since PowerVS and VPC region are from different region")
	}
	return nil
}

func validateIBMPowerVSClusterCreateInfraPrereq(cluster *infrav1.IBMPowerVSCluster) (allErrs field.ErrorList) {
	annotations := cluster.GetAnnotations()
	if len(annotations) == 0 {
		return nil
	}

	value, found := annotations[infrav1.CreateInfrastructureAnnotation]
	if !found {
		return nil
	}

	createInfra, err := strconv.ParseBool(value)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("annotations"), cluster.Annotations, "value of powervs.cluster.x-k8s.io/create-infra should be boolean"))
	}

	if !createInfra {
		return nil
	}

	if cluster.Spec.Zone == nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec.zone"), cluster.Spec.Zone, "value of zone is empty"))
	}

	if cluster.Spec.Zone != nil && !regionUtil.ValidateZone(*cluster.Spec.Zone) {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec.zone"), cluster.Spec.Zone, fmt.Sprintf("zone '%s' is not supported", *cluster.Spec.Zone)))
	}

	if cluster.Spec.VPC == nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec.vpc"), cluster.Spec.VPC, "value of VPC is empty"))
	}

	if cluster.Spec.VPC != nil && cluster.Spec.VPC.Region == nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec.vpc.region"), cluster.Spec.VPC.Region, "value of VPC region is empty"))
	}

	if cluster.Spec.VPC != nil && cluster.Spec.VPC.Region != nil && !regionUtil.ValidateVPCRegion(*cluster.Spec.VPC.Region) {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec.vpc.region"), cluster.Spec.VPC.Region, fmt.Sprintf("vpc region '%s' is not supported", *cluster.Spec.VPC.Region)))
	}

	if cluster.Spec.ResourceGroup == nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec.resourceGroup"), cluster.Spec.ResourceGroup, "value of resource group is empty"))
	}
	if err := validateIBMPowerVSClusterVPCSubnetNames(cluster); err != nil {
		allErrs = append(allErrs, err...)
	}

	if err := validateIBMPowerVSClusterLoadBalancers(cluster); err != nil {
		allErrs = append(allErrs, err...)
	}

	if err := validateIBMPowerVSClusterTransitGateway(cluster); err != nil {
		allErrs = append(allErrs, err)
	}

	return allErrs
}

func validateAdditionalListenerSelector(newCluster, oldCluster *infrav1.IBMPowerVSCluster) (allErrs field.ErrorList) {
	newLoadBalancerListeners := map[string]metav1.LabelSelector{}
	for _, loadbalancer := range newCluster.Spec.LoadBalancers {
		for _, additionalListener := range loadbalancer.AdditionalListeners {
			newLoadBalancerListeners[fmt.Sprintf("%d-%s", additionalListener.Port, *additionalListener.Protocol)] = additionalListener.Selector
		}
	}
	for _, loadbalancer := range oldCluster.Spec.LoadBalancers {
		for _, additionalListener := range loadbalancer.AdditionalListeners {
			if selector, ok := newLoadBalancerListeners[fmt.Sprintf("%d-%s", additionalListener.Port, *additionalListener.Protocol)]; ok && !reflect.DeepEqual(selector, additionalListener.Selector) {
				allErrs = append(allErrs, field.Forbidden(field.NewPath("selector"), "Selector is immutable"))
			}
		}
	}
	return allErrs
}
