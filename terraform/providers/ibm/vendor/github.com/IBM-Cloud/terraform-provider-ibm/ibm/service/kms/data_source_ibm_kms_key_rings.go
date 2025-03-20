// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kms

import (
	"context"

	//kp "github.com/IBM/keyprotect-go-client"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMKMSkeyRings() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMKMSKeyRingsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Key protect or hpcs instance GUID",
				DiffSuppressFunc: suppressKMSInstanceIDDiff,
			},
			"endpoint_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"public", "private"}),
				Description:  "public or private",
				Default:      "public",
			},
			"key_rings": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Key Rings for a particualer instance",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMKMSKeyRingsRead(d *schema.ResourceData, meta interface{}) error {
	instanceID := getInstanceIDFromCRN(d.Get("instance_id").(string))
	api, _, err := populateKPClient(d, meta, instanceID)
	if err != nil {
		return err
	}
	endpointType := d.Get("endpoint_type").(string)
	keys, err := api.GetKeyRings(context.Background())
	if err != nil || keys == nil {
		return flex.FmtErrorf("[ERROR] Get Key Rings failed with error: %s", err)
	}
	if keys.KeyRings == nil || len(keys.KeyRings) == 0 {
		return flex.FmtErrorf("[ERROR] No key Rings in instance  %s", instanceID)
	}

	keyRingMap := make([]map[string]interface{}, 0, len(keys.KeyRings))

	for _, keyRing := range keys.KeyRings {
		keyInstance := make(map[string]interface{})

		keyInstance["id"] = keyRing.ID
		keyInstance["created_by"] = keyRing.CreatedBy
		if keyRing.CreationDate != nil {
			keyInstance["creation_date"] = keyRing.CreationDate.String()
		}
		keyRingMap = append(keyRingMap, keyInstance)

	}

	d.SetId(instanceID)
	d.Set("key_rings", keyRingMap)
	d.Set("instance_id", instanceID)
	d.Set("endpoint_type", endpointType)

	return nil

}
