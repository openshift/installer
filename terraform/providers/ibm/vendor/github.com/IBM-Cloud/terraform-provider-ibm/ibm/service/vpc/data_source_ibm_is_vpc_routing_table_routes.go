// Copyright IBM Corp. 2017, 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func DataSourceIBMISVPCRoutingTableRoutes() *schema.Resource {
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
						"advertise": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether this route will be advertised to the ingress sources specified by the `advertise_routes_to` routing table property.",
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
		start = flex.GetNext(result.Next)
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
		// creator changes
		creator := []map[string]interface{}{}
		if instance.Creator != nil {
			mm, err := dataSourceIBMIsRouteCreatorToMap(instance.Creator)
			if err != nil {
				log.Printf("Error reading list of VPC Routing Table Routes' creator:%s", err)
				return err
			}
			creator = append(creator, mm)

		}
		route["creator"] = creator
		if instance.LifecycleState != nil {
			route[isRoutingTableRouteLifecycleState] = *instance.LifecycleState
		}
		if instance.Action != nil {
			route[isRoutingTableRouteAction] = *instance.Action
		}
		if instance.Advertise != nil {
			route["advertise"] = *instance.Advertise
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
		//orgin
		if instance.Origin != nil {
			route["origin"] = *instance.Origin
		}
		// priority
		if instance.Priority != nil {
			route["priority"] = *instance.Priority
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

func dataSourceIBMIsRouteCreatorToMap(model vpcv1.RouteCreatorIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.RouteCreatorVPNGatewayReference); ok {
		return DataSourceIBMIsRouteCreatorVPNGatewayReferenceToMap(model.(*vpcv1.RouteCreatorVPNGatewayReference))
	} else if _, ok := model.(*vpcv1.RouteCreatorVPNServerReference); ok {
		return DataSourceIBMIsRouteCreatorVPNServerReferenceToMap(model.(*vpcv1.RouteCreatorVPNServerReference))
	} else if _, ok := model.(*vpcv1.RouteCreator); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.RouteCreator)
		if model.CRN != nil {
			modelMap["crn"] = *model.CRN
		}
		if model.Deleted != nil {
			deletedMap, err := DataSourceIBMIsRouteVPNGatewayReferenceDeletedToMap(model.Deleted)
			if err != nil {
				return modelMap, err
			}
			modelMap["deleted"] = []map[string]interface{}{deletedMap}
		}
		if model.Href != nil {
			modelMap["href"] = *model.Href
		}
		if model.ID != nil {
			modelMap["id"] = *model.ID
		}
		if model.Name != nil {
			modelMap["name"] = *model.Name
		}
		if model.ResourceType != nil {
			modelMap["resource_type"] = *model.ResourceType
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("[Error] unrecognized vpcv1.RouteCreatorIntf subtype encountered")
	}
}

func DataSourceIBMIsRouteCreatorVPNGatewayReferenceToMap(model *vpcv1.RouteCreatorVPNGatewayReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CRN != nil {
		modelMap["crn"] = *model.CRN
	}
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsRouteVPNGatewayReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.ResourceType != nil {
		modelMap["resource_type"] = *model.ResourceType
	}
	return modelMap, nil
}

func DataSourceIBMIsRouteCreatorVPNServerReferenceToMap(model *vpcv1.RouteCreatorVPNServerReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CRN != nil {
		modelMap["crn"] = *model.CRN
	}
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsRouteVPNServerReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.ResourceType != nil {
		modelMap["resource_type"] = *model.ResourceType
	}
	return modelMap, nil
}

func DataSourceIBMIsRouteVPNGatewayReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}

func DataSourceIBMIsRouteVPNServerReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}
