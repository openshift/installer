// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package catalogmanagement

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/catalogmanagementv1"
)

func ResourceIBMCmValidation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMCmValidationCreate,
		ReadContext:   resourceIBMCmValidationRead,
		UpdateContext: resourceIBMCmValidationUpdate,
		DeleteContext: resourceIBMCmValidationDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"version_locator": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Version locator - the version that will be validated.",
			},
			"region": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Validation region.",
			},
			"override_values": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Override values during validation.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"environment_variables": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Environment variables to include in the schematics workspace.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the environment variable.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Value of the environment variable.",
						},
						"secure": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "If the environment variablel should be secure.",
						},
					},
				},
			},
			"schematics": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: "Other values to pass to the schematics workspace.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name for the schematics workspace.",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description for the schematics workspace.",
						},
						"resource_group_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The resource group ID.",
						},
						"terraform_version": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Version of terraform to use in schematics.",
						},
						"region": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Region to use for the schematics installation.",
						},
						"tags": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of tags for the schematics workspace.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"validated": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Data and time of last successful validation.",
			},
			"requested": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Data and time of last validation request.",
			},
			"state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current validation state - <empty>, in_progress, valid, invalid, expired.",
			},
			"last_operation": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last operation (e.g. submit_deployment, generate_installer, install_offering.",
			},
			// "target": &schema.Schema{
			// 	Type:        schema.TypeMap,
			// 	Computed:    true,
			// 	Description: "Validation target information (e.g. cluster_id, region, namespace, etc).  Values will vary by Content type.",
			// 	Elem:        &schema.Schema{Type: schema.TypeString},
			// },
			"message": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Any message needing to be conveyed as part of the validation job.",
			},
			"x_auth_refresh_token": &schema.Schema{
				Deprecated:  "This argument is deprecated because it is now retrieved automatically.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Authentication token used to submit validation job.",
			},
			"revalidate_if_validated": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "If the version should be revalidated if it is already validated.",
			},
			"mark_version_consumable": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If the version should be marked as consumable or \"ready to share\".",
			},
		},
	}
}

func resourceIBMCmValidationCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_validation", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	validateInstallOptions := &catalogmanagementv1.ValidateInstallOptions{}
	var version catalogmanagementv1.Version

	if _, ok := d.GetOk("version_locator"); ok {
		validateInstallOptions.SetVersionLocatorID(d.Get("version_locator").(string))
		validateInstallOptions.SetVersionLocID(d.Get("version_locator").(string))

		getVersionOptions := &catalogmanagementv1.GetVersionOptions{}

		getVersionOptions.SetVersionLocID(d.Get("version_locator").(string))

		offering, response, err := catalogManagementClient.GetVersionWithContext(context, getVersionOptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetVersionWithContext failed %s\n%s", err, response), "ibm_cm_validation", "create")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}

		version = offering.Kinds[0].Versions[0]
	}

	mk := fmt.Sprintf("%s.%s", *version.CatalogID, *version.OfferingID)
	conns.IbmMutexKV.Lock(mk)
	defer conns.IbmMutexKV.Unlock(mk)

	valid := "valid"
	if version.Validation.State == &valid && d.Get("revalidate_if_validated") != true {
		// version already validated and do not wish to revalidate
		d.SetId(*validateInstallOptions.VersionLocID)
		if _, ok := d.GetOk("mark_version_consumable"); ok && d.Get("mark_version_consumable").(bool) {
			err = markVersionAsConsumable(version, context, meta)
			if err != nil {
				d.SetId("")
				tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_validation", "create")
				log.Printf("[DEBUG]\\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
		}

		return resourceIBMCmValidationRead(context, d, meta)
	}

	if _, ok := d.GetOk("region"); ok {
		validateInstallOptions.SetRegion(d.Get("region").(string))
	}
	if _, ok := d.GetOk("override_values"); ok {
		overridesModel, err := configureOverrides(d.Get("override_values").(map[string]interface{}))
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_validation", "create")
			log.Printf("[DEBUG]\\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		validateInstallOptions.SetOverrideValues(&overridesModel)
	}
	if _, ok := d.GetOk("environment_variables"); ok {
		envsModel, err := envVariablesToDeployRequestBodyEnvVariables(d.Get("environment_variables").([]interface{}))
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_validation", "create")
			log.Printf("[DEBUG]\\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		validateInstallOptions.SetEnvironmentVariables(envsModel)
	}
	if _, ok := d.GetOk("schematics.0"); ok {
		schematicsModel, err := schematicsMapToDeployRequestBodySchematics(d.Get("schematics.0").(map[string]interface{}))
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_validation", "create")
			log.Printf("[DEBUG]\\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		validateInstallOptions.SetSchematics(&schematicsModel)
	}

	bxSession, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_validation", "create")
		log.Printf("[DEBUG]\\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	validateInstallOptions.SetXAuthRefreshToken(bxSession.Config.IAMRefreshToken)

	response, err := catalogManagementClient.ValidateInstallWithContext(context, validateInstallOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ValidateInstallWithContext failed %s\n%s", err, response), "ibm_cm_validation", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*validateInstallOptions.VersionLocID)

	validationStatusOptions := &catalogmanagementv1.GetValidationStatusOptions{}
	validationStatusOptions.SetVersionLocID(*validateInstallOptions.VersionLocID)
	validationStatusOptions.SetXAuthRefreshToken(bxSession.Config.IAMRefreshToken)
	result, response, err := catalogManagementClient.GetValidationStatusWithContext(context, validationStatusOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetValidationStatusWithContext failed %s\n%s", err, response), "ibm_cm_validation", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	status := *result.State

	// Track progress of validation through schematics workspace status
	// Do a GET every 5 seconds to check for an updated status
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		newStatus := *result.State
		if status != newStatus {
			status = newStatus
		}
		log.Printf("[DEBUG] Status is %s\n", status)
		// Break from loop if installation is successful or fails
		if status == "valid" || status == "invalid" || status == "expired" {
			ticker.Stop()
			break
		}

		result, response, err = catalogManagementClient.GetValidationStatusWithContext(context, validationStatusOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetValidationStatusWithContext failed %s\n%s", err, response), "ibm_cm_validation", "create")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	// mark consumable if specified and validation passed
	if _, ok := d.GetOk("mark_version_consumable"); ok && d.Get("mark_version_consumable").(bool) && status == "valid" {
		err = markVersionAsConsumable(version, context, meta)
		if err != nil {
			d.SetId("")
			tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_validation", "create")
			log.Printf("[DEBUG]\\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIBMCmValidationRead(context, d, meta)
}

func resourceIBMCmValidationRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_validation", "read")
		log.Printf("[DEBUG]\\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getVersionOptions := &catalogmanagementv1.GetVersionOptions{}

	getVersionOptions.SetVersionLocID(d.Id())

	offering, response, err := catalogManagementClient.GetVersionWithContext(context, getVersionOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetVersionWithContext failed %s\n%s", err, response), "ibm_cm_validation", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	version := offering.Kinds[0].Versions[0]

	if err = d.Set("version_locator", version.VersionLocator); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting version_locator: %s", err), "ibm_cm_validation", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if version.Validation != nil && version.Validation.Validated != nil {
		if err = d.Set("validated", version.Validation.Validated.String()); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting validation: %s", err), "ibm_cm_validation", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	if version.Validation != nil && version.Validation.Requested != nil {
		if err = d.Set("requested", version.Validation.Requested.String()); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting requested: %s", err), "ibm_cm_validation", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	if version.Validation != nil && version.Validation.State != nil {
		if err = d.Set("state", version.Validation.State); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting state: %s", err), "ibm_cm_validation", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	if version.Validation != nil && version.Validation.LastOperation != nil {
		if err = d.Set("last_operation", version.Validation.LastOperation); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting last_operation: %s", err), "ibm_cm_validation", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	if version.Validation != nil && version.Validation.Message != nil {
		if err = d.Set("message", version.Validation.Message); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting message: %s", err), "ibm_cm_validation", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	// if version.Validation != nil && version.Validation.Target != nil {
	// 	if err = d.Set("target", version.Validation.Target); err != nil {
	// 		return diag.FromErr(fmt.Errorf("Error setting target: %s", err))
	// 	}
	// }

	return nil
}

func resourceIBMCmValidationUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceIBMCmValidationDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func configureOverrides(overrides map[string]interface{}) (catalogmanagementv1.DeployRequestBodyOverrideValues, error) {
	overridesModel := catalogmanagementv1.DeployRequestBodyOverrideValues{}
	overridesModel.SetProperties(overrides)
	if overrides["vsi_instance_name"] != nil && overrides["vsi_instance_name"].(string) != "" {
		overridesModel.VsiInstanceName = core.StringPtr(overrides["vsi_instance_name"].(string))
	}
	if overrides["vpc_profile"] != nil && overrides["vpc_profile"].(string) != "" {
		overridesModel.VPCProfile = core.StringPtr(overrides["vpc_profile"].(string))
	}
	if overrides["subnet_id"] != nil && overrides["subnet_id"].(string) != "" {
		overridesModel.SubnetID = core.StringPtr(overrides["subnet_id"].(string))
	}
	if overrides["vpc_id"] != nil && overrides["vpc_id"].(string) != "" {
		overridesModel.VPCID = core.StringPtr(overrides["vpc_id"].(string))
	}
	if overrides["subnet_zone"] != nil && overrides["subnet_zone"].(string) != "" {
		overridesModel.SubnetZone = core.StringPtr(overrides["subnet_zone"].(string))
	}
	if overrides["ssh_key_id"] != nil && overrides["ssh_key_id"].(string) != "" {
		overridesModel.SSHKeyID = core.StringPtr(overrides["ssh_key_id"].(string))
	}
	if overrides["vpc_region"] != nil && overrides["vpc_region"].(string) != "" {
		overridesModel.VPCRegion = core.StringPtr(overrides["vpc_region"].(string))
	}
	return overridesModel, nil
}

func markVersionAsConsumable(version catalogmanagementv1.Version, context context.Context, meta interface{}) error {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		return err
	}

	consumableVersionOptions := catalogmanagementv1.ConsumableVersionOptions{}
	consumableVersionOptions.SetVersionLocID(*version.VersionLocator)

	_, err = catalogManagementClient.ConsumableVersionWithContext(context, &consumableVersionOptions)
	if err != nil {
		return err
	}

	return nil
}

func envVariablesToDeployRequestBodyEnvVariables(envVariables []interface{}) ([]catalogmanagementv1.DeployRequestBodyEnvironmentVariablesItem, error) {
	var modelArr []catalogmanagementv1.DeployRequestBodyEnvironmentVariablesItem
	for _, envVar := range envVariables {
		if envVar != nil {
			model := catalogmanagementv1.DeployRequestBodyEnvironmentVariablesItem{}
			if envVar.(map[string]interface{})["name"] != nil && envVar.(map[string]interface{})["name"].(string) != "" {
				model.Name = core.StringPtr(envVar.(map[string]interface{})["name"].(string))
			}
			if envVar.(map[string]interface{})["value"] != nil {
				model.Value = envVar.(map[string]interface{})["value"]
			}
			if envVar.(map[string]interface{})["secure"] != nil {
				model.Secure = core.BoolPtr(envVar.(map[string]interface{})["secure"].(bool))
			}
			modelArr = append(modelArr, model)
		}
	}
	return modelArr, nil
}

func schematicsMapToDeployRequestBodySchematics(schematicsMap map[string]interface{}) (catalogmanagementv1.DeployRequestBodySchematics, error) {
	model := catalogmanagementv1.DeployRequestBodySchematics{}
	if schematicsMap["name"] != nil && schematicsMap["name"].(string) != "" {
		model.Name = core.StringPtr(schematicsMap["name"].(string))
	}
	if schematicsMap["description"] != nil && schematicsMap["description"].(string) != "" {
		model.Description = core.StringPtr(schematicsMap["description"].(string))
	}
	if schematicsMap["resource_group_id"] != nil && schematicsMap["resource_group_id"].(string) != "" {
		model.ResourceGroupID = core.StringPtr(schematicsMap["resource_group_id"].(string))
	}
	if schematicsMap["terraform_version"] != nil && schematicsMap["terraform_version"].(string) != "" {
		model.TerraformVersion = core.StringPtr(schematicsMap["terraform_version"].(string))
	}
	if schematicsMap["region"] != nil && schematicsMap["region"].(string) != "" {
		model.Region = core.StringPtr(schematicsMap["region"].(string))
	}
	if schematicsMap["tags"] != nil && len(schematicsMap["tags"].([]interface{})) > 0 {
		model.Tags = schematicsMap["name"].([]string)
	}
	return model, nil
}
