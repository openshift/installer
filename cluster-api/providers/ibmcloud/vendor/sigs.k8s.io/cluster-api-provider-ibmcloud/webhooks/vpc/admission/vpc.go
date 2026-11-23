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

package vpc

import (
	"fmt"
	"regexp"

	"k8s.io/apimachinery/pkg/util/validation/field"

	infrav1 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/vpc/v1beta2"
)

// IBM Cloud CRN validation regex.
var crnRegex = regexp.MustCompile(`^crn:v[0-9]+:[a-z0-9-]+:[a-z0-9-]+:[a-z0-9-]+:[a-z0-9-]*:([a-z]\/[a-z0-9-]+)?:[a-z0-9-]*:[a-z0-9-]*:[a-zA-Z0-9-_\.\/]*$`)

const (
	// customProfile is the first-generation volume profile with user-defined iops.
	customProfile = "custom"
	// sdpProfile is the second-generation volume profile with independently adjustable iops and bandwidth.
	sdpProfile = "sdp"
)

func defaultIBMVPCMachineSpec(spec *infrav1.IBMVPCMachineSpec) {
	if spec.Profile == "" {
		spec.Profile = "bx2-2x8"
	}
}

func validateVolumes(spec infrav1.IBMVPCMachineSpec) field.ErrorList {
	var allErrs field.ErrorList

	allErrs = append(allErrs, validateAdditionalVolumes(spec)...)

	if spec.BootVolume != nil {
		allErrs = append(allErrs, validateBootVolume(spec)...)
	}

	return allErrs
}

// validateAdditionalVolumes validates the additional volumes configuration.
func validateAdditionalVolumes(spec infrav1.IBMVPCMachineSpec) field.ErrorList {
	var allErrs field.ErrorList

	for i := range spec.AdditionalVolumes {
		// A check is required for SizeGiB here but not in BootVolumes because BootVolumes have a default size of 100GiB that is allocated when the size is missing.
		// The same is not true for AdditionalVolumes, therefore it is a mandatory field here.
		if spec.AdditionalVolumes[i].SizeGiB == 0 {
			allErrs = append(allErrs, field.Invalid(field.NewPath(fmt.Sprintf("spec.AdditionalVolumes[%d]", i)), spec, "sizeGiB has to be specified"))
		}
		if spec.AdditionalVolumes[i].Profile == customProfile && spec.AdditionalVolumes[i].Iops == 0 {
			allErrs = append(allErrs, field.Invalid(field.NewPath(fmt.Sprintf("spec.AdditionalVolumes[%d]", i)), spec, "iops has to be specified when profile is set to `custom` "))
		}
		if spec.AdditionalVolumes[i].Iops != 0 && !volumeProfileSupportsIops(spec.AdditionalVolumes[i].Profile) {
			allErrs = append(allErrs, field.Invalid(field.NewPath(fmt.Sprintf("spec.AdditionalVolumes[%d]", i)), spec, "iops applicable only to volumes using a profile of type `custom` or `sdp`"))
		}
		if spec.AdditionalVolumes[i].Bandwidth != 0 && spec.AdditionalVolumes[i].Profile != sdpProfile {
			allErrs = append(allErrs, field.Invalid(field.NewPath(fmt.Sprintf("spec.AdditionalVolumes[%d]", i)), spec, "bandwidth applicable only to volumes using a profile of type `sdp`"))
		}
		if spec.AdditionalVolumes[i].EncryptionKeyCRN != "" && !isValidCRN(spec.AdditionalVolumes[i].EncryptionKeyCRN) {
			allErrs = append(allErrs, field.Invalid(field.NewPath(fmt.Sprintf("spec.AdditionalVolumes[%d]", i)), spec, "encryptionKeyCRN not in proper IBM Cloud CRN format"))
		}
	}

	return allErrs
}

// validateBootVolume validates the boot volume configuration.
func validateBootVolume(spec infrav1.IBMVPCMachineSpec) field.ErrorList {
	var allErrs field.ErrorList

	// Second-generation (sdp) boot volumes support larger capacities than first-generation profiles.
	maxBootVolumeSizeGiB := int64(250)
	if spec.BootVolume.Profile == sdpProfile {
		maxBootVolumeSizeGiB = 32000
	}
	if spec.BootVolume.SizeGiB < 10 || spec.BootVolume.SizeGiB > maxBootVolumeSizeGiB {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec.bootVolume.sizeGiB"), spec, fmt.Sprintf("valid Boot VPCVolume size is 10 - %d GB", maxBootVolumeSizeGiB)))
	}

	if spec.BootVolume.Iops != 0 && !volumeProfileSupportsIops(spec.BootVolume.Profile) {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec.bootVolume.iops"), spec, "iops applicable only to volumes using a profile of type `custom` or `sdp`"))
	}

	if spec.BootVolume.Bandwidth != 0 && spec.BootVolume.Profile != sdpProfile {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec.bootVolume.bandwidth"), spec, "bandwidth applicable only to volumes using a profile of type `sdp`"))
	}

	//  Validate spec.BootVolume.EncryptionKeyCRN to ensure its in proper IBM Cloud CRN format
	if spec.BootVolume.EncryptionKeyCRN != "" && !isValidCRN(spec.BootVolume.EncryptionKeyCRN) {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec.bootVolume.encryptionKeyCRN"), spec, "encryptionKeyCRN not in proper IBM Cloud CRN format"))
	}

	return allErrs
}

// volumeProfileSupportsIops reports whether the volume profile allows the iops to be specified by the user.
func volumeProfileSupportsIops(profile string) bool {
	return profile == customProfile || profile == sdpProfile
}

// isValidCRN checks whether the provided string is a valid IBM Cloud CRN.
func isValidCRN(crn string) bool {
	return crnRegex.MatchString(crn)
}
