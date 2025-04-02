// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.94.1-71478489-20240820-161623
 */

package codeengine

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/code-engine-go-sdk/codeenginev2"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIbmCodeEngineFunction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmCodeEngineFunctionCreate,
		ReadContext:   resourceIbmCodeEngineFunctionRead,
		UpdateContext: resourceIbmCodeEngineFunctionUpdate,
		DeleteContext: resourceIbmCodeEngineFunctionDelete,
		Importer:      &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_function", "project_id"),
				Description:  "The ID of the project.",
			},
			"code_binary": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies whether the code is binary or not. Defaults to false when `code_reference` is set to a data URL. When `code_reference` is set to a code bundle URL, this field is always true.",
			},
			"code_main": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "main",
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_function", "code_main"),
				Description:  "Specifies the name of the function that should be invoked.",
			},
			"code_reference": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_function", "code_reference"),
				Description:  "Specifies either a reference to a code bundle or the source code itself. To specify the source code, use the data URL scheme and include the source code as base64 encoded. The data URL scheme is defined in [RFC 2397](https://tools.ietf.org/html/rfc2397).",
			},
			"code_secret": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_function", "code_secret"),
				Description:  "The name of the secret that is used to access the specified `code_reference`. The secret is used to authenticate with a non-public endpoint that is specified as`code_reference`.",
			},
			"managed_domain_mappings": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "local_public",
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_function", "managed_domain_mappings"),
				Description:  "Optional value controlling which of the system managed domain mappings will be setup for the function. Valid values are 'local_public', 'local_private' and 'local'. Visibility can only be 'local_private' if the project supports function private visibility.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_function", "name"),
				Description:  "The name of the function.",
			},
			"run_env_variables": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "References to config maps, secrets or literal values, which are defined by the function owner and are exposed as environment variables in the function.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The key to reference as environment variable.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the environment variable.",
						},
						"prefix": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A prefix that can be added to all keys of a full secret or config map reference.",
						},
						"reference": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the secret or config map.",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "literal",
							Description: "Specify the type of the environment variable.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The literal value of the environment variable.",
						},
					},
				},
			},
			"runtime": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_function", "runtime"),
				Description:  "The managed runtime used to execute the injected code.",
			},
			"scale_concurrency": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_function", "scale_concurrency"),
				Description:  "Number of parallel requests handled by a single instance, supported only by Node.js, default is `1`.",
			},
			"scale_cpu_limit": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "1",
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_function", "scale_cpu_limit"),
				Description:  "Optional amount of CPU set for the instance of the function. For valid values see [Supported memory and CPU combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo).",
			},
			"scale_down_delay": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_function", "scale_down_delay"),
				Description:  "Optional amount of time in seconds that delays the scale down behavior for a function.",
			},
			"scale_max_execution_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      60,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_function", "scale_max_execution_time"),
				Description:  "Timeout in secs after which the function is terminated.",
			},
			"scale_memory_limit": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "4G",
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_function", "scale_memory_limit"),
				Description:  "Optional amount of memory set for the instance of the function. For valid values see [Supported memory and CPU combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo). The units for specifying memory are Megabyte (M) or Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information see [Units of measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).",
			},
			"computed_env_variables": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "References to config maps, secrets or literal values, which are defined and set by Code Engine and are exposed as environment variables in the function.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The key to reference as environment variable.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The name of the environment variable.",
						},
						"prefix": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "A prefix that can be added to all keys of a full secret or config map reference.",
						},
						"reference": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The name of the secret or config map.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specify the type of the environment variable.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The literal value of the environment variable.",
						},
					},
				},
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the resource was created.",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL to invoke the function.",
			},
			"endpoint_internal": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL to function that is only visible within the project.",
			},
			"entity_tag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the function instance, which is used to achieve optimistic locking.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When you provision a new function, a relative URL path is created identifying the location of the instance.",
			},
			"function_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The identifier of the resource.",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The region of the project the resource is located in. Possible values: 'au-syd', 'br-sao', 'ca-tor', 'eu-de', 'eu-gb', 'jp-osa', 'jp-tok', 'us-east', 'us-south'.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the function.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the function.",
			},
			"status_details": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The detailed status of the function.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Provides additional information about the status of the function.",
						},
					},
				},
			},
			"etag": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceIbmCodeEngineFunctionValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "project_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[0-9a-z]{8}-[0-9a-z]{4}-[0-9a-z]{4}-[0-9a-z]{4}-[0-9a-z]{12}$`,
			MinValueLength:             36,
			MaxValueLength:             36,
		},
		validate.ValidateSchema{
			Identifier:                 "code_main",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[a-zA-Z_][a-zA-Z0-9_]*$`,
			MinValueLength:             1,
			MaxValueLength:             63,
		},
		validate.ValidateSchema{
			Identifier:                 "code_reference",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^([a-z0-9][a-z0-9\-_.]+[a-z0-9][\/])?([a-z0-9][a-z0-9\-_]+[a-z0-9][\/])?[a-z0-9][a-z0-9\-_.\/]+[a-z0-9](:[\w][\w.\-]{0,127})?(@sha256:[a-fA-F0-9]{64})?$|data:([-\w]+\/[-+\w.]+)?(;?\w+=[-\w]+)*;base64,.*`,
			MinValueLength:             1,
			MaxValueLength:             1048576,
		},
		validate.ValidateSchema{
			Identifier:                 "code_secret",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[a-z0-9]([\-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([\-a-z0-9]*[a-z0-9])?)*$`,
			MinValueLength:             1,
			MaxValueLength:             253,
		},
		validate.ValidateSchema{
			Identifier:                 "managed_domain_mappings",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "local, local_private, local_public",
		},
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-z]([-a-z0-9]*[a-z0-9])?$`,
			MinValueLength:             1,
			MaxValueLength:             63,
		},
		validate.ValidateSchema{
			Identifier:                 "runtime",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-z]*\-[0-9]*(\.[0-9]*)?$`,
			MinValueLength:             1,
			MaxValueLength:             63,
		},
		validate.ValidateSchema{
			Identifier:                 "scale_concurrency",
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			Optional:                   true,
			MinValue:                   "1",
			MaxValue:                   "100",
		},
		validate.ValidateSchema{
			Identifier:                 "scale_cpu_limit",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([0-9.]+)([eEinumkKMGTPB]*)$`,
			MinValueLength:             0,
			MaxValueLength:             10,
		},
		validate.ValidateSchema{
			Identifier:                 "scale_down_delay",
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			Optional:                   true,
			MinValue:                   "0",
			MaxValue:                   "600",
		},
		validate.ValidateSchema{
			Identifier:                 "scale_max_execution_time",
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			Optional:                   true,
			MinValue:                   "1",
			MaxValue:                   "120",
		},
		validate.ValidateSchema{
			Identifier:                 "scale_memory_limit",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([0-9.]+)([eEinumkKMGTPB]*)$`,
			MinValueLength:             0,
			MaxValueLength:             10,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_code_engine_function", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmCodeEngineFunctionCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createFunctionOptions := &codeenginev2.CreateFunctionOptions{}

	createFunctionOptions.SetProjectID(d.Get("project_id").(string))
	createFunctionOptions.SetCodeReference(d.Get("code_reference").(string))
	createFunctionOptions.SetName(d.Get("name").(string))
	createFunctionOptions.SetRuntime(d.Get("runtime").(string))
	if _, ok := d.GetOk("code_binary"); ok {
		createFunctionOptions.SetCodeBinary(d.Get("code_binary").(bool))
	}
	if _, ok := d.GetOk("code_main"); ok {
		createFunctionOptions.SetCodeMain(d.Get("code_main").(string))
	}
	if _, ok := d.GetOk("code_secret"); ok {
		createFunctionOptions.SetCodeSecret(d.Get("code_secret").(string))
	}
	if _, ok := d.GetOk("managed_domain_mappings"); ok {
		createFunctionOptions.SetManagedDomainMappings(d.Get("managed_domain_mappings").(string))
	}
	if _, ok := d.GetOk("run_env_variables"); ok {
		var runEnvVariables []codeenginev2.EnvVarPrototype
		for _, v := range d.Get("run_env_variables").([]interface{}) {
			value := v.(map[string]interface{})
			runEnvVariablesItem, err := ResourceIbmCodeEngineFunctionMapToEnvVarPrototype(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "create", "parse-run_env_variables").GetDiag()
			}
			runEnvVariables = append(runEnvVariables, *runEnvVariablesItem)
		}
		createFunctionOptions.SetRunEnvVariables(runEnvVariables)
	}
	if _, ok := d.GetOk("scale_concurrency"); ok {
		createFunctionOptions.SetScaleConcurrency(int64(d.Get("scale_concurrency").(int)))
	}
	if _, ok := d.GetOk("scale_cpu_limit"); ok {
		createFunctionOptions.SetScaleCpuLimit(d.Get("scale_cpu_limit").(string))
	}
	if _, ok := d.GetOk("scale_down_delay"); ok {
		createFunctionOptions.SetScaleDownDelay(int64(d.Get("scale_down_delay").(int)))
	}
	if _, ok := d.GetOk("scale_max_execution_time"); ok {
		createFunctionOptions.SetScaleMaxExecutionTime(int64(d.Get("scale_max_execution_time").(int)))
	}
	if _, ok := d.GetOk("scale_memory_limit"); ok {
		createFunctionOptions.SetScaleMemoryLimit(d.Get("scale_memory_limit").(string))
	}

	function, _, err := codeEngineClient.CreateFunctionWithContext(context, createFunctionOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateFunctionWithContext failed: %s", err.Error()), "ibm_code_engine_function", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *createFunctionOptions.ProjectID, *function.Name))

	_, err = waitForIbmCodeEngineFunctionCreate(d, meta)
	if err != nil {
		errMsg := fmt.Sprintf("Error waiting for resource IbmCodeEngineFunction (%s) to be created: %s", d.Id(), err)
		return flex.DiscriminatedTerraformErrorf(err, errMsg, "ibm_code_engine_function", "create", "wait-for-state").GetDiag()
	}

	return resourceIbmCodeEngineFunctionRead(context, d, meta)
}

func waitForIbmCodeEngineFunctionCreate(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		return false, err
	}
	getFunctionOptions := &codeenginev2.GetFunctionOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return false, err
	}

	getFunctionOptions.SetProjectID(parts[0])
	getFunctionOptions.SetName(parts[1])

	stateConf := &resource.StateChangeConf{
		Pending: []string{"deploying"},
		Target:  []string{"ready", "failed", "offline"},
		Refresh: func() (interface{}, string, error) {
			stateObj, response, err := codeEngineClient.GetFunction(getFunctionOptions)
			if err != nil {
				if sdkErr, ok := err.(*core.SDKProblem); ok && response.GetStatusCode() == 404 {
					sdkErr.Summary = fmt.Sprintf("The instance %s does not exist anymore: %s", "getFunctionOptions", err)
					return nil, "", sdkErr
				}
				return nil, "", err
			}
			failStates := map[string]bool{"failure": true, "failed": true}
			if failStates[*stateObj.Status] {
				return stateObj, *stateObj.Status, fmt.Errorf("The instance %s failed: %s", "getFunctionOptions", err)
			}
			return stateObj, *stateObj.Status, nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      60 * time.Second,
		MinTimeout: 60 * time.Second,
	}

	return stateConf.WaitForStateContext(context.Background())
}

func resourceIbmCodeEngineFunctionRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getFunctionOptions := &codeenginev2.GetFunctionOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "read", "sep-id-parts").GetDiag()
	}

	getFunctionOptions.SetProjectID(parts[0])
	getFunctionOptions.SetName(parts[1])

	function, response, err := codeEngineClient.GetFunctionWithContext(context, getFunctionOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetFunctionWithContext failed: %s", err.Error()), "ibm_code_engine_function", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("project_id", function.ProjectID); err != nil {
		err = fmt.Errorf("Error setting project_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "read", "set-project_id").GetDiag()
	}
	if !core.IsNil(function.CodeBinary) {
		if err = d.Set("code_binary", function.CodeBinary); err != nil {
			err = fmt.Errorf("Error setting code_binary: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "read", "set-code_binary").GetDiag()
		}
	}
	if !core.IsNil(function.CodeMain) {
		if err = d.Set("code_main", function.CodeMain); err != nil {
			err = fmt.Errorf("Error setting code_main: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "read", "set-code_main").GetDiag()
		}
	}
	if err = d.Set("code_reference", function.CodeReference); err != nil {
		err = fmt.Errorf("Error setting code_reference: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "read", "set-code_reference").GetDiag()
	}
	if !core.IsNil(function.CodeSecret) {
		if err = d.Set("code_secret", function.CodeSecret); err != nil {
			err = fmt.Errorf("Error setting code_secret: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "read", "set-code_secret").GetDiag()
		}
	}
	if !core.IsNil(function.ManagedDomainMappings) {
		if err = d.Set("managed_domain_mappings", function.ManagedDomainMappings); err != nil {
			err = fmt.Errorf("Error setting managed_domain_mappings: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "read", "set-managed_domain_mappings").GetDiag()
		}
	}
	if err = d.Set("name", function.Name); err != nil {
		err = fmt.Errorf("Error setting name: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "read", "set-name").GetDiag()
	}
	if !core.IsNil(function.RunEnvVariables) {
		runEnvVariables := []map[string]interface{}{}
		for _, runEnvVariablesItem := range function.RunEnvVariables {
			runEnvVariablesItemMap, err := ResourceIbmCodeEngineFunctionEnvVarToMap(&runEnvVariablesItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "read", "run_env_variables-to-map").GetDiag()
			}
			runEnvVariables = append(runEnvVariables, runEnvVariablesItemMap)
		}
		if err = d.Set("run_env_variables", runEnvVariables); err != nil {
			err = fmt.Errorf("Error setting run_env_variables: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "read", "set-run_env_variables").GetDiag()
		}
	}
	if err = d.Set("runtime", function.Runtime); err != nil {
		err = fmt.Errorf("Error setting runtime: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "read", "set-runtime").GetDiag()
	}
	if !core.IsNil(function.ScaleConcurrency) {
		if err = d.Set("scale_concurrency", flex.IntValue(function.ScaleConcurrency)); err != nil {
			err = fmt.Errorf("Error setting scale_concurrency: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "read", "set-scale_concurrency").GetDiag()
		}
	}
	if !core.IsNil(function.ScaleCpuLimit) {
		if err = d.Set("scale_cpu_limit", function.ScaleCpuLimit); err != nil {
			err = fmt.Errorf("Error setting scale_cpu_limit: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "read", "set-scale_cpu_limit").GetDiag()
		}
	}
	if !core.IsNil(function.ScaleDownDelay) {
		if err = d.Set("scale_down_delay", flex.IntValue(function.ScaleDownDelay)); err != nil {
			err = fmt.Errorf("Error setting scale_down_delay: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "read", "set-scale_down_delay").GetDiag()
		}
	}
	if !core.IsNil(function.ScaleMaxExecutionTime) {
		if err = d.Set("scale_max_execution_time", flex.IntValue(function.ScaleMaxExecutionTime)); err != nil {
			err = fmt.Errorf("Error setting scale_max_execution_time: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "read", "set-scale_max_execution_time").GetDiag()
		}
	}
	if !core.IsNil(function.ScaleMemoryLimit) {
		if err = d.Set("scale_memory_limit", function.ScaleMemoryLimit); err != nil {
			err = fmt.Errorf("Error setting scale_memory_limit: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "read", "set-scale_memory_limit").GetDiag()
		}
	}
	if !core.IsNil(function.ComputedEnvVariables) {
		computedEnvVariables := []map[string]interface{}{}
		for _, computedEnvVariablesItem := range function.ComputedEnvVariables {
			computedEnvVariablesItemMap, err := ResourceIbmCodeEngineFunctionEnvVarToMap(&computedEnvVariablesItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "read", "computed_env_variables-to-map").GetDiag()
			}
			computedEnvVariables = append(computedEnvVariables, computedEnvVariablesItemMap)
		}
		if err = d.Set("computed_env_variables", computedEnvVariables); err != nil {
			err = fmt.Errorf("Error setting computed_env_variables: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "read", "set-computed_env_variables").GetDiag()
		}
	}
	if !core.IsNil(function.CreatedAt) {
		if err = d.Set("created_at", function.CreatedAt); err != nil {
			err = fmt.Errorf("Error setting created_at: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "read", "set-created_at").GetDiag()
		}
	}
	if !core.IsNil(function.Endpoint) {
		if err = d.Set("endpoint", function.Endpoint); err != nil {
			err = fmt.Errorf("Error setting endpoint: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "read", "set-endpoint").GetDiag()
		}
	}
	if !core.IsNil(function.EndpointInternal) {
		if err = d.Set("endpoint_internal", function.EndpointInternal); err != nil {
			err = fmt.Errorf("Error setting endpoint_internal: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "read", "set-endpoint_internal").GetDiag()
		}
	}
	if err = d.Set("entity_tag", function.EntityTag); err != nil {
		err = fmt.Errorf("Error setting entity_tag: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "read", "set-entity_tag").GetDiag()
	}
	if !core.IsNil(function.Href) {
		if err = d.Set("href", function.Href); err != nil {
			err = fmt.Errorf("Error setting href: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "read", "set-href").GetDiag()
		}
	}
	if !core.IsNil(function.ID) {
		if err = d.Set("function_id", function.ID); err != nil {
			err = fmt.Errorf("Error setting function_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "read", "set-function_id").GetDiag()
		}
	}
	if !core.IsNil(function.Region) {
		if err = d.Set("region", function.Region); err != nil {
			err = fmt.Errorf("Error setting region: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "read", "set-region").GetDiag()
		}
	}
	if !core.IsNil(function.ResourceType) {
		if err = d.Set("resource_type", function.ResourceType); err != nil {
			err = fmt.Errorf("Error setting resource_type: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "read", "set-resource_type").GetDiag()
		}
	}
	if !core.IsNil(function.Status) {
		if err = d.Set("status", function.Status); err != nil {
			err = fmt.Errorf("Error setting status: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "read", "set-status").GetDiag()
		}
	}
	statusDetailsMap, err := ResourceIbmCodeEngineFunctionFunctionStatusToMap(function.StatusDetails)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "read", "status_details-to-map").GetDiag()
	}
	if err = d.Set("status_details", []map[string]interface{}{statusDetailsMap}); err != nil {
		err = fmt.Errorf("Error setting status_details: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "read", "set-status_details").GetDiag()
	}
	if err = d.Set("etag", response.Headers.Get("Etag")); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting etag: %s", err), "ibm_code_engine_function", "read", "set-etag").GetDiag()
	}

	return nil
}

func resourceIbmCodeEngineFunctionUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateFunctionOptions := &codeenginev2.UpdateFunctionOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "update", "sep-id-parts").GetDiag()
	}

	updateFunctionOptions.SetProjectID(parts[0])
	updateFunctionOptions.SetName(parts[1])

	hasChange := false

	patchVals := &codeenginev2.FunctionPatch{}
	if d.HasChange("project_id") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "project_id")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_code_engine_function", "update", "project_id-forces-new").GetDiag()
	}
	if d.HasChange("code_binary") {
		newCodeBinary := d.Get("code_binary").(bool)
		patchVals.CodeBinary = &newCodeBinary
		hasChange = true
	}
	if d.HasChange("code_main") {
		newCodeMain := d.Get("code_main").(string)
		patchVals.CodeMain = &newCodeMain
		hasChange = true
	}
	if d.HasChange("code_reference") {
		newCodeReference := d.Get("code_reference").(string)
		patchVals.CodeReference = &newCodeReference
		hasChange = true
	}
	if d.HasChange("code_secret") {
		newCodeSecret := d.Get("code_secret").(string)
		patchVals.CodeSecret = &newCodeSecret
		hasChange = true
	}
	if d.HasChange("managed_domain_mappings") {
		newManagedDomainMappings := d.Get("managed_domain_mappings").(string)
		patchVals.ManagedDomainMappings = &newManagedDomainMappings
		hasChange = true
	}
	if d.HasChange("run_env_variables") {
		var runEnvVariables []codeenginev2.EnvVarPrototype
		for _, v := range d.Get("run_env_variables").([]interface{}) {
			value := v.(map[string]interface{})
			runEnvVariablesItem, err := ResourceIbmCodeEngineFunctionMapToEnvVarPrototype(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "update", "parse-run_env_variables").GetDiag()
			}
			runEnvVariables = append(runEnvVariables, *runEnvVariablesItem)
		}
		patchVals.RunEnvVariables = runEnvVariables
		hasChange = true
	}
	if d.HasChange("runtime") {
		newRuntime := d.Get("runtime").(string)
		patchVals.Runtime = &newRuntime
		hasChange = true
	}
	if d.HasChange("scale_concurrency") {
		newScaleConcurrency := int64(d.Get("scale_concurrency").(int))
		patchVals.ScaleConcurrency = &newScaleConcurrency
		hasChange = true
	}
	if d.HasChange("scale_cpu_limit") {
		newScaleCpuLimit := d.Get("scale_cpu_limit").(string)
		patchVals.ScaleCpuLimit = &newScaleCpuLimit
		hasChange = true
	}
	if d.HasChange("scale_down_delay") {
		newScaleDownDelay := int64(d.Get("scale_down_delay").(int))
		patchVals.ScaleDownDelay = &newScaleDownDelay
		hasChange = true
	}
	if d.HasChange("scale_max_execution_time") {
		newScaleMaxExecutionTime := int64(d.Get("scale_max_execution_time").(int))
		patchVals.ScaleMaxExecutionTime = &newScaleMaxExecutionTime
		hasChange = true
	}
	if d.HasChange("scale_memory_limit") {
		newScaleMemoryLimit := d.Get("scale_memory_limit").(string)
		patchVals.ScaleMemoryLimit = &newScaleMemoryLimit
		hasChange = true
	}
	updateFunctionOptions.SetIfMatch(d.Get("etag").(string))

	if hasChange {
		// Fields with `nil` values are omitted from the generic map,
		// so we need to re-add them to support removing arguments
		// in merge-patch operations sent to the service.
		updateFunctionOptions.Function = ResourceIbmCodeEngineFunctionFunctionPatchAsPatch(patchVals, d)

		_, _, err = codeEngineClient.UpdateFunctionWithContext(context, updateFunctionOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateFunctionWithContext failed: %s", err.Error()), "ibm_code_engine_function", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIbmCodeEngineFunctionRead(context, d, meta)
}

func resourceIbmCodeEngineFunctionDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteFunctionOptions := &codeenginev2.DeleteFunctionOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_function", "delete", "sep-id-parts").GetDiag()
	}

	deleteFunctionOptions.SetProjectID(parts[0])
	deleteFunctionOptions.SetName(parts[1])

	_, err = codeEngineClient.DeleteFunctionWithContext(context, deleteFunctionOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteFunctionWithContext failed: %s", err.Error()), "ibm_code_engine_function", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIbmCodeEngineFunctionMapToEnvVarPrototype(modelMap map[string]interface{}) (*codeenginev2.EnvVarPrototype, error) {
	model := &codeenginev2.EnvVarPrototype{}
	if modelMap["key"] != nil && modelMap["key"].(string) != "" {
		model.Key = core.StringPtr(modelMap["key"].(string))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["prefix"] != nil && modelMap["prefix"].(string) != "" {
		model.Prefix = core.StringPtr(modelMap["prefix"].(string))
	}
	if modelMap["reference"] != nil && modelMap["reference"].(string) != "" {
		model.Reference = core.StringPtr(modelMap["reference"].(string))
	}
	if modelMap["type"] != nil && modelMap["type"].(string) != "" {
		model.Type = core.StringPtr(modelMap["type"].(string))
	}
	if modelMap["value"] != nil && modelMap["value"].(string) != "" {
		model.Value = core.StringPtr(modelMap["value"].(string))
	}
	return model, nil
}

func ResourceIbmCodeEngineFunctionEnvVarToMap(model *codeenginev2.EnvVar) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Key != nil {
		modelMap["key"] = *model.Key
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Prefix != nil {
		modelMap["prefix"] = *model.Prefix
	}
	if model.Reference != nil {
		modelMap["reference"] = *model.Reference
	}
	modelMap["type"] = *model.Type
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	return modelMap, nil
}

func ResourceIbmCodeEngineFunctionFunctionStatusToMap(model *codeenginev2.FunctionStatus) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Reason != nil {
		modelMap["reason"] = *model.Reason
	}
	return modelMap, nil
}

func ResourceIbmCodeEngineFunctionFunctionPatchAsPatch(patchVals *codeenginev2.FunctionPatch, d *schema.ResourceData) map[string]interface{} {
	patch, _ := patchVals.AsPatch()
	var path string

	path = "code_binary"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["code_binary"] = nil
	}
	path = "code_main"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["code_main"] = nil
	}
	path = "code_reference"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["code_reference"] = nil
	}
	path = "code_secret"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["code_secret"] = nil
	}
	path = "managed_domain_mappings"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["managed_domain_mappings"] = nil
	}
	path = "run_env_variables"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		runEnvVariables := []map[string]interface{}{}
		patch["run_env_variables"] = runEnvVariables
	}
	path = "runtime"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["runtime"] = nil
	}
	path = "scale_concurrency"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["scale_concurrency"] = nil
	}
	path = "scale_cpu_limit"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["scale_cpu_limit"] = nil
	}
	path = "scale_down_delay"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["scale_down_delay"] = nil
	}
	path = "scale_max_execution_time"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["scale_max_execution_time"] = nil
	}
	path = "scale_memory_limit"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["scale_memory_limit"] = nil
	}

	return patch
}

func ResourceIbmCodeEngineFunctionEnvVarPrototypeAsPatch(patch map[string]interface{}, d *schema.ResourceData) {
	var path string

	path = "run_env_variables.0.key"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["key"] = nil
	}
	path = "run_env_variables.0.name"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["name"] = nil
	}
	path = "run_env_variables.0.prefix"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["prefix"] = nil
	}
	path = "run_env_variables.0.reference"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["reference"] = nil
	}
	path = "run_env_variables.0.type"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["type"] = nil
	}
	path = "run_env_variables.0.value"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["value"] = nil
	}
}
