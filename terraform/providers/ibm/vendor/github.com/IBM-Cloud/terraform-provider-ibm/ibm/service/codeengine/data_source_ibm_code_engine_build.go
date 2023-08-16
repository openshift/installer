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
	"github.com/IBM/code-engine-go-sdk/codeenginev2"
)

func DataSourceIbmCodeEngineBuild() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmCodeEngineBuildRead,

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the project.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of your build.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the resource was created.",
			},
			"entity_tag": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the build instance, which is used to achieve optimistic locking.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When you provision a new build,  a URL is created identifying the location of the instance.",
			},
			"build_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The identifier of the resource.",
			},
			"output_image": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the image.",
			},
			"output_secret": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The secret that is required to access the image registry. Make sure that the secret is granted with push permissions towards the specified container registry namespace.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the build.",
			},
			"source_context_dir": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Option directory in the repository that contains the buildpacks file or the Dockerfile.",
			},
			"source_revision": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Commit, tag, or branch in the source repository to pull. This field is optional if the `source_type` is `git` and uses the HEAD of default branch if not specified. If the `source_type` value is `local`, this field must be omitted.",
			},
			"source_secret": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the secret that is used access the repository source. This field is optional if the `source_type` is `git`. Additionally, if the `source_url` points to a repository that requires authentication, the build will be created but cannot access any source code, until this property is provided, too. If the `source_type` value is `local`, this field must be omitted.",
			},
			"source_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the type of source to determine if your build source is in a repository or based on local source code.* local - For builds from local source code.* git - For builds from git version controlled source code.",
			},
			"source_url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the code repository. This field is required if the `source_type` is `git`. If the `source_type` value is `local`, this field must be omitted. If the repository is publicly available you can provide a 'https' URL like `https://github.com/IBM/CodeEngine`. If the repository requires authentication, you need to provide a 'ssh' URL like `git@github.com:IBM/CodeEngine.git` along with a `source_secret` that points to a secret of format `ssh_auth`.",
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the build.",
			},
			"status_details": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The detailed status of the build.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reason": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Optional information to provide more context in case of a 'failed' or 'warning' status.",
						},
					},
				},
			},
			"strategy_size": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Optional size for the build, which determines the amount of resources used. Build sizes are `small`, `medium`, `large`, `xlarge`.",
			},
			"strategy_spec_file": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Optional path to the specification file that is used for build strategies for building an image.",
			},
			"strategy_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The strategy to use for building the image.",
			},
			"timeout": &schema.Schema{
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
		return diag.FromErr(err)
	}

	getBuildOptions := &codeenginev2.GetBuildOptions{}

	getBuildOptions.SetProjectID(d.Get("project_id").(string))
	getBuildOptions.SetName(d.Get("name").(string))

	build, response, err := codeEngineClient.GetBuildWithContext(context, getBuildOptions)
	if err != nil {
		log.Printf("[DEBUG] GetBuildWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetBuildWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *getBuildOptions.ProjectID, *getBuildOptions.Name))

	if err = d.Set("created_at", build.CreatedAt); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}

	if err = d.Set("entity_tag", build.EntityTag); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting entity_tag: %s", err))
	}

	if err = d.Set("href", build.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}

	if err = d.Set("build_id", build.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting build_id: %s", err))
	}

	if err = d.Set("output_image", build.OutputImage); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting output_image: %s", err))
	}

	if err = d.Set("output_secret", build.OutputSecret); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting output_secret: %s", err))
	}

	if err = d.Set("resource_type", build.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}

	if err = d.Set("source_context_dir", build.SourceContextDir); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting source_context_dir: %s", err))
	}

	if err = d.Set("source_revision", build.SourceRevision); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting source_revision: %s", err))
	}

	if err = d.Set("source_secret", build.SourceSecret); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting source_secret: %s", err))
	}

	if err = d.Set("source_type", build.SourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting source_type: %s", err))
	}

	if err = d.Set("source_url", build.SourceURL); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting source_url: %s", err))
	}

	if err = d.Set("status", build.Status); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting status: %s", err))
	}

	statusDetails := []map[string]interface{}{}
	if build.StatusDetails != nil {
		modelMap, err := dataSourceIbmCodeEngineBuildBuildStatusToMap(build.StatusDetails)
		if err != nil {
			return diag.FromErr(err)
		}
		statusDetails = append(statusDetails, modelMap)
	}
	if err = d.Set("status_details", statusDetails); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting status_details %s", err))
	}

	if err = d.Set("strategy_size", build.StrategySize); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting strategy_size: %s", err))
	}

	if err = d.Set("strategy_spec_file", build.StrategySpecFile); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting strategy_spec_file: %s", err))
	}

	if err = d.Set("strategy_type", build.StrategyType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting strategy_type: %s", err))
	}

	if err = d.Set("timeout", flex.IntValue(build.Timeout)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting timeout: %s", err))
	}

	return nil
}

func dataSourceIbmCodeEngineBuildBuildStatusToMap(model *codeenginev2.BuildStatus) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Reason != nil {
		modelMap["reason"] = model.Reason
	}
	return modelMap, nil
}
