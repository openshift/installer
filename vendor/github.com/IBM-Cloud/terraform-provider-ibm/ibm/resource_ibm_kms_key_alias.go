package ibm

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	kp "github.com/IBM/keyprotect-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIBMKmskeyAlias() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMKmsKeyAliasCreate,
		Delete:   resourceIBMKmsKeyAliasDelete,
		Read:     resourceIBMKmsKeyAliasRead,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Key ID",
				ForceNew:    true,
			},
			"alias": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Key protect or hpcs key alias name",
			},
			"key_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Key ID",
				ForceNew:    true,
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

func resourceIBMKmsKeyAliasCreate(d *schema.ResourceData, meta interface{}) error {
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

	aliasName := d.Get("alias").(string)
	keyID := d.Get("key_id").(string)
	stkey, err := kpAPI.CreateKeyAlias(context.Background(), aliasName, keyID)
	if err != nil {
		return fmt.Errorf(
			"Error while creating alias name for the key: %s", err)
	}
	key, err := kpAPI.GetKey(context.Background(), stkey.KeyID)
	if err != nil {
		return fmt.Errorf("Get Key failed with error: %s", err)
	}
	d.SetId(fmt.Sprintf("%s:alias:%s", stkey.Alias, key.CRN))

	return resourceIBMKmsKeyAliasRead(d, meta)
}

func resourceIBMKmsKeyAliasRead(d *schema.ResourceData, meta interface{}) error {
	kpAPI, err := meta.(ClientSession).keyManagementAPI()
	if err != nil {
		return err
	}
	id := strings.Split(d.Id(), ":alias:")
	crn := id[1]
	crnData := strings.Split(crn, ":")
	endpointType := crnData[3]
	instanceID := crnData[len(crnData)-3]
	keyid := crnData[len(crnData)-1]

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
	key, err := kpAPI.GetKey(context.Background(), keyid)
	if err != nil {
		kpError := err.(*kp.Error)
		if kpError.StatusCode == 404 {
			d.SetId("")
			return nil
		} else {
			return fmt.Errorf("Get Key failed with error: %s", err)
		}
	}
	d.Set("alias", id[0])
	d.Set("key_id", key.ID)
	d.Set("instance_id", instanceID)
	d.Set("endpoint_type", endpointType)

	return nil
}

func resourceIBMKmsKeyAliasDelete(d *schema.ResourceData, meta interface{}) error {
	kpAPI, err := meta.(ClientSession).keyManagementAPI()
	if err != nil {
		return err
	}
	id := strings.Split(d.Id(), ":alias:")
	crn := id[1]
	crnData := strings.Split(crn, ":")
	endpointType := crnData[3]
	instanceID := crnData[len(crnData)-3]
	keyid := crnData[len(crnData)-1]

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
	err1 := kpAPI.DeleteKeyAlias(context.Background(), id[0], keyid)
	if err1 != nil {
		kpError := err1.(*kp.Error)
		if kpError.StatusCode == 404 {
			return nil
		} else {
			return fmt.Errorf(" failed to Destroy alias with error: %s", err1)
		}
	}
	return nil

}
