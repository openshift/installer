// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMDirectLinkPort() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMDirectLinkPortRead,
		Schema: map[string]*schema.Schema{
			dlPortID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Port ID",
			},
			dlCount: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Count of existing Direct Link gateways in this account on this port",
			},
			dlLabel: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Port Label",
			},
			dlLocationDisplayName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Port location long name",
			},
			dlLocationName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Port location name identifier",
			},
			dlProviderName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Port's provider name",
			},
			dlSupportedLinkSpeeds: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Computed:    true,
				Description: "Port's supported speeds in megabits per second",
			},
		},
	}
}

func dataSourceIBMDirectLinkPortRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := meta.(ClientSession).DirectlinkV1API()
	if err != nil {
		return err
	}

	getPortsOptions := sess.NewGetPortOptions(d.Get(dlPortID).(string))
	response, resp, err := sess.GetPort(getPortsOptions)
	if err != nil {
		log.Println("[WARN] Error getting port", resp, err)
		return err
	}

	d.SetId(*response.ID)
	d.Set(dlPortID, *response.ID)
	d.Set(dlCount, *response.DirectLinkCount)
	d.Set(dlLabel, *response.Label)
	d.Set(dlLocationDisplayName, *response.LocationDisplayName)
	d.Set(dlLocationName, *response.LocationName)
	d.Set(dlProviderName, *response.ProviderName)
	speed := make([]interface{}, 0)
	for _, s := range response.SupportedLinkSpeeds {
		speed = append(speed, s)
	}
	d.Set(dlSupportedLinkSpeeds, speed)
	return nil
}
