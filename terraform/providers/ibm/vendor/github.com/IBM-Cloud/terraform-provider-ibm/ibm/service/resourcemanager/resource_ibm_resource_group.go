// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package resourcemanager

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	rg "github.com/IBM/platform-services-go-sdk/resourcemanagerv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMResourceGroup() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMResourceGroupCreate,
		Read:     resourceIBMResourceGroupRead,
		Update:   resourceIBMResourceGroupUpdate,
		Delete:   resourceIBMResourceGroupDelete,
		Exists:   resourceIBMResourceGroupExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the resource group",
			},
			"default": {
				Description: "Specifies whether its default resource group or not",
				Type:        schema.TypeBool,
				Computed:    true,
			},

			"state": {
				Type:        schema.TypeString,
				Description: "State of the resource group",
				Computed:    true,
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"crn": {
				Type:        schema.TypeString,
				Description: "The full CRN associated with the resource group",
				Computed:    true,
			},
			"created_at": {
				Type:        schema.TypeString,
				Description: "The date when the resource group was initially created.",
				Computed:    true,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Description: "The date when the resource group was last updated.",
				Computed:    true,
			},
			"teams_url": {
				Type:        schema.TypeString,
				Description: "The URL to access the team details that associated with the resource group.",
				Computed:    true,
			},
			"payment_methods_url": {
				Type:        schema.TypeString,
				Description: "The URL to access the payment methods details that associated with the resource group.",
				Computed:    true,
			},
			"quota_url": {
				Type:        schema.TypeString,
				Description: "The URL to access the quota details that associated with the resource group.",
				Computed:    true,
			},
			"quota_id": {
				Type:        schema.TypeString,
				Description: "An alpha-numeric value identifying the quota ID associated with the resource group.",
				Computed:    true,
			},
			"resource_linkages": {
				Type:        schema.TypeSet,
				Description: "An array of the resources that linked to the resource group",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
		},
	}
}

func resourceIBMResourceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	rMgtClient, err := meta.(conns.ClientSession).ResourceManagerV2API()
	if err != nil {
		return err
	}
	name := d.Get("name").(string)

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	accountID := userDetails.UserAccount

	resourceGroupCreate := rg.CreateResourceGroupOptions{
		Name:      &name,
		AccountID: &accountID,
	}

	resourceGroup, resp, err := rMgtClient.CreateResourceGroup(&resourceGroupCreate)
	if err != nil {
		return fmt.Errorf("[ERROR] Error creating resource group: %s with response code  %s", err, resp)
	}

	d.SetId(*resourceGroup.ID)

	return resourceIBMResourceGroupRead(d, meta)
}

func resourceIBMResourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	rMgtClient, err := meta.(conns.ClientSession).ResourceManagerV2API()
	if err != nil {
		return err
	}
	resourceGroupID := d.Id()
	resourceGroupGet := rg.GetResourceGroupOptions{
		ID: &resourceGroupID,
	}

	resourceGroup, resp, err := rMgtClient.GetResourceGroup(&resourceGroupGet)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] Resource Group is not found")
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error retrieving resource group: %s with response code  %s", err, resp)
	}

	d.Set("name", *resourceGroup.Name)
	if resourceGroup.State != nil {
		d.Set("state", *resourceGroup.State)
	}
	if resourceGroup.Default != nil {
		d.Set("default", *resourceGroup.Default)
	}
	if resourceGroup.CRN != nil {
		d.Set("crn", *resourceGroup.CRN)
	}
	if resourceGroup.CreatedAt != nil {
		createdAt := *resourceGroup.CreatedAt
		d.Set("created_at", createdAt.String())
	}
	if resourceGroup.UpdatedAt != nil {
		UpdatedAt := *resourceGroup.UpdatedAt
		d.Set("updated_at", UpdatedAt.String())
	}
	if resourceGroup.TeamsURL != nil {
		d.Set("teams_url", *resourceGroup.TeamsURL)
	}
	if resourceGroup.PaymentMethodsURL != nil {
		d.Set("payment_methods_url", *resourceGroup.PaymentMethodsURL)
	}
	if resourceGroup.QuotaURL != nil {
		d.Set("quota_url", *resourceGroup.QuotaURL)
	}
	if resourceGroup.QuotaID != nil {
		d.Set("quota_id", *resourceGroup.QuotaID)
	}
	if resourceGroup.ResourceLinkages != nil {
		rl := make([]string, 0)
		for _, r := range resourceGroup.ResourceLinkages {
			rl = append(rl, r.(string))
		}
		d.Set("resource_linkages", rl)
	}
	return nil
}

func resourceIBMResourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	rMgtClient, err := meta.(conns.ClientSession).ResourceManagerV2API()
	if err != nil {
		return err
	}

	resourceGroupID := d.Id()
	resourceGroupUpdate := rg.UpdateResourceGroupOptions{
		ID: &resourceGroupID,
	}
	hasChange := false
	if d.HasChange("name") {
		name := d.Get("name").(string)
		resourceGroupUpdate.Name = &name
		hasChange = true
	}

	if hasChange {
		_, resp, err := rMgtClient.UpdateResourceGroup(&resourceGroupUpdate)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating resource group: %s with response code  %s", err, resp)
		}

	}
	return resourceIBMResourceGroupRead(d, meta)
}

func resourceIBMResourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	rMgtClient, err := meta.(conns.ClientSession).ResourceManagerV2API()
	if err != nil {
		return err
	}

	resourceGroupID := d.Id()
	resourceGroupDelete := rg.DeleteResourceGroupOptions{
		ID: &resourceGroupID,
	}

	resp, err := rMgtClient.DeleteResourceGroup(&resourceGroupDelete)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] Resource Group is not found")
			return nil
		}
		if resp != nil && resp.StatusCode == 500 {
			err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
				resp, err = rMgtClient.DeleteResourceGroup(&resourceGroupDelete)
				if err != nil {
					if resp != nil && resp.StatusCode == 500 {
						return resource.RetryableError(err)
					}
					if resp != nil && resp.StatusCode == 404 {
						log.Printf("[WARN] Resource Group is not found")
						return nil
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			if conns.IsResourceTimeoutError(err) {
				resp, err = rMgtClient.DeleteResourceGroup(&resourceGroupDelete)
			}
			if err != nil {
				if resp != nil && resp.StatusCode == 404 {
					log.Printf("[WARN] Resource Group is not found")
					return nil
				}
				return fmt.Errorf("[ERROR] Error Deleting resource group: %s with response code  %s", err, resp)
			}
		} else {

			return fmt.Errorf("[ERROR] Error Deleting resource group: %s with response code  %s", err, resp)
		}
	}

	d.SetId("")

	return nil
}

func resourceIBMResourceGroupExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	rMgtClient, err := meta.(conns.ClientSession).ResourceManagerV2API()
	if err != nil {
		return false, err
	}
	resourceGroupID := d.Id()
	resourceGroupGet := rg.GetResourceGroupOptions{
		ID: &resourceGroupID,
	}

	resourceGroup, resp, err := rMgtClient.GetResourceGroup(&resourceGroupGet)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error getting resource group: %s with response code  %s", err, resp)
	}

	return *resourceGroup.ID == resourceGroupID, nil
}
