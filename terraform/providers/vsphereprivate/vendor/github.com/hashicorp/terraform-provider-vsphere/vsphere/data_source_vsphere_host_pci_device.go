package vsphere

import (
	"log"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/hostsystem"
	"github.com/vmware/govmomi/vim25/types"
)

func dataSourceVSphereHostPciDevice() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVSphereHostPciDeviceRead,

		Schema: map[string]*schema.Schema{
			"host_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Managed Object ID of the host system.",
			},
			"name_regex": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A regular expression used to match the PCI device name.",
			},
			"class_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The hexadecimal value of the PCI device's class ID.",
			},
			"vendor_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The hexadecimal value of the PCI device's vendor ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the PCI device.",
			},
		},
	}
}

func dataSourceVSphereHostPciDeviceRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] DataHostPCIDev: Beginning PCI device lookup on %s", d.Get("host_id").(string))
	client := meta.(*Client).vimClient
	host, err := hostsystem.FromID(client, d.Get("host_id").(string))
	if err != nil {
		return err
	}
	hprops, err := hostsystem.Properties(host)
	if err != nil {
		return err
	}
	devices, err := matchName(d, hprops.Hardware.PciDevice)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] DataHostPCIDev: Looking for a device with matching class_id and vendor_id")
	for _, device := range devices {
		// Match the class_id if it is set.
		if class, exists := d.GetOk("class_id"); exists {
			classInt, err := strconv.ParseInt(class.(string), 16, 16)
			if err != nil {
				return err
			}
			if device.ClassId != int16(classInt) {
				continue
			}
		}
		// Now match the vendor_id if it is set.
		if vendor, exists := d.GetOk("vendor_id"); exists {
			vendorInt, err := strconv.ParseInt(vendor.(string), 16, 16)
			if err != nil {
				return err
			}
			if device.VendorId != int16(vendorInt) {
				continue
			}
		}
		classHex := strconv.FormatInt(int64(device.ClassId), 16)
		vendorHex := strconv.FormatInt(int64(device.VendorId), 16)
		d.SetId(device.Id)
		_ = d.Set("name", device.DeviceName)
		_ = d.Set("class_id", classHex)
		_ = d.Set("vendor_id", vendorHex)
		log.Printf("[DEBUG] DataHostPCIDev: Matching PCI device found: %s", device.DeviceName)
		return nil
	}
	return nil
}

func matchName(d *schema.ResourceData, devices []types.HostPciDevice) ([]types.HostPciDevice, error) {
	log.Printf("[DEBUG] DataHostPCIDev: Selecting devices which match name regex")
	matches := []types.HostPciDevice{}
	re, err := regexp.Compile(d.Get("name_regex").(string))
	if err != nil {
		return nil, err
	}
	for _, device := range devices {
		if re.Match([]byte(device.DeviceName)) {
			matches = append(matches, device)
		}
	}
	return matches, nil
}
