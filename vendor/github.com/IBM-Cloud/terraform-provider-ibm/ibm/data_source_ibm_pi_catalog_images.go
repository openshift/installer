// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
)

/*
Datasource to get the list of images that are available when a power instance is created

*/
func dataSourceIBMPICatalogImages() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceIBMPICatalogImagesRead,
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

func dataSourceIBMPICatalogImagesRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := meta.(ClientSession).IBMPISession()

	if err != nil {
		return err
	}
	sap := false
	powerinstanceid := d.Get(helpers.PICloudInstanceId).(string)
	if v, ok := d.GetOk("sap"); ok {
		sap = v.(bool)
	}

	imageC := instance.NewIBMPIImageClient(sess, powerinstanceid)
	result, err := imageC.GetSAPImages(powerinstanceid, sap)
	if err != nil {
		return err
	}
	imageData := result.Images
	images := make([]map[string]interface{}, 0)
	for _, i := range imageData {
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
			if &s.ImageType != nil {
				image["image_type"] = s.ImageType
			}
			if &s.ContainerFormat != nil {
				image["container_format"] = s.ContainerFormat
			}
			if &s.DiskFormat != nil {
				image["disk_format"] = s.DiskFormat
			}
			if &s.OperatingSystem != nil {
				image["operating_system"] = s.OperatingSystem
			}
			if &s.HypervisorType != nil {
				image["hypervisor_type"] = s.HypervisorType
			}
			if &s.Architecture != nil {
				image["architecture"] = s.Architecture
			}
			if &s.Endianness != nil {
				image["endianness"] = s.Endianness
			}
		}
		images = append(images, image)
	}
	d.SetId(time.Now().UTC().String())
	d.Set("images", images)
	d.Set(helpers.PICloudInstanceId, powerinstanceid)
	return nil

}
