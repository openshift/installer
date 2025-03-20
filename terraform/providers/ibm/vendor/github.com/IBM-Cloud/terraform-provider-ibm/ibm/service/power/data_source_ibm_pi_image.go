// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"log"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPIImage() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIImagesRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_ImageName: {
				Description:  "The ID of the image.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_Architecture: {
				Computed:    true,
				Description: "The CPU architecture that the image is designed for. ",
				Type:        schema.TypeString,
			},
			Attr_CRN: {
				Computed:    true,
				Description: "The CRN of this resource.",
				Type:        schema.TypeString,
			},
			Attr_Hypervisor: {
				Computed:    true,
				Description: "Hypervision Type.",
				Type:        schema.TypeString,
			},
			Attr_ImageType: {
				Computed:    true,
				Description: "The identifier of this image type.",
				Type:        schema.TypeString,
			},
			Attr_OperatingSystem: {
				Computed:    true,
				Description: "The operating system that is installed with the image.",
				Type:        schema.TypeString,
			},
			Attr_Size: {
				Computed:    true,
				Description: "The size of the image in megabytes.",
				Type:        schema.TypeInt,
			},
			Attr_SourceChecksum: {
				Computed:    true,
				Description: "Checksum of the image.",
				Type:        schema.TypeString,
			},
			Attr_State: {
				Computed:    true,
				Description: "The state for this image. ",
				Type:        schema.TypeString,
			},
			Attr_StoragePool: {
				Computed:    true,
				Description: "Storage pool where image resides.",
				Type:        schema.TypeString,
			},
			Attr_StorageType: {
				Computed:    true,
				Description: "The storage type for this image.",
				Type:        schema.TypeString,
			},
			Attr_UserTags: {
				Computed:    true,
				Description: "List of user tags attached to the resource.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Type:        schema.TypeSet,
			},
		},
	}
}

func dataSourceIBMPIImagesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)

	imageC := instance.NewIBMPIImageClient(ctx, sess, cloudInstanceID)
	imagedata, err := imageC.Get(d.Get(Arg_ImageName).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*imagedata.ImageID)
	d.Set(Attr_Architecture, imagedata.Specifications.Architecture)
	if imagedata.Crn != "" {
		d.Set(Attr_CRN, imagedata.Crn)
		tags, err := flex.GetTagsUsingCRN(meta, string(imagedata.Crn))
		if err != nil {
			log.Printf("Error on get of pi image (%s) user_tags: %s", *imagedata.ImageID, err)
		}
		d.Set(Attr_UserTags, tags)
	}
	d.Set(Attr_Hypervisor, imagedata.Specifications.HypervisorType)
	d.Set(Attr_ImageType, imagedata.Specifications.ImageType)
	d.Set(Attr_OperatingSystem, imagedata.Specifications.OperatingSystem)
	d.Set(Attr_Size, imagedata.Size)
	d.Set(Attr_SourceChecksum, imagedata.Specifications.SourceChecksum)
	d.Set(Attr_State, imagedata.State)
	d.Set(Attr_StoragePool, imagedata.StoragePool)
	d.Set(Attr_StorageType, imagedata.StorageType)

	return nil
}
