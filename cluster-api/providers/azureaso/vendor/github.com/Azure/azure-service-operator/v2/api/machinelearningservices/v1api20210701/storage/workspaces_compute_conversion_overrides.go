// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package storage

import (
	v20240401 "github.com/Azure/azure-service-operator/v2/api/machinelearningservices/v1api20240401/storage"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

var _ augmentConversionForAKS_Properties = &AKS_Properties{}

func (aks *AKS_Properties) AssignPropertiesFrom(src *v20240401.AKS_Properties) error {
	if src.LoadBalancerSubnetReference != nil {
		isNotKubeRef := !src.LoadBalancerSubnetReference.IsKubernetesReference()
		// Note that using isNotKubeRef is a bit awkward because in reality it shouldn't be possible to have a genruntime.ResourceReference with no
		// kube ref AND no ARM ref, but if that does happen we pass the empty-string along to maintain the round-trip invariant.
		if len(src.LoadBalancerSubnetReference.ARMID) > 0 || isNotKubeRef {
			aks.LoadBalancerSubnet = &src.LoadBalancerSubnetReference.ARMID
			aks.PropertyBag.Remove("LoadBalancerSubnetReference") // Remove it from property bag added by code generated code
		}
	}

	return nil
}

func (aks *AKS_Properties) AssignPropertiesTo(dst *v20240401.AKS_Properties) error {
	// LoadBalancerSubnet
	if aks.LoadBalancerSubnet != nil {
		dst.LoadBalancerSubnetReference = &genruntime.ResourceReference{
			ARMID: *aks.LoadBalancerSubnet,
		}
	}

	return nil
}
