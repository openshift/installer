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

package optional

import (
	"k8s.io/apimachinery/pkg/conversion"
)

func RestoreString(previous, dst *String) {
	if *dst == nil || **dst == "" {
		*dst = *previous
	}
}

func RestoreInt(previous, dst *Int) {
	if *dst == nil || **dst == 0 {
		*dst = *previous
	}
}

func RestoreUInt16(previous, dst *UInt16) {
	if *dst == nil || **dst == 0 {
		*dst = *previous
	}
}

func RestoreBool(previous, dst *Bool) {
	if *dst == nil || !**dst {
		*dst = *previous
	}
}

func Convert_string_To_optional_String(in *string, out *String, _ conversion.Scope) error {
	// NOTE: This function has the opposite defaulting behaviour to
	// Convert_string_to_Pointer_string defined in apimachinery: it converts
	// empty strings to nil instead of a pointer to an empty string.

	if *in == "" {
		*out = nil
	} else {
		*out = in
	}
	return nil
}

func Convert_optional_String_To_string(in *String, out *string, _ conversion.Scope) error {
	if *in == nil {
		*out = ""
	} else {
		*out = **in
	}
	return nil
}

func Convert_int_To_optional_Int(in *int, out *Int, _ conversion.Scope) error {
	if *in == 0 {
		*out = nil
	} else {
		*out = in
	}
	return nil
}

func Convert_optional_Int_To_int(in *Int, out *int, _ conversion.Scope) error {
	if *in == nil {
		*out = 0
	} else {
		*out = **in
	}
	return nil
}

func Convert_uint16_To_optional_UInt16(in *uint16, out *UInt16, _ conversion.Scope) error {
	if *in == 0 {
		*out = nil
	} else {
		*out = in
	}
	return nil
}

func Convert_optional_UInt16_To_uint16(in *UInt16, out *uint16, _ conversion.Scope) error {
	if *in == nil {
		*out = 0
	} else {
		*out = **in
	}
	return nil
}

func Convert_bool_To_optional_Bool(in *bool, out *Bool, _ conversion.Scope) error {
	if !*in {
		*out = nil
	} else {
		*out = in
	}
	return nil
}

func Convert_optional_Bool_To_bool(in *Bool, out *bool, _ conversion.Scope) error {
	if *in == nil {
		*out = false
	} else {
		*out = **in
	}
	return nil
}
