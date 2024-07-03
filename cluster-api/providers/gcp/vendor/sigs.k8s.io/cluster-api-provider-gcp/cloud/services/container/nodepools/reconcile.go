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

package nodepools

import (
	"context"
	"fmt"
	"strings"

	"sigs.k8s.io/cluster-api-provider-gcp/util/resourceurl"

	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"

	"cloud.google.com/go/compute/apiv1/computepb"
	"cloud.google.com/go/container/apiv1/containerpb"
	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/googleapis/gax-go/v2/apierror"
	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud/providerid"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud/services/shared"
	infrav1exp "sigs.k8s.io/cluster-api-provider-gcp/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-gcp/util/reconciler"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// setReadyStatusFromConditions updates the GCPManagedMachinePool's ready status based on its conditions.
func (s *Service) setReadyStatusFromConditions() {
	machinePool := s.scope.GCPManagedMachinePool
	if conditions.IsTrue(machinePool, clusterv1.ReadyCondition) || conditions.IsTrue(machinePool, infrav1exp.GKEMachinePoolUpdatingCondition) {
		s.scope.GCPManagedMachinePool.Status.Ready = true
		return
	}

	s.scope.GCPManagedMachinePool.Status.Ready = false
}

// Reconcile reconcile GKE node pool.
func (s *Service) Reconcile(ctx context.Context) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Reconciling node pool resources")

	// Update GCPManagedMachinePool ready status based on conditions
	defer s.setReadyStatusFromConditions()

	nodePool, err := s.describeNodePool(ctx, &log)
	if err != nil {
		conditions.MarkFalse(s.scope.ConditionSetter(), clusterv1.ReadyCondition, infrav1exp.GKEMachinePoolReconciliationFailedReason, clusterv1.ConditionSeverityError, err.Error())
		return ctrl.Result{}, err
	}
	if nodePool == nil {
		log.Info("Node pool not found, creating", "cluster", s.scope.Cluster.Name)
		if err = s.createNodePool(ctx, &log); err != nil {
			conditions.MarkFalse(s.scope.ConditionSetter(), clusterv1.ReadyCondition, infrav1exp.GKEMachinePoolReconciliationFailedReason, clusterv1.ConditionSeverityError, err.Error())
			conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolReadyCondition, infrav1exp.GKEMachinePoolReconciliationFailedReason, clusterv1.ConditionSeverityError, err.Error())
			conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolCreatingCondition, infrav1exp.GKEMachinePoolReconciliationFailedReason, clusterv1.ConditionSeverityError, err.Error())
			return ctrl.Result{}, err
		}
		log.Info("Node pool provisioning in progress")
		conditions.MarkFalse(s.scope.ConditionSetter(), clusterv1.ReadyCondition, infrav1exp.GKEMachinePoolCreatingReason, clusterv1.ConditionSeverityInfo, "")
		conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolReadyCondition, infrav1exp.GKEMachinePoolCreatingReason, clusterv1.ConditionSeverityInfo, "")
		conditions.MarkTrue(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolCreatingCondition)
		return ctrl.Result{RequeueAfter: reconciler.DefaultRetryTime}, nil
	}
	log.V(2).Info("Node pool found", "cluster", s.scope.Cluster.Name, "nodepool", nodePool.GetName())

	instances, err := s.getInstances(ctx, nodePool)
	if err != nil {
		conditions.MarkFalse(s.scope.ConditionSetter(), clusterv1.ReadyCondition, infrav1exp.GKEMachinePoolReconciliationFailedReason, clusterv1.ConditionSeverityError, err.Error())
		return ctrl.Result{}, err
	}
	providerIDList := []string{}
	for _, instance := range instances {
		log.V(4).Info("parsing gce instance url", "url", instance.GetInstance())
		providerID, err := providerid.NewFromResourceURL(instance.GetInstance())
		if err != nil {
			log.Error(err, "parsing instance url", "url", instance.GetInstance())
			conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolReadyCondition, infrav1exp.GKEMachinePoolErrorReason, clusterv1.ConditionSeverityError, "")
			return ctrl.Result{}, err
		}
		providerIDList = append(providerIDList, providerID.String())
	}
	s.scope.GCPManagedMachinePool.Spec.ProviderIDList = providerIDList
	s.scope.GCPManagedMachinePool.Status.Replicas = int32(len(providerIDList))

	// Update GKEManagedMachinePool conditions based on GKE node pool status
	switch nodePool.GetStatus() {
	case containerpb.NodePool_PROVISIONING:
		// node pool is creating
		log.Info("Node pool provisioning in progress")
		conditions.MarkFalse(s.scope.ConditionSetter(), clusterv1.ReadyCondition, infrav1exp.GKEMachinePoolCreatingReason, clusterv1.ConditionSeverityInfo, "")
		conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolReadyCondition, infrav1exp.GKEMachinePoolCreatingReason, clusterv1.ConditionSeverityInfo, "")
		conditions.MarkTrue(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolCreatingCondition)
		return ctrl.Result{RequeueAfter: reconciler.DefaultRetryTime}, nil
	case containerpb.NodePool_RECONCILING:
		// node pool is updating/reconciling
		log.Info("Node pool reconciling in progress")
		conditions.MarkTrue(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolUpdatingCondition)
		return ctrl.Result{RequeueAfter: reconciler.DefaultRetryTime}, nil
	case containerpb.NodePool_STOPPING:
		// node pool is deleting
		log.Info("Node pool stopping in progress")
		conditions.MarkFalse(s.scope.ConditionSetter(), clusterv1.ReadyCondition, infrav1exp.GKEMachinePoolDeletingReason, clusterv1.ConditionSeverityInfo, "")
		conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolReadyCondition, infrav1exp.GKEMachinePoolDeletingReason, clusterv1.ConditionSeverityInfo, "")
		conditions.MarkTrue(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolDeletingCondition)
		return ctrl.Result{}, nil
	case containerpb.NodePool_ERROR, containerpb.NodePool_RUNNING_WITH_ERROR:
		// node pool is in error or degraded state
		var msg string
		if len(nodePool.GetConditions()) > 0 {
			msg = nodePool.GetConditions()[0].GetMessage()
		}
		log.Error(errors.New("Node pool in error/degraded state"), msg, "name", s.scope.GCPManagedMachinePool.Name)
		conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolReadyCondition, infrav1exp.GKEMachinePoolErrorReason, clusterv1.ConditionSeverityError, "")
		return ctrl.Result{}, nil
	case containerpb.NodePool_RUNNING:
		// node pool is ready and running
		conditions.MarkTrue(s.scope.ConditionSetter(), clusterv1.ReadyCondition)
		conditions.MarkTrue(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolReadyCondition)
		conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolCreatingCondition, infrav1exp.GKEMachinePoolCreatedReason, clusterv1.ConditionSeverityInfo, "")
		log.Info("Node pool running")
	default:
		log.Error(errors.New("Unhandled node pool status"), fmt.Sprintf("Unhandled node pool status %s", nodePool.GetStatus()), "name", s.scope.GCPManagedMachinePool.Name)
		return ctrl.Result{}, nil
	}

	needUpdateConfig, nodePoolUpdateConfigRequest := s.checkDiffAndPrepareUpdateConfig(nodePool)
	if needUpdateConfig {
		log.Info("Node pool config update required", "request", nodePoolUpdateConfigRequest)
		err = s.updateNodePoolConfig(ctx, nodePoolUpdateConfigRequest)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("node pool config update (either version/labels/taints/locations/image type/network tag/linux node config or all) failed: %s", err)
		}
		log.Info("Node pool config updating in progress")
		s.scope.GCPManagedMachinePool.Status.Ready = true
		conditions.MarkTrue(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolUpdatingCondition)
		return ctrl.Result{RequeueAfter: reconciler.DefaultRetryTime}, nil
	}

	needUpdateAutoscaling, setNodePoolAutoscalingRequest := s.checkDiffAndPrepareUpdateAutoscaling(nodePool)
	if needUpdateAutoscaling {
		log.Info("Auto scaling update required")
		err = s.updateNodePoolAutoscaling(ctx, setNodePoolAutoscalingRequest)
		if err != nil {
			return ctrl.Result{}, err
		}
		log.Info("Node pool auto scaling updating in progress")
		conditions.MarkTrue(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolUpdatingCondition)
		return ctrl.Result{RequeueAfter: reconciler.DefaultRetryTime}, nil
	}

	needUpdateSize, setNodePoolSizeRequest := s.checkDiffAndPrepareUpdateSize(nodePool)
	if needUpdateSize {
		log.Info("Size update required")
		err = s.updateNodePoolSize(ctx, setNodePoolSizeRequest)
		if err != nil {
			return ctrl.Result{}, err
		}
		log.Info("Node pool size updating in progress")
		conditions.MarkTrue(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolUpdatingCondition)
		return ctrl.Result{RequeueAfter: reconciler.DefaultRetryTime}, nil
	}

	conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolUpdatingCondition, infrav1exp.GKEMachinePoolUpdatedReason, clusterv1.ConditionSeverityInfo, "")

	s.scope.SetReplicas(int32(len(s.scope.GCPManagedMachinePool.Spec.ProviderIDList)))
	log.Info("Node pool reconciled")
	s.scope.GCPManagedMachinePool.Status.Ready = true
	conditions.MarkTrue(s.scope.ConditionSetter(), clusterv1.ReadyCondition)
	conditions.MarkTrue(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolReadyCondition)
	conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolCreatingCondition, infrav1exp.GKEMachinePoolCreatedReason, clusterv1.ConditionSeverityInfo, "")

	return ctrl.Result{}, nil
}

// Delete delete GKE node pool.
func (s *Service) Delete(ctx context.Context) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Deleting node pool resources")

	defer s.setReadyStatusFromConditions()

	nodePool, err := s.describeNodePool(ctx, &log)
	if err != nil {
		return ctrl.Result{}, err
	}
	if nodePool == nil {
		log.Info("Node pool already deleted")
		conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolDeletingCondition, infrav1exp.GKEMachinePoolDeletedReason, clusterv1.ConditionSeverityInfo, "")
		return ctrl.Result{}, err
	}

	switch nodePool.GetStatus() {
	case containerpb.NodePool_PROVISIONING:
		log.Info("Node pool provisioning in progress")
		return ctrl.Result{RequeueAfter: reconciler.DefaultRetryTime}, nil
	case containerpb.NodePool_RECONCILING:
		log.Info("Node pool reconciling in progress")
		return ctrl.Result{RequeueAfter: reconciler.DefaultRetryTime}, nil
	case containerpb.NodePool_STOPPING:
		log.Info("Node pool stopping in progress")
		conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolReadyCondition, infrav1exp.GKEMachinePoolDeletingReason, clusterv1.ConditionSeverityInfo, "")
		conditions.MarkTrue(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolDeletingCondition)
		return ctrl.Result{RequeueAfter: reconciler.DefaultRetryTime}, nil
	default:
		break
	}

	if err = s.deleteNodePool(ctx); err != nil {
		conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolDeletingCondition, infrav1exp.GKEMachinePoolReconciliationFailedReason, clusterv1.ConditionSeverityError, err.Error())
		return ctrl.Result{}, err
	}
	log.Info("Node pool deleting in progress")
	conditions.MarkFalse(s.scope.ConditionSetter(), clusterv1.ReadyCondition, infrav1exp.GKEMachinePoolDeletingReason, clusterv1.ConditionSeverityInfo, "")
	conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolReadyCondition, infrav1exp.GKEMachinePoolDeletingReason, clusterv1.ConditionSeverityInfo, "")
	conditions.MarkTrue(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolDeletingCondition)

	return ctrl.Result{}, nil
}

func (s *Service) describeNodePool(ctx context.Context, log *logr.Logger) (*containerpb.NodePool, error) {
	getNodePoolRequest := &containerpb.GetNodePoolRequest{
		Name: s.scope.NodePoolFullName(),
	}
	nodePool, err := s.scope.ManagedMachinePoolClient().GetNodePool(ctx, getNodePoolRequest)
	if err != nil {
		var e *apierror.APIError
		if ok := errors.As(err, &e); ok {
			if e.GRPCStatus().Code() == codes.NotFound {
				return nil, nil
			}
		}
		log.Error(err, "Error getting GKE node pool", "name", s.scope.GCPManagedMachinePool.Name)
		return nil, err
	}

	return nodePool, nil
}

func (s *Service) getInstances(ctx context.Context, nodePool *containerpb.NodePool) ([]*computepb.ManagedInstance, error) {
	instances := []*computepb.ManagedInstance{}

	for _, url := range nodePool.GetInstanceGroupUrls() {
		resourceURL, err := resourceurl.Parse(url)
		if err != nil {
			return nil, errors.Wrap(err, "error parsing instance group url")
		}
		listManagedInstancesRequest := &computepb.ListManagedInstancesInstanceGroupManagersRequest{
			InstanceGroupManager: resourceURL.Name,
			Project:              resourceURL.Project,
			Zone:                 resourceURL.Location,
		}
		iter := s.scope.InstanceGroupManagersClient().ListManagedInstances(ctx, listManagedInstancesRequest)
		for {
			resp, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, err
			}
			instances = append(instances, resp)
		}
	}

	return instances, nil
}

func (s *Service) createNodePool(ctx context.Context, log *logr.Logger) error {
	log.V(2).Info("Running pre-flight checks on machine pool before creation")
	if err := shared.ManagedMachinePoolPreflightCheck(s.scope.GCPManagedMachinePool, s.scope.MachinePool, s.scope.Region()); err != nil {
		return fmt.Errorf("preflight checks on machine pool before creating: %w", err)
	}

	isRegional := shared.IsRegional(s.scope.Region())

	createNodePoolRequest := &containerpb.CreateNodePoolRequest{
		NodePool: scope.ConvertToSdkNodePool(*s.scope.GCPManagedMachinePool, *s.scope.MachinePool, isRegional, s.scope.GCPManagedControlPlane.Spec.ClusterName),
		Parent:   s.scope.NodePoolLocation(),
	}
	_, err := s.scope.ManagedMachinePoolClient().CreateNodePool(ctx, createNodePoolRequest)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) updateNodePoolConfig(ctx context.Context, updateNodePoolRequest *containerpb.UpdateNodePoolRequest) error {
	_, err := s.scope.ManagedMachinePoolClient().UpdateNodePool(ctx, updateNodePoolRequest)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) updateNodePoolAutoscaling(ctx context.Context, setNodePoolAutoscalingRequest *containerpb.SetNodePoolAutoscalingRequest) error {
	_, err := s.scope.ManagedMachinePoolClient().SetNodePoolAutoscaling(ctx, setNodePoolAutoscalingRequest)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) updateNodePoolSize(ctx context.Context, setNodePoolSizeRequest *containerpb.SetNodePoolSizeRequest) error {
	_, err := s.scope.ManagedMachinePoolClient().SetNodePoolSize(ctx, setNodePoolSizeRequest)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) deleteNodePool(ctx context.Context) error {
	deleteNodePoolRequest := &containerpb.DeleteNodePoolRequest{
		Name: s.scope.NodePoolFullName(),
	}
	_, err := s.scope.ManagedMachinePoolClient().DeleteNodePool(ctx, deleteNodePoolRequest)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) checkDiffAndPrepareUpdateConfig(existingNodePool *containerpb.NodePool) (bool, *containerpb.UpdateNodePoolRequest) {
	needUpdate := false
	updateNodePoolRequest := containerpb.UpdateNodePoolRequest{
		Name: s.scope.NodePoolFullName(),
	}

	isRegional := shared.IsRegional(s.scope.Region())
	desiredNodePool := scope.ConvertToSdkNodePool(*s.scope.GCPManagedMachinePool, *s.scope.MachinePool, isRegional, s.scope.GCPManagedControlPlane.Spec.ClusterName)

	// Node version
	if s.scope.NodePoolVersion() != nil {
		desiredNodePoolVersion := infrav1exp.ConvertFromSdkNodeVersion(*s.scope.NodePoolVersion())
		if desiredNodePoolVersion != infrav1exp.ConvertFromSdkNodeVersion(existingNodePool.GetVersion()) {
			needUpdate = true
			updateNodePoolRequest.NodeVersion = desiredNodePoolVersion
		}
	}
	// Kubernetes labels
	if !cmp.Equal(desiredNodePool.GetConfig().GetLabels(), existingNodePool.GetConfig().GetLabels()) {
		needUpdate = true
		updateNodePoolRequest.Labels = &containerpb.NodeLabels{
			Labels: desiredNodePool.GetConfig().GetLabels(),
		}
	}
	// Kubernetes taints
	if !cmp.Equal(desiredNodePool.GetConfig().GetTaints(), existingNodePool.GetConfig().GetTaints()) {
		needUpdate = true
		updateNodePoolRequest.Taints = &containerpb.NodeTaints{
			Taints: desiredNodePool.GetConfig().GetTaints(),
		}
	}
	// Node image type
	// GCP API returns image type string in all uppercase, we can do a case-insensitive check here.
	if desiredNodePool.GetConfig().GetImageType() != "" && !strings.EqualFold(desiredNodePool.GetConfig().GetImageType(), existingNodePool.GetConfig().GetImageType()) {
		needUpdate = true
		updateNodePoolRequest.ImageType = desiredNodePool.GetConfig().GetImageType()
	}
	// Additional resource labels
	if !cmp.Equal(desiredNodePool.GetConfig().GetResourceLabels(), existingNodePool.GetConfig().GetResourceLabels()) {
		needUpdate = true
		updateNodePoolRequest.ResourceLabels = &containerpb.ResourceLabels{
			Labels: desiredNodePool.GetConfig().GetResourceLabels(),
		}
	}
	// Locations
	desiredLocations := s.scope.GCPManagedMachinePool.Spec.NodeLocations
	if desiredLocations != nil && !cmp.Equal(desiredLocations, existingNodePool.GetLocations()) {
		needUpdate = true
		updateNodePoolRequest.Locations = desiredLocations
	}
	// Network tags
	desiredNetworkTags := s.scope.GCPManagedMachinePool.Spec.NodeNetwork.Tags
	if existingNodePool.GetConfig() != nil && !cmp.Equal(desiredNetworkTags, existingNodePool.GetConfig().GetTags()) {
		needUpdate = true
		updateNodePoolRequest.Tags = &containerpb.NetworkTags{
			Tags: desiredNetworkTags,
		}
	}
	// LinuxNodeConfig
	desiredLinuxNodeConfig := infrav1exp.ConvertToSdkLinuxNodeConfig(s.scope.GCPManagedMachinePool.Spec.LinuxNodeConfig)
	if !cmp.Equal(desiredLinuxNodeConfig, existingNodePool.GetConfig().GetLinuxNodeConfig(), cmpopts.IgnoreUnexported(containerpb.LinuxNodeConfig{})) {
		needUpdate = true
		updateNodePoolRequest.LinuxNodeConfig = desiredLinuxNodeConfig
	}

	return needUpdate, &updateNodePoolRequest
}

func (s *Service) checkDiffAndPrepareUpdateAutoscaling(existingNodePool *containerpb.NodePool) (bool, *containerpb.SetNodePoolAutoscalingRequest) {
	needUpdate := false
	desiredAutoscaling := infrav1exp.ConvertToSdkAutoscaling(s.scope.GCPManagedMachinePool.Spec.Scaling)

	setNodePoolAutoscalingRequest := containerpb.SetNodePoolAutoscalingRequest{
		Name: s.scope.NodePoolFullName(),
	}

	if !cmp.Equal(desiredAutoscaling, existingNodePool.GetAutoscaling(), cmpopts.IgnoreUnexported(containerpb.NodePoolAutoscaling{})) {
		needUpdate = true
		setNodePoolAutoscalingRequest.Autoscaling = desiredAutoscaling
	}
	return needUpdate, &setNodePoolAutoscalingRequest
}

func (s *Service) checkDiffAndPrepareUpdateSize(existingNodePool *containerpb.NodePool) (bool, *containerpb.SetNodePoolSizeRequest) {
	needUpdate := false
	desiredAutoscaling := infrav1exp.ConvertToSdkAutoscaling(s.scope.GCPManagedMachinePool.Spec.Scaling)

	if desiredAutoscaling.GetEnabled() {
		// Do not update node pool size if autoscaling is enabled.
		return false, nil
	}

	setNodePoolSizeRequest := containerpb.SetNodePoolSizeRequest{
		Name: s.scope.NodePoolFullName(),
	}

	replicas := *s.scope.MachinePool.Spec.Replicas
	if shared.IsRegional(s.scope.Region()) {
		replicas /= int32(len(existingNodePool.GetLocations()))
	}

	if replicas != existingNodePool.GetInitialNodeCount() {
		needUpdate = true
		setNodePoolSizeRequest.NodeCount = replicas
	}
	return needUpdate, &setNodePoolSizeRequest
}
