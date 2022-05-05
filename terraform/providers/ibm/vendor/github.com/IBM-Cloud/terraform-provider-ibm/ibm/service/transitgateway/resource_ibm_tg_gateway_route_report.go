// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package transitgateway

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/networking-go-sdk/transitgatewayapisv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	tgRouteReportId = "route_report_id"

	isTransitGatewayRouteReportPending = "pending"
	isTransitGatewayRouteReportDone    = "complete"
)

func ResourceIBMTransitGatewayRouteReport() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMTransitGatewayRouteReportCreate,
		Read:     resourceIBMTransitGatewayRouteReportRead,
		Delete:   resourceIBMTransitGatewayRouteReportDelete,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			tgGatewayId: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Transit Gateway identifier",
			},
			tgRouteReportId: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Transit Gateway Route Report identifier",
			},
			tgRouteReportConnections: {
				Type:        schema.TypeList,
				Description: "Collection of transit gateway connections",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						tgRouteReportConnectionBgps: {
							Type:        schema.TypeList,
							Description: "Collection of transit gateway connection's bgps",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									tgRouteReportConnectionBgpAsPath: {
										Type:     schema.TypeString,
										Computed: true,
									},
									tgRouteReportConnectionBgpIsUsed: {
										Type:     schema.TypeBool,
										Computed: true,
									},
									tgRouteReportConnectionBgpLocalPreference: {
										Type:     schema.TypeString,
										Computed: true,
									},
									tgRouteReportConnectionBgpPrefix: {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						ID: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgConnName: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgRouteReportConnectionType: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgRouteReportConnectionRoutes: {
							Type:        schema.TypeList,
							Description: "Collection of transit gateway connection's used routes",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									tgRouteReportConnectionRoutePrefix: {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			tgCreatedAt: {
				Type:     schema.TypeString,
				Computed: true,
			},
			tgRouteReportOverlappingRoutes: {
				Type:        schema.TypeList,
				Description: "Collection of transit gateway overlapping routes",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						tgRouteReportOverlappingRoutesDetail: {
							Type:        schema.TypeList,
							Description: "Collection of transit gateway overlapping route's details",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									tgConnectionId: {
										Type:     schema.TypeString,
										Computed: true,
									},
									tgRouteReportOverlappingPrefix: {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			tgStatus: {
				Type:     schema.TypeString,
				Computed: true,
			},
			tgUpdatedAt: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIBMTransitGatewayRouteReportCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := transitgatewayClient(meta)
	if err != nil {
		return err
	}

	createTransitGatewayRouteReportOptions := &transitgatewayapisv1.CreateTransitGatewayRouteReportOptions{}

	gatewayId := d.Get(tgGatewayId).(string)
	createTransitGatewayRouteReportOptions.SetTransitGatewayID(gatewayId)

	tgRouteReport, response, err := client.CreateTransitGatewayRouteReport(createTransitGatewayRouteReportOptions)
	if err != nil {
		return fmt.Errorf("Create Transit Gateway Route Report err %s\n%s", err, response)
	}

	d.SetId(fmt.Sprintf("%s/%s", gatewayId, *tgRouteReport.ID))
	d.Set(tgRouteReportId, *tgRouteReport.ID)

	_, err = isWaitForTransitGatewayRouteReportAvailable(client, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	return resourceIBMTransitGatewayRouteReportRead(d, meta)
}

func isWaitForTransitGatewayRouteReportAvailable(client *transitgatewayapisv1.TransitGatewayApisV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for transit gateway route report (%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isTransitGatewayRouteReportPending},
		Target:     []string{isTransitGatewayRouteReportDone, ""},
		Refresh:    isTransitGatewayRouteReportRefreshFunc(client, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isTransitGatewayRouteReportRefreshFunc(client *transitgatewayapisv1.TransitGatewayApisV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		parts, err := flex.IdParts(id)
		if err != nil {
			return nil, "", fmt.Errorf("Error Getting Transit Gateway Route Report: %s", err)
		}

		gatewayId := parts[0]
		ID := parts[1]

		getTransitGatewayRouteReportOptions := &transitgatewayapisv1.GetTransitGatewayRouteReportOptions{}
		getTransitGatewayRouteReportOptions.SetTransitGatewayID(gatewayId)
		getTransitGatewayRouteReportOptions.SetID(ID)

		routeReport, response, getRouteErr := client.GetTransitGatewayRouteReport(getTransitGatewayRouteReportOptions)
		if getRouteErr != nil {
			return nil, "", fmt.Errorf("Error Getting Transit Gateway Route Report: %s\n%s", err, response)
		}

		if *routeReport.Status == "complete" {
			return routeReport, isTransitGatewayRouteReportDone, nil
		}

		return routeReport, isTransitGatewayRouteReportPending, nil
	}
}

func resourceIBMTransitGatewayRouteReportRead(d *schema.ResourceData, meta interface{}) error {

	client, err := transitgatewayClient(meta)
	if err != nil {
		return err
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}

	gatewayId := parts[0]
	ID := parts[1]

	getTransitGatewayRouteReportOptions := &transitgatewayapisv1.GetTransitGatewayRouteReportOptions{}
	getTransitGatewayRouteReportOptions.SetTransitGatewayID(gatewayId)
	getTransitGatewayRouteReportOptions.SetID(ID)
	instance, response, err := client.GetTransitGatewayRouteReport(getTransitGatewayRouteReportOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting Transit Gateway Route Report (%s): %s\n%s", ID, err, response)
	}

	d.Set(tgRouteReportId, *instance.ID)
	d.Set(tgGatewayId, gatewayId)

	if instance.Status != nil {
		d.Set(tgStatus, *instance.Status)
	}
	if instance.CreatedAt != nil {
		d.Set(tgCreatedAt, instance.CreatedAt.String())
	}
	if instance.UpdatedAt != nil {
		d.Set(tgUpdatedAt, instance.UpdatedAt.String())
	}

	connections := make([]map[string]interface{}, 0)
	for _, connection := range instance.Connections {
		tgConn := map[string]interface{}{}

		// Set connection info
		if connection.ID != nil {
			tgConn[ID] = *connection.ID
		}
		if connection.Name != nil {
			tgConn[tgConnName] = *connection.Name
		}
		if connection.Type != nil {
			tgConn[tgRouteReportConnectionType] = *connection.Type
		}

		// set bgps
		bgps := make([]map[string]interface{}, 0)
		for _, bgp := range connection.Bgps {
			tgConnBgp := map[string]interface{}{}

			if bgp.AsPath != nil {
				tgConnBgp[tgRouteReportConnectionBgpAsPath] = *bgp.AsPath
			}
			if bgp.IsUsed != nil {
				tgConnBgp[tgRouteReportConnectionBgpIsUsed] = *bgp.IsUsed
			}
			if bgp.LocalPreference != nil {
				tgConnBgp[tgRouteReportConnectionBgpLocalPreference] = *bgp.LocalPreference
			}
			if bgp.Prefix != nil {
				tgConnBgp[tgRouteReportConnectionBgpPrefix] = *bgp.Prefix
			}

			bgps = append(bgps, tgConnBgp)
		}
		tgConn[tgRouteReportConnectionBgps] = bgps

		// Set connection routes
		routes := make([]map[string]interface{}, 0)
		for _, route := range connection.Routes {
			tgConnRoute := map[string]interface{}{}

			if route.Prefix != nil {
				tgConnRoute[tgRouteReportConnectionRoutePrefix] = *route.Prefix
			}

			routes = append(routes, tgConnRoute)
		}
		tgConn[tgRouteReportConnectionRoutes] = routes

		connections = append(connections, tgConn)
	}
	d.Set(tgRouteReportConnections, connections)

	// Set overlapping route details
	overlappingRoutes := make([]map[string]interface{}, 0)
	for _, overlap := range instance.OverlappingRoutes {
		tgRoutes := map[string]interface{}{}

		routes := make([]map[string]interface{}, 0)
		for _, routeDetail := range overlap.Routes {
			tgRoutesDetail := map[string]interface{}{}

			if routeDetail.ConnectionID != nil {
				tgRoutesDetail[tgConnectionId] = *routeDetail.ConnectionID
			}
			if routeDetail.Prefix != nil {
				tgRoutesDetail[tgRouteReportOverlappingPrefix] = *routeDetail.Prefix
			}

			routes = append(routes, tgRoutesDetail)
		}
		tgRoutes[tgRouteReportOverlappingRoutesDetail] = routes

		overlappingRoutes = append(overlappingRoutes, tgRoutes)
	}
	d.Set(tgRouteReportOverlappingRoutes, overlappingRoutes)

	return nil
}

func resourceIBMTransitGatewayRouteReportDelete(d *schema.ResourceData, meta interface{}) error {

	client, err := transitgatewayClient(meta)
	if err != nil {
		return err
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}

	gatewayId := parts[0]
	ID := parts[1]
	deleteTransitGatewayRouteReportOptions := &transitgatewayapisv1.DeleteTransitGatewayRouteReportOptions{
		ID: &ID,
	}
	deleteTransitGatewayRouteReportOptions.SetTransitGatewayID(gatewayId)
	response, err := client.DeleteTransitGatewayRouteReport(deleteTransitGatewayRouteReportOptions)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return fmt.Errorf("Error deleting Transit Gateway Route Report(%s): %s\n%s", ID, err, response)
	}

	d.SetId("")
	return nil
}
