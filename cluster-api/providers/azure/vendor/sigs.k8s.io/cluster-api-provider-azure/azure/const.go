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

package azure

const (
	// VMTagsLastAppliedAnnotation is the key for the machine object annotation
	// which tracks the AdditionalTags in the Machine Provider Config.
	// See https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/
	// for annotation formatting rules.
	VMTagsLastAppliedAnnotation = "sigs.k8s.io/cluster-api-provider-azure-last-applied-tags-vm"

	// RGTagsLastAppliedAnnotation is the key for the Azure Cluster object annotation
	// which tracks the AdditionalTags for Resource Group which is part in the Azure Cluster.
	// See https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/
	// for annotation formatting rules.
	RGTagsLastAppliedAnnotation = "sigs.k8s.io/cluster-api-provider-azure-last-applied-tags-rg"

	// ManagedClusterTagsLastAppliedAnnotation is the key for the AzureManagedControlPlane
	// object annotation which tracks the AdditionalTags for managed clusters.
	// See https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/
	// for annotation formatting rules.
	ManagedClusterTagsLastAppliedAnnotation = "sigs.k8s.io/cluster-api-provider-azure-last-applied-tags-managedcluster"

	// SecurityRuleLastAppliedAnnotation is the key for the Azure Cluster
	// object annotation which tracks the security rules for security groups.
	// See https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/
	// for annotation formatting rules.
	SecurityRuleLastAppliedAnnotation = "sigs.k8s.io/cluster-api-provider-azure-last-applied-security-rules"

	// CustomDataHashAnnotation is the key for the machine object annotation
	// which tracks the hash of the custom data.
	// See https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/
	// for annotation formatting rules.
	CustomDataHashAnnotation = "sigs.k8s.io/cluster-api-provider-azure-vmss-custom-data-hash"
)
