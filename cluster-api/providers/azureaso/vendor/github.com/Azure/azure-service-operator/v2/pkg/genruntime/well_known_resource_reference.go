/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package genruntime

// WellKnownResourceReference is a variation of ResourceReference that also permits the referenced resource to be
// identified by a well known name.
type WellKnownResourceReference struct {
	ResourceReference `json:",inline"`

	WellKnownName string `json:"wellKnownName,omitempty"`
}

// Copy makes an independent copy of the WellKnownResourceReference
func (ref *WellKnownResourceReference) Copy() WellKnownResourceReference {
	return *ref
}

// CreateWellKnownResourceReferenceFromARMID creates a new WellKnownResourceReference from a string representing an ARM ID
func CreateWellKnownResourceReferenceFromARMID(armID string) WellKnownResourceReference {
	return WellKnownResourceReference{
		ResourceReference: ResourceReference{
			ARMID: armID,
		},
	}
}
