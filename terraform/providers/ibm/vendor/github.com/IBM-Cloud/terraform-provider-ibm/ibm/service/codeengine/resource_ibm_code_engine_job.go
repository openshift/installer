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

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/code-engine-go-sdk/codeenginev2"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIbmCodeEngineJob() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmCodeEngineJobCreate,
		ReadContext:   resourceIbmCodeEngineJobRead,
		UpdateContext: resourceIbmCodeEngineJobUpdate,
		DeleteContext: resourceIbmCodeEngineJobDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_job", "project_id"),
				Description:  "The ID of the project.",
			},
			"image_reference": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_job", "image_reference"),
				Description:  "The name of the image that is used for this job. The format is `REGISTRY/NAMESPACE/REPOSITORY:TAG` where `REGISTRY` and `TAG` are optional. If `REGISTRY` is not specified, the default is `docker.io`. If `TAG` is not specified, the default is `latest`. If the image reference points to a registry that requires authentication, make sure to also specify the property `image_secret`.",
			},
			"image_secret": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_job", "image_secret"),
				Description:  "The name of the image registry access secret. The image registry access secret is used to authenticate with a private registry when you download the container image. If the image reference points to a registry that requires authentication, the job / job runs will be created but submitted job runs will fail, until this property is provided, too. This property must not be set on a job run, which references a job template.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_job", "name"),
				Description:  "The name of the job.",
			},
			"run_arguments": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Set arguments for the job that are passed to start job run containers. If not specified an empty string array will be applied and the arguments specified by the container image, will be used to start the container.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"run_as_user": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The user ID (UID) to run the job.",
			},
			"run_commands": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Set commands for the job that are passed to start job run containers. If not specified an empty string array will be applied and the command specified by the container image, will be used to start the container.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"run_env_variables": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "References to config maps, secrets or literal values, which are exposed as environment variables in the job run.",
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
			"run_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "task",
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_job", "run_mode"),
				Description:  "The mode for runs of the job. Valid values are `task` and `daemon`. In `task` mode, the `max_execution_time` and `retry_limit` properties apply. In `daemon` mode, since there is no timeout and failed instances are restarted indefinitely, the `max_execution_time` and `retry_limit` properties are not allowed.",
			},
			"run_service_account": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "default",
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_job", "run_service_account"),
				Description:  "The name of the service account. For built-in service accounts, you can use the shortened names `manager`, `none`, `reader`, and `writer`. This property must not be set on a job run, which references a job template.",
			},
			"run_volume_mounts": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional mounts of config maps or secrets.",
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
			"scale_array_spec": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "0",
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_job", "scale_array_spec"),
				Description:  "Define a custom set of array indices as a comma-separated list containing single values and hyphen-separated ranges, such as  5,12-14,23,27. Each instance gets its array index value from the environment variable JOB_INDEX. The number of unique array indices that you specify with this parameter determines the number of job instances to run.",
			},
			"scale_cpu_limit": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "1",
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_job", "scale_cpu_limit"),
				Description:  "Optional amount of CPU set for the instance of the job. For valid values see [Supported memory and CPU combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo).",
			},
			"scale_ephemeral_storage_limit": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "400M",
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_job", "scale_ephemeral_storage_limit"),
				Description:  "Optional amount of ephemeral storage to set for the instance of the job. The amount specified as ephemeral storage, must not exceed the amount of `scale_memory_limit`. The units for specifying ephemeral storage are Megabyte (M) or Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information see [Units of measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).",
			},
			"scale_max_execution_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     7200,
				Description: "The maximum execution time in seconds for runs of the job. This property can only be specified if `run_mode` is `task`.",
			},
			"scale_memory_limit": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "4G",
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_job", "scale_memory_limit"),
				Description:  "Optional amount of memory set for the instance of the job. For valid values see [Supported memory and CPU combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo). The units for specifying memory are Megabyte (M) or Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information see [Units of measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).",
			},
			"scale_retry_limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3,
				Description: "The number of times to rerun an instance of the job before the job is marked as failed. This property can only be specified if `run_mode` is `task`.",
			},
			"build": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reference to a build that is associated with the job.",
			},
			"build_run": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reference to a build run that is associated with the job.",
			},
			"computed_env_variables": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "References to config maps, secrets or literal values, which are defined and set by Code Engine and are exposed as environment variables in the job run.",
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
			"entity_tag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the job instance, which is used to achieve optimistic locking.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When you provision a new job,  a URL is created identifying the location of the instance.",
			},
			"job_id": {
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
				Description: "The type of the job.",
			},
			"etag": {
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
			Identifier:                 "image_secret",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[a-z0-9]([\-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([\-a-z0-9]*[a-z0-9])?)*$`,
			MinValueLength:             1,
			MaxValueLength:             253,
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
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
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
			runEnvVariablesItem, err := ResourceIbmCodeEngineJobMapToEnvVarPrototype(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "create", "parse-run_env_variables").GetDiag()
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
			runVolumeMountsItem, err := ResourceIbmCodeEngineJobMapToVolumeMountPrototype(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "create", "parse-run_volume_mounts").GetDiag()
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

	job, _, err := codeEngineClient.CreateJobWithContext(context, createJobOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateJobWithContext failed: %s", err.Error()), "ibm_code_engine_job", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *createJobOptions.ProjectID, *job.Name))

	return resourceIbmCodeEngineJobRead(context, d, meta)
}

func resourceIbmCodeEngineJobRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getJobOptions := &codeenginev2.GetJobOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "read", "sep-id-parts").GetDiag()
	}

	getJobOptions.SetProjectID(parts[0])
	getJobOptions.SetName(parts[1])

	job, response, err := codeEngineClient.GetJobWithContext(context, getJobOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetJobWithContext failed: %s", err.Error()), "ibm_code_engine_job", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("project_id", job.ProjectID); err != nil {
		err = fmt.Errorf("Error setting project_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "read", "set-project_id").GetDiag()
	}
	if err = d.Set("image_reference", job.ImageReference); err != nil {
		err = fmt.Errorf("Error setting image_reference: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "read", "set-image_reference").GetDiag()
	}
	if !core.IsNil(job.ImageSecret) {
		if err = d.Set("image_secret", job.ImageSecret); err != nil {
			err = fmt.Errorf("Error setting image_secret: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "read", "set-image_secret").GetDiag()
		}
	}
	if err = d.Set("name", job.Name); err != nil {
		err = fmt.Errorf("Error setting name: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "read", "set-name").GetDiag()
	}
	if !core.IsNil(job.RunArguments) {
		if err = d.Set("run_arguments", job.RunArguments); err != nil {
			err = fmt.Errorf("Error setting run_arguments: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "read", "set-run_arguments").GetDiag()
		}
	}
	if !core.IsNil(job.RunAsUser) {
		if err = d.Set("run_as_user", flex.IntValue(job.RunAsUser)); err != nil {
			err = fmt.Errorf("Error setting run_as_user: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "read", "set-run_as_user").GetDiag()
		}
	}
	if !core.IsNil(job.RunCommands) {
		if err = d.Set("run_commands", job.RunCommands); err != nil {
			err = fmt.Errorf("Error setting run_commands: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "read", "set-run_commands").GetDiag()
		}
	}
	if !core.IsNil(job.RunEnvVariables) {
		runEnvVariables := []map[string]interface{}{}
		for _, runEnvVariablesItem := range job.RunEnvVariables {
			runEnvVariablesItemMap, err := ResourceIbmCodeEngineJobEnvVarToMap(&runEnvVariablesItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "read", "run_env_variables-to-map").GetDiag()
			}
			runEnvVariables = append(runEnvVariables, runEnvVariablesItemMap)
		}
		if err = d.Set("run_env_variables", runEnvVariables); err != nil {
			err = fmt.Errorf("Error setting run_env_variables: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "read", "set-run_env_variables").GetDiag()
		}
	}
	if !core.IsNil(job.RunMode) {
		if err = d.Set("run_mode", job.RunMode); err != nil {
			err = fmt.Errorf("Error setting run_mode: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "read", "set-run_mode").GetDiag()
		}
	}
	if !core.IsNil(job.RunServiceAccount) {
		if err = d.Set("run_service_account", job.RunServiceAccount); err != nil {
			err = fmt.Errorf("Error setting run_service_account: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "read", "set-run_service_account").GetDiag()
		}
	}
	if !core.IsNil(job.RunVolumeMounts) {
		runVolumeMounts := []map[string]interface{}{}
		for _, runVolumeMountsItem := range job.RunVolumeMounts {
			runVolumeMountsItemMap, err := ResourceIbmCodeEngineJobVolumeMountToMap(&runVolumeMountsItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "read", "run_volume_mounts-to-map").GetDiag()
			}
			runVolumeMounts = append(runVolumeMounts, runVolumeMountsItemMap)
		}
		if err = d.Set("run_volume_mounts", runVolumeMounts); err != nil {
			err = fmt.Errorf("Error setting run_volume_mounts: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "read", "set-run_volume_mounts").GetDiag()
		}
	}
	if !core.IsNil(job.ScaleArraySpec) {
		if err = d.Set("scale_array_spec", job.ScaleArraySpec); err != nil {
			err = fmt.Errorf("Error setting scale_array_spec: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "read", "set-scale_array_spec").GetDiag()
		}
	}
	if !core.IsNil(job.ScaleCpuLimit) {
		if err = d.Set("scale_cpu_limit", job.ScaleCpuLimit); err != nil {
			err = fmt.Errorf("Error setting scale_cpu_limit: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "read", "set-scale_cpu_limit").GetDiag()
		}
	}
	if !core.IsNil(job.ScaleEphemeralStorageLimit) {
		if err = d.Set("scale_ephemeral_storage_limit", job.ScaleEphemeralStorageLimit); err != nil {
			err = fmt.Errorf("Error setting scale_ephemeral_storage_limit: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "read", "set-scale_ephemeral_storage_limit").GetDiag()
		}
	}
	if !core.IsNil(job.ScaleMaxExecutionTime) {
		if err = d.Set("scale_max_execution_time", flex.IntValue(job.ScaleMaxExecutionTime)); err != nil {
			err = fmt.Errorf("Error setting scale_max_execution_time: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "read", "set-scale_max_execution_time").GetDiag()
		}
	}
	if !core.IsNil(job.ScaleMemoryLimit) {
		if err = d.Set("scale_memory_limit", job.ScaleMemoryLimit); err != nil {
			err = fmt.Errorf("Error setting scale_memory_limit: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "read", "set-scale_memory_limit").GetDiag()
		}
	}
	if !core.IsNil(job.ScaleRetryLimit) {
		if err = d.Set("scale_retry_limit", flex.IntValue(job.ScaleRetryLimit)); err != nil {
			err = fmt.Errorf("Error setting scale_retry_limit: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "read", "set-scale_retry_limit").GetDiag()
		}
	}
	if !core.IsNil(job.Build) {
		if err = d.Set("build", job.Build); err != nil {
			err = fmt.Errorf("Error setting build: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "read", "set-build").GetDiag()
		}
	}
	if !core.IsNil(job.BuildRun) {
		if err = d.Set("build_run", job.BuildRun); err != nil {
			err = fmt.Errorf("Error setting build_run: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "read", "set-build_run").GetDiag()
		}
	}
	if !core.IsNil(job.ComputedEnvVariables) {
		computedEnvVariables := []map[string]interface{}{}
		for _, computedEnvVariablesItem := range job.ComputedEnvVariables {
			computedEnvVariablesItemMap, err := ResourceIbmCodeEngineJobEnvVarToMap(&computedEnvVariablesItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "read", "computed_env_variables-to-map").GetDiag()
			}
			computedEnvVariables = append(computedEnvVariables, computedEnvVariablesItemMap)
		}
		if err = d.Set("computed_env_variables", computedEnvVariables); err != nil {
			err = fmt.Errorf("Error setting computed_env_variables: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "read", "set-computed_env_variables").GetDiag()
		}
	}
	if !core.IsNil(job.CreatedAt) {
		if err = d.Set("created_at", job.CreatedAt); err != nil {
			err = fmt.Errorf("Error setting created_at: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "read", "set-created_at").GetDiag()
		}
	}
	if err = d.Set("entity_tag", job.EntityTag); err != nil {
		err = fmt.Errorf("Error setting entity_tag: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "read", "set-entity_tag").GetDiag()
	}
	if !core.IsNil(job.Href) {
		if err = d.Set("href", job.Href); err != nil {
			err = fmt.Errorf("Error setting href: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "read", "set-href").GetDiag()
		}
	}
	if !core.IsNil(job.ID) {
		if err = d.Set("job_id", job.ID); err != nil {
			err = fmt.Errorf("Error setting job_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "read", "set-job_id").GetDiag()
		}
	}
	if !core.IsNil(job.Region) {
		if err = d.Set("region", job.Region); err != nil {
			err = fmt.Errorf("Error setting region: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "read", "set-region").GetDiag()
		}
	}
	if !core.IsNil(job.ResourceType) {
		if err = d.Set("resource_type", job.ResourceType); err != nil {
			err = fmt.Errorf("Error setting resource_type: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "read", "set-resource_type").GetDiag()
		}
	}
	if err = d.Set("etag", response.Headers.Get("Etag")); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting etag: %s", err), "ibm_code_engine_job", "read", "set-etag").GetDiag()
	}

	return nil
}

func resourceIbmCodeEngineJobUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateJobOptions := &codeenginev2.UpdateJobOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "update", "sep-id-parts").GetDiag()
	}

	updateJobOptions.SetProjectID(parts[0])
	updateJobOptions.SetName(parts[1])

	hasChange := false

	patchVals := &codeenginev2.JobPatch{}
	if d.HasChange("project_id") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "project_id")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_code_engine_job", "update", "project_id-forces-new").GetDiag()
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
			runEnvVariablesItem, err := ResourceIbmCodeEngineJobMapToEnvVarPrototype(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "update", "parse-run_env_variables").GetDiag()
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
			runVolumeMountsItem, err := ResourceIbmCodeEngineJobMapToVolumeMountPrototype(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "update", "parse-run_volume_mounts").GetDiag()
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
		// Fields with `nil` values are omitted from the generic map,
		// so we need to re-add them to support removing arguments
		// in merge-patch operations sent to the service.
		updateJobOptions.Job = ResourceIbmCodeEngineJobJobPatchAsPatch(patchVals, d)

		_, _, err = codeEngineClient.UpdateJobWithContext(context, updateJobOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateJobWithContext failed: %s", err.Error()), "ibm_code_engine_job", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIbmCodeEngineJobRead(context, d, meta)
}

func resourceIbmCodeEngineJobDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteJobOptions := &codeenginev2.DeleteJobOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_job", "delete", "sep-id-parts").GetDiag()
	}

	deleteJobOptions.SetProjectID(parts[0])
	deleteJobOptions.SetName(parts[1])

	_, err = codeEngineClient.DeleteJobWithContext(context, deleteJobOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteJobWithContext failed: %s", err.Error()), "ibm_code_engine_job", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIbmCodeEngineJobMapToEnvVarPrototype(modelMap map[string]interface{}) (*codeenginev2.EnvVarPrototype, error) {
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

func ResourceIbmCodeEngineJobMapToVolumeMountPrototype(modelMap map[string]interface{}) (*codeenginev2.VolumeMountPrototype, error) {
	model := &codeenginev2.VolumeMountPrototype{}
	model.MountPath = core.StringPtr(modelMap["mount_path"].(string))
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	model.Reference = core.StringPtr(modelMap["reference"].(string))
	model.Type = core.StringPtr(modelMap["type"].(string))
	return model, nil
}

func ResourceIbmCodeEngineJobEnvVarToMap(model *codeenginev2.EnvVar) (map[string]interface{}, error) {
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

func ResourceIbmCodeEngineJobVolumeMountToMap(model *codeenginev2.VolumeMount) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["mount_path"] = *model.MountPath
	modelMap["name"] = *model.Name
	modelMap["reference"] = *model.Reference
	modelMap["type"] = *model.Type
	return modelMap, nil
}

func ResourceIbmCodeEngineJobJobPatchAsPatch(patchVals *codeenginev2.JobPatch, d *schema.ResourceData) map[string]interface{} {
	patch, _ := patchVals.AsPatch()
	var path string

	path = "image_reference"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["image_reference"] = nil
	}
	path = "image_secret"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["image_secret"] = nil
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
	path = "run_mode"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["run_mode"] = nil
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
	path = "scale_array_spec"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["scale_array_spec"] = nil
	}
	path = "scale_cpu_limit"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["scale_cpu_limit"] = nil
	}
	path = "scale_ephemeral_storage_limit"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["scale_ephemeral_storage_limit"] = nil
	}
	path = "scale_max_execution_time"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["scale_max_execution_time"] = nil
	}
	path = "scale_memory_limit"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["scale_memory_limit"] = nil
	}
	path = "scale_retry_limit"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["scale_retry_limit"] = nil
	}

	return patch
}

func ResourceIbmCodeEngineJobVolumeMountPrototypeAsPatch(patch map[string]interface{}, d *schema.ResourceData) {
	var path string

	path = "run_volume_mounts.0.name"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["name"] = nil
	}
}

func ResourceIbmCodeEngineJobEnvVarPrototypeAsPatch(patch map[string]interface{}, d *schema.ResourceData) {
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
