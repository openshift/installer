/*
Copyright 2020 The Kubernetes Authors.

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
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v2"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	expv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/labels/format"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	machinepool "sigs.k8s.io/cluster-api-provider-azure/azure/scope/strategies/machinepool_deployments"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/resourceskus"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/roleassignments"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/scalesets"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/virtualmachineimages"
	infrav1exp "sigs.k8s.io/cluster-api-provider-azure/exp/api/v1beta1"
	azureutil "sigs.k8s.io/cluster-api-provider-azure/util/azure"
	"sigs.k8s.io/cluster-api-provider-azure/util/futures"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

// ScalesetsServiceName is the name of the scalesets service.
// TODO: move this to scalesets.go once we remove the usage in this package,
// added here to avoid a circular dependency.
const ScalesetsServiceName = "scalesets"

type (
	// MachinePoolScopeParams defines the input parameters used to create a new MachinePoolScope.
	MachinePoolScopeParams struct {
		Client           client.Client
		MachinePool      *expv1.MachinePool
		AzureMachinePool *infrav1exp.AzureMachinePool
		ClusterScope     azure.ClusterScoper
		Cache            *MachinePoolCache
	}

	// MachinePoolScope defines a scope defined around a machine pool and its cluster.
	MachinePoolScope struct {
		azure.ClusterScoper
		AzureMachinePool           *infrav1exp.AzureMachinePool
		MachinePool                *expv1.MachinePool
		client                     client.Client
		patchHelper                *patch.Helper
		capiMachinePoolPatchHelper *patch.Helper
		vmssState                  *azure.VMSS
		cache                      *MachinePoolCache
		skuCache                   *resourceskus.Cache
	}

	// NodeStatus represents the status of a Kubernetes node.
	NodeStatus struct {
		Ready   bool
		Version string
	}

	// MachinePoolCache stores common machine pool information so we don't have to hit the API multiple times within the same reconcile loop.
	MachinePoolCache struct {
		BootstrapData           string
		HasBootstrapDataChanges bool
		VMImage                 *infrav1.Image
		VMSKU                   resourceskus.SKU
		MaxSurge                int
	}
)

// NewMachinePoolScope creates a new MachinePoolScope from the supplied parameters.
// This is meant to be called for each reconcile iteration.
func NewMachinePoolScope(params MachinePoolScopeParams) (*MachinePoolScope, error) {
	if params.Client == nil {
		return nil, errors.New("client is required when creating a MachinePoolScope")
	}

	if params.MachinePool == nil {
		return nil, errors.New("machine pool is required when creating a MachinePoolScope")
	}

	if params.AzureMachinePool == nil {
		return nil, errors.New("azure machine pool is required when creating a MachinePoolScope")
	}

	helper, err := patch.NewHelper(params.AzureMachinePool, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}

	capiMachinePoolPatchHelper, err := patch.NewHelper(params.MachinePool, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init capi patch helper")
	}

	return &MachinePoolScope{
		client:                     params.Client,
		MachinePool:                params.MachinePool,
		AzureMachinePool:           params.AzureMachinePool,
		patchHelper:                helper,
		capiMachinePoolPatchHelper: capiMachinePoolPatchHelper,
		ClusterScoper:              params.ClusterScope,
	}, nil
}

// InitMachinePoolCache sets cached information about the machine pool to be used in the scope.
func (m *MachinePoolScope) InitMachinePoolCache(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "scope.MachinePoolScope.InitMachinePoolCache")
	defer done()

	if m.cache == nil {
		var err error
		m.cache = &MachinePoolCache{}

		m.cache.BootstrapData, err = m.GetBootstrapData(ctx)
		if err != nil {
			return err
		}

		m.cache.HasBootstrapDataChanges, err = m.HasBootstrapDataChanges(ctx)
		if err != nil {
			return err
		}

		m.cache.VMImage, err = m.GetVMImage(ctx)
		if err != nil {
			return err
		}
		m.SaveVMImageToStatus(m.cache.VMImage)

		m.cache.MaxSurge, err = m.MaxSurge()
		if err != nil {
			return err
		}

		if m.skuCache == nil {
			skuCache, err := resourceskus.GetCache(m, m.Location())
			if err != nil {
				return errors.Wrap(err, "failed to init resourceskus cache")
			}
			m.skuCache = skuCache
		}

		m.cache.VMSKU, err = m.skuCache.Get(ctx, m.AzureMachinePool.Spec.Template.VMSize, resourceskus.VirtualMachines)
		if err != nil {
			return errors.Wrapf(err, "failed to get VM SKU %s in compute api", m.AzureMachinePool.Spec.Template.VMSize)
		}
	}

	return nil
}

// ScaleSetSpec returns the scale set spec.
func (m *MachinePoolScope) ScaleSetSpec(ctx context.Context) azure.ResourceSpecGetter {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "scope.MachinePoolScope.ScaleSetSpec")
	defer done()

	spec := &scalesets.ScaleSetSpec{
		Name:                         m.Name(),
		ResourceGroup:                m.NodeResourceGroup(),
		Size:                         m.AzureMachinePool.Spec.Template.VMSize,
		Capacity:                     int64(ptr.Deref[int32](m.MachinePool.Spec.Replicas, 0)),
		SSHKeyData:                   m.AzureMachinePool.Spec.Template.SSHPublicKey,
		OSDisk:                       m.AzureMachinePool.Spec.Template.OSDisk,
		DataDisks:                    m.AzureMachinePool.Spec.Template.DataDisks,
		SubnetName:                   m.AzureMachinePool.Spec.Template.NetworkInterfaces[0].SubnetName,
		VNetName:                     m.Vnet().Name,
		VNetResourceGroup:            m.Vnet().ResourceGroup,
		PublicLBName:                 m.OutboundLBName(infrav1.Node),
		PublicLBAddressPoolName:      m.OutboundPoolName(infrav1.Node),
		AcceleratedNetworking:        m.AzureMachinePool.Spec.Template.NetworkInterfaces[0].AcceleratedNetworking,
		Identity:                     m.AzureMachinePool.Spec.Identity,
		UserAssignedIdentities:       m.AzureMachinePool.Spec.UserAssignedIdentities,
		DiagnosticsProfile:           m.AzureMachinePool.Spec.Template.Diagnostics,
		SecurityProfile:              m.AzureMachinePool.Spec.Template.SecurityProfile,
		SpotVMOptions:                m.AzureMachinePool.Spec.Template.SpotVMOptions,
		FailureDomains:               m.MachinePool.Spec.FailureDomains,
		TerminateNotificationTimeout: m.AzureMachinePool.Spec.Template.TerminateNotificationTimeout,
		NetworkInterfaces:            m.AzureMachinePool.Spec.Template.NetworkInterfaces,
		IPv6Enabled:                  m.IsIPv6Enabled(),
		OrchestrationMode:            m.AzureMachinePool.Spec.OrchestrationMode,
		Location:                     m.AzureMachinePool.Spec.Location,
		SubscriptionID:               m.SubscriptionID(),
		HasReplicasExternallyManaged: m.HasReplicasExternallyManaged(ctx),
		ClusterName:                  m.ClusterName(),
		AdditionalTags:               m.AdditionalTags(),
		PlatformFaultDomainCount:     m.AzureMachinePool.Spec.PlatformFaultDomainCount,
		ZoneBalance:                  m.AzureMachinePool.Spec.ZoneBalance,
	}

	if m.AzureMachinePool.Spec.ZoneBalance != nil && len(m.MachinePool.Spec.FailureDomains) <= 1 {
		log.V(4).Info("zone balance is enabled but one or less failure domains are specified, zone balance will be disabled")
		spec.ZoneBalance = nil
	}

	if m.cache != nil {
		spec.ShouldPatchCustomData = m.cache.HasBootstrapDataChanges
		log.V(4).Info("has bootstrap data changed?", "shouldPatchCustomData", spec.ShouldPatchCustomData)
		spec.VMSSExtensionSpecs = m.VMSSExtensionSpecs()
		spec.SKU = m.cache.VMSKU
		spec.VMImage = m.cache.VMImage
		spec.BootstrapData = m.cache.BootstrapData
		spec.MaxSurge = m.cache.MaxSurge
	} else {
		log.V(4).Info("machinepool cache is nil, this is only expected when deleting a machinepool")
	}

	return spec
}

// Name returns the Azure Machine Pool Name.
func (m *MachinePoolScope) Name() string {
	// Windows Machine pools names cannot be longer than 9 chars
	if m.AzureMachinePool.Spec.Template.OSDisk.OSType == azure.WindowsOS && len(m.AzureMachinePool.Name) > 9 {
		return "win-" + m.AzureMachinePool.Name[len(m.AzureMachinePool.Name)-5:]
	}
	return m.AzureMachinePool.Name
}

// SetInfrastructureMachineKind sets the infrastructure machine kind in the status if it is not set already, returning
// `true` if the status was updated. This supports MachinePool Machines.
func (m *MachinePoolScope) SetInfrastructureMachineKind() bool {
	if m.AzureMachinePool.Status.InfrastructureMachineKind != infrav1exp.AzureMachinePoolMachineKind {
		m.AzureMachinePool.Status.InfrastructureMachineKind = infrav1exp.AzureMachinePoolMachineKind

		return true
	}

	return false
}

// ProviderID returns the AzureMachinePool ID by parsing Spec.ProviderID.
func (m *MachinePoolScope) ProviderID() string {
	resourceID, err := azureutil.ParseResourceID(m.AzureMachinePool.Spec.ProviderID)
	if err != nil {
		return ""
	}
	return resourceID.Name
}

// SetProviderID sets the AzureMachinePool providerID in spec.
func (m *MachinePoolScope) SetProviderID(v string) {
	m.AzureMachinePool.Spec.ProviderID = v
}

// SystemAssignedIdentityName returns the scope for the system assigned identity.
func (m *MachinePoolScope) SystemAssignedIdentityName() string {
	if m.AzureMachinePool.Spec.SystemAssignedIdentityRole != nil {
		return m.AzureMachinePool.Spec.SystemAssignedIdentityRole.Name
	}
	return ""
}

// SystemAssignedIdentityScope returns the scope for the system assigned identity.
func (m *MachinePoolScope) SystemAssignedIdentityScope() string {
	if m.AzureMachinePool.Spec.SystemAssignedIdentityRole != nil {
		return m.AzureMachinePool.Spec.SystemAssignedIdentityRole.Scope
	}
	return ""
}

// SystemAssignedIdentityDefinitionID returns the role definition ID for the system assigned identity.
func (m *MachinePoolScope) SystemAssignedIdentityDefinitionID() string {
	if m.AzureMachinePool.Spec.SystemAssignedIdentityRole != nil {
		return m.AzureMachinePool.Spec.SystemAssignedIdentityRole.DefinitionID
	}
	return ""
}

// ProvisioningState returns the AzureMachinePool provisioning state.
func (m *MachinePoolScope) ProvisioningState() infrav1.ProvisioningState {
	if m.AzureMachinePool.Status.ProvisioningState != nil {
		return *m.AzureMachinePool.Status.ProvisioningState
	}
	return ""
}

// SetVMSSState updates the machine pool scope with the current state of the VMSS.
func (m *MachinePoolScope) SetVMSSState(vmssState *azure.VMSS) {
	m.vmssState = vmssState
}

// NeedsRequeue return true if any machines are not on the latest model or the VMSS is not in a terminal provisioning
// state.
func (m *MachinePoolScope) NeedsRequeue() bool {
	state := m.AzureMachinePool.Status.ProvisioningState
	if m.vmssState == nil {
		return state != nil && infrav1.IsTerminalProvisioningState(*state)
	}

	if !m.vmssState.HasLatestModelAppliedToAll() {
		return true
	}

	desiredMatchesActual := len(m.vmssState.Instances) == int(m.DesiredReplicas())
	return !(state != nil && infrav1.IsTerminalProvisioningState(*state) && desiredMatchesActual)
}

// DesiredReplicas returns the replica count on machine pool or 0 if machine pool replicas is nil.
func (m MachinePoolScope) DesiredReplicas() int32 {
	return ptr.Deref[int32](m.MachinePool.Spec.Replicas, 0)
}

// MaxSurge returns the number of machines to surge, or 0 if the deployment strategy does not support surge.
func (m MachinePoolScope) MaxSurge() (int, error) {
	if surger, ok := m.getDeploymentStrategy().(machinepool.Surger); ok {
		surgeCount, err := surger.Surge(int(m.DesiredReplicas()))
		if err != nil {
			return 0, errors.Wrap(err, "failed to calculate surge for the machine pool")
		}

		return surgeCount, nil
	}

	return 0, nil
}

// updateReplicasAndProviderIDs ties the Azure VMSS instance data and the Node status data together to build and update
// the AzureMachinePool replica count and providerIDList.
func (m *MachinePoolScope) updateReplicasAndProviderIDs(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "scope.MachinePoolScope.UpdateInstanceStatuses")
	defer done()

	machines, err := m.GetMachinePoolMachines(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get machine pool machines")
	}

	var readyReplicas int32
	providerIDs := make([]string, len(machines))
	for i, machine := range machines {
		if machine.Status.Ready {
			readyReplicas++
		}
		providerIDs[i] = machine.Spec.ProviderID
	}

	m.AzureMachinePool.Status.Replicas = readyReplicas
	m.AzureMachinePool.Spec.ProviderIDList = providerIDs
	return nil
}

func (m *MachinePoolScope) getMachinePoolMachineLabels() map[string]string {
	return map[string]string{
		clusterv1.ClusterNameLabel:      m.ClusterName(),
		infrav1exp.MachinePoolNameLabel: m.AzureMachinePool.Name,
		clusterv1.MachinePoolNameLabel:  format.MustFormatValue(m.MachinePool.Name),
		m.ClusterName():                 string(infrav1.ResourceLifecycleOwned),
	}
}

// GetMachinePoolMachines returns the list of AzureMachinePoolMachines associated with this AzureMachinePool.
func (m *MachinePoolScope) GetMachinePoolMachines(ctx context.Context) ([]infrav1exp.AzureMachinePoolMachine, error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "scope.MachinePoolScope.getMachinePoolMachines")
	defer done()

	labels := m.getMachinePoolMachineLabels()
	ampml := &infrav1exp.AzureMachinePoolMachineList{}
	if err := m.client.List(ctx, ampml, client.InNamespace(m.AzureMachinePool.Namespace), client.MatchingLabels(labels)); err != nil {
		return nil, errors.Wrap(err, "failed to list AzureMachinePoolMachines")
	}

	return ampml.Items, nil
}

func (m *MachinePoolScope) applyAzureMachinePoolMachines(ctx context.Context) error {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "scope.MachinePoolScope.applyAzureMachinePoolMachines")
	defer done()

	if m.vmssState == nil {
		return nil
	}

	ampms, err := m.GetMachinePoolMachines(ctx)
	if err != nil {
		return err
	}

	existingMachinesByProviderID := make(map[string]infrav1exp.AzureMachinePoolMachine, len(ampms))
	for _, ampm := range ampms {
		machine, err := util.GetOwnerMachine(ctx, m.client, ampm.ObjectMeta)
		if err != nil {
			return fmt.Errorf("failed to find owner machine for %s/%s: %w", ampm.Namespace, ampm.Name, err)
		}

		if _, ampmHasDeleteAnnotation := ampm.Annotations[clusterv1.DeleteMachineAnnotation]; !ampmHasDeleteAnnotation {
			// fetch Machine delete annotation from owner machine to AzureMachinePoolMachine.
			// This ensures setting a deleteMachine annotation on the Machine has an effect on the AzureMachinePoolMachine
			// and the deployment strategy, in case the automatic propagation of the annotation from Machine to AzureMachinePoolMachine
			// hasn't been done yet.
			if machine != nil && machine.Annotations != nil {
				if _, hasDeleteAnnotation := machine.Annotations[clusterv1.DeleteMachineAnnotation]; hasDeleteAnnotation {
					log.V(4).Info("fetched DeleteMachineAnnotation", "machine", ampm.Spec.ProviderID)
					if ampm.Annotations == nil {
						ampm.Annotations = make(map[string]string)
					}
					ampm.Annotations[clusterv1.DeleteMachineAnnotation] = machine.Annotations[clusterv1.DeleteMachineAnnotation]
				}
			}
		} else {
			log.V(4).Info("DeleteMachineAnnotation already set")
		}

		existingMachinesByProviderID[ampm.Spec.ProviderID] = ampm
	}

	// determine which machines need to be created to reflect the current state in Azure
	azureMachinesByProviderID := m.vmssState.InstancesByProviderID(m.AzureMachinePool.Spec.OrchestrationMode)
	for key, val := range azureMachinesByProviderID {
		if val.State == infrav1.Deleting || val.State == infrav1.Deleted {
			log.V(4).Info("not recreating AzureMachinePoolMachine because VMSS VM is deleting", "providerID", key)
			continue
		}
		if _, ok := existingMachinesByProviderID[key]; !ok {
			log.V(4).Info("creating AzureMachinePoolMachine", "providerID", key)
			if err := m.createMachine(ctx, val); err != nil {
				return errors.Wrap(err, "failed creating AzureMachinePoolMachine")
			}
			continue
		}
	}

	deleted := false
	// Delete MachinePool Machines for instances that no longer exist in Azure, i.e. deleted out-of-band
	for key, ampm := range existingMachinesByProviderID {
		if _, ok := azureMachinesByProviderID[key]; !ok {
			deleted = true
			log.V(4).Info("deleting AzureMachinePoolMachine because it no longer exists in the VMSS", "providerID", key)
			delete(existingMachinesByProviderID, key)
			if err := m.DeleteMachine(ctx, ampm); err != nil {
				return errors.Wrap(err, "failed deleting AzureMachinePoolMachine no longer existing in Azure")
			}
		}
	}

	if deleted {
		log.V(4).Info("exiting early due to finding AzureMachinePoolMachine(s) that were deleted because they no longer exist in the VMSS")
		// exit early to be less greedy about delete
		return nil
	}

	if futures.Has(m.AzureMachinePool, m.Name(), ScalesetsServiceName, infrav1.PatchFuture) ||
		futures.Has(m.AzureMachinePool, m.Name(), ScalesetsServiceName, infrav1.PutFuture) ||
		futures.Has(m.AzureMachinePool, m.Name(), ScalesetsServiceName, infrav1.DeleteFuture) {
		log.V(4).Info("exiting early due an in-progress long running operation on the ScaleSet")
		// exit early to be less greedy about delete
		return nil
	}

	// when replicas are externally managed, we do not want to scale down manually since that is handled by the external scaler.
	if m.HasReplicasExternallyManaged(ctx) {
		log.V(4).Info("exiting early due to replicas externally managed")
		return nil
	}

	deleteSelector := m.getDeploymentStrategy()
	if deleteSelector == nil {
		log.V(4).Info("can not select AzureMachinePoolMachines to delete because no deployment strategy is specified")
		return nil
	}

	// Select Machines to delete to lower the replica count
	toDelete, err := deleteSelector.SelectMachinesToDelete(ctx, m.DesiredReplicas(), existingMachinesByProviderID)
	if err != nil {
		return errors.Wrap(err, "failed selecting AzureMachinePoolMachine(s) to delete")
	}

	// Delete MachinePool Machines as a part of scaling down
	for i := range toDelete {
		ampm := toDelete[i]
		log.Info("deleting selected AzureMachinePoolMachine", "providerID", ampm.Spec.ProviderID)
		if err := m.DeleteMachine(ctx, ampm); err != nil {
			return errors.Wrap(err, "failed deleting AzureMachinePoolMachine to reduce replica count")
		}
	}

	log.V(4).Info("done reconciling AzureMachinePoolMachine(s)")
	return nil
}

func (m *MachinePoolScope) createMachine(ctx context.Context, machine azure.VMSSVM) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "scope.MachinePoolScope.createMachine")
	defer done()

	parsed, err := azureutil.ParseResourceID(machine.ID)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to parse resource id %q", machine.ID))
	}
	instanceID := strings.ReplaceAll(parsed.Name, "_", "-")

	ampm := infrav1exp.AzureMachinePoolMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.AzureMachinePool.Name + "-" + instanceID,
			Namespace: m.AzureMachinePool.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         infrav1exp.GroupVersion.String(),
					Kind:               infrav1.AzureMachinePoolKind,
					Name:               m.AzureMachinePool.Name,
					BlockOwnerDeletion: ptr.To(true),
					UID:                m.AzureMachinePool.UID,
				},
			},
			Annotations: map[string]string{},
		},
		Spec: infrav1exp.AzureMachinePoolMachineSpec{
			ProviderID: machine.ProviderID(),
			InstanceID: machine.InstanceID,
		},
	}

	labels := m.getMachinePoolMachineLabels()
	ampm.Labels = labels

	controllerutil.AddFinalizer(&ampm, infrav1exp.AzureMachinePoolMachineFinalizer)
	conditions.MarkFalse(&ampm, infrav1.VMRunningCondition, string(infrav1.Creating), clusterv1.ConditionSeverityInfo, "")
	if err := m.client.Create(ctx, &ampm); err != nil {
		return errors.Wrapf(err, "failed creating AzureMachinePoolMachine %s in AzureMachinePool %s", machine.ID, m.AzureMachinePool.Name)
	}

	return nil
}

// DeleteMachine deletes an AzureMachinePoolMachine by fetching its owner Machine and deleting it. This ensures that the node cordon/drain happens before deleting the infrastructure.
func (m *MachinePoolScope) DeleteMachine(ctx context.Context, ampm infrav1exp.AzureMachinePoolMachine) error {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "scope.MachinePoolScope.DeleteMachine")
	defer done()

	machine, err := util.GetOwnerMachine(ctx, m.client, ampm.ObjectMeta)
	if err != nil {
		return errors.Wrapf(err, "error getting owner Machine for AzureMachinePoolMachine %s/%s", ampm.Namespace, ampm.Name)
	}
	if machine == nil {
		log.V(2).Info("No owner Machine exists for AzureMachinePoolMachine", "ampm", klog.KObj(&ampm))
		// If the AzureMachinePoolMachine does not have an owner Machine, do not attempt to delete the AzureMachinePoolMachine as the MachinePool controller will create the
		// Machine and we want to let it catch up. If we are too hasty to delete, that introduces a race condition where the AzureMachinePoolMachine could be deleted
		// just as the Machine comes online.

		// In the case where the MachinePool is being deleted and the Machine will never come online, the AzureMachinePoolMachine will be deleted via its ownerRef to the
		// AzureMachinePool, so that is covered as well.

		return nil
	}

	if err := m.client.Delete(ctx, machine); err != nil {
		return errors.Wrapf(err, "failed to delete Machine %s for AzureMachinePoolMachine %s in MachinePool %s", machine.Name, ampm.Name, m.MachinePool.Name)
	}

	return nil
}

// SetLongRunningOperationState will set the future on the AzureMachinePool status to allow the resource to continue
// in the next reconciliation.
func (m *MachinePoolScope) SetLongRunningOperationState(future *infrav1.Future) {
	futures.Set(m.AzureMachinePool, future)
}

// GetLongRunningOperationState will get the future on the AzureMachinePool status.
func (m *MachinePoolScope) GetLongRunningOperationState(name, service, futureType string) *infrav1.Future {
	return futures.Get(m.AzureMachinePool, name, service, futureType)
}

// DeleteLongRunningOperationState will delete the future from the AzureMachinePool status.
func (m *MachinePoolScope) DeleteLongRunningOperationState(name, service, futureType string) {
	futures.Delete(m.AzureMachinePool, name, service, futureType)
}

// setProvisioningStateAndConditions sets the AzureMachinePool provisioning state and conditions.
func (m *MachinePoolScope) setProvisioningStateAndConditions(v infrav1.ProvisioningState) {
	m.AzureMachinePool.Status.ProvisioningState = &v
	switch {
	case v == infrav1.Succeeded && *m.MachinePool.Spec.Replicas == m.AzureMachinePool.Status.Replicas:
		// vmss is provisioned with enough ready replicas
		conditions.MarkTrue(m.AzureMachinePool, infrav1.ScaleSetRunningCondition)
		conditions.MarkTrue(m.AzureMachinePool, infrav1.ScaleSetModelUpdatedCondition)
		conditions.MarkTrue(m.AzureMachinePool, infrav1.ScaleSetDesiredReplicasCondition)
		m.SetReady()
	case v == infrav1.Succeeded && *m.MachinePool.Spec.Replicas != m.AzureMachinePool.Status.Replicas:
		// not enough ready or too many ready replicas we must still be scaling up or down
		updatingState := infrav1.Updating
		m.AzureMachinePool.Status.ProvisioningState = &updatingState
		if *m.MachinePool.Spec.Replicas > m.AzureMachinePool.Status.Replicas {
			conditions.MarkFalse(m.AzureMachinePool, infrav1.ScaleSetDesiredReplicasCondition, infrav1.ScaleSetScaleUpReason, clusterv1.ConditionSeverityInfo, "")
		} else {
			conditions.MarkFalse(m.AzureMachinePool, infrav1.ScaleSetDesiredReplicasCondition, infrav1.ScaleSetScaleDownReason, clusterv1.ConditionSeverityInfo, "")
		}
		m.SetNotReady()
	case v == infrav1.Updating:
		conditions.MarkFalse(m.AzureMachinePool, infrav1.ScaleSetModelUpdatedCondition, infrav1.ScaleSetModelOutOfDateReason, clusterv1.ConditionSeverityInfo, "")
		m.SetNotReady()
	case v == infrav1.Creating:
		conditions.MarkFalse(m.AzureMachinePool, infrav1.ScaleSetRunningCondition, infrav1.ScaleSetCreatingReason, clusterv1.ConditionSeverityInfo, "")
		m.SetNotReady()
	case v == infrav1.Deleting:
		conditions.MarkFalse(m.AzureMachinePool, infrav1.ScaleSetRunningCondition, infrav1.ScaleSetDeletingReason, clusterv1.ConditionSeverityInfo, "")
		m.SetNotReady()
	default:
		conditions.MarkFalse(m.AzureMachinePool, infrav1.ScaleSetRunningCondition, string(v), clusterv1.ConditionSeverityInfo, "")
		m.SetNotReady()
	}
}

// SetReady sets the AzureMachinePool Ready Status to true.
func (m *MachinePoolScope) SetReady() {
	m.AzureMachinePool.Status.Ready = true
}

// SetNotReady sets the AzureMachinePool Ready Status to false.
func (m *MachinePoolScope) SetNotReady() {
	m.AzureMachinePool.Status.Ready = false
}

// SetFailureMessage sets the AzureMachinePool status failure message.
func (m *MachinePoolScope) SetFailureMessage(v error) {
	m.AzureMachinePool.Status.FailureMessage = ptr.To(v.Error())
}

// SetFailureReason sets the AzureMachinePool status failure reason.
func (m *MachinePoolScope) SetFailureReason(v string) {
	m.AzureMachinePool.Status.FailureReason = &v
}

// AdditionalTags merges AdditionalTags from the scope's AzureCluster and AzureMachinePool. If the same key is present in both,
// the value from AzureMachinePool takes precedence.
func (m *MachinePoolScope) AdditionalTags() infrav1.Tags {
	tags := make(infrav1.Tags)
	// Start with the cluster-wide tags...
	tags.Merge(m.ClusterScoper.AdditionalTags())
	// ... and merge in the Machine Pool's
	tags.Merge(m.AzureMachinePool.Spec.AdditionalTags)
	// Set the cloud provider tag
	tags[infrav1.ClusterAzureCloudProviderTagKey(m.ClusterName())] = string(infrav1.ResourceLifecycleOwned)

	return tags
}

// SetAnnotation sets a key value annotation on the AzureMachinePool.
func (m *MachinePoolScope) SetAnnotation(key, value string) {
	if m.AzureMachinePool.Annotations == nil {
		m.AzureMachinePool.Annotations = map[string]string{}
	}
	m.AzureMachinePool.Annotations[key] = value
}

// PatchObject persists the AzureMachinePool spec and status.
func (m *MachinePoolScope) PatchObject(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "scope.MachinePoolScope.PatchObject")
	defer done()

	conditions.SetSummary(m.AzureMachinePool)
	return m.patchHelper.Patch(
		ctx,
		m.AzureMachinePool,
		patch.WithOwnedConditions{Conditions: []clusterv1.ConditionType{
			clusterv1.ReadyCondition,
			infrav1.BootstrapSucceededCondition,
			infrav1.ScaleSetDesiredReplicasCondition,
			infrav1.ScaleSetModelUpdatedCondition,
			infrav1.ScaleSetRunningCondition,
		}})
}

// Close the MachinePoolScope by updating the AzureMachinePool spec and AzureMachinePool status.
func (m *MachinePoolScope) Close(ctx context.Context) error {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "scope.MachinePoolScope.Close")
	defer done()

	if m.vmssState != nil {
		if err := m.applyAzureMachinePoolMachines(ctx); err != nil {
			log.Error(err, "failed to apply changes to the AzureMachinePoolMachines")
			return errors.Wrap(err, "failed to apply changes to AzureMachinePoolMachines")
		}

		m.setProvisioningStateAndConditions(m.vmssState.State)
		if err := m.updateReplicasAndProviderIDs(ctx); err != nil {
			return errors.Wrap(err, "failed to update replicas and providerIDs")
		}
		if err := m.updateCustomDataHash(ctx); err != nil {
			// ignore errors to calculating the custom data hash since it's not absolutely crucial.
			log.V(4).Error(err, "unable to update custom data hash, ignoring.")
		}
	}

	if err := m.PatchObject(ctx); err != nil {
		return errors.Wrap(err, "unable to patch AzureMachinePool")
	}
	if err := m.PatchCAPIMachinePoolObject(ctx); err != nil {
		return errors.Wrap(err, "unable to patch CAPI MachinePool")
	}
	return nil
}

// GetBootstrapData returns the bootstrap data from the secret in the MachinePool's bootstrap.dataSecretName.
func (m *MachinePoolScope) GetBootstrapData(ctx context.Context) (string, error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "scope.MachinePoolScope.GetBootstrapData")
	defer done()

	dataSecretName := m.MachinePool.Spec.Template.Spec.Bootstrap.DataSecretName
	if dataSecretName == nil {
		return "", errors.New("error retrieving bootstrap data: linked MachinePool Spec's bootstrap.dataSecretName is nil")
	}
	secret := &corev1.Secret{}
	key := types.NamespacedName{Namespace: m.AzureMachinePool.Namespace, Name: *dataSecretName}
	if err := m.client.Get(ctx, key, secret); err != nil {
		return "", errors.Wrapf(err, "failed to retrieve bootstrap data secret for AzureMachinePool %s/%s", m.AzureMachinePool.Namespace, m.Name())
	}

	value, ok := secret.Data["value"]
	if !ok {
		return "", errors.New("error retrieving bootstrap data: secret value key is missing")
	}
	return base64.StdEncoding.EncodeToString(value), nil
}

// calculateBootstrapDataHash calculates the sha256 hash of the bootstrap data.
func (m *MachinePoolScope) calculateBootstrapDataHash(_ context.Context) (string, error) {
	bootstrapData := m.cache.BootstrapData
	h := sha256.New()
	n, err := io.WriteString(h, bootstrapData)
	if err != nil || n == 0 {
		return "", fmt.Errorf("unable to write custom data (bytes written: %q): %w", n, err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// HasBootstrapDataChanges calculates the sha256 hash of the bootstrap data and compares it with the saved hash in AzureMachinePool.Status.
func (m *MachinePoolScope) HasBootstrapDataChanges(ctx context.Context) (bool, error) {
	newHash, err := m.calculateBootstrapDataHash(ctx)
	if err != nil {
		return false, err
	}
	return m.AzureMachinePool.GetAnnotations()[azure.CustomDataHashAnnotation] != newHash, nil
}

// updateCustomDataHash calculates the sha256 hash of the bootstrap data and saves it in AzureMachinePool.Status.
func (m *MachinePoolScope) updateCustomDataHash(ctx context.Context) error {
	newHash, err := m.calculateBootstrapDataHash(ctx)
	if err != nil {
		return err
	}
	m.SetAnnotation(azure.CustomDataHashAnnotation, newHash)
	return nil
}

// GetVMImage picks an image from the AzureMachinePool configuration, or uses a default one.
func (m *MachinePoolScope) GetVMImage(ctx context.Context) (*infrav1.Image, error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "scope.MachinePoolScope.GetVMImage")
	defer done()

	// Use custom Marketplace image, Image ID or a Shared Image Gallery image if provided
	if m.AzureMachinePool.Spec.Template.Image != nil {
		return m.AzureMachinePool.Spec.Template.Image, nil
	}

	var (
		err          error
		defaultImage *infrav1.Image
	)

	svc, err := virtualmachineimages.New(m)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create virtualmachineimages service")
	}

	if m.AzureMachinePool.Spec.Template.OSDisk.OSType == azure.WindowsOS {
		runtime := m.AzureMachinePool.Annotations["runtime"]
		windowsServerVersion := m.AzureMachinePool.Annotations["windowsServerVersion"]
		log.V(4).Info("No image specified for machine, using default Windows Image", "machine", m.MachinePool.GetName(), "runtime", runtime, "windowsServerVersion", windowsServerVersion)
		defaultImage, err = svc.GetDefaultWindowsImage(ctx, m.Location(), ptr.Deref(m.MachinePool.Spec.Template.Spec.Version, ""), runtime, windowsServerVersion)
	} else {
		defaultImage, err = svc.GetDefaultLinuxImage(ctx, m.Location(), ptr.Deref(m.MachinePool.Spec.Template.Spec.Version, ""))
	}

	if err != nil {
		return defaultImage, errors.Wrap(err, "failed to get default OS image")
	}

	return defaultImage, nil
}

// SaveVMImageToStatus persists the AzureMachinePool image to the status.
func (m *MachinePoolScope) SaveVMImageToStatus(image *infrav1.Image) {
	m.AzureMachinePool.Status.Image = image
}

// RoleAssignmentSpecs returns the role assignment specs.
func (m *MachinePoolScope) RoleAssignmentSpecs(principalID *string) []azure.ResourceSpecGetter {
	roles := make([]azure.ResourceSpecGetter, 1)
	if m.HasSystemAssignedIdentity() {
		roles[0] = &roleassignments.RoleAssignmentSpec{
			Name:             m.SystemAssignedIdentityName(),
			MachineName:      m.Name(),
			ResourceGroup:    m.NodeResourceGroup(),
			ResourceType:     azure.VirtualMachineScaleSet,
			Scope:            m.SystemAssignedIdentityScope(),
			RoleDefinitionID: m.SystemAssignedIdentityDefinitionID(),
			PrincipalID:      principalID,
			PrincipalType:    armauthorization.PrincipalTypeServicePrincipal,
		}
		return roles
	}
	return []azure.ResourceSpecGetter{}
}

// RoleAssignmentResourceType returns the role assignment resource type.
func (m *MachinePoolScope) RoleAssignmentResourceType() string {
	return azure.VirtualMachineScaleSet
}

// HasSystemAssignedIdentity returns true if the azure machine pool has system
// assigned identity.
func (m *MachinePoolScope) HasSystemAssignedIdentity() bool {
	return m.AzureMachinePool.Spec.Identity == infrav1.VMIdentitySystemAssigned
}

// VMSSExtensionSpecs returns the VMSS extension specs.
func (m *MachinePoolScope) VMSSExtensionSpecs() []azure.ResourceSpecGetter {
	var extensionSpecs = []azure.ResourceSpecGetter{}

	for _, extension := range m.AzureMachinePool.Spec.Template.VMExtensions {
		extensionSpecs = append(extensionSpecs, &scalesets.VMSSExtensionSpec{
			ExtensionSpec: azure.ExtensionSpec{
				Name:              extension.Name,
				VMName:            m.Name(),
				Publisher:         extension.Publisher,
				Version:           extension.Version,
				Settings:          extension.Settings,
				ProtectedSettings: extension.ProtectedSettings,
			},
			ResourceGroup: m.NodeResourceGroup(),
		})
	}

	cpuArchitectureType, _ := m.cache.VMSKU.GetCapability(resourceskus.CPUArchitectureType)
	bootstrapExtensionSpec := azure.GetBootstrappingVMExtension(m.AzureMachinePool.Spec.Template.OSDisk.OSType, m.CloudEnvironment(), m.Name(), cpuArchitectureType)

	if bootstrapExtensionSpec != nil {
		extensionSpecs = append(extensionSpecs, &scalesets.VMSSExtensionSpec{
			ExtensionSpec: *bootstrapExtensionSpec,
			ResourceGroup: m.NodeResourceGroup(),
		})
	}

	return extensionSpecs
}

func (m *MachinePoolScope) getDeploymentStrategy() machinepool.TypedDeleteSelector {
	if m.AzureMachinePool == nil {
		return nil
	}

	return machinepool.NewMachinePoolDeploymentStrategy(m.AzureMachinePool.Spec.Strategy)
}

// SetSubnetName defaults the AzureMachinePool subnet name to the name of the subnet with role 'node' when there is only one of them.
// Note: this logic exists only for purposes of ensuring backwards compatibility for old clusters created without the `subnetName` field being
// set, and should be removed in the future when this field is no longer optional.
func (m *MachinePoolScope) SetSubnetName() error {
	if m.AzureMachinePool.Spec.Template.NetworkInterfaces[0].SubnetName == "" {
		subnetName := ""
		for _, subnet := range m.NodeSubnets() {
			subnetName = subnet.Name
		}
		if len(m.NodeSubnets()) == 0 || len(m.NodeSubnets()) > 1 || subnetName == "" {
			return errors.New("a subnet name must be specified when no subnets are specified or more than 1 subnet of role 'node' exist")
		}

		m.AzureMachinePool.Spec.Template.NetworkInterfaces[0].SubnetName = subnetName
	}

	return nil
}

// UpdateDeleteStatus updates a condition on the AzureMachinePool status after a DELETE operation.
func (m *MachinePoolScope) UpdateDeleteStatus(condition clusterv1.ConditionType, service string, err error) {
	switch {
	case err == nil:
		conditions.MarkFalse(m.AzureMachinePool, condition, infrav1.DeletedReason, clusterv1.ConditionSeverityInfo, "%s successfully deleted", service)
	case azure.IsOperationNotDoneError(err):
		conditions.MarkFalse(m.AzureMachinePool, condition, infrav1.DeletingReason, clusterv1.ConditionSeverityInfo, "%s deleting", service)
	default:
		conditions.MarkFalse(m.AzureMachinePool, condition, infrav1.DeletionFailedReason, clusterv1.ConditionSeverityError, "%s failed to delete. err: %s", service, err.Error())
	}
}

// UpdatePutStatus updates a condition on the AzureMachinePool status after a PUT operation.
func (m *MachinePoolScope) UpdatePutStatus(condition clusterv1.ConditionType, service string, err error) {
	switch {
	case err == nil:
		conditions.MarkTrue(m.AzureMachinePool, condition)
	case azure.IsOperationNotDoneError(err):
		conditions.MarkFalse(m.AzureMachinePool, condition, infrav1.CreatingReason, clusterv1.ConditionSeverityInfo, "%s creating or updating", service)
	default:
		conditions.MarkFalse(m.AzureMachinePool, condition, infrav1.FailedReason, clusterv1.ConditionSeverityError, "%s failed to create or update. err: %s", service, err.Error())
	}
}

// UpdatePatchStatus updates a condition on the AzureMachinePool status after a PATCH operation.
func (m *MachinePoolScope) UpdatePatchStatus(condition clusterv1.ConditionType, service string, err error) {
	switch {
	case err == nil:
		conditions.MarkTrue(m.AzureMachinePool, condition)
	case azure.IsOperationNotDoneError(err):
		conditions.MarkFalse(m.AzureMachinePool, condition, infrav1.UpdatingReason, clusterv1.ConditionSeverityInfo, "%s updating", service)
	default:
		conditions.MarkFalse(m.AzureMachinePool, condition, infrav1.FailedReason, clusterv1.ConditionSeverityError, "%s failed to update. err: %s", service, err.Error())
	}
}

// PatchCAPIMachinePoolObject persists the capi machinepool configuration and status.
func (m *MachinePoolScope) PatchCAPIMachinePoolObject(ctx context.Context) error {
	return m.capiMachinePoolPatchHelper.Patch(
		ctx,
		m.MachinePool,
	)
}

// UpdateCAPIMachinePoolReplicas updates the associated MachinePool replica count.
func (m *MachinePoolScope) UpdateCAPIMachinePoolReplicas(_ context.Context, replicas *int32) {
	m.MachinePool.Spec.Replicas = replicas
}

// HasReplicasExternallyManaged returns true if the externally managed annotation is set on the CAPI MachinePool resource.
func (m *MachinePoolScope) HasReplicasExternallyManaged(_ context.Context) bool {
	return annotations.ReplicasManagedByExternalAutoscaler(m.MachinePool)
}

// ReconcileReplicas ensures MachinePool replicas match VMSS capacity if replicas are externally managed by an autoscaler.
func (m *MachinePoolScope) ReconcileReplicas(ctx context.Context, vmss *azure.VMSS) error {
	if !m.HasReplicasExternallyManaged(ctx) {
		return nil
	}

	var replicas int32
	if m.MachinePool.Spec.Replicas != nil {
		replicas = *m.MachinePool.Spec.Replicas
	}

	if capacity := int32(vmss.Capacity); capacity != replicas {
		m.UpdateCAPIMachinePoolReplicas(ctx, &capacity)
	}

	return nil
}

// AnnotationJSON returns a map[string]interface from a JSON annotation.
func (m *MachinePoolScope) AnnotationJSON(annotation string) (map[string]interface{}, error) {
	out := map[string]interface{}{}
	jsonAnnotation := m.AzureMachinePool.GetAnnotations()[annotation]
	if jsonAnnotation == "" {
		return out, nil
	}
	err := json.Unmarshal([]byte(jsonAnnotation), &out)
	if err != nil {
		return out, err
	}
	return out, nil
}

// UpdateAnnotationJSON updates the `annotation` with
// `content`. `content` in this case should be a `map[string]interface{}`
// suitable for turning into JSON. This `content` map will be marshalled into a
// JSON string before being set as the given `annotation`.
func (m *MachinePoolScope) UpdateAnnotationJSON(annotation string, content map[string]interface{}) error {
	b, err := json.Marshal(content)
	if err != nil {
		return err
	}
	m.SetAnnotation(annotation, string(b))
	return nil
}

// TagsSpecs returns the tags for the AzureMachinePool.
func (m *MachinePoolScope) TagsSpecs() []azure.TagsSpec {
	return []azure.TagsSpec{
		{
			Scope:      azure.VMSSID(m.SubscriptionID(), m.NodeResourceGroup(), m.Name()),
			Tags:       m.AdditionalTags(),
			Annotation: azure.VMSSTagsLastAppliedAnnotation,
		},
	}
}
