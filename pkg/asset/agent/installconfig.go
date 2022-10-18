package agent

import (
	"fmt"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

const (
	installConfigFilename = "install-config.yaml"
)

// OptionalInstallConfig is an InstallConfig where the default is empty, rather
// than generated from running the survey.
type OptionalInstallConfig struct {
	installconfig.InstallConfig
	Supplied bool
}

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

	var foundValidatedDefaultInstallConfig bool

	foundValidatedDefaultInstallConfig, err := a.InstallConfig.Load(f)
	if foundValidatedDefaultInstallConfig && err == nil {
		a.Supplied = true
		if err := a.validateInstallConfig(a.Config).ToAggregate(); err != nil {
			return false, errors.Wrapf(err, "invalid install-config configuration")
		}
	}
	return foundValidatedDefaultInstallConfig, err
}

func (a *OptionalInstallConfig) validateInstallConfig(installConfig *types.InstallConfig) field.ErrorList {
	var allErrs field.ErrorList

	if err := a.validateSupportedPlatforms(installConfig); err != nil {
		allErrs = append(allErrs, err...)
	}

	warnUnusedConfig(installConfig)

	if err := a.validateSNOConfiguration(installConfig); err != nil {
		allErrs = append(allErrs, err...)
	}

	return allErrs
}

func (a *OptionalInstallConfig) validateSupportedPlatforms(installConfig *types.InstallConfig) field.ErrorList {
	var allErrs field.ErrorList

	fieldPath := field.NewPath("Platform")

	if installConfig.Platform.Name() != "" && !IsSupportedPlatform(installConfig.Platform.Name()) {
		allErrs = append(allErrs, field.NotSupported(fieldPath, installConfig.Platform.Name(), SupportedPlatforms))
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

	//  platform None always imply SNO cluster
	if installConfig.Platform.Name() == none.Name {
		if installConfig.Networking.NetworkType != "OVNKubernetes" {
			fieldPath = field.NewPath("Networking", "NetworkType")
			allErrs = append(allErrs, field.Invalid(fieldPath, installConfig.Networking.NetworkType, "Only OVNKubernetes network type is allowed for Single Node OpenShift (SNO) cluster"))
		}

		if *installConfig.ControlPlane.Replicas != 1 {
			fieldPath = field.NewPath("ControlPlane", "Replicas")
			allErrs = append(allErrs, field.Required(fieldPath, fmt.Sprintf("ControlPlane.Replicas must be 1 for %s platform. Found %v", none.Name, *installConfig.ControlPlane.Replicas)))
		}

		if workers != 0 {
			fieldPath = field.NewPath("Compute", "Replicas")
			allErrs = append(allErrs, field.Required(fieldPath, fmt.Sprintf("Total number of Compute.Replicas must be 0 for %s platform. Found %v", none.Name, workers)))
		}
	}

	if installConfig.ControlPlane != nil && *installConfig.ControlPlane.Replicas == 1 && workers == 0 && installConfig.Platform.Name() != none.Name {
		fieldPath = field.NewPath("Platform")
		allErrs = append(allErrs, field.Invalid(fieldPath, installConfig.Platform.Name(), fmt.Sprintf("Platform should be set to %s if the ControlPlane.Replicas is %d and total number of Compute.Replicas is %d", none.Name, *installConfig.ControlPlane.Replicas, workers)))
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
		// kubebuilder:default=amd64. Set from generic install config code
		if compute.Architecture != types.ArchitectureAMD64 && compute.Architecture != "x86_64" {
			fieldPath := field.NewPath(fmt.Sprintf("Compute[%d]", i), "Architecture")
			logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, compute.Architecture))
		}

		if compute.Platform.AWS != nil || compute.Platform.AlibabaCloud != nil || compute.Platform.Azure != nil ||
			compute.Platform.BareMetal != nil || compute.Platform.GCP != nil || compute.Platform.IBMCloud != nil ||
			compute.Platform.Libvirt != nil || compute.Platform.Nutanix != nil || compute.Platform.OpenStack != nil ||
			compute.Platform.Ovirt != nil || compute.Platform.PowerVS != nil || compute.Platform.VSphere != nil {
			fieldPath := field.NewPath(fmt.Sprintf("Compute[%d]", i), "Platform")
			logrus.Warnf(fmt.Sprintf("%s is ignored", fieldPath))
		}
	}

	if installConfig.ControlPlane.Hyperthreading != "Enabled" {
		fieldPath := field.NewPath("ControlPlane", "Hyperthreading")
		logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, installConfig.ControlPlane.Hyperthreading))
	}

	// kubebuilder:default=amd64. Set from generic install config code
	if installConfig.ControlPlane.Architecture != types.ArchitectureAMD64 && installConfig.ControlPlane.Architecture != "x86_64" {
		fieldPath := field.NewPath("ControlPlane", "Architecture")
		logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, installConfig.ControlPlane.Architecture))
	}

	if installConfig.ControlPlane.Platform.AWS != nil || installConfig.ControlPlane.Platform.AlibabaCloud != nil ||
		installConfig.ControlPlane.Platform.Azure != nil || installConfig.ControlPlane.Platform.BareMetal != nil ||
		installConfig.ControlPlane.Platform.GCP != nil || installConfig.ControlPlane.Platform.IBMCloud != nil ||
		installConfig.ControlPlane.Platform.Libvirt != nil || installConfig.ControlPlane.Platform.Nutanix != nil ||
		installConfig.ControlPlane.Platform.OpenStack != nil || installConfig.ControlPlane.Platform.Ovirt != nil ||
		installConfig.ControlPlane.Platform.PowerVS != nil || installConfig.ControlPlane.Platform.VSphere != nil {
		fieldPath := field.NewPath("ControlPlane", "Platform")
		logrus.Warnf(fmt.Sprintf("%s is ignored", fieldPath))
	}

	switch installConfig.Platform.Name() {
	case baremetal.Name:
		baremetal := installConfig.Platform.BareMetal
		// +kubebuilder:default="qemu:///system". Set from generic install config code
		if baremetal.LibvirtURI != "qemu:///system" {
			fieldPath := field.NewPath("Platform", "Baremetal", "LibvirtURI")
			logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, baremetal.LibvirtURI))
		}
		if baremetal.ClusterProvisioningIP != "" {
			fieldPath := field.NewPath("Platform", "Baremetal", "ClusterProvisioningIP")
			logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, baremetal.ClusterProvisioningIP))
		}
		if baremetal.DeprecatedProvisioningHostIP != "" {
			fieldPath := field.NewPath("Platform", "Baremetal", "ProvisioningHostIP")
			logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, baremetal.DeprecatedProvisioningHostIP))
		}
		if baremetal.BootstrapProvisioningIP != "" {
			fieldPath := field.NewPath("Platform", "Baremetal", "BootstrapProvisioningIP")
			logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, baremetal.BootstrapProvisioningIP))
		}
		if baremetal.ExternalBridge != "" {
			fieldPath := field.NewPath("Platform", "Baremetal", "ExternalBridge")
			logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, baremetal.ExternalBridge))
		}
		if baremetal.ExternalMACAddress != "" {
			fieldPath := field.NewPath("Platform", "Baremetal", "ExternalMACAddress")
			logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, baremetal.ExternalMACAddress))
		}
		// +kubebuilder:default=Managed
		if baremetal.ProvisioningNetwork != "Managed" {
			fieldPath := field.NewPath("Platform", "Baremetal", "ProvisioningNetwork")
			logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, baremetal.ProvisioningNetwork))
		}
		if baremetal.ProvisioningBridge != "" {
			fieldPath := field.NewPath("Platform", "Baremetal", "ProvisioningBridge")
			logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, baremetal.ProvisioningBridge))
		}
		if baremetal.ProvisioningMACAddress != "" {
			fieldPath := field.NewPath("Platform", "Baremetal", "ProvisioningMACAddress")
			logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, baremetal.ProvisioningMACAddress))
		}
		if baremetal.ProvisioningNetworkInterface != "" {
			fieldPath := field.NewPath("Platform", "Baremetal", "ProvisioningNetworkInterface")
			logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, baremetal.ProvisioningNetworkInterface))
		}
		if baremetal.ProvisioningNetworkCIDR.String() != "" {
			fieldPath := field.NewPath("Platform", "Baremetal", "ProvisioningNetworkCIDR")
			logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, baremetal.ProvisioningNetworkCIDR))
		}
		if baremetal.DeprecatedProvisioningDHCPExternal {
			fieldPath := field.NewPath("Platform", "Baremetal", "ProvisioningDHCPExternal")
			logrus.Warnf(fmt.Sprintf("%s: true is ignored", fieldPath))
		}
		if baremetal.ProvisioningDHCPRange != "" {
			fieldPath := field.NewPath("Platform", "Baremetal", "ProvisioningDHCPRange")
			logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, baremetal.ProvisioningDHCPRange))
		}

		for i, host := range baremetal.Hosts {
			if host.BMC.Username != "" {
				fieldPath := field.NewPath("Platform", "Baremetal", fmt.Sprintf("Hosts[%d]", i), "BMC", "Username")
				logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, host.BMC.Username))
			}
			if host.BMC.Password != "" {
				fieldPath := field.NewPath("Platform", "Baremetal", fmt.Sprintf("Hosts[%d]", i), "BMC", "Password")
				logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, host.BMC.Password))
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
			logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, baremetal.BootstrapOSImage))
		}
		if baremetal.ClusterOSImage != "" {
			fieldPath := field.NewPath("Platform", "Baremetal", "ClusterOSImage")
			logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, baremetal.ClusterOSImage))
		}
		if baremetal.BootstrapExternalStaticIP != "" {
			fieldPath := field.NewPath("Platform", "Baremetal", "BootstrapExternalStaticIP")
			logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, baremetal.BootstrapExternalStaticIP))
		}
		if baremetal.BootstrapExternalStaticGateway != "" {
			fieldPath := field.NewPath("Platform", "Baremetal", "BootstrapExternalStaticGateway")
			logrus.Warnf(fmt.Sprintf("%s: %s is ignored", fieldPath, baremetal.BootstrapExternalStaticGateway))
		}
	}
}
