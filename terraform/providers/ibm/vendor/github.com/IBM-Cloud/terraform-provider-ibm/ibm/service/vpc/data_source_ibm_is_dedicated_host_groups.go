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

func DataSourceIbmIsDedicatedHostGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmIsDedicatedHostGroupsRead,

		Schema: map[string]*schema.Schema{
			"resource_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique identifier of the resource group this dedicated host group belongs to",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The zone name this dedicated host group is in",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the dedicated host group",
			},
			"host_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of dedicated host groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"class": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The dedicated host profile class for hosts in this group.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the dedicated host group was created.",
						},
						"crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this dedicated host group.",
						},
						"dedicated_hosts": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The dedicated hosts that are in this dedicated host group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this dedicated host.",
									},
									"deleted": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									"href": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this dedicated host.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this dedicated host.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique user-defined name for this dedicated host. If unspecified, the name will be a hyphenated list of randomly-selected words.",
									},
									"resource_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of resource referenced.",
									},
								},
							},
						},
						"family": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The dedicated host profile family for hosts in this group.",
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this dedicated host group.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this dedicated host group.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique user-defined name for this dedicated host group. If unspecified, the name will be a hyphenated list of randomly-selected words.",
						},
						"resource_group": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource group for this dedicated host group.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of resource referenced.",
						},
						"supported_instance_profiles": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Array of instance profiles that can be used by instances placed on this dedicated host group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this virtual server instance profile.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The globally unique name for this virtual server instance profile.",
									},
								},
							},
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique name of the zone this dedicated host group resides in.",
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

func dataSourceIbmIsDedicatedHostGroupsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}
	listDedicatedHostGroupsOptions := &vpcv1.ListDedicatedHostGroupsOptions{}

	if resgroupintf, ok := d.GetOk("resource_group"); ok {
		resGroup := resgroupintf.(string)
		listDedicatedHostGroupsOptions.ResourceGroupID = &resGroup
	}
	if zoneintf, ok := d.GetOk("zone"); ok {
		zoneName := zoneintf.(string)
		listDedicatedHostGroupsOptions.ZoneName = &zoneName
	}
	if nameintf, ok := d.GetOk("name"); ok {
		name := nameintf.(string)
		listDedicatedHostGroupsOptions.Name = &name
	}

	start := ""
	allrecs := []vpcv1.DedicatedHostGroup{}
	for {
		if start != "" {
			listDedicatedHostGroupsOptions.Start = &start
		}
		listDedicatedHostGroupsOptions, response, err := vpcClient.ListDedicatedHostGroupsWithContext(context, listDedicatedHostGroupsOptions)
		if err != nil {
			log.Printf("[DEBUG] ListDedicatedHostGroupsWithContext failed %s\n%s", err, response)
			return diag.FromErr(err)
		}
		start = flex.GetNext(listDedicatedHostGroupsOptions.Next)
		allrecs = append(allrecs, listDedicatedHostGroupsOptions.Groups...)
		if start == "" {
			break
		}
	}

	if len(allrecs) != 0 {

		d.SetId(dataSourceIbmIsDedicatedHostGroupsID(d))
		err = d.Set("host_groups", dataSourceDedicatedHostGroupCollectionFlattenGroups(allrecs))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting groups %s", err))
		}

		if err = d.Set("total_count", len(allrecs)); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting total_count: %s", err))
		}

	}
	return nil
}

// dataSourceIbmIsDedicatedHostGroupsID returns a reasonable ID for the list.
func dataSourceIbmIsDedicatedHostGroupsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceDedicatedHostGroupCollectionFirstToMap(firstItem vpcv1.DedicatedHostGroupCollectionFirst) (firstMap map[string]interface{}) {
	firstMap = map[string]interface{}{}

	if firstItem.Href != nil {
		firstMap["href"] = firstItem.Href
	}

	return firstMap
}

func dataSourceDedicatedHostGroupCollectionFlattenGroups(result []vpcv1.DedicatedHostGroup) (groups []map[string]interface{}) {
	for _, groupsItem := range result {
		groups = append(groups, dataSourceDedicatedHostGroupCollectionGroupsToMap(groupsItem))
	}

	return groups
}

func dataSourceDedicatedHostGroupCollectionGroupsToMap(groupsItem vpcv1.DedicatedHostGroup) (groupsMap map[string]interface{}) {
	groupsMap = map[string]interface{}{}

	if groupsItem.Class != nil {
		groupsMap["class"] = groupsItem.Class
	}
	if groupsItem.CreatedAt != nil {
		groupsMap["created_at"] = groupsItem.CreatedAt.String()
	}
	if groupsItem.CRN != nil {
		groupsMap["crn"] = groupsItem.CRN
	}
	if groupsItem.DedicatedHosts != nil {
		dedicatedHostsList := []map[string]interface{}{}
		for _, dedicatedHostsItem := range groupsItem.DedicatedHosts {
			dedicatedHostsList = append(dedicatedHostsList, dataSourceDedicatedHostGroupCollectionGroupsDedicatedHostsToMap(dedicatedHostsItem))
		}
		groupsMap["dedicated_hosts"] = dedicatedHostsList
	}
	if groupsItem.Family != nil {
		groupsMap["family"] = groupsItem.Family
	}
	if groupsItem.Href != nil {
		groupsMap["href"] = groupsItem.Href
	}
	if groupsItem.ID != nil {
		groupsMap["id"] = groupsItem.ID
	}
	if groupsItem.Name != nil {
		groupsMap["name"] = groupsItem.Name
	}
	if groupsItem.ResourceGroup != nil {
		groupsMap["resource_group"] = *groupsItem.ResourceGroup.ID
	}
	if groupsItem.ResourceType != nil {
		groupsMap["resource_type"] = groupsItem.ResourceType
	}
	if groupsItem.SupportedInstanceProfiles != nil {
		supportedInstanceProfilesList := []map[string]interface{}{}
		for _, supportedInstanceProfilesItem := range groupsItem.SupportedInstanceProfiles {
			supportedInstanceProfilesList = append(supportedInstanceProfilesList, dataSourceDedicatedHostGroupCollectionGroupsSupportedInstanceProfilesToMap(supportedInstanceProfilesItem))
		}
		groupsMap["supported_instance_profiles"] = supportedInstanceProfilesList
	}
	if groupsItem.Zone != nil {
		groupsMap["zone"] = *groupsItem.Zone.Name
	}

	return groupsMap
}

func dataSourceDedicatedHostGroupCollectionDedicatedHostsDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceDedicatedHostGroupCollectionGroupsDedicatedHostsToMap(dedicatedHostsItem vpcv1.DedicatedHostReference) (dedicatedHostsMap map[string]interface{}) {
	dedicatedHostsMap = map[string]interface{}{}

	if dedicatedHostsItem.CRN != nil {
		dedicatedHostsMap["crn"] = dedicatedHostsItem.CRN
	}
	if dedicatedHostsItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceDedicatedHostGroupCollectionDedicatedHostsDeletedToMap(*dedicatedHostsItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		dedicatedHostsMap["deleted"] = deletedList
	}
	if dedicatedHostsItem.Href != nil {
		dedicatedHostsMap["href"] = dedicatedHostsItem.Href
	}
	if dedicatedHostsItem.ID != nil {
		dedicatedHostsMap["id"] = dedicatedHostsItem.ID
	}
	if dedicatedHostsItem.Name != nil {
		dedicatedHostsMap["name"] = dedicatedHostsItem.Name
	}
	if dedicatedHostsItem.ResourceType != nil {
		dedicatedHostsMap["resource_type"] = dedicatedHostsItem.ResourceType
	}

	return dedicatedHostsMap
}

func dataSourceDedicatedHostGroupCollectionGroupsSupportedInstanceProfilesToMap(supportedInstanceProfilesItem vpcv1.InstanceProfileReference) (supportedInstanceProfilesMap map[string]interface{}) {
	supportedInstanceProfilesMap = map[string]interface{}{}

	if supportedInstanceProfilesItem.Href != nil {
		supportedInstanceProfilesMap["href"] = supportedInstanceProfilesItem.Href
	}
	if supportedInstanceProfilesItem.Name != nil {
		supportedInstanceProfilesMap["name"] = supportedInstanceProfilesItem.Name
	}

	return supportedInstanceProfilesMap
}

func dataSourceDedicatedHostGroupCollectionFlattenFirst(result vpcv1.DedicatedHostGroupCollectionFirst) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceDedicatedHostGroupCollectionFirstToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceDedicatedHostGroupCollectionFlattenNext(result vpcv1.DedicatedHostGroupCollectionNext) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceDedicatedHostGroupCollectionNextToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceDedicatedHostGroupCollectionNextToMap(nextItem vpcv1.DedicatedHostGroupCollectionNext) (nextMap map[string]interface{}) {
	nextMap = map[string]interface{}{}

	if nextItem.Href != nil {
		nextMap["href"] = nextItem.Href
	}

	return nextMap
}
