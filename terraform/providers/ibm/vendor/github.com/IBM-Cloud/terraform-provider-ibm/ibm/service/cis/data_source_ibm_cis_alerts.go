// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/networking-go-sdk/alertsv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	cisAlerts = "alert_policies"
)

func DataSourceIBMCISAlert() *schema.Resource {
	return &schema.Resource{
		Read: dataIBMCISAlertPolicyRead,
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
			},
			cisAlerts: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Container for response information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						cisAlertID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy ID",
						},
						cisAlertName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy name",
						},
						cisAlertDescription: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy Description",
						},
						cisAlertEnabled: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Is the alert policy active",
						},
						cisAlertType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Condition for the alert",
						},
						cisAlertMechanisms: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Delivery mechanisms for the alert, can include an email, a webhook, or both.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									cisAlertEmail: {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Set: schema.HashString,
									},
									cisAlertWebhook: {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Set: schema.HashString,
									},
								},
							},
						},
						cisAlertFilters: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Filters based on filter type",
						},
						cisAlertConditions: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Conditions based on filter type",
						},
					},
				},
			},
		},
	}
}
func dataIBMCISAlertPolicyRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(conns.ClientSession).CisAlertsSession()
	if err != nil {
		return err
	}
	crn := d.Get(cisID).(string)
	sess.Crn = core.StringPtr(crn)

	opt := sess.NewGetAlertPoliciesOptions()
	result, resp, err := sess.GetAlertPolicies(opt)
	if err != nil {
		log.Printf("[WARN] List all alerts failed: %v\n", resp)
		return err
	}
	alertList := make([]map[string]interface{}, 0)
	for _, alertObj := range result.Result {
		alertOutput := map[string]interface{}{}
		alertOutput[cisAlertID] = *alertObj.ID
		alertOutput[cisAlertName] = *alertObj.Name
		alertOutput[cisAlertDescription] = *alertObj.Description
		alertOutput[cisAlertEnabled] = *alertObj.Enabled
		alertOutput[cisAlertType] = *alertObj.AlertType
		filterOpt, err := json.Marshal(alertObj.Filters)
		if err != nil {
			return fmt.Errorf("[ERROR] Error marshalling the created filters: %s", err)
		}
		alertOutput[cisAlertFilters] = string(filterOpt)
		conditionsOpt, err := json.Marshal(alertObj.Conditions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error marshalling the created Conditions: %s", err)
		}
		alertOutput[cisAlertConditions] = string(conditionsOpt)
		alertOutput[cisAlertMechanisms] = dataflattenCISMechanism(*alertObj.Mechanisms)
		alertList = append(alertList, alertOutput)

	}
	d.SetId(dataSourceCISAlertsCheckID(d))
	d.Set(cisID, crn)

	d.Set(cisAlerts, alertList)
	return nil
}
func dataSourceCISAlertsCheckID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
func dataflattenCISMechanism(Mechanism alertsv1.ListAlertPoliciesRespResultItemMechanisms) interface{} {
	emailoutput := []string{}
	webhookoutput := []string{}

	output := map[string]interface{}{}
	flatten := []map[string]interface{}{}

	for _, mech := range Mechanism.Email {
		emailoutput = append(emailoutput, *mech.ID)
	}

	for _, mech := range Mechanism.Webhooks {
		webhookoutput = append(webhookoutput, *mech.ID)
	}

	output[cisAlertEmail] = emailoutput
	output[cisAlertWebhook] = webhookoutput

	flatten = append(flatten, output)

	return flatten
}
