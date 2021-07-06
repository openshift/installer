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

func dataSourceIbmEnterpriseAccounts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIbmEnterpriseAccountsRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The name of the account.",
				ValidateFunc: validateAllowedEnterpriseNameValue(),
			},
			"accounts": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of accounts.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL of the account.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The account ID.",
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Cloud Resource Name (CRN) of the account.",
						},
						"parent": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN of the parent of the account.",
						},
						"enterprise_account_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The enterprise account ID.",
						},
						"enterprise_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The enterprise ID that the account is a part of.",
						},
						"enterprise_path": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The path from the enterprise to this particular account.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the account.",
						},
						"state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The state of the account.",
						},
						"owner_iam_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IAM ID of the owner of the account.",
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
				},
			},
		},
	}
}

func dataSourceIbmEnterpriseAccountsRead(d *schema.ResourceData, meta interface{}) error {
	enterpriseManagementClient, err := meta.(ClientSession).EnterpriseManagementV1()
	if err != nil {
		return err
	}
	next_docid := ""
	var allRecs []enterprisemanagementv1.Account
	for {
		listAccountsOptions := &enterprisemanagementv1.ListAccountsOptions{}
		if next_docid != "" {
			listAccountsOptions.NextDocid = &next_docid
		}
		listAccountsResponse, response, err := enterpriseManagementClient.ListAccountsWithContext(context.TODO(), listAccountsOptions)
		if err != nil {
			log.Printf("[DEBUG] ListAccountsWithContext failed %s\n%s", err, response)
			return err
		}
		next_docid, err = getEnterpriseNext(listAccountsResponse.NextURL)
		if err != nil {
			log.Printf("[DEBUG] ListAccountsWithContext failed. Error occurred while parsing NextURL: %s", err)
			return err
		}
		allRecs = append(allRecs, listAccountsResponse.Resources...)
		if next_docid == "" {
			break
		}
	}

	// Use the provided filter argument and construct a new list with only the requested resource(s)
	var matchResources []enterprisemanagementv1.Account
	var name string
	var suppliedFilter bool

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
		suppliedFilter = true
		for _, data := range allRecs {
			if *data.Name == name {
				matchResources = append(matchResources, data)
			}
		}
	} else {
		matchResources = allRecs
	}
	allRecs = matchResources

	if len(allRecs) == 0 {
		return fmt.Errorf("no Resources found with name %s\nIf not specified, please specify more filters", name)
	}

	if suppliedFilter {
		d.SetId(name)
	} else {
		d.SetId(dataSourceIbmEnterpriseAccountsID(d))
	}

	if allRecs != nil {
		err = d.Set("accounts", dataSourceListEnterpriseAccountsResponseFlattenResources(allRecs))
		if err != nil {
			return fmt.Errorf("Error setting resources %s", err)
		}
	}

	return nil
}

// dataSourceIbmAccountsID returns a reasonable ID for the list.
func dataSourceIbmEnterpriseAccountsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceListEnterpriseAccountsResponseFlattenResources(result []enterprisemanagementv1.Account) (resources []map[string]interface{}) {
	for _, resourcesItem := range result {
		resources = append(resources, dataSourceListEnterpriseAccountsResponseResourcesToMap(resourcesItem))
	}

	return resources
}

func dataSourceListEnterpriseAccountsResponseResourcesToMap(resourcesItem enterprisemanagementv1.Account) (resourcesMap map[string]interface{}) {
	resourcesMap = map[string]interface{}{}

	if resourcesItem.URL != nil {
		resourcesMap["url"] = resourcesItem.URL
	}
	if resourcesItem.ID != nil {
		resourcesMap["id"] = resourcesItem.ID
	}
	if resourcesItem.CRN != nil {
		resourcesMap["crn"] = resourcesItem.CRN
	}
	if resourcesItem.Parent != nil {
		resourcesMap["parent"] = resourcesItem.Parent
	}
	if resourcesItem.EnterpriseAccountID != nil {
		resourcesMap["enterprise_account_id"] = resourcesItem.EnterpriseAccountID
	}
	if resourcesItem.EnterpriseID != nil {
		resourcesMap["enterprise_id"] = resourcesItem.EnterpriseID
	}
	if resourcesItem.EnterprisePath != nil {
		resourcesMap["enterprise_path"] = resourcesItem.EnterprisePath
	}
	if resourcesItem.Name != nil {
		resourcesMap["name"] = resourcesItem.Name
	}
	if resourcesItem.State != nil {
		resourcesMap["state"] = resourcesItem.State
	}
	if resourcesItem.OwnerIamID != nil {
		resourcesMap["owner_iam_id"] = resourcesItem.OwnerIamID
	}
	if resourcesItem.Paid != nil {
		resourcesMap["paid"] = resourcesItem.Paid
	}
	if resourcesItem.OwnerEmail != nil {
		resourcesMap["owner_email"] = resourcesItem.OwnerEmail
	}
	if resourcesItem.IsEnterpriseAccount != nil {
		resourcesMap["is_enterprise_account"] = resourcesItem.IsEnterpriseAccount
	}
	if resourcesItem.CreatedAt != nil {
		resourcesMap["created_at"] = resourcesItem.CreatedAt.String()
	}
	if resourcesItem.CreatedBy != nil {
		resourcesMap["created_by"] = resourcesItem.CreatedBy
	}
	if resourcesItem.UpdatedAt != nil {
		resourcesMap["updated_at"] = resourcesItem.UpdatedAt.String()
	}
	if resourcesItem.UpdatedBy != nil {
		resourcesMap["updated_by"] = resourcesItem.UpdatedBy
	}

	return resourcesMap
}
