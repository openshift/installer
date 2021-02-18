package vsphere

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/virtualmachine"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/virtualdevice"
	"github.com/vmware/govmomi/object"
)

func dataSourceVSphereVirtualMachine() *schema.Resource {
	s := map[string]*schema.Schema{
		"datacenter_id": {
			Type:        schema.TypeString,
			Description: "The managed object ID of the datacenter the virtual machine is in. This is not required when using ESXi directly, or if there is only one datacenter in your infrastructure.",
			Optional:    true,
		},
		"scsi_controller_scan_count": {
			Type:        schema.TypeInt,
			Description: "The number of SCSI controllers to scan for disk sizes and controller types on.",
			Optional:    true,
			Default:     1,
		},
		"sata_controller_scan_count": {
			Type:        schema.TypeInt,
			Description: "The number of SATA controllers to scan for disk sizes and controller types on.",
			Optional:    true,
			Default:     0,
		},
		"ide_controller_scan_count": {
			Type:        schema.TypeInt,
			Description: "The number of IDE controllers to scan for disk sizes and controller types on.",
			Optional:    true,
			Default:     2,
		},
		"scsi_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The common SCSI bus type of all controllers on the virtual machine.",
		},
		"scsi_bus_sharing": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Mode for sharing the SCSI bus.",
		},
		"disks": {
			Type:        schema.TypeList,
			Description: "Select configuration attributes from the disks on this virtual machine, sorted by bus and unit number.",
			Computed:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"size": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"eagerly_scrub": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"thin_provisioned": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"label": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"unit_number": {
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
		"network_interface_types": {
			Type:        schema.TypeList,
			Description: "The types of network interfaces found on the virtual machine, sorted by unit number.",
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"guest_ip_addresses": {
			Type:        schema.TypeList,
			Description: "The current list of IP addresses on this virtual machine.",
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}

	// Merge the VirtualMachineConfig structure so that we can include the number of
	// include the number of cpus, memory, firmware, disks, etc.
	structure.MergeSchema(s, schemaVirtualMachineConfigSpec())

	// Now that the schema has been composed and merged, we can attach our reader and
	// return the resource back to our host process.
	return &schema.Resource{
		Read:   dataSourceVSphereVirtualMachineRead,
		Schema: s,
	}
}

func dataSourceVSphereVirtualMachineRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*VSphereClient).vimClient

	name := d.Get("name").(string)
	log.Printf("[DEBUG] Looking for VM or template by name/path %q", name)
	var dc *object.Datacenter
	if dcID, ok := d.GetOk("datacenter_id"); ok {
		var err error
		dc, err = datacenterFromID(client, dcID.(string))
		if err != nil {
			return fmt.Errorf("cannot locate datacenter: %s", err)
		}
		log.Printf("[DEBUG] Datacenter for VM/template search: %s", dc.InventoryPath)
	}
	vm, err := virtualmachine.FromPath(client, name, dc)
	if err != nil {
		return fmt.Errorf("error fetching virtual machine: %s", err)
	}
	props, err := virtualmachine.Properties(vm)
	if err != nil {
		return fmt.Errorf("error fetching virtual machine properties: %s", err)
	}

	if props.Config == nil {
		return fmt.Errorf("no configuration returned for virtual machine %q", vm.InventoryPath)
	}

	if props.Config.Uuid == "" {
		return fmt.Errorf("virtual machine %q does not have a UUID", vm.InventoryPath)
	}

	// Read general VM config info
	if err := flattenVirtualMachineConfigInfo(d, props.Config); err != nil {
		return fmt.Errorf("error reading virtual machine configuration: %s", err)
	}

	d.SetId(props.Config.Uuid)
	d.Set("guest_id", props.Config.GuestId)
	d.Set("alternate_guest_name", props.Config.AlternateGuestName)
	d.Set("scsi_type", virtualdevice.ReadSCSIBusType(object.VirtualDeviceList(props.Config.Hardware.Device), d.Get("scsi_controller_scan_count").(int)))
	d.Set("scsi_bus_sharing", virtualdevice.ReadSCSIBusSharing(object.VirtualDeviceList(props.Config.Hardware.Device), d.Get("scsi_controller_scan_count").(int)))
	d.Set("firmware", props.Config.Firmware)
	disks, err := virtualdevice.ReadDiskAttrsForDataSource(object.VirtualDeviceList(props.Config.Hardware.Device), d)
	if err != nil {
		return fmt.Errorf("error reading disk sizes: %s", err)
	}
	nics, err := virtualdevice.ReadNetworkInterfaceTypes(object.VirtualDeviceList(props.Config.Hardware.Device))
	if err != nil {
		return fmt.Errorf("error reading network interface types: %s", err)
	}
	if err := d.Set("disks", disks); err != nil {
		return fmt.Errorf("error setting disk sizes: %s", err)
	}
	if err := d.Set("network_interface_types", nics); err != nil {
		return fmt.Errorf("error setting network interface types: %s", err)
	}
	if props.Guest != nil {
		if err := buildAndSelectGuestIPs(d, *props.Guest); err != nil {
			return fmt.Errorf("error setting guest IP addresses: %s", err)
		}
	}
	log.Printf("[DEBUG] VM search for %q completed successfully (UUID %q)", name, props.Config.Uuid)
	return nil
}
