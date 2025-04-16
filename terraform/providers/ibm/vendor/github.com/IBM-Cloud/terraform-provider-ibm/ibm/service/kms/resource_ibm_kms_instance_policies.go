// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kms

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	kp "github.com/IBM/keyprotect-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMKmsInstancePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMKmsInstancePolicyCreate,
		ReadContext:   resourceIBMKmsInstancePoliciesRead,
		UpdateContext: resourceIBMKmsInstancePolicyUpdate,
		DeleteContext: resourceIBMKmsInstancePolicyDelete,
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
				Description:      "Key protect or hpcs instance GUID or CRN",
				DiffSuppressFunc: suppressKMSInstanceIDDiff,
			},
			"dual_auth_delete": {
				Type:         schema.TypeList,
				Optional:     true,
				AtLeastOneOf: []string{"rotation", "dual_auth_delete", "metrics", "key_create_import_access"},
				MaxItems:     1,
				Description:  "Data associated with the dual authorization delete policy for instance",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "If set to true, Key Protect enables a dual authorization policy for the instance.",
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
						"last_updated": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Updates when the policy is replaced or modified. The date format follows RFC 3339.",
						},
					},
				},
			},
			"rotation": {
				Type:         schema.TypeList,
				Optional:     true,
				AtLeastOneOf: []string{"rotation", "dual_auth_delete", "metrics", "key_create_import_access"},
				MaxItems:     1,
				Description:  "Data associated with the rotation policy for instance",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "If set to true, Key Protect enables a rotation policy for the instance.",
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
						"last_updated": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Updates when the policy is replaced or modified. The date format follows RFC 3339.",
						},
						"interval_month": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "Specifies the rotation time interval in months for the instance.",
							ValidateFunc: validate.ValidateAllowedRangeInt(1, 12),
						},
					},
				},
			},
			"key_create_import_access": {
				Type:         schema.TypeList,
				Optional:     true,
				AtLeastOneOf: []string{"rotation", "dual_auth_delete", "metrics", "key_create_import_access"},
				MaxItems:     1,
				Description:  "Data associated with the key create import access policy for the instance",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "If set to true, Key Protect enables a KCIA policy for the instance.",
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
						"last_updated": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Updates when the policy is replaced or modified. The date format follows RFC 3339.",
						},
						"create_root_key": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "If set to true, Key Protect allows you or any authorized users to create root keys in the instance.",
							Default:     true,
						},
						"create_standard_key": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "If set to true, Key Protect allows you or any authorized users to create standard keys in the instance.",
							Default:     true,
						},
						"import_root_key": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "If set to true, Key Protect allows you or any authorized users to import root keys into the instance.",
							Default:     true,
						},
						"import_standard_key": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "If set to true, Key Protect allows you or any authorized users to import standard keys into the instance.",
							Default:     true,
						},
						"enforce_token": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "If set to true, the service prevents you or any authorized users from importing key material into the specified service instance without using an import token.",
							Default:     false,
						},
					},
				},
			},
			"metrics": {
				Type:         schema.TypeList,
				Optional:     true,
				AtLeastOneOf: []string{"rotation", "dual_auth_delete", "metrics", "key_create_import_access"},
				MaxItems:     1,
				Description:  "Data associated with the metric policy for instance",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "If set to true, Key Protect enables a metrics policy on the instance.",
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
						"last_updated": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Updates when the policy is replaced or modified. The date format follows RFC 3339.",
						},
					},
				},
			},
		},
	}
}

func resourceIBMKmsInstancePolicyCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	instanceID := getInstanceIDFromCRN(d.Get("instance_id").(string))
	kpAPI, instanceCRN, err := populateKPClient(d, meta, instanceID)
	if err != nil {
		return diag.FromErr(err)
	}
	policyCreateOrUpdate(context, d, kpAPI)
	d.SetId(*instanceCRN)
	return resourceIBMKmsInstancePoliciesRead(context, d, meta)
}

func resourceIBMKmsInstancePoliciesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, instanceID, _ := getInstanceAndKeyDataFromCRN(d.Id())
	kpAPI, _, err := populateKPClient(d, meta, instanceID)
	if err != nil {
		return diag.FromErr(err)
	}
	instancePolicies, err := kpAPI.GetInstancePolicies(context)
	if err != nil {
		return diag.Errorf("[ERROR] Get Policies failed with error : %s", err)
	}
	d.Set("instance_id", instanceID)
	d.Set("dual_auth_delete", flex.FlattenInstancePolicy("dual_auth_delete", instancePolicies))
	d.Set("rotation", flex.FlattenInstancePolicy("rotation", instancePolicies))
	d.Set("metrics", flex.FlattenInstancePolicy("metrics", instancePolicies))
	d.Set("key_create_import_access", flex.FlattenInstancePolicy("key_create_import_access", instancePolicies))
	return nil

}

func resourceIBMKmsInstancePolicyUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	if d.HasChange("rotation") || d.HasChange("dual_auth_delete") || d.HasChange("metric") || d.HasChange("key_create_import_access") {

		instanceID := getInstanceIDFromCRN(d.Get("instance_id").(string))
		kpAPI, _, err := populateKPClient(d, meta, instanceID)
		if err != nil {
			return diag.FromErr(err)
		}

		err = policyCreateOrUpdate(context, d, kpAPI)
		if err != nil {
			return diag.Errorf("Could not update the policies: %s", err)
		}
	}
	return resourceIBMKmsInstancePoliciesRead(context, d, meta)
}

func resourceIBMKmsInstancePolicyDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//Do not support delete Policies
	log.Println("Warning:  `terraform destroy` does not remove the policies of the Instance but only clears the state file. Instance Policies get deleted when the associated instance resource is destroyed.")
	d.SetId("")
	return nil

}

func policyCreateOrUpdate(context context.Context, d *schema.ResourceData, kpAPI *kp.Client) error {
	var mulPolicy kp.MultiplePolicies
	if dualAuthDeleteInstancePolicy, ok := d.GetOk("dual_auth_delete"); ok {
		dualAuthDeleteInstancePolicyList := dualAuthDeleteInstancePolicy.([]interface{})
		if len(dualAuthDeleteInstancePolicyList) != 0 {
			mulPolicy.DualAuthDelete = &kp.BasicPolicyData{
				Enabled: dualAuthDeleteInstancePolicyList[0].(map[string]interface{})["enabled"].(bool),
			}
		}
	}
	if rotationInstancePolicy, ok := d.GetOk("rotation"); ok {
		rotationInstancePolicyList := rotationInstancePolicy.([]interface{})
		if len(rotationInstancePolicyList) != 0 {
			iM := rotationInstancePolicyList[0].(map[string]interface{})["interval_month"].(int)
			enabled := rotationInstancePolicyList[0].(map[string]interface{})["enabled"].(bool)
			//For case when enabled = false && no input to interval month.
			if iM == 0 {
				mulPolicy.Rotation = &kp.RotationPolicyData{
					Enabled:       enabled,
					IntervalMonth: nil,
				}
			} else {
				mulPolicy.Rotation = &kp.RotationPolicyData{
					Enabled:       enabled,
					IntervalMonth: &iM,
				}
			}

		}
	}
	if metricsInstancePolicy, ok := d.GetOk("metrics"); ok {
		metricsInstancePolicyList := metricsInstancePolicy.([]interface{})
		if len(metricsInstancePolicyList) != 0 {
			mulPolicy.Metrics = &kp.BasicPolicyData{
				Enabled: metricsInstancePolicyList[0].(map[string]interface{})["enabled"].(bool),
			}
		}
	}
	if kciaip, ok := d.GetOk("key_create_import_access"); ok {
		kciaipList := kciaip.([]interface{})
		if len(kciaipList) != 0 {
			mulPolicy.KeyCreateImportAccess = &kp.KeyCreateImportAccessInstancePolicy{
				Enabled:           kciaipList[0].(map[string]interface{})["enabled"].(bool),
				CreateRootKey:     kciaipList[0].(map[string]interface{})["create_root_key"].(bool),
				CreateStandardKey: kciaipList[0].(map[string]interface{})["create_standard_key"].(bool),
				ImportRootKey:     kciaipList[0].(map[string]interface{})["import_root_key"].(bool),
				ImportStandardKey: kciaipList[0].(map[string]interface{})["import_standard_key"].(bool),
				EnforceToken:      kciaipList[0].(map[string]interface{})["enforce_token"].(bool),
			}
		}
	}
	err := kpAPI.SetInstancePolicies(context, mulPolicy)
	if err != nil {
		return fmt.Errorf("[ERROR] Error while setting instance policies: %s", err)
	}
	return nil

}
