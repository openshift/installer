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

func ResourceIBMEnSource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMEnSourceCreate,
		ReadContext:   resourceIBMEnSourceRead,
		UpdateContext: resourceIBMEnSourceUpdate,
		DeleteContext: resourceIBMEnSourceDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"instance_guid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Unique identifier for IBM Cloud Event Notifications instance.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Source name.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Source description.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "The enabled flag for source",
			},
			"source_id": {
				Type:        schema.TypeString,
				Computed:    true,
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

func resourceIBMEnSourceCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_source", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	options := &en.CreateSourcesOptions{}

	options.SetInstanceID(d.Get("instance_guid").(string))
	options.SetName(d.Get("name").(string))

	if _, ok := d.GetOk("description"); ok {
		options.SetDescription(d.Get("description").(string))
	}

	options.SetEnabled(d.Get("enabled").(bool))

	result, _, err := enClient.CreateSourcesWithContext(context, options)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateSourcesWithContext failed: %s", err.Error()), "ibm_en_source", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *options.InstanceID, *result.ID))

	return resourceIBMEnSourceRead(context, d, meta)
}

func resourceIBMEnSourceRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_source", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	options := &en.GetSourceOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_source", "read")
		return tfErr.GetDiag()
	}

	options.SetInstanceID(parts[0])
	options.SetID(parts[1])

	result, response, err := enClient.GetSourceWithContext(context, options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetSourceWithContext failed: %s", err.Error()), "ibm_en_source", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("instance_guid", options.InstanceID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting instance_guid: %s", err))
	}

	if err = d.Set("source_id", options.ID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting source_id: %s", err))
	}

	if err = d.Set("name", result.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}

	if err = d.Set("description", result.Description); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting description: %s", err))
	}

	if err = d.Set("enabled", result.Enabled); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting enabled: %s", err))
	}

	if err = d.Set("updated_at", flex.DateTimeToString(result.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting updated_at: %s", err))
	}

	return nil
}

func resourceIBMEnSourceUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_source", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	options := &en.UpdateSourceOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_source", "update")
		return tfErr.GetDiag()
	}

	options.SetInstanceID(parts[0])
	options.SetID(parts[1])

	if ok := d.HasChanges("name", "description", "enabled"); ok {
		options.SetName(d.Get("name").(string))

		if _, ok := d.GetOk("description"); ok {
			options.SetDescription(d.Get("description").(string))
		}

		options.SetEnabled(d.Get("enabled").(bool))

		_, _, err := enClient.UpdateSourceWithContext(context, options)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateSourceWithContext failed: %s", err.Error()), "ibm_en_source", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}

		return resourceIBMEnSourceRead(context, d, meta)
	}

	return nil
}

func resourceIBMEnSourceDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_source", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	options := &en.DeleteSourceOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_source", "delete")
		return tfErr.GetDiag()
	}

	options.SetInstanceID(parts[0])
	options.SetID(parts[1])

	response, err := enClient.DeleteSourceWithContext(context, options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteSourceWithContext: failed: %s", err.Error()), "ibm_en_source", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}
