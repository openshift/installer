// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package hpcs

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/ibm-hpcs-uko-sdk/ukov4"
)

func DataSourceIbmVault() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIbmVaultRead,

		Schema: map[string]*schema.Schema{
			"instance_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the UKO instance this resource exists in.",
			},
			"region": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The region of the UKO instance this resource exists in.",
			},
			"vault_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "UUID of the vault.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the vault.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the vault.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date and time when the vault was created.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date and time when the vault was last updated.",
			},
			"created_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the user that created the vault.",
			},
			"updated_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the user that last updated the vault.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A URL that uniquely identifies your cloud resource.",
			},
		},
	}
}

func DataSourceIbmVaultRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ukoClient, err := meta.(conns.ClientSession).UkoV4()
	if err != nil {
		return diag.FromErr(err)
	}

	getVaultOptions := &ukov4.GetVaultOptions{}

	region := d.Get("region").(string)
	instance_id := d.Get("instance_id").(string)
	vault_id := d.Get("vault_id").(string)

	getVaultOptions.SetID(vault_id)

	url, err := getUkoUrl(context, region, instance_id, ukoClient)
	if err != nil {
		return diag.FromErr(err)
	}
	ukoClient.SetServiceURL(url)

	vault, response, err := ukoClient.GetVaultWithContext(context, getVaultOptions)
	if err != nil {
		log.Printf("[DEBUG] GetVaultWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetVaultWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s/%s/%s", region, instance_id, vault_id, *getVaultOptions.ID))

	if err = d.Set("name", vault.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}

	if err = d.Set("description", vault.Description); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
	}

	if err = d.Set("created_at", flex.DateTimeToString(vault.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}

	if err = d.Set("updated_at", flex.DateTimeToString(vault.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
	}

	if err = d.Set("created_by", vault.CreatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_by: %s", err))
	}

	if err = d.Set("updated_by", vault.UpdatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_by: %s", err))
	}

	if err = d.Set("href", vault.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}

	return nil
}
