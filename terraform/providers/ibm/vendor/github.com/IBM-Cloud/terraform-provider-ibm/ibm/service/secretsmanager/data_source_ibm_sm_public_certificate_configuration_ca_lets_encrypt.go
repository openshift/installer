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
	"github.com/IBM/secrets-manager-go-sdk/v2/secretsmanagerv2"
)

func DataSourceIbmSmPublicCertificateConfigurationCALetsEncrypt() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSmPublicCertificateConfigurationCALetsEncryptRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the configuration.",
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
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when a resource was recently modified. The date format follows RFC 3339.",
			},
			"lets_encrypt_environment": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The configuration of the Let's Encrypt CA environment.",
			},
			"lets_encrypt_private_key": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The PEM encoded private key of your Lets Encrypt account.",
			},
			"lets_encrypt_preferred_chain": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Prefer the chain with an issuer matching this Subject Common Name.",
			},
		},
	}
}

func dataSourceIbmSmPublicCertificateConfigurationCALetsEncryptRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		return diag.FromErr(err)
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	getConfigurationOptions := &secretsmanagerv2.GetConfigurationOptions{}

	getConfigurationOptions.SetName(d.Get("name").(string))

	publicCertificateConfigurationCALetsEncryptIntf, response, err := secretsManagerClient.GetConfigurationWithContext(context, getConfigurationOptions)
	if err != nil {
		log.Printf("[DEBUG] GetConfigurationWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetConfigurationWithContext failed %s\n%s", err, response))
	}
	publicCertificateConfigurationCALetsEncrypt := publicCertificateConfigurationCALetsEncryptIntf.(*secretsmanagerv2.PublicCertificateConfigurationCALetsEncrypt)

	d.SetId(fmt.Sprintf("%s/%s/%s", region, instanceId, *getConfigurationOptions.Name))

	if err = d.Set("region", region); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting region: %s", err))
	}

	if err = d.Set("created_by", publicCertificateConfigurationCALetsEncrypt.CreatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_by: %s", err))
	}

	if err = d.Set("created_at", DateTimeToRFC3339(publicCertificateConfigurationCALetsEncrypt.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}

	if err = d.Set("updated_at", DateTimeToRFC3339(publicCertificateConfigurationCALetsEncrypt.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
	}

	if err = d.Set("lets_encrypt_environment", publicCertificateConfigurationCALetsEncrypt.LetsEncryptEnvironment); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lets_encrypt_environment: %s", err))
	}

	if err = d.Set("lets_encrypt_private_key", publicCertificateConfigurationCALetsEncrypt.LetsEncryptPrivateKey); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lets_encrypt_private_key: %s", err))
	}

	if err = d.Set("lets_encrypt_preferred_chain", publicCertificateConfigurationCALetsEncrypt.LetsEncryptPreferredChain); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lets_encrypt_preferred_chain: %s", err))
	}

	return nil
}
