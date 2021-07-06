// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceIBMPIInstanceVolumes() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceIBMPIInstanceVolumesRead,
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

func dataSourceIBMPIInstanceVolumesRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return err
	}

	powerinstanceid := d.Get(helpers.PICloudInstanceId).(string)
	volumeC := instance.NewIBMPIVolumeClient(sess, powerinstanceid)
	volumedata, err := volumeC.GetAll(d.Get(helpers.PIInstanceName).(string), powerinstanceid, getTimeOut)

	if err != nil {
		return err
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
			"shareable": *i.Shareable,
			"bootable":  *i.Bootable,
		}

		result = append(result, l)
	}
	return result
}
