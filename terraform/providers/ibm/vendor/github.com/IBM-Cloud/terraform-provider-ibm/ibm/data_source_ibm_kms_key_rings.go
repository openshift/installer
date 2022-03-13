// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"strings"

	//kp "github.com/IBM/keyprotect-go-client"
	rc "github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIBMKMSkeyRings() *schema.Resource {
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
				ValidateFunc: validateAllowedStringValue([]string{"public", "private"}),
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
	api, err := meta.(ClientSession).keyManagementAPI()
	if err != nil {
		return err
	}

	instanceID := d.Get("instance_id").(string)
	CrnInstanceID := strings.Split(instanceID, ":")
	if len(CrnInstanceID) > 3 {
		instanceID = CrnInstanceID[len(CrnInstanceID)-3]
	}
	endpointType := d.Get("endpoint_type").(string)

	rsConClient, err := meta.(ClientSession).ResourceControllerV2API()
	if err != nil {
		return err
	}
	resourceInstanceGet := rc.GetResourceInstanceOptions{
		ID: &instanceID,
	}

	instanceData, resp, err := rsConClient.GetResourceInstance(&resourceInstanceGet)
	if err != nil || instanceData == nil {
		return fmt.Errorf("[ERROR] Error retrieving resource instance: %s with resp code: %s", err, resp)
	}
	extensions := instanceData.Extensions
	URL, err := KmsEndpointURL(api, endpointType, extensions)
	if err != nil {
		return err
	}
	api.URL = URL

	api.Config.InstanceID = instanceID
	keys, err := api.GetKeyRings(context.Background())
	if err != nil {
		return fmt.Errorf(
			"Get Key Rings failed with error: %s", err)
	}
	retreivedKeyRings := keys.KeyRings
	if keys == nil || len(retreivedKeyRings) == 0 {
		return fmt.Errorf("No key Rings in instance  %s", instanceID)
	}
	var keyRingName string

	if len(retreivedKeyRings) == 0 {
		return fmt.Errorf("No key Ring with name %s in instance  %s", keyRingName, instanceID)
	}

	keyRingMap := make([]map[string]interface{}, 0, len(retreivedKeyRings))

	for _, keyRing := range retreivedKeyRings {
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
