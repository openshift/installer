package vsphere

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/terraform"
	gathervsphere "github.com/openshift/installer/pkg/terraform/gather/vsphere"
	"github.com/openshift/installer/pkg/terraform/stages"
	"github.com/openshift/installer/pkg/types"
)

// PlatformStages are the stages to run to provision the infrastructure in vsphere.
var PlatformStages = []terraform.Stage{
	stages.NewStage("vsphere", "pre-bootstrap"),
	stages.NewStage("vsphere", "bootstrap", stages.WithNormalDestroy(), stages.WithCustomExtractHostAddresses(extractOutputHostAddresses)),
	stages.NewStage("vsphere", "master", stages.WithCustomExtractHostAddresses(extractOutputHostAddresses)),
}

func extractOutputHostAddresses(s stages.SplitStage, directory string, config *types.InstallConfig) (bootstrap string, port int, masters []string, err error) {
	port = 22
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

	var bootstrapMoid string
	if bootstrapRaw, ok := outputs["bootstrap_moid"]; ok {
		bootstrapMoid, ok = bootstrapRaw.(string)
		if !ok {
			return "", 0, nil, errors.Errorf("could not read bootstrap MOID from outputs file %q", outputsFilePath)
		}
	}

	var mastersMoids []string
	if mastersRaw, ok := outputs["control_plane_moids"]; ok {
		mastersSlice, ok := mastersRaw.([]interface{})
		if !ok {
			return "", 0, nil, errors.Errorf("could not read control plane MOIDs from outputs file %q", outputsFilePath)
		}
		mastersMoids = make([]string, len(mastersSlice))
		for i, moidRaw := range mastersSlice {
			moid, ok := moidRaw.(string)
			if !ok {
				return "", 0, nil, errors.Errorf("could not read control plane MOIDs from outputs file %q", outputsFilePath)
			}
			mastersMoids[i] = moid
		}
	}

	bootstrap, err = gathervsphere.HostIP(config, bootstrapMoid)
	if err != nil {
		return "", 0, nil, errors.Errorf("could not extract IP with bootstrap MOID: %s", bootstrapMoid)
	}

	masters = make([]string, len(mastersMoids))
	for i, moid := range mastersMoids {
		masters[i], err = gathervsphere.HostIP(config, moid)
		if err != nil {
			return "", 0, nil, errors.Errorf("could not extract IP with control node MOID: %s", moid)
		}
	}

	return bootstrap, port, masters, nil
}
