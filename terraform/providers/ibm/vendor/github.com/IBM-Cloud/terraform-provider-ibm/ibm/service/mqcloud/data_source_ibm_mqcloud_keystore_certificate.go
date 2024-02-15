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

func DataSourceIbmMqcloudKeystoreCertificate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmMqcloudKeystoreCertificateRead,

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
				Description: "The total count of key store certificates.",
			},
			"key_store": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of key store certificates.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the certificate.",
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
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this key store certificate.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmMqcloudKeystoreCertificateRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		return diag.FromErr(err)
	}
	err = checkSIPlan(d, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Read Keystore Certificate failed %s", err))
	}

	listKeyStoreCertificatesOptions := &mqcloudv1.ListKeyStoreCertificatesOptions{}

	listKeyStoreCertificatesOptions.SetServiceInstanceGuid(d.Get("service_instance_guid").(string))
	listKeyStoreCertificatesOptions.SetQueueManagerID(d.Get("queue_manager_id").(string))

	keyStoreCertificateDetailsCollection, response, err := mqcloudClient.ListKeyStoreCertificatesWithContext(context, listKeyStoreCertificatesOptions)
	if err != nil {
		log.Printf("[DEBUG] ListKeyStoreCertificatesWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ListKeyStoreCertificatesWithContext failed %s\n%s", err, response))
	}

	// Use the provided filter argument and construct a new list with only the requested resource(s)
	var matchKeyStore []mqcloudv1.KeyStoreCertificateDetails
	var label string
	var suppliedFilter bool

	if v, ok := d.GetOk("label"); ok {
		label = v.(string)
		suppliedFilter = true
		for _, data := range keyStoreCertificateDetailsCollection.KeyStore {
			if *data.Label == label {
				matchKeyStore = append(matchKeyStore, data)
			}
		}
	} else {
		matchKeyStore = keyStoreCertificateDetailsCollection.KeyStore
	}
	keyStoreCertificateDetailsCollection.KeyStore = matchKeyStore

	if suppliedFilter {
		if len(keyStoreCertificateDetailsCollection.KeyStore) == 0 {
			return diag.FromErr(fmt.Errorf("No Key Store Certificate found with label: \"%s\"", label))
		}
		d.SetId(label)
	} else {
		d.SetId(dataSourceIbmMqcloudKeystoreCertificateID(d))
	}

	if err = d.Set("total_count", flex.IntValue(keyStoreCertificateDetailsCollection.TotalCount)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting total_count: %s", err))
	}

	keyStore := []map[string]interface{}{}
	if keyStoreCertificateDetailsCollection.KeyStore != nil {
		for _, modelItem := range keyStoreCertificateDetailsCollection.KeyStore {
			modelItem := modelItem
			modelMap, err := dataSourceIbmMqcloudKeystoreCertificateKeyStoreCertificateDetailsToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			keyStore = append(keyStore, modelMap)
		}
	}
	if err = d.Set("key_store", keyStore); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting key_store: %s", err))
	}

	return nil
}

// dataSourceIbmMqcloudKeystoreCertificateID returns a reasonable ID for the list.
func dataSourceIbmMqcloudKeystoreCertificateID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIbmMqcloudKeystoreCertificateKeyStoreCertificateDetailsToMap(model *mqcloudv1.KeyStoreCertificateDetails) (map[string]interface{}, error) {
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
	modelMap["is_default"] = model.IsDefault
	modelMap["dns_names_total_count"] = flex.IntValue(model.DnsNamesTotalCount)
	modelMap["dns_names"] = model.DnsNames
	modelMap["href"] = model.Href
	return modelMap, nil
}
