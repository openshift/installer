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

func DataSourceIbmSmIamCredentialsConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSmIamCredentialsConfigurationRead,

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
			"api_key": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "An IBM Cloud API key that can create and manage service IDs. The API key must be assigned the Editor platform role on the Access Groups Service and the Operator platform role on the IAM Identity Service. For more information, see the [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-configure-iam-engine).",
			},
			"disabled": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "This attribute indicates whether the API key configuration is disabled. If it is set to `true`, the IAM credentials engine doesn't use the configured API key for credentials management.",
			},
		},
	}
}

func dataSourceIbmSmIamCredentialsConfigurationRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", fmt.Sprintf("(Data) %s", IAMCredentialsConfigResourceName), "read")
		return tfErr.GetDiag()
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	getConfigurationOptions := &secretsmanagerv2.GetConfigurationOptions{}

	getConfigurationOptions.SetName(d.Get("name").(string))

	iAMCredentialsConfigurationIntf, response, err := secretsManagerClient.GetConfigurationWithContext(context, getConfigurationOptions)
	if err != nil {
		log.Printf("[DEBUG] GetConfigurationWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetConfigurationWithContext failed %s\n%s", err, response), fmt.Sprintf("(Data) %s", IAMCredentialsConfigResourceName), "read")
		return tfErr.GetDiag()
	}
	iAMCredentialsConfiguration := iAMCredentialsConfigurationIntf.(*secretsmanagerv2.IAMCredentialsConfiguration)

	d.SetId(fmt.Sprintf("%s/%s/%s", region, instanceId, *getConfigurationOptions.Name))

	if err = d.Set("region", region); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting region"), fmt.Sprintf("(Data) %s", IAMCredentialsConfigResourceName), "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("config_type", iAMCredentialsConfiguration.ConfigType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting config_type"), fmt.Sprintf("(Data) %s", IAMCredentialsConfigResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("secret_type", iAMCredentialsConfiguration.SecretType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_type"), fmt.Sprintf("(Data) %s", IAMCredentialsConfigResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("created_by", iAMCredentialsConfiguration.CreatedBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_by"), fmt.Sprintf("(Data) %s", IAMCredentialsConfigResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("created_at", DateTimeToRFC3339(iAMCredentialsConfiguration.CreatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_at"), fmt.Sprintf("(Data) %s", IAMCredentialsConfigResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("updated_at", DateTimeToRFC3339(iAMCredentialsConfiguration.UpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated_at"), fmt.Sprintf("(Data) %s", IAMCredentialsConfigResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("api_key", iAMCredentialsConfiguration.ApiKey); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting api_key"), fmt.Sprintf("(Data) %s", IAMCredentialsConfigResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("disabled", iAMCredentialsConfiguration.Disabled); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting disabled"), fmt.Sprintf("(Data) %s", IAMCredentialsConfigResourceName), "read")
		return tfErr.GetDiag()
	}

	return nil
}
