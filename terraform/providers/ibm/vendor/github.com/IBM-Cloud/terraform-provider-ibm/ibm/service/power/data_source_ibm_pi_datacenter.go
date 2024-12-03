// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPIDatacenter() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIDatacenterRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_DatacenterZone: {
				Description:  "Datacenter zone you want to retrieve.",
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_DatacenterCapabilities: {
				Computed:    true,
				Description: "Datacenter Capabilities.",
				Elem: &schema.Schema{
					Type: schema.TypeBool,
				},
				Type: schema.TypeMap,
			},
			Attr_DatacenterHref: {
				Computed:    true,
				Description: "Datacenter href.",
				Type:        schema.TypeString,
			},
			Attr_DatacenterLocation: {
				Computed:    true,
				Description: "Datacenter location.",
				Type:        schema.TypeMap,
			},
			Attr_DatacenterStatus: {
				Computed:    true,
				Description: "Datacenter status, active,maintenance or down.",
				Type:        schema.TypeString,
			},
			Attr_DatacenterType: {
				Computed:    true,
				Description: "Datacenter type, off-premises or on-premises.",
				Type:        schema.TypeString,
			},
		},
	}
}

func dataSourceIBMPIDatacenterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	datacenterZone := sess.Options.Zone
	if region, ok := d.GetOk(Arg_DatacenterZone); ok {
		datacenterZone = region.(string)
	}
	client := instance.NewIBMPIDatacenterClient(ctx, sess, "")
	dcData, err := client.Get(datacenterZone)
	if err != nil {
		return diag.FromErr(err)
	}
	var genID, _ = uuid.GenerateUUID()
	d.SetId(genID)
	d.Set(Attr_DatacenterCapabilities, dcData.Capabilities)
	dclocation := map[string]interface{}{
		Attr_Region: *dcData.Location.Region,
		Attr_Type:   *dcData.Location.Type,
		Attr_URL:    *dcData.Location.URL,
	}
	d.Set(Attr_DatacenterHref, dcData.Href)
	d.Set(Attr_DatacenterLocation, flex.Flatten(dclocation))
	d.Set(Attr_DatacenterStatus, dcData.Status)
	d.Set(Attr_DatacenterType, dcData.Type)

	return nil
}
