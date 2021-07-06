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

func dataSourceIbmIsInstanceDisks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIbmIsInstanceDisksRead,

		Schema: map[string]*schema.Schema{
			"instance": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The instance identifier.",
			},
			isInstanceDisks: &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of the instance's disks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the disk was created.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this instance disk.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this instance disk.",
						},
						"interface_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The disk interface used for attaching the disk.The enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected property value was encountered.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this disk.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						"size": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The size of the disk in GB (gigabytes).",
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmIsInstanceDisksRead(d *schema.ResourceData, meta interface{}) error {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return err
	}
	listInstanceDisksOptions := &vpcv1.ListInstanceDisksOptions{}

	listInstanceDisksOptions.SetInstanceID(d.Get("instance").(string))

	instanceDiskCollection, response, err := vpcClient.ListInstanceDisksWithContext(context.TODO(), listInstanceDisksOptions)
	if err != nil {
		log.Printf("[DEBUG] ListInstanceDisksWithContext failed %s\n%s", err, response)
		return err
	}

	d.SetId(dataSourceIbmIsInstanceDisksID(d))

	if instanceDiskCollection.Disks != nil {
		err = d.Set(isInstanceDisks, dataSourceInstanceDiskCollectionFlattenDisks(instanceDiskCollection.Disks))
		if err != nil {
			return fmt.Errorf("Error setting disks %s", err)
		}
	}

	return nil
}

// dataSourceIbmIsInstanceDisksID returns a reasonable ID for the list.
func dataSourceIbmIsInstanceDisksID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceInstanceDiskCollectionFlattenDisks(result []vpcv1.InstanceDisk) (disks []map[string]interface{}) {
	for _, disksItem := range result {
		disks = append(disks, dataSourceInstanceDiskCollectionDisksToMap(disksItem))
	}

	return disks
}

func dataSourceInstanceDiskCollectionDisksToMap(disksItem vpcv1.InstanceDisk) (disksMap map[string]interface{}) {
	disksMap = map[string]interface{}{}

	if disksItem.CreatedAt != nil {
		disksMap["created_at"] = disksItem.CreatedAt.String()
	}
	if disksItem.Href != nil {
		disksMap["href"] = disksItem.Href
	}
	if disksItem.ID != nil {
		disksMap["id"] = disksItem.ID
	}
	if disksItem.InterfaceType != nil {
		disksMap["interface_type"] = disksItem.InterfaceType
	}
	if disksItem.Name != nil {
		disksMap["name"] = disksItem.Name
	}
	if disksItem.ResourceType != nil {
		disksMap["resource_type"] = disksItem.ResourceType
	}
	if disksItem.Size != nil {
		disksMap["size"] = disksItem.Size
	}

	return disksMap
}
