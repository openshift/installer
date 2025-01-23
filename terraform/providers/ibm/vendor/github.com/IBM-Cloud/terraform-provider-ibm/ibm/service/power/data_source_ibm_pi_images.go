// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Datasource to list images that are available when a power instance is created
func DataSourceIBMPIImages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIImagesAllRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_ImageInfo: {
				Computed:    true,
				Description: "List of all supported images.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_Href: {
							Computed:    true,
							Description: "The hyper link of an image.",
							Type:        schema.TypeString,
						},
						Attr_ID: {
							Computed:    true,
							Description: "The unique identifier of an image.",
							Type:        schema.TypeString,
						},
						Attr_ImageType: {
							Computed:    true,
							Description: "The identifier of this image type.",
							Type:        schema.TypeString,
						},
						Attr_Name: {
							Computed:    true,
							Description: "The name of an image.",
							Type:        schema.TypeString,
						},
						Attr_State: {
							Computed:    true,
							Description: "The state of an image.",
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

func dataSourceIBMPIImagesAllRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)

	imageC := instance.NewIBMPIImageClient(ctx, sess, cloudInstanceID)
	imagedata, err := imageC.GetAll()
	if err != nil {
		return diag.FromErr(err)
	}

	var clientgenU, _ = uuid.GenerateUUID()
	d.SetId(clientgenU)
	d.Set(Attr_ImageInfo, flattenStockImages(imagedata.Images))

	return nil
}

func flattenStockImages(list []*models.ImageReference) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		l := map[string]interface{}{
			Attr_Href:        *i.Href,
			Attr_ID:          *i.ImageID,
			Attr_ImageType:   i.Specifications.ImageType,
			Attr_Name:        *i.Name,
			Attr_State:       *i.State,
			Attr_StoragePool: *i.StoragePool,
			Attr_StorageType: *i.StorageType,
		}
		result = append(result, l)
	}
	return result
}
