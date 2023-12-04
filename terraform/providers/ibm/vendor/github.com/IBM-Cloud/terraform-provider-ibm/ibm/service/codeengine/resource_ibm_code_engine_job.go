// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package codeengine

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/code-engine-go-sdk/codeenginev2"
	"github.com/IBM/go-sdk-core/v5/core"
)

func ResourceIbmCodeEngineJob() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmCodeEngineJobCreate,
		ReadContext:   resourceIbmCodeEngineJobRead,
		UpdateContext: resourceIbmCodeEngineJobUpdate,
		DeleteContext: resourceIbmCodeEngineJobDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_job", "project_id"),
				Description:  "The ID of the project.",
			},
			"image_reference": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_job", "image_reference"),
				Description:  "The name of the image that is used for this job. The format is `REGISTRY/NAMESPACE/REPOSITORY:TAG` where `REGISTRY` and `TAG` are optional. If `REGISTRY` is not specified, the default is `docker.io`. If `TAG` is not specified, the default is `latest`. If the image reference points to a registry that requires authentication, make sure to also specify the property `image_secret`.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_job", "name"),
				Description:  "The name of the job. Use a name that is unique within the project.",
			},
			"image_secret": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_job", "image_secret"),
				Description:  "The name of the image registry access secret. The image registry access secret is used to authenticate with a private registry when you download the container image. If the image reference points to a registry that requires authentication, the job / job runs will be created but submitted job runs will fail, until this property is provided, too. This property must not be set on a job run, which references a job template.",
			},
			"run_arguments": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				MinItems:    0,
				Description: "Set arguments for the job that are passed to start job run containers. If not specified an empty string array will be applied and the arguments specified by the container image, will be used to start the container.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"run_as_user": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The user ID (UID) to run the application (e.g., 1001).",
			},
			"run_commands": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				MinItems:    0,
				Description: "Set commands for the job that are passed to start job run containers. If not specified an empty string array will be applied and the command specified by the container image, will be used to start the container.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"run_env_variables": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				MinItems:    0,
				Description: "Optional references to config maps, secrets or a literal values.",
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
			"run_mode": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "task",
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_job", "run_mode"),
				Description:  "The mode for runs of the job. Valid values are `task` and `daemon`. In `task` mode, the `scale_max_execution_time` and `scale_retry_limit` properties apply. In `daemon` mode, since there is no timeout and failed instances are restarted indefinitely, the `scale_max_execution_time` and `scale_retry_limit` properties are not allowed.",
			},
			"run_service_account": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "default",
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_job", "run_service_account"),
				Description:  "The name of the service account. For built-in service accounts, you can use the shortened names `manager`, `none`, `reader`, and `writer`. This property must not be set on a job run, which references a job template.",
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
			"scale_array_spec": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "0",
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_job", "scale_array_spec"),
				Description:  "Define a custom set of array indices as comma-separated list containing single values and hyphen-separated ranges like `5,12-14,23,27`. Each instance can pick up its array index via environment variable `JOB_INDEX`. The number of unique array indices specified here determines the number of job instances to run.",
			},
			"scale_cpu_limit": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "1",
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_job", "scale_cpu_limit"),
				Description:  "Optional amount of CPU set for the instance of the job. For valid values see [Supported memory and CPU combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo).",
			},
			"scale_ephemeral_storage_limit": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "400M",
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_job", "scale_ephemeral_storage_limit"),
				Description:  "Optional amount of ephemeral storage to set for the instance of the job. The amount specified as ephemeral storage, must not exceed the amount of `scale_memory_limit`. The units for specifying ephemeral storage are Megabyte (M) or Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information see [Units of measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).",
			},
			"scale_max_execution_time": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     7200,
				Description: "The maximum execution time in seconds for runs of the job. This property can only be specified if `run_mode` is `task`.",
			},
			"scale_memory_limit": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "4G",
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_job", "scale_memory_limit"),
				Description:  "Optional amount of memory set for the instance of the job. For valid values see [Supported memory and CPU combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo). The units for specifying memory are Megabyte (M) or Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information see [Units of measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).",
			},
			"scale_retry_limit": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3,
				Description: "The number of times to rerun an instance of the job before the job is marked as failed. This property can only be specified if `run_mode` is `task`.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the resource was created.",
			},
			"entity_tag": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the job instance, which is used to achieve optimistic locking.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When you provision a new job,  a URL is created identifying the location of the instance.",
			},
			"job_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The identifier of the resource.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the job.",
			},
			"etag": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceIbmCodeEngineJobValidator() *validate.ResourceValidator {
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
			Regexp:                     `^[a-z0-9]([\-a-z0-9]*[a-z0-9])?$`,
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
			Identifier:                 "run_mode",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "daemon, task",
			Regexp:                     `^(task|daemon)$`,
			MinValueLength:             0,
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
			Identifier:                 "scale_array_spec",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^(?:[1-9]\d\d\d\d\d\d|[1-9]\d\d\d\d\d|[1-9]\d\d\d\d|[1-9]\d\d\d|[1-9]\d\d|[1-9]?\d)(?:-(?:[1-9]\d\d\d\d\d\d|[1-9]\d\d\d\d\d|[1-9]\d\d\d\d|[1-9]\d\d\d|[1-9]\d\d|[1-9]?\d))?(?:,(?:[1-9]\d\d\d\d\d\d|[1-9]\d\d\d\d\d|[1-9]\d\d\d\d|[1-9]\d\d\d|[1-9]\d\d|[1-9]?\d)(?:-(?:[1-9]\d\d\d\d\d\d|[1-9]\d\d\d\d\d|[1-9]\d\d\d\d|[1-9]\d\d\d|[1-9]\d\d|[1-9]?\d))?)*$`,
			MinValueLength:             1,
			MaxValueLength:             253,
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

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_code_engine_job", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmCodeEngineJobCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	createJobOptions := &codeenginev2.CreateJobOptions{}

	createJobOptions.SetProjectID(d.Get("project_id").(string))
	createJobOptions.SetImageReference(d.Get("image_reference").(string))
	createJobOptions.SetName(d.Get("name").(string))
	if _, ok := d.GetOk("image_secret"); ok {
		createJobOptions.SetImageSecret(d.Get("image_secret").(string))
	}
	if _, ok := d.GetOk("run_arguments"); ok {
		var runArguments []string
		for _, v := range d.Get("run_arguments").([]interface{}) {
			runArgumentsItem := v.(string)
			runArguments = append(runArguments, runArgumentsItem)
		}
		createJobOptions.SetRunArguments(runArguments)
	}
	if _, ok := d.GetOk("run_as_user"); ok {
		createJobOptions.SetRunAsUser(int64(d.Get("run_as_user").(int)))
	}
	if _, ok := d.GetOk("run_commands"); ok {
		var runCommands []string
		for _, v := range d.Get("run_commands").([]interface{}) {
			runCommandsItem := v.(string)
			runCommands = append(runCommands, runCommandsItem)
		}
		createJobOptions.SetRunCommands(runCommands)
	}
	if _, ok := d.GetOk("run_env_variables"); ok {
		var runEnvVariables []codeenginev2.EnvVarPrototype
		for _, v := range d.Get("run_env_variables").([]interface{}) {
			value := v.(map[string]interface{})
			runEnvVariablesItem, err := resourceIbmCodeEngineJobMapToEnvVarPrototype(value)
			if err != nil {
				return diag.FromErr(err)
			}
			runEnvVariables = append(runEnvVariables, *runEnvVariablesItem)
		}
		createJobOptions.SetRunEnvVariables(runEnvVariables)
	}
	if _, ok := d.GetOk("run_mode"); ok {
		createJobOptions.SetRunMode(d.Get("run_mode").(string))
	}
	if _, ok := d.GetOk("run_service_account"); ok {
		createJobOptions.SetRunServiceAccount(d.Get("run_service_account").(string))
	}
	if _, ok := d.GetOk("run_volume_mounts"); ok {
		var runVolumeMounts []codeenginev2.VolumeMountPrototype
		for _, v := range d.Get("run_volume_mounts").([]interface{}) {
			value := v.(map[string]interface{})
			runVolumeMountsItem, err := resourceIbmCodeEngineJobMapToVolumeMountPrototype(value)
			if err != nil {
				return diag.FromErr(err)
			}
			runVolumeMounts = append(runVolumeMounts, *runVolumeMountsItem)
		}
		createJobOptions.SetRunVolumeMounts(runVolumeMounts)
	}
	if _, ok := d.GetOk("scale_array_spec"); ok {
		createJobOptions.SetScaleArraySpec(d.Get("scale_array_spec").(string))
	}
	if _, ok := d.GetOk("scale_cpu_limit"); ok {
		createJobOptions.SetScaleCpuLimit(d.Get("scale_cpu_limit").(string))
	}
	if _, ok := d.GetOk("scale_ephemeral_storage_limit"); ok {
		createJobOptions.SetScaleEphemeralStorageLimit(d.Get("scale_ephemeral_storage_limit").(string))
	}
	if _, ok := d.GetOk("scale_max_execution_time"); ok {
		createJobOptions.SetScaleMaxExecutionTime(int64(d.Get("scale_max_execution_time").(int)))
	}
	if _, ok := d.GetOk("scale_memory_limit"); ok {
		createJobOptions.SetScaleMemoryLimit(d.Get("scale_memory_limit").(string))
	}
	if _, ok := d.GetOk("scale_retry_limit"); ok {
		createJobOptions.SetScaleRetryLimit(int64(d.Get("scale_retry_limit").(int)))
	}

	job, response, err := codeEngineClient.CreateJobWithContext(context, createJobOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateJobWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateJobWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *createJobOptions.ProjectID, *job.Name))

	return resourceIbmCodeEngineJobRead(context, d, meta)
}

func resourceIbmCodeEngineJobRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	getJobOptions := &codeenginev2.GetJobOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getJobOptions.SetProjectID(parts[0])
	getJobOptions.SetName(parts[1])

	job, response, err := codeEngineClient.GetJobWithContext(context, getJobOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetJobWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetJobWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("project_id", job.ProjectID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting project_id: %s", err))
	}
	if err = d.Set("image_reference", job.ImageReference); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting image_reference: %s", err))
	}
	if err = d.Set("name", job.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if !core.IsNil(job.ImageSecret) {
		if err = d.Set("image_secret", job.ImageSecret); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting image_secret: %s", err))
		}
	}
	if !core.IsNil(job.RunArguments) {
		if err = d.Set("run_arguments", job.RunArguments); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting run_arguments: %s", err))
		}
	}
	if !core.IsNil(job.RunAsUser) {
		if err = d.Set("run_as_user", flex.IntValue(job.RunAsUser)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting run_as_user: %s", err))
		}
	}
	if !core.IsNil(job.RunCommands) {
		if err = d.Set("run_commands", job.RunCommands); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting run_commands: %s", err))
		}
	}
	if !core.IsNil(job.RunEnvVariables) {
		runEnvVariables := []map[string]interface{}{}
		for _, runEnvVariablesItem := range job.RunEnvVariables {
			runEnvVariablesItemMap, err := resourceIbmCodeEngineJobEnvVarToMap(&runEnvVariablesItem)
			if err != nil {
				return diag.FromErr(err)
			}
			runEnvVariables = append(runEnvVariables, runEnvVariablesItemMap)
		}
		if err = d.Set("run_env_variables", runEnvVariables); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting run_env_variables: %s", err))
		}
	}
	if !core.IsNil(job.RunMode) {
		if err = d.Set("run_mode", job.RunMode); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting run_mode: %s", err))
		}
	}
	if !core.IsNil(job.RunServiceAccount) {
		if err = d.Set("run_service_account", job.RunServiceAccount); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting run_service_account: %s", err))
		}
	}
	if !core.IsNil(job.RunVolumeMounts) {
		runVolumeMounts := []map[string]interface{}{}
		for _, runVolumeMountsItem := range job.RunVolumeMounts {
			runVolumeMountsItemMap, err := resourceIbmCodeEngineJobVolumeMountToMap(&runVolumeMountsItem)
			if err != nil {
				return diag.FromErr(err)
			}
			runVolumeMounts = append(runVolumeMounts, runVolumeMountsItemMap)
		}
		if err = d.Set("run_volume_mounts", runVolumeMounts); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting run_volume_mounts: %s", err))
		}
	}
	if !core.IsNil(job.ScaleArraySpec) {
		if err = d.Set("scale_array_spec", job.ScaleArraySpec); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting scale_array_spec: %s", err))
		}
	}
	if !core.IsNil(job.ScaleCpuLimit) {
		if err = d.Set("scale_cpu_limit", job.ScaleCpuLimit); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting scale_cpu_limit: %s", err))
		}
	}
	if !core.IsNil(job.ScaleEphemeralStorageLimit) {
		if err = d.Set("scale_ephemeral_storage_limit", job.ScaleEphemeralStorageLimit); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting scale_ephemeral_storage_limit: %s", err))
		}
	}
	if !core.IsNil(job.ScaleMaxExecutionTime) {
		if err = d.Set("scale_max_execution_time", flex.IntValue(job.ScaleMaxExecutionTime)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting scale_max_execution_time: %s", err))
		}
	}
	if !core.IsNil(job.ScaleMemoryLimit) {
		if err = d.Set("scale_memory_limit", job.ScaleMemoryLimit); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting scale_memory_limit: %s", err))
		}
	}
	if !core.IsNil(job.ScaleRetryLimit) {
		if err = d.Set("scale_retry_limit", flex.IntValue(job.ScaleRetryLimit)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting scale_retry_limit: %s", err))
		}
	}
	if !core.IsNil(job.CreatedAt) {
		if err = d.Set("created_at", job.CreatedAt); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
		}
	}
	if err = d.Set("entity_tag", job.EntityTag); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting entity_tag: %s", err))
	}
	if !core.IsNil(job.Href) {
		if err = d.Set("href", job.Href); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
		}
	}
	if !core.IsNil(job.ID) {
		if err = d.Set("job_id", job.ID); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting job_id: %s", err))
		}
	}
	if !core.IsNil(job.ResourceType) {
		if err = d.Set("resource_type", job.ResourceType); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
		}
	}
	if err = d.Set("etag", response.Headers.Get("Etag")); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting etag: %s", err))
	}

	return nil
}

func resourceIbmCodeEngineJobUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	updateJobOptions := &codeenginev2.UpdateJobOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	updateJobOptions.SetProjectID(parts[0])
	updateJobOptions.SetName(parts[1])

	hasChange := false

	patchVals := &codeenginev2.JobPatch{}
	if d.HasChange("image_reference") || d.HasChange("name") {
		newImageReference := d.Get("image_reference").(string)
		patchVals.ImageReference = &newImageReference
		updateJobOptions.SetName(d.Get("name").(string))
		hasChange = true
	}
	if d.HasChange("image_secret") {
		newImageSecret := d.Get("image_secret").(string)
		patchVals.ImageSecret = &newImageSecret
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
			runEnvVariablesItem, err := resourceIbmCodeEngineJobMapToEnvVarPrototype(value)
			if err != nil {
				return diag.FromErr(err)
			}
			runEnvVariables = append(runEnvVariables, *runEnvVariablesItem)
		}
		patchVals.RunEnvVariables = runEnvVariables
		hasChange = true
	}
	if d.HasChange("run_mode") {
		newRunMode := d.Get("run_mode").(string)
		patchVals.RunMode = &newRunMode
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
			runVolumeMountsItem, err := resourceIbmCodeEngineJobMapToVolumeMountPrototype(value)
			if err != nil {
				return diag.FromErr(err)
			}
			runVolumeMounts = append(runVolumeMounts, *runVolumeMountsItem)
		}
		patchVals.RunVolumeMounts = runVolumeMounts
		hasChange = true
	}
	if d.HasChange("scale_array_spec") {
		newScaleArraySpec := d.Get("scale_array_spec").(string)
		patchVals.ScaleArraySpec = &newScaleArraySpec
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
	if d.HasChange("scale_retry_limit") {
		newScaleRetryLimit := int64(d.Get("scale_retry_limit").(int))
		patchVals.ScaleRetryLimit = &newScaleRetryLimit
		hasChange = true
	}
	updateJobOptions.SetIfMatch(d.Get("etag").(string))

	if hasChange {
		updateJobOptions.Job, _ = patchVals.AsPatch()
		_, response, err := codeEngineClient.UpdateJobWithContext(context, updateJobOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateJobWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateJobWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIbmCodeEngineJobRead(context, d, meta)
}

func resourceIbmCodeEngineJobDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteJobOptions := &codeenginev2.DeleteJobOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteJobOptions.SetProjectID(parts[0])
	deleteJobOptions.SetName(parts[1])

	response, err := codeEngineClient.DeleteJobWithContext(context, deleteJobOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteJobWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteJobWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIbmCodeEngineJobMapToEnvVarPrototype(modelMap map[string]interface{}) (*codeenginev2.EnvVarPrototype, error) {
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

func resourceIbmCodeEngineJobMapToVolumeMountPrototype(modelMap map[string]interface{}) (*codeenginev2.VolumeMountPrototype, error) {
	model := &codeenginev2.VolumeMountPrototype{}
	model.MountPath = core.StringPtr(modelMap["mount_path"].(string))
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	model.Reference = core.StringPtr(modelMap["reference"].(string))
	model.Type = core.StringPtr(modelMap["type"].(string))
	return model, nil
}

func resourceIbmCodeEngineJobEnvVarToMap(model *codeenginev2.EnvVar) (map[string]interface{}, error) {
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

func resourceIbmCodeEngineJobVolumeMountToMap(model *codeenginev2.VolumeMount) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["mount_path"] = model.MountPath
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	modelMap["reference"] = model.Reference
	modelMap["type"] = model.Type
	return modelMap, nil
}
