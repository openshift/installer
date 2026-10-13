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

package openstack

import (
	version "github.com/hashicorp/go-version"
	klog "k8s.io/klog/v2"
)

const (
	OctaviaFeatureTags              = 0
	OctaviaFeatureVIPACL            = 1
	OctaviaFeatureFlavors           = 2
	OctaviaFeatureTimeout           = 3
	OctaviaFeatureAvailabilityZones = 4
	lbProviderOVN                   = "ovn"
)

// IsOctaviaFeatureSupported returns true if the given feature is supported in the deployed Octavia version.
// copied from https://github.com/kubernetes/cloud-provider-openstack/blob/master/pkg/util/openstack/loadbalancer.go#L95-L148
func IsOctaviaFeatureSupported(octaviaVer string, feature int, lbProvider string) bool {
	currentVer, _ := version.NewVersion(octaviaVer)

	switch feature {
	case OctaviaFeatureVIPACL:
		if lbProvider == lbProviderOVN {
			return false
		}
		verACL, _ := version.NewVersion("v2.12")
		if currentVer.GreaterThanOrEqual(verACL) {
			return true
		}
	case OctaviaFeatureTags:
		verTags, _ := version.NewVersion("v2.5")
		if currentVer.GreaterThanOrEqual(verTags) {
			return true
		}
	case OctaviaFeatureFlavors:
		if lbProvider == lbProviderOVN {
			return false
		}
		verFlavors, _ := version.NewVersion("v2.6")
		if currentVer.GreaterThanOrEqual(verFlavors) {
			return true
		}
	case OctaviaFeatureTimeout:
		if lbProvider == lbProviderOVN {
			return false
		}
		verFlavors, _ := version.NewVersion("v2.1")
		if currentVer.GreaterThanOrEqual(verFlavors) {
			return true
		}
	case OctaviaFeatureAvailabilityZones:
		if lbProvider == lbProviderOVN {
			return false
		}
		verAvailabilityZones, _ := version.NewVersion("v2.14")
		if currentVer.GreaterThanOrEqual(verAvailabilityZones) {
			return true
		}
	default:
		klog.Warningf("Feature %d not recognized", feature)
	}

	return false
}
