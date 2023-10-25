/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package genruntime

import (
	"strings"
)

// ConvertToARMResolvedDetails contains resolved references and names for use in
// converting a Kubernetes type to an ARM type.
type ConvertToARMResolvedDetails struct {
	// Name is the name of the resource
	// TODO: We might be able to remove this in favor of using AzureName() everywhere in the future
	Name string

	// ResolvedReferences is a set of references which have been resolved to their ARM IDs.
	ResolvedReferences Resolved[ResourceReference]

	// ResolvedSecrets is a set of secret references which have been resolved to the corresponding
	// secret value.
	ResolvedSecrets Resolved[SecretReference]

	// ResolvedConfigMaps is a set of config map references which have been resolved to the corresponding
	// config map value.
	ResolvedConfigMaps Resolved[ConfigMapReference]
}

type ToARMConverter interface {
	// ConvertToARM converts this to an ARM resource.
	ConvertToARM(resolved ConvertToARMResolvedDetails) (interface{}, error)
}

type FromARMConverter interface {
	NewEmptyARMValue() ARMResourceStatus
	PopulateFromARM(owner ArbitraryOwnerReference, input interface{}) error
}

// TODO: Consider ArmSpecTransformer and ARMTransformer, so we don't have to pass owningName/name through all the calls
// ARMTransformer is a type which can be converted to/from an Arm object shape.
// Each CRD resource must implement these methods.
type ARMTransformer interface {
	ToARMConverter
	FromARMConverter
}

// ExtractKubernetesResourceNameFromARMName extracts the Kubernetes resource name from an ARM name.
// See https://docs.microsoft.com/en-us/azure/azure-resource-manager/templates/child-resource-name-type#outside-parent-resource
// for details on the format of the name field in ARM templates.
func ExtractKubernetesResourceNameFromARMName(armName string) string {
	if len(armName) == 0 {
		return ""
	}

	// TODO: Possibly need to worry about preserving case here, although ARM should be already
	strs := strings.Split(armName, "/")
	return strs[len(strs)-1]
}
