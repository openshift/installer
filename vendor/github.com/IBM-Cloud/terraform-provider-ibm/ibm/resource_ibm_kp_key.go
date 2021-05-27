// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	kp "github.com/IBM/keyprotect-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIBMkey() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMKeyCreate,
		Read:     resourceIBMKeyRead,
		Update:   resourceIBMKeyUpdate,
		Delete:   resourceIBMKeyDelete,
		Exists:   resourceIBMKeyExists,
		Importer: &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"key_protect_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Key protect instance ID",
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
			},
			"force_delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "set to true to force delete the key",
				ForceNew:    false,
				Default:     false,
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
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Crn of the key",
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

func resourceIBMKeyCreate(d *schema.ResourceData, meta interface{}) error {
	api, err := meta.(ClientSession).keyProtectAPI()
	if err != nil {
		return err
	}
	instanceID := d.Get("key_protect_id").(string)
	api.Config.InstanceID = instanceID
	name := d.Get("key_name").(string)
	standardKey := d.Get("standard_key").(bool)

	var keyCRN string
	if standardKey {
		if v, ok := d.GetOk("payload"); ok {
			//import standard key
			payload := v.(string)
			stkey, err := api.CreateImportedStandardKey(context.Background(), name, nil, payload)
			if err != nil {
				return fmt.Errorf(
					"Error while creating standard key: %s", err)
			}
			keyCRN = stkey.CRN
		} else {
			//create standard key
			stkey, err := api.CreateStandardKey(context.Background(), name, nil)
			if err != nil {
				return fmt.Errorf(
					"Error while creating standard key: %s", err)
			}
			keyCRN = stkey.CRN
		}
		d.SetId(keyCRN)
	} else {
		if v, ok := d.GetOk("payload"); ok {
			payload := v.(string)
			encryptedNonce := d.Get("encrypted_nonce").(string)
			iv := d.Get("iv_value").(string)
			stkey, err := api.CreateImportedRootKey(context.Background(), name, nil, payload, encryptedNonce, iv)
			if err != nil {
				return fmt.Errorf(
					"Error while creating Root key: %s", err)
			}
			keyCRN = stkey.CRN
		} else {
			stkey, err := api.CreateRootKey(context.Background(), name, nil)
			if err != nil {
				return fmt.Errorf(
					"Error while creating Root key: %s", err)
			}
			keyCRN = stkey.CRN
		}

		d.SetId(keyCRN)

	}
	d.Set("force_delete", d.Get("force_delete").(bool))

	return resourceIBMKeyRead(d, meta)
}

func resourceIBMKeyRead(d *schema.ResourceData, meta interface{}) error {
	api, err := meta.(ClientSession).keyProtectAPI()
	if err != nil {
		return err
	}
	crn := d.Id()
	crnData := strings.Split(crn, ":")

	instanceID := crnData[len(crnData)-3]
	keyid := crnData[len(crnData)-1]
	api.Config.InstanceID = instanceID
	// keyid := d.Id()
	key, err := api.GetKey(context.Background(), keyid)
	if err != nil {
		return fmt.Errorf(
			"Get Key failed with error: %s", err)
	}
	d.Set("key_id", keyid)
	d.Set("standard_key", key.Extractable)
	d.Set("payload", key.Payload)
	d.Set("encrypted_nonce", key.EncryptedNonce)
	d.Set("iv_value", key.IV)
	d.Set("key_name", key.Name)
	d.Set("crn", key.CRN)

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

func resourceIBMKeyUpdate(d *schema.ResourceData, meta interface{}) error {

	if d.HasChange("force_delete") {
		d.Set("force_delete", d.Get("force_delete").(bool))
	}
	return resourceIBMKeyRead(d, meta)

}

func resourceIBMKeyDelete(d *schema.ResourceData, meta interface{}) error {
	api, err := meta.(ClientSession).keyProtectAPI()
	if err != nil {
		return err
	}
	crn := d.Id()
	crnData := strings.Split(crn, ":")

	instanceID := crnData[len(crnData)-3]
	keyid := crnData[len(crnData)-1]
	api.Config.InstanceID = instanceID
	force := d.Get("force_delete").(bool)
	f := kp.ForceOpt{
		Force: force,
	}
	_, err1 := api.DeleteKey(context.Background(), keyid, kp.ReturnRepresentation, f)
	if err1 != nil {
		return fmt.Errorf(
			"Error while deleting: %s", err1)
	}
	d.SetId("")
	return nil

}

func resourceIBMKeyExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	api, err := meta.(ClientSession).keyProtectAPI()
	if err != nil {
		return false, err
	}
	crn := d.Id()
	crnData := strings.Split(crn, ":")

	instanceID := crnData[len(crnData)-3]
	keyid := crnData[len(crnData)-1]
	api.Config.InstanceID = instanceID
	// keyid := d.Id()
	_, err = api.GetKey(context.Background(), keyid)
	if err != nil {
		kpError := err.(*kp.Error)
		if kpError.StatusCode == 404 {
			return false, nil
		}
		return false, err
	}
	return true, nil

}
