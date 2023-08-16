// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package directlink

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/networking-go-sdk/directlinkv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMDLGatewayRouteReport() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMdlGatewayRouteReportCreate,
		Read:     resourceIBMDLRouteReportRead,
		Delete:   resourceIBMdlGatewayRouteReportDelete,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			dlGatewayId: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Direct Link gateway identifier",
			},
			dlRouteReportId: {
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

func resourceIBMdlGatewayRouteReportCreate(d *schema.ResourceData, meta interface{}) error {
	directLink, err := directlinkClient(meta)
	if err != nil {
		return err
	}

	gatewayId := d.Get(dlGatewayId).(string)
	createGatewayRouteReportOptionsModel := &directlinkv1.CreateGatewayRouteReportOptions{GatewayID: &gatewayId}
	routeReport, response, err := directLink.CreateGatewayRouteReport(createGatewayRouteReportOptionsModel)
	if err != nil {
		log.Println("[DEBUG] Create Route Report for DirectLink gateway", gatewayId, "err: ", err, " with response code:", response.StatusCode)
		return fmt.Errorf("[ERROR] Create Route Report for DirectLink gateway(%s) err: %s with response code: %d", gatewayId, err, response.StatusCode)
	}

	if routeReport == nil {
		return fmt.Errorf("error creating route report for gateway: %s with response code: %d", gatewayId, response.StatusCode)
	} else if routeReport.ID != nil {
		d.SetId(fmt.Sprintf("%s/%s", gatewayId, *routeReport.ID))
		d.Set(dlRouteReportId, *routeReport.ID)
	}

	isWaitForDirectLinkGatewayRouteReportCompleted(directLink, d.Id(), d.Timeout(schema.TimeoutCreate))

	return resourceIBMDLRouteReportRead(d, meta)
}

func isWaitForDirectLinkGatewayRouteReportCompleted(client *directlinkv1.DirectLinkV1, ID string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for direct link route report to be completed for  (%s) ", ID)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", dlRouteReportPending},
		Target:     []string{dlRouteReportComplete, ""},
		Refresh:    isDirectLinkGatewayRouteReportRefreshFunc(client, ID),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForState()
}
func isDirectLinkGatewayRouteReportRefreshFunc(client *directlinkv1.DirectLinkV1, ID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		parts, err := flex.IdParts(ID)
		if err != nil {
			return nil, "", fmt.Errorf("error getting ID for directlink route report: %s", err)
		}

		gatewayId := parts[0]
		routeReportId := parts[1]

		getOptions := &directlinkv1.GetGatewayRouteReportOptions{
			GatewayID: &gatewayId,
			ID:        &routeReportId,
		}
		routeReport, response, err := client.GetGatewayRouteReport(getOptions)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] error fetching directlink route report: %s\n%s", err, response)
		}
		if *routeReport.Status == "complete" {
			return routeReport, dlRouteReportComplete, nil
		}
		return routeReport, dlRouteReportPending, nil
	}
}

func resourceIBMDLRouteReportRead(d *schema.ResourceData, meta interface{}) error {
	directLink, err := directlinkClient(meta)
	if err != nil {
		return err
	}

	log.Println("[Info] Fetching DL Route Report GW ID:")

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}

	gatewayId := parts[0]
	routeReportId := parts[1]

	log.Println("[Info] Fetching DL Route Report GW ID:", gatewayId, " and Report ID: ", routeReportId)

	getGatewayRouteReportOptionsModel := &directlinkv1.GetGatewayRouteReportOptions{GatewayID: &gatewayId, ID: &routeReportId}
	report, response, err := directLink.GetGatewayRouteReport(getGatewayRouteReportOptionsModel)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			log.Println("[ERROR] Unable to fetch route report for directlink gateway with err:", err)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error fetching DL Route Reports for gateway(%s) err: %s with response code  %d", gatewayId, err, response.StatusCode)
	}

	if report == nil {
		return fmt.Errorf("error fetching route report for gateway: %s and route report: %s with response code: %d", gatewayId, routeReportId, response.StatusCode)
	}

	if report.Status != nil {
		log.Println("[Info] Fetching DL Route Reports status:", *report.Status)
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

func resourceIBMdlGatewayRouteReportDelete(d *schema.ResourceData, meta interface{}) error {

	directLink, err := directlinkClient(meta)
	if err != nil {
		return err
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}

	gatewayId := parts[0]
	routeReportId := parts[1]

	delOptions := directLink.NewDeleteGatewayRouteReportOptions(gatewayId, routeReportId)
	response, err := directLink.DeleteGatewayRouteReport(delOptions)

	if err != nil {
		if response.StatusCode == 404 {
			return nil
		}

		log.Printf("Error deleting Direct Link Route Report : %s", response)
		return err
	}

	d.SetId("")
	return nil
}
