// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package secretsmanager

import (
	"context"
	"fmt"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/secrets-manager-go-sdk/v2/secretsmanagerv2"
)

func ResourceIbmSmConfigurationPublicCertificateDNSCis() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmSmConfigurationPublicCertificateDNSCisCreate,
		ReadContext:   resourceIbmSmConfigurationPublicCertificateDNSCisRead,
		UpdateContext: resourceIbmSmConfigurationPublicCertificateDNSCisUpdate,
		DeleteContext: resourceIbmSmConfigurationPublicCertificateDNSCisDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A human-readable unique name to assign to your configuration.To protect your privacy, do not use personal data, such as your name or location, as an name for your secret.",
			},
			"config_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The configuration type.",
			},
			"cloud_internet_services_apikey": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator(PublicCertConfigDnsCISResourceName, "cloud_internet_services_apikey"),
				Description:  "An IBM Cloud API key that can to list domains in your Cloud Internet Services instance.To grant Secrets Manager the ability to view the Cloud Internet Services instance and all of its domains, the API key must be assigned the Reader service role on Internet Services (`internet-svcs`).If you need to manage specific domains, you can assign the Manager role. For production environments, it is recommended that you assign the Reader access role, and then use the[IAM Policy Management API](https://cloud.ibm.com/apidocs/iam-policy-management#create-policy) to control specific domains. For more information, see the [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-prepare-order-certificates#authorize-specific-domains).",
			},
			"cloud_internet_services_crn": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator(PublicCertConfigDnsCISResourceName, "cloud_internet_services_crn"),
				Description:  "A CRN that uniquely identifies an IBM Cloud resource.",
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
		},
	}
}

func ResourceIbmSmConfigurationPublicCertificateDNSCisValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cloud_internet_services_apikey",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `(.*?)`,
			MinValueLength:             4,
			MaxValueLength:             128,
		},
		validate.ValidateSchema{
			Identifier:                 "cloud_internet_services_crn",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^crn:v[0-9](:([A-Za-z0-9-._~!$&'()*+,;=@\/]|%[0-9A-Z]{2})*){8}$`,
			MinValueLength:             9,
			MaxValueLength:             512,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: PublicCertConfigDnsCISResourceName, Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmSmConfigurationPublicCertificateDNSCisCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PublicCertConfigDnsCISResourceName, "create")
		return tfErr.GetDiag()
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))
	bodyModelMap := map[string]interface{}{}
	createConfigurationOptions := &secretsmanagerv2.CreateConfigurationOptions{}

	if _, ok := d.GetOk("config_type"); ok {
		bodyModelMap["config_type"] = d.Get("config_type")
	}

	if _, ok := d.GetOk("name"); ok {
		bodyModelMap["name"] = d.Get("name")
	}

	if _, ok := d.GetOk("cloud_internet_services_apikey"); ok {
		bodyModelMap["cloud_internet_services_apikey"] = d.Get("cloud_internet_services_apikey")
	}
	if _, ok := d.GetOk("cloud_internet_services_crn"); ok {
		bodyModelMap["cloud_internet_services_crn"] = d.Get("cloud_internet_services_crn")
	}
	convertedModel, err := resourceIbmSmConfigurationPublicCertificateCisMapToPublicCertificateConfigurationDNSCloudInternetServicesPrototype(bodyModelMap)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PublicCertConfigDnsCISResourceName, "create")
		return tfErr.GetDiag()
	}
	createConfigurationOptions.SetConfigurationPrototype(convertedModel)

	configurationIntf, response, err := secretsManagerClient.CreateConfigurationWithContext(context, createConfigurationOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateConfigurationWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateConfigurationWithContext failed: %s\n%s", err.Error(), response), PublicCertConfigDnsCISResourceName, "create")
		return tfErr.GetDiag()
	}

	configuration := configurationIntf.(*secretsmanagerv2.PublicCertificateConfigurationDNSCloudInternetServices)
	d.SetId(fmt.Sprintf("%s/%s/%s", region, instanceId, *configuration.Name))

	return resourceIbmSmConfigurationPublicCertificateDNSCisRead(context, d, meta)
}

func resourceIbmSmConfigurationPublicCertificateDNSCisRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PublicCertConfigDnsCISResourceName, "read")
		return tfErr.GetDiag()
	}

	id := strings.Split(d.Id(), "/")
	if len(id) != 3 {
		tfErr := flex.TerraformErrorf(nil, "Wrong format of resource ID. To import a DNS configuration use the format `<region>/<instance_id>/<name>`", PublicCertConfigDnsCISResourceName, "read")
		return tfErr.GetDiag()
	}
	region := id[0]
	instanceId := id[1]
	configName := id[2]
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))
	getConfigurationOptions := &secretsmanagerv2.GetConfigurationOptions{}

	getConfigurationOptions.SetName(configName)

	configurationIntf, response, err := secretsManagerClient.GetConfigurationWithContext(context, getConfigurationOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetConfigurationWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetConfigurationWithContext failed %s\n%s", err, response), PublicCertConfigDnsCISResourceName, "read")
		return tfErr.GetDiag()
	}

	configuration := configurationIntf.(*secretsmanagerv2.PublicCertificateConfigurationDNSCloudInternetServices)
	if err = d.Set("instance_id", instanceId); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting instance_id"), PublicCertConfigDnsCISResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("region", region); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting region"), PublicCertConfigDnsCISResourceName, "read")
		return tfErr.GetDiag()
	}
	if !core.IsNil(configuration.ConfigType) {
		if err = d.Set("config_type", configuration.ConfigType); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting config_type"), PublicCertConfigDnsCISResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if !core.IsNil(configuration.CloudInternetServicesApikey) {
		if err = d.Set("cloud_internet_services_apikey", configuration.CloudInternetServicesApikey); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting cloud_internet_services_apikey"), PublicCertConfigDnsCISResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if !core.IsNil(configuration.CloudInternetServicesCrn) {
		if err = d.Set("cloud_internet_services_crn", configuration.CloudInternetServicesCrn); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting cloud_internet_services_crn"), PublicCertConfigDnsCISResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("name", configuration.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name"), PublicCertConfigDnsCISResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("secret_type", configuration.SecretType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_type"), PublicCertConfigDnsCISResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("created_by", configuration.CreatedBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_by"), PublicCertConfigDnsCISResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("created_at", DateTimeToRFC3339(configuration.CreatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_at"), PublicCertConfigDnsCISResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("updated_at", DateTimeToRFC3339(configuration.UpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated_at"), PublicCertConfigDnsCISResourceName, "read")
		return tfErr.GetDiag()
	}

	return nil
}

func resourceIbmSmConfigurationPublicCertificateDNSCisUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PublicCertConfigDnsCISResourceName, "update")
		return tfErr.GetDiag()
	}

	id := strings.Split(d.Id(), "/")
	region := id[0]
	instanceId := id[1]
	configName := id[2]
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))
	updateConfigurationOptions := &secretsmanagerv2.UpdateConfigurationOptions{}

	updateConfigurationOptions.SetName(configName)

	hasChange := false

	patchVals := &secretsmanagerv2.ConfigurationPatch{}
	if d.HasChange("cloud_internet_services_apikey") {
		newCloudInternetServicesApikey := d.Get("cloud_internet_services_apikey").(string)
		patchVals.CloudInternetServicesApikey = &newCloudInternetServicesApikey
		hasChange = true
	}
	if d.HasChange("cloud_internet_services_crn") {
		newCloudInternetServicesCrn := d.Get("cloud_internet_services_crn").(string)
		patchVals.CloudInternetServicesCrn = &newCloudInternetServicesCrn
		hasChange = true
	}

	if hasChange {
		updateConfigurationOptions.ConfigurationPatch, _ = patchVals.AsPatch()
		_, response, err := secretsManagerClient.UpdateConfigurationWithContext(context, updateConfigurationOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateConfigurationWithContext failed %s\n%s", err, response)
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateConfigurationWithContext failed %s\n%s", err, response), PublicCertConfigDnsCISResourceName, "update")
			return tfErr.GetDiag()
		}
	}

	return resourceIbmSmConfigurationPublicCertificateDNSCisRead(context, d, meta)
}

func resourceIbmSmConfigurationPublicCertificateDNSCisDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PublicCertConfigDnsCISResourceName, "delete")
		return tfErr.GetDiag()
	}

	id := strings.Split(d.Id(), "/")
	region := id[0]
	instanceId := id[1]
	configName := id[2]
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))
	deleteConfigurationOptions := &secretsmanagerv2.DeleteConfigurationOptions{}

	deleteConfigurationOptions.SetName(configName)

	response, err := secretsManagerClient.DeleteConfigurationWithContext(context, deleteConfigurationOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteConfigurationWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteConfigurationWithContext failed %s\n%s", err, response), PublicCertConfigDnsCISResourceName, "delete")
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func resourceIbmSmConfigurationPublicCertificateCisMapToPublicCertificateConfigurationDNSCloudInternetServicesPrototype(modelMap map[string]interface{}) (*secretsmanagerv2.PublicCertificateConfigurationDNSCloudInternetServicesPrototype, error) {
	model := &secretsmanagerv2.PublicCertificateConfigurationDNSCloudInternetServicesPrototype{}
	model.ConfigType = core.StringPtr("public_cert_configuration_dns_cloud_internet_services")
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["cloud_internet_services_apikey"] != nil && modelMap["cloud_internet_services_apikey"].(string) != "" {
		model.CloudInternetServicesApikey = core.StringPtr(modelMap["cloud_internet_services_apikey"].(string))
	}
	model.CloudInternetServicesCrn = core.StringPtr(modelMap["cloud_internet_services_crn"].(string))
	return model, nil
}
