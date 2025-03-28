// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.95.2-120e65bc-20240924-152329
 */

package mqcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/mqcloud-go-sdk/mqcloudv1"
)

func DataSourceIbmMqcloudTruststoreCertificate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmMqcloudTruststoreCertificateRead,

		Schema: map[string]*schema.Schema{
			"service_instance_guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The GUID that uniquely identifies the MQaaS service instance.",
			},
			"queue_manager_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the queue manager to retrieve its full details.",
			},
			"label": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Certificate label in queue manager store.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total count of trust store certificates.",
			},
			"trust_store": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of trust store certificates.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the certificate.",
						},
						"label": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate label in queue manager store.",
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
					},
				},
			},
		},
	}
}

func dataSourceIbmMqcloudTruststoreCertificateRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		// Error is coming from SDK client, so it doesn't need to be discriminated.
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_mqcloud_truststore_certificate", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	err = checkSIPlan(d, meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Read Truststore Certificate failed: %s", err.Error()), "(Data) ibm_mqcloud_truststore_certificate", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	listTrustStoreCertificatesOptions := &mqcloudv1.ListTrustStoreCertificatesOptions{}

	listTrustStoreCertificatesOptions.SetServiceInstanceGuid(d.Get("service_instance_guid").(string))
	listTrustStoreCertificatesOptions.SetQueueManagerID(d.Get("queue_manager_id").(string))

	trustStoreCertificateDetailsCollection, _, err := mqcloudClient.ListTrustStoreCertificatesWithContext(context, listTrustStoreCertificatesOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListTrustStoreCertificatesWithContext failed: %s", err.Error()), "(Data) ibm_mqcloud_truststore_certificate", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	// Use the provided filter argument and construct a new list with only the requested resource(s)
	var matchTrustStore []mqcloudv1.TrustStoreCertificateDetails
	var label string
	var suppliedFilter bool

	if v, ok := d.GetOk("label"); ok {
		label = v.(string)
		suppliedFilter = true
		for _, data := range trustStoreCertificateDetailsCollection.TrustStore {
			if *data.Label == label {
				matchTrustStore = append(matchTrustStore, data)
			}
		}
	} else {
		matchTrustStore = trustStoreCertificateDetailsCollection.TrustStore
	}
	trustStoreCertificateDetailsCollection.TrustStore = matchTrustStore

	if suppliedFilter {
		if len(trustStoreCertificateDetailsCollection.TrustStore) == 0 {
			return flex.DiscriminatedTerraformErrorf(nil, fmt.Sprintf("no TrustStore found with label %s", label), "(Data) ibm_mqcloud_truststore_certificate", "read", "no-collection-found").GetDiag()
		}
		d.SetId(label)
	} else {
		d.SetId(dataSourceIbmMqcloudTruststoreCertificateID(d))
	}

	if err = d.Set("total_count", flex.IntValue(trustStoreCertificateDetailsCollection.TotalCount)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting total_count: %s", err), "(Data) ibm_mqcloud_truststore_certificate", "read", "set-total_count").GetDiag()
	}
	trustStore := []map[string]interface{}{}
	if trustStoreCertificateDetailsCollection.TrustStore != nil {
		for _, modelItem := range trustStoreCertificateDetailsCollection.TrustStore {
			modelItem := modelItem
			modelMap, err := DataSourceIbmMqcloudTruststoreCertificateTrustStoreCertificateDetailsToMap(&modelItem)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_mqcloud_truststore_certificate", "read", "trust_store-to-map").GetDiag()
			}
			trustStore = append(trustStore, modelMap)
		}
	}
	if err = d.Set("trust_store", trustStore); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting trust_store: %s", err), "(Data) ibm_mqcloud_truststore_certificate", "read", "set-trust_store").GetDiag()
	}

	return nil
}

// dataSourceIbmMqcloudTruststoreCertificateID returns a reasonable ID for the list.
func dataSourceIbmMqcloudTruststoreCertificateID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIbmMqcloudTruststoreCertificateTrustStoreCertificateDetailsToMap(model *mqcloudv1.TrustStoreCertificateDetails) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	modelMap["label"] = *model.Label
	modelMap["certificate_type"] = *model.CertificateType
	modelMap["fingerprint_sha256"] = *model.FingerprintSha256
	modelMap["subject_dn"] = *model.SubjectDn
	modelMap["subject_cn"] = *model.SubjectCn
	modelMap["issuer_dn"] = *model.IssuerDn
	modelMap["issuer_cn"] = *model.IssuerCn
	modelMap["issued"] = model.Issued.String()
	modelMap["expiry"] = model.Expiry.String()
	modelMap["trusted"] = *model.Trusted
	modelMap["href"] = *model.Href
	return modelMap, nil
}
