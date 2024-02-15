// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package mqcloud

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/mqcloud-go-sdk/mqcloudv1"
)

func ResourceIbmMqcloudKeystoreCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmMqcloudKeystoreCertificateCreate,
		ReadContext:   resourceIbmMqcloudKeystoreCertificateRead,
		DeleteContext: resourceIbmMqcloudKeystoreCertificateDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"service_instance_guid": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_mqcloud_keystore_certificate", "service_instance_guid"),
				Description:  "The GUID that uniquely identifies the MQ on Cloud service instance.",
			},
			"queue_manager_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_mqcloud_keystore_certificate", "queue_manager_id"),
				Description:  "The id of the queue manager to retrieve its full details.",
			},
			"label": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_mqcloud_keystore_certificate", "label"),
				Description:  "The label to use for the certificate to be uploaded.",
			},
			"certificate_file": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The filename and path of the certificate to be uploaded.",
			},
			"certificate_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of certificate.",
			},
			"fingerprint_sha256": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Fingerprint SHA256.",
			},
			"subject_dn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Subject's Distinguished Name.",
			},
			"subject_cn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Subject's Common Name.",
			},
			"issuer_dn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Issuer's Distinguished Name.",
			},
			"issuer_cn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Issuer's Common Name.",
			},
			"issued": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date certificate was issued.",
			},
			"expiry": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expiry date for the certificate.",
			},
			"is_default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether it is the queue manager's default certificate.",
			},
			"dns_names_total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total count of dns names.",
			},
			"dns_names": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of DNS names.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this key store certificate.",
			},
			"certificate_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the certificate.",
			},
		},
	}
}

func ResourceIbmMqcloudKeystoreCertificateValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "service_instance_guid",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
			MinValueLength:             36,
			MaxValueLength:             36,
		},
		validate.ValidateSchema{
			Identifier:                 "queue_manager_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[0-9a-fA-F]{32}$`,
			MinValueLength:             32,
			MaxValueLength:             32,
		},
		validate.ValidateSchema{
			Identifier:                 "label",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-zA-Z0-9_.]*$`,
			MinValueLength:             1,
			MaxValueLength:             64,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_mqcloud_keystore_certificate", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmMqcloudKeystoreCertificateCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		return diag.FromErr(err)
	}
	err = checkSIPlan(d, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Create Keystore Certificate failed %s", err))
	}
	createKeyStorePemCertificateOptions := &mqcloudv1.CreateKeyStorePemCertificateOptions{}

	createKeyStorePemCertificateOptions.SetServiceInstanceGuid(d.Get("service_instance_guid").(string))
	createKeyStorePemCertificateOptions.SetQueueManagerID(d.Get("queue_manager_id").(string))
	createKeyStorePemCertificateOptions.SetLabel(d.Get("label").(string))
	certificateFileBytes, err := base64.StdEncoding.DecodeString(d.Get("certificate_file").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	createKeyStorePemCertificateOptions.SetCertificateFile(io.NopCloser(bytes.NewReader(certificateFileBytes)))

	keyStoreCertificateDetails, response, err := mqcloudClient.CreateKeyStorePemCertificateWithContext(context, createKeyStorePemCertificateOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateKeyStorePemCertificateWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateKeyStorePemCertificateWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", *createKeyStorePemCertificateOptions.ServiceInstanceGuid, *createKeyStorePemCertificateOptions.QueueManagerID, *keyStoreCertificateDetails.ID))

	return resourceIbmMqcloudKeystoreCertificateRead(context, d, meta)
}

func resourceIbmMqcloudKeystoreCertificateRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getKeyStoreCertificateOptions := &mqcloudv1.GetKeyStoreCertificateOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getKeyStoreCertificateOptions.SetServiceInstanceGuid(parts[0])
	getKeyStoreCertificateOptions.SetQueueManagerID(parts[1])
	getKeyStoreCertificateOptions.SetCertificateID(parts[2])

	keyStoreCertificateDetails, response, err := mqcloudClient.GetKeyStoreCertificateWithContext(context, getKeyStoreCertificateOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetKeyStoreCertificateWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetKeyStoreCertificateWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("service_instance_guid", parts[0]); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting service_instance_guid: %s", err))
	}
	if err = d.Set("queue_manager_id", parts[1]); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting queue_manager_id: %s", err))
	}
	if err = d.Set("label", keyStoreCertificateDetails.Label); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting label: %s", err))
	}
	if err = d.Set("certificate_type", keyStoreCertificateDetails.CertificateType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting certificate_type: %s", err))
	}
	if err = d.Set("fingerprint_sha256", keyStoreCertificateDetails.FingerprintSha256); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting fingerprint_sha256: %s", err))
	}
	if err = d.Set("subject_dn", keyStoreCertificateDetails.SubjectDn); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting subject_dn: %s", err))
	}
	if err = d.Set("subject_cn", keyStoreCertificateDetails.SubjectCn); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting subject_cn: %s", err))
	}
	if err = d.Set("issuer_dn", keyStoreCertificateDetails.IssuerDn); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting issuer_dn: %s", err))
	}
	if err = d.Set("issuer_cn", keyStoreCertificateDetails.IssuerCn); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting issuer_cn: %s", err))
	}
	if err = d.Set("issued", flex.DateTimeToString(keyStoreCertificateDetails.Issued)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting issued: %s", err))
	}
	if err = d.Set("expiry", flex.DateTimeToString(keyStoreCertificateDetails.Expiry)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting expiry: %s", err))
	}
	if err = d.Set("is_default", keyStoreCertificateDetails.IsDefault); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting is_default: %s", err))
	}
	if err = d.Set("dns_names_total_count", flex.IntValue(keyStoreCertificateDetails.DnsNamesTotalCount)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting dns_names_total_count: %s", err))
	}
	if err = d.Set("dns_names", keyStoreCertificateDetails.DnsNames); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting dns_names: %s", err))
	}
	if err = d.Set("href", keyStoreCertificateDetails.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("certificate_id", keyStoreCertificateDetails.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting certificate_id: %s", err))
	}

	return nil
}

func resourceIbmMqcloudKeystoreCertificateDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		return diag.FromErr(err)
	}
	err = checkSIPlan(d, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Delete Keystore Certificate failed %s", err))
	}

	deleteKeyStoreCertificateOptions := &mqcloudv1.DeleteKeyStoreCertificateOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteKeyStoreCertificateOptions.SetServiceInstanceGuid(parts[0])
	deleteKeyStoreCertificateOptions.SetQueueManagerID(parts[1])
	deleteKeyStoreCertificateOptions.SetCertificateID(parts[2])

	response, err := mqcloudClient.DeleteKeyStoreCertificateWithContext(context, deleteKeyStoreCertificateOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteKeyStoreCertificateWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteKeyStoreCertificateWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}
