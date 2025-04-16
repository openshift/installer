// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
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
	ibmCISCustomPage             = "ibm_cis_custom_page"
	cisCustomPageIdentifier      = "page_id"
	cisCustomPageURL             = "url"
	cisCustomPageState           = "state"
	cisCustomPageStateDefault    = "default"
	cisCustomPageStateCustomized = "customized"
	cisCustomPageDesc            = "description"
	cisCustomPageRequiredTokens  = "required_tokens"
	cisCustomPagePreviewTarget   = "preview_target"
	cisCustomPageCreatedOn       = "created_on"
	cisCustomPageModifiedOn      = "modified_on"
)

func ResourceIBMCISCustomPage() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
				ValidateFunc: validate.InvokeValidator(ibmCISCustomPage,
					"cis_id"),
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Description:      "Associated CIS domain",
				Required:         true,
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisCustomPageIdentifier: {
				Type:        schema.TypeString,
				Description: "Custom page identifier",
				ForceNew:    true,
				Required:    true,
				ValidateFunc: validate.InvokeValidator(ibmCISCustomPage,
					cisCustomPageIdentifier),
			},
			cisCustomPageURL: {
				Type:        schema.TypeString,
				Description: "Custom page url",
				Required:    true,
			},
			cisCustomPageState: {
				Type:        schema.TypeString,
				Description: "Custom page state",
				Computed:    true,
			},
			cisCustomPageDesc: {
				Type:        schema.TypeString,
				Description: "Free text",
				Computed:    true,
			},
			cisCustomPageRequiredTokens: {
				Type:        schema.TypeList,
				Description: "Custom page state",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			cisCustomPagePreviewTarget: {
				Type:        schema.TypeString,
				Description: "Custom page preview target",
				Computed:    true,
			},
			cisCustomPageCreatedOn: {
				Type:        schema.TypeString,
				Description: "Custom page created date",
				Computed:    true,
			},
			cisCustomPageModifiedOn: {
				Type:        schema.TypeString,
				Description: "Custom page modified date",
				Computed:    true,
			},
		},
		Create:   resourceCISCustomPageUpdate,
		Read:     resourceCISCustomPageRead,
		Update:   resourceCISCustomPageUpdate,
		Delete:   resourceCISCustomPageDelete,
		Importer: &schema.ResourceImporter{},
	}
}

func ResourceIBMCISCustomPageValidator() *validate.ResourceValidator {
	customPageIDs := "basic_challenge, waf_challenge, waf_block, ratelimit_block," +
		"country_challenge, ip_block, under_attack, 500_errors, 1000_errors, always_online"
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
			Identifier:                 cisCustomPageIdentifier,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              customPageIDs})
	ibmCISCustomPageResourceValidator := validate.ResourceValidator{
		ResourceName: ibmCISCustomPage,
		Schema:       validateSchema}
	return &ibmCISCustomPageResourceValidator
}

func resourceCISCustomPageUpdate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisCustomPageClientSession()
	if err != nil {
		return err
	}
	crn := d.Get(cisID).(string)
	zoneID, _, _ := flex.ConvertTftoCisTwoVar(d.Get(cisDomainID).(string))
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)
	pageID := d.Get(cisCustomPageIdentifier).(string)

	if d.HasChange(cisCustomPageURL) {

		url := d.Get(cisCustomPageURL).(string)
		state := cisCustomPageStateDefault
		if len(url) > 0 {
			state = cisCustomPageStateCustomized
		}
		opt := cisClient.NewUpdateZoneCustomPageOptions(pageID)
		opt.SetURL(url)
		opt.SetState(state)

		result, response, err := cisClient.UpdateZoneCustomPage(opt)
		if err != nil {
			log.Printf("Update custom page failed : %v", response)
			return err
		}
		d.SetId(flex.ConvertCisToTfThreeVar(*result.Result.ID, zoneID, crn))
	}
	return resourceCISCustomPageRead(d, meta)
}

func resourceCISCustomPageRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisCustomPageClientSession()
	if err != nil {
		return err
	}
	pageID, zoneID, crn, _ := flex.ConvertTfToCisThreeVar(d.Id())
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	opt := cisClient.NewGetZoneCustomPageOptions(pageID)

	result, response, err := cisClient.GetZoneCustomPage(opt)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			log.Printf("Custom Page has some error: %v", response)
			d.SetId("")
			return nil
		}
		log.Printf("Get custom page failed : %v", response)
		return err
	}

	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisCustomPageIdentifier, result.Result.ID)
	d.Set(cisCustomPageURL, result.Result.URL)
	d.Set(cisCustomPageState, result.Result.State)
	d.Set(cisCustomPageRequiredTokens, flex.FlattenStringList(result.Result.RequiredTokens))
	d.Set(cisCustomPageDesc, result.Result.Description)
	d.Set(cisCustomPagePreviewTarget, result.Result.PreviewTarget)
	d.Set(cisCustomPageCreatedOn, (*result.Result.CreatedOn).String())
	d.Set(cisCustomPageModifiedOn, (*result.Result.ModifiedOn).String())
	return nil
}

func resourceCISCustomPageDelete(d *schema.ResourceData, meta interface{}) error {
	// Nothing to delete on CIS resource
	d.SetId("")
	return nil
}
