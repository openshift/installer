// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package secretsmanager

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/secrets-manager-go-sdk/v2/secretsmanagerv2"
)

func DataSourceIbmSmSecretGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSmSecretGroupRead,

		Schema: map[string]*schema.Schema{
			"secret_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the secret group.",
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
	}
}

func dataSourceIbmSmSecretGroupRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		return diag.FromErr(err)
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	getSecretGroupOptions := &secretsmanagerv2.GetSecretGroupOptions{}

	secretGroupId := d.Get("secret_group_id").(string)
	getSecretGroupOptions.SetID(secretGroupId)

	secretGroup, response, err := secretsManagerClient.GetSecretGroupWithContext(context, getSecretGroupOptions)
	if err != nil {
		log.Printf("[DEBUG] GetSecretGroupWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetSecretGroupWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", region, instanceId, secretGroupId))

	if err = d.Set("region", region); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting region: %s", err))
	}
	if err = d.Set("name", secretGroup.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}

	if err = d.Set("description", secretGroup.Description); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
	}

	if err = d.Set("created_at", DateTimeToRFC3339(secretGroup.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}

	if err = d.Set("updated_at", DateTimeToRFC3339(secretGroup.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
	}

	return nil
}
