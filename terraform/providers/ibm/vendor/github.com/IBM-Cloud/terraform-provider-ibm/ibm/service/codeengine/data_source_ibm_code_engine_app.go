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

func DataSourceIbmCodeEngineApp() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmCodeEngineAppRead,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the project.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of your application.",
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
			"image_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Optional port the app listens on. While the app will always be exposed via port `443` for end users, this port is used to connect to the port that is exposed by the container image.",
			},
			"image_reference": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the image that is used for this app. The format is `REGISTRY/NAMESPACE/REPOSITORY:TAG` where `REGISTRY` and `TAG` are optional. If `REGISTRY` is not specified, the default is `docker.io`. If `TAG` is not specified, the default is `latest`. If the image reference points to a registry that requires authentication, make sure to also specify the property `image_secret`.",
			},
			"image_secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Optional name of the image registry access secret. The image registry access secret is used to authenticate with a private registry when you download the container image. If the image reference points to a registry that requires authentication, the app will be created but cannot reach the ready status, until this property is provided, too.",
			},
			"managed_domain_mappings": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Optional value controlling which of the system managed domain mappings will be setup for the application. Valid values are 'local_public', 'local_private' and 'local'. Visibility can only be 'local_private' if the project supports application private visibility.",
			},
			"probe_liveness": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Response model for probes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"failure_threshold": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of consecutive, unsuccessful checks for the probe to be considered failed.",
						},
						"initial_delay": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The amount of time in seconds to wait before the first probe check is performed.",
						},
						"interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The amount of time in seconds between probe checks.",
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The path of the HTTP request to the resource. A path is only supported for a probe with a `type` of `http`.",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The port on which to probe the resource.",
						},
						"timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The amount of time in seconds that the probe waits for a response from the application before it times out and fails.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies whether to use HTTP or TCP for the probe checks. The default is TCP.",
						},
					},
				},
			},
			"probe_readiness": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Response model for probes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"failure_threshold": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of consecutive, unsuccessful checks for the probe to be considered failed.",
						},
						"initial_delay": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The amount of time in seconds to wait before the first probe check is performed.",
						},
						"interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The amount of time in seconds between probe checks.",
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The path of the HTTP request to the resource. A path is only supported for a probe with a `type` of `http`.",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The port on which to probe the resource.",
						},
						"timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The amount of time in seconds that the probe waits for a response from the application before it times out and fails.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies whether to use HTTP or TCP for the probe checks. The default is TCP.",
						},
					},
				},
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
			"run_arguments": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Optional arguments for the app that are passed to start the container. If not specified an empty string array will be applied and the arguments specified by the container image, will be used to start the container.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"run_as_user": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Optional user ID (UID) to run the app.",
			},
			"run_commands": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Optional commands for the app that are passed to start the container. If not specified an empty string array will be applied and the command specified by the container image, will be used to start the container.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"run_env_variables": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "References to config maps, secrets or literal values, which are exposed as environment variables in the application.",
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
			"run_service_account": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Optional name of the service account. For built-in service accounts, you can use the shortened names `manager` , `none`, `reader`, and `writer`.",
			},
			"run_volume_mounts": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Mounts of config maps or secrets.",
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
			"scale_concurrency": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Optional maximum number of requests that can be processed concurrently per instance.",
			},
			"scale_concurrency_target": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Optional threshold of concurrent requests per instance at which one or more additional instances are created. Use this value to scale up instances based on concurrent number of requests. This option defaults to the value of the `scale_concurrency` option, if not specified.",
			},
			"scale_cpu_limit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Optional number of CPU set for the instance of the app. For valid values see [Supported memory and CPU combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo).",
			},
			"scale_down_delay": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Optional amount of time in seconds that delays the scale-down behavior for an app instance.",
			},
			"scale_ephemeral_storage_limit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Optional amount of ephemeral storage to set for the instance of the app. The amount specified as ephemeral storage, must not exceed the amount of `scale_memory_limit`. The units for specifying ephemeral storage are Megabyte (M) or Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information see [Units of measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).",
			},
			"scale_initial_instances": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Optional initial number of instances that are created upon app creation or app update.",
			},
			"scale_max_instances": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Optional maximum number of instances for this app. If you set this value to `0`, this property does not set a upper scaling limit. However, the app scaling is still limited by the project quota for instances. See [Limits and quotas for Code Engine](https://cloud.ibm.com/docs/codeengine?topic=codeengine-limits).",
			},
			"scale_memory_limit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Optional amount of memory set for the instance of the app. For valid values see [Supported memory and CPU combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo). The units for specifying memory are Megabyte (M) or Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information see [Units of measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).",
			},
			"scale_min_instances": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Optional minimum number of instances for this app. If you set this value to `0`, the app will scale down to zero, if not hit by any request for some time.",
			},
			"scale_request_timeout": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Optional amount of time in seconds that is allowed for a running app to respond to a request.",
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
		},
	}
}

func dataSourceIbmCodeEngineAppRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_code_engine_app", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getAppOptions := &codeenginev2.GetAppOptions{}

	getAppOptions.SetProjectID(d.Get("project_id").(string))
	getAppOptions.SetName(d.Get("name").(string))

	app, _, err := codeEngineClient.GetAppWithContext(context, getAppOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetAppWithContext failed: %s", err.Error()), "(Data) ibm_code_engine_app", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *getAppOptions.ProjectID, *getAppOptions.Name))

	if err = d.Set("build", app.Build); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting build: %s", err), "(Data) ibm_code_engine_app", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("build_run", app.BuildRun); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting build_run: %s", err), "(Data) ibm_code_engine_app", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("created_at", app.CreatedAt); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_at: %s", err), "(Data) ibm_code_engine_app", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("endpoint", app.Endpoint); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting endpoint: %s", err), "(Data) ibm_code_engine_app", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("endpoint_internal", app.EndpointInternal); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting endpoint_internal: %s", err), "(Data) ibm_code_engine_app", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("entity_tag", app.EntityTag); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting entity_tag: %s", err), "(Data) ibm_code_engine_app", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("href", app.Href); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_code_engine_app", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("app_id", app.ID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting app_id: %s", err), "(Data) ibm_code_engine_app", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("image_port", flex.IntValue(app.ImagePort)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting image_port: %s", err), "(Data) ibm_code_engine_app", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("image_reference", app.ImageReference); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting image_reference: %s", err), "(Data) ibm_code_engine_app", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("image_secret", app.ImageSecret); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting image_secret: %s", err), "(Data) ibm_code_engine_app", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("managed_domain_mappings", app.ManagedDomainMappings); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting managed_domain_mappings: %s", err), "(Data) ibm_code_engine_app", "read")
		return tfErr.GetDiag()
	}

	probeLiveness := []map[string]interface{}{}
	if app.ProbeLiveness != nil {
		modelMap, err := dataSourceIbmCodeEngineAppProbeToMap(app.ProbeLiveness)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_code_engine_app", "read")
			return tfErr.GetDiag()
		}
		probeLiveness = append(probeLiveness, modelMap)
	}
	if err = d.Set("probe_liveness", probeLiveness); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting probe_liveness: %s", err), "(Data) ibm_code_engine_app", "read")
		return tfErr.GetDiag()
	}

	probeReadiness := []map[string]interface{}{}
	if app.ProbeReadiness != nil {
		modelMap, err := dataSourceIbmCodeEngineAppProbeToMap(app.ProbeReadiness)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_code_engine_app", "read")
			return tfErr.GetDiag()
		}
		probeReadiness = append(probeReadiness, modelMap)
	}
	if err = d.Set("probe_readiness", probeReadiness); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting probe_readiness: %s", err), "(Data) ibm_code_engine_app", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("region", app.Region); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting region: %s", err), "(Data) ibm_code_engine_app", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("resource_type", app.ResourceType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting resource_type: %s", err), "(Data) ibm_code_engine_app", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("run_as_user", flex.IntValue(app.RunAsUser)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting run_as_user: %s", err), "(Data) ibm_code_engine_app", "read")
		return tfErr.GetDiag()
	}

	runEnvVariables := []map[string]interface{}{}
	if app.RunEnvVariables != nil {
		for _, modelItem := range app.RunEnvVariables {
			modelMap, err := dataSourceIbmCodeEngineAppEnvVarToMap(&modelItem) /* #nosec G601 */
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_code_engine_app", "read")
				return tfErr.GetDiag()
			}
			runEnvVariables = append(runEnvVariables, modelMap)
		}
	}
	if err = d.Set("run_env_variables", runEnvVariables); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting run_env_variables: %s", err), "(Data) ibm_code_engine_app", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("run_service_account", app.RunServiceAccount); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting run_service_account: %s", err), "(Data) ibm_code_engine_app", "read")
		return tfErr.GetDiag()
	}

	runVolumeMounts := []map[string]interface{}{}
	if app.RunVolumeMounts != nil {
		for _, modelItem := range app.RunVolumeMounts {
			modelMap, err := dataSourceIbmCodeEngineAppVolumeMountToMap(&modelItem) /* #nosec G601 */
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_code_engine_app", "read")
				return tfErr.GetDiag()
			}
			runVolumeMounts = append(runVolumeMounts, modelMap)
		}
	}
	if err = d.Set("run_volume_mounts", runVolumeMounts); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting run_volume_mounts: %s", err), "(Data) ibm_code_engine_app", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("scale_concurrency", flex.IntValue(app.ScaleConcurrency)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting scale_concurrency: %s", err), "(Data) ibm_code_engine_app", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("scale_concurrency_target", flex.IntValue(app.ScaleConcurrencyTarget)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting scale_concurrency_target: %s", err), "(Data) ibm_code_engine_app", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("scale_cpu_limit", app.ScaleCpuLimit); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting scale_cpu_limit: %s", err), "(Data) ibm_code_engine_app", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("scale_down_delay", flex.IntValue(app.ScaleDownDelay)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting scale_down_delay: %s", err), "(Data) ibm_code_engine_app", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("scale_ephemeral_storage_limit", app.ScaleEphemeralStorageLimit); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting scale_ephemeral_storage_limit: %s", err), "(Data) ibm_code_engine_app", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("scale_initial_instances", flex.IntValue(app.ScaleInitialInstances)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting scale_initial_instances: %s", err), "(Data) ibm_code_engine_app", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("scale_max_instances", flex.IntValue(app.ScaleMaxInstances)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting scale_max_instances: %s", err), "(Data) ibm_code_engine_app", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("scale_memory_limit", app.ScaleMemoryLimit); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting scale_memory_limit: %s", err), "(Data) ibm_code_engine_app", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("scale_min_instances", flex.IntValue(app.ScaleMinInstances)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting scale_min_instances: %s", err), "(Data) ibm_code_engine_app", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("scale_request_timeout", flex.IntValue(app.ScaleRequestTimeout)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting scale_request_timeout: %s", err), "(Data) ibm_code_engine_app", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("status", app.Status); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting status: %s", err), "(Data) ibm_code_engine_app", "read")
		return tfErr.GetDiag()
	}

	statusDetails := []map[string]interface{}{}
	if app.StatusDetails != nil {
		modelMap, err := dataSourceIbmCodeEngineAppAppStatusToMap(app.StatusDetails)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_code_engine_app", "read")
			return tfErr.GetDiag()
		}
		statusDetails = append(statusDetails, modelMap)
	}
	if err = d.Set("status_details", statusDetails); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting status_details: %s", err), "(Data) ibm_code_engine_app", "read")
		return tfErr.GetDiag()
	}

	return nil
}

func dataSourceIbmCodeEngineAppProbeToMap(model *codeenginev2.Probe) (map[string]interface{}, error) {
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
		modelMap["path"] = model.Path
	}
	if model.Port != nil {
		modelMap["port"] = flex.IntValue(model.Port)
	}
	if model.Timeout != nil {
		modelMap["timeout"] = flex.IntValue(model.Timeout)
	}
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	return modelMap, nil
}

func dataSourceIbmCodeEngineAppEnvVarToMap(model *codeenginev2.EnvVar) (map[string]interface{}, error) {
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

func dataSourceIbmCodeEngineAppVolumeMountToMap(model *codeenginev2.VolumeMount) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["mount_path"] = model.MountPath
	modelMap["name"] = model.Name
	modelMap["reference"] = model.Reference
	modelMap["type"] = model.Type
	return modelMap, nil
}

func dataSourceIbmCodeEngineAppAppStatusToMap(model *codeenginev2.AppStatus) (map[string]interface{}, error) {
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
