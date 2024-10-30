// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"log"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPIVolumeOnboardings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIVolumeOnboardingsReads,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_Onboardings: {
				Computed:    true,
				Description: "List of volume onboardings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_Description: {
							Computed:    true,
							Description: "The description of the volume onboarding operation.",
							Type:        schema.TypeString,
						},
						Attr_ID: {
							Computed:    true,
							Description: "The type of cycling mode used.",
							Type:        schema.TypeString,
						},
						Attr_InputVolumes: {
							Computed:    true,
							Description: "List of volumes requested to be onboarded.",
							Elem:        &schema.Schema{Type: schema.TypeString},
							Type:        schema.TypeList,
						},
						Attr_Status: {
							Computed:    true,
							Description: "The status of volume onboarding operation.",
							Type:        schema.TypeString,
						},
					},
				},
				Type: schema.TypeList,
			},
		},
	}
}

func dataSourceIBMPIVolumeOnboardingsReads(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	volOnboardClient := instance.NewIBMPIVolumeOnboardingClient(ctx, sess, cloudInstanceID)
	volOnboardings, err := volOnboardClient.GetAll()
	if err != nil {
		return diag.FromErr(err)
	}

	var clientgenU, _ = uuid.GenerateUUID()
	d.SetId(clientgenU)
	d.Set(Attr_Onboardings, flattenVolumeOnboardings(volOnboardings.Onboardings))

	return nil
}

func flattenVolumeOnboardings(list []*models.VolumeOnboardingCommon) (networks []map[string]interface{}) {
	log.Printf("Calling the flattenVolumeOnboardings call with list %d", len(list))
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		l := map[string]interface{}{
			Attr_ID:           *i.ID,
			Attr_Description:  i.Description,
			Attr_InputVolumes: i.InputVolumes,
			Attr_Status:       i.Status,
		}
		result = append(result, l)
	}
	return result
}
