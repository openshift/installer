package stages

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/types"
)

// StageOption is an option for configuring a split stage.
type StageOption func(*SplitStage)

// NewStage creates a new split stage.
// The default behavior is the following. The behavior can be changed by providing StageOptions.
// - The resources of the stage will not be deleted as part of destroying the bootstrap.
// - The IP addresses for the bootstrap and control plane VMs will be output from the stage as bootstrap_ip and
//   control_plane_ips, respectively. Only one stage for the platform should output a particular variable. This will
//   likely be the same stage that creates the VM.
func NewStage(platform, name string, opts ...StageOption) SplitStage {
	s := SplitStage{
		platform: platform,
		name:     name,
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

// SplitStage is a split stage.
type SplitStage struct {
	platform             string
	name                 string
	destroyWithBootstrap bool
	destroy              DestroyFunc
	extractHostAddresses ExtractFunc
}

// DestroyFunc is a function for destroying the stage.
type DestroyFunc func(s SplitStage, directory string, extraArgs []string) error

// ExtractFunc is a function for extracting host addresses.
type ExtractFunc func(s SplitStage, directory string, ic *types.InstallConfig) (string, int, []string, error)

// Name implements pkg/terraform/Stage.Name
func (s SplitStage) Name() string {
	return s.name
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
func (s SplitStage) Destroy(directory string, extraArgs []string) error {
	return s.destroy(s, directory, extraArgs)
}

// ExtractHostAddresses implements pkg/terraform/Stage.ExtractHostAddresses
func (s SplitStage) ExtractHostAddresses(directory string, ic *types.InstallConfig) (string, int, []string, error) {
	if s.extractHostAddresses != nil {
		return s.extractHostAddresses(s, directory, ic)
	}
	return normalExtractHostAddresses(s, directory, ic)
}

func normalExtractHostAddresses(s SplitStage, directory string, _ *types.InstallConfig) (string, int, []string, error) {
	outputsFilePath := filepath.Join(directory, s.OutputsFilename())
	if _, err := os.Stat(outputsFilePath); err != nil {
		return "", 0, nil, errors.Wrapf(err, "could not find outputs file %q", outputsFilePath)
	}

	outputsFile, err := ioutil.ReadFile(outputsFilePath)
	if err != nil {
		return "", 0, nil, errors.Wrapf(err, "failed to read outputs file %q", outputsFilePath)
	}

	outputs := map[string]interface{}{}
	if err := json.Unmarshal(outputsFile, &outputs); err != nil {
		return "", 0, nil, errors.Wrapf(err, "could not unmarshal outputs file %q", outputsFilePath)
	}

	var bootstrap string
	if bootstrapRaw, ok := outputs["bootstrap_ip"]; ok {
		bootstrap, ok = bootstrapRaw.(string)
		if !ok {
			return "", 0, nil, errors.Errorf("could not read bootstrap IP from outputs file %q", outputsFilePath)
		}
	}

	var masters []string
	if mastersRaw, ok := outputs["control_plane_ips"]; ok {
		mastersSlice, ok := mastersRaw.([]interface{})
		if !ok {
			return "", 0, nil, errors.Errorf("could not read control plane IPs from outputs file %q", outputsFilePath)
		}
		masters = make([]string, len(mastersSlice))
		for i, ipRaw := range mastersSlice {
			ip, ok := ipRaw.(string)
			if !ok {
				return "", 0, nil, errors.Errorf("could not read control plane IPs from outputs file %q", outputsFilePath)
			}
			masters[i] = ip
		}
	}

	return bootstrap, 0, masters, nil
}

func normalDestroy(s SplitStage, directory string, extraArgs []string) error {
	return errors.Wrap(terraform.Destroy(directory, s.platform, s, extraArgs...), "terraform destroy")
}
