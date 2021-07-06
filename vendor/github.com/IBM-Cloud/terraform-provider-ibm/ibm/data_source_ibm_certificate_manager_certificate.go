// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataIBMCertificateManagerCertificate() *schema.Resource {
	return &schema.Resource{
		Read: dataIBMCertificateManagerCertificateRead,
		Schema: map[string]*schema.Schema{
			"certificate_manager_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"certificate_details": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of certificate",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cert_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"issuer": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domains": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"begins_on": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"expires_on": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"imported": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"has_previous": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"key_algorithm": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"algorithm": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"issuance_info": {
							Type:     schema.TypeMap,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
func dataIBMCertificateManagerCertificateRead(d *schema.ResourceData, meta interface{}) error {
	cmService, err := meta.(ClientSession).CertificateManagerAPI()
	if err != nil {
		return err
	}
	instanceID := d.Get("certificate_manager_instance_id").(string)
	certName := d.Get("name").(string)

	certificateList, err := cmService.Certificate().ListCertificates(instanceID)
	if err != nil {
		return err
	}
	record := make([]map[string]interface{}, 0)
	for _, cert := range certificateList {
		if certName == cert.Name {
			certificate := make(map[string]interface{})
			certificatedata, err := cmService.Certificate().GetCertData(cert.ID)
			if err != nil {
				return err
			}
			certificate["cert_id"] = certificatedata.ID
			certificate["name"] = certificatedata.Name
			certificate["domains"] = certificatedata.Domains
			certificate["description"] = certificatedata.Description
			if certificatedata.Data != nil {
				data := map[string]interface{}{
					"content": certificatedata.Data.Content,
				}
				if certificatedata.Data.Privatekey != "" {
					data["priv_key"] = certificatedata.Data.Privatekey
				}
				if certificatedata.Data.IntermediateCertificate != "" {
					data["intermediate"] = certificatedata.Data.IntermediateCertificate
				}
				certificate["data"] = data
			}
			if &certificatedata.IssuanceInfo != nil {
				issuanceinfo := map[string]interface{}{}
				if certificatedata.IssuanceInfo.Status != "" {
					issuanceinfo["status"] = certificatedata.IssuanceInfo.Status
				}
				if certificatedata.IssuanceInfo.Code != "" {
					issuanceinfo["code"] = certificatedata.IssuanceInfo.Code
				}
				if certificatedata.IssuanceInfo.AdditionalInfo != "" {
					issuanceinfo["additional_info"] = certificatedata.IssuanceInfo.AdditionalInfo
				}
				if certificatedata.IssuanceInfo.OrderedOn != 0 {
					order := certificatedata.IssuanceInfo.OrderedOn
					orderedOn := strconv.FormatInt(order, 10)
					issuanceinfo["ordered_on"] = orderedOn
				}
				certificate["issuance_info"] = issuanceinfo
			}
			certificate["status"] = certificatedata.Status
			certificate["issuer"] = certificatedata.Issuer
			certificate["imported"] = certificatedata.Imported
			certificate["has_previous"] = certificatedata.HasPrevious
			certificate["key_algorithm"] = certificatedata.KeyAlgorithm
			certificate["algorithm"] = certificatedata.Algorithm
			certificate["begins_on"] = certificatedata.BeginsOn
			certificate["expires_on"] = certificatedata.ExpiresOn

			record = append(record, certificate)
			d.Set("certificate_details", record)
		}
	}
	d.SetId(fmt.Sprintf("%s:%s", certName, instanceID))
	d.Set("certificate_manager_instance_id", instanceID)
	d.Set("name", certName)
	return nil
}
