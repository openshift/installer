// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package atracker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/platform-services-go-sdk/atrackerv2"
)

func DataSourceIBMAtrackerRoutes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMAtrackerRoutesRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the route.",
			},
			"routes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of route resources.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The uuid of the route resource.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the route.",
						},
						"crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The crn of the route resource.",
						},
						"version": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The version of the route.",
						},
						"receive_global_events": {
							Type:        schema.TypeBool,
							Computed:    true,
							Deprecated:  "use rules.locations instead",
							Description: "Indicates whether or not all global events should be forwarded to this region.",
						},
						"rules": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The routing rules that will be evaluated in their order of the array. Once a rule is matched, the remaining rules in the route definition will be skipped.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_ids": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The target ID List. All the events will be send to all targets listed in the rule. You can include targets from other regions.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"locations": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Logs from these locations will be sent to the targets specified. Locations is a superset of regions including global and *.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"created": {
							Type:        schema.TypeString,
							Computed:    true,
							Deprecated:  "use created_at instead",
							Description: "The timestamp of the route creation time.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The timestamp of the route creation time.",
						},
						"updated": {
							Type:        schema.TypeString,
							Computed:    true,
							Deprecated:  "use updated_at instead",
							Description: "The timestamp of the target last updated time.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The timestamp of the route last updated time.",
						},
						"api_version": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The API version of the route.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMAtrackerRoutesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	atrackerClientv2, err := getAtrackerClients(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	listRoutesOptions := &atrackerv2.ListRoutesOptions{}

	routeList, response, err := atrackerClientv2.ListRoutesWithContext(context, listRoutesOptions)
	if err != nil {
		log.Printf("[DEBUG] ListRoutesWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ListRoutesWithContext failed %s\n%s", err, response))
	}

	// Use the provided filter argument and construct a new list with only the requested resource(s)
	var matchRoutes []atrackerv2.Route
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
		d.SetId(DataSourceIBMAtrackerRoutesID(d))
	}

	routes := []map[string]interface{}{}
	if routeList.Routes != nil {
		for _, modelItem := range routeList.Routes {
			modelMap, err := dataSourceIBMAtrackerRoutesRouteToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			routes = append(routes, modelMap)
		}
	}

	if err = d.Set("routes", routes); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting routes %s", err))
	}

	return nil
}

// DataSourceIBMAtrackerRoutesID returns a reasonable ID for the list.
func DataSourceIBMAtrackerRoutesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIBMAtrackerRoutesRouteToMap(model *atrackerv2.Route) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.CRN != nil {
		modelMap["crn"] = *model.CRN
	}
	if model.Version != nil {
		modelMap["version"] = *model.Version
	}
	if model.Rules != nil {
		rules := []map[string]interface{}{}
		for _, rulesItem := range model.Rules {
			rulesItemMap, err := dataSourceIBMAtrackerRoutesRuleToMap(&rulesItem)
			if err != nil {
				return modelMap, err
			}
			rules = append(rules, rulesItemMap)
		}
		modelMap["rules"] = rules
	}
	if model.CreatedAt != nil {
		modelMap["created_at"] = model.CreatedAt.String()
	}
	if model.UpdatedAt != nil {
		modelMap["updated_at"] = model.UpdatedAt.String()
	}
	if model.APIVersion != nil {
		modelMap["api_version"] = *model.APIVersion
	} else {
		modelMap["api_version"] = 1
	}
	return modelMap, nil
}

func dataSourceIBMAtrackerRoutesRuleToMap(model *atrackerv2.Rule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.TargetIds != nil {
		modelMap["target_ids"] = model.TargetIds
	}
	if model.Locations != nil {
		modelMap["locations"] = model.Locations
	}
	return modelMap, nil
}
