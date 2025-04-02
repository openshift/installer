// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsVPNServerRoutes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsVPNServerRoutesRead,

		Schema: map[string]*schema.Schema{
			"vpn_server": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The VPN server identifier.",
			},
			"routes": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of VPN routes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The action to perform with a packet matching the VPN route:- `translate`: translate the source IP address to one of the private IP addresses of the VPN server.- `deliver`: deliver the packet into the VPC.- `drop`: drop the packetThe enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the VPN route on which the unexpected property value was encountered.",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the VPN route was created.",
						},
						"destination": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The destination for this VPN route in the VPN server. If an incoming packet does not match any destination, it will be dropped.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this VPN route.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this VPN route.",
						},
						"health_state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The health of this resource.- `ok`: Healthy- `degraded`: Suffering from compromised performance, capacity, or connectivity- `faulted`: Completely unreachable, inoperative, or otherwise entirely incapacitated- `inapplicable`: The health state does not apply because of the current lifecycle state. A resource with a lifecycle state of `failed` or `deleting` will have a health state of `inapplicable`. A `pending` resource may also have this state.",
						},
						"health_reasons": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A snake case string succinctly identifying the reason for this health state.",
									},

									"message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "An explanation of the reason for this health state.",
									},

									"more_info": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about the reason for this health state.",
									},
								},
							},
						},
						"lifecycle_state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lifecycle state of the VPN route.",
						},
						"lifecycle_reasons": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The reasons for the current lifecycle_state (if any).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A snake case string succinctly identifying the reason for this lifecycle state.",
									},

									"message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "An explanation of the reason for this lifecycle state.",
									},

									"more_info": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about the reason for this lifecycle state.",
									},
								},
							},
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this VPN route.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIsVPNServerRoutesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	start := ""
	allrecs := []vpcv1.VPNServerRoute{}

	for {
		listVPNServerRoutesOptions := &vpcv1.ListVPNServerRoutesOptions{}
		listVPNServerRoutesOptions.SetVPNServerID(d.Get("vpn_server").(string))

		if start != "" {
			listVPNServerRoutesOptions.Start = &start
		}
		vpnServerRouteCollection, response, err := sess.ListVPNServerRoutesWithContext(context, listVPNServerRoutesOptions)
		if err != nil {
			log.Printf("[DEBUG] ListVPNServerRoutesWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("[ERROR] ListVPNServerRoutesWithContext failed %s\n%s", err, response))
		}
		start = flex.GetNext(vpnServerRouteCollection.Next)
		allrecs = append(allrecs, vpnServerRouteCollection.Routes...)
		if start == "" {
			break
		}
	}

	d.SetId(dataSourceIBMIsVPNServerRoutesID(d))

	if allrecs != nil {
		err = d.Set("routes", dataSourceVPNServerRouteCollectionFlattenRoutes(allrecs))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting routes %s", err))
		}
	}

	return nil
}

// dataSourceIBMIsVPNServerRoutesID returns a reasonable ID for the list.
func dataSourceIBMIsVPNServerRoutesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceVPNServerRouteCollectionFlattenFirst(result vpcv1.PageLink) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceVPNServerRouteCollectionFirstToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceVPNServerRouteCollectionFirstToMap(firstItem vpcv1.PageLink) (firstMap map[string]interface{}) {
	firstMap = map[string]interface{}{}

	if firstItem.Href != nil {
		firstMap["href"] = firstItem.Href
	}

	return firstMap
}

func dataSourceVPNServerRouteCollectionFlattenNext(result vpcv1.PageLink) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceVPNServerRouteCollectionNextToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceVPNServerRouteCollectionNextToMap(nextItem vpcv1.PageLink) (nextMap map[string]interface{}) {
	nextMap = map[string]interface{}{}

	if nextItem.Href != nil {
		nextMap["href"] = nextItem.Href
	}

	return nextMap
}

func dataSourceVPNServerRouteCollectionFlattenRoutes(result []vpcv1.VPNServerRoute) (routes []map[string]interface{}) {
	for _, routesItem := range result {
		routes = append(routes, dataSourceVPNServerRouteCollectionRoutesToMap(routesItem))
	}

	return routes
}

func dataSourceVPNServerRouteCollectionRoutesToMap(routesItem vpcv1.VPNServerRoute) (routesMap map[string]interface{}) {
	routesMap = map[string]interface{}{}

	if routesItem.Action != nil {
		routesMap["action"] = routesItem.Action
	}
	if routesItem.CreatedAt != nil {
		routesMap["created_at"] = routesItem.CreatedAt.String()
	}
	if routesItem.Destination != nil {
		routesMap["destination"] = routesItem.Destination
	}
	if routesItem.Href != nil {
		routesMap["href"] = routesItem.Href
	}
	if routesItem.ID != nil {
		routesMap["id"] = routesItem.ID
	}
	if routesItem.HealthState != nil {
		routesMap["health_state"] = routesItem.HealthState
	}
	if routesItem.HealthReasons != nil {
		routesMap["health_reasons"] = resourceVPNServerRouteFlattenHealthReasons(routesItem.HealthReasons)
	}
	if routesItem.LifecycleState != nil {
		routesMap["lifecycle_state"] = routesItem.LifecycleState
	}
	if routesItem.LifecycleReasons != nil {
		routesMap["lifecycle_reasons"] = resourceVPNServerRouteFlattenLifecycleReasons(routesItem.LifecycleReasons)

	}
	if routesItem.Name != nil {
		routesMap["name"] = routesItem.Name
	}
	if routesItem.ResourceType != nil {
		routesMap["resource_type"] = routesItem.ResourceType
	}

	return routesMap
}
