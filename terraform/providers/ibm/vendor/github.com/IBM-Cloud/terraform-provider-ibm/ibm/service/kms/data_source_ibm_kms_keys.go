// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kms

import (
	"context"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	kp "github.com/IBM/keyprotect-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMKMSkeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMKMSKeysRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Key protect or hpcs instance GUID",
				DiffSuppressFunc: suppressKMSInstanceIDDiff,
			},
			"key_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The name of the key to be fetched",
				ConflictsWith: []string{"alias", "key_id"},
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Limit till the keys to be fetched",
			},
			"alias": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The name of the key to be fetched",
				ConflictsWith: []string{"key_name", "key_id"},
			},
			"key_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"alias", "key_name"},
			},
			"endpoint_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"public", "private"}),
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
						"description": {
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
												"crn": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Cloud Resource Name (CRN) that uniquely identifies your cloud resources.",
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
												"enabled": {
													Type:     schema.TypeBool,
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
												"crn": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Cloud Resource Name (CRN) that uniquely identifies your cloud resources.",
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
	instanceID := getInstanceIDFromCRN(d.Get("instance_id").(string))
	api, _, err := populateKPClient(d, meta, instanceID)
	if err != nil {
		return err
	}
	var totalKeys []kp.Key
	if v, ok := d.GetOk("alias"); ok {
		aliasName := v.(string)
		key, err := api.GetKey(context.Background(), aliasName)
		if err != nil {
			return flex.FmtErrorf("[ERROR] Get Keys failed with error: %s", err)
		}
		keyMap := make([]map[string]interface{}, 0, 1)
		keyInstance := make(map[string]interface{})
		keyInstance["id"] = key.ID
		keyInstance["name"] = key.Name
		keyInstance["crn"] = key.CRN
		keyInstance["standard_key"] = key.Extractable
		keyInstance["description"] = key.Description
		keyInstance["aliases"] = key.Aliases
		keyInstance["key_ring_id"] = key.KeyRingID
		policies, err := api.GetPolicies(context.Background(), key.ID)
		if err != nil {
			return flex.FmtErrorf("[ERROR] Failed to read policies: %s", err)
		}
		if len(policies) == 0 {
			log.Printf("No Policy Configurations read\n")
		} else {
			keyInstance["policies"] = flex.FlattenKeyPolicies(policies)
		}
		keyMap = append(keyMap, keyInstance)
		d.Set("keys", keyMap)

	} else if v, ok := d.GetOk("key_id"); ok {
		key, err := api.GetKey(context.Background(), v.(string))
		if err != nil {
			return flex.FmtErrorf("[ERROR] Get Keys failed with error: %s", err)
		}
		keyMap := make([]map[string]interface{}, 0, 1)
		keyInstance := make(map[string]interface{})
		keyInstance["id"] = key.ID
		keyInstance["name"] = key.Name
		keyInstance["crn"] = key.CRN
		keyInstance["standard_key"] = key.Extractable
		keyInstance["description"] = key.Description
		keyInstance["aliases"] = key.Aliases
		keyInstance["key_ring_id"] = key.KeyRingID
		policies, err := api.GetPolicies(context.Background(), key.ID)
		if err != nil {
			return flex.FmtErrorf("[ERROR] Failed to read policies: %s", err)
		}
		if len(policies) == 0 {
			log.Printf("No Policy Configurations read\n")
		} else {
			keyInstance["policies"] = flex.FlattenKeyPolicies(policies)
		}
		keyMap = append(keyMap, keyInstance)

		d.SetId(instanceID)
		d.Set("keys", keyMap)
		d.Set("instance_id", instanceID)
	} else {
		limit := d.Get("limit")
		limitVal := limit.(int)
		offset := 0
		//default page size of API is 200 as stated
		pageSize := 200

		// when the limit is not passed, the api works in default way to avoid backward compatibility issues

		if limitVal == 0 {
			{
				keys, err := api.GetKeys(context.Background(), 0, offset)
				if err != nil {
					return flex.FmtErrorf("[ERROR] Get Keys failed with error: %s", err)
				}
				retreivedKeys := keys.Keys
				totalKeys = append(totalKeys, retreivedKeys...)
			}
		} else {
			// when the limit is passed by the user
			for {
				if offset < limitVal {
					if (limitVal - offset) < pageSize {
						keys, err := api.GetKeys(context.Background(), (limitVal - offset), offset)
						if err != nil {
							return flex.FmtErrorf("[ERROR] Get Keys failed with error: %s", err)
						}
						retreivedKeys := keys.Keys
						totalKeys = append(totalKeys, retreivedKeys...)
						break
					} else {
						keys, err := api.GetKeys(context.Background(), pageSize, offset)
						if err != nil {
							return flex.FmtErrorf("[ERROR] Get Keys failed with error: %s", err)
						}
						numOfKeysFetched := keys.Metadata.NumberOfKeys
						retreivedKeys := keys.Keys
						totalKeys = append(totalKeys, retreivedKeys...)
						if numOfKeysFetched < pageSize || offset+pageSize == limitVal {
							break
						}

						offset = offset + pageSize
					}
				}
			}
		}
		if len(totalKeys) == 0 {
			return flex.FmtErrorf("[ERROR] No keys in instance %s", instanceID)
		}
		var keyName string
		var matchKeys []kp.Key
		if v, ok := d.GetOk("key_name"); ok {
			keyName = v.(string)
			for _, keyData := range totalKeys {
				if keyData.Name == keyName {
					matchKeys = append(matchKeys, keyData)
				}
			}
		} else {
			matchKeys = totalKeys
		}

		if len(matchKeys) == 0 {
			return flex.FmtErrorf("[ERROR] No keys with name %s in instance  %s", keyName, instanceID)
		}

		keyMap := make([]map[string]interface{}, 0, len(matchKeys))

		for _, key := range matchKeys {
			keyInstance := make(map[string]interface{})
			keyInstance["id"] = key.ID
			keyInstance["name"] = key.Name
			keyInstance["crn"] = key.CRN
			keyInstance["standard_key"] = key.Extractable
			keyInstance["description"] = key.Description
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
