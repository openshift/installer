// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"log"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
)

func DataSourceIBMPIVolumeOnboardings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIVolumeOnboardingsReads,
		Schema: map[string]*schema.Schema{
			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			// Computed Attributes
			"onboardings": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of volume onboardings",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the volume onboarding operation",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates the volume onboarding operation id",
						},
						"input_volumes": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of volumes requested to be onboarded",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates the status of volume onboarding operation",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPIVolumeOnboardingsReads(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)
	volOnboardClient := instance.NewIBMPIVolumeOnboardingClient(ctx, sess, cloudInstanceID)
	volOnboardings, err := volOnboardClient.GetAll()
	if err != nil {
		return diag.FromErr(err)
	}

	var clientgenU, _ = uuid.GenerateUUID()
	d.SetId(clientgenU)
	d.Set("onboardings", flattenVolumeOnboardings(volOnboardings.Onboardings))

	return nil
}

func flattenVolumeOnboardings(list []*models.VolumeOnboardingCommon) (networks []map[string]interface{}) {
	log.Printf("Calling the flattenVolumeOnboardings call with list %d", len(list))
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		l := map[string]interface{}{
			"id":            *i.ID,
			"description":   i.Description,
			"input_volumes": i.InputVolumes,
			"status":        i.Status,
		}

		result = append(result, l)
	}

	return result
}
