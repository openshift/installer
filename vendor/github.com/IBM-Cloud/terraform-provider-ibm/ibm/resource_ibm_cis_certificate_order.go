// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"log"
	"time"

	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

func resourceIBMCISCertificateOrder() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMCISCertificateOrderCreate,
		Update:   resourceIBMCISCertificateOrderRead,
		Read:     resourceIBMCISCertificateOrderRead,
		Delete:   resourceIBMCISCertificateOrderDelete,
		Exists:   resourceIBMCISCertificateOrderExist,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS object id or CRN",
				Required:    true,
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

func resourceIBMCISCertificateOrderValidator() *ResourceValidator {
	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisCertificateOrderType,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              cisCertificateOrderTypeDedicated})

	cisCertificateOrderValidator := ResourceValidator{
		ResourceName: ibmCISCertificateOrder,
		Schema:       validateSchema}
	return &cisCertificateOrderValidator
}

func resourceIBMCISCertificateOrderCreate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisSSLClientSession()
	if err != nil {
		return err
	}
	crn := d.Get(cisID).(string)
	zoneID, _, err := convertTftoCisTwoVar(d.Get(cisDomainID).(string))
	certType := d.Get(cisCertificateOrderType).(string)
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	hosts := d.Get(cisCertificateOrderHosts)
	hostsList := expandStringList(hosts.([]interface{}))
	opt := cisClient.NewOrderCertificateOptions()
	opt.SetType(certType)
	opt.SetHosts(hostsList)

	result, resp, err := cisClient.OrderCertificate(opt)
	if err != nil {
		log.Printf("Certificate order failed: %v", resp)
		return err
	}

	d.SetId(convertCisToTfThreeVar(*result.Result.ID, zoneID, crn))
	return resourceIBMCISCertificateOrderRead(d, meta)
}

func resourceIBMCISCertificateOrderRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisSSLClientSession()
	if err != nil {
		return err
	}
	certificateID, zoneID, crn, err := convertTfToCisThreeVar(d.Id())
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
	d.Set(cisCertificateOrderHosts, flattenStringList(result.Result.Hosts))
	d.Set(cisCertificateOrderStatus, result.Result.Status)
	return nil
}

func resourceIBMCISCertificateOrderDelete(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisSSLClientSession()
	if err != nil {
		return err
	}
	certificateID, zoneID, crn, err := convertTfToCisThreeVar(d.Id())
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

func resourceIBMCISCertificateOrderExist(d *schema.ResourceData, meta interface{}) (bool, error) {
	cisClient, err := meta.(ClientSession).CisSSLClientSession()
	if err != nil {
		return false, err
	}
	certificateID, zoneID, crn, err := convertTfToCisThreeVar(d.Id())
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
	cisClient, err := meta.(ClientSession).CisSSLClientSession()
	if err != nil {
		return nil, err
	}
	certificateID, zoneID, crn, err := convertTfToCisThreeVar(d.Id())
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
