// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package eventnotification

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	en "github.com/IBM/event-notifications-go-admin-sdk/eventnotificationsv1"
)

func DataSourceIBMEnSources() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMEnSourcesRead,

		Schema: map[string]*schema.Schema{
			"instance_guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier for IBM Cloud Event Notifications instance.",
			},
			"search_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter the source by name or type.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of sources.",
			},
			"sources": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of sources.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Source ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Source name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "description of the source.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Source type.",
						},
						"topic_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Topic count.",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Source is enabled or not.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Updated at.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMEnSourcesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_en_smtp_users", "list")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	options := &en.ListSourcesOptions{}

	options.SetInstanceID(d.Get("instance_guid").(string))

	if _, ok := d.GetOk("search_key"); ok {
		options.SetSearch(d.Get("search_key").(string))
	}
	var sourceList *en.SourceList

	finalList := []en.SourceListItem{}

	var offset int64 = 0
	var limit int64 = 100

	options.SetLimit(limit)

	for {
		options.SetOffset(offset)

		result, _, err := enClient.ListSourcesWithContext(context, options)

		sourceList = result

		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListSourcesWithContext failed: %s", err.Error()), "(Data) ibm_en_smtp_users", "list")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}

		offset = offset + limit

		finalList = append(finalList, result.Sources...)

		if offset > *result.TotalCount {
			break
		}
	}

	sourceList.Sources = finalList

	d.SetId(fmt.Sprintf("Sources/%s", *options.InstanceID))

	if err = d.Set("total_count", flex.IntValue(sourceList.TotalCount)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting total_count: %s", err), "(Data) ibm_en_sources", "list")
		return tfErr.GetDiag()
	}

	if sourceList.Sources != nil {
		if err = d.Set("sources", enFlattenSourcesList(sourceList.Sources)); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("SMTPUsers: %s", err), "(Data) ibm_en_sources", "list")
			return tfErr.GetDiag()
		}
	}

	return nil
}

func enFlattenSourcesList(result []en.SourceListItem) (sources []map[string]interface{}) {
	for _, sourcesItem := range result {
		sources = append(sources, enSourceListToMap(sourcesItem))
	}

	return sources
}

func enSourceListToMap(sourceItem en.SourceListItem) (source map[string]interface{}) {
	source = map[string]interface{}{}

	if sourceItem.ID != nil {
		source["id"] = sourceItem.ID
	}
	if sourceItem.Name != nil {
		source["name"] = sourceItem.Name
	}
	if sourceItem.Description != nil {
		source["description"] = sourceItem.Description
	}
	if sourceItem.Type != nil {
		source["type"] = sourceItem.Type
	}
	if sourceItem.TopicCount != nil {
		source["topic_count"] = sourceItem.TopicCount
	}
	if sourceItem.Enabled != nil {
		source["enabled"] = sourceItem.Enabled
	}
	if sourceItem.UpdatedAt != nil {
		source["updated_at"] = flex.DateTimeToString(sourceItem.UpdatedAt)
	}

	return source
}
