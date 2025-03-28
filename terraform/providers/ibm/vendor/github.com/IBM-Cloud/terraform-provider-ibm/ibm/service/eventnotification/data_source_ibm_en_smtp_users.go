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

func DataSourceIBMEnSMTPUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMEnSMTPUsersRead,

		Schema: map[string]*schema.Schema{
			"instance_guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier for IBM Cloud Event Notifications instance.",
			},
			"smtp_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier for SMTP Configuration.",
			},
			"search_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter the user by name or type.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of smtp users.",
			},
			"smtp_users": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of smtp users.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "UserID.",
						},
						"username": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SMTP user name.",
						},
						"smtp_config_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SMTP confg Id.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SMTP User description.",
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
						"created_at": {
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

func dataSourceIBMEnSMTPUsersRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_en_smtp_users", "list")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	options := &en.ListSMTPUsersOptions{}

	options.SetInstanceID(d.Get("instance_guid").(string))
	options.SetID(d.Get("smtp_id").(string))

	if _, ok := d.GetOk("search_key"); ok {
		options.SetSearch(d.Get("search_key").(string))
	}
	var smtpuserList *en.SMTPUsersList

	finalList := []en.SMTPUser{}

	var offset int64 = 0
	var limit int64 = 100

	options.SetLimit(limit)

	for {
		options.SetOffset(offset)

		result, _, err := enClient.ListSMTPUsersWithContext(context, options)

		smtpuserList = result

		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListSMTPUsersWithContext failed: %s", err.Error()), "(Data) ibm_en_smtp_users", "list")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}

		offset = offset + limit

		finalList = append(finalList, result.Users...)

		if offset > *result.TotalCount {
			break
		}
	}

	smtpuserList.Users = finalList

	d.SetId(fmt.Sprintf("SMTPConfigurations/%s", *options.InstanceID))

	if err = d.Set("total_count", flex.IntValue(smtpuserList.TotalCount)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting total_count: %s", err), "(Data) ibm_en_smtp_users", "list")
		return tfErr.GetDiag()
	}

	if smtpuserList.Users != nil {
		if err = d.Set("smtp_users", enFlattenSMTPUsersList(smtpuserList.Users)); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("SMTPUsers: %s", err), "(Data) ibm_en_smtp_users", "list")
			return tfErr.GetDiag()
		}
	}

	return nil
}

func enFlattenSMTPUsersList(result []en.SMTPUser) (smtpusers []map[string]interface{}) {
	for _, smtpuserItem := range result {
		smtpusers = append(smtpusers, enSMTPUserListToMap(smtpuserItem))
	}

	return smtpusers
}

func enSMTPUserListToMap(smtpuserItem en.SMTPUser) (smtpuser map[string]interface{}) {
	smtpuser = map[string]interface{}{}

	if smtpuserItem.ID != nil {
		smtpuser["id"] = smtpuserItem.ID
	}
	if smtpuserItem.Username != nil {
		smtpuser["username"] = smtpuserItem.Username
	}
	if smtpuserItem.SMTPConfigID != nil {
		smtpuser["smtp_config_id"] = smtpuserItem.SMTPConfigID
	}
	if smtpuserItem.Description != nil {
		smtpuser["description"] = smtpuserItem.Description
	}
	if smtpuserItem.Domain != nil {
		smtpuser["domain"] = smtpuserItem.Domain
	}
	if smtpuserItem.UpdatedAt != nil {
		smtpuser["updated_at"] = flex.DateTimeToString(smtpuserItem.UpdatedAt)
	}
	if smtpuserItem.CreatedAt != nil {
		smtpuser["created_at"] = flex.DateTimeToString(smtpuserItem.CreatedAt)
	}

	return smtpuser
}
