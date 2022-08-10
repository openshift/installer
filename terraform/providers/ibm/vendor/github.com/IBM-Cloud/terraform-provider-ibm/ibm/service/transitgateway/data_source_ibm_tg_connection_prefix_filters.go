// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package transitgateway

import (
	"fmt"
	"github.com/IBM/networking-go-sdk/transitgatewayapisv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

const (
	tgPrefixFilters = "prefix_filters"
	tgAction        = "action"
	tgBefore        = "before"
	tgGe            = "ge"
	tgLe            = "le"
	tgPrefix        = "prefix"
)

func DataSourceIBMTransitGatewayConnectionPrefixFilters() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceIBMTransitGatewayConnectionPrefixFiltersRead,
		Schema: map[string]*schema.Schema{
			tgGatewayId: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Transit Gateway identifier",
			},
			tgConnectionId: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Transit Gateway Connection identifier",
			},
			tgPrefixFilters: {
				Type:        schema.TypeList,
				Description: "Collection of prefix filters",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						tgID: {
							Type:     schema.TypeString,
							Computed: true,
						},
						tgAction: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether to permit or deny the prefix filter",
						},
						tgBefore: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identifier of prefix filter that handles ordering",
						},
						tgCreatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that this prefix filter was created",
						},
						tgGe: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "IP Prefix GE",
						},
						tgLe: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "IP Prefix LE",
						},
						tgPrefix: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP Prefix",
						},
						tgUpdatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that this prefix filter was last updated",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMTransitGatewayConnectionPrefixFiltersRead(d *schema.ResourceData, meta interface{}) error {
	client, err := transitgatewayClient(meta)
	if err != nil {
		return err
	}

	gatewayId := d.Get(tgGatewayId).(string)
	connectionId := d.Get(tgConnectionId).(string)

	listTransitGatewayConnectionPrefixFiltersOptionsModel := &transitgatewayapisv1.ListTransitGatewayConnectionPrefixFiltersOptions{}
	listTransitGatewayConnectionPrefixFiltersOptionsModel.SetTransitGatewayID(gatewayId)
	listTransitGatewayConnectionPrefixFiltersOptionsModel.SetID(connectionId)
	listPrefixFilters, response, err := client.ListTransitGatewayConnectionPrefixFilters(listTransitGatewayConnectionPrefixFiltersOptionsModel)
	if err != nil {
		return fmt.Errorf("Error while listing transit gateway connection prefix filters %s\n%s", err, response)
	}

	prefixFiltersCollection := make([]map[string]interface{}, 0)
	for _, prefixFilter := range listPrefixFilters.PrefixFilters {
		tgPrefixFilter := map[string]interface{}{}
		tgPrefixFilter[tgID] = prefixFilter.ID
		tgPrefixFilter[tgAction] = prefixFilter.Action
		tgPrefixFilter[tgCreatedAt] = prefixFilter.CreatedAt.String()
		tgPrefixFilter[tgPrefix] = prefixFilter.Prefix

		if prefixFilter.UpdatedAt != nil {
			tgPrefixFilter[tgUpdatedAt] = prefixFilter.UpdatedAt.String()
		}
		if prefixFilter.Before != nil {
			tgPrefixFilter[tgBefore] = prefixFilter.Before
		}
		if prefixFilter.Ge != nil {
			tgPrefixFilter[tgGe] = prefixFilter.Ge
		}
		if prefixFilter.Le != nil {
			tgPrefixFilter[tgLe] = prefixFilter.Le
		}

		prefixFiltersCollection = append(prefixFiltersCollection, tgPrefixFilter)
	}

	d.Set(tgPrefixFilters, prefixFiltersCollection)
	d.SetId(dataSourceIBMTransitGatewayConnectionPrefixFiltersID(d))
	return nil
}

// dataSourceIBMTransitGatewayRouteReportsID returns a reasonable ID for a transit gateways list.
func dataSourceIBMTransitGatewayConnectionPrefixFiltersID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
