// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
	"github.com/IBM-Cloud/bluemix-go/bmxerror"
)

func resourceIBMContainerALBCert() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMContainerALBCertCreate,
		Read:     resourceIBMContainerALBCertRead,
		Update:   resourceIBMContainerALBCertUpdate,
		Delete:   resourceIBMContainerALBCertDelete,
		Exists:   resourceIBMContainerALBCertExists,
		Importer: &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"cert_crn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "Certificate CRN id",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cluster ID",
			},
			"secret_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Secret name",
			},
			"namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "ibm-cert-store",
				ForceNew:    true,
				Description: "Namespace of the secret",
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
				Description: "Certificate expaire on date",
			},
			"issuer_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "certificate issuer name",
				Deprecated:  "This field is depricated and is not available in v2 version of ingress api",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Secret Status",
			},
			"cloud_cert_instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "cloud cert instance ID",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Deprecated:  "This field is deprecated",
				Description: "region name",
			},
		},
	}
}

func resourceIBMContainerALBCertCreate(d *schema.ResourceData, meta interface{}) error {
	ingressClient, err := meta.(ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}

	certCRN := d.Get("cert_crn").(string)
	cluster := d.Get("cluster_id").(string)
	secretName := d.Get("secret_name").(string)
	namespace := d.Get("namespace").(string)

	params := v2.SecretCreateConfig{
		CRN:       certCRN,
		Cluster:   cluster,
		Name:      secretName,
		Namespace: namespace,
	}
	// params.State = "update_false"
	if v, ok := d.GetOk("persistence"); ok {
		params.Persistence = v.(bool)
	}

	ingressAPI := ingressClient.Ingresses()
	response, err := ingressAPI.CreateIngressSecret(params)

	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%s/%s/%s", cluster, secretName, response.Namespace))
	_, err = waitForContainerALBCert(d, meta, schema.TimeoutCreate)
	if err != nil {
		return fmt.Errorf(
			"Error waiting for create resource alb cert (%s) : %s", d.Id(), err)
	}

	return resourceIBMContainerALBCertRead(d, meta)
}

func resourceIBMContainerALBCertRead(d *schema.ResourceData, meta interface{}) error {
	ingressClient, err := meta.(ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	clusterID := parts[0]
	secretName := parts[1]
	namespace := "ibm-cert-store"
	if len(parts) > 2 && len(parts[2]) > 0 {
		namespace = parts[2]
	}

	ingressAPI := ingressClient.Ingresses()
	ingressSecretConfig, err := ingressAPI.GetIngressSecret(clusterID, secretName, namespace)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%s/%s/%s", clusterID, secretName, namespace))
	d.Set("cluster_id", ingressSecretConfig.Cluster)
	d.Set("secret_name", ingressSecretConfig.Name)
	d.Set("namespace", ingressSecretConfig.Namespace)
	d.Set("cert_crn", ingressSecretConfig.CRN)
	instancecrn := strings.Split(ingressSecretConfig.CRN, ":certificate:")
	d.Set("cloud_cert_instance_id", fmt.Sprintf("%s::", instancecrn[0]))
	d.Set("domain_name", ingressSecretConfig.Domain)
	d.Set("expires_on", ingressSecretConfig.ExpiresOn)
	d.Set("status", ingressSecretConfig.Status)
	d.Set("persistence", ingressSecretConfig.Persistence)

	return nil
}

func resourceIBMContainerALBCertDelete(d *schema.ResourceData, meta interface{}) error {
	ingressClient, err := meta.(ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}

	ingressAPI := ingressClient.Ingresses()

	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	clusterID := parts[0]
	secretName := parts[1]
	namespace := "ibm-cert-store"
	if len(parts) > 2 && len(parts[2]) > 0 {
		namespace = parts[2]
	}
	params := v2.SecretDeleteConfig{
		Cluster:   clusterID,
		Name:      secretName,
		Namespace: namespace,
	}

	err = ingressAPI.DeleteIngressSecret(params)
	if err != nil {
		return err
	}
	_, albCertDeletionError := waitForALBCertDelete(d, meta, schema.TimeoutDelete)
	if albCertDeletionError != nil {
		return albCertDeletionError
	}
	d.SetId("")
	return nil
}

func waitForALBCertDelete(d *schema.ResourceData, meta interface{}, timeout string) (interface{}, error) {
	ingressClient, err := meta.(ClientSession).VpcContainerAPI()
	if err != nil {
		return false, err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return false, err
	}
	clusterID := parts[0]
	secretName := parts[1]
	namespace := "ibm-cert-store"
	if len(parts) > 2 && len(parts[2]) > 0 {
		namespace = parts[2]
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"deleting"},
		Target:  []string{"deleted"},
		Refresh: func() (interface{}, string, error) {

			secret, err := ingressClient.Ingresses().GetIngressSecret(clusterID, secretName, namespace)
			if err != nil {
				if apiErr, ok := err.(bmxerror.RequestFailure); ok && apiErr.StatusCode() == 404 {
					return secret, "deleted", nil
				}
				return nil, "", err
			}
			if secret.Status != "deleted" {
				return secret, "deleting", nil
			}
			return secret, "deleted", nil
		},
		Timeout:    d.Timeout(timeout),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func resourceIBMContainerALBCertUpdate(d *schema.ResourceData, meta interface{}) error {
	ingressClient, err := meta.(ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	cluster := parts[0]
	secretName := parts[1]
	namespace := "ibm-cert-store"
	if len(parts) > 2 && len(parts[2]) > 0 {
		namespace = parts[2]
	}

	if d.HasChange("cert_crn") {
		crn := d.Get("cert_crn").(string)
		params := v2.SecretUpdateConfig{
			CRN:       crn,
			Cluster:   cluster,
			Name:      secretName,
			Namespace: namespace,
		}
		// params.State = "update_true"

		ingressAPI := ingressClient.Ingresses()
		_, err = ingressAPI.UpdateIngressSecret(params)
		if err != nil {
			return err
		}

		_, err = waitForContainerALBCert(d, meta, schema.TimeoutUpdate)
		if err != nil {
			return fmt.Errorf(
				"Error waiting for updating resource alb cert (%s) : %s", d.Id(), err)
		}
	}
	return resourceIBMContainerALBCertRead(d, meta)
}

func resourceIBMContainerALBCertExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	ingressClient, err := meta.(ClientSession).VpcContainerAPI()
	if err != nil {
		return false, err
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return false, err
	}
	clusterID := parts[0]
	secretName := parts[1]
	namespace := "ibm-cert-store"
	if len(parts) > 2 && len(parts[2]) > 0 {
		namespace = parts[2]
	}

	ingressAPI := ingressClient.Ingresses()
	ingressSecretConfig, err := ingressAPI.GetIngressSecret(clusterID, secretName, namespace)

	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}

	return ingressSecretConfig.Cluster == clusterID && ingressSecretConfig.Name == secretName, nil
}

func waitForContainerALBCert(d *schema.ResourceData, meta interface{}, timeout string) (interface{}, error) {
	ingressClient, err := meta.(ClientSession).VpcContainerAPI()
	if err != nil {
		return false, err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return false, err
	}
	clusterID := parts[0]
	secretName := parts[1]
	namespace := "ibm-cert-store"
	if len(parts) > 2 && len(parts[2]) > 0 {
		namespace = parts[2]
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"creating"},
		Target:  []string{"done"},
		Refresh: func() (interface{}, string, error) {

			alb, err := ingressClient.Ingresses().GetIngressSecret(clusterID, secretName, namespace)
			if err != nil {
				if apiErr, ok := err.(bmxerror.RequestFailure); ok && apiErr.StatusCode() == 404 {
					return alb, "creating", nil
				}
				return nil, "", err
			}
			if alb.Status != "created" {
				if strings.Contains(alb.Status, "failed") {
					return alb, "failed", fmt.Errorf("The resource alb cert %s failed: %v", d.Id(), err)
				}

				if alb.Status == "updated" {
					return alb, "done", nil
				}
				return alb, "creating", nil
			}
			return alb, "done", nil
		},
		Timeout:    d.Timeout(timeout),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}
