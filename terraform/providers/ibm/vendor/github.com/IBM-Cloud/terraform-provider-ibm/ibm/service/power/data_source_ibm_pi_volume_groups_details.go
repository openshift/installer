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

func DataSourceIBMPIVolumeGroupsDetails() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIVolumeGroupsDetailsRead,
		Schema: map[string]*schema.Schema{
			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			// Computed Attributes
			"volume_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"replication_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"consistency_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status_description_errors": vgStatusDescriptionErrors(),
						"volume_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				}},
		},
	}
}

func dataSourceIBMPIVolumeGroupsDetailsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)
	vgClient := instance.NewIBMPIVolumeGroupClient(ctx, sess, cloudInstanceID)
	vgData, err := vgClient.GetAllDetails()
	if err != nil {
		return diag.FromErr(err)
	}

	var clientgenU, _ = uuid.GenerateUUID()
	d.SetId(clientgenU)
	d.Set("volume_groups", flattenVolumeGroupsDetails(vgData.VolumeGroups))

	return nil
}

func flattenVolumeGroupsDetails(list []*models.VolumeGroupDetails) []map[string]interface{} {
	log.Printf("Calling the flattenVolumeGroups call with list %d", len(list))
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		l := map[string]interface{}{
			"id":                        *i.ID,
			"replication_status":        i.ReplicationStatus,
			"consistency_group_name":    i.ConsistencyGroupName,
			"status":                    i.Status,
			"status_description_errors": flattenVolumeGroupStatusDescription(i.StatusDescription.Errors),
			"volume_group_name":         i.Name,
			"volume_ids":                i.VolumeIDs,
		}

		result = append(result, l)
	}

	return result
}
