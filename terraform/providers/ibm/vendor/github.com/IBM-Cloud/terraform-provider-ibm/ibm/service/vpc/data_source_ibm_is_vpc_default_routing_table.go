// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isDefaultRoutingTableID             = "default_routing_table"
	isDefaultRoutingTableHref           = "href"
	isDefaultRoutingTableCrn            = "crn"
	isDefaultRoutingTableName           = "name"
	isDefaultRoutingTableResourceType   = "resource_type"
	isDefaultRoutingTableCreatedAt      = "created_at"
	isDefaultRoutingTableLifecycleState = "lifecycle_state"
	isDefaultRoutingTableRoutesList     = "routes"
	isDefaultRoutingTableSubnetsList    = "subnets"
	isDefaultRTVpcID                    = "vpc"
	isDefaultRTDirectLinkIngress        = "route_direct_link_ingress"
	isDefaultRTInternetIngress          = "route_internet_ingress"
	isDefaultRTTransitGatewayIngress    = "route_transit_gateway_ingress"
	isDefaultRTVPCZoneIngress           = "route_vpc_zone_ingress"
	isDefaultRTDefault                  = "is_default"
	isDefaultRTResourceGroup            = "resource_group"
	isDefaultRTResourceGroupHref        = "href"
	isDefaultRTResourceGroupId          = "id"
	isDefaultRTResourceGroupName        = "name"
	isDefaultRTTags                     = "tags"
	isDefaultRTAccessTags               = "access_tags"
	isDefaultRTAccessTagType            = "access"
	isDefaultRTUserTagType              = "user"
)

func DataSourceIBMISVPCDefaultRoutingTable() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISVPCDefaultRoutingTableGet,
		Schema: map[string]*schema.Schema{
			isDefaultRTVpcID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPC identifier",
			},
			isDefaultRoutingTableID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default Routing Table ID",
			},
			isDefaultRoutingTableHref: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default Routing table Href",
			},
			isDefaultRoutingTableName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default Routing table Name",
			},
			isDefaultRoutingTableCrn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default Routing table Crn",
			},
			isDefaultRoutingTableResourceType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default Routing table Resource Type",
			},
			isDefaultRoutingTableCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default Routing table Created At",
			},
			isDefaultRoutingTableLifecycleState: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default Routing table Lifecycle State",
			},
			isDefaultRTDirectLinkIngress: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If set to true, this routing table will be used to route traffic that originates from Direct Link to this VPC.",
			},
			isDefaultRTInternetIngress: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If set to true, this routing table will be used to route traffic that originates from the internet. For this to succeed, the VPC must not already have a routing table with this property set to true.",
			},
			isDefaultRTTransitGatewayIngress: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If set to true, this routing table will be used to route traffic that originates from Transit Gateway to this VPC.",
			},
			isDefaultRTVPCZoneIngress: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If set to true, this routing table will be used to route traffic that originates from subnets in other zones in this VPC.",
			},
			isDefaultRTDefault: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this is the default routing table for this VPC",
			},
			isDefaultRoutingTableRoutesList: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Route name",
						},

						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Route ID",
						},
					},
				},
			},
			isDefaultRoutingTableSubnetsList: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet name",
						},

						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet ID",
						},
					},
				},
			},
			isDefaultRTResourceGroup: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The resource group for this volume.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isDefaultRTResourceGroupHref: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this resource group.",
						},
						isDefaultRTResourceGroupId: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this resource group.",
						},
						isDefaultRTResourceGroupName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this resource group.",
						},
					},
				},
			},
			isDefaultRTTags: {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      flex.ResourceIBMVPCHash,
			},
			isDefaultRTAccessTags: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of access tags",
			},
		},
	}
}

func dataSourceIBMISVPCDefaultRoutingTableGet(d *schema.ResourceData, meta interface{}) error {

	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	vpcID := d.Get(isDefaultRTVpcID).(string)

	getVpcDefaultRoutingTableOptions := sess.NewGetVPCDefaultRoutingTableOptions(vpcID)
	result, detail, err := sess.GetVPCDefaultRoutingTable(getVpcDefaultRoutingTableOptions)
	if err != nil || result == nil {
		log.Printf("Error reading details of VPC Default Routing Table:%s", detail)
		return err
	}
	d.Set(isDefaultRoutingTableID, *result.ID)
	d.Set(isDefaultRoutingTableHref, *result.Href)
	d.Set(isDefaultRoutingTableName, *result.Name)
	d.Set(isDefaultRoutingTableCrn, *result.CRN)
	d.Set(isDefaultRoutingTableResourceType, *result.ResourceType)
	createdAt := *result.CreatedAt
	d.Set(isDefaultRoutingTableCreatedAt, createdAt.String())
	d.Set(isDefaultRoutingTableLifecycleState, *result.LifecycleState)
	d.Set(isDefaultRTDirectLinkIngress, *result.RouteDirectLinkIngress)
	d.Set(isDefaultRTInternetIngress, *result.RouteInternetIngress)
	d.Set(isDefaultRTTransitGatewayIngress, *result.RouteTransitGatewayIngress)
	d.Set(isDefaultRTVPCZoneIngress, *result.RouteVPCZoneIngress)
	d.Set(isDefaultRTDefault, *result.IsDefault)
	subnetsInfo := make([]map[string]interface{}, 0)
	for _, subnet := range result.Subnets {
		if subnet.Name != nil && subnet.ID != nil {
			l := map[string]interface{}{
				"name": *subnet.Name,
				"id":   *subnet.ID,
			}
			subnetsInfo = append(subnetsInfo, l)
		}
	}
	d.Set(isDefaultRoutingTableSubnetsList, subnetsInfo)
	routesInfo := make([]map[string]interface{}, 0)
	for _, route := range result.Routes {
		if route.Name != nil && route.ID != nil {
			k := map[string]interface{}{
				"name": *route.Name,
				"id":   *route.ID,
			}
			routesInfo = append(routesInfo, k)
		}
	}
	d.Set(isDefaultRoutingTableRoutesList, routesInfo)
	resourceGroupList := []map[string]interface{}{}
	if result.ResourceGroup != nil {
		resourceGroupMap := routingTableResourceGroupToMap(*result.ResourceGroup)
		resourceGroupList = append(resourceGroupList, resourceGroupMap)
	}
	d.Set(isDefaultRTResourceGroup, resourceGroupList)
	d.Set(isDefaultRTVpcID, vpcID)
	d.SetId(*result.ID)

	tags, err := flex.GetGlobalTagsUsingCRN(meta, *result.CRN, "", isDefaultRTUserTagType)
	if err != nil {
		log.Printf(
			"An error occured during reading of default routing table (%s) tags : %s", d.Id(), err)
	}
	d.Set(isDefaultRTTags, tags)

	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *result.CRN, "", isDefaultRTAccessTagType)
	if err != nil {
		log.Printf(
			"An error occured during reading of default routing table (%s) access tags: %s", d.Id(), err)
	}
	d.Set(isDefaultRTAccessTags, accesstags)
	return nil
}
