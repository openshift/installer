// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package directlink

import (
	"fmt"
	"log"

	"github.com/IBM/networking-go-sdk/directlinkv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMDLExportRouteFilter() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMDLExportRouteFilterRead,
		Schema: map[string]*schema.Schema{
			dlGatewayId: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Direct Link gateway identifier",
			},
			dlExportRouteFilterId: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Export route Filter identifier",
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
				Description: "The date and time of the export route filter was created",
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
				Description: "The date and time of the export route filter was last updated",
			},
		},
	}
}

func dataSourceIBMDLExportRouteFilterRead(d *schema.ResourceData, meta interface{}) error {
	directLink, err := directlinkClient(meta)
	if err != nil {
		return err
	}
	gatewayId := d.Get(dlGatewayId).(string)
	exportRouteFilterId := d.Get(dlExportRouteFilterId).(string)
	getGatewayExportRouteFilterOptionsModel := &directlinkv1.GetGatewayExportRouteFilterOptions{GatewayID: &gatewayId, ID: &exportRouteFilterId}
	exportRouteFilter, response, err := directLink.GetGatewayExportRouteFilter(getGatewayExportRouteFilterOptionsModel)
	if err != nil {
		log.Println("[ERROR] Error while listing the DL Export Route Filter", response, err)
		return err
	}
	if exportRouteFilter == nil {
		return fmt.Errorf("error while reading the Export Route filter for gateway: %s and Export route FilterId: %s with response code: %d", gatewayId, exportRouteFilterId, response.StatusCode)
	} else if exportRouteFilter.ID != nil {
		d.SetId(*exportRouteFilter.ID)
	}
	if exportRouteFilter.Action != nil {
		d.Set(dlAction, *exportRouteFilter.Action)
	}
	if exportRouteFilter.Before != nil {
		d.Set(dlBefore, *exportRouteFilter.Before)
	}
	if exportRouteFilter.CreatedAt != nil {
		d.Set(dlCreatedAt, exportRouteFilter.CreatedAt.String())
	}
	if exportRouteFilter.Prefix != nil {
		d.Set(dlPrefix, *exportRouteFilter.Prefix)
	}
	if exportRouteFilter.UpdatedAt != nil {
		d.Set(dlUpdatedAt, exportRouteFilter.UpdatedAt.String())
	}
	if exportRouteFilter.Ge != nil {
		d.Set(dlGe, *exportRouteFilter.Ge)
	}
	if exportRouteFilter.Le != nil {
		d.Set(dlLe, *exportRouteFilter.Le)
	}
	return nil
}
