// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"strings"

	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/helpers/order"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/sl"
)

func resourceIBMObjectStorageAccount() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMObjectStorageAccountCreate,
		Read:     resourceIBMObjectStorageAccountRead,
		Update:   resourceIBMObjectStorageAccountUpdate,
		Delete:   resourceIBMObjectStorageAccountDelete,
		Exists:   resourceIBMObjectStorageAccountExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"local_note": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func resourceIBMObjectStorageAccountCreate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	accountService := services.GetAccountService(sess)

	// Check if an object storage account exists
	objectStorageAccounts, err := accountService.GetHubNetworkStorage()
	if err != nil {
		return fmt.Errorf("resource_ibm_object_storage_account: Error on create: %s", err)
	}

	if len(objectStorageAccounts) == 0 {
		// Order the account
		productOrderService := services.GetProductOrderService(sess.SetRetries(0))

		receipt, err := productOrderService.PlaceOrder(&datatypes.Container_Product_Order{
			Quantity:  sl.Int(1),
			PackageId: sl.Int(0),
			Prices: []datatypes.Product_Item_Price{
				{Id: sl.Int(30920)},
			},
		}, sl.Bool(false))
		if err != nil {
			return fmt.Errorf(
				"resource_ibm_object_storage_account: Error ordering account: %s", err)
		}

		// Wait for the object storage account order to complete.
		billingOrderItem, err := WaitForOrderCompletion(&receipt, meta)
		if err != nil {
			return fmt.Errorf(
				"Error waiting for object storage account order (%d) to complete: %s", receipt.OrderId, err)
		}

		// Get accountName using filter on hub network storage
		objectStorageAccounts, err = accountService.Filter(
			filter.Path("billingItem.id").Eq(billingOrderItem.BillingItem.Id).Build(),
		).GetNetworkStorage()
		if err != nil {
			return fmt.Errorf("resource_ibm_object_storage_account: Error on retrieving new: %s", err)
		}

		if len(objectStorageAccounts) == 0 {
			return fmt.Errorf("resource_ibm_object_storage_account: Failed to create object storage account.")
		}
	}

	// Get account name and set as the Id
	d.SetId(*objectStorageAccounts[0].Username)
	d.Set("name", *objectStorageAccounts[0].Username)

	return nil
}

func WaitForOrderCompletion(
	receipt *datatypes.Container_Product_Order_Receipt, meta interface{}) (datatypes.Billing_Order_Item, error) {

	log.Printf("Waiting for billing order %d to have zero active transactions", receipt.OrderId)
	var billingOrderItem *datatypes.Billing_Order_Item

	stateConf := &resource.StateChangeConf{
		Pending: []string{"", "in progress"},
		Target:  []string{"complete"},
		Refresh: func() (interface{}, string, error) {
			var err error
			var completed bool

			completed, billingOrderItem, err = order.CheckBillingOrderComplete(meta.(ClientSession).SoftLayerSession(), receipt)
			if err != nil {
				return nil, "", err
			}

			if completed {
				return billingOrderItem, "complete", nil
			} else {
				return billingOrderItem, "in progress", nil
			}
		},
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return *billingOrderItem, err
}

func resourceIBMObjectStorageAccountRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	accountService := services.GetAccountService(sess)
	accountName := d.Id()
	d.Set("name", accountName)

	// Check if an object storage account exists
	objectStorageAccounts, err := accountService.Filter(
		filter.Path("username").Eq(accountName).Build(),
	).GetHubNetworkStorage()
	if err != nil {
		return fmt.Errorf("resource_ibm_object_storage_account: Error on Read: %s", err)
	}

	for _, objectStorageAccount := range objectStorageAccounts {
		if *objectStorageAccount.Username == accountName {
			return nil
		}
	}

	return fmt.Errorf("resource_ibm_object_storage_account: Could not find account %s", accountName)
}

func resourceIBMObjectStorageAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	// Nothing to update for now. Not supported.
	return nil
}

func resourceIBMObjectStorageAccountDelete(d *schema.ResourceData, meta interface{}) error {
	// Delete is not supported for now.
	return nil
}

func resourceIBMObjectStorageAccountExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	err := resourceIBMObjectStorageAccountRead(d, meta)
	if err != nil {
		if strings.Contains(err.Error(), "Could not find account") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
