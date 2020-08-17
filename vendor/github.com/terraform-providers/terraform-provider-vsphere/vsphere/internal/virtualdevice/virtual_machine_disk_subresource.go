package virtualdevice

import (
	"errors"
	"fmt"
	"log"
	"math"
	"path"
	"reflect"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/mitchellh/copystructure"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/datastore"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/spbm"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/storagepod"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/viapi"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/virtualmachine"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
)

// diskNameDeprecationNotice is the deprecation warning for the "name"
// attribute, which will removed in 2.0. The notice is verbose, so we format it
// so it looks a little better over CLI.
//
// TODO: Remove this in 2.0.
const diskNameDeprecationNotice = `
The name attribute for virtual disks will be removed in favor of "label" in
future releases. To transition existing disks, rename the "name" attribute to
"label". When doing so, ensure the value of the attribute stays the same.

Note that "label" does not control the name of a VMDK and does not need to bear
the name of one on new disks or virtual machines. For more information, see the
documentation for the label attribute at: 

https://www.terraform.io/docs/providers/vsphere/r/virtual_machine.html#label
`

// diskDatastoreComputedName is a friendly display for disks with datastores
// marked as computed. This happens in datastore cluster workflows.
const diskDatastoreComputedName = "<computed>"

// diskDeletedName is a placeholder name for deleted disks. This is to assist
// with user-friendliness in the diff.
const diskDeletedName = "<deleted>"

// diskDetachedName is a placeholder name for disks that are getting detached,
// either because they have keep_on_remove set or are external disks attached
// with "attach".
const diskDetachedName = "<remove, keep disk>"

// diskOrphanedPrefix is a placeholder name for disks that have been discovered
// as not being tracked by Terraform. These disks are assigned
// "orphaned_disk_0", "orphaned_disk_1", and so on.
const diskOrphanedPrefix = "orphaned_disk_"

var diskSubresourceModeAllowedValues = []string{
	string(types.VirtualDiskModePersistent),
	string(types.VirtualDiskModeNonpersistent),
	string(types.VirtualDiskModeUndoable),
	string(types.VirtualDiskModeIndependent_persistent),
	string(types.VirtualDiskModeIndependent_nonpersistent),
	string(types.VirtualDiskModeAppend),
}

var diskSubresourceSharingAllowedValues = []string{
	string(types.VirtualDiskSharingSharingNone),
	string(types.VirtualDiskSharingSharingMultiWriter),
}

// DiskSubresourceSchema represents the schema for the disk sub-resource.
func DiskSubresourceSchema() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		// VirtualDiskFlatVer2BackingInfo
		"datastore_id": {
			Type:          schema.TypeString,
			Optional:      true,
			ConflictsWith: []string{"datastore_cluster_id"},
			Description:   "The datastore ID for this virtual disk, if different than the virtual machine.",
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The file name of the disk. This can be either a name or path relative to the root of the datastore. If simply a name, the disk is located with the virtual machine.",
			ValidateFunc: func(v interface{}, _ string) ([]string, []error) {
				if path.Ext(v.(string)) != ".vmdk" {
					return nil, []error{fmt.Errorf("disk name %s must end in .vmdk", v.(string))}
				}
				return nil, nil
			},
			Deprecated: diskNameDeprecationNotice,
		},
		"path": {
			Type:          schema.TypeString,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"datastore_cluster_id"},
			Description:   "The full path of the virtual disk. This can only be provided if attach is set to true, otherwise it is a read-only value.",
			ValidateFunc: func(v interface{}, _ string) ([]string, []error) {
				if path.Ext(v.(string)) != ".vmdk" {
					return nil, []error{fmt.Errorf("disk path %s must end in .vmdk", v.(string))}
				}
				return nil, nil
			},
		},
		"disk_mode": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      string(types.VirtualDiskModePersistent),
			Description:  "The mode of this this virtual disk for purposes of writes and snapshotting. Can be one of append, independent_nonpersistent, independent_persistent, nonpersistent, persistent, or undoable.",
			ValidateFunc: validation.StringInSlice(diskSubresourceModeAllowedValues, false),
		},
		"eagerly_scrub": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "The virtual disk file zeroing policy when thin_provision is not true. The default is false, which lazily-zeros the disk, speeding up thick-provisioned disk creation time.",
		},
		"disk_sharing": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      string(types.VirtualDiskSharingSharingNone),
			Description:  "The sharing mode of this virtual disk. Can be one of sharingMultiWriter or sharingNone.",
			ValidateFunc: validation.StringInSlice(diskSubresourceSharingAllowedValues, false),
		},
		"thin_provisioned": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "If true, this disk is thin provisioned, with space for the file being allocated on an as-needed basis.",
		},
		"write_through": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "If true, writes for this disk are sent directly to the filesystem immediately instead of being buffered.",
		},
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The UUID of the virtual disk.",
		},

		// StorageIOAllocationInfo
		"io_limit": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      -1,
			Description:  "The upper limit of IOPS that this disk can use.",
			ValidateFunc: validation.IntAtLeast(-1),
		},
		"io_reservation": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      0,
			Description:  "The I/O guarantee that this disk has, in IOPS.",
			ValidateFunc: validation.IntAtLeast(0),
		},
		"io_share_level": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      string(types.SharesLevelNormal),
			Description:  "The share allocation level for this disk. Can be one of low, normal, high, or custom.",
			ValidateFunc: validation.StringInSlice(sharesLevelAllowedValues, false),
		},
		"io_share_count": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      0,
			Description:  "The share count for this disk when the share level is custom.",
			ValidateFunc: validation.IntAtLeast(0),
		},

		// VirtualDisk
		"size": {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "The size of the disk, in GB.",
			ValidateFunc: validation.IntAtLeast(1),
		},

		// Complex terraform-local things
		"label": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A unique label for this disk.",
			ValidateFunc: func(v interface{}, _ string) ([]string, []error) {
				if strings.HasPrefix(v.(string), diskOrphanedPrefix) {
					return nil, []error{fmt.Errorf("disk label %q cannot start with %q", v.(string), diskOrphanedPrefix)}
				}
				return nil, nil
			},
		},
		"unit_number": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      0,
			Description:  "The unique device number for this disk. This number determines where on the SCSI bus this device will be attached.",
			ValidateFunc: validation.IntBetween(0, 59),
		},
		"keep_on_remove": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Set to true to keep the underlying VMDK file when removing this virtual disk from configuration.",
		},
		"attach": {
			Type:          schema.TypeBool,
			Optional:      true,
			Default:       false,
			ConflictsWith: []string{"datastore_cluster_id"},
			Description:   "If this is true, the disk is attached instead of created. Implies keep_on_remove.",
		},
		"storage_policy_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The ID of the storage policy to assign to the virtual disk in VM.",
		},
	}
	structure.MergeSchema(s, subresourceSchema())
	return s
}

// DiskSubresource represents a vsphere_virtual_machine disk sub-resource, with
// a complex device lifecycle.
type DiskSubresource struct {
	*Subresource

	// The set hash for the device as it exists when NewDiskSubresource is
	// called.
	ID int
}

// NewDiskSubresource returns a subresource populated with all of the necessary
// fields.
func NewDiskSubresource(client *govmomi.Client, rdd resourceDataDiff, d, old map[string]interface{}, idx int) *DiskSubresource {
	sr := &DiskSubresource{
		Subresource: &Subresource{
			schema:  DiskSubresourceSchema(),
			client:  client,
			srtype:  subresourceTypeDisk,
			data:    d,
			olddata: old,
			rdd:     rdd,
		},
	}
	sr.Index = idx
	return sr
}

// DiskApplyOperation processes an apply operation for all disks in the
// resource.
//
// The function takes the root resource's ResourceData, the provider
// connection, and the device list as known to vSphere at the start of this
// operation. All disk operations are carried out, with both the complete,
// updated, VirtualDeviceList, and the complete list of changes returned as a
// slice of BaseVirtualDeviceConfigSpec.
func DiskApplyOperation(d *schema.ResourceData, c *govmomi.Client, l object.VirtualDeviceList) (object.VirtualDeviceList, []types.BaseVirtualDeviceConfigSpec, error) {
	log.Printf("[DEBUG] DiskApplyOperation: Beginning apply operation")
	o, n := d.GetChange(subresourceTypeDisk)
	ods := o.([]interface{})
	nds := n.([]interface{})

	var spec []types.BaseVirtualDeviceConfigSpec

	// Our old and new sets now have an accurate description of devices that may
	// have been added, removed, or changed. Look for removed devices first.
	log.Printf("[DEBUG] DiskApplyOperation: Looking for resources to delete")
	for oi, oe := range ods {
		om := oe.(map[string]interface{})
		if err := diskApplyOperationDelete(oi, om, nds, c, d, &l, &spec); err != nil {
			return nil, nil, err
		}
	}

	// Now check for creates and updates.  The results of this operation are
	// committed to state after the operation completes, on top of the items that
	// have not changed.
	var updates []interface{}
	log.Printf("[DEBUG] DiskApplyOperation: Looking for resources to create or update")
	log.Printf("[DEBUG] DiskApplyOperation: Resources not being changed: %s", subresourceListString(updates))
	for ni, ne := range nds {
		nm := ne.(map[string]interface{})
		if err := diskApplyOperationCreateUpdate(ni, nm, ods, c, d, &l, &spec, &updates); err != nil {
			return nil, nil, err
		}
	}

	log.Printf("[DEBUG] DiskApplyOperation: Post-apply final resource list: %s", subresourceListString(updates))
	// We are now done! Return the updated device list and config spec. Save updates as well.
	if err := d.Set(subresourceTypeDisk, updates); err != nil {
		return nil, nil, err
	}
	log.Printf("[DEBUG] DiskApplyOperation: Device list at end of operation: %s", DeviceListString(l))
	log.Printf("[DEBUG] DiskApplyOperation: Device config operations from apply: %s", DeviceChangeString(spec))
	log.Printf("[DEBUG] DiskApplyOperation: Apply complete, returning updated spec")
	return l, spec, nil
}

// diskApplyOperationDelete is an inner-loop helper for disk deletion
// operations.
func diskApplyOperationDelete(
	index int,
	oldData map[string]interface{},
	newDataSet []interface{},
	c *govmomi.Client,
	d *schema.ResourceData,
	l *object.VirtualDeviceList,
	spec *[]types.BaseVirtualDeviceConfigSpec,
) error {
	didx := -1
	for ni, ne := range newDataSet {
		newData := ne.(map[string]interface{})
		var name string
		var err error
		if name, err = diskLabelOrName(newData); err != nil {
			return err
		}
		if (name == diskDeletedName || name == diskDetachedName) && oldData["uuid"] == newData["uuid"] {
			didx = ni
			break
		}
	}
	if didx < 0 {
		// Deleted entry not found
		return nil
	}
	r := NewDiskSubresource(c, d, oldData, nil, index)
	dspec, err := r.Delete(*l)
	if err != nil {
		return fmt.Errorf("%s: %s", r.Addr(), err)
	}
	*l = applyDeviceChange(*l, dspec)
	*spec = append(*spec, dspec...)
	return nil
}

// diskApplyOperationCreateUpdate is an inner-loop helper for disk creation and
// update operations.
func diskApplyOperationCreateUpdate(
	index int,
	newData map[string]interface{},
	oldDataSet []interface{},
	c *govmomi.Client,
	d *schema.ResourceData,
	l *object.VirtualDeviceList,
	spec *[]types.BaseVirtualDeviceConfigSpec,
	updates *[]interface{},
) error {
	var name string
	var err error
	if name, err = diskLabelOrName(newData); err != nil {
		return err
	}
	if name == diskDeletedName || name == diskDetachedName {
		// This is a "dummy" deleted resource and should be skipped over
		return nil
	}
	for _, oe := range oldDataSet {
		oldData := oe.(map[string]interface{})
		if newData["uuid"] == oldData["uuid"] {
			// This is an update
			r := NewDiskSubresource(c, d, newData, oldData, index)
			// If the only thing changing here is the datastore, or keep_on_remove,
			// this is a no-op as far as a device change is concerned. Datastore
			// changes are handled during storage vMotion later on during the
			// update phase. keep_on_remove is a Terraform-only attribute and only
			// needs to be committed to state.
			omc, err := copystructure.Copy(oldData)
			if err != nil {
				return fmt.Errorf("%s: error generating copy of old disk data: %s", r.Addr(), err)
			}
			oldCopy := omc.(map[string]interface{})
			oldCopy["datastore_id"] = newData["datastore_id"]
			oldCopy["keep_on_remove"] = newData["keep_on_remove"]
			// TODO: Remove these in 2.0, when all attributes should bear a label and
			// name is gone, and we won't need to exempt transitions.
			oldCopy["label"] = newData["label"]
			oldCopy["name"] = newData["name"]
			if reflect.DeepEqual(oldCopy, newData) {
				*updates = append(*updates, r.Data())
				return nil
			}
			uspec, err := r.Update(*l)
			if err != nil {
				return fmt.Errorf("%s: %s", r.Addr(), err)
			}
			*l = applyDeviceChange(*l, uspec)
			*spec = append(*spec, uspec...)
			*updates = append(*updates, r.Data())
			return nil
		}
	}
	// New data was not found - this is a create operation
	r := NewDiskSubresource(c, d, newData, nil, index)
	cspec, err := r.Create(*l)
	if err != nil {
		return fmt.Errorf("%s: %s", r.Addr(), err)
	}
	*l = applyDeviceChange(*l, cspec)
	*spec = append(*spec, cspec...)
	*updates = append(*updates, r.Data())
	return nil
}

// DiskRefreshOperation processes a refresh operation for all of the disks in
// the resource.
//
// This functions similar to DiskApplyOperation, but nothing to change is
// returned, all necessary values are just set and committed to state.
func DiskRefreshOperation(d *schema.ResourceData, c *govmomi.Client, l object.VirtualDeviceList) error {
	log.Printf("[DEBUG] DiskRefreshOperation: Beginning refresh")
	devices := SelectDisks(l, d.Get("scsi_controller_count").(int))
	log.Printf("[DEBUG] DiskRefreshOperation: Disk devices located: %s", DeviceListString(devices))
	curSet := d.Get(subresourceTypeDisk).([]interface{})
	log.Printf("[DEBUG] DiskRefreshOperation: Current resource set from state: %s", subresourceListString(curSet))
	var newSet []interface{}
	// First check for negative keys. These are freshly added devices that are
	// usually coming into read post-create.
	//
	// If we find what we are looking for, we remove the device from the working
	// set so that we don't try and process it in the next few passes.
	log.Printf("[DEBUG] DiskRefreshOperation: Looking for freshly-created or re-assigned resources to read in")
	for i, item := range curSet {
		m := item.(map[string]interface{})
		if m["key"].(int) < 1 {
			r := NewDiskSubresource(c, d, m, nil, i)
			if err := r.Read(l); err != nil {
				return fmt.Errorf("%s: %s", r.Addr(), err)
			}
			if r.Get("key").(int) < 1 {
				// This should not have happened - if it did, our device
				// creation/update logic failed somehow that we were not able to track.
				return fmt.Errorf("device %d with address %s still unaccounted for after update/read", r.Get("key").(int), r.Get("device_address").(string))
			}
			newSet = append(newSet, r.Data())
			for i := 0; i < len(devices); i++ {
				device := devices[i]
				if device.GetVirtualDevice().Key == int32(r.Get("key").(int)) {
					devices = append(devices[:i], devices[i+1:]...)
					i--
				}
			}
		}
	}
	log.Printf("[DEBUG] DiskRefreshOperation: Disk devices after created/re-assigned device search: %s", DeviceListString(devices))
	log.Printf("[DEBUG] DiskRefreshOperation: Resource set to write after created/re-assigned device search: %s", subresourceListString(newSet))

	// Go over the remaining devices, refresh via key, and then remove their
	// entries as well.
	log.Printf("[DEBUG] DiskRefreshOperation: Looking for devices known in state")
	for i := 0; i < len(devices); i++ {
		device := devices[i]
		for n, item := range curSet {
			m := item.(map[string]interface{})
			if m["key"].(int) < 1 {
				// Skip any of these keys as we won't be matching any of those anyway here
				continue
			}

			if !diskUUIDMatch(device, m["uuid"].(string)) {
				// Skip any device that doesn't match UUID
				continue
			}
			// We should have our device -> resource match, so read now.
			r := NewDiskSubresource(c, d, m, nil, n)
			if err := r.Read(l); err != nil {
				return fmt.Errorf("%s: %s", r.Addr(), err)
			}

			if strings.HasPrefix(r.Get("label").(string), diskOrphanedPrefix) {
				// Skip if it's previously discovered orphaned device
				continue
			}

			// Done reading, push this onto our new set and remove the device from
			// the list
			newSet = append(newSet, r.Data())
			devices = append(devices[:i], devices[i+1:]...)
			i--
		}
	}
	log.Printf("[DEBUG] DiskRefreshOperation: Resource set to write after known device search: %s", subresourceListString(newSet))
	log.Printf("[DEBUG] DiskRefreshOperation: Probable orphaned disk devices: %s", DeviceListString(devices))

	// Finally, any device that is still here is orphaned. They should be added
	// as new devices.
	log.Printf("[DEBUG] DiskRefreshOperation: Adding orphaned devices")
	for i, device := range devices {
		m := make(map[string]interface{})
		vd := device.GetVirtualDevice()
		ctlr := l.FindByKey(vd.ControllerKey)
		if ctlr == nil {
			return fmt.Errorf("could not find controller with key %d", vd.Key)
		}
		m["key"] = int(vd.Key)
		var err error
		m["device_address"], err = computeDevAddr(vd, ctlr.(types.BaseVirtualController))
		if err != nil {
			return fmt.Errorf("error computing device address: %s", err)
		}
		// We want to set keep_on_remove for these disks as well so that they are
		// not destroyed when we remove them in the next TF run.
		m["keep_on_remove"] = true
		r := NewDiskSubresource(c, d, m, nil, len(newSet))
		if err := r.Read(l); err != nil {
			return fmt.Errorf("%s: %s", r.Addr(), err)
		}
		// Add a generic label indicating that this disk is orphaned.
		r.Set("label", fmt.Sprintf("%s%d", diskOrphanedPrefix, i))
		newSet = append(newSet, r.Data())
	}

	log.Printf("[DEBUG] DiskRefreshOperation: Resource set to write after adding orphaned devices: %s", subresourceListString(newSet))
	// Sort the device list by unit number. This provides some semblance of order
	// in the state as devices are added and removed.
	sort.Sort(virtualDiskSubresourceSorter(newSet))
	log.Printf("[DEBUG] DiskRefreshOperation: Final (sorted) resource set to write: %s", subresourceListString(newSet))
	log.Printf("[DEBUG] DiskRefreshOperation: Refresh operation complete, sending new resource set")
	return d.Set(subresourceTypeDisk, newSet)
}

// DiskDestroyOperation process the destroy operation for virtual disks.
//
// Disks are the only real operation that require special destroy logic, and
// that's because we want to check to make sure that we detach any disks that
// need to be simply detached (not deleted) before we destroy the entire
// virtual machine, as that would take those disks with it.
func DiskDestroyOperation(d *schema.ResourceData, c *govmomi.Client, l object.VirtualDeviceList) ([]types.BaseVirtualDeviceConfigSpec, error) {
	log.Printf("[DEBUG] DiskDestroyOperation: Beginning destroy")
	// All we are doing here is getting a config spec for detaching the disks
	// that we need to detach, so we don't need the vast majority of the stateful
	// logic that is in deviceApplyOperation.
	ds := d.Get(subresourceTypeDisk).([]interface{})

	var spec []types.BaseVirtualDeviceConfigSpec

	log.Printf("[DEBUG] DiskDestroyOperation: Detaching devices with keep_on_remove enabled")
	for oi, oe := range ds {
		m := oe.(map[string]interface{})
		if !m["keep_on_remove"].(bool) && !m["attach"].(bool) {
			// We don't care about disks we haven't set to keep
			continue
		}
		r := NewDiskSubresource(c, d, m, nil, oi)
		dspec, err := r.Delete(l)
		if err != nil {
			return nil, fmt.Errorf("%s: %s", r.Addr(), err)
		}
		l = applyDeviceChange(l, dspec)
		spec = append(spec, dspec...)
	}

	log.Printf("[DEBUG] DiskDestroyOperation: Device config operations from destroy: %s", DeviceChangeString(spec))
	return spec, nil
}

// DiskDiffOperation performs operations relevant to managing the diff on disk
// sub-resources.
//
// Most importantly, this works to prevent spurious diffs by extrapolating the
// correct correlation between the old and new sets using the name as a primary
// key, and then normalizing the two diffs so that computed data is properly
// set.
//
// The following validation operations are also carried out on the set as a
// whole:
//
// * Ensuring all names are unique across the set.
// * Ensuring that at least one element in the set has a unit_number of 0.
func DiskDiffOperation(d *schema.ResourceDiff, c *govmomi.Client) error {
	log.Printf("[DEBUG] DiskDiffOperation: Beginning disk diff customization")
	o, n := d.GetChange(subresourceTypeDisk)
	// Some global validation first. We handle individual validation later.
	log.Printf("[DEBUG] DiskDiffOperation: Beginning collective diff validation (indexes aligned to new config)")
	names := make(map[string]struct{})
	attachments := make(map[string]struct{})
	units := make(map[int]struct{})
	if len(n.([]interface{})) < 1 {
		return errors.New("there must be at least one disk specified")
	}
	for ni, ne := range n.([]interface{}) {
		nm := ne.(map[string]interface{})
		name, err := diskLabelOrName(nm)
		if err != nil {
			return fmt.Errorf("disk.%d: %s", ni, err)
		}
		if _, ok := names[name]; ok {
			return fmt.Errorf("disk: duplicate name %s", name)
		}
		// If attach is set, we need to validate that there's no other duplicate paths.
		curDiskPath := fmt.Sprintf("disk.%d.path", ni)
		pathKnown := d.NewValueKnown(curDiskPath)
		if nm["attach"].(bool) {
			path := diskPathOrName(nm)
			if pathKnown {
				if path == "" {
					return fmt.Errorf("disk.%d: path or name cannot be empty when using attach", ni)
				}
				if _, ok := attachments[path]; ok {
					return fmt.Errorf("disk: multiple entries trying to attach external disk %s", path)
				}
				attachments[path] = struct{}{}
			} else {
				log.Printf("[DEBUG] Disk path for disk %d is not known yet.", ni)
			}
		}

		if _, ok := units[nm["unit_number"].(int)]; ok {
			return fmt.Errorf("disk: duplicate unit_number %d", nm["unit_number"].(int))
		}
		names[name] = struct{}{}
		units[nm["unit_number"].(int)] = struct{}{}
		r := NewDiskSubresource(c, d, nm, nil, ni)
		if err := r.DiffGeneral(); err != nil {
			return fmt.Errorf("%s: %s", r.Addr(), err)
		}
	}
	if _, ok := units[0]; !ok {
		return errors.New("at least one disk must have a unit_number of 0")
	}

	// Perform the normalization here.
	log.Printf("[DEBUG] DiskDiffOperation: Beginning diff validation and normalization (indexes aligned to old state)")
	ods := o.([]interface{})
	nds := n.([]interface{})

	normalized := make([]interface{}, len(ods))
nextNew:
	for _, ne := range nds {
		nm := ne.(map[string]interface{})
		for oi, oe := range ods {
			om := oe.(map[string]interface{})
			var oname, nname string
			var err error
			if oname, err = diskLabelOrName(om); err != nil {
				return fmt.Errorf("disk.%d: %s", oi, err)
			}
			if nname, err = diskLabelOrName(nm); err != nil {
				return fmt.Errorf("disk.%d: %s", oi, err)
			}
			// We extrapolate using the label as a "primary key" of sorts.
			if nname == oname {
				r := NewDiskSubresource(c, d, nm, om, oi)
				if err := r.DiffExisting(); err != nil {
					return fmt.Errorf("%s: %s", r.Addr(), err)
				}
				normalized[oi] = r.Data()
				continue nextNew
			}
		}
		// We didn't find a match for this resource, it could be a new resource or
		// significantly altered. Put it back on the list in the same form we got
		// it in, but all computed data first, just in case it was in a position
		// previously occupied by an existing resource.
		nm["key"] = 0
		nm["device_address"] = ""
		nm["uuid"] = ""
		if a, ok := nm["attach"]; !ok || !a.(bool) {
			nm["path"] = ""
		} else {
			_, ok := nm["size"]
			if ok {
				delete(nm, "size")
			}
		}
		if dsID, ok := nm["datastore_id"]; !ok || dsID == "" {
			nm["datastore_id"] = diskDatastoreComputedName
		}
		normalized = append(normalized, nm)
	}

	// Go thru the new list, and replace any nils with a "deleted" copy of the
	// old resource. This is basically a copy of the old entry with <deleted> in
	// place of the label, so it shows up nicely in the diff.
	for ni, ne := range normalized {
		if ne != nil {
			continue
		}
		nv, err := copystructure.Copy(ods[ni])
		if err != nil {
			return fmt.Errorf("disk.%d: error making updated diff of deleted entry: %s", ni, err)
		}
		nm := nv.(map[string]interface{})
		// Clear out the name. We put the message in label now, even if name was
		// the item defined.  TODO: Remove this after 2.0.
		nm["name"] = ""
		switch {
		case nm["keep_on_remove"].(bool):
			fallthrough
		case nm["attach"].(bool):
			nm["label"] = diskDetachedName
		default:
			nm["label"] = diskDeletedName
		}
		normalized[ni] = nm
	}

	// All done. We can end the customization off by setting the new, normalized diff.
	log.Printf("[DEBUG] DiskDiffOperation: New resource set post-normalization: %s", subresourceListString(normalized))
	log.Printf("[DEBUG] DiskDiffOperation: Disk diff customization complete, sending new diff")
	return d.SetNew(subresourceTypeDisk, normalized)
}

// DiskCloneValidateOperation takes the VirtualDeviceList, which should come
// from a source VM or template, and validates the following:
//
// * There are at least as many disks defined in the configuration as there are
// in the source VM or template.
// * All disks survive a disk sub-resource read operation.
//
// This function is meant to be called during diff customization. It is a
// subset of the normal refresh behaviour as we don't worry about checking
// existing state.
func DiskCloneValidateOperation(d *schema.ResourceDiff, c *govmomi.Client, l object.VirtualDeviceList, linked bool) error {
	log.Printf("[DEBUG] DiskCloneValidateOperation: Checking existing virtual disk configuration")
	devices := SelectDisks(l, d.Get("scsi_controller_count").(int))
	// Sort the device list, in case it's not sorted already.
	devSort := virtualDeviceListSorter{
		Sort:       devices,
		DeviceList: l,
	}
	log.Printf("[DEBUG] DiskCloneValidateOperation: Disk devices order before sort: %s", DeviceListString(devices))
	sort.Sort(devSort)
	devices = devSort.Sort
	log.Printf("[DEBUG] DiskCloneValidateOperation: Disk devices order after sort: %s", DeviceListString(devices))
	// Do the same for our listed disks.
	curSet := d.Get(subresourceTypeDisk).([]interface{})
	log.Printf("[DEBUG] DiskCloneValidateOperation: Current resource set: %s", subresourceListString(curSet))
	sort.Sort(virtualDiskSubresourceSorter(curSet))
	log.Printf("[DEBUG] DiskCloneValidateOperation: Resource set order after sort: %s", subresourceListString(curSet))

	// Quickly validate length. If there are more disks in the template than
	// there is in the configuration, kick out an error.
	if len(devices) > len(curSet) {
		return fmt.Errorf("not enough disks in configuration - you need at least %d to use this template (current: %d)", len(devices), len(curSet))
	}

	// Do test read operations on all disks.
	log.Printf("[DEBUG] DiskCloneValidateOperation: Running test read operations on all disks")
	for i, device := range devices {
		m := make(map[string]interface{})
		vd := device.GetVirtualDevice()
		ctlr := l.FindByKey(vd.ControllerKey)
		if ctlr == nil {
			return fmt.Errorf("could not find controller with key %d", vd.Key)
		}
		m["key"] = int(vd.Key)
		var err error
		m["device_address"], err = computeDevAddr(vd, ctlr.(types.BaseVirtualController))
		if err != nil {
			return fmt.Errorf("error computing device address: %s", err)
		}
		r := NewDiskSubresource(c, d, m, nil, i)
		if err := r.Read(l); err != nil {
			return fmt.Errorf("%s: validation failed (%s)", r.Addr(), err)
		}
		// Load the target resource to do a few comparisons for correctness in config.
		targetM := curSet[i].(map[string]interface{})
		tr := NewDiskSubresource(c, d, targetM, nil, i)

		// Do some pre-clone validation. This is mainly to make sure that the disks
		// clone in a way that is consistent with configuration.
		targetName, err := diskLabelOrName(tr.Data())
		if err != nil {
			return fmt.Errorf("%s: %s", tr.Addr(), err)
		}
		targetPath := r.Get("path").(string)
		sourceSize := r.Get("size").(int)
		targetSize := tr.Get("size").(int)
		targetThin := tr.Get("thin_provisioned").(bool)
		targetEager := tr.Get("eagerly_scrub").(bool)

		var sourceThin, sourceEager bool
		if b := r.Get("thin_provisioned"); b != nil {
			sourceThin = b.(bool)
		}
		if b := r.Get("eagerly_scrub"); b != nil {
			sourceEager = b.(bool)
		}

		switch {
		case linked:
			switch {
			case sourceSize != targetSize:
				return fmt.Errorf("%s: disk name %s must be the exact size of source when using linked_clone (expected: %d GiB)", tr.Addr(), targetName, sourceSize)
			case sourceThin != targetThin:
				return fmt.Errorf("%s: disk name %s must have same value for thin_provisioned as source when using linked_clone (expected: %t)", tr.Addr(), targetName, sourceThin)
			case sourceEager != targetEager:
				return fmt.Errorf("%s: disk name %s must have same value for eagerly_scrub as source when using linked_clone (expected: %t)", tr.Addr(), targetName, sourceEager)
			}
		default:
			if sourceSize > targetSize {
				return fmt.Errorf("%s: disk name %s must be at least the same size of source when cloning (expected: >= %d GiB)", tr.Addr(), targetName, sourceSize)
			}
		}

		// Finally, we don't support non-SCSI (ie: SATA, IDE, NVMe) disks, so kick
		// back an error if we see one of those.
		ct, _, _, err := splitDevAddr(r.DevAddr())
		if err != nil {
			return fmt.Errorf("%s: error parsing device address after reading disk %q: %s", tr.Addr(), targetPath, err)
		}
		if ct != SubresourceControllerTypeSCSI {
			return fmt.Errorf("%s: unsupported controller type %s for disk %q. Please use a template with SCSI disks only", tr.Addr(), ct, targetPath)
		}
	}
	log.Printf("[DEBUG] DiskCloneValidateOperation: All disks in source validated successfully")
	return nil
}

// DiskMigrateRelocateOperation assembles the
// VirtualMachineRelocateSpecDiskLocator slice for a virtual machine migration
// operation, otherwise known as storage vMotion.
func DiskMigrateRelocateOperation(d *schema.ResourceData, c *govmomi.Client, l object.VirtualDeviceList) ([]types.VirtualMachineRelocateSpecDiskLocator, bool, error) {
	log.Printf("[DEBUG] DiskMigrateRelocateOperation: Generating any necessary disk relocate specs")
	ods, nds := d.GetChange(subresourceTypeDisk)

	var relocators []types.VirtualMachineRelocateSpecDiskLocator
	var relocateOK bool

	// We are only concerned with resources that would normally be updated, as
	// incoming or outgoing disks obviously won't need migrating. Hence, this is
	// a simplified subset of the normal apply logic.
	for ni, ne := range nds.([]interface{}) {
		nm := ne.(map[string]interface{})
		var name string
		var err error
		if name, err = diskLabelOrName(nm); err != nil {
			return nil, false, fmt.Errorf("disk.%d: %s", ni, err)
		}
		if name == diskDeletedName || name == diskDetachedName {
			continue
		}
		for _, oe := range ods.([]interface{}) {
			om := oe.(map[string]interface{})
			if nm["uuid"] == om["uuid"] {
				// No change in datastore is a no-op, unless we are changing default datastores
				if nm["datastore_id"] == om["datastore_id"] && !d.HasChange("datastore_id") {
					break
				}
				// If we got this far, some sort of datastore migration will be
				// necessary. Flag this now.
				relocateOK = true

				// A disk locator is only useful if a target datastore is available. If we
				// don't have a datastore specified (ie: when Storage DRS is in use), then
				// we just need to skip this disk. The disk will be migrated properly
				// through the SDRS API.
				if nm["datastore_id"] == "" || nm["datastore_id"] == diskDatastoreComputedName {
					break
				}
				r := NewDiskSubresource(c, d, nm, om, ni)
				relocator, err := r.Relocate(l, false)
				if err != nil {
					return nil, false, fmt.Errorf("%s: %s", r.Addr(), err)
				}
				if d.Get("datastore_id").(string) == relocator.Datastore.Value {
					log.Printf("[DEBUG] %s: Datastore in spec is same as default, dropping in favor of implicit relocation", r.Addr())
					break
				}
				relocators = append(relocators, relocator)
			}
		}
	}

	if !relocateOK {
		log.Printf("[DEBUG] DiskMigrateRelocateOperation: Disk relocation not necessary")
		return nil, false, nil
	}

	log.Printf("[DEBUG] DiskMigrateRelocateOperation: Disk relocator list: %s", diskRelocateListString(relocators))
	log.Printf("[DEBUG] DiskMigrateRelocateOperation: Disk relocator generation complete")
	return relocators, true, nil
}

// DiskCloneRelocateOperation assembles the
// VirtualMachineRelocateSpecDiskLocator slice for a virtual machine clone
// operation.
//
// This differs from a regular storage vMotion in that we have no existing
// devices in the resource to work off of - the disks in the source virtual
// machine is our source of truth. These disks are assigned to our disk
// sub-resources in config and the relocate specs are generated off of the
// backing data defined in config, taking on these filenames when cloned. After
// the clone is complete, natural re-configuration happens to bring the disk
// configurations fully in sync with what is defined.
func DiskCloneRelocateOperation(d *schema.ResourceData, c *govmomi.Client, l object.VirtualDeviceList) ([]types.VirtualMachineRelocateSpecDiskLocator, error) {
	log.Printf("[DEBUG] DiskCloneRelocateOperation: Generating full disk relocate spec list")
	devices := SelectDisks(l, d.Get("scsi_controller_count").(int))
	log.Printf("[DEBUG] DiskCloneRelocateOperation: Disk devices located: %s", DeviceListString(devices))
	// Sort the device list, in case it's not sorted already.
	devSort := virtualDeviceListSorter{
		Sort:       devices,
		DeviceList: l,
	}
	sort.Sort(devSort)
	devices = devSort.Sort
	log.Printf("[DEBUG] DiskCloneRelocateOperation: Disk devices order after sort: %s", DeviceListString(devices))
	// Do the same for our listed disks.
	curSet := d.Get(subresourceTypeDisk).([]interface{})
	log.Printf("[DEBUG] DiskCloneRelocateOperation: Current resource set: %s", subresourceListString(curSet))
	sort.Sort(virtualDiskSubresourceSorter(curSet))
	log.Printf("[DEBUG] DiskCloneRelocateOperation: Resource set order after sort: %s", subresourceListString(curSet))

	log.Printf("[DEBUG] DiskCloneRelocateOperation: Generating relocators for source disks")
	var relocators []types.VirtualMachineRelocateSpecDiskLocator
	for i, device := range devices {
		m := curSet[i].(map[string]interface{})
		vd := device.GetVirtualDevice()
		ctlr := l.FindByKey(vd.ControllerKey)
		if ctlr == nil {
			return nil, fmt.Errorf("could not find controller with key %d", vd.Key)
		}
		m["key"] = int(vd.Key)
		var err error
		m["device_address"], err = computeDevAddr(vd, ctlr.(types.BaseVirtualController))
		if err != nil {
			return nil, fmt.Errorf("error computing device address: %s", err)
		}
		r := NewDiskSubresource(c, d, m, nil, i)
		// A disk locator is only useful if a target datastore is available. If we
		// don't have a datastore specified (ie: when Storage DRS is in use), then
		// we just need to skip this disk. The disk will be migrated properly
		// through the SDRS API.
		if dsID := r.Get("datastore_id"); dsID == "" || dsID == diskDatastoreComputedName {
			continue
		}
		// Otherwise, proceed with generating and appending the locator.
		relocator, err := r.Relocate(l, true)
		if err != nil {
			return nil, fmt.Errorf("%s: %s", r.Addr(), err)
		}
		relocators = append(relocators, relocator)
	}

	log.Printf("[DEBUG] DiskCloneRelocateOperation: Disk relocator list: %s", diskRelocateListString(relocators))
	log.Printf("[DEBUG] DiskCloneRelocateOperation: Disk relocator generation complete")
	return relocators, nil
}

// DiskPostCloneOperation normalizes the virtual disks on a freshly-cloned
// virtual machine and outputs any necessary device change operations. It also
// sets the state in advance of the post-create read.
//
// This differs from a regular apply operation in that a configuration is
// already present, but we don't have any existing state, which the standard
// virtual device operations rely pretty heavily on.
func DiskPostCloneOperation(d *schema.ResourceData, c *govmomi.Client, l object.VirtualDeviceList) (object.VirtualDeviceList, []types.BaseVirtualDeviceConfigSpec, error) {
	log.Printf("[DEBUG] DiskPostCloneOperation: Looking for disk device changes post-clone")
	devices := SelectDisks(l, d.Get("scsi_controller_count").(int))
	log.Printf("[DEBUG] DiskPostCloneOperation: Disk devices located: %s", DeviceListString(devices))
	// Sort the device list, in case it's not sorted already.
	devSort := virtualDeviceListSorter{
		Sort:       devices,
		DeviceList: l,
	}
	sort.Sort(devSort)
	devices = devSort.Sort
	log.Printf("[DEBUG] DiskPostCloneOperation: Disk devices order after sort: %s", DeviceListString(devices))
	// Do the same for our listed disks.
	curSet := d.Get(subresourceTypeDisk).([]interface{})
	log.Printf("[DEBUG] DiskPostCloneOperation: Current resource set: %s", subresourceListString(curSet))
	sort.Sort(virtualDiskSubresourceSorter(curSet))
	log.Printf("[DEBUG] DiskPostCloneOperation: Resource set order after sort: %s", subresourceListString(curSet))

	var spec []types.BaseVirtualDeviceConfigSpec
	var updates []interface{}

	log.Printf("[DEBUG] DiskPostCloneOperation: Looking for and applying device changes in source disks")
	for i, device := range devices {
		src := curSet[i].(map[string]interface{})
		vd := device.GetVirtualDevice()
		ctlr := l.FindByKey(vd.ControllerKey)
		if ctlr == nil {
			return nil, nil, fmt.Errorf("could not find controller with key %d", vd.Key)
		}
		src["key"] = int(vd.Key)
		var err error
		src["device_address"], err = computeDevAddr(vd, ctlr.(types.BaseVirtualController))
		if err != nil {
			return nil, nil, fmt.Errorf("error computing device address: %s", err)
		}
		// Copy the source set into old. This allows us to patch a copy of the
		// product of this set with the source, creating a diff.
		old, err := copystructure.Copy(src)
		if err != nil {
			return nil, nil, fmt.Errorf("error copying source set for disk at unit_number %d: %s", src["unit_number"].(int), err)
		}
		rOld := NewDiskSubresource(c, d, old.(map[string]interface{}), nil, i)
		if err := rOld.Read(l); err != nil {
			return nil, nil, fmt.Errorf("%s: %s", rOld.Addr(), err)
		}
		new, err := copystructure.Copy(rOld.Data())
		if err != nil {
			return nil, nil, fmt.Errorf("error copying current device state for disk at unit_number %d: %s", src["unit_number"].(int), err)
		}
		for k, v := range src {
			// Skip label, path (path will always be computed here as cloned disks
			// are not being attached externally), name, datastore_id, and uuid. Also
			// skip share_count if we the share level isn't custom.
			//
			// TODO: Remove "name" after 2.0.
			switch k {
			case "label", "path", "name", "datastore_id", "uuid", "thin_provisioned", "eagerly_scrub":
				continue
			case "io_share_count":
				if src["io_share_level"] != string(types.SharesLevelCustom) {
					continue
				}
			}
			new.(map[string]interface{})[k] = v
		}
		rNew := NewDiskSubresource(c, d, new.(map[string]interface{}), rOld.Data(), i)
		if !reflect.DeepEqual(rNew.Data(), rOld.Data()) {
			uspec, err := rNew.Update(l)
			if err != nil {
				return nil, nil, fmt.Errorf("%s: %s", rNew.Addr(), err)
			}
			l = applyDeviceChange(l, uspec)
			spec = append(spec, uspec...)
		}
		updates = append(updates, rNew.Data())
	}

	// Any disk past the current device list is a new device. Create those now.
	for _, ni := range curSet[len(devices):] {
		r := NewDiskSubresource(c, d, ni.(map[string]interface{}), nil, len(updates))
		cspec, err := r.Create(l)
		if err != nil {
			return nil, nil, fmt.Errorf("%s: %s", r.Addr(), err)
		}
		l = applyDeviceChange(l, cspec)
		spec = append(spec, cspec...)
		updates = append(updates, r.Data())
	}

	log.Printf("[DEBUG] DiskPostCloneOperation: Post-clone final resource list: %s", subresourceListString(updates))
	if err := d.Set(subresourceTypeDisk, updates); err != nil {
		return nil, nil, err
	}
	log.Printf("[DEBUG] DiskPostCloneOperation: Device list at end of operation: %s", DeviceListString(l))
	log.Printf("[DEBUG] DiskPostCloneOperation: Device config operations from post-clone: %s", DeviceChangeString(spec))
	log.Printf("[DEBUG] DiskPostCloneOperation: Operation complete, returning updated spec")
	return l, spec, nil
}

// DiskImportOperation validates the disk configuration of the virtual
// machine's VirtualDeviceList to ensure it will be imported properly, and also
// saves device addresses into state for disks defined in config. Both the
// imported device list is sorted by the device's unit number on the SCSI bus.
func DiskImportOperation(d *schema.ResourceData, c *govmomi.Client, l object.VirtualDeviceList) error {
	log.Printf("[DEBUG] DiskImportOperation: Performing pre-read import and validation of virtual disks")
	devices := SelectDisks(l, d.Get("scsi_controller_count").(int))
	// Sort the device list, in case it's not sorted already.
	devSort := virtualDeviceListSorter{
		Sort:       devices,
		DeviceList: l,
	}
	log.Printf("[DEBUG] DiskImportOperation: Disk devices order before sort: %s", DeviceListString(devices))
	sort.Sort(devSort)
	devices = devSort.Sort
	log.Printf("[DEBUG] DiskImportOperation: Disk devices order after sort: %s", DeviceListString(devices))

	// Read in the disks. We don't do anything with the results here other than
	// validate that the disks are SCSI disks. The read operation validates the rest.
	var curSet []interface{}
	log.Printf("[DEBUG] DiskImportOperation: Validating disk type and saving ")
	for i, device := range devices {
		vd := device.GetVirtualDevice()
		ctlr := l.FindByKey(vd.ControllerKey)
		if ctlr == nil {
			return fmt.Errorf("could not find controller with key %d", vd.Key)
		}
		addr, err := computeDevAddr(vd, ctlr.(types.BaseVirtualController))
		if err != nil {
			return fmt.Errorf("error computing device address: %s", err)
		}
		ct, _, _, err := splitDevAddr(addr)
		if err != nil {
			return fmt.Errorf("disk.%d: error parsing device address %s: %s", i, addr, err)
		}
		if ct != SubresourceControllerTypeSCSI {
			return fmt.Errorf("disk.%d: unsupported controller type %s for disk %s. The VM resource supports SCSI disks only", i, ct, addr)
		}
		// As one final validation, as we are no longer reading here, validate that
		// this is a VMDK-backed virtual disk to make sure we aren't importing RDM
		// disks or what not. The device should have already been validated as a
		// virtual disk via SelectDisks.
		if _, ok := device.(*types.VirtualDisk).Backing.(*types.VirtualDiskFlatVer2BackingInfo); !ok {
			return fmt.Errorf(
				"disk.%d: unsupported disk type at %s (expected flat VMDK version 2, got %T)",
				i,
				addr,
				device.(*types.VirtualDisk).Backing,
			)
		}
		m := make(map[string]interface{})
		// Save information so that the next DiskRefreshOperation can pick this
		// disk up as if it was newly added and not attempt to try and line up
		// UUIDs. We use a negative key for this reason, in addition to assigning
		// the device address.
		m["key"] = (i + 1) * -1
		m["device_address"] = addr
		// Assign a computed label. This label *needs* be the label this disk is
		// assigned in config, or you risk service interruptions or data corruption.
		m["label"] = fmt.Sprintf("disk%d", i)
		// Set keep_on_remove to ensure that if labels are assigned incorrectly,
		// all that happens is that the disk is removed. The comments above
		// regarding the risk of incorrect label assignment are still true, but
		// this greatly reduces the risk of data loss.
		m["keep_on_remove"] = true

		curSet = append(curSet, m)
	}
	log.Printf("[DEBUG] DiskImportOperation: Discovered disks from import: %s", subresourceListString(curSet))
	if err := d.Set(subresourceTypeDisk, curSet); err != nil {
		return err
	}
	log.Printf("[DEBUG] DiskImportOperation: Pre-read import and validation complete")
	return nil
}

// ReadDiskAttrsForDataSource returns select attributes from the list of disks
// on a virtual machine. This is used in the VM data source to discover
// specific options of all of the disks on the virtual machine sorted by the
// order that they would be added in if a clone were to be done.
func ReadDiskAttrsForDataSource(l object.VirtualDeviceList, count int) ([]map[string]interface{}, error) {
	log.Printf("[DEBUG] ReadDiskAttrsForDataSource: Fetching select attributes for disks across %d SCSI controllers", count)
	devices := SelectDisks(l, count)
	log.Printf("[DEBUG] ReadDiskAttrsForDataSource: Disk devices located: %s", DeviceListString(devices))
	// Sort the device list, in case it's not sorted already.
	devSort := virtualDeviceListSorter{
		Sort:       devices,
		DeviceList: l,
	}
	sort.Sort(devSort)
	devices = devSort.Sort
	log.Printf("[DEBUG] ReadDiskAttrsForDataSource: Disk devices order after sort: %s", DeviceListString(devices))
	var out []map[string]interface{}
	for i, device := range devices {
		disk := device.(*types.VirtualDisk)
		backing, ok := disk.Backing.(*types.VirtualDiskFlatVer2BackingInfo)
		if !ok {
			return nil, fmt.Errorf("disk number %d has an unsupported backing type (expected flat VMDK version 2, got %T)", i, disk.Backing)
		}
		m := make(map[string]interface{})
		var eager, thin bool
		if backing.EagerlyScrub != nil {
			eager = *backing.EagerlyScrub
		}
		if backing.ThinProvisioned != nil {
			thin = *backing.ThinProvisioned
		}
		m["size"] = diskCapacityInGiB(disk)
		m["eagerly_scrub"] = eager
		m["thin_provisioned"] = thin
		out = append(out, m)
	}
	log.Printf("[DEBUG] ReadDiskAttrsForDataSource: Attributes returned: %+v", out)
	return out, nil
}

// Create creates a vsphere_virtual_machine disk sub-resource.
func (r *DiskSubresource) Create(l object.VirtualDeviceList) ([]types.BaseVirtualDeviceConfigSpec, error) {
	log.Printf("[DEBUG] %s: Creating disk", r)
	var spec []types.BaseVirtualDeviceConfigSpec

	disk, err := r.createDisk(l)
	if err != nil {
		return nil, fmt.Errorf("error creating disk: %s", err)
	}
	// We now have the controller on which we can create our device on.
	// Assign the disk to a controller.
	ctlr, err := r.assignDisk(l, disk)
	if err != nil {
		return nil, fmt.Errorf("cannot assign disk: %s", err)
	}

	if err := r.expandDiskSettings(disk); err != nil {
		return nil, err
	}

	// Done here. Save ID, push the device to the new device list and return.
	if err := r.SaveDevIDs(disk, ctlr); err != nil {
		return nil, err
	}
	dspec, err := object.VirtualDeviceList{disk}.ConfigSpec(types.VirtualDeviceConfigSpecOperationAdd)
	if err != nil {
		return nil, err
	}
	if len(dspec) != 1 {
		return nil, fmt.Errorf("incorrect number of config spec items returned - expected 1, got %d", len(dspec))
	}
	// Clear the file operation if we are attaching.
	if r.Get("attach").(bool) {
		dspec[0].GetVirtualDeviceConfigSpec().FileOperation = ""
	}

	// Attach the SPBM storage policy if specified
	if policyID := r.Get("storage_policy_id").(string); policyID != "" {
		dspec[0].GetVirtualDeviceConfigSpec().Profile = spbm.PolicySpecByID(policyID)
	}

	spec = append(spec, dspec...)
	log.Printf("[DEBUG] %s: Device config operations from create: %s", r, DeviceChangeString(spec))
	log.Printf("[DEBUG] %s: Create finished", r)
	return spec, nil
}

// Read reads a vsphere_virtual_machine disk sub-resource and commits the data
// to the newData layer.
func (r *DiskSubresource) Read(l object.VirtualDeviceList) error {
	log.Printf("[DEBUG] %s: Reading state", r)
	disk, err := r.findVirtualDisk(l, true)
	if err != nil {
		return fmt.Errorf("cannot find disk device: %s", err)
	}
	unit, ctlr, err := r.findControllerInfo(l, disk)
	if err != nil {
		return err
	}
	r.Set("unit_number", unit)
	if err := r.SaveDevIDs(disk, ctlr); err != nil {
		return err
	}

	// Fetch disk attachment state in config
	var attach bool
	if r.Get("attach") != nil {
		attach = r.Get("attach").(bool)
	}
	// Save disk backing settings
	b, ok := disk.Backing.(*types.VirtualDiskFlatVer2BackingInfo)
	if !ok {
		return fmt.Errorf("disk backing at %s is of an unsupported type (type %T)", r.Get("device_address").(string), disk.Backing)
	}
	r.Set("uuid", b.Uuid)
	r.Set("disk_mode", b.DiskMode)
	r.Set("write_through", b.WriteThrough)

	// Only use disk_sharing if we are on vSphere 6.0 and higher. In addition,
	// skip if the value is unset - this prevents spurious diffs during upgrade
	// situations where the VM hardware version does not actually allow disk
	// sharing. In this situation, the value will be blank, and setting it will
	// actually result in an error.
	version := viapi.ParseVersionFromClient(r.client)
	if version.Newer(viapi.VSphereVersion{Product: version.Product, Major: 6}) && b.Sharing != "" {
		r.Set("disk_sharing", b.Sharing)
	}

	if !attach {
		r.Set("thin_provisioned", b.ThinProvisioned)
		r.Set("eagerly_scrub", b.EagerlyScrub)
	}
	r.Set("datastore_id", b.Datastore.Value)

	// Disk settings
	if !attach {
		dp := &object.DatastorePath{}
		if ok := dp.FromString(b.FileName); !ok {
			return fmt.Errorf("could not parse path from filename: %s", b.FileName)
		}
		r.Set("path", dp.Path)
		r.Set("size", diskCapacityInGiB(disk))
	}

	if allocation := disk.StorageIOAllocation; allocation != nil {
		r.Set("io_limit", allocation.Limit)
		r.Set("io_reservation", allocation.Reservation)
		if shares := allocation.Shares; shares != nil {
			r.Set("io_share_level", string(shares.Level))
			r.Set("io_share_count", shares.Shares)
		}
	}

	// Set storage policy if the VM exists.
	vmUUID := r.rdd.Id()
	if vmUUID != "" {
		result, err := virtualmachine.MOIDForUUID(r.client, vmUUID)
		if err != nil {
			return err
		}
		polID, err := spbm.PolicyIDByVirtualDisk(r.client, result.MOID, r.Get("key").(int))
		if err != nil {
			return err
		}
		r.Set("storage_policy_id", polID)
	}

	log.Printf("[DEBUG] %s: Read finished (key and device address may have changed)", r)
	return nil
}

// Update updates a vsphere_virtual_machine disk sub-resource.
func (r *DiskSubresource) Update(l object.VirtualDeviceList) ([]types.BaseVirtualDeviceConfigSpec, error) {
	log.Printf("[DEBUG] %s: Beginning update", r)
	disk, err := r.findVirtualDisk(l, false)
	if err != nil {
		return nil, fmt.Errorf("cannot find disk device: %s", err)
	}

	// Has the unit number changed?
	if r.HasChange("unit_number") {
		ctlr, err := r.assignDisk(l, disk)
		if err != nil {
			return nil, fmt.Errorf("cannot assign disk: %s", err)
		}
		r.SetRestart("unit_number")
		if err := r.SaveDevIDs(disk, ctlr); err != nil {
			return nil, fmt.Errorf("error saving device address: %s", err)
		}
		// A change in disk unit number forces a device key change after the
		// reconfigure. We need to keep the key in the device change spec we send
		// along, but we can reset it here safely. Set it to 0, which will send it
		// though the new device loop, but will distinguish it from newly-created
		// devices.
		r.Set("key", 0)
	}

	// We can now expand the rest of the settings.
	if err := r.expandDiskSettings(disk); err != nil {
		return nil, err
	}

	dspec, err := object.VirtualDeviceList{disk}.ConfigSpec(types.VirtualDeviceConfigSpecOperationEdit)
	if err != nil {
		return nil, err
	}
	if len(dspec) != 1 {
		return nil, fmt.Errorf("incorrect number of config spec items returned - expected 1, got %d", len(dspec))
	}
	// Clear file operation - VirtualDeviceList currently sets this to replace, which is invalid
	dspec[0].GetVirtualDeviceConfigSpec().FileOperation = ""

	// Attach the SPBM storage policy if specified
	if policyID := r.Get("storage_policy_id").(string); policyID != "" {
		dspec[0].GetVirtualDeviceConfigSpec().Profile = spbm.PolicySpecByID(policyID)
	}

	log.Printf("[DEBUG] %s: Device config operations from update: %s", r, DeviceChangeString(dspec))
	log.Printf("[DEBUG] %s: Update complete", r)
	return dspec, nil
}

// Delete deletes a vsphere_virtual_machine disk sub-resource.
func (r *DiskSubresource) Delete(l object.VirtualDeviceList) ([]types.BaseVirtualDeviceConfigSpec, error) {
	log.Printf("[DEBUG] %s: Beginning delete", r)
	disk, err := r.findVirtualDisk(l, false)
	if err != nil {
		return nil, fmt.Errorf("cannot find disk device: %s", err)
	}
	deleteSpec, err := object.VirtualDeviceList{disk}.ConfigSpec(types.VirtualDeviceConfigSpecOperationRemove)
	if err != nil {
		return nil, err
	}
	if len(deleteSpec) != 1 {
		return nil, fmt.Errorf("incorrect number of config spec items returned - expected 1, got %d", len(deleteSpec))
	}
	if r.Get("keep_on_remove").(bool) || r.Get("attach").(bool) {
		// Clear file operation so that the disk is kept on remove.
		deleteSpec[0].GetVirtualDeviceConfigSpec().FileOperation = ""
	}
	log.Printf("[DEBUG] %s: Device config operations from update: %s", r, DeviceChangeString(deleteSpec))
	log.Printf("[DEBUG] %s: Delete completed", r)
	return deleteSpec, nil
}

// DiffExisting validates and normalizes the fields for an existing disk
// sub-resource.  It handles carrying over existing values, so this should not
// be used on disks that have not been successfully matched up between current
// and old diffs.
func (r *DiskSubresource) DiffExisting() error {
	log.Printf("[DEBUG] %s: Beginning normalization of existing disk", r)
	name, err := diskLabelOrName(r.data)
	if err != nil {
		return err
	}
	// Prevent a backward migration of label -> name. TODO: Remove this after
	// 2.0.
	olabel, nlabel := r.GetChange("label")
	if olabel != "" && nlabel == "" {
		return errors.New("cannot migrate from label to name")
	}
	// Carry forward the name attribute like we used to if no label is defined.
	// TODO: Remove this after 2.0.
	if nlabel == "" {
		oname, _ := r.GetChange("name")
		r.Set("name", oname.(string))
	}

	// set some computed fields: key, device_address, and uuid will always be
	// non-populated, so copy those.
	okey, _ := r.GetChange("key")
	odaddr, _ := r.GetChange("device_address")
	ouuid, _ := r.GetChange("uuid")
	r.Set("key", okey)
	r.Set("device_address", odaddr)
	r.Set("uuid", ouuid)

	if !r.Get("attach").(bool) {
		// Carry forward path when attach is not set
		opath, _ := r.GetChange("path")
		r.Set("path", opath.(string))
	}

	// Set the datastore if it's missing as we infer this from the default
	// datastore in that case
	if r.Get("datastore_id") == "" {
		switch {
		case r.rdd.HasChange("datastore_id"):
			// If the default datastore is changing and we don't have a default
			// datastore here, we need to use the implicit setting here to indicate
			// that we may need to migrate. This allows us to differentiate between a
			// full storage vMotion no-op, an implicit migration, and a migration
			// where we will need to generate a relocate spec for the individual disk
			// to ensure it stays at a datastore it might be pinned on.
			dsID := r.rdd.Get("datastore_id").(string)
			if dsID == "" {
				r.Set("datastore_id", diskDatastoreComputedName)
			} else {
				r.Set("datastore_id", dsID)
			}
		default:
			if err = r.normalizeDiskDatastore(); err != nil {
				return err
			}
		}
	}

	// Preserve the share value if we don't have custom shares set
	osc, _ := r.GetChange("io_share_count")
	if r.Get("io_share_level").(string) != string(types.SharesLevelCustom) {
		r.Set("io_share_count", osc)
	}

	// Ensure that the user is not attempting to shrink the disk. If we do more
	// we might want to change the name of this method, but we want to check this
	// here as CustomizeDiff is meant for vetoing.
	osize, nsize := r.GetChange("size")
	if osize.(int) > nsize.(int) {
		return fmt.Errorf("virtual disk %q: virtual disks cannot be shrunk (old: %d new: %d)", name, osize.(int), nsize.(int))
	}

	// Ensure that there is no change in either eagerly_scrub or thin_provisioned
	// - these values cannot be changed once set.
	if _, err = r.GetWithVeto("eagerly_scrub"); err != nil {
		return fmt.Errorf("virtual disk %q: %s", name, err)
	}
	if _, err = r.GetWithVeto("thin_provisioned"); err != nil {
		return fmt.Errorf("virtual disk %q: %s", name, err)
	}
	// Same with attach
	if _, err = r.GetWithVeto("attach"); err != nil {
		return fmt.Errorf("virtual disk %q: %s", name, err)
	}

	// Validate storage vMotion if the datastore is changing
	if r.HasChange("datastore_id") {
		if err = r.validateStorageRelocateDiff(); err != nil {
			return err
		}
	}
	log.Printf("[DEBUG] %s: Normalization of existing disk diff complete", r)
	return nil
}

// DiffGeneral performs complex validation of an individual disk sub-resource
// that can't be done in schema alone. Should be run on new and existing disks.
func (r *DiskSubresource) DiffGeneral() error {
	log.Printf("[DEBUG] %s: Beginning diff validation", r)
	name, err := diskLabelOrName(r.data)
	if err != nil {
		return err
	}

	// Enforce the maximum unit number, which is the current value of
	// scsi_controller_count * 15 - 1.
	ctlrCount := r.rdd.Get("scsi_controller_count").(int)
	maxUnit := ctlrCount*15 - 1
	currentUnit := r.Get("unit_number").(int)
	if currentUnit > maxUnit {
		return fmt.Errorf("unit_number on disk %q too high (%d) - maximum value is %d with %d SCSI controller(s)", name, currentUnit, maxUnit, ctlrCount)
	}

	if r.Get("attach").(bool) {
		switch {
		case r.Get("datastore_id").(string) == "":
			return fmt.Errorf("datastore_id for disk %q is required when attach is set", name)
		case r.Get("size").(int) > 0:
			return fmt.Errorf("size for disk %q cannot be defined when attach is set", name)
		case r.Get("eagerly_scrub").(bool):
			return fmt.Errorf("eagerly_scrub for disk %q cannot be defined when attach is set", name)
		case r.Get("keep_on_remove").(bool):
			return fmt.Errorf("keep_on_remove for disk %q is implicit when attach is set, please remove this setting", name)
		}
	} else {
		// Enforce size as a required field when attach is not set
		if r.Get("size").(int) < 1 {
			return fmt.Errorf("size for disk %q: required option not set", name)
		}
	}
	// Block certain options from being set depending on the vSphere version.
	version := viapi.ParseVersionFromClient(r.client)
	if r.Get("disk_sharing").(string) != string(types.VirtualDiskSharingSharingNone) {
		if version.Older(viapi.VSphereVersion{Product: version.Product, Major: 6}) {
			return fmt.Errorf("multi-writer disk_sharing is only supported on vSphere 6 and higher")
		}
	}
	// Prevent eagerly_scrub and thin_provisioned from both being set to true. A
	// thin_provisioned disk cannot be eagerly scrubbed since it would then be
	// allocating the entire disk.
	if r.Get("eagerly_scrub").(bool) && r.Get("thin_provisioned").(bool) {
		return fmt.Errorf("%s: eagerly_scrub and thin_provisioned cannot both be set to true", name)
	}
	log.Printf("[DEBUG] %s: Diff validation complete", r)
	return nil
}

// normalizeDiskDatastore normalizes the datastore_id field in a disk
// sub-resource. If the VM has a datastore cluster defined, it checks to make
// sure the datastore in the current state of the disk is a member of the
// currently defined datastore cluster, and if it is not, it marks the disk as
// computed so that it can be migrated back to the datastore cluster on the
// next update.
func (r *DiskSubresource) normalizeDiskDatastore() error {
	podID := r.rdd.Get("datastore_cluster_id").(string)
	dsID, _ := r.GetChange("datastore_id")

	if podID == "" {
		// We don't have a storage pod, just set the old ID and exit. We don't need
		// to worry about whether or not the storage pod is computed here as if it
		// is, the VM datastore will have been marked as computed and this function
		// will have never ran.
		r.Set("datastore_id", dsID)
		return nil
	}

	log.Printf("[DEBUG] %s: Checking datastore cluster membership of disk", r)

	pod, err := storagepod.FromID(r.client, podID)
	if err != nil {
		return fmt.Errorf("error fetching datastore cluster ID %q: %s", podID, err)
	}

	ds, err := datastore.FromID(r.client, dsID.(string))
	if err != nil {
		return fmt.Errorf("error fetching datastore ID %q: %s", dsID, err)
	}

	isMember, err := storagepod.IsMember(pod, ds)
	if err != nil {
		return fmt.Errorf("error checking storage pod membership: %s", err)
	}
	if !isMember {
		log.Printf(
			"[DEBUG] %s: Disk's datastore %q not a member of cluster %q, marking datastore ID as computed",
			r,
			ds.Name(),
			pod.Name(),
		)
		dsID = diskDatastoreComputedName
	}

	r.Set("datastore_id", dsID)
	return nil
}

// validateStorageRelocateDiff validates certain storage vMotion diffs to make
// sure they are functional. These mainly have to do with limitations
// associated with our tracking of virtual disks via their names.
//
// The current limitations are:
//
// * Externally-attached virtual disks are not allowed to be vMotioned.
// * Disks must match the vSphere naming convention, where the first disk is
// named VMNAME.vmdk, and all other disks are named VMNAME_INDEX.vmdk This is a
// validation we use for cloning as well.
// * Any VM that has been created by a linked clone is blocked from storage
// vMotion full stop.
//
// TODO: Once we have solved the disk tracking issue and are no longer tracking
// disks via their file names, the only restriction that should remain is for
// externally attached disks. That restriction will go away once we figure out
// a strategy for handling when said disks have been moved OOB of the VM
// workflow.
func (r *DiskSubresource) validateStorageRelocateDiff() error {
	log.Printf("[DEBUG] %s: Validating storage vMotion eligibility", r)
	if err := r.blockRelocateAttachedDisks(); err != nil {
		return err
	}
	log.Printf("[DEBUG] %s: Storage vMotion validation successful", r)
	return nil
}

func (r *DiskSubresource) blockRelocateAttachedDisks() error {
	attach := r.Get("attach")
	if attach == nil {
		return nil
	}
	if attach.(bool) {
		return fmt.Errorf("externally attached disk %q cannot be migrated", diskPathOrName(r.data))
	}
	return nil
}

// Relocate produces a VirtualMachineRelocateSpecDiskLocator for this resource
// and is used for both cloning and storage vMotion.
func (r *DiskSubresource) Relocate(l object.VirtualDeviceList, clone bool) (types.VirtualMachineRelocateSpecDiskLocator, error) {
	log.Printf("[DEBUG] %s: Starting relocate generation", r)
	disk, err := r.findVirtualDisk(l, clone)
	var relocate types.VirtualMachineRelocateSpecDiskLocator
	if err != nil {
		return relocate, fmt.Errorf("cannot find disk device: %s", err)
	}

	// Expand all of the necessary disk settings first. This ensures all backing
	// data is properly populate and updated.
	if err := r.expandDiskSettings(disk); err != nil {
		return relocate, err
	}

	relocate.DiskId = disk.Key

	// Set the datastore for the relocation
	dsID := r.Get("datastore_id").(string)
	if dsID == "" {
		// Default to the default datastore
		dsID = r.rdd.Get("datastore_id").(string)
	}
	ds, err := datastore.FromID(r.client, dsID)
	if err != nil {
		return relocate, err
	}
	dsref := ds.Reference()
	relocate.Datastore = dsref

	// Add additional backing options if we are cloning.
	if r.rdd.Id() == "" {
		log.Printf("[DEBUG] %s: Adding additional options to relocator for cloning", r)

		backing := disk.Backing.(*types.VirtualDiskFlatVer2BackingInfo)
		backing.FileName = ds.Path("")
		backing.Datastore = &dsref
		relocate.DiskBackingInfo = backing
	}

	// Attach the SPBM storage policy if specified
	if policyID := r.Get("storage_policy_id").(string); policyID != "" {
		relocate.Profile = spbm.PolicySpecByID(policyID)
	}

	// Done!
	log.Printf("[DEBUG] %s: Generated disk locator: %s", r, diskRelocateString(relocate))
	log.Printf("[DEBUG] %s: Relocate generation complete", r)
	return relocate, nil
}

// String prints out the disk sub-resource's information including the ID at
// time of instantiation, the path of the disk, and the current device
// key and address.
func (r *DiskSubresource) String() string {
	p := diskPathOrName(r.data)
	if p == "" {
		p = "<unknown>"
	}
	return fmt.Sprintf("%s (%s)", r.Subresource.String(), p)
}

// expandDiskSettings sets appropriate fields on an existing disk - this is
// used during Create and Update to set attributes to those found in
// configuration.
func (r *DiskSubresource) expandDiskSettings(disk *types.VirtualDisk) error {
	// Backing settings
	b := disk.Backing.(*types.VirtualDiskFlatVer2BackingInfo)
	b.DiskMode = r.GetWithRestart("disk_mode").(string)
	b.WriteThrough = structure.BoolPtr(r.GetWithRestart("write_through").(bool))

	// Only use disk_sharing if we are on vSphere 6.0 and higher
	version := viapi.ParseVersionFromClient(r.client)
	if version.Newer(viapi.VSphereVersion{Product: version.Product, Major: 6}) {
		b.Sharing = r.GetWithRestart("disk_sharing").(string)
	}

	// This settings are only set for internal disks
	if !r.Get("attach").(bool) {
		var err error
		var v interface{}
		if v, err = r.GetWithVeto("thin_provisioned"); err != nil {
			return err
		}
		b.ThinProvisioned = structure.BoolPtr(v.(bool))

		if v, err = r.GetWithVeto("eagerly_scrub"); err != nil {
			return err
		}
		b.EagerlyScrub = structure.BoolPtr(v.(bool))

		// Disk settings
		os, ns := r.GetChange("size")
		if os.(int) > ns.(int) {
			return fmt.Errorf("virtual disks cannot be shrunk")
		}
		disk.CapacityInBytes = structure.GiBToByte(ns.(int))
		disk.CapacityInKB = disk.CapacityInBytes / 1024
	}

	alloc := &types.StorageIOAllocationInfo{
		Limit:       structure.Int64Ptr(int64(r.Get("io_limit").(int))),
		Reservation: structure.Int32Ptr(int32(r.Get("io_reservation").(int))),
		Shares: &types.SharesInfo{
			Shares: int32(r.Get("io_share_count").(int)),
			Level:  types.SharesLevel(r.Get("io_share_level").(string)),
		},
	}
	disk.StorageIOAllocation = alloc

	return nil
}

// createDisk performs all of the logic for a base virtual disk creation.
func (r *DiskSubresource) createDisk(l object.VirtualDeviceList) (*types.VirtualDisk, error) {
	disk := new(types.VirtualDisk)
	disk.Backing = new(types.VirtualDiskFlatVer2BackingInfo)

	// Only assign backing info if a datastore cluster is not specified. If one
	// is, skip this step.
	if r.rdd.Get("datastore_cluster_id").(string) == "" {
		if err := r.assignBackingInfo(disk); err != nil {
			return nil, err
		}
	}

	// Set a new device key for this device
	disk.Key = l.NewKey()
	return disk, nil
}

func (r *DiskSubresource) assignBackingInfo(disk *types.VirtualDisk) error {
	dsID := r.Get("datastore_id").(string)
	if dsID == "" || dsID == diskDatastoreComputedName {
		// Default to the default datastore
		dsID = r.rdd.Get("datastore_id").(string)

		if dsID == "" {
			vmObj, err := virtualmachine.FromUUID(r.client, r.rdd.Id())
			if err != nil {
				return err
			}

			vmprops, err := virtualmachine.Properties(vmObj)
			if err != nil {
				return err
			}
			if len(vmprops.Datastore) == 0 {
				return fmt.Errorf("no datastore was set and was unable to find a default to fall back to")
			}
			dsID = vmprops.Datastore[0].Value

		}
	}
	ds, err := datastore.FromID(r.client, dsID)
	if err != nil {
		return err
	}
	dsref := ds.Reference()

	var diskName string
	if r.Get("attach").(bool) {
		// No path interpolation is performed any more for attached disks - the
		// provided path must be the full path to the virtual disk you want to
		// attach.
		diskName = diskPathOrName(r.data)
	}

	backing := disk.Backing.(*types.VirtualDiskFlatVer2BackingInfo)
	backing.FileName = ds.Path(diskName)
	backing.Datastore = &dsref

	return nil
}

// assignDisk takes a unit number and assigns it correctly to a controller on
// the SCSI bus. An error is returned if the assigned unit number is taken.
func (r *DiskSubresource) assignDisk(l object.VirtualDeviceList, disk *types.VirtualDisk) (types.BaseVirtualController, error) {
	number := r.Get("unit_number").(int)
	// Figure out the bus number, and look up the SCSI controller that matches
	// that. You can attach 15 disks to a SCSI controller, and we allow a maximum
	// of 30 devices.
	bus := number / 15
	// Also determine the unit number on that controller.
	unit := int32(math.Mod(float64(number), 15))

	// Find the controller.
	ctlr, err := r.ControllerForCreateUpdate(l, SubresourceControllerTypeSCSI, bus)
	if err != nil {
		return nil, err
	}

	// Build the unit list.
	units := make([]bool, 16)
	// Reserve the SCSI unit number
	scsiUnit := ctlr.(types.BaseVirtualSCSIController).GetVirtualSCSIController().ScsiCtlrUnitNumber
	units[scsiUnit] = true

	ckey := ctlr.GetVirtualController().Key

	for _, device := range l {
		d := device.GetVirtualDevice()
		if d.ControllerKey != ckey || d.UnitNumber == nil {
			continue
		}
		units[*d.UnitNumber] = true
	}

	// We now have a valid list of units. If we need to, shift up the desired
	// unit number so it's not taking the unit of the controller itself.
	if unit >= scsiUnit {
		unit++
	}

	if units[unit] {
		return nil, fmt.Errorf("unit number %d on SCSI bus %d is in use", unit, bus)
	}

	// If we made it this far, we are good to go!
	disk.ControllerKey = ctlr.GetVirtualController().Key
	disk.UnitNumber = &unit
	return ctlr, nil
}

// findControllerInfo determines the normalized unit number for the disk device
// based on the SCSI controller and unit number it's connected to. The
// controller is also returned.
func (r *Subresource) findControllerInfo(l object.VirtualDeviceList, disk *types.VirtualDisk) (int, types.BaseVirtualController, error) {
	ctlr := l.FindByKey(disk.ControllerKey)
	if ctlr == nil {
		return -1, nil, fmt.Errorf("could not find disk controller with key %d for disk key %d", disk.ControllerKey, disk.Key)
	}
	if disk.UnitNumber == nil {
		return -1, nil, fmt.Errorf("unit number on disk key %d is unset", disk.Key)
	}
	sc, ok := ctlr.(types.BaseVirtualSCSIController)
	if !ok {
		return -1, nil, fmt.Errorf("controller at key %d is not a SCSI controller (actual: %T)", ctlr.GetVirtualDevice().Key, ctlr)
	}
	unit := *disk.UnitNumber
	if unit > sc.GetVirtualSCSIController().ScsiCtlrUnitNumber {
		unit--
	}
	unit = unit + 15*sc.GetVirtualSCSIController().BusNumber
	return int(unit), ctlr.(types.BaseVirtualController), nil
}

// diskRelocateListString pretty-prints a list of
// VirtualMachineRelocateSpecDiskLocator.
func diskRelocateListString(relocators []types.VirtualMachineRelocateSpecDiskLocator) string {
	var out []string
	for _, relocate := range relocators {
		out = append(out, diskRelocateString(relocate))
	}
	return strings.Join(out, ",")
}

// diskRelocateString prints out information from a
// VirtualMachineRelocateSpecDiskLocator in a friendly way.
//
// The format depends on whether or not a backing has been defined.
func diskRelocateString(relocate types.VirtualMachineRelocateSpecDiskLocator) string {
	key := relocate.DiskId
	var locstring string
	if backing, ok := relocate.DiskBackingInfo.(*types.VirtualDiskFlatVer2BackingInfo); ok && backing != nil {
		locstring = backing.FileName
	} else {
		locstring = relocate.Datastore.Value
	}
	return fmt.Sprintf("(%d => %s)", key, locstring)
}

// virtualDeviceListSorter is an internal type to facilitate sorting of a BaseVirtualDeviceList.
type virtualDeviceListSorter struct {
	Sort       object.VirtualDeviceList
	DeviceList object.VirtualDeviceList
}

// Len implements sort.Interface for virtualDeviceListSorter.
func (l virtualDeviceListSorter) Len() int {
	return len(l.Sort)
}

// Less helps implement sort.Interface for virtualDeviceListSorter. A
// BaseVirtualDevice is "less" than another device if its controller's bus
// number and unit number combination are earlier in the order than the other.
func (l virtualDeviceListSorter) Less(i, j int) bool {
	li := l.Sort[i]
	lj := l.Sort[j]
	liCtlr := l.DeviceList.FindByKey(li.GetVirtualDevice().ControllerKey)
	ljCtlr := l.DeviceList.FindByKey(lj.GetVirtualDevice().ControllerKey)
	if liCtlr == nil || ljCtlr == nil {
		panic(errors.New("virtualDeviceListSorter cannot be used with devices that are not assigned to a controller"))
	}
	liCtlrBus := liCtlr.(types.BaseVirtualController).GetVirtualController().BusNumber
	ljCtlrBus := ljCtlr.(types.BaseVirtualController).GetVirtualController().BusNumber
	if liCtlrBus != ljCtlrBus {
		return liCtlrBus < ljCtlrBus
	}
	liUnit := li.GetVirtualDevice().UnitNumber
	ljUnit := lj.GetVirtualDevice().UnitNumber
	if liUnit == nil || ljUnit == nil {
		panic(errors.New("virtualDeviceListSorter cannot be used with devices that do not have unit numbers set"))
	}
	return *liUnit < *ljUnit
}

// Swap helps implement sort.Interface for virtualDeviceListSorter.
func (l virtualDeviceListSorter) Swap(i, j int) {
	l.Sort[i], l.Sort[j] = l.Sort[j], l.Sort[i]
}

// virtualDiskSubresourceSorter sorts a list of disk sub-resources, based on unit number.
type virtualDiskSubresourceSorter []interface{}

// Len implements sort.Interface for virtualDiskSubresourceSorter.
func (s virtualDiskSubresourceSorter) Len() int {
	return len(s)
}

// Less helps implement sort.Interface for virtualDiskSubresourceSorter.
func (s virtualDiskSubresourceSorter) Less(i, j int) bool {
	mi := s[i].(map[string]interface{})
	mj := s[j].(map[string]interface{})
	return mi["unit_number"].(int) < mj["unit_number"].(int)
}

// Swap helps implement sort.Interface for virtualDiskSubresourceSorter.
func (s virtualDiskSubresourceSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// datastorePathHasBase is a helper to check if a datastore path's file matches
// a supplied file name.
func datastorePathHasBase(p, b string) bool {
	dp := &object.DatastorePath{}
	if ok := dp.FromString(p); !ok {
		return false
	}
	return path.Base(dp.Path) == path.Base(b)
}

// SelectDisks looks for disks that Terraform is supposed to manage. count is
// the number of controllers that Terraform is managing and serves as an upper
// limit (count - 1) of the SCSI bus number for a controller that eligible
// disks need to be attached to.
func SelectDisks(l object.VirtualDeviceList, count int) object.VirtualDeviceList {
	devices := l.Select(func(device types.BaseVirtualDevice) bool {
		if disk, ok := device.(*types.VirtualDisk); ok {
			ctlr, err := findControllerForDevice(l, disk)
			if err != nil {
				log.Printf("[DEBUG] DiskRefreshOperation: Error looking for controller for device %q: %s", l.Name(disk), err)
				return false
			}
			if sc, ok := ctlr.(types.BaseVirtualSCSIController); ok && sc.GetVirtualSCSIController().BusNumber < int32(count) {
				cd := sc.(types.BaseVirtualDevice)
				log.Printf("[DEBUG] DiskRefreshOperation: Found controller %q for device %q", l.Name(cd), l.Name(disk))
				return true
			}
		}
		return false
	})
	return devices
}

// diskLabelOrName is a helper method that returns the unique label for a disk
// - either its label or name. An error is returned if both are defined.
//
// TODO: This method will be removed in future releases.
func diskLabelOrName(data map[string]interface{}) (string, error) {
	var label, name string
	if v, ok := data["label"]; ok && v != nil {
		label = v.(string)
	}
	if v, ok := data["name"]; ok && v != nil {
		name = v.(string)
	}
	if name != "" {
		name = path.Base(name)
	}

	log.Printf("[DEBUG] diskLabelOrName: label: %q name: %q", label, name)
	switch {
	case label == "" && name == "":
		return "", errors.New("disk label or name must be defined and cannot be computed")
	case label != "" && name != "":
		return "", errors.New("disk label and name cannot be defined at the same time")
	case label != "":
		log.Printf("[DEBUG] diskLabelOrName: Using defined label value %q", label)
		return label, nil
	}
	log.Printf("[DEBUG] diskLabelOrName: Using defined name value as fallback %q", name)
	return name, nil
}

// diskPathOrName is a helper method that returns the path for a disk - either
// its path attribute or name as a fallback.
//
// TODO: This method will be removed in future releases.
func diskPathOrName(data map[string]interface{}) string {
	var path, name string
	if v, ok := data["path"]; ok && v != nil {
		path = v.(string)
	}
	if v, ok := data["name"]; ok && v != nil {
		name = v.(string)
	}
	if path != "" {
		log.Printf("[DEBUG] diskPathOrName: Using defined path value %q", path)
		return path
	}
	log.Printf("[DEBUG] diskPathOrName: Using defined name value as fallback %q", name)
	return name
}

// findVirtualDisk locates a virtual disk by it UUID, or by its device address
// if UUID is missing.
//
// The device address search is only used if fallback is true - this is so that
// we can distinguish situations where it should be used, such as a read,
// versus situations where it should never be used, such as an update or
// delete.
func (r *DiskSubresource) findVirtualDisk(l object.VirtualDeviceList, fallback bool) (*types.VirtualDisk, error) {
	device, err := r.findVirtualDiskByUUIDOrAddress(l, fallback)
	if err != nil {
		return nil, err
	}
	return device.(*types.VirtualDisk), nil
}

func (r *DiskSubresource) findVirtualDiskByUUIDOrAddress(l object.VirtualDeviceList, fallback bool) (types.BaseVirtualDevice, error) {
	var uuid string
	if v := r.Get("uuid"); v != nil {
		uuid = v.(string)
	}
	switch {
	case uuid == "" && fallback:
		return r.FindVirtualDevice(l)
	case uuid == "" && !fallback:
		return nil, errors.New("disk UUID is missing")
	}
	devices := l.Select(func(device types.BaseVirtualDevice) bool {
		return diskUUIDMatch(device, uuid)
	})
	switch {
	case len(devices) < 1:
		return nil, fmt.Errorf("virtual disk with UUID %s not found", uuid)
	case len(devices) > 1:
		// This is an edge/should never happen case
		return nil, fmt.Errorf("multiple virtual disks with UUID %s found", uuid)
	}
	return devices[0], nil
}

func diskUUIDMatch(device types.BaseVirtualDevice, uuid string) bool {
	disk, ok := device.(*types.VirtualDisk)
	if !ok {
		return false
	}
	backing, ok := disk.Backing.(*types.VirtualDiskFlatVer2BackingInfo)
	if !ok {
		return false
	}
	if backing.Uuid != uuid {
		return false
	}
	return true
}

// diskCapacityInGiB reports the supplied disk's capacity, by first checking
// CapacityInBytes, and then falling back to CapacityInKB if that value is
// unavailable. This helps correct some situations where the former value's
// data gets cleared, which seems to happen on upgrades.
func diskCapacityInGiB(disk *types.VirtualDisk) int {
	if disk.CapacityInBytes > 0 {
		return int(structure.ByteToGiB(disk.CapacityInBytes).(int64))
	}
	log.Printf(
		"[DEBUG] diskCapacityInGiB: capacityInBytes missing for for %s, falling back to capacityInKB",
		object.VirtualDeviceList{}.Name(disk),
	)
	return int(structure.ByteToGiB(disk.CapacityInKB * 1024).(int64))
}
