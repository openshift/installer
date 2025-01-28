package configimage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/yaml"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/conversion"
	"github.com/openshift/installer/pkg/types/defaults"
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

// loadFromFile method to load the install-config.yaml file.
// Default one adds many default values that are not needed for image-based install config and can break our logic such as machine network for example.
func (i *InstallConfig) loadFromFile(f asset.FileFetcher) (found bool, err error) {
	file, err := f.FetchByName(InstallConfigFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("%s, err: %w", asset.InstallConfigError, err)
	}

	config := &types.InstallConfig{}
	if err := yaml.UnmarshalStrict(file.Data, config, yaml.DisallowUnknownFields); err != nil {
		err = fmt.Errorf("failed to unmarshal %s, err: %w", InstallConfigFilename, err)
		if !strings.Contains(err.Error(), "unknown field") {
			return false, fmt.Errorf("%s, err: %w", asset.InstallConfigError, err)
		}
		err = fmt.Errorf("failed to parse first occurrence of unknown field, err: %w", err)
		logrus.Warn(err.Error())
		logrus.Info("Attempting to unmarshal while ignoring unknown keys because strict unmarshaling failed")
		if err = yaml.Unmarshal(file.Data, config); err != nil {
			err = fmt.Errorf("failed to unmarshal %s, err: %w", InstallConfigFilename, err)
			return false, fmt.Errorf("%s, err: %w", asset.InstallConfigError, err)
		}
	}
	i.Config = config

	// Upconvert any deprecated fields
	if err := conversion.ConvertInstallConfig(i.Config); err != nil {
		return false, fmt.Errorf("%s, failed to upconvert install config: %w", asset.InstallConfigError, err)
	}

	return true, nil
}

// Load returns the installconfig from disk.
func (i *InstallConfig) Load(f asset.FileFetcher) (bool, error) {
	found, err := i.loadFromFile(f)
	if found && err == nil {
		installConfig := &types.InstallConfig{}
		if err := deepCopy(i.Config, installConfig); err != nil {
			return false, fmt.Errorf("invalid install-config configuration: %w", err)
		}
		if err := i.validateInstallConfig(installConfig).ToAggregate(); err != nil {
			return false, fmt.Errorf("invalid install-config configuration: %w", err)
		}
		if err := i.RecordFile(); err != nil {
			return false, err
		}
	}
	return found, err
}

// in order to avoid the validation errors, we need to set the defaults and validate the configuration
// though those defaults are not used in the image-based install config.
func (i *InstallConfig) validateInstallConfig(installConfig *types.InstallConfig) field.ErrorList {
	var allErrs field.ErrorList

	defaults.SetInstallConfigDefaults(installConfig)
	if err := validation.ValidateInstallConfig(installConfig, true); err != nil {
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
	if machineNetworksCount > 1 {
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
		logrus.Warnf("%s: %s is ignored", fieldPath, installConfig.AdditionalTrustBundlePolicy)
	}

	for i, compute := range installConfig.Compute {
		if compute.Hyperthreading != "Enabled" {
			fieldPath := field.NewPath(fmt.Sprintf("Compute[%d]", i), "Hyperthreading")
			logrus.Warnf("%s: %s is ignored", fieldPath, compute.Hyperthreading)
		}

		if compute.Platform != (types.MachinePoolPlatform{}) {
			fieldPath := field.NewPath(fmt.Sprintf("Compute[%d]", i), "Platform")
			logrus.Warnf("%s is ignored", fieldPath)
		}
	}

	if installConfig.ControlPlane.Hyperthreading != "Enabled" {
		fieldPath := field.NewPath("ControlPlane", "Hyperthreading")
		logrus.Warnf("%s: %s is ignored", fieldPath, installConfig.ControlPlane.Hyperthreading)
	}

	if installConfig.ControlPlane.Platform != (types.MachinePoolPlatform{}) {
		fieldPath := field.NewPath("ControlPlane", "Platform")
		logrus.Warnf("%s is ignored", fieldPath)
	}
}

func deepCopy(src, dst interface{}) error {
	bytes, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, dst)
}
