// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func dataSourceIbmIsInstanceDisk() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIbmIsInstanceDiskRead,

		Schema: map[string]*schema.Schema{
			"instance": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The instance identifier.",
			},
			"disk": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The instance disk identifier.",
			},
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
	}
}

func dataSourceIbmIsInstanceDiskRead(d *schema.ResourceData, meta interface{}) error {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return err
	}

	getInstanceDiskOptions := &vpcv1.GetInstanceDiskOptions{}

	getInstanceDiskOptions.SetInstanceID(d.Get("instance").(string))
	getInstanceDiskOptions.SetID(d.Get("disk").(string))

	instanceDisk, response, err := vpcClient.GetInstanceDiskWithContext(context.TODO(), getInstanceDiskOptions)
	if err != nil {
		log.Printf("[DEBUG] GetInstanceDiskWithContext failed %s\n%s", err, response)
		return err
	}

	d.SetId(*instanceDisk.ID)
	if err = d.Set("created_at", instanceDisk.CreatedAt.String()); err != nil {
		return fmt.Errorf("Error setting created_at: %s", err)
	}
	if err = d.Set("href", instanceDisk.Href); err != nil {
		return fmt.Errorf("Error setting href: %s", err)
	}
	if err = d.Set("interface_type", instanceDisk.InterfaceType); err != nil {
		return fmt.Errorf("Error setting interface_type: %s", err)
	}
	if err = d.Set("name", instanceDisk.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err = d.Set("resource_type", instanceDisk.ResourceType); err != nil {
		return fmt.Errorf("Error setting resource_type: %s", err)
	}
	if err = d.Set("size", instanceDisk.Size); err != nil {
		return fmt.Errorf("Error setting size: %s", err)
	}

	return nil
}
