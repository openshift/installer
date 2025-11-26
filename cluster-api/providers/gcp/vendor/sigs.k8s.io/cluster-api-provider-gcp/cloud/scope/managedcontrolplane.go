/*
Copyright 2023 The Kubernetes Authors.

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

package scope

import (
	"context"
	"fmt"

	"sigs.k8s.io/cluster-api-provider-gcp/util/location"

	"sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"

	container "cloud.google.com/go/container/apiv1"
	credentials "cloud.google.com/go/iam/credentials/apiv1"
	resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
	"github.com/pkg/errors"
	infrav1exp "sigs.k8s.io/cluster-api-provider-gcp/exp/api/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	patch "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// APIServerPort is the port of the GKE api server.
	APIServerPort = 443
)

// ManagedControlPlaneScopeParams defines the input parameters used to create a new Scope.
type ManagedControlPlaneScopeParams struct {
	CredentialsClient      *credentials.IamCredentialsClient
	ManagedClusterClient   *container.ClusterManagerClient
	TagBindingsClient      *resourcemanager.TagBindingsClient
	Client                 client.Client
	Cluster                *clusterv1.Cluster
	GCPManagedCluster      *infrav1exp.GCPManagedCluster
	GCPManagedControlPlane *infrav1exp.GCPManagedControlPlane
}

// NewManagedControlPlaneScope creates a new Scope from the supplied parameters.
// This is meant to be called for each reconcile iteration.
func NewManagedControlPlaneScope(ctx context.Context, params ManagedControlPlaneScopeParams) (*ManagedControlPlaneScope, error) {
	if params.Cluster == nil {
		return nil, errors.New("failed to generate new scope from nil Cluster")
	}
	if params.GCPManagedCluster == nil {
		return nil, errors.New("failed to generate new scope from nil GCPManagedCluster")
	}
	if params.GCPManagedControlPlane == nil {
		return nil, errors.New("failed to generate new scope from nil GCPManagedControlPlane")
	}

	credential, err := getCredentials(ctx, params.GCPManagedCluster.Spec.CredentialsRef, params.Client)
	if err != nil {
		return nil, fmt.Errorf("getting gcp credentials: %w", err)
	}

	if params.ManagedClusterClient == nil {
		managedClusterClient, err := newClusterManagerClient(ctx, params.GCPManagedCluster.Spec.CredentialsRef, params.Client, params.GCPManagedCluster.Spec.ServiceEndpoints)
		if err != nil {
			return nil, errors.Errorf("failed to create gcp managed cluster client: %v", err)
		}
		params.ManagedClusterClient = managedClusterClient
	}
	if params.TagBindingsClient == nil {
		tagBindingsClient, err := newTagBindingsClient(ctx, params.GCPManagedCluster.Spec.CredentialsRef, params.Client, params.GCPManagedCluster.Spec.Region, params.GCPManagedCluster.Spec.ServiceEndpoints)
		if err != nil {
			return nil, errors.Errorf("failed to create gcp tag bindings client: %v", err)
		}
		params.TagBindingsClient = tagBindingsClient
	}
	if params.CredentialsClient == nil {
		var credentialsClient *credentials.IamCredentialsClient
		credentialsClient, err = newIamCredentialsClient(ctx, params.GCPManagedCluster.Spec.CredentialsRef, params.Client, params.GCPManagedCluster.Spec.ServiceEndpoints)
		if err != nil {
			return nil, errors.Errorf("failed to create gcp credentials client: %v", err)
		}
		params.CredentialsClient = credentialsClient
	}

	helper, err := patch.NewHelper(params.GCPManagedControlPlane, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}

	return &ManagedControlPlaneScope{
		client:                 params.Client,
		Cluster:                params.Cluster,
		GCPManagedCluster:      params.GCPManagedCluster,
		GCPManagedControlPlane: params.GCPManagedControlPlane,
		mcClient:               params.ManagedClusterClient,
		tagBindingsClient:      params.TagBindingsClient,
		credentialsClient:      params.CredentialsClient,
		credential:             credential,
		patchHelper:            helper,
	}, nil
}

// ManagedControlPlaneScope defines the basic context for an actuator to operate upon.
type ManagedControlPlaneScope struct {
	client      client.Client
	patchHelper *patch.Helper

	Cluster                *clusterv1.Cluster
	GCPManagedCluster      *infrav1exp.GCPManagedCluster
	GCPManagedControlPlane *infrav1exp.GCPManagedControlPlane
	mcClient               *container.ClusterManagerClient
	tagBindingsClient      *resourcemanager.TagBindingsClient
	credentialsClient      *credentials.IamCredentialsClient
	credential             *Credential

	AllMachinePools        []clusterv1.MachinePool
	AllManagedMachinePools []infrav1exp.GCPManagedMachinePool
}

// PatchObject persists the managed control plane configuration and status.
func (s *ManagedControlPlaneScope) PatchObject() error {
	return s.patchHelper.Patch(
		context.TODO(),
		s.GCPManagedControlPlane,
		patch.WithOwnedConditions{Conditions: []clusterv1.ConditionType{
			infrav1exp.GKEControlPlaneReadyCondition,
			infrav1exp.GKEControlPlaneCreatingCondition,
			infrav1exp.GKEControlPlaneUpdatingCondition,
			infrav1exp.GKEControlPlaneDeletingCondition,
		}})
}

// Close closes the current scope persisting the managed control plane configuration and status.
func (s *ManagedControlPlaneScope) Close() error {
	s.mcClient.Close()
	s.tagBindingsClient.Close()
	s.credentialsClient.Close()
	return s.PatchObject()
}

// ConditionSetter return a condition setter (which is GCPManagedControlPlane itself).
func (s *ManagedControlPlaneScope) ConditionSetter() conditions.Setter {
	return s.GCPManagedControlPlane
}

// Client returns a k8s client.
func (s *ManagedControlPlaneScope) Client() client.Client {
	return s.client
}

// ManagedControlPlaneClient returns a client used to interact with GKE.
func (s *ManagedControlPlaneScope) ManagedControlPlaneClient() *container.ClusterManagerClient {
	return s.mcClient
}

// TagBindingsClient returns a client used to interact with resource manager tags.
func (s *ManagedControlPlaneScope) TagBindingsClient() *resourcemanager.TagBindingsClient {
	return s.tagBindingsClient
}

// CredentialsClient returns a client used to interact with IAM.
func (s *ManagedControlPlaneScope) CredentialsClient() *credentials.IamCredentialsClient {
	return s.credentialsClient
}

// GetCredential returns the credential data.
func (s *ManagedControlPlaneScope) GetCredential() *Credential {
	return s.credential
}

// GetAllNodePools gets all node pools for the control plane.
func (s *ManagedControlPlaneScope) GetAllNodePools(ctx context.Context) ([]infrav1exp.GCPManagedMachinePool, []clusterv1.MachinePool, error) {
	if len(s.AllManagedMachinePools) == 0 {
		listOptions := []client.ListOption{
			client.InNamespace(s.GCPManagedControlPlane.Namespace),
			client.MatchingLabels(map[string]string{clusterv1.ClusterNameLabel: s.Cluster.Name}),
		}

		machinePoolList := &clusterv1.MachinePoolList{}
		if err := s.client.List(ctx, machinePoolList, listOptions...); err != nil {
			return nil, nil, err
		}
		managedMachinePoolList := &infrav1exp.GCPManagedMachinePoolList{}
		if err := s.client.List(ctx, managedMachinePoolList, listOptions...); err != nil {
			return nil, nil, err
		}
		if len(machinePoolList.Items) != len(managedMachinePoolList.Items) {
			return nil, nil, fmt.Errorf("machinePoolList length (%d) != managedMachinePoolList length (%d)", len(machinePoolList.Items), len(managedMachinePoolList.Items))
		}
		s.AllMachinePools = machinePoolList.Items
		s.AllManagedMachinePools = managedMachinePoolList.Items
	}

	return s.AllManagedMachinePools, s.AllMachinePools, nil
}

// Region returns the region of the GKE cluster.
func (s *ManagedControlPlaneScope) Region() string {
	loc, _ := location.Parse(s.GCPManagedControlPlane.Spec.Location)
	return loc.Region
}

// ClusterLocation returns the location of the cluster.
func (s *ManagedControlPlaneScope) ClusterLocation() string {
	return fmt.Sprintf("projects/%s/locations/%s", s.GCPManagedControlPlane.Spec.Project, s.GCPManagedControlPlane.Spec.Location)
}

// ClusterFullName returns the full name of the cluster.
func (s *ManagedControlPlaneScope) ClusterFullName() string {
	return fmt.Sprintf("%s/clusters/%s", s.ClusterLocation(), s.GCPManagedControlPlane.Spec.ClusterName)
}

// ClusterName returns the name of the cluster.
func (s *ManagedControlPlaneScope) ClusterName() string {
	return s.GCPManagedControlPlane.Spec.ClusterName
}

// SetEndpoint sets the Endpoint of GCPManagedControlPlane.
func (s *ManagedControlPlaneScope) SetEndpoint(host string) {
	s.GCPManagedControlPlane.Spec.Endpoint = clusterv1.APIEndpoint{
		Host: host,
		Port: APIServerPort,
	}
}

// IsAutopilotCluster returns true if this is an autopilot cluster.
func (s *ManagedControlPlaneScope) IsAutopilotCluster() bool {
	return s.GCPManagedControlPlane.Spec.EnableAutopilot
}

// GetControlPlaneVersion returns the control plane version from the specification.
func (s *ManagedControlPlaneScope) GetControlPlaneVersion() *string {
	if s.GCPManagedControlPlane.Spec.Version != nil {
		return s.GCPManagedControlPlane.Spec.Version
	}
	if s.GCPManagedControlPlane.Spec.ControlPlaneVersion != nil {
		return s.GCPManagedControlPlane.Spec.ControlPlaneVersion
	}
	return nil
}
