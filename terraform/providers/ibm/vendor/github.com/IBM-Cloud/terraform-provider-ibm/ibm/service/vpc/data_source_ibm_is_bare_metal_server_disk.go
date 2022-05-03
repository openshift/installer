// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isBareMetalServerDisk = "disk"
)

func DataSourceIBMIsBareMetalServerDisk() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISBareMetalServerDiskRead,

		Schema: map[string]*schema.Schema{
			isBareMetalServerID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The bare metal server identifier",
			},

			isBareMetalServerDisk: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The bare metal server disk identifier",
			},
			//disks

			isBareMetalServerDiskHref: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this bare metal server disk",
			},
			isBareMetalServerDiskID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this bare metal server disk",
			},
			isBareMetalServerDiskInterfaceType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The disk interface used for attaching the disk. Supported values are [ nvme, sata ]",
			},
			isBareMetalServerDiskName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user-defined name for this disk",
			},
			isBareMetalServerDiskResourceType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type",
			},
			isBareMetalServerDiskSize: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The size of the disk in GB (gigabytes)",
			},
		},
	}
}

func dataSourceIBMISBareMetalServerDiskRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	bareMetalServerID := d.Get(isBareMetalServerID).(string)
	bareMetalServerDiskID := d.Get(isBareMetalServerDisk).(string)
	sess, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}
	options := &vpcv1.GetBareMetalServerDiskOptions{
		BareMetalServerID: &bareMetalServerID,
		ID:                &bareMetalServerDiskID,
	}

	disk, response, err := sess.GetBareMetalServerDiskWithContext(context, options)
	if err != nil || disk == nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting Bare Metal Server (%s) disk (%s): %s\n%s", bareMetalServerID, bareMetalServerDiskID, err, response))
	}
	d.SetId(*disk.ID)
	d.Set(isBareMetalServerDiskHref, *disk.Href)
	d.Set(isBareMetalServerDiskInterfaceType, *disk.InterfaceType)
	d.Set(isBareMetalServerDiskName, *disk.Name)
	d.Set(isBareMetalServerDiskResourceType, *disk.ResourceType)
	d.Set(isBareMetalServerDiskSize, *disk.Size)

	return nil
}
