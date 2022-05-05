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

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/scc-go-sdk/v3/posturemanagementv2"
)

func DataSourceIBMSccPostureProfiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMSccPostureListProfilesRead,

		Schema: map[string]*schema.Schema{
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
			"profiles": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Profiles.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the profile.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A description of the profile.",
						},
						"version": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The version of the profile.",
						},
						"created_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user who created the profile.",
						},
						"modified_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user who last modified the profile.",
						},
						"reason_for_delete": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A reason that you want to delete a profile.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An auto-generated unique identifying number of the profile.",
						},
						"base_profile": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The base profile that the controls are pulled from.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of profile.",
						},
						"no_of_controls": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "no of Controls.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time that the profile was created in UTC.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time that the profile was most recently modified in UTC.",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "The profile status. If the profile is enabled, the value is true. If the profile is disabled, the value is false.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMSccPostureListProfilesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	postureManagementClient, err := meta.(conns.ClientSession).PostureManagementV2()
	if err != nil {
		return diag.FromErr(err)
	}

	listProfilesOptions := &posturemanagementv2.ListProfilesOptions{}
	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting userDetails %s", err))
	}

	accountID := userDetails.UserAccount
	listProfilesOptions.SetAccountID(accountID)

	var profileList *posturemanagementv2.ProfileList
	var offset int64
	finalList := []posturemanagementv2.Profile{}

	for {
		listProfilesOptions.Offset = &offset

		listProfilesOptions.Limit = core.Int64Ptr(int64(100))
		result, response, err := postureManagementClient.ListProfilesWithContext(context, listProfilesOptions)
		profileList = result
		if err != nil {
			log.Printf("[DEBUG] ListProfilesWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("ListProfilesWithContext failed %s\n%s", err, response))
		}
		offset = dataSourceProfileListGetNext(result.Next)
		finalList = append(finalList, result.Profiles...)
		if offset == 0 {
			break
		}
	}

	profileList.Profiles = finalList

	d.SetId(dataSourceIBMSccPostureListProfilesID(d))

	if profileList.First != nil {
		err = d.Set("first", dataSourceProfileListFlattenFirst(*profileList.First))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting first %s", err))
		}
	}

	if profileList.Last != nil {
		err = d.Set("last", dataSourceProfileListFlattenLast(*profileList.Last))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting last %s", err))
		}
	}

	if profileList.Previous != nil {
		err = d.Set("previous", dataSourceProfileListFlattenPrevious(*profileList.Previous))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting previous %s", err))
		}
	}

	if profileList.Profiles != nil {
		err = d.Set("profiles", dataSourceProfileListFlattenProfiles(profileList.Profiles))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting profiles %s", err))
		}
	}

	return nil
}

// dataSourceIBMListProfilesID returns a reasonable ID for the list.
func dataSourceIBMSccPostureListProfilesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceProfileListFlattenFirst(result posturemanagementv2.PageLink) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceProfileListFirstToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceProfileListFirstToMap(firstItem posturemanagementv2.PageLink) (firstMap map[string]interface{}) {
	firstMap = map[string]interface{}{}

	if firstItem.Href != nil {
		firstMap["href"] = firstItem.Href
	}

	return firstMap
}

func dataSourceProfileListFlattenLast(result posturemanagementv2.PageLink) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceProfileListLastToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceProfileListLastToMap(lastItem posturemanagementv2.PageLink) (lastMap map[string]interface{}) {
	lastMap = map[string]interface{}{}

	if lastItem.Href != nil {
		lastMap["href"] = lastItem.Href
	}

	return lastMap
}

func dataSourceProfileListFlattenPrevious(result posturemanagementv2.PageLink) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceProfileListPreviousToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceProfileListPreviousToMap(previousItem posturemanagementv2.PageLink) (previousMap map[string]interface{}) {
	previousMap = map[string]interface{}{}

	if previousItem.Href != nil {
		previousMap["href"] = previousItem.Href
	}

	return previousMap
}

func dataSourceProfileListFlattenProfiles(result []posturemanagementv2.Profile) (profiles []map[string]interface{}) {
	for _, profilesItem := range result {
		profiles = append(profiles, dataSourceProfileListProfilesToMap(profilesItem))
	}

	return profiles
}

func dataSourceProfileListProfilesToMap(profilesItem posturemanagementv2.Profile) (profilesMap map[string]interface{}) {
	profilesMap = map[string]interface{}{}

	if profilesItem.Name != nil {
		profilesMap["name"] = profilesItem.Name
	}
	if profilesItem.Description != nil {
		profilesMap["description"] = profilesItem.Description
	}
	if profilesItem.Version != nil {
		profilesMap["version"] = profilesItem.Version
	}
	if profilesItem.CreatedBy != nil {
		profilesMap["created_by"] = profilesItem.CreatedBy
	}
	if profilesItem.ModifiedBy != nil {
		profilesMap["modified_by"] = profilesItem.ModifiedBy
	}
	if profilesItem.ReasonForDelete != nil {
		profilesMap["reason_for_delete"] = profilesItem.ReasonForDelete
	}
	if profilesItem.ID != nil {
		profilesMap["id"] = profilesItem.ID
	}
	if profilesItem.BaseProfile != nil {
		profilesMap["base_profile"] = profilesItem.BaseProfile
	}
	if profilesItem.Type != nil {
		profilesMap["type"] = profilesItem.Type
	}
	if profilesItem.NoOfControls != nil {
		profilesMap["no_of_controls"] = profilesItem.NoOfControls
	}
	if profilesItem.CreatedAt != nil {
		profilesMap["created_at"] = profilesItem.CreatedAt.String()
	}
	if profilesItem.UpdatedAt != nil {
		profilesMap["updated_at"] = profilesItem.UpdatedAt.String()
	}
	if profilesItem.Enabled != nil {
		profilesMap["enabled"] = profilesItem.Enabled
	}

	return profilesMap
}

func dataSourceProfileListGetNext(next interface{}) int64 {
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
