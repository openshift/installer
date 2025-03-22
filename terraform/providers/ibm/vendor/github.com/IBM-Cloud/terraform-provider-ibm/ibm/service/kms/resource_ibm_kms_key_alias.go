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

func ResourceIBMKmskeyAlias() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMKmsKeyAliasCreate,
		Delete:   resourceIBMKmsKeyAliasDelete,
		Read:     resourceIBMKmsKeyAliasRead,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Key ID",
				ForceNew:         true,
				DiffSuppressFunc: suppressKMSInstanceIDDiff,
			},
			"alias": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Key protect or hpcs key alias name",
			},
			"key_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Key ID",
				ForceNew:     true,
				ExactlyOneOf: []string{"key_id", "existing_alias"},
			},
			"existing_alias": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Existing Alias of the Key",
				ForceNew:     true,
				ExactlyOneOf: []string{"key_id", "existing_alias"},
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

func resourceIBMKmsKeyAliasCreate(d *schema.ResourceData, meta interface{}) error {
	instanceID := getInstanceIDFromCRN(d.Get("instance_id").(string))
	kpAPI, _, err := populateKPClient(d, meta, instanceID)
	if err != nil {
		return err
	}

	aliasName := d.Get("alias").(string)
	var id string
	if v, ok := d.GetOk("key_id"); ok {
		id = v.(string)
		d.Set("key_id", id)
	}
	if v, ok := d.GetOk("existing_alias"); ok {
		id = v.(string)
	}
	stkey, err := kpAPI.CreateKeyAlias(context.Background(), aliasName, id)
	if err != nil {
		return flex.FmtErrorf("[ERROR] Error while creating alias name for the key: %s", err)
	}
	key, err := kpAPI.GetKey(context.Background(), stkey.KeyID)
	if err != nil {
		return flex.FmtErrorf("[ERROR] Get Key failed with error: %s", err)
	}
	d.SetId(fmt.Sprintf("%s:alias:%s", stkey.Alias, key.CRN))

	return resourceIBMKmsKeyAliasRead(d, meta)
}

func resourceIBMKmsKeyAliasRead(d *schema.ResourceData, meta interface{}) error {
	id := strings.Split(d.Id(), ":alias:")
	if len(id) < 2 {
		return flex.FmtErrorf("[ERROR] Incorrect ID %s: Id should be a combination of keyAlias:alias:keyCRN", d.Id())
	}
	_, instanceID, keyid := getInstanceAndKeyDataFromCRN(id[1])
	kpAPI, _, err := populateKPClient(d, meta, instanceID)
	if err != nil {
		return err
	}
	key, err := kpAPI.GetKey(context.Background(), keyid)
	if err != nil {
		if kpError, ok := err.(*kp.Error); ok {
			if kpError.StatusCode == 404 || kpError.StatusCode == 409 {
				d.SetId("")
				return nil
			}
		}
		return flex.FmtErrorf("[ERROR] Get Key failed with error while reading policies: %s", err)
	} else if key.State == 5 { //Refers to Deleted state of the Key
		d.SetId("")
		return nil
	}
	d.Set("alias", id[0])
	d.Set("key_id", key.ID)
	d.Set("instance_id", instanceID)
	if strings.Contains((kpAPI.URL).String(), "private") || strings.Contains(kpAPI.Config.BaseURL, "private") {
		d.Set("endpoint_type", "private")
	} else {
		d.Set("endpoint_type", "public")
	}

	return nil
}

func resourceIBMKmsKeyAliasDelete(d *schema.ResourceData, meta interface{}) error {
	id := strings.Split(d.Id(), ":alias:")
	_, instanceID, keyid := getInstanceAndKeyDataFromCRN(id[1])
	kpAPI, _, err := populateKPClient(d, meta, instanceID)
	if err != nil {
		return err
	}
	err1 := kpAPI.DeleteKeyAlias(context.Background(), id[0], keyid)
	if err1 != nil {
		if kpError, ok := err1.(*kp.Error); ok {
			if kpError.StatusCode == 404 {
				return nil
			}
		}
		return flex.FmtErrorf(" failed to Destroy alias with error: %s", err1)
	}
	return nil
}
