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
	"reflect"
	"strings"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/utils/ptr"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/controllers/noderefutil"
	"sigs.k8s.io/cluster-api/controllers/remote"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
	v1beta1patch "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/converters"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/scalesetvms"
	infrav1exp "sigs.k8s.io/cluster-api-provider-azure/exp/api/v1beta1"
	azureutil "sigs.k8s.io/cluster-api-provider-azure/util/azure"
	"sigs.k8s.io/cluster-api-provider-azure/util/futures"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
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
		MachinePool             *clusterv1.MachinePool
		Machine                 *clusterv1.Machine

		// workloadNodeGetter is only used for testing purposes and provides a way for mocking requests to the workload cluster
		workloadNodeGetter nodeGetter
	}

	// MachinePoolMachineScope defines a scope defined around a machine pool machine.
	MachinePoolMachineScope struct {
		azure.ClusterScoper
		AzureMachinePoolMachine *infrav1exp.AzureMachinePoolMachine
		AzureMachinePool        *infrav1exp.AzureMachinePool
		MachinePool             *clusterv1.MachinePool
		Machine                 *clusterv1.Machine
		MachinePoolScope        *MachinePoolScope
		client                  client.Client
		patchHelper             *v1beta1patch.Helper
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

	if params.Machine == nil {
		return nil, errors.New("machine is required when creating a MachinePoolScope")
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

	helper, err := v1beta1patch.NewHelper(params.AzureMachinePoolMachine, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}

	return &MachinePoolMachineScope{
		AzureMachinePool:        params.AzureMachinePool,
		AzureMachinePoolMachine: params.AzureMachinePoolMachine,
		ClusterScoper:           params.ClusterScope,
		MachinePool:             params.MachinePool,
		Machine:                 params.Machine,
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
		ResourceGroup: s.NodeResourceGroup(),
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
func (s *MachinePoolMachineScope) UpdateDeleteStatus(condition clusterv1beta1.ConditionType, service string, err error) {
	switch {
	case err == nil:
		v1beta1conditions.MarkFalse(s.AzureMachinePoolMachine, condition, infrav1.DeletedReason, clusterv1beta1.ConditionSeverityInfo, "%s successfully deleted", service)
	case azure.IsOperationNotDoneError(err):
		v1beta1conditions.MarkFalse(s.AzureMachinePoolMachine, condition, infrav1.DeletingReason, clusterv1beta1.ConditionSeverityInfo, "%s deleting", service)
	default:
		v1beta1conditions.MarkFalse(s.AzureMachinePoolMachine, condition, infrav1.DeletionFailedReason, clusterv1beta1.ConditionSeverityError, "%s failed to delete. err: %s", service, err.Error())
	}
}

// UpdatePutStatus updates a condition on the AzureMachinePoolMachine status after a PUT operation.
func (s *MachinePoolMachineScope) UpdatePutStatus(condition clusterv1beta1.ConditionType, service string, err error) {
	switch {
	case err == nil:
		v1beta1conditions.MarkTrue(s.AzureMachinePoolMachine, condition)
	case azure.IsOperationNotDoneError(err):
		v1beta1conditions.MarkFalse(s.AzureMachinePoolMachine, condition, infrav1.CreatingReason, clusterv1beta1.ConditionSeverityInfo, "%s creating or updating", service)
	default:
		v1beta1conditions.MarkFalse(s.AzureMachinePoolMachine, condition, infrav1.FailedReason, clusterv1beta1.ConditionSeverityError, "%s failed to create or update. err: %s", service, err.Error())
	}
}

// UpdatePatchStatus updates a condition on the AzureMachinePoolMachine status after a PATCH operation.
func (s *MachinePoolMachineScope) UpdatePatchStatus(condition clusterv1beta1.ConditionType, service string, err error) {
	switch {
	case err == nil:
		v1beta1conditions.MarkTrue(s.AzureMachinePoolMachine, condition)
	case azure.IsOperationNotDoneError(err):
		v1beta1conditions.MarkFalse(s.AzureMachinePoolMachine, condition, infrav1.UpdatingReason, clusterv1beta1.ConditionSeverityInfo, "%s updating", service)
	default:
		v1beta1conditions.MarkFalse(s.AzureMachinePoolMachine, condition, infrav1.FailedReason, clusterv1beta1.ConditionSeverityError, "%s failed to update. err: %s", service, err.Error())
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
func (s *MachinePoolMachineScope) SetFailureReason(v string) {
	s.AzureMachinePoolMachine.Status.FailureReason = &v
}

// ProviderID returns the AzureMachinePool ID by parsing Spec.FakeProviderID.
func (s *MachinePoolMachineScope) ProviderID() string {
	return s.AzureMachinePoolMachine.Spec.ProviderID
}

// updateDeleteMachineAnnotation sets the clusterv1.DeleteMachineAnnotation on the AzureMachinePoolMachine if it exists on the owner Machine.
func (s *MachinePoolMachineScope) updateDeleteMachineAnnotation() {
	if s.Machine.Annotations != nil {
		if _, ok := s.Machine.Annotations[clusterv1.DeleteMachineAnnotation]; ok {
			if s.AzureMachinePoolMachine.Annotations == nil {
				s.AzureMachinePoolMachine.Annotations = map[string]string{}
			}

			s.AzureMachinePoolMachine.Annotations[clusterv1.DeleteMachineAnnotation] = "true"
		}
	}
}

// PatchObject persists the MachinePoolMachine spec and status.
func (s *MachinePoolMachineScope) PatchObject(ctx context.Context) error {
	v1beta1conditions.SetSummary(s.AzureMachinePoolMachine)

	return s.patchHelper.Patch(
		ctx,
		s.AzureMachinePoolMachine,
		v1beta1patch.WithOwnedConditions{Conditions: []clusterv1beta1.ConditionType{
			clusterv1beta1.ReadyCondition,
			clusterv1.MachineNodeHealthyCondition,
		}})
}

// Close updates the state of MachinePoolMachine.
func (s *MachinePoolMachineScope) Close(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(
		ctx,
		"scope.MachinePoolMachineScope.Close",
	)
	defer done()

	s.updateDeleteMachineAnnotation()

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
			v1beta1conditions.MarkFalse(s.AzureMachinePoolMachine, infrav1.BootstrapSucceededCondition, infrav1.BootstrapInProgressReason, clusterv1beta1.ConditionSeverityInfo, "VM bootstrapping")
		case infrav1.Failed:
			log.Info("VM bootstrapping failed")
			v1beta1conditions.MarkFalse(s.AzureMachinePoolMachine, infrav1.BootstrapSucceededCondition, infrav1.BootstrapFailedReason, clusterv1beta1.ConditionSeverityInfo, "VM bootstrapping failed")
		case infrav1.Succeeded:
			log.Info("VM bootstrapping succeeded")
			v1beta1conditions.MarkTrue(s.AzureMachinePoolMachine, infrav1.BootstrapSucceededCondition)
		}
	}

	var node *corev1.Node
	nodeRef := s.AzureMachinePoolMachine.Status.NodeRef

	// See if we can fetch a node using either the providerID or the nodeRef
	node, found, err := s.GetNode(ctx)
	switch {
	case err != nil && apierrors.IsNotFound(err) && nodeRef != nil && nodeRef.Name != "":
		// Node was not found due to 404 when finding by ObjectReference.
		v1beta1conditions.MarkFalse(s.AzureMachinePoolMachine, clusterv1.MachineNodeHealthyCondition, clusterv1beta1.NodeNotFoundReason, clusterv1beta1.ConditionSeverityError, "")
	case err != nil:
		// Failed due to an unexpected error
		return err
	case !found && s.ProviderID() == "":
		// Node was not found due to not having a providerID set
		v1beta1conditions.MarkFalse(s.AzureMachinePoolMachine, clusterv1.MachineNodeHealthyCondition, clusterv1beta1.WaitingForNodeRefReason, clusterv1beta1.ConditionSeverityInfo, "")
	case !found && s.ProviderID() != "":
		// Node was not found due to not finding a matching node by providerID
		v1beta1conditions.MarkFalse(s.AzureMachinePoolMachine, clusterv1.MachineNodeHealthyCondition, clusterv1beta1.NodeProvisioningReason, clusterv1beta1.ConditionSeverityInfo, "")
	default:
		// Node was found. Check if it is ready.
		nodeReady := noderefutil.IsNodeReady(node)
		s.AzureMachinePoolMachine.Status.Ready = nodeReady
		if nodeReady {
			v1beta1conditions.MarkTrue(s.AzureMachinePoolMachine, clusterv1.MachineNodeHealthyCondition)
		} else {
			v1beta1conditions.MarkFalse(s.AzureMachinePoolMachine, clusterv1.MachineNodeHealthyCondition, clusterv1beta1.NodeConditionsFailedReason, clusterv1beta1.ConditionSeverityWarning, "")
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
