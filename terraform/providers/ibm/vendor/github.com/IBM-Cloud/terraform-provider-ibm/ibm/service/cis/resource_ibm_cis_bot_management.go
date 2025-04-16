// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMCISBotManagement() *schema.Resource {
	return &schema.Resource{
		Read:     dataSourceIBMCISBotAnalyticsRead,
		Create:   ResourceIBMCISBotManagementCreate,
		Update:   ResourceIBMCISBotManagementUpdate,
		Delete:   ResourceIBMCISBotManagementDelete,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
				ValidateFunc: validate.InvokeValidator("ibm_cis_bot_management",
					"cis_id"),
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Description:      "Associated CIS domain",
				Required:         true,
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisBotManagementFightMode: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Fight Mode",
			},
			cisBotManagementSessionScore: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Session Score",
			},
			cisBotManagementEnableJs: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Enable JS",
			},
			cisBotManagementAuthIdLogging: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Auth ID Logging",
			},
			cisBotManagementUseLatestModel: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Use Latest Model",
			},
		},
	}
}

func ResourceIBMCISBotManagementCreate(d *schema.ResourceData, meta interface{}) error {

	return ResourceIBMCISBotManagementUpdate(d, meta)
}

func ResourceIBMCISBotManagementUpdate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(conns.ClientSession).CisBotManagementSession()
	if err != nil {
		return fmt.Errorf("[ERROR] Error while Getting IAM Access Token using BluemixSession %s", err)
	}
	crn := d.Get(cisID).(string)
	sess.Crn = core.StringPtr(crn)

	cisClient, err := meta.(conns.ClientSession).CisBotManagementSession()
	if err != nil {
		return fmt.Errorf("[ERROR] Error while getting the CisBotManagementSession %s", err)
	}

	if d.HasChange(cisBotManagementFightMode) ||
		d.HasChange(cisBotManagementSessionScore) ||
		d.HasChange(cisBotManagementEnableJs) ||
		d.HasChange(cisBotManagementAuthIdLogging) ||
		d.HasChange(cisBotManagementUseLatestModel) {

		opt := cisClient.NewUpdateBotManagementOptions()

		if f, ok := d.GetOk(cisBotManagementFightMode); ok {
			fightmode := f.(bool)
			opt.SetFightMode(fightmode)
		}
		if s, ok := d.GetOk(cisBotManagementSessionScore); ok {
			sessionscore := s.(bool)
			opt.SetSessionScore(sessionscore)
		}
		if e, ok := d.GetOk(cisBotManagementEnableJs); ok {
			enablejs := e.(bool)
			opt.SetEnableJs(enablejs)
		}
		if a, ok := d.GetOk(cisBotManagementAuthIdLogging); ok {
			authidlogging := a.(bool)
			opt.SetAuthIdLogging(authidlogging)
		}
		if sl, ok := d.GetOk(cisBotManagementUseLatestModel); ok {
			uselatestmodel := sl.(bool)
			opt.SetUseLatestModel(uselatestmodel)
		}

		_, resp, err := cisClient.UpdateBotManagement(opt)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating BotManagement with error: %s %s", err, resp)
		}
	}
	return dataSourceIBMCISBotManagementRead(d, meta)
}

func ResourceIBMCISBotManagementValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})

	ibmCISBotManagementResourceValidator := validate.ResourceValidator{ResourceName: "ibm_cis_bot_management", Schema: validateSchema}
	return &ibmCISBotManagementResourceValidator
}

func ResourceIBMCISBotManagementDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
