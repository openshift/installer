// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package eventnotification

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	en "github.com/IBM/event-notifications-go-admin-sdk/eventnotificationsv1"
)

func ResourceIBMEnIBMSource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMEnIBMSourceCreate,
		ReadContext:   resourceIBMEnIBMSourceRead,
		UpdateContext: resourceIBMEnIBMSourceUpdate,
		DeleteContext: resourceIBMEnIBMSourceDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"instance_guid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Unique identifier for IBM Cloud Event Notifications instance.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "The enabled flag for source",
			},
			"source_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Destination ID",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last updated time.",
			},
		},
	}
}

func resourceIBMEnIBMSourceCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_ibmsource", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
		// return diag.FromErr(err)
	}

	options := &en.UpdateSourceOptions{}

	options.SetInstanceID(d.Get("instance_guid").(string))
	options.SetID(d.Get("source_id").(string))

	options.SetEnabled(d.Get("enabled").(bool))

	result, _, err := enClient.UpdateSourceWithContext(context, options)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpadteIBMSourceWithContext failed: %s", err.Error()), "ibm_en_ibmsource", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
		// return diag.FromErr(fmt.Errorf("CreateSourceWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *options.InstanceID, *result.ID))

	return resourceIBMEnIBMSourceRead(context, d, meta)
}

func resourceIBMEnIBMSourceRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_ibmsource", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
		// return diag.FromErr(err)
	}

	options := &en.GetSourceOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_ibmsource", "read")
		return tfErr.GetDiag()
		// return diag.FromErr(err)
	}

	options.SetInstanceID(parts[0])
	options.SetID(parts[1])

	options.SetInstanceID(d.Get("instance_guid").(string))
	options.SetID(d.Get("source_id").(string))

	result, response, err := enClient.GetSourceWithContext(context, options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId(d.Get("source_id").(string))
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetSourceWithContext failed: %s", err.Error()), "ibm_en_ibmsource", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
		// return diag.FromErr(fmt.Errorf("GetSourceWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("instance_guid", options.InstanceID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting instance_guid: %s", err))
	}

	if err = d.Set("source_id", options.ID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting source_id: %s", err))

	}

	// if err = d.Set("name", result.Name); err != nil {
	// 	return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	// }

	// if err = d.Set("description", result.Description); err != nil {
	// 	return diag.FromErr(fmt.Errorf("[ERROR] Error setting description: %s", err))
	// }

	if err = d.Set("enabled", result.Enabled); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting enabled: %s", err))
	}

	if err = d.Set("updated_at", flex.DateTimeToString(result.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting updated_at: %s", err))
	}

	return nil
}

func resourceIBMEnIBMSourceUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_ibmsource", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
		// return diag.FromErr(err)
	}

	options := &en.UpdateSourceOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_ibmsource", "update")
		return tfErr.GetDiag()
		// return diag.FromErr(err)
	}

	options.SetInstanceID(parts[0])
	options.SetID(parts[1])

	options.SetInstanceID(d.Get("instance_guid").(string))
	options.SetID(d.Get("source_id").(string))

	if ok := d.HasChanges("enabled"); ok {

		options.SetEnabled(d.Get("enabled").(bool))

		_, _, err := enClient.UpdateSourceWithContext(context, options)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateSourceWithContext failed: %s", err.Error()), "ibm_en_ibmsource", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
			// return diag.FromErr(fmt.Errorf("UpdateSourceWithContext failed %s\n%s", err, response))
		}

		return resourceIBMEnIBMSourceRead(context, d, meta)
	}

	return nil
}

func resourceIBMEnIBMSourceDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return nil
}
