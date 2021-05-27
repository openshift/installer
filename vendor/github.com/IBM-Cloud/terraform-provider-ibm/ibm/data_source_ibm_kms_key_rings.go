// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	//kp "github.com/IBM/keyprotect-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMKMSkeyRings() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMKMSKeyRingsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Key protect or hpcs instance GUID",
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

	rContollerClient, err := meta.(ClientSession).ResourceControllerAPIV2()
	if err != nil {
		return err
	}

	instanceID := d.Get("instance_id").(string)
	endpointType := d.Get("endpoint_type").(string)

	rContollerApi := rContollerClient.ResourceServiceInstanceV2()

	instanceData, err := rContollerApi.GetInstance(instanceID)
	if err != nil {
		return err
	}
	instanceCRN := instanceData.Crn.String()

	var hpcsEndpointURL string
	crnData := strings.Split(instanceCRN, ":")

	if crnData[4] == "hs-crypto" {

		hpcsEndpointApi, err := meta.(ClientSession).HpcsEndpointAPI()
		if err != nil {
			return err
		}
		resp, err := hpcsEndpointApi.Endpoint().GetAPIEndpoint(instanceID)
		if err != nil {
			return err
		}

		if endpointType == "public" {
			hpcsEndpointURL = "https://" + resp.Kms.Public + "/api/v2/keys"
		} else {
			hpcsEndpointURL = "https://" + resp.Kms.Private + "/api/v2/keys"
		}

		u, err := url.Parse(hpcsEndpointURL)
		if err != nil {
			return fmt.Errorf("Error Parsing hpcs EndpointURL")
		}
		api.URL = u
	} else if crnData[4] == "kms" {
		if endpointType == "private" {
			if !strings.HasPrefix(api.Config.BaseURL, "private") {
				api.Config.BaseURL = "private." + api.Config.BaseURL
			}
		}
	} else {
		return fmt.Errorf("Invalid or unsupported service Instance")
	}

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
