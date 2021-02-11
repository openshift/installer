package vsphere

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/datastore"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/folder"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/storagepod"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
)

// schemaDatastoreSummary returns schema items for resources that
// need to work with a DatastoreSummary.
func schemaDatastoreSummary() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Note that the following fields are not represented in the schema here:
		// * Name (more than likely the ID attribute and will be represented in
		// resource schema)
		// * Type (redundant attribute as the datastore type will be represented by
		// the resource)
		"accessible": {
			Type:        schema.TypeBool,
			Description: "The connectivity status of the datastore. If this is false, some other computed attributes may be out of date.",
			Computed:    true,
		},
		"capacity": {
			Type:        schema.TypeInt,
			Description: "Maximum capacity of the datastore, in MB.",
			Computed:    true,
		},
		"free_space": {
			Type:        schema.TypeInt,
			Description: "Available space of this datastore, in MB.",
			Computed:    true,
		},
		"maintenance_mode": {
			Type:        schema.TypeString,
			Description: "The current maintenance mode state of the datastore.",
			Computed:    true,
		},
		"multiple_host_access": {
			Type:        schema.TypeBool,
			Description: "If true, more than one host in the datacenter has been configured with access to the datastore.",
			Computed:    true,
		},
		"uncommitted_space": {
			Type:        schema.TypeInt,
			Description: "Total additional storage space, in MB, potentially used by all virtual machines on this datastore.",
			Computed:    true,
		},
		"url": {
			Type:        schema.TypeString,
			Description: "The unique locator for the datastore.",
			Computed:    true,
		},
	}
}

// flattenDatastoreSummary reads various fields from a DatastoreSummary into
// the passed in ResourceData.
func flattenDatastoreSummary(d *schema.ResourceData, obj *types.DatastoreSummary) error {
	d.Set("accessible", obj.Accessible)
	d.Set("capacity", structure.ByteToMB(obj.Capacity))
	d.Set("free_space", structure.ByteToMB(obj.FreeSpace))
	d.Set("maintenance_mode", obj.MaintenanceMode)
	d.Set("multiple_host_access", obj.MultipleHostAccess)
	d.Set("uncommitted_space", structure.ByteToMB(obj.Uncommitted))
	d.Set("url", obj.Url)

	// Set the name attribute off of the name here - since we do not track this
	// here we check for errors
	if err := d.Set("name", obj.Name); err != nil {
		return err
	}
	return nil
}

// resourceVSphereDatastoreApplyFolderOrStorageClusterPath returns a path to a
// folder or a datastore cluster, depending on what has been selected in the
// resource.
func resourceVSphereDatastoreApplyFolderOrStorageClusterPath(d *schema.ResourceData, meta interface{}) (string, error) {
	var path string
	fvalue, fok := d.GetOk("folder")
	cvalue, cok := d.GetOk("datastore_cluster_id")
	switch {
	case fok:
		path = fvalue.(string)
	case cok:
		return resourceVSphereDatastoreStorageClusterPathNormalized(meta, cvalue.(string))
	}
	return path, nil
}

func resourceVSphereDatastoreStorageClusterPathNormalized(meta interface{}, id string) (string, error) {
	client := meta.(*VSphereClient).vimClient
	pod, err := storagepod.FromID(client, id)
	if err != nil {
		return "", err
	}
	return folder.RootPathParticleDatastore.SplitRelative(pod.InventoryPath)
}

// resourceVSphereDatastoreReadFolderOrStorageClusterPath checks the inventory
// path of the supplied datastore and checks to see if it is a normal folder or
// if it's a datastore cluster, and saves the attributes accordingly.
func resourceVSphereDatastoreReadFolderOrStorageClusterPath(d *schema.ResourceData, ds *object.Datastore) error {
	props, err := datastore.Properties(ds)
	if err != nil {
		return fmt.Errorf("error fetching datastore properties while parsing path: %s", err)
	}
	switch props.Parent.Type {
	case "Folder":
		return resourceVSphereDatastoreReadFolderOrStorageClusterPathAsFolder(d, ds)
	case "StoragePod":
		return resourceVSphereDatastoreReadFolderOrStorageClusterPathSetAttributes(d, "", props.Parent.Value)
	}
	return fmt.Errorf("unknown datastore parent type %q while parsing inventory path", props.Parent.Type)
}

func resourceVSphereDatastoreReadFolderOrStorageClusterPathAsFolder(d *schema.ResourceData, ds *object.Datastore) error {
	f, err := folder.RootPathParticleDatastore.SplitRelativeFolder(ds.InventoryPath)
	if err != nil {
		return fmt.Errorf("error parsing datastore path %q: %s", ds.InventoryPath, err)
	}
	return resourceVSphereDatastoreReadFolderOrStorageClusterPathSetAttributes(d, folder.NormalizePath(f), "")
}

func resourceVSphereDatastoreReadFolderOrStorageClusterPathSetAttributes(d *schema.ResourceData, f, c string) error {
	if err := d.Set("folder", f); err != nil {
		return fmt.Errorf("error setting folder attribute: %s", err)
	}
	if err := d.Set("datastore_cluster_id", c); err != nil {
		return fmt.Errorf("error setting datastore_cluster_id attribute: %s", err)
	}
	return nil
}
