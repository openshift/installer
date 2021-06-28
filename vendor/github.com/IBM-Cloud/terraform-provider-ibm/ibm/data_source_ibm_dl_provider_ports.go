// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	dlProviderV2 "github.com/IBM/networking-go-sdk/directlinkproviderv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"time"
)

func dataSourceIBMDirectLinkProviderPorts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMDirectLinkProviderPortsRead,

		Schema: map[string]*schema.Schema{

			dlPorts: {

				Type:        schema.TypeList,
				Description: "Collection of direct link provider ports",
				Computed:    true,
				Elem: &schema.Resource{

					Schema: map[string]*schema.Schema{

						dlPortID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Port ID",
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

func directlinkProviderClient(meta interface{}) (*dlProviderV2.DirectLinkProviderV2, error) {
	sess, err := meta.(ClientSession).DirectlinkProviderV2API()
	return sess, err
}
func dataSourceIBMDirectLinkProviderPortsRead(d *schema.ResourceData, meta interface{}) error {
	directLinkProvider, err := directlinkProviderClient(meta)
	if err != nil {
		return err
	}
	start := ""
	allrecs := []dlProviderV2.ProviderPort{}
	for {
		listPortsProviderOptions := directLinkProvider.NewListProviderPortsOptions()
		if start != "" {
			listPortsProviderOptions.Start = &start
		}

		ports, resp, err := directLinkProvider.ListProviderPorts(listPortsProviderOptions)
		if err != nil {
			log.Println("[WARN] Error listing dl provider ports", ports, resp, err)
			return err
		}
		start = GetNext(ports.Next)
		allrecs = append(allrecs, ports.Ports...)
		if start == "" {
			break
		}
	}
	portCollections := make([]map[string]interface{}, 0)
	for _, port := range allrecs {
		portCollection := map[string]interface{}{}
		portCollection[dlPortID] = *port.ID
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
	d.SetId(dataSourceIBMDirectLinkProviderPortsReadID(d))
	d.Set(dlPorts, portCollections)
	return nil
}

func dataSourceIBMDirectLinkProviderPortsReadID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
