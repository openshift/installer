// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"time"

	"github.com/IBM/networking-go-sdk/transitgatewayapisv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMTransitGateways() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceIBMTransitGatewaysRead,
		Schema: map[string]*schema.Schema{
			tgGateways: {
				Type:        schema.TypeList,
				Description: "Collection of transit gateways",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						tgID: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgCrn: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgName: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgLocation: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgCreatedAt: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgGlobal: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						tgStatus: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgUpdatedAt: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgResourceGroup: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMTransitGatewaysRead(d *schema.ResourceData, meta interface{}) error {

	client, err := transitgatewayClient(meta)
	if err != nil {
		return err
	}

	listTransitGatewaysOptionsModel := &transitgatewayapisv1.ListTransitGatewaysOptions{}
	listTransitGateways, response, err := client.ListTransitGateways(listTransitGatewaysOptionsModel)
	if err != nil {
		return fmt.Errorf("Error while listing transit gateways %s\n%s", err, response)
	}

	tgws := make([]map[string]interface{}, 0)
	for _, instance := range listTransitGateways.TransitGateways {

		transitgateway := map[string]interface{}{}
		transitgateway[tgID] = instance.ID
		transitgateway[tgName] = instance.Name
		transitgateway[tgCreatedAt] = instance.CreatedAt.String()
		transitgateway[tgLocation] = instance.Location
		transitgateway[tgStatus] = instance.Status

		if instance.UpdatedAt != nil {
			transitgateway[tgUpdatedAt] = instance.UpdatedAt.String()
		}
		transitgateway[tgGlobal] = instance.Global
		transitgateway[tgCrn] = instance.Crn

		if instance.ResourceGroup != nil {
			rg := instance.ResourceGroup
			transitgateway[tgResourceGroup] = *rg.ID
		}

		tgws = append(tgws, transitgateway)
	}
	d.Set(tgGateways, tgws)
	d.SetId(dataSourceIBMTransitGatewaysID(d))
	return nil
}

// dataSourceIBMTransitGatewaysID returns a reasonable ID for a transit gateways list.
func dataSourceIBMTransitGatewaysID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
