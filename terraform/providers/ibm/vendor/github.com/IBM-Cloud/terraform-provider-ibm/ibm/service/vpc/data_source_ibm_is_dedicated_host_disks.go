// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIbmIsDedicatedHostDisks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmIsDedicatedHostDisksRead,

		Schema: map[string]*schema.Schema{
			"dedicated_host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The dedicated host identifier.",
			},
			"disks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of the dedicated host's disks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"available": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The remaining space left for instance placement in GB (gigabytes).",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the disk was created.",
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this disk.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this disk.",
						},
						"instance_disks": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Instance disks that are on this dedicated host disk.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
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
										Description: "The URL for this instance disk.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this instance disk.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this disk.",
									},
									"resource_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"interface_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The disk interface used for attaching the diskThe enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected property value was encountered.",
						},
						"lifecycle_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lifecycle state of this dedicated host disk.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined or system-provided name for this disk.",
						},
						"provisionable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether this dedicated host disk is available for instance disk creation.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of resource referenced.",
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The size of the disk in GB (gigabytes).",
						},
						"supported_instance_interface_types": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The instance disk interfaces supported for this dedicated host disk.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmIsDedicatedHostDisksRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	listDedicatedHostDisksOptions := &vpcv1.ListDedicatedHostDisksOptions{}

	listDedicatedHostDisksOptions.SetDedicatedHostID(d.Get("dedicated_host").(string))

	dedicatedHostDiskCollection, response, err := vpcClient.ListDedicatedHostDisksWithContext(context, listDedicatedHostDisksOptions)
	if err != nil {
		log.Printf("[DEBUG] ListDedicatedHostDisksWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	d.SetId(dataSourceIbmIsDedicatedHostDisksID(d))

	if dedicatedHostDiskCollection.Disks != nil {
		err = d.Set("disks", dataSourceDedicatedHostDiskCollectionFlattenDisks(dedicatedHostDiskCollection.Disks))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting disks %s", err))
		}
	}

	return nil
}

// dataSourceIbmIsDedicatedHostDisksID returns a reasonable ID for the list.
func dataSourceIbmIsDedicatedHostDisksID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceDedicatedHostDiskCollectionFlattenDisks(result []vpcv1.DedicatedHostDisk) (disks []map[string]interface{}) {
	for _, disksItem := range result {
		disks = append(disks, dataSourceDedicatedHostDiskCollectionDisksToMap(disksItem))
	}

	return disks
}

func dataSourceDedicatedHostDiskCollectionDisksToMap(disksItem vpcv1.DedicatedHostDisk) (disksMap map[string]interface{}) {
	disksMap = map[string]interface{}{}

	if disksItem.Available != nil {
		disksMap["available"] = disksItem.Available
	}
	if disksItem.CreatedAt != nil {
		disksMap["created_at"] = disksItem.CreatedAt.String()
	}
	if disksItem.Href != nil {
		disksMap["href"] = disksItem.Href
	}
	if disksItem.ID != nil {
		disksMap["id"] = disksItem.ID
	}
	if disksItem.InstanceDisks != nil {
		instanceDisksList := []map[string]interface{}{}
		for _, instanceDisksItem := range disksItem.InstanceDisks {
			instanceDisksList = append(instanceDisksList, dataSourceDedicatedHostDiskCollectionDisksInstanceDisksToMap(instanceDisksItem))
		}
		disksMap["instance_disks"] = instanceDisksList
	}
	if disksItem.InterfaceType != nil {
		disksMap["interface_type"] = disksItem.InterfaceType
	}
	if disksItem.LifecycleState != nil {
		disksMap["lifecycle_state"] = disksItem.LifecycleState
	}
	if disksItem.Name != nil {
		disksMap["name"] = disksItem.Name
	}
	if disksItem.Provisionable != nil {
		disksMap["provisionable"] = disksItem.Provisionable
	}
	if disksItem.ResourceType != nil {
		disksMap["resource_type"] = disksItem.ResourceType
	}
	if disksItem.Size != nil {
		disksMap["size"] = disksItem.Size
	}
	if disksItem.SupportedInstanceInterfaceTypes != nil {
		disksMap["supported_instance_interface_types"] = disksItem.SupportedInstanceInterfaceTypes
	}

	return disksMap
}

func dataSourceDedicatedHostDiskCollectionDisksInstanceDisksToMap(instanceDisksItem vpcv1.InstanceDiskReference) (instanceDisksMap map[string]interface{}) {
	instanceDisksMap = map[string]interface{}{}

	if instanceDisksItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceDedicatedHostDiskCollectionInstanceDisksDeletedToMap(*instanceDisksItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		instanceDisksMap["deleted"] = deletedList
	}
	if instanceDisksItem.Href != nil {
		instanceDisksMap["href"] = instanceDisksItem.Href
	}
	if instanceDisksItem.ID != nil {
		instanceDisksMap["id"] = instanceDisksItem.ID
	}
	if instanceDisksItem.Name != nil {
		instanceDisksMap["name"] = instanceDisksItem.Name
	}
	if instanceDisksItem.ResourceType != nil {
		instanceDisksMap["resource_type"] = instanceDisksItem.ResourceType
	}

	return instanceDisksMap
}

func dataSourceDedicatedHostDiskCollectionInstanceDisksDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}
