// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	//"fmt"
	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceIBMPITenant() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceIBMPITenantRead,
		Schema: map[string]*schema.Schema{
			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			// Computed Attributes
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"tenant_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloud_instances": {

				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPITenantRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := meta.(ClientSession).IBMPISession()

	if err != nil {
		return err
	}

	powerinstanceid := d.Get(helpers.PICloudInstanceId).(string)
	//tenantid := d.Get("tenantid").(string)

	tenantC := instance.NewIBMPITenantClient(sess, powerinstanceid)
	tenantData, err := tenantC.Get(powerinstanceid)

	if err != nil {
		return err
	}

	d.SetId(*tenantData.TenantID)
	d.Set("creation_date", tenantData.CreationDate)
	d.Set("enabled", tenantData.Enabled)

	if tenantData.CloudInstances != nil {

		d.Set("tenant_name", tenantData.CloudInstances[0].Name)
	}

	if tenantData.CloudInstances != nil {
		tenants := make([]map[string]interface{}, len(tenantData.CloudInstances))
		for i, cloudinstance := range tenantData.CloudInstances {
			j := make(map[string]interface{})
			j["region"] = cloudinstance.Region
			j["cloud_instance_id"] = cloudinstance.CloudInstanceID
			tenants[i] = j
		}

		d.Set("cloud_instances", tenants)
	}

	return nil

}
