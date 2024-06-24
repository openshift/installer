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

package clusters

import (
	"context"
	"fmt"
	"strings"

	"sigs.k8s.io/cluster-api-provider-gcp/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud/services/shared"

	"cloud.google.com/go/container/apiv1/containerpb"
	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	"github.com/googleapis/gax-go/v2/apierror"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	infrav1exp "sigs.k8s.io/cluster-api-provider-gcp/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-gcp/util/reconciler"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// Reconcile reconcile GKE cluster.
func (s *Service) Reconcile(ctx context.Context) (ctrl.Result, error) {
	log := log.FromContext(ctx).WithValues("service", "container.clusters")
	log.Info("Reconciling cluster resources")

	cluster, err := s.describeCluster(ctx, &log)
	if err != nil {
		s.scope.GCPManagedControlPlane.Status.Initialized = false
		s.scope.GCPManagedControlPlane.Status.Ready = false
		conditions.MarkFalse(s.scope.ConditionSetter(), clusterv1.ReadyCondition, infrav1exp.GKEControlPlaneReconciliationFailedReason, clusterv1.ConditionSeverityError, err.Error())
		return ctrl.Result{}, err
	}
	if cluster == nil {
		log.Info("Cluster not found, creating")
		s.scope.GCPManagedControlPlane.Status.Initialized = false
		s.scope.GCPManagedControlPlane.Status.Ready = false

		nodePools, _, err := s.scope.GetAllNodePools(ctx)
		if err != nil {
			conditions.MarkFalse(s.scope.ConditionSetter(), clusterv1.ReadyCondition, infrav1exp.GKEControlPlaneReconciliationFailedReason, clusterv1.ConditionSeverityError, err.Error())
			conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEControlPlaneReadyCondition, infrav1exp.GKEControlPlaneReconciliationFailedReason, clusterv1.ConditionSeverityError, err.Error())
			conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEControlPlaneCreatingCondition, infrav1exp.GKEControlPlaneReconciliationFailedReason, clusterv1.ConditionSeverityError, err.Error())
			return ctrl.Result{}, err
		}
		if s.scope.IsAutopilotCluster() {
			if len(nodePools) > 0 {
				log.Error(ErrAutopilotClusterMachinePoolsNotAllowed, fmt.Sprintf("%d machine pools defined", len(nodePools)))
				conditions.MarkFalse(s.scope.ConditionSetter(), clusterv1.ReadyCondition, infrav1exp.GKEControlPlaneRequiresAtLeastOneNodePoolReason, clusterv1.ConditionSeverityInfo, "")
				conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEControlPlaneReadyCondition, infrav1exp.GKEControlPlaneRequiresAtLeastOneNodePoolReason, clusterv1.ConditionSeverityInfo, "")
				conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEControlPlaneCreatingCondition, infrav1exp.GKEControlPlaneRequiresAtLeastOneNodePoolReason, clusterv1.ConditionSeverityInfo, "")
				return ctrl.Result{}, ErrAutopilotClusterMachinePoolsNotAllowed
			}
		} else {
			if len(nodePools) == 0 {
				log.Info("At least 1 node pool is required to create GKE cluster with autopilot disabled")
				conditions.MarkFalse(s.scope.ConditionSetter(), clusterv1.ReadyCondition, infrav1exp.GKEControlPlaneRequiresAtLeastOneNodePoolReason, clusterv1.ConditionSeverityInfo, "")
				conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEControlPlaneReadyCondition, infrav1exp.GKEControlPlaneRequiresAtLeastOneNodePoolReason, clusterv1.ConditionSeverityInfo, "")
				conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEControlPlaneCreatingCondition, infrav1exp.GKEControlPlaneRequiresAtLeastOneNodePoolReason, clusterv1.ConditionSeverityInfo, "")
				return ctrl.Result{RequeueAfter: reconciler.DefaultRetryTime}, nil
			}
		}

		if err = s.createCluster(ctx, &log); err != nil {
			log.Error(err, "failed creating cluster")
			conditions.MarkFalse(s.scope.ConditionSetter(), clusterv1.ReadyCondition, infrav1exp.GKEControlPlaneReconciliationFailedReason, clusterv1.ConditionSeverityError, err.Error())
			conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEControlPlaneReadyCondition, infrav1exp.GKEControlPlaneReconciliationFailedReason, clusterv1.ConditionSeverityError, err.Error())
			conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEControlPlaneCreatingCondition, infrav1exp.GKEControlPlaneReconciliationFailedReason, clusterv1.ConditionSeverityError, err.Error())
			return ctrl.Result{}, err
		}
		log.Info("Cluster created provisioning in progress")
		conditions.MarkFalse(s.scope.ConditionSetter(), clusterv1.ReadyCondition, infrav1exp.GKEControlPlaneCreatingReason, clusterv1.ConditionSeverityInfo, "")
		conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEControlPlaneReadyCondition, infrav1exp.GKEControlPlaneCreatingReason, clusterv1.ConditionSeverityInfo, "")
		conditions.MarkTrue(s.scope.ConditionSetter(), infrav1exp.GKEControlPlaneCreatingCondition)
		return ctrl.Result{RequeueAfter: reconciler.DefaultRetryTime}, nil
	}

	log.V(2).Info("gke cluster found", "status", cluster.GetStatus())
	s.scope.GCPManagedControlPlane.Status.CurrentVersion = convertToSdkMasterVersion(cluster.GetCurrentMasterVersion())

	switch cluster.GetStatus() {
	case containerpb.Cluster_PROVISIONING:
		log.Info("Cluster provisioning in progress")
		conditions.MarkFalse(s.scope.ConditionSetter(), clusterv1.ReadyCondition, infrav1exp.GKEControlPlaneCreatingReason, clusterv1.ConditionSeverityInfo, "")
		conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEControlPlaneReadyCondition, infrav1exp.GKEControlPlaneCreatingReason, clusterv1.ConditionSeverityInfo, "")
		conditions.MarkTrue(s.scope.ConditionSetter(), infrav1exp.GKEControlPlaneCreatingCondition)
		s.scope.GCPManagedControlPlane.Status.Initialized = false
		s.scope.GCPManagedControlPlane.Status.Ready = false
		return ctrl.Result{RequeueAfter: reconciler.DefaultRetryTime}, nil
	case containerpb.Cluster_RECONCILING:
		log.Info("Cluster reconciling in progress")
		conditions.MarkTrue(s.scope.ConditionSetter(), infrav1exp.GKEControlPlaneUpdatingCondition)
		s.scope.GCPManagedControlPlane.Status.Initialized = true
		s.scope.GCPManagedControlPlane.Status.Ready = true
		return ctrl.Result{RequeueAfter: reconciler.DefaultRetryTime}, nil
	case containerpb.Cluster_STOPPING:
		log.Info("Cluster stopping in progress")
		conditions.MarkFalse(s.scope.ConditionSetter(), clusterv1.ReadyCondition, infrav1exp.GKEControlPlaneDeletingReason, clusterv1.ConditionSeverityInfo, "")
		conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEControlPlaneReadyCondition, infrav1exp.GKEControlPlaneDeletingReason, clusterv1.ConditionSeverityInfo, "")
		conditions.MarkTrue(s.scope.ConditionSetter(), infrav1exp.GKEControlPlaneDeletingCondition)
		s.scope.GCPManagedControlPlane.Status.Initialized = false
		s.scope.GCPManagedControlPlane.Status.Ready = false
		return ctrl.Result{RequeueAfter: reconciler.DefaultRetryTime}, nil
	case containerpb.Cluster_ERROR, containerpb.Cluster_DEGRADED:
		var msg string
		if len(cluster.GetConditions()) > 0 {
			msg = cluster.GetConditions()[0].GetMessage()
		}
		log.Error(errors.New("Cluster in error/degraded state"), msg, "name", s.scope.ClusterName())
		conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEControlPlaneReadyCondition, infrav1exp.GKEControlPlaneErrorReason, clusterv1.ConditionSeverityError, "")
		s.scope.GCPManagedControlPlane.Status.Ready = false
		s.scope.GCPManagedControlPlane.Status.Initialized = false
		return ctrl.Result{}, nil
	case containerpb.Cluster_RUNNING:
		log.Info("Cluster running")
	default:
		statusErr := NewErrUnexpectedClusterStatus(string(cluster.GetStatus()))
		log.Error(statusErr, fmt.Sprintf("Unhandled cluster status %s", cluster.GetStatus()), "name", s.scope.ClusterName())
		return ctrl.Result{}, statusErr
	}

	needUpdate, updateClusterRequest := s.checkDiffAndPrepareUpdate(cluster, &log)
	if needUpdate {
		log.Info("Update required")
		err = s.updateCluster(ctx, updateClusterRequest, &log)
		if err != nil {
			return ctrl.Result{}, err
		}
		log.Info("Cluster updating in progress")
		conditions.MarkTrue(s.scope.ConditionSetter(), infrav1exp.GKEControlPlaneUpdatingCondition)
		s.scope.GCPManagedControlPlane.Status.Initialized = true
		s.scope.GCPManagedControlPlane.Status.Ready = true
		return ctrl.Result{}, nil
	}
	conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEControlPlaneUpdatingCondition, infrav1exp.GKEControlPlaneUpdatedReason, clusterv1.ConditionSeverityInfo, "")

	// Reconcile kubeconfig
	err = s.reconcileKubeconfig(ctx, cluster, &log)
	if err != nil {
		log.Error(err, "Failed to reconcile CAPI kubeconfig")
		return ctrl.Result{}, err
	}
	err = s.reconcileAdditionalKubeconfigs(ctx, cluster, &log)
	if err != nil {
		log.Error(err, "Failed to reconcile additional kubeconfig")
		return ctrl.Result{}, err
	}

	s.scope.SetEndpoint(cluster.GetEndpoint())
	conditions.MarkTrue(s.scope.ConditionSetter(), clusterv1.ReadyCondition)
	conditions.MarkTrue(s.scope.ConditionSetter(), infrav1exp.GKEControlPlaneReadyCondition)
	conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEControlPlaneCreatingCondition, infrav1exp.GKEControlPlaneCreatedReason, clusterv1.ConditionSeverityInfo, "")
	s.scope.GCPManagedControlPlane.Status.Ready = true
	s.scope.GCPManagedControlPlane.Status.Initialized = true

	log.Info("Cluster reconciled")

	return ctrl.Result{}, nil
}

// Delete delete GKE cluster.
func (s *Service) Delete(ctx context.Context) (ctrl.Result, error) {
	log := log.FromContext(ctx).WithValues("service", "container.clusters")
	log.Info("Deleting cluster resources")

	cluster, err := s.describeCluster(ctx, &log)
	if err != nil {
		return ctrl.Result{}, err
	}
	if cluster == nil {
		log.Info("Cluster already deleted")
		conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEControlPlaneDeletingCondition, infrav1exp.GKEControlPlaneDeletedReason, clusterv1.ConditionSeverityInfo, "")
		return ctrl.Result{}, nil
	}

	switch cluster.GetStatus() {
	case containerpb.Cluster_PROVISIONING:
		log.Info("Cluster provisioning in progress")
		return ctrl.Result{}, nil
	case containerpb.Cluster_RECONCILING:
		log.Info("Cluster reconciling in progress")
		return ctrl.Result{}, nil
	case containerpb.Cluster_STOPPING:
		log.Info("Cluster stopping in progress")
		conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEControlPlaneReadyCondition, infrav1exp.GKEControlPlaneDeletingReason, clusterv1.ConditionSeverityInfo, "")
		conditions.MarkTrue(s.scope.ConditionSetter(), infrav1exp.GKEControlPlaneDeletingCondition)
		return ctrl.Result{}, nil
	default:
		break
	}

	if err = s.deleteCluster(ctx, &log); err != nil {
		conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEControlPlaneDeletingCondition, infrav1exp.GKEControlPlaneReconciliationFailedReason, clusterv1.ConditionSeverityError, err.Error())
		return ctrl.Result{}, err
	}
	log.Info("Cluster deleting in progress")
	s.scope.GCPManagedControlPlane.Status.Initialized = false
	s.scope.GCPManagedControlPlane.Status.Ready = false
	conditions.MarkFalse(s.scope.ConditionSetter(), clusterv1.ReadyCondition, infrav1exp.GKEControlPlaneDeletingReason, clusterv1.ConditionSeverityInfo, "")
	conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEControlPlaneReadyCondition, infrav1exp.GKEControlPlaneDeletingReason, clusterv1.ConditionSeverityInfo, "")
	conditions.MarkTrue(s.scope.ConditionSetter(), infrav1exp.GKEControlPlaneDeletingCondition)

	return ctrl.Result{}, nil
}

func (s *Service) describeCluster(ctx context.Context, log *logr.Logger) (*containerpb.Cluster, error) {
	getClusterRequest := &containerpb.GetClusterRequest{
		Name: s.scope.ClusterFullName(),
	}
	cluster, err := s.scope.ManagedControlPlaneClient().GetCluster(ctx, getClusterRequest)
	if err != nil {
		var e *apierror.APIError
		if ok := errors.As(err, &e); ok {
			if e.GRPCStatus().Code() == codes.NotFound {
				return nil, nil
			}
		}
		log.Error(err, "Error getting GKE cluster", "name", s.scope.ClusterName())
		return nil, err
	}

	return cluster, nil
}

func (s *Service) createCluster(ctx context.Context, log *logr.Logger) error {
	nodePools, machinePools, _ := s.scope.GetAllNodePools(ctx)

	log.V(2).Info("Running pre-flight checks on machine pools before cluster creation")
	if err := shared.ManagedMachinePoolsPreflightCheck(nodePools, machinePools, s.scope.Region()); err != nil {
		return fmt.Errorf("preflight checks on machine pools before cluster create: %w", err)
	}

	isRegional := shared.IsRegional(s.scope.Region())
	cluster := &containerpb.Cluster{
		Name:       s.scope.ClusterName(),
		Network:    *s.scope.GCPManagedCluster.Spec.Network.Name,
		Subnetwork: s.getSubnetNameInClusterRegion(),
		Autopilot: &containerpb.Autopilot{
			Enabled: s.scope.GCPManagedControlPlane.Spec.EnableAutopilot,
		},
		ReleaseChannel: &containerpb.ReleaseChannel{
			Channel: convertToSdkReleaseChannel(s.scope.GCPManagedControlPlane.Spec.ReleaseChannel),
		},
		MasterAuthorizedNetworksConfig: convertToSdkMasterAuthorizedNetworksConfig(s.scope.GCPManagedControlPlane.Spec.MasterAuthorizedNetworksConfig),
	}
	if s.scope.GCPManagedControlPlane.Spec.ControlPlaneVersion != nil {
		cluster.InitialClusterVersion = convertToSdkMasterVersion(*s.scope.GCPManagedControlPlane.Spec.ControlPlaneVersion)
	}
	if !s.scope.IsAutopilotCluster() {
		cluster.NodePools = scope.ConvertToSdkNodePools(nodePools, machinePools, isRegional, cluster.GetName())
	}

	createClusterRequest := &containerpb.CreateClusterRequest{
		Cluster: cluster,
		Parent:  s.scope.ClusterLocation(),
	}

	log.V(2).Info("Creating GKE cluster")
	_, err := s.scope.ManagedControlPlaneClient().CreateCluster(ctx, createClusterRequest)
	if err != nil {
		log.Error(err, "Error creating GKE cluster", "name", s.scope.ClusterName())
		return err
	}

	err = shared.ResourceTagBinding(
		ctx,
		s.scope.TagBindingsClient(),
		s.scope.GCPManagedCluster.Spec,
		s.scope.ClusterName(),
	)
	if err != nil {
		log.Error(err, "Error binding tags to cluster resources", "name", s.scope.ClusterName())
		return err
	}

	return nil
}

// getSubnetNameInClusterRegion returns the subnet which is in the same region as cluster. If not found it returns empty string.
func (s *Service) getSubnetNameInClusterRegion() string {
	for _, subnet := range s.scope.GCPManagedCluster.Spec.Network.Subnets {
		if subnet.Region == s.scope.Region() {
			return subnet.Name
		}
	}
	return ""
}

func (s *Service) updateCluster(ctx context.Context, updateClusterRequest *containerpb.UpdateClusterRequest, log *logr.Logger) error {
	_, err := s.scope.ManagedControlPlaneClient().UpdateCluster(ctx, updateClusterRequest)
	if err != nil {
		log.Error(err, "Error updating GKE cluster", "name", s.scope.ClusterName())
		return err
	}

	return nil
}

func (s *Service) deleteCluster(ctx context.Context, log *logr.Logger) error {
	deleteClusterRequest := &containerpb.DeleteClusterRequest{
		Name: s.scope.ClusterFullName(),
	}
	_, err := s.scope.ManagedControlPlaneClient().DeleteCluster(ctx, deleteClusterRequest)
	if err != nil {
		log.Error(err, "Error deleting GKE cluster", "name", s.scope.ClusterName())
		return err
	}

	return nil
}

func convertToSdkReleaseChannel(channel *infrav1exp.ReleaseChannel) containerpb.ReleaseChannel_Channel {
	if channel == nil {
		return containerpb.ReleaseChannel_UNSPECIFIED
	}
	switch *channel {
	case infrav1exp.Rapid:
		return containerpb.ReleaseChannel_RAPID
	case infrav1exp.Regular:
		return containerpb.ReleaseChannel_REGULAR
	case infrav1exp.Stable:
		return containerpb.ReleaseChannel_STABLE
	default:
		return containerpb.ReleaseChannel_UNSPECIFIED
	}
}

func convertToSdkMasterVersion(masterVersion string) string {
	// For example, the master version returned from GCP SDK can be 1.27.2-gke.2100, we want to convert it to 1.27.2
	return strings.Replace(strings.Split(masterVersion, "-")[0], "v", "", 1)
}

// convertToSdkMasterAuthorizedNetworksConfig converts the MasterAuthorizedNetworksConfig defined in CRs to the SDK version.
func convertToSdkMasterAuthorizedNetworksConfig(config *infrav1exp.MasterAuthorizedNetworksConfig) *containerpb.MasterAuthorizedNetworksConfig {
	// if config is nil, it means that the user wants to disable the feature.
	if config == nil {
		return &containerpb.MasterAuthorizedNetworksConfig{
			Enabled:                     false,
			CidrBlocks:                  []*containerpb.MasterAuthorizedNetworksConfig_CidrBlock{},
			GcpPublicCidrsAccessEnabled: new(bool),
		}
	}

	// Convert the CidrBlocks slice.
	cidrBlocks := make([]*containerpb.MasterAuthorizedNetworksConfig_CidrBlock, len(config.CidrBlocks))
	for i, cidrBlock := range config.CidrBlocks {
		cidrBlocks[i] = &containerpb.MasterAuthorizedNetworksConfig_CidrBlock{
			CidrBlock:   cidrBlock.CidrBlock,
			DisplayName: cidrBlock.DisplayName,
		}
	}

	return &containerpb.MasterAuthorizedNetworksConfig{
		Enabled:                     true,
		CidrBlocks:                  cidrBlocks,
		GcpPublicCidrsAccessEnabled: config.GcpPublicCidrsAccessEnabled,
	}
}

func (s *Service) checkDiffAndPrepareUpdate(existingCluster *containerpb.Cluster, log *logr.Logger) (bool, *containerpb.UpdateClusterRequest) {
	log.V(4).Info("Checking diff and preparing update.")

	needUpdate := false
	clusterUpdate := containerpb.ClusterUpdate{}
	// Release channel
	desiredReleaseChannel := convertToSdkReleaseChannel(s.scope.GCPManagedControlPlane.Spec.ReleaseChannel)
	if desiredReleaseChannel != existingCluster.GetReleaseChannel().GetChannel() {
		log.V(2).Info("Release channel update required", "current", existingCluster.GetReleaseChannel().GetChannel(), "desired", desiredReleaseChannel)
		needUpdate = true
		clusterUpdate.DesiredReleaseChannel = &containerpb.ReleaseChannel{
			Channel: desiredReleaseChannel,
		}
	}
	// Master version
	if s.scope.GCPManagedControlPlane.Spec.ControlPlaneVersion != nil {
		desiredMasterVersion := convertToSdkMasterVersion(*s.scope.GCPManagedControlPlane.Spec.ControlPlaneVersion)
		existingClusterMasterVersion := convertToSdkMasterVersion(existingCluster.GetCurrentMasterVersion())
		if desiredMasterVersion != existingClusterMasterVersion {
			needUpdate = true
			clusterUpdate.DesiredMasterVersion = desiredMasterVersion
			log.V(2).Info("Master version update required", "current", existingClusterMasterVersion, "desired", desiredMasterVersion)
		}
	}

	// DesiredMasterAuthorizedNetworksConfig
	// When desiredMasterAuthorizedNetworksConfig is nil, it means that the user wants to disable the feature.
	desiredMasterAuthorizedNetworksConfig := convertToSdkMasterAuthorizedNetworksConfig(s.scope.GCPManagedControlPlane.Spec.MasterAuthorizedNetworksConfig)
	if !compareMasterAuthorizedNetworksConfig(desiredMasterAuthorizedNetworksConfig, existingCluster.GetMasterAuthorizedNetworksConfig()) {
		needUpdate = true
		clusterUpdate.DesiredMasterAuthorizedNetworksConfig = desiredMasterAuthorizedNetworksConfig
		log.V(2).Info("Master authorized networks config update required", "current", existingCluster.GetMasterAuthorizedNetworksConfig(), "desired", desiredMasterAuthorizedNetworksConfig)
	}
	log.V(4).Info("Master authorized networks config update check", "current", existingCluster.GetMasterAuthorizedNetworksConfig())
	if desiredMasterAuthorizedNetworksConfig != nil {
		log.V(4).Info("Master authorized networks config update check", "desired", desiredMasterAuthorizedNetworksConfig)
	}

	updateClusterRequest := containerpb.UpdateClusterRequest{
		Name:   s.scope.ClusterFullName(),
		Update: &clusterUpdate,
	}
	log.V(4).Info("Update cluster request. ", "needUpdate", needUpdate, "updateClusterRequest", &updateClusterRequest)
	return needUpdate, &updateClusterRequest
}

// compare if two MasterAuthorizedNetworksConfig are equal.
func compareMasterAuthorizedNetworksConfig(a, b *containerpb.MasterAuthorizedNetworksConfig) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	if a.GetEnabled() != b.GetEnabled() {
		return false
	}
	if (a.GcpPublicCidrsAccessEnabled == nil && b.GcpPublicCidrsAccessEnabled != nil) || (a.GcpPublicCidrsAccessEnabled != nil && b.GcpPublicCidrsAccessEnabled == nil) {
		return false
	}
	if a.GcpPublicCidrsAccessEnabled != nil && b.GcpPublicCidrsAccessEnabled != nil && a.GetGcpPublicCidrsAccessEnabled() != b.GetGcpPublicCidrsAccessEnabled() {
		return false
	}
	// if one cidrBlocks is nil, but the other is empty, they are equal.
	if (a.CidrBlocks == nil && b.CidrBlocks != nil && len(b.GetCidrBlocks()) == 0) || (b.CidrBlocks == nil && a.CidrBlocks != nil && len(a.GetCidrBlocks()) == 0) {
		return true
	}
	if !cmp.Equal(a.GetCidrBlocks(), b.GetCidrBlocks()) {
		return false
	}
	return true
}
