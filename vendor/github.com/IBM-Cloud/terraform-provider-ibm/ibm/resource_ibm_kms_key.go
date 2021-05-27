// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	kp "github.com/IBM/keyprotect-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIBMKmskey() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMKmsKeyCreate,
		Read:     resourceIBMKmsKeyRead,
		Update:   resourceIBMKmsKeyUpdate,
		Delete:   resourceIBMKmsKeyDelete,
		Exists:   resourceIBMKmsKeyExists,
		Importer: &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Key protect or hpcs instance GUID",
			},
			"key_ring_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "default",
				Description: "Key Ring for the Key",
			},
			"key_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Key ID",
			},
			"key_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Key name",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "type of service hs-crypto or kms",
			},
			"endpoint_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{"public", "private"}),
				Description:  "public or private",
				ForceNew:     true,
				Default:      "public",
			},
			"standard_key": {
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
				ForceNew:    true,
				Description: "Standard key type",
			},
			"payload": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			"encrypted_nonce": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Only for imported root key",
			},
			"iv_value": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Only for imported root key",
			},
			"force_delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "set to true to force delete the key",
				ForceNew:    false,
				Default:     false,
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Crn of the key",
			},
			"expiration_date": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The date the key material expires. The date format follows RFC 3339. You can set an expiration date on any key on its creation. A key moves into the Deactivated state within one hour past its expiration date, if one is assigned. If you create a key without specifying an expiration date, the key does not expire",
				ForceNew:    true,
			},
			"policies": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Creates or updates one or more policies for the specified key",
				MinItems:    1,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rotation": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: []string{"policies.0.rotation", "policies.0.dual_auth_delete"},
							Description:  "Specifies the key rotation time interval in months, with a minimum of 1, and a maximum of 12",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The v4 UUID used to uniquely identify the policy resource, as specified by RFC 4122.",
									},
									"crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cloud Resource Name (CRN) that uniquely identifies your cloud resources.",
									},
									"created_by": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for the resource that created the policy.",
									},
									"creation_date": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The date the policy was created. The date format follows RFC 3339.",
									},
									"updated_by": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for the resource that updated the policy.",
									},
									"last_update_date": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Updates when the policy is replaced or modified. The date format follows RFC 3339.",
									},
									"interval_month": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validateAllowedRangeInt(1, 12),
										Description:  "Specifies the key rotation time interval in months",
									},
								},
							},
						},
						"dual_auth_delete": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: []string{"policies.0.rotation", "policies.0.dual_auth_delete"},
							Description:  "Data associated with the dual authorization delete policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The v4 UUID used to uniquely identify the policy resource, as specified by RFC 4122.",
									},
									"crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cloud Resource Name (CRN) that uniquely identifies your cloud resources.",
									},
									"created_by": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for the resource that created the policy.",
									},
									"creation_date": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The date the policy was created. The date format follows RFC 3339.",
									},
									"updated_by": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for the resource that updated the policy.",
									},
									"last_update_date": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Updates when the policy is replaced or modified. The date format follows RFC 3339.",
									},
									"enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "If set to true, Key Protect enables a dual authorization policy on a single key.",
									},
								},
							},
						},
					},
				},
			},
			ResourceName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource",
			},

			ResourceCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			ResourceStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the resource",
			},

			ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},
			ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about the resource",
			},
		},
	}
}

func resourceIBMKmsKeyCreate(d *schema.ResourceData, meta interface{}) error {
	kpAPI, err := meta.(ClientSession).keyManagementAPI()
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
	crnData := strings.Split(instanceCRN, ":")

	var hpcsEndpointURL string

	if crnData[4] == "hs-crypto" {
		hpcsEndpointAPI, err := meta.(ClientSession).HpcsEndpointAPI()
		if err != nil {
			return err
		}

		resp, err := hpcsEndpointAPI.Endpoint().GetAPIEndpoint(instanceID)
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
		kpAPI.URL = u
	} else if crnData[4] == "kms" {
		if endpointType == "private" {
			if !strings.Contains(kpAPI.Config.BaseURL, "private") {
				kmsEndpURL := strings.SplitAfter(kpAPI.Config.BaseURL, "https://")
				if len(kmsEndpURL) == 2 {
					kmsEndpointURL := kmsEndpURL[0] + "private." + kmsEndpURL[1]
					u, err := url.Parse(kmsEndpointURL)
					if err != nil {
						return fmt.Errorf("Error Parsing kms EndpointURL")
					}
					kpAPI.URL = u
				} else {
					return fmt.Errorf("Error in Kms EndPoint URL ")
				}
			}
		}
	} else {
		return fmt.Errorf("Invalid or unsupported service Instance")
	}
	kpAPI.Config.InstanceID = instanceID

	kpAPI.Config.KeyRing = d.Get("key_ring_id").(string)

	name := d.Get("key_name").(string)
	standardKey := d.Get("standard_key").(bool)

	var expiration *time.Time
	if es, ok := d.GetOk("expiration_date"); ok {
		expiration_string := es.(string)
		// parse string to required time format
		expiration_time, err := time.Parse(time.RFC3339, expiration_string)
		if err != nil {
			return fmt.Errorf("Invalid time format (the date format follows RFC 3339): %s", err)
		}
		expiration = &expiration_time
	} else {
		expiration = nil
	}

	var keyCRN string
	if standardKey {
		if v, ok := d.GetOk("payload"); ok {
			//import standard key
			payload := v.(string)
			stkey, err := kpAPI.CreateImportedStandardKey(context.Background(), name, expiration, payload)
			if err != nil {
				return fmt.Errorf(
					"Error while creating standard key with payload: %s", err)
			}
			keyCRN = stkey.CRN
			d.SetId(keyCRN)

		} else {
			//create standard key
			stkey, err := kpAPI.CreateStandardKey(context.Background(), name, expiration)
			if err != nil {
				return fmt.Errorf(
					"Error while creating standard key: %s", err)
			}
			keyCRN = stkey.CRN
			d.SetId(keyCRN)

		}
	} else {
		if v, ok := d.GetOk("payload"); ok {
			payload := v.(string)
			encryptedNonce := d.Get("encrypted_nonce").(string)
			iv := d.Get("iv_value").(string)
			stkey, err := kpAPI.CreateImportedRootKey(context.Background(), name, expiration, payload, encryptedNonce, iv)
			if err != nil {
				return fmt.Errorf(
					"Error while creating Root key with payload: %s", err)
			}
			keyCRN = stkey.CRN
			d.SetId(keyCRN)

		} else {
			stkey, err := kpAPI.CreateRootKey(context.Background(), name, expiration)
			if err != nil {
				return fmt.Errorf(
					"Error while creating Root key: %s", err)
			}
			keyCRN = stkey.CRN
			d.SetId(keyCRN)
		}
	}
	return resourceIBMKmsKeyUpdate(d, meta)
}

func resourceIBMKmsKeyRead(d *schema.ResourceData, meta interface{}) error {
	kpAPI, err := meta.(ClientSession).keyManagementAPI()
	if err != nil {
		return err
	}
	crn := d.Id()
	crnData := strings.Split(crn, ":")
	endpointType := crnData[3]
	instanceID := crnData[len(crnData)-3]
	keyid := crnData[len(crnData)-1]

	var instanceType string
	var hpcsEndpointURL string

	if crnData[4] == "hs-crypto" {
		instanceType = "hs-crypto"
		hpcsEndpointAPI, err := meta.(ClientSession).HpcsEndpointAPI()
		if err != nil {
			return err
		}

		resp, err := hpcsEndpointAPI.Endpoint().GetAPIEndpoint(instanceID)
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
		kpAPI.URL = u
	} else if crnData[4] == "kms" {
		instanceType = "kms"
		if endpointType == "private" {
			if !strings.Contains(kpAPI.Config.BaseURL, "private") {
				kmsEndpURL := strings.SplitAfter(kpAPI.Config.BaseURL, "https://")
				if len(kmsEndpURL) == 2 {
					kmsEndpointURL := kmsEndpURL[0] + "private." + kmsEndpURL[1]
					u, err := url.Parse(kmsEndpointURL)
					if err != nil {
						return fmt.Errorf("Error Parsing kms EndpointURL")
					}
					kpAPI.URL = u
				} else {
					return fmt.Errorf("Error in Kms EndPoint URL ")
				}
			}
		}
	} else {
		return fmt.Errorf("Invalid or unsupported service Instance")
	}

	kpAPI.Config.InstanceID = instanceID
	// keyid := d.Id()
	key, err := kpAPI.GetKey(context.Background(), keyid)
	if err != nil {
		return fmt.Errorf("Get Key failed with error: %s", err)
	}

	policies, err := kpAPI.GetPolicies(context.Background(), keyid)

	if err != nil {
		return fmt.Errorf("Failed to read policies: %s", err)
	}
	if len(policies) == 0 {
		log.Printf("No Policy Configurations read\n")
	} else {
		d.Set("policies", flattenKeyPolicies(policies))
	}
	d.Set("instance_id", instanceID)
	d.Set("key_id", keyid)
	d.Set("standard_key", key.Extractable)
	d.Set("payload", key.Payload)
	d.Set("encrypted_nonce", key.EncryptedNonce)
	d.Set("iv_value", key.IV)
	d.Set("key_name", key.Name)
	d.Set("crn", key.CRN)
	d.Set("endpoint_type", endpointType)
	d.Set("type", instanceType)
	d.Set("force_delete", d.Get("force_delete").(bool))
	d.Set("key_ring_id", key.KeyRingID)
	if key.Expiration != nil {
		expiration := key.Expiration
		d.Set("expiration_date", expiration.Format(time.RFC3339))
	} else {
		d.Set("expiration_date", "")
	}
	d.Set(ResourceName, key.Name)
	d.Set(ResourceCRN, key.CRN)
	state := key.State
	d.Set(ResourceStatus, strconv.Itoa(state))
	rcontroller, err := getBaseController(meta)
	if err != nil {
		return err
	}
	id := key.ID
	crn1 := strings.TrimSuffix(key.CRN, ":key:"+id)

	d.Set(ResourceControllerURL, rcontroller+"/services/kms/"+url.QueryEscape(crn1)+"%3A%3A")

	return nil

}

func resourceIBMKmsKeyUpdate(d *schema.ResourceData, meta interface{}) error {

	if d.HasChange("force_delete") {
		d.Set("force_delete", d.Get("force_delete").(bool))
	}
	if d.HasChange("policies") {

		kpAPI, err := meta.(ClientSession).keyManagementAPI()
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
		crnData := strings.Split(instanceCRN, ":")

		var hpcsEndpointURL string

		if crnData[4] == "hs-crypto" {
			hpcsEndpointAPI, err := meta.(ClientSession).HpcsEndpointAPI()
			if err != nil {
				return err
			}

			resp, err := hpcsEndpointAPI.Endpoint().GetAPIEndpoint(instanceID)
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
			kpAPI.URL = u
		} else if crnData[4] == "kms" {
			if endpointType == "private" {
				if !strings.Contains(kpAPI.Config.BaseURL, "private") {
					kmsEndpURL := strings.SplitAfter(kpAPI.Config.BaseURL, "https://")
					if len(kmsEndpURL) == 2 {
						kmsEndpointURL := kmsEndpURL[0] + "private." + kmsEndpURL[1]
						u, err := url.Parse(kmsEndpointURL)
						if err != nil {
							return fmt.Errorf("Error Parsing kms EndpointURL")
						}
						kpAPI.URL = u
					} else {
						return fmt.Errorf("Error in Kms EndPoint URL ")
					}
				}
			}
		} else {
			return fmt.Errorf("Invalid or unsupported service Instance")
		}

		kpAPI.Config.InstanceID = instanceID

		crn := d.Id()
		crnData = strings.Split(crn, ":")
		key_id := crnData[len(crnData)-1]

		err = handlePolicies(d, kpAPI, meta, key_id)
		if err != nil {
			resourceIBMKmsKeyRead(d, meta)
			return fmt.Errorf("Could not create policies: %s", err)
		}
	}
	return resourceIBMKmsKeyRead(d, meta)

}

func resourceIBMKmsKeyDelete(d *schema.ResourceData, meta interface{}) error {
	kpAPI, err := meta.(ClientSession).keyManagementAPI()
	if err != nil {
		return err
	}
	crn := d.Id()
	crnData := strings.Split(crn, ":")
	endpointType := crnData[3]
	instanceID := crnData[len(crnData)-3]
	keyid := crnData[len(crnData)-1]
	kpAPI.Config.InstanceID = instanceID

	var hpcsEndpointURL string

	if crnData[4] == "hs-crypto" {
		hpcsEndpointAPI, err := meta.(ClientSession).HpcsEndpointAPI()
		if err != nil {
			return err
		}

		resp, err := hpcsEndpointAPI.Endpoint().GetAPIEndpoint(instanceID)
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
		kpAPI.URL = u
	} else if crnData[4] == "kms" {
		if endpointType == "private" {
			if !strings.Contains(kpAPI.Config.BaseURL, "private") {
				kmsEndpURL := strings.SplitAfter(kpAPI.Config.BaseURL, "https://")
				if len(kmsEndpURL) == 2 {
					kmsEndpointURL := kmsEndpURL[0] + "private." + kmsEndpURL[1]
					u, err := url.Parse(kmsEndpointURL)
					if err != nil {
						return fmt.Errorf("Error Parsing kms EndpointURL")
					}
					kpAPI.URL = u
				} else {
					return fmt.Errorf("Error in Kms EndPoint URL ")
				}
			}
		}
	} else {
		return fmt.Errorf("Invalid or unsupported service Instance")
	}

	force := d.Get("force_delete").(bool)
	f := kp.ForceOpt{
		Force: force,
	}

	_, err1 := kpAPI.DeleteKey(context.Background(), keyid, kp.ReturnRepresentation, f)
	if err1 != nil {
		return fmt.Errorf(
			"Error while deleting: %s", err1)
	}
	d.SetId("")
	return nil

}

func resourceIBMKmsKeyExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	kpAPI, err := meta.(ClientSession).keyManagementAPI()
	if err != nil {
		return false, err
	}

	crn := d.Id()
	crnData := strings.Split(crn, ":")
	endpointType := crnData[3]
	instanceID := crnData[len(crnData)-3]
	keyid := crnData[len(crnData)-1]
	kpAPI.Config.InstanceID = instanceID

	var hpcsEndpointURL string

	if crnData[4] == "hs-crypto" {
		hpcsEndpointAPI, err := meta.(ClientSession).HpcsEndpointAPI()
		if err != nil {
			return false, err
		}

		resp, err := hpcsEndpointAPI.Endpoint().GetAPIEndpoint(instanceID)
		if err != nil {
			return false, err
		}

		if endpointType == "public" {
			hpcsEndpointURL = "https://" + resp.Kms.Public + "/api/v2/keys"
		} else {
			hpcsEndpointURL = "https://" + resp.Kms.Private + "/api/v2/keys"
		}

		u, err := url.Parse(hpcsEndpointURL)
		if err != nil {
			return false, fmt.Errorf("Error Parsing hpcs EndpointURL")

		}
		kpAPI.URL = u
	} else if crnData[4] == "kms" {
		if endpointType == "private" {
			if !strings.Contains(kpAPI.Config.BaseURL, "private") {
				kmsEndpURL := strings.SplitAfter(kpAPI.Config.BaseURL, "https://")
				if len(kmsEndpURL) == 2 {
					kmsEndpointURL := kmsEndpURL[0] + "private." + kmsEndpURL[1]
					u, err := url.Parse(kmsEndpointURL)
					if err != nil {
						return false, fmt.Errorf("Error Parsing kms EndpointURL")
					}
					kpAPI.URL = u
				} else {
					return false, fmt.Errorf("Error in Kms EndPoint URL ")
				}
			}
		}
	} else {
		return false, fmt.Errorf("Invalid or unsupported service Instance")
	}

	_, err = kpAPI.GetKey(context.Background(), keyid)
	if err != nil {
		kpError := err.(*kp.Error)
		if kpError.StatusCode == 404 {
			return false, nil
		}
		return false, err
	}
	return true, nil

}

func handlePolicies(d *schema.ResourceData, kpAPI *kp.Client, meta interface{}, key_id string) error {
	var setRotation, setDualAuthDelete, dualAuthEnable bool
	var rotationInterval int

	if policyInfo, ok := d.GetOk("policies"); ok {

		policyDataList := policyInfo.([]interface{})
		policyData := policyDataList[0].(map[string]interface{})

		if rpd, ok := policyData["rotation"]; ok {
			rpdList := rpd.([]interface{})
			if len(rpdList) != 0 {
				rotationInterval = rpdList[0].(map[string]interface{})["interval_month"].(int)
				setRotation = true
			}
		}
		if dadp, ok := policyData["dual_auth_delete"]; ok {
			dadpList := dadp.([]interface{})
			if len(dadpList) != 0 {
				dualAuthEnable = dadpList[0].(map[string]interface{})["enabled"].(bool)
				setDualAuthDelete = true
			}
		}

		_, err := kpAPI.SetPolicies(context.Background(), key_id, setRotation, rotationInterval, setDualAuthDelete, dualAuthEnable)
		if err != nil {
			return fmt.Errorf("Error while creating policies: %s", err)
		}
	}
	return nil
}
