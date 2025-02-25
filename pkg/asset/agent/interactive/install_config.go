package interactive

import (
	"context"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/yaml"

	"github.com/openshift/installer/pkg/asset"
)

var (
	interactiveInstallConfigFilename = "interactive-install-config.yaml"

	defaultPullSecret = `{"auths":{"":{"auth":"dXNlcjpwYXNz"}}}` //nolint:gosec // no sensitive info
)

// InstallConfig defines the default setting for the cluster
// when using the interfactive disconnected workflow.
type InstallConfig struct {
	config Config
}

// Config is used to read the content of the configuration file.
type Config struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	RendezvousIP string   `json:"rendezvousIP,omitempty"`
	SSHKey       string   `json:"sshKey,omitempty"`
	PullSecret   string   `json:"pullSecret,omitempty"`
	Operators    []string `json:"operators,omitempty"`
}

var _ asset.WritableAsset = (*InstallConfig)(nil)

// Name returns the human-friendly name of the asset.
func (ic *InstallConfig) Name() string {
	return "Agent Installer Interactive InstallConfig"
}

// Dependencies returns all of the dependencies directly needed to generate
// the asset.
func (*InstallConfig) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates the InteractiveInstallConfig.
func (ic *InstallConfig) Generate(ctx context.Context, dependencies asset.Parents) error {
	return nil
}

// Files is not used for this asset.
func (*InstallConfig) Files() []*asset.File {
	return []*asset.File{}
}

// Load returns the config asset from the disk.
func (ic *InstallConfig) Load(f asset.FileFetcher) (bool, error) {
	file, err := f.FetchByName(interactiveInstallConfigFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to load %s file: %w", interactiveInstallConfigFilename, err)
	}

	config := Config{}
	if err = yaml.Unmarshal(file.Data, &config); err != nil {
		return false, fmt.Errorf("failed to unmarshal %s: %w", interactiveInstallConfigFilename, err)
	}
	ic.config = config

	if err = ic.finish(); err != nil {
		return false, err
	}
	return true, nil
}

func (*InstallConfig) finish() error {
	allErrs := field.ErrorList{}

	// TODO: validations

	return allErrs.ToAggregate()
}

// ClusterName returns the cluster name.
func (*InstallConfig) ClusterName() string {
	return "generic-cluster"
}

// ClusterNamespace returns the cluster namespace.
func (*InstallConfig) ClusterNamespace() string {
	return "generic"
}

// PullSecret returns the pull-secret if configured,
// otherwise a dummy (but valid) one.
func (ic *InstallConfig) PullSecret() string {
	pullSecret := defaultPullSecret
	if ic.config.PullSecret != "" {
		pullSecret = ic.config.PullSecret
	}
	return pullSecret
}

// SSHKey return the configure ssh key.
func (ic *InstallConfig) SSHKey() string {
	return ic.config.SSHKey
}

// RendezvousIP returns the configured rendezvous IP.
func (ic *InstallConfig) RendezvousIP() string {
	return ic.config.RendezvousIP
}

// Operators returns the list of OLM operators to be installed.
func (ic *InstallConfig) Operators() []string {
	return ic.config.Operators
}
