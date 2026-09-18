/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package crdmanagement

import (
	"regexp"
	"slices"

	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

var versionRegexp = regexp.MustCompile(`^v(1api|1alpha1api|1beta)?(\d{8}|\d)(preview)?(storage)?$`)

// GetDeprecatedStorageVersions returns a new list of storedVersions by removing the deprecated versions
// The first return value is the list of versions that are NOT deprecated.
// The second return value is the list of versions that ARE deprecated.
func GetDeprecatedStorageVersions(crd apiextensions.CustomResourceDefinition) ([]string, []string) {
	storedVersions := crd.Status.StoredVersions

	if len(storedVersions) <= 1 {
		// No versions to remove
		return storedVersions, nil
	}

	maxVersion := slices.MaxFunc(storedVersions, CompareVersions)
	newStoredVersions := make([]string, 0, len(storedVersions))
	var removedVersions []string
	for _, version := range storedVersions {
		if version == maxVersion {
			newStoredVersions = append(newStoredVersions, version)
		} else {
			removedVersions = append(removedVersions, version)
		}
	}

	return newStoredVersions, removedVersions
}

// GetDeprecatedStorageVersions returns a map of CRD names to their deprecated storage versions.
func GetAllDeprecatedStorageVersions(crds map[string]apiextensions.CustomResourceDefinition) map[string][]string {
	result := make(map[string][]string)
	for _, crd := range crds {
		_, deprecatedVersions := GetDeprecatedStorageVersions(crd)
		if len(deprecatedVersions) > 0 {
			result[crd.Name] = deprecatedVersions
		}
	}
	return result
}

// CompareVersions compares two group version strings, for example: v1api20230101preview, v1api20220101, v20220101, v1api20220101storage, and v20220101storage
// and returns -1 if a < b, 0 if a == b, and 1 if a > b. The storage suffix is ignored in the comparison.
// Preview versions are considered less than non-preview versions.
func CompareVersions(a string, b string) int {
	matchA := versionRegexp.FindStringSubmatch(a)
	matchB := versionRegexp.FindStringSubmatch(b)

	if matchA == nil || matchB == nil {
		// If either version doesn't match the expected pattern, fall back to string comparison
		// TODO: Do we actually want to do this?
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	}

	// Prefer non-preview over preview
	isAPreview := matchA[3] == "preview"
	isBPreview := matchB[3] == "preview"

	if isAPreview && !isBPreview {
		return -1
	} else if !isAPreview && isBPreview {
		return 1
	}

	// Prefer everything over v1alpha1apiYYYYMMDD
	hasV1Alpha1A := matchA[1] == "1alpha1api"
	hasV1Alpha1B := matchB[1] == "1alpha1api"

	if hasV1Alpha1A && !hasV1Alpha1B {
		return -1
	} else if !hasV1Alpha1A && hasV1Alpha1B {
		return 1
	}

	// Prefer everything over v1betaYYYYMMDD
	hasV1BetaA := matchA[1] == "1beta"
	hasV1BetaB := matchB[1] == "1beta"

	if hasV1BetaA && !hasV1BetaB {
		return -1
	} else if !hasV1BetaA && hasV1BetaB {
		return 1
	}

	// Prefer newer date over older date
	dateA := matchA[2]
	dateB := matchB[2]

	if dateA < dateB {
		return -1
	} else if dateA > dateB {
		return 1
	}

	// Prefer vYYYYMMDD over v1apiYYYYMMDD
	hasV1ApiA := matchA[1] == "1api"
	hasV1ApiB := matchB[1] == "1api"

	if hasV1ApiA && !hasV1ApiB {
		return -1
	} else if !hasV1ApiA && hasV1ApiB {
		return 1
	}

	return 0
}
