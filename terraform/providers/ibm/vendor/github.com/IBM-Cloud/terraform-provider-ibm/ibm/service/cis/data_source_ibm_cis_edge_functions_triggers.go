// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"fmt"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const cisEdgeFunctionsTriggers = "cis_edge_functions_triggers"

func DataSourceIBMCISEdgeFunctionsTriggers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMCISEdgeFunctionsTriggerRead,
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CIS Intance CRN",
				ValidateFunc: validate.InvokeDataSourceValidator(
					"ibm_cis_edge_functions_actions",
					"cis_id"),
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "CIS Domain ID",
				DiffSuppressFunc: suppressDataDiff,
			},
			cisEdgeFunctionsTriggers: {
				Type:        schema.TypeList,
				Description: "List of edge functions triggers",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Edge function trigger id",
						},
						cisEdgeFunctionsTriggerID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Edge function trigger route id",
						},
						cisEdgeFunctionsTriggerPattern: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Edge function trigger pattern",
						},
						cisEdgeFunctionsTriggerActionName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Edge function trigger script name",
						},
						cisEdgeFunctionsTriggerRequestLimitFailOpen: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Edge function trigger request limit fail open",
						},
					},
				},
			},
		},
	}
}

func DataSourceIBMCISEdgeFunctionsTriggersValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})

	iBMCISEdgeFunctionsTriggersValidator := validate.ResourceValidator{
		ResourceName: "ibm_cis_edge_functions_triggers",
		Schema:       validateSchema}
	return &iBMCISEdgeFunctionsTriggersValidator
}

func dataSourceIBMCISEdgeFunctionsTriggerRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisEdgeFunctionClientSession()
	if err != nil {
		return err
	}
	crn := d.Get(cisID).(string)
	zoneID, _, _ := flex.ConvertTftoCisTwoVar(d.Get(cisDomainID).(string))
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	opt := cisClient.NewListEdgeFunctionsTriggersOptions()
	result, _, err := cisClient.ListEdgeFunctionsTriggers(opt)
	if err != nil {
		return fmt.Errorf("[ERROR] Error listing edge functions triggers: %v", err)
	}
	triggerInfo := make([]map[string]interface{}, 0)
	for _, trigger := range result.Result {
		l := map[string]interface{}{}
		l["id"] = flex.ConvertCisToTfThreeVar(*trigger.ID, zoneID, crn)
		l[cisEdgeFunctionsTriggerID] = *trigger.ID
		l[cisEdgeFunctionsTriggerPattern] = *trigger.Pattern
		l[cisEdgeFunctionsTriggerRequestLimitFailOpen] = *trigger.RequestLimitFailOpen
		if trigger.Script != nil {
			l[cisEdgeFunctionsTriggerActionName] = *trigger.Script
		}
		triggerInfo = append(triggerInfo, l)
	}
	d.SetId(dataSourceIBMCISEdgeFunctionsTriggersID(d))
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisEdgeFunctionsTriggers, triggerInfo)
	return nil
}

func dataSourceIBMCISEdgeFunctionsTriggersID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
