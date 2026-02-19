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
// Package bigqueryreservation defines types and methods for working with bigqueryreservation GCP resources.
package bigqueryreservation

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

// This resource has a custom matcher so that either name or assignee + job_type can be used depending on which is available.
func (r *Assignment) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalAssignment(b, c, r)
		if err != nil {
			c.Config.Logger.Warning("failed to unmarshal provided resource in matcher.")
			return false
		}
		nr := r.urlNormalized()
		ncr := cr.urlNormalized()
		c.Config.Logger.Infof("looking for %v\nin %v", nr, ncr)

		if dcl.IsEmptyValueIndirect(nr.Name) || dcl.IsEmptyValueIndirect(ncr.Name) {
			// Name field not available for matching, use job type and assignee.
			if nr.Assignee == nil && ncr.Assignee == nil {
				c.Config.Logger.Info("Both Assignee fields null - considering equal.")
			} else if nr.Assignee == nil || ncr.Assignee == nil {
				c.Config.Logger.Info("Only one Assignee field is null - considering unequal.")
				return false
			} else if *nr.Assignee != *ncr.Assignee {
				return false
			}
			if nr.JobType == nil && ncr.JobType == nil {
				c.Config.Logger.Info("Both JobType fields null - considering equal.")
			} else if nr.JobType == nil || ncr.JobType == nil {
				c.Config.Logger.Info("Only one JobType field is null - considering unequal.")
				return false
			} else if *nr.JobType != *ncr.JobType {
				return false
			}
		} else {
			// Name field is available.
			if nr.Name == nil && ncr.Name == nil {
				c.Config.Logger.Info("Both Name fields null - considering equal.")
			} else if nr.Name == nil || ncr.Name == nil {
				c.Config.Logger.Info("Only one Name field is null - considering unequal.")
				return false
			} else if *nr.Name != *ncr.Name {
				return false
			}
		}
		return true
	}
}
