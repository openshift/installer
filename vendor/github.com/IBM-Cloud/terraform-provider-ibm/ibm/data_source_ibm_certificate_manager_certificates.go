// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataIBMCertificateManagerCertificates() *schema.Resource {
	return &schema.Resource{
		Read: dataIBMCertificateManagerCertificatesRead,
		Schema: map[string]*schema.Schema{
			"certificate_manager_instance_id": {
				Type:        schema.TypeString,
				Description: "Certificate Manager Instance ID",
				Required:    true,
			},
			"certificates": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of certificates",
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
						"domains": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"issuer": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"begins_on": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"expires_on": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"algorithm": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_algorithm": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"serial_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"imported": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"has_previous": {
							Type:     schema.TypeBool,
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
func dataIBMCertificateManagerCertificatesRead(d *schema.ResourceData, meta interface{}) error {
	cmService, err := meta.(ClientSession).CertificateManagerAPI()
	if err != nil {
		return err
	}
	instanceID := d.Get("certificate_manager_instance_id").(string)
	result, err := cmService.Certificate().ListCertificates(instanceID)
	if err != nil {
		return err
	}
	record := make([]map[string]interface{}, len(result))
	for i, c := range result {
		certificate := make(map[string]interface{})
		certificate["cert_id"] = c.ID
		certificate["name"] = c.Name
		certificate["domains"] = c.Domains
		certificate["status"] = c.Status
		certificate["issuer"] = c.Issuer
		certificate["begins_on"] = c.BeginsOn
		certificate["expires_on"] = c.ExpiresOn
		certificate["algorithm"] = c.Algorithm
		certificate["key_algorithm"] = c.KeyAlgorithm
		certificate["serial_number"] = c.SerialNumber
		certificate["imported"] = c.Imported
		certificate["has_previous"] = c.HasPrevious
		if c.IssuanceInfo != nil {
			issuanceinfo := make(map[string]interface{})
			if c.IssuanceInfo.Status != "" {
				issuanceinfo["status"] = c.IssuanceInfo.Status
			}
			if c.IssuanceInfo.Code != "" {
				issuanceinfo["code"] = c.IssuanceInfo.Code
			}
			if c.IssuanceInfo.AdditionalInfo != "" {
				issuanceinfo["additional_info"] = c.IssuanceInfo.AdditionalInfo
			}
			if c.IssuanceInfo.OrderedOn != 0 {
				order := c.IssuanceInfo.OrderedOn
				orderedOn := strconv.FormatInt(order, 10)
				issuanceinfo["ordered_on"] = orderedOn
			}
			certificate["issuance_info"] = issuanceinfo
		}
		record[i] = certificate
	}
	d.SetId(instanceID)
	d.Set("certificate_manager_instance_id", instanceID)
	d.Set("certificates", record)

	return nil
}
