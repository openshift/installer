// Copyright IBM Corp. 2017, 2021, 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	ibmCISOriginCertificateOrder     = "ibm_cis_origin_certificate_order"
	cisOriginCertificate             = "certificate"
	cisOriginCertificateID           = "certificate_id"
	cisOriginCertificateHosts        = "hostnames"
	cisOriginCertificateType         = "request_type"
	cisOriginCertificateValidityDays = "requested_validity"
	cisOriginCertificateCSR          = "csr"
	cisOriginCertificateExpiresOn    = "expires_on"
	cisOriginCertificatePrivateKey   = "private_key"
)

func ResourceIBMCISOriginCertificateOrder() *schema.Resource {
	return &schema.Resource{
		Create:   ResourceIBMCISOriginCertificateCreate,
		Update:   ResourceIBMCISOriginCertificateRead,
		Read:     ResourceIBMCISOriginCertificateRead,
		Delete:   ResourceIBMCISOriginCertificateDelete,
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
				Computed:    true,
			},
			cisOriginCertificateType: {
				Type:        schema.TypeString,
				Description: "Certificate type",
				Required:    true,
			},
			cisOriginCertificateHosts: {
				Type:        schema.TypeList,
				Description: "Hosts for which certificates need to be ordered",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			cisOriginCertificateValidityDays: {
				Type:        schema.TypeInt,
				Description: "Calidity days",
				Required:    true,
			},
			cisOriginCertificateCSR: {
				Type:        schema.TypeString,
				Description: "CSR",
				Required:    true,
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
	}
}

func ResourceIBMCISOriginCertificateOrderValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})

	cisCertificateOrderValidator := validate.ResourceValidator{
		ResourceName: ibmCISOriginCertificateOrder,
		Schema:       validateSchema}
	return &cisCertificateOrderValidator
}

func ResourceIBMCISOriginCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisSSLClientSession()
	if err != nil {
		return err
	}
	crn := d.Get(cisID).(string)
	zoneID := d.Get(cisDomainID).(string)
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	certType := d.Get(cisOriginCertificateType).(string)
	hosts := d.Get(cisOriginCertificateHosts)
	hostsList := flex.ExpandStringList(hosts.([]interface{}))
	validityDays := int64(d.Get(cisOriginCertificateValidityDays).(int))
	csr := d.Get(cisOriginCertificateCSR).(string)

	opt := cisClient.NewCreateOriginCertificateOptions(crn, zoneID)
	opt.SetHostnames(hostsList)
	opt.SetCsr(csr)
	opt.SetRequestType(certType)
	opt.SetRequestedValidity(validityDays)

	result, resp, err := cisClient.CreateOriginCertificate(opt)
	if err != nil {
		log.Printf("Origin Certificate order failed: %v", resp)
		return err
	}

	d.SetId(flex.ConvertCisToTfThreeVar(*result.Result.ID, zoneID, crn))
	return ResourceIBMCISOriginCertificateRead(d, meta)
}

func ResourceIBMCISOriginCertificateRead(d *schema.ResourceData, meta interface{}) error {

	cisClient, err := meta.(conns.ClientSession).CisSSLClientSession()
	if err != nil {
		return err
	}
	certificateID, zoneID, crn, err := flex.ConvertTfToCisThreeVar(d.Id())
	if err != nil {
		log.Println("Error in reading certificate ID")
		return err
	}
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)
	opt := cisClient.NewGetOriginCertificateOptions(crn, zoneID, certificateID)
	result, resp, err := cisClient.GetOriginCertificate(opt)
	if err != nil {
		log.Printf("Certificate read failed: %v", resp)
		return err
	}
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisOriginCertificateID, result.Result.ID)
	d.Set(cisOriginCertificate, result.Result.Certificate)
	d.Set(cisOriginCertificateHosts, flex.FlattenStringList(result.Result.Hostnames))
	d.Set(cisOriginCertificateExpiresOn, result.Result.ExpiresOn)
	d.Set(cisOriginCertificateType, result.Result.RequestType)
	d.Set(cisOriginCertificateValidityDays, result.Result.RequestedValidity)
	d.Set(cisOriginCertificateCSR, result.Result.Csr)
	d.Set(cisOriginCertificatePrivateKey, result.Result.PrivateKey)
	return nil
}

func ResourceIBMCISOriginCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisSSLClientSession()
	if err != nil {
		return err
	}
	certificateID, zoneID, crn, err := flex.ConvertTfToCisThreeVar(d.Id())
	if err != nil {
		log.Println("Error in reading certificate Id")
		return err
	}
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)
	opt := cisClient.NewRevokeOriginCertificateOptions(crn, zoneID, certificateID)
	resp, _, err := cisClient.RevokeOriginCertificate(opt)
	if err != nil {
		log.Printf("Origin Certificate delete failed: %v", resp)
		return err
	}

	return nil
}
