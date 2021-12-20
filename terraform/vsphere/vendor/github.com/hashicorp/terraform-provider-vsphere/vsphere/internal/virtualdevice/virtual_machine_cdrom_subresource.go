package virtualdevice

import (
	"fmt"
	"log"
	"reflect"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/datastore"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/mitchellh/copystructure"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
)

const vAppTransportIso = "iso"

// CdromSubresourceSchema represents the schema for the cdrom sub-resource.
func CdromSubresourceSchema() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		// VirtualDeviceFileBackingInfo
		"datastore_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The datastore ID the ISO is located on.",
		},
		"path": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The path to the ISO file on the datastore.",
		},
		// VirtualCdromRemoteAtapiBackingInfo
		"client_device": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Indicates whether the device should be mapped to a remote client device",
		},
	}
	structure.MergeSchema(s, subresourceSchema())
	return s
}

// CdromSubresource represents a vsphere_virtual_machine cdrom sub-resource,
// with a complex device lifecycle.
type CdromSubresource struct {
	*Subresource
}

// NewCdromSubresource returns a subresource populated with all of the necessary
// fields.
func NewCdromSubresource(client *govmomi.Client, rdd resourceDataDiff, d, old map[string]interface{}, idx int) *CdromSubresource {
	sr := &CdromSubresource{
		Subresource: &Subresource{
			schema:  CdromSubresourceSchema(),
			client:  client,
			srtype:  subresourceTypeCdrom,
			data:    d,
			olddata: old,
			rdd:     rdd,
		},
	}
	sr.Index = idx
	return sr
}

// CdromApplyOperation processes an apply operation for all disks in the
// resource.
//
// The function takes the root resource's ResourceData, the provider
// connection, and the device list as known to vSphere at the start of this
// operation. All disk operations are carried out, with both the complete,
// updated, VirtualDeviceList, and the complete list of changes returned as a
// slice of BaseVirtualDeviceConfigSpec.
func CdromApplyOperation(d *schema.ResourceData, c *govmomi.Client, l object.VirtualDeviceList) (object.VirtualDeviceList, []types.BaseVirtualDeviceConfigSpec, error) {
	log.Printf("[DEBUG] CdromApplyOperation: Beginning apply operation")
	// While we are currently only restricting CD devices to one device, we have
	// to actually account for the fact that someone could add multiple CD drives
	// out of band. So this workflow is similar to the multi-device workflow that
	// exists for network devices.
	o, n := d.GetChange(subresourceTypeCdrom)
	ods := o.([]interface{})
	nds := n.([]interface{})

	var spec []types.BaseVirtualDeviceConfigSpec

	// Our old and new sets now have an accurate description of devices that may
	// have been added, removed, or changed. Look for removed devices first.
	log.Printf("[DEBUG] CdromApplyOperation: Looking for resources to delete")
nextOld:
	for n, oe := range ods {
		om := oe.(map[string]interface{})
		for _, ne := range nds {
			nm := ne.(map[string]interface{})
			if om["key"] == nm["key"] {
				continue nextOld
			}
		}
		r := NewCdromSubresource(c, d, om, nil, n)
		dspec, err := r.Delete(l)
		if err != nil {
			return nil, nil, fmt.Errorf("%s: %s", r.Addr(), err)
		}
		l = applyDeviceChange(l, dspec)
		spec = append(spec, dspec...)
	}

	// Now check for creates and updates. The results of this operation are
	// committed to state after the operation completes.
	var updates []interface{}
	log.Printf("[DEBUG] CdromApplyOperation: Looking for resources to create or update")
	for n, ne := range nds {
		nm := ne.(map[string]interface{})
		if n < len(ods) {
			// This is an update
			oe := ods[n]
			om := oe.(map[string]interface{})
			if nm["key"] != om["key"] {
				return nil, nil, fmt.Errorf("key mismatch on %s.%d (old: %d, new: %d). This is a bug with the provider, please report it", subresourceTypeCdrom, n, nm["key"].(int), om["key"].(int))
			}
			if reflect.DeepEqual(nm, om) {
				// no change is a no-op
				updates = append(updates, nm)
				log.Printf("[DEBUG] CdromApplyOperation: No-op resource: key %d", nm["key"].(int))
				continue
			}
			r := NewCdromSubresource(c, d, nm, om, n)
			uspec, err := r.Update(l)
			if err != nil {
				return nil, nil, fmt.Errorf("%s: %s", r.Addr(), err)
			}
			l = applyDeviceChange(l, uspec)
			spec = append(spec, uspec...)
			updates = append(updates, r.Data())
			continue
		}
		// New device
		r := NewCdromSubresource(c, d, nm, nil, n)
		cspec, err := r.Create(l)
		if err != nil {
			return nil, nil, fmt.Errorf("%s: %s", r.Addr(), err)
		}
		l = applyDeviceChange(l, cspec)
		spec = append(spec, cspec...)
		updates = append(updates, r.Data())
	}

	log.Printf("[DEBUG] CdromApplyOperation: Post-apply final resource list: %s", subresourceListString(updates))
	// We are now done! Return the updated device list and config spec. Save updates as well.
	if err := d.Set(subresourceTypeCdrom, updates); err != nil {
		return nil, nil, err
	}
	log.Printf("[DEBUG] CdromApplyOperation: Device list at end of operation: %s", DeviceListString(l))
	log.Printf("[DEBUG] CdromApplyOperation: Device config operations from apply: %s", DeviceChangeString(spec))
	log.Printf("[DEBUG] CdromApplyOperation: Apply complete, returning updated spec")
	return l, spec, nil
}

// CdromRefreshOperation processes a refresh operation for all of the disks in
// the resource.
//
// This functions similar to CdromApplyOperation, but nothing to change is
// returned, all necessary values are just set and committed to state.
func CdromRefreshOperation(d *schema.ResourceData, c *govmomi.Client, l object.VirtualDeviceList) error {
	log.Printf("[DEBUG] CdromRefreshOperation: Beginning refresh")
	// While we are currently only restricting CD devices to one device, we have
	// to actually account for the fact that someone could add multiple CD drives
	// out of band. So this workflow is similar to the multi-device workflow that
	// exists for network devices.
	devices := l.Select(func(device types.BaseVirtualDevice) bool {
		if _, ok := device.(*types.VirtualCdrom); ok {
			return true
		}
		return false
	})
	log.Printf("[DEBUG] CdromRefreshOperation: CDROM devices located: %s", DeviceListString(devices))
	curSet := d.Get(subresourceTypeCdrom).([]interface{})
	log.Printf("[DEBUG] CdromRefreshOperation: Current resource set from state: %s", subresourceListString(curSet))
	var newSet []interface{}
	// First check for negative keys. These are freshly added devices that are
	// usually coming into read post-create.
	//
	// If we find what we are looking for, we remove the device from the working
	// set so that we don't try and process it in the next few passes.
	log.Printf("[DEBUG] CdromRefreshOperation: Looking for freshly-created resources to read in")
	for n, item := range curSet {
		m := item.(map[string]interface{})
		if m["key"].(int) < 1 {
			r := NewCdromSubresource(c, d, m, nil, n)
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
	log.Printf("[DEBUG] CdromRefreshOperation: CDROM devices after freshly-created device search: %s", DeviceListString(devices))
	log.Printf("[DEBUG] CdromRefreshOperation: Resource set to write after freshly-created device search: %s", subresourceListString(newSet))

	// Go over the remaining devices, refresh via key, and then remove their
	// entries as well.
	log.Printf("[DEBUG] CdromRefreshOperation: Looking for devices known in state")
	for i := 0; i < len(devices); i++ {
		device := devices[i]
		for n, item := range curSet {
			m := item.(map[string]interface{})
			if m["key"].(int) < 0 {
				// Skip any of these keys as we won't be matching any of those anyway here
				continue
			}
			if device.GetVirtualDevice().Key != int32(m["key"].(int)) {
				// Skip any device that doesn't match key as well
				continue
			}
			// We should have our device -> resource match, so read now.
			r := NewCdromSubresource(c, d, m, nil, n)
			vApp, err := verifyVAppCdromIso(d, device.(*types.VirtualCdrom), l, c)
			if err != nil {
				return err
			}
			if vApp == true && r.Get("client_device") == true {
				log.Printf("[DEBUG] CdromRefreshOperation: %s: Skipping read since CDROM is in use for vApp ISO transport", r)
				// Set the CDROM properties to match a client device so there won't be a diff.
				r.Set("client_device", true)
				r.Set("datastore_id", "")
				r.Set("path", "")
			} else {
				if err := r.Read(l); err != nil {
					return fmt.Errorf("%s: %s", r.Addr(), err)
				}
			}
			// Done reading, push this onto our new set and remove the device from
			// the list
			newSet = append(newSet, r.Data())
			devices = append(devices[:i], devices[i+1:]...)
			i--
		}
	}
	log.Printf("[DEBUG] CdromRefreshOperation: Resource set to write after known device search: %s", subresourceListString(newSet))
	log.Printf("[DEBUG] CdromRefreshOperation: Probable orphaned CDROM devices: %s", DeviceListString(devices))

	// Finally, any device that is still here is orphaned. They should be added
	// as new devices.
	for n, device := range devices {
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
		r := NewCdromSubresource(c, d, m, nil, n)
		if err := r.Read(l); err != nil {
			return fmt.Errorf("%s: %s", r.Addr(), err)
		}
		newSet = append(newSet, r.Data())
	}

	log.Printf("[DEBUG] CdromRefreshOperation: Resource set to write after adding orphaned devices: %s", subresourceListString(newSet))
	log.Printf("[DEBUG] CdromRefreshOperation: Refresh operation complete, sending new resource set")
	return d.Set(subresourceTypeCdrom, newSet)
}

// CdromPostCloneOperation normalizes CDROM devices on a freshly-cloned virtual
// machine and outputs any necessary device change operations. It also sets the
// state in advance of the post-create read.
//
// This differs from a regular apply operation in that a configuration is
// already present, but we don't have any existing state, which the standard
// virtual device operations rely pretty heavily on.
func CdromPostCloneOperation(d *schema.ResourceData, c *govmomi.Client, l object.VirtualDeviceList) (object.VirtualDeviceList, []types.BaseVirtualDeviceConfigSpec, error) {
	log.Printf("[DEBUG] CdromPostCloneOperation: Looking for post-clone device changes")
	// While we are currently only restricting CD devices to one device, we have
	// to actually account for the fact that someone could add multiple CD drives
	// out of band. So this workflow is similar to the multi-device workflow that
	// exists for network devices.
	devices := l.Select(func(device types.BaseVirtualDevice) bool {
		if _, ok := device.(*types.VirtualCdrom); ok {
			return true
		}
		return false
	})
	log.Printf("[DEBUG] CdromPostCloneOperation: CDROM devices located: %s", DeviceListString(devices))
	curSet := d.Get(subresourceTypeCdrom).([]interface{})
	log.Printf("[DEBUG] CdromPostCloneOperation: Current resource set from configuration: %s", subresourceListString(curSet))
	var srcSet []interface{}

	// Populate the source set as if the devices were orphaned. This give us a
	// base to diff off of.
	log.Printf("[DEBUG] CdromPostCloneOperation: Reading existing devices")
	for n, device := range devices {
		m := make(map[string]interface{})
		vd := device.GetVirtualDevice()
		ctlr := l.FindByKey(vd.ControllerKey)
		if ctlr == nil {
			return nil, nil, fmt.Errorf("could not find controller with key %d", vd.Key)
		}
		m["key"] = int(vd.Key)
		var err error
		m["device_address"], err = computeDevAddr(vd, ctlr.(types.BaseVirtualController))
		if err != nil {
			return nil, nil, fmt.Errorf("error computing device address: %s", err)
		}
		r := NewCdromSubresource(c, d, m, nil, n)
		if err := r.Read(l); err != nil {
			return nil, nil, fmt.Errorf("%s: %s", r.Addr(), err)
		}
		srcSet = append(srcSet, r.Data())
	}

	// Now go over our current set, kind of treating it like an apply:
	//
	// * Device past the boundaries of existing devices are created
	// * Devices within the bounds are changed changed
	// * Data at the source with the same data after patching config data is a
	// no-op, but we still push the device's state
	var spec []types.BaseVirtualDeviceConfigSpec
	var updates []interface{}
	for i, ci := range curSet {
		cm := ci.(map[string]interface{})
		if i > len(srcSet)-1 {
			// New device
			r := NewCdromSubresource(c, d, cm, nil, i)
			cspec, err := r.Create(l)
			if err != nil {
				return nil, nil, fmt.Errorf("%s: %s", r.Addr(), err)
			}
			l = applyDeviceChange(l, cspec)
			spec = append(spec, cspec...)
			updates = append(updates, r.Data())
			continue
		}
		sm := srcSet[i].(map[string]interface{})
		nm, err := copystructure.Copy(sm)
		if err != nil {
			return nil, nil, fmt.Errorf("error copying source CDROM device state data at index %d: %s", i, err)
		}
		for k, v := range cm {
			// Skip key and device_address here
			switch k {
			case "key", "device_address":
				continue
			}
			nm.(map[string]interface{})[k] = v
		}
		r := NewCdromSubresource(c, d, nm.(map[string]interface{}), sm, i)
		if !reflect.DeepEqual(sm, nm) {
			// Update
			cspec, err := r.Update(l)
			if err != nil {
				return nil, nil, fmt.Errorf("%s: %s", r.Addr(), err)
			}
			l = applyDeviceChange(l, cspec)
			spec = append(spec, cspec...)
		}
		updates = append(updates, r.Data())
	}

	// Any other device past the end of the CDROM devices listed in config needs
	// to be removed.
	if len(curSet) < len(srcSet) {
		for i, si := range srcSet[len(curSet):] {
			sm := si.(map[string]interface{})
			r := NewCdromSubresource(c, d, sm, nil, i+len(curSet))
			dspec, err := r.Delete(l)
			if err != nil {
				return nil, nil, fmt.Errorf("%s: %s", r.Addr(), err)
			}
			l = applyDeviceChange(l, dspec)
			spec = append(spec, dspec...)
		}
	}

	log.Printf("[DEBUG] CdromPostCloneOperation: Post-clone final resource list: %s", subresourceListString(updates))
	// We are now done! Return the updated device list and config spec. Save updates as well.
	if err := d.Set(subresourceTypeCdrom, updates); err != nil {
		return nil, nil, err
	}
	log.Printf("[DEBUG] CdromPostCloneOperation: Device list at end of operation: %s", DeviceListString(l))
	log.Printf("[DEBUG] CdromPostCloneOperation: Device config operations from post-clone: %s", DeviceChangeString(spec))
	log.Printf("[DEBUG] CdromPostCloneOperation: Operation complete, returning updated spec")
	return l, spec, nil
}

// CdromDiffOperation performs operations relevant to managing the
// diff on cdrom sub-resources
func CdromDiffOperation(d *schema.ResourceDiff, c *govmomi.Client) error {
	log.Printf("[DEBUG] CdromDiffOperation: Beginning diff validation")
	cr := d.Get(subresourceTypeCdrom)
	for ci, ce := range cr.([]interface{}) {
		cm := ce.(map[string]interface{})
		r := NewCdromSubresource(c, d, cm, nil, ci)
		if !structure.ValuesAvailable(fmt.Sprintf("%s.%d.", subresourceTypeCdrom, ci), []string{"datastore_id", "path"}, d) {
			log.Printf("[DEBUG] CdromDiffOperation: Cdrom contains a value that depends on a computed value from another resource. Skipping validation")
			return nil
		}
		if err := r.ValidateDiff(); err != nil {
			return fmt.Errorf("%s: %s", r.Addr(), err)
		}
	}
	log.Printf("[DEBUG] CdromDiffOperation: Diff validation complete")
	return nil
}

// ValidateDiff performs any complex validation of an individual
// cdrom sub-resource that can't be done in schema alone.
func (r *CdromSubresource) ValidateDiff() error {
	log.Printf("[DEBUG] %s: Beginning CDROM configuration validation", r)
	dsID := r.Get("datastore_id").(string)
	path := r.Get("path").(string)
	clientDevice := r.Get("client_device").(bool)
	switch {
	case clientDevice && (dsID != "" || path != ""):
		return fmt.Errorf("Cannot have both client_device parameter and ISO file parameters (datastore_id, path) set")
	case !clientDevice && (dsID == "" || path == ""):
		return fmt.Errorf("Either client_device or datastore_id and path must be set")
	}
	log.Printf("[DEBUG] %s: Config validation complete", r)
	return nil
}

// Create creates a vsphere_virtual_machine cdrom sub-resource.
func (r *CdromSubresource) Create(l object.VirtualDeviceList) ([]types.BaseVirtualDeviceConfigSpec, error) {
	log.Printf("[DEBUG] %s: Running create", r)
	var spec []types.BaseVirtualDeviceConfigSpec
	var ctlr types.BaseVirtualController
	ctlr, err := r.ControllerForCreateUpdate(l, SubresourceControllerTypeIDE, 0)
	if err != nil {
		return nil, err
	}

	// We now have the controller on which we can create our device on.
	device, err := l.CreateCdrom(ctlr.(*types.VirtualIDEController))
	if err != nil {
		return nil, err
	}
	// Map the CDROM to the correct device
	r.mapCdrom(device, l)
	// Done here. Save IDs, push the device to the new device list and return.
	if err := r.SaveDevIDs(device, ctlr); err != nil {
		return nil, err
	}
	dspec, err := object.VirtualDeviceList{device}.ConfigSpec(types.VirtualDeviceConfigSpecOperationAdd)
	if err != nil {
		return nil, err
	}
	spec = append(spec, dspec...)
	log.Printf("[DEBUG] %s: Device config operations from create: %s", r, DeviceChangeString(spec))
	log.Printf("[DEBUG] %s: Create finished", r)
	return spec, nil
}

// Read reads a vsphere_virtual_machine cdrom sub-resource.
func (r *CdromSubresource) Read(l object.VirtualDeviceList) error {
	log.Printf("[DEBUG] %s: Reading state", r)
	d, err := r.FindVirtualDevice(l)
	if err != nil {
		return fmt.Errorf("cannot find disk device: %s", err)
	}
	device, ok := d.(*types.VirtualCdrom)
	if !ok {
		return fmt.Errorf("device at %q is not a virtual CDROM device", l.Name(d))
	}
	// Only read backing info if it's available.
	switch backing := device.Backing.(type) {
	case *types.VirtualCdromRemoteAtapiBackingInfo:
		r.Set("client_device", true)
	case *types.VirtualCdromIsoBackingInfo:
		dp := &object.DatastorePath{}
		if ok := dp.FromString(backing.FileName); !ok {
			return fmt.Errorf("could not read datastore path in backing %q", backing.FileName)
		}
		// If a vApp ISO was inserted, it will be removed if the VM is powered off
		// and cause backing.Datastore to be nil.
		if backing.Datastore != nil {
			r.Set("datastore_id", backing.Datastore.Value)
		}
		r.Set("path", dp.Path)
	default:
		// This is an unsupported entry, so we clear all attributes in the
		// subresource (except for the device address and key, of course).  In
		// addition to making sure correct diffs get created for these kinds of
		// devices, this ensures we don't fail on CDROM device types we don't
		// support right now, such as passthrough devices. We might support these
		// later.
		log.Printf("%s: [DEBUG] Unknown CDROM type %T, clearing all attributes", r, backing)
		r.Set("datastore_id", "")
		r.Set("path", "")
		r.Set("client_device", false)
	}
	// Save the device key and address data
	ctlr, err := findControllerForDevice(l, d)
	if err != nil {
		return err
	}
	if err := r.SaveDevIDs(d, ctlr); err != nil {
		return err
	}
	log.Printf("[DEBUG] %s: Read finished (key and device address may have changed)", r)
	return nil
}

// Update updates a vsphere_virtual_machine cdrom sub-resource.
func (r *CdromSubresource) Update(l object.VirtualDeviceList) ([]types.BaseVirtualDeviceConfigSpec, error) {
	log.Printf("[DEBUG] %s: Beginning update", r)
	d, err := r.FindVirtualDevice(l)
	if err != nil {
		return nil, fmt.Errorf("cannot find disk device: %s", err)
	}
	device, ok := d.(*types.VirtualCdrom)
	if !ok {
		return nil, fmt.Errorf("device at %q is not a virtual CDROM device", l.Name(d))
	}

	// Map the CDROM to the correct device
	r.mapCdrom(device, l)
	spec, err := object.VirtualDeviceList{device}.ConfigSpec(types.VirtualDeviceConfigSpecOperationEdit)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Device config operations from update: %s", r, DeviceChangeString(spec))
	log.Printf("[DEBUG] %s: Update complete", r)
	return spec, nil
}

// Delete deletes a vsphere_virtual_machine cdrom sub-resource.
func (r *CdromSubresource) Delete(l object.VirtualDeviceList) ([]types.BaseVirtualDeviceConfigSpec, error) {
	log.Printf("[DEBUG] %s: Beginning delete", r)
	d, err := r.FindVirtualDevice(l)
	if err != nil {
		return nil, fmt.Errorf("cannot find disk device: %s", err)
	}
	device, ok := d.(*types.VirtualCdrom)
	if !ok {
		return nil, fmt.Errorf("device at %q is not a virtual CDROM device", l.Name(d))
	}
	deleteSpec, err := object.VirtualDeviceList{device}.ConfigSpec(types.VirtualDeviceConfigSpecOperationRemove)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Device config operations from update: %s", r, DeviceChangeString(deleteSpec))
	log.Printf("[DEBUG] %s: Delete completed", r)
	return deleteSpec, nil
}

// mapCdrom takes a CdromSubresource and attaches either a client device or a datastore ISO.
func (r *CdromSubresource) mapCdrom(device *types.VirtualCdrom, l object.VirtualDeviceList) error {
	dsID := r.Get("datastore_id").(string)
	path := r.Get("path").(string)
	clientDevice := r.Get("client_device").(bool)
	switch {
	case dsID != "" && path != "":
		// If the datastore ID and path are both set, the CDROM will be mapped to a file on a datastore.
		ds, err := datastore.FromID(r.client, dsID)
		if err != nil {
			return fmt.Errorf("cannot find datastore: %s", err)
		}
		dsProps, err := datastore.Properties(ds)
		if err != nil {
			return fmt.Errorf("could not get properties for datastore: %s", err)
		}
		dsName := dsProps.Name
		dsPath := &object.DatastorePath{
			Datastore: dsName,
			Path:      path,
		}
		device = l.InsertIso(device, dsPath.String())
		l.Connect(device)
		return nil
	case clientDevice == true:
		// If set to use the client device, then the CDROM will be mapped to a remote device.
		device.Backing = &types.VirtualCdromRemoteAtapiBackingInfo{
			VirtualDeviceRemoteDeviceBackingInfo: types.VirtualDeviceRemoteDeviceBackingInfo{},
		}
		return nil
	}
	panic(fmt.Sprintf("%s: no CDROM types specified", r))
}

// VerifyVAppTransport validates that all the required components are included in
// the virtual machine configuration if vApp properties are set.
func VerifyVAppTransport(d *schema.ResourceDiff, c *govmomi.Client) error {
	log.Printf("[DEBUG] VAppDiffOperation: Verifying configuration meets requirements for vApp transport")
	// Check if there is a client CDROM device configured.
	cl := d.Get("cdrom")
	for _, c := range cl.([]interface{}) {
		if c.(map[string]interface{})["client_device"].(bool) == true {
			// There is a device configured that can support vApp ISO transport if needed
			log.Printf("[DEBUG] VAppDiffOperation: Client CDROM device exists which can support ISO transport")
			return nil
		}
	}
	// Iterate over each transport and see if ISO transport is supported.
	tm := d.Get("vapp_transport").([]interface{})
	for _, m := range tm {
		if m.(string) == vAppTransportIso && len(tm) == 1 {
			return fmt.Errorf("this virtual machine requires a client CDROM device to deliver vApp properties")
		}
	}
	log.Printf("[DEBUG] VAppDiffOperation: ISO transport is not supported on this virtual machine or multiple transport options exist")
	return nil
}

// verifyVAppCdromIso takes VirtualCdrom and determines if it is needed for
// vApp ISO transport. It does this by first checking if it has an ISO inserted
// that matches the vApp ISO naming pattern. If it does, then the next step is
// to see if vApp ISO transport is supported on the VM. If both of those
// conditions are met, then the CDROM is considered in use for vApp transport.
func verifyVAppCdromIso(d *schema.ResourceData, device *types.VirtualCdrom, l object.VirtualDeviceList, c *govmomi.Client) (bool, error) {
	log.Printf("[DEBUG] IsVAppCdrom: Checking if CDROM is using a vApp ISO")
	// If the CDROM is using VirtualCdromIsoBackingInfo and matches the ISO
	// naming pattern, it has been used as a vApp CDROM, and we can move on to
	// checking if the parent VM supports ISO transport.
	if backing, ok := device.Backing.(*types.VirtualCdromIsoBackingInfo); ok {
		dp := &object.DatastorePath{}
		if ok := dp.FromString(backing.FileName); !ok {
			// If the ISO path can not be read, we can't tell if a vApp ISO is
			// connected.
			log.Printf("[DEBUG] IsVAppCdrom: Cannot read ISO path, cannot determine if CDROM is used for vApp")
			return false, nil
		}
		// The pattern used for vApp ISO naming is
		// "<vmname>/_ovfenv-<vmname>.iso"
		re := regexp.MustCompile(".*/_ovfenv-.*.iso")
		if !re.MatchString(dp.Path) {
			log.Printf("[DEBUG] IsVAppCdrom: ISO is name does not match vApp ISO naming pattern (<vmname>/_ovfenv-<vmname>.iso): %s", dp.Path)
			return false, nil
		}
	} else {
		// vApp CDROMs must be backed by an ISO.
		log.Printf("[DEBUG] IsVAppCdrom: CDROM is not backed by an ISO")
		return false, nil
	}
	log.Printf("[DEBUG] IsVAppCdrom: CDROM has a vApp ISO inserted")
	// Set the vApp transport methods
	tm := d.Get("vapp_transport").([]interface{})
	for _, t := range tm {
		if t.(string) == "iso" {
			log.Printf("[DEBUG] IsVAppCdrom: vApp ISO transport is supported")
			return true, nil
		}
	}
	log.Printf("[DEBUG] IsVAppCdrom: vApp ISO transport is not required")
	return false, nil
}
