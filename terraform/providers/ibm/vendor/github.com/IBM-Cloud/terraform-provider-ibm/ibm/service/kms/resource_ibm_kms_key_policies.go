// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kms

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	kp "github.com/IBM/keyprotect-go-client"
	rc "github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMKmskeyPolicies() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMKmsKeyPolicyCreate,
		ReadContext:   resourceIBMKmsKeyPolicyRead,
		UpdateContext: resourceIBMKmsKeyPolicyUpdate,
		DeleteContext: resourceIBMKmsKeyPolicyDelete,
		Importer:      &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "Key protect or hpcs instance GUID",
				DiffSuppressFunc: suppressKMSInstanceIDDiff,
			},
			"key_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Key ID",
				ExactlyOneOf: []string{"key_id", "alias"},
			},
			"alias": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"key_id", "alias"},
			},
			"endpoint_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"public", "private"}),
				Description:  "public or private",
				ForceNew:     true,
				Default:      "public",
			},
			"rotation": {
				Type:         schema.TypeList,
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"rotation", "dual_auth_delete"},
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
							ValidateFunc: validate.ValidateAllowedRangeInt(1, 12),
							Description:  "Specifies the key rotation time interval in months",
						},
					},
				},
			},
			"dual_auth_delete": {
				Type:         schema.TypeList,
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"rotation", "dual_auth_delete"},
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
			flex.ResourceName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource",
			},
			flex.ResourceCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},
			flex.ResourceStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the resource",
			},
		},
	}
}
func resourceIBMKmsKeyPolicyCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	kpAPI, err := meta.(conns.ClientSession).KeyManagementAPI()
	if err != nil {
		return diag.FromErr(err)
	}
	instanceID := d.Get("instance_id").(string)
	CrnInstanceID := strings.Split(instanceID, ":")
	if len(CrnInstanceID) > 3 {
		instanceID = CrnInstanceID[len(CrnInstanceID)-3]
	}
	endpointType := d.Get("endpoint_type").(string)
	var id string
	if v, ok := d.GetOk("key_id"); ok {
		id = v.(string)
		d.Set("key_id", id)
	}
	if v, ok := d.GetOk("alias"); ok {
		id = v.(string)
	}
	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return diag.FromErr(err)
	}
	resourceInstanceGet := rc.GetResourceInstanceOptions{
		ID: &instanceID,
	}
	instanceData, resp, err := rsConClient.GetResourceInstance(&resourceInstanceGet)
	if err != nil || instanceData == nil {
		return diag.Errorf("[ERROR] Error retrieving resource instance: %s with resp code: %s", err, resp)
	}
	extensions := instanceData.Extensions
	URL, err := KmsEndpointURL(kpAPI, endpointType, extensions)
	if err != nil {
		return diag.FromErr(err)
	}
	kpAPI.URL = URL

	kpAPI.Config.InstanceID = instanceID

	key, err := kpAPI.GetKey(context, id)
	if err != nil {
		return diag.Errorf("Get Key failed with error while creating policies: %s", err)
	}
	err = resourceHandlePolicies(context, d, kpAPI, meta, id)
	if err != nil {
		return diag.Errorf("Could not create policies: %s", err)
	}
	d.SetId(key.CRN)
	return resourceIBMKmsKeyPolicyRead(context, d, meta)
}

func resourceIBMKmsKeyPolicyRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	kpAPI, err := meta.(conns.ClientSession).KeyManagementAPI()
	if err != nil {
		return diag.FromErr(err)
	}
	crn := d.Id()
	crnData := strings.Split(crn, ":")
	endpointType := d.Get("endpoint_type").(string)
	instanceID := crnData[len(crnData)-3]
	keyid := crnData[len(crnData)-1]

	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return diag.FromErr(err)
	}
	resourceInstanceGet := rc.GetResourceInstanceOptions{
		ID: &instanceID,
	}
	instanceData, resp, err := rsConClient.GetResourceInstance(&resourceInstanceGet)
	if err != nil || instanceData == nil {
		return diag.Errorf("[ERROR] Error retrieving resource instance: %s with resp code: %s", err, resp)
	}
	extensions := instanceData.Extensions
	URL, err := KmsEndpointURL(kpAPI, endpointType, extensions)
	if err != nil {
		return diag.FromErr(err)
	}
	kpAPI.URL = URL

	kpAPI.Config.InstanceID = instanceID
	key, err := kpAPI.GetKey(context, keyid)
	if err != nil {
		kpError := err.(*kp.Error)
		if kpError.StatusCode == 404 || kpError.StatusCode == 409 {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Get Key failed with error while reading policies: %s", err)
	} else if key.State == 5 { //Refers to Deleted state of the Key
		d.SetId("")
		return nil
	}

	d.Set("instance_id", instanceID)
	d.Set("key_id", keyid)
	if strings.Contains((kpAPI.URL).String(), "private") {
		d.Set("endpoint_type", "private")
	} else {
		d.Set("endpoint_type", "public")
	}
	d.Set(flex.ResourceName, key.Name)
	d.Set(flex.ResourceCRN, key.CRN)
	state := key.State
	d.Set(flex.ResourceStatus, strconv.Itoa(state))
	rcontroller, err := flex.GetBaseController(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	id := key.ID
	crn1 := strings.TrimSuffix(key.CRN, ":key:"+id)

	d.Set(flex.ResourceControllerURL, rcontroller+"/services/kms/"+url.QueryEscape(crn1)+"%3A%3A")

	policies, err := kpAPI.GetPolicies(context, keyid)

	if err != nil {
		return diag.Errorf("Failed to read policies: %s", err)
	}
	if len(policies) == 0 {
		log.Printf("No Policy Configurations read\n")
	} else {
		d.Set("rotation", flex.FlattenKeyIndividualPolicy("rotation", policies))
		d.Set("dual_auth_delete", flex.FlattenKeyIndividualPolicy("dual_auth_delete", policies))
	}

	return nil

}

func resourceIBMKmsKeyPolicyUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	if d.HasChange("rotation") || d.HasChange("dual_auth_delete") {

		kpAPI, err := meta.(conns.ClientSession).KeyManagementAPI()
		if err != nil {
			return diag.FromErr(err)
		}

		instanceID := d.Get("instance_id").(string)
		CrnInstanceID := strings.Split(instanceID, ":")
		if len(CrnInstanceID) > 3 {
			instanceID = CrnInstanceID[len(CrnInstanceID)-3]
		}
		endpointType := d.Get("endpoint_type").(string)

		rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
		if err != nil {
			return diag.FromErr(err)
		}
		resourceInstanceGet := rc.GetResourceInstanceOptions{
			ID: &instanceID,
		}
		instanceData, resp, err := rsConClient.GetResourceInstance(&resourceInstanceGet)
		if err != nil || instanceData == nil {
			return diag.Errorf("[ERROR] Error retrieving resource instance: %s with resp code: %s", err, resp)
		}
		extensions := instanceData.Extensions
		URL, err := KmsEndpointURL(kpAPI, endpointType, extensions)
		if err != nil {
			return diag.FromErr(err)
		}
		kpAPI.URL = URL
		kpAPI.Config.InstanceID = instanceID

		crn := d.Id()
		crnData := strings.Split(crn, ":")
		key_id := crnData[len(crnData)-1]

		err = resourceHandlePolicies(context, d, kpAPI, meta, key_id)
		if err != nil {
			resourceIBMKmsKeyRead(d, meta)
			return diag.Errorf("Could not create policies: %s", err)
		}
	}
	return resourceIBMKmsKeyPolicyRead(context, d, meta)

}

func resourceIBMKmsKeyPolicyDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//Do not support delete Policies
	log.Println("Warning:  `terraform destroy` does not remove the policies of the Key but only clears the state file. Key Policies get deleted when the associated key resource is destroyed.")
	d.SetId("")
	return nil

}

func resourceHandlePolicies(context context.Context, d *schema.ResourceData, kpAPI *kp.Client, meta interface{}, key_id string) error {
	var setRotation, setDualAuthDelete, dualAuthEnable bool
	var rotationInterval int

	if policyInfo, ok := d.GetOk("rotation"); ok {
		rpdList := policyInfo.([]interface{})
		if len(rpdList) != 0 {
			rotationInterval = rpdList[0].(map[string]interface{})["interval_month"].(int)
			setRotation = true
		}
	}
	if dadp, ok := d.GetOk("dual_auth_delete"); ok {
		dadpList := dadp.([]interface{})
		if len(dadpList) != 0 {
			dualAuthEnable = dadpList[0].(map[string]interface{})["enabled"].(bool)
			setDualAuthDelete = true
		}
	}
	_, err := kpAPI.SetPolicies(context, key_id, setRotation, rotationInterval, setDualAuthDelete, dualAuthEnable)
	if err != nil {
		return fmt.Errorf("[ERROR] Error while creating policies: %s", err)
	}
	return nil
}
