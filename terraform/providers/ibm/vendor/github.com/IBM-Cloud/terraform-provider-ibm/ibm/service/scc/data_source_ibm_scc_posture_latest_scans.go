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

func DataSourceIBMSccPostureLatestScans() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMSccPostureListLatestScansRead,

		Schema: map[string]*schema.Schema{
			"scan_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the scan.",
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
			"latest_scans": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The details of a scan.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scan_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the scan.",
						},
						"scan_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A system generated name that is the combination of 12 characters in the scope name and 12 characters of a profile name.",
						},
						"scope_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The scope ID of the scan.",
						},
						"scope_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the scope.",
						},
						"profiles": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Profiles array.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the profile.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "An auto-generated unique identifier for the scope.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of profile.",
									},
								},
							},
						},
						"group_profile_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The group ID of profile.",
						},
						"group_profile_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The group name of the profile.",
						},
						"report_run_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The entity that ran the report.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time the scan was run.",
						},
						"report_setting_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique ID for Scan that is created.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time the scan completed.",
						},
						"result": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The result of a scan.The above values will not be avaialble if no scopes are available.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"goals_pass_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of goals that passed the scan.",
									},
									"goals_unable_to_perform_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of goals that could not be validated. A control is listed as 'Unable to perform' when information about its associated resource can't be collected.",
									},
									"goals_not_applicable_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of goals that are not relevant to the current scan. A scan is listed as 'Not applicable' when information about its associated resource can't be found.",
									},
									"goals_fail_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of goals that failed the scan.",
									},
									"goals_total_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The total number of goals that were included in the scan.",
									},
									"controls_pass_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of controls that passed the scan.",
									},
									"controls_fail_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of controls that failed the scan.",
									},
									"controls_not_applicable_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of controls that are not relevant to the current scan. A scan is listed as 'Not applicable' when information about its associated resource can't be found.",
									},
									"controls_unable_to_perform_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of controls that could not be validated. A control is listed as 'Unable to perform' when information about its associated resource can't be collected.",
									},
									"controls_total_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The total number of controls that were included in the scan.",
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

func dataSourceIBMSccPostureListLatestScansRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	postureManagementClient, err := meta.(conns.ClientSession).PostureManagementV2()
	if err != nil {
		return diag.FromErr(err)
	}

	listLatestScansOptions := &posturemanagementv2.ListLatestScansOptions{}
	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting userDetails %s", err))
	}

	accountID := userDetails.UserAccount
	listLatestScansOptions.SetAccountID(accountID)

	var scanList *posturemanagementv2.ScanList
	var offset int64
	finalList := []posturemanagementv2.ScanItem{}
	var scanID string
	var suppliedFilter bool

	if v, ok := d.GetOk("scan_id"); ok {
		scanID = v.(string)
		suppliedFilter = true
	}

	for {
		listLatestScansOptions.Offset = &offset

		listLatestScansOptions.Limit = core.Int64Ptr(int64(100))
		result, response, err := postureManagementClient.ListLatestScansWithContext(context, listLatestScansOptions)
		scanList = result
		if err != nil {
			log.Printf("[DEBUG] ListLatestScansWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("ListLatestScansWithContext failed %s\n%s", err, response))
		}
		offset = dataSourceScanListGetNext(result.Next)
		if suppliedFilter {
			for _, data := range result.LatestScans {
				if *data.ScanID == scanID {
					finalList = append(finalList, data)
				}
			}
		} else {
			finalList = append(finalList, result.LatestScans...)
		}
		if offset == 0 {
			break
		}
	}

	scanList.LatestScans = finalList

	if suppliedFilter {
		if len(scanList.LatestScans) == 0 {
			return diag.FromErr(fmt.Errorf("no LatestScans found with scanID %s", scanID))
		}
		d.SetId(scanID)
	} else {
		d.SetId(dataSourceIBMListLatestScansID(d))
	}

	if scanList.First != nil {
		err = d.Set("first", dataSourceScanListFlattenFirst(*scanList.First))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting first %s", err))
		}
	}

	if scanList.Last != nil {
		err = d.Set("last", dataSourceScanListFlattenLast(*scanList.Last))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting last %s", err))
		}
	}

	if scanList.Previous != nil {
		err = d.Set("previous", dataSourceScanListFlattenPrevious(*scanList.Previous))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting previous %s", err))
		}
	}

	if scanList.LatestScans != nil {
		err = d.Set("latest_scans", dataSourceScanListFlattenLatestScans(scanList.LatestScans))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting latest_scans %s", err))
		}
	}

	return nil
}

// dataSourceIBMListLatestScansID returns a reasonable ID for the list.
func dataSourceIBMListLatestScansID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceScanListFlattenFirst(result posturemanagementv2.PageLink) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceScanListFirstToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceScanListFirstToMap(firstItem posturemanagementv2.PageLink) (firstMap map[string]interface{}) {
	firstMap = map[string]interface{}{}

	if firstItem.Href != nil {
		firstMap["href"] = firstItem.Href
	}

	return firstMap
}

func dataSourceScanListFlattenLast(result posturemanagementv2.PageLink) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceScanListLastToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceScanListLastToMap(lastItem posturemanagementv2.PageLink) (lastMap map[string]interface{}) {
	lastMap = map[string]interface{}{}

	if lastItem.Href != nil {
		lastMap["href"] = lastItem.Href
	}

	return lastMap
}

func dataSourceScanListFlattenPrevious(result posturemanagementv2.PageLink) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceScanListPreviousToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceScanListPreviousToMap(previousItem posturemanagementv2.PageLink) (previousMap map[string]interface{}) {
	previousMap = map[string]interface{}{}

	if previousItem.Href != nil {
		previousMap["href"] = previousItem.Href
	}

	return previousMap
}

func dataSourceScanListFlattenLatestScans(result []posturemanagementv2.ScanItem) (latestScans []map[string]interface{}) {
	for _, latestScansItem := range result {
		latestScans = append(latestScans, dataSourceScanListLatestScansToMap(latestScansItem))
	}

	return latestScans
}

func dataSourceScanListLatestScansToMap(latestScansItem posturemanagementv2.ScanItem) (latestScansMap map[string]interface{}) {
	latestScansMap = map[string]interface{}{}

	if latestScansItem.ScanID != nil {
		latestScansMap["scan_id"] = latestScansItem.ScanID
	}
	if latestScansItem.ScanName != nil {
		latestScansMap["scan_name"] = latestScansItem.ScanName
	}
	if latestScansItem.ScopeID != nil {
		latestScansMap["scope_id"] = latestScansItem.ScopeID
	}
	if latestScansItem.ScopeName != nil {
		latestScansMap["scope_name"] = latestScansItem.ScopeName
	}
	if latestScansItem.Profiles != nil {
		profilesList := []map[string]interface{}{}
		for _, profilesItem := range latestScansItem.Profiles {
			profilesList = append(profilesList, dataSourceScanListLatestScansProfilesToMap(profilesItem))
		}
		latestScansMap["profiles"] = profilesList
	}
	if latestScansItem.GroupProfileID != nil {
		latestScansMap["group_profile_id"] = latestScansItem.GroupProfileID
	}
	if latestScansItem.GroupProfileName != nil {
		latestScansMap["group_profile_name"] = latestScansItem.GroupProfileName
	}
	if latestScansItem.ReportRunBy != nil {
		latestScansMap["report_run_by"] = latestScansItem.ReportRunBy
	}
	if latestScansItem.StartTime != nil {
		latestScansMap["start_time"] = latestScansItem.StartTime.String()
	}
	if latestScansItem.ReportSettingID != nil {
		latestScansMap["report_setting_id"] = latestScansItem.ReportSettingID
	}
	if latestScansItem.EndTime != nil {
		latestScansMap["end_time"] = latestScansItem.EndTime.String()
	}
	if latestScansItem.Result != nil {
		resultList := []map[string]interface{}{}
		resultMap := dataSourceScanListLatestScansResultToMap(*latestScansItem.Result)
		resultList = append(resultList, resultMap)
		latestScansMap["result"] = resultList
	}

	return latestScansMap
}

func dataSourceScanListLatestScansProfilesToMap(profilesItem posturemanagementv2.ProfileItem) (profilesMap map[string]interface{}) {
	profilesMap = map[string]interface{}{}

	if profilesItem.Name != nil {
		profilesMap["name"] = profilesItem.Name
	}
	if profilesItem.ID != nil {
		profilesMap["id"] = profilesItem.ID
	}
	if profilesItem.Type != nil {
		profilesMap["type"] = profilesItem.Type
	}

	return profilesMap
}

func dataSourceScanListLatestScansResultToMap(resultItem posturemanagementv2.ScanResult) (resultMap map[string]interface{}) {
	resultMap = map[string]interface{}{}

	if resultItem.GoalsPassCount != nil {
		resultMap["goals_pass_count"] = resultItem.GoalsPassCount
	}
	if resultItem.GoalsUnableToPerformCount != nil {
		resultMap["goals_unable_to_perform_count"] = resultItem.GoalsUnableToPerformCount
	}
	if resultItem.GoalsNotApplicableCount != nil {
		resultMap["goals_not_applicable_count"] = resultItem.GoalsNotApplicableCount
	}
	if resultItem.GoalsFailCount != nil {
		resultMap["goals_fail_count"] = resultItem.GoalsFailCount
	}
	if resultItem.GoalsTotalCount != nil {
		resultMap["goals_total_count"] = resultItem.GoalsTotalCount
	}
	if resultItem.ControlsPassCount != nil {
		resultMap["controls_pass_count"] = resultItem.ControlsPassCount
	}
	if resultItem.ControlsFailCount != nil {
		resultMap["controls_fail_count"] = resultItem.ControlsFailCount
	}
	if resultItem.ControlsNotApplicableCount != nil {
		resultMap["controls_not_applicable_count"] = resultItem.ControlsNotApplicableCount
	}
	if resultItem.ControlsUnableToPerformCount != nil {
		resultMap["controls_unable_to_perform_count"] = resultItem.ControlsUnableToPerformCount
	}
	if resultItem.ControlsTotalCount != nil {
		resultMap["controls_total_count"] = resultItem.ControlsTotalCount
	}

	return resultMap
}

func dataSourceScanListGetNext(next interface{}) int64 {
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
