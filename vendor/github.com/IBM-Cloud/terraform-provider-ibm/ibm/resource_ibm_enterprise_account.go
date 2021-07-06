// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"errors"
	"fmt"

	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/IBM/platform-services-go-sdk/enterprisemanagementv1"
)

func resourceIbmEnterpriseAccount() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIbmEnterpriseAccountCreate,
		Read:     resourceIbmEnterpriseAccountRead,
		Update:   resourceIbmEnterpriseAccountUpdate,
		Delete:   resourceIbmEnterpriseAccountDelete,
		Importer: &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"parent": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The CRN of the parent under which the account will be created. The parent can be an existing account group or the enterprise itself.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The name of the account. This field must have 3 - 60 characters.",
				ForceNew:     true,
				ValidateFunc: validateAllowedEnterpriseNameValue(),
			},
			"owner_iam_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The IAM ID of the account owner, such as `IBMid-0123ABC`. The IAM ID must already exist.",
				ForceNew:    true,
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the account.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Cloud Resource Name (CRN) of the account.",
			},
			"enterprise_account_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The enterprise account ID.",
			},
			"enterprise_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The enterprise ID that the account is a part of.",
			},
			"account_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The source account id of account to be imported",
			},
			"enterprise_path": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The path from the enterprise to this particular account.",
			},
			"state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The state of the account.",
			},
			"paid": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The type of account - whether it is free or paid.",
			},
			"owner_email": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The email address of the owner of the account.",
			},
			"is_enterprise_account": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The flag to indicate whether the account is an enterprise account or not.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time stamp at which the account was created.",
			},
			"created_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IAM ID of the user or service that created the account.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time stamp at which the account was last updated.",
			},
			"updated_by": &schema.Schema{
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

func resourceIbmEnterpriseAccountCreate(d *schema.ResourceData, meta interface{}) error {
	enterpriseManagementClient, err := meta.(ClientSession).EnterpriseManagementV1()
	if err != nil {
		return err
	}

	if checkImportAccount(d) {
		importAccountToEnterpriseOptions := &enterprisemanagementv1.ImportAccountToEnterpriseOptions{}
		importAccountToEnterpriseOptions.SetEnterpriseID(d.Get("enterprise_id").(string))
		importAccountToEnterpriseOptions.SetAccountID(d.Get("account_id").(string))
		response, err := enterpriseManagementClient.ImportAccountToEnterpriseWithContext(context.TODO(), importAccountToEnterpriseOptions)
		if err != nil {
			log.Printf("[DEBUG] ImportAccountToEnterpriseWithContext failed %s\n%s", err, response)
			return err
		}
		d.SetId(d.Get("account_id").(string))
	} else if checkCreateAccount(d) {
		createAccountOptions := &enterprisemanagementv1.CreateAccountOptions{}
		createAccountOptions.SetParent(d.Get("parent").(string))
		createAccountOptions.SetName(d.Get("name").(string))
		createAccountOptions.SetOwnerIamID(d.Get("owner_iam_id").(string))
		createAccountResponse, response, err := enterpriseManagementClient.CreateAccountWithContext(context.TODO(), createAccountOptions)
		if err != nil {
			log.Printf("[DEBUG] CreateAccountWithContext failed %s\n%s", err, response)
			return err
		}
		d.SetId(*createAccountResponse.AccountID)
	} else {

		err := errors.New("Required Parameters are missing." +
			"Please input parent,name,owner_iam_id for creating a new account in enterprise." +
			"Input enterprise_id and enterprise_account_id for importing an existing account to enterprise.")
		return err
	}
	return resourceIbmEnterpriseAccountRead(d, meta)
}

func resourceIbmEnterpriseAccountRead(d *schema.ResourceData, meta interface{}) error {
	enterpriseManagementClient, err := meta.(ClientSession).EnterpriseManagementV1()
	if err != nil {
		return err
	}

	getAccountOptions := &enterprisemanagementv1.GetAccountOptions{}

	getAccountOptions.SetAccountID(d.Id())

	account, response, err := enterpriseManagementClient.GetAccountWithContext(context.TODO(), getAccountOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetAccountWithContext failed %s\n%s", err, response)
		return err
	}

	if err = d.Set("parent", account.Parent); err != nil {
		return fmt.Errorf("Error setting parent: %s", err)
	}
	if err = d.Set("name", account.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err = d.Set("owner_iam_id", account.OwnerIamID); err != nil {
		return fmt.Errorf("Error setting owner_iam_id: %s", err)
	}
	if err = d.Set("account_id", account.ID); err != nil {
		return fmt.Errorf("Error setting account_id: %s", err)
	}
	if err = d.Set("url", account.URL); err != nil {
		return fmt.Errorf("Error setting url: %s", err)
	}

	if err = d.Set("crn", account.CRN); err != nil {
		return fmt.Errorf("Error setting crn: %s", err)
	}
	if err = d.Set("enterprise_account_id", account.EnterpriseAccountID); err != nil {
		return fmt.Errorf("Error setting enterprise_account_id: %s", err)
	}
	if err = d.Set("enterprise_id", account.EnterpriseID); err != nil {
		return fmt.Errorf("Error setting enterprise_id: %s", err)
	}
	if err = d.Set("enterprise_path", account.EnterprisePath); err != nil {
		return fmt.Errorf("Error setting enterprise_path: %s", err)
	}
	if err = d.Set("state", account.State); err != nil {
		return fmt.Errorf("Error setting state: %s", err)
	}
	if err = d.Set("paid", account.Paid); err != nil {
		return fmt.Errorf("Error setting paid: %s", err)
	}
	if err = d.Set("owner_email", account.OwnerEmail); err != nil {
		return fmt.Errorf("Error setting owner_email: %s", err)
	}
	if err = d.Set("is_enterprise_account", account.IsEnterpriseAccount); err != nil {
		return fmt.Errorf("Error setting is_enterprise_account: %s", err)
	}
	if err = d.Set("created_at", account.CreatedAt.String()); err != nil {
		return fmt.Errorf("Error setting created_at: %s", err)
	}
	if err = d.Set("created_by", account.CreatedBy); err != nil {
		return fmt.Errorf("Error setting created_by: %s", err)
	}
	if account.UpdatedAt != nil {
		if err = d.Set("updated_at", account.UpdatedAt.String()); err != nil {
			return fmt.Errorf("Error setting updated_at: %s", err)
		}
	}
	if account.UpdatedBy != nil {
		if err = d.Set("updated_by", account.UpdatedBy); err != nil {
			return fmt.Errorf("Error setting updated_by: %s", err)
		}
	}
	return nil
}

func resourceIbmEnterpriseAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	enterpriseManagementClient, err := meta.(ClientSession).EnterpriseManagementV1()
	if err != nil {
		return err
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
		response, err := enterpriseManagementClient.UpdateAccountWithContext(context.TODO(), updateAccountOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateAccountWithContext failed %s\n%s", err, response)
			return err
		}
	}

	return resourceIbmEnterpriseAccountRead(d, meta)
}

func resourceIbmEnterpriseAccountDelete(d *schema.ResourceData, meta interface{}) error {

	d.SetId("")

	return nil
}
