package ibm

import (
	"context"
	"fmt"
	"strings"

	kp "github.com/IBM/keyprotect-go-client"
	rc "github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIBMKmskeyAlias() *schema.Resource {
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
				Type:        schema.TypeString,
				Required:    true,
				Description: "Key ID",
				ForceNew:    true,
			},
			"endpoint_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateAllowedStringValue([]string{"public", "private"}),
				Description:  "public or private",
				ForceNew:     true,
			},
		},
	}
}

func resourceIBMKmsKeyAliasCreate(d *schema.ResourceData, meta interface{}) error {
	kpAPI, err := meta.(ClientSession).keyManagementAPI()
	if err != nil {
		return err
	}

	instanceID := d.Get("instance_id").(string)
	CrnInstanceID := strings.Split(instanceID, ":")
	if len(CrnInstanceID) > 3 {
		instanceID = CrnInstanceID[len(CrnInstanceID)-3]
	}
	endpointType := d.Get("endpoint_type").(string)

	rsConClient, err := meta.(ClientSession).ResourceControllerV2API()
	if err != nil {
		return err
	}
	resourceInstanceGet := rc.GetResourceInstanceOptions{
		ID: &instanceID,
	}

	instanceData, resp, err := rsConClient.GetResourceInstance(&resourceInstanceGet)
	if err != nil || instanceData == nil {
		return fmt.Errorf("[ERROR] Error retrieving resource instance: %s with resp code: %s", err, resp)
	}
	extensions := instanceData.Extensions
	URL, err := KmsEndpointURL(kpAPI, endpointType, extensions)
	if err != nil {
		return err
	}
	kpAPI.URL = URL
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
	if len(id) < 2 {
		return fmt.Errorf("Incorrect ID %s: Id should be a combination of keyAlias:alias:keyCRN", d.Id())
	}
	crn := id[1]
	crnData := strings.Split(crn, ":")
	endpointType := d.Get("endpoint_type").(string)
	instanceID := crnData[len(crnData)-3]
	keyid := crnData[len(crnData)-1]

	rsConClient, err := meta.(ClientSession).ResourceControllerV2API()
	if err != nil {
		return err
	}
	resourceInstanceGet := rc.GetResourceInstanceOptions{
		ID: &instanceID,
	}

	instanceData, resp, err := rsConClient.GetResourceInstance(&resourceInstanceGet)
	if err != nil || instanceData == nil {
		return fmt.Errorf("[ERROR] Error retrieving resource instance: %s with resp code: %s", err, resp)
	}
	extensions := instanceData.Extensions
	URL, err := KmsEndpointURL(kpAPI, endpointType, extensions)
	if err != nil {
		return err
	}
	kpAPI.URL = URL
	kpAPI.Config.InstanceID = instanceID
	key, err := kpAPI.GetKey(context.Background(), keyid)
	if err != nil {
		kpError := err.(*kp.Error)
		if kpError.StatusCode == 404 || kpError.StatusCode == 409 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Get Key failed with error while reading policies: %s", err)
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
	kpAPI, err := meta.(ClientSession).keyManagementAPI()
	if err != nil {
		return err
	}
	id := strings.Split(d.Id(), ":alias:")
	crn := id[1]
	crnData := strings.Split(crn, ":")
	endpointType := d.Get("endpoint_type").(string)
	instanceID := crnData[len(crnData)-3]
	keyid := crnData[len(crnData)-1]

	rsConClient, err := meta.(ClientSession).ResourceControllerV2API()
	if err != nil {
		return err
	}
	resourceInstanceGet := rc.GetResourceInstanceOptions{
		ID: &instanceID,
	}

	instanceData, resp, err := rsConClient.GetResourceInstance(&resourceInstanceGet)
	if err != nil || instanceData == nil {
		return fmt.Errorf("[ERROR] Error retrieving resource instance: %s with resp code: %s", err, resp)
	}
	extensions := instanceData.Extensions
	URL, err := KmsEndpointURL(kpAPI, endpointType, extensions)
	if err != nil {
		return err
	}
	kpAPI.URL = URL
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
