package machines

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"

	machinev1 "github.com/openshift/api/machine/v1"
	machinev1beta1 "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	awstypes "github.com/openshift/installer/pkg/types/aws"
)

// MasterCPMSSync detects drift between master machine providerSpecs and the
// CPMS template providerSpec, syncing provisioning-relevant fields so the CPMS
// does not trigger an unintended rolling update post-install.
type MasterCPMSSync struct{}

var _ asset.Asset = (*MasterCPMSSync)(nil)

// Name returns the human-friendly name of the asset.
func (a *MasterCPMSSync) Name() string {
	return "Master CPMS Template Sync"
}

// Dependencies returns the assets upon which this asset directly depends.
func (a *MasterCPMSSync) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
		&Master{},
	}
}

// Generate compares MAPI master machine providerSpecs to the CPMS template
// providerSpec and syncs provisioning-relevant fields when drift is detected.
func (a *MasterCPMSSync) Generate(ctx context.Context, dependencies asset.Parents) error {
	ic := &installconfig.InstallConfig{}
	mastersAsset := &Master{}
	dependencies.Get(ic, mastersAsset)

	if !capiutils.IsEnabled(ic) {
		return nil
	}

	platform := ic.Config.Platform.Name()
	if platform != awstypes.Name {
		return nil
	}

	if mastersAsset.ControlPlaneMachineSet == nil {
		return nil
	}

	masters, err := mastersAsset.Machines()
	if err != nil {
		logrus.Debugf("MasterCPMSSync: skipping, could not parse MAPI machines: %v", err)
		return nil
	}
	if len(masters) == 0 {
		return nil
	}

	cpms := &machinev1.ControlPlaneMachineSet{}
	if err := yaml.Unmarshal(mastersAsset.ControlPlaneMachineSet.Data, cpms); err != nil {
		logrus.Debugf("MasterCPMSSync: skipping, could not parse CPMS: %v", err)
		return nil
	}

	tmpl := cpms.Spec.Template.OpenShiftMachineV1Beta1Machine
	if tmpl == nil {
		return nil
	}

	cpmsProviderSpec, err := decodeCPMSProviderSpec(tmpl.Spec.ProviderSpec.Value)
	if cpmsProviderSpec == nil || err != nil {
		logrus.Debugf("MasterCPMSSync: skipping, could not decode CPMS providerSpec: %v", err)
		return nil
	}

	firstMaster := masters[0]
	mapiConfig, ok := firstMaster.Spec.ProviderSpec.Value.Object.(*machinev1beta1.AWSMachineProviderConfig)
	if !ok || mapiConfig == nil {
		return nil
	}

	drifts := syncCPMSAWSFields(mapiConfig, cpmsProviderSpec)

	if len(drifts) == 0 {
		return nil
	}

	logrus.Warnf("Detected drift between MAPI master machine providerSpec (openshift/) and CPMS template.\n%s\n"+
		"  Syncing master machine values to CPMS template to prevent unintended rolling update.\n"+
		"  To avoid this warning, also edit 99_openshift-machine-api_master-control-plane-machine-set.yaml when customizing masters.",
		strings.Join(drifts, "\n"))

	rawPS, err := json.Marshal(cpmsProviderSpec)
	if err != nil {
		logrus.Debugf("MasterCPMSSync: failed to marshal synced providerSpec: %v", err)
		return nil
	}
	tmpl.Spec.ProviderSpec.Value = &runtime.RawExtension{Raw: rawPS}
	cpmsData, err := yaml.Marshal(cpms)
	if err != nil {
		logrus.Debugf("MasterCPMSSync: failed to re-serialize CPMS: %v", err)
		return nil
	}

	mastersAsset.ControlPlaneMachineSet.Data = cpmsData

	return nil
}

// syncCPMSAWSFields compares provisioning-relevant fields (excluding zone-specific
// fields like Placement.AvailabilityZone and Subnet which are handled by CPMS
// FailureDomains) from MAPI machine to CPMS template. Returns drift descriptions.
func syncCPMSAWSFields(mapi *machinev1beta1.AWSMachineProviderConfig, cpms *machinev1beta1.AWSMachineProviderConfig) []string {
	var drifts []string

	if mapi.InstanceType != "" && mapi.InstanceType != cpms.InstanceType {
		drifts = append(drifts, fmt.Sprintf("    instanceType: machine has %q, CPMS template has %q → syncing", mapi.InstanceType, cpms.InstanceType))
		cpms.InstanceType = mapi.InstanceType
	}

	if mapi.AMI.ID != nil && *mapi.AMI.ID != "" {
		cpmsAMI := ""
		if cpms.AMI.ID != nil {
			cpmsAMI = *cpms.AMI.ID
		}
		if *mapi.AMI.ID != cpmsAMI {
			drifts = append(drifts, fmt.Sprintf("    ami.id: machine has %q, CPMS template has %q → syncing", *mapi.AMI.ID, cpmsAMI))
			cpms.AMI.ID = mapi.AMI.ID
		}
	}

	if len(mapi.BlockDevices) > 0 && mapi.BlockDevices[0].EBS != nil {
		ebs := mapi.BlockDevices[0].EBS
		if len(cpms.BlockDevices) == 0 {
			cpms.BlockDevices = []machinev1beta1.BlockDeviceMappingSpec{{EBS: &machinev1beta1.EBSBlockDeviceSpec{}}}
		}
		if cpms.BlockDevices[0].EBS == nil {
			cpms.BlockDevices[0].EBS = &machinev1beta1.EBSBlockDeviceSpec{}
		}
		cpmsEBS := cpms.BlockDevices[0].EBS

		if ebs.VolumeSize != nil {
			cpmsSize := int64(0)
			if cpmsEBS.VolumeSize != nil {
				cpmsSize = *cpmsEBS.VolumeSize
			}
			if *ebs.VolumeSize != cpmsSize {
				drifts = append(drifts, fmt.Sprintf("    rootVolume.size: machine has %d, CPMS template has %d → syncing", *ebs.VolumeSize, cpmsSize))
				cpmsEBS.VolumeSize = ebs.VolumeSize
			}
		}

		if ebs.VolumeType != nil && *ebs.VolumeType != "" {
			cpmsType := ""
			if cpmsEBS.VolumeType != nil {
				cpmsType = *cpmsEBS.VolumeType
			}
			if *ebs.VolumeType != cpmsType {
				drifts = append(drifts, fmt.Sprintf("    rootVolume.type: machine has %q, CPMS template has %q → syncing", *ebs.VolumeType, cpmsType))
				cpmsEBS.VolumeType = ebs.VolumeType
			}
		}

		if ebs.Iops != nil {
			cpmsIOPS := int64(0)
			if cpmsEBS.Iops != nil {
				cpmsIOPS = *cpmsEBS.Iops
			}
			if *ebs.Iops != cpmsIOPS {
				drifts = append(drifts, fmt.Sprintf("    rootVolume.iops: machine has %d, CPMS template has %d → syncing", *ebs.Iops, cpmsIOPS))
				cpmsEBS.Iops = ebs.Iops
			}
		}

		if ebs.Encrypted != nil {
			cpmsEncrypted := false
			if cpmsEBS.Encrypted != nil {
				cpmsEncrypted = *cpmsEBS.Encrypted
			}
			if *ebs.Encrypted != cpmsEncrypted {
				drifts = append(drifts, fmt.Sprintf("    rootVolume.encrypted: machine has %v, CPMS template has %v → syncing", *ebs.Encrypted, cpmsEncrypted))
				cpmsEBS.Encrypted = ebs.Encrypted
			}
		}

		mapiKMS := ""
		if ebs.KMSKey.ARN != nil {
			mapiKMS = *ebs.KMSKey.ARN
		}
		if mapiKMS == "" && ebs.KMSKey.ID != nil {
			mapiKMS = *ebs.KMSKey.ID
		}
		cpmsKMS := ""
		if cpmsEBS.KMSKey.ARN != nil {
			cpmsKMS = *cpmsEBS.KMSKey.ARN
		}
		if cpmsKMS == "" && cpmsEBS.KMSKey.ID != nil {
			cpmsKMS = *cpmsEBS.KMSKey.ID
		}
		if mapiKMS != "" && mapiKMS != cpmsKMS {
			drifts = append(drifts, fmt.Sprintf("    rootVolume.kmsKey: machine has %q, CPMS template has %q → syncing", mapiKMS, cpmsKMS))
			cpmsEBS.KMSKey = ebs.KMSKey
		}
	}

	if mapi.IAMInstanceProfile != nil && mapi.IAMInstanceProfile.ID != nil && *mapi.IAMInstanceProfile.ID != "" {
		cpmsProfile := ""
		if cpms.IAMInstanceProfile != nil && cpms.IAMInstanceProfile.ID != nil {
			cpmsProfile = *cpms.IAMInstanceProfile.ID
		}
		if *mapi.IAMInstanceProfile.ID != cpmsProfile {
			drifts = append(drifts, fmt.Sprintf("    iamInstanceProfile: machine has %q, CPMS template has %q → syncing", *mapi.IAMInstanceProfile.ID, cpmsProfile))
			cpms.IAMInstanceProfile = mapi.IAMInstanceProfile
		}
	}

	if mapi.MetadataServiceOptions.Authentication != "" {
		mapiAuth := string(mapi.MetadataServiceOptions.Authentication)
		cpmsAuth := string(cpms.MetadataServiceOptions.Authentication)
		if mapiAuth != cpmsAuth {
			drifts = append(drifts, fmt.Sprintf("    metadataServiceOptions.authentication: machine has %q, CPMS template has %q → syncing", mapiAuth, cpmsAuth))
			cpms.MetadataServiceOptions.Authentication = mapi.MetadataServiceOptions.Authentication
		}
	}

	if mapi.PublicIP != nil {
		cpmsPublicIP := false
		if cpms.PublicIP != nil {
			cpmsPublicIP = *cpms.PublicIP
		}
		if *mapi.PublicIP != cpmsPublicIP {
			drifts = append(drifts, fmt.Sprintf("    publicIP: machine has %v, CPMS template has %v → syncing", *mapi.PublicIP, cpmsPublicIP))
			cpms.PublicIP = mapi.PublicIP
		}
	}

	return drifts
}

func decodeCPMSProviderSpec(raw *runtime.RawExtension) (*machinev1beta1.AWSMachineProviderConfig, error) {
	if raw == nil {
		return nil, nil
	}

	if raw.Object != nil {
		if cfg, ok := raw.Object.(*machinev1beta1.AWSMachineProviderConfig); ok {
			return cfg, nil
		}
	}

	if raw.Raw == nil {
		return nil, nil
	}

	cfg := &machinev1beta1.AWSMachineProviderConfig{}
	if err := json.Unmarshal(raw.Raw, cfg); err != nil {
		return nil, fmt.Errorf("decode CPMS providerSpec: %w", err)
	}
	return cfg, nil
}
