// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMResourceQuota() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMResourceQuotaRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Resource quota name, for example Trial Quota",
				Type:        schema.TypeString,
				Required:    true,
			},
			"type": {
				Description: "Type of the quota.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"max_apps": {
				Description: "Defines the total app limit.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"max_instances_per_app": {
				Description: "Defines the total instances limit per app.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"max_app_instance_memory": {
				Description: "Defines the total memory of app instance.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"total_app_memory": {
				Description: "Defines the total memory for app.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"max_service_instances": {
				Description: "Defines the total service instances limit.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"vsi_limit": {
				Description: "Defines the VSI limit.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func dataSourceIBMResourceQuotaRead(d *schema.ResourceData, meta interface{}) error {
	rsManagementAPI, err := meta.(ClientSession).ResourceManagementAPIv2()
	if err != nil {
		return err
	}
	rsQuota := rsManagementAPI.ResourceQuota()
	rsQuotaName := d.Get("name").(string)
	rsQuotas, err := rsQuota.FindByName(rsQuotaName)
	if err != nil {
		return fmt.Errorf("Error retrieving resource quota: %s", err)
	}

	if len(rsQuotas) == 0 {
		return fmt.Errorf("Error retrieving resource quota: %s", err)
	}

	rsQuotaFields := rsQuotas[0]
	d.SetId(rsQuotaFields.ID)
	d.Set("type", rsQuotaFields.Type)
	d.Set("max_apps", rsQuotaFields.AppCountLimit)
	d.Set("max_instances_per_app", rsQuotaFields.AppInstanceCountLimit)
	d.Set("max_app_instance_memory", rsQuotaFields.AppInstanceMemoryLimit)
	d.Set("total_app_memory", rsQuotaFields.TotalAppMemoryLimit)
	d.Set("max_service_instances", rsQuotaFields.ServiceInstanceCountLimit)
	d.Set("vsi_limit", rsQuotaFields.VSICountLimit)
	return nil
}
