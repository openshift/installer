// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package satellite

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/container-services-go-sdk/satellitelinkv1"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMSatelliteLink() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSatelliteLinkRead,

		Schema: map[string]*schema.Schema{
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Location ID.",
			},
			"ws_endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ws endpoint of the location.",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Service instance associated with this location.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the location.",
			},
			"satellite_link_host": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Satellite Link hostname of the location.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Enabled/Disabled.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp of creation of location.",
			},
			"last_change": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp of latest modification of location.",
			},
			"performance": {
				Type: schema.TypeList,
				//MaxItems:    1,
				Computed:    true,
				Description: "The last performance data of the Location.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tunnels": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Tunnels number estbalished from the Location.",
						},
						"health_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tunnels health status based on the Tunnels number established. Down(0)/Critical(1)/Up(>=2).",
						},
						"avg_latency": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Average latency calculated form latency of each Connector between Tunnel Server, unit is ms. -1 means no Connector established Tunnel.",
						},
						"rx_bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Average Receive (to Cloud) Bandwidth of last two minutes, unit is Byte/s.",
						},
						"tx_bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Average Transmitted (to Location) Bandwidth of last two minutes, unit is Byte/s.",
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Average Tatal Bandwidth of last two minutes, unit is Byte/s.",
						},
						"connectors": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The last performance data of the Location read from each Connector.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"connector": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the connector reported the performance data.",
									},
									"latency": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Latency between Connector and the Tunnel Server it connected.",
									},
									"rx_bw": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Average Transmitted (to Location) Bandwidth of last two minutes read from the Connector, unit is Byte/s.",
									},
									"tx_bw": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Average Transmitted (to Location) Bandwidth of last two minutes read from the Connector, unit is Byte/s.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmSatelliteLinkRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	satelliteLinkClient, err := meta.(conns.ClientSession).SatellitLinkClientSession()
	if err != nil {
		return diag.FromErr(err)
	}

	getLinkOptions := &satellitelinkv1.GetLinkOptions{}

	getLinkOptions.SetLocationID(d.Get("location").(string))

	location, response, err := satelliteLinkClient.GetLinkWithContext(context, getLinkOptions)
	if err != nil {
		log.Printf("[DEBUG] GetLinkWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetLinkWithContext failed %s\n%s", err, response))
	}

	d.SetId(*location.LocationID)
	if err = d.Set("ws_endpoint", location.WsEndpoint); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting ws_endpoint: %s", err))
	}
	if err = d.Set("crn", location.Crn); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting crn: %s", err))
	}
	if err = d.Set("description", location.Desc); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting description: %s", err))
	}
	if err = d.Set("satellite_link_host", location.SatelliteLinkHost); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting satellite_link_host: %s", err))
	}
	if err = d.Set("status", location.Status); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting status: %s", err))
	}
	if err = d.Set("created_at", location.CreatedAt); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
	}
	if err = d.Set("last_change", location.LastChange); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting last_change: %s", err))
	}

	if location.Performance != nil {
		err = d.Set("performance", dataSourceLocationFlattenPerformance(*location.Performance))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting performance %s", err))
		}
	}

	return nil
}

func dataSourceLocationFlattenPerformance(result satellitelinkv1.LocationPerformance) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceLocationPerformanceToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceLocationPerformanceToMap(performanceItem satellitelinkv1.LocationPerformance) (performanceMap map[string]interface{}) {
	performanceMap = map[string]interface{}{}

	if performanceItem.Tunnels != nil {
		performanceMap["tunnels"] = performanceItem.Tunnels
	}
	if performanceItem.HealthStatus != nil {
		performanceMap["health_status"] = performanceItem.HealthStatus
	}
	if performanceItem.AvgLatency != nil {
		performanceMap["avg_latency"] = performanceItem.AvgLatency
	}
	if performanceItem.RxBandwidth != nil {
		performanceMap["rx_bandwidth"] = performanceItem.RxBandwidth
	}
	if performanceItem.TxBandwidth != nil {
		performanceMap["tx_bandwidth"] = performanceItem.TxBandwidth
	}
	if performanceItem.Bandwidth != nil {
		performanceMap["bandwidth"] = performanceItem.Bandwidth
	}
	if performanceItem.Connectors != nil {
		connectorsList := []map[string]interface{}{}
		for _, connectorsItem := range performanceItem.Connectors {
			connectorsList = append(connectorsList, dataSourceLocationPerformanceConnectorsToMap(connectorsItem))
		}
		performanceMap["connectors"] = connectorsList
	}

	return performanceMap
}

func dataSourceLocationPerformanceConnectorsToMap(connectorsItem satellitelinkv1.LocationPerformanceConnectorsItem) (connectorsMap map[string]interface{}) {
	connectorsMap = map[string]interface{}{}

	if connectorsItem.Connector != nil {
		connectorsMap["connector"] = connectorsItem.Connector
	}
	if connectorsItem.Latency != nil {
		connectorsMap["latency"] = connectorsItem.Latency
	}
	if connectorsItem.RxBW != nil {
		connectorsMap["rx_bw"] = connectorsItem.RxBW
	}
	if connectorsItem.TxBW != nil {
		connectorsMap["tx_bw"] = connectorsItem.TxBW
	}

	return connectorsMap
}
