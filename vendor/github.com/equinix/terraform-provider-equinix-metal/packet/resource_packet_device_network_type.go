package packet

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/packethost/packngo"
)

func resourcePacketDeviceNetworkType() *schema.Resource {
	return &schema.Resource{
		Create: resourcePacketDeviceNetworkTypeCreate,
		Read:   resourcePacketDeviceNetworkTypeRead,
		Delete: resourcePacketDeviceNetworkTypeDelete,
		Update: resourcePacketDeviceNetworkTypeUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"device_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"layer3", "layer2-bonded", "layer2-individual", "hybrid"}, false),
			},
		},
	}
}

func getDevIDandNetworkType(d *schema.ResourceData, c *packngo.Client) (string, string, error) {
	deviceID := d.Id()
	if len(deviceID) == 0 {
		deviceID = d.Get("device_id").(string)
	}

	dev, _, err := c.Devices.Get(deviceID, nil)
	if err != nil {
		return "", "", err
	}
	devType := dev.GetNetworkType()

	return dev.ID, devType, nil
}

func getAndPossiblySetNetworkType(d *schema.ResourceData, c *packngo.Client, targetType string) error {
	devID, devType, err := getDevIDandNetworkType(d, c)
	if err != nil {
		return err
	}

	if devType != targetType {
		_, err := c.DevicePorts.DeviceToNetworkType(devID, targetType)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourcePacketDeviceNetworkTypeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	ntype := d.Get("type").(string)

	err := getAndPossiblySetNetworkType(d, client, ntype)
	if err != nil {
		return err
	}
	d.SetId(d.Get("device_id").(string))
	return resourcePacketDeviceNetworkTypeRead(d, meta)
}

func resourcePacketDeviceNetworkTypeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	_, devNType, err := getDevIDandNetworkType(d, client)

	if err != nil {
		err = friendlyError(err)

		if isNotFound(err) {
			log.Printf("[WARN] Device (%s) for Network Type request not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("type", devNType)
	return nil
}

func resourcePacketDeviceNetworkTypeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	ntype := d.Get("type").(string)
	if d.HasChange("type") {
		err := getAndPossiblySetNetworkType(d, client, ntype)
		if err != nil {
			return err
		}
	}

	return resourcePacketDeviceNetworkTypeRead(d, meta)
}

func resourcePacketDeviceNetworkTypeDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
