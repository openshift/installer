// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	//"fmt"
	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

/*
Datasource to get the list of images that are available when a power instance is created

*/
func dataSourceIBMPIImages() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceIBMPIImagesAllRead,
		Schema: map[string]*schema.Schema{

			helpers.PIImageName: {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Imagename Name to be used for pvminstances",
				ValidateFunc: validation.NoZeroValues,
				Deprecated:   "This field is deprectaed.",
			},
			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			// Computed Attributes

			"image_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPIImagesAllRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := meta.(ClientSession).IBMPISession()

	if err != nil {
		return err
	}

	powerinstanceid := d.Get(helpers.PICloudInstanceId).(string)

	imageC := instance.NewIBMPIImageClient(sess, powerinstanceid)

	imagedata, err := imageC.GetAll(powerinstanceid)

	if err != nil {
		return err
	}

	var clientgenU, _ = uuid.GenerateUUID()
	d.SetId(clientgenU)
	_ = d.Set("image_info", flattenStockImages(imagedata.Images))

	return nil

}

func flattenStockImages(list []*models.ImageReference) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {

		l := map[string]interface{}{
			"id":           *i.ImageID,
			"state":        *i.State,
			"href":         *i.Href,
			"name":         *i.Name,
			"storage_type": *i.StorageType,
		}

		result = append(result, l)

	}
	return result
}
