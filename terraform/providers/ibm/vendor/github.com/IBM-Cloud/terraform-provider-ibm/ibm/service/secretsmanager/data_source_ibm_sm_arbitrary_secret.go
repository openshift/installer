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

func DataSourceIbmSmArbitrarySecret() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSmArbitrarySecretRead,

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
			"expiration_date": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date a secret is expired. The date format follows RFC 3339.",
			},
			"payload": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The arbitrary secret's data payload.",
			},
		},
	}
}

func dataSourceIbmSmArbitrarySecretRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secret, region, instanceId, diagError := getSecretByIdOrByName(context, d, meta, ArbitrarySecretType, ArbitrarySecretResourceName)
	if diagError != nil {
		return diagError
	}

	arbitrarySecret := secret.(*secretsmanagerv2.ArbitrarySecret)
	d.SetId(fmt.Sprintf("%s/%s/%s", region, instanceId, *arbitrarySecret.ID))

	var err error
	if err = d.Set("region", region); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting region"), fmt.Sprintf("(Data) %s", ArbitrarySecretResourceName), "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("created_by", arbitrarySecret.CreatedBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_by"), fmt.Sprintf("(Data) %s", ArbitrarySecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("created_at", DateTimeToRFC3339(arbitrarySecret.CreatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_at"), fmt.Sprintf("(Data) %s", ArbitrarySecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("crn", arbitrarySecret.Crn); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting crn"), fmt.Sprintf("(Data) %s", ArbitrarySecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if arbitrarySecret.CustomMetadata != nil {
		convertedMap := make(map[string]interface{}, len(arbitrarySecret.CustomMetadata))
		for k, v := range arbitrarySecret.CustomMetadata {
			convertedMap[k] = v
		}

		if err = d.Set("custom_metadata", flex.Flatten(convertedMap)); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting custom_metadata"), fmt.Sprintf("(Data) %s", ArbitrarySecretResourceName), "read")
			return tfErr.GetDiag()
		}
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting custom_metadata"), fmt.Sprintf("(Data) %s", ArbitrarySecretResourceName), "read")
			return tfErr.GetDiag()
		}
	}

	if err = d.Set("description", arbitrarySecret.Description); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting description"), fmt.Sprintf("(Data) %s", ArbitrarySecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("downloaded", arbitrarySecret.Downloaded); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting downloaded"), fmt.Sprintf("(Data) %s", ArbitrarySecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("locks_total", flex.IntValue(arbitrarySecret.LocksTotal)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting locks_total"), fmt.Sprintf("(Data) %s", ArbitrarySecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("name", arbitrarySecret.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name"), fmt.Sprintf("(Data) %s", ArbitrarySecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("secret_group_id", arbitrarySecret.SecretGroupID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_group_id"), fmt.Sprintf("(Data) %s", ArbitrarySecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("secret_type", arbitrarySecret.SecretType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_type"), fmt.Sprintf("(Data) %s", ArbitrarySecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("state", flex.IntValue(arbitrarySecret.State)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting state"), fmt.Sprintf("(Data) %s", ArbitrarySecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("state_description", arbitrarySecret.StateDescription); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting state_description"), fmt.Sprintf("(Data) %s", ArbitrarySecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("updated_at", DateTimeToRFC3339(arbitrarySecret.UpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated_at"), fmt.Sprintf("(Data) %s", ArbitrarySecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("versions_total", flex.IntValue(arbitrarySecret.VersionsTotal)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting versions_total"), fmt.Sprintf("(Data) %s", ArbitrarySecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("expiration_date", DateTimeToRFC3339(arbitrarySecret.ExpirationDate)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting expiration_date"), fmt.Sprintf("(Data) %s", ArbitrarySecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("payload", arbitrarySecret.Payload); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting payload"), fmt.Sprintf("(Data) %s", ArbitrarySecretResourceName), "read")
		return tfErr.GetDiag()
	}

	return nil
}
