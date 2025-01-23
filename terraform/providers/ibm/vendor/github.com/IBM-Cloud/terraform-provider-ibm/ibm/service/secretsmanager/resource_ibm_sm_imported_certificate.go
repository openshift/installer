// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package secretsmanager

import (
	"context"
	"fmt"
	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/secrets-manager-go-sdk/v2/secretsmanagerv2"
)

func ResourceIbmSmImportedCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmSmImportedCertificateCreate,
		ReadContext:   resourceIbmSmImportedCertificateRead,
		UpdateContext: resourceIbmSmImportedCertificateUpdate,
		DeleteContext: resourceIbmSmImportedCertificateDelete,
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
				Computed:    true,
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
			"secret_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A v4 UUID identifier.",
			},
			"secret_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "A v4 UUID identifier, or `default` secret group.",
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
			"certificate": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					if removeNewLineFromCertificate(oldValue) == removeNewLineFromCertificate(newValue) {
						return true
					}
					return false
				},
				Description: "The PEM-encoded contents of your certificate.",
			},
			"intermediate": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					if removeNewLineFromCertificate(oldValue) == removeNewLineFromCertificate(newValue) {
						return true
					}
					return false
				},
				Description: "(Optional) The PEM-encoded intermediate certificate to associate with the root certificate.",
			},
			"private_key": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					if removeNewLineFromCertificate(oldValue) == removeNewLineFromCertificate(newValue) {
						return true
					}
					return false
				},
				Description: "(Optional) The PEM-encoded private key to associate with the certificate.",
			},
			"common_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Common Name (AKA CN) represents the server name that is protected by the SSL certificate.",
			},
			"key_algorithm": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The identifier for the cryptographic algorithm to be used to generate the public key that is associated with the certificate.The algorithm that you select determines the encryption algorithm (`RSA` or `ECDSA`) and key size to be used to generate keys and sign certificates. For longer living certificates, it is recommended to use longer keys to provide more encryption protection. Allowed values:  RSA2048, RSA4096, EC256, EC384.",
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
			"signing_algorithm": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The identifier for the cryptographic algorithm that was used by the issuing certificate authority to sign a certificate.",
			},
			"intermediate_included": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the certificate was imported with an associated intermediate certificate.",
			},
			"issuer": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The distinguished name that identifies the entity that signed and issued the certificate.",
			},
			"private_key_included": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the certificate was imported with an associated private key.",
			},
			"serial_number": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique serial number that was assigned to a certificate by the issuing certificate authority.",
			},
			"validity": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The date and time that the certificate validity period begins and ends.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"not_before": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date-time format follows RFC 3339.",
						},
						"not_after": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date-time format follows RFC 3339.",
						},
					},
				},
			},
		},
	}
}

func resourceIbmSmImportedCertificateCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", ImportedCertSecretResourceName, "create")
		return tfErr.GetDiag()
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	createSecretOptions := &secretsmanagerv2.CreateSecretOptions{}

	secretPrototypeModel, err := resourceIbmSmImportedCertificateMapToSecretPrototype(d)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", ImportedCertSecretResourceName, "create")
		return tfErr.GetDiag()
	}
	createSecretOptions.SetSecretPrototype(secretPrototypeModel)

	secretIntf, response, err := secretsManagerClient.CreateSecretWithContext(context, createSecretOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateSecretWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateSecretWithContext failed: %s\n%s", err.Error(), response), ImportedCertSecretResourceName, "create")
		return tfErr.GetDiag()
	}
	secret := secretIntf.(*secretsmanagerv2.ImportedCertificate)

	d.SetId(fmt.Sprintf("%s/%s/%s", region, instanceId, *secret.ID))
	d.Set("secret_id", *secret.ID)

	_, err = waitForIbmSmImportedCertificateCreate(secretsManagerClient, d)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error waiting for resource IbmSmImportedCertificate (%s) to be created: %s", d.Id(), err.Error()), ImportedCertSecretResourceName, "create")
		return tfErr.GetDiag()
	}

	return resourceIbmSmImportedCertificateRead(context, d, meta)
}

func waitForIbmSmImportedCertificateCreate(secretsManagerClient *secretsmanagerv2.SecretsManagerV2, d *schema.ResourceData) (interface{}, error) {
	getSecretOptions := &secretsmanagerv2.GetSecretOptions{}

	id := strings.Split(d.Id(), "/")
	secretId := id[2]

	getSecretOptions.SetID(secretId)

	stateConf := &resource.StateChangeConf{
		Pending: []string{"pre_activation"},
		Target:  []string{"active"},
		Refresh: func() (interface{}, string, error) {
			stateObjIntf, response, err := secretsManagerClient.GetSecret(getSecretOptions)
			stateObj := stateObjIntf.(*secretsmanagerv2.ImportedCertificate)
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

func resourceIbmSmImportedCertificateRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", ImportedCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}

	id := strings.Split(d.Id(), "/")
	if len(id) != 3 {
		tfErr := flex.TerraformErrorf(nil, "Wrong format of resource ID. To import a secret use the format `<region>/<instance_id>/<secret_id>`", ImportedCertSecretResourceName, "read")
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
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetSecretWithContext failed %s\n%s", err, response), ImportedCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	secret := secretIntf.(*secretsmanagerv2.ImportedCertificate)

	if err = d.Set("secret_id", secretId); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_id"), ImportedCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("instance_id", instanceId); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting instance_id"), ImportedCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("region", region); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting region"), ImportedCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("created_by", secret.CreatedBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_by"), ImportedCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("created_at", DateTimeToRFC3339(secret.CreatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_at"), ImportedCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("crn", secret.Crn); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting crn"), ImportedCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if secret.CustomMetadata != nil {
		d.Set("custom_metadata", secret.CustomMetadata)
	}
	if err = d.Set("description", secret.Description); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting description"), ImportedCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("downloaded", secret.Downloaded); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting downloaded"), ImportedCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if secret.Labels != nil {
		if err = d.Set("labels", secret.Labels); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting labels"), ImportedCertSecretResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("locks_total", flex.IntValue(secret.LocksTotal)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting locks_total"), ImportedCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("name", secret.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name"), ImportedCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("secret_group_id", secret.SecretGroupID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_group_id"), ImportedCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("secret_type", secret.SecretType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_type"), ImportedCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("state", flex.IntValue(secret.State)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting state"), ImportedCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("state_description", secret.StateDescription); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting state_description"), ImportedCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("updated_at", DateTimeToRFC3339(secret.UpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated_at"), ImportedCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("versions_total", flex.IntValue(secret.VersionsTotal)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting versions_total"), ImportedCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("signing_algorithm", secret.SigningAlgorithm); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting signing_algorithm"), ImportedCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("common_name", secret.CommonName); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting common_name"), ImportedCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("expiration_date", DateTimeToRFC3339(secret.ExpirationDate)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting expiration_date"), ImportedCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("intermediate_included", secret.IntermediateIncluded); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting intermediate_included"), ImportedCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("issuer", secret.Issuer); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting issuer"), ImportedCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("key_algorithm", secret.KeyAlgorithm); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting key_algorithm"), ImportedCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("private_key_included", secret.PrivateKeyIncluded); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting private_key_included"), ImportedCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("serial_number", secret.SerialNumber); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting serial_number"), ImportedCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	validityMap, err := resourceIbmSmImportedCertificateCertificateValidityToMap(secret.Validity)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", ImportedCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("validity", []map[string]interface{}{validityMap}); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting validity"), ImportedCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("certificate", secret.Certificate); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting certificate"), ImportedCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("intermediate", secret.Intermediate); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting intermediate"), ImportedCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("private_key", secret.PrivateKey); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting private_key"), ImportedCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}

	// Call get version metadata API to get the current version_custom_metadata
	getVersionMetdataOptions := &secretsmanagerv2.GetSecretVersionMetadataOptions{}
	getVersionMetdataOptions.SetSecretID(secretId)
	getVersionMetdataOptions.SetID("current")

	versionMetadataIntf, response, err := secretsManagerClient.GetSecretVersionMetadataWithContext(context, getVersionMetdataOptions)
	if err != nil {
		log.Printf("[DEBUG] GetSecretVersionMetadataWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetSecretVersionMetadataWithContext failed %s\n%s", err, response), ImportedCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}

	versionMetadata := versionMetadataIntf.(*secretsmanagerv2.ImportedCertificateVersionMetadata)
	if versionMetadata.VersionCustomMetadata != nil {
		if err = d.Set("version_custom_metadata", versionMetadata.VersionCustomMetadata); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting version_custom_metadata"), ImportedCertSecretResourceName, "read")
			return tfErr.GetDiag()
		}
	}

	return nil
}

func resourceIbmSmImportedCertificateUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", ImportedCertSecretResourceName, "update")
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

	patchVals := &secretsmanagerv2.ImportedCertificateMetadataPatch{}

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

	// Apply change in metadata (if changed)
	if hasChange {
		updateSecretMetadataOptions.SecretMetadataPatch, _ = patchVals.AsPatch()
		_, response, err := secretsManagerClient.UpdateSecretMetadataWithContext(context, updateSecretMetadataOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateSecretMetadataWithContext failed %s\n%s", err, response)
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateSecretMetadataWithContext failed %s\n%s", err, response), ImportedCertSecretResourceName, "update")
			return tfErr.GetDiag()
		}
	}

	// Apply change in secret data (if changed)
	if d.HasChange("certificate") || d.HasChange("intermediate") || d.HasChange("private_key") {
		versionModel := &secretsmanagerv2.ImportedCertificateVersionPrototype{}
		versionModel.Certificate = core.StringPtr(d.Get("certificate").(string))
		if _, ok := d.GetOk("intermediate"); ok {
			versionModel.Intermediate = core.StringPtr(formatCertificate(d.Get("intermediate").(string)))
		}
		if _, ok := d.GetOk("private_key"); ok {
			versionModel.PrivateKey = core.StringPtr(formatCertificate(d.Get("private_key").(string)))
		}
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
				resourceIbmSmImportedCertificateRead(context, d, meta)
			}
			log.Printf("[DEBUG] CreateSecretVersionWithContext failed %s\n%s", err, response)
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateSecretVersionWithContext failed %s\n%s", err, response), ImportedCertSecretResourceName, "update")
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
				resourceIbmSmImportedCertificateRead(context, d, meta)
			}
			log.Printf("[DEBUG] UpdateSecretVersionMetadataWithContext failed %s\n%s", err, response)
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateSecretVersionMetadataWithContext failed %s\n%s", err, response), ImportedCertSecretResourceName, "update")
			return tfErr.GetDiag()
		}
	}
	return resourceIbmSmImportedCertificateRead(context, d, meta)
}

func resourceIbmSmImportedCertificateDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", ImportedCertSecretResourceName, "delete")
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
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteSecretWithContext failed %s\n%s", err, response), ImportedCertSecretResourceName, "delete")
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func resourceIbmSmImportedCertificateMapToSecretPrototype(d *schema.ResourceData) (secretsmanagerv2.SecretPrototypeIntf, error) {
	model := &secretsmanagerv2.ImportedCertificatePrototype{}
	model.SecretType = core.StringPtr("imported_cert")

	if _, ok := d.GetOk("name"); ok {
		model.Name = core.StringPtr(d.Get("name").(string))
	}
	if _, ok := d.GetOk("custom_metadata"); ok {
		model.CustomMetadata = d.Get("custom_metadata").(map[string]interface{})
	}
	if _, ok := d.GetOk("description"); ok {
		model.Description = core.StringPtr(d.Get("description").(string))
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
	if _, ok := d.GetOk("certificate"); ok {
		model.Certificate = core.StringPtr(formatCertificate(d.Get("certificate").(string)))
	}

	if _, ok := d.GetOk("intermediate"); ok {
		model.Intermediate = core.StringPtr(formatCertificate(d.Get("intermediate").(string)))
	}

	if _, ok := d.GetOk("private_key"); ok {
		model.PrivateKey = core.StringPtr(formatCertificate(d.Get("private_key").(string)))
	}

	return model, nil
}

func resourceIbmSmImportedCertificateImportedCertificatePrototypeToMap(model *secretsmanagerv2.ImportedCertificatePrototype) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["secret_type"] = model.SecretType
	modelMap["name"] = model.Name
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.SecretGroupID != nil {
		modelMap["secret_group_id"] = model.SecretGroupID
	}
	if model.Labels != nil {
		modelMap["labels"] = model.Labels
	}
	modelMap["certificate"] = model.Certificate
	if model.Intermediate != nil {
		modelMap["intermediate"] = model.Intermediate
	}
	if model.PrivateKey != nil {
		modelMap["private_key"] = model.PrivateKey
	}
	if model.CustomMetadata != nil {
		// TODO: handle CustomMetadata of type TypeMap -- container, not list
	}
	if model.VersionCustomMetadata != nil {
		// TODO: handle VersionCustomMetadata of type TypeMap -- container, not list
	}
	return modelMap, nil
}

func resourceIbmSmImportedCertificateCertificateValidityToMap(model *secretsmanagerv2.CertificateValidity) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["not_before"] = model.NotBefore.String()
	modelMap["not_after"] = model.NotAfter.String()
	return modelMap, nil
}

func removeNewLineFromCertificate(originalCert string) string {
	if originalCert == "" {
		return originalCert
	}
	noR := strings.ReplaceAll(originalCert, "\r", "")
	noNnoR := strings.ReplaceAll(noR, "\n", "")
	return noNnoR
}

func formatCertificate(originalCert string) string {
	if originalCert == "" {
		return originalCert
	}
	noR := strings.ReplaceAll(originalCert, "\r", "")
	noNnoR := strings.SplitN(noR, "\n", -1)
	certParsed := ""
	i := 0
	for i < len(noNnoR) {
		certParsed += noNnoR[i]
		if i < len(noNnoR)-1 && len(noNnoR[i+1]) > 0 {
			certParsed += "\r\n"
		} else {
			break
		}
		i++
	}
	return certParsed
}
