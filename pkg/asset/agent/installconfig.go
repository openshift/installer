package agent

import (
	"context"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/validation/field"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/api/features"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/releaseimage"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
	baremetaldefaults "github.com/openshift/installer/pkg/types/baremetal/defaults"
	baremetalvalidation "github.com/openshift/installer/pkg/types/baremetal/validation"
	"github.com/openshift/installer/pkg/types/external"
	"github.com/openshift/installer/pkg/types/none"
	nonevalidation "github.com/openshift/installer/pkg/types/none/validation"
	"github.com/openshift/installer/pkg/types/validation"
	"github.com/openshift/installer/pkg/types/vsphere"
)

const (
	// InstallConfigFilename is the file containing the install-config.
	InstallConfigFilename = "install-config.yaml"
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
func (a *OptionalInstallConfig) Generate(_ context.Context, parents asset.Parents) error {
	// Just generate an empty install config, since we have no dependencies.
	return nil
}

// Load returns the installconfig from disk.
func (a *OptionalInstallConfig) Load(f asset.FileFetcher) (bool, error) {
	ctx := context.TODO()
	found, err := a.LoadFromFile(f)
	if found && err == nil {
		a.Supplied = true
		if err := a.validateInstallConfig(ctx, a.Config).ToAggregate(); err != nil {
			return false, errors.Wrapf(err, "invalid install-config configuration")
		}
		if err := a.RecordFile(); err != nil {
			return false, err
		}
	}
	return found, err
}

func (a *OptionalInstallConfig) validateInstallConfig(ctx context.Context, installConfig *types.InstallConfig) field.ErrorList {
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
	if err := a.validateReleaseArch(ctx, installConfig); err != nil {
		allErrs = append(allErrs, err...)
	}

	if installConfig.FeatureSet != configv1.Default {
		allErrs = append(allErrs, field.NotSupported(field.NewPath("FeatureSet"), installConfig.FeatureSet, []string{string(configv1.Default)}))
	}

	warnUnusedConfig(installConfig)

	numMasters, numWorkers := GetReplicaCount(installConfig)
	logrus.Infof(fmt.Sprintf("Configuration has %d master replicas and %d worker replicas", numMasters, numWorkers))

	if err := a.validateControlPlaneConfiguration(installConfig); err != nil {
		allErrs = append(allErrs, err...)
	}

	if err := a.validateSNOConfiguration(installConfig); err != nil {
		allErrs = append(allErrs, err...)
	}

	return allErrs
}

func (a *OptionalInstallConfig) validateSupportedPlatforms(installConfig *types.InstallConfig) field.ErrorList {
	allErrs := ValidateSupportedPlatforms(installConfig.Platform, string(installConfig.ControlPlane.Architecture))
	return append(allErrs, a.validatePlatformsByName(installConfig)...)
}

// ValidateSupportedPlatforms verifies if the specified platform/arch is supported or not.
func ValidateSupportedPlatforms(platform types.Platform, controlPlaneArch string) field.ErrorList {
	var allErrs field.ErrorList

	fieldPath := field.NewPath("Platform")

	if platform.Name() != "" && !IsSupportedPlatform(HivePlatformType(platform)) {
		allErrs = append(allErrs, field.NotSupported(fieldPath, platform.Name(), SupportedInstallerPlatforms()))
	}
	if platform.Name() != none.Name && controlPlaneArch == types.ArchitecturePPC64LE {
		allErrs = append(allErrs, field.Invalid(fieldPath, platform.Name(), fmt.Sprintf("CPU architecture \"%s\" only supports platform \"%s\".", types.ArchitecturePPC64LE, none.Name)))
	}
	if platform.Name() != none.Name && controlPlaneArch == types.ArchitectureS390X {
		allErrs = append(allErrs, field.Invalid(fieldPath, platform.Name(), fmt.Sprintf("CPU architecture \"%s\" only supports platform \"%s\".", types.ArchitectureS390X, none.Name)))
	}
	return allErrs
}

func (a *OptionalInstallConfig) validatePlatformsByName(installConfig *types.InstallConfig) field.ErrorList {
	var allErrs field.ErrorList

	if installConfig.Platform.Name() == external.Name {
		if installConfig.Platform.External.PlatformName == ExternalPlatformNameOci &&
			installConfig.Platform.External.CloudControllerManager != external.CloudControllerManagerTypeExternal {
			fieldPath := field.NewPath("Platform", "External", "CloudControllerManager")
			allErrs = append(allErrs, field.Invalid(fieldPath, installConfig.Platform.External.CloudControllerManager,
				fmt.Sprintf("When using external %s platform, %s must be set to %s",
					ExternalPlatformNameOci, fieldPath, external.CloudControllerManagerTypeExternal)))
		}
	}

	if installConfig.Platform.Name() == vsphere.Name {
		allErrs = append(allErrs, a.validateVSpherePlatform(installConfig)...)
	}

	if installConfig.Platform.Name() == baremetal.Name {
		allErrs = append(allErrs, a.validateBMCConfig(installConfig)...)
	}

	if installConfig.Platform.Name() == none.Name {
		allErrs = append(allErrs, a.validateFencingCredentials(installConfig)...)
	}

	return allErrs
}

func (a *OptionalInstallConfig) validateReleaseArch(ctx context.Context, installConfig *types.InstallConfig) field.ErrorList {
	var allErrs field.ErrorList

	fieldPath := field.NewPath("ControlPlane", "Architecture")
	releaseImage := &releaseimage.Image{}
	asseterr := releaseImage.Generate(ctx, asset.Parents{})
	if asseterr != nil {
		allErrs = append(allErrs, field.InternalError(fieldPath, asseterr))
	}
	releaseArch, err := DetermineReleaseImageArch(installConfig.PullSecret, releaseImage.PullSpec)
	if err != nil {
		logrus.Warnf("Unable to validate the release image architecture, skipping validation")
	} else {
		// Validate that the release image supports the install-config architectures.
		switch releaseArch {
		// Check the release image to see if it is multi.
		case "multi":
			logrus.Debugf("multi architecture release image %s found, all archs supported", releaseImage.PullSpec)
		// If the release image isn't multi, then its single arch, and it must match the cpu architecture.
		case string(installConfig.ControlPlane.Architecture):
			logrus.Debugf("Supported architecture %s found for the release image: %s", installConfig.ControlPlane.Architecture, releaseImage.PullSpec)
		default:
			allErrs = append(allErrs, field.Forbidden(fieldPath, fmt.Sprintf("unsupported release image architecture. ControlPlane Arch: %s doesn't match Release Image Arch: %s", installConfig.ControlPlane.Architecture, releaseArch)))
		}
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
	case types.ArchitectureS390X:
	default:
		allErrs = append(allErrs, field.NotSupported(fieldPath, installConfig.ControlPlane.Architecture, []string{types.ArchitectureAMD64, types.ArchitectureARM64, types.ArchitecturePPC64LE, types.ArchitectureS390X}))
	}

	for i, compute := range installConfig.Compute {
		fieldPath := field.NewPath(fmt.Sprintf("Compute[%d]", i), "Architecture")

		switch string(compute.Architecture) {
		case types.ArchitectureAMD64:
		case types.ArchitectureARM64:
		case types.ArchitecturePPC64LE:
		case types.ArchitectureS390X:
		default:
			allErrs = append(allErrs, field.NotSupported(fieldPath, compute.Architecture, []string{types.ArchitectureAMD64, types.ArchitectureARM64, types.ArchitecturePPC64LE, types.ArchitectureS390X}))
		}
	}

	return allErrs
}

func (a *OptionalInstallConfig) validateControlPlaneConfiguration(installConfig *types.InstallConfig) field.ErrorList {
	var allErrs field.ErrorList
	var fieldPath *field.Path

	if installConfig.ControlPlane != nil {
		if *installConfig.ControlPlane.Replicas < 1 || *installConfig.ControlPlane.Replicas > 5 || (installConfig.Arbiter == nil && *installConfig.ControlPlane.Replicas == 2) {
			fieldPath = field.NewPath("ControlPlane", "Replicas")
			supportedControlPlaneRange := "to 5, 4, 3, or 1"
			if installConfig.EnabledFeatureGates().Enabled(features.FeatureGateHighlyAvailableArbiter) {
				supportedControlPlaneRange = "between 1 and 5"
			}
			allErrs = append(allErrs, field.Invalid(fieldPath, installConfig.ControlPlane.Replicas, fmt.Sprintf("ControlPlane.Replicas can only be set %s. Found %v", supportedControlPlaneRange, *installConfig.ControlPlane.Replicas)))
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
	vcenterServers := map[string]bool{}
	userProvidedCredentials := false
	for _, vcenter := range vspherePlatform.VCenters {
		vcenterServers[vcenter.Server] = true

		// If any one of the required credential values is entered, then the user is choosing to enter credentials
		if VCenterCredentialsAreProvided(vcenter) {
			// Then check all required credential values are filled
			userProvidedCredentials = true
			message := "All credential fields are required if any one is specified"
			if vcenter.Server == "" {
				fieldPath := field.NewPath("Platform", "VSphere", "vcenter")
				allErrs = append(allErrs, field.Required(fieldPath, message))
			}
			if vcenter.Username == "" {
				fieldPath := field.NewPath("Platform", "VSphere", "user")
				if vspherePlatform.DeprecatedVCenter != "" || vspherePlatform.DeprecatedPassword != "" || vspherePlatform.DeprecatedDatacenter != "" {
					fieldPath = field.NewPath("Platform", "VSphere", "username")
				}
				allErrs = append(allErrs, field.Required(fieldPath, message))
			}
			if vcenter.Password == "" {
				fieldPath := field.NewPath("Platform", "VSphere", "password")
				allErrs = append(allErrs, field.Required(fieldPath, message))
			}
			if len(vcenter.Datacenters) == 0 {
				fieldPath := field.NewPath("Platform", "VSphere", "datacenter")
				allErrs = append(allErrs, field.Required(fieldPath, message))
			}
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

// ClusterNamespace returns the namespace of the cluster.
func (a *OptionalInstallConfig) ClusterNamespace() string {
	if a.Config != nil && a.Config.ObjectMeta.Namespace != "" {
		return a.Config.ObjectMeta.Namespace
	}
	return ""
}

// GetBaremetalHosts gets the hosts defined for a baremetal platform.
func (a *OptionalInstallConfig) GetBaremetalHosts() []*baremetal.Host {
	if a.Config != nil && a.Config.Platform.Name() == baremetal.Name {
		return a.Config.Platform.BareMetal.Hosts
	}
	return nil
}

func (a *OptionalInstallConfig) validateBMCConfig(installConfig *types.InstallConfig) field.ErrorList {
	var allErrs field.ErrorList

	bmcConfigured := false
	for _, host := range installConfig.Platform.BareMetal.Hosts {
		if host.BMC.Address == "" {
			continue
		}
		bmcConfigured = true
	}

	if bmcConfigured {
		fieldPath := field.NewPath("Platform", "BareMetal")
		allErrs = append(allErrs, baremetalvalidation.ValidateProvisioningNetworking(installConfig.Platform.BareMetal, installConfig.Networking, fieldPath)...)
	}

	return allErrs
}

func (a *OptionalInstallConfig) validateFencingCredentials(installConfig *types.InstallConfig) field.ErrorList {
	var allErrs field.ErrorList

	fieldPath := field.NewPath("Platform", "None")
	allErrs = append(allErrs, nonevalidation.ValidateFencingCredentials(installConfig.Platform.None.FencingCredentials, fieldPath)...)

	return allErrs
}

func warnUnusedConfig(installConfig *types.InstallConfig) {
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
		if baremetal.BootstrapProvisioningIP != defaultBM.BootstrapProvisioningIP {
			fieldPath := field.NewPath("Platform", "Baremetal", "BootstrapProvisioningIP")
			logrus.Debugf(fmt.Sprintf("%s: %s is ignored", fieldPath, baremetal.BootstrapProvisioningIP))
		}
		if baremetal.ExternalBridge != defaultBM.ExternalBridge {
			fieldPath := field.NewPath("Platform", "Baremetal", "ExternalBridge")
			logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, baremetal.ExternalBridge))
		}
		if baremetal.ProvisioningBridge != defaultBM.ProvisioningBridge {
			fieldPath := field.NewPath("Platform", "Baremetal", "ProvisioningBridge")
			logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, baremetal.ProvisioningBridge))
		}

		for i, host := range baremetal.Hosts {
			// The default is UEFI. +kubebuilder:validation:Enum="";UEFI;UEFISecureBoot;legacy. Set from generic install config code
			if host.BootMode != "UEFI" {
				fieldPath := field.NewPath("Platform", "Baremetal", fmt.Sprintf("Hosts[%d]", i), "BootMode")
				logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, host.BootMode))
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
