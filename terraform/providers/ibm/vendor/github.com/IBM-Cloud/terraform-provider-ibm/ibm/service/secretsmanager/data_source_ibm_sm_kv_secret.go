// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package secretsmanager

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/secrets-manager-go-sdk/v2/secretsmanagerv2"
)

func DataSourceIbmSmKvSecret() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSmKvSecretRead,

		Schema: map[string]*schema.Schema{
			"secret_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the secret.",
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
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The human-readable name of your secret.",
			},
			"secret_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A v4 UUID identifier, or `default` secret group.",
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
			"data": &schema.Schema{
				Type:        schema.TypeMap,
				Computed:    true,
				Sensitive:   true,
				Description: "The payload data of a key-value secret.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceIbmSmKvSecretRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		return diag.FromErr(err)
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	getSecretOptions := &secretsmanagerv2.GetSecretOptions{}

	secretId := d.Get("secret_id").(string)
	getSecretOptions.SetID(secretId)

	kVSecretIntf, response, err := secretsManagerClient.GetSecretWithContext(context, getSecretOptions)
	if err != nil {
		log.Printf("[DEBUG] GetSecretWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetSecretWithContext failed %s\n%s", err, response))
	}
	kVSecret := kVSecretIntf.(*secretsmanagerv2.KVSecret)

	d.SetId(fmt.Sprintf("%s/%s/%s", region, instanceId, secretId))

	if err = d.Set("region", region); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting region: %s", err))
	}

	if err = d.Set("created_by", kVSecret.CreatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_by: %s", err))
	}

	if err = d.Set("created_at", DateTimeToRFC3339(kVSecret.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}

	if err = d.Set("crn", kVSecret.Crn); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}

	if kVSecret.CustomMetadata != nil {
		convertedMap := make(map[string]interface{}, len(kVSecret.CustomMetadata))
		for k, v := range kVSecret.CustomMetadata {
			convertedMap[k] = v
		}

		if err = d.Set("custom_metadata", flex.Flatten(convertedMap)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting custom_metadata: %s", err))
		}
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting custom_metadata %s", err))
		}
	}

	if err = d.Set("description", kVSecret.Description); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
	}

	if err = d.Set("downloaded", kVSecret.Downloaded); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting downloaded: %s", err))
	}

	if err = d.Set("locks_total", flex.IntValue(kVSecret.LocksTotal)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting locks_total: %s", err))
	}

	if err = d.Set("name", kVSecret.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}

	if err = d.Set("secret_group_id", kVSecret.SecretGroupID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting secret_group_id: %s", err))
	}

	if err = d.Set("secret_type", kVSecret.SecretType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting secret_type: %s", err))
	}

	if err = d.Set("state", flex.IntValue(kVSecret.State)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting state: %s", err))
	}

	if err = d.Set("state_description", kVSecret.StateDescription); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting state_description: %s", err))
	}

	if err = d.Set("updated_at", DateTimeToRFC3339(kVSecret.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
	}

	if err = d.Set("versions_total", flex.IntValue(kVSecret.VersionsTotal)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting versions_total: %s", err))
	}

	if kVSecret.Data != nil {
		convertedMap := make(map[string]interface{}, len(kVSecret.Data))
		for k, v := range kVSecret.Data {
			convertedMap[k] = v
		}

		if err = d.Set("data", flex.Flatten(convertedMap)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting data: %s", err))
		}
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting data %s", err))
		}
	}

	return nil
}
