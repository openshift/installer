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

func DataSourceIBMIsClusterNetworks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsClusterNetworksRead,

		Schema: map[string]*schema.Schema{
			"resource_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to resources with a `resource_group.id` property matching the specified identifier.",
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
			"vpc_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to cluster networks with a `vpc.id` property matching the specified id.",
			},
			"vpc_crn": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to cluster networks with a `vpc.crn` property matching the specified CRN.",
			},
			"vpc_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to cluster networks with a `vpc.name` property matching the specified name.",
			},
			"cluster_networks": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A page of cluster networks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the cluster network was created.",
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this cluster network.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this cluster network.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this cluster network.",
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
							Description: "The lifecycle state of the cluster network.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this cluster network. The name must not be used by another cluster network in the region.",
						},
						"profile": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The profile for this cluster network.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this cluster network profile.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The globally unique name for this cluster network profile.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"resource_group": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The resource group for this cluster network.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this resource group.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this resource group.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this resource group.",
									},
								},
							},
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						"subnet_prefixes": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The IP address ranges available for subnets for this cluster network.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allocation_policy": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The allocation policy for this subnet prefix:- `auto`: Subnets created by total count in this cluster network can use this prefix.",
									},
									"cidr": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CIDR block for this prefix.",
									},
								},
							},
						},
						"vpc": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The VPC this cluster network resides in.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this VPC.",
									},
									"deleted": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this VPC.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this VPC.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this VPC. The name is unique across all VPCs in the region.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"zone": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The zone this cluster network resides in.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this zone.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The globally unique name for this zone.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIsClusterNetworksRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_cluster_networks", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	listClusterNetworksOptions := &vpcv1.ListClusterNetworksOptions{}

	if _, ok := d.GetOk("resource_group_id"); ok {
		listClusterNetworksOptions.SetResourceGroupID(d.Get("resource_group_id").(string))
	}
	if _, ok := d.GetOk("name"); ok {
		listClusterNetworksOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("sort"); ok {
		listClusterNetworksOptions.SetSort(d.Get("sort").(string))
	}
	if _, ok := d.GetOk("vpc_id"); ok {
		listClusterNetworksOptions.SetVPCID(d.Get("vpc_id").(string))
	}
	if _, ok := d.GetOk("vpc_crn"); ok {
		listClusterNetworksOptions.SetVPCCRN(d.Get("vpc_crn").(string))
	}
	if _, ok := d.GetOk("vpc_name"); ok {
		listClusterNetworksOptions.SetVPCName(d.Get("vpc_name").(string))
	}

	var pager *vpcv1.ClusterNetworksPager
	pager, err = vpcClient.NewClusterNetworksPager(listClusterNetworksOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_cluster_networks", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	allItems, err := pager.GetAll()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ClusterNetworksPager.GetAll() failed %s", err), "(Data) ibm_is_cluster_networks", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIBMIsClusterNetworksID(d))

	mapSlice := []map[string]interface{}{}
	for _, modelItem := range allItems {
		modelMap, err := DataSourceIBMIsClusterNetworksClusterNetworkToMap(&modelItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_cluster_networks", "read", "ClusterNetworks-to-map").GetDiag()
		}
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("cluster_networks", mapSlice); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting cluster_networks %s", err), "(Data) ibm_is_cluster_networks", "read", "cluster_networks-set").GetDiag()
	}

	return nil
}

// dataSourceIBMIsClusterNetworksID returns a reasonable ID for the list.
func dataSourceIBMIsClusterNetworksID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIBMIsClusterNetworksClusterNetworkToMap(model *vpcv1.ClusterNetwork) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["created_at"] = model.CreatedAt.String()
	modelMap["crn"] = *model.CRN
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	lifecycleReasons := []map[string]interface{}{}
	for _, lifecycleReasonsItem := range model.LifecycleReasons {
		lifecycleReasonsItemMap, err := DataSourceIBMIsClusterNetworksClusterNetworkLifecycleReasonToMap(&lifecycleReasonsItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		lifecycleReasons = append(lifecycleReasons, lifecycleReasonsItemMap)
	}
	modelMap["lifecycle_reasons"] = lifecycleReasons
	modelMap["lifecycle_state"] = *model.LifecycleState
	modelMap["name"] = *model.Name
	profileMap, err := DataSourceIBMIsClusterNetworksClusterNetworkProfileReferenceToMap(model.Profile)
	if err != nil {
		return modelMap, err
	}
	modelMap["profile"] = []map[string]interface{}{profileMap}
	resourceGroupMap, err := DataSourceIBMIsClusterNetworksResourceGroupReferenceToMap(model.ResourceGroup)
	if err != nil {
		return modelMap, err
	}
	modelMap["resource_group"] = []map[string]interface{}{resourceGroupMap}
	modelMap["resource_type"] = *model.ResourceType
	subnetPrefixes := []map[string]interface{}{}
	for _, subnetPrefixesItem := range model.SubnetPrefixes {
		subnetPrefixesItemMap, err := DataSourceIBMIsClusterNetworksClusterNetworkSubnetPrefixToMap(&subnetPrefixesItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		subnetPrefixes = append(subnetPrefixes, subnetPrefixesItemMap)
	}
	modelMap["subnet_prefixes"] = subnetPrefixes
	vpcMap, err := DataSourceIBMIsClusterNetworksVPCReferenceToMap(model.VPC)
	if err != nil {
		return modelMap, err
	}
	modelMap["vpc"] = []map[string]interface{}{vpcMap}
	zoneMap, err := DataSourceIBMIsClusterNetworksZoneReferenceToMap(model.Zone)
	if err != nil {
		return modelMap, err
	}
	modelMap["zone"] = []map[string]interface{}{zoneMap}
	return modelMap, nil
}

func DataSourceIBMIsClusterNetworksClusterNetworkLifecycleReasonToMap(model *vpcv1.ClusterNetworkLifecycleReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = *model.Code
	modelMap["message"] = *model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}

func DataSourceIBMIsClusterNetworksClusterNetworkProfileReferenceToMap(model *vpcv1.ClusterNetworkProfileReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIBMIsClusterNetworksResourceGroupReferenceToMap(model *vpcv1.ResourceGroupReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func DataSourceIBMIsClusterNetworksClusterNetworkSubnetPrefixToMap(model *vpcv1.ClusterNetworkSubnetPrefix) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["allocation_policy"] = *model.AllocationPolicy
	modelMap["cidr"] = *model.CIDR
	return modelMap, nil
}

func DataSourceIBMIsClusterNetworksVPCReferenceToMap(model *vpcv1.VPCReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsClusterNetworksDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIBMIsClusterNetworksDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}

func DataSourceIBMIsClusterNetworksZoneReferenceToMap(model *vpcv1.ZoneReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["name"] = *model.Name
	return modelMap, nil
}
