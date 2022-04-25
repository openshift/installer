package vsphere

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25"
	vmwaretypes "github.com/vmware/govmomi/vim25/types"

	"github.com/openshift/installer/pkg/client/vsphere"
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/providers"
	"github.com/openshift/installer/pkg/terraform/stages"
	"github.com/openshift/installer/pkg/types"
)

// PlatformStages are the stages to run to provision the infrastructure in vsphere.
var PlatformStages = []terraform.Stage{
	stages.NewStage(
		"vsphere",
		"pre-bootstrap",
		[]providers.Provider{providers.VSphere, providers.VSpherePrivate},
	),
	stages.NewStage(
		"vsphere",
		"bootstrap",
		[]providers.Provider{providers.VSphere},
		stages.WithNormalBootstrapDestroy(),
		stages.WithCustomExtractHostAddresses(extractOutputHostAddresses),
	),
	stages.NewStage(
		"vsphere",
		"master",
		[]providers.Provider{providers.VSphere},
		stages.WithCustomExtractHostAddresses(extractOutputHostAddresses),
	),
}

func extractOutputHostAddresses(s stages.SplitStage, directory string, config *types.InstallConfig) (bootstrap string, port int, masters []string, err error) {
	port = 22

	outputs, err := stages.GetTerraformOutputs(s, directory)
	if err != nil {
		return "", 0, nil, err
	}

	var bootstrapMoid string
	if bootstrapRaw, ok := outputs["bootstrap_moid"]; ok {
		bootstrapMoid, ok = bootstrapRaw.(string)
		if !ok {
			return "", 0, nil, errors.New("could not read bootstrap MOID from terraform outputs")
		}
	}

	var mastersMoids []string
	if mastersRaw, ok := outputs["control_plane_moids"]; ok {
		mastersSlice, ok := mastersRaw.([]interface{})
		if !ok {
			return "", 0, nil, errors.New("could not read control plane MOIDs from terraform outputs")
		}
		mastersMoids = make([]string, len(mastersSlice))
		for i, moidRaw := range mastersSlice {
			moid, ok := moidRaw.(string)
			if !ok {
				return "", 0, nil, errors.New("could not read control plane MOIDs from terraform outputs")
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
	client, _, cleanup, err := vsphere.CreateVSphereClients(context.TODO(), config.VSphere.VCenter, config.VSphere.Username, config.VSphere.Password)
	if err != nil {
		return "", err
	}
	defer cleanup()

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
