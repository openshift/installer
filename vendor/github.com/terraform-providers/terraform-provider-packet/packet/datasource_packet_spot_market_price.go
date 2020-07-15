package packet

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/packethost/packngo"
)

func dataSourceSpotMarketPrice() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePacketSpotMarketPriceRead,
		Schema: map[string]*schema.Schema{
			"facility": {
				Type:     schema.TypeString,
				Required: true,
			},
			"plan": {
				Type:     schema.TypeString,
				Required: true,
			},
			"price": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
		},
	}
}

func dataSourcePacketSpotMarketPriceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	facility := d.Get("facility").(string)
	plan := d.Get("plan").(string)

	prices, _, err := client.SpotMarket.Prices()
	if err != nil {
		return err
	}

	var price float64
	if fac, ok := prices[facility]; ok {
		if pri, ok := fac[plan]; ok {
			price = pri
		} else {
			return fmt.Errorf("Facility %s does not have prices for plan %s", facility, plan)
		}
	} else {
		return fmt.Errorf("There is no facility %s", facility)
	}
	d.Set("price", price)
	d.SetId(facility)
	return nil
}
