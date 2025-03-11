// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package eventnotification

import (
	"context"
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	en "github.com/IBM/event-notifications-go-admin-sdk/eventnotificationsv1"
)

func DataSourceIBMEnTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMEnEmailTemplatesRead,

		Schema: map[string]*schema.Schema{
			"instance_guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier for IBM Cloud Event Notifications instance.",
			},
			"search_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter the template by name or type.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of templates.",
			},
			"templates": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of templates.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Template ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Template name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "description of the template.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "template type smtp_custom.notification/smtp_custom.invitation.",
						},
						"subscription_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Subscription count.",
						},
						"subscription_names": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of subscriptions.",
							Elem:        &schema.Schema{Type: schema.TypeString},
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

func dataSourceIBMEnEmailTemplatesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		return diag.FromErr(err)
	}

	options := &en.ListTemplatesOptions{}

	options.SetInstanceID(d.Get("instance_guid").(string))

	if _, ok := d.GetOk("search_key"); ok {
		options.SetSearch(d.Get("search_key").(string))
	}
	var templateList *en.TemplateList

	finalList := []en.Template{}

	var offset int64 = 0
	var limit int64 = 100

	options.SetLimit(limit)

	for {
		options.SetOffset(offset)

		result, response, err := enClient.ListTemplatesWithContext(context, options)

		templateList = result

		if err != nil {
			return diag.FromErr(fmt.Errorf("ListTemplatesWithContext failed %s\n%s", err, response))
		}

		offset = offset + limit

		finalList = append(finalList, result.Templates...)

		if offset > *result.TotalCount {
			break
		}
	}

	templateList.Templates = finalList

	d.SetId(fmt.Sprintf("Templates/%s", *options.InstanceID))

	if err = d.Set("total_count", flex.IntValue(templateList.TotalCount)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting total_count: %s", err))
	}

	if templateList.Templates != nil {
		if err = d.Set("templates", enFlattentemplatesList(templateList.Templates)); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting Templates %s", err))
		}
	}

	return nil
}

func enFlattentemplatesList(result []en.Template) (templates []map[string]interface{}) {
	for _, templateItem := range result {
		templates = append(templates, enTemplateListToMap(templateItem))
	}

	return templates
}

func enTemplateListToMap(templateItem en.Template) (template map[string]interface{}) {
	template = map[string]interface{}{}

	if templateItem.ID != nil {
		template["id"] = templateItem.ID
	}
	if templateItem.Name != nil {
		template["name"] = templateItem.Name
	}
	if templateItem.Description != nil {
		template["description"] = templateItem.Description
	}
	if templateItem.Type != nil {
		template["type"] = templateItem.Type
	}
	if templateItem.SubscriptionCount != nil {
		template["subscription_count"] = templateItem.SubscriptionCount
	}
	if templateItem.SubscriptionNames != nil {
		template["subscription_names"] = templateItem.SubscriptionNames
	}
	if templateItem.UpdatedAt != nil {
		template["updated_at"] = flex.DateTimeToString(templateItem.UpdatedAt)
	}

	return template
}
