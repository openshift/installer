package stages

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/providers"
	"github.com/openshift/installer/pkg/types"
)

// StageOption is an option for configuring a split stage.
type StageOption func(*SplitStage)

// NewStage creates a new split stage.
// The default behavior is the following. The behavior can be changed by providing StageOptions.
//   - The resources of the stage will not be deleted as part of destroying the bootstrap.
//   - The IP addresses for the bootstrap and control plane VMs will be output from the stage as bootstrap_ip and
//     control_plane_ips, respectively. Only one stage for the platform should output a particular variable. This will
//     likely be the same stage that creates the VM.
func NewStage(platform, name string, providers []providers.Provider, opts ...StageOption) SplitStage {
	s := SplitStage{
		platform:  platform,
		name:      name,
		providers: providers,
	}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}

// WithNormalBootstrapDestroy returns an option for specifying that a split stage should use the normal bootstrap
// destroy process. The normal process is to fully delete all of the resources created in the stage.
func WithNormalBootstrapDestroy() StageOption {
	return WithCustomBootstrapDestroy(normalDestroy)
}

// WithCustomBootstrapDestroy returns an option for specifying that a split stage should use a custom bootstrap
// destroy process.
func WithCustomBootstrapDestroy(destroy DestroyFunc) StageOption {
	return func(s *SplitStage) {
		s.destroyWithBootstrap = true
		s.destroy = destroy
	}
}

// WithCustomExtractHostAddresses returns an option for specifying that a split stage should use a custom extract host addresses process.
func WithCustomExtractHostAddresses(extractHostAddresses ExtractFunc) StageOption {
	return func(s *SplitStage) {
		s.extractHostAddresses = extractHostAddresses
	}
}

// WithCustomExtractLBConfig returns an option for specifying that a split stage
// should use a custom method to extract load balancer DNS names.
func WithCustomExtractLBConfig(extractLBConfig ExtractLBConfigFunc) StageOption {
	return func(s *SplitStage) {
		s.extractLBConfig = extractLBConfig
	}
}

// SplitStage is a split stage.
type SplitStage struct {
	platform             string
	name                 string
	providers            []providers.Provider
	destroyWithBootstrap bool
	destroy              DestroyFunc
	extractHostAddresses ExtractFunc
	extractLBConfig      ExtractLBConfigFunc
}

// DestroyFunc is a function for destroying the stage.
type DestroyFunc func(s SplitStage, directory string, terraformDir string, varFiles []string) error

// ExtractFunc is a function for extracting host addresses.
type ExtractFunc func(s SplitStage, directory string, ic *types.InstallConfig) (string, int, []string, error)

// ExtractLBConfigFunc is a function for extracting LB DNS Names.
type ExtractLBConfigFunc func(s SplitStage, directory string, terraformDir string, file *asset.File, tfvarsFile *asset.File) (string, error)

// Name implements pkg/terraform/Stage.Name
func (s SplitStage) Name() string {
	return s.name
}

// Providers is the list of providers that are used for the stage.
func (s SplitStage) Providers() []providers.Provider {
	return s.providers
}

// StateFilename implements pkg/terraform/Stage.StateFilename
func (s SplitStage) StateFilename() string {
	return fmt.Sprintf("terraform.%s.tfstate", s.name)
}

// OutputsFilename implements pkg/terraform/Stage.OutputsFilename
func (s SplitStage) OutputsFilename() string {
	return fmt.Sprintf("%s.tfvars.json", s.name)
}

// DestroyWithBootstrap implements pkg/terraform/Stage.DestroyWithBootstrap
func (s SplitStage) DestroyWithBootstrap() bool {
	return s.destroyWithBootstrap
}

// Destroy implements pkg/terraform/Stage.Destroy
func (s SplitStage) Destroy(directory string, terraformDir string, varFiles []string) error {
	return s.destroy(s, directory, terraformDir, varFiles)
}

// Platform implements pkg/terraform/Stage.Platform.
func (s SplitStage) Platform() string {
	return s.platform
}

// ExtractHostAddresses implements pkg/terraform/Stage.ExtractHostAddresses
func (s SplitStage) ExtractHostAddresses(directory string, ic *types.InstallConfig) (string, int, []string, error) {
	if s.extractHostAddresses != nil {
		return s.extractHostAddresses(s, directory, ic)
	}
	return normalExtractHostAddresses(s, directory, ic)
}

// ExtractLBConfig implements pkg/terraform/Stage.ExtractLBConfig.
func (s SplitStage) ExtractLBConfig(directory string, terraformDir string, file *asset.File, tfvarsFile *asset.File) (string, error) {
	if s.extractLBConfig != nil {
		return s.extractLBConfig(s, directory, terraformDir, file, tfvarsFile)
	}
	return normalExtractLBConfig(s, directory, terraformDir, file, tfvarsFile)
}

// GetTerraformOutputs reads the terraform outputs file for the stage and parses it into a map of outputs.
func GetTerraformOutputs(s SplitStage, directory string) (map[string]interface{}, error) {
	outputsFilePath := filepath.Join(directory, s.OutputsFilename())
	if _, err := os.Stat(outputsFilePath); err != nil {
		return nil, errors.Wrapf(err, "could not find outputs file %q", outputsFilePath)
	}

	outputsFile, err := os.ReadFile(outputsFilePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read outputs file %q", outputsFilePath)
	}

	outputs := map[string]interface{}{}
	if err := json.Unmarshal(outputsFile, &outputs); err != nil {
		return nil, errors.Wrapf(err, "could not unmarshal outputs file %q", outputsFilePath)
	}

	return outputs, nil
}

func normalExtractHostAddresses(s SplitStage, directory string, _ *types.InstallConfig) (string, int, []string, error) {
	outputs, err := GetTerraformOutputs(s, directory)
	if err != nil {
		return "", 0, nil, err
	}

	var bootstrap string
	if bootstrapRaw, ok := outputs["bootstrap_ip"]; ok {
		bootstrap, ok = bootstrapRaw.(string)
		if !ok {
			return "", 0, nil, errors.New("could not read bootstrap IP from terraform outputs")
		}
	}

	var masters []string
	if mastersRaw, ok := outputs["control_plane_ips"]; ok {
		mastersSlice, ok := mastersRaw.([]interface{})
		if !ok {
			return "", 0, nil, errors.New("could not read control plane IPs from terraform outputs")
		}
		masters = make([]string, len(mastersSlice))
		for i, ipRaw := range mastersSlice {
			ip, ok := ipRaw.(string)
			if !ok {
				return "", 0, nil, errors.New("could not read control plane IPs from terraform outputs")
			}
			masters[i] = ip
		}
	}

	return bootstrap, 0, masters, nil
}

func normalDestroy(s SplitStage, directory string, terraformDir string, varFiles []string) error {
	opts := make([]tfexec.DestroyOption, len(varFiles))
	for i, varFile := range varFiles {
		opts[i] = tfexec.VarFile(varFile)
	}
	return errors.Wrap(terraform.Destroy(directory, s.platform, s, terraformDir, opts...), "terraform destroy")
}

func normalExtractLBConfig(s SplitStage, directory string, terraformDir string, file *asset.File, tfvarsFile *asset.File) (string, error) {
	return "", nil
}
