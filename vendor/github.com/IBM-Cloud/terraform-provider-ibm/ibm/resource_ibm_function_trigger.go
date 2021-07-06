// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/apache/openwhisk-client-go/whisk"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	funcTriggerNamespace    = "namespace"
	funcTriggerName         = "name"
	funcTriggerParams       = "parameters"
	funcTriggerUsrDefAnnots = "user_defined_annotations"
	funcTriggerUsrDefParams = "user_defined_parameters"

	feedLifeCycleEvent = "lifecycleEvent"
	feedTriggerName    = "triggerName"
	feedAuthKey        = "authKey"
	feedCreate         = "CREATE"
	feedDelete         = "DELETE"
)

func resourceIBMFunctionTrigger() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMFunctionTriggerCreate,
		Read:     resourceIBMFunctionTriggerRead,
		Update:   resourceIBMFunctionTriggerUpdate,
		Delete:   resourceIBMFunctionTriggerDelete,
		Exists:   resourceIBMFunctionTriggerExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			funcTriggerNamespace: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "IBM Cloud function namespace.",
				ValidateFunc: InvokeValidator("ibm_function_trigger", funcTriggerNamespace),
			},
			funcTriggerName: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Name of Trigger.",
				ValidateFunc: InvokeValidator("ibm_function_trigger", funcTriggerName),
			},
			"feed": {
				Type:        schema.TypeList,
				ForceNew:    true,
				Optional:    true,
				MaxItems:    1,
				Description: "Trigger feed",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Trigger feed ACTION_NAME.",
						},
						funcTriggerParams: {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "[]",
							Description:  "Parameters values in KEY VALUE format. Parameter bindings included in the context.TODO() passed to the action invoke.",
							ValidateFunc: InvokeValidator("ibm_function_trigger", funcTriggerParams),
							DiffSuppressFunc: func(k, o, n string, d *schema.ResourceData) bool {
								if o == "" && n == "" {
									return false
								}
								if o == "[]" {
									return true
								}
								return false
							},
							StateFunc: func(v interface{}) string {
								json, _ := normalizeJSONString(v)
								return json
							},
						},
					},
				},
			},
			"publish": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Trigger visbility.",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Semantic version of the item.",
			},
			funcTriggerUsrDefAnnots: {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Annotation values in KEY VALUE format.",
				Default:          "[]",
				ValidateFunc:     InvokeValidator("ibm_function_trigger", funcTriggerUsrDefAnnots),
				DiffSuppressFunc: suppressEquivalentJSON,
				StateFunc: func(v interface{}) string {
					json, _ := normalizeJSONString(v)
					return json
				},
			},
			funcTriggerUsrDefParams: {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "[]",
				Description:      "Parameters values in KEY VALUE format. Parameter bindings included in the context.TODO() passed to the trigger.",
				ValidateFunc:     InvokeValidator("ibm_function_trigger", funcTriggerUsrDefParams),
				DiffSuppressFunc: suppressEquivalentJSON,
				StateFunc: func(v interface{}) string {
					json, _ := normalizeJSONString(v)
					return json
				},
			},
			"annotations": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "All annotations set on trigger by user and those set by the IBM Cloud Function backend/API.",
			},
			"parameters": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "All parameters set on trigger by user and those set by the IBM Cloud Function backend/API.",
			},
			"trigger_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIBMFuncTriggerValidator() *ResourceValidator {
	validateSchema := make([]ValidateSchema, 1)

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 funcTriggerName,
			ValidateFunctionIdentifier: ValidateRegexp,
			Type:                       TypeString,
			Regexp:                     `\A([\w]|[\w][\w@ .-]*[\w@.-]+)\z`,
			Required:                   true})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 funcTriggerNamespace,
			ValidateFunctionIdentifier: ValidateNoZeroValues,
			Type:                       TypeString,
			Required:                   true})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 funcTriggerUsrDefAnnots,
			ValidateFunctionIdentifier: ValidateJSONString,
			Type:                       TypeString,
			Default:                    "[]",
			Optional:                   true})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 funcTriggerUsrDefParams,
			ValidateFunctionIdentifier: ValidateJSONString,
			Type:                       TypeString,
			Optional:                   true})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 funcTriggerParams,
			ValidateFunctionIdentifier: ValidateJSONString,
			Type:                       TypeString,
			Default:                    "[]",
			Optional:                   true})

	ibmFuncTriggerResourceValidator := ResourceValidator{ResourceName: "ibm_function_trigger", Schema: validateSchema}
	return &ibmFuncTriggerResourceValidator
}

func resourceIBMFunctionTriggerCreate(d *schema.ResourceData, meta interface{}) error {
	functionNamespaceAPI, err := meta.(ClientSession).FunctionIAMNamespaceAPI()
	if err != nil {
		return err
	}

	bxSession, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return err
	}
	namespace := d.Get("namespace").(string)
	wskClient, err := setupOpenWhiskClientConfig(namespace, bxSession, functionNamespaceAPI)
	if err != nil {
		return err

	}

	triggerService := wskClient.Triggers
	feed := false
	feedPayload := map[string]interface{}{}
	name := d.Get("name").(string)

	var qualifiedName = new(QualifiedName)

	if qualifiedName, err = NewQualifiedName(name); err != nil {
		return NewQualifiedNameError(name, err)
	}

	payload := whisk.Trigger{
		Name:      qualifiedName.GetEntityName(),
		Namespace: qualifiedName.GetNamespace(),
	}

	userDefinedAnnotations := d.Get("user_defined_annotations").(string)
	payload.Annotations, err = expandAnnotations(userDefinedAnnotations)
	if err != nil {
		return err
	}

	userDefinedParameters := d.Get("user_defined_parameters").(string)
	payload.Parameters, err = expandParameters(userDefinedParameters)
	if err != nil {
		return err
	}

	if v, ok := d.GetOk("feed"); ok {
		feed = true
		value := v.([]interface{})[0].(map[string]interface{})
		feedPaylod := whisk.KeyValue{
			Key:   "feed",
			Value: value["name"],
		}
		feedArray := make([]whisk.KeyValue, 0, 1)
		feedArray = append(feedArray, feedPaylod)
		payload.Annotations = payload.Annotations.AppendKeyValueArr(feedArray)
	}

	log.Println("[INFO] Creating IBM Cloud Function trigger")
	result, _, err := triggerService.Insert(&payload, false)
	if err != nil {
		return fmt.Errorf("Error creating IBM Cloud Function trigger: %s", err)
	}

	d.SetId(fmt.Sprintf("%s:%s", namespace, result.Name))

	if feed {
		feed := d.Get("feed").([]interface{})[0].(map[string]interface{})
		actionName := feed["name"].(string)
		parameters := feed["parameters"].(string)
		var err error
		feedParameters, err := expandParameters(parameters)
		if err != nil {
			return err
		}
		for _, value := range feedParameters {
			feedPayload[value.Key] = value.Value
		}
		var feedQualifiedName = new(QualifiedName)

		if feedQualifiedName, err = NewQualifiedName(actionName); err != nil {
			_, _, delerr := triggerService.Delete(name)
			if delerr != nil {
				return fmt.Errorf("Error creating IBM Cloud Function trigger with feed: %s", err)
			}
			return NewQualifiedNameError(actionName, err)
		}

		feedPayload[feedLifeCycleEvent] = feedCreate
		feedPayload[feedAuthKey] = wskClient.Config.AuthToken
		feedPayload[feedTriggerName] = fmt.Sprintf("/%s/%s", qualifiedName.GetNamespace(), name)

		c, err := whisk.NewClient(http.DefaultClient, &whisk.Config{
			AuthToken:         wskClient.AuthToken,
			Host:              wskClient.Host,
			AdditionalHeaders: wskClient.AdditionalHeaders,
		})

		if feedQualifiedName.GetNamespace() != namespace {
			c.Config.Namespace = feedQualifiedName.GetNamespace()
		}
		actionService := c.Actions
		_, _, err = actionService.Invoke(feedQualifiedName.GetEntityName(), feedPayload, true, true)
		if err != nil {
			_, _, delerr := triggerService.Delete(name)
			if delerr != nil {
				return fmt.Errorf("Error creating IBM Cloud Function trigger with feed: %s", err)
			}
			d.SetId("")
			return fmt.Errorf("Error creating IBM Cloud Function trigger with feed: %s", err)
		}
	}

	d.SetId(fmt.Sprintf("%s:%s", namespace, result.Name))

	return resourceIBMFunctionTriggerRead(d, meta)
}

func resourceIBMFunctionTriggerRead(d *schema.ResourceData, meta interface{}) error {
	parts, err := cfIdParts(d.Id())
	if err != nil {
		return err
	}

	namespace := ""
	triggerID := ""
	if len(parts) == 2 {
		namespace = parts[0]
		triggerID = parts[1]
	} else {
		namespace = os.Getenv("FUNCTION_NAMESPACE")
		triggerID = parts[0]
		d.SetId(fmt.Sprintf("%s:%s", namespace, triggerID))
	}

	functionNamespaceAPI, err := meta.(ClientSession).FunctionIAMNamespaceAPI()
	if err != nil {
		return err
	}

	bxSession, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return err
	}

	wskClient, err := setupOpenWhiskClientConfig(namespace, bxSession, functionNamespaceAPI)
	if err != nil {
		return err

	}

	triggerService := wskClient.Triggers

	trigger, _, err := triggerService.Get(triggerID)
	if err != nil {
		return fmt.Errorf("Error retrieving IBM Cloud Function Trigger %s : %s", triggerID, err)
	}
	d.Set("trigger_id", trigger.Name)
	d.Set("namespace", namespace)
	d.Set("name", trigger.Name)
	d.Set("publish", trigger.Publish)
	d.Set("version", trigger.Version)
	annotations, err := flattenAnnotations(trigger.Annotations)
	if err != nil {
		return err
	}
	d.Set("annotations", annotations)
	parameters, err := flattenParameters(trigger.Parameters)
	if err != nil {
		return err
	}
	d.Set("parameters", parameters)
	d.Set("user_defined_parameters", parameters)

	userDefinedAnnotations, err := filterTriggerAnnotations(trigger.Annotations)
	if err != nil {
		return err
	}
	d.Set("user_defined_annotations", userDefinedAnnotations)

	found := trigger.Annotations.FindKeyValue("feed")

	if found >= 0 {
		d.Set("feed", flattenFeed(trigger.Annotations.GetValue("feed").(string)))
	}

	return nil
}

func resourceIBMFunctionTriggerUpdate(d *schema.ResourceData, meta interface{}) error {
	functionNamespaceAPI, err := meta.(ClientSession).FunctionIAMNamespaceAPI()
	if err != nil {
		return err
	}

	bxSession, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return err
	}
	parts, err := cfIdParts(d.Id())
	if err != nil {
		return err
	}

	namespace := parts[0]
	wskClient, err := setupOpenWhiskClientConfig(namespace, bxSession, functionNamespaceAPI)
	if err != nil {
		return err

	}

	triggerService := wskClient.Triggers

	name := d.Get("name").(string)

	var qualifiedName = new(QualifiedName)

	if qualifiedName, err = NewQualifiedName(name); err != nil {
		return NewQualifiedNameError(name, err)
	}

	payload := whisk.Trigger{
		Name:      qualifiedName.GetEntityName(),
		Namespace: qualifiedName.GetNamespace(),
	}
	ischanged := false

	if d.HasChange("user_defined_parameters") {
		var err error
		payload.Parameters, err = expandParameters(d.Get("user_defined_parameters").(string))
		if err != nil {
			return err
		}
		ischanged = true
	}

	if d.HasChange("user_defined_annotations") {
		var err error
		payload.Annotations, err = expandAnnotations(d.Get("user_defined_annotations").(string))
		if err != nil {
			return err
		}
		ischanged = true
	}

	if ischanged {
		log.Println("[INFO] Update IBM Cloud Function Trigger")

		_, _, err = triggerService.Insert(&payload, true)
		if err != nil {
			return fmt.Errorf("Error updating IBM Cloud Function Trigger: %s", err)
		}
	}

	return resourceIBMFunctionTriggerRead(d, meta)
}

func resourceIBMFunctionTriggerDelete(d *schema.ResourceData, meta interface{}) error {
	functionNamespaceAPI, err := meta.(ClientSession).FunctionIAMNamespaceAPI()
	if err != nil {
		return err
	}

	bxSession, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return err
	}

	parts, err := cfIdParts(d.Id())
	if err != nil {
		return err
	}
	namespace := parts[0]
	triggerID := parts[1]

	wskClient, err := setupOpenWhiskClientConfig(namespace, bxSession, functionNamespaceAPI)
	if err != nil {
		return err

	}

	triggerService := wskClient.Triggers
	var qualifiedName = new(QualifiedName)
	fmt.Println(qualifiedName)
	if qualifiedName, err = NewQualifiedName(triggerID); err != nil {
		return NewQualifiedNameError(triggerID, err)
	}
	trigger, _, err := triggerService.Get(triggerID)
	if err != nil {
		return fmt.Errorf("Error retrieving IBM Cloud Function Trigger %s : %s", triggerID, err)
	}
	found := trigger.Annotations.FindKeyValue("feed")
	if found >= 0 {
		actionName := trigger.Annotations.GetValue("feed").(string)
		var feedQualifiedName = new(QualifiedName)

		if feedQualifiedName, err = NewQualifiedName(actionName); err != nil {
			return NewQualifiedNameError(actionName, err)
		}

		feedPayload := map[string]interface{}{
			feedLifeCycleEvent: feedDelete,
			feedAuthKey:        wskClient.Config.AuthToken,
			feedTriggerName:    fmt.Sprintf("/%s/%s", qualifiedName.GetNamespace(), triggerID),
		}

		c, err := whisk.NewClient(http.DefaultClient, &whisk.Config{
			AuthToken:         wskClient.AuthToken,
			Host:              wskClient.Host,
			AdditionalHeaders: wskClient.AdditionalHeaders,
		})
		if feedQualifiedName.GetNamespace() != namespace {
			c.Config.Namespace = feedQualifiedName.GetNamespace()
		}

		actionService := c.Actions
		_, _, err = actionService.Invoke(feedQualifiedName.GetEntityName(), feedPayload, true, true)
		if err != nil {
			return fmt.Errorf("Error deleting IBM Cloud Function trigger with feed: %s", err)

		}
	}

	_, _, err = triggerService.Delete(triggerID)
	if err != nil {
		return fmt.Errorf("Error deleting IBM Cloud Function Trigger: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceIBMFunctionTriggerExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	parts, err := cfIdParts(d.Id())
	if err != nil {
		return false, err
	}

	namespace := ""
	triggerID := ""
	if len(parts) == 2 {
		namespace = parts[0]
		triggerID = parts[1]
	} else {
		namespace = os.Getenv("FUNCTION_NAMESPACE")
		triggerID = parts[0]
		d.SetId(fmt.Sprintf("%s:%s", namespace, triggerID))
	}

	functionNamespaceAPI, err := meta.(ClientSession).FunctionIAMNamespaceAPI()
	if err != nil {
		return false, err
	}

	bxSession, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return false, err
	}

	wskClient, err := setupOpenWhiskClientConfig(namespace, bxSession, functionNamespaceAPI)
	if err != nil {
		return false, err

	}

	triggerService := wskClient.Triggers
	trigger, resp, err := triggerService.Get(triggerID)
	if err != nil {
		if resp.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error communicating with IBM Cloud Function Client : %s", err)
	}
	return trigger.Name == triggerID, nil
}
