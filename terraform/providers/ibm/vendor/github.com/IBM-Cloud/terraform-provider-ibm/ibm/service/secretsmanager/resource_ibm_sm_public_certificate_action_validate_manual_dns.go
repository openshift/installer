// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package secretsmanager

import (
	"context"
	"fmt"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/secrets-manager-go-sdk/v2/secretsmanagerv2"
)

func ResourceIbmSmPublicCertificateActionValidateManualDns() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmSmPublicCertificateActionValidateManualDnsCreateOrUpdate,
		ReadContext:   resourceIbmSmPublicCertificateActionValidateManualDnsRead,
		UpdateContext: resourceIbmSmPublicCertificateActionValidateManualDnsCreateOrUpdate,
		DeleteContext: resourceIbmSmPublicCertificateActionValidateManualDnsDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"secret_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "A v4 UUID identifier.",
				ForceNew:    true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(35 * time.Minute),
		},
	}
}

func resourceIbmSmPublicCertificateActionValidateManualDnsCreateOrUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PublicCertConfigActionValidateManualDNSResourceName, "create/update")
		return tfErr.GetDiag()
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	d.SetId(fmt.Sprintf("%s/%s/%s/validate_manual_dns", region, instanceId, d.Get("secret_id").(string)))

	validateManualDns(context, d, secretsManagerClient)

	return nil
}

func resourceIbmSmPublicCertificateActionValidateManualDnsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceIbmSmPublicCertificateActionValidateManualDnsDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}

func validateManualDns(context context.Context, d *schema.ResourceData, secretsManagerClient *secretsmanagerv2.SecretsManagerV2) diag.Diagnostics {
	createActionOptions := &secretsmanagerv2.CreateSecretActionOptions{}

	actionType := "public_cert_action_validate_dns_challenge"
	createActionOptions.SetSecretActionPrototype(&secretsmanagerv2.SecretActionPrototype{
		ActionType: &actionType,
	})
	createActionOptions.SetID(d.Get("secret_id").(string))

	_, response, err := secretsManagerClient.CreateSecretActionWithContext(context, createActionOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateSecretActionWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateSecretActionWithContext failed: %s\n%s", err.Error(), response), PublicCertConfigActionValidateManualDNSResourceName, "create")
		return tfErr.GetDiag()
	}

	_, err = waitForIbmSmPublicCertificateCreate(secretsManagerClient, d, "pre_activation", "active")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error waiting for resource IbmSmPublicCertificateActionValidateManualDns (%s) to be created: %s", d.Id(), err.Error()), PublicCertConfigActionValidateManualDNSResourceName, "create")
		return tfErr.GetDiag()
	}
	return nil
}
