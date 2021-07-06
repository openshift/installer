// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/sl"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/internal/hashcode"
)

const (
	IBMComputeTimeFormat                                 = string("2006-01-02T15:04:05-07:00")
	IBMCOMPUTE_SCALE_POLICY_TRIGGER_TYPE_ID_RESOURCE_USE = 1
	IBMCOMPUTE_SCALE_POLICY_TRIGGER_TYPE_ID_REPEATING    = 2
	IBMCOMPUTE_SCALE_POLICY_TRIGGER_TYPE_ID_ONE_TIME     = 3
)

var IBMComputeAutoScalePolicyObjectMask = []string{
	"cooldown",
	"id",
	"name",
	"scaleActions",
	"scaleGroupId",
	"oneTimeTriggers",
	"repeatingTriggers",
	"resourceUseTriggers.watches",
	"triggers",
}

func resourceIBMComputeAutoScalePolicy() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMComputeAutoScalePolicyCreate,
		Read:     resourceIBMComputeAutoScalePolicyRead,
		Update:   resourceIBMComputeAutoScalePolicyUpdate,
		Delete:   resourceIBMComputeAutoScalePolicyDelete,
		Exists:   resourceIBMComputeAutoScalePolicyExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name",
			},
			"scale_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "scale type",
			},
			"scale_amount": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Scale amount",
			},
			"cooldown": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "cooldown value",
			},
			"scale_group_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "scale group ID",
			},
			"triggers": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},

						// Conditionally-required fields, based on value of "type"
						"watches": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"metric": {
										Type:     schema.TypeString,
										Required: true,
									},
									"operator": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
									"period": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
							Set: resourceIBMComputeAutoScalePolicyHandlerHash,
						},

						"date": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"schedule": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Set: resourceIBMComputeAutoScalePolicyTriggerHash,
			},

			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "List of tags",
			},
		},
	}
}

func resourceIBMComputeAutoScalePolicyCreate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetScalePolicyService(sess.SetRetries(0))

	var err error

	// Build up creation options
	opts := datatypes.Scale_Policy{
		Name:         sl.String(d.Get("name").(string)),
		ScaleGroupId: sl.Int(d.Get("scale_group_id").(int)),
		Cooldown:     sl.Int(d.Get("cooldown").(int)),
	}

	if *opts.Cooldown < 0 || *opts.Cooldown > 864000 {
		return fmt.Errorf("Error retrieving scalePolicy: %s", "cooldown must be between 0 seconds and 10 days.")
	}

	opts.ScaleActions = []datatypes.Scale_Policy_Action_Scale{{
		Amount:    sl.Int(d.Get("scale_amount").(int)),
		ScaleType: sl.String(d.Get("scale_type").(string)),
	},
	}
	opts.ScaleActions[0].TypeId = sl.Int(1)

	if *opts.ScaleActions[0].Amount <= 0 {
		return fmt.Errorf("Error retrieving scalePolicy: %s", "scale_amount should be greater than 0.")
	}
	if *opts.ScaleActions[0].ScaleType != "ABSOLUTE" && *opts.ScaleActions[0].ScaleType != "RELATIVE" && *opts.ScaleActions[0].ScaleType != "PERCENT" {
		return fmt.Errorf("Error retrieving scalePolicy: %s", "scale_type should be ABSOLUTE, RELATIVE, or PERCENT.")
	}

	if _, ok := d.GetOk("triggers"); ok {
		err = validateTriggerTypes(d)
		if err != nil {
			return fmt.Errorf("Error retrieving scalePolicy: %s", err)
		}

		opts.OneTimeTriggers, err = prepareOneTimeTriggers(d)
		if err != nil {
			return fmt.Errorf("Error retrieving scalePolicy: %s", err)
		}

		opts.RepeatingTriggers, err = prepareRepeatingTriggers(d)
		if err != nil {
			return fmt.Errorf("Error retrieving scalePolicy: %s", err)
		}

		opts.ResourceUseTriggers, err = prepareResourceUseTriggers(d)
		if err != nil {
			return fmt.Errorf("Error retrieving scalePolicy: %s", err)
		}
	}

	res, err := service.CreateObject(&opts)
	if err != nil {
		return fmt.Errorf("Error creating Scale Policy: %s $s", err)
	}

	d.SetId(strconv.Itoa(*res.Id))
	log.Printf("[INFO] Scale Polocy: %d", res.Id)

	return resourceIBMComputeAutoScalePolicyRead(d, meta)
}

func resourceIBMComputeAutoScalePolicyRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetScalePolicyService(sess)

	scalePolicyId, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid scale policy ID, must be an integer: %s", err)
	}

	log.Printf("[INFO] Reading Scale Polocy: %d", scalePolicyId)
	scalePolicy, err := service.Id(scalePolicyId).Mask(strings.Join(IBMComputeAutoScalePolicyObjectMask, ";")).GetObject()
	if err != nil {
		return fmt.Errorf("Error retrieving Scale Policy: %s", err)
	}

	d.Set("name", scalePolicy.Name)
	d.Set("cooldown", scalePolicy.Cooldown)
	d.Set("scale_group_id", scalePolicy.ScaleGroupId)
	d.Set("scale_type", scalePolicy.ScaleActions[0].ScaleType)
	d.Set("scale_amount", scalePolicy.ScaleActions[0].Amount)
	triggers := make([]map[string]interface{}, 0)
	triggers = append(triggers, readOneTimeTriggers(scalePolicy.OneTimeTriggers)...)
	triggers = append(triggers, readRepeatingTriggers(scalePolicy.RepeatingTriggers)...)
	triggers = append(triggers, readResourceUseTriggers(scalePolicy.ResourceUseTriggers)...)

	d.Set("triggers", triggers)

	return nil
}

func resourceIBMComputeAutoScalePolicyUpdate(d *schema.ResourceData, meta interface{}) error {

	sess := meta.(ClientSession).SoftLayerSession()
	scalePolicyService := services.GetScalePolicyService(sess)
	scalePolicyTriggerService := services.GetScalePolicyTriggerService(sess)
	scalePolicyServiceNoRetry := services.GetScalePolicyService(sess.SetRetries(0))

	scalePolicyId, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid scale policy ID, must be an integer: %s", err)
	}

	scalePolicy, err := scalePolicyService.Id(scalePolicyId).Mask(strings.Join(IBMComputeAutoScalePolicyObjectMask, ";")).GetObject()
	if err != nil {
		return fmt.Errorf("Error retrieving scalePolicy: %s", err)
	}

	var template datatypes.Scale_Policy

	template.Id = sl.Int(scalePolicyId)

	if d.HasChange("name") {
		template.Name = sl.String(d.Get("name").(string))
	}

	if d.HasChange("scale_type") || d.HasChange("scale_amount") {
		template.ScaleActions = make([]datatypes.Scale_Policy_Action_Scale, 1)
		template.ScaleActions[0].Id = scalePolicy.ScaleActions[0].Id
		template.ScaleActions[0].TypeId = sl.Int(1)
	}

	if d.HasChange("scale_type") {
		template.ScaleActions[0].ScaleType = sl.String(d.Get("scale_type").(string))
		if *template.ScaleActions[0].ScaleType != "ABSOLUTE" && *template.ScaleActions[0].ScaleType != "RELATIVE" && *template.ScaleActions[0].ScaleType != "PERCENT" {
			return fmt.Errorf("Error retrieving scalePolicy: %s", "scale_type should be ABSOLUTE, RELATIVE, or PERCENT.")
		}
	}

	if d.HasChange("scale_amount") {
		template.ScaleActions[0].Amount = sl.Int(d.Get("scale_amount").(int))
		if *template.ScaleActions[0].Amount <= 0 {
			return fmt.Errorf("Error retrieving scalePolicy: %s", "scale_amount should be greater than 0.")
		}
	}

	if d.HasChange("cooldown") {
		template.Cooldown = sl.Int(d.Get("cooldown").(int))
		if *template.Cooldown <= 0 || *template.Cooldown > 864000 {
			return fmt.Errorf("Error retrieving scalePolicy: %s", "cooldown must be between 0 seconds and 10 days.")
		}
	}

	if _, ok := d.GetOk("triggers"); ok {
		template.OneTimeTriggers, err = prepareOneTimeTriggers(d)
		if err != nil {
			return fmt.Errorf("Error retrieving scalePolicy: %s", err)
		}
		template.RepeatingTriggers, err = prepareRepeatingTriggers(d)
		if err != nil {
			return fmt.Errorf("Error retrieving scalePolicy: %s", err)
		}
		template.ResourceUseTriggers, err = prepareResourceUseTriggers(d)
		if err != nil {
			return fmt.Errorf("Error retrieving scalePolicy: %s", err)
		}
	}

	for _, triggerList := range scalePolicy.Triggers {
		log.Printf("[INFO] DELETE TRIGGERS %d", *triggerList.Id)
		scalePolicyTriggerService.Id(*triggerList.Id).DeleteObject()
	}

	time.Sleep(60)
	log.Printf("[INFO] Updating scale policy: %d", scalePolicyId)
	_, err = scalePolicyServiceNoRetry.Id(scalePolicyId).EditObject(&template)

	if err != nil {
		return fmt.Errorf("Error updating scalie policy: %s", err)
	}

	return nil
}

func resourceIBMComputeAutoScalePolicyDelete(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetScalePolicyService(sess)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting scale policy: %s", err)
	}

	log.Printf("[INFO] Deleting scale policy: %d", id)
	_, err = service.Id(id).DeleteObject()
	if err != nil {
		return fmt.Errorf("Error deleting scale policy: %s", err)
	}

	d.SetId("")

	return nil
}

func resourceIBMComputeAutoScalePolicyExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetScalePolicyService(sess)

	policyId, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}
	result, err := service.Id(policyId).Mask("id").GetObject()
	if err != nil {
		if apiErr, ok := err.(sl.Error); ok {
			if apiErr.StatusCode == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}
	return result.Id != nil && *result.Id == policyId, nil

}

func validateTriggerTypes(d *schema.ResourceData) error {
	triggerLists := d.Get("triggers").(*schema.Set).List()
	for _, triggerList := range triggerLists {
		trigger := triggerList.(map[string]interface{})
		trigger_type := trigger["type"].(string)
		if trigger_type != "ONE_TIME" && trigger_type != "REPEATING" && trigger_type != "RESOURCE_USE" {
			return fmt.Errorf("Invalid trigger type: %s", trigger_type)
		}
	}
	return nil
}

func prepareOneTimeTriggers(d *schema.ResourceData) ([]datatypes.Scale_Policy_Trigger_OneTime, error) {
	triggerLists := d.Get("triggers").(*schema.Set).List()
	triggers := make([]datatypes.Scale_Policy_Trigger_OneTime, 0)

	portalTimeZone := time.FixedZone("PortalTimeZone", -5*60*60)

	for _, triggerList := range triggerLists {
		trigger := triggerList.(map[string]interface{})

		if trigger["type"].(string) == "ONE_TIME" {
			var oneTimeTrigger datatypes.Scale_Policy_Trigger_OneTime
			oneTimeTrigger.TypeId = sl.Int(IBMCOMPUTE_SCALE_POLICY_TRIGGER_TYPE_ID_ONE_TIME)
			timeStampString := trigger["date"].(string)

			// Use UTC time zone for a terraform configuration
			isUTC := strings.HasSuffix(timeStampString, "+00:00")
			if !isUTC {
				return nil, errors.New("The time zone should be an UTC(+00:00).")
			}

			timeStamp, err := time.Parse(IBMComputeTimeFormat, timeStampString)
			if err != nil {
				return nil, err
			}
			oneTimeTrigger.Date = &datatypes.Time{Time: timeStamp.In(portalTimeZone)}
			triggers = append(triggers, oneTimeTrigger)
		}
	}
	return triggers, nil
}

func prepareRepeatingTriggers(d *schema.ResourceData) ([]datatypes.Scale_Policy_Trigger_Repeating, error) {
	triggerLists := d.Get("triggers").(*schema.Set).List()
	triggers := make([]datatypes.Scale_Policy_Trigger_Repeating, 0)
	for _, triggerList := range triggerLists {
		trigger := triggerList.(map[string]interface{})

		if trigger["type"].(string) == "REPEATING" {
			var repeatingTrigger datatypes.Scale_Policy_Trigger_Repeating
			repeatingTrigger.TypeId = sl.Int(IBMCOMPUTE_SCALE_POLICY_TRIGGER_TYPE_ID_REPEATING)
			repeatingTrigger.Schedule = sl.String(trigger["schedule"].(string))
			triggers = append(triggers, repeatingTrigger)
		}
	}
	return triggers, nil
}

func prepareResourceUseTriggers(d *schema.ResourceData) ([]datatypes.Scale_Policy_Trigger_ResourceUse, error) {
	triggerLists := d.Get("triggers").(*schema.Set).List()
	triggers := make([]datatypes.Scale_Policy_Trigger_ResourceUse, 0)
	for _, triggerList := range triggerLists {
		trigger := triggerList.(map[string]interface{})

		if trigger["type"].(string) == "RESOURCE_USE" {
			var resourceUseTrigger datatypes.Scale_Policy_Trigger_ResourceUse
			var err error
			resourceUseTrigger.TypeId = sl.Int(IBMCOMPUTE_SCALE_POLICY_TRIGGER_TYPE_ID_RESOURCE_USE)
			resourceUseTrigger.Watches, err = prepareWatches(trigger["watches"].(*schema.Set))
			if err != nil {
				return nil, err
			}
			triggers = append(triggers, resourceUseTrigger)
		}
	}
	return triggers, nil
}

func prepareWatches(d *schema.Set) ([]datatypes.Scale_Policy_Trigger_ResourceUse_Watch, error) {
	watchLists := d.List()
	watches := make([]datatypes.Scale_Policy_Trigger_ResourceUse_Watch, 0)
	for _, watcheList := range watchLists {
		var watch datatypes.Scale_Policy_Trigger_ResourceUse_Watch
		watchMap := watcheList.(map[string]interface{})

		watch.Metric = sl.String(watchMap["metric"].(string))
		if *watch.Metric != "host.cpu.percent" && *watch.Metric != "host.network.backend.in.rate" && *watch.Metric != "host.network.backend.out.rate" && *watch.Metric != "host.network.frontend.in.rate" && *watch.Metric != "host.network.frontend.out.rate" {
			return nil, fmt.Errorf("Invalid metric : %s", *watch.Metric)
		}

		watch.Operator = sl.String(watchMap["operator"].(string))
		if *watch.Operator != ">" && *watch.Operator != "<" {
			return nil, fmt.Errorf("Invalid operator : %s", *watch.Operator)
		}

		watch.Period = sl.Int(watchMap["period"].(int))
		if *watch.Period <= 0 {
			return nil, errors.New("period shoud be greater than 0.")
		}

		watch.Value = sl.String(watchMap["value"].(string))

		// Autoscale only support EWMA algorithm.
		watch.Algorithm = sl.String("EWMA")

		watches = append(watches, watch)
	}
	return watches, nil
}

func readOneTimeTriggers(list []datatypes.Scale_Policy_Trigger_OneTime) []map[string]interface{} {
	triggers := make([]map[string]interface{}, 0, len(list))
	UTCZone, _ := time.LoadLocation("UTC")

	for _, trigger := range list {
		t := make(map[string]interface{})
		t["id"] = *trigger.Id
		t["type"] = "ONE_TIME"
		t["date"] = trigger.Date.In(UTCZone).Format(IBMComputeTimeFormat)
		triggers = append(triggers, t)
	}
	return triggers
}

func readRepeatingTriggers(list []datatypes.Scale_Policy_Trigger_Repeating) []map[string]interface{} {
	triggers := make([]map[string]interface{}, 0, len(list))
	for _, trigger := range list {
		t := make(map[string]interface{})
		t["id"] = *trigger.Id
		t["type"] = "REPEATING"
		t["schedule"] = *trigger.Schedule
		triggers = append(triggers, t)
	}
	return triggers
}

func readResourceUseTriggers(list []datatypes.Scale_Policy_Trigger_ResourceUse) []map[string]interface{} {
	triggers := make([]map[string]interface{}, 0, len(list))
	for _, trigger := range list {
		t := make(map[string]interface{})
		t["id"] = *trigger.Id
		t["type"] = "RESOURCE_USE"
		t["watches"] = schema.NewSet(resourceIBMComputeAutoScalePolicyHandlerHash,
			readResourceUseWatches(trigger.Watches))
		triggers = append(triggers, t)
	}
	return triggers
}

func readResourceUseWatches(list []datatypes.Scale_Policy_Trigger_ResourceUse_Watch) []interface{} {
	watches := make([]interface{}, 0, len(list))
	for _, watch := range list {
		w := make(map[string]interface{})
		w["id"] = *watch.Id
		w["metric"] = *watch.Metric
		w["operator"] = *watch.Operator
		w["period"] = *watch.Period
		w["value"] = *watch.Value
		watches = append(watches, w)
	}
	return watches
}

func resourceIBMComputeAutoScalePolicyTriggerHash(v interface{}) int {
	var buf bytes.Buffer
	trigger := v.(map[string]interface{})
	if trigger["type"].(string) == "ONE_TIME" {
		buf.WriteString(fmt.Sprintf("%s-", trigger["type"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", trigger["date"].(string)))
	}
	if trigger["type"].(string) == "REPEATING" {
		buf.WriteString(fmt.Sprintf("%s-", trigger["type"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", trigger["schedule"].(string)))
	}
	if trigger["type"].(string) == "RESOURCE_USE" {
		buf.WriteString(fmt.Sprintf("%s-", trigger["type"].(string)))
		for _, watchList := range trigger["watches"].(*schema.Set).List() {
			watch := watchList.(map[string]interface{})
			buf.WriteString(fmt.Sprintf("%s-", watch["metric"].(string)))
			buf.WriteString(fmt.Sprintf("%s-", watch["operator"].(string)))
			buf.WriteString(fmt.Sprintf("%s-", watch["value"].(string)))
			buf.WriteString(fmt.Sprintf("%d-", watch["period"].(int)))
		}
	}
	return hashcode.String(buf.String())
}

func resourceIBMComputeAutoScalePolicyHandlerHash(v interface{}) int {
	var buf bytes.Buffer
	watch := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", watch["metric"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", watch["operator"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", watch["value"].(string)))
	buf.WriteString(fmt.Sprintf("%d-", watch["period"].(int)))
	return hashcode.String(buf.String())
}
