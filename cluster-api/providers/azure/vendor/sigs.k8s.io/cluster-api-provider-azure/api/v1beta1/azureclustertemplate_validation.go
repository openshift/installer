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

package v1beta1

import (
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func (c *AzureClusterTemplate) validateClusterTemplate() (admission.Warnings, error) {
	var allErrs field.ErrorList
	allErrs = append(allErrs, c.validateClusterTemplateSpec()...)

	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(
		schema.GroupKind{Group: "infrastructure.cluster.x-k8s.io", Kind: "AzureClusterTemplate"},
		c.Name, allErrs)
}

func (c *AzureClusterTemplate) validateClusterTemplateSpec() field.ErrorList {
	var allErrs field.ErrorList

	allErrs = append(allErrs, validateVnetCIDR(
		c.Spec.Template.Spec.NetworkSpec.Vnet.CIDRBlocks,
		field.NewPath("spec").Child("template").Child("spec").
			Child("networkSpec").Child("vnet").Child("cidrBlocks"))...)

	allErrs = append(allErrs, validateSubnetTemplates(
		c.Spec.Template.Spec.NetworkSpec.Subnets,
		c.Spec.Template.Spec.NetworkSpec.Vnet,
		field.NewPath("spec").Child("template").Child("spec").Child("networkSpec").Child("subnets"),
	)...)

	allErrs = append(allErrs, c.validateAPIServerLB(
		field.NewPath("spec").Child("template").Child("spec").Child("networkSpec").Child("apiServerLB"),
	)...)

	allErrs = append(allErrs, c.validateNetworkSpec()...)

	allErrs = append(allErrs, c.validateControlPlaneOutboundLB()...)

	allErrs = append(allErrs, c.validatePrivateDNSZoneName()...)

	return allErrs
}

func (c *AzureClusterTemplate) validateNetworkSpec() field.ErrorList {
	var allErrs field.ErrorList

	var needOutboundLB bool
	networkSpec := c.Spec.Template.Spec.NetworkSpec
	for _, subnet := range networkSpec.Subnets {
		if subnet.Role == SubnetNode && subnet.IsIPv6Enabled() {
			needOutboundLB = true
			break
		}
	}
	if needOutboundLB {
		allErrs = append(allErrs, c.validateNodeOutboundLB()...)
	}

	return allErrs
}

func validateSubnetTemplates(subnets SubnetTemplatesSpec, vnet VnetTemplateSpec, fld *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	subnetNames := make(map[string]bool, len(subnets))
	requiredSubnetRoles := map[string]bool{
		"control-plane": false,
		"node":          false,
	}

	for i, subnet := range subnets {
		if err := validateSubnetName(subnet.Name, fld.Index(i).Child("name")); err != nil {
			allErrs = append(allErrs, err)
		}
		if _, ok := subnetNames[subnet.Name]; ok {
			allErrs = append(allErrs, field.Duplicate(fld, subnet.Name))
		}
		subnetNames[subnet.Name] = true
		for role := range requiredSubnetRoles {
			if role == string(subnet.Role) {
				requiredSubnetRoles[role] = true
			}
		}
		for j, rule := range subnet.SecurityGroup.SecurityRules {
			if err := validateSecurityRule(
				rule,
				fld.Index(i).Child("securityGroup").Child("securityGroup").Child("securityRules").Index(j),
			); err != nil {
				allErrs = append(allErrs, err...)
			}
		}
		allErrs = append(allErrs, validateSubnetCIDR(subnet.CIDRBlocks, vnet.CIDRBlocks, fld.Index(i).Child("cidrBlocks"))...)
	}
	for k, v := range requiredSubnetRoles {
		if !v {
			allErrs = append(allErrs, field.Required(fld,
				fmt.Sprintf("required role %s not included in provided subnets", k)))
		}
	}
	return allErrs
}

func (c *AzureClusterTemplate) validateAPIServerLB(apiServerLBPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	lb := c.Spec.Template.Spec.NetworkSpec.APIServerLB
	allErrs = append(allErrs, validateClassSpecForAPIServerLB(lb, nil, apiServerLBPath)...)
	return allErrs
}

func (c *AzureClusterTemplate) validateNodeOutboundLB() field.ErrorList {
	var allErrs field.ErrorList

	fldPath := field.NewPath("spec").Child("template").Child("spec").Child("networkSpec").Child("nodeOutboundLB")
	apiserverLB := c.Spec.Template.Spec.NetworkSpec.APIServerLB
	lb := c.Spec.Template.Spec.NetworkSpec.NodeOutboundLB

	allErrs = append(allErrs, validateClassSpecForNodeOutboundLB(lb, nil, apiserverLB, fldPath)...)

	return allErrs
}

func (c *AzureClusterTemplate) validateControlPlaneOutboundLB() field.ErrorList {
	var allErrs field.ErrorList

	fldPath := field.NewPath("spec").Child("template").Child("spec").Child("networkSpec").Child("controlPlaneOutboundLB")
	apiserverLB := c.Spec.Template.Spec.NetworkSpec.APIServerLB
	lb := c.Spec.Template.Spec.NetworkSpec.ControlPlaneOutboundLB

	allErrs = append(allErrs, validateClassSpecForControlPlaneOutboundLB(lb, apiserverLB, fldPath)...)

	return allErrs
}

func (c *AzureClusterTemplate) validatePrivateDNSZoneName() field.ErrorList {
	var allErrs field.ErrorList

	fldPath := field.NewPath("spec").Child("template").Child("spec").Child("networkSpec").Child("privateDNSZoneName")
	networkSpec := c.Spec.Template.Spec.NetworkSpec

	allErrs = append(allErrs, validatePrivateDNSZoneName(
		networkSpec.PrivateDNSZoneName,
		true,
		networkSpec.APIServerLB.Type,
		fldPath,
	)...)

	return allErrs
}
