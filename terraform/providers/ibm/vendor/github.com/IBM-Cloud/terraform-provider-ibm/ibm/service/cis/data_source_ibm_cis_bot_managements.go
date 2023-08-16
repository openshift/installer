// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	cisBotManagementFightMode      = "fight_mode"
	cisBotManagementSessionScore   = "session_score"
	cisBotManagementEnableJs       = "enable_js"
	cisBotManagementAuthIdLogging  = "auth_id_logging"
	cisBotManagementUseLatestModel = "use_latest_model"
)

func DataSourceIBMCISBotManagement() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMCISBotManagementRead,
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
				ValidateFunc: validate.InvokeDataSourceValidator(
					"ibm_cis_bot_managements",
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

func DataSourceIBMCISBotManagementValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})

	iBMCISBotManagementValidator := validate.ResourceValidator{
		ResourceName: "ibm_cis_bot_managements",
		Schema:       validateSchema}
	return &iBMCISBotManagementValidator
}

func dataSourceIBMCISBotManagementRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisBotManagementSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	zoneName := d.Get(cisDomainID).(string)
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneName)
	opt := cisClient.NewGetBotManagementOptions()

	result, resp, err := cisClient.GetBotManagement(opt)
	if err != nil {
		log.Printf("dataSourceIBMCISBotManagementRead - GetBotManagement Failed %s\n", resp)
		return err
	}

	res := result.Result
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneName)
	d.Set(cisBotManagementFightMode, res.FightMode)
	d.Set(cisBotManagementSessionScore, res.SessionScore)
	d.Set(cisBotManagementEnableJs, res.EnableJs)
	d.Set(cisBotManagementAuthIdLogging, res.AuthIdLogging)
	d.Set(cisBotManagementUseLatestModel, res.UseLatestModel)

	return nil
}
