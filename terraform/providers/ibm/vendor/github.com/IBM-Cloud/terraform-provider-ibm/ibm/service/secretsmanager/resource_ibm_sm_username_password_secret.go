// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package secretsmanager

import (
	"context"
	"fmt"
	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/secrets-manager-go-sdk/v2/secretsmanagerv2"
)

func ResourceIbmSmUsernamePasswordSecret() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmSmUsernamePasswordSecretCreate,
		ReadContext:   resourceIbmSmUsernamePasswordSecretRead,
		UpdateContext: resourceIbmSmUsernamePasswordSecretUpdate,
		DeleteContext: resourceIbmSmUsernamePasswordSecretDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"custom_metadata": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Description: "The secret metadata that a user can customize.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An extended description of your secret.To protect your privacy, do not use personal data, such as your name or location, as a description for your secret group.",
			},
			"expiration_date": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The date a secret is expired. The date format follows RFC 3339.",
			},
			"labels": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Labels that you can use to search for secrets in your instance.Up to 30 labels can be created.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "A human-readable name to assign to your secret.To protect your privacy, do not use personal data, such as your name or location, as a name for your secret.",
			},
			"secret_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "A v4 UUID identifier, or `default` secret group.",
			},
			"secret_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A v4 UUID identifier.",
			},
			"secret_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The secret type. Supported types are arbitrary, certificates (imported, public, and private), IAM credentials, key-value, and user credentials.",
			},
			"version_custom_metadata": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The secret version metadata that a user can customize.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"rotation": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "Determines whether Secrets Manager rotates your secrets automatically.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_rotate": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Determines whether Secrets Manager rotates your secret automatically.Default is `false`. If `auto_rotate` is set to `true` the service rotates your secret based on the defined interval.",
						},
						"interval": &schema.Schema{
							Type:             schema.TypeInt,
							Optional:         true,
							Computed:         true,
							Description:      "The length of the secret rotation time interval.",
							DiffSuppressFunc: rotationAttributesDiffSuppress,
						},
						"unit": &schema.Schema{
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							Description:      "The units for the secret rotation time interval.",
							DiffSuppressFunc: rotationAttributesDiffSuppress,
						},
					},
				},
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The username that is assigned to the secret.",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: "The password that is assigned to the secret.",
			},
			"password_generation_policy": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "Policy for auto-generated passwords.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"length": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     32,
							Description: "The length of auto-generated passwords.",
						},
						"include_digits": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Include digits in auto-generated passwords.",
						},
						"include_symbols": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Include symbols in auto-generated passwords.",
						},
						"include_uppercase": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Include uppercase letters in auto-generated passwords.",
						},
					},
				},
			},
			"created_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier that is associated with the entity that created the secret.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when a resource was created. The date format follows RFC 3339.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A CRN that uniquely identifies an IBM Cloud resource.",
			},
			"downloaded": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the secret data that is associated with a secret version was retrieved in a call to the service API.",
			},
			"locks_total": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of locks of the secret.",
			},
			"state": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The secret state that is based on NIST SP 800-57. States are integers and correspond to the `Pre-activation = 0`, `Active = 1`,  `Suspended = 2`, `Deactivated = 3`, and `Destroyed = 5` values.",
			},
			"state_description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A text representation of the secret state.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when a resource was recently modified. The date format follows RFC 3339.",
			},
			"versions_total": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of versions of the secret.",
			},
			"next_rotation_date": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date that the secret is scheduled for automatic rotation.The service automatically creates a new version of the secret on its next rotation date. This field exists only for secrets that have an existing rotation policy.",
			},
		},
	}
}

func resourceIbmSmUsernamePasswordSecretCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", UsernamePasswordSecretResourceName, "create")
		return tfErr.GetDiag()
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	createSecretOptions := &secretsmanagerv2.CreateSecretOptions{}

	secretPrototypeModel, err := resourceIbmSmUsernamePasswordSecretMapToSecretPrototype(d)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", UsernamePasswordSecretResourceName, "create")
		return tfErr.GetDiag()
	}
	createSecretOptions.SetSecretPrototype(secretPrototypeModel)

	secretIntf, response, err := secretsManagerClient.CreateSecretWithContext(context, createSecretOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateSecretWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateSecretWithContext failed: %s\n%s", err.Error(), response), UsernamePasswordSecretResourceName, "create")
		return tfErr.GetDiag()
	}
	secret := secretIntf.(*secretsmanagerv2.UsernamePasswordSecret)
	d.SetId(fmt.Sprintf("%s/%s/%s", region, instanceId, *secret.ID))
	d.Set("secret_id", *secret.ID)

	_, err = waitForIbmSmUsernamePasswordSecretCreate(secretsManagerClient, d)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error waiting for resource IbmSmUsernamePasswordSecret (%s) to be created: %s", d.Id(), err.Error()), UsernamePasswordSecretResourceName, "create")
		return tfErr.GetDiag()
	}

	return resourceIbmSmUsernamePasswordSecretRead(context, d, meta)
}

func waitForIbmSmUsernamePasswordSecretCreate(secretsManagerClient *secretsmanagerv2.SecretsManagerV2, d *schema.ResourceData) (interface{}, error) {
	getSecretOptions := &secretsmanagerv2.GetSecretOptions{}
	id := strings.Split(d.Id(), "/")
	secretId := id[2]
	getSecretOptions.SetID(secretId)

	stateConf := &resource.StateChangeConf{
		Pending: []string{"pre_activation"},
		Target:  []string{"active"},
		Refresh: func() (interface{}, string, error) {
			stateObjIntf, response, err := secretsManagerClient.GetSecret(getSecretOptions)
			stateObj := stateObjIntf.(*secretsmanagerv2.UsernamePasswordSecret)
			if err != nil {
				if apiErr, ok := err.(bmxerror.RequestFailure); ok && apiErr.StatusCode() == 404 {
					return nil, "", fmt.Errorf("The instance %s does not exist anymore: %s\n%s", "getSecretOptions", err, response)
				}
				return nil, "", err
			}
			failStates := map[string]bool{"destroyed": true}
			if failStates[*stateObj.StateDescription] {
				return stateObj, *stateObj.StateDescription, fmt.Errorf("The instance %s failed: %s\n%s", "getSecretOptions", err, response)
			}
			return stateObj, *stateObj.StateDescription, nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	return stateConf.WaitForState()
}

func resourceIbmSmUsernamePasswordSecretRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", UsernamePasswordSecretResourceName, "read")
		return tfErr.GetDiag()
	}

	id := strings.Split(d.Id(), "/")
	if len(id) != 3 {
		tfErr := flex.TerraformErrorf(nil, "Wrong format of resource ID. To import a secret use the format `<region>/<instance_id>/<secret_id>`", UsernamePasswordSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	region := id[0]
	instanceId := id[1]
	secretId := id[2]
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	getSecretOptions := &secretsmanagerv2.GetSecretOptions{}

	getSecretOptions.SetID(secretId)

	secretIntf, response, err := secretsManagerClient.GetSecretWithContext(context, getSecretOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetSecretWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetSecretWithContext failed %s\n%s", err, response), UsernamePasswordSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	secret := secretIntf.(*secretsmanagerv2.UsernamePasswordSecret)

	if err = d.Set("secret_id", secretId); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_id"), UsernamePasswordSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("instance_id", instanceId); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting instance_id"), UsernamePasswordSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("region", region); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting region"), UsernamePasswordSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("created_by", secret.CreatedBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_by"), UsernamePasswordSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("created_at", DateTimeToRFC3339(secret.CreatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_at"), UsernamePasswordSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("crn", secret.Crn); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting crn"), UsernamePasswordSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if secret.CustomMetadata != nil {
		d.Set("custom_metadata", secret.CustomMetadata)
	}
	if err = d.Set("description", secret.Description); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting description"), UsernamePasswordSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("downloaded", secret.Downloaded); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting downloaded"), UsernamePasswordSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if secret.Labels != nil {
		if err = d.Set("labels", secret.Labels); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting labels"), UsernamePasswordSecretResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("locks_total", flex.IntValue(secret.LocksTotal)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting locks_total"), UsernamePasswordSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("name", secret.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name"), UsernamePasswordSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("secret_group_id", secret.SecretGroupID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_group_id"), UsernamePasswordSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("secret_type", secret.SecretType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_type"), UsernamePasswordSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("state", flex.IntValue(secret.State)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting state"), UsernamePasswordSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("state_description", secret.StateDescription); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting state_description"), UsernamePasswordSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("updated_at", DateTimeToRFC3339(secret.UpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated_at"), UsernamePasswordSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("versions_total", flex.IntValue(secret.VersionsTotal)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting versions_total"), UsernamePasswordSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	rotationMap, err := resourceIbmSmUsernamePasswordSecretRotationPolicyToMap(secret.Rotation)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", UsernamePasswordSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("rotation", []map[string]interface{}{rotationMap}); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting rotation"), UsernamePasswordSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("expiration_date", DateTimeToRFC3339(secret.ExpirationDate)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting expiration_date"), UsernamePasswordSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("next_rotation_date", DateTimeToRFC3339(secret.NextRotationDate)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting next_rotation_date"), UsernamePasswordSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("username", secret.Username); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting username"), UsernamePasswordSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("password", secret.Password); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting password"), UsernamePasswordSecretResourceName, "read")
		return tfErr.GetDiag()
	}

	passwordPolicyMap, err := passwordGenerationPolicyToMap(secret.PasswordGenerationPolicy)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", UsernamePasswordSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("password_generation_policy", []map[string]interface{}{passwordPolicyMap}); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting password_generation_policy"), UsernamePasswordSecretResourceName, "read")
		return tfErr.GetDiag()
	}

	// Call get version metadata API to get the current version_custom_metadata
	getVersionMetdataOptions := &secretsmanagerv2.GetSecretVersionMetadataOptions{}
	getVersionMetdataOptions.SetSecretID(secretId)
	getVersionMetdataOptions.SetID("current")

	versionMetadataIntf, response, err := secretsManagerClient.GetSecretVersionMetadataWithContext(context, getVersionMetdataOptions)
	if err != nil {
		log.Printf("[DEBUG] GetSecretVersionMetadataWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetSecretVersionMetadataWithContext failed %s\n%s", err, response), UsernamePasswordSecretResourceName, "read")
		return tfErr.GetDiag()
	}

	versionMetadata := versionMetadataIntf.(*secretsmanagerv2.UsernamePasswordSecretVersionMetadata)
	if versionMetadata.VersionCustomMetadata != nil {
		if err = d.Set("version_custom_metadata", versionMetadata.VersionCustomMetadata); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting version_custom_metadata"), UsernamePasswordSecretResourceName, "read")
			return tfErr.GetDiag()
		}
	}

	return nil
}

func resourceIbmSmUsernamePasswordSecretUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", UsernamePasswordSecretResourceName, "update")
		return tfErr.GetDiag()
	}

	id := strings.Split(d.Id(), "/")
	region := id[0]
	instanceId := id[1]
	secretId := id[2]
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	updateSecretMetadataOptions := &secretsmanagerv2.UpdateSecretMetadataOptions{}

	updateSecretMetadataOptions.SetID(secretId)

	hasChange := false

	patchVals := &secretsmanagerv2.UsernamePasswordSecretMetadataPatch{}

	if d.HasChange("name") {
		patchVals.Name = core.StringPtr(d.Get("name").(string))
		hasChange = true
	}
	if d.HasChange("description") {
		patchVals.Description = core.StringPtr(d.Get("description").(string))
		hasChange = true
	}
	if d.HasChange("labels") {
		labels := d.Get("labels").([]interface{})
		labelsParsed := make([]string, len(labels))
		for i, v := range labels {
			labelsParsed[i] = fmt.Sprint(v)
		}
		patchVals.Labels = labelsParsed
		hasChange = true
	}
	if d.HasChange("custom_metadata") {
		patchVals.CustomMetadata = d.Get("custom_metadata").(map[string]interface{})
		hasChange = true
	}
	if d.HasChange("rotation") {
		RotationModel, err := resourceIbmSmUsernamePasswordSecretMapToRotationPolicy(d.Get("rotation").([]interface{})[0].(map[string]interface{}))
		if err != nil {
			log.Printf("[DEBUG] UpdateSecretMetadataWithContext failed: Reading Rotation parameter failed: %s", err)
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateSecretMetadataWithContext failed: Reading Rotation parameter failed: %s", err), UsernamePasswordSecretResourceName, "update")
			return tfErr.GetDiag()
		}
		patchVals.Rotation = RotationModel
		hasChange = true
	}

	if d.HasChange("expiration_date") {
		if _, ok := d.GetOk("expiration_date"); ok {
			layout := time.RFC3339
			parseToTime, err := time.Parse(layout, d.Get("expiration_date").(string))
			if err != nil {
				tfErr := flex.TerraformErrorf(err, fmt.Sprintf(`Failed to get "expiration_date". Error: %s`, err.Error()), UsernamePasswordSecretResourceName, "update")
				return tfErr.GetDiag()
			}
			parseToDateTime := strfmt.DateTime(parseToTime)
			patchVals.ExpirationDate = &parseToDateTime
			hasChange = true
		} else {
			tfErr := flex.TerraformErrorf(nil, `The "expiration_date" field cannot be removed. To disable expiration set expiration date to a far future date'`, UsernamePasswordSecretResourceName, "update")
			return tfErr.GetDiag()
		}
	}

	if d.HasChange("password_generation_policy") {
		passwordPolicyModel, err := mapToPasswordGenerationPolicyPatch(d.Get("password_generation_policy").([]interface{})[0].(map[string]interface{}))
		if err != nil {
			log.Printf("[DEBUG] UpdateSecretMetadataWithContext failed: Reading password_generation_policy parameter failed: %s", err)
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateSecretMetadataWithContext failed: Reading password_generation_policy parameter failed: %s", err), UsernamePasswordSecretResourceName, "update")
			return tfErr.GetDiag()
		}
		patchVals.PasswordGenerationPolicy = passwordPolicyModel
		hasChange = true
	}

	if hasChange {
		updateSecretMetadataOptions.SecretMetadataPatch, _ = patchVals.AsPatch()
		_, response, err := secretsManagerClient.UpdateSecretMetadataWithContext(context, updateSecretMetadataOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateSecretMetadataWithContext failed %s\n%s", err, response)
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateSecretMetadataWithContext failed %s\n%s", err, response), UsernamePasswordSecretResourceName, "update")
			return tfErr.GetDiag()
		}
	}

	// Apply change in payload (if changed)
	if d.HasChange("password") {
		versionModel := &secretsmanagerv2.UsernamePasswordSecretVersionPrototype{}
		versionModel.Password = core.StringPtr(d.Get("password").(string))
		if _, ok := d.GetOk("version_custom_metadata"); ok {
			versionModel.VersionCustomMetadata = d.Get("version_custom_metadata").(map[string]interface{})
		}
		if _, ok := d.GetOk("custom_metadata"); ok {
			versionModel.CustomMetadata = d.Get("custom_metadata").(map[string]interface{})
		}

		createSecretVersionOptions := &secretsmanagerv2.CreateSecretVersionOptions{}
		createSecretVersionOptions.SetSecretID(secretId)
		createSecretVersionOptions.SetSecretVersionPrototype(versionModel)
		_, response, err := secretsManagerClient.CreateSecretVersionWithContext(context, createSecretVersionOptions)
		if err != nil {
			if hasChange {
				// Before returning an error, call the read function to update the Terraform state with the change
				// that was already applied to the metadata
				resourceIbmSmUsernamePasswordSecretRead(context, d, meta)
			}
			log.Printf("[DEBUG] CreateSecretVersionWithContext failed %s\n%s", err, response)
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateSecretVersionWithContext failed %s\n%s", err, response), UsernamePasswordSecretResourceName, "update")
			return tfErr.GetDiag()
		}
	} else if d.HasChange("version_custom_metadata") {
		// Apply change to version_custom_metadata in current version
		secretVersionMetadataPatchModel := new(secretsmanagerv2.SecretVersionMetadataPatch)
		secretVersionMetadataPatchModel.VersionCustomMetadata = d.Get("version_custom_metadata").(map[string]interface{})
		secretVersionMetadataPatchModelAsPatch, _ := secretVersionMetadataAsPatchFunction(secretVersionMetadataPatchModel)

		updateSecretVersionOptions := &secretsmanagerv2.UpdateSecretVersionMetadataOptions{}
		updateSecretVersionOptions.SetSecretID(secretId)
		updateSecretVersionOptions.SetID("current")
		updateSecretVersionOptions.SetSecretVersionMetadataPatch(secretVersionMetadataPatchModelAsPatch)
		_, response, err := secretsManagerClient.UpdateSecretVersionMetadataWithContext(context, updateSecretVersionOptions)
		if err != nil {
			if hasChange {
				// Call the read function to update the Terraform state with the change already applied to the metadata
				resourceIbmSmUsernamePasswordSecretRead(context, d, meta)
			}
			log.Printf("[DEBUG] UpdateSecretVersionMetadataWithContext failed %s\n%s", err, response)
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateSecretVersionMetadataWithContext failed %s\n%s", err, response), UsernamePasswordSecretResourceName, "update")
			return tfErr.GetDiag()
		}
	}

	return resourceIbmSmUsernamePasswordSecretRead(context, d, meta)
}

func resourceIbmSmUsernamePasswordSecretDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", UsernamePasswordSecretResourceName, "delete")
		return tfErr.GetDiag()
	}

	id := strings.Split(d.Id(), "/")
	region := id[0]
	instanceId := id[1]
	secretId := id[2]
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	deleteSecretOptions := &secretsmanagerv2.DeleteSecretOptions{}

	deleteSecretOptions.SetID(secretId)

	response, err := secretsManagerClient.DeleteSecretWithContext(context, deleteSecretOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteSecretWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteSecretWithContext failed %s\n%s", err, response), UsernamePasswordSecretResourceName, "delete")
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func resourceIbmSmUsernamePasswordSecretMapToSecretPrototype(d *schema.ResourceData) (*secretsmanagerv2.UsernamePasswordSecretPrototype, error) {
	model := &secretsmanagerv2.UsernamePasswordSecretPrototype{}
	model.SecretType = core.StringPtr("username_password")

	if _, ok := d.GetOk("name"); ok {
		model.Name = core.StringPtr(d.Get("name").(string))
	}
	if _, ok := d.GetOk("custom_metadata"); ok {
		model.CustomMetadata = d.Get("custom_metadata").(map[string]interface{})
	}
	if _, ok := d.GetOk("description"); ok {
		model.Description = core.StringPtr(d.Get("description").(string))
	}
	if _, ok := d.GetOk("expiration_date"); ok {
		layout := time.RFC3339
		parseToTime, err := time.Parse(layout, d.Get("expiration_date").(string))
		if err != nil {
			return nil, fmt.Errorf(`Failed to get "expiration_date". Error: ` + err.Error())
		}
		parseToDateTime := strfmt.DateTime(parseToTime)
		model.ExpirationDate = &parseToDateTime
	}
	if _, ok := d.GetOk("labels"); ok {
		labels := d.Get("labels").([]interface{})
		labelsParsed := make([]string, len(labels))
		for i, v := range labels {
			labelsParsed[i] = fmt.Sprint(v)
		}
		model.Labels = labelsParsed
	}
	if _, ok := d.GetOk("secret_group_id"); ok {
		model.SecretGroupID = core.StringPtr(d.Get("secret_group_id").(string))
	}
	if _, ok := d.GetOk("version_custom_metadata"); ok {
		model.VersionCustomMetadata = d.Get("version_custom_metadata").(map[string]interface{})
	}
	if _, ok := d.GetOk("username"); ok {
		model.Username = core.StringPtr(d.Get("username").(string))
	}
	if _, ok := d.GetOk("password"); ok {
		model.Password = core.StringPtr(d.Get("password").(string))
	}
	if _, ok := d.GetOk("rotation"); ok {
		RotationModel, err := resourceIbmSmUsernamePasswordSecretMapToRotationPolicy(d.Get("rotation").([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Rotation = RotationModel
	}
	if _, ok := d.GetOk("password_generation_policy"); ok {
		passwordPolicyModel, err := mapToPasswordGenerationPolicy(d.Get("password_generation_policy").([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.PasswordGenerationPolicy = passwordPolicyModel
	}

	return model, nil
}

func resourceIbmSmUsernamePasswordSecretMapToRotationPolicy(modelMap map[string]interface{}) (secretsmanagerv2.RotationPolicyIntf, error) {
	model := &secretsmanagerv2.RotationPolicy{}
	if modelMap["auto_rotate"] != nil {
		model.AutoRotate = core.BoolPtr(modelMap["auto_rotate"].(bool))
	}
	if modelMap["interval"].(int) == 0 {
		model.Interval = nil
	} else {
		model.Interval = core.Int64Ptr(int64(modelMap["interval"].(int)))
	}
	if modelMap["unit"] != nil && modelMap["unit"].(string) != "" {
		model.Unit = core.StringPtr(modelMap["unit"].(string))
	}
	return model, nil
}

func resourceIbmSmUsernamePasswordSecretRotationPolicyToMap(model secretsmanagerv2.RotationPolicyIntf) (map[string]interface{}, error) {
	if _, ok := model.(*secretsmanagerv2.CommonRotationPolicy); ok {
		return resourceIbmSmUsernamePasswordSecretCommonRotationPolicyToMap(model.(*secretsmanagerv2.CommonRotationPolicy))
	} else if _, ok := model.(*secretsmanagerv2.RotationPolicy); ok {
		modelMap := make(map[string]interface{})
		model := model.(*secretsmanagerv2.RotationPolicy)
		if model.AutoRotate != nil {
			modelMap["auto_rotate"] = model.AutoRotate
		}
		if model.Interval != nil {
			modelMap["interval"] = flex.IntValue(model.Interval)
		}
		if model.Unit != nil {
			modelMap["unit"] = model.Unit
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized secretsmanagerv2.RotationPolicyIntf subtype encountered")
	}
}

func resourceIbmSmUsernamePasswordSecretCommonRotationPolicyToMap(model *secretsmanagerv2.CommonRotationPolicy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["auto_rotate"] = model.AutoRotate
	if model.Interval != nil {
		modelMap["interval"] = flex.IntValue(model.Interval)
	}
	if model.Unit != nil {
		modelMap["unit"] = model.Unit
	}
	return modelMap, nil
}

func rotationAttributesDiffSuppress(key, oldValue, newValue string, d *schema.ResourceData) bool {
	autoRotate := d.Get("rotation").([]interface{})[0].(map[string]interface{})["auto_rotate"].(bool)
	if autoRotate == false {
		return true
	}

	return oldValue == newValue
}

func passwordGenerationPolicyToMap(model *secretsmanagerv2.PasswordGenerationPolicyRO) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Length != nil {
		modelMap["length"] = model.Length
	}
	if model.IncludeDigits != nil {
		modelMap["include_digits"] = model.IncludeDigits
	}
	if model.IncludeSymbols != nil {
		modelMap["include_symbols"] = model.IncludeSymbols
	}
	if model.IncludeUppercase != nil {
		modelMap["include_uppercase"] = model.IncludeUppercase
	}
	return modelMap, nil
}

func mapToPasswordGenerationPolicy(modelMap map[string]interface{}) (*secretsmanagerv2.PasswordGenerationPolicy, error) {
	model := &secretsmanagerv2.PasswordGenerationPolicy{}
	if modelMap["length"].(int) == 0 {
		model.Length = nil
	} else {
		model.Length = core.Int64Ptr(int64(modelMap["length"].(int)))
	}
	if modelMap["include_digits"] != nil {
		model.IncludeDigits = core.BoolPtr(modelMap["include_digits"].(bool))
	}
	if modelMap["include_symbols"] != nil {
		model.IncludeSymbols = core.BoolPtr(modelMap["include_symbols"].(bool))
	}
	if modelMap["include_uppercase"] != nil {
		model.IncludeUppercase = core.BoolPtr(modelMap["include_uppercase"].(bool))
	}
	return model, nil
}

func mapToPasswordGenerationPolicyPatch(modelMap map[string]interface{}) (*secretsmanagerv2.PasswordGenerationPolicyPatch, error) {
	model := &secretsmanagerv2.PasswordGenerationPolicyPatch{}
	if modelMap["length"].(int) == 0 {
		model.Length = nil
	} else {
		model.Length = core.Int64Ptr(int64(modelMap["length"].(int)))
	}
	if modelMap["include_digits"] != nil {
		model.IncludeDigits = core.BoolPtr(modelMap["include_digits"].(bool))
	}
	if modelMap["include_symbols"] != nil {
		model.IncludeSymbols = core.BoolPtr(modelMap["include_symbols"].(bool))
	}
	if modelMap["include_uppercase"] != nil {
		model.IncludeUppercase = core.BoolPtr(modelMap["include_uppercase"].(bool))
	}
	return model, nil
}
