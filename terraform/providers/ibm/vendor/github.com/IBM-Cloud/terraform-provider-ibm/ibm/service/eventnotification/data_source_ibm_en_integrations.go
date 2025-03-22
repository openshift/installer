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

func DataSourceIBMEnIntegrations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMEnIntegrationsRead,

		Schema: map[string]*schema.Schema{
			"instance_guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier for IBM Cloud Event Notifications instance.",
			},
			"search_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter the destinations by type.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of destinations.",
			},
			"integrations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of destinations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Destination ID.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of integrayion kms/hpcs",
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

func dataSourceIBMEnIntegrationsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_en_integrations", "list")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	options := &en.ListIntegrationsOptions{}

	options.SetInstanceID(d.Get("instance_guid").(string))

	if _, ok := d.GetOk("search_key"); ok {
		options.SetSearch(d.Get("search_key").(string))
	}
	var integrationList *en.IntegrationList

	finalList := []en.IntegrationListItem{}

	var offset int64 = 0
	var limit int64 = 100

	options.SetLimit(limit)

	for {
		options.SetOffset(offset)

		result, _, err := enClient.ListIntegrationsWithContext(context, options)

		integrationList = result

		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListIntegrationsWithContext failed: %s", err.Error()), "(Data) ibm_en_integrations", "list")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}

		offset = offset + limit

		finalList = append(finalList, result.Integrations...)

		if offset > *result.TotalCount {
			break
		}
	}

	integrationList.Integrations = finalList

	d.SetId(fmt.Sprintf("integrations/%s", *options.InstanceID))

	if err = d.Set("total_count", flex.IntValue(integrationList.TotalCount)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting total_count: %s", err), "(Data) ibm_en_integrations", "list")
		return tfErr.GetDiag()
	}

	if integrationList.Integrations != nil {
		if err = d.Set("integrations", enFlattenIntegrationsList(integrationList.Integrations)); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting integrations: %s", err), "(Data) ibm_en_integrations", "list")
			return tfErr.GetDiag()
		}
	}

	return nil
}

func enFlattenIntegrationsList(result []en.IntegrationListItem) (integrations []map[string]interface{}) {
	for _, integrationsItem := range result {
		integrations = append(integrations, enIntegrationsListToMap(integrationsItem))
	}

	return integrations
}

func enIntegrationsListToMap(integrationItem en.IntegrationListItem) (integration map[string]interface{}) {
	integration = map[string]interface{}{}

	if integrationItem.ID != nil {
		integration["id"] = integrationItem.ID
	}
	if integrationItem.Type != nil {
		integration["type"] = integrationItem.Type
	}
	if integrationItem.UpdatedAt != nil {
		integration["updated_at"] = flex.DateTimeToString(integrationItem.UpdatedAt)
	}

	return integration
}
