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

func ResourceIbmMqcloudTruststoreCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmMqcloudTruststoreCertificateCreate,
		ReadContext:   resourceIbmMqcloudTruststoreCertificateRead,
		DeleteContext: resourceIbmMqcloudTruststoreCertificateDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"service_instance_guid": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_mqcloud_truststore_certificate", "service_instance_guid"),
				Description:  "The GUID that uniquely identifies the MQ on Cloud service instance.",
			},
			"queue_manager_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_mqcloud_truststore_certificate", "queue_manager_id"),
				Description:  "The id of the queue manager to retrieve its full details.",
			},
			"label": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_mqcloud_truststore_certificate", "label"),
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
				Description: "The Date the certificate was issued.",
			},
			"expiry": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expiry date for the certificate.",
			},
			"trusted": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether a certificate is trusted.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this trust store certificate.",
			},
			"certificate_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Id of the certificate.",
			},
		},
	}
}

func ResourceIbmMqcloudTruststoreCertificateValidator() *validate.ResourceValidator {
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

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_mqcloud_truststore_certificate", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmMqcloudTruststoreCertificateCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		return diag.FromErr(err)
	}
	err = checkSIPlan(d, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Create Truststore Certificate failed %s", err))
	}

	createTrustStorePemCertificateOptions := &mqcloudv1.CreateTrustStorePemCertificateOptions{}

	createTrustStorePemCertificateOptions.SetServiceInstanceGuid(d.Get("service_instance_guid").(string))
	createTrustStorePemCertificateOptions.SetQueueManagerID(d.Get("queue_manager_id").(string))
	createTrustStorePemCertificateOptions.SetLabel(d.Get("label").(string))
	certificateFileBytes, err := base64.StdEncoding.DecodeString(d.Get("certificate_file").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	createTrustStorePemCertificateOptions.SetCertificateFile(io.NopCloser(bytes.NewReader(certificateFileBytes)))

	trustStoreCertificateDetails, response, err := mqcloudClient.CreateTrustStorePemCertificateWithContext(context, createTrustStorePemCertificateOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateTrustStorePemCertificateWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateTrustStorePemCertificateWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", *createTrustStorePemCertificateOptions.ServiceInstanceGuid, *createTrustStorePemCertificateOptions.QueueManagerID, *trustStoreCertificateDetails.ID))

	return resourceIbmMqcloudTruststoreCertificateRead(context, d, meta)
}

func resourceIbmMqcloudTruststoreCertificateRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getTrustStoreCertificateOptions := &mqcloudv1.GetTrustStoreCertificateOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getTrustStoreCertificateOptions.SetServiceInstanceGuid(parts[0])
	getTrustStoreCertificateOptions.SetQueueManagerID(parts[1])
	getTrustStoreCertificateOptions.SetCertificateID(parts[2])

	trustStoreCertificateDetails, response, err := mqcloudClient.GetTrustStoreCertificateWithContext(context, getTrustStoreCertificateOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetTrustStoreCertificateWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetTrustStoreCertificateWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("service_instance_guid", parts[0]); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting service_instance_guid: %s", err))
	}
	if err = d.Set("queue_manager_id", parts[1]); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting queue_manager_id: %s", err))
	}
	if err = d.Set("label", trustStoreCertificateDetails.Label); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting label: %s", err))
	}
	if err = d.Set("certificate_type", trustStoreCertificateDetails.CertificateType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting certificate_type: %s", err))
	}
	if err = d.Set("fingerprint_sha256", trustStoreCertificateDetails.FingerprintSha256); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting fingerprint_sha256: %s", err))
	}
	if err = d.Set("subject_dn", trustStoreCertificateDetails.SubjectDn); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting subject_dn: %s", err))
	}
	if err = d.Set("subject_cn", trustStoreCertificateDetails.SubjectCn); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting subject_cn: %s", err))
	}
	if err = d.Set("issuer_dn", trustStoreCertificateDetails.IssuerDn); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting issuer_dn: %s", err))
	}
	if err = d.Set("issuer_cn", trustStoreCertificateDetails.IssuerCn); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting issuer_cn: %s", err))
	}
	if err = d.Set("issued", flex.DateTimeToString(trustStoreCertificateDetails.Issued)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting issued: %s", err))
	}
	if err = d.Set("expiry", flex.DateTimeToString(trustStoreCertificateDetails.Expiry)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting expiry: %s", err))
	}
	if err = d.Set("trusted", trustStoreCertificateDetails.Trusted); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting trusted: %s", err))
	}
	if err = d.Set("href", trustStoreCertificateDetails.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("certificate_id", trustStoreCertificateDetails.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting certificate_id: %s", err))
	}

	return nil
}

func resourceIbmMqcloudTruststoreCertificateDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		return diag.FromErr(err)
	}
	err = checkSIPlan(d, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Delete Truststore Certificate failed %s", err))
	}
	deleteTrustStoreCertificateOptions := &mqcloudv1.DeleteTrustStoreCertificateOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteTrustStoreCertificateOptions.SetServiceInstanceGuid(parts[0])
	deleteTrustStoreCertificateOptions.SetQueueManagerID(parts[1])
	deleteTrustStoreCertificateOptions.SetCertificateID(parts[2])

	response, err := mqcloudClient.DeleteTrustStoreCertificateWithContext(context, deleteTrustStoreCertificateOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteTrustStoreCertificateWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteTrustStoreCertificateWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}
