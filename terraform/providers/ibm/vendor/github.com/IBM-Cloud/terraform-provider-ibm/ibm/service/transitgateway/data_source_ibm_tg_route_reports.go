// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package transitgateway

import (
	"fmt"
	"time"

	"github.com/IBM/networking-go-sdk/transitgatewayapisv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	tgRouteReports                            = "route_reports"
	tgRouteReportConnections                  = "connections"
	tgRouteReportOverlappingRoutes            = "overlapping_routes"
	tgRouteReportOverlappingRoutesDetail      = "routes"
	tgRouteReportConnectionBgps               = "bgps"
	tgRouteReportConnectionBgpAsPath          = "as_path"
	tgRouteReportConnectionBgpIsUsed          = "is_used"
	tgRouteReportConnectionBgpLocalPreference = "local_preference"
	tgRouteReportConnectionBgpPrefix          = "prefix"
	tgRouteReportConnectionRoutes             = "routes"
	tgRouteReportConnectionRoutePrefix        = "prefix"
	tgRouteReportConnectionType               = "type"
	tgRouteReportOverlappingPrefix            = "prefix"
)

func DataSourceIBMTransitGatewayRouteReports() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceIBMTransitGatewayRouteReportsRead,
		Schema: map[string]*schema.Schema{

			tgGatewayId: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Transit Gateway identifier",
			},

			tgRouteReports: {
				Type:        schema.TypeList,
				Description: "Collection of transit gateway route reports",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						tgID: {
							Type:     schema.TypeString,
							Computed: true,
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
				},
			},
		},
	}
}

func dataSourceIBMTransitGatewayRouteReportsRead(d *schema.ResourceData, meta interface{}) error {

	client, err := transitgatewayClient(meta)
	if err != nil {
		return err
	}

	gatewayId := d.Get(tgGatewayId).(string)

	listTransitGatewayRouteReportsOptionsModel := &transitgatewayapisv1.ListTransitGatewayRouteReportsOptions{}
	listTransitGatewayRouteReportsOptionsModel.SetTransitGatewayID(gatewayId)
	listTransitGatewayRouteReports, response, err := client.ListTransitGatewayRouteReports(listTransitGatewayRouteReportsOptionsModel)
	if err != nil {
		return fmt.Errorf("Error while listing transit gateway route reports %s\n%s", err, response)
	}

	reports := make([]map[string]interface{}, 0)
	for _, routeReport := range listTransitGatewayRouteReports.RouteReports {
		tgReport := map[string]interface{}{}
		tgReport[tgID] = routeReport.ID
		tgReport[tgCreatedAt] = routeReport.CreatedAt.String()
		tgReport[tgStatus] = routeReport.Status

		if routeReport.UpdatedAt != nil {
			tgReport[tgUpdatedAt] = routeReport.UpdatedAt.String()
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
		tgReport[tgRouteReportConnections] = connections

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
		tgReport[tgRouteReportOverlappingRoutes] = overlappingRoutes

		reports = append(reports, tgReport)
	}
	d.Set(tgRouteReports, reports)
	d.SetId(dataSourceIBMTransitGatewayRouteReportsID(d))
	return nil
}

// dataSourceIBMTransitGatewayRouteReportsID returns a reasonable ID for a transit gateways list.
func dataSourceIBMTransitGatewayRouteReportsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
