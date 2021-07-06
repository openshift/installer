// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"net/url"
	"reflect"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/IBM/appconfiguration-go-admin-sdk/appconfigurationv1"
)

func dataSourceIbmAppConfigEnvironments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIbmAppConfigEnvironmentsRead,

		Schema: map[string]*schema.Schema{
			"guid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "GUID of the App Configuration service. Get it from the service instance credentials section of the dashboard.",
			},
			"tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "filter the resources to be returned based on the associated tags. Returns resources associated with any of the specified tags.",
			},
			"expand": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If set to `true`, returns expanded view of the resource details.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The number of records to retrieve. By default, the list operation return the first 10 records. To retrieve different set of records, use `limit` with `offset` to page through the available records.",
			},
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The number of records to skip. By specifying `offset`, you retrieve a subset of items that starts with the `offset` value. Use `offset` with `limit` to page through the available records.",
			},
			"environments": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Array of environments.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Environment name.",
						},
						"environment_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Environment id.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Environment description.",
						},
						"tags": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tags associated with the environment.",
						},
						"color_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Color code to distinguish the environment. The Hex code for the color. For example `#FF0000` for `red`.",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of the environment.",
						},
						"updated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last modified time of the environment data.",
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Environment URL.",
						},
					},
				},
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of records.",
			},
			"next": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "URL to navigate to the next list of records.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URL of the response.",
						},
					},
				},
			},
			"first": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "URL to navigate to the first page of records.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URL of the response.",
						},
					},
				},
			},
			"previous": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "URL to navigate to the previous list of records.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URL of the response.",
						},
					},
				},
			},
			"last": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "URL to navigate to the last page of records.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URL of the response.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmAppConfigEnvironmentsRead(d *schema.ResourceData, meta interface{}) error {
	guid := d.Get("guid").(string)

	appconfigClient, err := getAppConfigClient(meta, guid)
	if err != nil {
		return err
	}

	options := &appconfigurationv1.ListEnvironmentsOptions{}

	if _, ok := d.GetOk("expand"); ok {
		options.SetExpand(d.Get("expand").(bool))
	}

	if _, ok := d.GetOk("tags"); ok {
		options.SetTags(d.Get("tags").(string))
	}

	var environmentList *appconfigurationv1.EnvironmentList
	var offset int64 = 0
	var limit int64 = 10
	finalList := []appconfigurationv1.Environment{}

	var isLimit bool
	if _, ok := d.GetOk("limit"); ok {
		isLimit = true
		limit = int64(d.Get("limit").(int))
	}
	options.SetLimit(limit)

	if _, ok := d.GetOk("offset"); ok {
		offset = int64(d.Get("offset").(int))
	}
	for {
		options.SetOffset(offset)
		result, response, err := appconfigClient.ListEnvironments(options)
		environmentList = result
		if err != nil {
			log.Printf("[DEBUG] ListEnvironments failed %s\n%s", err, response)
			return err
		}
		if isLimit {
			offset = 0
		} else {
			offset = dataSourceEnvironmentListGetNext(result.Next)
		}
		finalList = append(finalList, result.Environments...)
		if offset == 0 {
			break
		}
	}

	environmentList.Environments = finalList

	d.SetId(guid)

	if environmentList.Environments != nil {
		err = d.Set("environments", dataSourceEnvironmentListFlattenEnvironments(environmentList.Environments))
		if err != nil {
			return fmt.Errorf("error setting environments %s", err)
		}
	}
	if environmentList.TotalCount != nil {
		if err = d.Set("total_count", environmentList.TotalCount); err != nil {
			return fmt.Errorf("error setting total_count: %s", err)
		}
	}
	if environmentList.Limit != nil {
		if err = d.Set("limit", environmentList.Limit); err != nil {
			return fmt.Errorf("error setting limit: %s", err)
		}
	}
	if environmentList.Offset != nil {
		if err = d.Set("offset", environmentList.Offset); err != nil {
			return fmt.Errorf("error setting offset: %s", err)
		}
	}
	if environmentList.First != nil {
		err = d.Set("first", dataSourceEnvironmentListFlattenPagination(*environmentList.First))
		if err != nil {
			return fmt.Errorf("error setting first %s", err)
		}
	}

	if environmentList.Previous != nil {
		err = d.Set("previous", dataSourceEnvironmentListFlattenPagination(*environmentList.Previous))
		if err != nil {
			return fmt.Errorf("error setting previous %s", err)
		}
	}

	if environmentList.Last != nil {
		err = d.Set("last", dataSourceEnvironmentListFlattenPagination(*environmentList.Last))
		if err != nil {
			return fmt.Errorf("error setting last %s", err)
		}
	}
	if environmentList.Next != nil {
		err = d.Set("next", dataSourceEnvironmentListFlattenPagination(*environmentList.Next))
		if err != nil {
			return fmt.Errorf("error setting next %s", err)
		}
	}

	return nil
}

func dataSourceEnvironmentListFlattenEnvironments(result []appconfigurationv1.Environment) (environments []map[string]interface{}) {
	for _, environmentsItem := range result {
		environments = append(environments, dataSourceEnvironmentListEnvironmentsToMap(environmentsItem))
	}

	return environments
}

func dataSourceEnvironmentListEnvironmentsToMap(environmentsItem appconfigurationv1.Environment) (environmentsMap map[string]interface{}) {
	environmentsMap = map[string]interface{}{}

	if environmentsItem.Name != nil {
		environmentsMap["name"] = environmentsItem.Name
	}
	if environmentsItem.EnvironmentID != nil {
		environmentsMap["environment_id"] = environmentsItem.EnvironmentID
	}
	if environmentsItem.Description != nil {
		environmentsMap["description"] = environmentsItem.Description
	}
	if environmentsItem.Tags != nil {
		environmentsMap["tags"] = environmentsItem.Tags
	}
	if environmentsItem.ColorCode != nil {
		environmentsMap["color_code"] = environmentsItem.ColorCode
	}
	if environmentsItem.CreatedTime != nil {
		environmentsMap["created_time"] = environmentsItem.CreatedTime.String()
	}
	if environmentsItem.UpdatedTime != nil {
		environmentsMap["updated_time"] = environmentsItem.UpdatedTime.String()
	}
	if environmentsItem.Href != nil {
		environmentsMap["href"] = environmentsItem.Href
	}

	return environmentsMap
}

func dataSourceEnvironmentListGetNext(next interface{}) int64 {
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

func dataSourceEnvironmentListFlattenPagination(result appconfigurationv1.PageHrefResponse) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceEnvironmentListURLToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceEnvironmentListURLToMap(urlItem appconfigurationv1.PageHrefResponse) (urlMap map[string]interface{}) {
	urlMap = map[string]interface{}{}

	if urlItem.Href != nil {
		urlMap["href"] = urlItem.Href
	}

	return urlMap
}
