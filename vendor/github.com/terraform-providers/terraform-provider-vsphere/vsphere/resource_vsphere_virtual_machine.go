package vsphere

import (
	"context"
	"errors"
	"fmt"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/contentlibrary"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/ovfdeploy"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/customattribute"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/datastore"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/folder"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/hostsystem"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/resourcepool"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/spbm"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/storagepod"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/vappcontainer"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/viapi"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/virtualmachine"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/virtualdevice"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/vmworkflow"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/ovf"
	"github.com/vmware/govmomi/vim25/types"
)

// formatVirtualMachinePostCloneRollbackError defines the verbose error when
// rollback fails on a post-clone virtual machine operation.
const formatVirtualMachinePostCloneRollbackError = `
WARNING:
There was an error performing post-clone changes to virtual machine %q:
%s
Additionally, there was an error removing the cloned virtual machine:
%s

The virtual machine may still exist in Terraform state. If it does, the
resource will need to be tainted before trying again. For more information on
how to do this, see the following page:
https://www.terraform.io/docs/commands/taint.html

If the virtual machine does not exist in state, manually delete it to try again.
`

// formatVirtualMachineCustomizationWaitError defines the verbose error that is
// sent when the customization waiter returns an error. This can either be due
// to timeout waiting for respective events or a guest-specific customization
// error. The resource does not roll back in this case, to assist with
// troubleshooting.
const formatVirtualMachineCustomizationWaitError = `
Virtual machine customization failed on %q:

%s

The virtual machine has not been deleted to assist with troubleshooting. If
corrective steps are taken without modifying the "customize" block of the
resource configuration, the resource will need to be tainted before trying
again. For more information on how to do this, see the following page:
https://www.terraform.io/docs/commands/taint.html
`

const questionCheckIntervalSecs = 5

func resourceVSphereVirtualMachine() *schema.Resource {
	s := map[string]*schema.Schema{
		"resource_pool_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The ID of a resource pool to put the virtual machine in.",
		},
		"datastore_id": {
			Type:          schema.TypeString,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"datastore_cluster_id"},
			Description:   "The ID of the virtual machine's datastore. The virtual machine configuration is placed here, along with any virtual disks that are created without datastores.",
		},
		"datastore_cluster_id": {
			Type:          schema.TypeString,
			Optional:      true,
			ConflictsWith: []string{"datastore_id"},
			Description:   "The ID of a datastore cluster to put the virtual machine in.",
		},
		"datacenter_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The ID of the datacenter where the VM is to be created.",
		},
		"folder": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The name of the folder to locate the virtual machine in.",
			StateFunc:   folder.NormalizePath,
		},
		"host_system_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The ID of an optional host system to pin the virtual machine to.",
		},
		"wait_for_guest_ip_timeout": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
			Description: "The amount of time, in minutes, to wait for an available IP address on this virtual machine. A value less than 1 disables the waiter.",
		},
		"wait_for_guest_net_timeout": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     5,
			Description: "The amount of time, in minutes, to wait for an available IP address on this virtual machine. A value less than 1 disables the waiter.",
		},
		"wait_for_guest_net_routable": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Controls whether or not the guest network waiter waits for a routable address. When false, the waiter does not wait for a default gateway, nor are IP addresses checked against any discovered default gateways as part of its success criteria.",
		},
		"ignored_guest_ips": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of IP addresses and CIDR networks to ignore while waiting for an IP",
			Elem: &schema.Schema{
				Type: schema.TypeString,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if strings.Contains(v, "/") {
						_, _, err := net.ParseCIDR(v)
						if err != nil {
							errs = append(errs, fmt.Errorf("%q contains invalid CIDR address: %s", key, v))
						}
					} else if net.ParseIP(v) == nil {
						errs = append(errs, fmt.Errorf("%q contains invalid IP address: %s", key, v))
					}
					return
				},
			},
		},
		"shutdown_wait_timeout": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      3,
			Description:  "The amount of time, in minutes, to wait for shutdown when making necessary updates to the virtual machine.",
			ValidateFunc: validation.IntBetween(1, 10),
		},
		"migrate_wait_timeout": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      30,
			Description:  "The amount of time, in minutes, to wait for a vMotion operation to complete before failing.",
			ValidateFunc: validation.IntAtLeast(10),
		},
		"poweron_timeout": {
			Type:         schema.TypeInt,
			Description:  "The amount of time, in seconds, that we will be trying to power on a VM",
			Default:      300,
			ValidateFunc: validation.IntAtLeast(300),
			Optional:     true,
		},
		"force_power_off": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Set to true to force power-off a virtual machine if a graceful guest shutdown failed for a necessary operation.",
		},
		"scsi_controller_count": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      1,
			Description:  "The number of SCSI controllers that Terraform manages on this virtual machine. This directly affects the amount of disks you can add to the virtual machine and the maximum disk unit number. Note that lowering this value does not remove controllers.",
			ValidateFunc: validation.IntBetween(1, 4),
		},
		"scsi_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      virtualdevice.SubresourceControllerTypeParaVirtual,
			Description:  "The type of SCSI bus this virtual machine will have. Can be one of lsilogic, lsilogic-sas or pvscsi.",
			ValidateFunc: validation.StringInSlice(virtualdevice.SCSIBusTypeAllowedValues, false),
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				if len(d.Get("ovf_deploy").([]interface{})) > 0 {
					return true
				}
				return false
			},
		},
		"scsi_bus_sharing": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      string(types.VirtualSCSISharingNoSharing),
			Description:  "Mode for sharing the SCSI bus. The modes are physicalSharing, virtualSharing, and noSharing.",
			ValidateFunc: validation.StringInSlice(virtualdevice.SCSIBusSharingAllowedValues, false),
		},
		// NOTE: disk is only optional so that we can flag it as computed and use
		// it in ResourceDiff. We validate this field in ResourceDiff to enforce it
		// having a minimum count of 1 for now - but may support diskless VMs
		// later.
		"disk": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "A specification for a virtual disk device on this virtual machine.",
			MaxItems:    60,
			Elem:        &schema.Resource{Schema: virtualdevice.DiskSubresourceSchema()},
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				if len(d.Get("ovf_deploy").([]interface{})) > 0 {
					return true
				}
				return false
			},
		},
		"network_interface": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "A specification for a virtual NIC on this virtual machine.",
			MaxItems:    10,
			Elem:        &schema.Resource{Schema: virtualdevice.NetworkInterfaceSubresourceSchema()},
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				//no network_interface diff for ovf template deployment
				if len(d.Get("ovf_deploy").([]interface{})) > 0 {
					return true
				}
				return false
			},
		},
		"cdrom": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "A specification for a CDROM device on this virtual machine.",
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: virtualdevice.CdromSubresourceSchema()},
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				if len(d.Get("ovf_deploy").([]interface{})) > 0 {
					return true
				}
				return false
			},
		},
		"clone": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "A specification for cloning a virtual machine from template.",
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: vmworkflow.VirtualMachineCloneSchema()},
		},
		"ovf_deploy": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "A specification for deploying a virtual machine from OVF template.",
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: vmworkflow.VirtualMachineOVFDeploySchema()},
		},
		"reboot_required": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Value internal to Terraform used to determine if a configuration set change requires a reboot.",
		},
		"vmware_tools_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The state of VMware tools in the guest. This will determine the proper course of action for some device operations.",
		},
		"vmx_path": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The path of the virtual machine's configuration file in the VM's datastore.",
		},
		"imported": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "A flag internal to Terraform that indicates that this resource was either imported or came from a earlier major version of this resource. Reset after the first post-import or post-upgrade apply.",
		},
		"moid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The machine object ID from VMWare",
		},
		vSphereTagAttributeKey:    tagsSchema(),
		customattribute.ConfigKey: customattribute.ConfigSchema(),
	}
	structure.MergeSchema(s, schemaVirtualMachineConfigSpec())
	structure.MergeSchema(s, schemaVirtualMachineGuestInfo())

	return &schema.Resource{
		Create:        resourceVSphereVirtualMachineCreate,
		Read:          resourceVSphereVirtualMachineRead,
		Update:        resourceVSphereVirtualMachineUpdate,
		Delete:        resourceVSphereVirtualMachineDelete,
		CustomizeDiff: resourceVSphereVirtualMachineCustomizeDiff,
		Importer: &schema.ResourceImporter{
			State: resourceVSphereVirtualMachineImport,
		},
		SchemaVersion: 3,
		MigrateState:  resourceVSphereVirtualMachineMigrateState,
		Schema:        s,
	}
}

func resourceVSphereVirtualMachineCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning create", resourceVSphereVirtualMachineIDString(d))
	client := meta.(*VSphereClient).vimClient
	tagsClient, err := tagsManagerIfDefined(d, meta)
	if err != nil {
		return err
	}
	// Verify a proper vCenter before proceeding if custom attributes are defined
	attrsProcessor, err := customattribute.GetDiffProcessorIfAttributesDefined(client, d)
	if err != nil {
		return err
	}

	var vm *object.VirtualMachine
	// This is where we process our various VM deploy workflows. We expect the ID
	// of the resource to be set in the workflow to ensure that any post-create
	// operations that fail during this process don't create a dangling resource.
	// The VM should also be returned powered on.
	switch {
	case len(d.Get("clone").([]interface{})) > 0:
		vm, err = resourceVSphereVirtualMachineCreateClone(d, meta)
	case len(d.Get("ovf_deploy").([]interface{})) > 0:
		vm, err = resourceVsphereMachineDeployOVF(d, meta)
	default:
		vm, err = resourceVSphereVirtualMachineCreateBare(d, meta)
	}

	if err != nil {
		return err
	}

	// Tag the VM
	if tagsClient != nil {
		if err := processTagDiff(tagsClient, d, vm); err != nil {
			return err
		}
	}

	// Set custom attributes
	if attrsProcessor != nil {
		if err := attrsProcessor.ProcessDiff(vm); err != nil {
			return err
		}
	}

	// The host attribute of CreateVM_Task seems to be ignored in vCenter 6.7.
	// Ensure that VMs are on the correct host and relocate if necessary. Do this
	// near the end of the VM creation since it involves updating the
	// ResourceData.
	vprops, err := virtualmachine.Properties(vm)
	if err != nil {
		return err
	}
	if hid, ok := d.GetOk("host_system_id"); hid.(string) != vprops.Runtime.Host.Reference().Value && ok {
		err = resourceVSphereVirtualMachineRead(d, meta)
		if err != nil {
			return err
		}
		// Restore the old host_system_id so we can still tell if a relocation is
		// necessary.
		err = d.Set("host_system_id", hid.(string))
		if err != nil {
			return err
		}
		if err = resourceVSphereVirtualMachineUpdateLocation(d, meta); err != nil {
			return err
		}
	}

	// Wait for guest IP address if we have been set to wait for one
	err = virtualmachine.WaitForGuestIP(
		client,
		vm,
		d.Get("wait_for_guest_ip_timeout").(int),
		d.Get("ignored_guest_ips").([]interface{}),
	)
	if err != nil {
		return err
	}

	// Wait for a routable address if we have been set to wait for one
	err = virtualmachine.WaitForGuestNet(
		client,
		vm,
		d.Get("wait_for_guest_net_routable").(bool),
		d.Get("wait_for_guest_net_timeout").(int),
		d.Get("ignored_guest_ips").([]interface{}),
	)
	if err != nil {
		return err
	}

	// All done!
	log.Printf("[DEBUG] %s: Create complete", resourceVSphereVirtualMachineIDString(d))
	return resourceVSphereVirtualMachineRead(d, meta)
}

func resourceVSphereVirtualMachineRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Reading state of virtual machine", resourceVSphereVirtualMachineIDString(d))
	client := meta.(*VSphereClient).vimClient
	id := d.Id()
	vm, err := virtualmachine.FromUUID(client, id)
	if err != nil {
		if _, ok := err.(*virtualmachine.UUIDNotFoundError); ok {
			log.Printf("[DEBUG] %s: Virtual machine not found, marking resource as gone: %s", resourceVSphereVirtualMachineIDString(d), err)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error searching for with UUID %q: %s", id, err)
	}

	vprops, err := virtualmachine.Properties(vm)
	if err != nil {
		return fmt.Errorf("error fetching VM properties: %s", err)
	}

	// Set the managed object id.
	moid := vm.Reference().Value
	d.Set("moid", moid)
	log.Printf("[DEBUG] MOID for VM %q is %q", vm.InventoryPath, moid)

	// Reset reboot_required. This is an update only variable and should not be
	// set across TF runs.
	d.Set("reboot_required", false)
	// Check to see if VMware tools is running.
	if vprops.Guest != nil {
		d.Set("vmware_tools_status", vprops.Guest.ToolsRunningStatus)
	}

	// Resource pool
	if vprops.ResourcePool != nil {
		d.Set("resource_pool_id", vprops.ResourcePool.Value)
	}
	// If the VM is part of a vApp, InventoryPath will point to a host path
	// rather than a VM path, so this step must be skipped.
	var vmContainer string
	if vprops.ParentVApp != nil {
		vmContainer = vprops.ParentVApp.Value
	} else {
		vmContainer = vprops.ResourcePool.Value
	}
	if !vappcontainer.IsVApp(client, vmContainer) {
		f, err := folder.RootPathParticleVM.SplitRelativeFolder(vm.InventoryPath)
		if err != nil {
			return fmt.Errorf("error parsing virtual machine path %q: %s", vm.InventoryPath, err)
		}
		d.Set("folder", folder.NormalizePath(f))
	}
	// Set VM's current host ID if available
	if vprops.Runtime.Host != nil {
		d.Set("host_system_id", vprops.Runtime.Host.Value)
	}

	// Set the VMX path and default datastore
	dp := &object.DatastorePath{}
	if ok := dp.FromString(vprops.Config.Files.VmPathName); !ok {
		return fmt.Errorf("could not parse VMX file path: %s", vprops.Config.Files.VmPathName)
	}
	// The easiest path for us to get an exact match on the datastore in use is
	// to look for the datastore name in the list of used datastores. This is
	// something we have access to from the VM's properties. This allows us to
	// get away with not having to have the datastore unnecessarily supplied to
	// the resource when it's not used by anything else.
	var ds *object.Datastore
	for _, dsRef := range vprops.Datastore {
		dsx, err := datastore.FromID(client, dsRef.Value)
		if err != nil {
			return fmt.Errorf("error locating VMX datastore: %s", err)
		}
		dsxProps, err := datastore.Properties(dsx)
		if err != nil {
			return fmt.Errorf("error fetching VMX datastore properties: %s", err)
		}
		if dsxProps.Summary.Name == dp.Datastore {
			ds = dsx
		}
	}
	if ds == nil {
		return fmt.Errorf("VMX datastore %s not found", dp.Datastore)
	}
	d.Set("datastore_id", ds.Reference().Value)
	d.Set("vmx_path", dp.Path)

	// Read general VM config info
	if err := flattenVirtualMachineConfigInfo(d, vprops.Config); err != nil {
		return fmt.Errorf("error reading virtual machine configuration: %s", err)
	}

	// Read the VM Home storage policy if associated
	polID, err := spbm.PolicyIDByVirtualMachine(client, moid)
	if err != nil {
		return err
	}
	d.Set("storage_policy_id", polID)

	// Perform pending device read operations.
	devices := object.VirtualDeviceList(vprops.Config.Hardware.Device)
	// Read the state of the SCSI bus.
	d.Set("scsi_type", virtualdevice.ReadSCSIBusType(devices, d.Get("scsi_controller_count").(int)))
	d.Set("scsi_bus_sharing", virtualdevice.ReadSCSIBusSharing(devices, d.Get("scsi_controller_count").(int)))
	// Disks first
	if err := virtualdevice.DiskRefreshOperation(d, client, devices); err != nil {
		return err
	}
	// Network devices
	if err := virtualdevice.NetworkInterfaceRefreshOperation(d, client, devices); err != nil {
		return err
	}
	// CDROM
	if err := virtualdevice.CdromRefreshOperation(d, client, devices); err != nil {
		return err
	}

	// Read tags if we have the ability to do so
	if tagsClient, _ := meta.(*VSphereClient).TagsManager(); tagsClient != nil {
		if err := readTagsForResource(tagsClient, vm, d); err != nil {
			return err
		}
	}

	// Read set custom attributes
	if customattribute.IsSupported(client) {
		customattribute.ReadFromResource(client, vprops.Entity(), d)
	}

	// Finally, select a valid IP address for use by the VM for purposes of
	// provisioning. This also populates some computed values to present to the
	// user.
	if vprops.Guest != nil {
		if err := buildAndSelectGuestIPs(d, *vprops.Guest); err != nil {
			return fmt.Errorf("error reading virtual machine guest data: %s", err)
		}
	}

	log.Printf("[DEBUG] %s: Read complete", resourceVSphereVirtualMachineIDString(d))
	return nil
}

func resourceVSphereVirtualMachineUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Performing update", resourceVSphereVirtualMachineIDString(d))
	client := meta.(*VSphereClient).vimClient
	tagsClient, err := tagsManagerIfDefined(d, meta)
	if err != nil {
		return err
	}
	// Verify a proper vCenter before proceeding if custom attributes are defined
	attrsProcessor, err := customattribute.GetDiffProcessorIfAttributesDefined(client, d)
	if err != nil {
		return err
	}

	id := d.Id()
	vm, err := virtualmachine.FromUUID(client, id)
	if err != nil {
		return fmt.Errorf("cannot locate virtual machine with UUID %q: %s", id, err)
	}

	if d.HasChange("resource_pool_id") {
		var rp *object.ResourcePool
		rp, err = resourcepool.FromID(client, d.Get("resource_pool_id").(string))
		if err != nil {
			return err
		}

		// Before we move the VM to the new RP we need to make sure the new one is on the same
		// cluster or host as the old one, otherwise vsphere will throw an error.
		dstRPProps, err := resourcepool.Properties(rp)
		if err != nil {
			return err
		}

		vmProps, err := virtualmachine.Properties(vm)
		if err != nil {
			return err
		}
		srcRPID := vmProps.ResourcePool.Value
		srcResourcePool, err := resourcepool.FromID(client, srcRPID)
		if err != nil {
			return err
		}
		srcRPProps, err := resourcepool.Properties(srcResourcePool)
		if err != nil {
			return err
		}

		// If the source and destination RPs have different owners (i.e hosts or clusters)
		// then it will be handled as a vMotion task
		if dstRPProps.Owner == srcRPProps.Owner {
			err = resourcepool.MoveIntoResourcePool(rp, vm.Reference())
			if err != nil {
				return err
			}
		} else {
			// If we're migrating away from the current host we're setting the host system ID
			// to nothing. It will be populated after the migration step, once we call Read().
			d.Set("host_system_id", "")
		}
		// If a VM is moved into or out of a vApp container, the VM's InventoryPath
		// will change. This can affect steps later in the update process such as
		// moving folders. To make sure the VM has the correct InventoryPath,
		// refresh the VM after moving into a new resource pool.
		vm, err = virtualmachine.FromMOID(client, vm.Reference().Value)
		if err != nil {
			return err
		}
	}

	// Update folder if necessary
	if d.HasChange("folder") && !vappcontainer.IsVApp(client, d.Get("resource_pool_id").(string)) {
		folder := d.Get("folder").(string)
		if err := virtualmachine.MoveToFolder(client, vm, folder); err != nil {
			return fmt.Errorf("could not move virtual machine to folder %q: %s", folder, err)
		}
	}

	// Apply any pending tags
	if tagsClient != nil {
		if err := processTagDiff(tagsClient, d, vm); err != nil {
			return err
		}
	}

	// Update custom attributes
	if attrsProcessor != nil {
		if err := attrsProcessor.ProcessDiff(vm); err != nil {
			return err
		}
	}

	// Ready to start the VM update. All changes from here, until the update
	// operation finishes successfully, need to be done in partial mode.
	d.Partial(true)

	vprops, err := virtualmachine.Properties(vm)
	if err != nil {
		return fmt.Errorf("error fetching VM properties: %s", err)
	}

	spec, changed, err := expandVirtualMachineConfigSpecChanged(d, client, vprops.Config)
	if err != nil {
		return fmt.Errorf("error in virtual machine configuration: %s", err)
	}

	devices := object.VirtualDeviceList(vprops.Config.Hardware.Device)
	if spec.DeviceChange, err = applyVirtualDevices(d, client, devices); err != nil {
		return err
	}
	// Only carry out the reconfigure if we actually have a change to process.
	cv := virtualmachine.GetHardwareVersionNumber(vprops.Config.Version)
	tv := d.Get("hardware_version").(int)
	if tv > cv {
		d.Set("reboot_required", true)
	}
	if changed || len(spec.DeviceChange) > 0 {
		//Check to see if we need to shutdown the VM for this process.
		if d.Get("reboot_required").(bool) && vprops.Runtime.PowerState != types.VirtualMachinePowerStatePoweredOff {
			// Attempt a graceful shutdown of this process. We wrap this in a VM helper.
			timeout := d.Get("shutdown_wait_timeout").(int)
			force := d.Get("force_power_off").(bool)
			if err := virtualmachine.GracefulPowerOff(client, vm, timeout, force); err != nil {
				return fmt.Errorf("error shutting down virtual machine: %s", err)
			}
		}

		// Start goroutine here that checks for questions
		gChan := make(chan bool)

		questions := map[string]string{
			"msg.cdromdisconnect.locked": "0",
		}
		go func() {
			// Sleep for a bit
			time.Sleep(questionCheckIntervalSecs * time.Second)
			for {
				select {
				case <-gChan:
					// We're done
					break
				default:
					vprops, err := virtualmachine.Properties(vm)
					if err != nil {
						log.Printf("[DEBUG] Error while retrieving VM properties. Error: %s", err)
						continue
					}
					q := vprops.Runtime.Question
					if q != nil {
						log.Printf("[DEBUG] Question: %#v", q)
						if len(q.Message) < 1 {
							log.Printf("[DEBUG] No messages found")
							continue
						}
						qMsg := q.Message[0].Id
						if response, ok := questions[qMsg]; ok {
							if err = vm.Answer(context.TODO(), q.Id, response); err != nil {
								log.Printf("[DEBUG] Failed to answer question. Error: %s", err)
								break
							}
						}
					} else {
						log.Printf("[DEBUG] No questions found")
					}
				}
			}
		}()

		// Perform updates.
		if _, ok := d.GetOk("datastore_cluster_id"); ok {
			err = resourceVSphereVirtualMachineUpdateReconfigureWithSDRS(d, meta, vm, spec)
		} else {
			err = virtualmachine.Reconfigure(vm, spec)
		}

		// Upgrade the VM's hardware version if needed.
		err = virtualmachine.SetHardwareVersion(vm, d.Get("hardware_version").(int))
		if err != nil {
			return err
		}

		// Regardless of the result we no longer need to watch for pending questions.
		gChan <- true

		if err != nil {
			return fmt.Errorf("error reconfiguring virtual machine: %s", err)
		}
		// Re-fetch properties
		vprops, err = virtualmachine.Properties(vm)
		if err != nil {
			return fmt.Errorf("error re-fetching VM properties after update: %s", err)
		}
		// Power back on the VM, and wait for network if necessary.
		if vprops.Runtime.PowerState != types.VirtualMachinePowerStatePoweredOn {
			pTimeoutStr := fmt.Sprintf("%ds", d.Get("poweron_timeout").(int))
			pTimeout, err := time.ParseDuration(pTimeoutStr)
			if err != nil {
				return fmt.Errorf("failed to parse poweron_timeout as a valid duration: %s", err)
			}
			// Start the virtual machine
			if err := virtualmachine.PowerOn(vm, pTimeout); err != nil {
				return fmt.Errorf("error powering on virtual machine: %s", err)
			}
			err = virtualmachine.WaitForGuestIP(
				client,
				vm,
				d.Get("wait_for_guest_ip_timeout").(int),
				d.Get("ignored_guest_ips").([]interface{}),
			)
			if err != nil {
				return err
			}
			err = virtualmachine.WaitForGuestNet(
				client,
				vm,
				d.Get("wait_for_guest_net_routable").(bool),
				d.Get("wait_for_guest_net_timeout").(int),
				d.Get("ignored_guest_ips").([]interface{}),
			)
			if err != nil {
				return err
			}
		}
	}

	// Now safe to turn off partial mode.
	d.Partial(false)
	d.Set("reboot_required", false)

	// Now that any pending changes have been done (namely, any disks that don't
	// need to be migrated have been deleted), proceed with vMotion if we have
	// one pending.
	if err := resourceVSphereVirtualMachineUpdateLocation(d, meta); err != nil {
		return fmt.Errorf("error running VM migration: %s", err)
	}

	// All done with updates.
	log.Printf("[DEBUG] %s: Update complete", resourceVSphereVirtualMachineIDString(d))
	return resourceVSphereVirtualMachineRead(d, meta)
}

// resourceVSphereVirtualMachineUpdateReconfigureWithSDRS runs the reconfigure
// part of resourceVSphereVirtualMachineUpdate through storage DRS. It's
// designed to be run when a storage cluster is specified, versus simply
// specifying datastores.
func resourceVSphereVirtualMachineUpdateReconfigureWithSDRS(
	d *schema.ResourceData,
	meta interface{},
	vm *object.VirtualMachine,
	spec types.VirtualMachineConfigSpec,
) error {
	// Check to see if we have any disk creation operations first, as sending an
	// update through SDRS without any disk creation operations will fail.
	if !storagepod.HasDiskCreationOperations(spec.DeviceChange) {
		log.Printf("[DEBUG] No disk operations for reconfiguration of VM %q, deferring to standard API", vm.InventoryPath)
		return virtualmachine.Reconfigure(vm, spec)
	}

	client := meta.(*VSphereClient).vimClient
	if err := viapi.ValidateVirtualCenter(client); err != nil {
		return fmt.Errorf("connection ineligible to use datastore_cluster_id: %s", err)
	}

	log.Printf("[DEBUG] %s: Reconfiguring virtual machine through Storage DRS API", resourceVSphereVirtualMachineIDString(d))
	pod, err := storagepod.FromID(client, d.Get("datastore_cluster_id").(string))
	if err != nil {
		return fmt.Errorf("error getting datastore cluster: %s", err)
	}

	err = storagepod.ReconfigureVM(client, vm, spec, pod)
	if err != nil {
		return fmt.Errorf("error reconfiguring VM on datastore cluster %q: %s", pod.Name(), err)
	}
	return nil
}

func resourceVSphereVirtualMachineDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Performing delete", resourceVSphereVirtualMachineIDString(d))
	client := meta.(*VSphereClient).vimClient
	id := d.Id()
	vm, err := virtualmachine.FromUUID(client, id)
	if err != nil {
		return fmt.Errorf("cannot locate virtual machine with UUID %q: %s", id, err)
	}
	vprops, err := virtualmachine.Properties(vm)
	if err != nil {
		return fmt.Errorf("error fetching VM properties: %s", err)
	}
	// Shutdown the VM first. We do attempt a graceful shutdown for the purpose
	// of catching any edge data issues with associated virtual disks that we may
	// need to retain on delete. However, we ignore the user-set force shutdown
	// flag.
	if vprops.Runtime.PowerState != types.VirtualMachinePowerStatePoweredOff {
		timeout := d.Get("shutdown_wait_timeout").(int)
		if err := virtualmachine.GracefulPowerOff(client, vm, timeout, true); err != nil {
			return fmt.Errorf("error shutting down virtual machine: %s", err)
		}
	}
	// Now attempt to detach any virtual disks that may need to be preserved.
	devices := object.VirtualDeviceList(vprops.Config.Hardware.Device)
	spec := types.VirtualMachineConfigSpec{}
	if spec.DeviceChange, err = virtualdevice.DiskDestroyOperation(d, client, devices); err != nil {
		return err
	}
	// Only run the reconfigure operation if there's actually disks in the spec.
	if len(spec.DeviceChange) > 0 {
		if err := virtualmachine.Reconfigure(vm, spec); err != nil {
			return fmt.Errorf("error detaching virtual disks: %s", err)
		}
	}

	// The final operation here is to destroy the VM.
	if err := virtualmachine.Destroy(vm); err != nil {
		return fmt.Errorf("error destroying virtual machine: %s", err)
	}
	d.SetId("")
	log.Printf("[DEBUG] %s: Delete complete", resourceVSphereVirtualMachineIDString(d))
	return nil
}

func resourceVSphereVirtualMachineCustomizeDiff(d *schema.ResourceDiff, meta interface{}) error {
	log.Printf("[DEBUG] %s: Performing diff customization and validation", resourceVSphereVirtualMachineIDString(d))
	client := meta.(*VSphereClient).vimClient

	// Block certain options from being set depending on the vSphere version.
	version := viapi.ParseVersionFromClient(client)
	if d.Get("efi_secure_boot_enabled").(bool) {
		if version.Older(viapi.VSphereVersion{Product: version.Product, Major: 6, Minor: 5}) {
			return fmt.Errorf("efi_secure_boot_enabled is only supported on vSphere 6.5 and higher")
		}
	}

	// Validate cdrom sub-resources when not deploying from ovf
	if len(d.Get("ovf_deploy").([]interface{})) == 0 {
		if err := virtualdevice.CdromDiffOperation(d, client); err != nil {
			return err
		}
	}

	if len(d.Get("ovf_deploy").([]interface{})) == 0 && len(d.Get("network_interface").([]interface{})) == 0 {
		return fmt.Errorf("network_interface parameter is required when not deploying from ovf template")

	}
	// Validate network device sub-resources
	if err := virtualdevice.NetworkInterfaceDiffOperation(d, client); err != nil {
		return err
	}

	// Process changes to resource pool
	if err := resourceVSphereVirtualMachineCustomizeDiffResourcePoolOperation(d); err != nil {
		return err
	}

	// Normalize datastore cluster vs datastore
	if err := datastoreClusterDiffOperation(d, client); err != nil {
		return err
	}

	// Validate and normalize disk sub-resources when not deploying from ovf
	if len(d.Get("ovf_deploy").([]interface{})) == 0 {
		if len(d.Get("clone").([]interface{})) == 0 {
			if err := virtualdevice.DiskDiffOperation(d, client); err != nil {
				return err
			}
		}
	}
	// When a VM is a member of a vApp container, it is no longer part of the VM
	// tree, and therefore cannot have its VM folder set.
	if _, ok := d.GetOk("folder"); ok && vappcontainer.IsVApp(client, d.Get("resource_pool_id").(string)) {
		return fmt.Errorf("cannot set folder while VM is in a vApp container")
	}

	if len(d.Get("ovf_deploy").([]interface{})) == 0 && d.Get("datacenter_id").(string) != "" {
		return fmt.Errorf("data center id is to be set only when deploying from ovf")
	}

	if len(d.Get("ovf_deploy").([]interface{})) > 0 {

		localOvfPath := d.Get("ovf_deploy.0.local_ovf_path").(string)
		remoteOvfUrl := d.Get("ovf_deploy.0.remote_ovf_url").(string)
		datacenterId := d.Get("datacenter_id").(string)
		dataStoreId := d.Get("datastore_id").(string)
		hostSystemId := d.Get("host_system_id").(string)

		if localOvfPath == "" && remoteOvfUrl == "" {
			return fmt.Errorf("either local ovf path or remote ovf url is required, both can't be empty")
		}
		if datacenterId == "" {
			return fmt.Errorf("data center ID is required for OVF deployment")
		}
		if dataStoreId == "" {
			return fmt.Errorf("data store ID is required for OVF deployment")
		}
		if hostSystemId == "" {
			return fmt.Errorf("host system ID is required for OVF deployment")
		}
		if localOvfPath != "" {
			if _, err := os.Stat(localOvfPath); os.IsNotExist(err) {
				return fmt.Errorf("file doesn't exist %s", localOvfPath)
			}
		}
	}

	// If this is a new resource and we are cloning, perform all clone validation
	// operations.
	if len(d.Get("clone").([]interface{})) > 0 {
		if err := viapi.ValidateVirtualCenter(client); err != nil {
			return errors.New("use of the clone sub-resource block requires vCenter")
		}

		switch {
		case d.Get("imported").(bool):
			// Imported workflows need to have the configuration of the clone
			// sub-resource block persisted to state without forcing a new resource.
			// Any changes after that will be properly tracked as a ForceNew, by
			// flagging the imported flag to off.
			d.SetNew("imported", false)
		case d.Id() == "":
			if contentlibrary.IsContentLibraryItem(meta.(*VSphereClient).restClient, d.Get("clone.0.template_uuid").(string)) {
				if _, ok := d.GetOk("datastore_cluster_id"); ok {
					return fmt.Errorf("Cannot use datastore_cluster_id with Content Library source")
				}
			} else if err := vmworkflow.ValidateVirtualMachineClone(d, client); err != nil {
				return err
			}
			fallthrough
		default:
			// For most cases (all non-imported workflows), any changed attribute in
			// the clone configuration namespace is a ForceNew. Flag those now.
			for _, k := range d.GetChangedKeysPrefix("clone.0") {
				if strings.HasSuffix(k, ".#") {
					k = strings.TrimSuffix(k, ".#")
				}
				// To maintain consistency with other timeout options, timeout does not
				// need to ForceNew
				if k == "clone.0.timeout" {
					continue
				}
				d.ForceNew(k)
			}
		}
	}

	// Validate hardware version changes.
	cv, tv := d.GetChange("hardware_version")
	virtualmachine.ValidateHardwareVersion(cv.(int), tv.(int))

	// Validate that the config has the necessary components for vApp support.
	// Note that for clones the data is prepopulated in
	// ValidateVirtualMachineClone.
	if err := virtualdevice.VerifyVAppTransport(d, client); err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Diff customization and validation complete", resourceVSphereVirtualMachineIDString(d))
	return nil
}

func resourceVSphereVirtualMachineCustomizeDiffResourcePoolOperation(d *schema.ResourceDiff) error {
	if d.HasChange("resource_pool_id") && !d.HasChange("host_system_id") {
		log.Printf(
			"[DEBUG] %s: resource_pool_id modified without change to host_system_id, marking as computed",
			resourceVSphereVirtualMachineIDString(d),
		)
		if err := d.SetNewComputed("host_system_id"); err != nil {
			return err
		}
	}
	return nil
}

func datastoreClusterDiffOperation(d *schema.ResourceDiff, client *govmomi.Client) error {
	if !structure.ValuesAvailable("", []string{"datastore_cluster_id", "datastore_id"}, d) {
		log.Printf("[DEBUG] DatastoreClusterDiffOperation: datastore_id or datastore_cluster_id value depends on a computed value from another resource. Skipping validation.")
		return nil
	}
	podID, podOk := d.GetOk("datastore_cluster_id")
	podKnown := d.NewValueKnown("datastore_cluster_id")
	dsID, dsOk := d.GetOk("datastore_id")
	dsKnown := d.NewValueKnown("datastore_id")

	switch {
	case podKnown && dsKnown && !podOk && !dsOk:
		// No root-level datastore option was available. This can happen on new
		// configs where the user has not supplied either option, so we need to
		// block this.
		return errors.New("one of datastore_id datastore_cluster_id must be specified")
	case podKnown && !podOk:
		// No datastore cluster
		return nil
	case !dsOk:
		// No datastore, we don't need to touch it
		return nil
	case !podKnown:
		// Datastore cluster ID changing but we don't know it yet. Mark the datastore ID as computed
		log.Printf("[DEBUG] %s: Datastore cluster ID unknown, marking VM datastore as computed", resourceVSphereVirtualMachineIDString(d))
		return d.SetNewComputed("datastore_id")
	}

	return datastoreClusterDiffOperationCheckMembership(d, client, podID.(string), dsID.(string))
}

func datastoreClusterDiffOperationCheckMembership(d *schema.ResourceDiff, client *govmomi.Client, podID, dsID string) error {
	log.Printf("[DEBUG] %s: Checking VM datastore cluster membership", resourceVSphereVirtualMachineIDString(d))

	// Determine if the current datastore from state is a member of the current
	// datastore cluster.
	pod, err := storagepod.FromID(client, podID)
	if err != nil {
		return fmt.Errorf("error fetching datastore cluster ID %q: %s", podID, err)
	}

	ds, err := datastore.FromID(client, dsID)
	if err != nil {
		return fmt.Errorf("error fetching datastore ID %q: %s", dsID, err)
	}

	isMember, err := storagepod.IsMember(pod, ds)
	if err != nil {
		return fmt.Errorf("error checking storage pod membership: %s", err)
	}
	if !isMember {
		// If the current datastore in state is not a member of the cluster, we
		// need to trigger a migration. Do this by setting the datastore ID to
		// computed so that it's picked up in the next update.
		log.Printf(
			"[DEBUG] %s: Datastore %q not a member of cluster %q, marking VM datastore as computed",
			resourceVSphereVirtualMachineIDString(d),
			ds.Name(),
			pod.Name(),
		)
		return d.SetNewComputed("datastore_id")
	}

	return nil
}

func resourceVSphereVirtualMachineImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*VSphereClient).vimClient

	name := d.Id()
	if name == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}

	log.Printf("[DEBUG] Looking for VM by name/path %q", name)
	vm, err := virtualmachine.FromPath(client, name, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching virtual machine: %s", err)
	}
	props, err := virtualmachine.Properties(vm)
	if err != nil {
		return nil, fmt.Errorf("error fetching virtual machine properties: %s", err)
	}

	// Block the import if the VM is a template.
	if props.Config.Template {
		return nil, fmt.Errorf("VM %q is a template and cannot be imported", name)
	}

	// Quickly walk the SCSI bus and determine the number of contiguous
	// controllers starting from bus number 0. This becomes the current SCSI
	// controller count. Anything past this is managed by config.
	log.Printf("[DEBUG] Determining number of SCSI controllers for VM %q", name)
	scsiBus := make([]bool, 4)
	for _, device := range props.Config.Hardware.Device {
		sc, ok := device.(types.BaseVirtualSCSIController)
		if !ok {
			continue
		}
		scsiBus[sc.GetVirtualSCSIController().BusNumber] = true
	}
	var ctlrCnt int
	for _, v := range scsiBus {
		if !v {
			break
		}
		ctlrCnt++
	}
	if ctlrCnt < 1 {
		return nil, fmt.Errorf("VM %q has no SCSI controllers", name)
	}
	d.Set("scsi_controller_count", ctlrCnt)

	// Validate the disks in the VM to make sure that they will work with the
	// resource. This is mainly ensuring that all disks are SCSI disks, but a
	// Read operation is attempted as well to make sure it will survive that.
	if err := virtualdevice.DiskImportOperation(d, client, object.VirtualDeviceList(props.Config.Hardware.Device)); err != nil {
		return nil, err
	}
	// The VM should be ready for reading now
	log.Printf("[DEBUG] VM UUID for %q is %q", name, props.Config.Uuid)
	d.SetId(props.Config.Uuid)
	d.Set("imported", true)

	// Set some defaults. This helps possibly prevent diffs where these values
	// have not been changed.
	rs := resourceVSphereVirtualMachine().Schema
	d.Set("force_power_off", rs["force_power_off"].Default)
	d.Set("migrate_wait_timeout", rs["migrate_wait_timeout"].Default)
	d.Set("shutdown_wait_timeout", rs["shutdown_wait_timeout"].Default)
	d.Set("wait_for_guest_ip_timeout", rs["wait_for_guest_ip_timeout"].Default)
	d.Set("wait_for_guest_net_timeout", rs["wait_for_guest_net_timeout"].Default)
	d.Set("wait_for_guest_net_routable", rs["wait_for_guest_net_routable"].Default)
	d.Set("poweron_timeout", rs["poweron_timeout"].Default)

	log.Printf("[DEBUG] %s: Import complete, resource is ready for read", resourceVSphereVirtualMachineIDString(d))
	return []*schema.ResourceData{d}, nil
}

// resourceVSphereVirtualMachineCreateBare contains the "bare metal" VM
// deploy path. The VM is returned.
func resourceVSphereVirtualMachineCreateBare(d *schema.ResourceData, meta interface{}) (*object.VirtualMachine, error) {
	log.Printf("[DEBUG] %s: VM being created from scratch", resourceVSphereVirtualMachineIDString(d))
	client := meta.(*VSphereClient).vimClient
	poolID := d.Get("resource_pool_id").(string)
	pool, err := resourcepool.FromID(client, poolID)
	if err != nil {
		return nil, fmt.Errorf("could not find resource pool ID %q: %s", poolID, err)
	}

	// Find the folder based off the path to the resource pool. Basically what we
	// are saying here is that the VM folder that we are placing this VM in needs
	// to be in the same hierarchy as the resource pool - so in other words, the
	// same datacenter.
	fo, err := folder.VirtualMachineFolderFromObject(client, pool, d.Get("folder").(string))
	if err != nil {
		return nil, err
	}
	var hs *object.HostSystem
	if v, ok := d.GetOk("host_system_id"); ok {
		hsID := v.(string)
		var err error
		if hs, err = hostsystem.FromID(client, hsID); err != nil {
			return nil, fmt.Errorf("error locating host system at ID %q: %s", hsID, err)
		}
	}

	// Validate that the host is part of the resource pool before proceeding
	if err := resourcepool.ValidateHost(client, pool, hs); err != nil {
		return nil, err
	}

	// Ready to start making the VM here. First expand our main config spec.
	spec, err := expandVirtualMachineConfigSpec(d, client)
	if err != nil {
		return nil, fmt.Errorf("error in virtual machine configuration: %s", err)
	}

	// Now we need to get the default device set - this is available in the
	// environment info in the resource pool, which we can then filter through
	// our device CRUD lifecycles to get a full deviceChange attribute for our
	// configspec.
	devices, err := resourcepool.DefaultDevices(client, pool, d.Get("guest_id").(string))
	if err != nil {
		return nil, fmt.Errorf("error loading default device list: %s", err)
	}
	log.Printf("[DEBUG] Default devices: %s", virtualdevice.DeviceListString(devices))

	if spec.DeviceChange, err = applyVirtualDevices(d, client, devices); err != nil {
		return nil, err
	}

	// Create the VM according the right API path - if we have a datastore
	// cluster, use the SDRS API, if not, use the standard API.
	var vm *object.VirtualMachine
	if _, ok := d.GetOk("datastore_cluster_id"); ok {
		vm, err = resourceVSphereVirtualMachineCreateBareWithSDRS(d, meta, fo, spec, pool, hs)
	} else {
		vm, err = resourceVSphereVirtualMachineCreateBareStandard(d, meta, fo, spec, pool, hs)
	}
	if err != nil {
		return nil, err
	}

	// VM is created. Set the ID now before proceeding, in case the rest of the
	// process here fails.
	vprops, err := virtualmachine.Properties(vm)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch properties of created virtual machine: %s", err)
	}
	log.Printf("[DEBUG] VM %q - UUID is %q", vm.InventoryPath, vprops.Config.Uuid)
	d.SetId(vprops.Config.Uuid)

	pTimeoutStr := fmt.Sprintf("%ds", d.Get("poweron_timeout").(int))
	pTimeout, err := time.ParseDuration(pTimeoutStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse poweron_timeout as a valid duration: %s", err)
	}
	// Start the virtual machine
	if err := virtualmachine.PowerOn(vm, pTimeout); err != nil {
		return nil, fmt.Errorf("error powering on virtual machine: %s", err)
	}
	return vm, nil
}

// resourceVSphereVirtualMachineCreateBareWithSDRS runs the creation part of
// resourceVSphereVirtualMachineCreateBare through storage DRS. It's designed
// to be run when a storage cluster is specified, versus simply specifying
// datastores.
func resourceVSphereVirtualMachineCreateBareWithSDRS(
	d *schema.ResourceData,
	meta interface{},
	fo *object.Folder,
	spec types.VirtualMachineConfigSpec,
	pool *object.ResourcePool,
	hs *object.HostSystem,
) (*object.VirtualMachine, error) {
	client := meta.(*VSphereClient).vimClient
	if err := viapi.ValidateVirtualCenter(client); err != nil {
		return nil, fmt.Errorf("connection ineligible to use datastore_cluster_id: %s", err)
	}

	log.Printf("[DEBUG] %s: Creating virtual machine through Storage DRS API", resourceVSphereVirtualMachineIDString(d))
	pod, err := storagepod.FromID(client, d.Get("datastore_cluster_id").(string))
	if err != nil {
		return nil, fmt.Errorf("error getting datastore cluster: %s", err)
	}

	vm, err := storagepod.CreateVM(client, fo, spec, pool, hs, pod)
	if err != nil {
		return nil, fmt.Errorf("error creating virtual machine on datastore cluster %q: %s", pod.Name(), err)
	}

	return vm, nil
}

// resourceVSphereVirtualMachineCreateBareStandard performs the steps necessary
// during resourceVSphereVirtualMachineCreateBare to create a virtual machine
// when a datastore cluster is not supplied.
func resourceVSphereVirtualMachineCreateBareStandard(
	d *schema.ResourceData,
	meta interface{},
	fo *object.Folder,
	spec types.VirtualMachineConfigSpec,
	pool *object.ResourcePool,
	hs *object.HostSystem,
) (*object.VirtualMachine, error) {
	client := meta.(*VSphereClient).vimClient

	// Set the datastore for the VM.
	ds, err := datastore.FromID(client, d.Get("datastore_id").(string))
	if err != nil {
		return nil, fmt.Errorf("error locating datastore for VM: %s", err)
	}
	spec.Files = &types.VirtualMachineFileInfo{
		VmPathName: fmt.Sprintf("[%s]", ds.Name()),
	}

	vm, err := virtualmachine.Create(client, fo, spec, pool, hs)
	if err != nil {
		return nil, fmt.Errorf("error creating virtual machine: %s", err)
	}
	return vm, nil
}

// Deploy vm from OVF template
func resourceVsphereMachineDeployOVF(d *schema.ResourceData, meta interface{}) (*object.VirtualMachine, error) {

	localOvfPath := d.Get("ovf_deploy.0.local_ovf_path").(string)
	remoteOvfUrl := d.Get("ovf_deploy.0.remote_ovf_url").(string)
	ovFromLocal := true
	if localOvfPath == "" {
		ovFromLocal = false
	}
	log.Printf("[DEBUG] %s:%s VM is being deployed from OVF template ", localOvfPath, remoteOvfUrl)

	client := meta.(*VSphereClient).vimClient
	name := d.Get("name").(string)
	var vm *object.VirtualMachine

	poolID := d.Get("resource_pool_id").(string)
	poolObj, err := resourcepool.FromID(client, poolID)
	if err != nil {
		return nil, fmt.Errorf("could not find resource pool ID %q: %s", poolID, err)
	}
	resourcePoolMor := poolObj.Reference()

	folderObj, err := folder.VirtualMachineFolderFromObject(client, poolObj, d.Get("folder").(string))
	if err != nil {
		return nil, err
	}

	hostId := d.Get("host_system_id").(string)
	hostObj, err := hostsystem.FromID(client, hostId)
	if err != nil {
		return nil, fmt.Errorf("could not find host with ID %q: %s", hostId, err)
	}
	hostMor := hostObj.Reference()

	dsId := d.Get("datastore_id").(string)
	dsObj, err := datastore.FromID(client, dsId)
	if err != nil {
		return nil, fmt.Errorf("could not find datastore with ID %q: %s", dsId, err)
	}
	dsMor := dsObj.Reference()

	networkMapping, err := ovfdeploy.GetNetworkMapping(client, d)
	if err != nil {
		return nil, err
	}
	importSpecParam := types.OvfCreateImportSpecParams{
		EntityName:         name,
		HostSystem:         &hostMor,
		NetworkMapping:     networkMapping,
		IpAllocationPolicy: d.Get("ovf_deploy.0.ip_allocation_policy").(string),
		IpProtocol:         d.Get("ovf_deploy.0.ip_protocol").(string),
		DiskProvisioning:   d.Get("ovf_deploy.0.disk_provisioning").(string),
	}

	ovfDescriptor := ""
	if ovFromLocal {
		fileBuffer, err := ioutil.ReadFile(localOvfPath)
		if err != nil {
			return nil, err
		}
		ovfDescriptor = string(fileBuffer)
	} else {
		resp, err := http.Get(remoteOvfUrl)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			ovfDescriptor = string(bodyBytes)
		}
	}
	if ovfDescriptor == "" {
		return nil, fmt.Errorf("error while reading the OVF File %s %s", localOvfPath, remoteOvfUrl)
	}

	ovfManager := ovf.NewManager(client.Client)
	ovfCreateImportSpecResult, err := ovfManager.CreateImportSpec(context.Background(), ovfDescriptor,
		resourcePoolMor, dsMor, importSpecParam)
	if err != nil {
		return nil, err
	}

	ovfPath := localOvfPath
	if !ovFromLocal {
		ovfPath = remoteOvfUrl
	}

	log.Print(" [DEBUG] start deploying from OVF Template")
	err = ovfdeploy.DeployOVFAndGetResult(ovfCreateImportSpecResult, poolObj, folderObj, hostObj, ovfPath, ovFromLocal)
	if err != nil {
		return nil, fmt.Errorf("error while importing OVF template %s", err)
	}

	dataCenterId := d.Get("datacenter_id").(string)
	datacenterObj, err := datacenterFromID(client, dataCenterId)
	if err != nil {
		return nil, fmt.Errorf("error while getting datacenter with id %s %s", dataCenterId, err)
	}
	vm, err = virtualmachine.FromPath(client, name, datacenterObj)
	if err != nil {
		return nil, fmt.Errorf("error while fetching the created vm %s", err)
	}

	// set ID for the vm
	vprops, err := virtualmachine.Properties(vm)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch properties of created virtual machine: %s", err)
	}

	log.Printf("[DEBUG] VM %q - UUID is %q", vm.InventoryPath, vprops.Config.Uuid)
	d.SetId(vprops.Config.Uuid)
	pTimeoutStr := fmt.Sprintf("%ds", d.Get("poweron_timeout").(int))
	pTimeout, err := time.ParseDuration(pTimeoutStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse poweron_timeout as a valid duration: %s", err)
	}
	// Start the virtual machine
	if err := virtualmachine.PowerOn(vm, pTimeout); err != nil {
		return nil, fmt.Errorf("error powering on virtual machine: %s", err)
	}
	return vm, nil
}

// resourceVSphereVirtualMachineCreateClone contains the clone VM deploy
// path. The VM is returned.
func resourceVSphereVirtualMachineCreateClone(d *schema.ResourceData, meta interface{}) (*object.VirtualMachine, error) {
	log.Printf("[DEBUG] %s: VM being created from clone", resourceVSphereVirtualMachineIDString(d))
	client := meta.(*VSphereClient).vimClient

	// Find the folder based off the path to the resource pool. Basically what we
	// are saying here is that the VM folder that we are placing this VM in needs
	// to be in the same hierarchy as the resource pool - so in other words, the
	// same datacenter.
	poolID := d.Get("resource_pool_id").(string)
	pool, err := resourcepool.FromID(client, poolID)
	if err != nil {
		return nil, fmt.Errorf("could not find resource pool ID %q: %s", poolID, err)
	}
	fo, err := folder.VirtualMachineFolderFromObject(client, pool, d.Get("folder").(string))
	if err != nil {
		return nil, err
	}

	// Start the clone
	name := d.Get("name").(string)
	timeout := d.Get("clone.0.timeout").(int)
	var vm *object.VirtualMachine
	if contentlibrary.IsContentLibraryItem(meta.(*VSphereClient).restClient, d.Get("clone.0.template_uuid").(string)) {
		// Clone source is an item from a Content Library. Prepare required resources.

		// First, if a host is set, check that it is valid.
		host := &object.HostSystem{}
		if hid := d.Get("host_system_id").(string); hid != "" {
			host, err = hostsystem.FromID(client, hid)
			if err != nil {
				return nil, err
			}
		}

		// Get OVF mappings for NICs
		nm := contentlibrary.MapNetworkDevices(d)

		// Validate the specified datastore
		dsID := d.Get("datastore_id")
		ds, err := datastore.FromID(client, dsID.(string))
		if err != nil {
			return nil, err
		}
		if err != nil {
			return nil, err
		}

		spId := d.Get("storage_policy_id").(string)

		dd := virtualmachine.DeployDest(name, d.Get("annotation").(string), pool, host, fo, ds, spId, nm)

		rclient := meta.(*VSphereClient).restClient
		item, err := contentlibrary.ItemFromID(rclient, d.Get("clone.0.template_uuid").(string))
		if err != nil {
			return nil, err
		}
		vmoid, err := virtualmachine.Deploy(rclient, item, dd, 30)
		if err != nil {
			return nil, err
		}
		vm, err = virtualmachine.FromMOID(client, vmoid.Value)
		if err != nil {
			return nil, err
		}
		// There is not currently a way to pull config values from Content Library items. If we do not send the values,
		// the defaults from the template will be used.
		d.Set("guest_id", "")
	} else {
		// Expand the clone spec. We get the source VM here too.
		cloneSpec, srcVM, err := vmworkflow.ExpandVirtualMachineCloneSpec(d, client)
		if err != nil {
			return nil, err
		}
		if _, ok := d.GetOk("datastore_cluster_id"); ok {
			vm, err = resourceVSphereVirtualMachineCreateCloneWithSDRS(d, meta, srcVM, fo, name, cloneSpec, timeout)
		} else {
			vm, err = virtualmachine.Clone(client, srcVM, fo, name, cloneSpec, timeout)
		}
		if err != nil {
			return nil, fmt.Errorf("error cloning virtual machine: %s", err)
		}
	}
	// The VM has been created. We still need to do post-clone configuration, and
	// while the resource should have an ID until this is done, we need it to go
	// through post-clone rollback workflows. All rollback functions will remove
	// the ID after it has done its rollback.
	//
	// It's generally safe to not rollback after the initial re-configuration is
	// fully complete and we move on to sending the customization spec.
	vprops, err := virtualmachine.Properties(vm)
	if err != nil {
		return nil, resourceVSphereVirtualMachineRollbackCreate(
			d,
			meta,
			vm,
			fmt.Errorf("cannot fetch properties of created virtual machine: %s", err),
		)
	}
	log.Printf("[DEBUG] VM %q - UUID is %q", vm.InventoryPath, vprops.Config.Uuid)
	d.SetId(vprops.Config.Uuid)

	// Before starting or proceeding any further, we need to normalize the
	// configuration of the newly cloned VM. This is basically a subset of update
	// with the stipulation that there is currently no state to help move this
	// along.
	cfgSpec, err := expandVirtualMachineConfigSpec(d, client)
	if err != nil {
		return nil, resourceVSphereVirtualMachineRollbackCreate(
			d,
			meta,
			vm,
			fmt.Errorf("error in virtual machine configuration: %s", err),
		)
	}

	// To apply device changes, we need the current devicecfgSpec from the config
	// info. We then filter this list through the same apply process we did for
	// create, which will apply the changes in an incremental fashion.
	devices := object.VirtualDeviceList(vprops.Config.Hardware.Device)
	var delta []types.BaseVirtualDeviceConfigSpec
	// First check the state of our SCSI bus. Normalize it if we need to.
	devices, delta, err = virtualdevice.NormalizeSCSIBus(devices, d.Get("scsi_type").(string), d.Get("scsi_controller_count").(int), d.Get("scsi_bus_sharing").(string))
	if err != nil {
		return nil, resourceVSphereVirtualMachineRollbackCreate(
			d,
			meta,
			vm,
			fmt.Errorf("error normalizing SCSI bus post-clone: %s", err),
		)
	}
	cfgSpec.DeviceChange = virtualdevice.AppendDeviceChangeSpec(cfgSpec.DeviceChange, delta...)
	// Disks
	devices, delta, err = virtualdevice.DiskPostCloneOperation(d, client, devices)
	if err != nil {
		return nil, resourceVSphereVirtualMachineRollbackCreate(
			d,
			meta,
			vm,
			fmt.Errorf("error processing disk changes post-clone: %s", err),
		)
	}
	cfgSpec.DeviceChange = virtualdevice.AppendDeviceChangeSpec(cfgSpec.DeviceChange, delta...)
	// Network devices
	devices, delta, err = virtualdevice.NetworkInterfacePostCloneOperation(d, client, devices)
	if err != nil {
		return nil, resourceVSphereVirtualMachineRollbackCreate(
			d,
			meta,
			vm,
			fmt.Errorf("error processing network device changes post-clone: %s", err),
		)
	}
	cfgSpec.DeviceChange = virtualdevice.AppendDeviceChangeSpec(cfgSpec.DeviceChange, delta...)
	// CDROM
	devices, delta, err = virtualdevice.CdromPostCloneOperation(d, client, devices)
	if err != nil {
		return nil, resourceVSphereVirtualMachineRollbackCreate(
			d,
			meta,
			vm,
			fmt.Errorf("error processing CDROM device changes post-clone: %s", err),
		)
	}
	cfgSpec.DeviceChange = virtualdevice.AppendDeviceChangeSpec(cfgSpec.DeviceChange, delta...)
	log.Printf("[DEBUG] %s: Final device list: %s", resourceVSphereVirtualMachineIDString(d), virtualdevice.DeviceListString(devices))
	log.Printf("[DEBUG] %s: Final device change cfgSpec: %s", resourceVSphereVirtualMachineIDString(d), virtualdevice.DeviceChangeString(cfgSpec.DeviceChange))

	// Perform updates
	if _, ok := d.GetOk("datastore_cluster_id"); ok {
		err = resourceVSphereVirtualMachineUpdateReconfigureWithSDRS(d, meta, vm, cfgSpec)
	} else {
		err = virtualmachine.Reconfigure(vm, cfgSpec)
	}
	if err != nil {
		return nil, resourceVSphereVirtualMachineRollbackCreate(
			d,
			meta,
			vm,
			fmt.Errorf("error reconfiguring virtual machine: %s", err),
		)
	}

	// Upgrade the VM's hardware version if needed.
	err = virtualmachine.SetHardwareVersion(vm, d.Get("hardware_version").(int))
	if err != nil {
		return nil, err
	}

	var cw *virtualMachineCustomizationWaiter
	// Send customization spec if any has been defined.
	if len(d.Get("clone.0.customize").([]interface{})) > 0 {
		family, err := resourcepool.OSFamily(client, pool, d.Get("guest_id").(string))
		if err != nil {
			return nil, fmt.Errorf("cannot find OS family for guest ID %q: %s", d.Get("guest_id").(string), err)
		}
		custSpec := vmworkflow.ExpandCustomizationSpec(d, family)
		cw = newVirtualMachineCustomizationWaiter(client, vm, d.Get("clone.0.customize.0.timeout").(int))
		if err := virtualmachine.Customize(vm, custSpec); err != nil {
			// Roll back the VMs as per the error handling in reconfigure.
			if derr := resourceVSphereVirtualMachineDelete(d, meta); derr != nil {
				return nil, fmt.Errorf(formatVirtualMachinePostCloneRollbackError, vm.InventoryPath, err, derr)
			}
			d.SetId("")
			return nil, fmt.Errorf("error sending customization spec: %s", err)
		}
	}
	// Finally time to power on the virtual machine!
	pTimeout := time.Duration(d.Get("poweron_timeout").(int)) * time.Second
	if err := virtualmachine.PowerOn(vm, pTimeout); err != nil {
		return nil, fmt.Errorf("error powering on virtual machine: %s", err)
	}
	// If we customized, wait on customization.
	if cw != nil {
		log.Printf("[DEBUG] %s: Waiting for VM customization to complete", resourceVSphereVirtualMachineIDString(d))
		<-cw.Done()
		if err := cw.Err(); err != nil {
			return nil, fmt.Errorf(formatVirtualMachineCustomizationWaitError, vm.InventoryPath, err)
		}
	}
	// Clone is complete and ready to return
	return vm, nil
}

// resourceVSphereVirtualMachineCreateCloneWithSDRS runs the clone part of
// resourceVSphereVirtualMachineCreateClone through storage DRS. It's designed
// to be run when a storage cluster is specified, versus simply specifying
// datastores.
func resourceVSphereVirtualMachineCreateCloneWithSDRS(
	d *schema.ResourceData,
	meta interface{},
	srcVM *object.VirtualMachine,
	fo *object.Folder,
	name string,
	spec types.VirtualMachineCloneSpec,
	timeout int,
) (*object.VirtualMachine, error) {
	client := meta.(*VSphereClient).vimClient
	if err := viapi.ValidateVirtualCenter(client); err != nil {
		return nil, fmt.Errorf("connection ineligible to use datastore_cluster_id: %s", err)
	}

	log.Printf("[DEBUG] %s: Cloning virtual machine through Storage DRS API", resourceVSphereVirtualMachineIDString(d))
	pod, err := storagepod.FromID(client, d.Get("datastore_cluster_id").(string))
	if err != nil {
		return nil, fmt.Errorf("error getting datastore cluster: %s", err)
	}

	vm, err := storagepod.CloneVM(client, srcVM, fo, name, spec, timeout, pod)
	if err != nil {
		return nil, fmt.Errorf("error cloning on datastore cluster %q: %s", pod.Name(), err)
	}

	return vm, nil
}

// resourceVSphereVirtualMachineRollbackCreate attempts to "roll back" a
// resource due to an error that happened post-create that will put the VM in a
// state where it cannot be worked with. This should only be done early on in
// the process, namely on clone operations between when the clone actually
// happens, and no later than after the initial post-clone update is complete.
//
// If the rollback fails, an error is displayed prompting the user to manually
// delete the virtual machine before trying again.
func resourceVSphereVirtualMachineRollbackCreate(
	d *schema.ResourceData,
	meta interface{},
	vm *object.VirtualMachine,
	origErr error,
) error {
	defer d.SetId("")
	// Updates are largely atomic, so more than likely no disks with
	// keep_on_remove were attached, but just in case, we run this through delete
	// to make sure to safely remove any disk that may have been attached as part
	// of this process if it was flagged as such.
	if err := resourceVSphereVirtualMachineDelete(d, meta); err != nil {
		return fmt.Errorf(formatVirtualMachinePostCloneRollbackError, vm.InventoryPath, origErr, err)
	}
	return fmt.Errorf("error reconfiguring virtual machine: %s", origErr)
}

// resourceVSphereVirtualMachineUpdateLocation manages vMotion. This includes
// the migration of a VM from one host to another, or from one datastore to
// another (storage vMotion).
//
// This function is responsible for building the top-level relocate spec. For
// disks, we call out to relocate functionality in the disk sub-resource.
func resourceVSphereVirtualMachineUpdateLocation(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Checking for pending migration operations", resourceVSphereVirtualMachineIDString(d))
	client := meta.(*VSphereClient).vimClient

	// A little bit of duplication of VM object data is done here to keep the
	// method signature lean.
	id := d.Id()
	vm, err := virtualmachine.FromUUID(client, id)
	if err != nil {
		return fmt.Errorf("cannot locate virtual machine with UUID %q: %s", id, err)
	}

	// Determine if we are performing any storage vMotion tasks. This will generate the relocators if there are any.
	vprops, err := virtualmachine.Properties(vm)
	if err != nil {
		return fmt.Errorf("error fetching VM properties: %s", err)
	}
	devices := object.VirtualDeviceList(vprops.Config.Hardware.Device)
	relocators, diskRelocateOK, err := virtualdevice.DiskMigrateRelocateOperation(d, client, devices)
	if err != nil {
		return err
	}
	// If we don't have any changes, stop here.
	if !d.HasChange("resource_pool_id") && !d.HasChange("host_system_id") && !d.HasChange("datastore_id") && !diskRelocateOK {
		log.Printf("[DEBUG] %s: No migration operations found", resourceVSphereVirtualMachineIDString(d))
		return nil
	}
	log.Printf("[DEBUG] %s: Migration operations found, proceeding with migration", resourceVSphereVirtualMachineIDString(d))

	// Fetch and validate pool and host
	poolID := d.Get("resource_pool_id").(string)
	pool, err := resourcepool.FromID(client, poolID)
	if err != nil {
		return fmt.Errorf("could not find resource pool ID %q: %s", poolID, err)
	}
	var hs *object.HostSystem
	if v, ok := d.GetOk("host_system_id"); ok {
		hsID := v.(string)
		var err error
		if hs, err = hostsystem.FromID(client, hsID); err != nil {
			return fmt.Errorf("error locating host system at ID %q: %s", hsID, err)
		}
		if err := resourcepool.ValidateHost(client, pool, hs); err != nil {
			return err
		}
	}

	// Start building the spec
	spec := types.VirtualMachineRelocateSpec{
		Pool: types.NewReference(pool.Reference()),
	}

	// Fetch the datastore only if a datastore_cluster is not set
	if _, ok := d.GetOk("datastore_cluster_id"); !ok {
		if dsID, ok := d.GetOk("datastore_id"); ok {
			ds, err := datastore.FromID(client, dsID.(string))
			if err != nil {
				return fmt.Errorf("error locating datastore for VM: %s", err)
			}
			spec.Datastore = types.NewReference(ds.Reference())
		}
	}

	if hs != nil {
		hsRef := hs.Reference()
		spec.Host = &hsRef
	}

	spec.Disk = relocators

	// Ready to perform migration
	timeout := d.Get("migrate_wait_timeout").(int)
	if _, ok := d.GetOk("datastore_cluster_id"); ok {
		err = resourceVSphereVirtualMachineUpdateLocationRelocateWithSDRS(d, meta, vm, spec, timeout)
	} else {
		err = virtualmachine.Relocate(vm, spec, timeout)
	}
	return err
}

// resourceVSphereVirtualMachineUpdateLocationRelocateWithSDRS runs the storage vMotion
// part of resourceVSphereVirtualMachineUpdateLocation through storage DRS.
// It's designed to be run when a storage cluster is specified, versus simply
// specifying datastores.
func resourceVSphereVirtualMachineUpdateLocationRelocateWithSDRS(
	d *schema.ResourceData,
	meta interface{},
	vm *object.VirtualMachine,
	spec types.VirtualMachineRelocateSpec,
	timeout int,
) error {
	client := meta.(*VSphereClient).vimClient
	if err := viapi.ValidateVirtualCenter(client); err != nil {
		return fmt.Errorf("connection ineligible to use datastore_cluster_id: %s", err)
	}

	log.Printf("[DEBUG] %s: Running virtual machine relocate Storage DRS API", resourceVSphereVirtualMachineIDString(d))
	pod, err := storagepod.FromID(client, d.Get("datastore_cluster_id").(string))
	if err != nil {
		return fmt.Errorf("error getting datastore cluster: %s", err)
	}

	err = storagepod.RelocateVM(client, vm, spec, timeout, pod)
	if err != nil {
		return fmt.Errorf("error running vMotion on datastore cluster %q: %s", pod.Name(), err)
	}
	return nil
}

// applyVirtualDevices is used by Create and Update to build a list of virtual
// device changes.
func applyVirtualDevices(d *schema.ResourceData, c *govmomi.Client, l object.VirtualDeviceList) ([]types.BaseVirtualDeviceConfigSpec, error) {
	// We filter this device list through each major device class' apply
	// operation. This will give us a final set of changes that will be our
	// deviceChange attribute.
	var spec, delta []types.BaseVirtualDeviceConfigSpec
	var err error
	// First check the state of our SCSI bus. Normalize it if we need to.
	l, delta, err = virtualdevice.NormalizeSCSIBus(l, d.Get("scsi_type").(string), d.Get("scsi_controller_count").(int), d.Get("scsi_bus_sharing").(string))
	if err != nil {
		return nil, err
	}
	if len(delta) > 0 {
		log.Printf("[DEBUG] %s: SCSI bus has changed and requires a VM restart", resourceVSphereVirtualMachineIDString(d))
		d.Set("reboot_required", true)
	}
	spec = virtualdevice.AppendDeviceChangeSpec(spec, delta...)
	// Disks
	l, delta, err = virtualdevice.DiskApplyOperation(d, c, l)
	if err != nil {
		return nil, err
	}
	spec = virtualdevice.AppendDeviceChangeSpec(spec, delta...)
	// Network devices
	l, delta, err = virtualdevice.NetworkInterfaceApplyOperation(d, c, l)
	if err != nil {
		return nil, err
	}
	spec = virtualdevice.AppendDeviceChangeSpec(spec, delta...)
	// CDROM
	l, delta, err = virtualdevice.CdromApplyOperation(d, c, l)
	if err != nil {
		return nil, err
	}
	spec = virtualdevice.AppendDeviceChangeSpec(spec, delta...)
	log.Printf("[DEBUG] %s: Final device list: %s", resourceVSphereVirtualMachineIDString(d), virtualdevice.DeviceListString(l))
	log.Printf("[DEBUG] %s: Final device change spec: %s", resourceVSphereVirtualMachineIDString(d), virtualdevice.DeviceChangeString(spec))
	return spec, nil
}

// resourceVSphereVirtualMachineIDString prints a friendly string for the
// vsphere_virtual_machine resource.
func resourceVSphereVirtualMachineIDString(d structure.ResourceIDStringer) string {
	return structure.ResourceIDString(d, "vsphere_virtual_machine")
}
