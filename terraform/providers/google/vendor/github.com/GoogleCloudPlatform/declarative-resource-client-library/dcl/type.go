// Copyright 2024 Google LLC. All Rights Reserved.
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
