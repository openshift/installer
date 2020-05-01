package vsphereprivate

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/pkg/errors"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/nfc"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/ovf"
	"github.com/vmware/govmomi/vapi/tags"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"

	"github.com/vmware/govmomi/govc/importx"
)

func resourceVSpherePrivateImportOva() *schema.Resource {
	return &schema.Resource{
		Create:        resourceVSpherePrivateImportOvaCreate,
		Read:          resourceVSpherePrivateImportOvaRead,
		Delete:        resourceVSpherePrivateImportOvaDelete,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "The name of the virtual machine that will be created.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"filename": {
				Type:         schema.TypeString,
				Description:  "The filename path to the ova file to be imported.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"datacenter": {
				Type:         schema.TypeString,
				Description:  "The name of the datacenter.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"cluster": {
				Type:         schema.TypeString,
				Description:  "The name of the cluster.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"network": {
				Type:         schema.TypeString,
				Description:  "The name of a network that the virtual machine will use.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"datastore": {
				Type:         schema.TypeString,
				Description:  "The name of the virtual machine's datastore.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"folder": {
				Type:         schema.TypeString,
				Description:  "The name of the folder to locate the virtual machine in.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"tag": {
				Type:         schema.TypeString,
				Description:  "The name of the tag to attach the virtual machine in.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}

// importOvaParams contains the vCenter objects required to import a OVA into vSphere.
type importOvaParams struct {
	ResourcePool *object.ResourcePool
	Datacenter   *object.Datacenter
	Datastore    *object.Datastore
	Network      *object.Network
	Host         *object.HostSystem
	Folder       *object.Folder
}

func findImportOvaParams(client *vim25.Client, datacenter, cluster, datastore, network, folder string) (*importOvaParams, error) {
	var ccrMo mo.ClusterComputeResource

	ctx, cancel := context.WithTimeout(context.TODO(), defaultAPITimeout)
	defer cancel()

	importOvaParams := &importOvaParams{}
	finder := find.NewFinder(client)

	// Find the object Datacenter by using its name provided by install-config
	dcObj, err := finder.Datacenter(ctx, datacenter)
	if err != nil {
		return nil, err
	}
	importOvaParams.Datacenter = dcObj

	// Create an absolute path to the folder in case the provided folder is nested.
	folderPath := fmt.Sprintf("/%s/vm/%s", datacenter, folder)
	folderObj, err := finder.Folder(ctx, folderPath)
	if err != nil {
		return nil, err
	}
	importOvaParams.Folder = folderObj

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
	err = clusterComputeResource.Properties(ctx, clusterComputeResource.Reference(), []string{"network"}, &ccrMo)
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
		return nil, errors.Errorf("failed to find a host in the cluster that contains the provided datastore")
	}
	if !foundNetwork {
		return nil, errors.Errorf("failed to find a host in the cluster that contains the provided network")
	}

	return importOvaParams, nil
}

func attachTag(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := context.WithTimeout(context.TODO(), defaultAPITimeout)
	defer cancel()
	tagManager := tags.NewManager(meta.(*VSphereClient).restClient)
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

func resourceVSpherePrivateImportOvaCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning import ova create", d.Get("filename").(string))

	ctx, cancel := context.WithTimeout(context.TODO(), defaultAPITimeout)
	defer cancel()
	client := meta.(*VSphereClient).vimClient.Client
	archive := &importx.ArchiveFlag{Archive: &importx.TapeArchive{Path: d.Get("filename").(string)}}

	importOvaParams, err := findImportOvaParams(client,
		d.Get("datacenter").(string),
		d.Get("cluster").(string),
		d.Get("datastore").(string),
		d.Get("network").(string),
		d.Get("folder").(string))
	if err != nil {
		return errors.Errorf("failed to find provided vSphere objects: %s", err)
	}

	ovfDescriptor, err := archive.ReadOvf("*.ovf")
	if err != nil {
		return errors.Errorf("failed to read ovf: %s", err)
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
		EntityName:     d.Get("name").(string),
		NetworkMapping: networkMappings,
	}

	m := ovf.NewManager(client)
	spec, err := m.CreateImportSpec(ctx,
		string(ovfDescriptor),
		importOvaParams.ResourcePool.Reference(),
		importOvaParams.Datastore.Reference(),
		cisp)

	if err != nil {
		return errors.Errorf("failed to create import spec: %s", err)
	}
	if spec.Error != nil {
		return errors.New(spec.Error[0].LocalizedMessage)
	}

	// The lease and upload cannot be used with a timeout
	// since we do not know how long it will take to upload
	// the ova to vSphere
	ctx = context.TODO()

	//Creates a new entity in this resource pool.
	//See VMware vCenter API documentation: Managed Object - ResourcePool - ImportVApp
	lease, err := importOvaParams.ResourcePool.ImportVApp(ctx,
		spec.ImportSpec,
		importOvaParams.Folder,
		importOvaParams.Host)

	if err != nil {
		return errors.Errorf("failed to import vapp: %s", err)
	}

	info, err := lease.Wait(ctx, spec.FileItem)
	if err != nil {
		return errors.Errorf("failed to lease wait: %s", err)
	}

	u := lease.StartUpdater(ctx, info)
	defer u.Done()

	for _, i := range info.Items {
		// upload the vmdk to which ever host that was first
		// available with the required network and datastore.
		err = upload(ctx, archive, lease, i)
		if err != nil {
			return errors.Errorf("failed to upload: %s", err)
		}
	}

	err = lease.Complete(ctx)
	if err != nil {
		return errors.Errorf("failed to lease complete: %s", err)
	}

	d.SetId(info.Entity.Value)

	err = attachTag(d, meta)
	if err != nil {
		return errors.Errorf("failed to attach tag to virtual machine: %s", err)
	}
	log.Printf("[DEBUG] %s: ova import complete", d.Get("name").(string))

	return resourceVSpherePrivateImportOvaRead(d, meta)
}

func resourceVSpherePrivateImportOvaRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*VSphereClient).vimClient.Client
	moRef := types.ManagedObjectReference{
		Value: d.Id(),
		Type:  "VirtualMachine",
	}

	vm := object.NewVirtualMachine(client, moRef)
	if vm == nil {
		return fmt.Errorf("error VirtualMachine not found, managed object id: %s", d.Id())
	}

	return nil
}

func resourceVSpherePrivateImportOvaDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning delete", d.Get("name").(string))
	ctx, cancel := context.WithTimeout(context.TODO(), defaultAPITimeout)
	defer cancel()

	client := meta.(*VSphereClient).vimClient.Client
	moRef := types.ManagedObjectReference{
		Value: d.Id(),
		Type:  "VirtualMachine",
	}

	vm := object.NewVirtualMachine(client, moRef)
	if vm == nil {
		return errors.Errorf("VirtualMachine not found")
	}

	task, err := vm.Destroy(ctx)
	if err != nil {
		return errors.Errorf("failed to destroy virtual machine %s", err)
	}

	err = task.Wait(ctx)
	if err != nil {
		return errors.Errorf("failed to destroy virtual machine %s", err)
	}

	d.SetId("")

	log.Printf("[DEBUG] %s: Delete complete", d.Get("name").(string))

	return nil
}
