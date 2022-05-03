// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPIInstanceVolumes() *schema.Resource {

	return &schema.Resource{
		ReadContext: dataSourceIBMPIInstanceVolumesRead,
		Schema: map[string]*schema.Schema{
			helpers.PIInstanceName: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Instance Name to be used for pvminstances",
				ValidateFunc: validation.NoZeroValues,
			},
			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			//Computed Attributes
			"boot_volume_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pool": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"shareable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"bootable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPIInstanceVolumesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)

	volumeC := instance.NewIBMPIVolumeClient(ctx, sess, cloudInstanceID)
	volumedata, err := volumeC.GetAllInstanceVolumes(d.Get(helpers.PIInstanceName).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	var clientgenU, _ = uuid.GenerateUUID()
	d.SetId(clientgenU)
	d.Set("boot_volume_id", *volumedata.Volumes[0].VolumeID)
	d.Set("instance_volumes", flattenVolumesInstances(volumedata.Volumes))

	return nil

}

func flattenVolumesInstances(list []*models.VolumeReference) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		l := map[string]interface{}{
			"id":        *i.VolumeID,
			"state":     *i.State,
			"href":      *i.Href,
			"name":      *i.Name,
			"size":      *i.Size,
			"type":      *i.DiskType,
			"pool":      i.VolumePool,
			"shareable": *i.Shareable,
			"bootable":  *i.Bootable,
		}

		result = append(result, l)
	}
	return result
}
