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

const (
	Datacenters = "datacenters"
)

func DataSourceIBMPIDatacenters() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIDatacentersRead,
		Schema: map[string]*schema.Schema{
			Datacenters: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

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
				},
			},
		},
	}
}
func dataSourceIBMPIDatacentersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// session
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
				Attr_DatacenterLocation: map[string]interface{}{
					DatacenterRegion: datacenter.Location.Region,
					DatacenterType:   datacenter.Location.Type,
					DatacenterUrl:    datacenter.Location.URL,
				},
				Attr_DatacenterStatus: datacenter.Status,
				Attr_DatacenterType:   datacenter.Type,
				Attr_DatacenterHref:   datacenter.Href,
			}
			datacenters = append(datacenters, dc)
		}

	}
	var clientgenU, _ = uuid.GenerateUUID()
	d.SetId(clientgenU)
	d.Set(Datacenters, datacenters)
	return nil
}
