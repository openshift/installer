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

func DataSourceIBMEnSMTPCOnfigurations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMEnSMTPConfigurationsRead,

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
				Description: "Total number of smtp configurations.",
			},
			"smtp_configurations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of smtp configurations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SMTP Configuration ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SMTP Configuration name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "description of the SMTP Configuration.",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain .",
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

func dataSourceIBMEnSMTPConfigurationsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		return diag.FromErr(err)
	}

	options := &en.ListSMTPConfigurationsOptions{}

	options.SetInstanceID(d.Get("instance_guid").(string))

	if _, ok := d.GetOk("search_key"); ok {
		options.SetSearch(d.Get("search_key").(string))
	}
	var smtpConfigurationList *en.SMTPConfigurationsList

	finalList := []en.SMTPConfiguration{}

	var offset int64 = 0
	var limit int64 = 100

	options.SetLimit(limit)

	for {
		options.SetOffset(offset)

		result, response, err := enClient.ListSMTPConfigurationsWithContext(context, options)

		smtpConfigurationList = result

		if err != nil {
			return diag.FromErr(fmt.Errorf("ListSMTPConfigurationsWithContext failed %s\n%s", err, response))
		}

		offset = offset + limit

		finalList = append(finalList, result.SMTPConfigurations...)

		if offset > *result.TotalCount {
			break
		}
	}

	smtpConfigurationList.SMTPConfigurations = finalList

	d.SetId(fmt.Sprintf("SMTPConfigurations/%s", *options.InstanceID))

	if err = d.Set("total_count", flex.IntValue(smtpConfigurationList.TotalCount)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting total_count: %s", err))
	}

	if smtpConfigurationList.SMTPConfigurations != nil {
		if err = d.Set("smtp_configurations", enFlattenSMTPConfigurationList(smtpConfigurationList.SMTPConfigurations)); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting SMTPConfigurations %s", err))
		}
	}

	return nil
}

func enFlattenSMTPConfigurationList(result []en.SMTPConfiguration) (smtpconfigurations []map[string]interface{}) {
	for _, smtpconfigItem := range result {
		smtpconfigurations = append(smtpconfigurations, enSMTPConfigurationListToMap(smtpconfigItem))
	}

	return smtpconfigurations
}

func enSMTPConfigurationListToMap(smtpconfigItem en.SMTPConfiguration) (smtpconfiguration map[string]interface{}) {
	smtpconfiguration = map[string]interface{}{}

	if smtpconfigItem.ID != nil {
		smtpconfiguration["id"] = smtpconfigItem.ID
	}
	if smtpconfigItem.Name != nil {
		smtpconfiguration["name"] = smtpconfigItem.Name
	}
	if smtpconfigItem.Description != nil {
		smtpconfiguration["description"] = smtpconfigItem.Description
	}
	if smtpconfigItem.Domain != nil {
		smtpconfiguration["domain"] = smtpconfigItem.Domain
	}
	if smtpconfigItem.UpdatedAt != nil {
		smtpconfiguration["updated_at"] = flex.DateTimeToString(smtpconfigItem.UpdatedAt)
	}

	return smtpconfiguration
}
