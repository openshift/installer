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

func DataSourceIbmSmConfigurationPublicCertificateDNSCis() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSmConfigurationPublicCertificateDNSCisRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the configuration.",
			},
			"config_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The configuration type.",
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
			"cloud_internet_services_apikey": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An IBM Cloud API key that can to list domains in your Cloud Internet Services instance.To grant Secrets Manager the ability to view the Cloud Internet Services instance and all of its domains, the API key must be assigned the Reader service role on Internet Services (`internet-svcs`).If you need to manage specific domains, you can assign the Manager role. For production environments, it is recommended that you assign the Reader access role, and then use the[IAM Policy Management API](https://cloud.ibm.com/apidocs/iam-policy-management#create-policy) to control specific domains. For more information, see the [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-prepare-order-certificates#authorize-specific-domains).",
			},
			"cloud_internet_services_crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A CRN that uniquely identifies an IBM Cloud resource.",
			},
		},
	}
}

func dataSourceIbmSmConfigurationPublicCertificateDNSCisRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", fmt.Sprintf("(Data) %s", PublicCertConfigDnsCISResourceName), "read")
		return tfErr.GetDiag()
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	getConfigurationOptions := &secretsmanagerv2.GetConfigurationOptions{}

	getConfigurationOptions.SetName(d.Get("name").(string))

	publicCertificateConfigurationDNSCloudInternetServicesIntf, response, err := secretsManagerClient.GetConfigurationWithContext(context, getConfigurationOptions)
	if err != nil {
		log.Printf("[DEBUG] GetConfigurationWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetConfigurationWithContext failed %s\n%s", err, response), fmt.Sprintf("(Data) %s", PublicCertConfigDnsCISResourceName), "read")
		return tfErr.GetDiag()
	}

	publicCertificateConfigurationDNSCloudInternetServices := publicCertificateConfigurationDNSCloudInternetServicesIntf.(*secretsmanagerv2.PublicCertificateConfigurationDNSCloudInternetServices)

	d.SetId(fmt.Sprintf("%s/%s/%s", region, instanceId, *getConfigurationOptions.Name))

	if err = d.Set("region", region); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting region"), fmt.Sprintf("(Data) %s", PublicCertConfigDnsCISResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("config_type", publicCertificateConfigurationDNSCloudInternetServices.ConfigType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting config_type"), fmt.Sprintf("(Data) %s", PublicCertConfigDnsCISResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("secret_type", publicCertificateConfigurationDNSCloudInternetServices.SecretType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_type"), fmt.Sprintf("(Data) %s", PublicCertConfigDnsCISResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("created_by", publicCertificateConfigurationDNSCloudInternetServices.CreatedBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_by"), fmt.Sprintf("(Data) %s", PublicCertConfigDnsCISResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("created_at", DateTimeToRFC3339(publicCertificateConfigurationDNSCloudInternetServices.CreatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_at"), fmt.Sprintf("(Data) %s", PublicCertConfigDnsCISResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("updated_at", DateTimeToRFC3339(publicCertificateConfigurationDNSCloudInternetServices.UpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated_at"), fmt.Sprintf("(Data) %s", PublicCertConfigDnsCISResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("cloud_internet_services_apikey", publicCertificateConfigurationDNSCloudInternetServices.CloudInternetServicesApikey); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting cloud_internet_services_apikey"), fmt.Sprintf("(Data) %s", PublicCertConfigDnsCISResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("cloud_internet_services_crn", publicCertificateConfigurationDNSCloudInternetServices.CloudInternetServicesCrn); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting cloud_internet_services_crn"), fmt.Sprintf("(Data) %s", PublicCertConfigDnsCISResourceName), "read")
		return tfErr.GetDiag()
	}

	return nil
}
