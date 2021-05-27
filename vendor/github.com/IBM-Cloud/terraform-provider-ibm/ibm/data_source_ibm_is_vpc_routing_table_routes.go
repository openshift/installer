// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"log"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isRoutingTableRouteID             = "route_id"
	isRoutingTableRouteHref           = "href"
	isRoutingTableRouteName           = "name"
	isRoutingTableRouteCreatedAt      = "created_at"
	isRoutingTableRouteLifecycleState = "lifecycle_state"
	isRoutingTableRouteAction         = "action"
	isRoutingTableRouteDestination    = "destination"
	isRoutingTableRouteNexthop        = "nexthop"
	isRoutingTableRouteZoneName       = "zone"
	isRoutingTableRouteVpcID          = "vpc"
	isRouteTableID                    = "routing_table"
	isRoutingTableRoutes              = "routes"
)

func dataSourceIBMISVPCRoutingTableRoutes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISVPCRoutingTableRoutesList,
		Schema: map[string]*schema.Schema{
			isRoutingTableRouteVpcID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPC identifier",
			},
			isRouteTableID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Routing table identifier",
			},
			isRoutingTableRoutes: {
				Type:        schema.TypeList,
				Description: "Collection of Routing Table Routes",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isRoutingTableRouteID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Routing Table Route ID",
						},
						isRoutingTableRouteHref: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Routing Table Route Href",
						},
						isRoutingTableRouteName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Routing Table Route Name",
						},
						isRoutingTableRouteCreatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Routing Table Route Created At",
						},
						isRoutingTableRouteLifecycleState: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Routing Table Route Lifecycle State",
						},
						isRoutingTableRouteAction: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Routing Table Route Action",
						},
						isRoutingTableRouteDestination: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Routing Table Route Destination",
						},
						isRoutingTableRouteNexthop: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Routing Table Route Nexthop Address or VPN Gateway Connection ID",
						},
						isRoutingTableRouteZoneName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Routing Table Route Zone Name",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMISVPCRoutingTableRoutesList(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	vpcID := d.Get(isRoutingTableRouteVpcID).(string)
	routingTableID := d.Get(isRouteTableID).(string)
	start := ""
	allrecs := []vpcv1.Route{}
	for {
		listVpcRoutingTablesRoutesOptions := sess.NewListVPCRoutingTableRoutesOptions(vpcID, routingTableID)
		if start != "" {
			listVpcRoutingTablesRoutesOptions.Start = &start
		}
		result, detail, err := sess.ListVPCRoutingTableRoutes(listVpcRoutingTablesRoutesOptions)
		if err != nil {
			log.Printf("Error reading list of VPC Routing Table Routes:%s\n%s", err, detail)
			return err
		}
		start = GetNext(result.Next)
		allrecs = append(allrecs, result.Routes...)
		if start == "" {
			break
		}
	}

	vpcRoutingTableRoutes := make([]map[string]interface{}, 0)

	for _, instance := range allrecs {
		route := map[string]interface{}{}
		if instance.ID != nil {
			route[isRoutingTableRouteID] = *instance.ID
		}
		if instance.Href != nil {
			route[isRoutingTableRouteHref] = *instance.Href
		}
		if instance.Name != nil {
			route[isRoutingTableRouteName] = *instance.Name
		}
		if instance.CreatedAt != nil {
			route[isRoutingTableRouteCreatedAt] = (*instance.CreatedAt).String()
		}
		if instance.LifecycleState != nil {
			route[isRoutingTableRouteLifecycleState] = *instance.LifecycleState
		}
		if instance.Destination != nil {
			route[isRoutingTableRouteDestination] = *instance.Destination
		}
		if instance.Zone != nil && instance.Zone.Name != nil {
			route[isRoutingTableRouteZoneName] = *instance.Zone.Name
		}
		if instance.NextHop != nil {
			nexthop := *instance.NextHop.(*vpcv1.RouteNextHop)
			if nexthop.Address != nil {
				route[isRoutingTableRouteNexthop] = *nexthop.Address
			} else {
				route[isRoutingTableRouteNexthop] = *nexthop.ID
			}
		}

		vpcRoutingTableRoutes = append(vpcRoutingTableRoutes, route)
	}
	d.SetId(dataSourceIBMISVPCRoutingTableRoutesID(d))
	d.Set(isRoutingTableRouteVpcID, vpcID)
	d.Set(isRouteTableID, routingTableID)
	d.Set(isRoutingTableRoutes, vpcRoutingTableRoutes)
	return nil
}

// dataSourceIBMISVPCRoutingTablesID returns a reasonable ID for dns zones list.
func dataSourceIBMISVPCRoutingTableRoutesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
