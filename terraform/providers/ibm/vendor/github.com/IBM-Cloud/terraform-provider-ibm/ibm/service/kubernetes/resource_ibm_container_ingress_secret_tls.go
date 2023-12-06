// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kubernetes

import (
	"fmt"
	"time"

	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMContainerIngressSecretTLS() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMContainerIngressSecretTLSCreate,
		Read:     resourceIBMContainerIngressSecretTLSRead,
		Update:   resourceIBMContainerIngressSecretTLSUpdate,
		Delete:   resourceIBMContainerIngressSecretTLSDelete,
		Exists:   resourceIBMContainerIngressSecretTLSExists,
		Importer: &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"cluster": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cluster ID or name",
				ValidateFunc: validate.InvokeValidator(
					"ibm_container_ingress_secret_tls",
					"cluster"),
			},
			"secret_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Secret name",
				ForceNew:    true,
				ValidateFunc: validate.InvokeValidator(
					"ibm_container_ingress_secret_tls",
					"secret_name",
				),
			},
			"secret_namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Secret namespace",
				ForceNew:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "TLS secret type",
			},
			"cert_crn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Certificate CRN",
			},
			"persistence": {
				Type:        schema.TypeBool,
				Optional:    true,
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
			"last_updated_timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp secret was last updated",
			},
			"update_secret": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Updates secret from secrets manager if value is changed (increment each usage)",
			},
		},
	}
}

func ResourceIBMContainerIngressSecretTLSValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cluster",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			Required:                   true,
			CloudDataType:              "cluster",
			CloudDataRange:             []string{"resolved_to:id"}})

	validateSchema = append(validateSchema, validate.ValidateSchema{
		Identifier:                 "secret_name",
		ValidateFunctionIdentifier: validate.ValidateRegexpLen,
		Type:                       validate.TypeString,
		Required:                   true,
		Regexp:                     `^([a-z0-9]([-a-z0-9]*[a-z0-9])?(.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*)$`,
		MinValueLength:             1,
		MaxValueLength:             63,
	})
	iBMContainerIngressInstanceValidator := validate.ResourceValidator{ResourceName: "ibm_container_ingress_secret_tls", Schema: validateSchema}
	return &iBMContainerIngressInstanceValidator
}

func resourceIBMContainerIngressSecretTLSCreate(d *schema.ResourceData, meta interface{}) error {
	ingressClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}

	cluster := d.Get("cluster").(string)
	secretName := d.Get("secret_name").(string)
	secretNamespace := d.Get("secret_namespace").(string)
	secretType := "TLS"
	certCRN := d.Get("cert_crn").(string)

	params := v2.SecretCreateConfig{
		Cluster:   cluster,
		Name:      secretName,
		Namespace: secretNamespace,
		Type:      secretType,
		CRN:       certCRN,
	}

	if persistence, ok := d.GetOk("persistence"); ok {
		params.Persistence = persistence.(bool)
	}

	ingressAPI := ingressClient.Ingresses()
	response, err := ingressAPI.CreateIngressSecret(params)

	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%s/%s/%s", cluster, response.Name, response.Namespace))

	return resourceIBMContainerIngressSecretTLSRead(d, meta)
}

func resourceIBMContainerIngressSecretTLSRead(d *schema.ResourceData, meta interface{}) error {
	ingressClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}
	cluster := parts[0]
	secretName := parts[1]
	secretNamespace := parts[2]

	ingressAPI := ingressClient.Ingresses()
	ingressSecretConfig, err := ingressAPI.GetIngressSecret(cluster, secretName, secretNamespace)
	if err != nil {
		return err
	}

	d.Set("cluster", cluster)
	d.Set("secret_name", ingressSecretConfig.Name)
	d.Set("secret_namespace", ingressSecretConfig.Namespace)
	d.Set("cert_crn", ingressSecretConfig.CRN)
	d.Set("persistence", ingressSecretConfig.Persistence)
	d.Set("domain_name", ingressSecretConfig.Domain)
	d.Set("expires_on", ingressSecretConfig.ExpiresOn)
	d.Set("status", ingressSecretConfig.Status)
	d.Set("type", ingressSecretConfig.Type)
	d.Set("user_managed", ingressSecretConfig.UserManaged)
	d.Set("last_updated_timestamp", ingressSecretConfig.LastUpdatedTimestamp)

	return nil
}

func resourceIBMContainerIngressSecretTLSDelete(d *schema.ResourceData, meta interface{}) error {
	ingressClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}

	ingressAPI := ingressClient.Ingresses()

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}
	cluster := parts[0]
	secretName := parts[1]
	secretNamespace := parts[2]

	params := v2.SecretDeleteConfig{
		Cluster:   cluster,
		Name:      secretName,
		Namespace: secretNamespace,
	}

	err = ingressAPI.DeleteIngressSecret(params)
	if err != nil {
		return err
	}

	return nil
}

func resourceIBMContainerIngressSecretTLSUpdate(d *schema.ResourceData, meta interface{}) error {
	ingressClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}
	cluster := parts[0]
	secretName := parts[1]
	secretNamespace := parts[2]

	params := v2.SecretUpdateConfig{
		Cluster:   cluster,
		Name:      secretName,
		Namespace: secretNamespace,
	}

	ingressAPI := ingressClient.Ingresses()
	if d.HasChange("cert_crn") {
		params.CRN = d.Get("cert_crn").(string)

		_, err := ingressAPI.UpdateIngressSecret(params)
		if err != nil {
			return err
		}
	} else if d.HasChange("update_secret") {
		// user wants to force an upstream secret update from secrets manager onto kube cluster w/out changing crn
		_, err = ingressAPI.UpdateIngressSecret(params)
		if err != nil {
			return err
		}
	}

	return resourceIBMContainerIngressSecretTLSRead(d, meta)
}

func resourceIBMContainerIngressSecretTLSExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	ingressClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return false, err
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return false, err
	}
	cluster := parts[0]
	secretName := parts[1]
	secretNamespace := parts[2]

	ingressAPI := ingressClient.Ingresses()
	ingressSecretConfig, err := ingressAPI.GetIngressSecret(cluster, secretName, secretNamespace)

	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("[ERROR] Error getting ingress secret: %s", err)
	}

	return ingressSecretConfig.Name == secretName && ingressSecretConfig.Namespace == secretNamespace && ingressSecretConfig.Status != "deleted", nil
}
