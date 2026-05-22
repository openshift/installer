package machines

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"k8s.io/utils/ptr"
	capa "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/yaml"

	machinev1beta1 "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	awstypes "github.com/openshift/installer/pkg/types/aws"
)

// MasterMAPISync validates and syncs provisioning-relevant fields from
// MAPI master machine manifests (openshift/) to CAPI AWSMachine manifests
// (cluster-api/machines/) when drift is detected. This prevents silent
// misconfiguration when users edit openshift/ master manifests after
// "create manifests" but before "create cluster".
type MasterMAPISync struct{}

var _ asset.Asset = (*MasterMAPISync)(nil)

// Name returns the human-friendly name of the asset.
func (a *MasterMAPISync) Name() string {
	return "Master Machine MAPI-to-CAPI Sync"
}

// Dependencies returns the assets upon which this asset directly depends.
func (a *MasterMAPISync) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
		&Master{},
		&ClusterAPI{},
	}
}

// Generate compares MAPI and CAPI master machine specs and syncs
// provisioning-relevant fields from MAPI to CAPI when drift is detected.
func (a *MasterMAPISync) Generate(ctx context.Context, dependencies asset.Parents) error {
	ic := &installconfig.InstallConfig{}
	mastersAsset := &Master{}
	capiAsset := &ClusterAPI{}
	dependencies.Get(ic, mastersAsset, capiAsset)

	if !capiutils.IsEnabled(ic) {
		return nil
	}

	platform := ic.Config.Platform.Name()
	if platform != awstypes.Name {
		return nil
	}

	masters, err := mastersAsset.Machines()
	if err != nil {
		logrus.Debugf("MasterMAPISync: skipping, could not parse MAPI machines: %v", err)
		return nil
	}
	if len(masters) == 0 {
		return nil
	}

	capiFiles := capiAsset.RuntimeFiles()
	if len(capiFiles) == 0 {
		return nil
	}

	awsMachines := indexAWSMachinesByName(capiFiles)
	if len(awsMachines) == 0 {
		return nil
	}

	var driftMessages []string

	for i, m := range masters {
		mapiConfig, ok := m.Spec.ProviderSpec.Value.Object.(*machinev1beta1.AWSMachineProviderConfig)
		if !ok || mapiConfig == nil {
			continue
		}

		machineName := m.Name
		awsMachine, found := awsMachines[machineName]
		if !found {
			awsMachine, found = findAWSMachineByIndex(awsMachines, i)
			if !found {
				continue
			}
		}

		drifts := syncAWSFields(mapiConfig, awsMachine)
		if len(drifts) > 0 {
			driftMessages = append(driftMessages, fmt.Sprintf("  Machine %s:", machineName))
			driftMessages = append(driftMessages, drifts...)
		}
	}

	if len(driftMessages) > 0 {
		logrus.Warnf("Detected drift between MAPI master machine manifests (openshift/) and CAPI machine manifests (cluster-api/machines/).\n%s\n"+
			"  Syncing MAPI values to CAPI manifests for provisioning.\n"+
			"  To avoid this warning, set values in install-config.yaml (controlPlane.platform.aws) before 'create manifests',\n"+
			"  or edit cluster-api/machines/10_inframachine_*-master-*.yaml directly.",
			strings.Join(driftMessages, "\n"))

		for _, rf := range capiFiles {
			if am, ok := rf.Object.(*capa.AWSMachine); ok {
				if !strings.Contains(am.Name, "-master-") {
					continue
				}
				objData, err := yaml.Marshal(rf.Object)
				if err != nil {
					logrus.Debugf("MasterMAPISync: failed to re-serialize AWSMachine %s: %v", am.Name, err)
					continue
				}
				rf.Data = objData
			}
		}
	}

	return nil
}

// syncAWSFields compares and syncs provisioning-relevant fields from MAPI
// AWSMachineProviderConfig to CAPI AWSMachine. Returns a list of drift
// description strings for logging.
func syncAWSFields(mapi *machinev1beta1.AWSMachineProviderConfig, capi *capa.AWSMachine) []string {
	var drifts []string

	if mapi.InstanceType != "" && mapi.InstanceType != capi.Spec.InstanceType {
		drifts = append(drifts, fmt.Sprintf("    instanceType: openshift/ has %q, cluster-api/ has %q → syncing", mapi.InstanceType, capi.Spec.InstanceType))
		capi.Spec.InstanceType = mapi.InstanceType
	}

	if mapi.AMI.ID != nil && *mapi.AMI.ID != "" {
		capiAMI := ptr.Deref(capi.Spec.AMI.ID, "")
		if *mapi.AMI.ID != capiAMI {
			drifts = append(drifts, fmt.Sprintf("    ami.id: openshift/ has %q, cluster-api/ has %q → syncing", *mapi.AMI.ID, capiAMI))
			capi.Spec.AMI.ID = mapi.AMI.ID
		}
	}

	if len(mapi.BlockDevices) > 0 && mapi.BlockDevices[0].EBS != nil {
		ebs := mapi.BlockDevices[0].EBS
		if capi.Spec.RootVolume == nil {
			capi.Spec.RootVolume = &capa.Volume{}
		}

		if ebs.VolumeSize != nil && *ebs.VolumeSize != capi.Spec.RootVolume.Size {
			drifts = append(drifts, fmt.Sprintf("    rootVolume.size: openshift/ has %d, cluster-api/ has %d → syncing", *ebs.VolumeSize, capi.Spec.RootVolume.Size))
			capi.Spec.RootVolume.Size = *ebs.VolumeSize
		}

		if ebs.VolumeType != nil && *ebs.VolumeType != "" {
			capiType := string(capi.Spec.RootVolume.Type)
			if *ebs.VolumeType != capiType {
				drifts = append(drifts, fmt.Sprintf("    rootVolume.type: openshift/ has %q, cluster-api/ has %q → syncing", *ebs.VolumeType, capiType))
				capi.Spec.RootVolume.Type = capa.VolumeType(*ebs.VolumeType)
			}
		}

		if ebs.Iops != nil {
			if *ebs.Iops != capi.Spec.RootVolume.IOPS {
				drifts = append(drifts, fmt.Sprintf("    rootVolume.iops: openshift/ has %d, cluster-api/ has %d → syncing", *ebs.Iops, capi.Spec.RootVolume.IOPS))
				capi.Spec.RootVolume.IOPS = *ebs.Iops
			}
		}

		if ebs.ThroughputMib != nil {
			capiThroughput := ptr.Deref(capi.Spec.RootVolume.Throughput, 0)
			mapiThroughput := int64(*ebs.ThroughputMib)
			if mapiThroughput != capiThroughput {
				drifts = append(drifts, fmt.Sprintf("    rootVolume.throughput: openshift/ has %d, cluster-api/ has %d → syncing", mapiThroughput, capiThroughput))
				capi.Spec.RootVolume.Throughput = ptr.To(mapiThroughput)
			}
		}

		kmsARN := ptr.Deref(ebs.KMSKey.ARN, "")
		kmsID := ptr.Deref(ebs.KMSKey.ID, "")
		mapiKMS := kmsARN
		if mapiKMS == "" {
			mapiKMS = kmsID
		}
		if mapiKMS != "" && mapiKMS != capi.Spec.RootVolume.EncryptionKey {
			drifts = append(drifts, fmt.Sprintf("    rootVolume.encryptionKey: openshift/ has %q, cluster-api/ has %q → syncing", mapiKMS, capi.Spec.RootVolume.EncryptionKey))
			capi.Spec.RootVolume.EncryptionKey = mapiKMS
		}
	}

	if mapi.IAMInstanceProfile != nil && mapi.IAMInstanceProfile.ID != nil && *mapi.IAMInstanceProfile.ID != "" {
		if *mapi.IAMInstanceProfile.ID != capi.Spec.IAMInstanceProfile {
			drifts = append(drifts, fmt.Sprintf("    iamInstanceProfile: openshift/ has %q, cluster-api/ has %q → syncing", *mapi.IAMInstanceProfile.ID, capi.Spec.IAMInstanceProfile))
			capi.Spec.IAMInstanceProfile = *mapi.IAMInstanceProfile.ID
		}
	}

	if mapi.MetadataServiceOptions.Authentication != "" {
		if capi.Spec.InstanceMetadataOptions == nil {
			capi.Spec.InstanceMetadataOptions = &capa.InstanceMetadataOptions{}
		}
		mapiTokens := strings.ToLower(string(mapi.MetadataServiceOptions.Authentication))
		capiTokens := string(capi.Spec.InstanceMetadataOptions.HTTPTokens)
		if mapiTokens != capiTokens {
			drifts = append(drifts, fmt.Sprintf("    instanceMetadataOptions.httpTokens: openshift/ has %q, cluster-api/ has %q → syncing", mapiTokens, capiTokens))
			capi.Spec.InstanceMetadataOptions.HTTPTokens = capa.HTTPTokensState(mapiTokens)
		}
	}

	// CPUOptions / ConfidentialCompute
	if mapi.CPUOptions != nil && mapi.CPUOptions.ConfidentialCompute != nil && *mapi.CPUOptions.ConfidentialCompute != "" {
		mapiCC := capa.AWSConfidentialComputePolicy(*mapi.CPUOptions.ConfidentialCompute)
		if mapiCC != capi.Spec.CPUOptions.ConfidentialCompute {
			drifts = append(drifts, fmt.Sprintf("    cpuOptions.confidentialCompute: openshift/ has %q, cluster-api/ has %q → syncing", mapiCC, capi.Spec.CPUOptions.ConfidentialCompute))
			capi.Spec.CPUOptions.ConfidentialCompute = mapiCC
		}
	}

	if mapi.PublicIP != nil {
		capiPublicIP := ptr.Deref(capi.Spec.PublicIP, false)
		if *mapi.PublicIP != capiPublicIP {
			drifts = append(drifts, fmt.Sprintf("    publicIP: openshift/ has %v, cluster-api/ has %v → syncing", *mapi.PublicIP, capiPublicIP))
			capi.Spec.PublicIP = mapi.PublicIP
		}
	}

	syncSubnet(mapi, capi, &drifts)

	return drifts
}

func syncSubnet(mapi *machinev1beta1.AWSMachineProviderConfig, capi *capa.AWSMachine, drifts *[]string) {
	if mapi.Subnet.ID != nil && *mapi.Subnet.ID != "" {
		if capi.Spec.Subnet == nil {
			capi.Spec.Subnet = &capa.AWSResourceReference{}
		}
		capiSubnetID := ptr.Deref(capi.Spec.Subnet.ID, "")
		if *mapi.Subnet.ID != capiSubnetID {
			*drifts = append(*drifts, fmt.Sprintf("    subnet.id: openshift/ has %q, cluster-api/ has %q → syncing", *mapi.Subnet.ID, capiSubnetID))
			capi.Spec.Subnet.ID = mapi.Subnet.ID
			capi.Spec.Subnet.Filters = nil
		}
	} else if len(mapi.Subnet.Filters) > 0 {
		if capi.Spec.Subnet == nil {
			capi.Spec.Subnet = &capa.AWSResourceReference{}
		}
		mapiFilterStr := formatMAPIFilters(mapi.Subnet.Filters)
		capiFilterStr := formatCAPIFilters(capi.Spec.Subnet.Filters)
		if mapiFilterStr != capiFilterStr {
			*drifts = append(*drifts, fmt.Sprintf("    subnet.filters: openshift/ has %q, cluster-api/ has %q → syncing", mapiFilterStr, capiFilterStr))
			capi.Spec.Subnet.ID = nil
			capi.Spec.Subnet.Filters = convertMAPIFiltersToCAPI(mapi.Subnet.Filters)
		}
	}
}

func indexAWSMachinesByName(files []*asset.RuntimeFile) map[string]*capa.AWSMachine {
	result := make(map[string]*capa.AWSMachine)
	for _, rf := range files {
		if am, ok := rf.Object.(*capa.AWSMachine); ok {
			if strings.Contains(am.Name, "-master-") {
				result[am.Name] = am
			}
		}
	}
	return result
}

func findAWSMachineByIndex(machines map[string]*capa.AWSMachine, idx int) (*capa.AWSMachine, bool) {
	suffix := fmt.Sprintf("-master-%d", idx)
	for name, m := range machines {
		if strings.HasSuffix(name, suffix) {
			return m, true
		}
	}
	return nil, false
}

func formatMAPIFilters(filters []machinev1beta1.Filter) string {
	var parts []string
	for _, f := range filters {
		parts = append(parts, fmt.Sprintf("%s=%s", f.Name, strings.Join(f.Values, ",")))
	}
	return strings.Join(parts, ";")
}

func formatCAPIFilters(filters []capa.Filter) string {
	var parts []string
	for _, f := range filters {
		parts = append(parts, fmt.Sprintf("%s=%s", f.Name, strings.Join(f.Values, ",")))
	}
	return strings.Join(parts, ";")
}

func convertMAPIFiltersToCAPI(filters []machinev1beta1.Filter) []capa.Filter {
	result := make([]capa.Filter, len(filters))
	for i, f := range filters {
		result[i] = capa.Filter{
			Name:   f.Name,
			Values: f.Values,
		}
	}
	return result
}
