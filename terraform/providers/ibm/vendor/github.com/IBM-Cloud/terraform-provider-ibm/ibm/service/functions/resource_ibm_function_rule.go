// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package functions

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/apache/openwhisk-client-go/whisk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	funcRuleNamespace = "namespace"
	funcRuleName      = "name"
)

func ResourceIBMFunctionRule() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMFunctionRuleCreate,
		Read:     resourceIBMFunctionRuleRead,
		Update:   resourceIBMFunctionRuleUpdate,
		Delete:   resourceIBMFunctionRuleDelete,
		Exists:   resourceIBMFunctionRuleExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			funcRuleNamespace: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "IBM Cloud function namespace.",
				ValidateFunc: validate.InvokeValidator("ibm_function_rule", funcRuleNamespace),
			},
			funcRuleName: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Name of rule.",
				ValidateFunc: validate.InvokeValidator("ibm_function_rule", funcRuleName),
			},
			"trigger_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of trigger.",
			},
			"action_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of action.",
				DiffSuppressFunc: func(k, o, n string, d *schema.ResourceData) bool {
					new := strings.Split(n, "/")
					old := strings.Split(o, "/")
					action_name_new := new[len(new)-1]
					action_name_old := old[len(old)-1]

					if o == "" {
						return false
					}
					if strings.HasPrefix(n, "/_") {
						temp := strings.Replace(n, "/_", "/"+d.Get("namespace").(string), 1)
						if strings.Compare(temp, o) == 0 {
							return true
						}
						if strings.Compare(action_name_old, action_name_new) == 0 {
							return true
						}

					}
					if !strings.HasPrefix(n, "/") {
						if strings.HasPrefix(o, "/"+d.Get("namespace").(string)) {
							return true
						}
						if strings.Compare(action_name_old, action_name_new) == 0 {
							return true
						}
					}
					return false
				},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the rule.",
			},
			"publish": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Rule visbility.",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Semantic version of the item.",
			},
			"rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceIBMFuncRuleValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 funcRuleName,
			ValidateFunctionIdentifier: validate.ValidateRegexp,
			Type:                       validate.TypeString,
			Regexp:                     `\A([\w]|[\w][\w@ .-]*[\w@.-]+)\z`,
			Required:                   true})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 funcRuleNamespace,
			ValidateFunctionIdentifier: validate.ValidateNoZeroValues,
			Type:                       validate.TypeString,
			Required:                   true})

	ibmFuncRuleResourceValidator := validate.ResourceValidator{ResourceName: "ibm_function_rule", Schema: validateSchema}
	return &ibmFuncRuleResourceValidator
}

func resourceIBMFunctionRuleCreate(d *schema.ResourceData, meta interface{}) error {
	functionNamespaceAPI, err := meta.(conns.ClientSession).FunctionIAMNamespaceAPI()
	if err != nil {
		return err
	}

	bxSession, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return err
	}
	namespace := d.Get("namespace").(string)
	wskClient, err := conns.SetupOpenWhiskClientConfig(namespace, bxSession, functionNamespaceAPI)
	if err != nil {
		return err

	}

	ruleService := wskClient.Rules

	name := d.Get("name").(string)

	var qualifiedName = new(QualifiedName)
	if qualifiedName, err = NewQualifiedName(name); err != nil {
		return NewQualifiedNameError(name, err)
	}
	trigger := d.Get("trigger_name").(string)
	action := d.Get("action_name").(string)

	triggerName := getQualifiedName(trigger, wskClient.Config.Namespace)
	actionName := getQualifiedName(action, wskClient.Config.Namespace)
	payload := whisk.Rule{
		Name:      qualifiedName.GetEntityName(),
		Namespace: qualifiedName.GetNamespace(),
		Trigger:   triggerName,
		Action:    actionName,
	}
	log.Println("[INFO] Creating IBM Cloud Function rule")
	result, _, err := ruleService.Insert(&payload, false)
	if err != nil {
		return fmt.Errorf("[ERROR] Error creating IBM Cloud Function rule: %s", err)
	}

	d.SetId(fmt.Sprintf("%s:%s", namespace, result.Name))

	return resourceIBMFunctionRuleRead(d, meta)
}

func resourceIBMFunctionRuleRead(d *schema.ResourceData, meta interface{}) error {
	parts, err := flex.CfIdParts(d.Id())
	if err != nil {
		return err
	}

	namespace := ""
	ruleID := ""
	if len(parts) == 2 {
		namespace = parts[0]
		ruleID = parts[1]
	} else {
		namespace = os.Getenv("FUNCTION_NAMESPACE")
		ruleID = parts[0]
		d.SetId(fmt.Sprintf("%s:%s", namespace, ruleID))
	}

	functionNamespaceAPI, err := meta.(conns.ClientSession).FunctionIAMNamespaceAPI()
	if err != nil {
		return err
	}

	bxSession, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return err
	}

	wskClient, err := conns.SetupOpenWhiskClientConfig(namespace, bxSession, functionNamespaceAPI)
	if err != nil {
		return err

	}

	ruleService := wskClient.Rules
	rule, _, err := ruleService.Get(ruleID)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving IBM Cloud Function rule %s : %s", ruleID, err)
	}

	d.Set("rule_id", rule.Name)
	d.Set("name", rule.Name)
	d.Set("publish", rule.Publish)
	d.Set("namespace", namespace)
	d.Set("version", rule.Version)
	d.Set("status", rule.Status)

	path := rule.Action.(map[string]interface{})["path"]
	d.Set("trigger_name", rule.Trigger.(map[string]interface{})["name"])
	actionName := rule.Action.(map[string]interface{})["name"]
	d.Set("action_name", fmt.Sprintf("/%s/%s", path, actionName))
	d.SetId(fmt.Sprintf("%s:%s", namespace, rule.Name))
	return nil
}

func resourceIBMFunctionRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	functionNamespaceAPI, err := meta.(conns.ClientSession).FunctionIAMNamespaceAPI()
	if err != nil {
		return err
	}

	bxSession, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return err
	}

	parts, err := flex.CfIdParts(d.Id())
	if err != nil {
		return err
	}

	namespace := parts[0]
	wskClient, err := conns.SetupOpenWhiskClientConfig(namespace, bxSession, functionNamespaceAPI)
	if err != nil {
		return err

	}

	ruleService := wskClient.Rules

	var qualifiedName = new(QualifiedName)

	if qualifiedName, err = NewQualifiedName(d.Get("name").(string)); err != nil {
		return NewQualifiedNameError(d.Get("name").(string), err)
	}

	payload := whisk.Rule{
		Name:      qualifiedName.GetEntityName(),
		Namespace: qualifiedName.GetNamespace(),
	}
	ischanged := false

	if d.HasChange("trigger_name") {
		trigger := d.Get("trigger_name").(string)
		payload.Trigger = getQualifiedName(trigger, wskClient.Config.Namespace)
		ischanged = true
	}

	if d.HasChange("action_name") {
		action := d.Get("action_name").(string)
		payload.Action = getQualifiedName(action, wskClient.Config.Namespace)
		ischanged = true
	}

	if ischanged {
		log.Println("[INFO] Update IBM Cloud Function Rule")
		result, _, err := ruleService.Insert(&payload, true)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating IBM Cloud Function Rule: %s", err)
		}
		_, _, err = ruleService.SetState(result.Name, "active")
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating IBM Cloud Function Rule: %s", err)
		}
	}

	return resourceIBMFunctionRuleRead(d, meta)
}

func resourceIBMFunctionRuleDelete(d *schema.ResourceData, meta interface{}) error {
	functionNamespaceAPI, err := meta.(conns.ClientSession).FunctionIAMNamespaceAPI()
	if err != nil {
		return err
	}

	bxSession, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return err
	}

	parts, err := flex.CfIdParts(d.Id())
	if err != nil {
		return err
	}

	namespace := parts[0]
	ruleID := parts[1]
	wskClient, err := conns.SetupOpenWhiskClientConfig(namespace, bxSession, functionNamespaceAPI)
	if err != nil {
		return err

	}

	ruleService := wskClient.Rules

	_, err = ruleService.Delete(ruleID)
	if err != nil {
		return fmt.Errorf("[ERROR] Error deleting IBM Cloud Function Rule: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceIBMFunctionRuleExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	parts, err := flex.CfIdParts(d.Id())
	if err != nil {
		return false, err
	}

	namespace := ""
	ruleID := ""
	if len(parts) == 2 {
		namespace = parts[0]
		ruleID = parts[1]
	} else {
		namespace = os.Getenv("FUNCTION_NAMESPACE")
		ruleID = parts[0]
		d.SetId(fmt.Sprintf("%s:%s", namespace, ruleID))
	}

	functionNamespaceAPI, err := meta.(conns.ClientSession).FunctionIAMNamespaceAPI()
	if err != nil {
		return false, err
	}

	bxSession, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return false, err
	}

	wskClient, err := conns.SetupOpenWhiskClientConfig(namespace, bxSession, functionNamespaceAPI)
	if err != nil {
		return false, err

	}

	ruleService := wskClient.Rules

	rule, resp, err := ruleService.Get(ruleID)
	if err != nil {
		if resp.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error communicating with IBM Cloud Function Client : %s", err)
	}
	return rule.Name == ruleID, nil
}
