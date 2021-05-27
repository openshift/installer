// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	dl "github.com/IBM/networking-go-sdk/directlinkv1"

	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	dlPorts  = "ports"
	dlPortID = "port_id"
	dlCount  = "direct_link_count"
	dlLabel  = "label"
	// dlLocationDisplayName = "location_display_name"
	// dlLocationName        = "location_name"
	dlSupportedLinkSpeeds = "supported_link_speeds"
	dlProviderName        = "provider_name"
)

func dataSourceIBMDirectLinkPorts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMDirectLinkPortsRead,

		Schema: map[string]*schema.Schema{
			dlLocationName: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Direct Link location short name",
			},
			dlPorts: {

				Type:        schema.TypeList,
				Description: "Collection of direct link ports",
				Computed:    true,
				Elem: &schema.Resource{

					Schema: map[string]*schema.Schema{

						dlPortID: {
							Type:        schema.TypeString,
							Computed:    true,
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
				},
			},
		},
	}
}

func dataSourceIBMDirectLinkPortsRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := meta.(ClientSession).DirectlinkV1API()
	if err != nil {
		return err
	}

	start := ""
	allrecs := []dl.Port{}
	for {
		listPortsOptions := sess.NewListPortsOptions()
		if _, ok := d.GetOk(dlLocationName); ok {
			dlLocationName := d.Get(dlLocationName).(string)
			listPortsOptions.SetLocationName(dlLocationName)
		}
		if start != "" {
			listPortsOptions.Start = &start

		}

		response, resp, err := sess.ListPorts(listPortsOptions)
		if err != nil {
			log.Println("[WARN] Error listing dl ports", resp, err)
			return err
		}
		start = GetNext(response.Next)
		allrecs = append(allrecs, response.Ports...)
		if start == "" {
			break
		}
	}

	portCollections := make([]map[string]interface{}, 0)
	for _, port := range allrecs {
		portCollection := map[string]interface{}{}
		portCollection[dlPortID] = *port.ID
		portCollection[dlCount] = *port.DirectLinkCount
		portCollection[dlLabel] = *port.Label
		portCollection[dlLocationDisplayName] = *port.LocationDisplayName
		portCollection[dlLocationName] = *port.LocationName
		portCollection[dlProviderName] = *port.ProviderName
		speed := make([]interface{}, 0)
		for _, s := range port.SupportedLinkSpeeds {
			speed = append(speed, s)
		}
		portCollection[dlSupportedLinkSpeeds] = speed
		portCollections = append(portCollections, portCollection)
	}
	d.SetId(dataSourceIBMDirectLinkPortsReadID(d))
	d.Set(dlPorts, portCollections)
	return nil
}

func dataSourceIBMDirectLinkPortsReadID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
