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
	"net"
	"strings"

	"github.com/google/go-cmp/cmp"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/annotations"
)

const (
	warningClassicELB                = "%s load balancer is using a classic elb which is deprecated & support will be removed in a future release, please consider using another type of load balancer instead"
	warningHealthCheckProtocolNotSet = "healthcheck protocol is not set, the default value has changed from SSL to TCP. Health checks for existing clusters will be updated to TCP"
)

// log is for logging in this package.
var _ = ctrl.Log.WithName("awscluster-resource")

func (r *AWSCluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-awscluster,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsclusters,versions=v1beta2,name=validation.awscluster.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-awscluster,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsclusters,versions=v1beta2,name=default.awscluster.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

var (
	_ webhook.Validator = &AWSCluster{}
	_ webhook.Defaulter = &AWSCluster{}
)

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *AWSCluster) ValidateCreate() (admission.Warnings, error) {
	var allErrs field.ErrorList
	var allWarnings admission.Warnings

	allErrs = append(allErrs, r.Spec.Bastion.Validate()...)
	allErrs = append(allErrs, r.validateSSHKeyName()...)
	allErrs = append(allErrs, r.Spec.AdditionalTags.Validate()...)
	allErrs = append(allErrs, r.Spec.S3Bucket.Validate()...)
	allErrs = append(allErrs, r.validateNetwork()...)

	warnings, errs := r.validateControlPlaneLBs()
	if len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}
	if len(warnings) > 0 {
		allWarnings = append(allWarnings, warnings...)
	}

	return allWarnings, aggregateObjErrors(r.GroupVersionKind().GroupKind(), r.Name, allErrs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (r *AWSCluster) ValidateDelete() (admission.Warnings, error) {
	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *AWSCluster) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	var allErrs field.ErrorList
	var allWarnings admission.Warnings

	allErrs = append(allErrs, r.validateGCTasksAnnotation()...)

	oldC, ok := old.(*AWSCluster)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected an AWSCluster but got a %T", old))
	}

	if r.Spec.Region != oldC.Spec.Region {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "region"), r.Spec.Region, "field is immutable"),
		)
	}

	// Validate the control plane load balancers.
	lbs := map[*AWSLoadBalancerSpec]*AWSLoadBalancerSpec{
		oldC.Spec.ControlPlaneLoadBalancer:          r.Spec.ControlPlaneLoadBalancer,
		oldC.Spec.SecondaryControlPlaneLoadBalancer: r.Spec.SecondaryControlPlaneLoadBalancer,
	}

	for oldLB, newLB := range lbs {
		if oldLB == nil && newLB == nil {
			continue
		}

		allErrs = append(allErrs, r.validateControlPlaneLoadBalancerUpdate(oldLB, newLB)...)
	}

	if !cmp.Equal(oldC.Spec.ControlPlaneEndpoint, clusterv1.APIEndpoint{}) &&
		!cmp.Equal(r.Spec.ControlPlaneEndpoint, oldC.Spec.ControlPlaneEndpoint) {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "controlPlaneEndpoint"), r.Spec.ControlPlaneEndpoint, "field is immutable"),
		)
	}

	// Modifying VPC id is not allowed because it will cause a new VPC creation if set to nil.
	if !cmp.Equal(oldC.Spec.NetworkSpec, NetworkSpec{}) &&
		!cmp.Equal(oldC.Spec.NetworkSpec.VPC, VPCSpec{}) &&
		oldC.Spec.NetworkSpec.VPC.ID != "" {
		if cmp.Equal(r.Spec.NetworkSpec, NetworkSpec{}) ||
			cmp.Equal(r.Spec.NetworkSpec.VPC, VPCSpec{}) ||
			oldC.Spec.NetworkSpec.VPC.ID != r.Spec.NetworkSpec.VPC.ID {
			allErrs = append(allErrs,
				field.Invalid(field.NewPath("spec", "network", "vpc", "id"),
					r.Spec.NetworkSpec.VPC.ID, "field cannot be modified once set"))
		}
	}

	// If a identityRef is already set, do not allow removal of it.
	if oldC.Spec.IdentityRef != nil && r.Spec.IdentityRef == nil {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "identityRef"),
				r.Spec.IdentityRef, "field cannot be set to nil"),
		)
	}

	if annotations.IsExternallyManaged(oldC) && !annotations.IsExternallyManaged(r) {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("metadata", "annotations"),
				r.Annotations, "removal of externally managed annotation is not allowed"),
		)
	}

	allErrs = append(allErrs, r.Spec.Bastion.Validate()...)
	allErrs = append(allErrs, r.Spec.AdditionalTags.Validate()...)
	allErrs = append(allErrs, r.Spec.S3Bucket.Validate()...)

	if r.Spec.ControlPlaneLoadBalancer != nil {
		if r.Spec.ControlPlaneLoadBalancer.LoadBalancerType == LoadBalancerTypeClassic {
			allWarnings = append(allWarnings, fmt.Sprintf(warningClassicELB, "primary control plane"))
		}
	}

	if r.Spec.SecondaryControlPlaneLoadBalancer != nil {
		if r.Spec.SecondaryControlPlaneLoadBalancer.LoadBalancerType == LoadBalancerTypeClassic {
			allWarnings = append(allWarnings, fmt.Sprintf(warningClassicELB, "secondary control plane"))
		}
	}

	if r.Spec.ControlPlaneLoadBalancer == nil || r.Spec.ControlPlaneLoadBalancer.HealthCheckProtocol == nil {
		allWarnings = append(allWarnings, fmt.Sprintf("%s. Existing load balancers will be updates", warningHealthCheckProtocolNotSet))
	}

	return allWarnings, aggregateObjErrors(r.GroupVersionKind().GroupKind(), r.Name, allErrs)
}

func (r *AWSCluster) validateControlPlaneLoadBalancerUpdate(oldlb, newlb *AWSLoadBalancerSpec) field.ErrorList {
	var allErrs field.ErrorList

	if oldlb == nil {
		// If old scheme was nil, the only value accepted here is the default value: internet-facing
		if newlb.Scheme != nil && newlb.Scheme.String() != ELBSchemeInternetFacing.String() {
			allErrs = append(allErrs,
				field.Invalid(field.NewPath("spec", "controlPlaneLoadBalancer", "scheme"),
					newlb.Scheme, "field is immutable, default value was set to internet-facing"),
			)
		}
	} else {
		// A disabled Load Balancer has many implications that must be treated as immutable/
		// this is mostly used by externally managed Control Plane, and there's no need to support type changes.
		// More info: https://kubernetes.slack.com/archives/CD6U2V71N/p1708983246100859?thread_ts=1708973478.410979&cid=CD6U2V71N
		if (oldlb.LoadBalancerType == LoadBalancerTypeDisabled && newlb.LoadBalancerType != LoadBalancerTypeDisabled) ||
			(newlb.LoadBalancerType == LoadBalancerTypeDisabled && oldlb.LoadBalancerType != LoadBalancerTypeDisabled) {
			allErrs = append(allErrs,
				field.Invalid(field.NewPath("spec", "controlPlaneLoadBalancer", "type"),
					newlb.Scheme, "field is immutable when created of disabled type"),
			)
		}
		// If old scheme was not nil, the new scheme should be the same.
		if !cmp.Equal(oldlb.Scheme, newlb.Scheme) {
			allErrs = append(allErrs,
				field.Invalid(field.NewPath("spec", "controlPlaneLoadBalancer", "scheme"),
					newlb.Scheme, "field is immutable"),
			)
		}
		// The name must be defined when the AWSCluster is created. If it is not defined,
		// then the controller generates a default name at runtime, but does not store it,
		// so the name remains nil. In either case, the name cannot be changed.
		if !cmp.Equal(oldlb.Name, newlb.Name) {
			allErrs = append(allErrs,
				field.Invalid(field.NewPath("spec", "controlPlaneLoadBalancer", "name"),
					newlb.Name, "field is immutable"),
			)
		}

		// Block the update for Protocol :
		// - if it was not set in old spec but added in new spec
		// - if it was set in old spec but changed in new spec
		if oldlb.LoadBalancerType != LoadBalancerTypeClassic {
			if !cmp.Equal(newlb.HealthCheckProtocol, oldlb.HealthCheckProtocol) {
				allErrs = append(allErrs,
					field.Invalid(field.NewPath("spec", "controlPlaneLoadBalancer", "healthCheckProtocol"),
						newlb.HealthCheckProtocol, "field is immutable once set"),
				)
			}
		}
	}

	return allErrs
}

// Default satisfies the defaulting webhook interface.
func (r *AWSCluster) Default() {
	SetObjectDefaults_AWSCluster(r)
}

func (r *AWSCluster) validateGCTasksAnnotation() field.ErrorList {
	var allErrs field.ErrorList

	annotations := r.GetAnnotations()
	if annotations == nil {
		return nil
	}

	if gcTasksAnnotationValue := annotations[ExternalResourceGCTasksAnnotation]; gcTasksAnnotationValue != "" {
		gcTasks := strings.Split(gcTasksAnnotationValue, ",")

		supportedGCTasks := []GCTask{GCTaskLoadBalancer, GCTaskTargetGroup, GCTaskSecurityGroup}

		for _, gcTask := range gcTasks {
			found := false

			for _, supportedGCTask := range supportedGCTasks {
				if gcTask == string(supportedGCTask) {
					found = true
					break
				}
			}

			if !found {
				allErrs = append(allErrs,
					field.Invalid(field.NewPath("metadata", "annotations"),
						r.Annotations,
						fmt.Sprintf("annotation %s contains unsupported GC task %s", ExternalResourceGCTasksAnnotation, gcTask)),
				)
			}
		}
	}

	return allErrs
}

func (r *AWSCluster) validateSSHKeyName() field.ErrorList {
	return validateSSHKeyName(r.Spec.SSHKeyName)
}

func (r *AWSCluster) validateNetwork() field.ErrorList {
	var allErrs field.ErrorList
	if r.Spec.NetworkSpec.VPC.IsIPv6Enabled() {
		allErrs = append(allErrs, field.Invalid(field.NewPath("ipv6"), r.Spec.NetworkSpec.VPC.IPv6, "IPv6 cannot be used with unmanaged clusters at this time."))
	}
	for _, subnet := range r.Spec.NetworkSpec.Subnets {
		if subnet.IsIPv6 || subnet.IPv6CidrBlock != "" {
			allErrs = append(allErrs, field.Invalid(field.NewPath("subnets"), r.Spec.NetworkSpec.Subnets, "IPv6 cannot be used with unmanaged clusters at this time."))
		}
		if subnet.ZoneType != nil && subnet.IsEdge() {
			if subnet.ParentZoneName == nil {
				allErrs = append(allErrs, field.Invalid(field.NewPath("subnets"), r.Spec.NetworkSpec.Subnets, "ParentZoneName must be set when ZoneType is 'local-zone'."))
			}
		}
	}

	if r.Spec.NetworkSpec.VPC.CidrBlock != "" && r.Spec.NetworkSpec.VPC.IPAMPool != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("cidrBlock"), r.Spec.NetworkSpec.VPC.CidrBlock, "cidrBlock and ipamPool cannot be used together"))
	}

	if r.Spec.NetworkSpec.VPC.IPAMPool != nil && r.Spec.NetworkSpec.VPC.IPAMPool.ID == "" && r.Spec.NetworkSpec.VPC.IPAMPool.Name == "" {
		allErrs = append(allErrs, field.Invalid(field.NewPath("ipamPool"), r.Spec.NetworkSpec.VPC.IPAMPool, "ipamPool must have either id or name"))
	}

	for _, rule := range r.Spec.NetworkSpec.AdditionalControlPlaneIngressRules {
		allErrs = append(allErrs, r.validateIngressRule(rule)...)
	}

	for cidrBlockIndex, cidrBlock := range r.Spec.NetworkSpec.NodePortIngressRuleCidrBlocks {
		if _, _, err := net.ParseCIDR(cidrBlock); err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "network", fmt.Sprintf("nodePortIngressRuleCidrBlocks[%d]", cidrBlockIndex)), r.Spec.NetworkSpec.NodePortIngressRuleCidrBlocks, "CIDR block is invalid"))
		}
	}

	if r.Spec.NetworkSpec.VPC.ElasticIPPool != nil {
		eipp := r.Spec.NetworkSpec.VPC.ElasticIPPool
		if eipp.PublicIpv4Pool != nil {
			if eipp.PublicIpv4PoolFallBackOrder == nil {
				return append(allErrs, field.Invalid(field.NewPath("elasticIpPool.publicIpv4PoolFallbackOrder"), r.Spec.NetworkSpec.VPC.ElasticIPPool, "publicIpv4PoolFallbackOrder must be set when publicIpv4Pool is defined."))
			}
			awsPublicIpv4PoolPrefix := "ipv4pool-ec2-"
			if !strings.HasPrefix(*eipp.PublicIpv4Pool, awsPublicIpv4PoolPrefix) {
				return append(allErrs, field.Invalid(field.NewPath("elasticIpPool.publicIpv4Pool"), r.Spec.NetworkSpec.VPC.ElasticIPPool, fmt.Sprintf("publicIpv4Pool must start with %s.", awsPublicIpv4PoolPrefix)))
			}
		}
		if eipp.PublicIpv4Pool == nil && eipp.PublicIpv4PoolFallBackOrder != nil {
			return append(allErrs, field.Invalid(field.NewPath("elasticIpPool.publicIpv4PoolFallbackOrder"), r.Spec.NetworkSpec.VPC.ElasticIPPool, "publicIpv4Pool must be set when publicIpv4PoolFallbackOrder is defined."))
		}
	}

	secondaryCidrBlocks := r.Spec.NetworkSpec.VPC.SecondaryCidrBlocks
	secondaryCidrBlocksField := field.NewPath("spec", "network", "vpc", "secondaryCidrBlocks")
	for _, cidrBlock := range secondaryCidrBlocks {
		if r.Spec.NetworkSpec.VPC.CidrBlock != "" && r.Spec.NetworkSpec.VPC.CidrBlock == cidrBlock.IPv4CidrBlock {
			allErrs = append(allErrs, field.Invalid(secondaryCidrBlocksField, secondaryCidrBlocks, fmt.Sprintf("AWSCluster.spec.network.vpc.secondaryCidrBlocks must not contain the primary AWSCluster.spec.network.vpc.cidrBlock %v", r.Spec.NetworkSpec.VPC.CidrBlock)))
		}
	}

	return allErrs
}

func (r *AWSCluster) validateControlPlaneLBs() (admission.Warnings, field.ErrorList) {
	var allErrs field.ErrorList
	var allWarnings admission.Warnings

	if r.Spec.ControlPlaneLoadBalancer != nil && r.Spec.ControlPlaneLoadBalancer.LoadBalancerType == LoadBalancerTypeClassic {
		allWarnings = append(allWarnings, fmt.Sprintf(warningClassicELB, "primary control plane"))

		if r.Spec.ControlPlaneLoadBalancer.HealthCheckProtocol == nil {
			allWarnings = append(allWarnings, warningHealthCheckProtocolNotSet)
		}

		if r.Spec.ControlPlaneLoadBalancer.HealthCheckProtocol != nil && *r.Spec.ControlPlaneLoadBalancer.HealthCheckProtocol == ELBProtocolSSL {
			allWarnings = append(allWarnings, "loadbalancer is using a classic elb with SSL health check, this causes issues with ciper suites with kubernetes v1.30+")
		}
	}

	// If the secondary is defined, check that the name is not empty and different from the primary.
	// Also, ensure that the secondary load balancer is an NLB
	if r.Spec.SecondaryControlPlaneLoadBalancer != nil {
		if r.Spec.SecondaryControlPlaneLoadBalancer.Name == nil || *r.Spec.SecondaryControlPlaneLoadBalancer.Name == "" {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "secondaryControlPlaneLoadBalancer", "name"), r.Spec.SecondaryControlPlaneLoadBalancer.Name, "secondary controlPlaneLoadBalancer.name cannot be empty"))
		}

		if r.Spec.SecondaryControlPlaneLoadBalancer.Name == r.Spec.ControlPlaneLoadBalancer.Name {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "secondaryControlPlaneLoadBalancer", "name"), r.Spec.SecondaryControlPlaneLoadBalancer.Name, "field must be different from controlPlaneLoadBalancer.name"))
		}

		if r.Spec.SecondaryControlPlaneLoadBalancer.Scheme.Equals(r.Spec.ControlPlaneLoadBalancer.Scheme) {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "secondaryControlPlaneLoadBalancer", "scheme"), r.Spec.SecondaryControlPlaneLoadBalancer.Scheme, "control plane load balancers must have different schemes"))
		}

		if r.Spec.SecondaryControlPlaneLoadBalancer.LoadBalancerType != LoadBalancerTypeNLB {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "secondaryControlPlaneLoadBalancer", "loadBalancerType"), r.Spec.SecondaryControlPlaneLoadBalancer.LoadBalancerType, "secondary control plane load balancer must be a Network Load Balancer"))
		}
		if r.Spec.SecondaryControlPlaneLoadBalancer.LoadBalancerType == LoadBalancerTypeClassic {
			allWarnings = append(allWarnings, fmt.Sprintf(warningClassicELB, "secondary control plane"))
		}
	}

	// Additional listeners are only supported for NLBs.
	// Validate the control plane load balancers.
	loadBalancers := []*AWSLoadBalancerSpec{
		r.Spec.ControlPlaneLoadBalancer,
		r.Spec.SecondaryControlPlaneLoadBalancer,
	}
	for _, cp := range loadBalancers {
		if cp == nil {
			continue
		}

		for _, rule := range cp.IngressRules {
			allErrs = append(allErrs, r.validateIngressRule(rule)...)
		}
	}

	if r.Spec.ControlPlaneLoadBalancer.LoadBalancerType == LoadBalancerTypeDisabled {
		if r.Spec.ControlPlaneLoadBalancer.Name != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "controlPlaneLoadBalancer", "name"), r.Spec.ControlPlaneLoadBalancer.Name, "cannot configure a name if the LoadBalancer reconciliation is disabled"))
		}

		if r.Spec.ControlPlaneLoadBalancer.CrossZoneLoadBalancing {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "controlPlaneLoadBalancer", "crossZoneLoadBalancing"), r.Spec.ControlPlaneLoadBalancer.CrossZoneLoadBalancing, "cross-zone load balancing cannot be set if the LoadBalancer reconciliation is disabled"))
		}

		if len(r.Spec.ControlPlaneLoadBalancer.Subnets) > 0 {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "controlPlaneLoadBalancer", "subnets"), r.Spec.ControlPlaneLoadBalancer.Subnets, "subnets cannot be set if the LoadBalancer reconciliation is disabled"))
		}

		if r.Spec.ControlPlaneLoadBalancer.HealthCheckProtocol != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "controlPlaneLoadBalancer", "healthCheckProtocol"), r.Spec.ControlPlaneLoadBalancer.HealthCheckProtocol, "healthcheck protocol cannot be set if the LoadBalancer reconciliation is disabled"))
		}

		if len(r.Spec.ControlPlaneLoadBalancer.AdditionalSecurityGroups) > 0 {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "controlPlaneLoadBalancer", "additionalSecurityGroups"), r.Spec.ControlPlaneLoadBalancer.AdditionalSecurityGroups, "additional Security Groups cannot be set if the LoadBalancer reconciliation is disabled"))
		}

		if len(r.Spec.ControlPlaneLoadBalancer.AdditionalListeners) > 0 {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "controlPlaneLoadBalancer", "additionalListeners"), r.Spec.ControlPlaneLoadBalancer.AdditionalListeners, "cannot set additional listeners if the LoadBalancer reconciliation is disabled"))
		}

		if len(r.Spec.ControlPlaneLoadBalancer.IngressRules) > 0 {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "controlPlaneLoadBalancer", "ingressRules"), r.Spec.ControlPlaneLoadBalancer.IngressRules, "ingress rules cannot be set if the LoadBalancer reconciliation is disabled"))
		}

		if r.Spec.ControlPlaneLoadBalancer.PreserveClientIP {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "controlPlaneLoadBalancer", "preserveClientIP"), r.Spec.ControlPlaneLoadBalancer.PreserveClientIP, "cannot preserve client IP if the LoadBalancer reconciliation is disabled"))
		}

		if r.Spec.ControlPlaneLoadBalancer.DisableHostsRewrite {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "controlPlaneLoadBalancer", "disableHostsRewrite"), r.Spec.ControlPlaneLoadBalancer.DisableHostsRewrite, "cannot disable hosts rewrite if the LoadBalancer reconciliation is disabled"))
		}
	}

	return allWarnings, allErrs
}

func (r *AWSCluster) validateIngressRule(rule IngressRule) field.ErrorList {
	var allErrs field.ErrorList
	if rule.NatGatewaysIPsSource {
		if rule.CidrBlocks != nil || rule.IPv6CidrBlocks != nil || rule.SourceSecurityGroupIDs != nil || rule.SourceSecurityGroupRoles != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("additionalControlPlaneIngressRules"), r.Spec.NetworkSpec.AdditionalControlPlaneIngressRules, "CIDR blocks and security group IDs or security group roles cannot be used together"))
		}
	} else {
		if (rule.CidrBlocks != nil || rule.IPv6CidrBlocks != nil) && (rule.SourceSecurityGroupIDs != nil || rule.SourceSecurityGroupRoles != nil) {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "controlPlaneLoadBalancer", "ingressRules"), r.Spec.ControlPlaneLoadBalancer.IngressRules, "CIDR blocks and security group IDs or security group roles cannot be used together"))
		}
	}
	return allErrs
}
