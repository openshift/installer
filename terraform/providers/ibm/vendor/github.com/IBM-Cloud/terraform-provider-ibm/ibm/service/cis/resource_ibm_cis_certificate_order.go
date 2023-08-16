// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	ibmCISCertificateOrder           = "ibm_cis_certificate_order"
	cisCertificateOrderID            = "certificate_id"
	cisCertificateOrderHosts         = "hosts"
	cisCertificateOrderType          = "type"
	cisCertificateOrderTypeDedicated = "dedicated"
	cisCertificateOrderStatus        = "status"
	cisCertificateOrderDeleted       = "deleted"
	cisCertificateOrderDeletePending = "deleting"
)

func ResourceIBMCISCertificateOrder() *schema.Resource {
	return &schema.Resource{
		Create:   ResourceIBMCISCertificateOrderCreate,
		Update:   ResourceIBMCISCertificateOrderRead,
		Read:     ResourceIBMCISCertificateOrderRead,
		Delete:   ResourceIBMCISCertificateOrderDelete,
		Exists:   ResourceIBMCISCertificateOrderExist,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS object id or CRN",
				Required:    true,
				ValidateFunc: validate.InvokeValidator(ibmCISCertificateOrder,
					"cis_id"),
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Description:      "Associated CIS domain",
				Required:         true,
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisCertificateOrderID: {
				Type:        schema.TypeString,
				Description: "certificate id",
				Computed:    true,
			},
			cisCertificateOrderType: {
				Type:        schema.TypeString,
				Description: "certificate type",
				Optional:    true,
				Default:     cisCertificateOrderTypeDedicated,
				ValidateFunc: validate.InvokeValidator(ibmCISCertificateOrder,
					cisCertificateOrderType),
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
		},
	}
}

func ResourceIBMCISCertificateOrderValidator() *validate.ResourceValidator {
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
			AllowedValues:              cisCertificateOrderTypeDedicated})

	cisCertificateOrderValidator := validate.ResourceValidator{
		ResourceName: ibmCISCertificateOrder,
		Schema:       validateSchema}
	return &cisCertificateOrderValidator
}

func ResourceIBMCISCertificateOrderCreate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisSSLClientSession()
	if err != nil {
		return err
	}
	crn := d.Get(cisID).(string)
	zoneID, _, _ := flex.ConvertTftoCisTwoVar(d.Get(cisDomainID).(string))
	certType := d.Get(cisCertificateOrderType).(string)
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	hosts := d.Get(cisCertificateOrderHosts)
	hostsList := flex.ExpandStringList(hosts.([]interface{}))
	opt := cisClient.NewOrderCertificateOptions()
	opt.SetType(certType)
	opt.SetHosts(hostsList)

	result, resp, err := cisClient.OrderCertificate(opt)
	if err != nil {
		log.Printf("Certificate order failed: %v", resp)
		return err
	}

	d.SetId(flex.ConvertCisToTfThreeVar(*result.Result.ID, zoneID, crn))
	return ResourceIBMCISCertificateOrderRead(d, meta)
}

func ResourceIBMCISCertificateOrderRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisSSLClientSession()
	if err != nil {
		return err
	}
	certificateID, zoneID, crn, _ := flex.ConvertTfToCisThreeVar(d.Id())
	if err != nil {
		log.Println("Error in reading certificate id")
		return err
	}
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)
	opt := cisClient.NewGetCustomCertificateOptions(certificateID)
	result, resp, err := cisClient.GetCustomCertificate(opt)
	if err != nil {
		log.Printf("Certificate read failed: %v", resp)
		return err
	}
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisCertificateOrderID, result.Result.ID)
	d.Set(cisCertificateOrderType, cisCertificateOrderTypeDedicated)
	d.Set(cisCertificateOrderHosts, flex.FlattenStringList(result.Result.Hosts))
	d.Set(cisCertificateOrderStatus, result.Result.Status)
	return nil
}

func ResourceIBMCISCertificateOrderDelete(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisSSLClientSession()
	if err != nil {
		return err
	}
	certificateID, zoneID, crn, _ := flex.ConvertTfToCisThreeVar(d.Id())
	if err != nil {
		log.Println("Error in reading certificate id")
		return err
	}
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)
	opt := cisClient.NewDeleteCertificateOptions(certificateID)
	resp, err := cisClient.DeleteCertificate(opt)
	if err != nil {
		log.Printf("Certificate delete failed: %v", resp)
		return err
	}

	_, err = waitForCISCertificateOrderDelete(d, meta)
	if err != nil {
		return err
	}

	return nil
}

func ResourceIBMCISCertificateOrderExist(d *schema.ResourceData, meta interface{}) (bool, error) {
	cisClient, err := meta.(conns.ClientSession).CisSSLClientSession()
	if err != nil {
		return false, err
	}
	certificateID, zoneID, crn, _ := flex.ConvertTfToCisThreeVar(d.Id())
	if err != nil {
		log.Println("Error in reading certificate id")
		return false, err
	}
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)
	opt := cisClient.NewGetCustomCertificateOptions(certificateID)
	_, response, err := cisClient.GetCustomCertificate(opt)
	if err != nil {
		if response != nil && response.StatusCode == 400 {
			log.Printf("Certificate is not found")
			return false, nil
		}
		log.Printf("Get Certificate failed: %v", response)
		return false, err
	}
	return true, nil
}

func waitForCISCertificateOrderDelete(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	cisClient, err := meta.(conns.ClientSession).CisSSLClientSession()
	if err != nil {
		return nil, err
	}
	certificateID, zoneID, crn, _ := flex.ConvertTfToCisThreeVar(d.Id())
	if err != nil {
		log.Println("Error in reading certificate id")
		return nil, err
	}
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)
	opt := cisClient.NewGetCustomCertificateOptions(certificateID)
	stateConf := &resource.StateChangeConf{
		Pending: []string{cisCertificateOrderDeletePending},
		Target:  []string{cisCertificateOrderDeleted},
		Refresh: func() (interface{}, string, error) {
			_, detail, err := cisClient.GetCustomCertificate(opt)
			if err != nil {
				if detail != nil && detail.StatusCode == 400 {
					return detail, cisCertificateOrderDeleted, nil
				}
				return nil, "", err
			}
			return detail, cisCertificateOrderDeletePending, nil
		},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		MinTimeout:   10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	return stateConf.WaitForState()
}
