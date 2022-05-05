// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	//"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/scc-go-sdk/v3/posturemanagementv2"
)

func DataSourceIBMSccPostureGroupProfileDetails() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMSccPostureGroupProfileDetailsRead,

		Schema: map[string]*schema.Schema{
			"profile_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The profile ID. This can be obtained from the Security and Compliance Center UI by clicking on the profile name. The URL contains the ID.",
			},
			"first": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The URL of a page.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL of a page.",
						},
					},
				},
			},
			"last": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The URL of a page.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL of a page.",
						},
					},
				},
			},
			"previous": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The URL of a page.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL of a page.",
						},
					},
				},
			},
			"controls": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Profiles array.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The identifier number of the control.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the control.",
						},
						"external_control_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The external identifier number of the control.",
						},
						"goals": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Mapped goals aganist the control identifier.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The description of the goal.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The goal ID.",
									},
									"severity": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The severity of the goal.",
									},
									"is_manual": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "The goal is manual check.",
									},
									"is_remediable": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "The goal is remediable or not.",
									},
									"is_reversible": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "The goal is reversible or not.",
									},
									"is_automatable": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "The goal is automatable or not.",
									},
									"is_auto_remediable": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "The goal is autoremediable or not.",
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

func dataSourceIBMSccPostureGroupProfileDetailsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	postureManagementClient, err := meta.(conns.ClientSession).PostureManagementV2()
	if err != nil {
		return diag.FromErr(err)
	}

	getProfileControlsOptions := &posturemanagementv2.GetProfileControlsOptions{}

	getProfileControlsOptions.SetProfileID(d.Get("profile_id").(string))
	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting userDetails %s", err))
	}

	accountID := userDetails.UserAccount
	getProfileControlsOptions.SetAccountID(accountID)

	var controlList *posturemanagementv2.ControlList
	var offset int64
	finalList := []posturemanagementv2.ControlItem{}

	for {

		result, response, err := postureManagementClient.GetProfileControlsWithContext(context, getProfileControlsOptions)
		controlList = result
		if err != nil {
			log.Printf("[DEBUG] GetProfileControlsWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("GetProfileControlsWithContext failed %s\n%s", err, response))
		}
		offset = dataSourceControlListGetNext(result.Next)
		finalList = append(finalList, result.Controls...)
		if offset == 0 {
			break
		}
	}

	controlList.Controls = finalList

	d.SetId(dataSourceIBMSccPostureGroupProfileDetailsID(d))

	if controlList.Controls != nil {
		err = d.Set("controls", dataSourceControlListFlattenControls(controlList.Controls))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting controls %s", err))
		}
	}

	return nil
}

// dataSourceIBMGroupProfileDetailsID returns a reasonable ID for the list.
func dataSourceIBMSccPostureGroupProfileDetailsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceControlListFlattenFirst(result posturemanagementv2.PageLink) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceControlListFirstToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceControlListFirstToMap(firstItem posturemanagementv2.PageLink) (firstMap map[string]interface{}) {
	firstMap = map[string]interface{}{}

	if firstItem.Href != nil {
		firstMap["href"] = firstItem.Href
	}

	return firstMap
}

func dataSourceControlListFlattenLast(result posturemanagementv2.PageLink) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceControlListLastToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceControlListLastToMap(lastItem posturemanagementv2.PageLink) (lastMap map[string]interface{}) {
	lastMap = map[string]interface{}{}

	if lastItem.Href != nil {
		lastMap["href"] = lastItem.Href
	}

	return lastMap
}

func dataSourceControlListFlattenPrevious(result posturemanagementv2.PageLink) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceControlListPreviousToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceControlListPreviousToMap(previousItem posturemanagementv2.PageLink) (previousMap map[string]interface{}) {
	previousMap = map[string]interface{}{}

	if previousItem.Href != nil {
		previousMap["href"] = previousItem.Href
	}

	return previousMap
}

func dataSourceControlListFlattenControls(result []posturemanagementv2.ControlItem) (controls []map[string]interface{}) {
	for _, controlsItem := range result {
		controls = append(controls, dataSourceControlListControlsToMap(controlsItem))
	}

	return controls
}

func dataSourceControlListControlsToMap(controlsItem posturemanagementv2.ControlItem) (controlsMap map[string]interface{}) {
	controlsMap = map[string]interface{}{}

	if controlsItem.ID != nil {
		controlsMap["id"] = controlsItem.ID
	}
	if controlsItem.Description != nil {
		controlsMap["description"] = controlsItem.Description
	}
	if controlsItem.ExternalControlID != nil {
		controlsMap["external_control_id"] = controlsItem.ExternalControlID
	}
	if controlsItem.Goals != nil {
		goalsList := []map[string]interface{}{}
		for _, goalsItem := range controlsItem.Goals {
			goalsList = append(goalsList, dataSourceControlListControlsGoalsToMap(goalsItem))
		}
		controlsMap["goals"] = goalsList
	}

	return controlsMap
}

func dataSourceControlListControlsGoalsToMap(goalsItem posturemanagementv2.GoalItem) (goalsMap map[string]interface{}) {
	goalsMap = map[string]interface{}{}

	if goalsItem.Description != nil {
		goalsMap["description"] = goalsItem.Description
	}
	if goalsItem.ID != nil {
		goalsMap["id"] = goalsItem.ID
	}
	if goalsItem.Severity != nil {
		goalsMap["severity"] = goalsItem.Severity
	}
	if goalsItem.IsManual != nil {
		goalsMap["is_manual"] = goalsItem.IsManual
	}
	if goalsItem.IsRemediable != nil {
		goalsMap["is_remediable"] = goalsItem.IsRemediable
	}
	if goalsItem.IsReversible != nil {
		goalsMap["is_reversible"] = goalsItem.IsReversible
	}
	if goalsItem.IsAutomatable != nil {
		goalsMap["is_automatable"] = goalsItem.IsAutomatable
	}
	if goalsItem.IsAutoRemediable != nil {
		goalsMap["is_auto_remediable"] = goalsItem.IsAutoRemediable
	}

	return goalsMap
}

func dataSourceControlListGetNext(next interface{}) int64 {
	if reflect.ValueOf(next).IsNil() {
		return 0
	}

	u, err := url.Parse(reflect.ValueOf(next).Elem().FieldByName("Href").Elem().String())
	if err != nil {
		return 0
	}

	q := u.Query()
	var page string

	if q.Get("start") != "" {
		page = q.Get("start")
	} else if q.Get("offset") != "" {
		page = q.Get("offset")
	}

	convertedVal, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		return 0
	}
	return convertedVal
}
