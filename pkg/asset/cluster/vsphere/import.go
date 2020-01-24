package vsphere

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	//_ "github.com/sirupsen/logrus"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/nfc"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/ovf"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"

	"github.com/vmware/govmomi/vim25/soap"

	installertypes "github.com/openshift/installer/pkg/types"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

// ImportOvaParams contains the vCenter objects required to import a OVA into vSphere.
type ImportOvaParams struct {
	ResourcePool *object.ResourcePool
	Datacenter   *object.Datacenter
	Datastore    *object.Datastore
	Network      *object.Network
	Host         *object.HostSystem
	Folder       *object.Folder
}

func findImportOvaParams(client *vim25.Client, datacenter, cluster, datastore, network string) (*ImportOvaParams, error) {
	var ccrMo mo.ClusterComputeResource
	ctx := context.TODO()
	importOvaParams := &ImportOvaParams{}

	finder := find.NewFinder(client)

	// Find the object Datacenter by using its name provided by install-config
	dcObj, err := finder.Datacenter(ctx, datacenter)
	if err != nil {
		return nil, err
	}
	importOvaParams.Datacenter = dcObj

	// Find the top-level (and hidden to view) folders in the
	// datacenter
	folders, err := importOvaParams.Datacenter.Folders(ctx)
	if err != nil {
		return nil, err
	}
	// The only folder we are interested in is VmFolder
	// Which can contain our template
	importOvaParams.Folder = folders.VmFolder

	clusterPath := fmt.Sprintf("/%s/host/%s", datacenter, cluster)

	// Find the cluster object by the datacenter and cluster name to
	// generate the path e.g. /datacenter/host/cluster
	clusterComputeResource, err := finder.ClusterComputeResource(ctx, clusterPath)
	if err != nil {
		return nil, err
	}

	// Get the network properties that is defined in ClusterComputeResource
	// We need to know if the network name provided exists in the cluster that was
	// also provided.
	err = clusterComputeResource.Properties(context.TODO(), clusterComputeResource.Reference(), []string{"network"}, &ccrMo)
	if err != nil {
		return nil, err
	}

	// Find the network object using the provided network name
	for _, networkMoRef := range ccrMo.Network {
		networkObj := object.NewNetwork(client, networkMoRef)
		networkObjectName, err := networkObj.ObjectName(ctx)
		if err != nil {
			return nil, err
		}
		if network == networkObjectName {
			importOvaParams.Network = networkObj
			break
		}
	}

	// Find all the datastores that are configured under the cluster
	datastores, err := clusterComputeResource.Datastores(ctx)
	if err != nil {
		return nil, err
	}

	// Find the specific datastore by the name provided
	for _, datastoreObj := range datastores {
		datastoreObjName, err := datastoreObj.ObjectName(ctx)
		if err != nil {
			return nil, err
		}
		if datastore == datastoreObjName {
			importOvaParams.Datastore = datastoreObj
			break
		}
	}

	// Find all the HostSystem(s) under cluster
	hosts, err := clusterComputeResource.Hosts(ctx)
	if err != nil {
		return nil, err
	}
	foundDatastore := false
	foundNetwork := false
	var hostSystemManagedObject mo.HostSystem

	// Confirm that the network and datastore that was provided is
	// available for use on the HostSystem we will import the
	// OVA to.
	for _, hostObj := range hosts {
		hostObj.Properties(ctx, hostObj.Reference(), []string{"network", "datastore"}, &hostSystemManagedObject)

		if err != nil {
			return nil, err
		}
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
			importOvaParams.Host = hostObj
			resourcePool, err := hostObj.ResourcePool(ctx)
			if err != nil {
				return nil, err
			}
			importOvaParams.ResourcePool = resourcePool
		}
	}
	if !foundDatastore {
		return nil, errors.Errorf("The hosts in the cluster do not have the datastore provided in install-config.yaml")
	}
	if !foundNetwork {
		return nil, errors.Errorf("The hosts in the cluster do not have the network provided in install-config.yaml")
	}

	return importOvaParams, nil
}

func importOva(config *installertypes.InstallConfig, infraID string) error {

	virtualMachineName := infraID + "-rhcos"
	ctx := context.TODO()
	// Login to vCenter, rest is not used in this case
	client, _, err := vspheretypes.CreateVSphereClients(ctx,
		config.VSphere.VCenter,
		config.VSphere.Username,
		config.VSphere.Password)

	if err != nil {
		return err
	}

	importOvaParams, err := findImportOvaParams(client,
		config.VSphere.Datacenter,
		config.VSphere.Cluster,
		config.VSphere.DefaultDatastore,
		config.VSphere.Network)

	if err != nil {
		return err
	}

	ovaTapeArchive := &TapeArchive{Path: config.VSphere.ClusterOSImage}
	ovaTapeArchive.Client = client

	archive := &ArchiveFlag{}
	archive.Archive = ovaTapeArchive

	ovfDescriptor, err := archive.ReadOvf("desc.ovf")
	if err != nil {
		return err
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
	// Create mapping between OVF and the network object
	// found by Name
	networkMappings := []types.OvfNetworkMapping{{
		Name:    ovfEnvelope.Network.Networks[0].Name,
		Network: importOvaParams.Network.Reference(),
	}}
	// This is a very minimal spec for importing
	// an OVF.
	cisp := types.OvfCreateImportSpecParams{
		EntityName:     virtualMachineName,
		NetworkMapping: networkMappings,
	}

	m := ovf.NewManager(client)
	spec, err := m.CreateImportSpec(ctx,
		string(ovfDescriptor),
		importOvaParams.ResourcePool.Reference(),
		importOvaParams.Datastore.Reference(),
		cisp)

	if err != nil {
		return err
	}
	if spec.Error != nil {
		return errors.New(spec.Error[0].LocalizedMessage)
	}

	//Creates a new entity in this resource pool.
	//See VMware vCenter API documentation: Managed Object - ResourcePool - ImportVApp
	lease, err := importOvaParams.ResourcePool.ImportVApp(ctx,
		spec.ImportSpec,
		importOvaParams.Folder,
		importOvaParams.Host)

	if err != nil {
		return err
	}

	info, err := lease.Wait(ctx, spec.FileItem)
	if err != nil {
		return err
	}

	u := lease.StartUpdater(ctx, info)
	defer u.Done()

	for _, i := range info.Items {
		// upload the vmdk to which ever host that was first
		// available with the required network and datastore.
		err = upload(ctx, archive, lease, i)
		if err != nil {
			return err
		}
	}
	err = lease.Complete(ctx)
	if err != nil {
		return err
	}

	return nil
}

func upload(ctx context.Context, archive *ArchiveFlag, lease *nfc.Lease, item nfc.FileItem) error {
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
