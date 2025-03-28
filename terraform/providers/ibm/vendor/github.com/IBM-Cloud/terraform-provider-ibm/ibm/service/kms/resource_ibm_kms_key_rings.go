// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kms

import (
	"context"
	"fmt"
	"strings"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	kp "github.com/IBM/keyprotect-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMKmskeyRings() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMKmsKeyRingCreate,
		Update:   resourceIBMKmsKeyRingUpdate,
		Delete:   resourceIBMKmsKeyRingDelete,
		Read:     resourceIBMKmsKeyRingRead,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Key protect Instance GUID",
				ForceNew:         true,
				DiffSuppressFunc: suppressKMSInstanceIDDiff,
			},
			"key_ring_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "User defined unique ID for the key ring",
				ValidateFunc: validate.InvokeValidator("ibm_kms_key_rings", "key_ring_id"),
			},
			"force_delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "(Deprecated) set to true to force delete this key ring. This allows key ring deletion as long as all keys inside have key state equals to 5 (destroyed). Keys are moved to the default key ring.",
				ForceNew:    false,
				Default:     false,
				Deprecated:  "force_delete is now deprecated. Please remove all references to this field.",
			},
			"endpoint_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"public", "private"}),
				Description:  "public or private",
				ForceNew:     true,
			},
		},
	}
}

func ResourceIBMKeyRingValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "key_ring_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-zA-Z0-9-]*$`,
			MinValueLength:             2,
			MaxValueLength:             100})

	ibmKeyRingResourceValidator := validate.ResourceValidator{ResourceName: "ibm_kms_key_rings", Schema: validateSchema}
	return &ibmKeyRingResourceValidator
}

func resourceIBMKmsKeyRingCreate(d *schema.ResourceData, meta interface{}) error {
	instanceID := getInstanceIDFromCRN(d.Get("instance_id").(string))
	keyRingID := d.Get("key_ring_id").(string)
	kpAPI, instanceCRN, err := populateKPClient(d, meta, instanceID)
	if err != nil {
		return err
	}

	err = kpAPI.CreateKeyRing(context.Background(), keyRingID)
	if err != nil {
		return flex.FmtErrorf("[ERROR] Error while creating key ring : %s", err)
	}
	var keyRing string
	keyRings, err := kpAPI.GetKeyRings(context.Background())
	if err != nil {
		return flex.FmtErrorf("[ERROR] Error while fetching key ring : %s", err)
	}
	for _, v := range keyRings.KeyRings {
		if v.ID == keyRingID {
			keyRing = v.ID
			break
		}
	}

	d.SetId(fmt.Sprintf("%s:keyRing:%s", keyRing, *instanceCRN))

	return resourceIBMKmsKeyRingRead(d, meta)
}

func resourceIBMKmsKeyRingUpdate(d *schema.ResourceData, meta interface{}) error {

	if d.HasChange("force_delete") {
		d.Set("force_delete", d.Get("force_delete").(bool))
	}
	return resourceIBMKmsKeyRingRead(d, meta)

}

func resourceIBMKmsKeyRingRead(d *schema.ResourceData, meta interface{}) error {
	id := strings.Split(d.Id(), ":keyRing:")
	if len(id) < 2 {
		return flex.FmtErrorf("[ERROR] Incorrect ID %s: Id should be a combination of keyRingID:keyRing:InstanceCRN", d.Id())
	}
	instanceID := getInstanceIDFromCRN(id[1])
	kpAPI, _, err := populateKPClient(d, meta, instanceID)
	if err != nil {
		return err
	}
	_, err = kpAPI.GetKeyRings(context.Background())
	if err != nil {
		if kpError, ok := err.(*kp.Error); ok {
			if kpError.StatusCode == 404 || kpError.StatusCode == 409 {
				d.SetId("")
				return nil
			}
		}
		return flex.FmtErrorf("[ERROR] Get Key Rings failed with error: %s", err)
	}

	d.Set("instance_id", instanceID)
	if strings.Contains((kpAPI.URL).String(), "private") || strings.Contains(kpAPI.Config.BaseURL, "private") {
		d.Set("endpoint_type", "private")
	} else {
		d.Set("endpoint_type", "public")
	}
	d.Set("key_ring_id", id[0])
	return nil
}

func resourceIBMKmsKeyRingDelete(d *schema.ResourceData, meta interface{}) error {
	id := strings.Split(d.Id(), ":keyRing:")
	instanceID := getInstanceIDFromCRN(id[1])
	kpAPI, _, err := populateKPClient(d, meta, instanceID)
	if err != nil {
		return err
	}

	err = kpAPI.DeleteKeyRing(context.Background(), id[0], kp.WithForce(true))
	if err != nil {
		if kpError, ok := err.(*kp.Error); ok {
			if kpError.StatusCode == 404 {
				return nil
			}
		}
	}
	return nil

}
