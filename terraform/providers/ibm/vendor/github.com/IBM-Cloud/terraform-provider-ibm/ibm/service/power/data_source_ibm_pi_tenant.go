// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	//"fmt"
	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPITenant() *schema.Resource {

	return &schema.Resource{
		ReadContext: dataSourceIBMPITenantRead,
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

func dataSourceIBMPITenantRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)
	//tenantid := d.Get("tenantid").(string)

	tenantC := instance.NewIBMPITenantClient(ctx, sess, cloudInstanceID)
	tenantData, err := tenantC.GetSelfTenant()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*tenantData.TenantID)
	d.Set("creation_date", tenantData.CreationDate.String())
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
