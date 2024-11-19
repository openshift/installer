// Copyright IBM Corp. 2017, 2021, 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"log"
	"reflect"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMCISOriginCertificateOrder() *schema.Resource {
	return &schema.Resource{
		Read:     DataIBMCISOriginCertificateRead,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS object ID or CRN",
				Required:    true,
				ValidateFunc: validate.InvokeValidator(ibmCISOriginCertificateOrder,
					"cis_id"),
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Description:      "Associated CIS domain",
				Required:         true,
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisOriginCertificateID: {
				Type:        schema.TypeString,
				Description: "Certificate ID",
				Optional:    true,
			},
			cisOriginCertificateList: {
				Type:        schema.TypeList,
				Description: "List of certificate",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						cisOriginCertificateID: {
							Type:        schema.TypeString,
							Description: "Certificate ID",
							Computed:    true,
						},
						cisOriginCertificateType: {
							Type:        schema.TypeString,
							Description: "Certificate type",
							Computed:    true,
						},
						cisOriginCertificateHosts: {
							Type:        schema.TypeList,
							Description: "Hosts for which certificates need to be ordered",
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						cisOriginCertificateValidityDays: {
							Type:        schema.TypeInt,
							Description: "Validity days",
							Computed:    true,
						},
						cisOriginCertificateCSR: {
							Type:        schema.TypeString,
							Description: "CSR",
							Computed:    true,
						},
						cisOriginCertificatePrivateKey: {
							Type:        schema.TypeString,
							Description: "Certificate private key",
							Computed:    true,
						},
						cisOriginCertificate: {
							Type:        schema.TypeString,
							Description: "Certificate",
							Computed:    true,
						},
						cisOriginCertificateExpiresOn: {
							Type:        schema.TypeString,
							Description: "Expiration date of the certificate",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func DataIBMCISOriginCertificateOrderValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "data_instance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})

	cisOriginCertificateOrderValidator := validate.ResourceValidator{
		ResourceName: ibmCISOriginCertificateOrder,
		Schema:       validateSchema}
	return &cisOriginCertificateOrderValidator
}

func DataIBMCISOriginCertificateRead(d *schema.ResourceData, meta interface{}) error {

	cisClient, err := meta.(conns.ClientSession).CisSSLClientSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	zoneID := d.Get(cisDomainID).(string)
	cert_id := d.Get(cisOriginCertificateID).(string)

	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	originCertList := make([]map[string]interface{}, 0)

	if cert_id != "" {

		opt := cisClient.NewGetOriginCertificateOptions(crn, zoneID, cert_id)
		result, resp, err := cisClient.GetOriginCertificate(opt)
		if err != nil {
			log.Printf("Get Certificate read failed: %v", resp)
			return err
		}

		certOutput := map[string]interface{}{}
		certOutput[cisOriginCertificateID] = *result.Result.ID
		if !reflect.ValueOf(result.Result.RequestType).IsNil() {
			certOutput[cisOriginCertificateType] = *result.Result.RequestType
		}
		certOutput[cisOriginCertificateHosts] = flex.FlattenStringList(result.Result.Hostnames)
		if !reflect.ValueOf(result.Result.RequestedValidity).IsNil() {
			certOutput[cisOriginCertificateValidityDays] = *result.Result.RequestedValidity
		}
		if !reflect.ValueOf(result.Result.Csr).IsNil() {
			certOutput[cisOriginCertificateCSR] = *result.Result.Csr
		}
		if !reflect.ValueOf(result.Result.PrivateKey).IsNil() {
			certOutput[cisOriginCertificatePrivateKey] = *result.Result.PrivateKey
		}
		certOutput[cisOriginCertificate] = *result.Result.Certificate
		certOutput[cisOriginCertificateExpiresOn] = *result.Result.ExpiresOn

		originCertList = append(originCertList, certOutput)

	} else {
		opt := cisClient.NewListOriginCertificatesOptions(crn, zoneID)
		result, resp, err := cisClient.ListOriginCertificates(opt)
		if err != nil {
			log.Printf("List Certificate read failed: %v", resp)
			return err
		}
		for _, certObj := range result.Result {
			certOutput := map[string]interface{}{}
			certOutput[cisOriginCertificateID] = *certObj.ID
			if !reflect.ValueOf(certObj.RequestType).IsNil() {
				certOutput[cisOriginCertificateType] = *certObj.RequestType
			}
			certOutput[cisOriginCertificateHosts] = flex.FlattenStringList(certObj.Hostnames)
			if !reflect.ValueOf(certObj.RequestedValidity).IsNil() {
				certOutput[cisOriginCertificateValidityDays] = *certObj.RequestedValidity
			}
			if !reflect.ValueOf(certObj.Csr).IsNil() {
				certOutput[cisOriginCertificateCSR] = *certObj.Csr
			}
			if !reflect.ValueOf(certObj.PrivateKey).IsNil() {
				certOutput[cisOriginCertificatePrivateKey] = *certObj.PrivateKey
			}
			certOutput[cisOriginCertificate] = *certObj.Certificate
			certOutput[cisOriginCertificateExpiresOn] = *certObj.ExpiresOn

			originCertList = append(originCertList, certOutput)
		}
	}

	d.SetId(dataSourceIBMCISOriginCertificatesID())
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisOriginCertificateList, originCertList)

	return nil
}

func dataSourceIBMCISOriginCertificatesID() string {
	return time.Now().UTC().String()
}
