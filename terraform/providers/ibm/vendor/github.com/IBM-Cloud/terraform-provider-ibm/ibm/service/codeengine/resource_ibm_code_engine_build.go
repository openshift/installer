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

func ResourceIbmCodeEngineBuild() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmCodeEngineBuildCreate,
		ReadContext:   resourceIbmCodeEngineBuildRead,
		UpdateContext: resourceIbmCodeEngineBuildUpdate,
		DeleteContext: resourceIbmCodeEngineBuildDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_build", "project_id"),
				Description:  "The ID of the project.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_build", "name"),
				Description:  "The name of the build. Use a name that is unique within the project.",
			},
			"output_image": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_build", "output_image"),
				Description:  "The name of the image.",
			},
			"output_secret": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_build", "output_secret"),
				Description:  "The secret that is required to access the image registry. Make sure that the secret is granted with push permissions towards the specified container registry namespace.",
			},
			"strategy_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_build", "strategy_type"),
				Description:  "The strategy to use for building the image.",
			},
			"source_context_dir": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_build", "source_context_dir"),
				Description:  "Option directory in the repository that contains the buildpacks file or the Dockerfile.",
			},
			"source_revision": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_build", "source_revision"),
				Description:  "Commit, tag, or branch in the source repository to pull. This field is optional if the `source_type` is `git` and uses the HEAD of default branch if not specified. If the `source_type` value is `local`, this field must be omitted.",
			},
			"source_secret": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_build", "source_secret"),
				Description:  "Name of the secret that is used access the repository source. This field is optional if the `source_type` is `git`. Additionally, if the `source_url` points to a repository that requires authentication, the build will be created but cannot access any source code, until this property is provided, too. If the `source_type` value is `local`, this field must be omitted.",
			},
			"source_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "git",
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_build", "source_type"),
				Description:  "Specifies the type of source to determine if your build source is in a repository or based on local source code.* local - For builds from local source code.* git - For builds from git version controlled source code.",
			},
			"source_url": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_build", "source_url"),
				Description:  "The URL of the code repository. This field is required if the `source_type` is `git`. If the `source_type` value is `local`, this field must be omitted. If the repository is publicly available you can provide a 'https' URL like `https://github.com/IBM/CodeEngine`. If the repository requires authentication, you need to provide a 'ssh' URL like `git@github.com:IBM/CodeEngine.git` along with a `source_secret` that points to a secret of format `ssh_auth`.",
			},
			"strategy_size": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "medium",
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_build", "strategy_size"),
				Description:  "Optional size for the build, which determines the amount of resources used. Build sizes are `small`, `medium`, `large`, `xlarge`.",
			},
			"strategy_spec_file": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Dockerfile",
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_build", "strategy_spec_file"),
				Description:  "Optional path to the specification file that is used for build strategies for building an image.",
			},
			"timeout": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      600,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_build", "timeout"),
				Description:  "The maximum amount of time, in seconds, that can pass before the build must succeed or fail.",
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
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the build.",
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
			"etag": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceIbmCodeEngineBuildValidator() *validate.ResourceValidator {
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
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-z0-9]([\-a-z0-9]*[a-z0-9])?$`,
			MinValueLength:             1,
			MaxValueLength:             63,
		},
		validate.ValidateSchema{
			Identifier:                 "output_image",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^([a-z0-9][a-z0-9\-_.]+[a-z0-9][\/])?([a-z0-9][a-z0-9\-_]+[a-z0-9][\/])?[a-z0-9][a-z0-9\-_.\/]+[a-z0-9](:[\w][\w.\-]{0,127})?(@sha256:[a-fA-F0-9]{64})?$`,
			MinValueLength:             1,
			MaxValueLength:             256,
		},
		validate.ValidateSchema{
			Identifier:                 "output_secret",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-z0-9]([\-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([\-a-z0-9]*[a-z0-9])?)*$`,
			MinValueLength:             1,
			MaxValueLength:             253,
		},
		validate.ValidateSchema{
			Identifier:                 "strategy_type",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `[\S]*`,
			MinValueLength:             1,
			MaxValueLength:             253,
		},
		validate.ValidateSchema{
			Identifier:                 "source_context_dir",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^(.*)+$`,
			MinValueLength:             0,
			MaxValueLength:             253,
		},
		validate.ValidateSchema{
			Identifier:                 "source_revision",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[\S]*$`,
			MinValueLength:             0,
			MaxValueLength:             253,
		},
		validate.ValidateSchema{
			Identifier:                 "source_secret",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[a-z0-9]([\-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([\-a-z0-9]*[a-z0-9])?)*$`,
			MinValueLength:             1,
			MaxValueLength:             253,
		},
		validate.ValidateSchema{
			Identifier:                 "source_type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "git, local",
		},
		validate.ValidateSchema{
			Identifier:                 "source_url",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^((https:\/\/[a-z0-9]([\-.]?[a-z0-9])+(:\d{1,5})?)|((ssh:\/\/)?git@[a-z0-9]([\-.]{0,1}[a-z0-9])+(:[a-zA-Z0-9\/][\w\-.]*)?))(\/([\w\-.]|%20)+)*$`,
			MinValueLength:             1,
			MaxValueLength:             253,
		},
		validate.ValidateSchema{
			Identifier:                 "strategy_size",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `[\S]*`,
			MinValueLength:             1,
			MaxValueLength:             253,
		},
		validate.ValidateSchema{
			Identifier:                 "strategy_spec_file",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[\S]*$`,
			MinValueLength:             1,
			MaxValueLength:             253,
		},
		validate.ValidateSchema{
			Identifier:                 "timeout",
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			Optional:                   true,
			MinValue:                   "1",
			MaxValue:                   "3600",
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_code_engine_build", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmCodeEngineBuildCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	createBuildOptions := &codeenginev2.CreateBuildOptions{}

	createBuildOptions.SetProjectID(d.Get("project_id").(string))
	createBuildOptions.SetName(d.Get("name").(string))
	createBuildOptions.SetOutputImage(d.Get("output_image").(string))
	createBuildOptions.SetOutputSecret(d.Get("output_secret").(string))
	createBuildOptions.SetStrategyType(d.Get("strategy_type").(string))
	if _, ok := d.GetOk("source_context_dir"); ok {
		createBuildOptions.SetSourceContextDir(d.Get("source_context_dir").(string))
	}
	if _, ok := d.GetOk("source_revision"); ok {
		createBuildOptions.SetSourceRevision(d.Get("source_revision").(string))
	}
	if _, ok := d.GetOk("source_secret"); ok {
		createBuildOptions.SetSourceSecret(d.Get("source_secret").(string))
	}
	if _, ok := d.GetOk("source_type"); ok {
		createBuildOptions.SetSourceType(d.Get("source_type").(string))
	}
	if _, ok := d.GetOk("source_url"); ok {
		createBuildOptions.SetSourceURL(d.Get("source_url").(string))
	}
	if _, ok := d.GetOk("strategy_size"); ok {
		createBuildOptions.SetStrategySize(d.Get("strategy_size").(string))
	}
	if _, ok := d.GetOk("strategy_spec_file"); ok {
		createBuildOptions.SetStrategySpecFile(d.Get("strategy_spec_file").(string))
	}
	if _, ok := d.GetOk("timeout"); ok {
		createBuildOptions.SetTimeout(int64(d.Get("timeout").(int)))
	}

	build, response, err := codeEngineClient.CreateBuildWithContext(context, createBuildOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateBuildWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateBuildWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *createBuildOptions.ProjectID, *build.Name))

	return resourceIbmCodeEngineBuildRead(context, d, meta)
}

func resourceIbmCodeEngineBuildRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	getBuildOptions := &codeenginev2.GetBuildOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getBuildOptions.SetProjectID(parts[0])
	getBuildOptions.SetName(parts[1])

	build, response, err := codeEngineClient.GetBuildWithContext(context, getBuildOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetBuildWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetBuildWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("project_id", build.ProjectID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting project_id: %s", err))
	}
	if err = d.Set("name", build.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("output_image", build.OutputImage); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting output_image: %s", err))
	}
	if err = d.Set("output_secret", build.OutputSecret); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting output_secret: %s", err))
	}
	if err = d.Set("strategy_type", build.StrategyType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting strategy_type: %s", err))
	}
	if !core.IsNil(build.SourceContextDir) {
		if err = d.Set("source_context_dir", build.SourceContextDir); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting source_context_dir: %s", err))
		}
	}
	if !core.IsNil(build.SourceRevision) {
		if err = d.Set("source_revision", build.SourceRevision); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting source_revision: %s", err))
		}
	}
	if !core.IsNil(build.SourceSecret) {
		if err = d.Set("source_secret", build.SourceSecret); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting source_secret: %s", err))
		}
	}
	if !core.IsNil(build.SourceType) {
		if err = d.Set("source_type", build.SourceType); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting source_type: %s", err))
		}
	}
	if !core.IsNil(build.SourceURL) {
		if err = d.Set("source_url", build.SourceURL); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting source_url: %s", err))
		}
	}
	if !core.IsNil(build.StrategySize) {
		if err = d.Set("strategy_size", build.StrategySize); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting strategy_size: %s", err))
		}
	}
	if !core.IsNil(build.StrategySpecFile) {
		if err = d.Set("strategy_spec_file", build.StrategySpecFile); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting strategy_spec_file: %s", err))
		}
	}
	if !core.IsNil(build.Timeout) {
		if err = d.Set("timeout", flex.IntValue(build.Timeout)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting timeout: %s", err))
		}
	}
	if !core.IsNil(build.CreatedAt) {
		if err = d.Set("created_at", build.CreatedAt); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
		}
	}
	if err = d.Set("entity_tag", build.EntityTag); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting entity_tag: %s", err))
	}
	if !core.IsNil(build.Href) {
		if err = d.Set("href", build.Href); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
		}
	}
	if !core.IsNil(build.ID) {
		if err = d.Set("build_id", build.ID); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting build_id: %s", err))
		}
	}
	if !core.IsNil(build.ResourceType) {
		if err = d.Set("resource_type", build.ResourceType); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
		}
	}
	if !core.IsNil(build.Status) {
		if err = d.Set("status", build.Status); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting status: %s", err))
		}
	}
	if !core.IsNil(build.StatusDetails) {
		statusDetailsMap, err := resourceIbmCodeEngineBuildBuildStatusToMap(build.StatusDetails)
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

func resourceIbmCodeEngineBuildUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	updateBuildOptions := &codeenginev2.UpdateBuildOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	updateBuildOptions.SetProjectID(parts[0])
	updateBuildOptions.SetName(parts[1])

	hasChange := false

	patchVals := &codeenginev2.BuildPatch{}
	if d.HasChange("project_id") {
		return diag.FromErr(fmt.Errorf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "project_id"))
	}
	if d.HasChange("name") {
		return diag.FromErr(fmt.Errorf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "name"))
	}
	if d.HasChange("output_image") || d.HasChange("output_secret") || d.HasChange("strategy_type") {
		newOutputImage := d.Get("output_image").(string)
		patchVals.OutputImage = &newOutputImage
		newOutputSecret := d.Get("output_secret").(string)
		patchVals.OutputSecret = &newOutputSecret
		newStrategyType := d.Get("strategy_type").(string)
		patchVals.StrategyType = &newStrategyType
		hasChange = true
	}
	if d.HasChange("source_context_dir") {
		newSourceContextDir := d.Get("source_context_dir").(string)
		patchVals.SourceContextDir = &newSourceContextDir
		hasChange = true
	}
	if d.HasChange("source_revision") {
		newSourceRevision := d.Get("source_revision").(string)
		patchVals.SourceRevision = &newSourceRevision
		hasChange = true
	}
	if d.HasChange("source_secret") {
		newSourceSecret := d.Get("source_secret").(string)
		patchVals.SourceSecret = &newSourceSecret
		hasChange = true
	}
	if d.HasChange("source_type") {
		newSourceType := d.Get("source_type").(string)
		patchVals.SourceType = &newSourceType
		hasChange = true
	}
	if d.HasChange("source_url") {
		newSourceURL := d.Get("source_url").(string)
		patchVals.SourceURL = &newSourceURL
		hasChange = true
	}
	if d.HasChange("strategy_size") {
		newStrategySize := d.Get("strategy_size").(string)
		patchVals.StrategySize = &newStrategySize
		hasChange = true
	}
	if d.HasChange("strategy_spec_file") {
		newStrategySpecFile := d.Get("strategy_spec_file").(string)
		patchVals.StrategySpecFile = &newStrategySpecFile
		hasChange = true
	}
	if d.HasChange("timeout") {
		newTimeout := int64(d.Get("timeout").(int))
		patchVals.Timeout = &newTimeout
		hasChange = true
	}
	updateBuildOptions.SetIfMatch(d.Get("etag").(string))

	if hasChange {
		updateBuildOptions.Build, _ = patchVals.AsPatch()
		_, response, err := codeEngineClient.UpdateBuildWithContext(context, updateBuildOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateBuildWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateBuildWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIbmCodeEngineBuildRead(context, d, meta)
}

func resourceIbmCodeEngineBuildDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteBuildOptions := &codeenginev2.DeleteBuildOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteBuildOptions.SetProjectID(parts[0])
	deleteBuildOptions.SetName(parts[1])

	response, err := codeEngineClient.DeleteBuildWithContext(context, deleteBuildOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteBuildWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteBuildWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIbmCodeEngineBuildBuildStatusToMap(model *codeenginev2.BuildStatus) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Reason != nil {
		modelMap["reason"] = model.Reason
	}
	return modelMap, nil
}
