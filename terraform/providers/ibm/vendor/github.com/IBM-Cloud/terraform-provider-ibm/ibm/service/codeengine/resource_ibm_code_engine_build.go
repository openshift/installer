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

func ResourceIbmCodeEngineBuild() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmCodeEngineBuildCreate,
		ReadContext:   resourceIbmCodeEngineBuildRead,
		UpdateContext: resourceIbmCodeEngineBuildUpdate,
		DeleteContext: resourceIbmCodeEngineBuildDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_build", "project_id"),
				Description:  "The ID of the project.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_build", "name"),
				Description:  "The name of the build.",
			},
			"output_image": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_build", "output_image"),
				Description:  "The name of the image.",
			},
			"output_secret": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_build", "output_secret"),
				Description:  "The secret that is required to access the image registry. Make sure that the secret is granted with push permissions towards the specified container registry namespace.",
			},
			"source_context_dir": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_build", "source_context_dir"),
				Description:  "Optional directory in the repository that contains the buildpacks file or the Dockerfile.",
			},
			"source_revision": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_build", "source_revision"),
				Description:  "Commit, tag, or branch in the source repository to pull. This field is optional if the `source_type` is `git` and uses the HEAD of default branch if not specified. If the `source_type` value is `local`, this field must be omitted.",
			},
			"source_secret": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_build", "source_secret"),
				Description:  "Name of the secret that is used access the repository source. This field is optional if the `source_type` is `git`. Additionally, if the `source_url` points to a repository that requires authentication, the build will be created but cannot access any source code, until this property is provided, too. If the `source_type` value is `local`, this field must be omitted.",
			},
			"source_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "git",
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_build", "source_type"),
				Description:  "Specifies the type of source to determine if your build source is in a repository or based on local source code.* local - For builds from local source code.* git - For builds from git version controlled source code.",
			},
			"source_url": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_build", "source_url"),
				Description:  "The URL of the code repository. This field is required if the `source_type` is `git`. If the `source_type` value is `local`, this field must be omitted. If the repository is publicly available you can provide a 'https' URL like `https://github.com/IBM/CodeEngine`. If the repository requires authentication, you need to provide a 'ssh' URL like `git@github.com:IBM/CodeEngine.git` along with a `source_secret` that points to a secret of format `ssh_auth`.",
			},
			"strategy_size": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "medium",
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_build", "strategy_size"),
				Description:  "Optional size for the build, which determines the amount of resources used. Build sizes are `small`, `medium`, `large`, `xlarge`, `xxlarge`.",
			},
			"strategy_spec_file": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_build", "strategy_spec_file"),
				Description:  "Optional path to the specification file that is used for build strategies for building an image.",
			},
			"strategy_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_build", "strategy_type"),
				Description:  "The strategy to use for building the image.",
			},
			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      600,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_build", "timeout"),
				Description:  "The maximum amount of time, in seconds, that can pass before the build must succeed or fail.",
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
			"etag": {
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
			Identifier:                 "strategy_type",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `[\S]*`,
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
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_build", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
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

	build, _, err := codeEngineClient.CreateBuildWithContext(context, createBuildOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateBuildWithContext failed: %s", err.Error()), "ibm_code_engine_build", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *createBuildOptions.ProjectID, *build.Name))

	return resourceIbmCodeEngineBuildRead(context, d, meta)
}

func resourceIbmCodeEngineBuildRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_build", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getBuildOptions := &codeenginev2.GetBuildOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_build", "read", "sep-id-parts").GetDiag()
	}

	getBuildOptions.SetProjectID(parts[0])
	getBuildOptions.SetName(parts[1])

	build, response, err := codeEngineClient.GetBuildWithContext(context, getBuildOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetBuildWithContext failed: %s", err.Error()), "ibm_code_engine_build", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("project_id", build.ProjectID); err != nil {
		err = fmt.Errorf("Error setting project_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_build", "read", "set-project_id").GetDiag()
	}
	if err = d.Set("name", build.Name); err != nil {
		err = fmt.Errorf("Error setting name: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_build", "read", "set-name").GetDiag()
	}
	if err = d.Set("output_image", build.OutputImage); err != nil {
		err = fmt.Errorf("Error setting output_image: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_build", "read", "set-output_image").GetDiag()
	}
	if err = d.Set("output_secret", build.OutputSecret); err != nil {
		err = fmt.Errorf("Error setting output_secret: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_build", "read", "set-output_secret").GetDiag()
	}
	if !core.IsNil(build.SourceContextDir) {
		if err = d.Set("source_context_dir", build.SourceContextDir); err != nil {
			err = fmt.Errorf("Error setting source_context_dir: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_build", "read", "set-source_context_dir").GetDiag()
		}
	}
	if !core.IsNil(build.SourceRevision) {
		if err = d.Set("source_revision", build.SourceRevision); err != nil {
			err = fmt.Errorf("Error setting source_revision: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_build", "read", "set-source_revision").GetDiag()
		}
	}
	if !core.IsNil(build.SourceSecret) {
		if err = d.Set("source_secret", build.SourceSecret); err != nil {
			err = fmt.Errorf("Error setting source_secret: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_build", "read", "set-source_secret").GetDiag()
		}
	}
	if !core.IsNil(build.SourceType) {
		if err = d.Set("source_type", build.SourceType); err != nil {
			err = fmt.Errorf("Error setting source_type: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_build", "read", "set-source_type").GetDiag()
		}
	}
	if !core.IsNil(build.SourceURL) {
		if err = d.Set("source_url", build.SourceURL); err != nil {
			err = fmt.Errorf("Error setting source_url: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_build", "read", "set-source_url").GetDiag()
		}
	}
	if !core.IsNil(build.StrategySize) {
		if err = d.Set("strategy_size", build.StrategySize); err != nil {
			err = fmt.Errorf("Error setting strategy_size: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_build", "read", "set-strategy_size").GetDiag()
		}
	}
	if !core.IsNil(build.StrategySpecFile) {
		if err = d.Set("strategy_spec_file", build.StrategySpecFile); err != nil {
			err = fmt.Errorf("Error setting strategy_spec_file: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_build", "read", "set-strategy_spec_file").GetDiag()
		}
	}
	if err = d.Set("strategy_type", build.StrategyType); err != nil {
		err = fmt.Errorf("Error setting strategy_type: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_build", "read", "set-strategy_type").GetDiag()
	}
	if !core.IsNil(build.Timeout) {
		if err = d.Set("timeout", flex.IntValue(build.Timeout)); err != nil {
			err = fmt.Errorf("Error setting timeout: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_build", "read", "set-timeout").GetDiag()
		}
	}
	if !core.IsNil(build.CreatedAt) {
		if err = d.Set("created_at", build.CreatedAt); err != nil {
			err = fmt.Errorf("Error setting created_at: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_build", "read", "set-created_at").GetDiag()
		}
	}
	if err = d.Set("entity_tag", build.EntityTag); err != nil {
		err = fmt.Errorf("Error setting entity_tag: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_build", "read", "set-entity_tag").GetDiag()
	}
	if !core.IsNil(build.Href) {
		if err = d.Set("href", build.Href); err != nil {
			err = fmt.Errorf("Error setting href: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_build", "read", "set-href").GetDiag()
		}
	}
	if !core.IsNil(build.ID) {
		if err = d.Set("build_id", build.ID); err != nil {
			err = fmt.Errorf("Error setting build_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_build", "read", "set-build_id").GetDiag()
		}
	}
	if !core.IsNil(build.Region) {
		if err = d.Set("region", build.Region); err != nil {
			err = fmt.Errorf("Error setting region: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_build", "read", "set-region").GetDiag()
		}
	}
	if !core.IsNil(build.ResourceType) {
		if err = d.Set("resource_type", build.ResourceType); err != nil {
			err = fmt.Errorf("Error setting resource_type: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_build", "read", "set-resource_type").GetDiag()
		}
	}
	if !core.IsNil(build.Status) {
		if err = d.Set("status", build.Status); err != nil {
			err = fmt.Errorf("Error setting status: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_build", "read", "set-status").GetDiag()
		}
	}
	if !core.IsNil(build.StatusDetails) {
		statusDetailsMap, err := ResourceIbmCodeEngineBuildBuildStatusToMap(build.StatusDetails)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_build", "read", "status_details-to-map").GetDiag()
		}
		if err = d.Set("status_details", []map[string]interface{}{statusDetailsMap}); err != nil {
			err = fmt.Errorf("Error setting status_details: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_build", "read", "set-status_details").GetDiag()
		}
	}
	if err = d.Set("etag", response.Headers.Get("Etag")); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting etag: %s", err), "ibm_code_engine_build", "read", "set-etag").GetDiag()
	}

	return nil
}

func resourceIbmCodeEngineBuildUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_build", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateBuildOptions := &codeenginev2.UpdateBuildOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_build", "update", "sep-id-parts").GetDiag()
	}

	updateBuildOptions.SetProjectID(parts[0])
	updateBuildOptions.SetName(parts[1])

	hasChange := false

	patchVals := &codeenginev2.BuildPatch{}
	if d.HasChange("project_id") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "project_id")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_code_engine_build", "update", "project_id-forces-new").GetDiag()
	}
	if d.HasChange("output_image") {
		newOutputImage := d.Get("output_image").(string)
		patchVals.OutputImage = &newOutputImage
		hasChange = true
	}
	if d.HasChange("output_secret") {
		newOutputSecret := d.Get("output_secret").(string)
		patchVals.OutputSecret = &newOutputSecret
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
	if d.HasChange("strategy_type") {
		newStrategyType := d.Get("strategy_type").(string)
		patchVals.StrategyType = &newStrategyType
		hasChange = true
	}
	if d.HasChange("timeout") {
		newTimeout := int64(d.Get("timeout").(int))
		patchVals.Timeout = &newTimeout
		hasChange = true
	}
	updateBuildOptions.SetIfMatch(d.Get("etag").(string))

	if hasChange {
		// Fields with `nil` values are omitted from the generic map,
		// so we need to re-add them to support removing arguments
		// in merge-patch operations sent to the service.
		updateBuildOptions.Build = ResourceIbmCodeEngineBuildBuildPatchAsPatch(patchVals, d)

		_, _, err = codeEngineClient.UpdateBuildWithContext(context, updateBuildOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateBuildWithContext failed: %s", err.Error()), "ibm_code_engine_build", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIbmCodeEngineBuildRead(context, d, meta)
}

func resourceIbmCodeEngineBuildDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_build", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteBuildOptions := &codeenginev2.DeleteBuildOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_build", "delete", "sep-id-parts").GetDiag()
	}

	deleteBuildOptions.SetProjectID(parts[0])
	deleteBuildOptions.SetName(parts[1])

	_, err = codeEngineClient.DeleteBuildWithContext(context, deleteBuildOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteBuildWithContext failed: %s", err.Error()), "ibm_code_engine_build", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIbmCodeEngineBuildBuildStatusToMap(model *codeenginev2.BuildStatus) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Reason != nil {
		modelMap["reason"] = *model.Reason
	}
	return modelMap, nil
}

func ResourceIbmCodeEngineBuildBuildPatchAsPatch(patchVals *codeenginev2.BuildPatch, d *schema.ResourceData) map[string]interface{} {
	patch, _ := patchVals.AsPatch()
	var path string

	path = "output_image"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["output_image"] = nil
	}
	path = "output_secret"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["output_secret"] = nil
	}
	path = "source_context_dir"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["source_context_dir"] = nil
	}
	path = "source_revision"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["source_revision"] = nil
	}
	path = "source_secret"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["source_secret"] = nil
	}
	path = "source_type"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["source_type"] = nil
	}
	path = "source_url"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["source_url"] = nil
	}
	path = "strategy_size"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["strategy_size"] = nil
	}
	path = "strategy_spec_file"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["strategy_spec_file"] = nil
	}
	path = "strategy_type"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["strategy_type"] = nil
	}
	path = "timeout"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["timeout"] = nil
	}

	return patch
}
