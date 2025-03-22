// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMPIVolumeOnboarding() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPIVolumeOnboardingCreate,
		ReadContext:   resourceIBMPIVolumeOnboardingRead,
		DeleteContext: resourceIBMPIVolumeOnboardingDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description: "The GUID of the service instance associated with an account.",
				ForceNew:    true,
				Required:    true,
				Type:        schema.TypeString,
			},
			Arg_Description: {
				Computed:    true,
				Description: "Description of the volume onboarding operation",
				Optional:    true,
				Type:        schema.TypeString,
			},
			Arg_OnboardingVolumes: {
				Description: "List of onboarding volumes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Arg_AuxiliaryVolumes: {
							Description: "List auxiliary volumes.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									Arg_AuxiliaryVolumeName: {
										Description: "The auxiliary volume name.",
										Required:    true,
										Type:        schema.TypeString,
									},
									Arg_DisplayName: {
										Description: "The display name of auxiliary volume which is to be onboarded.",
										Optional:    true,
										Type:        schema.TypeString,
									},
								},
							},
							MinItems: 1,
							Optional: true,
							Type:     schema.TypeList,
						},
						Arg_SourceCRN: {
							Description: "The crn of source service broker instance from where auxiliary volumes need to be onboarded.",
							Required:    true,
							Type:        schema.TypeString,
						},
					},
				},
				ForceNew: true,
				MinItems: 1,
				Required: true,
				Type:     schema.TypeList,
			},

			// Attributes
			Attr_CreateTime: {
				Computed:    true,
				Description: "The create time of volume onboarding operation.",
				Type:        schema.TypeString,
			},
			Attr_InputVolumes: {
				Computed:    true,
				Description: "List of volumes requested to be onboarded.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Type:        schema.TypeList,
			},
			Attr_OnboardingID: {
				Computed:    true,
				Description: "The volume onboarding ID.",
				Type:        schema.TypeString,
			},
			Attr_Progress: {
				Computed:    true,
				Description: "The progress of volume onboarding operation.",
				Type:        schema.TypeFloat,
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
							Description: "The failure reason for the volumes which have failed to be onboarded",
							Type:        schema.TypeString,
						},
						Attr_Volumes: {
							Computed:    true,
							Description: "List of volumes which have failed to be onboarded",
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

func resourceIBMPIVolumeOnboardingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	client := instance.NewIBMPIVolumeOnboardingClient(ctx, sess, cloudInstanceID)

	vol, err := expandCreateVolumeOnboarding(d.Get(Arg_OnboardingVolumes).([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	body := &models.VolumeOnboardingCreate{
		Volumes: vol,
	}

	if v, ok := d.GetOk(Arg_Description); ok {
		body.Description = v.(string)
	}

	resOnboarding, err := client.CreateVolumeOnboarding(body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, resOnboarding.ID))

	return resourceIBMPIVolumeOnboardingRead(ctx, d, meta)
}

func resourceIBMPIVolumeOnboardingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, onboardingID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := instance.NewIBMPIVolumeOnboardingClient(ctx, sess, cloudInstanceID)

	onboardingData, err := client.Get(onboardingID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set(Arg_Description, onboardingData.Description)
	d.Set(Attr_CreateTime, onboardingData.CreationTimestamp.String())
	d.Set(Attr_InputVolumes, onboardingData.InputVolumes)
	d.Set(Attr_OnboardingID, *onboardingData.ID)
	d.Set(Attr_Progress, onboardingData.Progress)
	d.Set(Attr_ResultsOnboardedVolumes, onboardingData.Results.OnboardedVolumes)
	d.Set(Attr_ResultsVolumeOnboardingFailures, flattenVolumeOnboardingFailures(onboardingData.Results.VolumeOnboardingFailures))
	d.Set(Attr_Status, onboardingData.Status)
	return nil
}

func resourceIBMPIVolumeOnboardingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// There is no delete or unset concept for instance action
	d.SetId("")
	return nil
}

// expandCreateVolumeOnboarding expands create volume onboarding resource
func expandCreateVolumeOnboarding(data []interface{}) ([]*models.AuxiliaryVolumesForOnboarding, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("[ERROR] no pi_onboarding_volumes received")
	}

	auxVolForOnboarding := make([]*models.AuxiliaryVolumesForOnboarding, 0)

	for _, d := range data {
		resource := d.(map[string]interface{})

		var crn string
		auxVolumes := make([]interface{}, 0)

		if v, ok := resource["pi_source_crn"]; ok && v != "" {
			crn = resource["pi_source_crn"].(string)
		}

		if v, ok := resource["pi_auxiliary_volumes"]; ok && len(v.([]interface{})) != 0 {
			auxVolumes = resource["pi_auxiliary_volumes"].([]interface{})
		}

		auxVolForOnboarding = append(auxVolForOnboarding, &models.AuxiliaryVolumesForOnboarding{
			SourceCRN:        &crn,
			AuxiliaryVolumes: expandAuxiliaryVolumeForOnboarding(auxVolumes),
		})

	}

	return auxVolForOnboarding, nil
}

func expandAuxiliaryVolumeForOnboarding(data []interface{}) []*models.AuxiliaryVolumeForOnboarding {
	auxVolumeForOnboarding := make([]*models.AuxiliaryVolumeForOnboarding, 0)

	for _, d := range data {
		var auxVolumeName, displayName string
		resource := d.(map[string]interface{})

		if v, ok := resource[Arg_AuxiliaryVolumeName]; ok && v != "" {
			auxVolumeName = resource[Arg_AuxiliaryVolumeName].(string)
		}

		if v, ok := resource[Arg_DisplayName]; ok && v != "" {
			displayName = resource[Arg_DisplayName].(string)
		}

		auxVolumeForOnboarding = append(auxVolumeForOnboarding, &models.AuxiliaryVolumeForOnboarding{
			AuxVolumeName: &auxVolumeName,
			Name:          displayName,
		})
	}

	return auxVolumeForOnboarding
}
