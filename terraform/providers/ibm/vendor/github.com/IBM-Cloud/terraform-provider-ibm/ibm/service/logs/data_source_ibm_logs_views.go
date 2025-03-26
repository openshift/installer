// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package logs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/logs-go-sdk/logsv0"
)

func DataSourceIbmLogsViews() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmLogsViewsRead,

		Schema: map[string]*schema.Schema{
			"views": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of views.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
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
				},
			},
		},
	}
}

func dataSourceIbmLogsViewsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_views", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	region := getLogsInstanceRegion(logsClient, d)
	instanceId := d.Get("instance_id").(string)
	logsClient = getClientWithLogsInstanceEndpoint(logsClient, instanceId, region, getLogsInstanceEndpointType(logsClient, d))

	listViewsOptions := &logsv0.ListViewsOptions{}

	viewCollection, _, err := logsClient.ListViewsWithContext(context, listViewsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListViewsWithContext failed: %s", err.Error()), "(Data) ibm_logs_views", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIbmLogsViewsID(d))

	views := []map[string]interface{}{}
	if viewCollection.Views != nil {
		for _, modelItem := range viewCollection.Views {
			modelMap, err := DataSourceIbmLogsViewsViewToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_views", "read")
				return tfErr.GetDiag()
			}
			views = append(views, modelMap)
		}
	}
	if err = d.Set("views", views); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting views: %s", err), "(Data) ibm_logs_views", "read")
		return tfErr.GetDiag()
	}

	return nil
}

// dataSourceIbmLogsViewsID returns a reasonable ID for the list.
func dataSourceIbmLogsViewsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIbmLogsViewsViewToMap(model *logsv0.View) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = flex.IntValue(model.ID)
	modelMap["name"] = *model.Name
	if model.SearchQuery != nil {
		searchQueryMap, err := DataSourceIbmLogsViewsApisViewsV1SearchQueryToMap(model.SearchQuery)
		if err != nil {
			return modelMap, err
		}
		modelMap["search_query"] = []map[string]interface{}{searchQueryMap}
	}
	timeSelectionMap, err := DataSourceIbmLogsViewsApisViewsV1TimeSelectionToMap(model.TimeSelection)
	if err != nil {
		return modelMap, err
	}
	modelMap["time_selection"] = []map[string]interface{}{timeSelectionMap}
	if model.Filters != nil {
		filtersMap, err := DataSourceIbmLogsViewsApisViewsV1SelectedFiltersToMap(model.Filters)
		if err != nil {
			return modelMap, err
		}
		modelMap["filters"] = []map[string]interface{}{filtersMap}
	}
	if model.FolderID != nil {
		modelMap["folder_id"] = model.FolderID.String()
	}
	return modelMap, nil
}

func DataSourceIbmLogsViewsApisViewsV1SearchQueryToMap(model *logsv0.ApisViewsV1SearchQuery) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["query"] = *model.Query
	return modelMap, nil
}

func DataSourceIbmLogsViewsApisViewsV1TimeSelectionToMap(model logsv0.ApisViewsV1TimeSelectionIntf) (map[string]interface{}, error) {
	if _, ok := model.(*logsv0.ApisViewsV1TimeSelectionSelectionTypeQuickSelection); ok {
		return DataSourceIbmLogsViewsApisViewsV1TimeSelectionSelectionTypeQuickSelectionToMap(model.(*logsv0.ApisViewsV1TimeSelectionSelectionTypeQuickSelection))
	} else if _, ok := model.(*logsv0.ApisViewsV1TimeSelectionSelectionTypeCustomSelection); ok {
		return DataSourceIbmLogsViewsApisViewsV1TimeSelectionSelectionTypeCustomSelectionToMap(model.(*logsv0.ApisViewsV1TimeSelectionSelectionTypeCustomSelection))
	} else if _, ok := model.(*logsv0.ApisViewsV1TimeSelection); ok {
		modelMap := make(map[string]interface{})
		model := model.(*logsv0.ApisViewsV1TimeSelection)
		if model.QuickSelection != nil {
			quickSelectionMap, err := DataSourceIbmLogsViewsApisViewsV1QuickTimeSelectionToMap(model.QuickSelection)
			if err != nil {
				return modelMap, err
			}
			modelMap["quick_selection"] = []map[string]interface{}{quickSelectionMap}
		}
		if model.CustomSelection != nil {
			customSelectionMap, err := DataSourceIbmLogsViewsApisViewsV1CustomTimeSelectionToMap(model.CustomSelection)
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

func DataSourceIbmLogsViewsApisViewsV1QuickTimeSelectionToMap(model *logsv0.ApisViewsV1QuickTimeSelection) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["caption"] = *model.Caption
	modelMap["seconds"] = flex.IntValue(model.Seconds)
	return modelMap, nil
}

func DataSourceIbmLogsViewsApisViewsV1CustomTimeSelectionToMap(model *logsv0.ApisViewsV1CustomTimeSelection) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["from_time"] = model.FromTime.String()
	modelMap["to_time"] = model.ToTime.String()
	return modelMap, nil
}

func DataSourceIbmLogsViewsApisViewsV1TimeSelectionSelectionTypeQuickSelectionToMap(model *logsv0.ApisViewsV1TimeSelectionSelectionTypeQuickSelection) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.QuickSelection != nil {
		quickSelectionMap, err := DataSourceIbmLogsViewsApisViewsV1QuickTimeSelectionToMap(model.QuickSelection)
		if err != nil {
			return modelMap, err
		}
		modelMap["quick_selection"] = []map[string]interface{}{quickSelectionMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsViewsApisViewsV1TimeSelectionSelectionTypeCustomSelectionToMap(model *logsv0.ApisViewsV1TimeSelectionSelectionTypeCustomSelection) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CustomSelection != nil {
		customSelectionMap, err := DataSourceIbmLogsViewsApisViewsV1CustomTimeSelectionToMap(model.CustomSelection)
		if err != nil {
			return modelMap, err
		}
		modelMap["custom_selection"] = []map[string]interface{}{customSelectionMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsViewsApisViewsV1SelectedFiltersToMap(model *logsv0.ApisViewsV1SelectedFilters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Filters != nil {
		filters := []map[string]interface{}{}
		for _, filtersItem := range model.Filters {
			filtersItemMap, err := DataSourceIbmLogsViewsApisViewsV1FilterToMap(&filtersItem)
			if err != nil {
				return modelMap, err
			}
			filters = append(filters, filtersItemMap)
		}
		modelMap["filters"] = filters
	}
	return modelMap, nil
}

func DataSourceIbmLogsViewsApisViewsV1FilterToMap(model *logsv0.ApisViewsV1Filter) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = *model.Name
	selectedValues := make(map[string]interface{})
	for k, v := range model.SelectedValues {
		selectedValues[k] = flex.Stringify(v)
	}
	modelMap["selected_values"] = selectedValues
	return modelMap, nil
}
