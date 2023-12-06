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

const (
	DatacenterRegion = "region"
	DatacenterType   = "type"
	DatacenterUrl    = "url"
)

func DataSourceIBMPIDatacenter() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIDatacenterRead,
		Schema: map[string]*schema.Schema{
			Arg_DatacenterZone: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			Attr_DatacenterCapabilities: {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Datacenter Capabilities",
				Elem: &schema.Schema{
					Type: schema.TypeBool,
				},
			},
			Attr_DatacenterHref: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Datacenter href",
			},
			Attr_DatacenterLocation: {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Datacenter location",
			},
			Attr_DatacenterStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Datacenter status",
			},
			Attr_DatacenterType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Datacenter type",
			},
		},
	}
}
func dataSourceIBMPIDatacenterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// session
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
		DatacenterRegion: *dcData.Location.Region,
		DatacenterType:   *dcData.Location.Type,
		DatacenterUrl:    *dcData.Location.URL,
	}
	d.Set(Attr_DatacenterLocation, flex.Flatten(dclocation))
	d.Set(Attr_DatacenterStatus, dcData.Status)
	d.Set(Attr_DatacenterType, dcData.Type)
	d.Set(Attr_DatacenterHref, dcData.Href)

	return nil
}
