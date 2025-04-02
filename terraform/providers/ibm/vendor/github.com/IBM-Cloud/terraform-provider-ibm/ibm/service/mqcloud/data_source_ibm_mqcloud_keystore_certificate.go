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

func DataSourceIbmMqcloudKeystoreCertificate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmMqcloudKeystoreCertificateRead,

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
						"config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The configuration details for this certificate.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ams": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "A list of channels that are configured with this certificate.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"channels": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "A list of channels that are configured with this certificate.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The name of the channel.",
															},
														},
													},
												},
											},
										},
									},
								},
							},
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
		// Error is coming from SDK client, so it doesn't need to be discriminated.
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_mqcloud_keystore_certificate", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	err = checkSIPlan(d, meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Read Keystore Certificate failed: %s", err.Error()), "(Data) ibm_mqcloud_keystore_certificate", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	listKeyStoreCertificatesOptions := &mqcloudv1.ListKeyStoreCertificatesOptions{}

	listKeyStoreCertificatesOptions.SetServiceInstanceGuid(d.Get("service_instance_guid").(string))
	listKeyStoreCertificatesOptions.SetQueueManagerID(d.Get("queue_manager_id").(string))

	keyStoreCertificateDetailsCollection, _, err := mqcloudClient.ListKeyStoreCertificatesWithContext(context, listKeyStoreCertificatesOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListKeyStoreCertificatesWithContext failed: %s", err.Error()), "(Data) ibm_mqcloud_keystore_certificate", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
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
			return flex.DiscriminatedTerraformErrorf(nil, fmt.Sprintf("no KeyStore found with label %s", label), "(Data) ibm_mqcloud_keystore_certificate", "read", "no-collection-found").GetDiag()
		}
		d.SetId(label)
	} else {
		d.SetId(dataSourceIbmMqcloudKeystoreCertificateID(d))
	}

	if err = d.Set("total_count", flex.IntValue(keyStoreCertificateDetailsCollection.TotalCount)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting total_count: %s", err), "(Data) ibm_mqcloud_keystore_certificate", "read", "set-total_count").GetDiag()
	}

	keyStore := []map[string]interface{}{}
	if keyStoreCertificateDetailsCollection.KeyStore != nil {
		for _, modelItem := range keyStoreCertificateDetailsCollection.KeyStore {
			modelItem := modelItem
			modelMap, err := DataSourceIbmMqcloudKeystoreCertificateKeyStoreCertificateDetailsToMap(&modelItem)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_mqcloud_keystore_certificate", "read", "key_store-to-map").GetDiag()
			}
			keyStore = append(keyStore, modelMap)
		}
	}
	if err = d.Set("key_store", keyStore); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting key_store: %s", err), "(Data) ibm_mqcloud_keystore_certificate", "read", "set-key_store").GetDiag()
	}

	return nil
}

// dataSourceIbmMqcloudKeystoreCertificateID returns a reasonable ID for the list.
func dataSourceIbmMqcloudKeystoreCertificateID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIbmMqcloudKeystoreCertificateKeyStoreCertificateDetailsToMap(model *mqcloudv1.KeyStoreCertificateDetails) (map[string]interface{}, error) {
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
	modelMap["is_default"] = *model.IsDefault
	modelMap["dns_names_total_count"] = flex.IntValue(model.DnsNamesTotalCount)
	modelMap["dns_names"] = model.DnsNames
	modelMap["href"] = *model.Href
	configMap, err := DataSourceIbmMqcloudKeystoreCertificateCertificateConfigurationToMap(model.Config)
	if err != nil {
		return modelMap, err
	}
	modelMap["config"] = []map[string]interface{}{configMap}
	return modelMap, nil
}

func DataSourceIbmMqcloudKeystoreCertificateCertificateConfigurationToMap(model *mqcloudv1.CertificateConfiguration) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	amsMap, err := DataSourceIbmMqcloudKeystoreCertificateChannelsDetailsToMap(model.Ams)
	if err != nil {
		return modelMap, err
	}
	modelMap["ams"] = []map[string]interface{}{amsMap}
	return modelMap, nil
}

func DataSourceIbmMqcloudKeystoreCertificateChannelsDetailsToMap(model *mqcloudv1.ChannelsDetails) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	channels := []map[string]interface{}{}
	for _, channelsItem := range model.Channels {
		channelsItem := channelsItem
		channelsItemMap, err := DataSourceIbmMqcloudKeystoreCertificateChannelDetailsToMap(&channelsItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		channels = append(channels, channelsItemMap)
	}
	modelMap["channels"] = channels
	return modelMap, nil
}

func DataSourceIbmMqcloudKeystoreCertificateChannelDetailsToMap(model *mqcloudv1.ChannelDetails) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	return modelMap, nil
}
