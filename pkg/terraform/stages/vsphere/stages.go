package vsphere

import (
	"context"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25"
	vmwaretypes "github.com/vmware/govmomi/vim25/types"

	"github.com/openshift/installer/pkg/asset/installconfig/vsphere"
	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/providers"
	"github.com/openshift/installer/pkg/terraform/stages"
	"github.com/openshift/installer/pkg/types"
)

func InitializeProvider(installDir string) ([]infrastructure.Stage, func() error, error) {
	terraformDir, err := terraform.Initialize(installDir)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error initializing terraform")
	}

	// PlatformStages are the stages to run to provision the infrastructure in a
	// multiple region and zone vsphere environment.
	var platformStages = []infrastructure.Stage{
		stages.NewStage(
			"vsphere",
			"pre-bootstrap",
			installDir,
			terraformDir,
			[]providers.Provider{providers.VSphere, providers.VSpherePrivate},
		),
		stages.NewStage(
			"vsphere",
			"bootstrap",
			installDir,
			terraformDir,
			[]providers.Provider{providers.VSphere},
			stages.WithNormalBootstrapDestroy(),
			stages.WithCustomExtractHostAddresses(extractOutputHostAddresses),
		),
		stages.NewStage(
			"vsphere",
			"master",
			installDir,
			terraformDir,
			[]providers.Provider{providers.VSphere},
			stages.WithCustomExtractHostAddresses(extractOutputHostAddresses),
		),
	}
	// It would be nice to not need to repeat this for each platform but at this stage
	// Perfect is the enemy of good
	terraform.UnpackTerraform(terraformDir, platformStages)

	cleanup := func() error {
		return os.RemoveAll(terraformDir)
	}

	return platformStages, cleanup, nil
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
	client, _, cleanup, err := vsphere.CreateVSphereClients(context.TODO(), config.VSphere.VCenters[0].Server, config.VSphere.VCenters[0].Username, config.VSphere.VCenters[0].Password)
	if err != nil {
		return "", err
	}
	defer cleanup()

	ip, err := waitForVirtualMachineIP(client, moid)
	if err != nil {
		return "", errors.Wrapf(err, "failed to lookup ipv4 address from given moid %s", moid)
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
