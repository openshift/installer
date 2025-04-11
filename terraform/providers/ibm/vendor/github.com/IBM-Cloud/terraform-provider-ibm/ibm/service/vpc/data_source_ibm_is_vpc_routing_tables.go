// Copyright IBM Corp. 2017, 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	//"encoding/json"

	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isRoutingTableAcceptRoutesFrom      = "accept_routes_from"
	isRoutingTableID                    = "routing_table"
	isRoutingTableCrn                   = "routing_table_crn"
	isRoutingTableHref                  = "href"
	isRoutingTableName                  = "name"
	isRoutingTableResourceType          = "resource_type"
	isRoutingTableCreatedAt             = "created_at"
	isRoutingTableLifecycleState        = "lifecycle_state"
	isRoutingTableRoutesList            = "routes"
	isRoutingTableSubnetsList           = "subnets"
	isRoutingTables                     = "routing_tables"
	isVpcID                             = "vpc"
	isRoutingTableDirectLinkIngress     = "route_direct_link_ingress"
	isRoutingTableInternetIngress       = "route_internet_ingress"
	isRoutingTableTransitGatewayIngress = "route_transit_gateway_ingress"
	isRoutingTableVPCZoneIngress        = "route_vpc_zone_ingress"
	isRoutingTableDefault               = "is_default"
)

func DataSourceIBMISVPCRoutingTables() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISVPCRoutingTablesList,
		Schema: map[string]*schema.Schema{
			isVpcID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPC identifier",
			},
			isRoutingTableDefault: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Filters the collection to routing tables with the specified is_default value",
			},
			isRoutingTables: {
				Type:        schema.TypeList,
				Description: "Collection of Routing tables",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						isRoutingTableID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Routing Table ID",
						},
						isRoutingTableCrn: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The crn of routing table",
						},
						"advertise_routes_to": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The ingress sources to advertise routes to. Routes in the table with `advertise` enabled will be advertised to these sources.The enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected property value was encountered.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						isRoutingTableHref: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Routing table Href",
						},
						isRoutingTableName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Routing table Name",
						},
						isRoutingTableResourceType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Routing table Resource Type",
						},
						isRoutingTableCreatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Routing table Created At",
						},
						isRoutingTableLifecycleState: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Routing table Lifecycle State",
						},
						isRoutingTableDirectLinkIngress: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If set to true, this routing table will be used to route traffic that originates from Direct Link to this VPC.",
						},
						isRoutingTableInternetIngress: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether this routing table is used to route traffic that originates from the internet.Incoming traffic will be routed according to the routing table with two exceptions:- Traffic destined for IP addresses associated with public gateways will not be  subject to routes in this routing table.- Routes with an action of deliver are treated as drop unless the `next_hop` is an  IP address bound to a network interface on a subnet in the route's `zone`.  Therefore, if an incoming packet matches a route with a `next_hop` of an  internet-bound IP address or a VPN gateway connection, the packet will be dropped.",
						},
						isRoutingTableTransitGatewayIngress: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If set to true, this routing table will be used to route traffic that originates from Transit Gateway to this VPC.",
						},
						isRoutingTableVPCZoneIngress: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If set to true, this routing table will be used to route traffic that originates from subnets in other zones in this VPC.",
						},
						isRoutingTableDefault: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether this is the default routing table for this VPC",
						},
						isRoutingTableRoutesList: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Route name",
									},

									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Route ID",
									},
								},
							},
						},
						isRoutingTableSubnetsList: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subnet name",
									},

									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subnet ID",
									},
								},
							},
						},
						rtResourceGroup: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The resource group for this volume.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									rtResourceGroupHref: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this resource group.",
									},
									rtResourceGroupId: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this resource group.",
									},
									rtResourceGroupName: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this resource group.",
									},
								},
							},
						},

						rtTags: {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      flex.ResourceIBMVPCHash,
						},

						rtAccessTags: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         flex.ResourceIBMVPCHash,
							Description: "List of access tags",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMISVPCRoutingTablesList(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	vpcID := d.Get(isVpcID).(string)
	listOptions := sess.NewListVPCRoutingTablesOptions(vpcID)

	if isDefaultIntf, ok := d.GetOk(isRoutingTableDefault); ok {
		isDefault := isDefaultIntf.(bool)
		listOptions.IsDefault = &isDefault
	}
	start := ""
	allrecs := []vpcv1.RoutingTable{}
	for {
		if start != "" {
			listOptions.Start = &start
		}
		result, detail, err := sess.ListVPCRoutingTables(listOptions)
		if err != nil {
			log.Printf("Error reading list of VPC Routing Tables:%s\n%s", err, detail)
			return err
		}
		start = flex.GetNext(result.Next)
		allrecs = append(allrecs, result.RoutingTables...)
		if start == "" {
			break
		}
	}

	vpcRoutingTables := make([]map[string]interface{}, 0)
	for _, routingTable := range allrecs {

		rtable := map[string]interface{}{}
		acceptRoutesFromInfo := make([]map[string]interface{}, 0)
		for _, AcceptRoutesFrom := range routingTable.AcceptRoutesFrom {
			if AcceptRoutesFrom.ResourceType != nil {
				l := map[string]interface{}{}
				l["resource_type"] = *AcceptRoutesFrom.ResourceType
				acceptRoutesFromInfo = append(acceptRoutesFromInfo, l)
			}
		}
		rtable[isRoutingTableAcceptRoutesFrom] = acceptRoutesFromInfo
		if routingTable.ID != nil {
			rtable[isRoutingTableID] = *routingTable.ID
		}
		if routingTable.CRN != nil {
			rtable[isRoutingTableCrn] = *routingTable.CRN
		}
		if routingTable.Href != nil {
			rtable[isRoutingTableHref] = *routingTable.Href
		}
		if routingTable.Name != nil {
			rtable[isRoutingTableName] = *routingTable.Name
		}
		if routingTable.ResourceType != nil {
			rtable[isRoutingTableResourceType] = *routingTable.ResourceType
		}
		if routingTable.CreatedAt != nil {
			rtable[isRoutingTableCreatedAt] = (*routingTable.CreatedAt).String()
		}
		if routingTable.LifecycleState != nil {
			rtable[isRoutingTableLifecycleState] = *routingTable.LifecycleState
		}
		if routingTable.RouteDirectLinkIngress != nil {
			rtable[isRoutingTableDirectLinkIngress] = *routingTable.RouteDirectLinkIngress
		}
		if routingTable.RouteInternetIngress != nil {
			rtable[isRoutingTableInternetIngress] = *&routingTable.RouteInternetIngress
		}
		if routingTable.RouteTransitGatewayIngress != nil {
			rtable[isRoutingTableTransitGatewayIngress] = *routingTable.RouteTransitGatewayIngress
		}
		if routingTable.RouteVPCZoneIngress != nil {
			rtable[isRoutingTableVPCZoneIngress] = *routingTable.RouteVPCZoneIngress
		}
		if routingTable.AdvertiseRoutesTo != nil {
			rtable["advertise_routes_to"] = routingTable.AdvertiseRoutesTo
		}
		if routingTable.IsDefault != nil {
			rtable[isRoutingTableDefault] = *routingTable.IsDefault
		}
		subnetsInfo := make([]map[string]interface{}, 0)
		for _, subnet := range routingTable.Subnets {
			if subnet.Name != nil && subnet.ID != nil {
				l := map[string]interface{}{
					"name": *subnet.Name,
					"id":   *subnet.ID,
				}
				subnetsInfo = append(subnetsInfo, l)
			}

		}
		rtable[isRoutingTableSubnetsList] = subnetsInfo
		routesInfo := make([]map[string]interface{}, 0)
		for _, route := range routingTable.Routes {
			if route.Name != nil && route.ID != nil {
				k := map[string]interface{}{
					"name": *route.Name,
					"id":   *route.ID,
				}
				routesInfo = append(routesInfo, k)
			}
		}
		rtable[isRoutingTableRoutesList] = routesInfo

		resourceGroupList := []map[string]interface{}{}
		if routingTable.ResourceGroup != nil {
			resourceGroupMap := routingTableResourceGroupToMap(*routingTable.ResourceGroup)
			resourceGroupList = append(resourceGroupList, resourceGroupMap)
		}
		rtable[rtResourceGroup] = resourceGroupList

		tags, err := flex.GetGlobalTagsUsingCRN(meta, *routingTable.CRN, "", rtUserTagType)
		if err != nil {
			log.Printf(
				"An error occured during reading of routing table (%s) tags : %s", d.Id(), err)
		}
		rtable[rtTags] = tags

		accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *routingTable.CRN, "", rtAccessTagType)
		if err != nil {
			log.Printf(
				"An error occured during reading of routing table (%s) access tags: %s", d.Id(), err)
		}
		rtable[rtAccessTags] = accesstags

		vpcRoutingTables = append(vpcRoutingTables, rtable)
	}

	d.SetId(dataSourceIBMISVPCRoutingTablesID(d))
	d.Set(isVpcID, vpcID)
	d.Set(isRoutingTables, vpcRoutingTables)
	return nil
}

// dataSourceIBMISVPCRoutingTablesID returns a reasonable ID for dns zones list.
func dataSourceIBMISVPCRoutingTablesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
