// Copyright 2021 Google LLC. All Rights Reserved.
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
//     http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package dcl

// Bool converts a bool to a *bool
func Bool(b bool) *bool {
	return &b
}

// Float64 converts a float64 to *float64
func Float64(f float64) *float64 {
	return &f
}

// Float64OrNil converts a float64 to *float64, returning nil if it's empty (0.0).
func Float64OrNil(f float64) *float64 {
	if f == 0.0 {
		return nil
	}
	return &f
}

// Int64 converts an int64 to *int64
func Int64(i int64) *int64 {
	return &i
}

// Int64OrNil converts an int64 to *int64, returning nil if it's empty (0).
func Int64OrNil(i int64) *int64 {
	if i == 0 {
		return nil
	}
	return &i
}

// String converts a string to a *string
func String(s string) *string {
	return &s
}

// StringWithError converts a string to a *string, returning a nil error to
// satisfy type signatures that expect one.
func StringWithError(s string) (*string, error) {
	return &s, nil
}

// StringOrNil converts a string to a *string, returning nil if it's empty ("").
func StringOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// Optional refers to all primitives that handle the proto3 optional keyword.
type Optional interface {
	IsUnset() bool
	GetValue() interface{}
}

// OptionalBool refers to a boolean value that is optional.
type OptionalBool struct {
	Unset bool
	Value *bool
}

// IsUnset returns if it is unset.
func (b *OptionalBool) IsUnset() bool {
	return b.Unset
}

// GetValue returns the value of this type.
func (b *OptionalBool) GetValue() interface{} {
	return b.Value
}

// OptionalInt64 refers to an int64 value that is optional.
type OptionalInt64 struct {
	Unset bool
	Value *int64
}

// IsUnset returns if it is unset.
func (b *OptionalInt64) IsUnset() bool {
	return b.Unset
}

// GetValue returns the value of this type.
func (b *OptionalInt64) GetValue() interface{} {
	return b.Value
}

// OptionalFloat64 refers to a float64 value that is optional.
type OptionalFloat64 struct {
	Unset bool
	Value *float64
}

// IsUnset returns if it is unset.
func (b *OptionalFloat64) IsUnset() bool {
	return b.Unset
}

// GetValue returns the value of this type.
func (b *OptionalFloat64) GetValue() interface{} {
	return b.Value
}

// OptionalString refers to a string value that is optional.
type OptionalString struct {
	Unset bool
	Value *string
}

// IsUnset returns if it is unset.
func (b *OptionalString) IsUnset() bool {
	return b.Unset
}

// GetValue returns the value of this type.
func (b *OptionalString) GetValue() interface{} {
	return b.Value
}

// SetOptionalString returns a optional string that is not unset.
func SetOptionalString(value string) *OptionalString {
	return &OptionalString{
		Unset: false,
		Value: String(value),
	}
}

// UnsetOptionalString returns a optional string that is unset.
func UnsetOptionalString() *OptionalString {
	return &OptionalString{
		Unset: true,
	}
}

// SetOptionalBool returns a optional bool that is not unset.
func SetOptionalBool(value bool) *OptionalBool {
	return &OptionalBool{
		Unset: false,
		Value: Bool(value),
	}
}

// UnsetOptionalBool returns a optional bool that is unset.
func UnsetOptionalBool() *OptionalBool {
	return &OptionalBool{
		Unset: true,
	}
}

// SetOptionalInt64 returns a optional int64 that is not unset.
func SetOptionalInt64(value int64) *OptionalInt64 {
	return &OptionalInt64{
		Unset: false,
		Value: Int64(value),
	}
}

// UnsetOptionalInt64 returns a optional int64 that is unset.
func UnsetOptionalInt64() *OptionalInt64 {
	return &OptionalInt64{
		Unset: true,
	}
}

// SetOptionalFloat64 returns a optional float64 that is not unset.
func SetOptionalFloat64(value float64) *OptionalFloat64 {
	return &OptionalFloat64{
		Unset: false,
		Value: Float64(value),
	}
}

// UnsetOptionalFloat64 returns a optional float64 that is unset.
func UnsetOptionalFloat64() *OptionalFloat64 {
	return &OptionalFloat64{
		Unset: true,
	}
}
