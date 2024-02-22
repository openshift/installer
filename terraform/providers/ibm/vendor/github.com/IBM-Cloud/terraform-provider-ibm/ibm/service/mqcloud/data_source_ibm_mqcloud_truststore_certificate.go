// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

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
				Description: "The GUID that uniquely identifies the MQ on Cloud service instance.",
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
		return diag.FromErr(err)
	}
	err = checkSIPlan(d, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Read Truststore Certificate failed %s", err))
	}

	listTrustStoreCertificatesOptions := &mqcloudv1.ListTrustStoreCertificatesOptions{}

	listTrustStoreCertificatesOptions.SetServiceInstanceGuid(d.Get("service_instance_guid").(string))
	listTrustStoreCertificatesOptions.SetQueueManagerID(d.Get("queue_manager_id").(string))

	trustStoreCertificateDetailsCollection, response, err := mqcloudClient.ListTrustStoreCertificatesWithContext(context, listTrustStoreCertificatesOptions)
	if err != nil {
		log.Printf("[DEBUG] ListTrustStoreCertificatesWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ListTrustStoreCertificatesWithContext failed %s\n%s", err, response))
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
			return diag.FromErr(fmt.Errorf("No Trust Store Certificate found with label: \"%s\"", label))
		}
		d.SetId(label)
	} else {
		d.SetId(dataSourceIbmMqcloudTruststoreCertificateID(d))
	}

	if err = d.Set("total_count", flex.IntValue(trustStoreCertificateDetailsCollection.TotalCount)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting total_count: %s", err))
	}

	trustStore := []map[string]interface{}{}
	if trustStoreCertificateDetailsCollection.TrustStore != nil {
		for _, modelItem := range trustStoreCertificateDetailsCollection.TrustStore {
			modelItem := modelItem
			modelMap, err := dataSourceIbmMqcloudTruststoreCertificateTrustStoreCertificateDetailsToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			trustStore = append(trustStore, modelMap)
		}
	}
	if err = d.Set("trust_store", trustStore); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting trust_store: %s", err))
	}

	return nil
}

// dataSourceIbmMqcloudTruststoreCertificateID returns a reasonable ID for the list.
func dataSourceIbmMqcloudTruststoreCertificateID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIbmMqcloudTruststoreCertificateTrustStoreCertificateDetailsToMap(model *mqcloudv1.TrustStoreCertificateDetails) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	modelMap["label"] = model.Label
	modelMap["certificate_type"] = model.CertificateType
	modelMap["fingerprint_sha256"] = model.FingerprintSha256
	modelMap["subject_dn"] = model.SubjectDn
	modelMap["subject_cn"] = model.SubjectCn
	modelMap["issuer_dn"] = model.IssuerDn
	modelMap["issuer_cn"] = model.IssuerCn
	modelMap["issued"] = model.Issued.String()
	modelMap["expiry"] = model.Expiry.String()
	modelMap["trusted"] = model.Trusted
	modelMap["href"] = model.Href
	return modelMap, nil
}
