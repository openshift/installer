/*
Copyright 2025 The Kubernetes Authors.

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
	"path"
	"strings"

	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/meta"
	"github.com/pkg/errors"
	"golang.org/x/mod/semver"
	"google.golang.org/api/compute/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud/services/shared"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-gcp/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-gcp/pkg/gcp"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// MachinePoolScope defines a scope defined around a machine and its cluster.
type MachinePoolScope struct {
	client                     client.Client
	patchHelper                *patch.Helper
	capiMachinePoolPatchHelper *patch.Helper

	ClusterGetter  cloud.ClusterGetter
	MachinePool    *clusterv1.MachinePool
	GCPMachinePool *expinfrav1.GCPMachinePool
}

// MachinePoolScopeParams defines a scope defined around a machine and its cluster.
type MachinePoolScopeParams struct {
	client.Client

	ClusterGetter  cloud.ClusterGetter
	MachinePool    *clusterv1.MachinePool
	GCPMachinePool *expinfrav1.GCPMachinePool
}

// NewMachinePoolScope creates a new MachinePoolScope from the supplied parameters.
// This is meant to be called for each reconcile iteration.
func NewMachinePoolScope(params MachinePoolScopeParams) (*MachinePoolScope, error) {
	if params.ClusterGetter == nil {
		return nil, errors.New("clusterGetter is required when creating a MachinePoolScope")
	}
	if params.Client == nil {
		return nil, errors.New("client is required when creating a MachinePoolScope")
	}
	if params.MachinePool == nil {
		return nil, errors.New("machinepool is required when creating a MachinePoolScope")
	}
	if params.GCPMachinePool == nil {
		return nil, errors.New("gcp machine pool is required when creating a MachinePoolScope")
	}

	ampHelper, err := patch.NewHelper(params.GCPMachinePool, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init GCPMachinePool patch helper")
	}
	mpHelper, err := patch.NewHelper(params.MachinePool, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init MachinePool patch helper")
	}

	return &MachinePoolScope{
		client:                     params.Client,
		patchHelper:                ampHelper,
		capiMachinePoolPatchHelper: mpHelper,

		ClusterGetter:  params.ClusterGetter,
		MachinePool:    params.MachinePool,
		GCPMachinePool: params.GCPMachinePool,
	}, nil
}

// Cloud returns initialized cloud.
func (m *MachinePoolScope) Cloud() cloud.Cloud {
	return m.ClusterGetter.Cloud()
}

// Name returns the GCPMachinePool name.
func (m *MachinePoolScope) Name() string {
	return m.GCPMachinePool.Name
}

// Namespace returns the namespace name.
func (m *MachinePoolScope) Namespace() string {
	return m.GCPMachinePool.Namespace
}

// ClusterName returns the cluster name.
func (m *MachinePoolScope) ClusterName() string {
	return m.ClusterGetter.Name()
}

// GetBootstrapData returns the bootstrap data from the secret in the Machine's bootstrap.dataSecretName.
func (m *MachinePoolScope) getBootstrapData(ctx context.Context) (string, error) {
	return GetBootstrapData(ctx, m.client, m.MachinePool, m.MachinePool.Spec.Template.Spec.Bootstrap)
}

// Zones returns the targeted zones for the machine pool
func (m *MachinePoolScope) Zones() []string {
	zones := m.MachinePool.Spec.FailureDomains
	if len(zones) == 0 {
		zones = append(zones, m.ClusterGetter.FailureDomains()...)
	}
	return zones
}

// Region returns the region for the GCP resources
func (m *MachinePoolScope) Region() string {
	return m.ClusterGetter.Region()
}

// PatchObject persists the machinepool spec and status.
func (m *MachinePoolScope) PatchObject(ctx context.Context) error {
	return m.patchHelper.Patch(
		ctx,
		m.GCPMachinePool,
		patch.WithOwnedConditions{Conditions: []string{
			string(expinfrav1.MIGReadyCondition),
			string(expinfrav1.InstanceTemplateReadyCondition),
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
	return m.PatchObject(context.TODO())
}

// InstanceGroupManagerResourceName is the name to use for the instanceGroupManager GCP resource
func (m *MachinePoolScope) InstanceGroupManagerResourceName() (*meta.Key, error) {
	igmName := m.ClusterName() + "-" + m.Name()

	zones := m.Zones()
	if len(zones) != 1 {
		return nil, fmt.Errorf("instanceGroupManager must be created in a single zone, got %d zones (%s)", len(zones), strings.Join(zones, ","))
	}
	zone := zones[0]
	igmKey := meta.ZonalKey(igmName, zone)

	return igmKey, nil
}

// InstanceGroupManagerResource is the desired state for the instanceGroupManager GCP resource
func (m *MachinePoolScope) InstanceGroupManagerResource(instanceTemplate *meta.Key) (*compute.InstanceGroupManager, error) {
	instanceTemplateSelfLink := gcp.SelfLink("instanceTemplates", instanceTemplate)
	baseInstanceName := m.Name()

	zones := m.Zones()
	if len(zones) == 0 {
		return nil, errors.New("must specify at least one zone")
	}

	replicas := int64(1)
	if p := m.MachinePool.Spec.Replicas; p != nil {
		replicas = int64(*p)
	}

	desired := &compute.InstanceGroupManager{
		BaseInstanceName: baseInstanceName,
		Description:      "",
		InstanceTemplate: instanceTemplateSelfLink,
		TargetSize:       replicas,
	}

	// DistributionPolicy can only be used if there are multiple zones
	if len(zones) > 1 {
		desired.DistributionPolicy = &compute.DistributionPolicy{}
		for _, zone := range zones {
			zoneSelfLink, err := buildZoneSelfLink(zone)
			if err != nil {
				return nil, err
			}
			desired.DistributionPolicy.Zones = append(desired.DistributionPolicy.Zones, &compute.DistributionPolicyZoneConfiguration{
				Zone: zoneSelfLink,
			})
		}
	} else {
		zoneSelfLink, err := buildZoneSelfLink(zones[0])
		if err != nil {
			return nil, err
		}
		desired.Zone = zoneSelfLink
	}

	return desired, nil
}

// buildZoneSelfLink returns a fully-qualified zone link from a user-provided zone
func buildZoneSelfLink(zone string) (string, error) {
	tokens := strings.Split(zone, "/")
	if len(tokens) == 1 {
		return "zones/" + tokens[0], nil
	}
	return "", fmt.Errorf("zone %q was not a recognized format", zone)
}

// BaseInstanceTemplateResourceName is the base name to use for the instanceTemplate GCP resource.
// The instance template is immutable, so we add a suffix that hash-encodes the version
func (m *MachinePoolScope) BaseInstanceTemplateResourceName() (*meta.Key, error) {
	name := m.Name()

	// We only use the first 46 characters, to leave room for a 16 character hash
	// 63 characters max, 16 character hash; 1 hyphen
	namePrefix := limitStringLength(name, 63-16-1) + "-"

	region := m.Region()
	return meta.RegionalKey(namePrefix, region), nil
}

// limitStringLength returns the string truncated to the specified maximum length.
func limitStringLength(s string, maxLength int) string {
	if len(s) > maxLength {
		return s[:maxLength]
	}
	return s
}

// InstanceTemplateResource is the desired state for the instanceTemplate GCP resource
func (m *MachinePoolScope) InstanceTemplateResource(ctx context.Context) (*compute.InstanceTemplate, error) {
	log := log.FromContext(ctx)

	bootstrapData, err := m.getBootstrapData(ctx)
	if err != nil {
		return nil, fmt.Errorf("retrieving bootstrap data for instanceTemplate: %w", err)
	}

	instance := &compute.InstanceProperties{
		MachineType: m.GCPMachinePool.Spec.InstanceType,
		Tags: &compute.Tags{
			Items: append(
				m.GCPMachinePool.Spec.AdditionalNetworkTags,
				fmt.Sprintf("%s-%s", m.ClusterGetter.Name(), m.Role()),
				m.ClusterGetter.Name(),
			),
		},
		ResourceManagerTags: shared.ResourceTagConvert(ctx, m.ResourceManagerTags()),
		Labels: infrav1.Build(infrav1.BuildParams{
			ClusterName: m.ClusterGetter.Name(),
			Lifecycle:   infrav1.ResourceLifecycleOwned,
			Role:        ptr.To[string](m.Role()),
			//nolint: godox
			// TODO: Check what needs to be added for the cloud provider label.
			Additional: m.ClusterGetter.AdditionalLabels().AddLabels(m.GCPMachinePool.Spec.AdditionalLabels),
		}),
		Scheduling: &compute.Scheduling{
			Preemptible: m.GCPMachinePool.Spec.Preemptible,
		},
	}

	if m.GCPMachinePool.Spec.ProvisioningModel != nil {
		switch *m.GCPMachinePool.Spec.ProvisioningModel {
		case infrav1.ProvisioningModelSpot:
			instance.Scheduling.ProvisioningModel = "SPOT"
		case infrav1.ProvisioningModelStandard:
			instance.Scheduling.ProvisioningModel = "STANDARD"
		default:
			return nil, fmt.Errorf("unknown ProvisioningModel value: %q", *m.GCPMachinePool.Spec.ProvisioningModel)
		}
	}

	instance.CanIpForward = true
	if m.GCPMachinePool.Spec.IPForwarding != nil && *m.GCPMachinePool.Spec.IPForwarding == infrav1.IPForwardingDisabled {
		instance.CanIpForward = false
	}
	if config := m.GCPMachinePool.Spec.ShieldedInstanceConfig; config != nil {
		instance.ShieldedInstanceConfig = &compute.ShieldedInstanceConfig{
			EnableSecureBoot:          false,
			EnableVtpm:                true,
			EnableIntegrityMonitoring: true,
		}
		if config.SecureBoot == infrav1.SecureBootPolicyEnabled {
			instance.ShieldedInstanceConfig.EnableSecureBoot = true
		}
		if config.VirtualizedTrustedPlatformModule == infrav1.VirtualizedTrustedPlatformModulePolicyDisabled {
			instance.ShieldedInstanceConfig.EnableVtpm = false
		}
		if config.IntegrityMonitoring == infrav1.IntegrityMonitoringPolicyDisabled {
			instance.ShieldedInstanceConfig.EnableIntegrityMonitoring = false
		}
	}
	if onHostMaintenance := ValueOf(m.GCPMachinePool.Spec.OnHostMaintenance); onHostMaintenance != "" {
		switch onHostMaintenance {
		case infrav1.HostMaintenancePolicyMigrate:
			instance.Scheduling.OnHostMaintenance = onHostMaintenanceMigrate
		case infrav1.HostMaintenancePolicyTerminate:
			instance.Scheduling.OnHostMaintenance = onHostMaintenanceTerminate
		default:
			log.Error(errors.New("Invalid value"), "Unknown OnHostMaintenance value", "Spec.OnHostMaintenance", onHostMaintenance)
			instance.Scheduling.OnHostMaintenance = strings.ToUpper(string(onHostMaintenance))
		}
	}

	if confidentialCompute := m.GCPMachinePool.Spec.ConfidentialCompute; confidentialCompute != nil {
		enabled := *confidentialCompute != infrav1.ConfidentialComputePolicyDisabled
		instance.ConfidentialInstanceConfig = &compute.ConfidentialInstanceConfig{
			EnableConfidentialCompute: enabled,
		}
		switch *confidentialCompute {
		case infrav1.ConfidentialComputePolicySEV:
			instance.ConfidentialInstanceConfig.ConfidentialInstanceType = "SEV"
		case infrav1.ConfidentialComputePolicySEVSNP:
			instance.ConfidentialInstanceConfig.ConfidentialInstanceType = "SEV_SNP"
		case infrav1.ConfidentialComputePolicyTDX:
			instance.ConfidentialInstanceConfig.ConfidentialInstanceType = "TDX"
		default:
		}
	}

	instance.Disks = append(instance.Disks, m.InstanceImageSpec(ctx))
	instance.Disks = append(instance.Disks, m.InstanceAdditionalDiskSpec()...)
	instance.Metadata = InstanceAdditionalMetadataSpec(m.GCPMachinePool.Spec.AdditionalMetadata)
	instance.ServiceAccounts = append(instance.ServiceAccounts, instanceServiceAccountsSpec(m.GCPMachinePool.Spec.ServiceAccount))
	var aliasIPRanges []infrav1.AliasIPRange // Not supported by MachinePool
	instance.NetworkInterfaces = append(instance.NetworkInterfaces, InstanceNetworkInterfaceSpec(m.ClusterGetter, m.GCPMachinePool.Spec.PublicIP, m.GCPMachinePool.Spec.Subnet, aliasIPRanges))
	instance.GuestAccelerators = instanceGuestAcceleratorsSpec(m.GCPMachinePool.Spec.GuestAccelerators)
	if len(instance.GuestAccelerators) > 0 {
		instance.Scheduling.OnHostMaintenance = onHostMaintenanceTerminate
	}

	instance.Metadata.Items = append(instance.Metadata.Items, &compute.MetadataItems{
		Key:   "user-data",
		Value: ptr.To[string](bootstrapData),
	})

	instanceTemplate := &compute.InstanceTemplate{
		Region:     m.Region(),
		Properties: instance,
	}

	return instanceTemplate, nil
}

// InstanceImageSpec returns compute instance image attched-disk spec.
func (m *MachinePoolScope) InstanceImageSpec(ctx context.Context) *compute.AttachedDisk {
	spec := m.GCPMachinePool.Spec

	version := m.MachinePool.Spec.Template.Spec.Version

	image := "capi-ubuntu-1804-k8s-" + strings.ReplaceAll(semver.MajorMinor(version), ".", "-")
	sourceImage := path.Join("projects", m.ClusterGetter.Project(), "global", "images", "family", image)
	if spec.Image != nil {
		sourceImage = *spec.Image
	} else if spec.ImageFamily != nil {
		sourceImage = *spec.ImageFamily
	}

	diskType := infrav1.PdStandardDiskType
	if t := spec.RootDeviceType; t != nil {
		diskType = *t
	}

	disk := &compute.AttachedDisk{
		AutoDelete: true,
		Boot:       true,
		InitializeParams: &compute.AttachedDiskInitializeParams{
			DiskSizeGb:          spec.RootDeviceSize,
			DiskType:            string(diskType),
			ResourceManagerTags: shared.ResourceTagConvert(ctx, spec.ResourceManagerTags),
			SourceImage:         sourceImage,
			Labels:              m.ClusterGetter.AdditionalLabels().AddLabels(spec.AdditionalLabels),
		},
	}

	if spec.RootDiskEncryptionKey != nil {
		if spec.RootDiskEncryptionKey.KeyType == infrav1.CustomerManagedKey && spec.RootDiskEncryptionKey.ManagedKey != nil {
			disk.DiskEncryptionKey = &compute.CustomerEncryptionKey{
				KmsKeyName: spec.RootDiskEncryptionKey.ManagedKey.KMSKeyName,
			}
			if spec.RootDiskEncryptionKey.KMSKeyServiceAccount != nil {
				disk.DiskEncryptionKey.KmsKeyServiceAccount = *spec.RootDiskEncryptionKey.KMSKeyServiceAccount
			}
		} else if spec.RootDiskEncryptionKey.KeyType == infrav1.CustomerSuppliedKey && spec.RootDiskEncryptionKey.SuppliedKey != nil {
			disk.DiskEncryptionKey = &compute.CustomerEncryptionKey{
				RawKey:          string(spec.RootDiskEncryptionKey.SuppliedKey.RawKey),
				RsaEncryptedKey: string(spec.RootDiskEncryptionKey.SuppliedKey.RSAEncryptedKey),
			}
			if spec.RootDiskEncryptionKey.KMSKeyServiceAccount != nil {
				disk.DiskEncryptionKey.KmsKeyServiceAccount = *spec.RootDiskEncryptionKey.KMSKeyServiceAccount
			}
		}
	}

	return disk
}

// InstanceAdditionalDiskSpec returns compute instance additional attched-disk spec.
func (m *MachinePoolScope) InstanceAdditionalDiskSpec() []*compute.AttachedDisk {
	spec := m.GCPMachinePool.Spec

	additionalDisks := make([]*compute.AttachedDisk, 0, len(spec.AdditionalDisks))
	for _, disk := range spec.AdditionalDisks {
		diskType := string(ValueOf(disk.DeviceType))

		additionalDisk := &compute.AttachedDisk{
			AutoDelete: true,
			InitializeParams: &compute.AttachedDiskInitializeParams{
				DiskSizeGb:          ptr.Deref(disk.Size, 30),
				DiskType:            diskType,
				ResourceManagerTags: shared.ResourceTagConvert(context.TODO(), spec.ResourceManagerTags),
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
			if spec.RootDiskEncryptionKey.KeyType == infrav1.CustomerManagedKey && spec.RootDiskEncryptionKey.ManagedKey != nil {
				additionalDisk.DiskEncryptionKey = &compute.CustomerEncryptionKey{
					KmsKeyName: spec.RootDiskEncryptionKey.ManagedKey.KMSKeyName,
				}
				if spec.RootDiskEncryptionKey.KMSKeyServiceAccount != nil {
					additionalDisk.DiskEncryptionKey.KmsKeyServiceAccount = *spec.RootDiskEncryptionKey.KMSKeyServiceAccount
				}
			} else if spec.RootDiskEncryptionKey.KeyType == infrav1.CustomerSuppliedKey && spec.RootDiskEncryptionKey.SuppliedKey != nil {
				additionalDisk.DiskEncryptionKey = &compute.CustomerEncryptionKey{
					RawKey:          string(spec.RootDiskEncryptionKey.SuppliedKey.RawKey),
					RsaEncryptedKey: string(spec.RootDiskEncryptionKey.SuppliedKey.RSAEncryptedKey),
				}
				if spec.RootDiskEncryptionKey.KMSKeyServiceAccount != nil {
					additionalDisk.DiskEncryptionKey.KmsKeyServiceAccount = *spec.RootDiskEncryptionKey.KMSKeyServiceAccount
				}
			}
		}

		additionalDisks = append(additionalDisks, additionalDisk)
	}

	return additionalDisks
}

// ResourceManagerTags merges ResourceManagerTags from the scope's GCPCluster and GCPMachine. If the same key is present in both,
// the value from GCPMachine takes precedence. The returned ResourceManagerTags will never be nil.
func (m *MachinePoolScope) ResourceManagerTags() infrav1.ResourceManagerTags {
	tags := infrav1.ResourceManagerTags{}

	// Start with the cluster-wide tags...
	tags.Merge(m.ClusterGetter.ResourceManagerTags())
	// ... and merge in the Machine's
	tags.Merge(m.GCPMachinePool.Spec.ResourceManagerTags)

	return tags
}

// Role returns the machine role from the labels.
func (m *MachinePoolScope) Role() string {
	_, isControlPlane := m.MachinePool.Labels[clusterv1.MachineControlPlaneLabel]

	if isControlPlane {
		return "control-plane"
	}

	return "node"
}

// ValueOf is a generic helper function that returns the value of a pointer, or the empty value if the pointer is nil.
func ValueOf[V any](v *V) V {
	if v != nil {
		return *v
	}
	var zero V
	return zero
}
