// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"

	kp "github.com/IBM/keyprotect-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMkey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMKeyRead,

		Schema: map[string]*schema.Schema{
			"key_protect_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"crn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"standard_key": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}

}

func dataSourceIBMKeyRead(d *schema.ResourceData, meta interface{}) error {
	api, err := meta.(ClientSession).keyProtectAPI()
	if err != nil {
		return err
	}

	instanceID := d.Get("key_protect_id").(string)
	api.Config.InstanceID = instanceID
	keys, err := api.GetKeys(context.Background(), 100, 0)
	if err != nil {
		return fmt.Errorf(
			"Get Keys failed with error: %s", err)
	}
	retreivedKeys := keys.Keys
	if len(retreivedKeys) == 0 {
		return fmt.Errorf("No keys in instance  %s", instanceID)
	}
	var keyName string
	var matchKeys []kp.Key
	if v, ok := d.GetOk("key_name"); ok {
		keyName = v.(string)
		for _, keyData := range retreivedKeys {
			if keyData.Name == keyName {
				matchKeys = append(matchKeys, keyData)
			}
		}
	} else {
		matchKeys = retreivedKeys
	}

	if len(matchKeys) == 0 {
		return fmt.Errorf("No keys with name %s in instance  %s", keyName, instanceID)
	}

	keyMap := make([]map[string]interface{}, 0, len(matchKeys))

	for _, key := range matchKeys {
		keyInstance := make(map[string]interface{})
		keyInstance["id"] = key.ID
		keyInstance["name"] = key.Name
		keyInstance["crn"] = key.CRN
		keyInstance["standard_key"] = key.Extractable
		keyMap = append(keyMap, keyInstance)

	}

	d.SetId(instanceID)
	d.Set("keys", keyMap)
	d.Set("key_protect_id", instanceID)

	return nil

}
