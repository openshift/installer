// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.95.2-120e65bc-20240924-152329
 */

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
				Description:  "The GUID that uniquely identifies the MQaaS service instance.",
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
		// Error is coming from SDK client, so it doesn't need to be discriminated.
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_mqcloud_truststore_certificate", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	err = checkSIPlan(d, meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Create Truststore Certificate failed: %s", err.Error()), "ibm_mqcloud_truststore_certificate", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createTrustStorePemCertificateOptions := &mqcloudv1.CreateTrustStorePemCertificateOptions{}

	createTrustStorePemCertificateOptions.SetServiceInstanceGuid(d.Get("service_instance_guid").(string))
	createTrustStorePemCertificateOptions.SetQueueManagerID(d.Get("queue_manager_id").(string))
	createTrustStorePemCertificateOptions.SetLabel(d.Get("label").(string))
	certificateFileBytes, err := base64.StdEncoding.DecodeString(d.Get("certificate_file").(string))
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_truststore_certificate", "create", "parse-certificate_file").GetDiag()
	}
	createTrustStorePemCertificateOptions.SetCertificateFile(io.NopCloser(bytes.NewReader(certificateFileBytes)))

	trustStoreCertificateDetails, _, err := mqcloudClient.CreateTrustStorePemCertificateWithContext(context, createTrustStorePemCertificateOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateTrustStorePemCertificateWithContext failed: %s", err.Error()), "ibm_mqcloud_truststore_certificate", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", *createTrustStorePemCertificateOptions.ServiceInstanceGuid, *createTrustStorePemCertificateOptions.QueueManagerID, *trustStoreCertificateDetails.ID))

	return resourceIbmMqcloudTruststoreCertificateRead(context, d, meta)
}

func resourceIbmMqcloudTruststoreCertificateRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_mqcloud_truststore_certificate", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getTrustStoreCertificateOptions := &mqcloudv1.GetTrustStoreCertificateOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_truststore_certificate", "read", "sep-id-parts").GetDiag()
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
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetTrustStoreCertificateWithContext failed: %s", err.Error()), "ibm_mqcloud_truststore_certificate", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("service_instance_guid", parts[0]); err != nil {
		err = fmt.Errorf("Error setting service_instance_guid: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_truststore_certificate", "read", "set-service_instance_guid").GetDiag()
	}
	if err = d.Set("queue_manager_id", parts[1]); err != nil {
		err = fmt.Errorf("Error setting service_instance_guid: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_truststore_certificate", "read", "set-queue_manager_id").GetDiag()
	}
	if err = d.Set("label", trustStoreCertificateDetails.Label); err != nil {
		err = fmt.Errorf("Error setting label: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_truststore_certificate", "read", "set-label").GetDiag()
	}
	if err = d.Set("certificate_type", trustStoreCertificateDetails.CertificateType); err != nil {
		err = fmt.Errorf("Error setting certificate_type: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_truststore_certificate", "read", "set-certificate_type").GetDiag()
	}
	if err = d.Set("fingerprint_sha256", trustStoreCertificateDetails.FingerprintSha256); err != nil {
		err = fmt.Errorf("Error setting fingerprint_sha256: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_truststore_certificate", "read", "set-fingerprint_sha256").GetDiag()
	}
	if err = d.Set("subject_dn", trustStoreCertificateDetails.SubjectDn); err != nil {
		err = fmt.Errorf("Error setting subject_dn: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_truststore_certificate", "read", "set-subject_dn").GetDiag()
	}
	if err = d.Set("subject_cn", trustStoreCertificateDetails.SubjectCn); err != nil {
		err = fmt.Errorf("Error setting subject_cn: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_truststore_certificate", "read", "set-subject_cn").GetDiag()
	}
	if err = d.Set("issuer_dn", trustStoreCertificateDetails.IssuerDn); err != nil {
		err = fmt.Errorf("Error setting issuer_dn: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_truststore_certificate", "read", "set-issuer_dn").GetDiag()
	}
	if err = d.Set("issuer_cn", trustStoreCertificateDetails.IssuerCn); err != nil {
		err = fmt.Errorf("Error setting issuer_cn: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_truststore_certificate", "read", "set-issuer_cn").GetDiag()
	}
	if err = d.Set("issued", flex.DateTimeToString(trustStoreCertificateDetails.Issued)); err != nil {
		err = fmt.Errorf("Error setting issued: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_truststore_certificate", "read", "set-issued").GetDiag()
	}
	if err = d.Set("expiry", flex.DateTimeToString(trustStoreCertificateDetails.Expiry)); err != nil {
		err = fmt.Errorf("Error setting expiry: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_truststore_certificate", "read", "set-expiry").GetDiag()
	}
	if err = d.Set("trusted", trustStoreCertificateDetails.Trusted); err != nil {
		err = fmt.Errorf("Error setting trusted: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_truststore_certificate", "read", "set-trusted").GetDiag()
	}
	if err = d.Set("href", trustStoreCertificateDetails.Href); err != nil {
		err = fmt.Errorf("Error setting href: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_truststore_certificate", "read", "set-href").GetDiag()
	}
	if err = d.Set("certificate_id", trustStoreCertificateDetails.ID); err != nil {
		err = fmt.Errorf("Error setting certificate_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_truststore_certificate", "read", "set-certificate_id").GetDiag()
	}

	return nil
}

func resourceIbmMqcloudTruststoreCertificateDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_mqcloud_truststore_certificate", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	err = checkSIPlan(d, meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Delete Truststore Certificate failed: %s", err.Error()), "ibm_mqcloud_truststore_certificate", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteTrustStoreCertificateOptions := &mqcloudv1.DeleteTrustStoreCertificateOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_truststore_certificate", "delete", "sep-id-parts").GetDiag()
	}

	deleteTrustStoreCertificateOptions.SetServiceInstanceGuid(parts[0])
	deleteTrustStoreCertificateOptions.SetQueueManagerID(parts[1])
	deleteTrustStoreCertificateOptions.SetCertificateID(parts[2])

	_, err = mqcloudClient.DeleteTrustStoreCertificateWithContext(context, deleteTrustStoreCertificateOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteTrustStoreCertificateWithContext failed: %s", err.Error()), "ibm_mqcloud_truststore_certificate", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}
