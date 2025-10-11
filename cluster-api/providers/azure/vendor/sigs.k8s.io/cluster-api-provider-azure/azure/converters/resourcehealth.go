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

package converters

import (
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resourcehealth/armresourcehealth"
	"k8s.io/utils/ptr"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
)

// SDKAvailabilityStatusToCondition converts an Azure Resource Health availability status to a status condition.
func SDKAvailabilityStatusToCondition(availStatus armresourcehealth.AvailabilityStatus) *clusterv1.Condition {
	if availStatus.Properties == nil {
		return conditions.FalseCondition(infrav1.AzureResourceAvailableCondition, "", "", "")
	}

	state := availStatus.Properties.AvailabilityState

	if ptr.Deref(state, "") == armresourcehealth.AvailabilityStateValuesAvailable {
		return conditions.TrueCondition(infrav1.AzureResourceAvailableCondition)
	}

	var reason strings.Builder
	if availStatus.Properties.ReasonType != nil {
		// CAPI specifies Reason should be CamelCase, though the Azure API
		// response may include spaces (e.g. "Customer Initiated")
		words := strings.Split(*availStatus.Properties.ReasonType, " ")
		for _, word := range words {
			if word != "" {
				reason.WriteString(strings.ToTitle(word[:1]))
			}
			if len(word) > 1 {
				reason.WriteString(word[1:])
			}
		}
	}

	var severity clusterv1.ConditionSeverity
	switch ptr.Deref(availStatus.Properties.AvailabilityState, "") {
	case armresourcehealth.AvailabilityStateValuesUnavailable:
		severity = clusterv1.ConditionSeverityError
	case armresourcehealth.AvailabilityStateValuesDegraded, armresourcehealth.AvailabilityStateValuesUnknown:
		severity = clusterv1.ConditionSeverityWarning
	}

	var message string
	if availStatus.Properties.Summary != nil {
		message = *availStatus.Properties.Summary
	}

	return conditions.FalseCondition(infrav1.AzureResourceAvailableCondition, reason.String(), severity, "%s", message)
}
