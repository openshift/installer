// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPIVolumeOnboarding() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIVolumeOnboardingReads,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_VolumeOnboardingID: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The ID of volume onboarding for which you want to retrieve detailed information.",
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			// TODO: Relabel this one "creation_date" to match literally every single other one
			Attr_CreateTime: {
				Computed:    true,
				Description: "The create-time of volume onboarding operation.",
				Type:        schema.TypeString,
			},
			Attr_Description: {
				Computed:    true,
				Description: "The description of the volume onboarding operation.",
				Type:        schema.TypeString,
			},
			Attr_InputVolumes: {
				Computed:    true,
				Description: "List of volumes requested to be onboarded.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Type:        schema.TypeList,
			},
			Attr_Progress: {
				Computed:    true,
				Description: "The progress of volume onboarding operation.",
				Type:        schema.TypeInt,
			},
			Attr_ResultsOnboardedVolumes: {
				Computed:    true,
				Description: "List of volumes which are onboarded successfully.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Type:        schema.TypeList,
			},
			Attr_ResultsVolumeOnboardingFailures: {
				Computed:    true,
				Description: "The volume onboarding failure details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_FailureMessage: {
							Computed:    true,
							Description: "The failure reason for the volumes which have failed to be onboarded.",
							Type:        schema.TypeString,
						},
						Attr_Volumes: {
							Computed:    true,
							Description: "List of volumes which have failed to be onboarded.",
							Elem:        &schema.Schema{Type: schema.TypeString},
							Type:        schema.TypeList,
						},
					},
				},
				Type: schema.TypeList,
			},
			Attr_Status: {
				Computed:    true,
				Description: "The status of volume onboarding operation.",
				Type:        schema.TypeString,
			},
		},
	}
}

func dataSourceIBMPIVolumeOnboardingReads(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	volOnboardClient := instance.NewIBMPIVolumeOnboardingClient(ctx, sess, cloudInstanceID)
	volOnboarding, err := volOnboardClient.Get(d.Get(Arg_VolumeOnboardingID).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*volOnboarding.ID)
	d.Set(Attr_CreateTime, volOnboarding.CreationTimestamp.String())
	d.Set(Attr_Description, volOnboarding.Description)
	d.Set(Attr_InputVolumes, volOnboarding.InputVolumes)
	d.Set(Attr_Progress, volOnboarding.Progress)
	d.Set(Attr_Status, volOnboarding.Status)
	d.Set(Attr_ResultsOnboardedVolumes, volOnboarding.Results.OnboardedVolumes)
	d.Set(Attr_ResultsVolumeOnboardingFailures, flattenVolumeOnboardingFailures(volOnboarding.Results.VolumeOnboardingFailures))

	return nil
}

func flattenVolumeOnboardingFailures(list []*models.VolumeOnboardingFailure) (result []map[string]interface{}) {
	if list != nil {
		result := make([]map[string]interface{}, len(list))
		for i, data := range list {
			l := map[string]interface{}{
				Attr_FailureMessage: data.FailureMessage,
				Attr_Volumes:        data.Volumes,
			}
			result[i] = l
		}
		return result
	}
	return
}
