package secretsmanager

import (
	"context"
	"fmt"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/secrets-manager-go-sdk/v2/secretsmanagerv2"
)

func ResourceIbmSmPrivateCertificateConfigurationActionSetSigned() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmSmPrivateCertificateConfigurationActionSetSignedCreateOrUpdate,
		ReadContext:   resourceIbmSmPrivateCertificateConfigurationActionSetSignedRead,
		UpdateContext: resourceIbmSmPrivateCertificateConfigurationActionSetSignedCreateOrUpdate,
		DeleteContext: resourceIbmSmPrivateCertificateConfigurationActionSetSignedDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name that uniquely identifies a configuration",
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
				ForceNew:    true,
				Description: "The PEM-encoded certificate.",
			},
		},
	}
}

func resourceIbmSmPrivateCertificateConfigurationActionSetSignedCreateOrUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PrivateCertConfigActionSetSigned, "create/update")
		return tfErr.GetDiag()
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	createConfigurationActionOptions := &secretsmanagerv2.CreateConfigurationActionOptions{}

	configurationActionPrototypeModel, err := resourceIbmSmPrivateCertificateConfigurationActionSetSignedPrototype(d)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PrivateCertConfigActionSetSigned, "create/update")
		return tfErr.GetDiag()
	}
	createConfigurationActionOptions.SetConfigActionPrototype(configurationActionPrototypeModel)
	createConfigurationActionOptions.SetName(d.Get("name").(string))

	_, response, err := secretsManagerClient.CreateConfigurationActionWithContext(context, createConfigurationActionOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateConfigurationActionWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateConfigurationActionWithContext failed: %s\n%s", err.Error(), response), PrivateCertConfigActionSetSigned, "create/update")
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s/%s/set_signed", region, instanceId, d.Get("name").(string)))

	return nil
}

func resourceIbmSmPrivateCertificateConfigurationActionSetSignedRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceIbmSmPrivateCertificateConfigurationActionSetSignedDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}

func resourceIbmSmPrivateCertificateConfigurationActionSetSignedPrototype(d *schema.ResourceData) (secretsmanagerv2.ConfigurationActionPrototypeIntf, error) {
	model := &secretsmanagerv2.PrivateCertificateConfigurationActionSetSignedPrototype{
		ActionType: core.StringPtr("private_cert_configuration_action_set_signed"),
	}
	if _, ok := d.GetOk("certificate"); ok {
		model.Certificate = core.StringPtr(formatCertificate(d.Get("certificate").(string)))
	}

	return model, nil
}
