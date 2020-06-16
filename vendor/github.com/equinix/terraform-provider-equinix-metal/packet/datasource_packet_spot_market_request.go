package packet

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/packethost/packngo"
)

func dataSourcePacketSpotMarketRequest() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePacketSpotMarketRequestRead,

		Schema: map[string]*schema.Schema{
			"request_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"device_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
		Timeouts: resourceDefaultTimeouts,
	}
}
func dataSourcePacketSpotMarketRequestRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	id := d.Get("request_id").(string)

	smr, _, err := client.SpotMarketRequests.Get(id, &packngo.GetOptions{Includes: []string{"project", "devices", "facilities"}})
	if err != nil {
		err = friendlyError(err)
		if isNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	deviceIDs := make([]string, len(smr.Devices))
	for i, d := range smr.Devices {
		deviceIDs[i] = d.ID
	}
	d.Set("device_ids", deviceIDs)
	d.SetId(id + strings.Join(deviceIDs, "-"))
	return nil
}
