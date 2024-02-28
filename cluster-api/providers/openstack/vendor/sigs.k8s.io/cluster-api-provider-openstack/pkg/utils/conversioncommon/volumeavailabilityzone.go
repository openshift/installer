/*
Copyright 2024 The Kubernetes Authors.

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

package conversioncommon

import (
	"k8s.io/apimachinery/pkg/conversion"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
)

func Convert_string_To_Pointer_v1beta1_VolumeAvailabilityZone(in *string, out **infrav1.VolumeAvailabilityZone, _ conversion.Scope) error {
	switch *in {
	case "":
		*out = &infrav1.VolumeAvailabilityZone{
			From: infrav1.VolumeAZFromMachine,
		}
	default:
		azName := infrav1.VolumeAZName(*in)
		*out = &infrav1.VolumeAvailabilityZone{
			From: infrav1.VolumeAZFromName,
			Name: &azName,
		}
	}

	return nil
}

func Convert_Pointer_v1beta1_VolumeAvailabilityZone_To_string(in **infrav1.VolumeAvailabilityZone, out *string, _ conversion.Scope) error {
	// This is a lossy: can't specify no AZ prior to v1beta1
	if *in == nil {
		*out = ""
		return nil
	}

	switch (*in).From {
	case "", infrav1.VolumeAZFromName:
		name := (*in).Name
		if name != nil {
			*out = string(*name)
		} else {
			*out = ""
		}
	case infrav1.VolumeAZFromMachine:
		*out = ""
	}

	return nil
}
