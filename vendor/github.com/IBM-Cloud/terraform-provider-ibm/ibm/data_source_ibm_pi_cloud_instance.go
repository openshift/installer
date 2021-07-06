// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/power/models"
)

func dataSourceIBMPICloudInstance() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceIBMPICloudInstanceRead,
		Schema: map[string]*schema.Schema{
			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			// Start of Computed Attributes

			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"capabilities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"total_processors_consumed": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"total_instances": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"total_memory_consumed": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"total_ssd_storage_consumed": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"total_standard_storage_consumed": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"pvm_instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"systype": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPICloudInstanceRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := meta.(ClientSession).IBMPISession()

	if err != nil {
		return err
	}

	powerinstanceid := d.Get(helpers.PICloudInstanceId).(string)
	cloud_instance := instance.NewIBMPICloudInstanceClient(sess, powerinstanceid)
	cloud_instance_data, err := cloud_instance.Get(powerinstanceid)

	if err != nil {
		return err
	}

	d.SetId(*cloud_instance_data.CloudInstanceID)
	d.Set("tenant_id", (cloud_instance_data.TenantID))
	d.Set("enabled", cloud_instance_data.Enabled)
	d.Set("region", cloud_instance_data.Region)
	d.Set("capabilities", cloud_instance_data.Capabilities)
	d.Set("pvm_instances", flattenpvminstances(cloud_instance_data.PvmInstances))
	d.Set("total_ssd_storage_consumed", cloud_instance_data.Usage.StorageSSD)
	d.Set("total_instances", cloud_instance_data.Usage.Instances)
	d.Set("total_standard_storage_consumed", cloud_instance_data.Usage.StorageStandard)
	d.Set("total_processors_consumed", cloud_instance_data.Usage.Processors)
	d.Set("total_memory_consumed", cloud_instance_data.Usage.Memory)

	return nil

}

func flattenpvminstances(list []*models.PVMInstanceReference) []map[string]interface{} {
	pvms := make([]map[string]interface{}, 0)
	for _, lpars := range list {

		l := map[string]interface{}{
			"id":            *lpars.PvmInstanceID,
			"name":          *lpars.ServerName,
			"href":          *lpars.Href,
			"status":        *lpars.Status,
			"systype":       lpars.SysType,
			"creation_date": lpars.CreationDate.String(),
		}
		pvms = append(pvms, l)

	}
	return pvms
}
