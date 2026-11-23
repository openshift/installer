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

	"cloud.google.com/go/compute/apiv1/computepb"
	"cloud.google.com/go/container/apiv1/containerpb"
	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud/providerid"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud/services/shared"
	infrav1exp "sigs.k8s.io/cluster-api-provider-gcp/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-gcp/util/reconciler"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// setReadyStatusFromConditions updates the GCPManagedMachinePool's ready status based on its conditions.
func (s *Service) setReadyStatusFromConditions() {
	machinePool := s.scope.GCPManagedMachinePool
	if v1beta1conditions.IsTrue(machinePool, clusterv1beta1.ReadyCondition) || v1beta1conditions.IsTrue(machinePool, infrav1exp.GKEMachinePoolUpdatingCondition) {
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
		v1beta1conditions.MarkFalse(s.scope.ConditionSetter(), clusterv1beta1.ReadyCondition, infrav1exp.GKEMachinePoolReconciliationFailedReason, clusterv1beta1.ConditionSeverityError, "reading node pool: %v", err)
		return ctrl.Result{}, err
	}
	if nodePool == nil {
		log.Info("Node pool not found, creating", "cluster", s.scope.Cluster.Name)
		if err = s.createNodePool(ctx, &log); err != nil {
			v1beta1conditions.MarkFalse(s.scope.ConditionSetter(), clusterv1beta1.ReadyCondition, infrav1exp.GKEMachinePoolReconciliationFailedReason, clusterv1beta1.ConditionSeverityError, "creating node pool: %v", err)
			v1beta1conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolReadyCondition, infrav1exp.GKEMachinePoolReconciliationFailedReason, clusterv1beta1.ConditionSeverityError, "creating node pool: %v", err)
			v1beta1conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolCreatingCondition, infrav1exp.GKEMachinePoolReconciliationFailedReason, clusterv1beta1.ConditionSeverityError, "creating node pool: %v", err)
			return ctrl.Result{}, err
		}
		log.Info("Node pool provisioning in progress")
		v1beta1conditions.MarkFalse(s.scope.ConditionSetter(), clusterv1beta1.ReadyCondition, infrav1exp.GKEMachinePoolCreatingReason, clusterv1beta1.ConditionSeverityInfo, "")
		v1beta1conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolReadyCondition, infrav1exp.GKEMachinePoolCreatingReason, clusterv1beta1.ConditionSeverityInfo, "")
		v1beta1conditions.MarkTrue(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolCreatingCondition)
		return ctrl.Result{RequeueAfter: reconciler.DefaultRetryTime}, nil
	}
	log.V(2).Info("Node pool found", "cluster", s.scope.Cluster.Name, "nodepool", nodePool.GetName())

	instances, err := s.getInstances(ctx, nodePool)
	if err != nil {
		v1beta1conditions.MarkFalse(s.scope.ConditionSetter(), clusterv1beta1.ReadyCondition, infrav1exp.GKEMachinePoolReconciliationFailedReason, clusterv1beta1.ConditionSeverityError, "reading instances: %v", err)
		return ctrl.Result{}, err
	}
	providerIDList := []string{}
	for _, instance := range instances {
		log.V(4).Info("parsing gce instance url", "url", instance.GetInstance())
		providerID, err := providerid.NewFromResourceURL(instance.GetInstance())
		if err != nil {
			log.Error(err, "parsing instance url", "url", instance.GetInstance())
			v1beta1conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolReadyCondition, infrav1exp.GKEMachinePoolErrorReason, clusterv1beta1.ConditionSeverityError, "")
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
		v1beta1conditions.MarkFalse(s.scope.ConditionSetter(), clusterv1beta1.ReadyCondition, infrav1exp.GKEMachinePoolCreatingReason, clusterv1beta1.ConditionSeverityInfo, "")
		v1beta1conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolReadyCondition, infrav1exp.GKEMachinePoolCreatingReason, clusterv1beta1.ConditionSeverityInfo, "")
		v1beta1conditions.MarkTrue(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolCreatingCondition)
		return ctrl.Result{RequeueAfter: reconciler.DefaultRetryTime}, nil
	case containerpb.NodePool_RECONCILING:
		// node pool is updating/reconciling
		log.Info("Node pool reconciling in progress")
		v1beta1conditions.MarkTrue(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolUpdatingCondition)
		return ctrl.Result{RequeueAfter: reconciler.DefaultRetryTime}, nil
	case containerpb.NodePool_STOPPING:
		// node pool is deleting
		log.Info("Node pool stopping in progress")
		v1beta1conditions.MarkFalse(s.scope.ConditionSetter(), clusterv1beta1.ReadyCondition, infrav1exp.GKEMachinePoolDeletingReason, clusterv1beta1.ConditionSeverityInfo, "")
		v1beta1conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolReadyCondition, infrav1exp.GKEMachinePoolDeletingReason, clusterv1beta1.ConditionSeverityInfo, "")
		v1beta1conditions.MarkTrue(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolDeletingCondition)
		return ctrl.Result{}, nil
	case containerpb.NodePool_ERROR, containerpb.NodePool_RUNNING_WITH_ERROR:
		// node pool is in error or degraded state
		var msg string
		if len(nodePool.GetConditions()) > 0 {
			msg = nodePool.GetConditions()[0].GetMessage()
		}
		log.Error(errors.New("Node pool in error/degraded state"), msg, "name", s.scope.GCPManagedMachinePool.Name)
		v1beta1conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolReadyCondition, infrav1exp.GKEMachinePoolErrorReason, clusterv1beta1.ConditionSeverityError, "")
		return ctrl.Result{}, nil
	case containerpb.NodePool_RUNNING:
		// node pool is ready and running
		v1beta1conditions.MarkTrue(s.scope.ConditionSetter(), clusterv1beta1.ReadyCondition)
		v1beta1conditions.MarkTrue(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolReadyCondition)
		v1beta1conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolCreatingCondition, infrav1exp.GKEMachinePoolCreatedReason, clusterv1beta1.ConditionSeverityInfo, "")
		log.Info("Node pool running")
	default:
		log.Error(errors.New("Unhandled node pool status"), fmt.Sprintf("Unhandled node pool status %s", nodePool.GetStatus()), "name", s.scope.GCPManagedMachinePool.Name)
		return ctrl.Result{}, nil
	}

	// Fetch instance template labels for AdditionalLabels sync. GKE does not echo
	// NodeConfig.ResourceLabels in GetNodePool responses; the instance template's
	// properties.labels is the authoritative source of what labels GKE stamps onto
	// every node in the pool. On error we proceed without label sync for this cycle.
	templateLabels, err := s.fetchInstanceTemplateLabels(ctx, nodePool)
	if err != nil {
		log.Error(err, "Failed to read instance template labels, skipping AdditionalLabels sync this cycle")
		templateLabels = nil
	}

	// Run all checks before applying any update so that a persistent diff in one
	// category cannot permanently starve updates in another. GKE only allows one
	// concurrent operation per node pool, so we apply one per cycle (with requeue)
	// in priority order: size first (user-initiated scaling is most time-sensitive),
	// then autoscaling config, then general config.
	needUpdateSize, setNodePoolSizeRequest := s.checkDiffAndPrepareUpdateSize(nodePool)
	needUpdateAutoscaling, setNodePoolAutoscalingRequest := s.checkDiffAndPrepareUpdateAutoscaling(nodePool)
	needUpdateConfig, nodePoolUpdateConfigRequest := s.checkDiffAndPrepareUpdateConfig(nodePool, templateLabels)

	if needUpdateSize {
		log.Info("Size update required")
		err = s.updateNodePoolSize(ctx, setNodePoolSizeRequest)
		if err != nil {
			return ctrl.Result{}, err
		}
		log.Info("Node pool size updating in progress")
		v1beta1conditions.MarkTrue(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolUpdatingCondition)
		return ctrl.Result{RequeueAfter: reconciler.DefaultRetryTime}, nil
	}

	if needUpdateAutoscaling {
		log.Info("Auto scaling update required")
		err = s.updateNodePoolAutoscaling(ctx, setNodePoolAutoscalingRequest)
		if err != nil {
			return ctrl.Result{}, err
		}
		log.Info("Node pool auto scaling updating in progress")
		v1beta1conditions.MarkTrue(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolUpdatingCondition)
		return ctrl.Result{RequeueAfter: reconciler.DefaultRetryTime}, nil
	}

	if needUpdateConfig {
		log.Info("Node pool config update required", "request", nodePoolUpdateConfigRequest)
		err = s.updateNodePoolConfig(ctx, nodePoolUpdateConfigRequest)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("node pool config update (either version/labels/taints/locations/image type/network tag/linux node config or all) failed: %w", err)
		}
		log.Info("Node pool config updating in progress")
		s.scope.GCPManagedMachinePool.Status.Ready = true
		v1beta1conditions.MarkTrue(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolUpdatingCondition)
		return ctrl.Result{RequeueAfter: reconciler.DefaultRetryTime}, nil
	}

	v1beta1conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolUpdatingCondition, infrav1exp.GKEMachinePoolUpdatedReason, clusterv1beta1.ConditionSeverityInfo, "")

	s.scope.SetReplicas(int32(len(s.scope.GCPManagedMachinePool.Spec.ProviderIDList)))
	log.Info("Node pool reconciled")
	s.scope.GCPManagedMachinePool.Status.Ready = true
	v1beta1conditions.MarkTrue(s.scope.ConditionSetter(), clusterv1beta1.ReadyCondition)
	v1beta1conditions.MarkTrue(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolReadyCondition)
	v1beta1conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolCreatingCondition, infrav1exp.GKEMachinePoolCreatedReason, clusterv1beta1.ConditionSeverityInfo, "")

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
		v1beta1conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolDeletingCondition, infrav1exp.GKEMachinePoolDeletedReason, clusterv1beta1.ConditionSeverityInfo, "")
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
		v1beta1conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolReadyCondition, infrav1exp.GKEMachinePoolDeletingReason, clusterv1beta1.ConditionSeverityInfo, "")
		v1beta1conditions.MarkTrue(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolDeletingCondition)
		return ctrl.Result{RequeueAfter: reconciler.DefaultRetryTime}, nil
	default:
		break
	}

	if err = s.deleteNodePool(ctx); err != nil {
		v1beta1conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolDeletingCondition, infrav1exp.GKEMachinePoolReconciliationFailedReason, clusterv1beta1.ConditionSeverityError, "deleting node pool: %v", err)
		return ctrl.Result{}, err
	}
	log.Info("Node pool deleting in progress")
	v1beta1conditions.MarkFalse(s.scope.ConditionSetter(), clusterv1beta1.ReadyCondition, infrav1exp.GKEMachinePoolDeletingReason, clusterv1beta1.ConditionSeverityInfo, "")
	v1beta1conditions.MarkFalse(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolReadyCondition, infrav1exp.GKEMachinePoolDeletingReason, clusterv1beta1.ConditionSeverityInfo, "")
	v1beta1conditions.MarkTrue(s.scope.ConditionSetter(), infrav1exp.GKEMachinePoolDeletingCondition)

	return ctrl.Result{}, nil
}

func (s *Service) describeNodePool(ctx context.Context, log *logr.Logger) (*containerpb.NodePool, error) {
	getNodePoolRequest := &containerpb.GetNodePoolRequest{
		Name: s.scope.NodePoolFullName(),
	}
	nodePool, err := s.scope.ManagedMachinePoolClient().GetNodePool(ctx, getNodePoolRequest)
	if err != nil {
		if shared.IsNotFound(err) {
			return nil, nil
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

// fetchInstanceTemplateLabels returns the properties.labels from the instance template
// used by the first instance group in the node pool. These are the labels GKE stamps
// onto every node in the pool and are the authoritative source for AdditionalLabels sync.
// Returns nil if the node pool has no instance groups yet.
func (s *Service) fetchInstanceTemplateLabels(ctx context.Context, nodePool *containerpb.NodePool) (map[string]string, error) {
	urls := nodePool.GetInstanceGroupUrls()
	if len(urls) == 0 {
		return nil, nil
	}
	igURL, err := resourceurl.Parse(urls[0])
	if err != nil {
		return nil, errors.Wrap(err, "parsing instance group url")
	}
	mig, err := s.scope.InstanceGroupManagersClient().Get(ctx, &computepb.GetInstanceGroupManagerRequest{
		InstanceGroupManager: igURL.Name,
		Project:              igURL.Project,
		Zone:                 igURL.Location,
	})
	if err != nil {
		return nil, errors.Wrap(err, "getting instance group manager")
	}
	// Instance template URL format:
	// https://www.googleapis.com/compute/v1/projects/{project}/global/instanceTemplates/{name}
	templateURL := mig.GetInstanceTemplate()
	if templateURL == "" {
		return nil, nil
	}
	parts := strings.SplitN(strings.TrimPrefix(templateURL, "https://www.googleapis.com/compute/v1/projects/"), "/", 4)
	if len(parts) != 4 {
		return nil, errors.Errorf("unexpected instance template url format: %s", templateURL)
	}
	tmpl, err := s.scope.InstanceTemplatesClient().Get(ctx, &computepb.GetInstanceTemplateRequest{
		Project:          parts[0],
		InstanceTemplate: parts[3],
	})
	if err != nil {
		return nil, errors.Wrap(err, "getting instance template")
	}
	return tmpl.GetProperties().GetLabels(), nil
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
		if shared.IsNotFound(err) {
			return nil
		}
		return err
	}

	return nil
}

func (s *Service) checkDiffAndPrepareUpdateConfig(existingNodePool *containerpb.NodePool, existingTemplateLabels map[string]string) (bool, *containerpb.UpdateNodePoolRequest) {
	needUpdate := false
	updateNodePoolRequest := containerpb.UpdateNodePoolRequest{
		Name: s.scope.NodePoolFullName(),
	}

	isRegional := shared.IsRegional(s.scope.Region())
	desiredNodePool := scope.ConvertToSdkNodePool(*s.scope.GCPManagedMachinePool, *s.scope.MachinePool, isRegional, s.scope.GCPManagedControlPlane.Spec.ClusterName)

	// Node version
	if s.scope.NodePoolVersion() != "" {
		desiredNodePoolVersion := infrav1exp.ConvertFromSdkNodeVersion(s.scope.NodePoolVersion())
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
	if !cmp.Equal(desiredNodePool.GetConfig().GetTaints(), existingNodePool.GetConfig().GetTaints(), cmpopts.IgnoreUnexported(containerpb.NodeTaint{})) {
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
	// Resource labels (AdditionalLabels) — GKE does not echo NodeConfig.ResourceLabels in
	// GetNodePool responses, so we compare against the instance template's properties.labels
	// instead (see fetchInstanceTemplateLabels). We use a subset check: all desired labels
	// must be present with the correct value; GKE-internal labels on the template are ignored.
	if existingTemplateLabels != nil {
		desiredResourceLabels := scope.NodePoolResourceLabels(s.scope.GCPManagedMachinePool.Spec.AdditionalLabels, s.scope.GCPManagedControlPlane.Spec.ClusterName)
		for k, v := range desiredResourceLabels {
			if existingTemplateLabels[k] != v {
				needUpdate = true
				updateNodePoolRequest.ResourceLabels = &containerpb.ResourceLabels{
					Labels: desiredResourceLabels,
				}
				break
			}
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
	// LinuxNodeConfig — only compare when the user has set it. ConvertToSdkLinuxNodeConfig
	// returns a non-nil pointer even for nil input, which always differs from the nil
	// GKE returns for pools with no linux node config, producing a spurious update loop.
	if s.scope.GCPManagedMachinePool.Spec.LinuxNodeConfig != nil {
		desiredLinuxNodeConfig := infrav1exp.ConvertToSdkLinuxNodeConfig(s.scope.GCPManagedMachinePool.Spec.LinuxNodeConfig)
		if !cmp.Equal(desiredLinuxNodeConfig, existingNodePool.GetConfig().GetLinuxNodeConfig(), cmpopts.IgnoreUnexported(containerpb.LinuxNodeConfig{})) {
			needUpdate = true
			updateNodePoolRequest.LinuxNodeConfig = desiredLinuxNodeConfig
		}
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

	// Only skip size update if autoscaling is explicitly configured and enabled.
	// ConvertToSdkAutoscaling returns {Enabled: true} when passed nil, so without this
	// guard any node pool with no Spec.Scaling set would have manual size updates silently
	// suppressed.
	if s.scope.GCPManagedMachinePool.Spec.Scaling != nil {
		desiredAutoscaling := infrav1exp.ConvertToSdkAutoscaling(s.scope.GCPManagedMachinePool.Spec.Scaling)
		if desiredAutoscaling.GetEnabled() {
			return false, nil
		}
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
