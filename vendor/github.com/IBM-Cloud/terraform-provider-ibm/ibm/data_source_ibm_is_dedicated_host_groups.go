// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func dataSourceIbmIsDedicatedHostGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIbmIsDedicatedHostGroupsRead,

		Schema: map[string]*schema.Schema{
			"first": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A link to the first page of resources.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for a page of resources.",
						},
					},
				},
			},
			"limit": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum number of resources that can be returned by the request.",
			},
			"next": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A link to the next page of resources. This property is present for all pagesexcept the last page.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for a page of resources.",
						},
					},
				},
			},
			"host_groups": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of dedicated host groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"class": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The dedicated host profile class for hosts in this group.",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the dedicated host group was created.",
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this dedicated host group.",
						},
						"dedicated_hosts": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The dedicated hosts that are in this dedicated host group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this dedicated host.",
									},
									"deleted": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
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
										Description: "The URL for this dedicated host.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this dedicated host.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique user-defined name for this dedicated host. If unspecified, the name will be a hyphenated list of randomly-selected words.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of resource referenced.",
									},
								},
							},
						},
						"family": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The dedicated host profile family for hosts in this group.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this dedicated host group.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this dedicated host group.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique user-defined name for this dedicated host group. If unspecified, the name will be a hyphenated list of randomly-selected words.",
						},
						"resource_group": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource group for this dedicated host group.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of resource referenced.",
						},
						"supported_instance_profiles": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Array of instance profiles that can be used by instances placed on this dedicated host group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this virtual server instance profile.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The globally unique name for this virtual server instance profile.",
									},
								},
							},
						},
						"zone": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique name of the zone this dedicated host group resides in.",
						},
					},
				},
			},
			"total_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of resources across all pages.",
			},
		},
	}
}

func dataSourceIbmIsDedicatedHostGroupsRead(d *schema.ResourceData, meta interface{}) error {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return err
	}
	listDedicatedHostGroupsOptions := &vpcv1.ListDedicatedHostGroupsOptions{}

	dedicatedHostGroupCollection, response, err := vpcClient.ListDedicatedHostGroupsWithContext(context.TODO(), listDedicatedHostGroupsOptions)
	if err != nil {
		log.Printf("[DEBUG] ListDedicatedHostGroupsWithContext failed %s\n%s", err, response)
		return err
	}

	if len(dedicatedHostGroupCollection.Groups) != 0 {

		d.SetId(dataSourceIbmIsDedicatedHostGroupsID(d))

		if dedicatedHostGroupCollection.First != nil {
			err = d.Set("first", dataSourceDedicatedHostGroupCollectionFlattenFirst(*dedicatedHostGroupCollection.First))
			if err != nil {
				return fmt.Errorf("Error setting first %s", err)
			}
		}

		if dedicatedHostGroupCollection.Groups != nil {
			err = d.Set("host_groups", dataSourceDedicatedHostGroupCollectionFlattenGroups(dedicatedHostGroupCollection.Groups))
			if err != nil {
				return fmt.Errorf("Error setting groups %s", err)
			}
		}
		if err = d.Set("limit", dedicatedHostGroupCollection.Limit); err != nil {
			return fmt.Errorf("Error setting limit: %s", err)
		}

		if dedicatedHostGroupCollection.Next != nil {
			err = d.Set("next", dataSourceDedicatedHostGroupCollectionFlattenNext(*dedicatedHostGroupCollection.Next))
			if err != nil {
				return fmt.Errorf("Error setting next %s", err)
			}
		}
		if err = d.Set("total_count", dedicatedHostGroupCollection.TotalCount); err != nil {
			return fmt.Errorf("Error setting total_count: %s", err)
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

func dataSourceDedicatedHostGroupCollectionDedicatedHostsDeletedToMap(deletedItem vpcv1.DedicatedHostReferenceDeleted) (deletedMap map[string]interface{}) {
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
