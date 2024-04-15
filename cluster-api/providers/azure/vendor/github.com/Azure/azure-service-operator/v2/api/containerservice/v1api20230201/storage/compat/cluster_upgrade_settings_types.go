/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package compat

import (
	"golang.org/x/exp/slices"

	v20231001s "github.com/Azure/azure-service-operator/v2/api/containerservice/v1api20231001/storage"
	"github.com/Azure/azure-service-operator/v2/internal/util/to"
)

var _ augmentConversionForUpgradeOverrideSettings = &UpgradeOverrideSettings{}

const ignoreKubernetesDeprecations = "IgnoreKubernetesDeprecations"

func (settings *UpgradeOverrideSettings) AssignPropertiesFrom(src *v20231001s.UpgradeOverrideSettings) error {
	// If the GA version has ForceUpgrade true, the preview version needs to add "IgnoreKubernetesDeprecations"
	if src.ForceUpgrade != nil && *src.ForceUpgrade {
		if !slices.Contains(settings.ControlPlaneOverrides, ignoreKubernetesDeprecations) {
			settings.ControlPlaneOverrides = append(settings.ControlPlaneOverrides, ignoreKubernetesDeprecations)
		}
	}

	settings.PropertyBag.Remove("ForceUpgrade")

	return nil
}

func (settings *UpgradeOverrideSettings) AssignPropertiesTo(dest *v20231001s.UpgradeOverrideSettings) error {
	// If the preview version has "IgnoreKubernetesDeprecations" in ControlPlaneOverrides, the GA version needs to have
	// ForceUpgrade true, otherwise false.
	dest.ForceUpgrade = to.Ptr(false)
	for _, override := range settings.ControlPlaneOverrides {
		if override == ignoreKubernetesDeprecations {
			dest.ForceUpgrade = to.Ptr(true)
			break
		}
	}

	return nil
}
