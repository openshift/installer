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
	"github.com/IBM/code-engine-go-sdk/codeenginev2"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIbmCodeEngineBuild() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmCodeEngineBuildRead,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the project.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of your build.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the resource was created.",
			},
			"entity_tag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the build instance, which is used to achieve optimistic locking.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When you provision a new build,  a URL is created identifying the location of the instance.",
			},
			"build_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The identifier of the resource.",
			},
			"output_image": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the image.",
			},
			"output_secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The secret that is required to access the image registry. Make sure that the secret is granted with push permissions towards the specified container registry namespace.",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The region of the project the resource is located in. Possible values: 'au-syd', 'br-sao', 'ca-tor', 'eu-de', 'eu-gb', 'jp-osa', 'jp-tok', 'us-east', 'us-south'.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the build.",
			},
			"source_context_dir": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Optional directory in the repository that contains the buildpacks file or the Dockerfile.",
			},
			"source_revision": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Commit, tag, or branch in the source repository to pull. This field is optional if the `source_type` is `git` and uses the HEAD of default branch if not specified. If the `source_type` value is `local`, this field must be omitted.",
			},
			"source_secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the secret that is used access the repository source. This field is optional if the `source_type` is `git`. Additionally, if the `source_url` points to a repository that requires authentication, the build will be created but cannot access any source code, until this property is provided, too. If the `source_type` value is `local`, this field must be omitted.",
			},
			"source_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the type of source to determine if your build source is in a repository or based on local source code.* local - For builds from local source code.* git - For builds from git version controlled source code.",
			},
			"source_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the code repository. This field is required if the `source_type` is `git`. If the `source_type` value is `local`, this field must be omitted. If the repository is publicly available you can provide a 'https' URL like `https://github.com/IBM/CodeEngine`. If the repository requires authentication, you need to provide a 'ssh' URL like `git@github.com:IBM/CodeEngine.git` along with a `source_secret` that points to a secret of format `ssh_auth`.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the build.",
			},
			"status_details": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The detailed status of the build.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Optional information to provide more context in case of a 'failed' or 'warning' status.",
						},
					},
				},
			},
			"strategy_size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Optional size for the build, which determines the amount of resources used. Build sizes are `small`, `medium`, `large`, `xlarge`, `xxlarge`.",
			},
			"strategy_spec_file": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Optional path to the specification file that is used for build strategies for building an image.",
			},
			"strategy_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The strategy to use for building the image.",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum amount of time, in seconds, that can pass before the build must succeed or fail.",
			},
		},
	}
}

func dataSourceIbmCodeEngineBuildRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_code_engine_build", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getBuildOptions := &codeenginev2.GetBuildOptions{}

	getBuildOptions.SetProjectID(d.Get("project_id").(string))
	getBuildOptions.SetName(d.Get("name").(string))

	build, _, err := codeEngineClient.GetBuildWithContext(context, getBuildOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetBuildWithContext failed: %s", err.Error()), "(Data) ibm_code_engine_build", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *getBuildOptions.ProjectID, *getBuildOptions.Name))

	if !core.IsNil(build.CreatedAt) {
		if err = d.Set("created_at", build.CreatedAt); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting created_at: %s", err), "(Data) ibm_code_engine_build", "read", "set-created_at").GetDiag()
		}
	}

	if err = d.Set("entity_tag", build.EntityTag); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting entity_tag: %s", err), "(Data) ibm_code_engine_build", "read", "set-entity_tag").GetDiag()
	}

	if !core.IsNil(build.Href) {
		if err = d.Set("href", build.Href); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_code_engine_build", "read", "set-href").GetDiag()
		}
	}

	if !core.IsNil(build.ID) {
		if err = d.Set("build_id", build.ID); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting build_id: %s", err), "(Data) ibm_code_engine_build", "read", "set-build_id").GetDiag()
		}
	}

	if err = d.Set("output_image", build.OutputImage); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting output_image: %s", err), "(Data) ibm_code_engine_build", "read", "set-output_image").GetDiag()
	}

	if err = d.Set("output_secret", build.OutputSecret); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting output_secret: %s", err), "(Data) ibm_code_engine_build", "read", "set-output_secret").GetDiag()
	}

	if !core.IsNil(build.Region) {
		if err = d.Set("region", build.Region); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting region: %s", err), "(Data) ibm_code_engine_build", "read", "set-region").GetDiag()
		}
	}

	if !core.IsNil(build.ResourceType) {
		if err = d.Set("resource_type", build.ResourceType); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting resource_type: %s", err), "(Data) ibm_code_engine_build", "read", "set-resource_type").GetDiag()
		}
	}

	if !core.IsNil(build.SourceContextDir) {
		if err = d.Set("source_context_dir", build.SourceContextDir); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting source_context_dir: %s", err), "(Data) ibm_code_engine_build", "read", "set-source_context_dir").GetDiag()
		}
	}

	if !core.IsNil(build.SourceRevision) {
		if err = d.Set("source_revision", build.SourceRevision); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting source_revision: %s", err), "(Data) ibm_code_engine_build", "read", "set-source_revision").GetDiag()
		}
	}

	if !core.IsNil(build.SourceSecret) {
		if err = d.Set("source_secret", build.SourceSecret); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting source_secret: %s", err), "(Data) ibm_code_engine_build", "read", "set-source_secret").GetDiag()
		}
	}

	if err = d.Set("source_type", build.SourceType); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting source_type: %s", err), "(Data) ibm_code_engine_build", "read", "set-source_type").GetDiag()
	}

	if !core.IsNil(build.SourceURL) {
		if err = d.Set("source_url", build.SourceURL); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting source_url: %s", err), "(Data) ibm_code_engine_build", "read", "set-source_url").GetDiag()
		}
	}

	if !core.IsNil(build.Status) {
		if err = d.Set("status", build.Status); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting status: %s", err), "(Data) ibm_code_engine_build", "read", "set-status").GetDiag()
		}
	}

	if !core.IsNil(build.StatusDetails) {
		statusDetails := []map[string]interface{}{}
		statusDetailsMap, err := DataSourceIbmCodeEngineBuildBuildStatusToMap(build.StatusDetails)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_code_engine_build", "read", "status_details-to-map").GetDiag()
		}
		statusDetails = append(statusDetails, statusDetailsMap)
		if err = d.Set("status_details", statusDetails); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting status_details: %s", err), "(Data) ibm_code_engine_build", "read", "set-status_details").GetDiag()
		}
	}

	if err = d.Set("strategy_size", build.StrategySize); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting strategy_size: %s", err), "(Data) ibm_code_engine_build", "read", "set-strategy_size").GetDiag()
	}

	if !core.IsNil(build.StrategySpecFile) {
		if err = d.Set("strategy_spec_file", build.StrategySpecFile); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting strategy_spec_file: %s", err), "(Data) ibm_code_engine_build", "read", "set-strategy_spec_file").GetDiag()
		}
	}

	if err = d.Set("strategy_type", build.StrategyType); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting strategy_type: %s", err), "(Data) ibm_code_engine_build", "read", "set-strategy_type").GetDiag()
	}

	if !core.IsNil(build.Timeout) {
		if err = d.Set("timeout", flex.IntValue(build.Timeout)); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting timeout: %s", err), "(Data) ibm_code_engine_build", "read", "set-timeout").GetDiag()
		}
	}

	return nil
}

func DataSourceIbmCodeEngineBuildBuildStatusToMap(model *codeenginev2.BuildStatus) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Reason != nil {
		modelMap["reason"] = *model.Reason
	}
	return modelMap, nil
}
