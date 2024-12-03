// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package logs

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/logs-go-sdk/logsv0"
)

func DataSourceIbmLogsView() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmLogsViewRead,

		Schema: map[string]*schema.Schema{
			"logs_view_id": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "View ID.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "View name.",
			},
			"search_query": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "View search query.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"query": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "View search query.",
						},
					},
				},
			},
			"time_selection": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "View time selection.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"quick_selection": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Quick time selection.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"caption": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Quick time selection caption.",
									},
									"seconds": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Quick time selection amount of seconds.",
									},
								},
							},
						},
						"custom_selection": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Custom time selection.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"from_time": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Custom time selection start timestamp.",
									},
									"to_time": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Custom time selection end timestamp.",
									},
								},
							},
						},
					},
				},
			},
			"filters": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "View selected filters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filters": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Selected filters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Filter name.",
									},
									"selected_values": &schema.Schema{
										Type:        schema.TypeMap,
										Computed:    true,
										Description: "Filter selected values.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"folder_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "View folder ID.",
			},
		},
	}
}

func dataSourceIbmLogsViewRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_view", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	region := getLogsInstanceRegion(logsClient, d)
	instanceId := d.Get("instance_id").(string)
	logsClient = getClientWithLogsInstanceEndpoint(logsClient, instanceId, region, getLogsInstanceEndpointType(logsClient, d))

	getViewOptions := &logsv0.GetViewOptions{}

	getViewOptions.SetID(int64(d.Get("logs_view_id").(int)))

	view, _, err := logsClient.GetViewWithContext(context, getViewOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetViewWithContext failed: %s", err.Error()), "(Data) ibm_logs_view", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%d", *getViewOptions.ID))

	if err = d.Set("name", view.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_logs_view", "read")
		return tfErr.GetDiag()
	}

	searchQuery := []map[string]interface{}{}
	if view.SearchQuery != nil {
		modelMap, err := DataSourceIbmLogsViewApisViewsV1SearchQueryToMap(view.SearchQuery)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_view", "read")
			return tfErr.GetDiag()
		}
		searchQuery = append(searchQuery, modelMap)
	}
	if err = d.Set("search_query", searchQuery); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting search_query: %s", err), "(Data) ibm_logs_view", "read")
		return tfErr.GetDiag()
	}

	timeSelection := []map[string]interface{}{}
	if view.TimeSelection != nil {
		modelMap, err := DataSourceIbmLogsViewApisViewsV1TimeSelectionToMap(view.TimeSelection)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_view", "read")
			return tfErr.GetDiag()
		}
		timeSelection = append(timeSelection, modelMap)
	}
	if err = d.Set("time_selection", timeSelection); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting time_selection: %s", err), "(Data) ibm_logs_view", "read")
		return tfErr.GetDiag()
	}

	filters := []map[string]interface{}{}
	if view.Filters != nil {
		modelMap, err := DataSourceIbmLogsViewApisViewsV1SelectedFiltersToMap(view.Filters)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_view", "read")
			return tfErr.GetDiag()
		}
		filters = append(filters, modelMap)
	}
	if err = d.Set("filters", filters); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting filters: %s", err), "(Data) ibm_logs_view", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("folder_id", view.FolderID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting folder_id: %s", err), "(Data) ibm_logs_view", "read")
		return tfErr.GetDiag()
	}

	return nil
}

func DataSourceIbmLogsViewApisViewsV1SearchQueryToMap(model *logsv0.ApisViewsV1SearchQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["query"] = *model.Query
	return modelMap, nil
}

func DataSourceIbmLogsViewApisViewsV1TimeSelectionToMap(model logsv0.ApisViewsV1TimeSelectionIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisViewsV1TimeSelectionSelectionTypeQuickSelection); ok {
		return DataSourceIbmLogsViewApisViewsV1TimeSelectionSelectionTypeQuickSelectionToMap(model.(*logsv0.ApisViewsV1TimeSelectionSelectionTypeQuickSelection))
	} else if _, ok := model.(*logsv0.ApisViewsV1TimeSelectionSelectionTypeCustomSelection); ok {
		return DataSourceIbmLogsViewApisViewsV1TimeSelectionSelectionTypeCustomSelectionToMap(model.(*logsv0.ApisViewsV1TimeSelectionSelectionTypeCustomSelection))
	} else if _, ok := model.(*logsv0.ApisViewsV1TimeSelection); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisViewsV1TimeSelection)
		if model.QuickSelection != nil {
			quickSelectionMap, err := DataSourceIbmLogsViewApisViewsV1QuickTimeSelectionToMap(model.QuickSelection)
			if err != nil {
				return modelMap, err
			}
			modelMap["quick_selection"] = []map[string]interface{}{quickSelectionMap}
		}
		if model.CustomSelection != nil {
			customSelectionMap, err := DataSourceIbmLogsViewApisViewsV1CustomTimeSelectionToMap(model.CustomSelection)
			if err != nil {
				return modelMap, err
			}
			modelMap["custom_selection"] = []map[string]interface{}{customSelectionMap}
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized logsv0.ApisViewsV1TimeSelectionIntf subtype encountered")
	}
}

func DataSourceIbmLogsViewApisViewsV1QuickTimeSelectionToMap(model *logsv0.ApisViewsV1QuickTimeSelection) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["caption"] = *model.Caption
	modelMap["seconds"] = flex.IntValue(model.Seconds)
	return modelMap, nil
}

func DataSourceIbmLogsViewApisViewsV1CustomTimeSelectionToMap(model *logsv0.ApisViewsV1CustomTimeSelection) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["from_time"] = model.FromTime.String()
	modelMap["to_time"] = model.ToTime.String()
	return modelMap, nil
}

func DataSourceIbmLogsViewApisViewsV1TimeSelectionSelectionTypeQuickSelectionToMap(model *logsv0.ApisViewsV1TimeSelectionSelectionTypeQuickSelection) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.QuickSelection != nil {
		quickSelectionMap, err := DataSourceIbmLogsViewApisViewsV1QuickTimeSelectionToMap(model.QuickSelection)
		if err != nil {
			return modelMap, err
		}
		modelMap["quick_selection"] = []map[string]interface{}{quickSelectionMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsViewApisViewsV1TimeSelectionSelectionTypeCustomSelectionToMap(model *logsv0.ApisViewsV1TimeSelectionSelectionTypeCustomSelection) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CustomSelection != nil {
		customSelectionMap, err := DataSourceIbmLogsViewApisViewsV1CustomTimeSelectionToMap(model.CustomSelection)
		if err != nil {
			return modelMap, err
		}
		modelMap["custom_selection"] = []map[string]interface{}{customSelectionMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsViewApisViewsV1SelectedFiltersToMap(model *logsv0.ApisViewsV1SelectedFilters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := DataSourceIbmLogsViewApisViewsV1FilterToMap(&filtersItem)
			if err != nil {
				return modelMap, err
			}
			filters = append(filters, filtersItemMap)
		}
		modelMap["filters"] = filters
	}
	return modelMap, nil
}

func DataSourceIbmLogsViewApisViewsV1FilterToMap(model *logsv0.ApisViewsV1Filter) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = *model.Name
	selectedValues := make(map[string]interface{})
	for k, v := range model.SelectedValues {
		selectedValues[k] = flex.Stringify(v)
	}
	modelMap["selected_values"] = selectedValues
	return modelMap, nil
}
