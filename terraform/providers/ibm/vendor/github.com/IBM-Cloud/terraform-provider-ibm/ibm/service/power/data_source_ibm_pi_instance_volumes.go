// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPIInstanceVolumes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIInstanceVolumesRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_InstanceName: {
				Description:  "The unique identifier or name of the instance.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attribute
			Attr_BootVolumeID: {
				Computed:    true,
				Description: "The unique identifier of the boot volume.",
				Type:        schema.TypeString,
			},
			Attr_InstanceVolumes: {
				Computed:    true,
				Description: "List of volumes attached to instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_Bootable: {
							Computed:    true,
							Description: "Indicates if the volume is boot capable.",
							Type:        schema.TypeBool,
						},
						Attr_Href: {
							Computed:    true,
							Description: "The hyper link of the volume.",
							Type:        schema.TypeString,
						},
						Attr_ID: {
							Computed:    true,
							Description: "The unique identifier of the volume.",
							Type:        schema.TypeString,
						},
						Attr_Name: {
							Computed:    true,
							Description: "The name of the volume.",
							Type:        schema.TypeString,
						},
						Attr_Pool: {
							Computed:    true,
							Description: "Volume pool, name of storage pool where the volume is located.",
							Type:        schema.TypeString,
						},
						Attr_Shareable: {
							Computed:    true,
							Description: "Indicates if the volume is shareable between VMs.",
							Type:        schema.TypeBool,
						},
						Attr_Size: {
							Computed:    true,
							Description: "The size of this volume in GB.",
							Type:        schema.TypeFloat,
						},
						Attr_State: {
							Computed:    true,
							Description: "The state of the volume.",
							Type:        schema.TypeString,
						},
						Attr_Type: {
							Computed:    true,
							Description: "The disk type that is used for this volume.",
							Type:        schema.TypeString,
						},
					},
				},
				Type: schema.TypeList,
			},
		},
	}
}

func dataSourceIBMPIInstanceVolumesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)

	volumeC := instance.NewIBMPIVolumeClient(ctx, sess, cloudInstanceID)
	volumedata, err := volumeC.GetAllInstanceVolumes(d.Get(Arg_InstanceName).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	var clientgenU, _ = uuid.GenerateUUID()
	d.SetId(clientgenU)
	d.Set(Attr_BootVolumeID, *volumedata.Volumes[0].VolumeID)
	d.Set(Attr_InstanceVolumes, flattenVolumesInstances(volumedata.Volumes))

	return nil
}

func flattenVolumesInstances(list []*models.VolumeReference) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		l := map[string]interface{}{
			Attr_Bootable:  *i.Bootable,
			Attr_Href:      *i.Href,
			Attr_ID:        *i.VolumeID,
			Attr_Name:      *i.Name,
			Attr_Pool:      i.VolumePool,
			Attr_Shareable: *i.Shareable,
			Attr_Size:      *i.Size,
			Attr_State:     *i.State,
			Attr_Type:      *i.DiskType,
		}
		result = append(result, l)
	}
	return result
}
