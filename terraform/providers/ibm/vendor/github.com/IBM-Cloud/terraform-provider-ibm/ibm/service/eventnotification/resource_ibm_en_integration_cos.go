// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package eventnotification

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	en "github.com/IBM/event-notifications-go-admin-sdk/eventnotificationsv1"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMEnCOSIntegration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMEnCOSIntegrationCreate,
		ReadContext:   resourceIBMEnCOSIntegrationRead,
		UpdateContext: resourceIBMEnCOSIntegrationUpdate,
		DeleteContext: resourceIBMEnCOSIntegrationDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"instance_guid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Unique identifier for IBM Cloud Event Notifications instance.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of integration is collect_failed_events.",
			},
			"metadata": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The public or private endpoint for COS bucket",
						},
						"crn": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The CRN of the Cloud Object Storage instance",
						},
						"bucket_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Cloud Object Storage bucket name",
						},
					},
				},
			},
			"integration_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier for Integration.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last updated time.",
			},
		},
	}
}

func resourceIBMEnCOSIntegrationCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_integration_cos", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
		// return diag.FromErr(err)
	}

	options := &en.CreateIntegrationOptions{}

	options.SetInstanceID(d.Get("instance_guid").(string))
	options.SetType(d.Get("type").(string))

	if _, ok := d.GetOk("metadata"); ok {
		metadata := ReplaceCOSIntegrationMapMetadata(d.Get("metadata.0").(map[string]interface{}))
		options.SetMetadata(&metadata)
	}

	result, _, err := enClient.CreateIntegrationWithContext(context, options)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateIntegrationWithContext failed: %s", err.Error()), "ibm_en_integration_cos", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
		// return diag.FromErr(fmt.Errorf("CreateIntegrationWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *options.InstanceID, *result.ID))

	return resourceIBMEnCOSIntegrationRead(context, d, meta)
}

func resourceIBMEnCOSIntegrationRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_integration_cos", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
		// return diag.FromErr(err)
	}

	options := &en.GetIntegrationOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_integration_cos", "read")
		return tfErr.GetDiag()
		// return diag.FromErr(err)
	}

	options.SetInstanceID(parts[0])
	options.SetID(parts[1])

	result, response, err := enClient.GetIntegrationWithContext(context, options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId(d.Get("integration_id").(string))
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetIntegrationWithContext failed: %s", err.Error()), "ibm_en_integration_cos", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
		// return diag.FromErr(fmt.Errorf("GetIntegrationWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("instance_guid", options.InstanceID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting instance_guid: %s", err))
	}

	if err = d.Set("integration_id", options.ID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error integration_id: %s", err))
	}
	if err = d.Set("type", result.Type); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting type: %s", err))
	}

	if result.Metadata != nil {
		err = d.Set("metadata", enCOSIntegrationFlattenMetadata(result.Metadata))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting Metadata %s", err))
		}
	}

	if err = d.Set("updated_at", flex.DateTimeToString(result.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting updated_at: %s", err))
	}

	return nil
}

func resourceIBMEnCOSIntegrationUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_integration_cos", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
		// return diag.FromErr(err)
	}

	options := &en.ReplaceIntegrationOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_integration_cos", "update")
		return tfErr.GetDiag()
		// return diag.FromErr(err)
	}

	options.SetInstanceID(parts[0])
	options.SetID(parts[1])
	options.SetType(d.Get("type").(string))

	if ok := d.HasChanges("metadata"); ok {

		if _, ok := d.GetOk("metadata"); ok {
			metadata := ReplaceCOSIntegrationMapMetadataUpdate(d.Get("metadata.0").(map[string]interface{}))
			options.SetMetadata(&metadata)
		}

		_, _, err := enClient.ReplaceIntegrationWithContext(context, options)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ReplaceIntegrationWithContext failed: %s", err.Error()), "ibm_en_integration_cos", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
			// return diag.FromErr(fmt.Errorf("ReplaceIntegrationWithContext failed %s\n%s", err, response))
		}

		return resourceIBMEnIntegrationRead(context, d, meta)
	}

	return nil
}

func ReplaceCOSIntegrationMapMetadata(metadataParams map[string]interface{}) en.IntegrationCreateMetadata {
	metadataconfigParams := new(en.IntegrationCreateMetadata)
	if metadataParams["endpoint"] != nil {
		metadataconfigParams.Endpoint = core.StringPtr(metadataParams["endpoint"].(string))
	}

	if metadataParams["crn"] != nil {
		metadataconfigParams.CRN = core.StringPtr(metadataParams["crn"].(string))
	}

	if metadataParams["bucket_name"] != nil {
		metadataconfigParams.BucketName = core.StringPtr(metadataParams["bucket_name"].(string))
	}

	return *metadataconfigParams
}

func ReplaceCOSIntegrationMapMetadataUpdate(metadataParams map[string]interface{}) en.IntegrationMetadata {
	metadataconfigParams := new(en.IntegrationMetadata)
	if metadataParams["endpoint"] != nil {
		metadataconfigParams.Endpoint = core.StringPtr(metadataParams["endpoint"].(string))
	}

	if metadataParams["crn"] != nil {
		metadataconfigParams.CRN = core.StringPtr(metadataParams["crn"].(string))
	}

	if metadataParams["bucket_name"] != nil {
		metadataconfigParams.BucketName = core.StringPtr(metadataParams["bucket_name"].(string))
	}

	return *metadataconfigParams
}

func resourceIBMEnCOSIntegrationDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return nil
}
