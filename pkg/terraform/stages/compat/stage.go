package compat

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/types"
)

// PlatformStages are the stages to run to provision the infrastructure used the legacy compat procedures.
func PlatformStages(platform string) []terraform.Stage {
	return []terraform.Stage{stage{platform: platform}}
}

type stage struct {
	platform string
}

func (s stage) Name() string {
	return ""
}

func (s stage) StateFilename() string {
	return "terraform.tfstate"
}

func (s stage) OutputsFilename() string {
	return "outputs.tfvars.json"
}

func (s stage) DestroyWithBootstrap() bool {
	return true
}

func (s stage) Destroy(directory string, extraArgs []string) error {
	extraArgs = append(extraArgs, "-target=module.bootstrap")

	return errors.Wrap(terraform.Destroy(directory, s.platform, s, extraArgs...), "terraform destroy")
}

func (s stage) ExtractHostAddresses(directory string, config *types.InstallConfig) (string, int, []string, error) {
	tfStateFilePath := filepath.Join(directory, s.StateFilename())
	_, err := os.Stat(tfStateFilePath)
	if os.IsNotExist(err) {
		return "", 0, nil, nil
	}
	if err != nil {
		return "", 0, nil, err
	}

	tfstate, err := terraform.ReadState(tfStateFilePath)
	if err != nil {
		return "", 0, nil, errors.Wrapf(err, "failed to read state from %q", tfStateFilePath)
	}
	bootstrap, port, masters, err := extractHostAddresses(config, tfstate)
	return bootstrap, port, masters, errors.Wrapf(err, "failed to get bootstrap and control plane host addresses from %q", tfStateFilePath)
}

func extractHostAddresses(_ *types.InstallConfig, _ *terraform.State) (bootstrap string, port int, masters []string, err error) {
	port = 22
	return bootstrap, port, masters, nil
}
