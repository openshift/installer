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
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_en_smtp_configurations", "list")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
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

		result, _, err := enClient.ListSMTPConfigurationsWithContext(context, options)

		smtpConfigurationList = result

		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListSMTPConfigurationsWithContext failed: %s", err.Error()), "(Data) ibm_en_smtp_configurations", "list")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
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
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting total_count: %s", err), "(Data) ibm_en_smtp_configurations", "list")
		return tfErr.GetDiag()
	}

	if smtpConfigurationList.SMTPConfigurations != nil {
		if err = d.Set("smtp_configurations", enFlattenSMTPConfigurationList(smtpConfigurationList.SMTPConfigurations)); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting SMTPConfigurations: %s", err), "(Data) ibm_en_smtp_configurations", "list")
			return tfErr.GetDiag()
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
