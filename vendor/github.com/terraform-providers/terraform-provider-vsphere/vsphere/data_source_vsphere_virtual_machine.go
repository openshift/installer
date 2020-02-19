package vsphere

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/virtualmachine"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/virtualdevice"
	"github.com/vmware/govmomi/object"
)

func dataSourceVSphereVirtualMachine() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVSphereVirtualMachineRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name or path of the virtual machine.",
				Required:    true,
			},
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
			"guest_id": {
				Type:        schema.TypeString,
				Description: "The guest ID of the virtual machine.",
				Computed:    true,
			},
			"firmware": {
				Type:        schema.TypeString,
				Description: "The firmware type for this virtual machine.",
				Computed:    true,
			},
			"alternate_guest_name": {
				Type:        schema.TypeString,
				Description: "The alternate guest name of the virtual machine when guest_id is a non-specific operating system, like otherGuest.",
				Computed:    true,
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
		},
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

	d.SetId(props.Config.Uuid)
	d.Set("guest_id", props.Config.GuestId)
	d.Set("alternate_guest_name", props.Config.AlternateGuestName)
	d.Set("scsi_type", virtualdevice.ReadSCSIBusType(object.VirtualDeviceList(props.Config.Hardware.Device), d.Get("scsi_controller_scan_count").(int)))
	d.Set("scsi_bus_sharing", virtualdevice.ReadSCSIBusSharing(object.VirtualDeviceList(props.Config.Hardware.Device), d.Get("scsi_controller_scan_count").(int)))
	d.Set("firmware", props.Config.Firmware)
	disks, err := virtualdevice.ReadDiskAttrsForDataSource(object.VirtualDeviceList(props.Config.Hardware.Device), d.Get("scsi_controller_scan_count").(int))
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
