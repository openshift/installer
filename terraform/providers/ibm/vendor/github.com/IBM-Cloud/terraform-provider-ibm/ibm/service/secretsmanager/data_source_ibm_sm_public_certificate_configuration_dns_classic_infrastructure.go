// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package secretsmanager

import (
	"context"
	"fmt"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/secrets-manager-go-sdk/v2/secretsmanagerv2"
)

func DataSourceIbmSmPublicCertificateConfigurationDNSClassicInfrastructure() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSmPublicCertificateConfigurationDNSClassicInfrastructureRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the configuration.",
			},
			"config_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Th configuration type.",
			},
			"secret_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The secret type. Supported types are arbitrary, certificates (imported, public, and private), IAM credentials, key-value, and user credentials.",
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
			"classic_infrastructure_username": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The username that is associated with your classic infrastructure account.In most cases, your classic infrastructure username is your `<account_id>_<email_address>`. For more information, see the [docs](https://cloud.ibm.com/docs/account?topic=account-classic_keys).",
			},
			"classic_infrastructure_password": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Your classic infrastructure API key.For information about viewing and accessing your classic infrastructure API key, see the [docs](https://cloud.ibm.com/docs/account?topic=account-classic_keys).",
			},
		},
	}
}

func dataSourceIbmSmPublicCertificateConfigurationDNSClassicInfrastructureRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", fmt.Sprintf("(Data) %s", PublicCertConfigDnsClassicInfrastructureResourceName), "read")
		return tfErr.GetDiag()
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	getConfigurationOptions := &secretsmanagerv2.GetConfigurationOptions{}

	getConfigurationOptions.SetName(d.Get("name").(string))

	publicCertificateConfigurationDNSClassicInfrastructureInf, response, err := secretsManagerClient.GetConfigurationWithContext(context, getConfigurationOptions)
	if err != nil {
		log.Printf("[DEBUG] GetConfigurationWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetConfigurationWithContext failed %s\n%s", err, response), fmt.Sprintf("(Data) %s", PublicCertConfigDnsClassicInfrastructureResourceName), "read")
		return tfErr.GetDiag()
	}

	publicCertificateConfigurationDNSClassicInfrastructure := publicCertificateConfigurationDNSClassicInfrastructureInf.(*secretsmanagerv2.PublicCertificateConfigurationDNSClassicInfrastructure)
	d.SetId(fmt.Sprintf("%s/%s/%s", region, instanceId, *getConfigurationOptions.Name))

	if err = d.Set("region", region); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting region"), fmt.Sprintf("(Data) %s", PublicCertConfigDnsClassicInfrastructureResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("config_type", publicCertificateConfigurationDNSClassicInfrastructure.ConfigType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting config_type"), fmt.Sprintf("(Data) %s", PublicCertConfigDnsClassicInfrastructureResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("secret_type", publicCertificateConfigurationDNSClassicInfrastructure.SecretType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_type"), fmt.Sprintf("(Data) %s", PublicCertConfigDnsClassicInfrastructureResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("created_by", publicCertificateConfigurationDNSClassicInfrastructure.CreatedBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_by"), fmt.Sprintf("(Data) %s", PublicCertConfigDnsClassicInfrastructureResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("created_at", DateTimeToRFC3339(publicCertificateConfigurationDNSClassicInfrastructure.CreatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_at"), fmt.Sprintf("(Data) %s", PublicCertConfigDnsClassicInfrastructureResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("updated_at", DateTimeToRFC3339(publicCertificateConfigurationDNSClassicInfrastructure.UpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated_at"), fmt.Sprintf("(Data) %s", PublicCertConfigDnsClassicInfrastructureResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("classic_infrastructure_username", publicCertificateConfigurationDNSClassicInfrastructure.ClassicInfrastructureUsername); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting classic_infrastructure_username"), fmt.Sprintf("(Data) %s", PublicCertConfigDnsClassicInfrastructureResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("classic_infrastructure_password", publicCertificateConfigurationDNSClassicInfrastructure.ClassicInfrastructurePassword); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting classic_infrastructure_password"), fmt.Sprintf("(Data) %s", PublicCertConfigDnsClassicInfrastructureResourceName), "read")
		return tfErr.GetDiag()
	}

	return nil
}
