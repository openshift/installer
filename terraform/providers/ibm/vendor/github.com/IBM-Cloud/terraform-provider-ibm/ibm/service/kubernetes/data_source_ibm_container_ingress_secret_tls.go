// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kubernetes

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMContainerIngressSecretTLS() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMContainerIngressSecretTLSRead,
		Schema: map[string]*schema.Schema{
			"cluster": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cluster ID or name",
				ValidateFunc: validate.InvokeValidator(
					"ibm_container_ingress_secret_tls",
					"cluster"),
			},
			"secret_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Secret name",
			},
			"secret_namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Secret namespace",
			},
			"cert_crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate CRN",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type TLS",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Secret Status",
			},
			"user_managed": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If the secret was created by the user",
			},
			"persistence": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Persistence of secret",
			},
			"domain_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Domain name",
			},
			"expires_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate expires on date",
			},
			"last_updated_timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp secret was last updated",
			},
		},
	}
}

func DataSourceIBMContainerIngressSecretTLSValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cluster",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			Required:                   true,
			CloudDataType:              "cluster",
			CloudDataRange:             []string{"resolved_to:id"}})

	iBMContainerIngressSecretValidator := validate.ResourceValidator{ResourceName: "ibm_container_ingress_secret_tls", Schema: validateSchema}
	return &iBMContainerIngressSecretValidator
}

func dataSourceIBMContainerIngressSecretTLSRead(d *schema.ResourceData, meta interface{}) error {
	ingressClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}

	clusterID := d.Get("cluster").(string)
	name := d.Get("secret_name").(string)
	namespace := d.Get("secret_namespace").(string)

	ingressAPI := ingressClient.Ingresses()
	ingressSecretConfig, err := ingressAPI.GetIngressSecret(clusterID, name, namespace)
	if err != nil {
		return err
	}

	d.Set("cluster", ingressSecretConfig.Cluster)
	d.Set("secret_name", ingressSecretConfig.Name)
	d.Set("secret_namespace", ingressSecretConfig.Namespace)
	d.Set("type", ingressSecretConfig.Type)
	d.Set("cert_crn", ingressSecretConfig.CRN)
	d.Set("persistence", ingressSecretConfig.Persistence)
	d.Set("domain_name", ingressSecretConfig.Domain)
	d.Set("expires_on", ingressSecretConfig.ExpiresOn)
	d.Set("status", ingressSecretConfig.Status)
	d.Set("user_managed", ingressSecretConfig.UserManaged)
	d.Set("last_updated_timestamp", ingressSecretConfig.LastUpdatedTimestamp)

	d.SetId(fmt.Sprintf("%s/%s/%s", clusterID, name, namespace))

	return nil
}
