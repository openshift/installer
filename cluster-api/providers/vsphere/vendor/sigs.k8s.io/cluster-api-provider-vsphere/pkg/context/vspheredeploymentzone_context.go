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

package context

import (
	"context"
	"fmt"

	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/session"
)

// VSphereDeploymentZoneContext contains information for the VSphereDeploymentZone reconciliation.
type VSphereDeploymentZoneContext struct {
	*ControllerManagerContext
	VSphereDeploymentZone *infrav1.VSphereDeploymentZone
	PatchHelper           *patch.Helper
	AuthSession           *session.Session
}

// Patch patches the VSphereDeploymentZone.
func (c *VSphereDeploymentZoneContext) Patch(ctx context.Context) error {
	conditions.SetSummary(c.VSphereDeploymentZone,
		conditions.WithConditions(
			infrav1.VCenterAvailableCondition,
			infrav1.VSphereFailureDomainValidatedCondition,
			infrav1.PlacementConstraintMetCondition,
		),
	)
	return c.PatchHelper.Patch(ctx, c.VSphereDeploymentZone)
}

// String returns a string with the GroupVersionKind and name of the VSphereDeploymentZone.
func (c *VSphereDeploymentZoneContext) String() string {
	return fmt.Sprintf("%s %s", c.VSphereDeploymentZone.GroupVersionKind(), c.VSphereDeploymentZone.Name)
}

// GetSession returns the session for the VSphereDeploymentZoneContext.
func (c *VSphereDeploymentZoneContext) GetSession() *session.Session {
	return c.AuthSession
}
