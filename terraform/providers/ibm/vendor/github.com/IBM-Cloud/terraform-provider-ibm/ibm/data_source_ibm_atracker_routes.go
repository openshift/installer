// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/platform-services-go-sdk/atrackerv1"
)

func dataSourceIBMAtrackerRoutes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMAtrackerRoutesRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the route.",
			},
			"routes": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of route resources.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The uuid of the route resource.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the route.",
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The crn of the route resource.",
						},
						"version": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The version of the route.",
						},
						"receive_global_events": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether or not all global events should be forwarded to this region.",
						},
						"rules": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The routing rules that will be evaluated in their order of the array.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_ids": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The target ID List. Only 1 target id is supported.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"created": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The timestamp of the route creation time.",
						},
						"updated": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The timestamp of the route last updated time.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMAtrackerRoutesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	atrackerClient, err := meta.(ClientSession).AtrackerV1()
	if err != nil {
		return diag.FromErr(err)
	}

	listRoutesOptions := &atrackerv1.ListRoutesOptions{}

	routeList, response, err := atrackerClient.ListRoutesWithContext(context, listRoutesOptions)
	if err != nil {
		log.Printf("[DEBUG] ListRoutesWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ListRoutesWithContext failed %s\n%s", err, response))
	}

	// Use the provided filter argument and construct a new list with only the requested resource(s)
	var matchRoutes []atrackerv1.Route
	var name string
	var suppliedFilter bool

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
		suppliedFilter = true
		for _, data := range routeList.Routes {
			if *data.Name == name {
				matchRoutes = append(matchRoutes, data)
			}
		}
	} else {
		matchRoutes = routeList.Routes
	}
	routeList.Routes = matchRoutes

	if suppliedFilter {
		if len(routeList.Routes) == 0 {
			return diag.FromErr(fmt.Errorf("no Routes found with name %s", name))
		}
		d.SetId(name)
	} else {
		d.SetId(dataSourceIBMAtrackerRoutesID(d))
	}

	if routeList.Routes != nil {
		err = d.Set("routes", dataSourceRouteListFlattenRoutes(routeList.Routes))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting routes %s", err))
		}
	}

	return nil
}

// dataSourceIBMAtrackerRoutesID returns a reasonable ID for the list.
func dataSourceIBMAtrackerRoutesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceRouteListFlattenRoutes(result []atrackerv1.Route) (routes []map[string]interface{}) {
	for _, routesItem := range result {
		routes = append(routes, dataSourceRouteListRoutesToMap(routesItem))
	}

	return routes
}

func dataSourceRouteListRoutesToMap(routesItem atrackerv1.Route) (routesMap map[string]interface{}) {
	routesMap = map[string]interface{}{}

	if routesItem.ID != nil {
		routesMap["id"] = routesItem.ID
	}
	if routesItem.Name != nil {
		routesMap["name"] = routesItem.Name
	}
	if routesItem.CRN != nil {
		routesMap["crn"] = routesItem.CRN
	}
	if routesItem.Version != nil {
		routesMap["version"] = routesItem.Version
	}
	if routesItem.ReceiveGlobalEvents != nil {
		routesMap["receive_global_events"] = routesItem.ReceiveGlobalEvents
	}
	if routesItem.Rules != nil {
		rulesList := []map[string]interface{}{}
		for _, rulesItem := range routesItem.Rules {
			rulesList = append(rulesList, dataSourceRouteListRoutesRulesToMap(rulesItem))
		}
		routesMap["rules"] = rulesList
	}
	if routesItem.Created != nil {
		routesMap["created"] = routesItem.Created.String()
	}
	if routesItem.Updated != nil {
		routesMap["updated"] = routesItem.Updated.String()
	}

	return routesMap
}

func dataSourceRouteListRoutesRulesToMap(rulesItem atrackerv1.Rule) (rulesMap map[string]interface{}) {
	rulesMap = map[string]interface{}{}

	if rulesItem.TargetIds != nil {
		rulesMap["target_ids"] = rulesItem.TargetIds
	}

	return rulesMap
}
