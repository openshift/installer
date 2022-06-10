package manifests

import (
	"fmt"
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/yaml"

	"github.com/openshift/installer/pkg/asset"
	"github.com/pkg/errors"
)

var (
	agentPullSecretName     = "pull-secret"
	agentPullSecretFilename = filepath.Join(clusterManifestDir, fmt.Sprintf("%s.yaml", agentPullSecretName))
)

// AgentPullSecret generates the pull-secret file used by the agent installer.
type AgentPullSecret struct {
	asset.DefaultFileWriter

	Config *corev1.Secret
}

var _ asset.WritableAsset = (*AgentPullSecret)(nil)

// Name returns a human friendly name for the asset.
func (*AgentPullSecret) Name() string {
	return "Agent PullSecret"
}

// Dependencies returns all of the dependencies directly needed to generate
// the asset.
func (*AgentPullSecret) Dependencies() []asset.Asset {
	return []asset.Asset{
		// &installconfig.InstallConfig{},
	}
}

// Generate generates the AgentPullSecret manifest.
func (a *AgentPullSecret) Generate(dependencies asset.Parents) error {

	// installConfigAsset := &installconfig.InstallConfig{}
	// dependencies.Get(installConfigAsset)

	// secret := &corev1.Secret{
	// 	TypeMeta: metav1.TypeMeta{
	// 		APIVersion: "v1",
	// 		Kind:       "Secret",
	// 	},
	// 	ObjectMeta: metav1.ObjectMeta{
	// 		Name:      a.ResourceName(),
	// 		Namespace: installConfigAsset.Config.Namespace,
	// 	},
	// 	StringData: map[string]string{
	// 		".dockerconfigjson": base64.StdEncoding.EncodeToString([]byte(installConfigAsset.Config.PullSecret)),
	// 	},
	// }
	// a.Config = secret

	// secretData, err := yaml.Marshal(secret)
	// if err != nil {
	// 	return errors.Wrap(err, "failed to marshal agent secret")
	// }

	// a.File = &asset.File{
	// 	Filename: agentPullSecretFilename,
	// 	Data:     secretData,
	// }

	return nil
}

// Load returns the asset from disk.
func (a *AgentPullSecret) Load(f asset.FileFetcher) (bool, error) {
	file, err := f.FetchByName(agentPullSecretFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.Wrap(err, fmt.Sprintf("failed to load %s file", agentPullSecretFilename))
	}

	config := &corev1.Secret{}
	if err := yaml.UnmarshalStrict(file.Data, config); err != nil {
		return false, errors.Wrapf(err, "failed to unmarshal %s", agentPullSecretFilename)
	}

	a.File, a.Config = file, config
	if err = a.finish(); err != nil {
		return false, err
	}

	return true, nil
}

// ResourceName return the name of the pull secret resource.
func (a *AgentPullSecret) ResourceName() string {
	return agentPullSecretName
}

func (a *AgentPullSecret) finish() error {
	if err := a.validatePullSecret().ToAggregate(); err != nil {
		return errors.Wrapf(err, "invalid PullSecret configuration")
	}

	return nil
}

func (a *AgentPullSecret) validatePullSecret() field.ErrorList {
	allErrs := field.ErrorList{}

	if err := a.validateSecretIsNotEmpty(); err != nil {
		allErrs = append(allErrs, err...)
	}

	return allErrs
}

func (a *AgentPullSecret) validateSecretIsNotEmpty() field.ErrorList {

	var allErrs field.ErrorList

	fieldPath := field.NewPath("StringData")

	if len(a.Config.StringData) == 0 {
		allErrs = append(allErrs, field.Required(fieldPath, "the pull secret is empty"))
		return allErrs
	}

	pullSecret, ok := a.Config.StringData[".dockerconfigjson"]
	if !ok {
		allErrs = append(allErrs, field.Required(fieldPath, "the pull secret key '.dockerconfigjson' is not defined"))
		return allErrs
	}

	if pullSecret == "" {
		allErrs = append(allErrs, field.Required(fieldPath, "the pull secret does not contain any data"))
		return allErrs
	}

	return allErrs
}
