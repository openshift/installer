// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	rg "github.com/IBM/platform-services-go-sdk/resourcemanagerv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIBMResourceGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMResourceGroupRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description:  "Resource group name",
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"is_default", "name"},
			},
			"is_default": {
				Description:  "Default Resource group",
				Type:         schema.TypeBool,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"is_default", "name"},
			},
			"state": {
				Type:        schema.TypeString,
				Description: "State of the resource group",
				Computed:    true,
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
			"account_id": {
				Type:        schema.TypeString,
				Description: "Account ID",
				Computed:    true,
			},
		},
	}
}

func dataSourceIBMResourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	rMgtClient, err := meta.(ClientSession).ResourceManagerV2API()
	if err != nil {
		return err
	}

	var defaultGrp bool
	if group, ok := d.GetOk("is_default"); ok {
		defaultGrp = group.(bool)
	}
	var name string
	if n, ok := d.GetOk("name"); ok {
		name = n.(string)
	}

	if !defaultGrp && name == "" {
		return fmt.Errorf("[ERROR] Missing required properties. Need a resource group name, or the is_default true")
	}
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	accountID := userDetails.userAccount

	resourceGroupList := rg.ListResourceGroupsOptions{
		AccountID: &accountID,
	}
	if defaultGrp {
		resourceGroupList.Default = &defaultGrp

	} else if name != "" {
		resourceGroupList.Name = &name
	}
	rg, resp, err := rMgtClient.ListResourceGroups(&resourceGroupList)
	if err != nil || rg == nil || rg.Resources == nil {
		return fmt.Errorf("[ERROR] Error retrieving resource group: %s %s", err, resp)
	}
	if len(rg.Resources) < 1 {
		return fmt.Errorf("[ERROR] Given Resource Group is not found in the account : %s %s", err, resp)
	}
	resourceGroup := rg.Resources[0]
	d.SetId(*resourceGroup.ID)
	if resourceGroup.Name != nil {
		d.Set("name", *resourceGroup.Name)
	}
	if resourceGroup.Default != nil {
		d.Set("is_default", *resourceGroup.Default)
	}
	if resourceGroup.State != nil {
		d.Set("state", *resourceGroup.State)
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
	if resourceGroup.QuotaID != nil {
		d.Set("account_id", *resourceGroup.AccountID)
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
