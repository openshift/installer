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

func dataSourceIbmEnterprises() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIbmEnterprisesRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The name of the enterprise.",
				ValidateFunc: validateAllowedEnterpriseNameValue(),
			},
			"enterprises": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of enterprise objects.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL of the enterprise.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The enterprise ID.",
						},
						"enterprise_account_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The enterprise account ID.",
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Cloud Resource Name (CRN) of the enterprise.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the enterprise.",
						},
						"domain": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The domain of the enterprise.",
						},
						"state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The state of the enterprise.",
						},
						"primary_contact_iam_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IAM ID of the primary contact of the enterprise, such as `IBMid-0123ABC`.",
						},
						"primary_contact_email": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The email of the primary contact of the enterprise.",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time stamp at which the enterprise was created.",
						},
						"created_by": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IAM ID of the user or service that created the enterprise.",
						},
						"updated_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time stamp at which the enterprise was last updated.",
						},
						"updated_by": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IAM ID of the user or service that updated the enterprise.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmEnterprisesRead(d *schema.ResourceData, meta interface{}) error {
	enterpriseManagementClient, err := meta.(ClientSession).EnterpriseManagementV1()
	if err != nil {
		return err
	}

	listEnterprisesOptions := &enterprisemanagementv1.ListEnterprisesOptions{}

	listEnterprisesResponse, response, err := enterpriseManagementClient.ListEnterprisesWithContext(context.TODO(), listEnterprisesOptions)
	if err != nil {
		log.Printf("[DEBUG] ListEnterprisesWithContext failed %s\n%s", err, response)
		return err
	}

	// Use the provided filter argument and construct a new list with only the requested resource(s)
	var matchResources []enterprisemanagementv1.Enterprise
	var name string
	var suppliedFilter bool

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
		suppliedFilter = true
		for _, data := range listEnterprisesResponse.Resources {
			if *data.Name == name {
				matchResources = append(matchResources, data)
			}
		}
	} else {
		matchResources = listEnterprisesResponse.Resources
	}
	listEnterprisesResponse.Resources = matchResources

	if len(listEnterprisesResponse.Resources) == 0 {
		return fmt.Errorf("no Resources found with name %s\nIf not specified, please specify more filters", name)
	}

	if suppliedFilter {
		d.SetId(name)
	} else {
		d.SetId(dataSourceIbmEnterprisesID(d))
	}

	if listEnterprisesResponse.Resources != nil {
		err = d.Set("enterprises", dataSourceListEnterprisesResponseFlattenResources(listEnterprisesResponse.Resources))
		if err != nil {
			return fmt.Errorf("Error setting resources %s", err)
		}
	}

	return nil
}

// dataSourceIbmEnterprisesID returns a reasonable ID for the list.
func dataSourceIbmEnterprisesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceListEnterprisesResponseFlattenResources(result []enterprisemanagementv1.Enterprise) (resources []map[string]interface{}) {
	for _, resourcesItem := range result {
		resources = append(resources, dataSourceListEnterprisesResponseResourcesToMap(resourcesItem))
	}

	return resources
}

func dataSourceListEnterprisesResponseResourcesToMap(resourcesItem enterprisemanagementv1.Enterprise) (resourcesMap map[string]interface{}) {
	resourcesMap = map[string]interface{}{}

	if resourcesItem.URL != nil {
		resourcesMap["url"] = resourcesItem.URL
	}
	if resourcesItem.ID != nil {
		resourcesMap["id"] = resourcesItem.ID
	}
	if resourcesItem.EnterpriseAccountID != nil {
		resourcesMap["enterprise_account_id"] = resourcesItem.EnterpriseAccountID
	}
	if resourcesItem.CRN != nil {
		resourcesMap["crn"] = resourcesItem.CRN
	}
	if resourcesItem.Name != nil {
		resourcesMap["name"] = resourcesItem.Name
	}
	if resourcesItem.Domain != nil {
		resourcesMap["domain"] = resourcesItem.Domain
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
