package vsphere

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/customattribute"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/datastore"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/folder"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/viapi"
	"github.com/vmware/govmomi/vim25/types"
)

const (
	retryDeletePending   = "retryDeletePending"
	retryDeleteCompleted = "retryDeleteCompleted"
	retryDeleteError     = "retryDeleteError"

	waitForDeletePending   = "waitForDeletePending"
	waitForDeleteCompleted = "waitForDeleteCompleted"
	waitForDeleteError     = "waitForDeleteError"
)

// formatVmfsDatastoreCreateRollbackErrorUpdate defines the verbose error for extending a
// disk on creation where rollback is not possible.
const formatVmfsDatastoreCreateRollbackErrorUpdate = `
WARNING: Dangling resource!
There was an error extending your datastore with disk: %q:
%s
Additionally, there was an error removing the created datastore:
%s
You will need to remove this datastore manually before trying again.
`

func resourceVSphereVmfsDatastore() *schema.Resource {
	s := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the datastore.",
			Required:    true,
		},
		"host_system_id": {
			Type:        schema.TypeString,
			Description: "The managed object ID of the host to set up the datastore on.",
			ForceNew:    true,
			Required:    true,
		},
		"folder": {
			Type:          schema.TypeString,
			Description:   "The path to the datastore folder to put the datastore in.",
			Optional:      true,
			ConflictsWith: []string{"datastore_cluster_id"},
			StateFunc:     folder.NormalizePath,
		},
		"datastore_cluster_id": {
			Type:          schema.TypeString,
			Description:   "The managed object ID of the datastore cluster to place the datastore in.",
			Optional:      true,
			ConflictsWith: []string{"folder"},
		},
		"disks": {
			Type:        schema.TypeList,
			Description: "The disks to add to the datastore.",
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
	structure.MergeSchema(s, schemaDatastoreSummary())

	// Add tags schema
	s[vSphereTagAttributeKey] = tagsSchema()
	// Add custom attributes schema
	s[customattribute.ConfigKey] = customattribute.ConfigSchema()

	return &schema.Resource{
		Create:        resourceVSphereVmfsDatastoreCreate,
		Read:          resourceVSphereVmfsDatastoreRead,
		Update:        resourceVSphereVmfsDatastoreUpdate,
		Delete:        resourceVSphereVmfsDatastoreDelete,
		CustomizeDiff: resourceVSphereVmfsDatastoreCustomizeDiff,
		Importer: &schema.ResourceImporter{
			State: resourceVSphereVmfsDatastoreImport,
		},
		Schema: s,
	}
}

func resourceVSphereVmfsDatastoreCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).vimClient

	// Load up the tags client, which will validate a proper vCenter before
	// attempting to proceed if we have tags defined.
	tagsClient, err := tagsManagerIfDefined(d, meta)
	if err != nil {
		return err
	}
	// Verify a proper vCenter before proceeding if custom attributes are defined
	attrsProcessor, err := customattribute.GetDiffProcessorIfAttributesDefined(client, d)
	if err != nil {
		return err
	}

	hsID := d.Get("host_system_id").(string)
	dss, err := hostDatastoreSystemFromHostSystemID(client, hsID)
	if err != nil {
		return fmt.Errorf("error loading host datastore system: %s", err)
	}

	// To ensure the datastore is fully created with all the disks that we want
	// to add to it, first we add the initial disk, then we expand the disk with
	// the rest of the extents.
	disks := d.Get("disks").([]interface{})
	disk := disks[0].(string)
	spec, err := diskSpecForCreate(dss, disk)
	if err != nil {
		return err
	}
	spec.Vmfs.VolumeName = d.Get("name").(string)
	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	ds, err := dss.CreateVmfsDatastore(ctx, *spec)
	if err != nil {
		return fmt.Errorf("error creating datastore with disk %s: %s", disk, err)
	}

	// Add any remaining disks.
	for _, disk := range disks[1:] {
		var extendSpec *types.VmfsDatastoreExtendSpec
		extendSpec, err = diskSpecForExtend(dss, ds, disk.(string))
		if err != nil {
			// We have to destroy the created datastore here.
			if remErr := removeDatastore(dss, ds); remErr != nil {
				// We could not destroy the created datastore and there is now a dangling
				// resource. We need to instruct the user to remove the datastore
				// manually.
				return fmt.Errorf(formatVmfsDatastoreCreateRollbackErrorUpdate, disk, err, remErr)
			}
			return fmt.Errorf("error fetching datastore extend spec for disk %q: %s", disk, err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
		defer cancel()
		if _, err = extendVmfsDatastore(ctx, dss, ds, *extendSpec); err != nil {
			if remErr := removeDatastore(dss, ds); remErr != nil {
				// We could not destroy the created datastore and there is now a dangling
				// resource. We need to instruct the user to remove the datastore
				// manually.
				return fmt.Errorf(formatVmfsDatastoreCreateRollbackErrorUpdate, disk, err, remErr)
			}
			return fmt.Errorf("error extending datastore with disk %q: %s", disk, err)
		}
	}

	// Set the ID here now as most other issues here can be applied on an update,
	// so we don't need to roll back on failure.
	d.SetId(ds.Reference().Value)

	// Move the datastore to the correct folder first, if specified.
	f, err := resourceVSphereDatastoreApplyFolderOrStorageClusterPath(d, meta)
	if err != nil {
		return err
	}
	if !folder.PathIsEmpty(f) {
		if err := datastore.MoveToFolderRelativeHostSystemID(client, ds, hsID, f); err != nil {
			return fmt.Errorf("could not move datastore to folder %q: %s", f, err)
		}
	}

	// Apply any pending tags now
	if tagsClient != nil {
		if err := processTagDiff(tagsClient, d, ds); err != nil {
			return err
		}
	}

	// Set custom attributes
	if attrsProcessor != nil {
		if err := attrsProcessor.ProcessDiff(ds); err != nil {
			return err
		}
	}

	// Done
	return resourceVSphereVmfsDatastoreRead(d, meta)
}

func resourceVSphereVmfsDatastoreRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).vimClient
	id := d.Id()
	ds, err := datastore.FromID(client, id)
	if err != nil {
		return fmt.Errorf("cannot find datastore: %s", err)
	}
	props, err := datastore.Properties(ds)
	if err != nil {
		return fmt.Errorf("could not get properties for datastore: %s", err)
	}
	if err := flattenDatastoreSummary(d, &props.Summary); err != nil {
		return err
	}

	// Set the folder
	if err := resourceVSphereDatastoreReadFolderOrStorageClusterPath(d, ds); err != nil {
		return err
	}

	// We also need to update the disk list from the summary.
	var disks []string
	for _, disk := range props.Info.(*types.VmfsDatastoreInfo).Vmfs.Extent {
		disks = append(disks, disk.DiskName)
	}
	if err := d.Set("disks", disks); err != nil {
		return err
	}

	// Read tags if we have the ability to do so
	if tagsClient, _ := meta.(*Client).TagsManager(); tagsClient != nil {
		if err := readTagsForResource(tagsClient, ds, d); err != nil {
			return err
		}
	}

	// Read custom attributes
	if customattribute.IsSupported(client) {
		customattribute.ReadFromResource(props.Entity(), d)
	}

	return nil
}

func resourceVSphereVmfsDatastoreUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).vimClient

	// Load up the tags client, which will validate a proper vCenter before
	// attempting to proceed if we have tags defined.
	tagsClient, err := tagsManagerIfDefined(d, meta)
	if err != nil {
		return err
	}
	// Verify a proper vCenter before proceeding if custom attributes are defined
	attrsProcessor, err := customattribute.GetDiffProcessorIfAttributesDefined(client, d)
	if err != nil {
		return err
	}

	hsID := d.Get("host_system_id").(string)
	dss, err := hostDatastoreSystemFromHostSystemID(client, hsID)
	if err != nil {
		return fmt.Errorf("error loading host datastore system: %s", err)
	}

	id := d.Id()
	ds, err := datastore.FromID(client, id)
	if err != nil {
		return fmt.Errorf("cannot find datastore: %s", err)
	}

	// Rename this datastore if our name has drifted.
	if d.HasChange("name") {
		if err := viapi.RenameObject(client, ds.Reference(), d.Get("name").(string)); err != nil {
			return err
		}
	}

	// Update folder or datastore cluster if necessary
	if d.HasChange("folder") || d.HasChange("datastore_cluster_id") {
		f, err := resourceVSphereDatastoreApplyFolderOrStorageClusterPath(d, meta)
		if err != nil {
			return err
		}
		if err := datastore.MoveToFolder(client, ds, f); err != nil {
			return fmt.Errorf("could not move datastore to folder %q: %s", f, err)
		}
	}

	// Apply any pending tags now
	if tagsClient != nil {
		if err := processTagDiff(tagsClient, d, ds); err != nil {
			return err
		}
	}

	// Apply custom attribute updates
	if attrsProcessor != nil {
		if err := attrsProcessor.ProcessDiff(ds); err != nil {
			return err
		}
	}

	// Veto this update if it means a disk was removed. Shrinking
	// datastores/removing extents is not supported.
	old, newValue := d.GetChange("disks")
	for _, v1 := range old.([]interface{}) {
		var found bool
		for _, v2 := range newValue.([]interface{}) {
			if v1.(string) == v2.(string) {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("disk %s found in state but not config (removal of disks is not supported)", v1)
		}
	}

	// Now we basically reverse what we did above when we were checking for
	// removed disks, and add any new disks that have been added.
	for _, v1 := range newValue.([]interface{}) {
		var found bool
		for _, v2 := range old.([]interface{}) {
			if v1.(string) == v2.(string) {
				found = true
			}
		}
		if !found {
			// Add the disk
			spec, err := diskSpecForExtend(dss, ds, v1.(string))
			if err != nil {
				return err
			}
			ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
			defer cancel()
			if _, err := extendVmfsDatastore(ctx, dss, ds, *spec); err != nil {
				return err
			}
		}
	}

	// Should be done with the update here.
	return resourceVSphereVmfsDatastoreRead(d, meta)
}

func resourceVSphereVmfsDatastoreDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).vimClient
	hsID := d.Get("host_system_id").(string)
	dss, err := hostDatastoreSystemFromHostSystemID(client, hsID)
	if err != nil {
		return fmt.Errorf("error loading host datastore system: %s", err)
	}

	id := d.Id()
	ds, err := datastore.FromID(client, id)
	if err != nil {
		return fmt.Errorf("cannot find datastore: %s", err)
	}

	// This is a race that more than likely will only come up during tests, but
	// we still want to guard against it - when working with datastores that end
	// up mounting across multiple hosts, removing the datastore will fail if
	// it's removed too quickly (like right away, for example). So we set up a
	// very short retry waiter to make sure if the first attempt fails, the
	// second one should probably succeed right away. We also insert a small
	// minimum delay to make an honest first attempt at trying to delete the
	// datastore without spamming the task log with errors.
	deleteRetryFunc := func() (interface{}, string, error) {
		err := removeDatastore(dss, ds)
		if err != nil {
			if viapi.IsResourceInUseError(err) {
				// Pending
				return struct{}{}, retryDeletePending, nil
			}
			// Some other error
			return struct{}{}, retryDeleteError, err
		}
		// Done
		return struct{}{}, retryDeleteCompleted, nil
	}

	deleteRetry := &resource.StateChangeConf{
		Pending:    []string{retryDeletePending},
		Target:     []string{retryDeleteCompleted},
		Refresh:    deleteRetryFunc,
		Timeout:    30 * time.Second,
		MinTimeout: 2 * time.Second,
		Delay:      2 * time.Second,
	}

	_, err = deleteRetry.WaitForState()
	if err != nil {
		return fmt.Errorf("could not delete datastore: %s", err)
	}

	// We need to make sure the datastore is completely removed. There appears to
	// be a bit of a delay sometimes on vCenter, and it causes issues in tests,
	// which means it could cause issues somewhere else too.
	waitForDeleteFunc := func() (interface{}, string, error) {
		_, err := datastore.FromID(client, id)
		if err != nil {
			if viapi.IsManagedObjectNotFoundError(err) {
				// Done
				return struct{}{}, waitForDeleteCompleted, nil
			}
			// Some other error
			return struct{}{}, waitForDeleteError, err
		}
		return struct{}{}, waitForDeletePending, nil
	}

	waitForDelete := &resource.StateChangeConf{
		Pending:        []string{waitForDeletePending},
		Target:         []string{waitForDeleteCompleted},
		Refresh:        waitForDeleteFunc,
		Timeout:        defaultAPITimeout,
		MinTimeout:     2 * time.Second,
		Delay:          1 * time.Second,
		NotFoundChecks: 35,
	}

	_, err = waitForDelete.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for datastore to delete: %s", err.Error())
	}

	return nil
}

func resourceVSphereVmfsDatastoreCustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	// Check all disks and make sure that the entries are not nil, empty, or duplicates.
	disks := make(map[string]struct{})
	for i, v := range d.Get("disks").([]interface{}) {
		if v == nil || v.(string) == "" {
			return fmt.Errorf("disk.%d: empty entry", i)
		}
		if _, ok := disks[v.(string)]; ok {
			return fmt.Errorf("disk.%d: duplicate name %q", i, v.(string))
		}
		disks[v.(string)] = struct{}{}
	}
	return nil
}

func resourceVSphereVmfsDatastoreImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// We support importing a MoRef - so we need to load the datastore and check
	// to make sure 1) it exists, and 2) it's a VMFS datastore. If it is, we are
	// good to go (rest of the stuff will be handled by read on refresh).
	ids := strings.SplitN(d.Id(), ":", 2)
	if len(ids) != 2 {
		return nil, errors.New("please supply the ID in the following format: DATASTOREID:HOSTID")
	}

	id := ids[0]
	hsID := ids[1]
	client := meta.(*Client).vimClient
	ds, err := datastore.FromID(client, id)
	if err != nil {
		return nil, fmt.Errorf("cannot find datastore: %s", err)
	}
	props, err := datastore.Properties(ds)
	if err != nil {
		return nil, fmt.Errorf("could not get properties for datastore: %s", err)
	}

	t := types.HostFileSystemVolumeFileSystemType(props.Summary.Type)
	if t != types.HostFileSystemVolumeFileSystemTypeVMFS {
		return nil, fmt.Errorf("datastore ID %q is not a VMFS datastore", id)
	}

	var found bool
	for _, mount := range props.Host {
		if mount.Key.Value == hsID {
			found = true
		}
	}
	if !found {
		return nil, fmt.Errorf("configured host_system_id %q not found as a mounted host on datastore", hsID)
	}
	d.SetId(id)
	_ = d.Set("host_system_id", hsID)

	return []*schema.ResourceData{d}, nil
}
