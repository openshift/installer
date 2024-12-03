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
	ibmCISAdvancedCertificatePackOrder            = "ibm_cis_advanced_certificate_pack_order"
	cisAdvancedCertificatePackOrderID             = "certificate_id"
	cisAdvancedCertificatePackOrderHosts          = "hosts"
	cisAdvancedCertificatePackOrderType           = "type"
	cisAdvancedCertificatePackOrderTypeDedicated  = "dedicated"
	cisAdvancedCertificatePackOrderStatus         = "status"
	cisAdvancedCertificatePackValidationMethod    = "validation_method"
	cisAdvancedCertificatePackValidityDays        = "validity"
	cisAdvancedCertificatePackCertificateAthority = "certificate_authority"
	cisAdvancedCertificatePackCloudflareBranding  = "cloudflare_branding"
	cisAdvancedCertificatePackOrderTypeAdvanced   = "advanced"
	cisOriginCertificateList                      = "origin_certificate_list"
)

func ResourceIBMCISAdvancedCertificatePackOrder() *schema.Resource {
	return &schema.Resource{
		Create:   ResourceIBMCISAdvancedCertificatePackOrderCreate,
		Update:   ResourceIBMCISAdvancedCertificatePackOrderRead,
		Read:     ResourceIBMCISAdvancedCertificatePackOrderRead,
		Delete:   ResourceIBMCISAdvancedCertificatePackOrderDelete,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS object ID or CRN",
				Required:    true,
				ValidateFunc: validate.InvokeValidator(ibmCISAdvancedCertificatePackOrder,
					"cis_id"),
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Description:      "Associated CIS domain",
				Required:         true,
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisAdvancedCertificatePackOrderID: {
				Type:        schema.TypeString,
				Description: "Certificate ID",
				Computed:    true,
			},
			cisAdvancedCertificatePackOrderType: {
				Type:        schema.TypeString,
				Description: "Certificate type",
				Optional:    true,
				Default:     cisAdvancedCertificatePackOrderTypeAdvanced,
				ValidateFunc: validate.InvokeValidator(ibmCISAdvancedCertificatePackOrder,
					cisAdvancedCertificatePackOrderType),
			},
			cisAdvancedCertificatePackOrderHosts: {
				Type:        schema.TypeList,
				Description: "Hosts for which certificates need to be ordered",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			cisAdvancedCertificatePackOrderStatus: {
				Type:        schema.TypeString,
				Description: "Certificate status",
				Computed:    true,
			},
			cisAdvancedCertificatePackValidationMethod: {
				Type:        schema.TypeString,
				Description: "Validation method",
				Required:    true,
			},
			cisAdvancedCertificatePackValidityDays: {
				Type:        schema.TypeInt,
				Description: "Validity days",
				Required:    true,
			},
			cisAdvancedCertificatePackCertificateAthority: {
				Type:        schema.TypeString,
				Description: "Certificate authority",
				Required:    true,
			},
			cisAdvancedCertificatePackCloudflareBranding: {
				Type:        schema.TypeBool,
				Description: "Cloudflare branding",
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func ResourceIBMCISAdvancedCertificatePackOrderValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 cisCertificateOrderType,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              cisAdvancedCertificatePackOrderTypeAdvanced})

	cisCertificateOrderValidator := validate.ResourceValidator{
		ResourceName: ibmCISAdvancedCertificatePackOrder,
		Schema:       validateSchema}
	return &cisCertificateOrderValidator
}

func ResourceIBMCISAdvancedCertificatePackOrderCreate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisSSLClientSession()
	if err != nil {
		return err
	}
	crn := d.Get(cisID).(string)
	zoneID := d.Get(cisDomainID).(string)
	certType := d.Get(cisAdvancedCertificatePackOrderType).(string)
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	hosts := d.Get(cisAdvancedCertificatePackOrderHosts)
	hostsList := flex.ExpandStringList(hosts.([]interface{}))
	validationMethod := d.Get(cisAdvancedCertificatePackValidationMethod).(string)
	validityDays := int64(d.Get(cisAdvancedCertificatePackValidityDays).(int))
	certificateAuthority := d.Get(cisAdvancedCertificatePackCertificateAthority).(string)
	cfBranding := d.Get(cisAdvancedCertificatePackCloudflareBranding).(bool)

	opt := cisClient.NewOrderAdvancedCertificateOptions()
	opt.SetType(certType)
	opt.SetHosts(hostsList)
	opt.SetValidationMethod(validationMethod)
	opt.SetValidityDays(validityDays)
	opt.SetCertificateAuthority(certificateAuthority)
	opt.SetCloudflareBranding(cfBranding)

	result, resp, err := cisClient.OrderAdvancedCertificate(opt)
	if err != nil {
		log.Printf("Advanced Certificate Pack order failed: %v", resp)
		return err
	}

	d.SetId(flex.ConvertCisToTfThreeVar(*result.Result.ID, zoneID, crn))
	d.Set(cisAdvancedCertificatePackOrderID, *result.Result.ID)
	d.Set(cisAdvancedCertificatePackOrderStatus, *result.Result.Status)

	return nil
}

func ResourceIBMCISAdvancedCertificatePackOrderRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func ResourceIBMCISAdvancedCertificatePackOrderDelete(d *schema.ResourceData, meta interface{}) error {
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
	opt := cisClient.NewDeleteCertificateV2Options(certificateID)
	resp, err := cisClient.DeleteCertificateV2(opt)
	if err != nil {
		log.Printf("Advanced Certificate Pack delete failed: %v", resp)
		return err
	}

	return nil
}
