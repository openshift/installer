// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package directlink

import (
	"log"
	"time"

	"github.com/IBM/networking-go-sdk/directlinkv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMDLRouteReports() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMDLRouteReportsRead,
		Schema: map[string]*schema.Schema{
			dlGatewayId: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Direct Link gateway identifier",
			},
			dlRouteReports: {
				Type:        schema.TypeList,
				Description: "List of route reports for a gateway",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						dlId: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the route report",
						},

						dlAdvertisedRoutes: {
							Type:        schema.TypeList,
							Description: "List of connection prefixes advertised to the on-prem network",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									dlAsPath: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The BGP AS path of the route",
									},
									dlPrefix: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Prefix for advertised routes",
									},
								},
							},
						},
						dlGatewayRoutes: {
							Type:        schema.TypeList,
							Description: "List of gateway routes",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									dlPrefix: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Prefix for gateway routes",
									},
								},
							},
						},
						dlOnPremRoutes: {
							Type:        schema.TypeList,
							Description: "List of onprem routes",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									dlAsPath: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The BGP AS path of the route",
									},
									dlPrefix: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Prefix for onprem routes",
									},
									dlRouteReportNextHop: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Next Hop address",
									},
								},
							},
						},
						dlOverlappingRoutes: {
							Type:        schema.TypeList,
							Description: "List of overlapping routes",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									dlRoutes: {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "overlapping routes",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												dlPrefix: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Prefix for overlapping routes",
												},
												dlType: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Type of route",
												},
												dlVirtualConnectionId: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Virtual connection ID",
												},
											},
										},
									},
								},
							},
						},
						dlVirtualConnectionRoutes: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Virtual Connection Routes",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									dlVirtualConnectionId: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Virtual connection ID",
									},
									dlVirtualConnectionType: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Virtual connection type",
									},
									dlVirtualConnectionName: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Virtual connection name",
									},
									dlRoutes: {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Virtual connection routes",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												dlPrefix: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Prefix for virtual connection routes",
												},
												dlActive: {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Indicates whether the route is the preferred path of the prefix",
												},
												dlLocalPreference: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The local preference of the route. This attribute can manipulate the chosen path on routes",
												},
											},
										},
									},
								},
							},
						},
						dlRouteReportStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Route report status",
						},
						dlCreatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time report was created",
						},
						dlUpdatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time resource was created",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMDLRouteReportsRead(d *schema.ResourceData, meta interface{}) error {
	directLink, err := directlinkClient(meta)
	if err != nil {
		return err
	}

	gatewayId := d.Get(dlGatewayId).(string)

	listGatewayRouteReportsOptionsModel := &directlinkv1.ListGatewayRouteReportsOptions{GatewayID: &gatewayId}
	routeReportsList, response, err := directLink.ListGatewayRouteReports(listGatewayRouteReportsOptionsModel)
	if err != nil {
		log.Println("[WARN] Error listing DL Route Reports", response, err)
		return err
	}

	routeReports := make([]map[string]interface{}, 0)
	for _, instance := range routeReportsList.RouteReports {
		routeReport := map[string]interface{}{}
		if instance.ID != nil {
			routeReport[dlId] = *instance.ID
		}
		if instance.Status != nil {
			routeReport[dlRouteReportStatus] = *instance.Status
		}

		// Build Advertised Routes
		advRoutes := make([]map[string]interface{}, 0)
		if instance.AdvertisedRoutes != nil {
			for _, adRoute := range instance.AdvertisedRoutes {
				route := map[string]interface{}{}
				route[dlAsPath] = adRoute.AsPath
				route[dlPrefix] = adRoute.Prefix
				advRoutes = append(advRoutes, route)
			}
		}

		log.Println("[Info] Length DL Route Reports advertised routes:", len(advRoutes))
		routeReport[dlAdvertisedRoutes] = advRoutes

		// Build Gateway Routes
		gatewayRoutes := make([]map[string]interface{}, 0)
		if instance.GatewayRoutes != nil {
			for _, gatewayRoute := range instance.GatewayRoutes {
				route := map[string]interface{}{}
				route[dlPrefix] = gatewayRoute.Prefix
				gatewayRoutes = append(gatewayRoutes, route)
			}

		}

		log.Println("[INFO] length DL Gateway Routes", len(gatewayRoutes))
		routeReport[dlGatewayRoutes] = gatewayRoutes

		// Build onPrem Routes
		onPremRoutes := make([]map[string]interface{}, 0)
		if instance.OnPremRoutes != nil {
			for _, onPremRoute := range instance.OnPremRoutes {
				route := map[string]interface{}{}
				route[dlPrefix] = onPremRoute.Prefix
				route[dlAsPath] = onPremRoute.AsPath
				route[dlRouteReportNextHop] = onPremRoute.NextHop
				onPremRoutes = append(onPremRoutes, route)
			}
		}

		log.Println("[INFO] length DL Onprem routes", len(onPremRoutes))
		routeReport[dlOnPremRoutes] = onPremRoutes

		// Build Overlapping Routes
		overlappingRoutesCollection := make([]map[string]interface{}, 0)
		if instance.OverlappingRoutes != nil && len(instance.OverlappingRoutes) > 0 {

			for _, o := range instance.OverlappingRoutes {
				overlappingRouteItem := map[string]interface{}{}
				routes := make([]map[string]interface{}, 0)
				for _, r := range o.Routes {
					overlappingRoute := map[string]interface{}{}
					route := r.(*directlinkv1.RouteReportOverlappingRoute)
					overlappingRoute[dlPrefix] = route.Prefix
					overlappingRoute[dlType] = route.Type
					overlappingRoute[dlVirtualConnectionId] = route.VirtualConnectionID

					routes = append(routes, overlappingRoute)
				}
				overlappingRouteItem[dlRoutes] = routes

				overlappingRoutesCollection = append(overlappingRoutesCollection, overlappingRouteItem)
			}
		}

		log.Println("[INFO] length DL overlapping routes", len(overlappingRoutesCollection))
		routeReport[dlOverlappingRoutes] = overlappingRoutesCollection

		// Build connection routes
		virtualConnectionRoutes := make([]map[string]interface{}, 0)
		if instance.VirtualConnectionRoutes != nil {
			for _, c := range instance.VirtualConnectionRoutes {
				conn := map[string]interface{}{}
				conn[dlVirtualConnectionId] = c.VirtualConnectionID
				conn[dlVirtualConnectionName] = c.VirtualConnectionName
				conn[dlVirtualConnectionType] = c.VirtualConnectionType

				connectionRoutes := make([]map[string]interface{}, 0)
				for _, r := range c.Routes {
					routes := map[string]interface{}{}
					routes[dlPrefix] = r.Prefix
					routes[dlLocalPreference] = r.LocalPreference
					routes[dlActive] = r.Active
					connectionRoutes = append(connectionRoutes, routes)
				}

				conn[dlRoutes] = connectionRoutes
				virtualConnectionRoutes = append(virtualConnectionRoutes, conn)
			}
		}

		log.Println("[INFO] length DL connection routes", len(virtualConnectionRoutes))
		routeReport[dlVirtualConnectionRoutes] = virtualConnectionRoutes

		// Add the created and updated dates
		if instance.CreatedAt != nil {
			routeReport[dlCreatedAt] = instance.CreatedAt.String()
		}
		if instance.UpdatedAt != nil {
			routeReport[dlUpdatedAt] = instance.UpdatedAt.String()
		}

		routeReports = append(routeReports, routeReport)

	}
	d.Set(dlRouteReports, routeReports)
	d.SetId(dataSourceIBMDirectLinkGatewayRouteReportsID(d))
	return nil
}

// dataSourceIBMDirectLinkGatewayRouteReportsID returns a reasonable ID for a directlink gateways list.
func dataSourceIBMDirectLinkGatewayRouteReportsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
