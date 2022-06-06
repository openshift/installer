// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package enterprise

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/platform-services-go-sdk/enterprisemanagementv1"
)

func DataSourceIBMEnterpriseAccountGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmEnterpriseAccountGroupsRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The name of the account group.",
				ValidateFunc: validate.ValidateAllowedEnterpriseNameValue(),
			},
			"account_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of account groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL of the account group.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The account group ID.",
						},
						"crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Cloud Resource Name (CRN) of the account group.",
						},
						"parent": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN of the parent of the account group.",
						},
						"enterprise_account_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The enterprise account ID.",
						},
						"enterprise_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The enterprise ID that the account group is a part of.",
						},
						"enterprise_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The path from the enterprise to this particular account group.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the account group.",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The state of the account group.",
						},
						"primary_contact_iam_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IAM ID of the primary contact of the account group.",
						},
						"primary_contact_email": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The email address of the primary contact of the account group.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time stamp at which the account group was created.",
						},
						"created_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IAM ID of the user or service that created the account group.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time stamp at which the account group was last updated.",
						},
						"updated_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IAM ID of the user or service that updated the account group.",
						},
					},
				},
			},
		},
	}
}

func getEnterpriseNext(next *string) (string, error) {
	if reflect.ValueOf(next).IsNil() {
		return "", nil
	}
	u, err := url.Parse(*next)
	if err != nil {
		return "", err
	}
	q := u.Query()
	return q.Get("next_docid"), nil
}

func dataSourceIbmEnterpriseAccountGroupsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enterpriseManagementClient, err := meta.(conns.ClientSession).EnterpriseManagementV1()
	if err != nil {
		return diag.FromErr(err)
	}
	next_docid := ""
	var allRecs []enterprisemanagementv1.AccountGroup
	for {
		listAccountGroupsOptions := &enterprisemanagementv1.ListAccountGroupsOptions{}
		if next_docid != "" {
			listAccountGroupsOptions.NextDocid = &next_docid
		}
		listAccountGroupsResponse, response, err := enterpriseManagementClient.ListAccountGroupsWithContext(context, listAccountGroupsOptions)
		if err != nil {
			log.Printf("[DEBUG] ListAccountGroupsWithContext failed %s\n%s", err, response)
			return diag.FromErr(err)
		}
		next_docid, err = getEnterpriseNext(listAccountGroupsResponse.NextURL)
		if err != nil {
			log.Printf("[DEBUG] ListAccountGroupsWithContext failed. Error occurred while parsing NextURL: %s", err)
			return diag.FromErr(err)
		}
		allRecs = append(allRecs, listAccountGroupsResponse.Resources...)
		if next_docid == "" {
			break
		}
	}

	// Use the provided filter argument and construct a new list with only the requested resource(s)
	var matchResources []enterprisemanagementv1.AccountGroup
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
		return diag.FromErr(fmt.Errorf("no Resources found with name %s", name))
	}

	if suppliedFilter {
		d.SetId(name)
	} else {
		d.SetId(dataSourceIbmAccountGroupsID(d))
	}
	if allRecs != nil {
		err = d.Set("account_groups", dataSourceListEnterpriseAccountGroupsResponseFlattenResources(allRecs))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting resources %s", err))
		}
	}

	return nil
}

// dataSourceIbmAccountGroupsID returns a reasonable ID for the list.
func dataSourceIbmAccountGroupsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceListEnterpriseAccountGroupsResponseFlattenResources(result []enterprisemanagementv1.AccountGroup) (resources []map[string]interface{}) {
	for _, resourcesItem := range result {
		resources = append(resources, dataSourceListEnterpriseAccountGroupsResponseResourcesToMap(resourcesItem))
	}

	return resources
}

func dataSourceListEnterpriseAccountGroupsResponseResourcesToMap(resourcesItem enterprisemanagementv1.AccountGroup) (resourcesMap map[string]interface{}) {
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
	if resourcesItem.PrimaryContactIamID != nil {
		resourcesMap["primary_contact_iam_id"] = resourcesItem.PrimaryContactIamID
	}
	if resourcesItem.PrimaryContactEmail != nil {
		resourcesMap["primary_contact_email"] = resourcesItem.PrimaryContactEmail
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
