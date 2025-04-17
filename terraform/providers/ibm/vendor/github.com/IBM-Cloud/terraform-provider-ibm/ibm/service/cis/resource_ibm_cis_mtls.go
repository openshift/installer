// Copyright IBM Corp. 2017, 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"context"
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	cisMtlsID           = "mtls_id"
	cisMtlsCert         = "certificate"
	cisMtlsHostNames    = "associated_hostnames"
	cisMtlsCertExpireOn = "expires_on"
)

func ResourceIBMCISMtls() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMCISMtlsCreate,
		ReadContext:   resourceIBMCISMtlsRead,
		UpdateContext: resourceIBMCISMtlsUpdate,
		DeleteContext: resourceIBMCISMtlsDelete,
		Importer:      &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
				ValidateFunc: validate.InvokeValidator("ibm_cis_mtls",
					"cis_id"),
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Description:      "Associated CIS domain",
				Required:         true,
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisMtlsID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Mtls transaction ID",
			},
			cisMtlsCert: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Certificate contents",
				Sensitive:   true,
			},
			cisMtlsCertName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Certificate name",
			},
			cisMtlsHostNames: {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "Host name list to be associated",
			},
			cisMtlsCertCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate Created At",
			},
			cisMtlsCertUpdatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate Updated At",
			},
			cisMtlsCertExpireOn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate Expires on",
			},
			cisMtlsCertID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate ID",
			},
		},
	}
}
func ResourceIBMCISMtlsValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})
	ibmCISMtlsValidator := validate.ResourceValidator{
		ResourceName: "ibm_cis_mtls",
		Schema:       validateSchema}
	return &ibmCISMtlsValidator
}
func resourceIBMCISMtlsCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).CisMtlsSession()
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error while getting the CisMtlsSession() %s %v", err, sess))
	}
	crn := d.Get(cisID).(string)
	zoneID := d.Get(cisDomainID).(string)
	sess.Crn = core.StringPtr(crn)

	options := sess.NewCreateAccessCertificateOptions(zoneID)
	if name, ok := d.GetOk(cisMtlsCertName); ok {
		options.SetName(name.(string))
	}

	if cert_val, ok := d.GetOk(cisMtlsCert); ok {
		options.SetCertificate(cert_val.(string))
	}

	if _, ok := d.GetOk(cisMtlsHostNames); ok {
		options.SetAssociatedHostnames(flex.ExpandStringList(d.Get(cisMtlsHostNames).([]interface{})))
		//options.SetAssociatedHostnames([]string{host_val.(string)})
	}

	result, resp, err := sess.CreateAccessCertificate(options)
	if err != nil || result == nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error creating MTLS access certificate %v", resp))
	}

	d.SetId(flex.ConvertCisToTfThreeVar(*result.Result.ID, zoneID, crn))
	return resourceIBMCISMtlsRead(context, d, meta)

}

func resourceIBMCISMtlsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).CisMtlsSession()
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error while getting the CisMtlsSession() %s %v", err, sess))
	}
	//crn := d.Get(cisID).(string)
	//sess.Crn = core.StringPtr(crn)

	certID, zoneID, crn, _ := flex.ConvertTfToCisThreeVar(d.Id())
	sess.Crn = core.StringPtr(crn)
	getOptions := sess.NewGetAccessCertificateOptions(zoneID, certID)
	result, response, err := sess.GetAccessCertificate(getOptions)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("[ERROR] Error While reading MTLS access certificate %v:%v:%v", err, response, result))
	}

	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisMtlsID, *result.Result.ID)
	d.Set(cisMtlsCertCreatedAt, *result.Result.CreatedAt)
	d.Set(cisMtlsCertUpdatedAt, *result.Result.UpdatedAt)
	d.Set(cisMtlsCertExpireOn, *result.Result.ExpiresOn)
	d.Set(cisMtlsCertID, *result.Result.ID)

	return nil
}

func resourceIBMCISMtlsUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).CisMtlsSession()
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error while getting the CisMtlsSession() %s %v", err, sess))
	}

	crn := d.Get(cisID).(string)
	sess.Crn = core.StringPtr(crn)

	certID, zoneID, _, _ := flex.ConvertTfToCisThreeVar(d.Id())

	if d.HasChange(cisMtlsCertName) ||
		d.HasChange(cisMtlsHostNames) {

		updateOption := sess.NewUpdateAccessCertificateOptions(zoneID, certID)
		if _, ok := d.GetOk(cisMtlsHostNames); ok {

			updateOption.SetAssociatedHostnames(flex.ExpandStringList(d.Get(cisMtlsHostNames).([]interface{})))
		}

		if name, ok := d.GetOk(cisMtlsCertName); ok {
			updateOption.SetName(name.(string))
		}

		_, updateResp, updateErr := sess.UpdateAccessCertificate(updateOption)
		if updateErr != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error while updating the MTLS cert options %v", updateResp))
		}
	}

	return resourceIBMCISMtlsRead(context, d, meta)
}

func resourceIBMCISMtlsDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).CisMtlsSession()
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error while getting the CisMtlsSession() %s %v", err, sess))
	}

	certID, zoneID, crn, _ := flex.ConvertTfToCisThreeVar(d.Id())
	//crn := d.Get(cisID).(string)
	sess.Crn = core.StringPtr(crn)
	// zoneID := d.Get(cisDomainID).(string)

	delOpt := sess.NewDeleteAccessCertificateOptions(zoneID, certID)
	_, delResp, delErr := sess.DeleteAccessCertificate(delOpt)
	if delErr != nil {

		return diag.FromErr(fmt.Errorf("[ERROR] Error While deleting the MTLS cert : %v", delResp))
	}

	return nil

}
