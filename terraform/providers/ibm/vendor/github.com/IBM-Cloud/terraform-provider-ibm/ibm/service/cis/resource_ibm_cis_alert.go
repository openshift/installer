// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/networking-go-sdk/alertsv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	ibmCISAlert           = "ibm_cis_alert"
	cisAlertID            = "policy_id"
	cisAlertName          = "name"
	cisAlertDescription   = "description"
	cisAlertEnabled       = "enabled"
	cisAlertType          = "alert_type"
	cisAlertMechanisms    = "mechanisms"
	cisAlertEmail         = "email"
	cisAlertEmailID       = "email_id"
	cisAlertWebhook       = "webhooks"
	cisAlertConditions    = "conditions"
	cisAlertFilters       = "filters"
	cisAlertFilterEnabled = "enabled"
	cisAlertFilterPoolID  = "pool_id"
	cisAlertType1         = "dos_attack_l7"
	cisAlertType2         = "g6_pool_toggle_alert"
	cisAlertType3         = "clickhouse_alert_fw_anomaly"
	cisAlertType4         = "clickhouse_alert_fw_ent_anomaly"
)

func ResourceIBMCISAlert() *schema.Resource {
	return &schema.Resource{
		Create:   ResourceIBMCISAlertPolicyCreate,
		Read:     ResourceIBMCISAlertPolicyRead,
		Update:   ResourceIBMCISAlertPolicyUpdate,
		Delete:   ResourceIBMCISAlertPolicyDelete,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
			},
			cisAlertID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identifier of the Alert Policy",
			},
			cisAlertName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Policy name",
			},
			cisAlertDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Policy Description",
			},
			cisAlertEnabled: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Is the alert policy active",
			},
			cisAlertType: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Condition for the alert",
			},
			cisAlertMechanisms: {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Delivery mechanisms for the alert, can include an email, a webhook, or both.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						cisAlertEmail: {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Set: schema.HashString,
						},
						cisAlertWebhook: {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Set: schema.HashString,
						},
					},
				},
			},
			cisAlertFilters: {
				Type:     schema.TypeString,
				Optional: true,
				StateFunc: func(v interface{}) string {
					json, err := flex.NormalizeJSONString(v)
					if err != nil {
						return fmt.Sprintf("%q", err.Error())
					}
					return json
				},
				Description: "Filters based on filter type",
			},
			cisAlertConditions: {
				Type:     schema.TypeString,
				Optional: true,
				StateFunc: func(v interface{}) string {
					json, err := flex.NormalizeJSONString(v)
					if err != nil {
						return fmt.Sprintf("%q", err.Error())
					}
					return json
				},
				Description: "Conditions based on filter type",
			},
		},
	}
}

func ResourceIBMCISAlertPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(conns.ClientSession).CisAlertsSession()
	if err != nil {
		return fmt.Errorf("[ERROR] Error while getting the CisAlertsSession %s", err)
	}
	crn := d.Get(cisID).(string)
	sess.Crn = core.StringPtr(crn)

	opt := sess.NewCreateAlertPolicyOptions()

	if name, ok := d.GetOk(cisAlertName); ok {
		opt.SetName(name.(string))
	}
	if description, ok := d.GetOk(cisAlertDescription); ok {
		opt.SetDescription(description.(string))
	}
	if enabled, ok := d.GetOk(cisAlertEnabled); ok {
		opt.SetEnabled(enabled.(bool))
	}
	if alertType, ok := d.GetOk(cisAlertType); ok {
		opt.SetAlertType(alertType.(string))
	}
	if retFilter, ok := d.GetOk(cisAlertFilters); ok {
		var filter interface{}
		json.Unmarshal([]byte(retFilter.(string)), &filter)
		opt.Filters = filter
	}
	mechanismsOpt := &alertsv1.CreateAlertPolicyInputMechanisms{}
	if mechanisms, ok := d.GetOk(cisAlertMechanisms); ok {
		mechanism := mechanisms.([]interface{})[0].(map[string]interface{})
		webhook, ok := mechanism[cisAlertWebhook]
		if ok {
			webhookString := webhook.(*schema.Set)
			if webhookString.Len() != 0 {
				var webhookarray = make([]alertsv1.CreateAlertPolicyInputMechanismsWebhooksItem, webhookString.Len())
				for k, w := range webhookString.List() {
					wString := w.(string)
					webhookarray[k] = alertsv1.CreateAlertPolicyInputMechanismsWebhooksItem{
						ID: &wString,
					}
				}
				mechanismsOpt.Webhooks = webhookarray
			}
		}
		email, ok := mechanism[cisAlertEmail]
		if ok {
			emailString := email.(*schema.Set)
			if emailString.Len() != 0 {
				var emailarray = make([]alertsv1.CreateAlertPolicyInputMechanismsEmailItem, emailString.Len())
				for k, w := range emailString.List() {
					wString := w.(string)
					emailarray[k] = alertsv1.CreateAlertPolicyInputMechanismsEmailItem{
						ID: &wString,
					}
				}
				mechanismsOpt.Email = emailarray
			}
		}
	}
	opt.Mechanisms = mechanismsOpt
	result, resp, err := sess.CreateAlertPolicy(opt)
	if err != nil || result == nil {
		return fmt.Errorf("[ERROR] Error creating Alert Policy %s %s", err, resp)
	}
	d.SetId(flex.ConvertCisToTfTwoVar(*result.Result.ID, crn))
	d.Set(cisAlertID, *result.Result.ID)

	return ResourceIBMCISAlertPolicyRead(d, meta)
}

func ResourceIBMCISAlertPolicyRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(conns.ClientSession).CisAlertsSession()
	if err != nil {
		return fmt.Errorf("[ERROR] Error while getting the CisAlertsSession %s", err)
	}

	alertID, crn, err := flex.ConvertTftoCisTwoVar(d.Id())
	if err != nil {
		return fmt.Errorf("[ERROR] Error while ConvertTftoCisTwoVar %s", err)
	}
	sess.Crn = core.StringPtr(crn)
	opt := sess.NewGetAlertPolicyOptions(alertID)
	result, resp, err := sess.GetAlertPolicy(opt)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error getting alert policy detail %s, %s", err, resp)
	}

	d.Set(cisID, crn)
	d.Set(cisAlertID, *result.Result.ID)
	d.Set(cisGLBPoolName, *result.Result.Name)
	d.Set(cisAlertDescription, *result.Result.Description)
	d.Set(cisAlertEnabled, *result.Result.Enabled)
	d.Set(cisAlertType, *result.Result.AlertType)
	d.Set(cisGLBPoolEnabled, *result.Result.Enabled)
	if err := d.Set(cisAlertMechanisms, flattenCISMechanism(*result.Result.Mechanisms)); err != nil {
		log.Printf("[WARN] Error setting mechanism for alert policies %q: %s", d.Id(), err)
	}

	filterOpt, err := json.Marshal(result.Result.Filters)
	if err != nil {
		return fmt.Errorf("[ERROR] Error marshalling the created filters: %s", err)
	}
	if err = d.Set(cisAlertFilters, string(filterOpt)); err != nil {
		return fmt.Errorf("[ERROR] Error setting the filters: %s", err)
	}
	conditionsOpt, err := json.Marshal(result.Result.Conditions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error marshalling the created Conditions: %s", err)
	}
	if err = d.Set(cisAlertConditions, string(conditionsOpt)); err != nil {
		return fmt.Errorf("[ERROR] Error setting the Conditions: %s", err)
	}
	return nil
}

func ResourceIBMCISAlertPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(conns.ClientSession).CisAlertsSession()
	if err != nil {
		return fmt.Errorf("[ERROR] Error while getting the CisAlertsSession %s", err)
	}

	alertID, crn, err := flex.ConvertTftoCisTwoVar(d.Id())

	if err != nil {
		return fmt.Errorf("[ERROR] Error while ConvertTftoCisTwoVar %s", err)
	}
	sess.Crn = core.StringPtr(crn)

	if d.HasChange(cisAlertName) ||
		d.HasChange(cisAlertEnabled) ||
		d.HasChange(cisAlertDescription) ||
		d.HasChange(cisAlertType) ||
		d.HasChange(cisAlertFilters) ||
		d.HasChange(cisAlertConditions) ||
		d.HasChange(cisAlertMechanisms) {

		opt := sess.NewUpdateAlertPolicyOptions(alertID)
		if name, ok := d.GetOk(cisAlertName); ok {
			opt.SetName(name.(string))
		}
		if description, ok := d.GetOk(cisAlertDescription); ok {
			opt.SetDescription(description.(string))
		}
		if enabled, ok := d.GetOk(cisAlertEnabled); ok {
			opt.SetEnabled(enabled.(bool))
		}
		if alertType, ok := d.GetOk(cisAlertType); ok {
			opt.SetAlertType(alertType.(string))

		}
		if retConditions, ok := d.GetOk(cisAlertConditions); ok {
			var condition interface{}
			json.Unmarshal([]byte(retConditions.(string)), &condition)
			opt.Conditions = condition
		}

		if retFilter, ok := d.GetOk(cisAlertFilters); ok {
			var filter interface{}
			json.Unmarshal([]byte(retFilter.(string)), &filter)
			opt.Filters = filter
		}

		mechanismsOpt := &alertsv1.UpdateAlertPolicyInputMechanisms{}
		if mechanisms, ok := d.GetOk(cisAlertMechanisms); ok {
			mechanism := mechanisms.([]interface{})[0].(map[string]interface{})
			webhook, ok := mechanism[cisAlertWebhook]
			if ok {
				webhookString := webhook.(*schema.Set)
				if webhookString.Len() != 0 {
					var webhookarray = make([]alertsv1.UpdateAlertPolicyInputMechanismsWebhooksItem, webhookString.Len())
					for k, w := range webhookString.List() {
						wString := w.(string)
						webhookarray[k] = alertsv1.UpdateAlertPolicyInputMechanismsWebhooksItem{
							ID: &wString,
						}
					}
					mechanismsOpt.Webhooks = webhookarray
				}
			}
			email, ok := mechanism[cisAlertEmail]
			if ok {
				emailString := email.(*schema.Set)
				if emailString.Len() != 0 {
					var emailarray = make([]alertsv1.UpdateAlertPolicyInputMechanismsEmailItem, emailString.Len())
					for k, w := range emailString.List() {
						wString := w.(string)
						emailarray[k] = alertsv1.UpdateAlertPolicyInputMechanismsEmailItem{
							ID: &wString,
						}
					}
					mechanismsOpt.Email = emailarray
				}
			}
		}
		opt.Mechanisms = mechanismsOpt

		result, resp, err := sess.UpdateAlertPolicy(opt)
		if err != nil || result == nil {
			return fmt.Errorf("[ERROR] Error while Update Alert Policy %s %s", err, resp)
		}
	}

	return nil
}
func ResourceIBMCISAlertPolicyDelete(d *schema.ResourceData, meta interface{}) error {

	sess, err := meta.(conns.ClientSession).CisAlertsSession()
	if err != nil {
		return fmt.Errorf("[ERROR] Error while getting the CisAlertsSession %s", err)
	}
	alertID, crn, err := flex.ConvertTftoCisTwoVar(d.Id())
	if err != nil {
		return err
	}
	sess.Crn = core.StringPtr(crn)
	opt := sess.NewDeleteAlertPolicyOptions(alertID)
	_, response, err := sess.DeleteAlertPolicy(opt)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return fmt.Errorf("[ERROR] Error deleting the alert %s:%s", err, response)
	}
	return nil
}
func getAlertMechanisms(s *schema.Set) interface{} {
	var alertMechanisms []interface{}
	for _, m := range s.List() {
		switch m.(type) {
		case alertsv1.CreateAlertPolicyInputMechanismsEmailItem:
			email := m.(alertsv1.CreateAlertPolicyInputMechanismsEmailItem)
			data := alertsv1.CreateAlertPolicyInputMechanismsEmailItem{
				ID: email.ID,
			}
			alertMechanisms = append(alertMechanisms, data)
		case alertsv1.CreateAlertPolicyInputMechanismsWebhooksItem:
			webhook := m.(alertsv1.CreateAlertPolicyInputMechanismsWebhooksItem)
			data := alertsv1.CreateAlertPolicyInputMechanismsWebhooksItem{
				ID: webhook.ID,
			}
			alertMechanisms = append(alertMechanisms, data)
		}
	}
	return alertMechanisms
}

func flattenCISMechanism(Mechanism alertsv1.GetAlertPolicyRespResultMechanisms) interface{} {
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
