package vsphere

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/govc/importx"
	"github.com/vmware/govmomi/nfc"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/ovf"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vapi/tags"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"

	machinev1beta1 "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/asset/installconfig"
	installertypes "github.com/openshift/installer/pkg/types"
	typesvsphere "github.com/openshift/installer/pkg/types/vsphere"
)

type VCenterConnection struct {
	Client     *govmomi.Client
	Finder     *find.Finder
	Context    context.Context
	RestClient *rest.Client
	Logout     func()

	Uri      string
	Username string
	Password string
}

func getVCenterClient(uri, username, password string) (*VCenterConnection, error) {
	logrus.Infof("In getVCenterClient")
	ctx := context.Background()

	connection := &VCenterConnection{
		Context: ctx,
	}

	u, err := soap.ParseURL(uri)
	if err != nil {
		return nil, err
	}
	connection.Username = username
	connection.Password = password
	connection.Uri = uri

	u.User = url.UserPassword(username, password)

	c, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		return nil, err
	}

	connection.RestClient = rest.NewClient(c.Client)

	err = connection.RestClient.Login(connection.Context, u.User)
	if err != nil {
		return nil, err
	}
	connection.Client = c

	connection.Finder = find.NewFinder(connection.Client.Client)

	connection.Logout = func() {
		connection.Client.Logout(connection.Context)
		connection.RestClient.Logout(connection.Context)
	}

	return connection, nil

}

// Metadata converts an install configuration to vSphere metadata.
func Metadata(config *installertypes.InstallConfig) *typesvsphere.Metadata {
	terraformPlatform := "vsphere"

	// Since currently we only support a single vCenter
	// just use the first entry in the VCenters slice.

	return &typesvsphere.Metadata{
		VCenter:           config.VSphere.VCenters[0].Server,
		Username:          config.VSphere.VCenters[0].Username,
		Password:          config.VSphere.VCenters[0].Password,
		TerraformPlatform: terraformPlatform,
	}
}

func PreTerraform(cachedImage, boostrapIgn, masterIgn string, controlPlaneMachines []machinev1beta1.Machine,
	clusterID string, installConfig *installconfig.InstallConfig) error {
	logrus.Infof("In PreTerraform")

	vmTemplateMap := make(map[string]*object.VirtualMachine)
	vconn, err := getVCenterClient(
		installConfig.Config.VSphere.VCenters[0].Server,
		installConfig.Config.VSphere.VCenters[0].Username,
		installConfig.Config.VSphere.VCenters[0].Password)

	if err != nil {
		return err
	}
	defer vconn.Logout()

	categoryId, err := createTagCategory(vconn, clusterID)
	if err != nil {
		return err
	}

	tagId, err := createTag(vconn, clusterID, categoryId)

	if err != nil {
		return err
	}

	for _, fd := range installConfig.Config.VSphere.FailureDomains {
		dc, err := vconn.Finder.Datacenter(vconn.Context, fd.Topology.Datacenter)
		if err != nil {
			return err
		}
		dcFolders, err := dc.Folders(vconn.Context)

		folderPath := path.Join(dcFolders.VmFolder.InventoryPath, clusterID)

		// we must set the Folder to the infraId somewhere, we will need to remove that.
		// if we are overwriting folderPath it needs to have a slash (path)
		if strings.Contains(fd.Topology.Folder, "/") {
			folderPath = fd.Topology.Folder
		}

		folder, err := createFolder(folderPath, vconn)
		if err != nil {
			return err
		}
		vmTemplate, err := importRhcosOva(vconn, folder, cachedImage, clusterID, tagId, string(installConfig.Config.VSphere.DiskType), fd)

		// This object.VirtualMachine is not fully defined
		vmName, err := vmTemplate.ObjectName(vconn.Context)

		if err != nil {
			return err
		}

		vmTemplateMap[vmName] = vmTemplate

		if err != nil {
			return err
		}

	}
	encodedMasterIgn := base64.StdEncoding.EncodeToString([]byte(masterIgn))
	encodedBootstrapIgn := base64.StdEncoding.EncodeToString([]byte(boostrapIgn))

	controlPlaneConfigs := make([]*machinev1beta1.VSphereMachineProviderSpec, len(controlPlaneMachines))

	bootstrap := true

	for i := 0; i < len(controlPlaneMachines); i++ {
		vmName := fmt.Sprintf("%s-master-%d", clusterID, i)

		encodedIgnition := encodedMasterIgn
		controlPlaneConfigs[i] = controlPlaneMachines[i].Spec.ProviderSpec.Value.Object.(*machinev1beta1.VSphereMachineProviderSpec)

		if bootstrap {
			if i == 0 {
				encodedIgnition = encodedBootstrapIgn
				vmName = fmt.Sprintf("%s-bootstrap", clusterID)
			}
		}

		task, err := clone(vconn, vmTemplateMap[controlPlaneConfigs[i].Template], controlPlaneConfigs[i], encodedIgnition, vmName)
		if err != nil {
			return err
		}

		taskInfo, err := task.WaitForResult(vconn.Context, nil)
		if err != nil {
			return err
		}

		vmMoRef := taskInfo.Result.(types.ManagedObjectReference)
		vm := object.NewVirtualMachine(vconn.Client.Client, vmMoRef)

		err = attachTag(vconn, vmMoRef.Value, tagId)
		if err != nil {
			return err
		}
		task, err = vm.PowerOn(vconn.Context)
		if err != nil {
			return err
		}

		task.WaitForResult(vconn.Context, nil)

		if bootstrap {
			if i == 0 {
				bootstrap = false
				i = -1
			}
		}
	}

	return nil
}

func createFolder(fullpath string, vconn *VCenterConnection) (*object.Folder, error) {
	logrus.Infof("In createFolder")
	dir := path.Dir(fullpath)
	base := path.Base(fullpath)

	folder, err := vconn.Finder.Folder(context.TODO(), fullpath)

	if folder == nil {
		folder, err = vconn.Finder.Folder(context.TODO(), dir)

		if _, ok := err.(*find.NotFoundError); ok {
			folder, err = createFolder(dir, vconn)
			if err != nil {
				return folder, err
			}
		}

		if folder != nil {
			folder, err = folder.CreateFolder(context.TODO(), base)
			if err != nil {
				return folder, err
			}
		}
	}
	return folder, err
}

func createTagCategory(vconn *VCenterConnection, clusterId string) (string, error) {
	logrus.Infof("In createTagCategory")
	categoryName := fmt.Sprintf("openshift-%s", clusterId)

	category := tags.Category{
		Name:        categoryName,
		Description: "Added by openshift-install do not remove",
		Cardinality: "SINGLE",
		AssociableTypes: []string{
			"VirtualMachine",
			"ResourcePool",
			"Folder",
			"Datastore",
			"StoragePod",
		},
	}

	return tags.NewManager(vconn.RestClient).CreateCategory(vconn.Context, &category)
}

func createTag(vconn *VCenterConnection, clusterId, categoryId string) (string, error) {
	logrus.Infof("In createTag")

	tag := tags.Tag{
		Description: "Added by openshift-install do not remove",
		Name:        clusterId,
		CategoryID:  categoryId,
	}

	return tags.NewManager(vconn.RestClient).CreateTag(vconn.Context, &tag)
}

func importRhcosOva(vconn *VCenterConnection, folder *object.Folder, cachedImage, clusterId, tagId, diskProvisioningType string, failureDomain typesvsphere.FailureDomain) (*object.VirtualMachine, error) {
	logrus.Infof("In importRhcosOva")
	name := fmt.Sprintf("%s-rhcos-%s-%s", clusterId, failureDomain.Region, failureDomain.Zone)

	archive := &importx.ArchiveFlag{Archive: &importx.TapeArchive{Path: cachedImage}}

	ovfDescriptor, err := archive.ReadOvf("*.ovf")
	if err != nil {
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

		return nil, errors.Errorf("ova %s has a sha256 of %x and a size of %d bytes, failed to read the ovf descriptor %s", cachedImage, h.Sum(nil), written, err)
	}

	ovfEnvelope, err := archive.ReadEnvelope(ovfDescriptor)
	if err != nil {
		return nil, errors.Errorf("failed to parse ovf: %s", err)
	}

	// The RHCOS OVA only has one network defined by default
	// The OVF envelope defines this.  We need a 1:1 mapping
	// between networks with the OVF and the host
	if len(ovfEnvelope.Network.Networks) != 1 {
		return nil, errors.Errorf("Expected the OVA to only have a single network adapter")
	}

	cluster, err := vconn.Finder.ClusterComputeResource(vconn.Context, failureDomain.Topology.ComputeCluster)

	if err != nil {
		return nil, err
	}
	clusterHostSystems, err := cluster.Hosts(vconn.Context)

	if err != nil {
		return nil, err
	}
	resourcePool, err := vconn.Finder.ResourcePool(vconn.Context, failureDomain.Topology.ResourcePool)

	networkPath := path.Join(cluster.InventoryPath, failureDomain.Topology.Networks[0])

	networkRef, err := vconn.Finder.Network(vconn.Context, networkPath)
	if err != nil {
		return nil, err
	}
	datastore, err := vconn.Finder.Datastore(vconn.Context, failureDomain.Topology.Datastore)

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

	m := ovf.NewManager(vconn.Client.Client)
	spec, err := m.CreateImportSpec(vconn.Context,
		string(ovfDescriptor),
		resourcePool.Reference(),
		datastore.Reference(),
		cisp)

	if err != nil {
		return nil, errors.Errorf("failed to create import spec: %s", err)
	}
	if spec.Error != nil {
		return nil, errors.New(spec.Error[0].LocalizedMessage)
	}

	hostSystem, err := findAvailableHostSystems(vconn, clusterHostSystems)

	if err != nil {
		return nil, err
	}
	//Creates a new entity in this resource pool.
	//See VMware vCenter API documentation: Managed Object - ResourcePool - ImportVApp
	lease, err := resourcePool.ImportVApp(vconn.Context, spec.ImportSpec, folder, hostSystem)

	if err != nil {
		return nil, errors.Errorf("failed to import vapp: %s", err)
	}

	info, err := lease.Wait(vconn.Context, spec.FileItem)
	if err != nil {
		return nil, errors.Errorf("failed to lease wait: %s", err)
	}

	if err != nil {
		return nil, errors.Errorf("failed to attach tag to virtual machine: %s", err)
	}

	u := lease.StartUpdater(vconn.Context, info)
	defer u.Done()

	for _, i := range info.Items {
		// upload the vmdk to which ever host that was first
		// available with the required network and datastore.
		err = upload(vconn.Context, archive, lease, i)
		if err != nil {
			return nil, errors.Errorf("failed to upload: %s", err)
		}
	}

	err = lease.Complete(vconn.Context)
	if err != nil {
		return nil, errors.Errorf("failed to lease complete: %s", err)
	}

	vm := object.NewVirtualMachine(vconn.Client.Client, info.Entity)
	if vm == nil {
		return nil, fmt.Errorf("error VirtualMachine not found, managed object id: %s", info.Entity.Value)
	}

	err = vm.MarkAsTemplate(vconn.Context)
	if err != nil {
		return nil, errors.Errorf("failed to mark vm as template: %s", err)
	}
	err = attachTag(vconn, vm.Reference().Value, tagId)
	if err != nil {
		return nil, err
	}

	return vm, nil
}

func findAvailableHostSystems(vconn *VCenterConnection, clusterHostSystems []*object.HostSystem) (*object.HostSystem, error) {
	logrus.Infof("In findAvailableHostSystems")
	var hostSystemManagedObject mo.HostSystem
	for _, hostObj := range clusterHostSystems {
		err := hostObj.Properties(vconn.Context, hostObj.Reference(), []string{"config.product", "network", "datastore", "runtime"}, &hostSystemManagedObject)
		if err != nil {
			return nil, err
		}
		if hostSystemManagedObject.Runtime.InMaintenanceMode {
			continue
		}
		return hostObj, nil
	}
	return nil, errors.New("all hosts unavailable")
}

// Used govc/importx/ovf.go as an example to implement
// resourceVspherePrivateImportOvaCreate and upload functions
// See: https://github.com/vmware/govmomi/blob/cc10a0758d5b4d4873388bcea417251d1ad03e42/govc/importx/ovf.go#L196-L324
func upload(ctx context.Context, archive *importx.ArchiveFlag, lease *nfc.Lease, item nfc.FileItem) error {
	logrus.Infof("In upload")
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

func attachTag(vconn *VCenterConnection, vmMoRefValue, tagId string) error {
	logrus.Infof("In attachTag")
	tagManager := tags.NewManager(vconn.RestClient)
	moRef := types.ManagedObjectReference{
		Value: vmMoRefValue,
		Type:  "VirtualMachine",
	}

	err := tagManager.AttachTag(vconn.Context, tagId, moRef)

	if err != nil {
		return err
	}
	return nil
}

const (
	GuestInfoIgnitionData     = "guestinfo.ignition.config.data"
	GuestInfoIgnitionEncoding = "guestinfo.ignition.config.data.encoding"
	GuestInfoHostname         = "guestinfo.hostname"

	// going to ignore for now...
	GuestInfoNetworkKargs = "guestinfo.afterburn.initrd.network-kargs"
	StealClock            = "stealclock.enable"
	ethCardType           = "vmxnet3"
)

func getExtraConfig(vmName, encodedIgnition string) []types.BaseOptionValue {
	return []types.BaseOptionValue{
		&types.OptionValue{
			Key:   GuestInfoIgnitionEncoding,
			Value: "base64",
		},
		&types.OptionValue{
			Key:   GuestInfoIgnitionData,
			Value: encodedIgnition,
		},
		&types.OptionValue{
			Key:   GuestInfoHostname,
			Value: vmName,
		},
		&types.OptionValue{
			Key:   StealClock,
			Value: "TRUE",
		},
	}
}

func clone(vconn *VCenterConnection,
	vmTemplate *object.VirtualMachine,
	machineProviderSpec *machinev1beta1.VSphereMachineProviderSpec,
	encodedIgnition, vmName string) (*object.Task, error) {

	extraConfig := getExtraConfig(vmName, encodedIgnition)

	deviceSpecs := []types.BaseVirtualDeviceConfigSpec{}
	virtualDeviceList, err := vmTemplate.Device(vconn.Context)
	if err != nil {
		return nil, err
	}

	networkDevices, err := getNetworkDevices(vconn, virtualDeviceList, machineProviderSpec)
	if err != nil {
		return nil, err
	}

	deviceSpecs = append(deviceSpecs, networkDevices...)

	diskSpec, err := getDiskSpec(virtualDeviceList, machineProviderSpec)
	if err != nil {
		return nil, err
	}
	deviceSpecs = append(deviceSpecs, diskSpec)

	datastore, err := vconn.Finder.Datastore(vconn.Context, machineProviderSpec.Workspace.Datastore)
	if err != nil {
		return nil, err
	}
	folder, err := vconn.Finder.Folder(vconn.Context, machineProviderSpec.Workspace.Folder)
	if err != nil {
		return nil, err
	}
	resourcepool, err := vconn.Finder.ResourcePool(vconn.Context, machineProviderSpec.Workspace.ResourcePool)

	spec := types.VirtualMachineCloneSpec{
		Config: &types.VirtualMachineConfigSpec{
			Flags:             newVMFlagInfo(),
			ExtraConfig:       extraConfig,
			DeviceChange:      deviceSpecs,
			NumCPUs:           machineProviderSpec.NumCPUs,
			NumCoresPerSocket: machineProviderSpec.NumCoresPerSocket,
			MemoryMB:          machineProviderSpec.MemoryMiB,
		},
		Location: types.VirtualMachineRelocateSpec{
			Datastore: types.NewReference(datastore.Reference()),
			Folder:    types.NewReference(folder.Reference()),
			Pool:      types.NewReference(resourcepool.Reference()),
		},
		PowerOn: false,
	}

	return vmTemplate.Clone(vconn.Context, folder, vmName, spec)
}

func getDiskSpec(devices object.VirtualDeviceList, machineProviderSpec *machinev1beta1.VSphereMachineProviderSpec) (types.BaseVirtualDeviceConfigSpec, error) {
	disks := devices.SelectByType((*types.VirtualDisk)(nil))

	disk := disks[0].(*types.VirtualDisk)
	cloneCapacityKB := int64(machineProviderSpec.DiskGiB) * 1024 * 1024
	disk.CapacityInKB = cloneCapacityKB

	return &types.VirtualDeviceConfigSpec{
		Operation: types.VirtualDeviceConfigSpecOperationEdit,
		Device:    disk,
	}, nil
}

func getNetworkDevices(vconn *VCenterConnection,
	devices object.VirtualDeviceList,
	machineProviderSpec *machinev1beta1.VSphereMachineProviderSpec) ([]types.BaseVirtualDeviceConfigSpec, error) {

	nics := devices.SelectByType((*types.VirtualEthernetCard)(nil))

	nic := nics[0].(*types.VirtualVmxnet3)

	var networkDevices []types.BaseVirtualDeviceConfigSpec

	resourcepool, err := vconn.Finder.ResourcePool(vconn.Context, machineProviderSpec.Workspace.ResourcePool)
	if err != nil {
		return nil, err
	}

	clusterObjRef, err := resourcepool.Owner(vconn.Context)
	if err != nil {
		return nil, err
	}

	computeCluster := clusterObjRef.(*object.ClusterComputeResource)

	netdev := machineProviderSpec.Network.Devices[0]
	networkPath := path.Join(computeCluster.InventoryPath, netdev.NetworkName)
	networkObject, err := vconn.Finder.Network(vconn.Context, networkPath)
	if err != nil {
		return nil, err
	}
	backing, err := networkObject.EthernetCardBackingInfo(vconn.Context)
	if err != nil {
		return nil, err
	}
	nic.Backing = backing

	networkDevices = append(networkDevices, &types.VirtualDeviceConfigSpec{
		Device:    nic,
		Operation: types.VirtualDeviceConfigSpecOperationEdit,
	})
	return networkDevices, nil
}

func newVMFlagInfo() *types.VirtualMachineFlagInfo {
	diskUUIDEnabled := true
	return &types.VirtualMachineFlagInfo{
		DiskUuidEnabled: &diskUUIDEnabled,
	}
}
