// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package schematics

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/schematics-go-sdk/schematicsv1"
)

func DataSourceIBMSchematicsInventory() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMSchematicsInventoryRead,

		Schema: map[string]*schema.Schema{
			"inventory_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Resource Inventory Id.  Use `GET /v2/inventories` API to look up the Resource Inventory definition Ids  in your IBM Cloud account.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique name of your Inventory.  The name can be up to 128 characters long and can include alphanumeric  characters, spaces, dashes, and underscores.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Inventory id.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of your Inventory.  The description can be up to 2048 characters long in size.",
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources, provisioned using Schematics.",
			},
			"resource_group": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource-group name for the Inventory definition.  By default, Inventory will be created in Default Resource Group.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Inventory creation time.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Email address of user who created the Inventory.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Inventory updation time.",
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Email address of user who updated the Inventory.",
			},
			"inventories_ini": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Input inventory of host and host group for the playbook,  in the .ini file format.",
			},
			"resource_queries": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Input resource queries that is used to dynamically generate  the inventory of host and host group for the playbook.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceIBMSchematicsInventoryRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsInventoryRead schematicsClient initialization failed: %s", err.Error()), "ibm_schematics_inventory", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if r, ok := d.GetOk("location"); ok {
		region := r.(string)
		schematicsURL, updatedURL, _ := SchematicsEndpointURL(region, meta)
		if updatedURL {
			schematicsClient.Service.Options.URL = schematicsURL
		}
	}

	getInventoryOptions := &schematicsv1.GetInventoryOptions{}

	getInventoryOptions.SetInventoryID(d.Get("inventory_id").(string))

	inventoryResourceRecord, response, err := schematicsClient.GetInventoryWithContext(context, getInventoryOptions)
	if err != nil {

		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsInventoryRead GetInventoryWithContext failed with error: %s and response:\n%s", err, response), "ibm_schematics_inventory", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*getInventoryOptions.InventoryID)
	if err = d.Set("name", inventoryResourceRecord.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsInventoryRead failed with error: %s", err), "ibm_schematics_inventory", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("id", inventoryResourceRecord.ID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsInventoryRead failed with error: %s", err), "ibm_schematics_inventory", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("description", inventoryResourceRecord.Description); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsInventoryRead failed with error: %s", err), "ibm_schematics_inventory", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("location", inventoryResourceRecord.Location); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsInventoryRead failed with error: %s", err), "ibm_schematics_inventory", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("resource_group", inventoryResourceRecord.ResourceGroup); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsInventoryRead failed with error: %s", err), "ibm_schematics_inventory", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("created_at", flex.DateTimeToString(inventoryResourceRecord.CreatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsInventoryRead failed with error: %s", err), "ibm_schematics_inventory", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("created_by", inventoryResourceRecord.CreatedBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsInventoryRead failed with error: %s", err), "ibm_schematics_inventory", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("updated_at", flex.DateTimeToString(inventoryResourceRecord.UpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsInventoryRead failed with error: %s", err), "ibm_schematics_inventory", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("updated_by", inventoryResourceRecord.UpdatedBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsInventoryRead failed with error: %s", err), "ibm_schematics_inventory", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("inventories_ini", inventoryResourceRecord.InventoriesIni); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsInventoryRead failed with error: %s", err), "ibm_schematics_inventory", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	return nil
}
