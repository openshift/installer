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

func resourceIBMKmskeyRings() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMKmsKeyRingCreate,
		Delete:   resourceIBMKmsKeyRingDelete,
		Read:     resourceIBMKmsKeyRingRead,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Key protect Instance GUID",
				ForceNew:    true,
			},
			"key_ring_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "User defined unique ID for the key ring",
				ValidateFunc: InvokeValidator("ibm_kms_key_rings", "key_ring_id"),
			},
			"endpoint_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{"public", "private"}),
				Description:  "public or private",
				ForceNew:     true,
				Default:      "public",
			},
		},
	}
}

func resourceIBMKeyRingValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 1)

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "key_ring_id",
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Required:                   true,
			Regexp:                     `^[a-zA-Z0-9-]*$`,
			MinValueLength:             2,
			MaxValueLength:             100})

	ibmKeyRingResourceValidator := ResourceValidator{ResourceName: "ibm_kms_key_rings", Schema: validateSchema}
	return &ibmKeyRingResourceValidator
}

func resourceIBMKmsKeyRingCreate(d *schema.ResourceData, meta interface{}) error {
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
	keyRingID := d.Get("key_ring_id").(string)

	rContollerAPI := rContollerClient.ResourceServiceInstanceV2()

	instanceData, err := rContollerAPI.GetInstance(instanceID)
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
			if !strings.HasPrefix(kpAPI.Config.BaseURL, "private") {
				kpAPI.Config.BaseURL = "private." + kpAPI.Config.BaseURL
			}
		}
	} else {
		return fmt.Errorf("Invalid or unsupported service Instance")
	}
	kpAPI.Config.InstanceID = instanceID

	err = kpAPI.CreateKeyRing(context.Background(), keyRingID)
	if err != nil {
		return fmt.Errorf(
			"Error while creating key ring : %s", err)
	}
	var keyRing string
	keyRings, err2 := kpAPI.GetKeyRings(context.Background())
	if err2 != nil {
		return fmt.Errorf(
			"Error while fetching key ring : %s", err2)
	}
	for _, v := range keyRings.KeyRings {
		if v.ID == keyRingID {
			keyRing = v.ID
			break
		}
	}

	d.SetId(fmt.Sprintf("%s:keyRing:%s", keyRing, instanceCRN))

	return resourceIBMKmsKeyRingRead(d, meta)
}

func resourceIBMKmsKeyRingRead(d *schema.ResourceData, meta interface{}) error {
	kpAPI, err := meta.(ClientSession).keyManagementAPI()
	if err != nil {
		return err
	}
	id := strings.Split(d.Id(), ":keyRing:")
	crn := id[1]
	crnData := strings.Split(crn, ":")
	endpointType := crnData[3]
	instanceID := crnData[len(crnData)-3]

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
			if !strings.HasPrefix(kpAPI.Config.BaseURL, "private") {
				kpAPI.Config.BaseURL = "private." + kpAPI.Config.BaseURL
			}
		}
	} else {
		return fmt.Errorf("Invalid or unsupported service Instance")
	}

	kpAPI.Config.InstanceID = instanceID
	_, err = kpAPI.GetKeyRings(context.Background())
	if err != nil {
		kpError := err.(*kp.Error)
		if kpError.StatusCode == 404 || kpError.StatusCode == 409 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Get Key Rings failed with error: %s", err)
	}

	d.Set("instance_id", instanceID)
	d.Set("endpoint_type", endpointType)
	d.Set("key_ring_id", id[0])
	return nil
}

func resourceIBMKmsKeyRingDelete(d *schema.ResourceData, meta interface{}) error {
	kpAPI, err := meta.(ClientSession).keyManagementAPI()
	if err != nil {
		return err
	}
	id := strings.Split(d.Id(), ":keyRing:")
	crn := id[1]
	crnData := strings.Split(crn, ":")
	endpointType := crnData[3]
	instanceID := crnData[len(crnData)-3]

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
			if !strings.HasPrefix(kpAPI.Config.BaseURL, "private") {
				kpAPI.Config.BaseURL = "private." + kpAPI.Config.BaseURL
			}
		}
	} else {
		return fmt.Errorf("Invalid or unsupported service Instance")
	}

	kpAPI.Config.InstanceID = instanceID
	err1 := kpAPI.DeleteKeyRing(context.Background(), id[0])
	if err1 != nil {
		kpError := err1.(*kp.Error)
		if kpError.StatusCode == 404 || kpError.StatusCode == 409 {
			return nil
		} else {
			return fmt.Errorf(" failed to Destroy key ring with error: %s", err1)
		}
	}
	return nil

}
