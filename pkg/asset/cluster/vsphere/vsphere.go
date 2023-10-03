package vsphere

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/govc/importx"
	"github.com/vmware/govmomi/nfc"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/ovf"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"

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

func PreTerraform(cachedImage string, clusterID string, installConfig *installconfig.InstallConfig) error {
	vconn, err := getVCenterClient(
		installConfig.Config.VSphere.VCenters[0].Server,
		installConfig.Config.VSphere.VCenters[0].Username,
		installConfig.Config.VSphere.VCenters[0].Password)

	if err != nil {
		return err
	}
	defer vconn.Logout()

	// before the ova is uploaded we need to create
	// oh and the folder...

	for _, fd := range installConfig.Config.VSphere.FailureDomains {
		//createTagCategory(rest)
		//createTag(rest)

		dc, err := vconn.Finder.Datacenter(vconn.Context, fd.Topology.Datacenter)
		if err != nil {
			return err
		}
		dcFolders, err := dc.Folders(vconn.Context)

		folderPath := path.Join(dcFolders.VmFolder.InventoryPath, clusterID)
		if fd.Topology.Folder != "" {
			folderPath = fd.Topology.Folder

		}

		folder, err := createFolder(folderPath, vconn)
		if err != nil {
			return err
		}
		err = importRhcosOva(vconn, folder, cachedImage, clusterID, fd)

		if err != nil {
			return err
		}
	}

	return nil
}

func createFolder(fullpath string, vconn *VCenterConnection) (*object.Folder, error) {
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

func createTagCategory(vconn *VCenterConnection, clusterId string, failureDomain typesvsphere.FailureDomain) error {

	// why is this file missing from the branch?
	return nil
}

func createTag(vconn *VCenterConnection, clusterId string, failureDomain typesvsphere.FailureDomain) error {

	return nil
}

func importRhcosOva(vconn *VCenterConnection, folder *object.Folder, cachedImage, clusterId string, failureDomain typesvsphere.FailureDomain) error {
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

		return errors.Errorf("ova %s has a sha256 of %x and a size of %d bytes, failed to read the ovf descriptor %s", cachedImage, h.Sum(nil), written, err)
	}

	ovfEnvelope, err := archive.ReadEnvelope(ovfDescriptor)
	if err != nil {
		return errors.Errorf("failed to parse ovf: %s", err)
	}

	// The RHCOS OVA only has one network defined by default
	// The OVF envelope defines this.  We need a 1:1 mapping
	// between networks with the OVF and the host
	if len(ovfEnvelope.Network.Networks) != 1 {
		return errors.Errorf("Expected the OVA to only have a single network adapter")
	}

	cluster, err := vconn.Finder.ClusterComputeResource(vconn.Context, failureDomain.Topology.ComputeCluster)

	if err != nil {
		return err
	}
	clusterHostSystems, err := cluster.Hosts(vconn.Context)

	if err != nil {
		return err
	}
	resourcePool, err := vconn.Finder.ResourcePool(vconn.Context, failureDomain.Topology.ResourcePool)

	networkPath := path.Join(cluster.InventoryPath, failureDomain.Topology.Networks[0])

	networkRef, err := vconn.Finder.Network(vconn.Context, networkPath)
	if err != nil {
		return err
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
	/*
		switch {
		case "":
			// Disk provisioning type will be set according to the default storage policy of vsphere.
		case "thin":
			cisp.DiskProvisioning = string(types.OvfCreateImportSpecParamsDiskProvisioningTypeThin)
		case "thick":
			cisp.DiskProvisioning = string(types.OvfCreateImportSpecParamsDiskProvisioningTypeThick)
		case "eagerZeroedThick":
			cisp.DiskProvisioning = string(types.OvfCreateImportSpecParamsDiskProvisioningTypeEagerZeroedThick)
		default:
			return errors.Errorf("Disk provisioning type %q is not supported.", diskType)
		}

	*/

	m := ovf.NewManager(vconn.Client.Client)
	spec, err := m.CreateImportSpec(vconn.Context,
		string(ovfDescriptor),
		resourcePool.Reference(),
		datastore.Reference(),
		cisp)

	if err != nil {
		return errors.Errorf("failed to create import spec: %s", err)
	}
	if spec.Error != nil {
		return errors.New(spec.Error[0].LocalizedMessage)
	}

	hostSystem, err := findAvailableHostSystems(vconn, clusterHostSystems)

	if err != nil {
		return err
	}
	//Creates a new entity in this resource pool.
	//See VMware vCenter API documentation: Managed Object - ResourcePool - ImportVApp
	lease, err := resourcePool.ImportVApp(vconn.Context, spec.ImportSpec, folder, hostSystem)

	if err != nil {
		return errors.Errorf("failed to import vapp: %s", err)
	}

	info, err := lease.Wait(vconn.Context, spec.FileItem)
	if err != nil {
		return errors.Errorf("failed to lease wait: %s", err)
	}

	if err != nil {
		return errors.Errorf("failed to attach tag to virtual machine: %s", err)
	}

	u := lease.StartUpdater(vconn.Context, info)
	defer u.Done()

	for _, i := range info.Items {
		// upload the vmdk to which ever host that was first
		// available with the required network and datastore.
		err = upload(vconn.Context, archive, lease, i)
		if err != nil {
			return errors.Errorf("failed to upload: %s", err)
		}
	}

	err = lease.Complete(vconn.Context)
	if err != nil {
		return errors.Errorf("failed to lease complete: %s", err)
	}

	vm := object.NewVirtualMachine(vconn.Client.Client, info.Entity)
	if vm == nil {
		return fmt.Errorf("error VirtualMachine not found, managed object id: %s", info.Entity.Value)
	}

	err = vm.MarkAsTemplate(vconn.Context)
	if err != nil {
		return errors.Errorf("failed to mark vm as template: %s", err)
	}

	return nil
}

func findAvailableHostSystems(vconn *VCenterConnection, clusterHostSystems []*object.HostSystem) (*object.HostSystem, error) {
	var hostSystemManagedObject mo.HostSystem
	for _, hostObj := range clusterHostSystems {
		err := hostObj.Properties(vconn.Context, hostObj.Reference(), []string{"config.product", "network", "datastore", "runtime"}, &hostSystemManagedObject)
		if err != nil {
			return nil, err
		}

		// Skip all hosts that are in maintenance mode.
		if hostSystemManagedObject.Runtime.InMaintenanceMode {
			continue
		}

		return hostObj, nil

		// these checks should not be here, if anything they should be in validation

		/*
			for _, dsMoRef := range hostSystemManagedObject.Datastore {

				if importOvaParams.Datastore.Reference().Value == dsMoRef.Value {
					foundDatastore = true
					break
				}
			}
			for _, nMoRef := range hostSystemManagedObject.Network {
				if importOvaParams.Network.Reference().Value == nMoRef.Value {
					foundNetwork = true
					break
				}
			}

			if foundDatastore && foundNetwork {
				return importOvaParams, nil
			}
		*/
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

/*
func attachTag(rest *rest.Client) error {
	ctx, cancel := context.WithTimeout(context.TODO(), defaultAPITimeout)
	defer cancel()
	tagManager := tags.NewManager(rest)
	moRef := types.ManagedObjectReference{
		Value: d.Id(),
		Type:  "VirtualMachine",
	}

	err := tagManager.AttachTag(ctx, d.Get("tag").(string), moRef)

	if err != nil {
		return err
	}
	return nil
}

*/
