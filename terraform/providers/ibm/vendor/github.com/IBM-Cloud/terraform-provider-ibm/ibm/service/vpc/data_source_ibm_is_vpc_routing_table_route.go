// Copyright IBM Corp. 2021, 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

const (
	rDeleted  = "deleted"
	rAddress  = "address"
	rMoreInfo = "more_info"
	rId       = "id"
)

func DataSourceIBMIBMIsVPCRoutingTableRoute() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIBMIsVPCRoutingTableRouteRead,

		Schema: map[string]*schema.Schema{
			isVpcID: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The VPC identifier.",
			},
			isRoutingTableID: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The routing table identifier.",
			},
			isRoutingTableRouteID: &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				AtLeastOneOf:  []string{rName, isRoutingTableRouteID},
				ConflictsWith: []string{rName},
				Description:   "The VPC routing table route identifier.",
			},
			rName: &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				AtLeastOneOf:  []string{rName, isRoutingTableRouteID},
				ConflictsWith: []string{isRoutingTableRouteID},
				Description:   "The user-defined name for this route.",
			},
			rAction: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The action to perform with a packet matching the route:- `delegate`: delegate to the system's built-in routes- `delegate_vpc`: delegate to the system's built-in routes, ignoring Internet-bound  routes- `deliver`: deliver the packet to the specified `next_hop`- `drop`: drop the packet.",
			},
			"advertise": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this route will be advertised to the ingress sources specified by the `advertise_routes_to` routing table property.",
			},
			rtCreateAt: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the route was created.",
			},
			"creator": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "If present, the resource that created the route. Routes with this property present cannot bedirectly deleted. All routes with an `origin` of `learned` or `service` will have thisproperty set, and future `origin` values may also have this property set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPN gateway's CRN.",
						},
						"deleted": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPN gateway's canonical URL.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this VPN gateway.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this VPN gateway.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			rDestination: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The destination of the route.",
			},
			rtHref: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this route.",
			},
			rtLifecycleState: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the route.",
			},
			rNextHop: &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "If `action` is `deliver`, the next hop that packets will be delivered to.  Forother `action` values, its `address` will be `0.0.0.0`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						rAddress: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address.This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.",
						},
						rDeleted: &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									rMoreInfo: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						rtHref: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPN connection's canonical URL.",
						},
						rId: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this VPN gateway connection.",
						},
						rName: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this VPN connection.",
						},
						rtResourceType: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"origin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The origin of this route:- `service`: route was directly created by a service- `user`: route was directly created by a userThe enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the route on which the unexpected property value was encountered.",
			},
			"priority": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The route's priority. Smaller values have higher priority. If a routing table contains routes with the same destination, the route with the highest priority (smallest value) is selected.",
			},
			rZone: &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The zone the route applies to. (Traffic from subnets in this zone will besubject to this route.).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						rtHref: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this zone.",
						},
						rName: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique name for this zone.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIBMIsVPCRoutingTableRouteRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	vpcID := d.Get(isVpcID).(string)
	routingTableId := d.Get("routing_table").(string)
	routeId := d.Get("route_id").(string)
	routeName := d.Get("name").(string)
	var route *vpcv1.Route

	if routeId != "" {
		getVPCRoutingTableRouteOptions := &vpcv1.GetVPCRoutingTableRouteOptions{}
		getVPCRoutingTableRouteOptions.SetVPCID(vpcID)
		getVPCRoutingTableRouteOptions.SetRoutingTableID(routingTableId)
		getVPCRoutingTableRouteOptions.SetID(routeId)

		r, response, err := vpcClient.GetVPCRoutingTableRouteWithContext(context, getVPCRoutingTableRouteOptions)
		if err != nil {
			log.Printf("[DEBUG] GetVPCRoutingTableRouteWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("[ERROR] GetVPCRoutingTableRouteWithContext failed %s\n%s", err, response))
		}
		route = r
	} else {
		start := ""
		allrecs := []vpcv1.Route{}
		for {
			listVpcRoutingTablesRoutesOptions := &vpcv1.ListVPCRoutingTableRoutesOptions{
				VPCID:          &vpcID,
				RoutingTableID: &routingTableId,
			}

			if start != "" {
				listVpcRoutingTablesRoutesOptions.Start = &start
			}
			result, detail, err := vpcClient.ListVPCRoutingTableRoutes(listVpcRoutingTablesRoutesOptions)
			if err != nil {
				log.Printf("Error reading list of VPC Routing Table Routes:%s\n%s", err, detail)
				return diag.FromErr(fmt.Errorf("[ERROR] GetVPCRoutingTableRouteWithContext failed %s\n%s", err, detail))
			}
			start = flex.GetNext(result.Next)
			allrecs = append(allrecs, result.Routes...)
			if start == "" {
				break
			}
		}

		for _, r := range allrecs {
			if *r.Name == routeName {
				route = &r
				break
			}
		}
		if route == nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Route not found with name: %s", routeName))
		}
	}

	d.SetId(*route.ID)

	if err = d.Set(rAction, route.Action); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting action: %s", err))
	}

	if err = d.Set("advertise", route.Advertise); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting advertise: %s", err))
	}

	if err = d.Set(rtCreateAt, flex.DateTimeToString(route.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
	}

	// creator changes
	creator := []map[string]interface{}{}
	if route.Creator != nil {
		mm, err := dataSourceIBMIsRouteCreatorToMap(route.Creator)
		if err != nil {
			log.Printf("Error reading list of VPC Routing Table Routes' creator:%s", err)
			return diag.FromErr(fmt.Errorf("[ERROR] Error fetching creator: %s", err))
		}
		creator = append(creator, mm)

	}
	if err = d.Set("creator", creator); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting creator: %s", err))
	}

	if err = d.Set(isRoutingTableRouteID, route.ID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting route_id: %s", err))
	}

	if err = d.Set(rDestination, route.Destination); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting destination: %s", err))
	}

	if err = d.Set(rtHref, route.Href); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting href: %s", err))
	}

	if err = d.Set(rtLifecycleState, route.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting lifecycle_state: %s", err))
	}

	if err = d.Set(rName, route.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}

	nextHop := []map[string]interface{}{}
	if route.NextHop != nil {
		modelMap, err := dataSourceIBMIBMIsVPCRoutingTableRouteRouteNextHopToMap(route.NextHop)
		if err != nil {
			return diag.FromErr(err)
		}
		nextHop = append(nextHop, modelMap)
	}
	if err = d.Set(rNextHop, nextHop); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting next_hop %s", err))
	}

	//orgin
	if err = d.Set("origin", route.Origin); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting origin %s", err))
	}
	// priority
	if err = d.Set("priority", route.Priority); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting priority, :%s", err))
	}
	zone := []map[string]interface{}{}
	if route.Zone != nil {
		modelMap, err := dataSourceIBMIBMIsVPCRoutingTableRouteZoneReferenceToMap(route.Zone)
		if err != nil {
			return diag.FromErr(err)
		}
		zone = append(zone, modelMap)
	}
	if err = d.Set(rZone, zone); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting zone %s", err))
	}

	return nil
}

func dataSourceIBMIBMIsVPCRoutingTableRouteRouteNextHopToMap(model vpcv1.RouteNextHopIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.RouteNextHopIP); ok {
		return dataSourceIBMIBMIsVPCRoutingTableRouteRouteNextHopIPToMap(model.(*vpcv1.RouteNextHopIP))
	} else if _, ok := model.(*vpcv1.RouteNextHopVPNGatewayConnectionReference); ok {
		return dataSourceIBMIBMIsVPCRoutingTableRouteRouteNextHopVPNGatewayConnectionReferenceToMap(model.(*vpcv1.RouteNextHopVPNGatewayConnectionReference))
	} else if _, ok := model.(*vpcv1.RouteNextHop); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.RouteNextHop)
		if model.Address != nil {
			modelMap[rAddress] = *model.Address
		}
		if model.Deleted != nil {
			deletedMap, err := dataSourceIBMIBMIsVPCRoutingTableRouteVPNGatewayConnectionReferenceDeletedToMap(model.Deleted)
			if err != nil {
				return modelMap, err
			}
			modelMap[rDeleted] = []map[string]interface{}{deletedMap}
		}
		if model.Href != nil {
			modelMap[rtHref] = *model.Href
		}
		if model.ID != nil {
			modelMap[rId] = *model.ID
		}
		if model.Name != nil {
			modelMap[rName] = *model.Name
		}
		if model.ResourceType != nil {
			modelMap[rtResourceType] = *model.ResourceType
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("[ERROR] Unrecognized vpcv1.RouteNextHopIntf subtype encountered")
	}
}

func dataSourceIBMIBMIsVPCRoutingTableRouteVPNGatewayConnectionReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := map[string]interface{}{}
	if model.MoreInfo != nil {
		modelMap[rMoreInfo] = *model.MoreInfo
	}
	return modelMap, nil
}

func dataSourceIBMIBMIsVPCRoutingTableRouteRouteNextHopIPToMap(model *vpcv1.RouteNextHopIP) (map[string]interface{}, error) {
	modelMap := map[string]interface{}{}
	if model.Address != nil {
		modelMap[rAddress] = *model.Address
	}
	return modelMap, nil
}

func dataSourceIBMIBMIsVPCRoutingTableRouteRouteNextHopVPNGatewayConnectionReferenceToMap(model *vpcv1.RouteNextHopVPNGatewayConnectionReference) (map[string]interface{}, error) {
	modelMap := map[string]interface{}{}
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIBMIsVPCRoutingTableRouteVPNGatewayConnectionReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap[rDeleted] = []map[string]interface{}{deletedMap}
	}
	if model.Href != nil {
		modelMap[rtHref] = *model.Href
	}
	if model.ID != nil {
		modelMap[rId] = *model.ID
	}
	if model.Name != nil {
		modelMap[rName] = *model.Name
	}
	if model.ResourceType != nil {
		modelMap[rtResourceType] = *model.ResourceType
	}
	return modelMap, nil
}

func dataSourceIBMIBMIsVPCRoutingTableRouteZoneReferenceToMap(model *vpcv1.ZoneReference) (map[string]interface{}, error) {
	modelMap := map[string]interface{}{}
	if model.Href != nil {
		modelMap[rtHref] = *model.Href
	}
	if model.Name != nil {
		modelMap[rName] = *model.Name
	}
	return modelMap, nil
}
