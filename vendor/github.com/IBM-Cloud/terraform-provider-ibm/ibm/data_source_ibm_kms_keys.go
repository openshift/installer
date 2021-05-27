// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	kp "github.com/IBM/keyprotect-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMKMSkeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMKMSKeysRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Key protect or hpcs instance GUID",
			},
			"key_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The name of the key to be fetched",
				ConflictsWith: []string{"alias"},
			},
			"alias": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The name of the key to be fetched",
				ConflictsWith: []string{"key_name"},
			},
			"endpoint_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{"public", "private"}),
				Description:  "public or private",
				ForceNew:     true,
				Default:      "public",
			},
			"keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"aliases": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
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
						"key_ring_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The key ring id of the key to be fetched",
						},
						"standard_key": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"policies": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rotation": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"created_by": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"creation_date": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"updated_by": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"last_update_date": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"interval_month": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
									"dual_auth_delete": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"created_by": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"creation_date": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"updated_by": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"last_update_date": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"enabled": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

}

func dataSourceIBMKMSKeysRead(d *schema.ResourceData, meta interface{}) error {
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
	if v, ok := d.GetOk("alias"); ok {
		aliasName := v.(string)
		key, err := api.GetKey(context.Background(), aliasName)
		if err != nil {
			return fmt.Errorf(
				"Get Keys failed with error: %s", err)
		} else {
			keyMap := make([]map[string]interface{}, 0, 1)
			keyInstance := make(map[string]interface{})
			keyInstance["id"] = key.ID
			keyInstance["name"] = key.Name
			keyInstance["crn"] = key.CRN
			keyInstance["standard_key"] = key.Extractable
			keyInstance["aliases"] = key.Aliases
			keyInstance["key_ring_id"] = key.KeyRingID
			keyMap = append(keyMap, keyInstance)
			d.Set("keys", keyMap)

		}
	} else {
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
			keyInstance["aliases"] = key.Aliases
			keyInstance["key_ring_id"] = key.KeyRingID
			keyMap = append(keyMap, keyInstance)

		}
		d.Set("keys", keyMap)
	}

	d.SetId(instanceID)
	d.Set("instance_id", instanceID)

	return nil

}
