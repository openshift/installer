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

func DataSourceIBMSccPostureScanSummaries() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMSccPostureScanSummariesRead,

		Schema: map[string]*schema.Schema{
			"report_setting_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The report setting ID. This can be obtained from the /validations/latest_scans API call.",
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
			"summaries": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Summaries.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the scan.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A system generated name that is the combination of 12 characters in the scope name and 12 characters of a profile name.",
						},
						"scope_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the scope.",
						},
						"scope_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the scope.",
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
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time the scan completed.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the collector as it completes a scan.",
						},
						"profiles": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The list of profiles.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the profile.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the profile.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of profile. To learn more about profile types, check out the [docs] (https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-profiles).",
									},
									"validation_result": {
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
						"group_profiles": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The list of group profiles.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the profile.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the profile.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of profile. To learn more about profile types, check out the [docs] (https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-profiles).",
									},
									"validation_result": {
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
				},
			},
		},
	}
}

func dataSourceIBMSccPostureScanSummariesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	postureManagementClient, err := meta.(conns.ClientSession).PostureManagementV2()
	if err != nil {
		return diag.FromErr(err)
	}

	scanSummariesOptions := &posturemanagementv2.ScanSummariesOptions{}
	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting userDetails %s", err))
	}

	accountID := userDetails.UserAccount
	scanSummariesOptions.SetAccountID(accountID)

	scanSummariesOptions.SetReportSettingID(d.Get("report_setting_id").(string))

	var summaryList *posturemanagementv2.SummaryList
	var offset int64
	finalList := []posturemanagementv2.SummaryItem{}

	for {
		scanSummariesOptions.Offset = &offset

		scanSummariesOptions.Limit = core.Int64Ptr(int64(100))
		result, response, err := postureManagementClient.ScanSummariesWithContext(context, scanSummariesOptions)
		summaryList = result
		if err != nil {
			log.Printf("[DEBUG] ScanSummariesWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("ScanSummariesWithContext failed %s\n%s", err, response))
		}
		offset = dataSourceSummaryListGetNext(result.Next)
		finalList = append(finalList, result.Summaries...)
		if offset == 0 {
			break
		}
	}

	summaryList.Summaries = finalList

	d.SetId(dataSourceIBMSccPostureScanSummariesID(d))

	if summaryList.First != nil {
		err = d.Set("first", dataSourceSummaryListFlattenFirst(*summaryList.First))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting first %s", err))
		}
	}

	if summaryList.Last != nil {
		err = d.Set("last", dataSourceSummaryListFlattenLast(*summaryList.Last))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting last %s", err))
		}
	}

	if summaryList.Previous != nil {
		err = d.Set("previous", dataSourceSummaryListFlattenPrevious(*summaryList.Previous))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting previous %s", err))
		}
	}

	if summaryList.Summaries != nil {
		err = d.Set("summaries", dataSourceSummaryListFlattenSummaries(summaryList.Summaries))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting summaries %s", err))
		}
	}

	return nil
}

// dataSourceIBMScanSummariesID returns a reasonable ID for the list.
func dataSourceIBMSccPostureScanSummariesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceSummaryListFlattenFirst(result posturemanagementv2.PageLink) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceSummaryListFirstToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceSummaryListFirstToMap(firstItem posturemanagementv2.PageLink) (firstMap map[string]interface{}) {
	firstMap = map[string]interface{}{}

	if firstItem.Href != nil {
		firstMap["href"] = firstItem.Href
	}

	return firstMap
}

func dataSourceSummaryListFlattenLast(result posturemanagementv2.PageLink) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceSummaryListLastToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceSummaryListLastToMap(lastItem posturemanagementv2.PageLink) (lastMap map[string]interface{}) {
	lastMap = map[string]interface{}{}

	if lastItem.Href != nil {
		lastMap["href"] = lastItem.Href
	}

	return lastMap
}

func dataSourceSummaryListFlattenPrevious(result posturemanagementv2.PageLink) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceSummaryListPreviousToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceSummaryListPreviousToMap(previousItem posturemanagementv2.PageLink) (previousMap map[string]interface{}) {
	previousMap = map[string]interface{}{}

	if previousItem.Href != nil {
		previousMap["href"] = previousItem.Href
	}

	return previousMap
}

func dataSourceSummaryListFlattenSummaries(result []posturemanagementv2.SummaryItem) (summaries []map[string]interface{}) {
	for _, summariesItem := range result {
		summaries = append(summaries, dataSourceSummaryListSummariesToMap(summariesItem))
	}

	return summaries
}

func dataSourceSummaryListSummariesToMap(summariesItem posturemanagementv2.SummaryItem) (summariesMap map[string]interface{}) {
	summariesMap = map[string]interface{}{}

	if summariesItem.ID != nil {
		summariesMap["id"] = summariesItem.ID
	}
	if summariesItem.Name != nil {
		summariesMap["name"] = summariesItem.Name
	}
	if summariesItem.ScopeID != nil {
		summariesMap["scope_id"] = summariesItem.ScopeID
	}
	if summariesItem.ScopeName != nil {
		summariesMap["scope_name"] = summariesItem.ScopeName
	}
	if summariesItem.ReportRunBy != nil {
		summariesMap["report_run_by"] = summariesItem.ReportRunBy
	}
	if summariesItem.StartTime != nil {
		summariesMap["start_time"] = summariesItem.StartTime.String()
	}
	if summariesItem.EndTime != nil {
		summariesMap["end_time"] = summariesItem.EndTime.String()
	}
	if summariesItem.Status != nil {
		summariesMap["status"] = summariesItem.Status
	}
	if summariesItem.Profiles != nil {
		profilesList := []map[string]interface{}{}
		for _, profilesItem := range summariesItem.Profiles {
			profilesList = append(profilesList, dataSourceSummaryListSummariesProfilesToMap(profilesItem))
		}
		summariesMap["profiles"] = profilesList
	}
	if summariesItem.GroupProfiles != nil {
		groupProfilesList := []map[string]interface{}{}
		for _, groupProfilesItem := range summariesItem.GroupProfiles {
			groupProfilesList = append(groupProfilesList, dataSourceSummaryListSummariesGroupProfilesToMap(groupProfilesItem))
		}
		summariesMap["group_profiles"] = groupProfilesList
	}

	return summariesMap
}

func dataSourceSummaryListSummariesProfilesToMap(profilesItem posturemanagementv2.ProfileResult) (profilesMap map[string]interface{}) {
	profilesMap = map[string]interface{}{}

	if profilesItem.ID != nil {
		profilesMap["id"] = profilesItem.ID
	}
	if profilesItem.Name != nil {
		profilesMap["name"] = profilesItem.Name
	}
	if profilesItem.Type != nil {
		profilesMap["type"] = profilesItem.Type
	}
	if profilesItem.ValidationResult != nil {
		validationResultList := []map[string]interface{}{}
		validationResultMap := dataSourceSummaryListProfilesValidationResultToMap(*profilesItem.ValidationResult)
		validationResultList = append(validationResultList, validationResultMap)
		profilesMap["validation_result"] = validationResultList
	}

	return profilesMap
}

func dataSourceSummaryListProfilesValidationResultToMap(validationResultItem posturemanagementv2.ScanResult) (validationResultMap map[string]interface{}) {
	validationResultMap = map[string]interface{}{}

	if validationResultItem.GoalsPassCount != nil {
		validationResultMap["goals_pass_count"] = validationResultItem.GoalsPassCount
	}
	if validationResultItem.GoalsUnableToPerformCount != nil {
		validationResultMap["goals_unable_to_perform_count"] = validationResultItem.GoalsUnableToPerformCount
	}
	if validationResultItem.GoalsNotApplicableCount != nil {
		validationResultMap["goals_not_applicable_count"] = validationResultItem.GoalsNotApplicableCount
	}
	if validationResultItem.GoalsFailCount != nil {
		validationResultMap["goals_fail_count"] = validationResultItem.GoalsFailCount
	}
	if validationResultItem.GoalsTotalCount != nil {
		validationResultMap["goals_total_count"] = validationResultItem.GoalsTotalCount
	}
	if validationResultItem.ControlsPassCount != nil {
		validationResultMap["controls_pass_count"] = validationResultItem.ControlsPassCount
	}
	if validationResultItem.ControlsFailCount != nil {
		validationResultMap["controls_fail_count"] = validationResultItem.ControlsFailCount
	}
	if validationResultItem.ControlsNotApplicableCount != nil {
		validationResultMap["controls_not_applicable_count"] = validationResultItem.ControlsNotApplicableCount
	}
	if validationResultItem.ControlsUnableToPerformCount != nil {
		validationResultMap["controls_unable_to_perform_count"] = validationResultItem.ControlsUnableToPerformCount
	}
	if validationResultItem.ControlsTotalCount != nil {
		validationResultMap["controls_total_count"] = validationResultItem.ControlsTotalCount
	}

	return validationResultMap
}

func dataSourceSummaryListSummariesGroupProfilesToMap(groupProfilesItem posturemanagementv2.ProfileResult) (groupProfilesMap map[string]interface{}) {
	groupProfilesMap = map[string]interface{}{}

	if groupProfilesItem.ID != nil {
		groupProfilesMap["id"] = groupProfilesItem.ID
	}
	if groupProfilesItem.Name != nil {
		groupProfilesMap["name"] = groupProfilesItem.Name
	}
	if groupProfilesItem.Type != nil {
		groupProfilesMap["type"] = groupProfilesItem.Type
	}
	if groupProfilesItem.ValidationResult != nil {
		validationResultList := []map[string]interface{}{}
		validationResultMap := dataSourceSummaryListGroupProfilesValidationResultToMap(*groupProfilesItem.ValidationResult)
		validationResultList = append(validationResultList, validationResultMap)
		groupProfilesMap["validation_result"] = validationResultList
	}

	return groupProfilesMap
}

func dataSourceSummaryListGroupProfilesValidationResultToMap(validationResultItem posturemanagementv2.ScanResult) (validationResultMap map[string]interface{}) {
	validationResultMap = map[string]interface{}{}

	if validationResultItem.GoalsPassCount != nil {
		validationResultMap["goals_pass_count"] = validationResultItem.GoalsPassCount
	}
	if validationResultItem.GoalsUnableToPerformCount != nil {
		validationResultMap["goals_unable_to_perform_count"] = validationResultItem.GoalsUnableToPerformCount
	}
	if validationResultItem.GoalsNotApplicableCount != nil {
		validationResultMap["goals_not_applicable_count"] = validationResultItem.GoalsNotApplicableCount
	}
	if validationResultItem.GoalsFailCount != nil {
		validationResultMap["goals_fail_count"] = validationResultItem.GoalsFailCount
	}
	if validationResultItem.GoalsTotalCount != nil {
		validationResultMap["goals_total_count"] = validationResultItem.GoalsTotalCount
	}
	if validationResultItem.ControlsPassCount != nil {
		validationResultMap["controls_pass_count"] = validationResultItem.ControlsPassCount
	}
	if validationResultItem.ControlsFailCount != nil {
		validationResultMap["controls_fail_count"] = validationResultItem.ControlsFailCount
	}
	if validationResultItem.ControlsNotApplicableCount != nil {
		validationResultMap["controls_not_applicable_count"] = validationResultItem.ControlsNotApplicableCount
	}
	if validationResultItem.ControlsUnableToPerformCount != nil {
		validationResultMap["controls_unable_to_perform_count"] = validationResultItem.ControlsUnableToPerformCount
	}
	if validationResultItem.ControlsTotalCount != nil {
		validationResultMap["controls_total_count"] = validationResultItem.ControlsTotalCount
	}

	return validationResultMap
}

func dataSourceSummaryListGetNext(next interface{}) int64 {
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
