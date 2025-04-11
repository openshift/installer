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

package v1beta2

import (
	"fmt"
	"strconv"

	regionUtil "github.com/ppc64le-cloud/powervs-utils"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"

	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	genUtil "sigs.k8s.io/cluster-api-provider-ibmcloud/util"
)

// log is for logging in this package.
var ibmpowervsclusterlog = logf.Log.WithName("ibmpowervscluster-resource")

func (r *IBMPowerVSCluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-ibmpowervscluster,mutating=true,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsclusters,verbs=create;update,versions=v1beta2,name=mibmpowervscluster.kb.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

var _ webhook.Defaulter = &IBMPowerVSCluster{}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (r *IBMPowerVSCluster) Default() {
	ibmpowervsclusterlog.Info("default", "name", r.Name)
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-ibmpowervscluster,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsclusters,versions=v1beta2,name=vibmpowervscluster.kb.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

var _ webhook.Validator = &IBMPowerVSCluster{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMPowerVSCluster) ValidateCreate() (admission.Warnings, error) {
	ibmpowervsclusterlog.Info("validate create", "name", r.Name)
	return r.validateIBMPowerVSCluster()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMPowerVSCluster) ValidateUpdate(_ runtime.Object) (admission.Warnings, error) {
	ibmpowervsclusterlog.Info("validate update", "name", r.Name)
	return r.validateIBMPowerVSCluster()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMPowerVSCluster) ValidateDelete() (admission.Warnings, error) {
	ibmpowervsclusterlog.Info("validate delete", "name", r.Name)
	return nil, nil
}

func (r *IBMPowerVSCluster) validateIBMPowerVSCluster() (admission.Warnings, error) {
	var allErrs field.ErrorList
	if err := r.validateIBMPowerVSClusterNetwork(); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := r.validateIBMPowerVSClusterCreateInfraPrereq(); err != nil {
		allErrs = append(allErrs, err...)
	}

	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(
		schema.GroupKind{Group: "infrastructure.cluster.x-k8s.io", Kind: "IBMPowerVSCluster"},
		r.Name, allErrs)
}

func (r *IBMPowerVSCluster) validateIBMPowerVSClusterNetwork() *field.Error {
	if res, err := validateIBMPowerVSNetworkReference(r.Spec.Network); !res {
		return err
	}
	return nil
}

func (r *IBMPowerVSCluster) validateIBMPowerVSClusterLoadBalancers() (allErrs field.ErrorList) {
	if err := r.validateIBMPowerVSClusterLoadBalancerNames(); err != nil {
		allErrs = append(allErrs, err...)
	}

	if len(r.Spec.LoadBalancers) == 0 {
		return allErrs
	}

	for _, loadbalancer := range r.Spec.LoadBalancers {
		if *loadbalancer.Public {
			return allErrs
		}
	}

	return append(allErrs, field.Invalid(field.NewPath("spec.LoadBalancers"), r.Spec.LoadBalancers, "Expect atleast one of the load balancer to be public"))
}

func (r *IBMPowerVSCluster) validateIBMPowerVSClusterLoadBalancerNames() (allErrs field.ErrorList) {
	found := make(map[string]bool)
	for i, loadbalancer := range r.Spec.LoadBalancers {
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

func (r *IBMPowerVSCluster) validateIBMPowerVSClusterVPCSubnetNames() (allErrs field.ErrorList) {
	found := make(map[string]bool)
	for i, subnet := range r.Spec.VPCSubnets {
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

func (r *IBMPowerVSCluster) validateIBMPowerVSClusterTransitGateway() *field.Error {
	if r.Spec.Zone == nil && r.Spec.VPC == nil {
		return nil
	}
	if r.Spec.TransitGateway == nil {
		return nil
	}
	if _, globalRouting, _ := genUtil.GetTransitGatewayLocationAndRouting(r.Spec.Zone, r.Spec.VPC.Region); r.Spec.TransitGateway.GlobalRouting != nil && !*r.Spec.TransitGateway.GlobalRouting && globalRouting != nil && *globalRouting {
		return field.Invalid(field.NewPath("spec.transitGateway.globalRouting"), r.Spec.TransitGateway.GlobalRouting, "global routing is required since PowerVS and VPC region are from different region")
	}
	return nil
}

func (r *IBMPowerVSCluster) validateIBMPowerVSClusterCreateInfraPrereq() (allErrs field.ErrorList) {
	annotations := r.GetAnnotations()
	if len(annotations) == 0 {
		return nil
	}

	value, found := annotations[CreateInfrastructureAnnotation]
	if !found {
		return nil
	}

	createInfra, err := strconv.ParseBool(value)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("annotations"), r.Annotations, "value of powervs.cluster.x-k8s.io/create-infra should be boolean"))
	}

	if !createInfra {
		return nil
	}

	if r.Spec.Zone == nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec.zone"), r.Spec.Zone, "value of zone is empty"))
	}

	if r.Spec.Zone != nil && !regionUtil.ValidateZone(*r.Spec.Zone) {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec.zone"), r.Spec.Zone, fmt.Sprintf("zone '%s' is not supported", *r.Spec.Zone)))
	}

	if r.Spec.VPC == nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec.vpc"), r.Spec.VPC, "value of VPC is empty"))
	}

	if r.Spec.VPC != nil && r.Spec.VPC.Region == nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec.vpc.region"), r.Spec.VPC.Region, "value of VPC region is empty"))
	}

	if r.Spec.VPC != nil && r.Spec.VPC.Region != nil && !regionUtil.ValidateVPCRegion(*r.Spec.VPC.Region) {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec.vpc.region"), r.Spec.VPC.Region, fmt.Sprintf("vpc region '%s' is not supported", *r.Spec.VPC.Region)))
	}

	if r.Spec.ResourceGroup == nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec.resourceGroup"), r.Spec.ResourceGroup, "value of resource group is empty"))
	}
	if err := r.validateIBMPowerVSClusterVPCSubnetNames(); err != nil {
		allErrs = append(allErrs, err...)
	}

	if err := r.validateIBMPowerVSClusterLoadBalancers(); err != nil {
		allErrs = append(allErrs, err...)
	}

	if err := r.validateIBMPowerVSClusterTransitGateway(); err != nil {
		allErrs = append(allErrs, err)
	}

	return allErrs
}
