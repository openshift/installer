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

// Package scope implements scope types.
package scope

import (
	"context"
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/go-logr/logr"

	"github.com/pkg/errors"
	"golang.org/x/mod/semver"
	"google.golang.org/api/compute/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	infrav1 "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud/providerid"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud/services/shared"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// MachineScopeParams defines the input parameters used to create a new MachineScope.
type MachineScopeParams struct {
	Client        client.Client
	ClusterGetter cloud.ClusterGetter
	Machine       *clusterv1.Machine
	GCPMachine    *infrav1.GCPMachine
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
	if params.GCPMachine == nil {
		return nil, errors.New("gcp machine is required when creating a MachineScope")
	}

	helper, err := patch.NewHelper(params.GCPMachine, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}

	return &MachineScope{
		client:        params.Client,
		Machine:       params.Machine,
		GCPMachine:    params.GCPMachine,
		ClusterGetter: params.ClusterGetter,
		patchHelper:   helper,
	}, nil
}

// MachineScope defines a scope defined around a machine and its cluster.
type MachineScope struct {
	client        client.Client
	patchHelper   *patch.Helper
	ClusterGetter cloud.ClusterGetter
	Machine       *clusterv1.Machine
	GCPMachine    *infrav1.GCPMachine
}

// ANCHOR: MachineGetter

// Cloud returns initialized cloud.
func (m *MachineScope) Cloud() cloud.Cloud {
	return m.ClusterGetter.Cloud()
}

// NetworkCloud returns initialized network cloud.
func (m *MachineScope) NetworkCloud() cloud.Cloud {
	return m.ClusterGetter.NetworkCloud()
}

// Zone returns the FailureDomain for the GCPMachine.
func (m *MachineScope) Zone() string {
	if m.Machine.Spec.FailureDomain == nil {
		fd := m.ClusterGetter.FailureDomains()
		if len(fd) == 0 {
			return ""
		}
		zones := make([]string, 0, len(fd))
		for zone := range fd {
			zones = append(zones, zone)
		}
		sort.Strings(zones)
		return zones[0]
	}
	return *m.Machine.Spec.FailureDomain
}

// Project return the project for the GCPMachine's cluster.
func (m *MachineScope) Project() string {
	return m.ClusterGetter.Project()
}

// Name returns the GCPMachine name.
func (m *MachineScope) Name() string {
	return m.GCPMachine.Name
}

// Namespace returns the namespace name.
func (m *MachineScope) Namespace() string {
	return m.GCPMachine.Namespace
}

// ControlPlaneGroupName returns the control-plane instance group name.
func (m *MachineScope) ControlPlaneGroupName() string {
	tag := ptr.Deref(m.ClusterGetter.LoadBalancer().APIServerInstanceGroupTagOverride, infrav1.APIServerRoleTagValue)
	return fmt.Sprintf("%s-%s-%s", m.ClusterGetter.Name(), tag, m.Zone())
}

// IsControlPlane returns true if the machine is a control plane.
func (m *MachineScope) IsControlPlane() bool {
	return IsControlPlaneMachine(m.Machine)
}

// Role returns the machine role from the labels.
func (m *MachineScope) Role() string {
	if IsControlPlaneMachine(m.Machine) {
		return "control-plane"
	}

	return "node"
}

// IsControlPlaneMachine checks machine is a control plane node.
func IsControlPlaneMachine(machine *clusterv1.Machine) bool {
	_, ok := machine.Labels[clusterv1.MachineControlPlaneLabel]
	return ok
}

// GetInstanceID returns the GCPMachine instance id by parsing Spec.ProviderID.
func (m *MachineScope) GetInstanceID() *string {
	parsed, err := NewProviderID(m.GetProviderID())
	if err != nil {
		return nil
	}

	return ptr.To[string](parsed.ID())
}

// GetProviderID returns the GCPMachine providerID from the spec.
func (m *MachineScope) GetProviderID() string {
	if m.GCPMachine.Spec.ProviderID != nil {
		return *m.GCPMachine.Spec.ProviderID
	}

	return ""
}

// ANCHOR_END: MachineGetter

// ANCHOR: MachineSetter

// SetProviderID sets the GCPMachine providerID in spec.
func (m *MachineScope) SetProviderID() {
	providerID, _ := providerid.New(m.ClusterGetter.Project(), m.Zone(), m.Name())
	m.GCPMachine.Spec.ProviderID = ptr.To[string](providerID.String())
}

// GetInstanceStatus returns the GCPMachine instance status.
func (m *MachineScope) GetInstanceStatus() *infrav1.InstanceStatus {
	return m.GCPMachine.Status.InstanceStatus
}

// SetInstanceStatus sets the GCPMachine instance status.
func (m *MachineScope) SetInstanceStatus(v infrav1.InstanceStatus) {
	m.GCPMachine.Status.InstanceStatus = &v
}

// SetReady sets the GCPMachine Ready Status.
func (m *MachineScope) SetReady() {
	m.GCPMachine.Status.Ready = true
}

// SetFailureMessage sets the GCPMachine status failure message.
func (m *MachineScope) SetFailureMessage(v error) {
	m.GCPMachine.Status.FailureMessage = ptr.To[string](v.Error())
}

// SetFailureReason sets the GCPMachine status failure reason.
func (m *MachineScope) SetFailureReason(v string) {
	m.GCPMachine.Status.FailureReason = &v
}

// SetAnnotation sets a key value annotation on the GCPMachine.
func (m *MachineScope) SetAnnotation(key, value string) {
	if m.GCPMachine.Annotations == nil {
		m.GCPMachine.Annotations = map[string]string{}
	}
	m.GCPMachine.Annotations[key] = value
}

// SetAddresses sets the addresses field on the GCPMachine.
func (m *MachineScope) SetAddresses(addressList []corev1.NodeAddress) {
	m.GCPMachine.Status.Addresses = addressList
}

// ANCHOR_END: MachineSetter

// ANCHOR: MachineInstanceSpec

// InstanceImageSpec returns compute instance image attched-disk spec.
func (m *MachineScope) InstanceImageSpec() *compute.AttachedDisk {
	version := ""
	if m.Machine.Spec.Version != nil {
		version = *m.Machine.Spec.Version
	}
	image := "capi-ubuntu-1804-k8s-" + strings.ReplaceAll(semver.MajorMinor(version), ".", "-")
	sourceImage := path.Join("projects", m.ClusterGetter.Project(), "global", "images", "family", image)
	if m.GCPMachine.Spec.Image != nil {
		sourceImage = *m.GCPMachine.Spec.Image
	} else if m.GCPMachine.Spec.ImageFamily != nil {
		sourceImage = *m.GCPMachine.Spec.ImageFamily
	}

	diskType := infrav1.PdStandardDiskType
	if t := m.GCPMachine.Spec.RootDeviceType; t != nil {
		diskType = *t
	}

	disk := &compute.AttachedDisk{
		AutoDelete: true,
		Boot:       true,
		InitializeParams: &compute.AttachedDiskInitializeParams{
			DiskSizeGb:          m.GCPMachine.Spec.RootDeviceSize,
			DiskType:            path.Join("zones", m.Zone(), "diskTypes", string(diskType)),
			ResourceManagerTags: shared.ResourceTagConvert(context.TODO(), m.GCPMachine.Spec.ResourceManagerTags),
			SourceImage:         sourceImage,
			Labels:              m.ClusterGetter.AdditionalLabels().AddLabels(m.GCPMachine.Spec.AdditionalLabels),
		},
	}

	if m.GCPMachine.Spec.RootDiskEncryptionKey != nil {
		if m.GCPMachine.Spec.RootDiskEncryptionKey.KeyType == infrav1.CustomerManagedKey && m.GCPMachine.Spec.RootDiskEncryptionKey.ManagedKey != nil {
			disk.DiskEncryptionKey = &compute.CustomerEncryptionKey{
				KmsKeyName: m.GCPMachine.Spec.RootDiskEncryptionKey.ManagedKey.KMSKeyName,
			}
			if m.GCPMachine.Spec.RootDiskEncryptionKey.KMSKeyServiceAccount != nil {
				disk.DiskEncryptionKey.KmsKeyServiceAccount = *m.GCPMachine.Spec.RootDiskEncryptionKey.KMSKeyServiceAccount
			}
		} else if m.GCPMachine.Spec.RootDiskEncryptionKey.KeyType == infrav1.CustomerSuppliedKey && m.GCPMachine.Spec.RootDiskEncryptionKey.SuppliedKey != nil {
			disk.DiskEncryptionKey = &compute.CustomerEncryptionKey{
				RawKey:          string(m.GCPMachine.Spec.RootDiskEncryptionKey.SuppliedKey.RawKey),
				RsaEncryptedKey: string(m.GCPMachine.Spec.RootDiskEncryptionKey.SuppliedKey.RSAEncryptedKey),
			}
			if m.GCPMachine.Spec.RootDiskEncryptionKey.KMSKeyServiceAccount != nil {
				disk.DiskEncryptionKey.KmsKeyServiceAccount = *m.GCPMachine.Spec.RootDiskEncryptionKey.KMSKeyServiceAccount
			}
		}
	}

	return disk
}

// instanceAdditionalDiskSpec returns compute instance additional attched-disk spec.
func instanceAdditionalDiskSpec(ctx context.Context, spec []infrav1.AttachedDiskSpec, rootDiskEncryptionKey *infrav1.CustomerEncryptionKey, zone string, resourceManagerTags infrav1.ResourceManagerTags) []*compute.AttachedDisk {
	additionalDisks := make([]*compute.AttachedDisk, 0, len(spec))
	for _, disk := range spec {
		additionalDisk := &compute.AttachedDisk{
			AutoDelete: true,
			InitializeParams: &compute.AttachedDiskInitializeParams{
				DiskSizeGb:          ptr.Deref(disk.Size, 30),
				DiskType:            path.Join("zones", zone, "diskTypes", string(*disk.DeviceType)),
				ResourceManagerTags: shared.ResourceTagConvert(ctx, resourceManagerTags),
			},
		}
		if strings.HasSuffix(additionalDisk.InitializeParams.DiskType, string(infrav1.LocalSsdDiskType)) {
			additionalDisk.Type = "SCRATCH" // Default is PERSISTENT.
			// Override the Disk size
			additionalDisk.InitializeParams.DiskSizeGb = 375
			// For local SSDs set interface to NVME (instead of default SCSI) which is faster.
			// Most OS images would work with both NVME and SCSI disks but some may work
			// considerably faster with NVME.
			// https://cloud.google.com/compute/docs/disks/local-ssd#choose_an_interface
			additionalDisk.Interface = "NVME"
		}
		if disk.EncryptionKey != nil {
			if rootDiskEncryptionKey.KeyType == infrav1.CustomerManagedKey && rootDiskEncryptionKey.ManagedKey != nil {
				additionalDisk.DiskEncryptionKey = &compute.CustomerEncryptionKey{
					KmsKeyName: rootDiskEncryptionKey.ManagedKey.KMSKeyName,
				}
				if rootDiskEncryptionKey.KMSKeyServiceAccount != nil {
					additionalDisk.DiskEncryptionKey.KmsKeyServiceAccount = *rootDiskEncryptionKey.KMSKeyServiceAccount
				}
			} else if rootDiskEncryptionKey.KeyType == infrav1.CustomerSuppliedKey && rootDiskEncryptionKey.SuppliedKey != nil {
				additionalDisk.DiskEncryptionKey = &compute.CustomerEncryptionKey{
					RawKey:          string(rootDiskEncryptionKey.SuppliedKey.RawKey),
					RsaEncryptedKey: string(rootDiskEncryptionKey.SuppliedKey.RSAEncryptedKey),
				}
				if rootDiskEncryptionKey.KMSKeyServiceAccount != nil {
					additionalDisk.DiskEncryptionKey.KmsKeyServiceAccount = *rootDiskEncryptionKey.KMSKeyServiceAccount
				}
			}
		}

		additionalDisks = append(additionalDisks, additionalDisk)
	}

	return additionalDisks
}

// InstanceNetworkInterfaceSpec returns compute network interface spec.
func (m *MachineScope) InstanceNetworkInterfaceSpec() *compute.NetworkInterface {
	networkInterface := &compute.NetworkInterface{
		Network: path.Join("projects", m.ClusterGetter.NetworkProject(), "global", "networks", m.ClusterGetter.NetworkName()),
	}

	if m.GCPMachine.Spec.PublicIP != nil && *m.GCPMachine.Spec.PublicIP {
		networkInterface.AccessConfigs = []*compute.AccessConfig{
			{
				Type: "ONE_TO_ONE_NAT",
				Name: "External NAT",
			},
		}
	}

	if m.GCPMachine.Spec.Subnet != nil {
		networkInterface.Subnetwork = path.Join("projects", m.ClusterGetter.NetworkProject(), "regions", m.ClusterGetter.Region(), "subnetworks", *m.GCPMachine.Spec.Subnet)
	}

	networkInterface.AliasIpRanges = m.InstanceNetworkInterfaceAliasIPRangesSpec()

	return networkInterface
}

// InstanceNetworkInterfaceAliasIPRangesSpec returns a slice of Alias IP Range specs.
func (m *MachineScope) InstanceNetworkInterfaceAliasIPRangesSpec() []*compute.AliasIpRange {
	if len(m.GCPMachine.Spec.AliasIPRanges) == 0 {
		return nil
	}
	aliasIPRanges := make([]*compute.AliasIpRange, 0, len(m.GCPMachine.Spec.AliasIPRanges))
	for _, alias := range m.GCPMachine.Spec.AliasIPRanges {
		aliasIPRange := &compute.AliasIpRange{
			IpCidrRange:         alias.IPCidrRange,
			SubnetworkRangeName: alias.SubnetworkRangeName,
		}
		aliasIPRanges = append(aliasIPRanges, aliasIPRange)
	}
	return aliasIPRanges
}

// instanceServiceAccountsSpec returns service-account spec.
func instanceServiceAccountsSpec(serviceAccount *infrav1.ServiceAccount) *compute.ServiceAccount {
	out := &compute.ServiceAccount{
		Email: "default",
		Scopes: []string{
			compute.CloudPlatformScope,
		},
	}

	if serviceAccount != nil {
		out.Email = serviceAccount.Email
		out.Scopes = serviceAccount.Scopes
	}

	return out
}

// InstanceAdditionalMetadataSpec returns additional metadata spec.
func (m *MachineScope) InstanceAdditionalMetadataSpec() *compute.Metadata {
	metadata := new(compute.Metadata)
	for _, additionalMetadata := range m.GCPMachine.Spec.AdditionalMetadata {
		metadata.Items = append(metadata.Items, &compute.MetadataItems{
			Key:   additionalMetadata.Key,
			Value: additionalMetadata.Value,
		})
	}

	return metadata
}

// instanceGuestAcceleratorsSpec returns a slice of Guest Accelerator Config specs.
func instanceGuestAcceleratorsSpec(guestAccelerators []infrav1.Accelerator) []*compute.AcceleratorConfig {
	if len(guestAccelerators) == 0 {
		return nil
	}
	accelConfigs := make([]*compute.AcceleratorConfig, 0, len(guestAccelerators))
	for _, accel := range guestAccelerators {
		accelConfig := &compute.AcceleratorConfig{
			AcceleratorType:  accel.Type,
			AcceleratorCount: accel.Count,
		}
		accelConfigs = append(accelConfigs, accelConfig)
	}
	return accelConfigs
}

// InstanceSpec returns instance spec.
func (m *MachineScope) InstanceSpec(log logr.Logger) *compute.Instance {
	ctx := context.TODO()

	instance := &compute.Instance{
		Name:        m.Name(),
		Zone:        m.Zone(),
		MachineType: path.Join("zones", m.Zone(), "machineTypes", m.GCPMachine.Spec.InstanceType),
		Tags: &compute.Tags{
			Items: append(
				m.GCPMachine.Spec.AdditionalNetworkTags,
				fmt.Sprintf("%s-%s", m.ClusterGetter.Name(), m.Role()),
				m.ClusterGetter.Name(),
			),
		},
		Params: &compute.InstanceParams{
			ResourceManagerTags: shared.ResourceTagConvert(context.TODO(), m.ResourceManagerTags()),
		},
		Labels: infrav1.Build(infrav1.BuildParams{
			ClusterName: m.ClusterGetter.Name(),
			Lifecycle:   infrav1.ResourceLifecycleOwned,
			Role:        ptr.To[string](m.Role()),
			//nolint: godox
			// TODO: Check what needs to be added for the cloud provider label.
			Additional: m.ClusterGetter.AdditionalLabels().AddLabels(m.GCPMachine.Spec.AdditionalLabels),
		}),
		Scheduling: &compute.Scheduling{
			Preemptible: m.GCPMachine.Spec.Preemptible,
		},
	}
	if m.GCPMachine.Spec.ProvisioningModel != nil {
		switch *m.GCPMachine.Spec.ProvisioningModel {
		case infrav1.ProvisioningModelSpot:
			instance.Scheduling.ProvisioningModel = "SPOT"
		case infrav1.ProvisioningModelStandard:
			instance.Scheduling.ProvisioningModel = "STANDARD"
		default:
			log.Error(errors.New("Invalid value"), "Unknown ProvisioningModel value", "Spec.ProvisioningModel", *m.GCPMachine.Spec.ProvisioningModel)
		}
	}

	instance.CanIpForward = true
	if m.GCPMachine.Spec.IPForwarding != nil && *m.GCPMachine.Spec.IPForwarding == infrav1.IPForwardingDisabled {
		instance.CanIpForward = false
	}
	if m.GCPMachine.Spec.ShieldedInstanceConfig != nil {
		instance.ShieldedInstanceConfig = &compute.ShieldedInstanceConfig{
			EnableSecureBoot:          false,
			EnableVtpm:                true,
			EnableIntegrityMonitoring: true,
		}
		if m.GCPMachine.Spec.ShieldedInstanceConfig.SecureBoot == infrav1.SecureBootPolicyEnabled {
			instance.ShieldedInstanceConfig.EnableSecureBoot = true
		}
		if m.GCPMachine.Spec.ShieldedInstanceConfig.VirtualizedTrustedPlatformModule == infrav1.VirtualizedTrustedPlatformModulePolicyDisabled {
			instance.ShieldedInstanceConfig.EnableVtpm = false
		}
		if m.GCPMachine.Spec.ShieldedInstanceConfig.IntegrityMonitoring == infrav1.IntegrityMonitoringPolicyDisabled {
			instance.ShieldedInstanceConfig.EnableIntegrityMonitoring = false
		}
	}
	if m.GCPMachine.Spec.OnHostMaintenance != nil {
		switch *m.GCPMachine.Spec.OnHostMaintenance {
		case infrav1.HostMaintenancePolicyMigrate:
			instance.Scheduling.OnHostMaintenance = "MIGRATE"
		case infrav1.HostMaintenancePolicyTerminate:
			instance.Scheduling.OnHostMaintenance = "TERMINATE"
		default:
			log.Error(errors.New("Invalid value"), "Unknown OnHostMaintenance value", "Spec.OnHostMaintenance", *m.GCPMachine.Spec.OnHostMaintenance)
		}

		instance.Scheduling.OnHostMaintenance = strings.ToUpper(string(*m.GCPMachine.Spec.OnHostMaintenance))
	}
	if m.GCPMachine.Spec.ConfidentialCompute != nil {
		enabled := *m.GCPMachine.Spec.ConfidentialCompute != infrav1.ConfidentialComputePolicyDisabled
		instance.ConfidentialInstanceConfig = &compute.ConfidentialInstanceConfig{
			EnableConfidentialCompute: enabled,
		}
		switch *m.GCPMachine.Spec.ConfidentialCompute {
		case infrav1.ConfidentialComputePolicySEV:
			instance.ConfidentialInstanceConfig.ConfidentialInstanceType = "SEV"
		case infrav1.ConfidentialComputePolicySEVSNP:
			instance.ConfidentialInstanceConfig.ConfidentialInstanceType = "SEV_SNP"
		case infrav1.ConfidentialComputePolicyTDX:
			instance.ConfidentialInstanceConfig.ConfidentialInstanceType = "TDX"
		default:
		}
	}

	instance.Disks = append(instance.Disks, m.InstanceImageSpec())
	instance.Disks = append(instance.Disks, instanceAdditionalDiskSpec(ctx, m.GCPMachine.Spec.AdditionalDisks, m.GCPMachine.Spec.RootDiskEncryptionKey, m.Zone(), m.ResourceManagerTags())...)
	instance.Metadata = m.InstanceAdditionalMetadataSpec()
	instance.ServiceAccounts = append(instance.ServiceAccounts, instanceServiceAccountsSpec(m.GCPMachine.Spec.ServiceAccount))
	instance.NetworkInterfaces = append(instance.NetworkInterfaces, m.InstanceNetworkInterfaceSpec())
	instance.GuestAccelerators = instanceGuestAcceleratorsSpec(m.GCPMachine.Spec.GuestAccelerators)
	if len(instance.GuestAccelerators) > 0 {
		instance.Scheduling.OnHostMaintenance = "TERMINATE"
	}

	return instance
}

// ANCHOR_END: MachineInstanceSpec

// GetBootstrapData returns the bootstrap data from the secret in the Machine's bootstrap.dataSecretName.
func (m *MachineScope) GetBootstrapData() (string, error) {
	if m.Machine.Spec.Bootstrap.DataSecretName == nil {
		return "", errors.New("error retrieving bootstrap data: linked Machine's bootstrap.dataSecretName is nil")
	}

	secret := &corev1.Secret{}
	key := types.NamespacedName{Namespace: m.Namespace(), Name: *m.Machine.Spec.Bootstrap.DataSecretName}
	if err := m.client.Get(context.TODO(), key, secret); err != nil {
		return "", errors.Wrapf(err, "failed to retrieve bootstrap data secret for GCPMachine %s/%s", m.Namespace(), m.Name())
	}

	value, ok := secret.Data["value"]
	if !ok {
		return "", errors.New("error retrieving bootstrap data: secret value key is missing")
	}

	return string(value), nil
}

// PatchObject persists the cluster configuration and status.
func (m *MachineScope) PatchObject() error {
	return m.patchHelper.Patch(context.TODO(), m.GCPMachine)
}

// Close closes the current scope persisting the cluster configuration and status.
func (m *MachineScope) Close() error {
	return m.PatchObject()
}

// ResourceManagerTags merges ResourceManagerTags from the scope's GCPCluster and GCPMachine. If the same key is present in both,
// the value from GCPMachine takes precedence. The returned ResourceManagerTags will never be nil.
func (m *MachineScope) ResourceManagerTags() infrav1.ResourceManagerTags {
	tags := infrav1.ResourceManagerTags{}

	// Start with the cluster-wide tags...
	tags.Merge(m.ClusterGetter.ResourceManagerTags())
	// ... and merge in the Machine's
	tags.Merge(m.GCPMachine.Spec.ResourceManagerTags)

	return tags
}
