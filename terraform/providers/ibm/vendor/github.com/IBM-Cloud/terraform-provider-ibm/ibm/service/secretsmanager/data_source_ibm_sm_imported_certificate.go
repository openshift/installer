// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package secretsmanager

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/secrets-manager-go-sdk/v2/secretsmanagerv2"
)

func DataSourceIbmSmImportedCertificate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSmImportedCertificateRead,

		Schema: map[string]*schema.Schema{
			"secret_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"secret_id", "name"},
				Description:  "The ID of the secret.",
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
			"custom_metadata": &schema.Schema{
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The secret metadata that a user can customize.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An extended description of your secret.To protect your privacy, do not use personal data, such as your name or location, as a description for your secret group.",
			},
			"downloaded": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the secret data that is associated with a secret version was retrieved in a call to the service API.",
			},
			"labels": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Labels that you can use to search for secrets in your instance.Up to 30 labels can be created.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"locks_total": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of locks of the secret.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"secret_id", "name"},
				RequiredWith: []string{"secret_group_name"},
				Description:  "The human-readable name of your secret.",
			},
			"secret_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A v4 UUID identifier, or `default` secret group.",
			},
			"secret_group_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"name"},
				Description:  "The human-readable name of your secret group.",
			},
			"secret_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The secret type. Supported types are arbitrary, certificates (imported, public, and private), IAM credentials, key-value, and user credentials.",
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
			"common_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Common Name (AKA CN) represents the server name protected by the SSL certificate.",
			},
			"expiration_date": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date a secret is expired. The date format follows RFC 3339.",
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
			"key_algorithm": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The identifier for the cryptographic algorithm used to generate the public key that is associated with the certificate.",
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
			"certificate": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The PEM-encoded contents of your certificate.",
			},
			"intermediate": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "(Optional) The PEM-encoded intermediate certificate to associate with the root certificate.",
			},
			"private_key": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "(Optional) The PEM-encoded private key to associate with the certificate.",
			},
		},
	}
}

func dataSourceIbmSmImportedCertificateRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	importedCertificateIntf, region, instanceId, diagError := getSecretByIdOrByName(context, d, meta, ImportedCertSecretType, ImportedCertSecretResourceName)
	if diagError != nil {
		return diagError
	}

	importedCertificate := importedCertificateIntf.(*secretsmanagerv2.ImportedCertificate)
	d.SetId(fmt.Sprintf("%s/%s/%s", region, instanceId, *importedCertificate.ID))

	var err error
	if err = d.Set("region", region); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting region"), fmt.Sprintf("(Data) %s", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("created_by", importedCertificate.CreatedBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_by"), fmt.Sprintf("(Data) %s", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("created_at", DateTimeToRFC3339(importedCertificate.CreatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_at"), fmt.Sprintf("(Data) %s", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("crn", importedCertificate.Crn); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting crn"), fmt.Sprintf("(Data) %s", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if importedCertificate.CustomMetadata != nil {
		convertedMap := make(map[string]interface{}, len(importedCertificate.CustomMetadata))
		for k, v := range importedCertificate.CustomMetadata {
			convertedMap[k] = v
		}

		if err = d.Set("custom_metadata", flex.Flatten(convertedMap)); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting custom_metadata"), fmt.Sprintf("(Data) %s", ImportedCertSecretResourceName), "read")
			return tfErr.GetDiag()
		}
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting custom_metadata"), fmt.Sprintf("(Data) %s", ImportedCertSecretResourceName), "read")
			return tfErr.GetDiag()
		}
	}

	if err = d.Set("description", importedCertificate.Description); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting description"), fmt.Sprintf("(Data) %s", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("downloaded", importedCertificate.Downloaded); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting downloaded"), fmt.Sprintf("(Data) %s", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("locks_total", flex.IntValue(importedCertificate.LocksTotal)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting locks_total"), fmt.Sprintf("(Data) %s", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("name", importedCertificate.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name"), fmt.Sprintf("(Data) %s", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("secret_group_id", importedCertificate.SecretGroupID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_group_id"), fmt.Sprintf("(Data) %s", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("secret_type", importedCertificate.SecretType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_type"), fmt.Sprintf("(Data) %s", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("state", flex.IntValue(importedCertificate.State)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting state"), fmt.Sprintf("(Data) %s", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("state_description", importedCertificate.StateDescription); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting state_description"), fmt.Sprintf("(Data) %s", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("updated_at", DateTimeToRFC3339(importedCertificate.UpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated_at"), fmt.Sprintf("(Data) %s", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("versions_total", flex.IntValue(importedCertificate.VersionsTotal)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting versions_total"), fmt.Sprintf("(Data) %s", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("signing_algorithm", importedCertificate.SigningAlgorithm); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting signing_algorithm"), fmt.Sprintf("(Data) %s", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("common_name", importedCertificate.CommonName); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting common_name"), fmt.Sprintf("(Data) %s", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("expiration_date", DateTimeToRFC3339(importedCertificate.ExpirationDate)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting expiration_date"), fmt.Sprintf("(Data) %s", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("intermediate_included", importedCertificate.IntermediateIncluded); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting intermediate_included"), fmt.Sprintf("(Data) %s", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("issuer", importedCertificate.Issuer); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting issuer"), fmt.Sprintf("(Data) %s", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("key_algorithm", importedCertificate.KeyAlgorithm); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting key_algorithm"), fmt.Sprintf("(Data) %s", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("private_key_included", importedCertificate.PrivateKeyIncluded); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting private_key_included"), fmt.Sprintf("(Data) %s", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("serial_number", importedCertificate.SerialNumber); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting serial_number"), fmt.Sprintf("(Data) %s", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	validity := []map[string]interface{}{}
	if importedCertificate.Validity != nil {
		modelMap, err := dataSourceIbmSmImportedCertificateCertificateValidityToMap(importedCertificate.Validity)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, "", fmt.Sprintf("(Data) %s", ImportedCertSecretResourceName), "read")
			return tfErr.GetDiag()
		}
		validity = append(validity, modelMap)
	}
	if err = d.Set("validity", validity); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting validity"), fmt.Sprintf("(Data) %s", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("certificate", importedCertificate.Certificate); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting certificate"), fmt.Sprintf("(Data) %s", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("intermediate", importedCertificate.Intermediate); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting intermediate"), fmt.Sprintf("(Data) %s", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("private_key", importedCertificate.PrivateKey); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting private_key"), fmt.Sprintf("(Data) %s", ImportedCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	return nil
}

func dataSourceIbmSmImportedCertificateCertificateValidityToMap(model *secretsmanagerv2.CertificateValidity) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.NotBefore != nil {
		modelMap["not_before"] = model.NotBefore.String()
	}
	if model.NotAfter != nil {
		modelMap["not_after"] = model.NotAfter.String()
	}
	return modelMap, nil
}
