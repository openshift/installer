package virtualdevice

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/computeresource"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/hostsystem"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/virtualmachine"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/mitchellh/copystructure"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
)

const (
	subresourceTypeDisk             = "disk"
	subresourceTypeNetworkInterface = "network_interface"
	subresourceTypeCdrom            = "cdrom"
)

const (
	// SubresourceControllerTypeIDE is a string representation of IDE controller
	// classes.
	SubresourceControllerTypeIDE = "ide"

	// SubresourceControllerTypeSATA is a string representation of SATA controller
	// classes.
	SubresourceControllerTypeSATA = "sata"

	// SubresourceControllerTypeSCSI is a string representation of all SCSI
	// controller types.
	//
	// This is mainly used when computing IDs so that we can use a more general
	// device search.
	SubresourceControllerTypeSCSI = "scsi"

	// SubresourceControllerTypeParaVirtual is a string representation of the
	// VMware PV SCSI controller type.
	SubresourceControllerTypeParaVirtual = "pvscsi"

	// SubresourceControllerTypeLsiLogic is a string representation of the
	// LSI Logic parallel virtual SCSI controller type.
	SubresourceControllerTypeLsiLogic = "lsilogic"

	// SubresourceControllerTypeLsiLogicSAS is a string representation of the
	// LSI Logic SAS virtual SCSI controller type.
	SubresourceControllerTypeLsiLogicSAS = "lsilogic-sas"

	// SubresourceControllerTypePCI is a string representation of PCI controller
	// classes.
	SubresourceControllerTypePCI = "pci"
)

const (
	subresourceControllerTypeMixed   = "mixed"
	subresourceControllerTypeUnknown = "unknown"
)

const (
	subresourceControllerSharingMixed   = "mixed"
	subresourceControllerSharingUnknown = "unknown"
)

var subresourceIDControllerTypeAllowedValues = []string{
	SubresourceControllerTypeIDE,
	SubresourceControllerTypeSCSI,
	SubresourceControllerTypePCI,
	SubresourceControllerTypeSATA,
}

var sharesLevelAllowedValues = []string{
	string(types.SharesLevelLow),
	string(types.SharesLevelNormal),
	string(types.SharesLevelHigh),
	string(types.SharesLevelCustom),
}

// SCSIBusSharingAllowedValues exports the list of supported SCSI bus sharing
// modes. These are the only modes that can be specified for scsi_bus_sharing
// and should be checked in a ValidateFunc.
var SCSIBusSharingAllowedValues = []string{
	string(types.VirtualSCSISharingNoSharing),
	string(types.VirtualSCSISharingPhysicalSharing),
	string(types.VirtualSCSISharingVirtualSharing),
}

// SCSIBusTypeAllowedValues exports the currently list of SCSI controller types
// that we support in the resource. The user is only allowed to select a type
// in this list, which should be used in a ValidateFunc on the appropriate
// field.
var SCSIBusTypeAllowedValues = []string{
	SubresourceControllerTypeParaVirtual,
	SubresourceControllerTypeLsiLogic,
	SubresourceControllerTypeLsiLogicSAS,
}

// newSubresourceFunc is a method signature for the wrapper methods that create
// a new instance of a specific subresource  that is derived from the base
// subresoruce object. It's used in the general apply and read operation
// methods, which themselves are called usually from higher-level apply
// functions for virtual devices.
type newSubresourceFunc func(*govmomi.Client, int, int, *schema.ResourceData) SubresourceInstance

// SubresourceInstance is an interface for derivative objects of Subresource.
// It's used on the general apply and read operation methods, and contains both
// exported methods of the base Subresource type and the CRUD methods that
// should be supplied by derivative objects.
//
// Note that this interface should be used sparingly - as such, only the
// methods that are needed by inparticular functions external to most virtual
// device workflows are exported into this interface.
type SubresourceInstance interface {
	Create(object.VirtualDeviceList) ([]types.BaseVirtualDeviceConfigSpec, error)
	Read(object.VirtualDeviceList) error
	Update(object.VirtualDeviceList) ([]types.BaseVirtualDeviceConfigSpec, error)
	Delete(object.VirtualDeviceList) ([]types.BaseVirtualDeviceConfigSpec, error)

	DevAddr() string
	Addr() string
	Set(string, interface{}) error
	Schema() map[string]*schema.Schema
	State() map[string]interface{}
}

// pciApplyConfig is used for making PCI device change functions easier to
// work with.
type pciApplyConfig struct {
	Client        *govmomi.Client
	ResourceData  *schema.ResourceData
	SystemId      string
	Spec          []types.BaseVirtualDeviceConfigSpec
	VirtualDevice object.VirtualDeviceList
}

// controllerTypeToClass converts a controller type to a specific short-form
// controller class, namely for use with working with IDs.
func controllerTypeToClass(c types.BaseVirtualController) (string, error) {
	var t string
	switch c.(type) {
	case *types.VirtualIDEController:
		t = SubresourceControllerTypeIDE
	case *types.VirtualAHCIController:
		t = SubresourceControllerTypeSATA
	case *types.VirtualPCIController:
		t = SubresourceControllerTypePCI
	case *types.ParaVirtualSCSIController, *types.VirtualBusLogicController,
		*types.VirtualLsiLogicController, *types.VirtualLsiLogicSASController:
		t = SubresourceControllerTypeSCSI
	default:
		return subresourceControllerTypeUnknown, fmt.Errorf("unsupported controller type %T", c)
	}
	return t, nil
}

// resourceDataDiff is an interface comprised of functions common to both
// ResourceData and ResourceDiff.
//
// During any inparticular CRUD or diff alteration call, either one of
// ResourceData or ResourceDiff will be available. Both will never be available
// at the same time. Having these underlying values exposed directly presents a
// potentially unsafe API where one of them will be nil at any given time.
// Having this as an interface allows common behavior to be exposed directly,
// while still offering the ability to type assert in certain situations.
//
// This is not an exhaustive list of methods - any missing ones should be added
// as needed.
type resourceDataDiff interface {
	Id() string
	Get(string) interface{}
	HasChange(string) bool
}

// Subresource defines a common interface for device sub-resources in the
// vsphere_virtual_machine resource.
//
// This object is designed to be used by parts of the resource with workflows
// that are so complex in their own right that probably the only way to handle
// their management is to treat them like resources themselves.
//
// This structure of this resource loosely follows schema.Resource with having
// CRUD and maintaining a set of resource data to work off of. However, since
// we are using schema.Resource, we take some liberties that we normally would
// not be able to take, or need to take considering the context of the data we
// are working with.
//
// Inparticular functions implement this structure by creating an instance into
// it, much like how a resource creates itself by creating an instance of
// schema.Resource.
type Subresource struct {
	// The index of this subresource - should either be an index or hash. It's up
	// to the upstream object to set this to something useful.
	Index int

	// The resource schema. This is an internal field as we build on this field
	// later on with common keys for all subresources, namely the internal ID.
	schema map[string]*schema.Schema

	// The client connection.
	client *govmomi.Client

	// The subresource type. This should match the key that the subresource is
	// named in the schema, such as "disk" or "network_interface".
	srtype string

	// The resource data - this should be loaded when the resource is created.
	data map[string]interface{}

	// The old resource data, if it exists.
	olddata map[string]interface{}

	// Either a root-level ResourceData or ResourceDiff. The one that is
	// specifically present will depend on the context the Subresource is being
	// used in.
	rdd resourceDataDiff
}

// subresourceSchema is a map[string]*schema.Schema of common schema fields.
// This includes the internal_id field, which is used as a unique ID for the
// lifecycle of this resource.
func subresourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "The unique device ID for this device within its virtual machine.",
		},
		"device_address": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The internally-computed address of this device, such as scsi:0:1, denoting scsi bus #0 and device unit 1.",
		},
	}
}

// Addr returns the resource address for this subresource.
func (r *Subresource) Addr() string {
	return fmt.Sprintf("%s.%d", r.srtype, r.Index)
}

// Get hands off to r.data.Get, with an address relative to this subresource.
func (r *Subresource) Get(key string) interface{} {
	return r.data[key]
}

// Set sets the specified key/value pair in the subresource.
func (r *Subresource) Set(key string, value interface{}) {
	if v := structure.NormalizeValue(value); v != nil {
		r.data[key] = v
	}
}

// HasChange checks to see if there has been a change in the resource data
// since the last update.
//
// Note that this operation may only be useful during update operations,
// depending on subresource-specific workflow.
func (r *Subresource) HasChange(key string) bool {
	o, n := r.GetChange(key)
	return !reflect.DeepEqual(o, n)
}

// GetChange gets the old and new values for the value specified by key.
func (r *Subresource) GetChange(key string) (interface{}, interface{}) {
	new := r.data[key]
	// No old data means no change,  so we use the new value as a placeholder.
	old := r.data[key]
	if r.olddata != nil {
		old = r.olddata[key]
	}
	return old, new
}

// GetWithRestart checks to see if a field has been modified, returns the new
// value, and sets restart if it has changed.
func (r *Subresource) GetWithRestart(key string) interface{} {
	if r.HasChange(key) {
		r.SetRestart(key)
	}
	return r.Get(key)
}

// GetWithVeto returns the value specified by key, but returns an error if it
// has changed. The intention here is to block changes to the resource in a
// fashion that would otherwise result in forcing a new resource.
func (r *Subresource) GetWithVeto(key string) (interface{}, error) {
	if r.HasChange(key) {
		old, new := r.GetChange(key)
		return r.Get(key), fmt.Errorf("cannot change the value of %q - (old: %v new: %v)", key, old, new)
	}
	return r.Get(key), nil
}

// SetRestart sets reboot_required in the global ResourceData. The key is only
// required for logging.
func (r *Subresource) SetRestart(key string) {
	log.Printf("[DEBUG] %s: Resource argument %q requires a VM restart", r, key)
	switch d := r.rdd.(type) {
	case *schema.ResourceData:
		d.Set("reboot_required", true)
	case *schema.ResourceDiff:
		d.SetNew("reboot_required", true)
	default:
		// This should never happen, but log if it does.
		log.Printf("[WARN] %s: Could not flag reboot_required: invalid type %T", r, r.rdd)
	}
}

// Data returns the underlying data map.
func (r *Subresource) Data() map[string]interface{} {
	return r.data
}

// Hash calculates a set hash for the current data. If you want a hash for
// error reporting a device address, it's probably a good idea to run this at
// the beginning of a run as any set calls will change the value this
// ultimately calculates.
func (r *Subresource) Hash() int {
	hf := schema.HashResource(&schema.Resource{Schema: r.schema})
	return hf(r.data)
}

// computeDevAddr handles the logic for SaveDevIDs and allows it to be used
// outside of a subresource.
func computeDevAddr(device types.BaseVirtualDevice, ctlr types.BaseVirtualController) (string, error) {
	vd := device.GetVirtualDevice()
	vc := ctlr.GetVirtualController()
	ctype, err := controllerTypeToClass(ctlr)
	if err != nil {
		return "", err
	}
	parts := []string{
		ctype,
		strconv.Itoa(int(vc.BusNumber)),
		strconv.Itoa(int(structure.DeRef(vd.UnitNumber).(int32))),
	}
	return strings.Join(parts, ":"), nil
}

// SaveDevIDs saves the device's current key, and also the device_address. The
// latter is a computed schema field that contains the controller type, the
// controller's bus number, and the device's unit number on that controller.
// This helps locate the device when the key is in flux (such as when devices
// are just being created).
func (r *Subresource) SaveDevIDs(device types.BaseVirtualDevice, ctlr types.BaseVirtualController) error {
	r.Set("key", device.GetVirtualDevice().Key)
	addr, err := computeDevAddr(device, ctlr)
	if err != nil {
		return err
	}
	r.Set("device_address", addr)
	return nil
}

// DevAddr returns the device_address attribute in the subresource. This
// function exists mainly as a functional counterpart to SaveDevIDs.
func (r *Subresource) DevAddr() string {
	return r.Get("device_address").(string)
}

// splitDevAddr splits an device addres into its inparticular parts and asserts
// that we have all the correct data.
func splitDevAddr(id string) (string, int, int, error) {
	parts := strings.Split(id, ":")
	if len(parts) < 3 {
		return "", 0, 0, fmt.Errorf("invalid ID %q", id)
	}
	ct, cbs, dus := parts[0], parts[1], parts[2]
	cb, cbe := strconv.Atoi(cbs)
	du, due := strconv.Atoi(dus)
	var found bool
	for _, v := range subresourceIDControllerTypeAllowedValues {
		if v == ct {
			found = true
		}
	}
	if !found {
		return ct, cb, du, fmt.Errorf("invalid controller type %q found in ID", ct)
	}
	if cbe != nil {
		return ct, cb, du, fmt.Errorf("invalid bus number %q found in ID", cbs)
	}
	if due != nil {
		return ct, cb, du, fmt.Errorf("invalid disk unit number %q found in ID", dus)
	}
	return ct, cb, du, nil
}

// findVirtualDeviceInListControllerSelectFunc returns a function that can be
// used with VirtualDeviceList.Select to locate a controller device based on
// the criteria that we have laid out.
func findVirtualDeviceInListControllerSelectFunc(ct string, cb int) func(types.BaseVirtualDevice) bool {
	return func(device types.BaseVirtualDevice) bool {
		switch ct {
		case SubresourceControllerTypeIDE:
			if _, ok := device.(*types.VirtualIDEController); !ok {
				return false
			}
		case SubresourceControllerTypeSATA:
			if _, ok := device.(*types.VirtualAHCIController); !ok {
				return false
			}
		case SubresourceControllerTypeSCSI:
			if _, ok := device.(types.BaseVirtualSCSIController); !ok {
				return false
			}
		case SubresourceControllerTypePCI:
			if _, ok := device.(*types.VirtualPCIController); !ok {
				return false
			}
		}
		vc := device.(types.BaseVirtualController).GetVirtualController()
		if vc.BusNumber == int32(cb) {
			return true
		}
		return false
	}
}

// findVirtualDeviceInListDeviceSelectFunc returns a function that can be used
// with VirtualDeviceList.Select to locate a virtual device based on its
// controller device key, and the unit number on the device.
func findVirtualDeviceInListDeviceSelectFunc(ckey int32, du int) func(types.BaseVirtualDevice) bool {
	return func(d types.BaseVirtualDevice) bool {
		vd := d.GetVirtualDevice()
		if vd.ControllerKey == ckey && vd.UnitNumber != nil && *vd.UnitNumber == int32(du) {
			return true
		}
		return false
	}
}

// findControllerForDevice locates a controller via its virtual device.
func findControllerForDevice(l object.VirtualDeviceList, bvd types.BaseVirtualDevice) (types.BaseVirtualController, error) {
	vd := bvd.GetVirtualDevice()
	ctlr := l.FindByKey(vd.ControllerKey)

	if ctlr == nil {
		return nil, fmt.Errorf("could not find controller key %d for device %d", vd.ControllerKey, vd.Key)
	}

	return ctlr.(types.BaseVirtualController), nil
}

// FindVirtualDeviceByAddr locates the subresource's virtual device in the
// supplied VirtualDeviceList by its device address.
func (r *Subresource) FindVirtualDeviceByAddr(l object.VirtualDeviceList) (types.BaseVirtualDevice, error) {
	log.Printf("[DEBUG] FindVirtualDevice: Looking for device with address %s", r.DevAddr())
	oldAddress, _ := r.GetChange("device_address")
	ct, cb, du, err := splitDevAddr(oldAddress.(string))
	if err != nil {
		return nil, err
	}

	// find the controller
	csf := findVirtualDeviceInListControllerSelectFunc(ct, cb)
	ctlrs := l.Select(csf)
	if len(ctlrs) != 1 {
		return nil, fmt.Errorf("invalid controller result - %d results returned (expected 1): type %q, bus number: %d", len(ctlrs), ct, cb)
	}
	ctlr := ctlrs[0]

	// find the device
	ckey := ctlr.GetVirtualDevice().Key
	dsf := findVirtualDeviceInListDeviceSelectFunc(ckey, du)
	devices := l.Select(dsf)
	if len(devices) != 1 {
		return nil, fmt.Errorf("invalid device result - %d results returned (expected 1): controller key %q, disk number: %d", len(devices), ckey, du)
	}
	device := devices[0]
	log.Printf("[DEBUG] FindVirtualDevice: Device found: %s", l.Name(device))
	return device, nil
}

// FindVirtualDevice will attempt to find an address by its device key if it is
// > 0, otherwise it will attempt to locate it by its device address.
func (r *Subresource) FindVirtualDevice(l object.VirtualDeviceList) (types.BaseVirtualDevice, error) {
	if key := r.Get("key").(int); key > 0 {
		log.Printf("[DEBUG] FindVirtualDevice: Looking for device with key %d", key)
		if dev := l.FindByKey(int32(key)); dev != nil {
			log.Printf("[DEBUG] FindVirtualDevice: Device found: %s", l.Name(dev))
			return dev, nil
		}
		return nil, fmt.Errorf("could not find device with key %d", key)
	}
	return r.FindVirtualDeviceByAddr(l)
}

// String prints out the device sub-resource's information including the ID at
// time of instantiation, the short name of the disk, and the current device
// key and address.
func (r *Subresource) String() string {
	devaddr := r.Get("device_address").(string)
	if devaddr == "" {
		devaddr = "<new device>"
	}
	return fmt.Sprintf("%s (key %d at %s)", r.Addr(), r.Get("key").(int), devaddr)
}

// swapSCSIDevice swaps out the supplied controller for a new one of the
// supplied controller type. Any connected devices are re-connected at the same
// device units on the new device. A list of changes is returned.
func swapSCSIDevice(l object.VirtualDeviceList, device types.BaseVirtualSCSIController, ct string, st string) ([]types.BaseVirtualDeviceConfigSpec, error) {
	log.Printf("[DEBUG] swapSCSIDevice: Swapping SCSI device for one of controller type %s: %s", ct, l.Name(device.(types.BaseVirtualDevice)))
	var spec []types.BaseVirtualDeviceConfigSpec
	bvd := device.(types.BaseVirtualDevice)
	cspec, err := object.VirtualDeviceList{bvd}.ConfigSpec(types.VirtualDeviceConfigSpecOperationRemove)
	if err != nil {
		return nil, err
	}
	spec = append(spec, cspec...)

	nsd, err := l.CreateSCSIController(ct)
	if err != nil {
		return nil, err
	}
	nsd.(types.BaseVirtualSCSIController).GetVirtualSCSIController().SharedBus = types.VirtualSCSISharing(st)
	nsd.(types.BaseVirtualSCSIController).GetVirtualSCSIController().BusNumber = device.GetVirtualSCSIController().BusNumber
	cspec, err = object.VirtualDeviceList{nsd}.ConfigSpec(types.VirtualDeviceConfigSpecOperationAdd)
	if err != nil {
		return nil, err
	}
	spec = append(spec, cspec...)
	ockey := device.GetVirtualSCSIController().Key
	nckey := nsd.GetVirtualDevice().Key
	for _, vd := range l {
		if vd.GetVirtualDevice().ControllerKey == ockey {
			vd.GetVirtualDevice().ControllerKey = nckey
			cspec, err := object.VirtualDeviceList{vd}.ConfigSpec(types.VirtualDeviceConfigSpecOperationEdit)
			if err != nil {
				return nil, err
			}
			if len(cspec) != 1 {
				return nil, fmt.Errorf("incorrect number of config spec items returned - expected 1, got %d", len(cspec))
			}
			// Clear the file operation
			cspec[0].GetVirtualDeviceConfigSpec().FileOperation = ""
			spec = append(spec, cspec...)
		}
	}
	log.Printf("[DEBUG] swapSCSIDevice: Outgoing device config spec: %s", DeviceChangeString(spec))
	return spec, nil
}

// NormalizeBus checks the storage controllers on the virtual machine and
// either creates them if they don't exist, or migrates them to the specified
// controller type. Devices are migrated to the new controller appropriately. A
// spec slice is returned with the changes.
//
// The first number of slots specified by count are normalized by this
// function. Any others are left unchanged.
func NormalizeBus(l object.VirtualDeviceList, d *schema.ResourceData) (object.VirtualDeviceList, []types.BaseVirtualDeviceConfigSpec, error) {
	scsiCount := d.Get("scsi_controller_count").(int)
	scsiType := d.Get("scsi_type").(string)
	scsiSharing := d.Get("scsi_bus_sharing").(string)
	sataCount := d.Get("sata_controller_count").(int)
	ideCount := d.Get("ide_controller_count").(int)
	var spec []types.BaseVirtualDeviceConfigSpec
	scsiCtlrs := make([]types.BaseVirtualSCSIController, scsiCount)
	sataCtlrs := make([]types.BaseVirtualSATAController, sataCount)
	ideCtlrs := make([]*types.VirtualIDEController, ideCount)
	// Don't worry about doing any fancy select stuff here, just go thru the
	// VirtualDeviceList and populate the controllers.
	log.Printf("[DEBUG] NormalizeBus: Normalizing first %d controllers on SCSI bus to device type %s", scsiCount, scsiType)
	log.Printf("[DEBUG] NormalizeBus: Normalizing first %d controllers on SATA bus", sataCount)
	log.Printf("[DEBUG] NormalizeBus: Normalizing first %d controllers on IDE bus", ideCount)
	for _, dev := range l {
		switch ctlr := dev.(type) {
		case types.BaseVirtualSCSIController:
			if busNumber := ctlr.GetVirtualSCSIController().BusNumber; busNumber < int32(scsiCount) {
				scsiCtlrs[busNumber] = ctlr
			}
		case types.BaseVirtualSATAController:
			if busNumber := ctlr.GetVirtualSATAController().BusNumber; busNumber < int32(sataCount) {
				sataCtlrs[busNumber] = ctlr
			}
		case *types.VirtualIDEController:
			if busNumber := ctlr.GetVirtualController().BusNumber; busNumber < int32(ideCount) {
				ideCtlrs[busNumber] = ctlr
			}
		}
	}
	log.Printf("[DEBUG] NormalizeBus: Current SCSI bus contents: %s", scsiControllerListString(scsiCtlrs))
	// Now iterate over the SCSI controllers
	for n, ctlr := range scsiCtlrs {
		if ctlr == nil {
			log.Printf("[DEBUG] NormalizeBus: Creating SCSI controller of type %s at bus number %d", scsiType, n)
			cspec, err := createSCSIController(&l, scsiType, scsiSharing)
			if err != nil {
				return nil, nil, err
			}
			spec = append(spec, cspec...)
			continue
		}
		if l.Type(ctlr.(types.BaseVirtualDevice)) == scsiType {
			cspec, err := setSCSIBusSharing(&l, ctlr, scsiSharing)
			if err != nil {
				return nil, nil, err
			}
			spec = append(spec, cspec...)
			continue
		}
		cspec, err := swapSCSIDevice(l, ctlr, scsiType, scsiSharing)
		if err != nil {
			return nil, nil, err
		}
		spec = append(spec, cspec...)
		l = applyDeviceChange(l, cspec)
		continue
	}
	log.Printf("[DEBUG] NormalizeBus: Current SATA bus contents: %s", sataControllerListString(sataCtlrs))
	// Now iterate over the SATA controllers
	for n, ctlr := range sataCtlrs {
		if ctlr == nil {
			log.Printf("[DEBUG] NormalizeBus: Creating SATA controller at bus number %d", n)
			cspec, err := createSATAController(&l, n)
			if err != nil {
				return nil, nil, err
			}
			spec = append(spec, cspec...)
		}
	}
	log.Printf("[DEBUG] NormalizeBus: Current IDE bus contents: %s", ideControllerListString(ideCtlrs))
	// Now iterate over the IDE controllers
	for n, ctlr := range ideCtlrs {
		if ctlr == nil {
			log.Printf("[DEBUG] NormalizeBus: Creating IDE controller at bus number %d", n)
			cspec, err := createIDEController(&l, n)
			if err != nil {
				return nil, nil, err
			}
			spec = append(spec, cspec...)
		}
	}
	log.Printf("[DEBUG] NormalizeBus: Outgoing device list: %s", DeviceListString(l))
	log.Printf("[DEBUG] NormalizeBus: Outgoing device config spec: %s", DeviceChangeString(spec))
	return l, spec, nil
}

// setSCSIBusSharing takes a BaseVirtualSCSIController, sets the sharing mode,
// and applies that change to the VirtualDeviceList.
func setSCSIBusSharing(l *object.VirtualDeviceList, ctlr types.BaseVirtualSCSIController, st string) ([]types.BaseVirtualDeviceConfigSpec, error) {
	var cspec []types.BaseVirtualDeviceConfigSpec
	if ctlr.GetVirtualSCSIController().SharedBus != types.VirtualSCSISharing(st) {
		ctlr.GetVirtualSCSIController().SharedBus = types.VirtualSCSISharing(st)
		var err error
		cspec, err = object.VirtualDeviceList{ctlr.(types.BaseVirtualDevice)}.ConfigSpec(types.VirtualDeviceConfigSpecOperationEdit)
		if err != nil {
			return nil, err
		}
		*l = applyDeviceChange(*l, cspec)
	}
	return cspec, nil
}

// createIDEController creates a new IDE controller.
func createIDEController(l *object.VirtualDeviceList, bus int) ([]types.BaseVirtualDeviceConfigSpec, error) {
	ide, _ := l.CreateIDEController()
	cspec, err := object.VirtualDeviceList{ide}.ConfigSpec(types.VirtualDeviceConfigSpecOperationAdd)
	*l = applyDeviceChange(*l, cspec)
	return cspec, err
}

// createSATAController creates a new SATA controller.
func createSATAController(l *object.VirtualDeviceList, bus int) ([]types.BaseVirtualDeviceConfigSpec, error) {
	ahci := &types.VirtualAHCIController{}
	ahci.Key = l.NewKey()
	ahci.BusNumber = int32(bus)
	cspec, err := object.VirtualDeviceList{ahci}.ConfigSpec(types.VirtualDeviceConfigSpecOperationAdd)
	*l = applyDeviceChange(*l, cspec)
	return cspec, err
}

// createSCSIController creates a new SCSI controller of the specified type and
// sharing mode.
func createSCSIController(l *object.VirtualDeviceList, ct string, st string) ([]types.BaseVirtualDeviceConfigSpec, error) {
	nc, err := l.CreateSCSIController(ct)
	if err != nil {
		return nil, err
	}
	nc.(types.BaseVirtualSCSIController).GetVirtualSCSIController().SharedBus = types.VirtualSCSISharing(st)
	cspec, err := object.VirtualDeviceList{nc}.ConfigSpec(types.VirtualDeviceConfigSpecOperationAdd)
	*l = applyDeviceChange(*l, cspec)
	return cspec, err
}

// ReadSCSIBusType checks the SCSI bus state and returns a device type
// depending on if all controllers are one specific kind or not. Only the first
// number of controllers specified by count are checked.
func ReadSCSIBusType(l object.VirtualDeviceList, count int) string {
	controllers := make([]types.BaseVirtualSCSIController, count)
	for _, dev := range l {
		if sc, ok := dev.(types.BaseVirtualSCSIController); ok && sc.GetVirtualSCSIController().BusNumber < int32(count) {
			controllers[sc.GetVirtualSCSIController().BusNumber] = sc
		}
	}
	log.Printf("[DEBUG] ReadSCSIBusType: SCSI controller layout for first %d controllers: %s", count, scsiControllerListString(controllers))
	if len(controllers) == 0 || controllers[0] == nil {
		return subresourceControllerTypeUnknown
	}
	last := l.Type(controllers[0].(types.BaseVirtualDevice))
	for _, ctlr := range controllers[1:] {
		if ctlr == nil || l.Type(ctlr.(types.BaseVirtualDevice)) != last {
			return subresourceControllerTypeMixed
		}
	}
	return last
}

// ReadSCSIBusSharing checks the SCSI bus sharing and returns a sharing type
// depending on if all controllers are one specific kind or not. Only the first
// number of controllers specified by count are checked.
func ReadSCSIBusSharing(l object.VirtualDeviceList, count int) string {
	ctlrs := make([]types.BaseVirtualSCSIController, count)
	for _, dev := range l {
		if sc, ok := dev.(types.BaseVirtualSCSIController); ok && sc.GetVirtualSCSIController().BusNumber < int32(count) {
			ctlrs[sc.GetVirtualSCSIController().BusNumber] = sc
		}
	}
	log.Printf("[DEBUG] ReadSCSIBusSharing: SCSI controller layout for first %d controllers: %s", count, scsiControllerListString(ctlrs))
	if len(ctlrs) == 0 || ctlrs[0] == nil {
		return subresourceControllerSharingUnknown
	}
	last := ctlrs[0].(types.BaseVirtualSCSIController).GetVirtualSCSIController().SharedBus
	for _, ctlr := range ctlrs[1:] {
		if ctlr == nil || ctlr.(types.BaseVirtualSCSIController).GetVirtualSCSIController().SharedBus != last {
			return subresourceControllerSharingMixed
		}
	}
	return string(last)
}

// pickController picks a controller at the specific bus number supplied.
func pickController(l object.VirtualDeviceList, bus int, ct string) (types.BaseVirtualController, error) {
	log.Printf("[DEBUG] pickController: Looking for %s controller at bus number %d", ct, bus)
	l = l.Select(func(device types.BaseVirtualDevice) bool {
		switch d := device.(type) {
		case types.BaseVirtualSCSIController:
			if ct == "scsi" {
				return d.GetVirtualSCSIController().BusNumber == int32(bus)
			}
		case types.BaseVirtualSATAController:
			if ct == "sata" {
				return d.GetVirtualSATAController().BusNumber == int32(bus)
			}
		case *types.VirtualIDEController:
			if ct == "ide" {
				return d.GetVirtualController().BusNumber == int32(bus)
			}
		case *types.VirtualPCIController:
			if ct == "pci" {
				return d.GetVirtualController().BusNumber == int32(bus)
			}
		}
		return false
	})

	if len(l) == 0 {
		return nil, fmt.Errorf("could not find controller at bus number %d", bus)
	}

	log.Printf("[DEBUG] pickSCSIController: Found controller: %s", l.Name(l[0]))
	return l[0].(types.BaseVirtualController), nil
}

// ControllerForCreateUpdate wraps the controller selection logic to make it
// easier to use in create or update operations. If the controller type is a
// SCSI device, the bus number is searched as well.
func (r *Subresource) ControllerForCreateUpdate(l object.VirtualDeviceList, ct string, bus int) (types.BaseVirtualController, error) {
	log.Printf("[DEBUG] ControllerForCreateUpdate: Looking for controller type %s", ct)
	var ctlr types.BaseVirtualController
	var err error
	ctlr, err = pickController(l, bus, ct)
	if err != nil {
		return nil, err
	}
	if ctlr == nil {
		return nil, fmt.Errorf("could not find an available %s controller", ct)
	}

	log.Printf("[DEBUG] ControllerForCreateUpdate: Found controller: %s", l.Name(ctlr.(types.BaseVirtualDevice)))

	return ctlr, nil
}

// ApplyDeviceChange applies a pending types.BaseVirtualDeviceConfigSpec to a
// working set to either add, remove, or update devices so that the working
// VirtualDeviceList is as up to date as possible.
func applyDeviceChange(l object.VirtualDeviceList, cs []types.BaseVirtualDeviceConfigSpec) object.VirtualDeviceList {
	log.Printf("[DEBUG] applyDeviceChange: Applying changes: %s", DeviceChangeString(cs))
	log.Printf("[DEBUG] applyDeviceChange: Device list before changes: %s", DeviceListString(l))
	for _, s := range cs {
		spec := s.GetVirtualDeviceConfigSpec()
		switch spec.Operation {
		case types.VirtualDeviceConfigSpecOperationAdd:
			l = append(l, spec.Device)
		case types.VirtualDeviceConfigSpecOperationEdit:
			// Edit operations may not be 100% necessary to apply. This is because
			// more often than not, the device will probably be edited in place off
			// of the original reference, meaning that the slice should actually
			// point to the updated item. However, the safer of the two options is to
			// assume that this may *not* be happening as we are not enforcing that
			// in implementation anywhere.
			for n, dev := range l {
				if dev.GetVirtualDevice().Key == spec.Device.GetVirtualDevice().Key {
					l[n] = spec.Device
				}
			}
		case types.VirtualDeviceConfigSpecOperationRemove:
			for i := 0; i < len(l); i++ {
				dev := l[i]
				if dev.GetVirtualDevice().Key == spec.Device.GetVirtualDevice().Key {
					l = append(l[:i], l[i+1:]...)
					i--
				}
			}
		default:
			panic("unknown op")
		}
	}
	log.Printf("[DEBUG] applyDeviceChange: Device list after changes: %s", DeviceListString(l))
	return l
}

// DeviceListString pretty-prints each device in a virtual device list, used
// for logging purposes and what not.
func DeviceListString(l object.VirtualDeviceList) string {
	var names []string
	for _, d := range l {
		if d == nil {
			names = append(names, "<nil>")
		} else {
			names = append(names, l.Name(d))
		}
	}
	return strings.Join(names, ",")
}

// DeviceChangeString pretty-prints a slice of VirtualDeviceConfigSpec.
func DeviceChangeString(specs []types.BaseVirtualDeviceConfigSpec) string {
	var strs []string
	for _, v := range specs {
		spec := v.GetVirtualDeviceConfigSpec()
		strs = append(strs, fmt.Sprintf("(%s: %T at key %d)", string(spec.Operation), spec.Device, spec.Device.GetVirtualDevice().Key))
	}
	return strings.Join(strs, ",")
}

// subresourceListString takes a list of sub-resources and pretty-prints the
// key and device address.
func subresourceListString(data []interface{}) string {
	var strs []string
	for _, v := range data {
		if v == nil {
			strs = append(strs, "(<nil>)")
			continue
		}
		m := v.(map[string]interface{})
		devaddr := m["device_address"].(string)
		if devaddr == "" {
			devaddr = "<new device>"
		}
		strs = append(strs, fmt.Sprintf("(key %d at %s)", m["key"].(int), devaddr))
	}
	return strings.Join(strs, ",")
}

// ideControllerListString pretty-prints a slice of IDE controllers.
func ideControllerListString(ctlrs []*types.VirtualIDEController) string {
	var l object.VirtualDeviceList
	for _, ctlr := range ctlrs {
		if ctlr == nil {
			l = append(l, types.BaseVirtualDevice(nil))
		} else {
			l = append(l, ctlr.GetVirtualDevice())
		}
	}
	return DeviceListString(l)
}

// sataControllerListString pretty-prints a slice of SATA controllers.
func sataControllerListString(ctlrs []types.BaseVirtualSATAController) string {
	var l object.VirtualDeviceList
	for _, ctlr := range ctlrs {
		if ctlr == nil {
			l = append(l, types.BaseVirtualDevice(nil))
		} else {
			l = append(l, ctlr.(types.BaseVirtualDevice))
		}
	}
	return DeviceListString(l)
}

// scsiControllerListString pretty-prints a slice of SCSI controllers.
func scsiControllerListString(ctlrs []types.BaseVirtualSCSIController) string {
	var l object.VirtualDeviceList
	for _, ctlr := range ctlrs {
		if ctlr == nil {
			l = append(l, types.BaseVirtualDevice(nil))
		} else {
			l = append(l, ctlr.(types.BaseVirtualDevice))
		}
	}
	return DeviceListString(l)
}

// AppendDeviceChangeSpec appends unique copies of the supplied device change
// operations and appends them to spec. The resulting list is returned.
//
// The object of this function is to provide deep copies of each virtual device
// to the spec as they looked like when the append operation was called,
// helping facilitate multiple update operations to the same device in a single
// reconfigure call.
func AppendDeviceChangeSpec(
	spec []types.BaseVirtualDeviceConfigSpec,
	ops ...types.BaseVirtualDeviceConfigSpec,
) []types.BaseVirtualDeviceConfigSpec {
	for _, op := range ops {
		c := copystructure.Must(copystructure.Copy(op)).(types.BaseVirtualDeviceConfigSpec)
		spec = append(spec, c)
	}
	return spec
}

// ApiToPciId is a helper to convert PCI DeviceIDs to their actual value.
// vSphere appears to store the PCI information in hex, but converts it to
// an int16 for the API. With large numbers, this overflows the int16.
func ApiToPciId(i int16) string {
	return strconv.FormatInt(int64(uint16(i)), 16)
}

// getHostPciDevice returns a HostPciDevice from a host based on the DeviceId.
func (c *pciApplyConfig) getHostPciDevice(id string) (*types.HostPciDevice, error) {
	host, err := hostsystem.FromID(c.Client, c.ResourceData.Get("host_system_id").(string))
	if err != nil {
		return nil, err
	}
	hprops, err := hostsystem.Properties(host)
	if err != nil {
		return nil, err
	}
	for _, hostPci := range hprops.Hardware.PciDevice {
		if id == hostPci.Id {
			return &hostPci, nil
		}
	}
	return nil, fmt.Errorf("Unable to locate PCI device: %s", id)
}

// getPciSysId fetchs the PCI SystemId of a host. The SystemId is required for
// PCI passthrough devices.
func (c *pciApplyConfig) getPciSysId() error {
	host, err := hostsystem.FromID(c.Client, c.ResourceData.Get("host_system_id").(string))
	if err != nil {
		return err
	}
	hostRef := host.Reference()
	e, err := computeresource.EnvironmentBrowserFromReference(c.Client, hostRef)
	if err != nil {
		return err
	}
	sysId, err := e.SystemId(context.TODO(), &hostRef)
	if err != nil {
		return err
	}
	c.SystemId = sysId
	return nil
}

// modifyVirtualPciDevices will take a list of devices and an operation and
// will create the appropriate config spec.
func (c *pciApplyConfig) modifyVirtualPciDevices(devList *schema.Set, op types.VirtualDeviceConfigSpecOperation) error {
	log.Printf("VirtualMachine: Creating PCI passthrough device specs %v", op)
	for _, addDev := range devList.List() {
		log.Printf("[DEBUG] modifyVirtualPciDevices: Appending %v spec for %s", op, addDev.(string))
		pciDev, err := c.getHostPciDevice(addDev.(string))
		if err != nil {
			return err
		}
		dev := &types.VirtualPCIPassthrough{
			VirtualDevice: types.VirtualDevice{
				DynamicData: types.DynamicData{},
				Backing: &types.VirtualPCIPassthroughDeviceBackingInfo{
					VirtualDeviceDeviceBackingInfo: types.VirtualDeviceDeviceBackingInfo{},
					Id:                             pciDev.Id,
					SystemId:                       c.SystemId,
					VendorId:                       pciDev.VendorId,
				},
			},
		}
		vm, err := virtualmachine.FromUUID(c.Client, c.ResourceData.Id())
		if err != nil {
			return err
		}
		vprops, err := virtualmachine.Properties(vm)
		if err != nil {
			return err
		}
		// This will only find a device for delete operations.
		for _, vmDevP := range vprops.Config.Hardware.Device {
			if vmDev, ok := vmDevP.(*types.VirtualPCIPassthrough); ok {
				if vmDev.Backing.(*types.VirtualPCIPassthroughDeviceBackingInfo).Id == pciDev.Id {
					dev = vmDev
				}
			}
		}
		dspec, err := object.VirtualDeviceList{dev}.ConfigSpec(op)
		if err != nil {
			return err
		}
		c.Spec = append(c.Spec, dspec...)
		c.VirtualDevice = applyDeviceChange(c.VirtualDevice, dspec)
	}
	log.Printf("VirtualMachine: PCI passthrough device specs created")
	return nil
}

// PciPassthroughApplyOperation checks for changes in a virtual machine's
// PCI passthrough devices and creates config specs to apply apply to the
// virtual machine.
func PciPassthroughApplyOperation(d *schema.ResourceData, c *govmomi.Client, l object.VirtualDeviceList) (object.VirtualDeviceList, []types.BaseVirtualDeviceConfigSpec, error) {
	old, new := d.GetChange("pci_device_id")
	oldDevIds := old.(*schema.Set)
	newDevIds := new.(*schema.Set)

	delDevs := oldDevIds.Difference(newDevIds)
	addDevs := newDevIds.Difference(oldDevIds)
	applyConfig := &pciApplyConfig{
		Client:        c,
		ResourceData:  d,
		Spec:          []types.BaseVirtualDeviceConfigSpec{},
		VirtualDevice: l,
	}
	if addDevs.Len() <= 0 && delDevs.Len() <= 0 {
		return applyConfig.VirtualDevice, applyConfig.Spec, nil
	}

	d.Set("reboot_required", true)
	err := applyConfig.getPciSysId()
	if err != nil {
		return nil, nil, err
	}

	// Add new PCI passthrough devices
	err = applyConfig.modifyVirtualPciDevices(addDevs, types.VirtualDeviceConfigSpecOperationAdd)
	if err != nil {
		return nil, nil, err
	}

	// Remove deleted PCI passthrough devices
	err = applyConfig.modifyVirtualPciDevices(delDevs, types.VirtualDeviceConfigSpecOperationRemove)
	if err != nil {
		return nil, nil, err
	}
	return applyConfig.VirtualDevice, applyConfig.Spec, nil
}
