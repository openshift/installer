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
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_app", "project_id"),
				Description:  "The ID of the project.",
			},
			"image_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     8080,
				Description: "Optional port the app listens on. While the app will always be exposed via port `443` for end users, this port is used to connect to the port that is exposed by the container image.",
			},
			"image_reference": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_app", "image_reference"),
				Description:  "The name of the image that is used for this app. The format is `REGISTRY/NAMESPACE/REPOSITORY:TAG` where `REGISTRY` and `TAG` are optional. If `REGISTRY` is not specified, the default is `docker.io`. If `TAG` is not specified, the default is `latest`. If the image reference points to a registry that requires authentication, make sure to also specify the property `image_secret`.",
			},
			"image_secret": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_app", "image_secret"),
				Description:  "Optional name of the image registry access secret. The image registry access secret is used to authenticate with a private registry when you download the container image. If the image reference points to a registry that requires authentication, the app will be created but cannot reach the ready status, until this property is provided, too.",
			},
			"managed_domain_mappings": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "local_public",
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_app", "managed_domain_mappings"),
				Description:  "Optional value controlling which of the system managed domain mappings will be setup for the application. Valid values are 'local_public', 'local_private' and 'local'. Visibility can only be 'local_private' if the project supports application private visibility.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_app", "name"),
				Description:  "The name of the app.",
			},
			"probe_liveness": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Response model for probes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"failure_threshold": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     1,
							Description: "The number of consecutive, unsuccessful checks for the probe to be considered failed.",
						},
						"initial_delay": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The amount of time in seconds to wait before the first probe check is performed.",
						},
						"interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     10,
							Description: "The amount of time in seconds between probe checks.",
						},
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The path of the HTTP request to the resource. A path is only supported for a probe with a `type` of `http`.",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The port on which to probe the resource.",
						},
						"timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     1,
							Description: "The amount of time in seconds that the probe waits for a response from the application before it times out and fails.",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies whether to use HTTP or TCP for the probe checks. The default is TCP.",
						},
					},
				},
			},
			"probe_readiness": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Response model for probes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"failure_threshold": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     1,
							Description: "The number of consecutive, unsuccessful checks for the probe to be considered failed.",
						},
						"initial_delay": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The amount of time in seconds to wait before the first probe check is performed.",
						},
						"interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     10,
							Description: "The amount of time in seconds between probe checks.",
						},
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The path of the HTTP request to the resource. A path is only supported for a probe with a `type` of `http`.",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The port on which to probe the resource.",
						},
						"timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     1,
							Description: "The amount of time in seconds that the probe waits for a response from the application before it times out and fails.",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies whether to use HTTP or TCP for the probe checks. The default is TCP.",
						},
					},
				},
			},
			"run_arguments": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional arguments for the app that are passed to start the container. If not specified an empty string array will be applied and the arguments specified by the container image, will be used to start the container.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"run_as_user": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Optional user ID (UID) to run the app.",
			},
			"run_commands": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional commands for the app that are passed to start the container. If not specified an empty string array will be applied and the command specified by the container image, will be used to start the container.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"run_env_variables": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "References to config maps, secrets or literal values, which are exposed as environment variables in the application.",
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
			"run_service_account": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "default",
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_app", "run_service_account"),
				Description:  "Optional name of the service account. For built-in service accounts, you can use the shortened names `manager` , `none`, `reader`, and `writer`.",
			},
			"run_volume_mounts": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Mounts of config maps or secrets.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mount_path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The path that should be mounted.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the mount.",
						},
						"reference": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the referenced secret or config map.",
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specify the type of the volume mount. Allowed types are: 'config_map', 'secret'.",
						},
					},
				},
			},
			"scale_concurrency": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     100,
				Description: "Optional maximum number of requests that can be processed concurrently per instance.",
			},
			"scale_concurrency_target": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      100,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_app", "scale_concurrency_target"),
				Description:  "Optional threshold of concurrent requests per instance at which one or more additional instances are created. Use this value to scale up instances based on concurrent number of requests. This option defaults to the value of the `scale_concurrency` option, if not specified.",
			},
			"scale_cpu_limit": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "1",
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_app", "scale_cpu_limit"),
				Description:  "Optional number of CPU set for the instance of the app. For valid values see [Supported memory and CPU combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo).",
			},
			"scale_down_delay": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_app", "scale_down_delay"),
				Description:  "Optional amount of time in seconds that delays the scale-down behavior for an app instance.",
			},
			"scale_ephemeral_storage_limit": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "400M",
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_app", "scale_ephemeral_storage_limit"),
				Description:  "Optional amount of ephemeral storage to set for the instance of the app. The amount specified as ephemeral storage, must not exceed the amount of `scale_memory_limit`. The units for specifying ephemeral storage are Megabyte (M) or Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information see [Units of measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).",
			},
			"scale_initial_instances": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Optional initial number of instances that are created upon app creation or app update.",
			},
			"scale_max_instances": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "Optional maximum number of instances for this app. If you set this value to `0`, this property does not set a upper scaling limit. However, the app scaling is still limited by the project quota for instances. See [Limits and quotas for Code Engine](https://cloud.ibm.com/docs/codeengine?topic=codeengine-limits).",
			},
			"scale_memory_limit": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "4G",
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_app", "scale_memory_limit"),
				Description:  "Optional amount of memory set for the instance of the app. For valid values see [Supported memory and CPU combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo). The units for specifying memory are Megabyte (M) or Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information see [Units of measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).",
			},
			"scale_min_instances": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Optional minimum number of instances for this app. If you set this value to `0`, the app will scale down to zero, if not hit by any request for some time.",
			},
			"scale_request_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     300,
				Description: "Optional amount of time in seconds that is allowed for a running app to respond to a request.",
			},
			"build": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reference to a build that is associated with the application.",
			},
			"build_run": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reference to a build run that is associated with the application.",
			},
			"computed_env_variables": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "References to config maps, secrets or literal values, which are defined and set by Code Engine and are exposed as environment variables in the application.",
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
				Description: "Optional URL to invoke the app. Depending on visibility,  this is accessible publicly or in the private network only. Empty in case 'managed_domain_mappings' is set to 'local'.",
			},
			"endpoint_internal": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL to the app that is only visible within the project.",
			},
			"entity_tag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the app instance, which is used to achieve optimistic locking.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When you provision a new app,  a URL is created identifying the location of the instance.",
			},
			"app_id": {
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
				Description: "The type of the app.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the app.",
			},
			"status_details": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The detailed status of the application.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"latest_created_revision": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Latest app revision that has been created.",
						},
						"latest_ready_revision": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Latest app revision that reached a ready state.",
						},
						"reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Optional information to provide more context in case of a 'failed' or 'warning' status.",
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
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-z]([-a-z0-9]*[a-z0-9])?$`,
			MinValueLength:             1,
			MaxValueLength:             63,
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
			Identifier:                 "scale_concurrency_target",
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			Optional:                   true,
			MinValue:                   "1",
			MaxValue:                   "1000",
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
			MaxValue:                   "3600",
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
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
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
	if _, ok := d.GetOk("probe_liveness"); ok {
		probeLivenessModel, err := ResourceIbmCodeEngineAppMapToProbePrototype(d.Get("probe_liveness.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "create", "parse-probe_liveness").GetDiag()
		}
		createAppOptions.SetProbeLiveness(probeLivenessModel)
	}
	if _, ok := d.GetOk("probe_readiness"); ok {
		probeReadinessModel, err := ResourceIbmCodeEngineAppMapToProbePrototype(d.Get("probe_readiness.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "create", "parse-probe_readiness").GetDiag()
		}
		createAppOptions.SetProbeReadiness(probeReadinessModel)
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
			runEnvVariablesItem, err := ResourceIbmCodeEngineAppMapToEnvVarPrototype(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "create", "parse-run_env_variables").GetDiag()
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
			runVolumeMountsItem, err := ResourceIbmCodeEngineAppMapToVolumeMountPrototype(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "create", "parse-run_volume_mounts").GetDiag()
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
	if _, ok := d.GetOk("scale_down_delay"); ok {
		createAppOptions.SetScaleDownDelay(int64(d.Get("scale_down_delay").(int)))
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

	app, _, err := codeEngineClient.CreateAppWithContext(context, createAppOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateAppWithContext failed: %s", err.Error()), "ibm_code_engine_app", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *createAppOptions.ProjectID, *app.Name))

	_, err = waitForIbmCodeEngineAppCreate(d, meta)
	if err != nil {
		errMsg := fmt.Sprintf("Error waiting for resource IbmCodeEngineApp (%s) to be created: %s", d.Id(), err)
		return flex.DiscriminatedTerraformErrorf(err, errMsg, "ibm_code_engine_app", "create", "wait-for-state").GetDiag()
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
				if sdkErr, ok := err.(*core.SDKProblem); ok && response.GetStatusCode() == 404 {
					sdkErr.Summary = fmt.Sprintf("The instance %s does not exist anymore: %s", "getAppOptions", err)
					return nil, "", sdkErr
				}
				return nil, "", err
			}
			failStates := map[string]bool{"failure": true, "failed": true}
			if failStates[*stateObj.Status] {
				return stateObj, *stateObj.Status, fmt.Errorf("The instance %s failed: %s", "getAppOptions", err)
			}
			return stateObj, *stateObj.Status, nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      60 * time.Second,
		MinTimeout: 60 * time.Second,
	}

	return stateConf.WaitForStateContext(context.Background())
}

func resourceIbmCodeEngineAppRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getAppOptions := &codeenginev2.GetAppOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "sep-id-parts").GetDiag()
	}

	getAppOptions.SetProjectID(parts[0])
	getAppOptions.SetName(parts[1])

	app, response, err := codeEngineClient.GetAppWithContext(context, getAppOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetAppWithContext failed: %s", err.Error()), "ibm_code_engine_app", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("project_id", app.ProjectID); err != nil {
		err = fmt.Errorf("Error setting project_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-project_id").GetDiag()
	}
	if !core.IsNil(app.ImagePort) {
		if err = d.Set("image_port", flex.IntValue(app.ImagePort)); err != nil {
			err = fmt.Errorf("Error setting image_port: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-image_port").GetDiag()
		}
	}
	if err = d.Set("image_reference", app.ImageReference); err != nil {
		err = fmt.Errorf("Error setting image_reference: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-image_reference").GetDiag()
	}
	if !core.IsNil(app.ImageSecret) {
		if err = d.Set("image_secret", app.ImageSecret); err != nil {
			err = fmt.Errorf("Error setting image_secret: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-image_secret").GetDiag()
		}
	}
	if !core.IsNil(app.ManagedDomainMappings) {
		if err = d.Set("managed_domain_mappings", app.ManagedDomainMappings); err != nil {
			err = fmt.Errorf("Error setting managed_domain_mappings: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-managed_domain_mappings").GetDiag()
		}
	}
	if err = d.Set("name", app.Name); err != nil {
		err = fmt.Errorf("Error setting name: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-name").GetDiag()
	}
	if !core.IsNil(app.ProbeLiveness) {
		probeLivenessMap, err := ResourceIbmCodeEngineAppProbeToMap(app.ProbeLiveness)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "probe_liveness-to-map").GetDiag()
		}
		if err = d.Set("probe_liveness", []map[string]interface{}{probeLivenessMap}); err != nil {
			err = fmt.Errorf("Error setting probe_liveness: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-probe_liveness").GetDiag()
		}
	}
	if !core.IsNil(app.ProbeReadiness) {
		probeReadinessMap, err := ResourceIbmCodeEngineAppProbeToMap(app.ProbeReadiness)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "probe_readiness-to-map").GetDiag()
		}
		if err = d.Set("probe_readiness", []map[string]interface{}{probeReadinessMap}); err != nil {
			err = fmt.Errorf("Error setting probe_readiness: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-probe_readiness").GetDiag()
		}
	}
	if !core.IsNil(app.RunArguments) {
		if err = d.Set("run_arguments", app.RunArguments); err != nil {
			err = fmt.Errorf("Error setting run_arguments: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-run_arguments").GetDiag()
		}
	}
	if !core.IsNil(app.RunAsUser) {
		if err = d.Set("run_as_user", flex.IntValue(app.RunAsUser)); err != nil {
			err = fmt.Errorf("Error setting run_as_user: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-run_as_user").GetDiag()
		}
	}
	if !core.IsNil(app.RunCommands) {
		if err = d.Set("run_commands", app.RunCommands); err != nil {
			err = fmt.Errorf("Error setting run_commands: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-run_commands").GetDiag()
		}
	}
	if !core.IsNil(app.RunEnvVariables) {
		runEnvVariables := []map[string]interface{}{}
		for _, runEnvVariablesItem := range app.RunEnvVariables {
			runEnvVariablesItemMap, err := ResourceIbmCodeEngineAppEnvVarToMap(&runEnvVariablesItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "run_env_variables-to-map").GetDiag()
			}
			runEnvVariables = append(runEnvVariables, runEnvVariablesItemMap)
		}
		if err = d.Set("run_env_variables", runEnvVariables); err != nil {
			err = fmt.Errorf("Error setting run_env_variables: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-run_env_variables").GetDiag()
		}
	}
	if !core.IsNil(app.RunServiceAccount) {
		if err = d.Set("run_service_account", app.RunServiceAccount); err != nil {
			err = fmt.Errorf("Error setting run_service_account: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-run_service_account").GetDiag()
		}
	}
	if !core.IsNil(app.RunVolumeMounts) {
		runVolumeMounts := []map[string]interface{}{}
		for _, runVolumeMountsItem := range app.RunVolumeMounts {
			runVolumeMountsItemMap, err := ResourceIbmCodeEngineAppVolumeMountToMap(&runVolumeMountsItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "run_volume_mounts-to-map").GetDiag()
			}
			runVolumeMounts = append(runVolumeMounts, runVolumeMountsItemMap)
		}
		if err = d.Set("run_volume_mounts", runVolumeMounts); err != nil {
			err = fmt.Errorf("Error setting run_volume_mounts: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-run_volume_mounts").GetDiag()
		}
	}
	if !core.IsNil(app.ScaleConcurrency) {
		if err = d.Set("scale_concurrency", flex.IntValue(app.ScaleConcurrency)); err != nil {
			err = fmt.Errorf("Error setting scale_concurrency: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-scale_concurrency").GetDiag()
		}
	}
	if !core.IsNil(app.ScaleConcurrencyTarget) {
		if err = d.Set("scale_concurrency_target", flex.IntValue(app.ScaleConcurrencyTarget)); err != nil {
			err = fmt.Errorf("Error setting scale_concurrency_target: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-scale_concurrency_target").GetDiag()
		}
	}
	if !core.IsNil(app.ScaleCpuLimit) {
		if err = d.Set("scale_cpu_limit", app.ScaleCpuLimit); err != nil {
			err = fmt.Errorf("Error setting scale_cpu_limit: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-scale_cpu_limit").GetDiag()
		}
	}
	if !core.IsNil(app.ScaleDownDelay) {
		if err = d.Set("scale_down_delay", flex.IntValue(app.ScaleDownDelay)); err != nil {
			err = fmt.Errorf("Error setting scale_down_delay: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-scale_down_delay").GetDiag()
		}
	}
	if !core.IsNil(app.ScaleEphemeralStorageLimit) {
		if err = d.Set("scale_ephemeral_storage_limit", app.ScaleEphemeralStorageLimit); err != nil {
			err = fmt.Errorf("Error setting scale_ephemeral_storage_limit: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-scale_ephemeral_storage_limit").GetDiag()
		}
	}
	if !core.IsNil(app.ScaleInitialInstances) {
		if err = d.Set("scale_initial_instances", flex.IntValue(app.ScaleInitialInstances)); err != nil {
			err = fmt.Errorf("Error setting scale_initial_instances: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-scale_initial_instances").GetDiag()
		}
	}
	if !core.IsNil(app.ScaleMaxInstances) {
		if err = d.Set("scale_max_instances", flex.IntValue(app.ScaleMaxInstances)); err != nil {
			err = fmt.Errorf("Error setting scale_max_instances: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-scale_max_instances").GetDiag()
		}
	}
	if !core.IsNil(app.ScaleMemoryLimit) {
		if err = d.Set("scale_memory_limit", app.ScaleMemoryLimit); err != nil {
			err = fmt.Errorf("Error setting scale_memory_limit: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-scale_memory_limit").GetDiag()
		}
	}
	if !core.IsNil(app.ScaleMinInstances) {
		if err = d.Set("scale_min_instances", flex.IntValue(app.ScaleMinInstances)); err != nil {
			err = fmt.Errorf("Error setting scale_min_instances: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-scale_min_instances").GetDiag()
		}
	}
	if !core.IsNil(app.ScaleRequestTimeout) {
		if err = d.Set("scale_request_timeout", flex.IntValue(app.ScaleRequestTimeout)); err != nil {
			err = fmt.Errorf("Error setting scale_request_timeout: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-scale_request_timeout").GetDiag()
		}
	}
	if !core.IsNil(app.Build) {
		if err = d.Set("build", app.Build); err != nil {
			err = fmt.Errorf("Error setting build: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-build").GetDiag()
		}
	}
	if !core.IsNil(app.BuildRun) {
		if err = d.Set("build_run", app.BuildRun); err != nil {
			err = fmt.Errorf("Error setting build_run: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-build_run").GetDiag()
		}
	}
	if !core.IsNil(app.ComputedEnvVariables) {
		computedEnvVariables := []map[string]interface{}{}
		for _, computedEnvVariablesItem := range app.ComputedEnvVariables {
			computedEnvVariablesItemMap, err := ResourceIbmCodeEngineAppEnvVarToMap(&computedEnvVariablesItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "computed_env_variables-to-map").GetDiag()
			}
			computedEnvVariables = append(computedEnvVariables, computedEnvVariablesItemMap)
		}
		if err = d.Set("computed_env_variables", computedEnvVariables); err != nil {
			err = fmt.Errorf("Error setting computed_env_variables: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-computed_env_variables").GetDiag()
		}
	}
	if !core.IsNil(app.CreatedAt) {
		if err = d.Set("created_at", app.CreatedAt); err != nil {
			err = fmt.Errorf("Error setting created_at: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-created_at").GetDiag()
		}
	}
	if !core.IsNil(app.Endpoint) {
		if err = d.Set("endpoint", app.Endpoint); err != nil {
			err = fmt.Errorf("Error setting endpoint: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-endpoint").GetDiag()
		}
	}
	if !core.IsNil(app.EndpointInternal) {
		if err = d.Set("endpoint_internal", app.EndpointInternal); err != nil {
			err = fmt.Errorf("Error setting endpoint_internal: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-endpoint_internal").GetDiag()
		}
	}
	if err = d.Set("entity_tag", app.EntityTag); err != nil {
		err = fmt.Errorf("Error setting entity_tag: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-entity_tag").GetDiag()
	}
	if !core.IsNil(app.Href) {
		if err = d.Set("href", app.Href); err != nil {
			err = fmt.Errorf("Error setting href: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-href").GetDiag()
		}
	}
	if !core.IsNil(app.ID) {
		if err = d.Set("app_id", app.ID); err != nil {
			err = fmt.Errorf("Error setting app_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-app_id").GetDiag()
		}
	}
	if !core.IsNil(app.Region) {
		if err = d.Set("region", app.Region); err != nil {
			err = fmt.Errorf("Error setting region: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-region").GetDiag()
		}
	}
	if !core.IsNil(app.ResourceType) {
		if err = d.Set("resource_type", app.ResourceType); err != nil {
			err = fmt.Errorf("Error setting resource_type: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-resource_type").GetDiag()
		}
	}
	if !core.IsNil(app.Status) {
		if err = d.Set("status", app.Status); err != nil {
			err = fmt.Errorf("Error setting status: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-status").GetDiag()
		}
	}
	if !core.IsNil(app.StatusDetails) {
		statusDetailsMap, err := ResourceIbmCodeEngineAppAppStatusToMap(app.StatusDetails)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "status_details-to-map").GetDiag()
		}
		if err = d.Set("status_details", []map[string]interface{}{statusDetailsMap}); err != nil {
			err = fmt.Errorf("Error setting status_details: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "read", "set-status_details").GetDiag()
		}
	}
	if err = d.Set("etag", response.Headers.Get("Etag")); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting etag: %s", err), "ibm_code_engine_app", "read", "set-etag").GetDiag()
	}

	return nil
}

func resourceIbmCodeEngineAppUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateAppOptions := &codeenginev2.UpdateAppOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "update", "sep-id-parts").GetDiag()
	}

	updateAppOptions.SetProjectID(parts[0])
	updateAppOptions.SetName(parts[1])

	hasChange := false

	patchVals := &codeenginev2.AppPatch{}
	if d.HasChange("project_id") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "project_id")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_code_engine_app", "update", "project_id-forces-new").GetDiag()
	}
	if d.HasChange("image_port") {
		newImagePort := int64(d.Get("image_port").(int))
		patchVals.ImagePort = &newImagePort
		hasChange = true
	}
	if d.HasChange("image_reference") {
		newImageReference := d.Get("image_reference").(string)
		patchVals.ImageReference = &newImageReference
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
	if d.HasChange("probe_liveness") {
		probeLiveness, err := ResourceIbmCodeEngineAppMapToProbePrototype(d.Get("probe_liveness.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "update", "parse-probe_liveness").GetDiag()
		}
		patchVals.ProbeLiveness = probeLiveness
		hasChange = true
	}
	if d.HasChange("probe_readiness") {
		probeReadiness, err := ResourceIbmCodeEngineAppMapToProbePrototype(d.Get("probe_readiness.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "update", "parse-probe_readiness").GetDiag()
		}
		patchVals.ProbeReadiness = probeReadiness
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
			runEnvVariablesItem, err := ResourceIbmCodeEngineAppMapToEnvVarPrototype(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "update", "parse-run_env_variables").GetDiag()
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
			runVolumeMountsItem, err := ResourceIbmCodeEngineAppMapToVolumeMountPrototype(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "update", "parse-run_volume_mounts").GetDiag()
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
	if d.HasChange("scale_down_delay") {
		newScaleDownDelay := int64(d.Get("scale_down_delay").(int))
		patchVals.ScaleDownDelay = &newScaleDownDelay
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
		// Fields with `nil` values are omitted from the generic map,
		// so we need to re-add them to support removing arguments
		// in merge-patch operations sent to the service.
		updateAppOptions.App = ResourceIbmCodeEngineAppAppPatchAsPatch(patchVals, d)

		_, _, err = codeEngineClient.UpdateAppWithContext(context, updateAppOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateAppWithContext failed: %s", err.Error()), "ibm_code_engine_app", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIbmCodeEngineAppRead(context, d, meta)
}

func resourceIbmCodeEngineAppDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteAppOptions := &codeenginev2.DeleteAppOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_app", "delete", "sep-id-parts").GetDiag()
	}

	deleteAppOptions.SetProjectID(parts[0])
	deleteAppOptions.SetName(parts[1])

	_, err = codeEngineClient.DeleteAppWithContext(context, deleteAppOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteAppWithContext failed: %s", err.Error()), "ibm_code_engine_app", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIbmCodeEngineAppMapToProbePrototype(modelMap map[string]interface{}) (*codeenginev2.ProbePrototype, error) {
	model := &codeenginev2.ProbePrototype{}
	if modelMap["failure_threshold"] != nil {
		model.FailureThreshold = core.Int64Ptr(int64(modelMap["failure_threshold"].(int)))
	}
	if modelMap["initial_delay"] != nil {
		model.InitialDelay = core.Int64Ptr(int64(modelMap["initial_delay"].(int)))
	}
	if modelMap["interval"] != nil {
		model.Interval = core.Int64Ptr(int64(modelMap["interval"].(int)))
	}
	if modelMap["path"] != nil && modelMap["path"].(string) != "" {
		model.Path = core.StringPtr(modelMap["path"].(string))
	}
	if modelMap["port"] != nil {
		model.Port = core.Int64Ptr(int64(modelMap["port"].(int)))
	}
	if modelMap["timeout"] != nil {
		model.Timeout = core.Int64Ptr(int64(modelMap["timeout"].(int)))
	}
	if modelMap["type"] != nil && modelMap["type"].(string) != "" {
		model.Type = core.StringPtr(modelMap["type"].(string))
	}
	return model, nil
}

func ResourceIbmCodeEngineAppMapToEnvVarPrototype(modelMap map[string]interface{}) (*codeenginev2.EnvVarPrototype, error) {
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

func ResourceIbmCodeEngineAppMapToVolumeMountPrototype(modelMap map[string]interface{}) (*codeenginev2.VolumeMountPrototype, error) {
	model := &codeenginev2.VolumeMountPrototype{}
	model.MountPath = core.StringPtr(modelMap["mount_path"].(string))
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	model.Reference = core.StringPtr(modelMap["reference"].(string))
	model.Type = core.StringPtr(modelMap["type"].(string))
	return model, nil
}

func ResourceIbmCodeEngineAppProbeToMap(model *codeenginev2.Probe) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.FailureThreshold != nil {
		modelMap["failure_threshold"] = flex.IntValue(model.FailureThreshold)
	}
	if model.InitialDelay != nil {
		modelMap["initial_delay"] = flex.IntValue(model.InitialDelay)
	}
	if model.Interval != nil {
		modelMap["interval"] = flex.IntValue(model.Interval)
	}
	if model.Path != nil {
		modelMap["path"] = *model.Path
	}
	if model.Port != nil {
		modelMap["port"] = flex.IntValue(model.Port)
	}
	if model.Timeout != nil {
		modelMap["timeout"] = flex.IntValue(model.Timeout)
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	return modelMap, nil
}

func ResourceIbmCodeEngineAppEnvVarToMap(model *codeenginev2.EnvVar) (map[string]interface{}, error) {
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

func ResourceIbmCodeEngineAppVolumeMountToMap(model *codeenginev2.VolumeMount) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["mount_path"] = *model.MountPath
	modelMap["name"] = *model.Name
	modelMap["reference"] = *model.Reference
	modelMap["type"] = *model.Type
	return modelMap, nil
}

func ResourceIbmCodeEngineAppAppStatusToMap(model *codeenginev2.AppStatus) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.LatestCreatedRevision != nil {
		modelMap["latest_created_revision"] = *model.LatestCreatedRevision
	}
	if model.LatestReadyRevision != nil {
		modelMap["latest_ready_revision"] = *model.LatestReadyRevision
	}
	if model.Reason != nil {
		modelMap["reason"] = *model.Reason
	}
	return modelMap, nil
}

func ResourceIbmCodeEngineAppAppPatchAsPatch(patchVals *codeenginev2.AppPatch, d *schema.ResourceData) map[string]interface{} {
	patch, _ := patchVals.AsPatch()
	var path string

	path = "image_port"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["image_port"] = nil
	}
	path = "image_reference"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["image_reference"] = nil
	}
	path = "image_secret"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["image_secret"] = nil
	}
	path = "managed_domain_mappings"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["managed_domain_mappings"] = nil
	}
	path = "probe_liveness"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["probe_liveness"] = nil
	}
	path = "probe_readiness"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["probe_readiness"] = nil
	}
	path = "run_arguments"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["run_arguments"] = nil
	}
	path = "run_as_user"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["run_as_user"] = nil
	}
	path = "run_commands"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["run_commands"] = nil
	}
	path = "run_env_variables"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		runEnvVariables := []map[string]interface{}{}
		patch["run_env_variables"] = runEnvVariables
	}
	path = "run_service_account"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["run_service_account"] = nil
	}
	path = "run_volume_mounts"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		runVolumeMounts := []map[string]interface{}{}
		patch["run_volume_mounts"] = runVolumeMounts
	}
	path = "scale_concurrency"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["scale_concurrency"] = nil
	}
	path = "scale_concurrency_target"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["scale_concurrency_target"] = nil
	}
	path = "scale_cpu_limit"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["scale_cpu_limit"] = nil
	}
	path = "scale_down_delay"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["scale_down_delay"] = nil
	}
	path = "scale_ephemeral_storage_limit"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["scale_ephemeral_storage_limit"] = nil
	}
	path = "scale_initial_instances"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["scale_initial_instances"] = nil
	}
	path = "scale_max_instances"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["scale_max_instances"] = nil
	}
	path = "scale_memory_limit"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["scale_memory_limit"] = nil
	}
	path = "scale_min_instances"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["scale_min_instances"] = nil
	}
	path = "scale_request_timeout"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["scale_request_timeout"] = nil
	}

	return patch
}

func ResourceIbmCodeEngineAppVolumeMountPrototypeAsPatch(patch map[string]interface{}, d *schema.ResourceData) {
	var path string

	path = "run_volume_mounts.0.name"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["name"] = nil
	}
}

func ResourceIbmCodeEngineAppEnvVarPrototypeAsPatch(patch map[string]interface{}, d *schema.ResourceData) {
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

func ResourceIbmCodeEngineAppProbePrototypeAsPatch(patch map[string]interface{}, d *schema.ResourceData) {
	var path string

	path = "probe_liveness.0.failure_threshold"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["failure_threshold"] = nil
	}
	path = "probe_liveness.0.initial_delay"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["initial_delay"] = nil
	}
	path = "probe_liveness.0.interval"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["interval"] = nil
	}
	path = "probe_liveness.0.path"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["path"] = nil
	}
	path = "probe_liveness.0.port"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["port"] = nil
	}
	path = "probe_liveness.0.timeout"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["timeout"] = nil
	}
	path = "probe_liveness.0.type"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["type"] = nil
	}
}
