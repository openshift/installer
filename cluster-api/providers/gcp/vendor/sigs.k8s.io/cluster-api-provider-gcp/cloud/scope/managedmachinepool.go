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
	"strings"

	"k8s.io/utils/ptr"
	infrav1 "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud"
	"sigs.k8s.io/cluster-api-provider-gcp/util/location"

	"sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"

	compute "cloud.google.com/go/compute/apiv1"
	container "cloud.google.com/go/container/apiv1"
	"cloud.google.com/go/container/apiv1/containerpb"
	"github.com/pkg/errors"
	infrav1exp "sigs.k8s.io/cluster-api-provider-gcp/exp/api/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	patch "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ManagedMachinePoolScopeParams defines the input parameters used to create a new Scope.
type ManagedMachinePoolScopeParams struct {
	ManagedClusterClient        *container.ClusterManagerClient
	InstanceGroupManagersClient *compute.InstanceGroupManagersClient
	Client                      client.Client
	Cluster                     *clusterv1.Cluster
	MachinePool                 *clusterv1.MachinePool
	GCPManagedCluster           *infrav1exp.GCPManagedCluster
	GCPManagedControlPlane      *infrav1exp.GCPManagedControlPlane
	GCPManagedMachinePool       *infrav1exp.GCPManagedMachinePool
}

// NewManagedMachinePoolScope creates a new Scope from the supplied parameters.
// This is meant to be called for each reconcile iteration.
func NewManagedMachinePoolScope(ctx context.Context, params ManagedMachinePoolScopeParams) (*ManagedMachinePoolScope, error) {
	if params.Cluster == nil {
		return nil, errors.New("failed to generate new scope from nil Cluster")
	}
	if params.MachinePool == nil {
		return nil, errors.New("failed to generate new scope from nil MachinePool")
	}
	if params.GCPManagedCluster == nil {
		return nil, errors.New("failed to generate new scope from nil GCPManagedCluster")
	}
	if params.GCPManagedControlPlane == nil {
		return nil, errors.New("failed to generate new scope from nil GCPManagedControlPlane")
	}
	if params.GCPManagedMachinePool == nil {
		return nil, errors.New("failed to generate new scope from nil GCPManagedMachinePool")
	}

	if params.ManagedClusterClient == nil {
		managedClusterClient, err := newClusterManagerClient(ctx, params.GCPManagedCluster.Spec.CredentialsRef, params.Client, params.GCPManagedCluster.Spec.ServiceEndpoints)
		if err != nil {
			return nil, errors.Errorf("failed to create gcp managed cluster client: %v", err)
		}
		params.ManagedClusterClient = managedClusterClient
	}
	if params.InstanceGroupManagersClient == nil {
		instanceGroupManagersClient, err := newInstanceGroupManagerClient(ctx, params.GCPManagedCluster.Spec.CredentialsRef, params.Client, params.GCPManagedCluster.Spec.ServiceEndpoints)
		if err != nil {
			return nil, errors.Errorf("failed to create gcp instance group manager client: %v", err)
		}
		params.InstanceGroupManagersClient = instanceGroupManagersClient
	}

	helper, err := patch.NewHelper(params.GCPManagedMachinePool, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}

	return &ManagedMachinePoolScope{
		client:                 params.Client,
		Cluster:                params.Cluster,
		MachinePool:            params.MachinePool,
		GCPManagedControlPlane: params.GCPManagedControlPlane,
		GCPManagedMachinePool:  params.GCPManagedMachinePool,
		mcClient:               params.ManagedClusterClient,
		migClient:              params.InstanceGroupManagersClient,
		patchHelper:            helper,
	}, nil
}

// ManagedMachinePoolScope defines the basic context for an actuator to operate upon.
type ManagedMachinePoolScope struct {
	client      client.Client
	patchHelper *patch.Helper

	Cluster                *clusterv1.Cluster
	MachinePool            *clusterv1.MachinePool
	GCPManagedCluster      *infrav1exp.GCPManagedCluster
	GCPManagedControlPlane *infrav1exp.GCPManagedControlPlane
	GCPManagedMachinePool  *infrav1exp.GCPManagedMachinePool
	mcClient               *container.ClusterManagerClient
	migClient              *compute.InstanceGroupManagersClient
}

// PatchObject persists the managed control plane configuration and status.
func (s *ManagedMachinePoolScope) PatchObject() error {
	return s.patchHelper.Patch(
		context.TODO(),
		s.GCPManagedMachinePool,
		patch.WithOwnedConditions{Conditions: []clusterv1.ConditionType{
			infrav1exp.GKEMachinePoolReadyCondition,
			infrav1exp.GKEMachinePoolCreatingCondition,
			infrav1exp.GKEMachinePoolUpdatingCondition,
			infrav1exp.GKEMachinePoolDeletingCondition,
		}})
}

// Close closes the current scope persisting the managed control plane configuration and status.
func (s *ManagedMachinePoolScope) Close() error {
	s.mcClient.Close()
	s.migClient.Close()
	return s.PatchObject()
}

// ConditionSetter return a condition setter (which is GCPManagedMachinePool itself).
func (s *ManagedMachinePoolScope) ConditionSetter() conditions.Setter {
	return s.GCPManagedMachinePool
}

// ManagedMachinePoolClient returns a client used to interact with GKE.
func (s *ManagedMachinePoolScope) ManagedMachinePoolClient() *container.ClusterManagerClient {
	return s.mcClient
}

// InstanceGroupManagersClient returns a client used to interact with GCP MIG.
func (s *ManagedMachinePoolScope) InstanceGroupManagersClient() *compute.InstanceGroupManagersClient {
	return s.migClient
}

// NodePoolVersion returns the k8s version of the node pool.
func (s *ManagedMachinePoolScope) NodePoolVersion() *string {
	return s.MachinePool.Spec.Template.Spec.Version
}

// NodePoolResourceLabels returns the resource labels of the node pool.
func NodePoolResourceLabels(additionalLabels infrav1.Labels, clusterName string) infrav1.Labels {
	if additionalLabels == nil {
		additionalLabels = infrav1.Labels{}
	}
	resourceLabels := additionalLabels.DeepCopy()
	resourceLabels[infrav1.ClusterTagKey(clusterName)] = string(infrav1.ResourceLifecycleOwned)
	return resourceLabels
}

// ConvertToSdkNodePool converts a node pool to format that is used by GCP SDK.
func ConvertToSdkNodePool(nodePool infrav1exp.GCPManagedMachinePool, machinePool clusterv1.MachinePool, regional bool, clusterName string) *containerpb.NodePool {
	replicas := *machinePool.Spec.Replicas
	if regional {
		if len(nodePool.Spec.NodeLocations) != 0 {
			replicas /= int32(len(nodePool.Spec.NodeLocations))
		} else {
			replicas /= cloud.DefaultNumRegionsPerZone
		}
	}
	nodePoolName := nodePool.Spec.NodePoolName
	if len(nodePoolName) == 0 {
		// Use the GCPManagedMachinePool CR name if nodePoolName is not specified
		nodePoolName = nodePool.Name
	}
	// build node pool in GCP SDK format using the GCPManagedMachinePool spec
	sdkNodePool := containerpb.NodePool{
		Name:             nodePoolName,
		InitialNodeCount: replicas,
		Config: &containerpb.NodeConfig{
			Labels: nodePool.Spec.KubernetesLabels,
			Taints: infrav1exp.ConvertToSdkTaint(nodePool.Spec.KubernetesTaints),
			ShieldedInstanceConfig: &containerpb.ShieldedInstanceConfig{
				EnableSecureBoot:          ptr.Deref(nodePool.Spec.NodeSecurity.EnableSecureBoot, false),
				EnableIntegrityMonitoring: ptr.Deref(nodePool.Spec.NodeSecurity.EnableIntegrityMonitoring, false),
			},
			ResourceLabels: NodePoolResourceLabels(nodePool.Spec.AdditionalLabels, clusterName),
		},
	}
	if nodePool.Spec.MachineType != nil {
		sdkNodePool.Config.MachineType = *nodePool.Spec.MachineType
	}
	if nodePool.Spec.DiskSizeGb != nil {
		sdkNodePool.Config.DiskSizeGb = *nodePool.Spec.DiskSizeGb
	}
	if nodePool.Spec.ImageType != nil {
		sdkNodePool.Config.ImageType = *nodePool.Spec.ImageType
	}
	if nodePool.Spec.LocalSsdCount != nil {
		sdkNodePool.Config.LocalSsdCount = *nodePool.Spec.LocalSsdCount
	}
	if nodePool.Spec.DiskType != nil {
		sdkNodePool.Config.DiskType = string(*nodePool.Spec.DiskType)
	}
	if nodePool.Spec.Scaling != nil {
		sdkNodePool.Autoscaling = infrav1exp.ConvertToSdkAutoscaling(nodePool.Spec.Scaling)
	}
	if nodePool.Spec.LinuxNodeConfig != nil {
		sdkNodePool.Config.LinuxNodeConfig = infrav1exp.ConvertToSdkLinuxNodeConfig(nodePool.Spec.LinuxNodeConfig)
	}
	if nodePool.Spec.Management != nil {
		sdkNodePool.Management = &containerpb.NodeManagement{
			AutoRepair:  nodePool.Spec.Management.AutoRepair,
			AutoUpgrade: nodePool.Spec.Management.AutoUpgrade,
		}
	}
	if nodePool.Spec.MaxPodsPerNode != nil {
		sdkNodePool.MaxPodsConstraint = &containerpb.MaxPodsConstraint{
			MaxPodsPerNode: *nodePool.Spec.MaxPodsPerNode,
		}
	}
	if nodePool.Spec.InstanceType != nil {
		sdkNodePool.Config.MachineType = *nodePool.Spec.InstanceType
	}
	if nodePool.Spec.ImageType != nil {
		sdkNodePool.Config.ImageType = *nodePool.Spec.ImageType
	}
	if nodePool.Spec.DiskType != nil {
		sdkNodePool.Config.DiskType = string(*nodePool.Spec.DiskType)
	}
	if nodePool.Spec.DiskSizeGB != nil {
		sdkNodePool.Config.DiskSizeGb = int32(*nodePool.Spec.DiskSizeGB) //nolint:gosec
	}
	if len(nodePool.Spec.NodeNetwork.Tags) != 0 {
		sdkNodePool.Config.Tags = nodePool.Spec.NodeNetwork.Tags
	}
	if nodePool.Spec.NodeSecurity.ServiceAccount.Email != nil {
		sdkNodePool.Config.ServiceAccount = *nodePool.Spec.NodeSecurity.ServiceAccount.Email
	}
	if len(nodePool.Spec.NodeSecurity.ServiceAccount.Scopes) != 0 {
		sdkNodePool.Config.OauthScopes = nodePool.Spec.NodeSecurity.ServiceAccount.Scopes
	}
	if len(nodePool.Spec.NodeLocations) != 0 {
		sdkNodePool.Locations = nodePool.Spec.NodeLocations
	}
	if nodePool.Spec.MaxPodsPerNode != nil {
		sdkNodePool.MaxPodsConstraint = &containerpb.MaxPodsConstraint{
			MaxPodsPerNode: *nodePool.Spec.MaxPodsPerNode,
		}
	}
	if nodePool.Spec.NodeNetwork.CreatePodRange != nil && nodePool.Spec.NodeNetwork.PodRangeName != nil && nodePool.Spec.NodeNetwork.PodRangeCidrBlock != nil {
		sdkNodePool.NetworkConfig = &containerpb.NodeNetworkConfig{
			CreatePodRange:   *nodePool.Spec.NodeNetwork.CreatePodRange,
			PodRange:         *nodePool.Spec.NodeNetwork.PodRangeName,
			PodIpv4CidrBlock: *nodePool.Spec.NodeNetwork.PodRangeCidrBlock,
		}
	}

	if ptr.Deref(nodePool.Spec.NodeSecurity.SandboxType, "") == "GVISOR" {
		sdkNodePool.Config.SandboxConfig = &containerpb.SandboxConfig{
			Type: containerpb.SandboxConfig_GVISOR,
		}
	}
	if machinePool.Spec.Template.Spec.Version != nil {
		sdkNodePool.Version = strings.Replace(*machinePool.Spec.Template.Spec.Version, "v", "", 1)
	}
	return &sdkNodePool
}

// ConvertToSdkNodePools converts node pools to format that is used by GCP SDK.
func ConvertToSdkNodePools(nodePools []infrav1exp.GCPManagedMachinePool, machinePools []clusterv1.MachinePool, regional bool, clusterName string) []*containerpb.NodePool {
	res := []*containerpb.NodePool{}
	for i := range nodePools {
		res = append(res, ConvertToSdkNodePool(nodePools[i], machinePools[i], regional, clusterName))
	}
	return res
}

// SetReplicas sets the replicas count in status.
func (s *ManagedMachinePoolScope) SetReplicas(replicas int32) {
	s.GCPManagedMachinePool.Status.Replicas = replicas
}

// NodePoolName returns the node pool name.
func (s *ManagedMachinePoolScope) NodePoolName() string {
	if len(s.GCPManagedMachinePool.Spec.NodePoolName) > 0 {
		return s.GCPManagedMachinePool.Spec.NodePoolName
	}
	return s.GCPManagedMachinePool.Name
}

// Region returns the region of the GKE node pool.
func (s *ManagedMachinePoolScope) Region() string {
	loc, _ := location.Parse(s.GCPManagedControlPlane.Spec.Location)
	return loc.Region
}

// NodePoolLocation returns the location of the node pool.
func (s *ManagedMachinePoolScope) NodePoolLocation() string {
	return fmt.Sprintf("projects/%s/locations/%s/clusters/%s", s.GCPManagedControlPlane.Spec.Project, s.Region(), s.GCPManagedControlPlane.Spec.ClusterName)
}

// NodePoolFullName returns the full name of the node pool.
func (s *ManagedMachinePoolScope) NodePoolFullName() string {
	return fmt.Sprintf("%s/nodePools/%s", s.NodePoolLocation(), s.NodePoolName())
}

// SetInfrastructureMachineKind sets the infrastructure machine kind in the status if it is not set already, returning
// `true` if the status was updated. This supports MachinePool Machines.
func (s *ManagedMachinePoolScope) SetInfrastructureMachineKind() bool {
	if s.GCPManagedMachinePool.Status.InfrastructureMachineKind != infrav1exp.GCPManagedMachinePoolMachineKind {
		s.GCPManagedMachinePool.Status.InfrastructureMachineKind = infrav1exp.GCPManagedMachinePoolMachineKind

		return true
	}

	return false
}
