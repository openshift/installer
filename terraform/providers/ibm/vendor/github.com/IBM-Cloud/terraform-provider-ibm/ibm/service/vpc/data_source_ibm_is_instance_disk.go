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

func DataSourceIbmIsInstanceDisk() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmIsInstanceDiskRead,

		Schema: map[string]*schema.Schema{
			"instance": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The instance identifier.",
			},
			"disk": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The instance disk identifier.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the disk was created.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this instance disk.",
			},
			"interface_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The disk interface used for attaching the disk.The enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected property value was encountered.",
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
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The size of the disk in GB (gigabytes).",
			},
		},
	}
}

func dataSourceIbmIsInstanceDiskRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getInstanceDiskOptions := &vpcv1.GetInstanceDiskOptions{}

	getInstanceDiskOptions.SetInstanceID(d.Get("instance").(string))
	getInstanceDiskOptions.SetID(d.Get("disk").(string))

	instanceDisk, response, err := vpcClient.GetInstanceDiskWithContext(context, getInstanceDiskOptions)
	if err != nil {
		log.Printf("[DEBUG] GetInstanceDiskWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	d.SetId(*instanceDisk.ID)
	if err = d.Set("created_at", instanceDisk.CreatedAt.String()); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
	}
	if err = d.Set("href", instanceDisk.Href); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting href: %s", err))
	}
	if err = d.Set("interface_type", instanceDisk.InterfaceType); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting interface_type: %s", err))
	}
	if err = d.Set("name", instanceDisk.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}
	if err = d.Set("resource_type", instanceDisk.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting resource_type: %s", err))
	}
	if err = d.Set("size", instanceDisk.Size); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting size: %s", err))
	}

	return nil
}
