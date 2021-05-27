// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0
package ibm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/IBM/platform-services-go-sdk/enterprisemanagementv1"
)

func resourceIbmEnterpriseAccountGroup() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIbmEnterpriseAccountGroupCreate,
		Read:     resourceIbmEnterpriseAccountGroupRead,
		Update:   resourceIbmEnterpriseAccountGroupUpdate,
		Delete:   resourceIbmEnterpriseAccountGroupDelete,
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
				Description: "The CRN of the parent under which the account group will be created. The parent can be an existing account group or the enterprise itself.",
				ForceNew:    true,
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the account group. This field must have 3 - 60 characters.",
				ValidateFunc: validateAllowedEnterpriseNameValue(),
			},
			"primary_contact_iam_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The IAM ID of the primary contact for this account group, such as `IBMid-0123ABC`. The IAM ID must already exist.",
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the account group.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Cloud Resource Name (CRN) of the account group.",
			},
			"enterprise_account_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The enterprise account ID.",
			},
			"enterprise_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The enterprise ID that the account group is a part of.",
			},
			"enterprise_path": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The path from the enterprise to this particular account group.",
			},
			"state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The state of the account group.",
			},
			"primary_contact_email": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The email address of the primary contact of the account group.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time stamp at which the account group was created.",
			},
			"created_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IAM ID of the user or service that created the account group.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time stamp at which the account group was last updated.",
			},
			"updated_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IAM ID of the user or service that updated the account group.",
			},
		},
	}
}

func resourceIbmEnterpriseAccountGroupCreate(d *schema.ResourceData, meta interface{}) error {
	enterpriseManagementClient, err := meta.(ClientSession).EnterpriseManagementV1()
	if err != nil {
		return err
	}

	createAccountGroupOptions := &enterprisemanagementv1.CreateAccountGroupOptions{}

	createAccountGroupOptions.SetParent(d.Get("parent").(string))
	createAccountGroupOptions.SetName(d.Get("name").(string))
	createAccountGroupOptions.SetPrimaryContactIamID(d.Get("primary_contact_iam_id").(string))

	createAccountGroupResponse, response, err := enterpriseManagementClient.CreateAccountGroupWithContext(context.TODO(), createAccountGroupOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateAccountGroupWithContext failed %s\n%s", err, response)
		return err
	}

	d.SetId(*createAccountGroupResponse.AccountGroupID)

	return resourceIbmEnterpriseAccountGroupRead(d, meta)
}

func resourceIbmEnterpriseAccountGroupRead(d *schema.ResourceData, meta interface{}) error {
	enterpriseManagementClient, err := meta.(ClientSession).EnterpriseManagementV1()
	if err != nil {
		return err
	}

	getAccountGroupOptions := &enterprisemanagementv1.GetAccountGroupOptions{}

	getAccountGroupOptions.SetAccountGroupID(d.Id())

	accountGroup, response, err := enterpriseManagementClient.GetAccountGroupWithContext(context.TODO(), getAccountGroupOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetAccountGroupWithContext failed %s\n%s", err, response)
		return err
	}
	log.Printf("[DEBUG] GetAccountGroupWithContext testing %s", response)
	if err = d.Set("parent", accountGroup.Parent); err != nil {
		return fmt.Errorf("Error setting parent: %s", err)
	}
	if err = d.Set("name", accountGroup.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err = d.Set("primary_contact_iam_id", accountGroup.PrimaryContactIamID); err != nil {
		return fmt.Errorf("Error setting primary_contact_iam_id: %s", err)
	}
	if err = d.Set("url", accountGroup.URL); err != nil {
		return fmt.Errorf("Error setting url: %s", err)
	}
	if err = d.Set("crn", accountGroup.CRN); err != nil {
		return fmt.Errorf("Error setting crn: %s", err)
	}
	if err = d.Set("enterprise_account_id", accountGroup.EnterpriseAccountID); err != nil {
		return fmt.Errorf("Error setting enterprise_account_id: %s", err)
	}
	if err = d.Set("enterprise_id", accountGroup.EnterpriseID); err != nil {
		return fmt.Errorf("Error setting enterprise_id: %s", err)
	}
	if err = d.Set("enterprise_path", accountGroup.EnterprisePath); err != nil {
		return fmt.Errorf("Error setting enterprise_path: %s", err)
	}
	if err = d.Set("state", accountGroup.State); err != nil {
		return fmt.Errorf("Error setting state: %s", err)
	}
	if err = d.Set("primary_contact_email", accountGroup.PrimaryContactEmail); err != nil {
		return fmt.Errorf("Error setting primary_contact_email: %s", err)
	}
	if err = d.Set("created_at", accountGroup.CreatedAt.String()); err != nil {
		return fmt.Errorf("Error setting created_at: %s", err)
	}
	if err = d.Set("created_by", accountGroup.CreatedBy); err != nil {
		return fmt.Errorf("Error setting created_by: %s", err)
	}
	if accountGroup.UpdatedAt != nil {
		if err = d.Set("updated_at", accountGroup.UpdatedAt.String()); err != nil {
			return fmt.Errorf("Error setting updated_at: %s", err)
		}
	}
	if accountGroup.UpdatedBy != nil {
		if err = d.Set("updated_by", accountGroup.UpdatedBy); err != nil {
			return fmt.Errorf("Error setting updated_by: %s", err)
		}
	}

	return nil
}

func resourceIbmEnterpriseAccountGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	enterpriseManagementClient, err := meta.(ClientSession).EnterpriseManagementV1()
	if err != nil {
		return err
	}

	updateAccountGroupOptions := &enterprisemanagementv1.UpdateAccountGroupOptions{}

	updateAccountGroupOptions.SetAccountGroupID(d.Id())

	hasChange := false

	// 	if d.HasChange("parent") {
	// 		updateAccountGroupOptions.SetParent(d.Get("parent").(string))
	// 		hasChange = true
	// 	}
	if d.HasChange("name") {
		updateAccountGroupOptions.SetName(d.Get("name").(string))
		hasChange = true
	}
	if d.HasChange("primary_contact_iam_id") {
		updateAccountGroupOptions.SetPrimaryContactIamID(d.Get("primary_contact_iam_id").(string))
		hasChange = true
	}

	if hasChange {
		response, err := enterpriseManagementClient.UpdateAccountGroupWithContext(context.TODO(), updateAccountGroupOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateAccountGroupWithContext failed %s\n%s", err, response)
			return err
		}
	}

	return resourceIbmEnterpriseAccountGroupRead(d, meta)
}

func resourceIbmEnterpriseAccountGroupDelete(d *schema.ResourceData, meta interface{}) error {

	d.SetId("")

	return nil
}
