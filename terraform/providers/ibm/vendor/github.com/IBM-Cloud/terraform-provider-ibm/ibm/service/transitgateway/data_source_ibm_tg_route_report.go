// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package transitgateway

import (
	"fmt"
	"github.com/IBM/networking-go-sdk/transitgatewayapisv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	tgRouteReport = "route_report"
)

func DataSourceIBMTransitGatewayRouteReport() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceIBMTransitGatewayRouteReportRead,
		Schema: map[string]*schema.Schema{
			tgGatewayId: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Transit Gateway identifier",
			},
			tgRouteReport: {
				Type:        schema.TypeString,
				Required:    true,
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

func dataSourceIBMTransitGatewayRouteReportRead(d *schema.ResourceData, meta interface{}) error {

	client, err := transitgatewayClient(meta)
	if err != nil {
		return err
	}

	gatewayId := d.Get(tgGatewayId).(string)
	routeReportId := d.Get(tgRouteReport).(string)

	getTransitGatewayRouteReportOptionsModel := &transitgatewayapisv1.GetTransitGatewayRouteReportOptions{}
	getTransitGatewayRouteReportOptionsModel.SetTransitGatewayID(gatewayId)
	getTransitGatewayRouteReportOptionsModel.SetID(routeReportId)
	routeReport, response, err := client.GetTransitGatewayRouteReport(getTransitGatewayRouteReportOptionsModel)
	if err != nil {
		return fmt.Errorf("Error while retrieving transit gateway route report %s\n%s", err, response)
	}

	d.Set(tgRouteReport, routeReport.ID)
	d.SetId(*routeReport.ID)
	d.Set(tgStatus, routeReport.Status)
	d.Set(tgCreatedAt, routeReport.CreatedAt.String())
	if routeReport.UpdatedAt != nil {
		d.Set(tgUpdatedAt, routeReport.UpdatedAt.String())
	}

	connections := make([]map[string]interface{}, 0)
	for _, connection := range routeReport.Connections {
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
	for _, overlap := range routeReport.OverlappingRoutes {
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
