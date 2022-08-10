// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIbmIsPlacementGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmIsPlacementGroupsRead,

		Schema: map[string]*schema.Schema{
			"placement_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of placement groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the placement group was created.",
						},
						"crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this placement group.",
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this placement group.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this placement group.",
						},
						"lifecycle_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lifecycle state of the placement group.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this placement group.",
						},
						"resource_group": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The resource group for this placement group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this resource group.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this resource group.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this resource group.",
									},
								},
							},
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						"strategy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The strategy for this placement group- `host_spread`: place on different compute hosts- `power_spread`: place on compute hosts that use different power sourcesThe enumerated values for this property may expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the placement group on which the unexpected strategy was encountered.",
						},
						isPlacementGroupTags: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         flex.ResourceIBMVPCHash,
							Description: "List of tags",
						},
						isPlacementGroupAccessTags: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         flex.ResourceIBMVPCHash,
							Description: "List of access management tags",
						},
					},
				},
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of resources across all pages.",
			},
		},
	}
}

func dataSourceIbmIsPlacementGroupsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	listPlacementGroupsOptions := &vpcv1.ListPlacementGroupsOptions{}
	start := ""
	allrecs := []vpcv1.PlacementGroup{}
	for {
		if start != "" {
			listPlacementGroupsOptions.Start = &start
		}
		placementGroupCollection, response, err := vpcClient.ListPlacementGroupsWithContext(context, listPlacementGroupsOptions)
		if err != nil {
			log.Printf("[DEBUG] ListPlacementGroupsWithContext failed %s\n%s", err, response)
			return diag.FromErr(err)
		}
		start = flex.GetNext(placementGroupCollection.Next)
		allrecs = append(allrecs, placementGroupCollection.PlacementGroups...)
		if start == "" {
			break
		}
	}

	d.SetId(dataSourceIbmIsPlacementGroupsID(d))
	err = d.Set("placement_groups", dataSourcePlacementGroupCollectionFlattenPlacementGroups(meta, allrecs))
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting placement_groups %s", err))
	}
	if err = d.Set("total_count", len(allrecs)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting total_count: %s", err))
	}
	return nil
}

// dataSourceIbmIsPlacementGroupsID returns a reasonable ID for the list.
func dataSourceIbmIsPlacementGroupsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourcePlacementGroupCollectionFlattenPlacementGroups(meta interface{}, result []vpcv1.PlacementGroup) (placementGroups []map[string]interface{}) {
	placementGroups = make([]map[string]interface{}, 0)
	for _, placementGroupsItem := range result {
		placementGroups = append(placementGroups, dataSourcePlacementGroupCollectionPlacementGroupsToMap(meta, placementGroupsItem))
	}

	return placementGroups
}

func dataSourcePlacementGroupCollectionPlacementGroupsToMap(meta interface{}, placementGroupsItem vpcv1.PlacementGroup) (placementGroupsMap map[string]interface{}) {
	placementGroupsMap = map[string]interface{}{}

	if placementGroupsItem.CreatedAt != nil {
		placementGroupsMap["created_at"] = placementGroupsItem.CreatedAt.String()
	}
	if placementGroupsItem.CRN != nil {
		placementGroupsMap["crn"] = placementGroupsItem.CRN
	}
	if placementGroupsItem.Href != nil {
		placementGroupsMap["href"] = placementGroupsItem.Href
	}
	if placementGroupsItem.ID != nil {
		placementGroupsMap["id"] = placementGroupsItem.ID
	}
	if placementGroupsItem.LifecycleState != nil {
		placementGroupsMap["lifecycle_state"] = placementGroupsItem.LifecycleState
	}
	if placementGroupsItem.Name != nil {
		placementGroupsMap["name"] = placementGroupsItem.Name
	}
	if placementGroupsItem.ResourceGroup != nil {
		resourceGroupList := []map[string]interface{}{}
		resourceGroupMap := dataSourcePlacementGroupCollectionPlacementGroupsResourceGroupToMap(*placementGroupsItem.ResourceGroup)
		resourceGroupList = append(resourceGroupList, resourceGroupMap)
		placementGroupsMap["resource_group"] = resourceGroupList
	}
	if placementGroupsItem.ResourceType != nil {
		placementGroupsMap["resource_type"] = placementGroupsItem.ResourceType
	}
	if placementGroupsItem.Strategy != nil {
		placementGroupsMap["strategy"] = placementGroupsItem.Strategy
	}
	tags, err := flex.GetGlobalTagsUsingCRN(meta, *placementGroupsItem.CRN, "", isUserTagType)
	if err != nil {
		log.Printf(
			"An error getting placement group (%s) tags : %s", *placementGroupsItem.ID, err)
	}

	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *placementGroupsItem.CRN, "", isAccessTagType)
	if err != nil {
		log.Printf(
			"Error getting placement group (%s) access tags: %s", *placementGroupsItem.ID, err)
	}

	placementGroupsMap[isPlacementGroupTags] = tags
	placementGroupsMap[isPlacementGroupAccessTags] = accesstags
	return placementGroupsMap
}

func dataSourcePlacementGroupCollectionPlacementGroupsResourceGroupToMap(resourceGroupItem vpcv1.ResourceGroupReference) (resourceGroupMap map[string]interface{}) {
	resourceGroupMap = map[string]interface{}{}

	if resourceGroupItem.Href != nil {
		resourceGroupMap["href"] = resourceGroupItem.Href
	}
	if resourceGroupItem.ID != nil {
		resourceGroupMap["id"] = resourceGroupItem.ID
	}
	if resourceGroupItem.Name != nil {
		resourceGroupMap["name"] = resourceGroupItem.Name
	}

	return resourceGroupMap
}
