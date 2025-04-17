// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"log"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	cissslv1 "github.com/IBM/networking-go-sdk/sslcertificateapiv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	ibmCISCertificateUpload             = "ibm_cis_certificate_upload"
	CisCertificateUploadCertificate     = "certificate"
	CisCertificateUploadPrivateKey      = "private_key"
	cisCertificateUploadBundleMethod    = "bundle_method"
	cisCertificateUploadGeoRestrictions = "geo_restrictions"
	cisCertificateUploadCustomCertID    = "custom_cert_id"
	cisCertificateUploadStatus          = "status"
	cisCertificateUploadPriority        = "priority"
	cisCertificateUploadHosts           = "hosts"
	cisCertificateUploadIssuer          = "issuer"
	cisCertificateUploadSignature       = "signature"
	cisCertificateUploadUploadedOn      = "uploaded_on"
	cisCertificateUploadModifiedOn      = "modified_on"
	cisCertificateUploadExpiresOn       = "expires_on"
	cisCertificateUploadUbiquitous      = "ubiquitous"
	cisCertificateUploadDeletePending   = "deleting"
	cisCertificateUploadDeleted         = "deleted"
)

func ResourceIBMCISCertificateUpload() *schema.Resource {
	return &schema.Resource{
		Create:   resourceCISCertificateUploadCreate,
		Read:     resourceCISCertificateUploadRead,
		Update:   resourceCISCertificateUploadUpdate,
		Delete:   resourceCISCertificateUploadDelete,
		Exists:   resourceCISCertificateUploadExists,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
				ValidateFunc: validate.InvokeValidator(ibmCISCertificateUpload,
					"cis_id"),
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Description:      "Associated CIS domain",
				Required:         true,
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisCertificateUploadCustomCertID: {
				Type:     schema.TypeString,
				Computed: true,
			},
			CisCertificateUploadCertificate: {
				Type:        schema.TypeString,
				Description: "Certificate key",
				Required:    true,
				Sensitive:   true,
			},
			CisCertificateUploadPrivateKey: {
				Type:        schema.TypeString,
				Description: "Certificate private key",
				Required:    true,
				Sensitive:   true,
			},
			cisCertificateUploadBundleMethod: {
				Type:        schema.TypeString,
				Description: "Certificate bundle method",
				Optional:    true,
				Default:     cisCertificateUploadUbiquitous,
				ValidateFunc: validate.InvokeValidator(
					ibmCISCertificateUpload,
					cisCertificateUploadBundleMethod),
			},
			cisCertificateUploadHosts: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "hosts which the certificate uploaded to",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			cisCertificateUploadPriority: {
				Type:        schema.TypeInt,
				Description: "Certificate priority",
				Optional:    true,
				Computed:    true,
			},
			cisCertificateUploadStatus: {
				Type:        schema.TypeString,
				Description: "certificate status",
				Computed:    true,
			},
			cisCertificateUploadIssuer: {
				Type:        schema.TypeString,
				Description: "certificate issuer",
				Computed:    true,
			},
			cisCertificateUploadSignature: {
				Type:        schema.TypeString,
				Description: "certificate signature",
				Computed:    true,
			},
			cisCertificateUploadUploadedOn: {
				Type:        schema.TypeString,
				Description: "certificate uploaded date",
				Computed:    true,
			},
			cisCertificateUploadModifiedOn: {
				Type:        schema.TypeString,
				Description: "certificate modified date",
				Computed:    true,
			},
			cisCertificateUploadExpiresOn: {
				Type:        schema.TypeString,
				Description: "certificate expires date",
				Computed:    true,
			},
		},
	}
}

func ResourceIBMCISCertificateUploadValidator() *validate.ResourceValidator {
	bundleMethod := "ubiquitous, optimal, force"
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
			Identifier:                 cisCertificateUploadBundleMethod,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              bundleMethod})

	cisCertificateUploadValidator := validate.ResourceValidator{ResourceName: ibmCISCertificateUpload, Schema: validateSchema}
	return &cisCertificateUploadValidator
}

func resourceCISCertificateUploadCreate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisSSLClientSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	zoneID, _, _ := flex.ConvertTftoCisTwoVar(d.Get(cisDomainID).(string))
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	// upload certificate
	opt := cisClient.NewUploadCustomCertificateOptions()
	opt.SetCertificate(d.Get(CisCertificateUploadCertificate).(string))
	opt.SetPrivateKey(d.Get(CisCertificateUploadPrivateKey).(string))
	if v, ok := d.GetOk(cisCertificateUploadBundleMethod); ok {
		opt.SetBundleMethod(v.(string))
	}

	result, response, err := cisClient.UploadCustomCertificate(opt)
	if err != nil {
		log.Printf("Upload custom certificate failed: %v", response)
		return err
	}
	certID := *result.Result.ID
	d.SetId(flex.ConvertCisToTfThreeVar(certID, zoneID, crn))

	// change priority of certificate
	certsList := []cissslv1.CertPriorityReqCertificatesItem{}
	id := certID
	var priority int64
	if v, ok := d.GetOk(cisCertificateUploadPriority); ok {
		priority = int64(v.(int))
		certsItem, _ := cisClient.NewCertPriorityReqCertificatesItem(id, priority)
		certsList = append(certsList, *certsItem)
		priorityOpt := cisClient.NewChangeCertificatePriorityOptions()
		priorityOpt.SetCertificates(certsList)
		priorityResponse, err := cisClient.ChangeCertificatePriority(priorityOpt)
		if err != nil {
			log.Printf("Change certificate priority failed: %v", priorityResponse)
			return err
		}
	}

	return resourceCISCertificateUploadRead(d, meta)
}
func resourceCISCertificateUploadRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisSSLClientSession()
	if err != nil {
		return err
	}

	certID, zoneID, crn, _ := flex.ConvertTfToCisThreeVar(d.Id())
	if err != nil {
		return err
	}
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	opt := cisClient.NewGetCustomCertificateOptions(certID)
	result, response, err := cisClient.GetCustomCertificate(opt)
	if err != nil {
		log.Printf("Get custom certificate failed: %v", response)
		return err
	}
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisCertificateUploadCustomCertID, result.Result.ID)
	d.Set(cisCertificateUploadBundleMethod, result.Result.BundleMethod)
	d.Set(cisCertificateUploadIssuer, result.Result.Issuer)
	d.Set(cisCertificateUploadHosts, flex.FlattenStringList(result.Result.Hosts))
	d.Set(cisCertificateUploadSignature, result.Result.Signature)
	d.Set(cisCertificateUploadPriority, result.Result.Priority)
	d.Set(cisCertificateUploadStatus, result.Result.Status)
	d.Set(cisCertificateUploadUploadedOn, result.Result.UploadedOn)
	d.Set(cisCertificateUploadModifiedOn, result.Result.ModifiedOn)
	d.Set(cisCertificateUploadExpiresOn, result.Result.ExpiresOn)
	return nil
}
func resourceCISCertificateUploadUpdate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisSSLClientSession()
	if err != nil {
		return err
	}

	certID, zoneID, crn, _ := flex.ConvertTfToCisThreeVar(d.Id())
	if err != nil {
		return err
	}
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	if d.HasChange(cisCertificateUploadBundleMethod) {

		opt := cisClient.NewUpdateCustomCertificateOptions(certID)
		opt.SetCertificate(d.Get(CisCertificateUploadCertificate).(string))
		opt.SetPrivateKey(d.Get(CisCertificateUploadPrivateKey).(string))
		if v, ok := d.GetOk(cisCertificateUploadBundleMethod); ok {
			opt.SetBundleMethod(v.(string))
		}
		_, response, err := cisClient.UpdateCustomCertificate(opt)
		if err != nil {
			log.Printf("Update custom certificate failed: %v", response)
			return err
		}
	}

	if d.HasChange(cisCertificateUploadPriority) {
		// change priority of certificate
		certsList := []cissslv1.CertPriorityReqCertificatesItem{}
		id := certID
		var priority int64
		if v, ok := d.GetOk(cisCertificateUploadPriority); ok {
			priority = int64(v.(int))
			certsItem, _ := cisClient.NewCertPriorityReqCertificatesItem(id, priority)
			certsList = append(certsList, *certsItem)
			priorityOpt := cisClient.NewChangeCertificatePriorityOptions()
			priorityOpt.SetCertificates(certsList)
			_, err := cisClient.ChangeCertificatePriority(priorityOpt)
			if err != nil {
				log.Printf("Change certificate priority failed: %v", err)
				return err
			}
		}
	}
	return resourceCISCertificateUploadRead(d, meta)
}

func resourceCISCertificateUploadDelete(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisSSLClientSession()
	if err != nil {
		return err
	}

	certID, zoneID, crn, _ := flex.ConvertTfToCisThreeVar(d.Id())
	if err != nil {
		return err
	}
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)
	opt := cisClient.NewDeleteCustomCertificateOptions(certID)
	_, err = cisClient.DeleteCustomCertificate(opt)
	if err != nil {
		log.Printf("Delete custom certificate failed: %v", err)
		return err
	}
	_, err = waitForCISCertificateUploadDelete(d, meta)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceCISCertificateUploadExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	cisClient, err := meta.(conns.ClientSession).CisSSLClientSession()
	if err != nil {
		return false, err
	}

	certID, zoneID, crn, _ := flex.ConvertTfToCisThreeVar(d.Id())
	if err != nil {
		return false, err
	}
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	opt := cisClient.NewGetCustomCertificateOptions(certID)
	_, detail, err := cisClient.GetCustomCertificate(opt)
	if err != nil {
		if detail != nil && strings.Contains(err.Error(), "Invalid certificate") {
			return false, nil
		}
		log.Printf("Get custom certificate failed: %v", err)
		return false, err
	}
	return true, nil
}

func waitForCISCertificateUploadDelete(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	cisClient, err := meta.(conns.ClientSession).CisSSLClientSession()
	if err != nil {
		return nil, err
	}

	certID, zoneID, crn, _ := flex.ConvertTfToCisThreeVar(d.Id())
	if err != nil {
		return nil, err
	}
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	opt := cisClient.NewGetCustomCertificateOptions(certID)
	stateConf := &resource.StateChangeConf{
		Pending: []string{cisCertificateUploadDeletePending},
		Target:  []string{cisCertificateUploadDeleted},
		Refresh: func() (interface{}, string, error) {
			_, detail, err := cisClient.GetCustomCertificate(opt)
			if err != nil {
				if detail != nil && strings.Contains(err.Error(), "Invalid certificate") {
					return detail, cisCertificateUploadDeleted, nil
				}
				return nil, "", err
			}
			return detail, cisCertificateUploadDeletePending, nil
		},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		MinTimeout:   5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	return stateConf.WaitForState()
}
