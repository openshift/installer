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
type StageOption func(SplitStage)

// NewStage creates a new split stage.
func NewStage(platform, name string, opts ...StageOption) SplitStage {
	s := SplitStage{
		platform: platform,
		name:     name,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// WithNormalDestroy returns an option for specifying that a split stage should use the normal destroy process.
func WithNormalDestroy() StageOption {
	return WithCustomDestroy(normalDestroy)
}

// WithCustomDestroy returns an option for specifying that a split stage should use a custom destroy process.
func WithCustomDestroy(destroy DestroyFunc) StageOption {
	return func(s SplitStage) {
		s.destroyWithBootstrap = true
		s.destroy = destroy
	}
}

// SplitStage is a split stage.
type SplitStage struct {
	platform             string
	name                 string
	destroyWithBootstrap bool
	destroy              DestroyFunc
}

// DestroyFunc is a function for destroying the stage.
type DestroyFunc func(s SplitStage, directory string, extraArgs []string) error

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
func (s SplitStage) ExtractHostAddresses(directory string, _ *types.InstallConfig) (string, int, []string, error) {
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
