// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
)

/*
Datasource to get the list of images that are available when a power instance is created
*/
func DataSourceIBMPICatalogImages() *schema.Resource {

	return &schema.Resource{
		ReadContext: dataSourceIBMPICatalogImagesRead,
		Schema: map[string]*schema.Schema{

			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"sap": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"vtl": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"images": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_pool": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_update_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"container_format": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_format": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operating_system": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hypervisor_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"architecture": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endianness": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPICatalogImagesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)
	includeSAP := false
	if s, ok := d.GetOk("sap"); ok {
		includeSAP = s.(bool)
	}
	includeVTL := false
	if v, ok := d.GetOk("vtl"); ok {
		includeVTL = v.(bool)
	}
	imageC := instance.NewIBMPIImageClient(ctx, sess, cloudInstanceID)
	stockImages, err := imageC.GetAllStockImages(includeSAP, includeVTL)
	if err != nil {
		return diag.FromErr(err)
	}

	images := make([]map[string]interface{}, 0)
	for _, i := range stockImages.Images {
		image := make(map[string]interface{})
		image["image_id"] = *i.ImageID
		image["name"] = *i.Name
		if i.State != nil {
			image["state"] = *i.State
		}
		if i.Description != nil {
			image["description"] = *i.Description
		}
		if i.StorageType != nil {
			image["storage_type"] = *i.StorageType
		}
		if i.StoragePool != nil {
			image["storage_pool"] = *i.StoragePool
		}
		if i.CreationDate != nil {
			image["creation_date"] = i.CreationDate.String()
		}
		if i.LastUpdateDate != nil {
			image["last_update_date"] = i.LastUpdateDate.String()
		}
		if i.Href != nil {
			image["href"] = *i.Href
		}
		if i.Specifications != nil {
			s := i.Specifications
			if s.ImageType != "" {
				image["image_type"] = s.ImageType
			}
			if s.ContainerFormat != "" {
				image["container_format"] = s.ContainerFormat
			}
			if s.DiskFormat != "" {
				image["disk_format"] = s.DiskFormat
			}
			if s.OperatingSystem != "" {
				image["operating_system"] = s.OperatingSystem
			}
			if s.HypervisorType != "" {
				image["hypervisor_type"] = s.HypervisorType
			}
			if s.Architecture != "" {
				image["architecture"] = s.Architecture
			}
			if s.Endianness != "" {
				image["endianness"] = s.Endianness
			}
		}
		images = append(images, image)
	}
	d.SetId(time.Now().UTC().String())
	d.Set("images", images)
	return nil

}
