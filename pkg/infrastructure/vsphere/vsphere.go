package vsphere

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

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
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/tfvars"
	"github.com/openshift/installer/pkg/tfvars/vsphere"
	typesinstall "github.com/openshift/installer/pkg/types"
	typesvsphere "github.com/openshift/installer/pkg/types/vsphere"
)

const (
	guestInfoIgnitionData     = "guestinfo.ignition.config.data"
	guestInfoIgnitionEncoding = "guestinfo.ignition.config.data.encoding"
	guestInfoHostname         = "guestinfo.hostname"
	guestInfoDomain           = "guestinfo.domain"
	guestInfoNetworkKargs     = "guestinfo.afterburn.initrd.network-kargs"
	stealClock                = "stealclock.enable"
)

type vCenterConnection struct {
	Client     *govmomi.Client
	Finder     *find.Finder
	Context    context.Context
	RestClient *rest.Client
	Logout     func() error

	URI      string
	Username string
	Password string
}

// InfrastructureProvider implements the govmomi-based installation method for vSphere IPI.
type InfrastructureProvider struct{}

// InitializeProvider implements the govmomi-based installation method for vSphere IPI.
func InitializeProvider() infrastructure.Provider {
	return &InfrastructureProvider{}
}

// Provision implements the govmomi-based installation virtual machine provisioning for vSphere IPI.
func (p *InfrastructureProvider) Provision(dir string, vars []*asset.File) ([]*asset.File, error) {
	vsphereConfig, clusterConfig, err := getTerraformVars(dir, vars)
	if err != nil {
		return nil, err
	}
	err = provision(vsphereConfig, clusterConfig)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func getTerraformVars(dir string, vars []*asset.File) (*vsphere.Config, *tfvars.Config, error) {
	vsphereConfig := &vsphere.Config{}
	clusterConfig := &tfvars.Config{}

	for _, v := range vars {
		var err error

		filePath := path.Join(dir, v.Filename)
		file, err := os.Open(filePath)

		if err != nil {
			return nil, nil, err
		}

		// decoder provides a rational error message if the json is screwed up.
		// whereas Unmarshal does not
		decoder := json.NewDecoder(file)
		decoder.DisallowUnknownFields()

		if v.Filename == "terraform.tfvars.json" {
			err = decoder.Decode(clusterConfig)
		}
		if v.Filename == "terraform.platform.auto.tfvars.json" {
			err = decoder.Decode(vsphereConfig)
		}

		if err != nil {
			return nil, nil, err
		}
	}
	return vsphereConfig, clusterConfig, nil
}

func getAssetFilesFromDir(dir string) ([]*asset.File, error) {
	var vars []*asset.File
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() {
			if strings.Contains(file.Name(), "tfvars.json") {
				vars = append(vars, &asset.File{Filename: file.Name()})
			}
		}
	}
	return vars, nil
}

// DestroyBootstrap implements the govmomi-based deletion of the bootstrap virtual machine for vSphere IPI.
func (p *InfrastructureProvider) DestroyBootstrap(dir string) error {
	vars, err := getAssetFilesFromDir(dir)
	if err != nil {
		return err
	}

	vsphereConfig, clusterConfig, err := getTerraformVars(dir, vars)
	if err != nil {
		return err
	}

	vcenterConnectionMap, err := getvCenterConnections(vsphereConfig)
	if err != nil {
		return err
	}

	for _, vconn := range vcenterConnectionMap {
		for _, fd := range vsphereConfig.FailureDomainMap {
			if vconn.URI == fd.Server {
				bootstrapName := fmt.Sprintf("%s-bootstrap", clusterConfig.ClusterID)

				bootstrapPath := path.Join(fd.Topology.Folder, bootstrapName)

				dc, err := vconn.Finder.Datacenter(vconn.Context, fd.Topology.Datacenter)

				if err != nil {
					return err
				}

				vconn.Finder = vconn.Finder.SetDatacenter(dc)

				virtualMachines, err := vconn.Finder.VirtualMachineList(vconn.Context, bootstrapPath)
				if err != nil {
					return err
				}

				for _, vm := range virtualMachines {
					if vm.Name() == bootstrapName {
						logrus.Debugf("Destroying bootstrap virtual machine: %s", bootstrapName)
						task, err := vm.PowerOff(vconn.Context)
						if err != nil {
							return err
						}

						err = task.Wait(vconn.Context)
						if err != nil {
							return err
						}

						// We don't need to wait for the destroy to finish.
						_, err = vm.Destroy(vconn.Context)
						if err != nil {
							return err
						}
						return nil
					}
				}
			}
		}
	}

	return nil
}

// ExtractHostAddresses is used for sdk (govmomi) vSphere IPI installations to determine the ip addresses for
// bootstrap and control plane nodes if possible.
func (p *InfrastructureProvider) ExtractHostAddresses(dir string, config *typesinstall.InstallConfig, ha *infrastructure.HostAddresses) error {
	vars, err := getAssetFilesFromDir(dir)
	if err != nil {
		return err
	}

	// there should be at least two terraform variable files there.
	if len(vars) >= 2 {
		vsphereConfig, clusterConfig, err := getTerraformVars(dir, vars)

		if err != nil {
			return err
		}

		vcenterConnectionMap, err := getvCenterConnections(vsphereConfig)
		if err != nil {
			return err
		}
		for _, vconn := range vcenterConnectionMap {
			for _, fd := range config.Platform.VSphere.FailureDomains { // we mess with FailureDomains within vsphereConfig
				if vconn.URI == fd.Server {
					dc, err := vconn.Finder.Datacenter(vconn.Context, fd.Topology.Datacenter)
					if err != nil {
						return err
					}
					dcFolders, err := dc.Folders(vconn.Context)
					if err != nil {
						return err
					}

					folderPath := fd.Topology.Folder
					if !path.IsAbs(folderPath) {
						// we futz with the folder path _somewhere_
						// probably because folder in terraform only takes -datacenter/vm/- <-- folder name
						folderPath = fmt.Sprintf("%s*", path.Join(dcFolders.VmFolder.InventoryPath, fd.Topology.Folder, clusterConfig.ClusterID))
					}

					vms, err := vconn.Finder.VirtualMachineList(vconn.Context, folderPath)
					if err != nil {
						return err
					}

					bootstrapName := fmt.Sprintf("%s-bootstrap", clusterConfig.ClusterID)
					masterBaseName := fmt.Sprintf("%s-master", clusterConfig.ClusterID)

					for _, vm := range vms {
						ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
						defer cancel()

						ip, err := vm.WaitForIP(ctx, true)
						if err != nil {
							// we don't care about the context exceeding because maybe the node doesn't have an IP address
							if !errors.Is(err, context.DeadlineExceeded) {
								return err
							}
						}

						if ip != "" {
							if strings.Contains(vm.Name(), bootstrapName) {
								ha.Bootstrap = ip
							}
							if strings.Contains(vm.Name(), masterBaseName) {
								ha.Masters = append(ha.Masters, ip)
							}
						}
					}
				}
			}
		}
	} else {
		return errors.Errorf("missing variable files to determine bootstrap and/or control plane ip addressing")
	}

	return nil
}

func getVCenterClient(uri, username, password string) (*vCenterConnection, error) {
	logrus.Debugf("Connecting to vCenter: %s with username: %s", uri, username)

	ctx := context.Background()

	connection := &vCenterConnection{
		Context: ctx,
	}

	u, err := soap.ParseURL(uri)
	if err != nil {
		return nil, err
	}
	connection.Username = username
	connection.Password = password
	connection.URI = uri

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
	connection.Logout = func() error {
		err := connection.Client.Logout(connection.Context)
		if err != nil {
			return err
		}
		err = connection.RestClient.Logout(connection.Context)
		if err != nil {
			return err
		}
		return nil
	}
	return connection, nil
}

func getvCenterConnections(vsphereConfig *vsphere.Config) (map[string]*vCenterConnection, error) {
	vcenterConnectionMap := make(map[string]*vCenterConnection)

	for _, v := range vsphereConfig.VCenters {
		tempvCenterConnection, err := getVCenterClient(
			v.Server,
			v.Username,
			v.Password)

		if err != nil {
			return nil, err
		}
		vcenterConnectionMap[v.Server] = tempvCenterConnection
	}

	return vcenterConnectionMap, nil
}

func provision(vsphereConfig *vsphere.Config, clusterConfig *tfvars.Config) error {
	vmTemplateMap := make(map[string]*object.VirtualMachine)
	tagMap := make(map[string]string)

	vcenterConnectionMap, err := getvCenterConnections(vsphereConfig)
	if err != nil {
		return err
	}

	// each vcenter needs a tag and tag category
	for _, v := range vcenterConnectionMap {
		categoryID, err := createTagCategory(v, clusterConfig.ClusterID)
		if err != nil {
			return err
		}

		tempTag, err := createTag(v, clusterConfig.ClusterID, categoryID)
		if err != nil {
			return err
		}

		tagMap[v.URI] = tempTag
	}

	for _, fd := range vsphereConfig.FailureDomainMap {
		var vmName string
		var vmTemplate *object.VirtualMachine
		vcenterConnection := vcenterConnectionMap[fd.Server]

		dc, err := vcenterConnection.Finder.Datacenter(vcenterConnection.Context, fd.Topology.Datacenter)
		if err != nil {
			return err
		}
		dcFolders, err := dc.Folders(vcenterConnection.Context)
		if err != nil {
			return err
		}

		folderPath := path.Join(dcFolders.VmFolder.InventoryPath, clusterConfig.ClusterID)

		// we must set the Folder to the infraId somewhere, we will need to remove that.
		// if we are overwriting folderPath it needs to have a slash (path)
		if strings.Contains(fd.Topology.Folder, "/") {
			folderPath = fd.Topology.Folder
		}

		folder, err := createFolder(folderPath, vcenterConnection)
		if err != nil {
			return err
		}

		// Not entirely fond of this being the switch between using existing template
		// and importing. I _think_ the better option would be to use the installConfig FailureDomain directly
		// though I guess that isn't available. Maybe add a terrform variable parameter for an unmodified
		// platform spec.

		if !path.IsAbs(fd.Topology.Template) { // scenario, if Template is just a name then upload
			vmTemplate, err = importRhcosOva(vcenterConnection, folder,
				vsphereConfig.OvaFilePath, clusterConfig.ClusterID, tagMap[fd.Server], string(vsphereConfig.DiskType), fd)
			if err != nil {
				return err
			}
			vmName, err = vmTemplate.ObjectName(vcenterConnection.Context)

			if err != nil {
				return err
			}
		} else { // scenario, if Template is a full path use existing
			vmTemplate, err = vcenterConnection.Finder.VirtualMachine(vcenterConnection.Context, fd.Topology.Template)

			if err != nil {
				return err
			}
			// if we use pre-existing template the full path is provided in the machine.
			vmName = vmTemplate.InventoryPath
		}
		vmTemplateMap[vmName] = vmTemplate
	}
	encodedMasterIgn := base64.StdEncoding.EncodeToString([]byte(clusterConfig.IgnitionMaster))
	encodedBootstrapIgn := base64.StdEncoding.EncodeToString([]byte(clusterConfig.IgnitionBootstrap))

	bootstrap := true

	for i := 0; i < len(vsphereConfig.ControlPlanes); i++ {
		var kargs string
		cp := vsphereConfig.ControlPlanes[i]
		vcenterConnection := vcenterConnectionMap[cp.Workspace.Server]

		vmName := fmt.Sprintf("%s-master-%d", clusterConfig.ClusterID, i)

		encodedIgnition := encodedMasterIgn

		if len(vsphereConfig.ControlPlaneNetworkKargs) > i {
			kargs = vsphereConfig.ControlPlaneNetworkKargs[i]
		}

		if bootstrap {
			if i == 0 {
				encodedIgnition = encodedBootstrapIgn
				vmName = fmt.Sprintf("%s-bootstrap", clusterConfig.ClusterID)
				kargs = vsphereConfig.BootStrapNetworkKargs
			}
		}

		task, err := clone(vcenterConnection, vmTemplateMap[cp.Template], cp, encodedIgnition, vmName, clusterConfig.ClusterDomain, kargs)
		if err != nil {
			return err
		}

		taskInfo, err := task.WaitForResult(vcenterConnection.Context, nil)
		if err != nil {
			return err
		}

		vmMoRef, ok := taskInfo.Result.(types.ManagedObjectReference)
		if !ok {
			return errors.New("unable to convert task info result into managed object reference")
		}
		vm := object.NewVirtualMachine(vcenterConnection.Client.Client, vmMoRef)

		err = attachTag(vcenterConnectionMap[cp.Workspace.Server], vmMoRef.Value, tagMap[cp.Workspace.Server])
		if err != nil {
			return err
		}

		datacenter, err := vcenterConnection.Finder.Datacenter(vcenterConnection.Context, cp.Workspace.Datacenter)
		if err != nil {
			return err
		}

		task, err = datacenter.PowerOnVM(vcenterConnection.Context, []types.ManagedObjectReference{vm.Reference()}, &types.OptionValue{
			Key:   string(types.ClusterPowerOnVmOptionOverrideAutomationLevel),
			Value: string(types.DrsBehaviorFullyAutomated),
		})
		if err != nil {
			return err
		}

		_, err = task.WaitForResult(vcenterConnection.Context, nil)
		if err != nil {
			return err
		}

		if bootstrap {
			if i == 0 {
				bootstrap = false
				i = -1
			}
		}
	}

	for _, v := range vcenterConnectionMap {
		err := v.Logout()
		if err != nil {
			return err
		}
	}
	return nil
}

func createFolder(fullpath string, vconn *vCenterConnection) (*object.Folder, error) {
	dir := path.Dir(fullpath)
	base := path.Base(fullpath)

	folder, err := vconn.Finder.Folder(context.TODO(), fullpath)

	if folder == nil {
		folder, err = vconn.Finder.Folder(context.TODO(), dir)

		var notFoundError *find.NotFoundError
		if errors.As(err, &notFoundError) {
			folder, err = createFolder(dir, vconn)
			if err != nil {
				return folder, err
			}
		}

		if folder != nil {
			logrus.Debugf("Creating vCenter folder: %s", base)
			folder, err = folder.CreateFolder(context.TODO(), base)
			if err != nil {
				return folder, err
			}
		}
	}
	return folder, err
}

func createTagCategory(vconn *vCenterConnection, clusterID string) (string, error) {
	categoryName := fmt.Sprintf("openshift-%s", clusterID)
	logrus.Debugf("Creating vCenter tag category: %s", categoryName)

	category := tags.Category{
		Name:        categoryName,
		Description: "Added by openshift-install do not remove",
		Cardinality: "SINGLE",
		AssociableTypes: []string{
			"urn:vim25:VirtualMachine",
			"urn:vim25:ResourcePool",
			"urn:vim25:Folder",
			"urn:vim25:Datastore",
			"urn:vim25:StoragePod",
		},
	}

	return tags.NewManager(vconn.RestClient).CreateCategory(vconn.Context, &category)
}

func createTag(vconn *vCenterConnection, clusterID, categoryID string) (string, error) {
	logrus.Debugf("Creating vCenter tag: %s", clusterID)

	tag := tags.Tag{
		Description: "Added by openshift-install do not remove",
		Name:        clusterID,
		CategoryID:  categoryID,
	}

	return tags.NewManager(vconn.RestClient).CreateTag(vconn.Context, &tag)
}

func importRhcosOva(vconn *vCenterConnection, folder *object.Folder, cachedImage, clusterID, tagID, diskProvisioningType string, failureDomain typesvsphere.FailureDomain) (*object.VirtualMachine, error) {
	name := fmt.Sprintf("%s-rhcos-%s-%s", clusterID, failureDomain.Region, failureDomain.Zone)
	logrus.Debugf("Importing RHCOS OVA: %s", name)

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
	if err != nil {
		return nil, err
	}

	networkPath := path.Join(cluster.InventoryPath, failureDomain.Topology.Networks[0])

	networkRef, err := vconn.Finder.Network(vconn.Context, networkPath)
	if err != nil {
		return nil, err
	}
	datastore, err := vconn.Finder.Datastore(vconn.Context, failureDomain.Topology.Datastore)
	if err != nil {
		return nil, err
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
	err = attachTag(vconn, vm.Reference().Value, tagID)
	if err != nil {
		return nil, err
	}

	return vm, nil
}

func findAvailableHostSystems(vconn *vCenterConnection, clusterHostSystems []*object.HostSystem) (*object.HostSystem, error) {
	logrus.Debug("Finding available ESXi hosts for OVA importing.")
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

func attachTag(vconn *vCenterConnection, vmMoRefValue, tagID string) error {
	logrus.Debugf("Attaching tag id: %s to virtual machine managed object id: %s", tagID, vmMoRefValue)

	tagManager := tags.NewManager(vconn.RestClient)
	moRef := types.ManagedObjectReference{
		Value: vmMoRefValue,
		Type:  "VirtualMachine",
	}

	err := tagManager.AttachTag(vconn.Context, tagID, moRef)

	if err != nil {
		return err
	}
	return nil
}

func getExtraConfig(vmName, clusterDomain, encodedIgnition, kargs string) []types.BaseOptionValue {
	extraConfig := []types.BaseOptionValue{
		&types.OptionValue{
			Key:   guestInfoIgnitionEncoding,
			Value: "base64",
		},
		&types.OptionValue{
			Key:   guestInfoIgnitionData,
			Value: encodedIgnition,
		},
		&types.OptionValue{
			Key:   guestInfoHostname,
			Value: vmName,
		},
		&types.OptionValue{
			Key:   stealClock,
			Value: "TRUE",
		},
		&types.OptionValue{
			Key:   guestInfoDomain,
			Value: clusterDomain,
		},
	}

	if kargs != "" {
		extraConfig = append(extraConfig, &types.OptionValue{
			Key:   guestInfoNetworkKargs,
			Value: kargs,
		})
	}

	return extraConfig
}

func clone(vconn *vCenterConnection, vmTemplate *object.VirtualMachine, machineProviderSpec *machinev1beta1.VSphereMachineProviderSpec, encodedIgnition, vmName, clusterDomain, kargs string) (*object.Task, error) {
	extraConfig := getExtraConfig(vmName, clusterDomain, encodedIgnition, kargs)

	var deviceSpecs []types.BaseVirtualDeviceConfigSpec
	virtualDeviceList, err := vmTemplate.Device(vconn.Context)
	if err != nil {
		return nil, err
	}

	diskSpec, err := getDiskSpec(virtualDeviceList, machineProviderSpec)
	if err != nil {
		return nil, err
	}
	deviceSpecs = append(deviceSpecs, diskSpec)

	networkDevices, err := getNetworkDevices(vconn, virtualDeviceList, machineProviderSpec)
	if err != nil {
		return nil, err
	}

	deviceSpecs = append(deviceSpecs, networkDevices...)

	datastore, err := vconn.Finder.Datastore(vconn.Context, machineProviderSpec.Workspace.Datastore)
	if err != nil {
		return nil, err
	}
	folder, err := vconn.Finder.Folder(vconn.Context, machineProviderSpec.Workspace.Folder)
	if err != nil {
		return nil, err
	}
	resourcepool, err := vconn.Finder.ResourcePool(vconn.Context, machineProviderSpec.Workspace.ResourcePool)
	if err != nil {
		return nil, err
	}

	diskUUIDEnabled := true
	spec := types.VirtualMachineCloneSpec{
		Config: &types.VirtualMachineConfigSpec{
			Flags: &types.VirtualMachineFlagInfo{
				DiskUuidEnabled: &diskUUIDEnabled,
			},
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

	disk, ok := disks[0].(*types.VirtualDisk)
	if !ok {
		return nil, errors.New("unable to convert disks to VirtualDisk type")
	}
	cloneCapacityKB := int64(machineProviderSpec.DiskGiB) * 1024 * 1024
	disk.CapacityInKB = cloneCapacityKB

	return &types.VirtualDeviceConfigSpec{
		Operation: types.VirtualDeviceConfigSpecOperationEdit,
		Device:    disk,
	}, nil
}

func getNetworkDevices(
	vconn *vCenterConnection,
	devices object.VirtualDeviceList,
	machineProviderSpec *machinev1beta1.VSphereMachineProviderSpec) ([]types.BaseVirtualDeviceConfigSpec, error) {
	var networkDevices []types.BaseVirtualDeviceConfigSpec

	nics := devices.SelectByType(&types.VirtualEthernetCard{})

	nic, ok := nics[0].(*types.VirtualVmxnet3)
	if !ok {
		return nil, errors.New("unable to convert nic to VirtualVmxnet3 type")
	}

	// I am sure there is a better way to do this...
	networkType := "Network"
	if strings.Contains(machineProviderSpec.Network.Devices[0].NetworkName, "dv") {
		networkType = "DistributedVirtualPortgroup"
	}
	networkObjRef := types.ManagedObjectReference{
		Value: machineProviderSpec.Network.Devices[0].NetworkName,
		Type:  networkType,
	}

	// if this doesn't error with NotFoundError, then the NetworkName
	// in the ManagedObjectReference is a Value string vs a path
	networkObject, err := vconn.Finder.ObjectReference(vconn.Context, networkObjRef)
	if err != nil {
		return nil, err
	}

	var backing types.BaseVirtualDeviceBackingInfo

	switch networkObject := networkObject.(type) {
	case object.DistributedVirtualPortgroup:
		backing, err = networkObject.EthernetCardBackingInfo(vconn.Context)
	case object.Network:
		backing, err = networkObject.EthernetCardBackingInfo(vconn.Context)
	}
	if err != nil {
		return nil, err
	}

	// These operations for network adapter is very similar to what govc does when cloning.
	newNicDevice, err := object.EthernetCardTypes().CreateEthernetCard("vmxnet3", backing)
	if err != nil {
		return nil, err
	}
	card := newNicDevice.(types.BaseVirtualEthernetCard).GetVirtualEthernetCard()
	card.Key = int32(1)

	card.MacAddress = ""
	card.AddressType = string(types.VirtualEthernetCardMacTypeGenerated)

	nic.Backing = card.Backing

	networkDevices = append(networkDevices, &types.VirtualDeviceConfigSpec{
		Device:    nic,
		Operation: types.VirtualDeviceConfigSpecOperationEdit,
	})
	return networkDevices, nil
}
