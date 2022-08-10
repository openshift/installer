// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package transitgateway

import (
	"fmt"
	"github.com/IBM/networking-go-sdk/transitgatewayapisv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	tgPrefixFilterId = "filter_id"
)

func DataSourceIBMTransitGatewayConnectionPrefixFilter() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceIBMTransitGatewayConnectionPrefixFilterRead,
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
			tgPrefixFilterId: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Transit Gateway Connection Prefix Filter identifier",
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
	}
}

func dataSourceIBMTransitGatewayConnectionPrefixFilterRead(d *schema.ResourceData, meta interface{}) error {
	client, err := transitgatewayClient(meta)
	if err != nil {
		return err
	}

	gatewayId := d.Get(tgGatewayId).(string)
	connectionId := d.Get(tgConnectionId).(string)
	filterId := d.Get(tgPrefixFilterId).(string)

	getTransitGatewayConnectionPrefixFilterOptionsModel := &transitgatewayapisv1.GetTransitGatewayConnectionPrefixFilterOptions{}
	getTransitGatewayConnectionPrefixFilterOptionsModel.SetTransitGatewayID(gatewayId)
	getTransitGatewayConnectionPrefixFilterOptionsModel.SetID(connectionId)
	getTransitGatewayConnectionPrefixFilterOptionsModel.SetFilterID(filterId)
	prefixFilter, response, err := client.GetTransitGatewayConnectionPrefixFilter(getTransitGatewayConnectionPrefixFilterOptionsModel)
	if err != nil {
		return fmt.Errorf("Error retrieving transit gateway connection prefix filter (%s): %s\n%s", filterId, err, response)
	}

	d.SetId(*prefixFilter.ID)
	d.Set(tgPrefixFilterId, prefixFilter.ID)
	d.Set(tgAction, prefixFilter.Action)
	d.Set(tgCreatedAt, prefixFilter.CreatedAt.String())
	d.Set(tgPrefix, prefixFilter.Prefix)

	if prefixFilter.UpdatedAt != nil {
		d.Set(tgUpdatedAt, prefixFilter.UpdatedAt.String())
	}
	if prefixFilter.Before != nil {
		d.Set(tgBefore, prefixFilter.Before)
	}
	if prefixFilter.Ge != nil {
		d.Set(tgGe, prefixFilter.Ge)
	}
	if prefixFilter.Le != nil {
		d.Set(tgLe, prefixFilter.Le)
	}

	return nil
}
