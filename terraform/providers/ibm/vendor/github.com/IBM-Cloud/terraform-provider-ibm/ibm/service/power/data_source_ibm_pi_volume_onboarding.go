// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
)

func DataSourceIBMPIVolumeOnboarding() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIVolumeOnboardingReads,
		Schema: map[string]*schema.Schema{
			PIVolumeOnboardingID: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Volume onboarding ID",
				ValidateFunc: validation.NoZeroValues,
			},
			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			// Computed Attributes
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the create time of volume onboarding operation",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the volume onboarding operation",
			},
			"input_volumes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of volumes requested to be onboarded",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"progress": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Indicates the progress of volume onboarding operation",
			},
			"results_onboarded_volumes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of volumes which are onboarded successfully",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"results_volume_onboarding_failures": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"failure_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The failure reason for the volumes which have failed to be onboarded",
						},
						"volumes": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of volumes which have failed to be onboarded",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					}},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the status of volume onboarding operation",
			},
		},
	}
}

func dataSourceIBMPIVolumeOnboardingReads(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)
	volOnboardClient := instance.NewIBMPIVolumeOnboardingClient(ctx, sess, cloudInstanceID)
	volOnboarding, err := volOnboardClient.Get(d.Get(PIVolumeOnboardingID).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*volOnboarding.ID)
	d.Set("create_time", volOnboarding.CreationTimestamp.String())
	d.Set("description", volOnboarding.Description)
	d.Set("input_volumes", volOnboarding.InputVolumes)
	d.Set("progress", volOnboarding.Progress)
	d.Set("status", volOnboarding.Status)
	d.Set("results_onboarded_volumes", volOnboarding.Results.OnboardedVolumes)
	d.Set("results_volume_onboarding_failures", flattenVolumeOnboardingFailures(volOnboarding.Results.VolumeOnboardingFailures))

	return nil
}

func flattenVolumeOnboardingFailures(list []*models.VolumeOnboardingFailure) (result []map[string]interface{}) {
	if list != nil {
		result := make([]map[string]interface{}, len(list))
		for i, data := range list {
			l := map[string]interface{}{
				"failure_message": data.FailureMessage,
				"volumes":         data.Volumes,
			}
			result[i] = l
		}
		return result
	}
	return
}
