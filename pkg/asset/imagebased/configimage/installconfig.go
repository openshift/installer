package configimage

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/validation/field"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/validation"
)

const (
	// InstallConfigFilename is the file containing the install-config.
	InstallConfigFilename = "install-config.yaml"
)

// InstallConfig is an InstallConfig where the default is empty, rather
// than generated from running the survey.
type InstallConfig struct {
	installconfig.AssetBase
}

var _ asset.WritableAsset = (*InstallConfig)(nil)

// Dependencies returns all of the dependencies directly needed by an
// InstallConfig asset.
func (i *InstallConfig) Dependencies() []asset.Asset {
	// Return no dependencies for the Image-based install config, because it is
	// optional. We don't need to run the survey if it doesn't exist, since it
	// does not support the `none` platform.
	return []asset.Asset{}
}

// Generate generates the install-config.yaml file.
func (i *InstallConfig) Generate(_ context.Context, parents asset.Parents) error {
	// Just generate an empty install config, since we have no dependencies.
	return nil
}

// Load returns the installconfig from disk.
func (i *InstallConfig) Load(f asset.FileFetcher) (bool, error) {
	found, err := i.LoadFromFile(f)
	if found && err == nil {
		if err := i.validateInstallConfig(i.Config).ToAggregate(); err != nil {
			return false, fmt.Errorf("invalid install-config configuration: %w", err)
		}
		if err := i.RecordFile(); err != nil {
			return false, err
		}
	}
	return found, err
}

func (i *InstallConfig) validateInstallConfig(installConfig *types.InstallConfig) field.ErrorList {
	var allErrs field.ErrorList
	if err := validation.ValidateInstallConfig(i.Config, true); err != nil {
		allErrs = append(allErrs, err...)
	}

	if err := i.validateSupportedPlatforms(installConfig); err != nil {
		allErrs = append(allErrs, err...)
	}

	if installConfig.FeatureSet != configv1.Default {
		allErrs = append(allErrs, field.NotSupported(field.NewPath("FeatureSet"), installConfig.FeatureSet, []string{string(configv1.Default)}))
	}

	warnUnusedConfig(installConfig)

	if err := i.validateSNOConfiguration(installConfig); err != nil {
		allErrs = append(allErrs, err...)
	}

	return allErrs
}

func (i *InstallConfig) validateSupportedPlatforms(installConfig *types.InstallConfig) field.ErrorList {
	var allErrs field.ErrorList

	fieldPath := field.NewPath("Platform")

	if installConfig.Platform.Name() != "" && installConfig.Platform.Name() != none.Name {
		allErrs = append(allErrs, field.NotSupported(fieldPath, installConfig.Platform.Name(), []string{none.Name}))
	}

	return allErrs
}

func (i *InstallConfig) validateSNOConfiguration(installConfig *types.InstallConfig) field.ErrorList {
	var allErrs field.ErrorList
	var fieldPath *field.Path

	controlPlaneReplicas := *installConfig.ControlPlane.Replicas
	if installConfig.ControlPlane != nil && controlPlaneReplicas != 1 {
		fieldPath = field.NewPath("ControlPlane", "Replicas")
		allErrs = append(allErrs, field.Required(fieldPath, fmt.Sprintf("Only Single Node OpenShift (SNO) is supported, total number of ControlPlane.Replicas must be 1. Found %v", controlPlaneReplicas)))
	}

	var workers int
	for _, worker := range installConfig.Compute {
		workers += int(*worker.Replicas)
	}
	if workers != 0 {
		fieldPath = field.NewPath("Compute", "Replicas")
		allErrs = append(allErrs, field.Required(fieldPath, fmt.Sprintf("Total number of Compute.Replicas must be 0 when ControlPlane.Replicas is 1 for platform %s. Found %v", none.Name, workers)))
	}

	if installConfig.Networking.NetworkType != "OVNKubernetes" {
		fieldPath = field.NewPath("Networking", "NetworkType")
		allErrs = append(allErrs, field.Invalid(fieldPath, installConfig.Networking.NetworkType, "Only OVNKubernetes network type is allowed for Single Node OpenShift (SNO) cluster"))
	}

	machineNetworksCount := len(installConfig.Networking.MachineNetwork)
	if machineNetworksCount != 1 {
		fieldPath = field.NewPath("Networking", "MachineNetwork")
		allErrs = append(allErrs, field.TooMany(fieldPath, machineNetworksCount, 1))
	}

	return allErrs
}

// ClusterName returns the name of the cluster, or a default name if no
// InstallConfig is supplied.
func (i *InstallConfig) ClusterName() string {
	if i.Config != nil && i.Config.ObjectMeta.Name != "" {
		return i.Config.ObjectMeta.Name
	}
	return "imagebased-sno-cluster"
}

// ClusterNamespace returns the namespace of the cluster.
func (i *InstallConfig) ClusterNamespace() string {
	if i.Config != nil && i.Config.ObjectMeta.Namespace != "" {
		return i.Config.ObjectMeta.Namespace
	}
	return ""
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
}
