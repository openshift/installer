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

// String is a string that can be unspecified. strings which are converted to
// optional.String during API conversion will be converted to nil if the value
// was previously the empty string.
// +optional.
type String *string

// Int is an int that can be unspecified. ints which are converted to
// optional.Int during API conversion will be converted to nil if the value
// was previously 0.
// +optional.
type Int *int

// UInt16 is a uint16 that can be unspecified. uint16s which are converted to
// optional.UInt16 during API conversion will be converted to nil if the value
// was previously 0.
// +kubebuilder:validation:Minimum:=0
// +kubebuilder:validation:Maximum:=65535
// +optional.
type UInt16 *uint16

// Bool is a bool that can be unspecified. bools which are converted to
// optional.Bool during API conversion will be converted to nil if the value
// was previously false.
type Bool *bool
