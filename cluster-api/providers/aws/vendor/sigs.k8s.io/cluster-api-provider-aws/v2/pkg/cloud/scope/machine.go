/*
Copyright 2018 The Kubernetes Authors.

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
	"encoding/base64"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
)

// MachineScopeParams defines the input parameters used to create a new MachineScope.
type MachineScopeParams struct {
	Client       client.Client
	Logger       *logger.Logger
	Cluster      *clusterv1.Cluster
	Machine      *clusterv1.Machine
	InfraCluster EC2Scope
	AWSMachine   *infrav1.AWSMachine
}

// NewMachineScope creates a new MachineScope from the supplied parameters.
// This is meant to be called for each reconcile iteration.
func NewMachineScope(params MachineScopeParams) (*MachineScope, error) {
	if params.Client == nil {
		return nil, errors.New("client is required when creating a MachineScope")
	}
	if params.Machine == nil {
		return nil, errors.New("machine is required when creating a MachineScope")
	}
	if params.Cluster == nil {
		return nil, errors.New("cluster is required when creating a MachineScope")
	}
	if params.AWSMachine == nil {
		return nil, errors.New("aws machine is required when creating a MachineScope")
	}
	if params.InfraCluster == nil {
		return nil, errors.New("aws cluster is required when creating a MachineScope")
	}

	if params.Logger == nil {
		log := klog.Background()
		params.Logger = logger.NewLogger(log)
	}

	helper, err := patch.NewHelper(params.AWSMachine, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}
	return &MachineScope{
		Logger:       *params.Logger,
		client:       params.Client,
		patchHelper:  helper,
		Cluster:      params.Cluster,
		Machine:      params.Machine,
		InfraCluster: params.InfraCluster,
		AWSMachine:   params.AWSMachine,
	}, nil
}

// MachineScope defines a scope defined around a machine and its cluster.
type MachineScope struct {
	logger.Logger
	client      client.Client
	patchHelper *patch.Helper

	Cluster      *clusterv1.Cluster
	Machine      *clusterv1.Machine
	InfraCluster EC2Scope
	AWSMachine   *infrav1.AWSMachine
}

// Name returns the AWSMachine name.
func (m *MachineScope) Name() string {
	return m.AWSMachine.Name
}

// Namespace returns the namespace name.
func (m *MachineScope) Namespace() string {
	return m.AWSMachine.Namespace
}

// IsControlPlane returns true if the machine is a control plane.
func (m *MachineScope) IsControlPlane() bool {
	return util.IsControlPlaneMachine(m.Machine)
}

// IsMachinePoolMachine returns true if the machine is created for a machinepool.
func (m *MachineScope) IsMachinePoolMachine() bool {
	if _, ok := m.Machine.GetLabels()[clusterv1.MachinePoolNameLabel]; ok {
		return true
	}
	for _, owner := range m.Machine.OwnerReferences {
		if owner.Kind == v1beta2.KindMachinePool {
			return true
		}
	}
	return false
}

// Role returns the machine role from the labels.
func (m *MachineScope) Role() string {
	if util.IsControlPlaneMachine(m.Machine) {
		return "control-plane"
	}
	return "node"
}

// GetInstanceID returns the AWSMachine instance id by parsing Spec.ProviderID.
func (m *MachineScope) GetInstanceID() *string {
	parsed, err := NewProviderID(m.GetProviderID())
	if err != nil {
		return nil
	}
	return ptr.To[string](parsed.ID())
}

// GetProviderID returns the AWSMachine providerID from the spec.
func (m *MachineScope) GetProviderID() string {
	if m.AWSMachine.Spec.ProviderID != nil {
		return *m.AWSMachine.Spec.ProviderID
	}
	return ""
}

// SetProviderID sets the AWSMachine providerID in spec.
func (m *MachineScope) SetProviderID(instanceID, availabilityZone string) {
	providerID := GenerateProviderID(availabilityZone, instanceID)
	m.AWSMachine.Spec.ProviderID = ptr.To[string](providerID)
}

// SetInstanceID sets the AWSMachine instanceID in spec.
func (m *MachineScope) SetInstanceID(instanceID string) {
	m.AWSMachine.Spec.InstanceID = ptr.To[string](instanceID)
}

// GetInstanceState returns the AWSMachine instance state from the status.
func (m *MachineScope) GetInstanceState() *infrav1.InstanceState {
	return m.AWSMachine.Status.InstanceState
}

// SetInstanceState sets the AWSMachine status instance state.
func (m *MachineScope) SetInstanceState(v infrav1.InstanceState) {
	m.AWSMachine.Status.InstanceState = &v
}

// SetReady sets the AWSMachine Ready Status.
func (m *MachineScope) SetReady() {
	m.AWSMachine.Status.Ready = true
}

// SetNotReady sets the AWSMachine Ready Status to false.
func (m *MachineScope) SetNotReady() {
	m.AWSMachine.Status.Ready = false
}

// SetFailureMessage sets the AWSMachine status failure message.
func (m *MachineScope) SetFailureMessage(v error) {
	m.AWSMachine.Status.FailureMessage = ptr.To[string](v.Error())
}

// SetFailureReason sets the AWSMachine status failure reason.
func (m *MachineScope) SetFailureReason(v string) {
	m.AWSMachine.Status.FailureReason = &v
}

// SetAnnotation sets a key value annotation on the AWSMachine.
func (m *MachineScope) SetAnnotation(key, value string) {
	if m.AWSMachine.Annotations == nil {
		m.AWSMachine.Annotations = map[string]string{}
	}
	m.AWSMachine.Annotations[key] = value
}

// UseSecretsManager returns the computed value of whether or not
// userdata should be stored using AWS Secrets Manager.
func (m *MachineScope) UseSecretsManager(userDataFormat string) bool {
	return !m.AWSMachine.Spec.CloudInit.InsecureSkipSecretsManager && !m.UseIgnition(userDataFormat)
}

// UseIgnition returns true if the AWSMachine should use Ignition.
func (m *MachineScope) UseIgnition(userDataFormat string) bool {
	return userDataFormat == "ignition" || (m.AWSMachine.Spec.Ignition != nil)
}

// SecureSecretsBackend returns the chosen secret backend.
func (m *MachineScope) SecureSecretsBackend() infrav1.SecretBackend {
	return m.AWSMachine.Spec.CloudInit.SecureSecretsBackend
}

// CompressUserData returns the computed value of whether or not
// userdata should be compressed using gzip.
func (m *MachineScope) CompressUserData(userDataFormat string) bool {
	if m.UseIgnition(userDataFormat) {
		return false
	}

	return m.AWSMachine.Spec.UncompressedUserData != nil && !*m.AWSMachine.Spec.UncompressedUserData
}

// GetSecretPrefix returns the prefix for the secrets belonging
// to the AWSMachine in AWS Secrets Manager.
func (m *MachineScope) GetSecretPrefix() string {
	return m.AWSMachine.Spec.CloudInit.SecretPrefix
}

// SetSecretPrefix sets the prefix for the secrets belonging
// to the AWSMachine in AWS Secrets Manager.
func (m *MachineScope) SetSecretPrefix(value string) {
	m.AWSMachine.Spec.CloudInit.SecretPrefix = value
}

// DeleteSecretPrefix deletes the prefix for the secret belonging
// to the AWSMachine in AWS Secrets Manager.
func (m *MachineScope) DeleteSecretPrefix() {
	m.AWSMachine.Spec.CloudInit.SecretPrefix = ""
}

// GetSecretCount returns the number of AWS Secret Manager entries making up
// the complete userdata.
func (m *MachineScope) GetSecretCount() int32 {
	return m.AWSMachine.Spec.CloudInit.SecretCount
}

// SetSecretCount sets the number of AWS Secret Manager entries making up
// the complete userdata.
func (m *MachineScope) SetSecretCount(i int32) {
	m.AWSMachine.Spec.CloudInit.SecretCount = i
}

// SetAddresses sets the AWSMachine address status.
func (m *MachineScope) SetAddresses(addrs []clusterv1.MachineAddress) {
	m.AWSMachine.Status.Addresses = addrs
}

// GetBootstrapData returns the bootstrap data from the secret in the Machine's bootstrap.dataSecretName as base64.
func (m *MachineScope) GetBootstrapData() (string, error) {
	value, err := m.GetRawBootstrapData()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(value), nil
}

// GetRawBootstrapData returns the bootstrap data from the secret in the Machine's bootstrap.dataSecretName.
func (m *MachineScope) GetRawBootstrapData() ([]byte, error) {
	data, _, err := m.GetRawBootstrapDataWithFormat()

	return data, err
}

// GetRawBootstrapDataWithFormat returns the bootstrap data from the secret in the Machine's bootstrap.dataSecretName.
func (m *MachineScope) GetRawBootstrapDataWithFormat() ([]byte, string, error) {
	if m.Machine.Spec.Bootstrap.DataSecretName == nil {
		return nil, "", errors.New("error retrieving bootstrap data: linked Machine's bootstrap.dataSecretName is nil")
	}

	secret := &corev1.Secret{}
	key := types.NamespacedName{Namespace: m.Namespace(), Name: *m.Machine.Spec.Bootstrap.DataSecretName}
	if err := m.client.Get(context.TODO(), key, secret); err != nil {
		return nil, "", errors.Wrapf(err, "failed to retrieve bootstrap data secret for AWSMachine %s/%s", m.Namespace(), m.Name())
	}

	value, ok := secret.Data["value"]
	if !ok {
		return nil, "", errors.New("error retrieving bootstrap data: secret value key is missing")
	}

	return value, string(secret.Data["format"]), nil
}

// PatchObject persists the machine spec and status.
func (m *MachineScope) PatchObject() error {
	// Always update the readyCondition by summarizing the state of other conditions.
	// A step counter is added to represent progress during the provisioning process (instead we are hiding during the deletion process).
	applicableConditions := []clusterv1.ConditionType{
		infrav1.InstanceReadyCondition,
		infrav1.SecurityGroupsReadyCondition,
	}

	if m.IsControlPlane() {
		applicableConditions = append(applicableConditions, infrav1.ELBAttachedCondition)
	}

	conditions.SetSummary(m.AWSMachine,
		conditions.WithConditions(applicableConditions...),
		conditions.WithStepCounterIf(m.AWSMachine.ObjectMeta.DeletionTimestamp.IsZero()),
		conditions.WithStepCounter(),
	)

	return m.patchHelper.Patch(
		context.TODO(),
		m.AWSMachine,
		patch.WithOwnedConditions{Conditions: []clusterv1.ConditionType{
			clusterv1.ReadyCondition,
			infrav1.InstanceReadyCondition,
			infrav1.SecurityGroupsReadyCondition,
			infrav1.ELBAttachedCondition,
		}})
}

// Close the MachineScope by updating the machine spec, machine status.
func (m *MachineScope) Close() error {
	return m.PatchObject()
}

// AdditionalTags merges AdditionalTags from the scope's AWSCluster and AWSMachine. If the same key is present in both,
// the value from AWSMachine takes precedence. The returned Tags will never be nil.
func (m *MachineScope) AdditionalTags() infrav1.Tags {
	tags := make(infrav1.Tags)

	// Start with the cluster-wide tags...
	tags.Merge(m.InfraCluster.AdditionalTags())
	// ... and merge in the Machine's
	tags.Merge(m.AWSMachine.Spec.AdditionalTags)

	return tags
}

// HasFailed returns the failure state of the machine scope.
func (m *MachineScope) HasFailed() bool {
	return m.AWSMachine.Status.FailureReason != nil || m.AWSMachine.Status.FailureMessage != nil
}

// InstanceIsRunning returns the instance state of the machine scope.
func (m *MachineScope) InstanceIsRunning() bool {
	state := m.GetInstanceState()
	return state != nil && infrav1.InstanceRunningStates.Has(string(*state))
}

// InstanceIsOperational returns the operational state of the machine scope.
func (m *MachineScope) InstanceIsOperational() bool {
	state := m.GetInstanceState()
	return state != nil && infrav1.InstanceOperationalStates.Has(string(*state))
}

// InstanceIsInKnownState checks if the machine scope's instance state is known.
func (m *MachineScope) InstanceIsInKnownState() bool {
	state := m.GetInstanceState()
	return state != nil && infrav1.InstanceKnownStates.Has(string(*state))
}

// AWSMachineIsDeleted checks if the AWS machine was deleted.
func (m *MachineScope) AWSMachineIsDeleted() bool {
	return !m.AWSMachine.ObjectMeta.DeletionTimestamp.IsZero()
}

// MachineIsDeleted checks if the machine was deleted.
func (m *MachineScope) MachineIsDeleted() bool {
	return !m.Machine.ObjectMeta.DeletionTimestamp.IsZero()
}

// IsEKSManaged checks if the machine is EKS managed.
func (m *MachineScope) IsEKSManaged() bool {
	return m.InfraCluster.InfraCluster().GetObjectKind().GroupVersionKind().Kind == ekscontrolplanev1.AWSManagedControlPlaneKind
}

// IsControlPlaneExternallyManaged checks if the control plane is externally managed.
//
// This is determined by the kind of the control plane object (EKS for example),
// or if the control plane referenced object is reporting as externally managed.
func (m *MachineScope) IsControlPlaneExternallyManaged() bool {
	if m.IsEKSManaged() {
		return true
	}

	// Check if the control plane is externally managed.
	u, err := m.InfraCluster.UnstructuredControlPlane()
	if err != nil {
		m.Error(err, "failed to get unstructured control plane")
		return false
	}
	return util.IsExternalManagedControlPlane(u)
}

// IsExternallyManaged checks if the machine is externally managed.
func (m *MachineScope) IsExternallyManaged() bool {
	return annotations.IsExternallyManaged(m.InfraCluster.InfraCluster())
}

// SetInterruptible sets the AWSMachine status Interruptible.
func (m *MachineScope) SetInterruptible() {
	if m.AWSMachine.Spec.SpotMarketOptions != nil {
		m.AWSMachine.Status.Interruptible = true
	}
}

// GetElasticIPPool returns the Elastic IP Pool for an machine, when exists.
func (m *MachineScope) GetElasticIPPool() *infrav1.ElasticIPPool {
	if m.AWSMachine == nil {
		return nil
	}
	return m.AWSMachine.Spec.ElasticIPPool
}
