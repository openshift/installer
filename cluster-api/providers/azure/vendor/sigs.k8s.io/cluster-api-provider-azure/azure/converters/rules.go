/*
Copyright 2020 The Kubernetes Authors.

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

package converters

import (
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v4"
	"k8s.io/utils/ptr"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
)

// SecurityRuleToSDK converts a CAPZ security rule to an Azure network security rule.
func SecurityRuleToSDK(rule infrav1.SecurityRule) *armnetwork.SecurityRule {
	secRule := &armnetwork.SecurityRule{
		Name: ptr.To(rule.Name),
		Properties: &armnetwork.SecurityRulePropertiesFormat{
			Description:              ptr.To(rule.Description),
			SourceAddressPrefix:      rule.Source,
			SourcePortRange:          rule.SourcePorts,
			DestinationAddressPrefix: rule.Destination,
			DestinationPortRange:     rule.DestinationPorts,
			Access:                   ptr.To(armnetwork.SecurityRuleAccess(rule.Action)),
			Priority:                 ptr.To[int32](rule.Priority),
		},
	}

	switch rule.Protocol {
	case infrav1.SecurityGroupProtocolAll:
		secRule.Properties.Protocol = ptr.To(armnetwork.SecurityRuleProtocolAsterisk)
	case infrav1.SecurityGroupProtocolTCP:
		secRule.Properties.Protocol = ptr.To(armnetwork.SecurityRuleProtocolTCP)
	case infrav1.SecurityGroupProtocolUDP:
		secRule.Properties.Protocol = ptr.To(armnetwork.SecurityRuleProtocolUDP)
	case infrav1.SecurityGroupProtocolICMP:
		secRule.Properties.Protocol = ptr.To(armnetwork.SecurityRuleProtocolIcmp)
	}

	switch rule.Direction {
	case infrav1.SecurityRuleDirectionOutbound:
		secRule.Properties.Direction = ptr.To(armnetwork.SecurityRuleDirectionOutbound)
	case infrav1.SecurityRuleDirectionInbound:
		secRule.Properties.Direction = ptr.To(armnetwork.SecurityRuleDirectionInbound)
	}

	return secRule
}
