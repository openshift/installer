// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	//"encoding/json"
	"log"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isRoutingTableID                    = "routing_table"
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
	isRoutingTableTransitGatewayIngress = "route_transit_gateway_ingress"
	isRoutingTableVPCZoneIngress        = "route_vpc_zone_ingress"
	isRoutingTableDefault               = "is_default"
)

func dataSourceIBMISVPCRoutingTables() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISVPCRoutingTablesList,
		Schema: map[string]*schema.Schema{
			isVpcID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPC identifier",
			},
			isRoutingTables: {
				Type:        schema.TypeList,
				Description: "Collection of Routing tables",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isRoutingTableID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Routing Table ID",
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

	start := ""
	allrecs := []vpcv1.RoutingTable{}
	for {
		listOptions := sess.NewListVPCRoutingTablesOptions(vpcID)
		if start != "" {
			listOptions.Start = &start
		}
		result, detail, err := sess.ListVPCRoutingTables(listOptions)
		if err != nil {
			log.Printf("Error reading list of VPC Routing Tables:%s\n%s", err, detail)
			return err
		}
		start = GetNext(result.Next)
		allrecs = append(allrecs, result.RoutingTables...)
		if start == "" {
			break
		}
	}

	vpcRoutingTables := make([]map[string]interface{}, 0)
	for _, routingTable := range allrecs {

		rtable := map[string]interface{}{}
		if routingTable.ID != nil {
			rtable[isRoutingTableID] = *routingTable.ID
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
		if routingTable.RouteTransitGatewayIngress != nil {
			rtable[isRoutingTableTransitGatewayIngress] = *routingTable.RouteTransitGatewayIngress
		}
		if routingTable.RouteVPCZoneIngress != nil {
			rtable[isRoutingTableVPCZoneIngress] = *routingTable.RouteVPCZoneIngress
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
