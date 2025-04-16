// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	cisCertificates                   = "certificates"
	cisCertificatesCertificates       = "certificates"
	cisCertificatesCertificatesID     = "id"
	cisCertificatesCertificatesHosts  = "hosts"
	cisCertificatesCertificatesStatus = "status"
	cisCertificatesPrimaryCertificate = "primary_certificate"
	cisCertificatesType               = "type"
	cisCertificateTypeDedicated       = "dedicated"
)

func DataSourceIBMCISCertificates() *schema.Resource {
	return &schema.Resource{
		Read: dataIBMCISCertificatesRead,
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
				ValidateFunc: validate.InvokeDataSourceValidator(
					"ibm_cis_certificates",
					"cis_id"),
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Description:      "Associated CIS domain",
				Required:         true,
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisCertificates: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of certificates",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "certificate identifier",
						},
						cisCertificateOrderID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "certificate id",
						},
						cisCertificatesType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "certificate type",
						},
						cisCertificateOrderHosts: {
							Type:        schema.TypeList,
							Description: "Hosts which certificate need to be ordered",
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						cisCertificateOrderStatus: {
							Type:        schema.TypeString,
							Description: "certificate status",
							Computed:    true,
						},
						cisCertificatesPrimaryCertificate: {
							Type:        schema.TypeString,
							Description: "Primary certificate id",
							Computed:    true,
						},
						cisCertificatesCertificates: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Collection of certificates associated with this certificates",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									cisCertificatesCertificatesID: {
										Type:        schema.TypeString,
										Description: "certificate id",
										Computed:    true,
									},
									cisCertificatesCertificatesHosts: {
										Type:        schema.TypeList,
										Description: "Hosts which certificates are ordered",
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									cisCertificatesCertificatesStatus: {
										Type:        schema.TypeString,
										Description: "certificate status",
										Computed:    true,
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
func DataSourceIBMCISCertificatesValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})

	iBMCISCertificatesValidator := validate.ResourceValidator{
		ResourceName: "ibm_cis_certificates",
		Schema:       validateSchema}
	return &iBMCISCertificatesValidator
}
func dataIBMCISCertificatesRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisSSLClientSession()
	if err != nil {
		return err
	}
	crn := d.Get(cisID).(string)
	zoneID, _, _ := flex.ConvertTftoCisTwoVar(d.Get(cisDomainID).(string))
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)
	opt := cisClient.NewListCertificatesOptions()
	result, response, err := cisClient.ListCertificates(opt)
	if err != nil {
		log.Printf("List all certificates failed: %v", response)
		return err
	}
	certificatesList := make([]interface{}, 0)
	for _, instance := range result.Result {
		certificate := map[string]interface{}{}
		certificate["id"] = flex.ConvertCisToTfThreeVar(*instance.ID, zoneID, crn)
		certificate[cisCertificateOrderID] = *instance.ID
		certificate[cisCertificateOrderStatus] = *instance.Status
		if instance.PrimaryCertificate != nil {
			certificate[cisCertificatesPrimaryCertificate] =
				convertCISCertificatesObj(*instance.Type, instance.PrimaryCertificate)
		}
		certificate[cisCertificateOrderHosts] = flex.FlattenStringList(instance.Hosts)

		certs := []interface{}{}
		for _, i := range instance.Certificates {
			cert := map[string]interface{}{}
			if i.ID != nil {
				cert[cisCertificatesCertificatesID] = convertCISCertificatesObj(*instance.Type, i.ID)
			}
			cert[cisCertificatesCertificatesStatus] = *i.Status
			cert[cisCertificatesCertificatesHosts] = flex.FlattenStringList(i.Hosts)
			certs = append(certs, cert)
		}
		certificate[cisCertificatesType] = *instance.Type
		certificate[cisCertificatesCertificates] = certs
		certificatesList = append(certificatesList, certificate)
	}
	d.SetId(dataSourceIBMCISCertificatesID(d))
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisCertificates, certificatesList)
	return nil
}

func dataSourceIBMCISCertificatesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func convertCISCertificatesObj(certType string, obj interface{}) (result string) {
	if certType == cisCertificateTypeDedicated {
		result = strings.TrimSpace(fmt.Sprintf("%32.f", obj))
	} else {
		result = fmt.Sprint(obj)
	}
	return result
}
