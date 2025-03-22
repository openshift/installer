// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kms

import (
	"context"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	kp "github.com/IBM/keyprotect-go-client"
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
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "If set to true, Key Protect enables a rotation policy on a single key.",
							Default:     true,
						},
						"interval_month": {
							Type:         schema.TypeInt,
							Optional:     true,
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
			flex.ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about the resource",
			},
		},
	}
}
func resourceIBMKmsKeyPolicyCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	instanceID := getInstanceIDFromCRN(d.Get("instance_id").(string))
	var id string
	if v, ok := d.GetOk("key_id"); ok {
		id = v.(string)
		// d.Set("key_id", id)
	}
	if v, ok := d.GetOk("alias"); ok {
		id = v.(string)
	}
	kpAPI, _, err := populateKPClient(d, meta, instanceID)
	if err != nil {
		return diag.FromErr(err)
	}
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
	_, instanceID, keyid := getInstanceAndKeyDataFromCRN(d.Id())
	kpAPI, _, err := populateKPClient(d, meta, instanceID)
	if err != nil {
		return diag.FromErr(err)
	}
	key, err := kpAPI.GetKey(context, keyid)
	if err != nil {
		if kpError, ok := err.(*kp.Error); ok {
			if kpError.StatusCode == 404 || kpError.StatusCode == 409 {
				d.SetId("")
				return nil
			}
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

		instanceID := getInstanceIDFromCRN(d.Get("instance_id").(string))
		kpAPI, _, err := populateKPClient(d, meta, instanceID)
		if err != nil {
			return diag.FromErr(err)
		}
		_, _, key_id := getInstanceAndKeyDataFromCRN(d.Id())

		err = resourceUpdatePolicies(context, d, kpAPI, meta, key_id)
		if err != nil {
			return diag.Errorf("Could not update policies: %s", err)
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

func resourceUpdatePolicies(context context.Context, d *schema.ResourceData, kpAPI *kp.Client, meta interface{}, key_id string) error {
	var dualAuthEnable, rotationEnable bool
	var rotationInterval int

	policy := getPolicyFromSchema(d)

	if policy.Rotation != nil {
		rotationInterval = policy.Rotation.Interval
		rotationEnable = *policy.Rotation.Enabled
		/* While updating a rotation policy, if the user does not set interval_month, it will be zero and the policy update
		will intend to only enable or disbale a policy. In case, the user inputs both values `enabled` and
		`interval_month`, policy update will update both with respective functions called from the SDK. */
		if rotationInterval == 0 {
			if rotationEnable {
				_, err := kpAPI.EnableRotationPolicy(context, key_id)
				if err != nil {
					return flex.FmtErrorf("[ERROR] Error while enabling key rotation policies: %s", err)
				}
			} else if !rotationEnable {
				_, err := kpAPI.DisableRotationPolicy(context, key_id)
				if err != nil {
					return flex.FmtErrorf("[ERROR] Error while disabling key rotation policies: %s", err)
				}
			}
		} else {
			_, err := kpAPI.SetRotationPolicy(context, key_id, rotationInterval, rotationEnable)
			if err != nil {
				return flex.FmtErrorf("[ERROR] Error while disabling key rotation policies: %s", err)
			}
		}
	}
	if policy.DualAuth != nil {
		dualAuthEnable = *policy.DualAuth.Enabled
		_, err := kpAPI.SetDualAuthDeletePolicy(context, key_id, dualAuthEnable)
		if err != nil {
			return flex.FmtErrorf("[ERROR] Error while setting dual_auth_delete policies: %s", err)
		}
	}
	return nil
}

func resourceHandlePolicies(context context.Context, d *schema.ResourceData, kpAPI *kp.Client, meta interface{}, key_id string) error {
	var setRotation, setDualAuthDelete, dualAuthEnable, rotationEnable bool
	var rotationInterval int

	policy := getPolicyFromSchema(d)

	if policy.Rotation != nil {
		setRotation = true
		rotationInterval = policy.Rotation.Interval
		rotationEnable = *policy.Rotation.Enabled
	}
	if policy.DualAuth != nil {
		setDualAuthDelete = true
		dualAuthEnable = *policy.DualAuth.Enabled
	}
	_, err := kpAPI.SetPolicies(context, key_id, setRotation, rotationInterval, setDualAuthDelete, dualAuthEnable, rotationEnable)
	if err != nil {
		return flex.FmtErrorf("[ERROR] Error while creating policies: %s", err)
	}
	return nil
}

func getPolicyFromSchema(d *schema.ResourceData) kp.Policy {
	var policy kp.Policy
	if rotationPolicyInfo, ok := d.GetOk("rotation"); ok {
		rotationPolicyList := rotationPolicyInfo.([]interface{})
		if len(rotationPolicyList) != 0 {
			rotationPolicyMap := rotationPolicyList[0].(map[string]interface{})
			policy.Rotation = &kp.Rotation{
				Interval: rotationPolicyMap["interval_month"].(int),
			}
			if _, ok := rotationPolicyMap["enabled"]; ok {
				enabled := rotationPolicyMap["enabled"].(bool)
				policy.Rotation.Enabled = &enabled
			}
		}
	}
	if dualAuthPolicyInfo, ok := d.GetOk("dual_auth_delete"); ok {
		dualAuthPolicyList := dualAuthPolicyInfo.([]interface{})
		if len(dualAuthPolicyList) != 0 {
			enabled := dualAuthPolicyList[0].(map[string]interface{})["enabled"].(bool)
			policy.DualAuth = &kp.DualAuth{
				Enabled: &enabled,
			}
		}
	}
	return policy
}
