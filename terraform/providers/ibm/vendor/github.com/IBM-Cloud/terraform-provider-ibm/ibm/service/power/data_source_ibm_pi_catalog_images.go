// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"time"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Datasource to list images that are available when a power instance is created
func DataSourceIBMPICatalogImages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPICatalogImagesRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_SAP: {
				Description: "Set true to include SAP images. The default value is false.",
				Optional:    true,
				Type:        schema.TypeBool,
			},
			Arg_VTL: {
				Description: "Set true to include VTL images. The default value is false.",
				Optional:    true,
				Type:        schema.TypeBool,
			},

			// Attributes
			Attr_Images: {
				Computed:    true,
				Description: "Lists all the images in the IBM Power Virtual Server Cloud.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_Architecture: {
							Computed:    true,
							Description: "The CPU architecture that the image is designed for.",
							Type:        schema.TypeString,
						},
						Attr_ContainerFormat: {
							Computed:    true,
							Description: "The container format.",
							Type:        schema.TypeString,
						},
						Attr_CreationDate: {
							Computed:    true,
							Description: "Date of image creation",
							Type:        schema.TypeString,
						},
						Attr_CRN: {
							Computed:    true,
							Description: "CRN of this resource.",
							Type:        schema.TypeString,
						},
						Attr_Description: {
							Computed:    true,
							Description: "The description of an image.",
							Type:        schema.TypeString,
						},
						Attr_DiskFormat: {
							Computed:    true,
							Description: "The disk format.",
							Type:        schema.TypeString,
						},
						Attr_Endianness: {
							Computed:    true,
							Description: "The Endianness order.",
							Type:        schema.TypeString,
						},
						Attr_Href: {
							Computed:    true,
							Description: "The href of an image.",
							Type:        schema.TypeString,
						},
						Attr_HypervisorType: {
							Computed:    true,
							Description: "Hypervisor type.",
							Type:        schema.TypeString,
						},
						Attr_ImageID: {
							Computed:    true,
							Description: "The unique identifier of an image.",
							Type:        schema.TypeString,
						},
						Attr_ImageType: {
							Computed:    true,
							Description: "The identifier of this image type.",
							Type:        schema.TypeString,
						},
						Attr_LastUpdateDate: {
							Computed:    true,
							Description: "The last updated date of an image.",
							Type:        schema.TypeString,
						},
						Attr_Name: {
							Computed:    true,
							Description: "The name of the image.",
							Type:        schema.TypeString,
						},
						Attr_OperatingSystem: {
							Computed:    true,
							Description: "Operating System.",
							Type:        schema.TypeString,
						},
						Attr_State: {
							Computed:    true,
							Description: "The state of an Operating System.",
							Type:        schema.TypeString,
						},
						Attr_StoragePool: {
							Computed:    true,
							Description: "Storage pool where image resides.",
							Type:        schema.TypeString,
						},
						Attr_StorageType: {
							Computed:    true,
							Description: "The storage type of an image.",
							Type:        schema.TypeString,
						},
					},
				},
				Type: schema.TypeList,
			},
		},
	}
}

func dataSourceIBMPICatalogImagesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	includeSAP := d.Get(Arg_SAP).(bool)
	includeVTL := d.Get(Arg_VTL).(bool)
	imageC := instance.NewIBMPIImageClient(ctx, sess, cloudInstanceID)
	stockImages, err := imageC.GetAllStockImages(includeSAP, includeVTL)
	if err != nil {
		return diag.FromErr(err)
	}

	images := make([]map[string]interface{}, 0)
	for _, i := range stockImages.Images {
		image := make(map[string]interface{})
		image[Attr_ImageID] = *i.ImageID
		image[Attr_Name] = *i.Name

		if i.Description != nil {
			image[Attr_Description] = *i.Description
		}
		if i.CreationDate != nil {
			image[Attr_CreationDate] = i.CreationDate.String()
		}
		if i.Crn != "" {
			image[Attr_CRN] = i.Crn
		}
		if i.Href != nil {
			image[Attr_Href] = *i.Href
		}
		if i.LastUpdateDate != nil {
			image[Attr_LastUpdateDate] = i.LastUpdateDate.String()
		}
		if i.Specifications != nil {
			s := i.Specifications
			if s.Architecture != "" {
				image[Attr_Architecture] = s.Architecture
			}
			if s.ContainerFormat != "" {
				image[Attr_ContainerFormat] = s.ContainerFormat
			}
			if s.DiskFormat != "" {
				image[Attr_DiskFormat] = s.DiskFormat
			}
			if s.Endianness != "" {
				image[Attr_Endianness] = s.Endianness
			}
			if s.HypervisorType != "" {
				image[Attr_HypervisorType] = s.HypervisorType
			}
			if s.ImageType != "" {
				image[Attr_ImageType] = s.ImageType
			}
			if s.OperatingSystem != "" {
				image[Attr_OperatingSystem] = s.OperatingSystem
			}
		}
		if i.State != nil {
			image[Attr_State] = *i.State
		}
		if i.StoragePool != nil {
			image[Attr_StoragePool] = *i.StoragePool
		}
		if i.StorageType != nil {
			image[Attr_StorageType] = *i.StorageType
		}
		images = append(images, image)
	}
	d.SetId(time.Now().UTC().String())
	d.Set(Attr_Images, images)

	return nil
}
