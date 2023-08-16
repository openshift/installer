// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package directlink

import (
	"fmt"
	"log"

	"github.com/IBM/networking-go-sdk/directlinkv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMDLImportRouteFilter() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMDLImportRouteFilterRead,
		Schema: map[string]*schema.Schema{
			dlGatewayId: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Direct Link gateway identifier",
			},
			dlImportRouteFilterId: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Import route Filter identifier",
			},
			dlAction: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Determines whether the  routes that match the prefix-set will be permit or deny",
			},
			dlBefore: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identifier of the next route filter to be considered",
			},
			dlCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of the import route filter was created",
			},
			dlGe: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The minimum matching length of the prefix-set",
			},
			dlLe: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum matching length of the prefix-set",
			},
			dlPrefix: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP prefix representing an address and mask length of the prefix-set",
			},
			dlUpdatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of the import route filter was last updated",
			},
		},
	}
}

func dataSourceIBMDLImportRouteFilterRead(d *schema.ResourceData, meta interface{}) error {
	directLink, err := directlinkClient(meta)
	if err != nil {
		return err
	}
	gatewayId := d.Get(dlGatewayId).(string)
	importRouteFilterId := d.Get(dlImportRouteFilterId).(string)
	getGatewayImportRouteFilterOptionsModel := &directlinkv1.GetGatewayImportRouteFilterOptions{GatewayID: &gatewayId, ID: &importRouteFilterId}
	importRouteFilter, response, err := directLink.GetGatewayImportRouteFilter(getGatewayImportRouteFilterOptionsModel)
	if err != nil {
		log.Println("[ERROR] Error  while listing Direct Link Import Route Filter", response, err)
		return err
	}
	if importRouteFilter == nil {
		return fmt.Errorf("error fetching  the Import route Filter for gateway: %s and Import Route FilterId: %s with response code: %d", gatewayId, importRouteFilterId, response.StatusCode)
	} else if importRouteFilter.ID != nil {
		d.SetId(*importRouteFilter.ID)
	}
	if importRouteFilter.Action != nil {
		d.Set(dlAction, *importRouteFilter.Action)
	}
	if importRouteFilter.Before != nil {
		d.Set(dlBefore, *importRouteFilter.Before)
	}
	if importRouteFilter.CreatedAt != nil {
		d.Set(dlCreatedAt, importRouteFilter.CreatedAt.String())
	}
	if importRouteFilter.Prefix != nil {
		d.Set(dlPrefix, *importRouteFilter.Prefix)
	}
	if importRouteFilter.UpdatedAt != nil {
		d.Set(dlUpdatedAt, importRouteFilter.UpdatedAt.String())
	}
	if importRouteFilter.Ge != nil {
		d.Set(dlGe, *importRouteFilter.Ge)
	}
	if importRouteFilter.Le != nil {
		d.Set(dlLe, *importRouteFilter.Le)
	}
	return nil
}
