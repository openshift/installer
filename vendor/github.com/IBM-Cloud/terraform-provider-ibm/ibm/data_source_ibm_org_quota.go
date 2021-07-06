// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMOrgQuota() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMOrgQuotaRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Org quota name, for example qIBM",
				Type:        schema.TypeString,
				Required:    true,
			},
			"non_basic_services_allowed": {
				Description: "Define non basic services are allowed for organization.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"total_services": {
				Description: "Defines the total services for organization.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"total_routes": {
				Description: "Defines the total route for organization.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"memory_limit": {
				Description: "Defines the total memory limit for organization.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"instance_memory_limit": {
				Description: "Defines the  total instance memory limit for organization.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"trial_db_allowed": {
				Description: "Defines trial db are allowed for organization.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"app_instance_limit": {
				Description: "Defines the total app instance limit for organization.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"total_private_domains": {
				Description: "Defines the total private domain limit for organization.v",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"app_tasks_limit": {
				Description: "Defines the total app task limit for organization.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"total_service_keys": {
				Description: "Defines the total service keys for organization.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"total_reserved_route_ports": {
				Description: "Defines the number of reserved route ports for organization. ",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func dataSourceIBMOrgQuotaRead(d *schema.ResourceData, meta interface{}) error {
	cfClient, err := meta.(ClientSession).MccpAPI()
	if err != nil {
		return err
	}
	orgQuotaAPI := cfClient.OrgQuotas()
	orgQuotaName := d.Get("name").(string)
	orgQuotaFields, err := orgQuotaAPI.FindByName(orgQuotaName)
	if err != nil {
		return fmt.Errorf("Error retrieving org quota: %s", err)
	}
	d.SetId(orgQuotaFields.GUID)
	d.Set("app_instance_limit", orgQuotaFields.AppInstanceLimit)
	d.Set("app_tasks_limit", orgQuotaFields.AppTasksLimit)
	d.Set("instance_memory_limit", orgQuotaFields.InstanceMemoryLimitInMB)
	d.Set("memory_limit", orgQuotaFields.MemoryLimitInMB)
	d.Set("non_basic_services_allowed", orgQuotaFields.NonBasicServicesAllowed)
	d.Set("total_private_domains", orgQuotaFields.PrivateDomainsLimit)
	d.Set("total_reserved_route_ports", orgQuotaFields.RoutePortsLimit)
	d.Set("total_routes", orgQuotaFields.RoutesLimit)
	d.Set("total_service_keys", orgQuotaFields.ServiceKeysLimit)
	d.Set("total_services", orgQuotaFields.ServicesLimit)
	d.Set("trial_db_allowed", orgQuotaFields.TrialDBAllowed)
	return nil
}
