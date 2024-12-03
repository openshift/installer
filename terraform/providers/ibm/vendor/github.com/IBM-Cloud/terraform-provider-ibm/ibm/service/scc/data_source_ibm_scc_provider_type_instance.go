// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/scc-go-sdk/v5/securityandcompliancecenterapiv3"
)

func DataSourceIbmSccProviderTypeInstance() *schema.Resource {
	return AddSchemaData(&schema.Resource{
		ReadContext: dataSourceIbmSccProviderTypeInstanceRead,

		Schema: map[string]*schema.Schema{
			"provider_type_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The provider type ID.",
			},
			"provider_type_instance_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The provider type instance ID.",
			},
			"provider_type_instance_item_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the provider type instance.",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the provider type.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the provider type instance.",
			},
			"attributes": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time at which resource was created.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time at which resource was updated.",
			},
		},
	})
}

func dataSourceIbmSccProviderTypeInstanceRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	securityAndComplianceCenterApIsClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	getProviderTypeInstanceOptions := &securityandcompliancecenterapiv3.GetProviderTypeInstanceOptions{}

	getProviderTypeInstanceOptions.SetProviderTypeID(d.Get("provider_type_id").(string))
	getProviderTypeInstanceOptions.SetProviderTypeInstanceID(d.Get("provider_type_instance_id").(string))
	getProviderTypeInstanceOptions.SetInstanceID(d.Get("instance_id").(string))

	providerTypeInstanceItem, response, err := securityAndComplianceCenterApIsClient.GetProviderTypeInstanceWithContext(context, getProviderTypeInstanceOptions)
	if err != nil {
		log.Printf("[DEBUG] GetProviderTypeInstanceWithContext failed %s\n%s", err, response)
		return diag.FromErr(flex.FmtErrorf("GetProviderTypeInstanceWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *getProviderTypeInstanceOptions.ProviderTypeID, *getProviderTypeInstanceOptions.ProviderTypeInstanceID))

	if err = d.Set("provider_type_instance_item_id", providerTypeInstanceItem.ID); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting provider_type_instance_item_id: %s", err))
	}

	if err = d.Set("type", providerTypeInstanceItem.Type); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting type: %s", err))
	}

	if err = d.Set("name", providerTypeInstanceItem.Name); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting name: %s", err))
	}

	attributes := map[string]interface{}{}
	if providerTypeInstanceItem.Attributes != nil {
		attributes = providerTypeInstanceItem.Attributes
	}
	if err = d.Set("attributes", attributes); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting attributes %s", err))
	}

	if err = d.Set("created_at", flex.DateTimeToString(providerTypeInstanceItem.CreatedAt)); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting created_at: %s", err))
	}

	if err = d.Set("updated_at", flex.DateTimeToString(providerTypeInstanceItem.UpdatedAt)); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting updated_at: %s", err))
	}

	return nil
}
