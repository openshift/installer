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
	rtRoutes = "routes"
	rtCrn    = "crn"
)

func DataSourceIBMIBMIsVPCRoutingTable() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIBMIsVPCRoutingTableRead,

		Schema: map[string]*schema.Schema{
			isVpcID: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The VPC identifier.",
			},
			rName: &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				AtLeastOneOf:  []string{rName, isRoutingTableID},
				ConflictsWith: []string{isRoutingTableID},
				Description:   "The user-defined name for this routing table.",
			},
			isRoutingTableAcceptRoutesFrom: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The filters specifying the resources that may create routes in this routing table.At present, only the `resource_type` filter is permitted, and only the `vpn_gateway` value is supported, but filter support is expected to expand in the future.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			isRoutingTableID: &schema.Schema{
				Type:          schema.TypeString,
				AtLeastOneOf:  []string{rName, isRoutingTableID},
				ConflictsWith: []string{rName},
				Optional:      true,
				Description:   "The routing table identifier.",
			},

			"advertise_routes_to": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The ingress sources to advertise routes to. Routes in the table with `advertise` enabled will be advertised to these sources.The enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected property value was encountered.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			rtCreateAt: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that this routing table was created.",
			},
			rtHref: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this routing table.",
			},
			rtIsDefault: &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this is the default routing table for this VPC.",
			},
			rtLifecycleState: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the routing table.",
			},
			rtResourceType: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			rtRouteDirectLinkIngress: &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this routing table is used to route traffic that originates from[Direct Link](https://cloud.ibm.com/docs/dl/) to this VPC.Incoming traffic will be routed according to the routing table with one exception: routes with an `action` of `deliver` are treated as `drop` unless the `next_hop` is an IP address within the VPC's address prefix ranges. Therefore, if an incoming packet matches a route with a `next_hop` of an internet-bound IP address or a VPN gateway connection, the packet will be dropped.",
			},
			rtRouteInternetIngress: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this routing table is used to route traffic that originates from the internet.Incoming traffic will be routed according to the routing table with two exceptions:- Traffic destined for IP addresses associated with public gateways will not be  subject to routes in this routing table.- Routes with an action of deliver are treated as drop unless the `next_hop` is an  IP address bound to a network interface on a subnet in the route's `zone`.  Therefore, if an incoming packet matches a route with a `next_hop` of an  internet-bound IP address or a VPN gateway connection, the packet will be dropped.",
			},
			rtRouteTransitGatewayIngress: &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this routing table is used to route traffic that originates from from [Transit Gateway](https://cloud.ibm.com/cloud/transit-gateway/) to this VPC.Incoming traffic will be routed according to the routing table with one exception: routes with an `action` of `deliver` are treated as `drop` unless the `next_hop` is an IP address within the VPC's address prefix ranges. Therefore, if an incoming packet matches a route with a `next_hop` of an internet-bound IP address or a VPN gateway connection, the packet will be dropped.",
			},
			rtRouteVPCZoneIngress: &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this routing table is used to route traffic that originates from subnets in other zones in this VPC.Incoming traffic will be routed according to the routing table with one exception: routes with an `action` of `deliver` are treated as `drop` unless the `next_hop` is an IP address within the VPC's address prefix ranges. Therefore, if an incoming packet matches a route with a `next_hop` of an internet-bound IP address or a VPN gateway connection, the packet will be dropped.",
			},
			rtRoutes: &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The routes for this routing table.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
							Description: "The URL for this route.",
						},
						rId: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this route.",
						},
						rName: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this route.",
						},
					},
				},
			},
			rtSubnets: &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The subnets to which this routing table is attached.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						rtCrn: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this subnet.",
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
							Description: "The URL for this subnet.",
						},
						rId: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this subnet.",
						},
						rName: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this subnet.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIBMIsVPCRoutingTableRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}
	vpcId := d.Get(isVpcID).(string)
	rtId := d.Get(isRoutingTableID).(string)
	routingTableName := d.Get(rtName).(string)
	var routingTable *vpcv1.RoutingTable
	if rtId != "" {
		getVPCRoutingTableOptions := &vpcv1.GetVPCRoutingTableOptions{
			VPCID: &vpcId,
			ID:    &rtId,
		}

		rt, response, err := vpcClient.GetVPCRoutingTableWithContext(context, getVPCRoutingTableOptions)
		if err != nil {
			log.Printf("[DEBUG] GetVPCRoutingTableWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("[ERROR] GetVPCRoutingTableWithContext failed %s\n%s", err, response))
		}
		routingTable = rt
	} else {
		start := ""
		allrecs := []vpcv1.RoutingTable{}
		for {
			listOptions := &vpcv1.ListVPCRoutingTablesOptions{
				VPCID: &vpcId,
			}
			if start != "" {
				listOptions.Start = &start
			}
			result, detail, err := vpcClient.ListVPCRoutingTables(listOptions)
			if err != nil {
				log.Printf("[ERROR] Error reading list of VPC Routing Tables:%s\n%s", err, detail)
				return diag.FromErr(fmt.Errorf("[ERROR] ListVPCRoutingTables failed %s\n%s", err, detail))
			}
			start = flex.GetNext(result.Next)
			allrecs = append(allrecs, result.RoutingTables...)
			if start == "" {
				break
			}
		}
		for _, r := range allrecs {
			if *r.Name == routingTableName {
				routingTable = &r
			}
		}
		if routingTable == nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Provided routing table %s cannot be found in the vpc %s", routingTableName, vpcId))
		}
	}

	d.SetId(*routingTable.ID)

	if err = d.Set(rtCreateAt, flex.DateTimeToString(routingTable.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
	}
	acceptRoutesFromInfo := make([]map[string]interface{}, 0)
	if routingTable.AcceptRoutesFrom != nil {
		for _, AcceptRoutesFrom := range routingTable.AcceptRoutesFrom {
			l := map[string]interface{}{}
			if AcceptRoutesFrom.ResourceType != nil {
				l["resource_type"] = *AcceptRoutesFrom.ResourceType
				acceptRoutesFromInfo = append(acceptRoutesFromInfo, l)
			}
		}
	}
	if err = d.Set(isRoutingTableAcceptRoutesFrom, acceptRoutesFromInfo); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting accept_routes_from %s", err))
	}

	if err = d.Set(isRoutingTableID, routingTable.ID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting routing_table: %s", err))
	}

	if err = d.Set(rtHref, routingTable.Href); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting href: %s", err))
	}

	if err = d.Set(rtIsDefault, routingTable.IsDefault); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting is_default: %s", err))
	}

	if err = d.Set(rtLifecycleState, routingTable.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting lifecycle_state: %s", err))
	}

	if err = d.Set(rName, routingTable.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}

	if err = d.Set(rtResourceType, routingTable.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting resource_type: %s", err))
	}

	if err = d.Set(rtRouteDirectLinkIngress, routingTable.RouteDirectLinkIngress); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting route_direct_link_ingress: %s", err))
	}

	if err = d.Set(rtRouteInternetIngress, routingTable.RouteInternetIngress); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting route_internet_ingress: %s", err))
	}
	if err = d.Set(rtRouteTransitGatewayIngress, routingTable.RouteTransitGatewayIngress); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting route_transit_gateway_ingress: %s", err))
	}

	if err = d.Set(rtRouteVPCZoneIngress, routingTable.RouteVPCZoneIngress); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting route_vpc_zone_ingress: %s", err))
	}

	if err = d.Set("advertise_routes_to", routingTable.AdvertiseRoutesTo); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting value of advertise_routes_to: %s", err))
	}
	routes := []map[string]interface{}{}
	if routingTable.Routes != nil {
		for _, modelItem := range routingTable.Routes {
			modelMap, err := dataSourceIBMIBMIsVPCRoutingTableRouteReferenceToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			routes = append(routes, modelMap)
		}
	}
	if err = d.Set(rtRoutes, routes); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting routes %s", err))
	}

	subnets := []map[string]interface{}{}
	if routingTable.Subnets != nil {
		for _, modelItem := range routingTable.Subnets {
			modelMap, err := dataSourceIBMIBMIsVPCRoutingTableSubnetReferenceToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			subnets = append(subnets, modelMap)
		}
	}
	if err = d.Set(rtSubnets, subnets); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting subnets %s", err))
	}

	return nil
}

func dataSourceIBMIBMIsVPCRoutingTableRouteReferenceToMap(model *vpcv1.RouteReference) (map[string]interface{}, error) {
	modelMap := map[string]interface{}{}
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIBMIsVPCRoutingTableRouteReferenceDeletedToMap(model.Deleted)
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
	return modelMap, nil
}

func dataSourceIBMIBMIsVPCRoutingTableRouteReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := map[string]interface{}{}
	if model.MoreInfo != nil {
		modelMap[rMoreInfo] = *model.MoreInfo
	}
	return modelMap, nil
}

func dataSourceIBMIBMIsVPCRoutingTableSubnetReferenceToMap(model *vpcv1.SubnetReference) (map[string]interface{}, error) {
	modelMap := map[string]interface{}{}
	if model.CRN != nil {
		modelMap[rtCrn] = *model.CRN
	}
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIBMIsVPCRoutingTableSubnetReferenceDeletedToMap(model.Deleted)
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
	return modelMap, nil
}

func dataSourceIBMIBMIsVPCRoutingTableSubnetReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := map[string]interface{}{}
	if model.MoreInfo != nil {
		modelMap[rMoreInfo] = *model.MoreInfo
	}
	return modelMap, nil
}
