// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"log"

	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	//"fmt"
	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPINetworkPort() *schema.Resource {

	return &schema.Resource{
		ReadContext: dataSourceIBMPINetworkPortsRead,
		Schema: map[string]*schema.Schema{
			helpers.PINetworkName: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Network Name to be used for pvminstances",
				ValidateFunc: validation.NoZeroValues,
			},
			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			// Computed Attributes
			"network_ports": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipaddress": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"macaddress": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"portid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Required: true,
						},
						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPINetworkPortsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)
	networkportC := instance.NewIBMPINetworkClient(ctx, sess, cloudInstanceID)
	networkportdata, err := networkportC.GetAllPorts(d.Get(helpers.PINetworkName).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	var clientgenU, _ = uuid.GenerateUUID()
	d.SetId(clientgenU)

	d.Set("network_ports", flattenNetworkPorts(networkportdata.Ports))

	return nil

}

func flattenNetworkPorts(networkPorts []*models.NetworkPort) interface{} {
	result := make([]map[string]interface{}, 0, len(networkPorts))
	log.Printf("the number of ports is %d", len(networkPorts))
	for _, i := range networkPorts {
		l := map[string]interface{}{
			"portid":     *i.PortID,
			"status":     *i.Status,
			"href":       i.Href,
			"ipaddress":  *i.IPAddress,
			"macaddress": *i.MacAddress,
			"public_ip":  i.ExternalIP,
		}

		result = append(result, l)
	}
	return result
}
