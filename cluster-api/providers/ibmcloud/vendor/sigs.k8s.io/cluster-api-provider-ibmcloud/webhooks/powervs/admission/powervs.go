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

package powervs

import (
	"strconv"

	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/validation/field"

	infrav1 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/powervs/v1beta3"
)

const (
	defaultSystemType = "s1022"
)

func defaultIBMPowerVSMachineSpec(spec *infrav1.IBMPowerVSMachineSpec) {
	if spec.MemoryGiB == 0 {
		spec.MemoryGiB = 2
	}
	if spec.SystemType == "" {
		spec.SystemType = defaultSystemType
	}
	if spec.ProcessorType == "" {
		spec.ProcessorType = infrav1.PowerVSProcessorTypeShared
	}
	// Default Processors based on ProcessorType: Dedicated requires a whole number (min 1),
	// Shared/Capped accept fractional values (min 0.25).
	if spec.Processors.StrVal == "" && spec.Processors.IntVal == 0 {
		if spec.ProcessorType == infrav1.PowerVSProcessorTypeDedicated {
			spec.Processors = intstr.FromInt32(1)
		} else {
			spec.Processors = intstr.FromString("0.25")
		}
	}
}

// validateIBMPowerVSProcessorValues validates the processors value against the
// rules that apply to the given ProcessorType:
//   - Dedicated: must be a whole number >= 1
//   - Shared / Capped: must be a number >= 0.25
func validateIBMPowerVSProcessorValues(procType infrav1.PowerVSProcessorType, resValue intstr.IntOrString) *field.Error {
	fldPath := field.NewPath("spec", "processors")

	var val float64
	switch resValue.Type {
	case intstr.Int:
		val = float64(resValue.IntVal)
	case intstr.String:
		var err error
		if val, err = strconv.ParseFloat(resValue.StrVal, 64); err != nil {
			return field.Invalid(fldPath, resValue, "processors must be a valid number")
		}
	}

	if procType == infrav1.PowerVSProcessorTypeDedicated {
		if val < 1 {
			return field.Invalid(fldPath, resValue, "processors must be at least 1 when processorType is Dedicated")
		}
		if val != float64(int64(val)) {
			return field.Invalid(fldPath, resValue, "processors cannot be fractional when processorType is Dedicated")
		}
	} else if val < 0.25 {
		return field.Invalid(fldPath, resValue, "processors must be at least 0.25")
	}

	return nil
}
