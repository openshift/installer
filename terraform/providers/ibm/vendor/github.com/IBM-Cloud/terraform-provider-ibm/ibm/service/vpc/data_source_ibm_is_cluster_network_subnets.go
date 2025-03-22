// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsClusterNetworkSubnets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsClusterNetworkSubnetsRead,

		Schema: map[string]*schema.Schema{
			"cluster_network_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The cluster network identifier.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to resources with a `name` property matching the exact specified name.",
			},
			"sort": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "-created_at",
				Description: "Sorts the returned collection by the specified property name in ascending order. A `-` may be prepended to the name to sort in descending order. For example, the value `-created_at` sorts the collection by the `created_at` property in descending order, and the value `name` sorts it by the `name` property in ascending order.",
			},
			"subnets": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A page of subnets for the cluster network.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"available_ipv4_address_count": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of IPv4 addresses in this cluster network subnet that are not in use, and have not been reserved by the user or the provider.",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the cluster network subnet was created.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this cluster network subnet.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this cluster network subnet.",
						},
						"ip_version": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP version for this cluster network subnet.The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
						},
						"ipv4_cidr_block": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IPv4 range of this cluster network subnet, expressed in CIDR format.",
						},
						"lifecycle_reasons": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The reasons for the current `lifecycle_state` (if any).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A reason code for this lifecycle state:- `internal_error`: internal error (contact IBM support)- `resource_suspended_by_provider`: The resource has been suspended (contact IBM  support)The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
									},
									"message": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "An explanation of the reason for this lifecycle state.",
									},
									"more_info": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about the reason for this lifecycle state.",
									},
								},
							},
						},
						"lifecycle_state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lifecycle state of the cluster network subnet.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this cluster network subnet. The name is unique across all cluster network subnets in the cluster network.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						"total_ipv4_address_count": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of IPv4 addresses in this cluster network subnet.Note: This is calculated as 2<sup>(32 - prefix length)</sup>. For example, the prefix length `/24` gives:<br> 2<sup>(32 - 24)</sup> = 2<sup>8</sup> = 256 addresses.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIsClusterNetworkSubnetsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_cluster_network_subnets", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	listClusterNetworkSubnetsOptions := &vpcv1.ListClusterNetworkSubnetsOptions{}

	listClusterNetworkSubnetsOptions.SetClusterNetworkID(d.Get("cluster_network_id").(string))
	if _, ok := d.GetOk("name"); ok {
		listClusterNetworkSubnetsOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("sort"); ok {
		listClusterNetworkSubnetsOptions.SetSort(d.Get("sort").(string))
	}

	var pager *vpcv1.ClusterNetworkSubnetsPager
	pager, err = vpcClient.NewClusterNetworkSubnetsPager(listClusterNetworkSubnetsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_cluster_network_subnets", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	allItems, err := pager.GetAll()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ClusterNetworkSubnetsPager.GetAll() failed %s", err), "(Data) ibm_is_cluster_network_subnets", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIBMIsClusterNetworkSubnetsID(d))

	mapSlice := []map[string]interface{}{}
	for _, modelItem := range allItems {
		modelMap, err := DataSourceIBMIsClusterNetworkSubnetsClusterNetworkSubnetToMap(&modelItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_cluster_network_subnets", "read", "ClusterNetworks-to-map").GetDiag()
		}
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("subnets", mapSlice); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting subnets %s", err), "(Data) ibm_is_cluster_network_subnets", "read", "subnets-set").GetDiag()
	}

	return nil
}

// dataSourceIBMIsClusterNetworkSubnetsID returns a reasonable ID for the list.
func dataSourceIBMIsClusterNetworkSubnetsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIBMIsClusterNetworkSubnetsClusterNetworkSubnetToMap(model *vpcv1.ClusterNetworkSubnet) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["available_ipv4_address_count"] = flex.IntValue(model.AvailableIpv4AddressCount)
	modelMap["created_at"] = model.CreatedAt.String()
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["ip_version"] = *model.IPVersion
	modelMap["ipv4_cidr_block"] = *model.Ipv4CIDRBlock
	lifecycleReasons := []map[string]interface{}{}
	for _, lifecycleReasonsItem := range model.LifecycleReasons {
		lifecycleReasonsItemMap, err := DataSourceIBMIsClusterNetworkSubnetsClusterNetworkSubnetLifecycleReasonToMap(&lifecycleReasonsItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		lifecycleReasons = append(lifecycleReasons, lifecycleReasonsItemMap)
	}
	modelMap["lifecycle_reasons"] = lifecycleReasons
	modelMap["lifecycle_state"] = *model.LifecycleState
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	modelMap["total_ipv4_address_count"] = flex.IntValue(model.TotalIpv4AddressCount)
	return modelMap, nil
}

func DataSourceIBMIsClusterNetworkSubnetsClusterNetworkSubnetLifecycleReasonToMap(model *vpcv1.ClusterNetworkSubnetLifecycleReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = *model.Code
	modelMap["message"] = *model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}
