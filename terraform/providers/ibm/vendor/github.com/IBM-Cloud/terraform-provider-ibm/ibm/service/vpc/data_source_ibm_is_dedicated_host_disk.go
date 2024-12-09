// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIbmIsDedicatedHostDisk() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmIsDedicatedHostDiskRead,

		Schema: map[string]*schema.Schema{
			"dedicated_host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The dedicated host identifier.",
			},
			"disk": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The dedicated host disk identifier.",
			},
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
	}
}

func dataSourceIbmIsDedicatedHostDiskRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getDedicatedHostDiskOptions := &vpcv1.GetDedicatedHostDiskOptions{}

	getDedicatedHostDiskOptions.SetDedicatedHostID(d.Get("dedicated_host").(string))
	getDedicatedHostDiskOptions.SetID(d.Get("disk").(string))

	dedicatedHostDisk, response, err := vpcClient.GetDedicatedHostDiskWithContext(context, getDedicatedHostDiskOptions)
	if err != nil {
		log.Printf("[DEBUG] GetDedicatedHostDiskWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	d.SetId(*dedicatedHostDisk.ID)
	if err = d.Set("available", dedicatedHostDisk.Available); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting available: %s", err))
	}
	if err = d.Set("created_at", dedicatedHostDisk.CreatedAt.String()); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
	}
	if err = d.Set("href", dedicatedHostDisk.Href); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting href: %s", err))
	}

	if dedicatedHostDisk.InstanceDisks != nil {
		err = d.Set("instance_disks", dataSourceDedicatedHostDiskFlattenInstanceDisks(dedicatedHostDisk.InstanceDisks))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting instance_disks %s", err))
		}
	}
	if err = d.Set("interface_type", dedicatedHostDisk.InterfaceType); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting interface_type: %s", err))
	}
	if dedicatedHostDisk.LifecycleState != nil {
		if err = d.Set("lifecycle_state", dedicatedHostDisk.LifecycleState); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting lifecycle_state: %s", err))
		}
	}
	if err = d.Set("name", dedicatedHostDisk.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}
	if err = d.Set("provisionable", dedicatedHostDisk.Provisionable); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting provisionable: %s", err))
	}
	if err = d.Set("resource_type", dedicatedHostDisk.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting resource_type: %s", err))
	}
	if err = d.Set("size", dedicatedHostDisk.Size); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting size: %s", err))
	}
	if err = d.Set("supported_instance_interface_types", dedicatedHostDisk.SupportedInstanceInterfaceTypes); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting supported_instance_interface_types: %s", err))
	}

	return nil
}

func dataSourceDedicatedHostDiskFlattenInstanceDisks(result []vpcv1.InstanceDiskReference) (instanceDisks []map[string]interface{}) {
	for _, instanceDisksItem := range result {
		instanceDisks = append(instanceDisks, dataSourceDedicatedHostDiskInstanceDisksToMap(instanceDisksItem))
	}

	return instanceDisks
}

func dataSourceDedicatedHostDiskInstanceDisksToMap(instanceDisksItem vpcv1.InstanceDiskReference) (instanceDisksMap map[string]interface{}) {
	instanceDisksMap = map[string]interface{}{}

	if instanceDisksItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceDedicatedHostDiskInstanceDisksDeletedToMap(*instanceDisksItem.Deleted)
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

func dataSourceDedicatedHostDiskInstanceDisksDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}
