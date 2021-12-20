// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/container-services-go-sdk/kubernetesserviceapiv1"
	"github.com/IBM-Cloud/container-services-go-sdk/satellitelinkv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIbmSatelliteLink() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmSatelliteLinkCreate,
		ReadContext:   resourceIbmSatelliteLinkRead,
		UpdateContext: resourceIbmSatelliteLinkUpdate,
		DeleteContext: resourceIbmSatelliteLinkDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "CRN of the Location.",
			},
			"location": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Location ID.",
			},
			"ws_endpoint": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ws endpoint of the location.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the location.",
			},
			"satellite_link_host": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Satellite Link hostname of the location.",
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Enabled/Disabled.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp of creation of location.",
			},
			"last_change": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp of latest modification of location.",
			},
			"performance": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The last performance data of the Location.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tunnels": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Tunnels number estbalished from the Location.",
						},
						"health_status": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tunnels health status based on the Tunnels number established. Down(0)/Critical(1)/Up(>=2).",
						},
						"avg_latency": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Average latency calculated form latency of each Connector between Tunnel Server, unit is ms. -1 means no Connector established Tunnel.",
						},
						"rx_bandwidth": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Average Receive (to Cloud) Bandwidth of last two minutes, unit is Byte/s.",
						},
						"tx_bandwidth": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Average Transmitted (to Location) Bandwidth of last two minutes, unit is Byte/s.",
						},
						"bandwidth": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Average Tatal Bandwidth of last two minutes, unit is Byte/s.",
						},
						"connectors": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The last performance data of the Location read from each Connector.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"connector": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The name of the connector reported the performance data.",
									},
									"latency": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Latency between Connector and the Tunnel Server it connected.",
									},
									"rx_bw": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Average Transmitted (to Location) Bandwidth of last two minutes read from the Connector, unit is Byte/s.",
									},
									"tx_bw": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
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

func resourceIbmSatelliteLinkCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	satelliteLinkClient, err := meta.(ClientSession).SatellitLinkClientSession()
	if err != nil {
		return diag.FromErr(err)
	}

	createLinkOptions := &satellitelinkv1.CreateLinkOptions{}

	if _, ok := d.GetOk("crn"); ok {
		createLinkOptions.SetCrn(d.Get("crn").(string))
	}
	if _, ok := d.GetOk("location"); ok {
		createLinkOptions.SetLocationID(d.Get("location").(string))
	}

	location, response, err := satelliteLinkClient.CreateLinkWithContext(context, createLinkOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateLinkWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateLinkWithContext failed %s\n%s", err, response))
	}

	d.SetId(*location.LocationID)

	return resourceIbmSatelliteLinkUpdate(context, d, meta)
}

func resourceIbmSatelliteLinkRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	satelliteLinkClient, err := meta.(ClientSession).SatellitLinkClientSession()
	if err != nil {
		return diag.FromErr(err)
	}

	satClient, err := meta.(ClientSession).SatelliteClientSession()
	if err != nil {
		return diag.FromErr(fmt.Errorf("SatelliteClientSession failed %s\n", err))
	}

	getLinkOptions := &satellitelinkv1.GetLinkOptions{}

	getLinkOptions.SetLocationID(d.Id())

	link, response, err := satelliteLinkClient.GetLinkWithContext(context, getLinkOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetLinkWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetLinkWithContext failed %s\n%s", err, response))
	}

	getSatLocOptions := &kubernetesserviceapiv1.GetSatelliteLocationOptions{
		Controller: ptrToString(d.Id()),
	}

	locInstance, response, err := satClient.GetSatelliteLocation(getSatLocOptions)
	if err != nil || locInstance == nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("GetSatelliteLocation failed %s\n%s", err, response))
	}

	d.Set("crn", *locInstance.Crn)
	if err = d.Set("location", link.LocationID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting location: %s", err))
	}
	if err = d.Set("ws_endpoint", link.WsEndpoint); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting ws_endpoint: %s", err))
	}
	if err = d.Set("description", link.Desc); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
	}
	if err = d.Set("satellite_link_host", link.SatelliteLinkHost); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting satellite_link_host: %s", err))
	}
	if err = d.Set("status", link.Status); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting status: %s", err))
	}
	if err = d.Set("created_at", link.CreatedAt); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("last_change", link.LastChange); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting last_change: %s", err))
	}
	if link.Performance != nil {
		performanceMap := resourceIbmSatelliteLinkLocationPerformanceToMap(*link.Performance)
		d.Set("performance", []map[string]interface{}{performanceMap})
	}

	return nil
}

func resourceIbmSatelliteLinkLocationPerformanceToMap(locationPerformance satellitelinkv1.LocationPerformance) map[string]interface{} {
	locationPerformanceMap := map[string]interface{}{}

	locationPerformanceMap["tunnels"] = intValue(locationPerformance.Tunnels)
	if locationPerformance.HealthStatus != nil {
		locationPerformanceMap["health_status"] = *locationPerformance.HealthStatus
	}

	locationPerformanceMap["avg_latency"] = intValue(locationPerformance.AvgLatency)
	locationPerformanceMap["rx_bandwidth"] = intValue(locationPerformance.RxBandwidth)
	locationPerformanceMap["tx_bandwidth"] = intValue(locationPerformance.TxBandwidth)
	locationPerformanceMap["bandwidth"] = intValue(locationPerformance.Bandwidth)
	if locationPerformance.Connectors != nil {
		connectors := []map[string]interface{}{}
		for _, connectorsItem := range locationPerformance.Connectors {
			connectorsItemMap := resourceIbmSatelliteLinkLocationPerformanceConnectorsItemToMap(connectorsItem)
			connectors = append(connectors, connectorsItemMap)
		}
		locationPerformanceMap["connectors"] = connectors
	}

	return locationPerformanceMap
}

func resourceIbmSatelliteLinkLocationPerformanceConnectorsItemToMap(locationPerformanceConnectorsItem satellitelinkv1.LocationPerformanceConnectorsItem) map[string]interface{} {
	locationPerformanceConnectorsItemMap := map[string]interface{}{}

	if locationPerformanceConnectorsItem.Connector != nil {
		locationPerformanceConnectorsItemMap["connector"] = *locationPerformanceConnectorsItem.Connector
	}
	if locationPerformanceConnectorsItem.Latency != nil {
		locationPerformanceConnectorsItemMap["latency"] = intValue(locationPerformanceConnectorsItem.Latency)
	}
	if locationPerformanceConnectorsItem.RxBW != nil {
		locationPerformanceConnectorsItemMap["rx_bw"] = intValue(locationPerformanceConnectorsItem.RxBW)
	}
	if locationPerformanceConnectorsItem.TxBW != nil {
		locationPerformanceConnectorsItemMap["tx_bw"] = intValue(locationPerformanceConnectorsItem.TxBW)
	}

	log.Println("locationPerformanceConnectorsItemMap ::", locationPerformanceConnectorsItemMap)
	return locationPerformanceConnectorsItemMap
}

func resourceIbmSatelliteLinkUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	satelliteLinkClient, err := meta.(ClientSession).SatellitLinkClientSession()
	if err != nil {
		return diag.FromErr(err)
	}

	updateLinkOptions := &satellitelinkv1.UpdateLinkOptions{}
	updateLinkOptions.SetLocationID(d.Id())

	hasChange := false

	if d.HasChange("ws_endpoint") {
		updateLinkOptions.SetLocationID(d.Get("location").(string))
		updateLinkOptions.SetWsEndpoint(d.Get("ws_endpoint").(string))
		hasChange = true
	}

	if hasChange {
		_, response, err := satelliteLinkClient.UpdateLinkWithContext(context, updateLinkOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateLinkWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateLinkWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIbmSatelliteLinkRead(context, d, meta)
}

func resourceIbmSatelliteLinkDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	satelliteLinkClient, err := meta.(ClientSession).SatellitLinkClientSession()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteLinkOptions := &satellitelinkv1.DeleteLinkOptions{}

	deleteLinkOptions.SetLocationID(d.Id())

	_, response, err := satelliteLinkClient.DeleteLinkWithContext(context, deleteLinkOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteLinkWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteLinkWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}
