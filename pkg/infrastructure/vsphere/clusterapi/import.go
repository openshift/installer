package clusterapi

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vmware/govmomi/nfc"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/ovf"
	"github.com/vmware/govmomi/ovf/importer"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/session"

	"github.com/openshift/installer/pkg/types/vsphere"
	"github.com/openshift/installer/pkg/utils"
)

const (
	vmware_serial_concentrator_config = "vmware-serial-port-concentrator.json"
)

func debugCorruptOva(cachedImage string, err error) error {
	// Open the corrupt OVA file
	f, ferr := os.Open(cachedImage)
	if ferr != nil {
		err = fmt.Errorf("%s, %w", err.Error(), ferr)
	}
	defer f.Close()

	// Get a sha256 on the corrupt OVA file
	// and the size of the file
	h := sha256.New()
	written, cerr := io.Copy(h, f)
	if cerr != nil {
		err = fmt.Errorf("%s, %w", err.Error(), cerr)
	}

	return fmt.Errorf("ova %s has a sha256 of %x and a size of %d bytes, failed to read the ovf descriptor %w", cachedImage, h.Sum(nil), written, err)
}

func checkOvaSecureBoot(ovfEnvelope *ovf.Envelope) bool {
	if ovfEnvelope.VirtualSystem != nil {
		for _, vh := range ovfEnvelope.VirtualSystem.VirtualHardware {
			for _, c := range vh.Config {
				if c.Key == "bootOptions.efiSecureBootEnabled" {
					if c.Value == "true" {
						return true
					}
				}
			}
		}
	}
	return false
}

func shouldAttachSerialPort() bool {
	if _, err := os.Stat(vmware_serial_concentrator_config); errors.Is(err, os.ErrNotExist) {
		return false

	} else {
		return true
	}
}

func attachVirtualSerialPort(ctx context.Context, vm *object.VirtualMachine) error {
	var concentratorConfig types.VirtualDeviceURIBackingInfo
	concentratorConfigFile, err := os.Open(vmware_serial_concentrator_config)
	if err != nil {
		return fmt.Errorf("unable to open concentrator configuration file: %v", err)
	}
	defer concentratorConfigFile.Close()

	concentratorConfigBytes, err := io.ReadAll(concentratorConfigFile)
	if err != nil {
		return fmt.Errorf("unable to read vmware-serial-port-concentrator.json: %v", err)
	}

	err = json.Unmarshal(concentratorConfigBytes, &concentratorConfig)
	if err != nil {
		return fmt.Errorf("unable to unmarshal vmware-serial-port-concentrator.json: %v", err)
	}

	devices, err := vm.Device(ctx)
	if err != nil {
		return fmt.Errorf("unable to get device list: %v", err)
	}

	debugSerialPort, err := devices.CreateSerialPort()
	if err != nil {
		return fmt.Errorf("unable to create serial port: %v", err)
	}

	isClient := concentratorConfig.Direction == "client"
	debugSerialPort = devices.ConnectSerialPort(debugSerialPort, concentratorConfig.ServiceURI, isClient, concentratorConfig.ProxyURI)
	err = vm.AddDevice(ctx, debugSerialPort)
	if err != nil {
		return fmt.Errorf("unable to add serial port device: %v", err)
	}

	return nil
}

func importRhcosOva(ctx context.Context, session *session.Session, folder *object.Folder, cachedImage, clusterID, tagID, diskProvisioningType string, failureDomain vsphere.FailureDomain) error {
	// Name originally was cluster id + fd.region + fd.zone.  This could cause length of ova to be longer than max allowed.
	// So for now, we are going to make cluster id  + fd.name
	name := utils.GenerateVSphereTemplateName(clusterID, failureDomain.Name)
	logrus.Infof("Importing OVA %v into failure domain %v.", name, failureDomain.Name)

	// OVA name must not exceed 80 characters
	if len(name) > 80 {
		logrus.Warningf("Unable to generate ova template name due to exceeding 80 characters. Cluster=\"%v\" Failure Domain=\"%v\" results in \"%v\"", clusterID, failureDomain.Name, name)
		return fmt.Errorf("ova name \"%v\" exceeed 80 characters (%d)", name, len(name))
	}

	archive := &importer.TapeArchive{Path: cachedImage}

	ovfDescriptor, err := importer.ReadOvf("*.ovf", archive)
	if err != nil {
		return debugCorruptOva(cachedImage, err)
	}

	ovfEnvelope, err := importer.ReadEnvelope(ovfDescriptor)
	if err != nil {
		return fmt.Errorf("failed to parse ovf: %w", err)
	}

	// The fcos ova enables secure boot by default, this causes
	// scos to fail once
	secureBoot := checkOvaSecureBoot(ovfEnvelope)

	// The RHCOS OVA only has one network defined by default
	// The OVF envelope defines this.  We need a 1:1 mapping
	// between networks with the OVF and the host
	if len(ovfEnvelope.Network.Networks) != 1 {
		return fmt.Errorf("expected the OVA to only have a single network adapter")
	}

	cluster, err := session.Finder.ClusterComputeResource(ctx, failureDomain.Topology.ComputeCluster)
	if err != nil {
		return fmt.Errorf("failed to find compute cluster: %w", err)
	}

	clusterHostSystems, err := cluster.Hosts(ctx)
	if err != nil {
		return fmt.Errorf("failed to get cluster hosts: %w", err)
	}

	if len(clusterHostSystems) == 0 {
		return fmt.Errorf("the vCenter cluster %s has no ESXi nodes", failureDomain.Topology.ComputeCluster)
	}

	resourcePool, err := session.Finder.ResourcePool(ctx, failureDomain.Topology.ResourcePool)
	if err != nil {
		return fmt.Errorf("failed to find resource pool: %w", err)
	}

	networkPath := path.Join(cluster.InventoryPath, failureDomain.Topology.Networks[0])

	networkRef, err := session.Finder.Network(ctx, networkPath)
	if err != nil {
		return fmt.Errorf("failed to find network: %w", err)
	}
	datastore, err := session.Finder.Datastore(ctx, failureDomain.Topology.Datastore)
	if err != nil {
		return fmt.Errorf("failed to find datastore: %w", err)
	}

	// Create mapping between OVF and the network object
	// found by Name
	networkMappings := []types.OvfNetworkMapping{{
		Name:    ovfEnvelope.Network.Networks[0].Name,
		Network: networkRef.Reference(),
	}}

	// This is a very minimal spec for importing an OVF.
	cisp := types.OvfCreateImportSpecParams{
		EntityName:     name,
		NetworkMapping: networkMappings,
	}

	switch diskProvisioningType {
	case "":
		// Disk provisioning type will be set according to the default storage policy of vsphere.
	case "thin":
		cisp.DiskProvisioning = string(types.OvfCreateImportSpecParamsDiskProvisioningTypeThin)
	case "thick":
		cisp.DiskProvisioning = string(types.OvfCreateImportSpecParamsDiskProvisioningTypeThick)
	case "eagerZeroedThick":
		cisp.DiskProvisioning = string(types.OvfCreateImportSpecParamsDiskProvisioningTypeEagerZeroedThick)
	default:
		return errors.Errorf("disk provisioning type %q is not supported", diskProvisioningType)
	}

	m := ovf.NewManager(session.Client.Client)
	spec, err := m.CreateImportSpec(ctx,
		string(ovfDescriptor),
		resourcePool.Reference(),
		datastore.Reference(),
		&cisp)

	if err != nil {
		return fmt.Errorf("failed to create import spec: %w", err)
	}
	if spec.Error != nil {
		return errors.New(spec.Error[0].LocalizedMessage)
	}

	hostSystem, err := findAvailableHostSystems(ctx, clusterHostSystems, networkRef, datastore)
	if err != nil {
		return fmt.Errorf("failed to find available host system: %w", err)
	}

	lease, err := resourcePool.ImportVApp(ctx, spec.ImportSpec, folder, hostSystem)
	if err != nil {
		return fmt.Errorf("failed to import vapp: %w", err)
	}

	info, err := lease.Wait(ctx, spec.FileItem)
	if err != nil {
		return fmt.Errorf("failed to lease wait: %w", err)
	}

	u := lease.StartUpdater(ctx, info)
	defer u.Done()

	for _, i := range info.Items {
		// upload the vmdk to which ever host that was first
		// available with the required network and datastore.
		err = upload(ctx, archive, lease, i)
		if err != nil {
			return fmt.Errorf("failed to upload: %w", err)
		}
	}

	err = lease.Complete(ctx)
	if err != nil {
		return fmt.Errorf("failed to lease complete: %w", err)
	}

	vm := object.NewVirtualMachine(session.Client.Client, info.Entity)
	if vm == nil {
		return fmt.Errorf("error VirtualMachine not found, managed object id: %s", info.Entity.Value)
	}
	if secureBoot {
		bootOptions, err := vm.BootOptions(ctx)
		if err != nil {
			return fmt.Errorf("failed to get boot options: %w", err)
		}
		bootOptions.EfiSecureBootEnabled = ptr.To(false)

		err = vm.SetBootOptions(ctx, bootOptions)
		if err != nil {
			return fmt.Errorf("failed to set boot options: %w", err)
		}
	}

	if shouldAttachSerialPort() {
		err = attachVirtualSerialPort(ctx, vm)
		if err != nil {
			return fmt.Errorf("failed to attach virtual serial port: %w", err)
		}
	}

	err = vm.MarkAsTemplate(ctx)
	if err != nil {
		return fmt.Errorf("failed to mark vm as template: %w", err)
	}

	err = attachTag(ctx, session, vm.Reference().Value, tagID)
	if err != nil {
		return fmt.Errorf("failed to attach tag: %w", err)
	}

	return nil
}

func findAvailableHostSystems(ctx context.Context, clusterHostSystems []*object.HostSystem, networkObjectRef object.NetworkReference, datastore *object.Datastore) (*object.HostSystem, error) {
	var hostSystemManagedObject mo.HostSystem
	for _, hostObj := range clusterHostSystems {
		err := hostObj.Properties(ctx, hostObj.Reference(), []string{"config.product", "network", "datastore", "runtime"}, &hostSystemManagedObject)
		if err != nil {
			return nil, err
		}

		// if distributed port group the cast will fail
		networkFound := isNetworkAvailable(networkObjectRef, hostSystemManagedObject.Network)
		datastoreFound := isDatastoreAvailable(datastore, hostSystemManagedObject.Datastore)
		hasUsablePowerState := hostSystemManagedObject.Runtime.PowerState != types.HostSystemPowerStatePoweredOff && hostSystemManagedObject.Runtime.PowerState != types.HostSystemPowerStateStandBy && !hostSystemManagedObject.Runtime.InMaintenanceMode

		// if the network or datastore is not found or the ESXi host is in maintenance mode, powered off or in StandBy (DPM) continue the loop
		if !networkFound || !datastoreFound || !hasUsablePowerState {
			continue
		}

		logrus.Debugf("using ESXi %s to import the OVA image", hostObj.Name())
		return hostObj, nil
	}
	return nil, errors.New("all hosts unavailable")
}

func isDatastoreAvailable(datastore *object.Datastore, hostDatastoreManagedObjectRefs []types.ManagedObjectReference) bool {
	for _, dsMoRef := range hostDatastoreManagedObjectRefs {
		if dsMoRef.Value == datastore.Reference().Value {
			return true
		}
	}
	return false
}

func isNetworkAvailable(networkObjectRef object.NetworkReference, hostNetworkManagedObjectRefs []types.ManagedObjectReference) bool {
	// If the object.NetworkReference is a standard portgroup make
	// sure that it exists on esxi host that the OVA will be imported to.
	if _, ok := networkObjectRef.(*object.Network); ok {
		for _, n := range hostNetworkManagedObjectRefs {
			if n.Value == networkObjectRef.Reference().Value {
				return true
			}
		}
	} else {
		// networkObjectReference is not a standard port group
		// and the other types are distributed so return true
		return true
	}
	return false
}

// Used govc/importx/ovf.go as an example to implement
// resourceVspherePrivateImportOvaCreate and upload functions
// See: https://github.com/vmware/govmomi/blob/cc10a0758d5b4d4873388bcea417251d1ad03e42/govc/importx/ovf.go#L196-L324
func upload(ctx context.Context, archive *importer.TapeArchive, lease *nfc.Lease, item nfc.FileItem) error {
	file := item.Path

	f, size, err := archive.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	opts := soap.Upload{
		ContentLength: size,
	}

	return lease.Upload(ctx, item, f, opts)
}
