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

package powervs

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	regionUtil "github.com/ppc64le-cloud/powervs-utils"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/powervs/v1beta3"
	ibmcloudutil "sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/util"
)

// Ensure IBMPowerVSCluster implements the typed webhook interfaces.
var (
	_ admission.Validator[*infrav1.IBMPowerVSCluster] = &IBMPowerVSCluster{}
	_ admission.Defaulter[*infrav1.IBMPowerVSCluster] = &IBMPowerVSCluster{}
)

const (
	infrastructureGroup = "infrastructure.cluster.x-k8s.io"

	// lbResourceNameMaxLen is the IBM Cloud resource name character limit.
	// Mirrors resourceNameMaxLen in pkg/cloud/scope/powervs/resource_name.go.
	lbResourceNameMaxLen = 63
)

//+kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta3-ibmpowervscluster,mutating=true,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsclusters,versions=v1beta3,name=mibmpowervscluster.kb.io,sideEffects=None,admissionReviewVersions=v1
//+kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta3-ibmpowervscluster,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsclusters,versions=v1beta3,name=vibmpowervscluster.kb.io,sideEffects=None,admissionReviewVersions=v1

func (r *IBMPowerVSCluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &infrav1.IBMPowerVSCluster{}).
		WithValidator(r).
		WithDefaulter(r).
		Complete()
}

// IBMPowerVSCluster implements a validation and defaulting webhook for IBMPowerVSCluster.
type IBMPowerVSCluster struct{}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (r *IBMPowerVSCluster) Default(_ context.Context, _ *infrav1.IBMPowerVSCluster) error {
	return nil
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMPowerVSCluster) ValidateCreate(_ context.Context, obj *infrav1.IBMPowerVSCluster) (admission.Warnings, error) {
	return validateIBMPowerVSCluster(nil, obj)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMPowerVSCluster) ValidateUpdate(_ context.Context, oldObj, newObj *infrav1.IBMPowerVSCluster) (warnings admission.Warnings, err error) {
	return validateIBMPowerVSCluster(oldObj, newObj)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMPowerVSCluster) ValidateDelete(_ context.Context, _ *infrav1.IBMPowerVSCluster) (admission.Warnings, error) {
	return nil, nil
}

func validateIBMPowerVSCluster(oldCluster, newCluster *infrav1.IBMPowerVSCluster) (admission.Warnings, error) {
	var allErrs field.ErrorList
	allErrs = append(allErrs, validateIBMPowerVSClusterCreateInfraPrereq(newCluster)...)
	// validateAdditionalListenerSelector is only meaningful on update, not create.
	if oldCluster != nil {
		allErrs = append(allErrs, validateAdditionalListenerSelector(newCluster, oldCluster)...)
	}

	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(
		schema.GroupKind{Group: infrastructureGroup, Kind: "IBMPowerVSCluster"},
		newCluster.Name, allErrs)
}

func validateIBMPowerVSClusterLoadBalancers(cluster *infrav1.IBMPowerVSCluster) (allErrs field.ErrorList) {
	allErrs = append(allErrs, validateIBMPowerVSClusterLoadBalancerNames(cluster)...)

	if len(cluster.Spec.LoadBalancers) == 0 {
		return allErrs
	}

	for _, loadBalancer := range cluster.Spec.LoadBalancers {
		if loadBalancer.Type == infrav1.SourceTypeProvision &&
			(loadBalancer.Provision.Type == infrav1.LoadBalancerTypePublic || loadBalancer.Provision.Type == "") {
			return allErrs
		}
	}

	return append(allErrs, field.Invalid(field.NewPath("spec").Child("loadBalancers"), cluster.Spec.LoadBalancers, "at least one load balancer must be public"))
}

func validateIBMPowerVSClusterLoadBalancerNames(cluster *infrav1.IBMPowerVSCluster) (allErrs field.ErrorList) {
	found := make(map[string]bool)
	for i, loadbalancer := range cluster.Spec.LoadBalancers {
		name := ""
		switch loadbalancer.Type {
		case infrav1.SourceTypeReference:
			name = loadbalancer.Reference.Name
		case infrav1.SourceTypeProvision:
			name = loadbalancer.Provision.Name
		}
		if name == "" {
			continue
		}

		if found[name] {
			allErrs = append(allErrs, field.Duplicate(field.NewPath("spec", "loadBalancers").Index(i).Child("name"), name))
			continue
		}
		found[name] = true
	}

	return allErrs
}

func validateIBMPowerVSClusterVPCSubnetNames(cluster *infrav1.IBMPowerVSCluster) (allErrs field.ErrorList) {
	found := make(map[string]bool)
	for i, subnet := range cluster.Spec.VPCSubnets {
		name := ""
		switch subnet.Type {
		case infrav1.SourceTypeReference:
			name = subnet.Reference.Name
		case infrav1.SourceTypeProvision:
			name = subnet.Provision.Name
		}
		if name == "" {
			continue
		}
		if found[name] {
			allErrs = append(allErrs, field.Duplicate(field.NewPath("spec", "subnets").Index(i).Child("name"), name))
			continue
		}
		found[name] = true
	}

	return allErrs
}

func validateIBMPowerVSClusterTransitGateway(cluster *infrav1.IBMPowerVSCluster) *field.Error {
	if cluster.Spec.Zone == "" || cluster.Spec.VPC.Region == "" {
		return nil
	}
	// TransitGateway is now a value type, check if Type is set to determine if it's configured
	if cluster.Spec.TransitGateway.Type == "" {
		return nil
	}
	// GlobalRouting is now in Provision field and is a string enum, not a bool pointer
	if cluster.Spec.TransitGateway.Type == infrav1.SourceTypeProvision {
		if _, globalRouting, _ := ibmcloudutil.GetTransitGatewayLocationAndRouting(&cluster.Spec.Zone, &cluster.Spec.VPC.Region); cluster.Spec.TransitGateway.Provision.GlobalRouting == infrav1.TransitGatewayRoutingLocal && globalRouting != nil && *globalRouting {
			return field.Invalid(field.NewPath("spec", "transitGateway", "provision", "globalRouting"), cluster.Spec.TransitGateway.Provision.GlobalRouting, "global routing is required since PowerVS and VPC region are from different region")
		}
	}
	return nil
}

// validateIBMPowerVSClusterCreateInfraPrereq validates the prerequisites required when
// Topology is LoadBalancer, which is the v1beta3 signal that infrastructure should be provisioned.
func validateIBMPowerVSClusterCreateInfraPrereq(cluster *infrav1.IBMPowerVSCluster) (allErrs field.ErrorList) {
	if cluster.Spec.Topology != infrav1.PowerVSLoadBalancerTopology {
		return nil
	}

	// CRD CEL rules enforce zone presence and resourceGroup validity.
	// ValidateZone additionally checks the value against the known PowerVS zone list, which CEL cannot express.
	if cluster.Spec.Zone != "" && !regionUtil.ValidateZone(cluster.Spec.Zone) {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "zone"), cluster.Spec.Zone, fmt.Sprintf("zone %q is not supported", cluster.Spec.Zone)))
	}

	// CRD field markers enforce VPC type and region presence.
	// ValidateVPCRegion additionally checks the value against the known IBM Cloud VPC region list, which CEL cannot express.
	if cluster.Spec.VPC.Region != "" && !regionUtil.ValidateVPCRegion(cluster.Spec.VPC.Region) {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "vpc", "region"), cluster.Spec.VPC.Region, fmt.Sprintf("vpc region %q is not supported", cluster.Spec.VPC.Region)))
	}

	allErrs = append(allErrs, validateIBMPowerVSClusterVPCSubnetNames(cluster)...)
	allErrs = append(allErrs, validateIBMPowerVSClusterLoadBalancers(cluster)...)

	if err := validateIBMPowerVSClusterTransitGateway(cluster); err != nil {
		allErrs = append(allErrs, err)
	}

	return allErrs
}

func validateAdditionalListenerSelector(newCluster, oldCluster *infrav1.IBMPowerVSCluster) (allErrs field.ErrorList) {
	// Build a map keyed by (lbName, port, protocol) so that listeners from
	// different load balancers with the same port are not confused with each other.
	// For nameless provision LBs the name is derived via resolvedLBName, which
	// mirrors the controller's own naming logic, ensuring the key is stable and unique.
	type listenerKey struct {
		lbName   string
		port     int64
		protocol infrav1.LoadBalancerListenerProtocol
	}
	newListeners := map[listenerKey]metav1.LabelSelector{}
	for i, lb := range newCluster.Spec.LoadBalancers {
		if lb.Type != infrav1.SourceTypeProvision {
			continue
		}
		lbName := resolvedLBName(newCluster.Name, lb, i)
		for _, al := range lb.Provision.AdditionalListeners {
			newListeners[listenerKey{lbName: lbName, port: al.Port, protocol: al.Protocol}] = al.Selector
		}
	}
	for i, lb := range oldCluster.Spec.LoadBalancers {
		if lb.Type != infrav1.SourceTypeProvision {
			continue
		}
		lbName := resolvedLBName(oldCluster.Name, lb, i)
		for _, al := range lb.Provision.AdditionalListeners {
			key := listenerKey{lbName: lbName, port: al.Port, protocol: al.Protocol}
			if selector, ok := newListeners[key]; ok && !reflect.DeepEqual(selector, al.Selector) {
				allErrs = append(allErrs, field.Forbidden(
					field.NewPath("spec", "loadBalancers").Index(i),
					fmt.Sprintf("selector is immutable for load balancer %q port %d", lbName, al.Port)))
			}
		}
	}
	return allErrs
}

// lbResourceName builds a deterministic IBM Cloud resource name.
// It mirrors ResourceName in pkg/cloud/scope/powervs/resource_name.go; we
// cannot import that package directly because its test suite imports this
// webhook package, which would create an import cycle.
func lbResourceName(clusterName, suffix, qualifier string) string {
	if qualifier != "" {
		suffix = suffix + "-" + qualifier
	}
	maxBase := lbResourceNameMaxLen - len(suffix) - 1
	base := clusterName
	if len(base) > maxBase {
		base = strings.TrimRight(base[:maxBase], "-")
	}
	return base + "-" + suffix
}

// resolvedLBName returns the effective load balancer name for a provision entry,
// mirroring the same logic the controller uses in ReconcileLoadBalancers:
// if the spec name is empty, a deterministic name is derived from the cluster
// name, the LB type (public/private), and the list index.
func resolvedLBName(clusterName string, lb infrav1.LoadBalancerSource, i int) string {
	if lb.Provision.Name != "" {
		return lb.Provision.Name
	}
	suffix := "loadbalancer-public"
	if lb.Provision.Type == infrav1.LoadBalancerTypePrivate {
		suffix = "loadbalancer-private"
	}
	qualifier := ""
	if i > 0 {
		qualifier = strconv.Itoa(i)
	}
	return lbResourceName(clusterName, suffix, qualifier)
}
