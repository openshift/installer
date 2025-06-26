/*
Copyright 2021 The Kubernetes Authors.

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

// Package addons provides a plan to manage EKS addons.
package addons

import (
	"context"
	"time"

	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/eks"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/planner"
)

// NewPlan creates a new Plan to manage EKS addons.
func NewPlan(clusterName string, desiredAddons, installedAddons []*EKSAddon, client eks.Client, maxWait time.Duration) planner.Plan {
	return &plan{
		installedAddons:           installedAddons,
		desiredAddons:             desiredAddons,
		eksClient:                 client,
		clusterName:               clusterName,
		maxWaitActiveUpdateDelete: maxWait,
	}
}

// Plan is a plan that will manage EKS addons.
type plan struct {
	installedAddons           []*EKSAddon
	desiredAddons             []*EKSAddon
	eksClient                 eks.Client
	clusterName               string
	maxWaitActiveUpdateDelete time.Duration
}

// Create will create the plan (i.e. list of procedures) for managing EKS addons.
func (a *plan) Create(_ context.Context) ([]planner.Procedure, error) {
	procedures := []planner.Procedure{}

	// Handle create and update
	for i := range a.desiredAddons {
		desired := a.desiredAddons[i]
		installed := a.getInstalled(*desired.Name)
		if installed == nil {
			// Need to add the addon
			procedures = append(procedures,
				&CreateAddonProcedure{plan: a, name: *desired.Name},
				&WaitAddonActiveProcedure{plan: a, name: *desired.Name, includeDegraded: true},
			)
		} else {
			// Check if its just the tags that need updating
			diffTags := desired.Tags.Difference(installed.Tags)
			if len(diffTags) > 0 {
				procedures = append(procedures, &UpdateAddonTagsProcedure{plan: a, name: *installed.Name})
			}
			// Check if we also need to update the addon
			if !desired.IsEqual(installed, false) {
				procedures = append(procedures,
					&UpdateAddonProcedure{plan: a, name: *installed.Name},
					&WaitAddonActiveProcedure{plan: a, name: *desired.Name, includeDegraded: true},
				)
			} else if *installed.Status != string(ekstypes.AddonStatusActive) {
				// If the desired and installed are the same make sure its active
				procedures = append(procedures, &WaitAddonActiveProcedure{plan: a, name: *desired.Name, includeDegraded: true})
			}
		}
	}

	// look for deletions & unchanged
	for i := range a.installedAddons {
		installed := a.installedAddons[i]
		desired := a.getDesired(*installed.Name)
		if desired == nil {
			if *installed.Status != string(ekstypes.AddonStatusDeleting) {
				procedures = append(procedures, &DeleteAddonProcedure{plan: a, name: *installed.Name})
			}
			procedures = append(procedures, &WaitAddonDeleteProcedure{plan: a, name: *installed.Name})
		}
	}

	return procedures, nil
}

func (a *plan) getInstalled(name string) *EKSAddon {
	for i := range a.installedAddons {
		installed := a.installedAddons[i]
		if *installed.Name == name {
			return installed
		}
	}

	return nil
}

func (a *plan) getDesired(name string) *EKSAddon {
	for i := range a.desiredAddons {
		desired := a.desiredAddons[i]
		if *desired.Name == name {
			return desired
		}
	}

	return nil
}
