/*
Copyright 2022 The Kubernetes Authors.

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

	"github.com/pkg/errors"
	"k8s.io/utils/ptr"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/agentpools"
	"sigs.k8s.io/cluster-api-provider-azure/util/futures"
	"sigs.k8s.io/cluster-api-provider-azure/util/maps"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	expv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ManagedMachinePoolScopeParams defines the input parameters used to create a new managed
// control plane.
type ManagedMachinePoolScopeParams struct {
	ManagedMachinePool
	Client                   client.Client
	Cluster                  *clusterv1.Cluster
	ControlPlane             *infrav1.AzureManagedControlPlane
	ManagedControlPlaneScope azure.ManagedClusterScoper
}

// ManagedMachinePool defines the scope interface for a managed machine pool.
type ManagedMachinePool struct {
	InfraMachinePool *infrav1.AzureManagedMachinePool
	MachinePool      *expv1.MachinePool
}

// NewManagedMachinePoolScope creates a new Scope from the supplied parameters.
// This is meant to be called for each reconcile iteration.
func NewManagedMachinePoolScope(ctx context.Context, params ManagedMachinePoolScopeParams) (*ManagedMachinePoolScope, error) {
	_, _, done := tele.StartSpanWithLogger(ctx, "scope.NewManagedMachinePoolScope")
	defer done()

	if params.Cluster == nil {
		return nil, errors.New("failed to generate new scope from nil Cluster")
	}

	if params.ControlPlane == nil {
		return nil, errors.New("failed to generate new scope from nil ControlPlane")
	}

	helper, err := patch.NewHelper(params.InfraMachinePool, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}

	capiMachinePoolPatchHelper, err := patch.NewHelper(params.MachinePool, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}

	return &ManagedMachinePoolScope{
		Client:                     params.Client,
		Cluster:                    params.Cluster,
		ControlPlane:               params.ControlPlane,
		MachinePool:                params.MachinePool,
		InfraMachinePool:           params.InfraMachinePool,
		patchHelper:                helper,
		capiMachinePoolPatchHelper: capiMachinePoolPatchHelper,
		ManagedClusterScoper:       params.ManagedControlPlaneScope,
	}, nil
}

// ManagedMachinePoolScope defines the basic context for an actuator to operate upon.
type ManagedMachinePoolScope struct {
	Client                     client.Client
	patchHelper                *patch.Helper
	capiMachinePoolPatchHelper *patch.Helper

	azure.ManagedClusterScoper
	Cluster          *clusterv1.Cluster
	MachinePool      *expv1.MachinePool
	ControlPlane     *infrav1.AzureManagedControlPlane
	InfraMachinePool *infrav1.AzureManagedMachinePool
}

// PatchObject persists the cluster configuration and status.
func (s *ManagedMachinePoolScope) PatchObject(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "scope.ManagedMachinePoolScope.PatchObject")
	defer done()

	conditions.SetSummary(s.InfraMachinePool)

	return s.patchHelper.Patch(
		ctx,
		s.InfraMachinePool,
		patch.WithOwnedConditions{Conditions: []clusterv1.ConditionType{
			clusterv1.ReadyCondition,
		}})
}

// Close closes the current scope persisting the cluster configuration and status.
func (s *ManagedMachinePoolScope) Close(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "scope.ManagedMachinePoolScope.Close")
	defer done()

	return s.PatchObject(ctx)
}

// AgentPoolAnnotations returns a map of annotations for the infra machine pool.
func (s *ManagedMachinePoolScope) AgentPoolAnnotations() map[string]string {
	return s.InfraMachinePool.Annotations
}

// Name returns the name of the infra machine pool.
func (s *ManagedMachinePoolScope) Name() string {
	return s.InfraMachinePool.Name
}

// SetSubnetName updates AzureManagedMachinePool.SubnetName if AzureManagedMachinePool.SubnetName is empty with s.ControlPlane.Spec.VirtualNetwork.Subnet.Name.
func (s *ManagedMachinePoolScope) SetSubnetName() {
	s.InfraMachinePool.Spec.SubnetName = getAgentPoolSubnet(s.ControlPlane, s.InfraMachinePool)
}

// AgentPoolSpec returns an azure.ResourceSpecGetter for currently reconciled AzureManagedMachinePool.
func (s *ManagedMachinePoolScope) AgentPoolSpec() azure.ResourceSpecGetter {
	return buildAgentPoolSpec(s.ControlPlane, s.MachinePool, s.InfraMachinePool, s.AgentPoolAnnotations())
}

func getAgentPoolSubnet(controlPlane *infrav1.AzureManagedControlPlane, infraMachinePool *infrav1.AzureManagedMachinePool) *string {
	if infraMachinePool.Spec.SubnetName == nil {
		return ptr.To(controlPlane.Spec.VirtualNetwork.Subnet.Name)
	}
	return infraMachinePool.Spec.SubnetName
}

func buildAgentPoolSpec(managedControlPlane *infrav1.AzureManagedControlPlane,
	machinePool *expv1.MachinePool,
	managedMachinePool *infrav1.AzureManagedMachinePool,
	agentPoolAnnotations map[string]string) azure.ResourceSpecGetter {
	var normalizedVersion *string
	if machinePool.Spec.Template.Spec.Version != nil {
		v := strings.TrimPrefix(*machinePool.Spec.Template.Spec.Version, "v")
		normalizedVersion = &v
	}

	replicas := int32(1)
	if machinePool.Spec.Replicas != nil {
		replicas = *machinePool.Spec.Replicas
	}

	agentPoolSpec := &agentpools.AgentPoolSpec{
		Name:          ptr.Deref(managedMachinePool.Spec.Name, ""),
		ResourceGroup: managedControlPlane.Spec.ResourceGroupName,
		Cluster:       managedControlPlane.Name,
		SKU:           managedMachinePool.Spec.SKU,
		Replicas:      replicas,
		Version:       normalizedVersion,
		OSType:        managedMachinePool.Spec.OSType,
		VnetSubnetID: azure.SubnetID(
			managedControlPlane.Spec.SubscriptionID,
			managedControlPlane.Spec.VirtualNetwork.ResourceGroup,
			managedControlPlane.Spec.VirtualNetwork.Name,
			ptr.Deref(getAgentPoolSubnet(managedControlPlane, managedMachinePool), ""),
		),
		Mode:                 managedMachinePool.Spec.Mode,
		MaxPods:              managedMachinePool.Spec.MaxPods,
		AvailabilityZones:    managedMachinePool.Spec.AvailabilityZones,
		OsDiskType:           managedMachinePool.Spec.OsDiskType,
		EnableUltraSSD:       managedMachinePool.Spec.EnableUltraSSD,
		Headers:              maps.FilterByKeyPrefix(agentPoolAnnotations, infrav1.CustomHeaderPrefix),
		EnableNodePublicIP:   managedMachinePool.Spec.EnableNodePublicIP,
		NodePublicIPPrefixID: managedMachinePool.Spec.NodePublicIPPrefixID,
		ScaleSetPriority:     managedMachinePool.Spec.ScaleSetPriority,
		ScaleDownMode:        managedMachinePool.Spec.ScaleDownMode,
		SpotMaxPrice:         managedMachinePool.Spec.SpotMaxPrice,
		AdditionalTags:       managedMachinePool.Spec.AdditionalTags,
		KubeletDiskType:      managedMachinePool.Spec.KubeletDiskType,
		LinuxOSConfig:        managedMachinePool.Spec.LinuxOSConfig,
		EnableFIPS:           managedMachinePool.Spec.EnableFIPS,
	}

	if managedMachinePool.Spec.OSDiskSizeGB != nil {
		agentPoolSpec.OSDiskSizeGB = *managedMachinePool.Spec.OSDiskSizeGB
	}

	if len(managedMachinePool.Spec.Taints) > 0 {
		nodeTaints := make([]string, 0, len(managedMachinePool.Spec.Taints))
		for _, t := range managedMachinePool.Spec.Taints {
			nodeTaints = append(nodeTaints, fmt.Sprintf("%s=%s:%s", t.Key, t.Value, t.Effect))
		}
		agentPoolSpec.NodeTaints = nodeTaints
	}

	if managedMachinePool.Spec.Scaling != nil {
		agentPoolSpec.EnableAutoScaling = true
		agentPoolSpec.MaxCount = managedMachinePool.Spec.Scaling.MaxSize
		agentPoolSpec.MinCount = managedMachinePool.Spec.Scaling.MinSize
	}

	if len(managedMachinePool.Spec.NodeLabels) > 0 {
		agentPoolSpec.NodeLabels = make(map[string]*string, len(managedMachinePool.Spec.NodeLabels))
		for k, v := range managedMachinePool.Spec.NodeLabels {
			agentPoolSpec.NodeLabels[k] = ptr.To(v)
		}
	}

	if managedMachinePool.Spec.KubeletConfig != nil {
		agentPoolSpec.KubeletConfig = &agentpools.KubeletConfig{
			CPUManagerPolicy:      (*string)(managedMachinePool.Spec.KubeletConfig.CPUManagerPolicy),
			CPUCfsQuota:           managedMachinePool.Spec.KubeletConfig.CPUCfsQuota,
			CPUCfsQuotaPeriod:     managedMachinePool.Spec.KubeletConfig.CPUCfsQuotaPeriod,
			ImageGcHighThreshold:  managedMachinePool.Spec.KubeletConfig.ImageGcHighThreshold,
			ImageGcLowThreshold:   managedMachinePool.Spec.KubeletConfig.ImageGcLowThreshold,
			TopologyManagerPolicy: (*string)(managedMachinePool.Spec.KubeletConfig.TopologyManagerPolicy),
			FailSwapOn:            managedMachinePool.Spec.KubeletConfig.FailSwapOn,
			ContainerLogMaxSizeMB: managedMachinePool.Spec.KubeletConfig.ContainerLogMaxSizeMB,
			ContainerLogMaxFiles:  managedMachinePool.Spec.KubeletConfig.ContainerLogMaxFiles,
			PodMaxPids:            managedMachinePool.Spec.KubeletConfig.PodMaxPids,
		}
		if len(managedMachinePool.Spec.KubeletConfig.AllowedUnsafeSysctls) > 0 {
			agentPoolSpec.KubeletConfig.AllowedUnsafeSysctls = &managedMachinePool.Spec.KubeletConfig.AllowedUnsafeSysctls
		}
	}

	return agentPoolSpec
}

// SetAgentPoolProviderIDList sets a list of agent pool's Azure VM IDs.
func (s *ManagedMachinePoolScope) SetAgentPoolProviderIDList(providerIDs []string) {
	s.InfraMachinePool.Spec.ProviderIDList = providerIDs
}

// SetAgentPoolReplicas sets the number of agent pool replicas.
func (s *ManagedMachinePoolScope) SetAgentPoolReplicas(replicas int32) {
	s.InfraMachinePool.Status.Replicas = replicas
}

// SetAgentPoolReady sets the flag that indicates if the agent pool is ready or not.
func (s *ManagedMachinePoolScope) SetAgentPoolReady(ready bool) {
	s.InfraMachinePool.Status.Ready = ready
}

// SetLongRunningOperationState will set the future on the AzureManagedMachinePool status to allow the resource to continue
// in the next reconciliation.
func (s *ManagedMachinePoolScope) SetLongRunningOperationState(future *infrav1.Future) {
	futures.Set(s.InfraMachinePool, future)
}

// GetLongRunningOperationState will get the future on the AzureManagedMachinePool status.
func (s *ManagedMachinePoolScope) GetLongRunningOperationState(name, service, futureType string) *infrav1.Future {
	return futures.Get(s.InfraMachinePool, name, service, futureType)
}

// DeleteLongRunningOperationState will delete the future from the AzureManagedMachinePool status.
func (s *ManagedMachinePoolScope) DeleteLongRunningOperationState(name, service, futureType string) {
	futures.Delete(s.InfraMachinePool, name, service, futureType)
}

// UpdateDeleteStatus updates a condition on the AzureManagedControlPlane status after a DELETE operation.
func (s *ManagedMachinePoolScope) UpdateDeleteStatus(condition clusterv1.ConditionType, service string, err error) {
	switch {
	case err == nil:
		conditions.MarkFalse(s.InfraMachinePool, condition, infrav1.DeletedReason, clusterv1.ConditionSeverityInfo, "%s successfully deleted", service)
	case azure.IsOperationNotDoneError(err):
		conditions.MarkFalse(s.InfraMachinePool, condition, infrav1.DeletingReason, clusterv1.ConditionSeverityInfo, "%s deleting", service)
	default:
		conditions.MarkFalse(s.InfraMachinePool, condition, infrav1.DeletionFailedReason, clusterv1.ConditionSeverityError, "%s failed to delete. err: %s", service, err.Error())
	}
}

// UpdatePutStatus updates a condition on the AzureManagedMachinePool status after a PUT operation.
func (s *ManagedMachinePoolScope) UpdatePutStatus(condition clusterv1.ConditionType, service string, err error) {
	switch {
	case err == nil:
		conditions.MarkTrue(s.InfraMachinePool, condition)
	case azure.IsOperationNotDoneError(err):
		conditions.MarkFalse(s.InfraMachinePool, condition, infrav1.CreatingReason, clusterv1.ConditionSeverityInfo, "%s creating or updating", service)
	default:
		conditions.MarkFalse(s.InfraMachinePool, condition, infrav1.FailedReason, clusterv1.ConditionSeverityError, "%s failed to create or update. err: %s", service, err.Error())
	}
}

// UpdatePatchStatus updates a condition on the AzureManagedMachinePool status after a PATCH operation.
func (s *ManagedMachinePoolScope) UpdatePatchStatus(condition clusterv1.ConditionType, service string, err error) {
	switch {
	case err == nil:
		conditions.MarkTrue(s.InfraMachinePool, condition)
	case azure.IsOperationNotDoneError(err):
		conditions.MarkFalse(s.InfraMachinePool, condition, infrav1.UpdatingReason, clusterv1.ConditionSeverityInfo, "%s updating", service)
	default:
		conditions.MarkFalse(s.InfraMachinePool, condition, infrav1.FailedReason, clusterv1.ConditionSeverityError, "%s failed to update. err: %s", service, err.Error())
	}
}

// PatchCAPIMachinePoolObject persists the capi machinepool configuration and status.
func (s *ManagedMachinePoolScope) PatchCAPIMachinePoolObject(ctx context.Context) error {
	return s.capiMachinePoolPatchHelper.Patch(
		ctx,
		s.MachinePool,
	)
}

// SetCAPIMachinePoolReplicas sets the associated MachinePool replica count.
func (s *ManagedMachinePoolScope) SetCAPIMachinePoolReplicas(replicas *int32) {
	s.MachinePool.Spec.Replicas = replicas
}

// SetCAPIMachinePoolAnnotation sets the specified annotation on the associated MachinePool.
func (s *ManagedMachinePoolScope) SetCAPIMachinePoolAnnotation(key, value string) {
	if s.MachinePool.Annotations == nil {
		s.MachinePool.Annotations = make(map[string]string)
	}
	s.MachinePool.Annotations[key] = value
}

// RemoveCAPIMachinePoolAnnotation removes the specified annotation on the associated MachinePool.
func (s *ManagedMachinePoolScope) RemoveCAPIMachinePoolAnnotation(key string) {
	delete(s.MachinePool.Annotations, key)
}

// GetCAPIMachinePoolAnnotation gets the specified annotation on the associated MachinePool.
func (s *ManagedMachinePoolScope) GetCAPIMachinePoolAnnotation(key string) (success bool, value string) {
	val, ok := s.MachinePool.Annotations[key]
	return ok, val
}
