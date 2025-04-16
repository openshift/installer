// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package directlink

import (
	"fmt"
	"log"

	"github.com/IBM/networking-go-sdk/directlinkv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	dlRouteReport = "route_report"
)

func DataSourceIBMDLRouteReport() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMDLRouteReportRead,
		Schema: map[string]*schema.Schema{
			dlGatewayId: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Direct Link gateway identifier",
			},
			dlRouteReport: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
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
							Description: "Prefix for gateway routes",
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
	}
}

func dataSourceIBMDLRouteReportRead(d *schema.ResourceData, meta interface{}) error {
	directLink, err := directlinkClient(meta)
	if err != nil {
		return err
	}

	gatewayId := d.Get(dlGatewayId).(string)
	routeReportId := d.Get(dlRouteReport).(string)

	log.Println("[Info]  fetching DL Route Reports GW ID:", gatewayId, " and report ID: ", routeReportId)

	getGatewayRouteReportOptionsModel := &directlinkv1.GetGatewayRouteReportOptions{GatewayID: &gatewayId, ID: &routeReportId}
	report, response, err := directLink.GetGatewayRouteReport(getGatewayRouteReportOptionsModel)
	if err != nil {
		log.Println("[DEBUG] Error fetching DL Route Reports for gateway:", gatewayId, "with response code:", response.StatusCode, " and err: ", err)
		return fmt.Errorf("[ERROR] Error fetching DL Route Reports: %s with response code  %d", err, response.StatusCode)
	}

	if report == nil {
		return fmt.Errorf("error fetching route report for gateway: %s and route report: %s with response code: %d", gatewayId, routeReportId, response.StatusCode)
	} else if report.ID != nil {
		d.SetId(*report.ID)
	}

	if report.Status != nil {
		log.Println("[Info]  fetching DL Route Reports status:", *report.Status)
		d.Set(dlRouteReportStatus, *report.Status)
	}

	// Build Advertised Routes
	advRoutes := make([]map[string]interface{}, 0)
	if report.AdvertisedRoutes != nil {
		for _, adRoute := range report.AdvertisedRoutes {
			route := map[string]interface{}{}
			route[dlAsPath] = adRoute.AsPath
			route[dlPrefix] = adRoute.Prefix
			advRoutes = append(advRoutes, route)
		}
	}

	log.Println("[Info] Length DL Route Reports advertised routes:", len(advRoutes))
	d.Set(dlAdvertisedRoutes, advRoutes)

	// Build Gateway Routes
	gatewayRoutes := make([]map[string]interface{}, 0)
	if report.GatewayRoutes != nil {
		for _, gatewayRoute := range report.GatewayRoutes {
			route := map[string]interface{}{}
			route[dlPrefix] = gatewayRoute.Prefix
			gatewayRoutes = append(gatewayRoutes, route)
		}

	}
	log.Println("[Info] Length DL Gateway Reports: ", len(gatewayRoutes))
	d.Set(dlGatewayRoutes, gatewayRoutes)

	// Build onPrem Routes
	onPremRoutes := make([]map[string]interface{}, 0)
	if report.OnPremRoutes != nil {
		for _, onPremRoute := range report.OnPremRoutes {
			route := map[string]interface{}{}
			route[dlPrefix] = onPremRoute.Prefix
			route[dlAsPath] = onPremRoute.AsPath
			route[dlRouteReportNextHop] = onPremRoute.NextHop
			onPremRoutes = append(onPremRoutes, route)
		}

	}

	log.Println("[Info] Length DL Route Reports onprem routes:", len(onPremRoutes))
	d.Set(dlOnPremRoutes, onPremRoutes)

	// Build Overlapping Routes
	overlappingRoutesCollection := make([]map[string]interface{}, 0)
	if report.OverlappingRoutes != nil && len(report.OverlappingRoutes) > 0 {

		for _, o := range report.OverlappingRoutes {
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

	log.Println("[INFO] Length DL overlapping routes", len(overlappingRoutesCollection))
	d.Set(dlOverlappingRoutes, overlappingRoutesCollection)

	// Build connection routes
	virtualConnectionRoutes := make([]map[string]interface{}, 0)
	if report.VirtualConnectionRoutes != nil {
		for _, c := range report.VirtualConnectionRoutes {
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

	log.Println("[Info] Length DL Route Reports connection routes:", len(virtualConnectionRoutes))
	d.Set(dlVirtualConnectionRoutes, virtualConnectionRoutes)

	// Add the created and updated dates
	if report.CreatedAt != nil {
		d.Set(dlCreatedAt, report.CreatedAt.String())
	}
	if report.UpdatedAt != nil {
		d.Set(dlUpdatedAt, report.UpdatedAt.String())
	}

	return nil
}
