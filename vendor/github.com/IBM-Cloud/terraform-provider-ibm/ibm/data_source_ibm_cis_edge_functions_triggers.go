// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"time"

	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const cisEdgeFunctionsTriggers = "cis_edge_functions_triggers"

func dataSourceIBMCISEdgeFunctionsTriggers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMCISEdgeFunctionsTriggerRead,
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CIS Intance CRN",
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

func dataSourceIBMCISEdgeFunctionsTriggerRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisEdgeFunctionClientSession()
	if err != nil {
		return err
	}
	crn := d.Get(cisID).(string)
	zoneID, _, err := convertTftoCisTwoVar(d.Get(cisDomainID).(string))
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	opt := cisClient.NewListEdgeFunctionsTriggersOptions()
	result, _, err := cisClient.ListEdgeFunctionsTriggers(opt)
	if err != nil {
		return fmt.Errorf("Error listing edge functions triggers: %v", err)
	}
	triggerInfo := make([]map[string]interface{}, 0)
	for _, trigger := range result.Result {
		l := map[string]interface{}{}
		l["id"] = convertCisToTfThreeVar(*trigger.ID, zoneID, crn)
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
