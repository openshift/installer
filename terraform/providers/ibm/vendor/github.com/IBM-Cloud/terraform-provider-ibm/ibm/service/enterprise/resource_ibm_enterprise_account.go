// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package enterprise

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/platform-services-go-sdk/enterprisemanagementv1"
)

func ResourceIBMEnterpriseAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmEnterpriseAccountCreate,
		ReadContext:   resourceIbmEnterpriseAccountRead,
		UpdateContext: resourceIbmEnterpriseAccountUpdate,
		DeleteContext: resourceIbmEnterpriseAccountDelete,
		Importer:      &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"parent": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The CRN of the parent under which the account will be created. The parent can be an existing account group or the enterprise itself.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The name of the account. This field must have 3 - 60 characters.",
				ForceNew:     true,
				ValidateFunc: validate.ValidateAllowedEnterpriseNameValue(),
			},
			"owner_iam_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The IAM ID of the account owner, such as `IBMid-0123ABC`. The IAM ID must already exist.",
				ForceNew:    true,
			},
			"traits": {
				Type:             schema.TypeSet,
				Description:      "The traits object can be used to set properties on child accounts of an enterprise. You can pass a field to opt-out of Multi-Factor Authentication setting or setup enterprise IAM settings when creating a child account in the enterprise. This is an optional field.",
				Optional:         true,
				DiffSuppressFunc: flex.ApplyOnce,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mfa": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "By default MFA will be enabled on a child account. To opt out, pass the traits object with the mfa field set to empty string. This is an optional field.",
						},
						"enterprise_iam_managed": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "The Enterprise IAM settings property will be turned off for a newly created child account by default. You can enable this property by passing 'true' in this boolean field. This is an optional field.",
						},
					},
				},
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the account.",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Cloud Resource Name (CRN) of the account.",
			},
			"enterprise_account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The enterprise account ID.",
			},
			"enterprise_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The enterprise ID that the account is a part of.",
			},
			"account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The source account id of account to be imported",
			},
			"enterprise_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The path from the enterprise to this particular account.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The state of the account.",
			},
			"paid": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The type of account - whether it is free or paid.",
			},
			"owner_email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The email address of the owner of the account.",
			},
			"is_enterprise_account": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The flag to indicate whether the account is an enterprise account or not.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time stamp at which the account was created.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IAM ID of the user or service that created the account.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time stamp at which the account was last updated.",
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IAM ID of the user or service that updated the account.",
			},
		},
	}
}

func checkImportAccount(d *schema.ResourceData) bool {
	_, validateEnterpriseAccountId := d.GetOk("account_id")
	_, validateEnterpriseId := d.GetOk("enterprise_id")
	if validateEnterpriseAccountId && validateEnterpriseId {
		return true
	}
	return false
}

func checkCreateAccount(d *schema.ResourceData) bool {
	_, validateParent := d.GetOk("parent")
	_, validateName := d.GetOk("name")
	_, validateOwnerIamId := d.GetOk("owner_iam_id")
	if validateParent && validateName && validateOwnerIamId {
		return true
	}
	return false
}

func resourceIbmEnterpriseAccountCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enterpriseManagementClient, err := meta.(conns.ClientSession).EnterpriseManagementV1()
	if err != nil {
		return diag.FromErr(err)
	}

	if checkImportAccount(d) {
		importAccountToEnterpriseOptions := &enterprisemanagementv1.ImportAccountToEnterpriseOptions{}
		importAccountToEnterpriseOptions.SetEnterpriseID(d.Get("enterprise_id").(string))
		importAccountToEnterpriseOptions.SetAccountID(d.Get("account_id").(string))
		response, err := enterpriseManagementClient.ImportAccountToEnterpriseWithContext(context, importAccountToEnterpriseOptions)
		if err != nil {
			log.Printf("[DEBUG] ImportAccountToEnterpriseWithContext failed %s\n%s", err, response)
			return diag.FromErr(err)
		}
		d.SetId(d.Get("account_id").(string))
	} else if checkCreateAccount(d) {
		createAccountOptions := &enterprisemanagementv1.CreateAccountOptions{}
		createAccountOptions.SetParent(d.Get("parent").(string))
		createAccountOptions.SetName(d.Get("name").(string))
		createAccountOptions.SetOwnerIamID(d.Get("owner_iam_id").(string))
		if _, ok := d.GetOk("Traits"); ok {
			createAccountOptions.SetTraits(d.Get("traits").(*enterprisemanagementv1.CreateAccountRequestTraits))
		}
		createAccountResponse, response, err := enterpriseManagementClient.CreateAccountWithContext(context, createAccountOptions)
		if err != nil {
			log.Printf("[DEBUG] CreateAccountWithContext failed %s\n%s", err, response)
			return diag.FromErr(err)
		}
		d.SetId(*createAccountResponse.AccountID)
	} else {

		err := errors.New("[ERROR] Required Parameters are missing." +
			"Please input parent,name,owner_iam_id for creating a new account in enterprise." +
			"Input enterprise_id and enterprise_account_id for importing an existing account to enterprise")
		return diag.FromErr(err)
	}
	return resourceIbmEnterpriseAccountRead(context, d, meta)
}

func resourceIbmEnterpriseAccountRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enterpriseManagementClient, err := meta.(conns.ClientSession).EnterpriseManagementV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getAccountOptions := &enterprisemanagementv1.GetAccountOptions{}

	getAccountOptions.SetAccountID(d.Id())

	account, response, err := enterpriseManagementClient.GetAccountWithContext(context, getAccountOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetAccountWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	if err = d.Set("parent", account.Parent); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting parent: %s", err))
	}
	if err = d.Set("name", account.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}
	if err = d.Set("owner_iam_id", account.OwnerIamID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting owner_iam_id: %s", err))
	}
	if err = d.Set("account_id", account.ID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting account_id: %s", err))
	}
	if err = d.Set("url", account.URL); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting url: %s", err))
	}

	if err = d.Set("crn", account.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting crn: %s", err))
	}
	if err = d.Set("enterprise_account_id", account.EnterpriseAccountID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting enterprise_account_id: %s", err))
	}
	if err = d.Set("enterprise_id", account.EnterpriseID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting enterprise_id: %s", err))
	}
	if err = d.Set("enterprise_path", account.EnterprisePath); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting enterprise_path: %s", err))
	}
	if err = d.Set("state", account.State); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting state: %s", err))
	}
	if err = d.Set("paid", account.Paid); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting paid: %s", err))
	}
	if err = d.Set("owner_email", account.OwnerEmail); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting owner_email: %s", err))
	}
	if err = d.Set("is_enterprise_account", account.IsEnterpriseAccount); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting is_enterprise_account: %s", err))
	}
	if err = d.Set("created_at", account.CreatedAt.String()); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
	}
	if err = d.Set("created_by", account.CreatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_by: %s", err))
	}
	if account.UpdatedAt != nil {
		if err = d.Set("updated_at", account.UpdatedAt.String()); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting updated_at: %s", err))
		}
	}
	if account.UpdatedBy != nil {
		if err = d.Set("updated_by", account.UpdatedBy); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting updated_by: %s", err))
		}
	}
	return nil
}

func resourceIbmEnterpriseAccountUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enterpriseManagementClient, err := meta.(conns.ClientSession).EnterpriseManagementV1()
	if err != nil {
		return diag.FromErr(err)
	}

	updateAccountOptions := &enterprisemanagementv1.UpdateAccountOptions{}

	updateAccountOptions.SetAccountID(d.Id())

	hasChange := false

	if d.HasChange("parent") {
		updateAccountOptions.SetParent(d.Get("parent").(string))
		hasChange = true
	}
	/** Removed as update call requires only parent **/
	//if d.HasChange("name") {
	//
	//	updateAccountOptions.SetName(d.Get("name").(string))
	//	hasChange = true
	//}
	//if d.HasChange("owner_iam_id") {
	//	updateAccountOptions.SetOwnerIamID(d.Get("owner_iam_id").(string))
	//	hasChange = true
	//}

	if hasChange {
		response, err := enterpriseManagementClient.UpdateAccountWithContext(context, updateAccountOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateAccountWithContext failed %s\n%s", err, response)
			return diag.FromErr(err)
		}
	}

	return resourceIbmEnterpriseAccountRead(context, d, meta)
}

func resourceIbmEnterpriseAccountDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	enterpriseManagementClient, err := meta.(conns.ClientSession).EnterpriseManagementV1()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteAccountOptions := &enterprisemanagementv1.DeleteAccountOptions{}

	deleteAccountOptions.SetAccountID(d.Id())

	response, err := enterpriseManagementClient.DeleteAccountWithContext(context, deleteAccountOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteAccountWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteAccountWithContext failed %s\n%s", err, response))
	}

	return nil
}
