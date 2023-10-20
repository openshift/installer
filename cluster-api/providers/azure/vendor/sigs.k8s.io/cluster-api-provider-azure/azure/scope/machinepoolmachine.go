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

package scope

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	kubedrain "k8s.io/kubectl/pkg/drain"
	"k8s.io/utils/ptr"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/converters"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/scalesetvms"
	infrav1exp "sigs.k8s.io/cluster-api-provider-azure/exp/api/v1beta1"
	azureutil "sigs.k8s.io/cluster-api-provider-azure/util/azure"
	"sigs.k8s.io/cluster-api-provider-azure/util/futures"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/controllers/noderefutil"
	"sigs.k8s.io/cluster-api/controllers/remote"
	capierrors "sigs.k8s.io/cluster-api/errors"
	expv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// MachinePoolMachineScopeName is the sourceName, or more specifically the UserAgent, of client used in cordon and drain.
	MachinePoolMachineScopeName = "azuremachinepoolmachine-scope"
)

type (
	nodeGetter interface {
		GetNodeByProviderID(ctx context.Context, providerID string) (*corev1.Node, error)
		GetNodeByObjectReference(ctx context.Context, nodeRef corev1.ObjectReference) (*corev1.Node, error)
	}

	workloadClusterProxy struct {
		Client  client.Client
		Cluster client.ObjectKey
	}

	// MachinePoolMachineScopeParams defines the input parameters used to create a new MachinePoolScope.
	MachinePoolMachineScopeParams struct {
		AzureMachinePool        *infrav1exp.AzureMachinePool
		AzureMachinePoolMachine *infrav1exp.AzureMachinePoolMachine
		Client                  client.Client
		ClusterScope            azure.ClusterScoper
		MachinePool             *expv1.MachinePool

		// workloadNodeGetter is only used for testing purposes and provides a way for mocking requests to the workload cluster
		workloadNodeGetter nodeGetter
	}

	// MachinePoolMachineScope defines a scope defined around a machine pool machine.
	MachinePoolMachineScope struct {
		azure.ClusterScoper
		AzureMachinePoolMachine *infrav1exp.AzureMachinePoolMachine
		AzureMachinePool        *infrav1exp.AzureMachinePool
		MachinePool             *expv1.MachinePool
		MachinePoolScope        *MachinePoolScope
		client                  client.Client
		patchHelper             *patch.Helper
		instance                *azure.VMSSVM

		// workloadNodeGetter is only used for testing purposes and provides a way for mocking requests to the workload cluster
		workloadNodeGetter nodeGetter
	}
)

// NewMachinePoolMachineScope creates a new MachinePoolMachineScope from the supplied parameters.
// This is meant to be called for each reconcile iteration.
func NewMachinePoolMachineScope(params MachinePoolMachineScopeParams) (*MachinePoolMachineScope, error) {
	if params.Client == nil {
		return nil, errors.New("client is required when creating a MachinePoolScope")
	}

	if params.ClusterScope == nil {
		return nil, errors.New("cluster scope is required when creating a MachinePoolScope")
	}

	if params.MachinePool == nil {
		return nil, errors.New("machine pool is required when creating a MachinePoolScope")
	}

	if params.AzureMachinePool == nil {
		return nil, errors.New("azure machine pool is required when creating a MachinePoolScope")
	}

	if params.AzureMachinePoolMachine == nil {
		return nil, errors.New("azure machine pool machine is required when creating a MachinePoolScope")
	}

	if params.workloadNodeGetter == nil {
		params.workloadNodeGetter = newWorkloadClusterProxy(
			params.Client,
			client.ObjectKey{
				Namespace: params.MachinePool.Namespace,
				Name:      params.ClusterScope.ClusterName(),
			},
		)
	}

	mpScope, err := NewMachinePoolScope(MachinePoolScopeParams{
		Client:           params.Client,
		MachinePool:      params.MachinePool,
		AzureMachinePool: params.AzureMachinePool,
		ClusterScope:     params.ClusterScope,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to build machine pool scope")
	}

	helper, err := patch.NewHelper(params.AzureMachinePoolMachine, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}

	return &MachinePoolMachineScope{
		AzureMachinePool:        params.AzureMachinePool,
		AzureMachinePoolMachine: params.AzureMachinePoolMachine,
		ClusterScoper:           params.ClusterScope,
		MachinePool:             params.MachinePool,
		MachinePoolScope:        mpScope,
		client:                  params.Client,
		patchHelper:             helper,
		workloadNodeGetter:      params.workloadNodeGetter,
	}, nil
}

// ScaleSetVMSpec returns the VMSS VM spec.
func (s *MachinePoolMachineScope) ScaleSetVMSpec() azure.ResourceSpecGetter {
	spec := &scalesetvms.ScaleSetVMSpec{
		Name:          s.Name(),
		InstanceID:    s.InstanceID(),
		ResourceGroup: s.ResourceGroup(),
		ScaleSetName:  s.ScaleSetName(),
		ProviderID:    s.ProviderID(),
		IsFlex:        s.OrchestrationMode() == infrav1.FlexibleOrchestrationMode,
	}

	if spec.IsFlex {
		spec.ResourceID = strings.TrimPrefix(spec.ProviderID, azureutil.ProviderIDPrefix)
	}

	return spec
}

// Name is the name of the Machine Pool Machine.
func (s *MachinePoolMachineScope) Name() string {
	return s.AzureMachinePoolMachine.Name
}

// InstanceID is the unique ID of the machine within the Machine Pool.
func (s *MachinePoolMachineScope) InstanceID() string {
	return s.AzureMachinePoolMachine.Spec.InstanceID
}

// ScaleSetName is the name of the VMSS.
func (s *MachinePoolMachineScope) ScaleSetName() string {
	return s.MachinePoolScope.Name()
}

// OrchestrationMode is the VMSS orchestration mode, either Uniform or Flexible.
func (s *MachinePoolMachineScope) OrchestrationMode() infrav1.OrchestrationModeType {
	return s.AzureMachinePool.Spec.OrchestrationMode
}

// SetLongRunningOperationState will set the future on the AzureMachinePoolMachine status to allow the resource to continue
// in the next reconciliation.
func (s *MachinePoolMachineScope) SetLongRunningOperationState(future *infrav1.Future) {
	futures.Set(s.AzureMachinePoolMachine, future)
}

// GetLongRunningOperationState will get the future on the AzureMachinePoolMachine status.
func (s *MachinePoolMachineScope) GetLongRunningOperationState(name, service, futureType string) *infrav1.Future {
	return futures.Get(s.AzureMachinePoolMachine, name, service, futureType)
}

// DeleteLongRunningOperationState will delete the future from the AzureMachinePoolMachine status.
func (s *MachinePoolMachineScope) DeleteLongRunningOperationState(name, service, futureType string) {
	futures.Delete(s.AzureMachinePoolMachine, name, service, futureType)
}

// UpdateDeleteStatus updates a condition on the AzureMachinePoolMachine status after a DELETE operation.
func (s *MachinePoolMachineScope) UpdateDeleteStatus(condition clusterv1.ConditionType, service string, err error) {
	switch {
	case err == nil:
		conditions.MarkFalse(s.AzureMachinePoolMachine, condition, infrav1.DeletedReason, clusterv1.ConditionSeverityInfo, "%s successfully deleted", service)
	case azure.IsOperationNotDoneError(err):
		conditions.MarkFalse(s.AzureMachinePoolMachine, condition, infrav1.DeletingReason, clusterv1.ConditionSeverityInfo, "%s deleting", service)
	default:
		conditions.MarkFalse(s.AzureMachinePoolMachine, condition, infrav1.DeletionFailedReason, clusterv1.ConditionSeverityError, "%s failed to delete. err: %s", service, err.Error())
	}
}

// UpdatePutStatus updates a condition on the AzureMachinePoolMachine status after a PUT operation.
func (s *MachinePoolMachineScope) UpdatePutStatus(condition clusterv1.ConditionType, service string, err error) {
	switch {
	case err == nil:
		conditions.MarkTrue(s.AzureMachinePoolMachine, condition)
	case azure.IsOperationNotDoneError(err):
		conditions.MarkFalse(s.AzureMachinePoolMachine, condition, infrav1.CreatingReason, clusterv1.ConditionSeverityInfo, "%s creating or updating", service)
	default:
		conditions.MarkFalse(s.AzureMachinePoolMachine, condition, infrav1.FailedReason, clusterv1.ConditionSeverityError, "%s failed to create or update. err: %s", service, err.Error())
	}
}

// UpdatePatchStatus updates a condition on the AzureMachinePoolMachine status after a PATCH operation.
func (s *MachinePoolMachineScope) UpdatePatchStatus(condition clusterv1.ConditionType, service string, err error) {
	switch {
	case err == nil:
		conditions.MarkTrue(s.AzureMachinePoolMachine, condition)
	case azure.IsOperationNotDoneError(err):
		conditions.MarkFalse(s.AzureMachinePoolMachine, condition, infrav1.UpdatingReason, clusterv1.ConditionSeverityInfo, "%s updating", service)
	default:
		conditions.MarkFalse(s.AzureMachinePoolMachine, condition, infrav1.FailedReason, clusterv1.ConditionSeverityError, "%s failed to update. err: %s", service, err.Error())
	}
}

// SetVMSSVM update the scope with the current state of the VMSS VM.
func (s *MachinePoolMachineScope) SetVMSSVM(instance *azure.VMSSVM) {
	s.instance = instance
}

// SetVMSSVMState update the scope with the current provisioning state of the VMSS VM.
func (s *MachinePoolMachineScope) SetVMSSVMState(state infrav1.ProvisioningState) {
	if s.instance != nil {
		s.instance.State = state
	}
}

// ProvisioningState returns the AzureMachinePoolMachine provisioning state.
func (s *MachinePoolMachineScope) ProvisioningState() infrav1.ProvisioningState {
	if s.AzureMachinePoolMachine.Status.ProvisioningState != nil {
		return *s.AzureMachinePoolMachine.Status.ProvisioningState
	}
	return ""
}

// IsReady indicates the machine has successfully provisioned and has a node ref associated.
func (s *MachinePoolMachineScope) IsReady() bool {
	state := s.AzureMachinePoolMachine.Status.ProvisioningState
	return s.AzureMachinePoolMachine.Status.Ready && state != nil && *state == infrav1.Succeeded
}

// SetFailureMessage sets the AzureMachinePoolMachine status failure message.
func (s *MachinePoolMachineScope) SetFailureMessage(v error) {
	s.AzureMachinePoolMachine.Status.FailureMessage = ptr.To(v.Error())
}

// SetFailureReason sets the AzureMachinePoolMachine status failure reason.
func (s *MachinePoolMachineScope) SetFailureReason(v capierrors.MachineStatusError) {
	s.AzureMachinePoolMachine.Status.FailureReason = &v
}

// ProviderID returns the AzureMachinePool ID by parsing Spec.FakeProviderID.
func (s *MachinePoolMachineScope) ProviderID() string {
	return s.AzureMachinePoolMachine.Spec.ProviderID
}

// PatchObject persists the MachinePoolMachine spec and status.
func (s *MachinePoolMachineScope) PatchObject(ctx context.Context) error {
	conditions.SetSummary(s.AzureMachinePoolMachine)

	return s.patchHelper.Patch(
		ctx,
		s.AzureMachinePoolMachine,
		patch.WithOwnedConditions{Conditions: []clusterv1.ConditionType{
			clusterv1.ReadyCondition,
			clusterv1.MachineNodeHealthyCondition,
			clusterv1.DrainingSucceededCondition,
		}})
}

// Close updates the state of MachinePoolMachine.
func (s *MachinePoolMachineScope) Close(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(
		ctx,
		"scope.MachinePoolMachineScope.Close",
	)
	defer done()

	return s.PatchObject(ctx)
}

// UpdateNodeStatus updates AzureMachinePoolMachine conditions and ready status. It will also update the node ref and the Kubernetes
// version of the VM instance if the node is found.
// Note: This func should be called at the end of a reconcile request and after updating the scope with the most recent Azure data.
func (s *MachinePoolMachineScope) UpdateNodeStatus(ctx context.Context) error {
	ctx, log, done := tele.StartSpanWithLogger(
		ctx,
		"scope.MachinePoolMachineScope.UpdateNodeStatus",
	)
	defer done()

	if s.instance != nil {
		switch s.instance.BootstrappingState {
		case infrav1.Creating:
			conditions.MarkFalse(s.AzureMachinePoolMachine, infrav1.BootstrapSucceededCondition, infrav1.BootstrapInProgressReason, clusterv1.ConditionSeverityInfo, "VM bootstrapping")
		case infrav1.Failed:
			log.Info("VM bootstrapping failed")
			conditions.MarkFalse(s.AzureMachinePoolMachine, infrav1.BootstrapSucceededCondition, infrav1.BootstrapFailedReason, clusterv1.ConditionSeverityInfo, "VM bootstrapping failed")
		case infrav1.Succeeded:
			log.Info("VM bootstrapping succeeded")
			conditions.MarkTrue(s.AzureMachinePoolMachine, infrav1.BootstrapSucceededCondition)
		}
	}

	var node *corev1.Node
	nodeRef := s.AzureMachinePoolMachine.Status.NodeRef

	// See if we can fetch a node using either the providerID or the nodeRef
	node, found, err := s.GetNode(ctx)
	switch {
	case err != nil && apierrors.IsNotFound(err) && nodeRef != nil && nodeRef.Name != "":
		// Node was not found due to 404 when finding by ObjectReference.
		conditions.MarkFalse(s.AzureMachinePoolMachine, clusterv1.MachineNodeHealthyCondition, clusterv1.NodeNotFoundReason, clusterv1.ConditionSeverityError, "")
	case err != nil:
		// Failed due to an unexpected error
		return err
	case !found && s.ProviderID() == "":
		// Node was not found due to not having a providerID set
		conditions.MarkFalse(s.AzureMachinePoolMachine, clusterv1.MachineNodeHealthyCondition, clusterv1.WaitingForNodeRefReason, clusterv1.ConditionSeverityInfo, "")
	case !found && s.ProviderID() != "":
		// Node was not found due to not finding a matching node by providerID
		conditions.MarkFalse(s.AzureMachinePoolMachine, clusterv1.MachineNodeHealthyCondition, clusterv1.NodeProvisioningReason, clusterv1.ConditionSeverityInfo, "")
	default:
		// Node was found. Check if it is ready.
		nodeReady := noderefutil.IsNodeReady(node)
		s.AzureMachinePoolMachine.Status.Ready = nodeReady
		if nodeReady {
			conditions.MarkTrue(s.AzureMachinePoolMachine, clusterv1.MachineNodeHealthyCondition)
		} else {
			conditions.MarkFalse(s.AzureMachinePoolMachine, clusterv1.MachineNodeHealthyCondition, clusterv1.NodeConditionsFailedReason, clusterv1.ConditionSeverityWarning, "")
		}

		s.AzureMachinePoolMachine.Status.NodeRef = &corev1.ObjectReference{
			Kind:       node.Kind,
			Namespace:  node.Namespace,
			Name:       node.Name,
			UID:        node.UID,
			APIVersion: node.APIVersion,
		}

		s.AzureMachinePoolMachine.Status.Version = node.Status.NodeInfo.KubeletVersion
	}

	return nil
}

// UpdateInstanceStatus updates the provisioning state of the AzureMachinePoolMachine and if it has the latest model applied
// using the VMSS VM instance.
// Note: This func should be called at the end of a reconcile request and after updating the scope with the most recent Azure data.
func (s *MachinePoolMachineScope) UpdateInstanceStatus(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(
		ctx,
		"scope.MachinePoolMachineScope.UpdateInstanceStatus",
	)
	defer done()

	if s.instance != nil {
		s.AzureMachinePoolMachine.Status.ProvisioningState = &s.instance.State
		hasLatestModel, err := s.hasLatestModelApplied(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to determine if the VMSS instance has the latest model")
		}

		s.AzureMachinePoolMachine.Status.LatestModelApplied = hasLatestModel
	}

	return nil
}

// CordonAndDrain will cordon and drain the Kubernetes node associated with this AzureMachinePoolMachine.
func (s *MachinePoolMachineScope) CordonAndDrain(ctx context.Context) error {
	ctx, log, done := tele.StartSpanWithLogger(
		ctx,
		"scope.MachinePoolMachineScope.CordonAndDrain",
	)
	defer done()

	// See if we can fetch a node using either the providerID or the nodeRef
	node, found, err := s.GetNode(ctx)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		// failed due to an unexpected error
		return errors.Wrap(err, "failed to get node")
	} else if !found {
		// node was not found due to not finding a nodes with the ProviderID
		return nil
	}

	// Drain node before deletion and issue a patch in order to make this operation visible to the users.
	if s.isNodeDrainAllowed() {
		patchHelper, err := patch.NewHelper(s.AzureMachinePoolMachine, s.client)
		if err != nil {
			return errors.Wrap(err, "failed to build a patchHelper when draining node")
		}

		log.V(4).Info("Draining node", "node", node.Name)
		// The DrainingSucceededCondition never exists before the node is drained for the first time,
		// so its transition time can be used to record the first time draining.
		// This `if` condition prevents the transition time to be changed more than once.
		if conditions.Get(s.AzureMachinePoolMachine, clusterv1.DrainingSucceededCondition) == nil {
			conditions.MarkFalse(s.AzureMachinePoolMachine, clusterv1.DrainingSucceededCondition, clusterv1.DrainingReason, clusterv1.ConditionSeverityInfo, "Draining the node before deletion")
		}

		if err := patchHelper.Patch(ctx, s.AzureMachinePoolMachine); err != nil {
			return errors.Wrap(err, "failed to patch AzureMachinePoolMachine")
		}

		if err := s.drainNode(ctx, node); err != nil {
			// Check for condition existence. If the condition exists, it may have a different severity or message, which
			// would cause the last transition time to be updated. The last transition time is used to determine how
			// long to wait to timeout the node drain operation. If we were to keep updating the last transition time,
			// a drain operation may never timeout.
			if conditions.Get(s.AzureMachinePoolMachine, clusterv1.DrainingSucceededCondition) == nil {
				conditions.MarkFalse(s.AzureMachinePoolMachine, clusterv1.DrainingSucceededCondition, clusterv1.DrainingFailedReason, clusterv1.ConditionSeverityWarning, err.Error())
			}
			return err
		}

		conditions.MarkTrue(s.AzureMachinePoolMachine, clusterv1.DrainingSucceededCondition)
	}

	return nil
}

func (s *MachinePoolMachineScope) drainNode(ctx context.Context, node *corev1.Node) error {
	ctx, log, done := tele.StartSpanWithLogger(
		ctx,
		"scope.MachinePoolMachineScope.drainNode",
	)
	defer done()

	restConfig, err := remote.RESTConfig(ctx, MachinePoolMachineScopeName, s.client, client.ObjectKey{
		Name:      s.ClusterName(),
		Namespace: s.AzureMachinePoolMachine.Namespace,
	})

	if err != nil {
		log.Error(err, "Error creating a remote client while deleting Machine, won't retry")
		return nil
	}

	kubeClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		log.Error(err, "Error creating a remote client while deleting Machine, won't retry")
		return nil
	}

	drainer := &kubedrain.Helper{
		Client:              kubeClient,
		Ctx:                 ctx,
		Force:               true,
		IgnoreAllDaemonSets: true,
		DeleteEmptyDirData:  true,
		GracePeriodSeconds:  -1,
		// If a pod is not evicted in 20 seconds, retry the eviction next time the
		// machine gets reconciled again (to allow other machines to be reconciled).
		Timeout: 20 * time.Second,
		OnPodDeletedOrEvicted: func(pod *corev1.Pod, usingEviction bool) {
			verbStr := "Deleted"
			if usingEviction {
				verbStr = "Evicted"
			}
			log.V(4).Info(fmt.Sprintf("%s pod from Node", verbStr),
				"pod", fmt.Sprintf("%s/%s", pod.Name, pod.Namespace))
		},
		Out:    writer{klog.Info},
		ErrOut: writer{klog.Error},
	}

	if noderefutil.IsNodeUnreachable(node) {
		// When the node is unreachable and some pods are not evicted for as long as this timeout, we ignore them.
		drainer.SkipWaitForDeleteTimeoutSeconds = 60 * 5 // 5 minutes
	}

	if err := kubedrain.RunCordonOrUncordon(drainer, node, true); err != nil {
		// Machine will be re-reconciled after a cordon failure.
		return azure.WithTransientError(errors.Errorf("unable to cordon node %s: %v", node.Name, err), 20*time.Second)
	}

	if err := kubedrain.RunNodeDrain(drainer, node.Name); err != nil {
		// Machine will be re-reconciled after a drain failure.
		return azure.WithTransientError(errors.Wrap(err, "Drain failed, retry in 20s"), 20*time.Second)
	}

	log.V(4).Info("Drain successful")
	return nil
}

// isNodeDrainAllowed checks to see the node is excluded from draining or if the NodeDrainTimeout has expired.
func (s *MachinePoolMachineScope) isNodeDrainAllowed() bool {
	if _, exists := s.AzureMachinePoolMachine.ObjectMeta.Annotations[clusterv1.ExcludeNodeDrainingAnnotation]; exists {
		return false
	}

	if s.nodeDrainTimeoutExceeded() {
		return false
	}

	return true
}

// nodeDrainTimeoutExceeded will check to see if the AzureMachinePool's NodeDrainTimeout is exceeded for the
// AzureMachinePoolMachine.
func (s *MachinePoolMachineScope) nodeDrainTimeoutExceeded() bool {
	// if the NodeDrainTineout type is not set by user
	pool := s.AzureMachinePool
	if pool == nil || pool.Spec.NodeDrainTimeout == nil || pool.Spec.NodeDrainTimeout.Seconds() <= 0 {
		return false
	}

	// if the draining succeeded condition does not exist
	if conditions.Get(s.AzureMachinePoolMachine, clusterv1.DrainingSucceededCondition) == nil {
		return false
	}

	now := time.Now()
	firstTimeDrain := conditions.GetLastTransitionTime(s.AzureMachinePoolMachine, clusterv1.DrainingSucceededCondition)
	diff := now.Sub(firstTimeDrain.Time)
	return diff.Seconds() >= s.AzureMachinePool.Spec.NodeDrainTimeout.Seconds()
}

func (s *MachinePoolMachineScope) hasLatestModelApplied(ctx context.Context) (bool, error) {
	ctx, _, done := tele.StartSpanWithLogger(
		ctx,
		"scope.MachinePoolMachineScope.hasLatestModelApplied",
	)
	defer done()

	if s.instance == nil {
		return false, errors.New("instance must not be nil")
	}

	image, err := s.MachinePoolScope.GetVMImage(ctx)
	if err != nil {
		return false, errors.Wrap(err, "unable to build vm image information from MachinePoolScope")
	}

	// this should never happen as GetVMImage should only return nil when err != nil. Just in case.
	if image == nil {
		return false, errors.New("machinepoolscope image must not be nil")
	}

	// check if image.ID is actually a compute gallery image
	if s.instance.Image.ComputeGallery != nil && image.ID != nil {
		newImage := converters.IDImageRefToImage(*image.ID)

		// this means the ID was a compute gallery image ID
		if newImage.ComputeGallery != nil {
			return reflect.DeepEqual(s.instance.Image, newImage), nil
		}
	}

	// if the images match, then the VM is of the same model
	return reflect.DeepEqual(s.instance.Image, *image), nil
}

func newWorkloadClusterProxy(c client.Client, cluster client.ObjectKey) *workloadClusterProxy {
	return &workloadClusterProxy{
		Client:  c,
		Cluster: cluster,
	}
}

// GetNode returns the node associated with the AzureMachinePoolMachine. Returns an error if one occurred, and a boolean
// indicating if the node was found if there was no error.
func (s *MachinePoolMachineScope) GetNode(ctx context.Context) (*corev1.Node, bool, error) {
	ctx, _, done := tele.StartSpanWithLogger(
		ctx,
		"scope.MachinePoolMachineScope.GetNode",
	)
	defer done()

	var (
		nodeRef = s.AzureMachinePoolMachine.Status.NodeRef
		node    *corev1.Node
		err     error
	)

	if nodeRef == nil || nodeRef.Name == "" {
		node, err = s.workloadNodeGetter.GetNodeByProviderID(ctx, s.ProviderID())
		if err != nil {
			return nil, false, errors.Wrap(err, "failed to get node by providerID")
		}
	} else {
		node, err = s.workloadNodeGetter.GetNodeByObjectReference(ctx, *nodeRef)
		if err != nil {
			return nil, false, errors.Wrap(err, "failed to get node by object reference")
		}
	}

	if node == nil {
		return nil, false, nil
	}

	return node, true, nil
}

// GetNodeByObjectReference will fetch a *corev1.Node via a node object reference.
func (np *workloadClusterProxy) GetNodeByObjectReference(ctx context.Context, nodeRef corev1.ObjectReference) (*corev1.Node, error) {
	workloadClient, err := getWorkloadClient(ctx, np.Client, np.Cluster)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create the workload cluster client")
	}

	var node corev1.Node
	err = workloadClient.Get(ctx, client.ObjectKey{
		Namespace: nodeRef.Namespace,
		Name:      nodeRef.Name,
	}, &node)

	return &node, err
}

// GetNodeByProviderID will fetch a node from the workload cluster by it's providerID.
func (np *workloadClusterProxy) GetNodeByProviderID(ctx context.Context, providerID string) (*corev1.Node, error) {
	ctx, _, done := tele.StartSpanWithLogger(
		ctx,
		"scope.MachinePoolMachineScope.getNode",
	)
	defer done()

	workloadClient, err := getWorkloadClient(ctx, np.Client, np.Cluster)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create the workload cluster client")
	}

	return getNodeByProviderID(ctx, workloadClient, providerID)
}

func getNodeByProviderID(ctx context.Context, workloadClient client.Client, providerID string) (*corev1.Node, error) {
	ctx, _, done := tele.StartSpanWithLogger(
		ctx,
		"scope.MachinePoolMachineScope.getNodeRefForProviderID",
	)
	defer done()

	nodeList := corev1.NodeList{}
	for {
		if err := workloadClient.List(ctx, &nodeList, client.Continue(nodeList.Continue)); err != nil {
			return nil, errors.Wrapf(err, "failed to List nodes")
		}

		for _, node := range nodeList.Items {
			if node.Spec.ProviderID == providerID {
				return &node, nil
			}
		}

		if nodeList.Continue == "" {
			break
		}
	}

	return nil, nil
}

func getWorkloadClient(ctx context.Context, c client.Client, cluster client.ObjectKey) (client.Client, error) {
	ctx, _, done := tele.StartSpanWithLogger(
		ctx,
		"scope.MachinePoolMachineScope.getWorkloadClient",
	)
	defer done()

	return remote.NewClusterClient(ctx, MachinePoolMachineScopeName, c, cluster)
}

// writer implements io.Writer interface as a pass-through for klog.
type writer struct {
	logFunc func(args ...interface{})
}

// Write passes string(p) into writer's logFunc and always returns len(p).
func (w writer) Write(p []byte) (n int, err error) {
	w.logFunc(string(p))
	return len(p), nil
}
