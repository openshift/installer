package agent

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/validation/field"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/assisted-service/models"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
	baremetaldefaults "github.com/openshift/installer/pkg/types/baremetal/defaults"
	"github.com/openshift/installer/pkg/types/external"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/validation"
	"github.com/openshift/installer/pkg/types/vsphere"
)

const (
	installConfigFilename = "install-config.yaml"
)

// OptionalInstallConfig is an InstallConfig where the default is empty, rather
// than generated from running the survey.
type OptionalInstallConfig struct {
	installconfig.AssetBase
	Supplied bool
}

var _ asset.WritableAsset = (*OptionalInstallConfig)(nil)

// Dependencies returns all of the dependencies directly needed by an
// InstallConfig asset.
func (a *OptionalInstallConfig) Dependencies() []asset.Asset {
	// Return no dependencies for the Agent install config, because it is
	// optional. We don't need to run the survey if it doesn't exist, since the
	// user may have supplied cluster-manifests that fully define the cluster.
	return []asset.Asset{}
}

// Generate generates the install-config.yaml file.
func (a *OptionalInstallConfig) Generate(parents asset.Parents) error {
	// Just generate an empty install config, since we have no dependencies.
	return nil
}

// Load returns the installconfig from disk.
func (a *OptionalInstallConfig) Load(f asset.FileFetcher) (bool, error) {
	found, err := a.LoadFromFile(f)
	if found && err == nil {
		a.Supplied = true
		if err := a.validateInstallConfig(a.Config).ToAggregate(); err != nil {
			return false, errors.Wrapf(err, "invalid install-config configuration")
		}
		if err := a.RecordFile(); err != nil {
			return false, err
		}
	}
	return found, err
}

func (a *OptionalInstallConfig) validateInstallConfig(installConfig *types.InstallConfig) field.ErrorList {
	var allErrs field.ErrorList
	if err := validation.ValidateInstallConfig(a.Config, true); err != nil {
		allErrs = append(allErrs, err...)
	}

	if err := a.validateSupportedPlatforms(installConfig); err != nil {
		allErrs = append(allErrs, err...)
	}

	if err := a.validateSupportedArchs(installConfig); err != nil {
		allErrs = append(allErrs, err...)
	}

	warnUnusedConfig(installConfig)

	numMasters, numWorkers := GetReplicaCount(installConfig)
	logrus.Infof(fmt.Sprintf("Configuration has %d master replicas and %d worker replicas", numMasters, numWorkers))

	if err := a.validateSNOConfiguration(installConfig); err != nil {
		allErrs = append(allErrs, err...)
	}

	return allErrs
}

func (a *OptionalInstallConfig) validateSupportedPlatforms(installConfig *types.InstallConfig) field.ErrorList {
	var allErrs field.ErrorList

	fieldPath := field.NewPath("Platform")

	if installConfig.Platform.Name() != "" && !IsSupportedPlatform(HivePlatformType(installConfig.Platform)) {
		allErrs = append(allErrs, field.NotSupported(fieldPath, installConfig.Platform.Name(), SupportedInstallerPlatforms()))
	}
	if installConfig.Platform.Name() != none.Name && installConfig.ControlPlane.Architecture == types.ArchitecturePPC64LE {
		allErrs = append(allErrs, field.Invalid(fieldPath, installConfig.Platform.Name(), fmt.Sprintf("CPU architecture \"%s\" only supports platform \"%s\".", types.ArchitecturePPC64LE, none.Name)))
	}
	if installConfig.Platform.Name() == external.Name {
		if installConfig.Platform.External.PlatformName != string(models.PlatformTypeOci) {
			fieldPath = field.NewPath("Platform", "External", "PlatformName")
			allErrs = append(allErrs, field.NotSupported(fieldPath, installConfig.Platform.External.PlatformName, []string{string(models.PlatformTypeOci)}))
		}
		if installConfig.Platform.External.PlatformName == string(models.PlatformTypeOci) &&
			installConfig.Platform.External.CloudControllerManager != external.CloudControllerManagerTypeExternal {
			fieldPath = field.NewPath("Platform", "External", "CloudControllerManager")
			allErrs = append(allErrs, field.Invalid(fieldPath, installConfig.Platform.External.CloudControllerManager, fmt.Sprintf("When using external %s platform, %s must be set to %s", string(models.PlatformTypeOci), fieldPath, external.CloudControllerManagerTypeExternal)))
		}
	}

	if installConfig.Platform.Name() == vsphere.Name {
		allErrs = append(allErrs, a.validateVSpherePlatform(installConfig)...)
	}

	return allErrs
}

func (a *OptionalInstallConfig) validateSupportedArchs(installConfig *types.InstallConfig) field.ErrorList {
	var allErrs field.ErrorList

	fieldPath := field.NewPath("ControlPlane", "Architecture")

	switch string(installConfig.ControlPlane.Architecture) {
	case types.ArchitectureAMD64:
	case types.ArchitectureARM64:
	case types.ArchitecturePPC64LE:
	default:
		allErrs = append(allErrs, field.NotSupported(fieldPath, installConfig.ControlPlane.Architecture, []string{types.ArchitectureAMD64, types.ArchitectureARM64, types.ArchitecturePPC64LE}))
	}

	for i, compute := range installConfig.Compute {
		fieldPath := field.NewPath(fmt.Sprintf("Compute[%d]", i), "Architecture")

		switch string(compute.Architecture) {
		case types.ArchitectureAMD64:
		case types.ArchitectureARM64:
		case types.ArchitecturePPC64LE:
		default:
			allErrs = append(allErrs, field.NotSupported(fieldPath, compute.Architecture, []string{types.ArchitectureAMD64, types.ArchitectureARM64, types.ArchitecturePPC64LE}))
		}
	}

	return allErrs
}

func (a *OptionalInstallConfig) validateSNOConfiguration(installConfig *types.InstallConfig) field.ErrorList {
	var allErrs field.ErrorList
	var fieldPath *field.Path

	var workers int
	for _, worker := range installConfig.Compute {
		workers = workers + int(*worker.Replicas)
	}

	if installConfig.ControlPlane != nil && *installConfig.ControlPlane.Replicas == 1 {
		if workers == 0 {
			if (installConfig.Platform.Name() == none.Name || installConfig.Platform.Name() == external.Name) && installConfig.Networking.NetworkType != "OVNKubernetes" {
				fieldPath = field.NewPath("Networking", "NetworkType")
				allErrs = append(allErrs, field.Invalid(fieldPath, installConfig.Networking.NetworkType, "Only OVNKubernetes network type is allowed for Single Node OpenShift (SNO) cluster"))
			}
			if installConfig.Platform.Name() != none.Name && installConfig.Platform.Name() != external.Name {
				fieldPath = field.NewPath("Platform")
				allErrs = append(allErrs, field.Invalid(fieldPath, installConfig.Platform.Name(), fmt.Sprintf("Only platform %s and %s supports 1 ControlPlane and 0 Compute nodes", none.Name, external.Name)))
			}
		} else {
			fieldPath = field.NewPath("Compute", "Replicas")
			allErrs = append(allErrs, field.Required(fieldPath, fmt.Sprintf("Total number of Compute.Replicas must be 0 when ControlPlane.Replicas is 1 for platform %s or %s. Found %v", none.Name, external.Name, workers)))
		}
	}
	return allErrs
}

// VCenterCredentialsAreProvided returns true if server, username, password, or at least one datacenter
// have been provided.
func VCenterCredentialsAreProvided(vcenter vsphere.VCenter) bool {
	if vcenter.Server != "" || vcenter.Username != "" || vcenter.Password != "" || len(vcenter.Datacenters) > 0 {
		return true
	}
	return false
}

func (a *OptionalInstallConfig) validateVSpherePlatform(installConfig *types.InstallConfig) field.ErrorList {
	var allErrs field.ErrorList
	vspherePlatform := installConfig.Platform.VSphere
	userProvidedCredentials := false
	for _, vcenter := range vspherePlatform.VCenters {
		if VCenterCredentialsAreProvided(vcenter) {
			userProvidedCredentials = true
		}
	}

	for _, failureDomain := range vspherePlatform.FailureDomains {
		// Although folder is optional in IPI/UPI, it must be set for agent-based installs.
		// If it is not set, assisted-service will set a placeholder value for folder:
		// "/datacenterplaceholder/vm/folderplaceholder"
		//
		// When assisted-service generates the install-config for the cluster, it will fail
		// validation because the placeholder value's datacenter name may not match
		// the datacenter set in the failureDomain in the install-config.yaml submitted
		// to the agent-based create image command.
		if failureDomain.Topology.Folder == "" && userProvidedCredentials {
			fieldPath := field.NewPath("Platform", "VSphere", "failureDomains", "topology", "folder")
			allErrs = append(allErrs, field.Required(fieldPath, "must specify a folder for agent-based installs"))
		}
	}

	return allErrs
}

// ClusterName returns the name of the cluster, or a default name if no
// InstallConfig is supplied.
func (a *OptionalInstallConfig) ClusterName() string {
	if a.Config != nil && a.Config.ObjectMeta.Name != "" {
		return a.Config.ObjectMeta.Name
	}
	return "agent-cluster"
}

func warnUnusedConfig(installConfig *types.InstallConfig) {
	// "Proxyonly" is the default set from generic install config code
	if installConfig.AdditionalTrustBundlePolicy != "Proxyonly" {
		fieldPath := field.NewPath("AdditionalTrustBundlePolicy")
		logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, installConfig.AdditionalTrustBundlePolicy))
	}

	for i, compute := range installConfig.Compute {
		if compute.Hyperthreading != "Enabled" {
			fieldPath := field.NewPath(fmt.Sprintf("Compute[%d]", i), "Hyperthreading")
			logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, compute.Hyperthreading))
		}

		if compute.Platform != (types.MachinePoolPlatform{}) {
			fieldPath := field.NewPath(fmt.Sprintf("Compute[%d]", i), "Platform")
			logrus.Warnf(fmt.Sprintf("%s is ignored", fieldPath))
		}
	}

	if installConfig.ControlPlane.Hyperthreading != "Enabled" {
		fieldPath := field.NewPath("ControlPlane", "Hyperthreading")
		logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, installConfig.ControlPlane.Hyperthreading))
	}

	if installConfig.ControlPlane.Platform != (types.MachinePoolPlatform{}) {
		fieldPath := field.NewPath("ControlPlane", "Platform")
		logrus.Warnf(fmt.Sprintf("%s is ignored", fieldPath))
	}

	switch installConfig.Platform.Name() {

	case baremetal.Name:
		defaultIc := &types.InstallConfig{Platform: types.Platform{BareMetal: &baremetal.Platform{}}}
		baremetaldefaults.SetPlatformDefaults(defaultIc.Platform.BareMetal, defaultIc)

		baremetal := installConfig.Platform.BareMetal
		defaultBM := defaultIc.Platform.BareMetal
		// Compare values from generic installconfig code to check for changes
		if baremetal.LibvirtURI != defaultBM.LibvirtURI {
			fieldPath := field.NewPath("Platform", "Baremetal", "LibvirtURI")
			logrus.Debugf(fmt.Sprintf("%s: %s is ignored", fieldPath, baremetal.LibvirtURI))
		}
		if baremetal.ClusterProvisioningIP != defaultBM.ClusterProvisioningIP {
			fieldPath := field.NewPath("Platform", "Baremetal", "ClusterProvisioningIP")
			logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, baremetal.ClusterProvisioningIP))
		}
		if baremetal.DeprecatedProvisioningHostIP != defaultBM.DeprecatedProvisioningHostIP {
			fieldPath := field.NewPath("Platform", "Baremetal", "ProvisioningHostIP")
			logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, baremetal.DeprecatedProvisioningHostIP))
		}
		if baremetal.BootstrapProvisioningIP != defaultBM.BootstrapProvisioningIP {
			fieldPath := field.NewPath("Platform", "Baremetal", "BootstrapProvisioningIP")
			logrus.Debugf(fmt.Sprintf("%s: %s is ignored", fieldPath, baremetal.BootstrapProvisioningIP))
		}
		if baremetal.ExternalBridge != defaultBM.ExternalBridge {
			fieldPath := field.NewPath("Platform", "Baremetal", "ExternalBridge")
			logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, baremetal.ExternalBridge))
		}
		if baremetal.ProvisioningNetwork != defaultBM.ProvisioningNetwork {
			fieldPath := field.NewPath("Platform", "Baremetal", "ProvisioningNetwork")
			logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, baremetal.ProvisioningNetwork))
		}
		if baremetal.ProvisioningBridge != defaultBM.ProvisioningBridge {
			fieldPath := field.NewPath("Platform", "Baremetal", "ProvisioningBridge")
			logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, baremetal.ProvisioningBridge))
		}
		if baremetal.ProvisioningNetworkInterface != defaultBM.ProvisioningNetworkInterface {
			fieldPath := field.NewPath("Platform", "Baremetal", "ProvisioningNetworkInterface")
			logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, baremetal.ProvisioningNetworkInterface))
		}
		if baremetal.ProvisioningNetworkCIDR.String() != defaultBM.ProvisioningNetworkCIDR.String() {
			fieldPath := field.NewPath("Platform", "Baremetal", "ProvisioningNetworkCIDR")
			logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, baremetal.ProvisioningNetworkCIDR))
		}
		if baremetal.DeprecatedProvisioningDHCPExternal != defaultBM.DeprecatedProvisioningDHCPExternal {
			fieldPath := field.NewPath("Platform", "Baremetal", "ProvisioningDHCPExternal")
			logrus.Warnf(fmt.Sprintf("%s: true is ignored", fieldPath))
		}
		if baremetal.ProvisioningDHCPRange != defaultBM.ProvisioningDHCPRange {
			fieldPath := field.NewPath("Platform", "Baremetal", "ProvisioningDHCPRange")
			logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, baremetal.ProvisioningDHCPRange))
		}

		for i, host := range baremetal.Hosts {
			if host.Name != "" {
				fieldPath := field.NewPath("Platform", "Baremetal", fmt.Sprintf("Hosts[%d]", i), "Name")
				logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, host.Name))
			}
			if host.BMC.Username != "" {
				fieldPath := field.NewPath("Platform", "Baremetal", fmt.Sprintf("Hosts[%d]", i), "BMC", "Username")
				logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, host.BMC.Username))
			}
			if host.BMC.Password != "" {
				fieldPath := field.NewPath("Platform", "Baremetal", fmt.Sprintf("Hosts[%d]", i), "BMC", "Password")
				logrus.Warnf(fmt.Sprintf("%s is ignored", fieldPath))
			}
			if host.BMC.Address != "" {
				fieldPath := field.NewPath("Platform", "Baremetal", fmt.Sprintf("Hosts[%d]", i), "BMC", "Address")
				logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, host.BMC.Address))
			}
			if host.BMC.DisableCertificateVerification {
				fieldPath := field.NewPath("Platform", "Baremetal", fmt.Sprintf("Hosts[%d]", i), "BMC", "DisableCertificateVerification")
				logrus.Warnf(fmt.Sprintf("%s: true is ignored", fieldPath))
			}
			if host.Role != "" {
				fieldPath := field.NewPath("Platform", "Baremetal", fmt.Sprintf("Hosts[%d]", i), "Role")
				logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, host.Role))
			}
			if host.BootMACAddress != "" {
				fieldPath := field.NewPath("Platform", "Baremetal", fmt.Sprintf("Hosts[%d]", i), "BootMACAddress")
				logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, host.BootMACAddress))
			}
			if host.HardwareProfile != "" {
				fieldPath := field.NewPath("Platform", "Baremetal", fmt.Sprintf("Hosts[%d]", i), "HardwareProfile")
				logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, host.HardwareProfile))
			}
			if host.RootDeviceHints != nil {
				fieldPath := field.NewPath("Platform", "Baremetal", fmt.Sprintf("Hosts[%d]", i), "RootDeviceHints")
				logrus.Warnf(fmt.Sprintf("%s is ignored", fieldPath))
			}
			// The default is UEFI. +kubebuilder:validation:Enum="";UEFI;UEFISecureBoot;legacy. Set from generic install config code
			if host.BootMode != "UEFI" {
				fieldPath := field.NewPath("Platform", "Baremetal", fmt.Sprintf("Hosts[%d]", i), "BootMode")
				logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, host.BootMode))
			}
			if host.NetworkConfig != nil {
				fieldPath := field.NewPath("Platform", "Baremetal", fmt.Sprintf("Hosts[%d]", i), "NetworkConfig")
				logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, host.NetworkConfig))
			}
		}

		if baremetal.DefaultMachinePlatform != nil {
			fieldPath := field.NewPath("Platform", "Baremetal", "DefaultMachinePlatform")
			logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, baremetal.DefaultMachinePlatform))
		}
		if baremetal.BootstrapOSImage != "" {
			fieldPath := field.NewPath("Platform", "Baremetal", "BootstrapOSImage")
			logrus.Debugf(fmt.Sprintf("%s: %s is ignored", fieldPath, baremetal.BootstrapOSImage))
		}
		// ClusterOSImage is ignored even in IPI now, so we probably don't need to check it at all.

		if baremetal.BootstrapExternalStaticIP != "" {
			fieldPath := field.NewPath("Platform", "Baremetal", "BootstrapExternalStaticIP")
			logrus.Debugf(fmt.Sprintf("%s: %s is ignored", fieldPath, baremetal.BootstrapExternalStaticIP))
		}
		if baremetal.BootstrapExternalStaticGateway != "" {
			fieldPath := field.NewPath("Platform", "Baremetal", "BootstrapExternalStaticGateway")
			logrus.Debugf(fmt.Sprintf("%s: %s is ignored", fieldPath, baremetal.BootstrapExternalStaticGateway))
		}
	case vsphere.Name:
		vspherePlatform := installConfig.Platform.VSphere

		if vspherePlatform.ClusterOSImage != "" {
			fieldPath := field.NewPath("Platform", "VSphere", "ClusterOSImage")
			logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, vspherePlatform.ClusterOSImage))
		}
		if vspherePlatform.DefaultMachinePlatform != nil && !reflect.DeepEqual(*vspherePlatform.DefaultMachinePlatform, vsphere.MachinePool{}) {
			fieldPath := field.NewPath("Platform", "VSphere", "DefaultMachinePlatform")
			logrus.Warnf(fmt.Sprintf("%s: %v is ignored", fieldPath, vspherePlatform.DefaultMachinePlatform))
		}
		if vspherePlatform.DiskType != "" {
			fieldPath := field.NewPath("Platform", "VSphere", "DiskType")
			logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, vspherePlatform.DiskType))
		}

		if vspherePlatform.LoadBalancer != nil && !reflect.DeepEqual(*vspherePlatform.LoadBalancer, configv1.VSpherePlatformLoadBalancer{}) {
			fieldPath := field.NewPath("Platform", "VSphere", "LoadBalancer")
			logrus.Warnf(fmt.Sprintf("%s: %v is ignored", fieldPath, vspherePlatform.LoadBalancer))
		}

		if len(vspherePlatform.Hosts) > 1 {
			fieldPath := field.NewPath("Platform", "VSphere", "Hosts")
			logrus.Warnf(fmt.Sprintf("%s: %v is ignored", fieldPath, vspherePlatform.Hosts))
		}
	}
	// "External" is the default set from generic install config code
	if installConfig.Publish != "External" {
		fieldPath := field.NewPath("Publish")
		logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, installConfig.Publish))
	}
	if installConfig.CredentialsMode != "" {
		fieldPath := field.NewPath("CredentialsMode")
		logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, installConfig.CredentialsMode))
	}
	if installConfig.BootstrapInPlace != nil && installConfig.BootstrapInPlace.InstallationDisk != "" {
		fieldPath := field.NewPath("BootstrapInPlace", "InstallationDisk")
		logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, installConfig.BootstrapInPlace.InstallationDisk))
	}
}

// GetReplicaCount gets the configured master and worker replicas.
func GetReplicaCount(installConfig *types.InstallConfig) (numMasters, numWorkers int64) {
	numRequiredMasters := int64(0)
	if installConfig.ControlPlane != nil && installConfig.ControlPlane.Replicas != nil {
		numRequiredMasters += *installConfig.ControlPlane.Replicas
	}

	numRequiredWorkers := int64(0)
	for _, worker := range installConfig.Compute {
		if worker.Replicas != nil {
			numRequiredWorkers += *worker.Replicas
		}
	}

	return numRequiredMasters, numRequiredWorkers
}
