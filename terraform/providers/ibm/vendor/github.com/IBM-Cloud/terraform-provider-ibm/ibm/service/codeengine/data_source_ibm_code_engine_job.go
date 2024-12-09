// Copyright IBM Corp. 2024 All Rights Reserved.
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
	"github.com/IBM/code-engine-go-sdk/codeenginev2"
)

func DataSourceIbmCodeEngineJob() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmCodeEngineJobRead,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the project.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of your job.",
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
			"image_reference": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the image that is used for this job. The format is `REGISTRY/NAMESPACE/REPOSITORY:TAG` where `REGISTRY` and `TAG` are optional. If `REGISTRY` is not specified, the default is `docker.io`. If `TAG` is not specified, the default is `latest`. If the image reference points to a registry that requires authentication, make sure to also specify the property `image_secret`.",
			},
			"image_secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the image registry access secret. The image registry access secret is used to authenticate with a private registry when you download the container image. If the image reference points to a registry that requires authentication, the job / job runs will be created but submitted job runs will fail, until this property is provided, too. This property must not be set on a job run, which references a job template.",
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
			"run_arguments": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Set arguments for the job that are passed to start job run containers. If not specified an empty string array will be applied and the arguments specified by the container image, will be used to start the container.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"run_as_user": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The user ID (UID) to run the job.",
			},
			"run_commands": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Set commands for the job that are passed to start job run containers. If not specified an empty string array will be applied and the command specified by the container image, will be used to start the container.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"run_env_variables": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "References to config maps, secrets or literal values, which are exposed as environment variables in the job run.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The key to reference as environment variable.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the environment variable.",
						},
						"prefix": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A prefix that can be added to all keys of a full secret or config map reference.",
						},
						"reference": {
							Type:        schema.TypeString,
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
							Computed:    true,
							Description: "The literal value of the environment variable.",
						},
					},
				},
			},
			"run_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The mode for runs of the job. Valid values are `task` and `daemon`. In `task` mode, the `max_execution_time` and `retry_limit` properties apply. In `daemon` mode, since there is no timeout and failed instances are restarted indefinitely, the `max_execution_time` and `retry_limit` properties are not allowed.",
			},
			"run_service_account": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the service account. For built-in service accounts, you can use the shortened names `manager`, `none`, `reader`, and `writer`. This property must not be set on a job run, which references a job template.",
			},
			"run_volume_mounts": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Optional mounts of config maps or secrets.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mount_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The path that should be mounted.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the mount.",
						},
						"reference": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the referenced secret or config map.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specify the type of the volume mount. Allowed types are: 'config_map', 'secret'.",
						},
					},
				},
			},
			"scale_array_spec": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Define a custom set of array indices as a comma-separated list containing single values and hyphen-separated ranges, such as  5,12-14,23,27. Each instance gets its array index value from the environment variable JOB_INDEX. The number of unique array indices that you specify with this parameter determines the number of job instances to run.",
			},
			"scale_cpu_limit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Optional amount of CPU set for the instance of the job. For valid values see [Supported memory and CPU combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo).",
			},
			"scale_ephemeral_storage_limit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Optional amount of ephemeral storage to set for the instance of the job. The amount specified as ephemeral storage, must not exceed the amount of `scale_memory_limit`. The units for specifying ephemeral storage are Megabyte (M) or Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information see [Units of measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).",
			},
			"scale_max_execution_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum execution time in seconds for runs of the job. This property can only be specified if `run_mode` is `task`.",
			},
			"scale_memory_limit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Optional amount of memory set for the instance of the job. For valid values see [Supported memory and CPU combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo). The units for specifying memory are Megabyte (M) or Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information see [Units of measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).",
			},
			"scale_retry_limit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of times to rerun an instance of the job before the job is marked as failed. This property can only be specified if `run_mode` is `task`.",
			},
		},
	}
}

func dataSourceIbmCodeEngineJobRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_code_engine_job", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getJobOptions := &codeenginev2.GetJobOptions{}

	getJobOptions.SetProjectID(d.Get("project_id").(string))
	getJobOptions.SetName(d.Get("name").(string))

	job, _, err := codeEngineClient.GetJobWithContext(context, getJobOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetJobWithContext failed: %s", err.Error()), "(Data) ibm_code_engine_job", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *getJobOptions.ProjectID, *getJobOptions.Name))

	if err = d.Set("build", job.Build); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting build: %s", err), "(Data) ibm_code_engine_job", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("build_run", job.BuildRun); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting build_run: %s", err), "(Data) ibm_code_engine_job", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("created_at", job.CreatedAt); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_at: %s", err), "(Data) ibm_code_engine_job", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("entity_tag", job.EntityTag); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting entity_tag: %s", err), "(Data) ibm_code_engine_job", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("href", job.Href); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_code_engine_job", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("job_id", job.ID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting job_id: %s", err), "(Data) ibm_code_engine_job", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("image_reference", job.ImageReference); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting image_reference: %s", err), "(Data) ibm_code_engine_job", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("image_secret", job.ImageSecret); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting image_secret: %s", err), "(Data) ibm_code_engine_job", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("region", job.Region); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting region: %s", err), "(Data) ibm_code_engine_job", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("resource_type", job.ResourceType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting resource_type: %s", err), "(Data) ibm_code_engine_job", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("run_as_user", flex.IntValue(job.RunAsUser)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting run_as_user: %s", err), "(Data) ibm_code_engine_job", "read")
		return tfErr.GetDiag()
	}

	runEnvVariables := []map[string]interface{}{}
	if job.RunEnvVariables != nil {
		for _, modelItem := range job.RunEnvVariables {
			modelMap, err := dataSourceIbmCodeEngineJobEnvVarToMap(&modelItem) /* #nosec G601 */
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_code_engine_job", "read")
				return tfErr.GetDiag()
			}
			runEnvVariables = append(runEnvVariables, modelMap)
		}
	}
	if err = d.Set("run_env_variables", runEnvVariables); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting run_env_variables: %s", err), "(Data) ibm_code_engine_job", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("run_mode", job.RunMode); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting run_mode: %s", err), "(Data) ibm_code_engine_job", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("run_service_account", job.RunServiceAccount); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting run_service_account: %s", err), "(Data) ibm_code_engine_job", "read")
		return tfErr.GetDiag()
	}

	runVolumeMounts := []map[string]interface{}{}
	if job.RunVolumeMounts != nil {
		for _, modelItem := range job.RunVolumeMounts {
			modelMap, err := dataSourceIbmCodeEngineJobVolumeMountToMap(&modelItem) /* #nosec G601 */
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_code_engine_job", "read")
				return tfErr.GetDiag()
			}
			runVolumeMounts = append(runVolumeMounts, modelMap)
		}
	}
	if err = d.Set("run_volume_mounts", runVolumeMounts); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting run_volume_mounts: %s", err), "(Data) ibm_code_engine_job", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("scale_array_spec", job.ScaleArraySpec); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting scale_array_spec: %s", err), "(Data) ibm_code_engine_job", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("scale_cpu_limit", job.ScaleCpuLimit); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting scale_cpu_limit: %s", err), "(Data) ibm_code_engine_job", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("scale_ephemeral_storage_limit", job.ScaleEphemeralStorageLimit); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting scale_ephemeral_storage_limit: %s", err), "(Data) ibm_code_engine_job", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("scale_max_execution_time", flex.IntValue(job.ScaleMaxExecutionTime)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting scale_max_execution_time: %s", err), "(Data) ibm_code_engine_job", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("scale_memory_limit", job.ScaleMemoryLimit); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting scale_memory_limit: %s", err), "(Data) ibm_code_engine_job", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("scale_retry_limit", flex.IntValue(job.ScaleRetryLimit)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting scale_retry_limit: %s", err), "(Data) ibm_code_engine_job", "read")
		return tfErr.GetDiag()
	}

	return nil
}

func dataSourceIbmCodeEngineJobEnvVarToMap(model *codeenginev2.EnvVar) (map[string]interface{}, error) {
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
	modelMap["type"] = model.Type
	if model.Value != nil {
		modelMap["value"] = model.Value
	}
	return modelMap, nil
}

func dataSourceIbmCodeEngineJobVolumeMountToMap(model *codeenginev2.VolumeMount) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["mount_path"] = model.MountPath
	modelMap["name"] = model.Name
	modelMap["reference"] = model.Reference
	modelMap["type"] = model.Type
	return modelMap, nil
}
