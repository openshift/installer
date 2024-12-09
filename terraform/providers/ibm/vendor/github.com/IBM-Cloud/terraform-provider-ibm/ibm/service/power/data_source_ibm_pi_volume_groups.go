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
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
)

func DataSourceIBMPIVolumeGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIVolumeGroupsRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_VolumeGroups: {
				Computed:    true,
				Description: "List of all volume groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_ConsistencyGroupName: {
							Computed:    true,
							Description: "The name of consistency group at storage controller level.",
							Type:        schema.TypeString,
						},
						Attr_ID: {
							Computed:    true,
							Description: "The unique identifier of the volume group.",
							Type:        schema.TypeString,
						},
						Attr_ReplicationStatus: {
							Computed:    true,
							Description: "The replication status of volume group.",
							Type:        schema.TypeString,
						},
						Attr_Status: {
							Computed:    true,
							Description: "The status of the volume group.",
							Type:        schema.TypeString,
						},
						Attr_StatusDescriptionErrors: {
							Computed:    true,
							Description: "The status details of the volume group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									Attr_Key: {
										Computed:    true,
										Description: "The volume group error key.",
										Type:        schema.TypeString,
									},
									Attr_Message: {
										Computed:    true,
										Description: "The failure message providing more details about the error key.",
										Type:        schema.TypeString,
									},
									Attr_VolumeIDs: {
										Computed:    true,
										Description: "List of volume IDs, which failed to be added/removed to/from the volume group, with the given error.",
										Elem:        &schema.Schema{Type: schema.TypeString},
										Type:        schema.TypeList,
									},
								},
							},
							Type: schema.TypeSet,
						},
						Attr_VolumeGroupName: {
							Computed:    true,
							Description: "The name of the volume group.",
							Type:        schema.TypeString,
						},
					},
				},
				Type: schema.TypeList,
			},
		},
	}
}

func dataSourceIBMPIVolumeGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	vgClient := instance.NewIBMPIVolumeGroupClient(ctx, sess, cloudInstanceID)
	vgData, err := vgClient.GetAll()
	if err != nil {
		return diag.FromErr(err)
	}

	var clientgenU, _ = uuid.GenerateUUID()
	d.SetId(clientgenU)
	d.Set(Attr_VolumeGroups, flattenVolumeGroups(vgData.VolumeGroups))

	return nil
}

func flattenVolumeGroups(list []*models.VolumeGroup) []map[string]interface{} {
	log.Printf("Calling the flattenVolumeGroups call with list %d", len(list))
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		l := map[string]interface{}{
			Attr_ConsistencyGroupName:    i.ConsistencyGroupName,
			Attr_ID:                      *i.ID,
			Attr_ReplicationStatus:       i.ReplicationStatus,
			Attr_Status:                  i.Status,
			Attr_StatusDescriptionErrors: flattenVolumeGroupStatusDescription(i.StatusDescription.Errors),
			Attr_VolumeGroupName:         i.Name,
		}
		result = append(result, l)
	}
	return result
}
