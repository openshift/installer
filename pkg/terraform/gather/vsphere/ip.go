// Package vsphere contains utilities that help gather vsphere specific
// information from terraform state.
package vsphere

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/types"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"

	"github.com/openshift/installer/pkg/terraform"
	installertypes "github.com/openshift/installer/pkg/types"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

func waitForVirtualMachineIP(client *vim25.Client, moRefValue string) (string, error) {
	moRef := types.ManagedObjectReference{
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

// BootstrapIP returns the ip address for bootstrap host.
func BootstrapIP(config *installertypes.InstallConfig, tfs *terraform.State) (string, error) {
	client, _, err := vspheretypes.CreateVSphereClients(context.TODO(), config.VSphere.VCenter, config.VSphere.Username, config.VSphere.Password)
	if err != nil {
		return "", err
	}

	br, err := terraform.LookupResource(tfs, "module.bootstrap", "vsphere_virtual_machine", "vm")

	if err != nil {
		return "", errors.Wrap(err, "failed to lookup bootstrap")
	}

	if len(br.Instances) == 0 {
		return "", errors.New("no bootstrap instance found")
	}

	moid, found, err := unstructured.NestedString(br.Instances[0].Attributes, "moid")
	if err != nil {
		return "", errors.Wrap(err, "failed to lookup bootstrap managed object reference")
	}
	if !found {
		return "", errors.Errorf("failed to lookup bootstrap managed object reference")
	}
	ip, err := waitForVirtualMachineIP(client, moid)
	if err != nil {
		return "", errors.Wrap(err, "failed to lookup bootstrap ipv4 address")
	}

	return ip, nil
}

// ControlPlaneIPs returns the ip addresses for control plane hosts.
func ControlPlaneIPs(config *installertypes.InstallConfig, tfs *terraform.State) ([]string, error) {
	client, _, err := vspheretypes.CreateVSphereClients(context.TODO(), config.VSphere.VCenter, config.VSphere.Username, config.VSphere.Password)
	if err != nil {
		return nil, err
	}

	mrs, err := terraform.LookupResource(tfs, "module.master", "vsphere_virtual_machine", "vm")
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup masters")
	}
	var errs []error
	var masters []string
	for idx, inst := range mrs.Instances {
		moid, _, err := unstructured.NestedString(inst.Attributes, "moid")
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "failed to lookup master.%d managed object reference", idx))
		}
		master, err := waitForVirtualMachineIP(client, moid)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "failed to lookup master.%d ipv4 address", idx))
		}

		masters = append(masters, master)
	}
	return masters, utilerrors.NewAggregate(errs)
}
