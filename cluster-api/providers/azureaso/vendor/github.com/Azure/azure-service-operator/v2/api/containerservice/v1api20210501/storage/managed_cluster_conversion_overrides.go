// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package storage

import (
	v20230201s "github.com/Azure/azure-service-operator/v2/api/containerservice/v1api20230201/storage"
	"github.com/Azure/azure-service-operator/v2/internal/util/to"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

var _ augmentConversionForManagedClusterAgentPoolProfile = &ManagedClusterAgentPoolProfile{}

func (profile *ManagedClusterAgentPoolProfile) AssignPropertiesFrom(src *v20230201s.ManagedClusterAgentPoolProfile) error {
	// Clone the existing property bag
	propertyBag := genruntime.NewPropertyBag(src.PropertyBag)

	// ProximityPlacementGroupIDReference
	if src.ProximityPlacementGroupReference != nil {
		isNotKubeRef := !src.ProximityPlacementGroupReference.IsKubernetesReference()
		// Note that using isNotKubeRef is a bit awkward because in reality it shouldn't be possible to have a genruntime.ResourceReference with no
		// kube ref AND no ARM ref, but if that does happen we pass the empty-string along to maintain the round-trip invariant.
		if len(src.ProximityPlacementGroupReference.ARMID) > 0 || isNotKubeRef {
			profile.ProximityPlacementGroupID = &src.ProximityPlacementGroupReference.ARMID
			propertyBag.Remove("ProximityPlacementGroupIDReference") // Remove it from property bag added by code generated code
		}
		// No need to handle the other case, as we would just put it into the property bag, which was already done by the generated code
	}

	// NodePublicIPPrefixReference
	if src.NodePublicIPPrefixReference != nil {
		nodePublicIPPrefixReference := src.NodePublicIPPrefixReference.Copy()
		profile.NodePublicIPPrefixIDReference = &nodePublicIPPrefixReference
	}
	propertyBag.Remove("NodePublicIPPrefixReference") // This should never be set

	// PodSubnetReference
	if src.PodSubnetReference != nil {
		podSubnetReference := src.PodSubnetReference.Copy()
		profile.PodSubnetIDReference = &podSubnetReference
	}
	propertyBag.Remove("PodSubnetReference") // This should never be set

	// VnetSubnetRefernece
	if src.VnetSubnetReference != nil {
		vnetSubnetReference := src.VnetSubnetReference.Copy()
		profile.VnetSubnetIDReference = &vnetSubnetReference
	}
	propertyBag.Remove("VnetSubnetReference") // This should never be set

	// Update the property bag
	if len(propertyBag) > 0 {
		profile.PropertyBag = propertyBag
	} else {
		profile.PropertyBag = nil
	}

	return nil
}

func (profile *ManagedClusterAgentPoolProfile) AssignPropertiesTo(dst *v20230201s.ManagedClusterAgentPoolProfile) error {
	// Clone the existing property bag
	dstPropertyBag := genruntime.NewPropertyBag(dst.PropertyBag)

	// ProximityPlacementGroupID
	if profile.ProximityPlacementGroupID != nil {
		dst.ProximityPlacementGroupReference = &genruntime.ResourceReference{
			ARMID: *profile.ProximityPlacementGroupID,
		}
	}
	// Ensure that this field is not set in the destination property bag (it shouldn't ever be there)
	dstPropertyBag.Remove("ProximityPlacementGroupID")

	// NodePublicIPPrefixID
	if profile.NodePublicIPPrefixIDReference != nil {
		nodePublicIPPrefixIDReference := profile.NodePublicIPPrefixIDReference.Copy()
		dst.NodePublicIPPrefixReference = &nodePublicIPPrefixIDReference
	}
	dstPropertyBag.Remove("NodePublicIPPrefixIDReference")

	// PodSubnetIDReference
	if profile.PodSubnetIDReference != nil {
		podSubnetIDReference := profile.PodSubnetIDReference.Copy()
		dst.PodSubnetReference = &podSubnetIDReference
	}
	dstPropertyBag.Remove("PodSubnetIDReference")

	// VnetSubnetIDReference
	if profile.VnetSubnetIDReference != nil {
		vnetSubnetIDReference := profile.VnetSubnetIDReference.Copy()
		dst.VnetSubnetReference = &vnetSubnetIDReference
	}
	dstPropertyBag.Remove("PodSubnetIDReference")

	// Update the property bag
	if len(dstPropertyBag) > 0 {
		dst.PropertyBag = dstPropertyBag
	} else {
		dst.PropertyBag = nil
	}

	return nil
}

// TODO: We can remove this interface implementation if we get config-based property rename handling
var _ augmentConversionForManagedCluster_Spec = &ManagedCluster_Spec{}

func (cluster *ManagedCluster_Spec) AssignPropertiesFrom(src *v20230201s.ManagedCluster_Spec) error {
	// Clone the existing property bag
	propertyBag := genruntime.NewPropertyBag(src.PropertyBag)

	// DiskEncryptionSetReference
	if src.DiskEncryptionSetReference != nil {
		diskEncryptionSetReference := src.DiskEncryptionSetReference.Copy()
		cluster.DiskEncryptionSetIDReference = &diskEncryptionSetReference
	}
	propertyBag.Remove("DiskEncryptionSetReference") // This should never be set

	// Update the property bag
	if len(propertyBag) > 0 {
		src.PropertyBag = propertyBag
	} else {
		src.PropertyBag = nil
	}

	return nil
}

func (cluster *ManagedCluster_Spec) AssignPropertiesTo(dst *v20230201s.ManagedCluster_Spec) error {
	// Clone the existing property bag
	dstPropertyBag := genruntime.NewPropertyBag(dst.PropertyBag)

	// DiskEncryptionSetIDReference
	if cluster.DiskEncryptionSetIDReference != nil {
		diskEncryptionSetIDReference := cluster.DiskEncryptionSetIDReference.Copy()
		dst.DiskEncryptionSetReference = &diskEncryptionSetIDReference
	}
	dstPropertyBag.Remove("DiskEncryptionSetIDReference") // This should never be set

	// Update the property bag
	if len(dstPropertyBag) > 0 {
		dst.PropertyBag = dstPropertyBag
	} else {
		dst.PropertyBag = nil
	}

	return nil
}

var _ augmentConversionForManagedClusterSKU = &ManagedClusterSKU{}

func (cluster *ManagedClusterSKU) AssignPropertiesFrom(_ *v20230201s.ManagedClusterSKU) error {
	// value will have already been set on cluster from code-generated conversion
	if to.Value(cluster.Name) == "Base" {
		cluster.Name = to.Ptr("Basic")
	}
	if to.Value(cluster.Tier) == "Standard" {
		cluster.Tier = to.Ptr("Paid")
	}

	return nil
}

func (*ManagedClusterSKU) AssignPropertiesTo(dst *v20230201s.ManagedClusterSKU) error {
	// value will have already been set on dst from code-generated conversion
	if to.Value(dst.Name) == "Basic" {
		dst.Name = to.Ptr("Base")
	}
	if to.Value(dst.Tier) == "Paid" {
		dst.Tier = to.Ptr("Standard")
	}

	return nil
}

var _ augmentConversionForManagedClusterSKU_STATUS = &ManagedClusterSKU_STATUS{}

func (cluster *ManagedClusterSKU_STATUS) AssignPropertiesFrom(_ *v20230201s.ManagedClusterSKU_STATUS) error {
	// value will have already been set on cluster from code-generated conversion
	if to.Value(cluster.Name) == "Base" {
		cluster.Name = to.Ptr("Basic")
	}
	if to.Value(cluster.Tier) == "Standard" {
		cluster.Tier = to.Ptr("Paid")
	}

	return nil
}

func (*ManagedClusterSKU_STATUS) AssignPropertiesTo(dst *v20230201s.ManagedClusterSKU_STATUS) error {
	// value will have already been set on dst from code-generated conversion
	if to.Value(dst.Name) == "Basic" {
		dst.Name = to.Ptr("Base")
	}
	if to.Value(dst.Tier) == "Paid" {
		dst.Tier = to.Ptr("Standard")
	}

	return nil
}
