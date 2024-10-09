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

import (
	"time"
)

// ProtoToTime converts a string from a DCL proto time string to a time.Time.
func ProtoToTime(s string) time.Time {
	// Invalid time values will be picked up downstream.
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

// TimeToProto converts a time.Time to a proto time string.
func TimeToProto(t time.Time) string {
	return t.Format(time.RFC3339)
}
