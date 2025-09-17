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
	"github.com/openshift/installer/pkg/types/nutanix"
	nutanixvalidation "github.com/openshift/installer/pkg/types/nutanix/validation"
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

	warnUnusedConfig(installConfig)

	numMasters, numArbiters, numWorkers := GetReplicaCount(installConfig)
	logrus.Infof("Configuration has %d master replicas, %d arbiter replicas, and %d worker replicas", numMasters, numArbiters, numWorkers)

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

	fieldPath := field.NewPath("platform")

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
			fieldPath := field.NewPath("platform", "external", "cloudControllerManager")
			allErrs = append(allErrs, field.Invalid(fieldPath, installConfig.Platform.External.CloudControllerManager,
				fmt.Sprintf("When using external %s platform, %s must be set to %s",
					ExternalPlatformNameOci, fieldPath, external.CloudControllerManagerTypeExternal)))
		}
	}

	if installConfig.Platform.Name() == vsphere.Name {
		allErrs = append(allErrs, a.validateVSpherePlatform(installConfig)...)
	}

	if installConfig.Platform.Name() == nutanix.Name {
		allErrs = append(allErrs, a.validateNutanixPlatform(installConfig)...)
	}

	if installConfig.Platform.Name() == baremetal.Name {
		allErrs = append(allErrs, a.validateBMCConfig(installConfig)...)
	}

	return allErrs
}

func (a *OptionalInstallConfig) validateReleaseArch(ctx context.Context, installConfig *types.InstallConfig) field.ErrorList {
	var allErrs field.ErrorList

	fieldPath := field.NewPath("controlPlane", "architecture")
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

	fieldPath := field.NewPath("controlPlane", "architecture")

	switch string(installConfig.ControlPlane.Architecture) {
	case types.ArchitectureAMD64:
	case types.ArchitectureARM64:
	case types.ArchitecturePPC64LE:
	case types.ArchitectureS390X:
	default:
		allErrs = append(allErrs, field.NotSupported(fieldPath, installConfig.ControlPlane.Architecture, []string{types.ArchitectureAMD64, types.ArchitectureARM64, types.ArchitecturePPC64LE, types.ArchitectureS390X}))
	}

	computePath := field.NewPath("compute")
	for i, compute := range installConfig.Compute {
		fieldPath := computePath.Index(i).Child("architecture")

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
func isValidTNFCluster(installConfig *types.InstallConfig) bool {
	// Check basic TNF requirements
	if installConfig.ControlPlane == nil ||
		installConfig.ControlPlane.Replicas == nil ||
		*installConfig.ControlPlane.Replicas != 2 ||
		!installConfig.EnabledFeatureGates().Enabled(features.FeatureGateDualReplica) {
		return false
	}

	// Check fencing credentials are provided
	if installConfig.ControlPlane.Fencing == nil ||
		len(installConfig.ControlPlane.Fencing.Credentials) == 0 {
		return false
	}

	// Platform validation is handled by validateFencingForPlatform() in general validation
	return true
}

func (a *OptionalInstallConfig) validateControlPlaneConfiguration(installConfig *types.InstallConfig) field.ErrorList {
	var allErrs field.ErrorList
	var fieldPath *field.Path

	if installConfig.ControlPlane != nil {
		if *installConfig.ControlPlane.Replicas < 1 ||
			*installConfig.ControlPlane.Replicas > 5 ||
			(installConfig.Arbiter == nil && *installConfig.ControlPlane.Replicas == 2 && !isValidTNFCluster(installConfig)) {

			fieldPath = field.NewPath("controlPlane", "replicas")
			supportedControlPlaneRange := []string{"3", "1", "4", "5"}
			if installConfig.EnabledFeatureGates().Enabled(features.FeatureGateHighlyAvailableArbiter) {
				supportedControlPlaneRange = append(supportedControlPlaneRange, "2 (with arbiter)")
			}
			if installConfig.EnabledFeatureGates().Enabled(features.FeatureGateDualReplica) {
				supportedControlPlaneRange = append(supportedControlPlaneRange, "2 (with fencing)")
			}
			allErrs = append(allErrs, field.NotSupported(fieldPath, installConfig.ControlPlane.Replicas, supportedControlPlaneRange))
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
				fieldPath = field.NewPath("networking", "networkType")
				allErrs = append(allErrs, field.Invalid(fieldPath, installConfig.Networking.NetworkType, "Only OVNKubernetes network type is allowed for Single Node OpenShift (SNO) cluster"))
			}
			if installConfig.Platform.Name() != none.Name && installConfig.Platform.Name() != external.Name {
				fieldPath = field.NewPath("platform")
				allErrs = append(allErrs, field.Invalid(fieldPath, installConfig.Platform.Name(), fmt.Sprintf("Only platform %s and %s supports 1 ControlPlane and 0 Compute nodes", none.Name, external.Name)))
			}
		} else {
			fieldPath = field.NewPath("compute", "replicas")
			allErrs = append(allErrs, field.Forbidden(fieldPath, fmt.Sprintf("Total number of compute replicas must be 0 when controlPlane.replicas is 1 for platform %s or %s. Found %v", none.Name, external.Name, workers)))
		}
		if installConfig.Arbiter != nil && installConfig.Arbiter.Replicas != nil && *installConfig.Arbiter.Replicas > 0 {
			fieldPath = field.NewPath("arbiter", "replicas")
			allErrs = append(allErrs, field.Forbidden(fieldPath, fmt.Sprintf("Total number of arbiter replicas must be 0 when controlPlane.replicas is 1 for Single Node Openshift (SNO) cluster. Found %d", *installConfig.Arbiter.Replicas)))
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
	platformPath := field.NewPath("platform", "vsphere")
	for i, vcenter := range vspherePlatform.VCenters {
		vcenterPath := platformPath.Child("vcenters").Index(i)
		vcenterServers[vcenter.Server] = true

		// If any one of the required credential values is entered, then the user is choosing to enter credentials
		if VCenterCredentialsAreProvided(vcenter) {
			// Then check all required credential values are filled
			userProvidedCredentials = true
			message := "All credential fields are required if any one is specified"
			if vcenter.Server == "" {
				fieldPath := vcenterPath.Child("server")
				allErrs = append(allErrs, field.Required(fieldPath, message))
			}
			if vcenter.Username == "" {
				fieldPath := vcenterPath.Child("user")
				if vspherePlatform.DeprecatedVCenter != "" || vspherePlatform.DeprecatedPassword != "" || vspherePlatform.DeprecatedDatacenter != "" {
					fieldPath = field.NewPath("platform", "vsphere", "username")
				}
				allErrs = append(allErrs, field.Required(fieldPath, message))
			}
			if vcenter.Password == "" {
				fieldPath := vcenterPath.Child("password")
				allErrs = append(allErrs, field.Required(fieldPath, message))
			}
			if len(vcenter.Datacenters) == 0 {
				fieldPath := vcenterPath.Child("datacenter")
				allErrs = append(allErrs, field.Required(fieldPath, message))
			}
		}
	}

	for i, failureDomain := range vspherePlatform.FailureDomains {
		// Although folder is optional in IPI/UPI, it must be set for agent-based installs.
		// If it is not set, assisted-service will set a placeholder value for folder:
		// "/datacenterplaceholder/vm/folderplaceholder"
		//
		// When assisted-service generates the install-config for the cluster, it will fail
		// validation because the placeholder value's datacenter name may not match
		// the datacenter set in the failureDomain in the install-config.yaml submitted
		// to the agent-based create image command.
		if failureDomain.Topology.Folder == "" && userProvidedCredentials {
			fieldPath := platformPath.Child("failureDomains").Index(i).Child("topology", "folder")
			allErrs = append(allErrs, field.Required(fieldPath, "must specify a folder for agent-based installs"))
		}
	}

	return allErrs
}

func (a *OptionalInstallConfig) validateNutanixPlatform(installConfig *types.InstallConfig) field.ErrorList {
	fldPath := field.NewPath("platform", "nutanix")
	return nutanixvalidation.ValidatePlatform(installConfig.Platform.Nutanix, fldPath, installConfig, true)
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
		fieldPath := field.NewPath("platform", "baremetal")
		allErrs = append(allErrs, baremetalvalidation.ValidateProvisioningNetworking(installConfig.Platform.BareMetal, installConfig.Networking, fieldPath)...)
	}

	return allErrs
}

// nolint:gocyclo
func warnUnusedConfig(installConfig *types.InstallConfig) {
	for i, compute := range installConfig.Compute {
		computePath := field.NewPath("compute").Index(i)
		if compute.Hyperthreading != "Enabled" {
			fieldPath := computePath.Child("hyperthreading")
			logrus.Warnf("%s: %s is ignored", fieldPath, compute.Hyperthreading)
		}

		if compute.Platform != (types.MachinePoolPlatform{}) {
			fieldPath := computePath.Child("platform")
			logrus.Warnf("%s is ignored", fieldPath)
		}
	}

	if installConfig.ControlPlane.Hyperthreading != "Enabled" {
		fieldPath := field.NewPath("controlPlane", "hyperthreading")
		logrus.Warnf("%s: %s is ignored", fieldPath, installConfig.ControlPlane.Hyperthreading)
	}

	if installConfig.ControlPlane.Platform != (types.MachinePoolPlatform{}) {
		fieldPath := field.NewPath("controlPlane", "platform")
		logrus.Warnf("%s is ignored", fieldPath)
	}

	switch installConfig.Platform.Name() {

	case baremetal.Name:
		defaultIc := &types.InstallConfig{Platform: types.Platform{BareMetal: &baremetal.Platform{}}}
		baremetaldefaults.SetPlatformDefaults(defaultIc.Platform.BareMetal, defaultIc)

		baremetal := installConfig.Platform.BareMetal
		defaultBM := defaultIc.Platform.BareMetal
		bmPath := field.NewPath("platform", "baremetal")
		// Compare values from generic installconfig code to check for changes
		if baremetal.LibvirtURI != defaultBM.LibvirtURI {
			fieldPath := bmPath.Child("libvirtURI")
			logrus.Debugf("%s: %s is ignored", fieldPath, baremetal.LibvirtURI)
		}
		if baremetal.BootstrapProvisioningIP != defaultBM.BootstrapProvisioningIP {
			fieldPath := bmPath.Child("bootstrapProvisioningIP")
			logrus.Debugf("%s: %s is ignored", fieldPath, baremetal.BootstrapProvisioningIP)
		}
		if baremetal.ExternalBridge != defaultBM.ExternalBridge {
			fieldPath := bmPath.Child("externalBridge")
			logrus.Warnf("%s: %s is ignored", fieldPath, baremetal.ExternalBridge)
		}
		if baremetal.ProvisioningBridge != defaultBM.ProvisioningBridge {
			fieldPath := bmPath.Child("provisioningBridge")
			logrus.Warnf("%s: %s is ignored", fieldPath, baremetal.ProvisioningBridge)
		}

		for i, host := range baremetal.Hosts {
			// The default is UEFI. +kubebuilder:validation:Enum="";UEFI;UEFISecureBoot;legacy. Set from generic install config code
			if host.BootMode != "UEFI" {
				fieldPath := bmPath.Child("hosts").Index(i).Child("bootMode")
				logrus.Warnf("%s: %s is ignored", fieldPath, host.BootMode)
			}
		}

		if baremetal.DefaultMachinePlatform != nil {
			fieldPath := bmPath.Child("defaultMachinePlatform")
			logrus.Warnf("%s: %s is ignored", fieldPath, baremetal.DefaultMachinePlatform)
		}
		if baremetal.BootstrapOSImage != "" {
			fieldPath := bmPath.Child("bootstrapOSImage")
			logrus.Debugf("%s: %s is ignored", fieldPath, baremetal.BootstrapOSImage)
		}
		// ClusterOSImage is ignored even in IPI now, so we probably don't need to check it at all.

		if baremetal.BootstrapExternalStaticIP != "" {
			fieldPath := bmPath.Child("bootstrapExternalStaticIP")
			logrus.Debugf("%s: %s is ignored", fieldPath, baremetal.BootstrapExternalStaticIP)
		}
		if baremetal.BootstrapExternalStaticGateway != "" {
			fieldPath := bmPath.Child("bootstrapExternalStaticGateway")
			logrus.Debugf("%s: %s is ignored", fieldPath, baremetal.BootstrapExternalStaticGateway)
		}
	case vsphere.Name:
		vspherePlatform := installConfig.Platform.VSphere
		vsPath := field.NewPath("platform", "vsphere")

		if vspherePlatform.ClusterOSImage != "" {
			fieldPath := vsPath.Child("clusterOSImage")
			logrus.Warnf("%s: %s is ignored", fieldPath, vspherePlatform.ClusterOSImage)
		}
		if vspherePlatform.DefaultMachinePlatform != nil && !reflect.DeepEqual(*vspherePlatform.DefaultMachinePlatform, vsphere.MachinePool{}) {
			fieldPath := vsPath.Child("defaultMachinePlatform")
			logrus.Warnf("%s: %v is ignored", fieldPath, vspherePlatform.DefaultMachinePlatform)
		}
		if vspherePlatform.DiskType != "" {
			fieldPath := vsPath.Child("diskType")
			logrus.Warnf("%s: %s is ignored", fieldPath, vspherePlatform.DiskType)
		}

		if vspherePlatform.LoadBalancer != nil && !reflect.DeepEqual(*vspherePlatform.LoadBalancer, configv1.VSpherePlatformLoadBalancer{}) {
			fieldPath := vsPath.Child("loadBalancer")
			logrus.Warnf("%s: %v is ignored", fieldPath, vspherePlatform.LoadBalancer)
		}

		if len(vspherePlatform.Hosts) > 1 {
			fieldPath := vsPath.Child("hosts")
			logrus.Warnf("%s: %v is ignored", fieldPath, vspherePlatform.Hosts)
		}
	case nutanix.Name:
		ntxPlatform := installConfig.Platform.Nutanix
		fieldPath := field.NewPath("Platform", "Nutanix")

		if ntxPlatform.ClusterOSImage != "" {
			logrus.Warnf("%s: %s is ignored", fieldPath.Child("ClusterOSImage").String(), ntxPlatform.ClusterOSImage)
		}
		if ntxPlatform.PreloadedOSImageName != "" {
			logrus.Warnf("%s: %s is ignored", fieldPath.Child("PreloadedOSImageName").String(), ntxPlatform.PreloadedOSImageName)
		}
		if ntxPlatform.DefaultMachinePlatform != nil && !reflect.DeepEqual(*ntxPlatform.DefaultMachinePlatform, nutanix.MachinePool{}) {
			logrus.Warnf("%s: %v is ignored", fieldPath.Child("DefaultMachinePlatform").String(), *ntxPlatform.DefaultMachinePlatform)
		}
		if ntxPlatform.LoadBalancer != nil && !reflect.DeepEqual(*ntxPlatform.LoadBalancer, configv1.NutanixPlatformLoadBalancer{}) {
			logrus.Warnf("%s: %v is ignored", fieldPath.Child("LoadBalancer").String(), *ntxPlatform.LoadBalancer)
		}
		if ntxPlatform.PrismAPICallTimeout != nil {
			logrus.Warnf("%s: %v is ignored", fieldPath.Child("PrismAPICallTimeout").String(), *ntxPlatform.PrismAPICallTimeout)
		}
		if ntxPlatform.FailureDomains != nil {
			logrus.Warnf("%s: %v is ignored", fieldPath.Child("FailureDomains").String(), ntxPlatform.FailureDomains)
		}
	}

	// "External" is the default set from generic install config code
	if installConfig.Publish != "External" {
		fieldPath := field.NewPath("publish")
		logrus.Warnf("%s: %s is ignored", fieldPath, installConfig.Publish)
	}
	if installConfig.CredentialsMode != "" {
		fieldPath := field.NewPath("credentialsMode")
		logrus.Warnf("%s: %s is ignored", fieldPath, installConfig.CredentialsMode)
	}
	if installConfig.BootstrapInPlace != nil && installConfig.BootstrapInPlace.InstallationDisk != "" {
		fieldPath := field.NewPath("bootstrapInPlace", "installationDisk")
		logrus.Warnf("%s: %s is ignored", fieldPath, installConfig.BootstrapInPlace.InstallationDisk)
	}
}

// GetReplicaCount gets the configured master, arbiter and worker replicas.
func GetReplicaCount(installConfig *types.InstallConfig) (numMasters, numArbiters, numWorkers int64) {
	numRequiredMasters := int64(0)
	if installConfig.ControlPlane != nil && installConfig.ControlPlane.Replicas != nil {
		numRequiredMasters += *installConfig.ControlPlane.Replicas
	}
	numRequiredArbiters := int64(0)
	if installConfig.Arbiter != nil && installConfig.Arbiter.Replicas != nil {
		numRequiredArbiters += *installConfig.Arbiter.Replicas
	}

	numRequiredWorkers := int64(0)
	for _, worker := range installConfig.Compute {
		if worker.Replicas != nil {
			numRequiredWorkers += *worker.Replicas
		}
	}

	return numRequiredMasters, numRequiredArbiters, numRequiredWorkers
}
