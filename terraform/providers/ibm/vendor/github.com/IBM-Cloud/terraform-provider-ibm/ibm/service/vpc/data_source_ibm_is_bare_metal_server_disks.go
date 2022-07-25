// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isBareMetalServerID               = "bare_metal_server"
	isBareMetalServerDiskHref         = "href"
	isBareMetalServerDiskResourceType = "resource_type"
)

func DataSourceIBMIsBareMetalServerDisks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISBareMetalServerDisksRead,

		Schema: map[string]*schema.Schema{
			isBareMetalServerID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The bare metal server identifier",
			},

			//disks

			isBareMetalServerDisks: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of bare metal server disks. Disk is a block device that is locally attached to the physical server. By default, the listed disks are sorted by their created_at property values, with the newest disk first.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
				},
			},
		},
	}
}

func dataSourceIBMISBareMetalServerDisksRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	bareMetalServerID := d.Get(isBareMetalServerID).(string)
	sess, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}
	options := &vpcv1.ListBareMetalServerDisksOptions{
		BareMetalServerID: &bareMetalServerID,
	}

	diskCollection, response, err := sess.ListBareMetalServerDisksWithContext(context, options)
	disks := diskCollection.Disks
	if err != nil || disks == nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting Bare Metal Server (%s) disks: %s\n%s", bareMetalServerID, err, response))
	}
	disksInfo := make([]map[string]interface{}, 0)
	for _, disk := range disks {
		l := map[string]interface{}{
			isBareMetalServerDiskHref:          disk.Href,
			isBareMetalServerDiskID:            disk.ID,
			isBareMetalServerDiskInterfaceType: disk.InterfaceType,
			isBareMetalServerDiskName:          disk.Name,
			isBareMetalServerDiskResourceType:  disk.ResourceType,
			isBareMetalServerDiskSize:          disk.Size,
		}
		disksInfo = append(disksInfo, l)
	}
	d.SetId(dataSourceIBMISBMSDisksID(d))
	d.Set(isBareMetalServerDisks, disksInfo)
	return nil
}

// dataSourceIBMISBMSProfilesID returns a reasonable ID for a Bare Metal Server Disks list.
func dataSourceIBMISBMSDisksID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
