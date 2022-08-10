// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package atracker

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/platform-services-go-sdk/atrackerv2"
)

func ResourceIBMAtrackerSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMAtrackerSettingsCreate,
		ReadContext:   resourceIBMAtrackerSettingsRead,
		UpdateContext: resourceIBMAtrackerSettingsUpdate,
		DeleteContext: resourceIBMAtrackerSettingsDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"metadata_region_primary": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_atracker_settings", "metadata_region_primary"),
				Description:  "To store all your meta data in a single region.",
			},
			"private_api_endpoint_only": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "If you set this true then you cannot access api through public network.",
			},
			"default_targets": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The target ID List. In the event that no routing rule causes the event to be sent to a target, these targets will receive the event.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"permitted_target_regions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "If present then only these regions may be used to define a target.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			// Future Planned support
			// "metadata_region_backup": &schema.Schema{
			// 	Type:         schema.TypeString,
			// 	Optional:     true,
			// 	ValidateFunc: validate.InvokeValidator("ibm_atracker_settings", "metadata_region_backup"),
			// 	Description:  "Provide a back up region to store meta data.",
			// },
			"api_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The lowest API version of targets or routes that customer might have under his or her account.",
			},
		},
	}
}

func ResourceIBMAtrackerSettingsValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "metadata_region_primary",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-zA-Z0-9 -_]`,
			MinValueLength:             3,
			MaxValueLength:             256,
		},
		// Future Planned support
		// validate.ValidateSchema{
		// 	Identifier:                 "metadata_region_backup",
		// 	ValidateFunctionIdentifier: validate.ValidateRegexpLen,
		// 	Type:                       validate.TypeString,
		// 	Optional:                   true,
		// 	Regexp:                     `^[a-zA-Z0-9 -_]`,
		// 	MinValueLength:             3,
		// 	MaxValueLength:             256,
		// },
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_atracker_settings", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMAtrackerSettingsCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	atrackerClient, err := meta.(conns.ClientSession).AtrackerV2()
	if err != nil {
		return diag.FromErr(err)
	}

	putSettingsOptions := &atrackerv2.PutSettingsOptions{}

	putSettingsOptions.SetMetadataRegionPrimary(d.Get("metadata_region_primary").(string))
	putSettingsOptions.SetPrivateAPIEndpointOnly(d.Get("private_api_endpoint_only").(bool))
	if _, ok := d.GetOk("default_targets"); ok {
		putSettingsOptions.SetDefaultTargets(resourceInterfaceToStringArray(d.Get("default_targets").([]interface{})))
	}
	if _, ok := d.GetOk("permitted_target_regions"); ok {
		putSettingsOptions.SetPermittedTargetRegions(resourceInterfaceToStringArray(d.Get("permitted_target_regions").([]interface{})))
	}
	// Future planned support
	// if _, ok := d.GetOk("metadata_region_backup"); ok {
	// 	putSettingsOptions.SetMetadataRegionBackup(d.Get("metadata_region_backup").(string))
	// }

	settings, response, err := atrackerClient.PutSettingsWithContext(context, putSettingsOptions)
	if err != nil {
		log.Printf("[DEBUG] PutSettingsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("PutSettingsWithContext failed %s\n%s", err, response))
	}

	d.SetId(*settings.MetadataRegionPrimary)

	return resourceIBMAtrackerSettingsRead(context, d, meta)
}

func resourceIBMAtrackerSettingsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, atrackerClient, err := getAtrackerClients(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	getSettingsOptions := &atrackerv2.GetSettingsOptions{}

	settings, response, err := atrackerClient.GetSettingsWithContext(context, getSettingsOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetSettingsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetSettingsWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("metadata_region_primary", settings.MetadataRegionPrimary); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting metadata_region_primary: %s", err))
	}
	if err = d.Set("private_api_endpoint_only", settings.PrivateAPIEndpointOnly); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting private_api_endpoint_only: %s", err))
	}
	if settings.DefaultTargets != nil {
		if err = d.Set("default_targets", settings.DefaultTargets); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting default_targets: %s", err))
		}
	}
	if settings.PermittedTargetRegions != nil {
		if err = d.Set("permitted_target_regions", settings.PermittedTargetRegions); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting permitted_target_regions: %s", err))
		}
	}
	// Future planned support
	// if err = d.Set("metadata_region_backup", settings.MetadataRegionBackup); err != nil {
	// 	return diag.FromErr(fmt.Errorf("Error setting metadata_region_backup: %s", err))
	// }
	if err = d.Set("api_version", flex.IntValue(settings.APIVersion)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting api_version: %s", err))
	}

	return nil
}

func resourceIBMAtrackerSettingsUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	atrackerClient, err := meta.(conns.ClientSession).AtrackerV2()
	if err != nil {
		return diag.FromErr(err)
	}

	putSettingsOptions := &atrackerv2.PutSettingsOptions{}

	hasChange := false
	newMetaDataRegionPrimary := d.Get("metadata_region_primary").(string)
	putSettingsOptions.SetMetadataRegionPrimary(newMetaDataRegionPrimary)
	putSettingsOptions.SetPrivateAPIEndpointOnly(d.Get("private_api_endpoint_only").(bool))
	hasChange = hasChange || d.HasChange("metadata_region_primary") || d.HasChange("private_api_endpoint_only") || d.HasChange("metadata_region_primary") || d.HasChange("permitted_target_regions") || d.HasChange("default_targets")

	if d.HasChange("metadata_region_primary") {
		d.SetId(newMetaDataRegionPrimary)
	}
	putSettingsOptions.DefaultTargets = resourceInterfaceToStringArray(d.Get("default_targets").([]interface{}))

	putSettingsOptions.PermittedTargetRegions = resourceInterfaceToStringArray(d.Get("permitted_target_regions").([]interface{}))

	// Future planned support
	// if d.HasChange("metadata_region_backup") {
	// 	putSettingsOptions.SetMetadataRegionBackup(d.Get("metadata_region_backup").(string))
	// 	hasChange = true
	// }

	if hasChange {
		setting, response, err := atrackerClient.PutSettingsWithContext(context, putSettingsOptions)
		if err != nil {
			log.Printf("[DEBUG] PutSettingsWithContext failed %s\n%s", err, response)
			log.Printf("[DEBUG] PutSettingsWithContext failed %v\n", putSettingsOptions)
			return diag.FromErr(fmt.Errorf("PutSettingsWithContext failed %s\n%s", err, response))
		}
		d.SetId(*setting.MetadataRegionPrimary)
	}

	return resourceIBMAtrackerSettingsRead(context, d, meta)
}

func resourceInterfaceToStringArray(resources []interface{}) (result []string) {
	result = make([]string, 0)
	for _, item := range resources {
		if item != nil {
			result = append(result, item.(string))
		} else {
			result = append(result, "")
		}
	}
	return result
}

func resourceIBMAtrackerSettingsDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	atrackerClient, err := meta.(conns.ClientSession).AtrackerV2()
	if err != nil {
		return diag.FromErr(err)
	}

	// Retrieve old settings and put them for required fields.  Remove all other fields
	settings, getResponse, err := atrackerClient.GetSettingsWithContext(context, &atrackerv2.GetSettingsOptions{})
	if err != nil {
		log.Printf("[DEBUG] PutSettingsWithContext with GetSettingsWithContext failed %s\n%s", err, getResponse)
		return diag.FromErr(fmt.Errorf("GetSettingsWithContext failed %s\n%s", err, getResponse))
	}
	putSettingsOptions := &atrackerv2.PutSettingsOptions{}

	putSettingsOptions.MetadataRegionPrimary = settings.MetadataRegionPrimary
	putSettingsOptions.PrivateAPIEndpointOnly = settings.PrivateAPIEndpointOnly
	putSettingsOptions.PermittedTargetRegions = []string{}
	putSettingsOptions.DefaultTargets = []string{}

	_, response, err := atrackerClient.PutSettingsWithContext(context, putSettingsOptions)
	if err != nil {
		log.Printf("[DEBUG] PutSettingsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("PutSettingsWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}
