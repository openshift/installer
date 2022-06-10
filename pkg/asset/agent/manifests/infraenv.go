package manifests

import (
	"fmt"
	"os"
	"path/filepath"

	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/yaml"

	"github.com/openshift/installer/pkg/asset"
)

var (
	infraEnvFilename = filepath.Join(clusterManifestDir, "infraenv.yaml")
)

// InfraEnv generates the infraenv.yaml file.
type InfraEnv struct {
	asset.DefaultFileWriter

	Config *aiv1beta1.InfraEnv
}

var _ asset.WritableAsset = (*InfraEnv)(nil)

// Name returns a human friendly name for the asset.
func (*InfraEnv) Name() string {
	return "InfraEnv Config"
}

// Dependencies returns all of the dependencies directly needed to generate
// the asset.
func (*InfraEnv) Dependencies() []asset.Asset {
	return []asset.Asset{
		// &installconfig.InstallConfig{},
		// &AgentPullSecret{},
	}
}

// Generate generates the InfraEnv manifest.
func (i *InfraEnv) Generate(dependencies asset.Parents) error {

	// installConfig := &installconfig.InstallConfig{}
	// agentPullSecret := &AgentPullSecret{}
	// dependencies.Get(installConfig, agentPullSecret)

	// infraEnv := &aiv1beta1.InfraEnv{
	// 	ObjectMeta: v1.ObjectMeta{
	// 		Name:      "infraEnv",
	// 		Namespace: installConfig.Config.Namespace,
	// 	},
	// 	Spec: aiv1beta1.InfraEnvSpec{
	// 		ClusterRef: &aiv1beta1.ClusterReference{
	// 			Name:      installConfig.Config.ObjectMeta.Name,
	// 			Namespace: installConfig.Config.ObjectMeta.Namespace,
	// 		},
	// 		SSHAuthorizedKey: installConfig.Config.SSHKey,
	// 		PullSecretRef: &corev1.LocalObjectReference{
	// 			Name: agentPullSecret.ResourceName(),
	// 		},
	// 		// NMStateConfigLabelSelector: v1.LabelSelector{
	// 		// 	MatchLabels: map[string]string{
	// 		// 		// fetch from NMStateConfig
	// 		// 	},
	// 		// },
	// 	},
	// }
	// i.Config = infraEnv

	// infraEnvData, err := yaml.Marshal(infraEnv)
	// if err != nil {
	// 	return errors.Wrap(err, "failed to marshal agent installer infraEnv")
	// }

	// i.File = &asset.File{
	// 	Filename: infraEnvFilename,
	// 	Data:     infraEnvData,
	// }

	return nil
}

// Load returns infraenv asset from the disk.
func (i *InfraEnv) Load(f asset.FileFetcher) (bool, error) {

	file, err := f.FetchByName(infraEnvFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.Wrap(err, fmt.Sprintf("failed to load %s file", infraEnvFilename))
	}

	config := &aiv1beta1.InfraEnv{}
	if err := yaml.UnmarshalStrict(file.Data, config); err != nil {
		return false, errors.Wrapf(err, "failed to unmarshal %s", infraEnvFilename)
	}

	i.File, i.Config = file, config
	if err = i.finish(); err != nil {
		return false, err
	}

	return true, nil
}

func (i *InfraEnv) finish() error {
	if err := i.validateInfraEnv().ToAggregate(); err != nil {
		return errors.Wrapf(err, "invalid InfraEnv configuration")
	}

	return nil
}

func (i *InfraEnv) validateInfraEnv() field.ErrorList {
	allErrs := field.ErrorList{}

	if err := i.validateNMStateLabelSelector(); err != nil {
		allErrs = append(allErrs, err...)
	}

	return allErrs
}

func (i *InfraEnv) validateNMStateLabelSelector() field.ErrorList {

	var allErrs field.ErrorList

	fieldPath := field.NewPath("Spec", "NMStateConfigLabelSelector", "MatchLabels")

	if len(i.Config.Spec.NMStateConfigLabelSelector.MatchLabels) == 0 {
		allErrs = append(allErrs, field.Required(fieldPath, "at least one label must be set"))
	}

	return allErrs
}
