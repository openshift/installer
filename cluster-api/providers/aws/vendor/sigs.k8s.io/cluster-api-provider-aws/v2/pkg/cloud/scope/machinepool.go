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
	"fmt"
	"strings"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/controllers/remote"
	capierrors "sigs.k8s.io/cluster-api/errors"
	expclusterv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
)

const (
	// ReplicasManagedByAnnotation is an annotation that indicates external (non-Cluster API) management of infra scaling.
	// The practical effect of this is that the capi "replica" count is derived from the number of observed infra machines,
	// instead of being a source of truth for eventual consistency.
	//
	// N.B. this is to be replaced by a direct reference to CAPI once https://github.com/kubernetes-sigs/cluster-api/pull/7107 is meged.
	ReplicasManagedByAnnotation = "cluster.x-k8s.io/replicas-managed-by"

	// ExternalAutoscalerReplicasManagedByAnnotationValue is used with the "cluster.x-k8s.io/replicas-managed-by" annotation
	// to indicate an external autoscaler enforces replica count.
	//
	// N.B. this is to be replaced by a direct reference to CAPI once https://github.com/kubernetes-sigs/cluster-api/pull/7107 is meged.
	ExternalAutoscalerReplicasManagedByAnnotationValue = "external-autoscaler"
)

// MachinePoolScope defines a scope defined around a machine and its cluster.
type MachinePoolScope struct {
	logger.Logger
	client.Client
	patchHelper                *patch.Helper
	capiMachinePoolPatchHelper *patch.Helper

	Cluster        *clusterv1.Cluster
	MachinePool    *expclusterv1.MachinePool
	InfraCluster   EC2Scope
	AWSMachinePool *expinfrav1.AWSMachinePool
}

// MachinePoolScopeParams defines a scope defined around a machine and its cluster.
type MachinePoolScopeParams struct {
	client.Client
	Logger *logger.Logger

	Cluster        *clusterv1.Cluster
	MachinePool    *expclusterv1.MachinePool
	InfraCluster   EC2Scope
	AWSMachinePool *expinfrav1.AWSMachinePool
}

// GetProviderID returns the AWSMachine providerID from the spec.
func (m *MachinePoolScope) GetProviderID() string {
	if m.AWSMachinePool.Spec.ProviderID != "" {
		return m.AWSMachinePool.Spec.ProviderID
	}
	return ""
}

// NewMachinePoolScope creates a new MachinePoolScope from the supplied parameters.
// This is meant to be called for each reconcile iteration.
func NewMachinePoolScope(params MachinePoolScopeParams) (*MachinePoolScope, error) {
	if params.Client == nil {
		return nil, errors.New("client is required when creating a MachinePoolScope")
	}
	if params.MachinePool == nil {
		return nil, errors.New("machinepool is required when creating a MachinePoolScope")
	}
	if params.Cluster == nil {
		return nil, errors.New("cluster is required when creating a MachinePoolScope")
	}
	if params.AWSMachinePool == nil {
		return nil, errors.New("aws machine pool is required when creating a MachinePoolScope")
	}
	if params.InfraCluster == nil {
		return nil, errors.New("aws cluster is required when creating a MachinePoolScope")
	}

	if params.Logger == nil {
		log := klog.Background()
		params.Logger = logger.NewLogger(log)
	}

	ampHelper, err := patch.NewHelper(params.AWSMachinePool, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init AWSMachinePool patch helper")
	}
	mpHelper, err := patch.NewHelper(params.MachinePool, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init MachinePool patch helper")
	}

	return &MachinePoolScope{
		Logger:                     *params.Logger,
		Client:                     params.Client,
		patchHelper:                ampHelper,
		capiMachinePoolPatchHelper: mpHelper,

		Cluster:        params.Cluster,
		MachinePool:    params.MachinePool,
		InfraCluster:   params.InfraCluster,
		AWSMachinePool: params.AWSMachinePool,
	}, nil
}

// Name returns the AWSMachinePool name.
func (m *MachinePoolScope) Name() string {
	return m.AWSMachinePool.Name
}

// Namespace returns the namespace name.
func (m *MachinePoolScope) Namespace() string {
	return m.AWSMachinePool.Namespace
}

// GetRawBootstrapData returns the bootstrap data from the secret in the Machine's bootstrap.dataSecretName.
// todo(rudoi): stolen from MachinePool - any way to reuse?
func (m *MachinePoolScope) GetRawBootstrapData() ([]byte, error) {
	data, _, err := m.getBootstrapData()

	return data, err
}

func (m *MachinePoolScope) GetRawBootstrapDataWithFormat() ([]byte, string, error) {
	return m.getBootstrapData()
}

func (m *MachinePoolScope) getBootstrapData() ([]byte, string, error) {
	if m.MachinePool.Spec.Template.Spec.Bootstrap.DataSecretName == nil {
		return nil, "", errors.New("error retrieving bootstrap data: linked Machine's bootstrap.dataSecretName is nil")
	}

	secret := &corev1.Secret{}
	key := types.NamespacedName{Namespace: m.Namespace(), Name: *m.MachinePool.Spec.Template.Spec.Bootstrap.DataSecretName}

	if err := m.Client.Get(context.TODO(), key, secret); err != nil {
		return nil, "", errors.Wrapf(err, "failed to retrieve bootstrap data secret for AWSMachine %s/%s", m.Namespace(), m.Name())
	}

	value, ok := secret.Data["value"]
	if !ok {
		return nil, "", errors.New("error retrieving bootstrap data: secret value key is missing")
	}

	return value, string(secret.Data["format"]), nil
}

// AdditionalTags merges AdditionalTags from the scope's AWSCluster and AWSMachinePool. If the same key is present in both,
// the value from AWSMachinePool takes precedence. The returned Tags will never be nil.
func (m *MachinePoolScope) AdditionalTags() infrav1.Tags {
	tags := make(infrav1.Tags)

	// Start with the cluster-wide tags...
	tags.Merge(m.InfraCluster.AdditionalTags())
	// ... and merge in the Machine's
	tags.Merge(m.AWSMachinePool.Spec.AdditionalTags)

	return tags
}

// PatchObject persists the machinepool spec and status.
func (m *MachinePoolScope) PatchObject() error {
	return m.patchHelper.Patch(
		context.TODO(),
		m.AWSMachinePool,
		patch.WithOwnedConditions{Conditions: []clusterv1.ConditionType{
			expinfrav1.ASGReadyCondition,
			expinfrav1.LaunchTemplateReadyCondition,
		}})
}

// PatchCAPIMachinePoolObject persists the capi machinepool configuration and status.
func (m *MachinePoolScope) PatchCAPIMachinePoolObject(ctx context.Context) error {
	return m.capiMachinePoolPatchHelper.Patch(
		ctx,
		m.MachinePool,
	)
}

// Close the MachinePoolScope by updating the machinepool spec, machine status.
func (m *MachinePoolScope) Close() error {
	return m.PatchObject()
}

// SetAnnotation sets a key value annotation on the AWSMachine.
func (m *MachinePoolScope) SetAnnotation(key, value string) {
	if m.AWSMachinePool.Annotations == nil {
		m.AWSMachinePool.Annotations = map[string]string{}
	}
	m.AWSMachinePool.Annotations[key] = value
}

// SetFailureMessage sets the AWSMachine status failure message.
func (m *MachinePoolScope) SetFailureMessage(v error) {
	m.AWSMachinePool.Status.FailureMessage = pointer.String(v.Error())
}

// SetFailureReason sets the AWSMachine status failure reason.
func (m *MachinePoolScope) SetFailureReason(v capierrors.MachineStatusError) {
	m.AWSMachinePool.Status.FailureReason = &v
}

// HasFailed returns true when the AWSMachinePool's Failure reason or Failure message is populated.
func (m *MachinePoolScope) HasFailed() bool {
	return m.AWSMachinePool.Status.FailureReason != nil || m.AWSMachinePool.Status.FailureMessage != nil
}

// SetNotReady sets the AWSMachinePool Ready Status to false.
func (m *MachinePoolScope) SetNotReady() {
	m.AWSMachinePool.Status.Ready = false
}

// GetASGStatus returns the AWSMachinePool instance state from the status.
func (m *MachinePoolScope) GetASGStatus() *expinfrav1.ASGStatus {
	return m.AWSMachinePool.Status.ASGStatus
}

// SetASGStatus sets the AWSMachinePool status instance state.
func (m *MachinePoolScope) SetASGStatus(v expinfrav1.ASGStatus) {
	m.AWSMachinePool.Status.ASGStatus = &v
}

func (m *MachinePoolScope) GetObjectMeta() *metav1.ObjectMeta {
	return &m.AWSMachinePool.ObjectMeta
}

func (m *MachinePoolScope) GetSetter() conditions.Setter {
	return m.AWSMachinePool
}

func (m *MachinePoolScope) GetEC2Scope() EC2Scope {
	return m.InfraCluster
}

func (m *MachinePoolScope) GetLaunchTemplateIDStatus() string {
	return m.AWSMachinePool.Status.LaunchTemplateID
}

func (m *MachinePoolScope) SetLaunchTemplateIDStatus(id string) {
	m.AWSMachinePool.Status.LaunchTemplateID = id
}

func (m *MachinePoolScope) GetLaunchTemplateLatestVersionStatus() string {
	if m.AWSMachinePool.Status.LaunchTemplateVersion != nil {
		return *m.AWSMachinePool.Status.LaunchTemplateVersion
	} else {
		return ""
	}
}

func (m *MachinePoolScope) SetLaunchTemplateLatestVersionStatus(version string) {
	m.AWSMachinePool.Status.LaunchTemplateVersion = &version
}

// IsEKSManaged checks if the AWSMachinePool is EKS managed.
func (m *MachinePoolScope) IsEKSManaged() bool {
	return m.InfraCluster.InfraCluster().GetObjectKind().GroupVersionKind().Kind == ekscontrolplanev1.AWSManagedControlPlaneKind
}

// SubnetIDs returns the machine pool subnet IDs.
func (m *MachinePoolScope) SubnetIDs(subnetIDs []string) ([]string, error) {
	strategy, err := newDefaultSubnetPlacementStrategy(&m.Logger)
	if err != nil {
		return subnetIDs, fmt.Errorf("getting subnet placement strategy: %w", err)
	}

	return strategy.Place(&placementInput{
		SpecSubnetIDs:           subnetIDs,
		SpecAvailabilityZones:   m.AWSMachinePool.Spec.AvailabilityZones,
		ParentAvailabilityZones: m.MachinePool.Spec.FailureDomains,
		ControlplaneSubnets:     m.InfraCluster.Subnets(),
	})
}

// NodeStatus represents the status of a Kubernetes node.
type NodeStatus struct {
	Ready   bool
	Version string
}

// UpdateInstanceStatuses ties ASG instances and Node status data together and updates AWSMachinePool
// This updates if ASG instances ready and kubelet version running on the node..
func (m *MachinePoolScope) UpdateInstanceStatuses(ctx context.Context, instances []infrav1.Instance) error {
	providerIDs := make([]string, len(instances))
	for i, instance := range instances {
		providerIDs[i] = fmt.Sprintf("aws:////%s", instance.ID)
	}

	nodeStatusByProviderID, err := m.getNodeStatusByProviderID(ctx, providerIDs)
	if err != nil {
		return errors.Wrap(err, "failed to get node status by provider id")
	}

	var readyReplicas int32
	instanceStatuses := make([]expinfrav1.AWSMachinePoolInstanceStatus, len(instances))
	for i, instance := range instances {
		instanceStatuses[i] = expinfrav1.AWSMachinePoolInstanceStatus{
			InstanceID: instance.ID,
		}

		instanceStatus := instanceStatuses[i]
		if nodeStatus, ok := nodeStatusByProviderID[fmt.Sprintf("aws:////%s", instanceStatus.InstanceID)]; ok {
			instanceStatus.Version = &nodeStatus.Version
			if nodeStatus.Ready {
				readyReplicas++
			}
		}
	}

	// TODO: readyReplicas can be used as status.replicas but this will delay machinepool to become ready. next reconcile updates this.
	m.AWSMachinePool.Status.Instances = instanceStatuses
	return nil
}

func (m *MachinePoolScope) getNodeStatusByProviderID(ctx context.Context, providerIDList []string) (map[string]*NodeStatus, error) {
	nodeStatusMap := map[string]*NodeStatus{}
	for _, id := range providerIDList {
		nodeStatusMap[id] = &NodeStatus{}
	}

	workloadClient, err := remote.NewClusterClient(ctx, "", m.Client, util.ObjectKey(m.Cluster))
	if err != nil {
		return nil, err
	}

	nodeList := corev1.NodeList{}
	for {
		if err := workloadClient.List(ctx, &nodeList, client.Continue(nodeList.Continue)); err != nil {
			return nil, errors.Wrapf(err, "failed to List nodes")
		}

		for _, node := range nodeList.Items {
			strList := strings.Split(node.Spec.ProviderID, "/")

			if status, ok := nodeStatusMap[fmt.Sprintf("aws:////%s", strList[len(strList)-1])]; ok {
				status.Ready = nodeIsReady(node)
				status.Version = node.Status.NodeInfo.KubeletVersion
			}
		}

		if nodeList.Continue == "" {
			break
		}
	}

	return nodeStatusMap, nil
}

func nodeIsReady(node corev1.Node) bool {
	for _, n := range node.Status.Conditions {
		if n.Type == corev1.NodeReady {
			return n.Status == corev1.ConditionTrue
		}
	}
	return false
}

func (m *MachinePoolScope) GetLaunchTemplate() *expinfrav1.AWSLaunchTemplate {
	return &m.AWSMachinePool.Spec.AWSLaunchTemplate
}

func (m *MachinePoolScope) GetMachinePool() *expclusterv1.MachinePool {
	return m.MachinePool
}

func (m *MachinePoolScope) LaunchTemplateName() string {
	return m.Name()
}

func (m *MachinePoolScope) GetRuntimeObject() runtime.Object {
	return m.AWSMachinePool
}

func ReplicasExternallyManaged(mp *expclusterv1.MachinePool) bool {
	val, ok := mp.Annotations[ReplicasManagedByAnnotation]
	return ok && val == ExternalAutoscalerReplicasManagedByAnnotationValue
}
