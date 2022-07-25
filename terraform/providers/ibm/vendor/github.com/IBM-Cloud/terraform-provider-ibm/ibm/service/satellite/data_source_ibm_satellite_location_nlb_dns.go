// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package satellite

import (
	"context"
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMSatelliteLocationNLBDNS() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMSatelliteLocationNLBDNSRead,

		Schema: map[string]*schema.Schema{
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A unique name of the Location",
			},
			"nlb_config": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of nlb config of Location",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"secret_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the secret.",
						},
						"secret_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of Secret.",
						},
						"cluster": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster Id.",
						},
						"dns_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of DNS.",
						},
						"lb_hostname": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Host Name of load Balancer.",
						},
						"nlb_ips": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: " NLB IPs.",
						},
						"nlb_sub_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "NLB Sub-Domain.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: " Nlb Type.",
						},
						"secret_namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Namespace of Secret.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMSatelliteLocationNLBDNSRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	location := d.Get("location").(string)

	satClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	nlbData, err := satClient.NlbDns().GetLocationNLBDNSList(location)
	if err != nil || nlbData == nil || len(nlbData) < 1 {
		return diag.FromErr(fmt.Errorf("[ERROR] Error Listing Satellite NLB DNS (%s): %s", location, err))
	}
	d.SetId(location)
	d.Set("location", location)
	d.Set("nlb_config", flex.FlattenNlbConfigs(nlbData))
	return nil
}
