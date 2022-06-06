// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package functions

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/apache/openwhisk-client-go/whisk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	funcActionName         = "name"
	funcActionNamespace    = "namespace"
	funcActionUsrDefAnnots = "user_defined_annotations"
	funcActionUsrDefParams = "user_defined_parameters"
)

func ResourceIBMFunctionAction() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMFunctionActionCreate,
		Read:     resourceIBMFunctionActionRead,
		Update:   resourceIBMFunctionActionUpdate,
		Delete:   resourceIBMFunctionActionDelete,
		Exists:   resourceIBMFunctionActionExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			funcActionName: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Name of action.",
				ValidateFunc: validate.InvokeValidator("ibm_function_action", funcActionName),
			},
			funcActionNamespace: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "IBM Cloud function namespace.",
				ValidateFunc: validate.InvokeValidator("ibm_function_action", funcActionNamespace),
			},
			"limits": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     60000,
							Description: "The timeout LIMIT in milliseconds after which the action is terminated.",
						},
						"memory": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     256,
							Description: "The maximum memory LIMIT in MB for the action (default 256.",
						},
						"log_size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     10,
							Description: "The maximum log size LIMIT in MB for the action.",
						},
					},
				},
			},
			"exec": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Execution info",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image": {
							Type:          schema.TypeString,
							Optional:      true,
							Description:   "Container image name when kind is 'blackbox'.",
							ConflictsWith: []string{"exec.0.components"},
						},
						"init": {
							Type:          schema.TypeString,
							Optional:      true,
							Description:   "Optional zipfile reference.",
							ConflictsWith: []string{"exec.0.image", "exec.0.components"},
						},
						"code": {
							Type:          schema.TypeString,
							Computed:      true,
							Optional:      true,
							Description:   "The code to execute.",
							ConflictsWith: []string{"exec.0.components", "exec.0.code_path"},
						},
						"code_path": {
							Type:          schema.TypeString,
							Optional:      true,
							Description:   "The file path of code to execute.",
							ConflictsWith: []string{"exec.0.components", "exec.0.code"},
						},
						"kind": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The type of action. Possible values can be found here (https://cloud.ibm.com/docs/openwhisk?topic=cloud-functions-runtimes)",
						},
						"main": {
							Type:          schema.TypeString,
							Optional:      true,
							Description:   "The name of the action entry point (function or fully-qualified method name when applicable).",
							ConflictsWith: []string{"exec.0.image", "exec.0.components"},
						},
						"components": {
							Type:          schema.TypeList,
							Optional:      true,
							Elem:          &schema.Schema{Type: schema.TypeString},
							Description:   "The List of fully qualified action.",
							ConflictsWith: []string{"exec.0.image", "exec.0.code", "exec.0.code_path"},
						},
					},
				},
			},
			"publish": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Action visibilty.",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Semantic version of the item.",
			},
			funcActionUsrDefAnnots: {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "[]",
				Description:      "Annotation values in KEY VALUE format.",
				ValidateFunc:     validate.InvokeValidator("ibm_function_action", funcActionUsrDefAnnots),
				DiffSuppressFunc: flex.SuppressEquivalentJSON,
				StateFunc: func(v interface{}) string {
					json, _ := flex.NormalizeJSONString(v)
					return json
				},
			},
			funcActionUsrDefParams: {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "[]",
				Description:      "Parameters values in KEY VALUE format. Parameter bindings included in the context passed to the action.",
				ValidateFunc:     validate.InvokeValidator("ibm_function_action", funcActionUsrDefParams),
				DiffSuppressFunc: flex.SuppressEquivalentJSON,
				StateFunc: func(v interface{}) string {
					json, _ := flex.NormalizeJSONString(v)
					return json
				},
			},
			"annotations": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "All annotations set on action by user and those set by the IBM Cloud Function backend/API.",
			},
			"parameters": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "All paramters set on action by user and those set by the IBM Cloud Function backend/API.",
			},
			"action_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_endpoint_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Action target endpoint URL.",
			},
		},
	}
}

func ResourceIBMFuncActionValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 funcActionName,
			ValidateFunctionIdentifier: validate.ValidateRegexp,
			Type:                       validate.TypeString,
			Regexp:                     `^[^/*][a-zA-Z0-9/_@.-]`,
			Required:                   true})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 funcActionNamespace,
			ValidateFunctionIdentifier: validate.ValidateNoZeroValues,
			Type:                       validate.TypeString,
			Required:                   true})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 funcActionUsrDefAnnots,
			ValidateFunctionIdentifier: validate.ValidateJSONString,
			Type:                       validate.TypeString,
			Default:                    "[]",
			Optional:                   true})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 funcActionUsrDefParams,
			ValidateFunctionIdentifier: validate.ValidateJSONString,
			Type:                       validate.TypeString,
			Optional:                   true})

	ibmFuncActionResourceValidator := validate.ResourceValidator{ResourceName: "ibm_function_action", Schema: validateSchema}
	return &ibmFuncActionResourceValidator
}

func resourceIBMFunctionActionCreate(d *schema.ResourceData, meta interface{}) error {
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

	actionService := wskClient.Actions
	name := d.Get("name").(string)

	var qualifiedName = new(QualifiedName)

	if qualifiedName, err = NewQualifiedName(name); err != nil {
		return NewQualifiedNameError(name, err)
	}

	payload := whisk.Action{
		Name:      qualifiedName.GetEntityName(),
		Namespace: namespace,
	}

	exec := d.Get("exec").([]interface{})
	payload.Exec = flex.ExpandExec(exec)

	userDefinedAnnotations := d.Get("user_defined_annotations").(string)
	payload.Annotations, err = flex.ExpandAnnotations(userDefinedAnnotations)
	if err != nil {
		return err
	}

	userDefinedParameters := d.Get("user_defined_parameters").(string)
	payload.Parameters, err = flex.ExpandParameters(userDefinedParameters)
	if err != nil {
		return err
	}

	if v, ok := d.GetOk("limits"); ok {
		payload.Limits = flex.ExpandLimits(v.([]interface{}))
	}

	if publish, ok := d.GetOk("publish"); ok {
		p := publish.(bool)
		payload.Publish = &p
	}

	log.Println("[INFO] Creating IBM Cloud Function Action")
	_, _, err = actionService.Insert(&payload, true)

	if err != nil {
		return fmt.Errorf("[ERROR] Error creating IBM Cloud Function Action: %s", err)
	}

	d.SetId(fmt.Sprintf("%s:%s", namespace, qualifiedName.GetEntityName()))

	return resourceIBMFunctionActionRead(d, meta)
}

func resourceIBMFunctionActionRead(d *schema.ResourceData, meta interface{}) error {
	parts, err := flex.CfIdParts(d.Id())
	if err != nil {
		return err
	}

	namespace := ""
	actionID := ""
	if len(parts) == 2 {
		namespace = parts[0]
		actionID = parts[1]
	} else {
		namespace = os.Getenv("FUNCTION_NAMESPACE")
		actionID = parts[0]
		d.SetId(fmt.Sprintf("%s:%s", namespace, actionID))
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

	actionService := wskClient.Actions
	action, _, err := actionService.Get(actionID, true)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving IBM Cloud Function Action %s : %s", actionID, err)
	}
	d.Set("namespace", namespace)
	d.Set("limits", flex.FlattenLimits(action.Limits))
	d.Set("exec", flex.FlattenExec(action.Exec, d))
	d.Set("publish", action.Publish)
	d.Set("version", action.Version)
	d.Set("action_id", action.Name)
	annotations, err := flex.FlattenAnnotations(action.Annotations)
	if err != nil {
		return err
	}

	d.Set("annotations", annotations)
	parameters, err := flex.FlattenParameters(action.Parameters)
	if err != nil {
		return err
	}
	d.Set("parameters", parameters)

	temp := strings.Split(action.Namespace, "/")
	pkgName := ""
	if len(temp) == 2 {
		pkgName = temp[1]
		d.Set("name", fmt.Sprintf("%s/%s", pkgName, action.Name))
		c, err := whisk.NewClient(http.DefaultClient, &whisk.Config{
			Namespace:         wskClient.Namespace,
			AuthToken:         wskClient.AuthToken,
			Host:              wskClient.Host,
			AdditionalHeaders: wskClient.AdditionalHeaders,
		})
		if err != nil {
			return err
		}
		pkg, _, err := c.Packages.Get(pkgName)
		if err != nil {
			return fmt.Errorf("[ERROR] Error retrieving package IBM Cloud Function package %s : %s", pkgName, err)
		}

		userAnnotations, err := flex.FlattenAnnotations(flex.FilterInheritedAnnotations(pkg.Annotations, action.Annotations))
		if err != nil {
			return err
		}

		d.Set("user_defined_annotations", userAnnotations)
		userParameters, err := flex.FlattenParameters(flex.FilterInheritedParameters(pkg.Parameters, action.Parameters))
		if err != nil {
			return err
		}
		d.Set("user_defined_parameters", userParameters)
	} else {
		d.Set("name", action.Name)
		userDefinedAnnotations, err := flex.FilterActionAnnotations(action.Annotations)
		if err != nil {
			return err
		}
		d.Set("user_defined_annotations", userDefinedAnnotations)

		userDefinedParameters, err := flex.FilterActionParameters(action.Parameters)
		if err != nil {
			return err
		}
		d.Set("user_defined_parameters", userDefinedParameters)
	}

	targetUrl, err := action.ActionURL(wskClient.Config.Host, "/api", wskClient.Config.Version, pkgName)
	if err != nil {
		log.Printf(
			"Error creating target endpoint URL for action (%s) targetURL : %s", d.Id(), err)

	}
	d.Set("target_endpoint_url", targetUrl)

	return nil
}

func resourceIBMFunctionActionUpdate(d *schema.ResourceData, meta interface{}) error {
	parts, err := flex.CfIdParts(d.Id())
	if err != nil {
		return err
	}

	namespace := parts[0]
	actionID := parts[1]

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

	actionService := wskClient.Actions

	var qualifiedName = new(QualifiedName)

	if qualifiedName, err = NewQualifiedName(actionID); err != nil {
		return NewQualifiedNameError(actionID, err)
	}

	payload := whisk.Action{
		Name:      qualifiedName.GetEntityName(),
		Namespace: namespace,
	}

	ischanged := false

	if d.HasChange("publish") {
		p := d.Get("publish").(bool)
		payload.Publish = &p
	}

	if d.HasChange("user_defined_parameters") {
		var err error
		payload.Parameters, err = flex.ExpandParameters(d.Get("user_defined_parameters").(string))
		if err != nil {
			return err
		}
		ischanged = true
	}

	if d.HasChange("user_defined_annotations") {
		var err error
		payload.Annotations, err = flex.ExpandAnnotations(d.Get("user_defined_annotations").(string))
		if err != nil {
			return err
		}
		ischanged = true
	}

	if d.HasChange("exec") {
		exec := d.Get("exec").([]interface{})
		payload.Exec = flex.ExpandExec(exec)
		ischanged = true
	}

	if d.HasChange("limits") {
		limits := d.Get("limits").([]interface{})
		payload.Limits = flex.ExpandLimits(limits)
		ischanged = true
	}

	if ischanged {
		log.Println("[INFO] Update IBM Cloud Function Action")
		_, _, err = actionService.Insert(&payload, true)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating IBM Cloud Function Action: %s", err)
		}
	}

	return resourceIBMFunctionActionRead(d, meta)
}

func resourceIBMFunctionActionDelete(d *schema.ResourceData, meta interface{}) error {
	parts, err := flex.CfIdParts(d.Id())
	if err != nil {
		return err
	}

	namespace := parts[0]
	actionID := parts[1]

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

	actionService := wskClient.Actions

	_, err = actionService.Delete(actionID)
	if err != nil {
		return fmt.Errorf("[ERROR] Error deleting IBM Cloud Function Action: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceIBMFunctionActionExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	parts, err := flex.CfIdParts(d.Id())
	if err != nil {
		return false, err
	}

	namespace := ""
	actionID := ""
	if len(parts) >= 2 {
		namespace = parts[0]
		actionID = parts[1]
	} else {
		namespace = os.Getenv("FUNCTION_NAMESPACE")
		actionID = parts[0]
		d.SetId(fmt.Sprintf("%s:%s", namespace, actionID))
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

	actionService := wskClient.Actions

	action, resp, err := actionService.Get(actionID, true)
	if err != nil {
		if resp.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error communicating with IBM Cloud Function Client : %s", err)
	}

	temp := strings.Split(action.Namespace, "/")
	var name string

	if len(temp) == 2 {
		name = fmt.Sprintf("%s/%s", temp[1], action.Name)
	} else {
		name = action.Name
	}

	return name == actionID, nil
}
