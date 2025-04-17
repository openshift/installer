// Copyright IBM Corp. 2017, 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0
package cis

import (
	"context"
	"fmt"

	"strings"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/networking-go-sdk/authenticatedoriginpullapiv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	cisOriginAuthID          = "auth_id"
	cisOriginAuthHost        = "hostname"
	cisOriginAuthEnable      = "enabled"
	cisOriginAuthLevel       = "level"
	cisOriginAuthCertContent = "certificate"
	cisOriginAuthCertKey     = "private_key"
	cisOriginAuthCertId      = "cert_id"
	CisOriginAuthStatus      = "status"
	cisOriginAuthExpiresOn   = "expires_on"
	cisOriginAuthUploadedOn  = "uploaded_on"
)

func ResourceIBMCISOriginAuthPull() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMCISOriginAuthPullCreate,
		ReadContext:   resourceIBMCISOriginAuthPullRead,
		UpdateContext: resourceIBMCISOriginAuthPullUpdate,
		DeleteContext: resourceIBMCISOriginAuthPullDelete,
		Importer:      &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
				ValidateFunc: validate.InvokeValidator("ibm_cis_origin_auth",
					"cis_id"),
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Description:      "Associated CIS domain",
				Required:         true,
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisOriginAuthLevel: {
				Type:        schema.TypeString,
				Description: "Origin auth level zone or hostname",
				Required:    true,
			},
			cisOriginAuthHost: {
				Type:        schema.TypeString,
				Description: "Host name needed for host level authentication",
				Optional:    true,
			},
			cisOriginAuthEnable: {
				Type:        schema.TypeBool,
				Description: "Enabel-disable origin auth for a zone or host",
				Optional:    true,
				Default:     true,
			},
			cisOriginAuthCertContent: {
				Type:        schema.TypeString,
				Description: "Certificate content which needs to be uploaded",
				Required:    true,
				Sensitive:   true,
			},
			cisOriginAuthCertKey: {
				Type:        schema.TypeString,
				Description: "Private key content which needs to be uploaded",
				Required:    true,
				Sensitive:   true,
			},
			CisOriginAuthStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Authentication status whether active or not",
			},
			cisOriginAuthCertId: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate ID which is uploaded",
			},
			cisOriginAuthExpiresOn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate expires on",
			},
			cisOriginAuthUploadedOn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate uploaded on",
			},
			cisOriginAuthID: {
				Type:        schema.TypeString,
				Description: "Associated CIS auth pull job id",
				Computed:    true,
			},
		},
	}

}
func ResourceIBMCISOriginAuthPullValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})
	ibmCISOriginAuthPullValidator := validate.ResourceValidator{
		ResourceName: "ibm_cis_origin_auth",
		Schema:       validateSchema}
	return &ibmCISOriginAuthPullValidator
}

func resourceIBMCISOriginAuthPullCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var cert_val string
	var key_val string
	var level_val string
	var zone_config bool

	sess, err := meta.(conns.ClientSession).CisOrigAuthSession()
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error while getting the CisOrigAuthSession %v", err))
	}

	crn := d.Get(cisID).(string)
	zoneID, _, _ := flex.ConvertTftoCisTwoVar(d.Get(cisDomainID).(string))

	sess.ZoneIdentifier = core.StringPtr(zoneID)
	sess.Crn = core.StringPtr(crn)

	if cert_content, ok := d.GetOk(cisOriginAuthCertContent); ok {
		cert_val = cert_content.(string)

	}

	if cert_key, ok := d.GetOk(cisOriginAuthCertKey); ok {
		key_val = cert_key.(string)

	}
	zone_config = true
	if lev_val, ok := d.GetOk(cisOriginAuthLevel); ok {
		level_val = lev_val.(string)
		if strings.ToLower(level_val) != "zone" {
			zone_config = false
		}
	}

	// Check host level certificate creation or zone level
	if zone_config {
		options := sess.NewUploadZoneOriginPullCertificateOptions()
		options.SetCertificate(cert_val)
		options.SetPrivateKey(key_val)

		result, resp, opErr := sess.UploadZoneOriginPullCertificate(options)
		if opErr != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error while uploading certificate zone level %v", resp))
		}

		d.SetId(flex.ConvertCisToTfFourVar(*result.Result.ID, level_val, zoneID, crn))

	} else {
		options := sess.NewUploadHostnameOriginPullCertificateOptions()
		options.SetCertificate(cert_val)
		options.SetPrivateKey(key_val)
		result, resp, opErr := sess.UploadHostnameOriginPullCertificate(options)
		if opErr != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error while uploading certificate host level %v", resp))
		}

		d.SetId(flex.ConvertCisToTfFourVar(*result.Result.ID, level_val, zoneID, crn))

	}

	return resourceIBMCISOriginAuthPullRead(context, d, meta)
}

func resourceIBMCISOriginAuthPullRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var zone_config bool
	sess, err := meta.(conns.ClientSession).CisOrigAuthSession()
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error while getting the CisOrigAuthSession %v", err))
	}

	certID, level_val, zoneID, crn, _ := flex.ConvertTfToCisFourVar(d.Id())
	sess.Crn = core.StringPtr(crn)
	sess.ZoneIdentifier = core.StringPtr(zoneID)

	zone_config = true
	if strings.ToLower(level_val) != "zone" {
		zone_config = false
	}

	if zone_config {
		getOptions := sess.NewGetZoneOriginPullCertificateOptions(certID)
		getOptions.SetCertIdentifier(certID)

		result, response, err := sess.GetZoneOriginPullCertificate(getOptions)

		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error while getting detail of zone origin auth pull %v:%v", err, response))
		}
		d.Set(cisOriginAuthID, *result.Result.ID)
		d.Set(cisOriginAuthCertContent, *result.Result.Certificate)
		d.Set(CisOriginAuthStatus, *result.Result.Status)
		d.Set(cisOriginAuthExpiresOn, *result.Result.ExpiresOn)
		d.Set(cisOriginAuthUploadedOn, *result.Result.UploadedOn)
		d.Set(cisOriginAuthCertId, *result.Result.ID)

	} else {
		getOptions := sess.NewGetHostnameOriginPullCertificateOptions(certID)
		getOptions.SetCertIdentifier(certID)

		result, response, err := sess.GetHostnameOriginPullCertificate(getOptions)

		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error while getting detail of host origin auth pull %v:%v", err, response))
		}
		d.Set(cisOriginAuthID, *result.Result.ID)
		d.Set(cisOriginAuthCertContent, *result.Result.Certificate)
		d.Set(CisOriginAuthStatus, *result.Result.Status)
		d.Set(cisOriginAuthExpiresOn, *result.Result.ExpiresOn)
		d.Set(cisOriginAuthUploadedOn, *result.Result.UploadedOn)
		d.Set(cisOriginAuthCertId, *result.Result.ID)
	}

	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)

	return nil
}

func resourceIBMCISOriginAuthPullUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var zone_config bool
	var host_name string
	sess, err := meta.(conns.ClientSession).CisOrigAuthSession()
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error while getting the CisOrigAuthSession %v", err))
	}

	certID, level_val, zoneID, crn, _ := flex.ConvertTfToCisFourVar(d.Id())
	sess.Crn = core.StringPtr(crn)
	sess.ZoneIdentifier = core.StringPtr(zoneID)

	zone_config = true
	if strings.ToLower(level_val) != "zone" {
		zone_config = false
	}

	if zone_config {

		if d.HasChange(cisOriginAuthEnable) {
			updateOption := sess.NewSetZoneOriginPullSettingsOptions()
			updateOption.SetEnabled(d.Get(cisOriginAuthEnable).(bool))
			_, response, err := sess.SetZoneOriginPullSettings(updateOption)

			if err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error while updaing the zone origin auth pull setting %v:%v", err, response))
			}

		}

	} else {

		if d.HasChange(cisOriginAuthEnable) {
			if host_val, ok := d.GetOk(cisOriginAuthHost); ok {
				host_name = host_val.(string)
			}

			model := &authenticatedoriginpullapiv1.HostnameOriginPullSettings{
				Hostname: core.StringPtr(host_name),
				CertID:   core.StringPtr(certID),
				Enabled:  core.BoolPtr(d.Get(cisOriginAuthEnable).(bool)),
			}
			setOption := sess.NewSetHostnameOriginPullSettingsOptions()
			setOption.SetConfig([]authenticatedoriginpullapiv1.HostnameOriginPullSettings{*model})
			_, setResp, setErr := sess.SetHostnameOriginPullSettings(setOption)
			if setErr != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error while updaing the host origin auth pull setting %v:%v", setErr, setResp))
			}

		}

	}
	return resourceIBMCISOriginAuthPullRead(context, d, meta)

}

func resourceIBMCISOriginAuthPullDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var zone_config bool
	sess, err := meta.(conns.ClientSession).CisOrigAuthSession()
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error while getting the CisOrigAuthSession %v", err))
	}

	certID, level_val, zoneID, crn, _ := flex.ConvertTfToCisFourVar(d.Id())
	sess.Crn = core.StringPtr(crn)
	sess.ZoneIdentifier = core.StringPtr(zoneID)

	zone_config = true
	if strings.ToLower(level_val) != "zone" {
		zone_config = false
	}

	if zone_config {
		delOpt := sess.NewDeleteZoneOriginPullCertificateOptions(certID)
		_, resp, err := sess.DeleteZoneOriginPullCertificate(delOpt)

		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error while deleting the certificate zone level %v: %v", certID, resp))
		}

	} else {
		delOpt := sess.NewDeleteHostnameOriginPullCertificateOptions(certID)
		_, resp, err := sess.DeleteHostnameOriginPullCertificate(delOpt)

		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error while deleting the certificate host level %v: %v", certID, resp))
		}

	}
	d.SetId("")
	return nil

}
