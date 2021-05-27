// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/IBM/container-registry-go-sdk/containerregistryv1"
)

func dataIBMContainerRegistryNamespaces() *schema.Resource {
	return &schema.Resource{
		Read: dataIBMContainerRegistryNamespacesRead,

		Schema: map[string]*schema.Schema{
			"namespaces": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Container Registry Namespaces",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Container Registry Namespace name",
						},
						"resource_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource Group to which namespace has to be assigned",
						},
						"crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CRN of the Namespace",
						},
						"created_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Created Date",
						},
						"updated_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Updated Date",
						},
						"resource_created_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "When the namespace was assigned to a resource group.",
						},
						"account": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IBM Cloud account that owns the namespace.",
						},
						// DEPRECATED FIELDS TO BE REMOVED IN FUTURE
						"created_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Created Date",
							Deprecated:  "This field is deprecated",
						},
						"updated_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Updated Date",
							Deprecated:  "This field is deprecated",
						},
					},
				},
			},
		},
	}
}

func dataIBMContainerRegistryNamespacesRead(d *schema.ResourceData, meta interface{}) error {
	containerRegistryClient, err := meta.(ClientSession).ContainerRegistryV1()
	if err != nil {
		return err
	}

	listNamespaceDetailsOptions := &containerregistryv1.ListNamespaceDetailsOptions{}

	namespaceDetailsList, _, err := containerRegistryClient.ListNamespaceDetails(listNamespaceDetailsOptions)
	if err != nil {
		return err
	}

	namespaces := []map[string]interface{}{}
	for _, namespaceDetails := range namespaceDetailsList {
		namespace := map[string]interface{}{}
		namespace["name"] = namespaceDetails.Name
		namespace["resource_group_id"] = namespaceDetails.ResourceGroup
		namespace["crn"] = namespaceDetails.CRN
		namespace["created_date"] = namespaceDetails.CreatedDate
		namespace["updated_date"] = namespaceDetails.UpdatedDate
		namespace["account"] = namespaceDetails.Account
		namespace["resource_created_date"] = namespaceDetails.ResourceCreatedDate
		// DEPRECATED FIELDS TO BE REMOVED IN FUTURE
		namespace["created_on"] = namespaceDetails.CreatedDate
		namespace["updated_on"] = namespaceDetails.UpdatedDate
		namespaces = append(namespaces, namespace)
	}
	if err = d.Set("namespaces", namespaces); err != nil {
		return fmt.Errorf("Error setting namespaces: %s", err)
	}
	d.SetId(time.Now().UTC().String())
	return nil
}
