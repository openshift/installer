package agent

import (
	"fmt"
	"os"
	"strings"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/vsphere"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/yaml"
)

const (
	installConfigFilename = "install-config.yaml"
)

// supportedPlatforms lists the supported platforms for agent installer
var supportedPlatforms = []string{baremetal.Name, vsphere.Name, none.Name}

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

	var found bool

	// First load the provided install config to early validate
	// as per agent installer specific requirements
	// Detailed generic validations of install config are
	// done by pkg/asset/installconfig/installconfig.go
	installConfig, err := a.loadEarly(f)
	if err != nil {
		return found, err
	}

	if err := a.validateInstallConfig(installConfig).ToAggregate(); err != nil {
		return found, errors.Wrapf(err, "invalid install-config configuration")
	}

	found, err = a.InstallConfig.Load(f)
	if found && err == nil {
		a.Supplied = true
	}
	return found, err
}

// loadEarly loads the install config from the disk
// to be able to validate early for agent installer
func (a *OptionalInstallConfig) loadEarly(f asset.FileFetcher) (*types.InstallConfig, error) {

	file, err := f.FetchByName(installConfigFilename)
	config := &types.InstallConfig{}
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil
		}
		return config, errors.Wrap(err, asset.InstallConfigError)
	}

	if err := yaml.UnmarshalStrict(file.Data, config, yaml.DisallowUnknownFields); err != nil {
		if strings.Contains(err.Error(), "unknown field") {
			err = errors.Wrapf(err, "failed to parse first occurence of unknown field")
		}
		err = errors.Wrapf(err, "failed to unmarshal %s", installConfigFilename)
		return config, errors.Wrap(err, asset.InstallConfigError)
	}
	return config, nil
}

func (a *OptionalInstallConfig) validateInstallConfig(installConfig *types.InstallConfig) field.ErrorList {
	var allErrs field.ErrorList

	if err := a.validateSupportedPlatforms(installConfig); err != nil {
		allErrs = append(allErrs, err...)
	}

	if err := a.validateVIPsAreSet(installConfig); err != nil {
		allErrs = append(allErrs, err...)
	}

	if err := a.validateSNOConfiguration(installConfig); err != nil {
		allErrs = append(allErrs, err...)
	}

	return allErrs
}

func (a *OptionalInstallConfig) validateSupportedPlatforms(installConfig *types.InstallConfig) field.ErrorList {
	var allErrs field.ErrorList

	fieldPath := field.NewPath("Platform")

	if installConfig.Platform.Name() != "" && !a.contains(installConfig.Platform.Name(), supportedPlatforms) {
		allErrs = append(allErrs, field.NotSupported(fieldPath, installConfig.Platform.Name(), supportedPlatforms))
	}
	return allErrs
}

func (a *OptionalInstallConfig) validateVIPsAreSet(installConfig *types.InstallConfig) field.ErrorList {
	var allErrs field.ErrorList
	var fieldPath *field.Path

	if installConfig.Platform.Name() == baremetal.Name {
		if len(installConfig.Platform.BareMetal.APIVIPs) == 0 {
			fieldPath = field.NewPath("Platform", "Baremetal", "ApiVips")
			allErrs = append(allErrs, field.Required(fieldPath, fmt.Sprintf("apiVips must be set for %s platform", baremetal.Name)))
		}
		if len(installConfig.Platform.BareMetal.IngressVIPs) == 0 {
			fieldPath = field.NewPath("Platform", "Baremetal", "IngressVips")
			allErrs = append(allErrs, field.Required(fieldPath, fmt.Sprintf("ingressVips must be set for %s platform", baremetal.Name)))
		}
	}

	if installConfig.Platform.Name() == vsphere.Name {
		if len(installConfig.Platform.VSphere.APIVIPs) == 0 {
			fieldPath = field.NewPath("Platform", "VSphere", "ApiVips")
			allErrs = append(allErrs, field.Required(fieldPath, fmt.Sprintf("apiVips must be set for %s platform", vsphere.Name)))
		}
		if len(installConfig.Platform.VSphere.IngressVIPs) == 0 {
			fieldPath = field.NewPath("Platform", "VSphere", "IngressVips")
			allErrs = append(allErrs, field.Required(fieldPath, fmt.Sprintf("ingressVips must be set for %s platform", vsphere.Name)))
		}
	}
	return allErrs
}

func (a *OptionalInstallConfig) validateSNOConfiguration(installConfig *types.InstallConfig) field.ErrorList {
	var allErrs field.ErrorList
	var fieldPath *field.Path

	//  platform None always imply SNO cluster
	if installConfig.Platform.Name() == none.Name {
		if *installConfig.ControlPlane.Replicas != 1 {
			fieldPath = field.NewPath("ControlPlane", "Replicas")
			allErrs = append(allErrs, field.Required(fieldPath, fmt.Sprintf("control plane replicas must be 1 for %s platform. Found %v", none.Name, *installConfig.ControlPlane.Replicas)))
		} else if len(installConfig.Compute) == 0 {
			fieldPath = field.NewPath("Compute", "Replicas")
			allErrs = append(allErrs, field.Required(fieldPath, "Installing a Single Node Openshift requires explicitly setting compute replicas to zero"))
		}

		var workers int
		for _, worker := range installConfig.Compute {
			workers = workers + int(*worker.Replicas)
		}
		if workers != 0 {
			fieldPath = field.NewPath("Compute", "Replicas")
			allErrs = append(allErrs, field.Required(fieldPath, fmt.Sprintf("total number of worker replicas must be 0 for %s platform. Found %v", none.Name, workers)))
		}
	}
	return allErrs
}

func (a *OptionalInstallConfig) contains(platform string, supportedPlatforms []string) bool {
	for _, p := range supportedPlatforms {
		if p == platform {
			return true
		}
	}
	return false
}
