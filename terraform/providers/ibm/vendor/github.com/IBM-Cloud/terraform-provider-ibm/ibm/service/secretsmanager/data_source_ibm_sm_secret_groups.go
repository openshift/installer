// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package secretsmanager

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/secrets-manager-go-sdk/v2/secretsmanagerv2"
)

func DataSourceIbmSmSecretGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSmSecretGroupsRead,

		Schema: map[string]*schema.Schema{
			"secret_groups": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A collection of secret groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A v4 UUID identifier.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of your secret group.",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An extended description of your secret group.To protect your privacy, do not use personal data, such as your name or location, as a description for your secret group.",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date that a resource was created. The date format follows RFC 3339.",
						},
						"updated_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date that a resource was recently modified. The date format follows RFC 3339.",
						},
					},
				},
			},
			"total_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of resources in a collection.",
			},
		},
	}
}

func dataSourceIbmSmSecretGroupsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", fmt.Sprintf("(Data) %s", SecretGroupsResourceName), "read")
		return tfErr.GetDiag()
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	listSecretGroupsOptions := &secretsmanagerv2.ListSecretGroupsOptions{}

	secretGroupCollection, response, err := secretsManagerClient.ListSecretGroupsWithContext(context, listSecretGroupsOptions)
	if err != nil {
		log.Printf("[DEBUG] ListSecretGroupsWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListSecretGroupsWithContext failed %s\n%s", err, response), fmt.Sprintf("(Data) %s", SecretGroupsResourceName), "read")
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", region, instanceId))

	secretGroups := []map[string]interface{}{}
	if secretGroupCollection.SecretGroups != nil {
		for _, modelItem := range secretGroupCollection.SecretGroups {
			modelMap, err := dataSourceIbmSmSecretGroupsSecretGroupToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, "", fmt.Sprintf("(Data) %s", SecretGroupsResourceName), "read")
				return tfErr.GetDiag()
			}
			secretGroups = append(secretGroups, modelMap)
		}
	}
	if err = d.Set("secret_groups", secretGroups); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_groups"), fmt.Sprintf("(Data) %s", SecretGroupsResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("total_count", flex.IntValue(secretGroupCollection.TotalCount)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting total_count"), fmt.Sprintf("(Data) %s", SecretGroupsResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("region", region); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting region"), fmt.Sprintf("(Data) %s", SecretGroupsResourceName), "read")
		return tfErr.GetDiag()
	}
	return nil
}

// dataSourceIbmSmSecretGroupsID returns a reasonable ID for the list.
func dataSourceIbmSmSecretGroupsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIbmSmSecretGroupsSecretGroupToMap(model *secretsmanagerv2.SecretGroup) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	if model.CreatedAt != nil {
		modelMap["created_at"] = model.CreatedAt.String()
	}
	if model.UpdatedAt != nil {
		modelMap["updated_at"] = model.UpdatedAt.String()
	}
	return modelMap, nil
}
