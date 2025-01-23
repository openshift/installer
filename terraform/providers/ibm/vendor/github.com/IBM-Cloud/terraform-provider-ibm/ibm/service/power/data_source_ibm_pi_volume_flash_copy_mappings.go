// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPIVolumeFlashCopyMappings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIVolumeFlashCopyMappings,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_VolumeID: {
				Description:  "The ID of the volume for which you want to retrieve detailed information.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_FlashCopyMappings: {
				Computed:    true,
				Description: "List of flash copy mappings details of a volume.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_CopyRate: {
							Computed:    true,
							Description: "The rate of flash copy operation of a volume.",
							Type:        schema.TypeInt,
						},
						Attr_FlashCopyName: {
							Computed:    true,
							Description: "The flash copy name of the volume.",
							Type:        schema.TypeString,
						},
						Attr_Progress: {
							Computed:    true,
							Description: "The progress of flash copy operation.",
							Type:        schema.TypeInt,
						},
						Attr_SourceVolumeName: {
							Computed:    true,
							Description: "The name of the source volume.",
							Type:        schema.TypeString,
						},
						Attr_StartTime: {
							Computed:    true,
							Description: "The start time of flash copy operation.",
							Type:        schema.TypeString,
						},
						Attr_Status: {
							Computed:    true,
							Description: "The copy status of a volume.",
							Type:        schema.TypeString,
						},
						Attr_TargetVolumeName: {
							Computed:    true,
							Description: "The name of the target volume.",
							Type:        schema.TypeString,
						},
					},
				},
				Type: schema.TypeList,
			},
		},
	}
}

func dataSourceIBMPIVolumeFlashCopyMappings(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	volClient := instance.NewIBMPIVolumeClient(ctx, sess, cloudInstanceID)
	volData, err := volClient.GetVolumeFlashCopyMappings(d.Get(Arg_VolumeID).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	results := make([]map[string]interface{}, 0, len(volData))
	for _, i := range volData {
		if i != nil {
			l := map[string]interface{}{
				Attr_CopyRate:         i.CopyRate,
				Attr_Progress:         i.Progress,
				Attr_SourceVolumeName: i.SourceVolumeName,
				Attr_StartTime:        i.StartTime.String(),
				Attr_Status:           i.Status,
				Attr_TargetVolumeName: i.TargetVolumeName,
			}
			if i.FlashCopyName != nil {
				l[Attr_FlashCopyName] = i.FlashCopyName
			}
			results = append(results, l)
		}
	}
	var clientgenU, _ = uuid.GenerateUUID()
	d.SetId(clientgenU)
	d.Set(Attr_FlashCopyMappings, results)

	return nil
}
