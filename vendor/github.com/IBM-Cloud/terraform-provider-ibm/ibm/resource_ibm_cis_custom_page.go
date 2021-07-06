// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"log"

	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

func resourceIBMCISCustomPage() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
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
				ValidateFunc: InvokeValidator(ibmCISCustomPage,
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

func resourceIBMCISCustomPageValidator() *ResourceValidator {
	customPageIDs := "basic_challenge, waf_challenge, waf_block, ratelimit_block," +
		"country_challenge, ip_block, under_attack, 500_errors, 1000_errors, always_online"
	validateSchema := make([]ValidateSchema, 1)

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisCustomPageIdentifier,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              customPageIDs})
	ibmCISCustomPageResourceValidator := ResourceValidator{
		ResourceName: ibmCISCustomPage,
		Schema:       validateSchema}
	return &ibmCISCustomPageResourceValidator
}

func resourceCISCustomPageUpdate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisCustomPageClientSession()
	if err != nil {
		return err
	}
	crn := d.Get(cisID).(string)
	zoneID, _, err := convertTftoCisTwoVar(d.Get(cisDomainID).(string))
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
		d.SetId(convertCisToTfThreeVar(*result.Result.ID, zoneID, crn))
	}
	return resourceCISCustomPageRead(d, meta)
}

func resourceCISCustomPageRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisCustomPageClientSession()
	if err != nil {
		return err
	}
	pageID, zoneID, crn, _ := convertTfToCisThreeVar(d.Id())
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
	d.Set(cisCustomPageRequiredTokens, flattenStringList(result.Result.RequiredTokens))
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
