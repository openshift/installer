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
	"strconv"

	regionUtil "github.com/ppc64le-cloud/powervs-utils"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1beta2 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
	genUtil "sigs.k8s.io/cluster-api-provider-ibmcloud/util"
)

//+kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-ibmpowervscluster,mutating=true,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsclusters,verbs=create;update,versions=v1beta2,name=mibmpowervscluster.kb.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
//+kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-ibmpowervscluster,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsclusters,versions=v1beta2,name=vibmpowervscluster.kb.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

func (r *IBMPowerVSCluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&infrav1beta2.IBMPowerVSCluster{}).
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
	objValue, ok := obj.(*infrav1beta2.IBMPowerVSCluster)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a IBMPowerVSCluster but got a %T", obj))
	}
	return validateIBMPowerVSCluster(objValue)
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type.
func (r *IBMPowerVSCluster) ValidateUpdate(_ context.Context, _, newObj runtime.Object) (warnings admission.Warnings, err error) {
	objValue, ok := newObj.(*infrav1beta2.IBMPowerVSCluster)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a IBMPowerVSCluster but got a %T", newObj))
	}
	return validateIBMPowerVSCluster(objValue)
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type.
func (r *IBMPowerVSCluster) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

func validateIBMPowerVSCluster(cluster *infrav1beta2.IBMPowerVSCluster) (admission.Warnings, error) {
	var allErrs field.ErrorList
	if err := validateIBMPowerVSClusterNetwork(cluster); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := validateIBMPowerVSClusterCreateInfraPrereq(cluster); err != nil {
		allErrs = append(allErrs, err...)
	}

	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(
		schema.GroupKind{Group: "infrastructure.cluster.x-k8s.io", Kind: "IBMPowerVSCluster"},
		cluster.Name, allErrs)
}

func validateIBMPowerVSClusterNetwork(cluster *infrav1beta2.IBMPowerVSCluster) *field.Error {
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

func validateIBMPowerVSClusterLoadBalancers(cluster *infrav1beta2.IBMPowerVSCluster) (allErrs field.ErrorList) {
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

func validateIBMPowerVSClusterLoadBalancerNames(cluster *infrav1beta2.IBMPowerVSCluster) (allErrs field.ErrorList) {
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

func validateIBMPowerVSClusterVPCSubnetNames(cluster *infrav1beta2.IBMPowerVSCluster) (allErrs field.ErrorList) {
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

func validateIBMPowerVSClusterTransitGateway(cluster *infrav1beta2.IBMPowerVSCluster) *field.Error {
	if cluster.Spec.Zone == nil && cluster.Spec.VPC == nil {
		return nil
	}
	if cluster.Spec.TransitGateway == nil {
		return nil
	}
	if _, globalRouting, _ := genUtil.GetTransitGatewayLocationAndRouting(cluster.Spec.Zone, cluster.Spec.VPC.Region); cluster.Spec.TransitGateway.GlobalRouting != nil && !*cluster.Spec.TransitGateway.GlobalRouting && globalRouting != nil && *globalRouting {
		return field.Invalid(field.NewPath("spec.transitGateway.globalRouting"), cluster.Spec.TransitGateway.GlobalRouting, "global routing is required since PowerVS and VPC region are from different region")
	}
	return nil
}

func validateIBMPowerVSClusterCreateInfraPrereq(cluster *infrav1beta2.IBMPowerVSCluster) (allErrs field.ErrorList) {
	annotations := cluster.GetAnnotations()
	if len(annotations) == 0 {
		return nil
	}

	value, found := annotations[infrav1beta2.CreateInfrastructureAnnotation]
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
