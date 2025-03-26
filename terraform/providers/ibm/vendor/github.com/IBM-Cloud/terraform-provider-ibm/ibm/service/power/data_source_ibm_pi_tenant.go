// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPITenant() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPITenantRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_CloudInstances: {
				Computed:    true,
				Description: "Set of regions and Power Systems Virtual Server instance IDs that the tenant owns.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_CloudInstanceID: {
							Computed:    true,
							Description: "The unique identifier of the cloud instance.",
							Type:        schema.TypeString,
						},
						Attr_Region: {
							Computed:    true,
							Description: "The region of the cloud instance.",
							Type:        schema.TypeString,
						},
					},
				},
				Type: schema.TypeSet,
			},
			Attr_CreationDate: {
				Computed:    true,
				Description: "Date of tenant creation.",
				Type:        schema.TypeString,
			},
			Attr_Enabled: {
				Computed:    true,
				Description: "Indicates if the tenant is enabled for the Power Systems Virtual Server instance ID.",
				Type:        schema.TypeBool,
			},
			Attr_TenantName: {
				Computed:    true,
				Description: "The name of the tenant.",
				Type:        schema.TypeString,
			},
		},
	}
}

func dataSourceIBMPITenantRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)

	tenantC := instance.NewIBMPITenantClient(ctx, sess, cloudInstanceID)
	tenantData, err := tenantC.GetSelfTenant()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*tenantData.TenantID)
	d.Set(Attr_CreationDate, tenantData.CreationDate.String())
	d.Set(Attr_Enabled, tenantData.Enabled)

	if tenantData.CloudInstances != nil {
		d.Set(Attr_TenantName, tenantData.CloudInstances[0].Name)
	}

	if tenantData.CloudInstances != nil {
		tenants := make([]map[string]interface{}, len(tenantData.CloudInstances))
		for i, cloudinstance := range tenantData.CloudInstances {
			j := make(map[string]interface{})
			j[Attr_CloudInstanceID] = cloudinstance.CloudInstanceID
			j[Attr_Region] = cloudinstance.Region
			tenants[i] = j
		}

		d.Set(Attr_CloudInstances, tenants)
	}

	return nil
}
