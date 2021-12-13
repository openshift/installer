package vsphere

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25"
	vmwaretypes "github.com/vmware/govmomi/vim25/types"

	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/stages"
	"github.com/openshift/installer/pkg/types"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
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

	bootstrap, err = hostIP(config, bootstrapMoid)
	if err != nil {
		return "", 0, nil, errors.Errorf("could not extract IP with bootstrap MOID: %s", bootstrapMoid)
	}

	masters = make([]string, len(mastersMoids))
	for i, moid := range mastersMoids {
		masters[i], err = hostIP(config, moid)
		if err != nil {
			return "", 0, nil, errors.Errorf("could not extract IP with control node MOID: %s", moid)
		}
	}

	return bootstrap, port, masters, nil
}

// hostIP returns the ip address for a host
func hostIP(config *types.InstallConfig, moid string) (string, error) {
	client, _, err := vspheretypes.CreateVSphereClients(context.TODO(), config.VSphere.VCenter, config.VSphere.Username, config.VSphere.Password)
	if err != nil {
		return "", err
	}

	var errs []error
	ip, err := waitForVirtualMachineIP(client, moid)
	if err != nil {
		errs = append(errs, errors.Wrapf(err, "failed to lookup ipv4 address from given moid %s", moid))
	}

	return ip, nil
}

func waitForVirtualMachineIP(client *vim25.Client, moRefValue string) (string, error) {
	moRef := vmwaretypes.ManagedObjectReference{
		Type:  "VirtualMachine",
		Value: moRefValue,
	}

	vm := object.NewVirtualMachine(client, moRef)
	if vm == nil {
		return "", errors.Errorf("VirtualMachine was not found")
	}
	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	ip, err := vm.WaitForIP(ctx, true)
	if err != nil {
		return "", err
	}
	return ip, nil
}
