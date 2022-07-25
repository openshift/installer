package vsphere

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/virtualmachine"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/virtualdevice"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
)

// resourceVSphereVirtualMachineMigrateState is the main state migration function for
// the vsphere_virtual_machine resource.
func resourceVSphereVirtualMachineMigrateState(version int, os *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	// Guard against a nil state.
	if os == nil {
		return nil, nil
	}

	// Guard against empty state, can't do anything with it
	if os.Empty() {
		return os, nil
	}

	var migrateFunc func(*terraform.InstanceState, interface{}) error
	switch version {
	case 2:
		log.Printf("[DEBUG] Migrating vsphere_virtual_machine state: old v%d state: %#v", version, os)
		migrateFunc = migrateVSphereVirtualMachineStateV3
	case 1:
		log.Printf("[DEBUG] Migrating vsphere_virtual_machine state: old v%d state: %#v", version, os)
		migrateFunc = migrateVSphereVirtualMachineStateV2
	case 0:
		log.Printf("[DEBUG] Migrating vsphere_virtual_machine state: old v%d state: %#v", version, os)
		migrateFunc = migrateVSphereVirtualMachineStateV1
	default:
		// Migration is complete
		log.Printf("[DEBUG] Migrating vsphere_virtual_machine state: completed v%d state: %#v", version, os)
		return os, nil
	}
	if err := migrateFunc(os, meta); err != nil {
		return nil, err
	}
	version++
	log.Printf("[DEBUG] Migrating vsphere_virtual_machine state: new v%d state: %#v", version, os)
	return resourceVSphereVirtualMachineMigrateState(version, os, meta)
}

// migrateVSphereVirtualMachineStateV3 migrates the state of the
// vsphere_virtual_machine from version 2 to version 3.
func migrateVSphereVirtualMachineStateV3(is *terraform.InstanceState, meta interface{}) error {
	// All we really preserve from the old state is the UUID of the virtual
	// machine. We leverage some of the special parts of the import functionality
	// - namely validating disks, and flagging the VM as imported in the state to
	// guard against someone adding customization to the configuration and
	// accidentally forcing a new resource.
	//
	// Read will handle most of the population post-migration as it does for
	// import, and there will be an unavoidable diff for TF-only options on the
	// next plan.
	client := meta.(*Client).vimClient
	id := is.ID

	log.Printf("[DEBUG] Migrate state for VM at UUID %q", id)
	vm, err := virtualmachine.FromUUID(client, id)
	if err != nil {
		return fmt.Errorf("error fetching virtual machine: %s", err)
	}
	props, err := virtualmachine.Properties(vm)
	if err != nil {
		return fmt.Errorf("error fetching virtual machine properties: %s", err)
	}

	// Populate the UUID field of all virtual disks in state.
	diskCnt, _ := strconv.Atoi(is.Attributes["disk.#"])
	l := object.VirtualDeviceList(props.Config.Hardware.Device)
	for i := 0; i < diskCnt; i++ {
		v, ok := is.Attributes[fmt.Sprintf("disk.%d.key", i)]
		if !ok {
			return fmt.Errorf("corrupt state: key disk.%d.key not found", i)
		}
		key, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return fmt.Errorf("while converting disk.%d.key key to int32: %s", i, err)
		}
		if key < 1 {
			// This is a possibility during v1 -> v3 migrations, and would fail to
			// find a device anyway, so we just ignore these.
			continue
		}
		device := l.FindByKey(int32(key))
		if device == nil {
			// Missing device, pass
			continue
		}
		disk, ok := device.(*types.VirtualDisk)
		if !ok {
			// Not the device we are looking for
			continue
		}
		backing, ok := disk.Backing.(*types.VirtualDiskFlatVer2BackingInfo)
		if !ok {
			// Someone has tampered with the VM to the point where we really should
			// not mess with it. We can't account for all cases, but if someone has
			// added something like an RDM disk or something else that is not
			// VMDK-backed, we don't want to continue. We have never supported this.
			return fmt.Errorf("disk device %s is not a VMDK-backed virtual disk and state import cannot continue", l.Name(disk))
		}
		is.Attributes[fmt.Sprintf("disk.%d.uuid", i)] = backing.Uuid
	}

	d := resourceVSphereVirtualMachine().Data(&terraform.InstanceState{})
	log.Printf("[DEBUG] %s: Migration to V3 complete", resourceVSphereVirtualMachineIDString(d))
	return nil
}

// migrateVSphereVirtualMachineStateV2 migrates the state of the
// vsphere_virtual_machine from version 1 to version 2.
func migrateVSphereVirtualMachineStateV2(is *terraform.InstanceState, meta interface{}) error {
	// All we really preserve from the old state is the UUID of the virtual
	// machine. We leverage some of the special parts of the import functionality
	// - namely validating disks, and flagging the VM as imported in the state to
	// guard against someone adding customization to the configuration and
	// accidentally forcing a new resource. To assist with the migration of state
	// from V1 to V3 as well, we now pull in disk attribute data that is
	// populated during the import process.
	//
	// Read will handle most of the population post-migration as it does for
	// import, and there will be an unavoidable diff for TF-only options on the
	// next plan. This diff should not require a reconfigure of the virtual
	// machine.
	client := meta.(*Client).vimClient
	name := is.ID
	id := is.Attributes["uuid"]
	if id == "" {
		return fmt.Errorf("resource ID %s has no UUID. State cannot be migrated", name)
	}

	log.Printf("[DEBUG] Migrate state for VM resource %q: UUID %q", name, id)
	vm, err := virtualmachine.FromUUID(client, id)
	if err != nil {
		return fmt.Errorf("error fetching virtual machine: %s", err)
	}
	props, err := virtualmachine.Properties(vm)
	if err != nil {
		return fmt.Errorf("error fetching virtual machine properties: %s", err)
	}

	// Validate the disks in the VM to make sure that they will work with the new
	// version of the resource. This is mainly ensuring that all disks are SCSI
	// disks.
	//
	// NOTE: This uses the current version of the resource to make this check,
	// which at some point in time may end up being a higher schema version than
	// version 2. At this point in time, there is nothing here that would cause
	// issues (nothing in the sub-resource read logic is reliant on schema
	// versions), and an empty ResourceData is sent anyway.
	diskCnt, _ := strconv.Atoi(is.Attributes["disk.#"])
	maxBus := diskCnt / 15
	l := object.VirtualDeviceList(props.Config.Hardware.Device)
	for k, v := range is.Attributes {
		if !regexp.MustCompile("disk\\.[0-9]+\\.key").MatchString(k) {
			continue
		}
		key, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return fmt.Errorf("while converting key %s to int32: %s", k, err)
		}
		if key < 1 {
			continue
		}
		device := l.FindByKey(int32(key))
		if device == nil {
			continue
		}
		ctlr := l.FindByKey(device.GetVirtualDevice().ControllerKey)
		if ctlr == nil {
			continue
		}
		if sc, ok := ctlr.(types.BaseVirtualSCSIController); ok && sc.GetVirtualSCSIController().BusNumber > int32(maxBus) {
			maxBus = int(sc.GetVirtualSCSIController().BusNumber)
		}
	}

	d := resourceVSphereVirtualMachine().Data(&terraform.InstanceState{})
	_ = d.Set("scsi_controller_count", maxBus+1)
	if err := virtualdevice.DiskImportOperation(d, object.VirtualDeviceList(props.Config.Hardware.Device)); err != nil {
		return err
	}

	rs := resourceVSphereVirtualMachine().Schema
	var guestNetTimeout string
	switch is.Attributes["wait_for_guest_net"] {
	case "false":
		guestNetTimeout = "-1"
	default:
		guestNetTimeout = fmt.Sprintf("%v", rs["wait_for_guest_net_timeout"].Default)
	}

	// The VM should be ready for reading now
	is.Attributes = make(map[string]string)
	is.ID = id
	is.Attributes["imported"] = "true"

	// Set some defaults. This helps possibly prevent diffs where these values
	// have not been changed.
	is.Attributes["scsi_controller_count"] = fmt.Sprintf("%v", rs["scsi_controller_count"].Default)
	is.Attributes["force_power_off"] = fmt.Sprintf("%v", rs["force_power_off"].Default)
	is.Attributes["migrate_wait_timeout"] = fmt.Sprintf("%v", rs["migrate_wait_timeout"].Default)
	is.Attributes["shutdown_wait_timeout"] = fmt.Sprintf("%v", rs["shutdown_wait_timeout"].Default)
	is.Attributes["wait_for_guest_net_timeout"] = guestNetTimeout
	is.Attributes["wait_for_guest_net_routable"] = fmt.Sprintf("%v", rs["wait_for_guest_net_routable"].Default)
	is.Attributes["scsi_controller_count"] = fmt.Sprintf("%v", maxBus+1)

	// Populate our disk data from the fake state.
	d.SetId(id)
	for k, v := range d.State().Attributes {
		if strings.HasPrefix(k, "disk.") {
			is.Attributes[k] = v
		}
	}

	log.Printf("[DEBUG] %s: Migration to V2 complete", resourceVSphereVirtualMachineIDString(d))
	return nil
}

func migrateVSphereVirtualMachineStateV1(is *terraform.InstanceState, _ interface{}) error {
	if is.Empty() || is.Attributes == nil {
		log.Println("[DEBUG] Empty VSphere Virtual Machine State; nothing to migrate.")
		return nil
	}

	log.Printf("[DEBUG] Attributes before migration: %#v", is.Attributes)

	if is.Attributes["skip_customization"] == "" {
		is.Attributes["skip_customization"] = "false"
	}

	if is.Attributes["enable_disk_uuid"] == "" {
		is.Attributes["enable_disk_uuid"] = "false"
	}

	for k := range is.Attributes {
		if strings.HasPrefix(k, "disk.") && strings.HasSuffix(k, ".size") {
			diskParts := strings.Split(k, ".")
			if len(diskParts) != 3 {
				continue
			}
			s := strings.Join([]string{diskParts[0], diskParts[1], "controller_type"}, ".")
			if _, ok := is.Attributes[s]; !ok {
				is.Attributes[s] = "scsi"
			}
		}
	}

	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)
	return nil
}
