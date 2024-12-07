// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMPIDatacenters() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIDatacentersRead,
		Schema: map[string]*schema.Schema{
			// Attributes
			Attr_Datacenters: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of Datacenters",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_DatacenterCapabilities: {
							Computed:    true,
							Description: "Datacenter Capabilities",
							Elem: &schema.Schema{
								Type: schema.TypeBool,
							},
							Type: schema.TypeMap,
						},
						Attr_DatacenterHref: {
							Computed:    true,
							Description: "Datacenter href",
							Type:        schema.TypeString,
						},
						Attr_DatacenterLocation: {
							Computed:    true,
							Description: "Datacenter location",
							Type:        schema.TypeMap,
						},
						Attr_DatacenterStatus: {
							Computed:    true,
							Description: "Datacenter status",
							Type:        schema.TypeString,
						},
						Attr_DatacenterType: {
							Computed:    true,
							Description: "Datacenter type",
							Type:        schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPIDatacentersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	client := instance.NewIBMPIDatacenterClient(ctx, sess, "")
	datacentersData, err := client.GetAll()
	if err != nil {
		return diag.FromErr(err)
	}
	datacenters := make([]map[string]interface{}, 0, len(datacentersData.Datacenters))
	for _, datacenter := range datacentersData.Datacenters {
		if datacenter != nil {
			dc := map[string]interface{}{
				Attr_DatacenterCapabilities: datacenter.Capabilities,
				Attr_DatacenterHref:         datacenter.Href,
				Attr_DatacenterLocation: map[string]interface{}{
					Attr_Region: datacenter.Location.Region,
					Attr_Type:   datacenter.Location.Type,
					Attr_URL:    datacenter.Location.URL,
				},
				Attr_DatacenterStatus: datacenter.Status,
				Attr_DatacenterType:   datacenter.Type,
			}
			datacenters = append(datacenters, dc)
		}
	}
	var clientgenU, _ = uuid.GenerateUUID()
	d.SetId(clientgenU)
	d.Set(Attr_Datacenters, datacenters)
	return nil
}
