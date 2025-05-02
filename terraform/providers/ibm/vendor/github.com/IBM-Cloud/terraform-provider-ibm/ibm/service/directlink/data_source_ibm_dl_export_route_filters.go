// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package directlink

import (
	"log"
	"time"

	"github.com/IBM/networking-go-sdk/directlinkv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMDLExportRouteFilters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMDLExportRouteFiltersRead,
		Schema: map[string]*schema.Schema{
			dlGatewayId: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Direct Link gateway identifier",
			},
			dlExportRouteFilters: {
				Type:        schema.TypeList,
				Description: "Collection of export route filters",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						dlExportRouteFilterId: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Export Route Filter identifier",
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
				},
			},
		},
	}
}

func dataSourceIBMDLExportRouteFiltersRead(d *schema.ResourceData, meta interface{}) error {
	directLink, err := directlinkClient(meta)
	if err != nil {
		return err
	}

	gatewayId := d.Get(dlGatewayId).(string)
	listGatewayExportRouteFiltersOptionsModel := &directlinkv1.ListGatewayExportRouteFiltersOptions{GatewayID: &gatewayId}
	exportRouteFilterList, response, err := directLink.ListGatewayExportRouteFilters(listGatewayExportRouteFiltersOptionsModel)
	if err != nil {
		log.Println("[ERROR] Error listing Direct Link Export Route Filters", response, err)
		return err
	}
	exportRouteFilters := make([]map[string]interface{}, 0)
	for _, instance := range exportRouteFilterList.ExportRouteFilters {
		routeFilter := map[string]interface{}{}
		if instance.ID != nil {
			routeFilter[dlExportRouteFilterId] = *instance.ID
		}
		if instance.Action != nil {
			routeFilter[dlAction] = *instance.Action
		}
		if instance.Before != nil {
			routeFilter[dlBefore] = *instance.Before
		}
		if instance.CreatedAt != nil {
			routeFilter[dlCreatedAt] = instance.CreatedAt.String()
		}
		if instance.Prefix != nil {
			routeFilter[dlPrefix] = *instance.Prefix
		}
		if instance.UpdatedAt != nil {
			routeFilter[dlUpdatedAt] = instance.UpdatedAt.String()
		}
		if instance.Ge != nil {
			routeFilter[dlGe] = *instance.Ge
		}
		if instance.Le != nil {
			routeFilter[dlLe] = *instance.Le
		}
		exportRouteFilters = append(exportRouteFilters, routeFilter)
	}
	d.Set(dlExportRouteFilters, exportRouteFilters)
	d.SetId(dataSourceIBMDirectLinkGatewayExportRouteFiltersID(d))
	return nil
}

// dataSourceIBMDirectLinkGatewayRouteReportsID returns a reasonable ID for a directlink gateways list.
func dataSourceIBMDirectLinkGatewayExportRouteFiltersID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
