// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package codeengine

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/code-engine-go-sdk/codeenginev2"
	"github.com/IBM/go-sdk-core/v5/core"
)

func ResourceIbmCodeEngineApp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmCodeEngineAppCreate,
		ReadContext:   resourceIbmCodeEngineAppRead,
		UpdateContext: resourceIbmCodeEngineAppUpdate,
		DeleteContext: resourceIbmCodeEngineAppDelete,
		Importer:      &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_app", "project_id"),
				Description:  "The ID of the project.",
			},
			"image_reference": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_app", "image_reference"),
				Description:  "The name of the image that is used for this app. The format is `REGISTRY/NAMESPACE/REPOSITORY:TAG` where `REGISTRY` and `TAG` are optional. If `REGISTRY` is not specified, the default is `docker.io`. If `TAG` is not specified, the default is `latest`. If the image reference points to a registry that requires authentication, make sure to also specify the property `image_secret`.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_app", "name"),
				Description:  "The name of the app. Use a name that is unique within the project.",
			},
			"image_port": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     8080,
				Description: "Optional port the app listens on. While the app will always be exposed via port `443` for end users, this port is used to connect to the port that is exposed by the container image.",
			},
			"image_secret": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_app", "image_secret"),
				Description:  "Optional name of the image registry access secret. The image registry access secret is used to authenticate with a private registry when you download the container image. If the image reference points to a registry that requires authentication, the app will be created but cannot reach the ready status, until this property is provided, too.",
			},
			"managed_domain_mappings": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "local_public",
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_app", "managed_domain_mappings"),
				Description:  "Optional value controlling which of the system managed domain mappings will be setup for the application. Valid values are 'local_public', 'local_private' and 'local'. Visibility can only be 'local_private' if the project supports application private visibility.",
			},
			"run_arguments": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				MinItems:    0,
				Description: "Optional arguments for the app that are passed to start the container. If not specified an empty string array will be applied and the arguments specified by the container image, will be used to start the container.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"run_as_user": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Optional user ID (UID) to run the app (e.g., `1001`).",
			},
			"run_commands": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				MinItems:    0,
				Description: "Optional commands for the app that are passed to start the container. If not specified an empty string array will be applied and the command specified by the container image, will be used to start the container.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"run_env_variables": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional references to config maps, secrets or literal values that are exposed as environment variables within the running application.",
				MinItems:    0,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The key to reference as environment variable.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the environment variable.",
						},
						"prefix": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A prefix that can be added to all keys of a full secret or config map reference.",
						},
						"reference": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the secret or config map.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "literal",
							Description: "Specify the type of the environment variable.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The literal value of the environment variable.",
						},
					},
				},
			},
			"run_service_account": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "default",
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_app", "run_service_account"),
				Description:  "Optional name of the service account. For built-in service accounts, you can use the shortened names `manager` , `none`, `reader`, and `writer`.",
			},
			"run_volume_mounts": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				MinItems:    0,
				Description: "Optional mounts of config maps or a secrets.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mount_path": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The path that should be mounted.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Optional name of the mount. If not set, it will be generated based on the `ref` and a random ID. In case the `ref` is longer than 58 characters, it will be cut off.",
						},
						"reference": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the referenced secret or config map.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specify the type of the volume mount. Allowed types are: 'config_map', 'secret'.",
						},
					},
				},
			},
			"scale_concurrency": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     100,
				Description: "Optional maximum number of requests that can be processed concurrently per instance.",
			},
			"scale_concurrency_target": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Optional threshold of concurrent requests per instance at which one or more additional instances are created. Use this value to scale up instances based on concurrent number of requests. This option defaults to the value of the `scale_concurrency` option, if not specified.",
			},
			"scale_cpu_limit": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "1",
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_app", "scale_cpu_limit"),
				Description:  "Optional number of CPU set for the instance of the app. For valid values see [Supported memory and CPU combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo).",
			},
			"scale_ephemeral_storage_limit": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "400M",
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_app", "scale_ephemeral_storage_limit"),
				Description:  "Optional amount of ephemeral storage to set for the instance of the app. The amount specified as ephemeral storage, must not exceed the amount of `scale_memory_limit`. The units for specifying ephemeral storage are Megabyte (M) or Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information see [Units of measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).",
			},
			"scale_initial_instances": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Optional initial number of instances that are created upon app creation or app update.",
			},
			"scale_max_instances": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "Optional maximum number of instances for this app. If you set this value to `0`, this property does not set a upper scaling limit. However, the app scaling is still limited by the project quota for instances. See [Limits and quotas for Code Engine](https://cloud.ibm.com/docs/codeengine?topic=codeengine-limits).",
			},
			"scale_memory_limit": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "4G",
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_app", "scale_memory_limit"),
				Description:  "Optional amount of memory set for the instance of the app. For valid values see [Supported memory and CPU combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo). The units for specifying memory are Megabyte (M) or Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information see [Units of measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).",
			},
			"scale_min_instances": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Optional minimum number of instances for this app. If you set this value to `0`, the app will scale down to zero, if not hit by any request for some time.",
			},
			"scale_request_timeout": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     300,
				Description: "Optional amount of time in seconds that is allowed for a running app to respond to a request.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the resource was created.",
			},
			"endpoint": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Optional URL to invoke app. Depending on visibility this is accessible publicly or in the private network only. Empty in case 'managed_domain_mappings' is set to 'local'.",
			},
			"endpoint_internal": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL to app that is only visible within the project.",
			},
			"entity_tag": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the app instance, which is used to achieve optimistic locking.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When you provision a new app,  a URL is created identifying the location of the instance.",
			},
			"app_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The identifier of the resource.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the app.",
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the app.",
			},
			"status_details": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The detailed status of the application.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"latest_created_revision": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Latest app revision that has been created.",
						},
						"latest_ready_revision": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Latest app revision that reached a ready state.",
						},
						"reason": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Optional information to provide more context in case of a 'failed' or 'warning' status.",
						},
					},
				},
			},
			"etag": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceIbmCodeEngineAppValidator() *validate.ResourceValidator {
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
			Identifier:                 "image_reference",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^([a-z0-9][a-z0-9\-_.]+[a-z0-9][\/])?([a-z0-9][a-z0-9\-_]+[a-z0-9][\/])?[a-z0-9][a-z0-9\-_.\/]+[a-z0-9](:[\w][\w.\-]{0,127})?(@sha256:[a-fA-F0-9]{64})?$`,
			MinValueLength:             1,
			MaxValueLength:             256,
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
			Identifier:                 "image_secret",
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
			Identifier:                 "run_service_account",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "default, manager, none, reader, writer",
			Regexp:                     `^(manager|reader|writer|none|default)$`,
			MinValueLength:             0,
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
			Identifier:                 "scale_ephemeral_storage_limit",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([0-9.]+)([eEinumkKMGTPB]*)$`,
			MinValueLength:             0,
			MaxValueLength:             10,
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

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_code_engine_app", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmCodeEngineAppCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	createAppOptions := &codeenginev2.CreateAppOptions{}

	createAppOptions.SetProjectID(d.Get("project_id").(string))
	createAppOptions.SetImageReference(d.Get("image_reference").(string))
	createAppOptions.SetName(d.Get("name").(string))
	if _, ok := d.GetOk("image_port"); ok {
		createAppOptions.SetImagePort(int64(d.Get("image_port").(int)))
	}
	if _, ok := d.GetOk("image_secret"); ok {
		createAppOptions.SetImageSecret(d.Get("image_secret").(string))
	}
	if _, ok := d.GetOk("managed_domain_mappings"); ok {
		createAppOptions.SetManagedDomainMappings(d.Get("managed_domain_mappings").(string))
	}
	if _, ok := d.GetOk("run_arguments"); ok {
		var runArguments []string
		for _, v := range d.Get("run_arguments").([]interface{}) {
			runArgumentsItem := v.(string)
			runArguments = append(runArguments, runArgumentsItem)
		}
		createAppOptions.SetRunArguments(runArguments)
	}
	if _, ok := d.GetOk("run_as_user"); ok {
		createAppOptions.SetRunAsUser(int64(d.Get("run_as_user").(int)))
	}
	if _, ok := d.GetOk("run_commands"); ok {
		var runCommands []string
		for _, v := range d.Get("run_commands").([]interface{}) {
			runCommandsItem := v.(string)
			runCommands = append(runCommands, runCommandsItem)
		}
		createAppOptions.SetRunCommands(runCommands)
	}
	if _, ok := d.GetOk("run_env_variables"); ok {
		var runEnvVariables []codeenginev2.EnvVarPrototype
		for _, v := range d.Get("run_env_variables").([]interface{}) {
			value := v.(map[string]interface{})
			runEnvVariablesItem, err := resourceIbmCodeEngineAppMapToEnvVarPrototype(value)
			if err != nil {
				return diag.FromErr(err)
			}
			runEnvVariables = append(runEnvVariables, *runEnvVariablesItem)
		}
		createAppOptions.SetRunEnvVariables(runEnvVariables)
	}
	if _, ok := d.GetOk("run_service_account"); ok {
		createAppOptions.SetRunServiceAccount(d.Get("run_service_account").(string))
	}
	if _, ok := d.GetOk("run_volume_mounts"); ok {
		var runVolumeMounts []codeenginev2.VolumeMountPrototype
		for _, v := range d.Get("run_volume_mounts").([]interface{}) {
			value := v.(map[string]interface{})
			runVolumeMountsItem, err := resourceIbmCodeEngineAppMapToVolumeMountPrototype(value)
			if err != nil {
				return diag.FromErr(err)
			}
			runVolumeMounts = append(runVolumeMounts, *runVolumeMountsItem)
		}
		createAppOptions.SetRunVolumeMounts(runVolumeMounts)
	}
	if _, ok := d.GetOk("scale_concurrency"); ok {
		createAppOptions.SetScaleConcurrency(int64(d.Get("scale_concurrency").(int)))
	}
	if _, ok := d.GetOk("scale_concurrency_target"); ok {
		createAppOptions.SetScaleConcurrencyTarget(int64(d.Get("scale_concurrency_target").(int)))
	}
	if _, ok := d.GetOk("scale_cpu_limit"); ok {
		createAppOptions.SetScaleCpuLimit(d.Get("scale_cpu_limit").(string))
	}
	if _, ok := d.GetOk("scale_ephemeral_storage_limit"); ok {
		createAppOptions.SetScaleEphemeralStorageLimit(d.Get("scale_ephemeral_storage_limit").(string))
	}
	if _, ok := d.GetOk("scale_initial_instances"); ok {
		createAppOptions.SetScaleInitialInstances(int64(d.Get("scale_initial_instances").(int)))
	}
	if _, ok := d.GetOk("scale_max_instances"); ok {
		createAppOptions.SetScaleMaxInstances(int64(d.Get("scale_max_instances").(int)))
	}
	if _, ok := d.GetOk("scale_memory_limit"); ok {
		createAppOptions.SetScaleMemoryLimit(d.Get("scale_memory_limit").(string))
	}
	if _, ok := d.GetOk("scale_min_instances"); ok {
		createAppOptions.SetScaleMinInstances(int64(d.Get("scale_min_instances").(int)))
	}
	if _, ok := d.GetOk("scale_request_timeout"); ok {
		createAppOptions.SetScaleRequestTimeout(int64(d.Get("scale_request_timeout").(int)))
	}

	app, response, err := codeEngineClient.CreateAppWithContext(context, createAppOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateAppWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateAppWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *createAppOptions.ProjectID, *app.Name))

	_, err = waitForIbmCodeEngineAppCreate(d, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf(
			"Error waiting for resource IbmCodeEngineApp (%s) to be created: %s", d.Id(), err))
	}

	return resourceIbmCodeEngineAppRead(context, d, meta)
}

func waitForIbmCodeEngineAppCreate(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		return false, err
	}
	getAppOptions := &codeenginev2.GetAppOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return false, err
	}

	getAppOptions.SetProjectID(parts[0])
	getAppOptions.SetName(parts[1])

	stateConf := &resource.StateChangeConf{
		Pending: []string{"deploying"},
		Target:  []string{"ready", "failed", "warning"},
		Refresh: func() (interface{}, string, error) {
			stateObj, response, err := codeEngineClient.GetApp(getAppOptions)
			if err != nil {
				if apiErr, ok := err.(bmxerror.RequestFailure); ok && apiErr.StatusCode() == 404 {
					return nil, "", fmt.Errorf("The instance %s does not exist anymore: %s\n%s", "getAppOptions", err, response)
				}
				return nil, "", err
			}
			failStates := map[string]bool{"failure": true, "failed": true}
			if failStates[*stateObj.Status] {
				return stateObj, *stateObj.Status, fmt.Errorf("The instance %s failed: %s\n%s", "getAppOptions", err, response)
			}
			return stateObj, *stateObj.Status, nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      60 * time.Second,
		MinTimeout: 60 * time.Second,
	}

	return stateConf.WaitForState()
}

func resourceIbmCodeEngineAppRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	getAppOptions := &codeenginev2.GetAppOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getAppOptions.SetProjectID(parts[0])
	getAppOptions.SetName(parts[1])

	app, response, err := codeEngineClient.GetAppWithContext(context, getAppOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetAppWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetAppWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("project_id", app.ProjectID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting project_id: %s", err))
	}
	if err = d.Set("image_reference", app.ImageReference); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting image_reference: %s", err))
	}
	if err = d.Set("name", app.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if !core.IsNil(app.ImagePort) {
		if err = d.Set("image_port", flex.IntValue(app.ImagePort)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting image_port: %s", err))
		}
	}
	if !core.IsNil(app.ImageSecret) {
		if err = d.Set("image_secret", app.ImageSecret); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting image_secret: %s", err))
		}
	}
	if !core.IsNil(app.ManagedDomainMappings) {
		if err = d.Set("managed_domain_mappings", app.ManagedDomainMappings); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting managed_domain_mappings: %s", err))
		}
	}
	if !core.IsNil(app.RunArguments) {
		if err = d.Set("run_arguments", app.RunArguments); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting run_arguments: %s", err))
		}
	}
	if !core.IsNil(app.RunAsUser) {
		if err = d.Set("run_as_user", flex.IntValue(app.RunAsUser)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting run_as_user: %s", err))
		}
	}
	if !core.IsNil(app.RunCommands) {
		if err = d.Set("run_commands", app.RunCommands); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting run_commands: %s", err))
		}
	}
	if !core.IsNil(app.RunEnvVariables) {
		runEnvVariables := []map[string]interface{}{}
		for _, runEnvVariablesItem := range app.RunEnvVariables {
			runEnvVariablesItemMap, err := resourceIbmCodeEngineAppEnvVarToMap(&runEnvVariablesItem)
			if err != nil {
				return diag.FromErr(err)
			}
			runEnvVariables = append(runEnvVariables, runEnvVariablesItemMap)
		}
		if err = d.Set("run_env_variables", runEnvVariables); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting run_env_variables: %s", err))
		}
	}
	if !core.IsNil(app.RunServiceAccount) {
		if err = d.Set("run_service_account", app.RunServiceAccount); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting run_service_account: %s", err))
		}
	}
	if !core.IsNil(app.RunVolumeMounts) {
		runVolumeMounts := []map[string]interface{}{}
		for _, runVolumeMountsItem := range app.RunVolumeMounts {
			runVolumeMountsItemMap, err := resourceIbmCodeEngineAppVolumeMountToMap(&runVolumeMountsItem)
			if err != nil {
				return diag.FromErr(err)
			}
			runVolumeMounts = append(runVolumeMounts, runVolumeMountsItemMap)
		}
		if err = d.Set("run_volume_mounts", runVolumeMounts); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting run_volume_mounts: %s", err))
		}
	}
	if !core.IsNil(app.ScaleConcurrency) {
		if err = d.Set("scale_concurrency", flex.IntValue(app.ScaleConcurrency)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting scale_concurrency: %s", err))
		}
	}
	if !core.IsNil(app.ScaleConcurrencyTarget) {
		if err = d.Set("scale_concurrency_target", flex.IntValue(app.ScaleConcurrencyTarget)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting scale_concurrency_target: %s", err))
		}
	}
	if !core.IsNil(app.ScaleCpuLimit) {
		if err = d.Set("scale_cpu_limit", app.ScaleCpuLimit); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting scale_cpu_limit: %s", err))
		}
	}
	if !core.IsNil(app.ScaleEphemeralStorageLimit) {
		if err = d.Set("scale_ephemeral_storage_limit", app.ScaleEphemeralStorageLimit); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting scale_ephemeral_storage_limit: %s", err))
		}
	}
	if !core.IsNil(app.ScaleInitialInstances) {
		if err = d.Set("scale_initial_instances", flex.IntValue(app.ScaleInitialInstances)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting scale_initial_instances: %s", err))
		}
	}
	if !core.IsNil(app.ScaleMaxInstances) {
		if err = d.Set("scale_max_instances", flex.IntValue(app.ScaleMaxInstances)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting scale_max_instances: %s", err))
		}
	}
	if !core.IsNil(app.ScaleMemoryLimit) {
		if err = d.Set("scale_memory_limit", app.ScaleMemoryLimit); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting scale_memory_limit: %s", err))
		}
	}
	if !core.IsNil(app.ScaleMinInstances) {
		if err = d.Set("scale_min_instances", flex.IntValue(app.ScaleMinInstances)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting scale_min_instances: %s", err))
		}
	}
	if !core.IsNil(app.ScaleRequestTimeout) {
		if err = d.Set("scale_request_timeout", flex.IntValue(app.ScaleRequestTimeout)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting scale_request_timeout: %s", err))
		}
	}
	if !core.IsNil(app.CreatedAt) {
		if err = d.Set("created_at", app.CreatedAt); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
		}
	}
	if !core.IsNil(app.Endpoint) {
		if err = d.Set("endpoint", app.Endpoint); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting endpoint: %s", err))
		}
	}
	if !core.IsNil(app.EndpointInternal) {
		if err = d.Set("endpoint_internal", app.EndpointInternal); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting endpoint_internal: %s", err))
		}
	}
	if err = d.Set("entity_tag", app.EntityTag); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting entity_tag: %s", err))
	}
	if !core.IsNil(app.Href) {
		if err = d.Set("href", app.Href); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
		}
	}
	if !core.IsNil(app.ID) {
		if err = d.Set("app_id", app.ID); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting app_id: %s", err))
		}
	}
	if !core.IsNil(app.ResourceType) {
		if err = d.Set("resource_type", app.ResourceType); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
		}
	}
	if !core.IsNil(app.Status) {
		if err = d.Set("status", app.Status); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting status: %s", err))
		}
	}
	if !core.IsNil(app.StatusDetails) {
		statusDetailsMap, err := resourceIbmCodeEngineAppAppStatusToMap(app.StatusDetails)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("status_details", []map[string]interface{}{statusDetailsMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting status_details: %s", err))
		}
	}
	if err = d.Set("etag", response.Headers.Get("Etag")); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting etag: %s", err))
	}

	return nil
}

func resourceIbmCodeEngineAppUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	updateAppOptions := &codeenginev2.UpdateAppOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	updateAppOptions.SetProjectID(parts[0])
	updateAppOptions.SetName(parts[1])

	hasChange := false

	patchVals := &codeenginev2.AppPatch{}
	if d.HasChange("image_reference") || d.HasChange("name") {
		newImageReference := d.Get("image_reference").(string)
		patchVals.ImageReference = &newImageReference
		updateAppOptions.SetName(d.Get("name").(string))
		hasChange = true
	}
	if d.HasChange("image_port") {
		newImagePort := int64(d.Get("image_port").(int))
		patchVals.ImagePort = &newImagePort
		hasChange = true
	}
	if d.HasChange("image_secret") {
		newImageSecret := d.Get("image_secret").(string)
		patchVals.ImageSecret = &newImageSecret
		hasChange = true
	}
	if d.HasChange("managed_domain_mappings") {
		newManagedDomainMappings := d.Get("managed_domain_mappings").(string)
		patchVals.ManagedDomainMappings = &newManagedDomainMappings
		hasChange = true
	}
	if d.HasChange("run_arguments") {
		var runArguments []string
		for _, v := range d.Get("run_arguments").([]interface{}) {
			runArgumentsItem := v.(string)
			runArguments = append(runArguments, runArgumentsItem)
		}
		patchVals.RunArguments = runArguments
		hasChange = true
	}
	if d.HasChange("run_as_user") {
		newRunAsUser := int64(d.Get("run_as_user").(int))
		patchVals.RunAsUser = &newRunAsUser
		hasChange = true
	}
	if d.HasChange("run_commands") {
		var runCommands []string
		for _, v := range d.Get("run_commands").([]interface{}) {
			runCommandsItem := v.(string)
			runCommands = append(runCommands, runCommandsItem)
		}
		patchVals.RunCommands = runCommands
		hasChange = true
	}
	if d.HasChange("run_env_variables") {
		var runEnvVariables []codeenginev2.EnvVarPrototype
		for _, v := range d.Get("run_env_variables").([]interface{}) {
			value := v.(map[string]interface{})
			runEnvVariablesItem, err := resourceIbmCodeEngineAppMapToEnvVarPrototype(value)
			if err != nil {
				return diag.FromErr(err)
			}
			runEnvVariables = append(runEnvVariables, *runEnvVariablesItem)
		}
		patchVals.RunEnvVariables = runEnvVariables
		hasChange = true
	}
	if d.HasChange("run_service_account") {
		newRunServiceAccount := d.Get("run_service_account").(string)
		patchVals.RunServiceAccount = &newRunServiceAccount
		hasChange = true
	}
	if d.HasChange("run_volume_mounts") {
		var runVolumeMounts []codeenginev2.VolumeMountPrototype
		for _, v := range d.Get("run_volume_mounts").([]interface{}) {
			value := v.(map[string]interface{})
			runVolumeMountsItem, err := resourceIbmCodeEngineAppMapToVolumeMountPrototype(value)
			if err != nil {
				return diag.FromErr(err)
			}
			runVolumeMounts = append(runVolumeMounts, *runVolumeMountsItem)
		}
		patchVals.RunVolumeMounts = runVolumeMounts
		hasChange = true
	}
	if d.HasChange("scale_concurrency") {
		newScaleConcurrency := int64(d.Get("scale_concurrency").(int))
		patchVals.ScaleConcurrency = &newScaleConcurrency
		hasChange = true
	}
	if d.HasChange("scale_concurrency_target") {
		newScaleConcurrencyTarget := int64(d.Get("scale_concurrency_target").(int))
		patchVals.ScaleConcurrencyTarget = &newScaleConcurrencyTarget
		hasChange = true
	}
	if d.HasChange("scale_cpu_limit") {
		newScaleCpuLimit := d.Get("scale_cpu_limit").(string)
		patchVals.ScaleCpuLimit = &newScaleCpuLimit
		hasChange = true
	}
	if d.HasChange("scale_ephemeral_storage_limit") {
		newScaleEphemeralStorageLimit := d.Get("scale_ephemeral_storage_limit").(string)
		patchVals.ScaleEphemeralStorageLimit = &newScaleEphemeralStorageLimit
		hasChange = true
	}
	if d.HasChange("scale_initial_instances") {
		newScaleInitialInstances := int64(d.Get("scale_initial_instances").(int))
		patchVals.ScaleInitialInstances = &newScaleInitialInstances
		hasChange = true
	}
	if d.HasChange("scale_max_instances") {
		newScaleMaxInstances := int64(d.Get("scale_max_instances").(int))
		patchVals.ScaleMaxInstances = &newScaleMaxInstances
		hasChange = true
	}
	if d.HasChange("scale_memory_limit") {
		newScaleMemoryLimit := d.Get("scale_memory_limit").(string)
		patchVals.ScaleMemoryLimit = &newScaleMemoryLimit
		hasChange = true
	}
	if d.HasChange("scale_min_instances") {
		newScaleMinInstances := int64(d.Get("scale_min_instances").(int))
		patchVals.ScaleMinInstances = &newScaleMinInstances
		hasChange = true
	}
	if d.HasChange("scale_request_timeout") {
		newScaleRequestTimeout := int64(d.Get("scale_request_timeout").(int))
		patchVals.ScaleRequestTimeout = &newScaleRequestTimeout
		hasChange = true
	}
	updateAppOptions.SetIfMatch(d.Get("etag").(string))

	if hasChange {
		updateAppOptions.App, _ = patchVals.AsPatch()
		_, response, err := codeEngineClient.UpdateAppWithContext(context, updateAppOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateAppWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateAppWithContext failed %s\n%s", err, response))
		}
	}

	_, err = waitForIbmCodeEngineAppUpdate(d, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf(
			"Error waiting for resource IbmCodeEngineApp (%s) to be updated: %s", d.Id(), err))
	}

	return resourceIbmCodeEngineAppRead(context, d, meta)
}

func waitForIbmCodeEngineAppUpdate(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		return false, err
	}
	getAppOptions := &codeenginev2.GetAppOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return false, err
	}

	getAppOptions.SetProjectID(parts[0])
	getAppOptions.SetName(parts[1])

	stateConf := &resource.StateChangeConf{
		Pending: []string{"deploying"},
		Target:  []string{"ready"},
		Refresh: func() (interface{}, string, error) {
			stateObj, response, err := codeEngineClient.GetApp(getAppOptions)
			if err != nil {
				if apiErr, ok := err.(bmxerror.RequestFailure); ok && apiErr.StatusCode() == 404 {
					return nil, "", fmt.Errorf("The instance %s does not exist anymore: %s\n%s", "getAppOptions", err, response)
				}
				return nil, "", err
			}
			failStates := map[string]bool{"failed": true, "warning": true}
			if failStates[*stateObj.Status] {
				return stateObj, *stateObj.Status, fmt.Errorf("The instance %s failed: %s\n%s", "getAppOptions", err, response)
			}
			return stateObj, *stateObj.Status, nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      20 * time.Second,
		MinTimeout: 20 * time.Second,
	}

	return stateConf.WaitForState()
}

func resourceIbmCodeEngineAppDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteAppOptions := &codeenginev2.DeleteAppOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteAppOptions.SetProjectID(parts[0])
	deleteAppOptions.SetName(parts[1])

	response, err := codeEngineClient.DeleteAppWithContext(context, deleteAppOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteAppWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteAppWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIbmCodeEngineAppMapToEnvVarPrototype(modelMap map[string]interface{}) (*codeenginev2.EnvVarPrototype, error) {
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

func resourceIbmCodeEngineAppMapToVolumeMountPrototype(modelMap map[string]interface{}) (*codeenginev2.VolumeMountPrototype, error) {
	model := &codeenginev2.VolumeMountPrototype{}
	model.MountPath = core.StringPtr(modelMap["mount_path"].(string))
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	model.Reference = core.StringPtr(modelMap["reference"].(string))
	model.Type = core.StringPtr(modelMap["type"].(string))
	return model, nil
}

func resourceIbmCodeEngineAppEnvVarToMap(model *codeenginev2.EnvVar) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Key != nil {
		modelMap["key"] = model.Key
	}
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Prefix != nil {
		modelMap["prefix"] = model.Prefix
	}
	if model.Reference != nil {
		modelMap["reference"] = model.Reference
	}
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	if model.Value != nil {
		modelMap["value"] = model.Value
	}
	return modelMap, nil
}

func resourceIbmCodeEngineAppVolumeMountToMap(model *codeenginev2.VolumeMount) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["mount_path"] = model.MountPath
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	modelMap["reference"] = model.Reference
	modelMap["type"] = model.Type
	return modelMap, nil
}

func resourceIbmCodeEngineAppAppStatusToMap(model *codeenginev2.AppStatus) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.LatestCreatedRevision != nil {
		modelMap["latest_created_revision"] = model.LatestCreatedRevision
	}
	if model.LatestReadyRevision != nil {
		modelMap["latest_ready_revision"] = model.LatestReadyRevision
	}
	if model.Reason != nil {
		modelMap["reason"] = model.Reason
	}
	return modelMap, nil
}
