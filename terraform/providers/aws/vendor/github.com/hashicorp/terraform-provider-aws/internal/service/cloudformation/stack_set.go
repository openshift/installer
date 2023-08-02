package cloudformation

import (
	"context"
	"log"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_cloudformation_stack_set", name="Stack Set")
// @Tags
func ResourceStackSet() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceStackSetCreate,
		ReadWithoutTimeout:   resourceStackSetRead,
		UpdateWithoutTimeout: resourceStackSetUpdate,
		DeleteWithoutTimeout: resourceStackSetDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"administration_role_arn": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"auto_deployment"},
				ValidateFunc:  verify.ValidARN,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_deployment": {
				Type:     schema.TypeList,
				MinItems: 1,
				MaxItems: 1,
				Optional: true,
				ForceNew: true,
				ConflictsWith: []string{
					"administration_role_arn",
					"execution_role_name",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"retain_stacks_on_account_removal": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"call_as": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(cloudformation.CallAs_Values(), false),
				Default:      cloudformation.CallAsSelf,
			},
			"capabilities": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(cloudformation.Capability_Values(), false),
				},
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 1024),
			},
			"execution_role_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"auto_deployment"},
			},
			"managed_execution": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"active": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
				DiffSuppressFunc: verify.SuppressMissingOptionalConfigurationBlock,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 128),
					validation.StringMatch(regexp.MustCompile(`^[a-zA-Z]`), "must begin with alphabetic character"),
					validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9-]+$`), "must contain only alphanumeric and hyphen characters"),
				),
			},
			"operation_preferences": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"failure_tolerance_count": {
							Type:          schema.TypeInt,
							Optional:      true,
							ValidateFunc:  validation.IntAtLeast(0),
							ConflictsWith: []string{"operation_preferences.0.failure_tolerance_percentage"},
						},
						"failure_tolerance_percentage": {
							Type:          schema.TypeInt,
							Optional:      true,
							ValidateFunc:  validation.IntBetween(0, 100),
							ConflictsWith: []string{"operation_preferences.0.failure_tolerance_count"},
						},
						"max_concurrent_count": {
							Type:          schema.TypeInt,
							Optional:      true,
							ValidateFunc:  validation.IntAtLeast(1),
							ConflictsWith: []string{"operation_preferences.0.max_concurrent_percentage"},
						},
						"max_concurrent_percentage": {
							Type:          schema.TypeInt,
							Optional:      true,
							ValidateFunc:  validation.IntBetween(1, 100),
							ConflictsWith: []string{"operation_preferences.0.max_concurrent_count"},
						},
						"region_concurrency_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(cloudformation.RegionConcurrencyType_Values(), false),
						},
						"region_order": {
							Type:     schema.TypeList,
							Optional: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9-]{1,128}$`), ""),
							},
						},
					},
				},
			},
			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"permission_model": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(cloudformation.PermissionModels_Values(), false),
				Default:      cloudformation.PermissionModelsSelfManaged,
			},
			"stack_set_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
			"template_body": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ConflictsWith:    []string{"template_url"},
				DiffSuppressFunc: verify.SuppressEquivalentJSONOrYAMLDiffs,
				ValidateFunc:     verify.ValidStringIsJSONOrYAML,
			},
			"template_url": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"template_body"},
			},
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceStackSetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CloudFormationConn(ctx)

	name := d.Get("name").(string)
	input := &cloudformation.CreateStackSetInput{
		ClientRequestToken: aws.String(id.UniqueId()),
		StackSetName:       aws.String(name),
		Tags:               GetTagsIn(ctx),
	}

	if v, ok := d.GetOk("administration_role_arn"); ok {
		input.AdministrationRoleARN = aws.String(v.(string))
	}

	if v, ok := d.GetOk("auto_deployment"); ok {
		input.AutoDeployment = expandAutoDeployment(v.([]interface{}))
	}

	if v, ok := d.GetOk("capabilities"); ok {
		input.Capabilities = flex.ExpandStringSet(v.(*schema.Set))
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = aws.String(v.(string))
	}

	if v, ok := d.GetOk("execution_role_name"); ok {
		input.ExecutionRoleName = aws.String(v.(string))
	}

	if v, ok := d.GetOk("managed_execution"); ok {
		input.ManagedExecution = expandManagedExecution(v.([]interface{}))
	}

	if v, ok := d.GetOk("parameters"); ok {
		input.Parameters = expandParameters(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("permission_model"); ok {
		input.PermissionModel = aws.String(v.(string))
	}

	if v, ok := d.GetOk("call_as"); ok {
		input.CallAs = aws.String(v.(string))
	}

	if v, ok := d.GetOk("template_body"); ok {
		input.TemplateBody = aws.String(v.(string))
	}

	if v, ok := d.GetOk("template_url"); ok {
		input.TemplateURL = aws.String(v.(string))
	}

	log.Printf("[DEBUG] Creating CloudFormation StackSet: %s", input)
	_, err := conn.CreateStackSetWithContext(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating CloudFormation StackSet (%s): %s", name, err)
	}

	d.SetId(name)

	return append(diags, resourceStackSetRead(ctx, d, meta)...)
}

func resourceStackSetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CloudFormationConn(ctx)

	callAs := d.Get("call_as").(string)
	stackSet, err := FindStackSetByName(ctx, conn, d.Id(), callAs)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] CloudFormation StackSet (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading CloudFormation StackSet (%s): %s", d.Id(), err)
	}

	d.Set("administration_role_arn", stackSet.AdministrationRoleARN)
	d.Set("arn", stackSet.StackSetARN)

	if err := d.Set("auto_deployment", flattenStackSetAutoDeploymentResponse(stackSet.AutoDeployment)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting auto_deployment: %s", err)
	}

	if err := d.Set("capabilities", aws.StringValueSlice(stackSet.Capabilities)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting capabilities: %s", err)
	}

	d.Set("description", stackSet.Description)
	d.Set("execution_role_name", stackSet.ExecutionRoleName)

	if err := d.Set("managed_execution", flattenStackSetManagedExecution(stackSet.ManagedExecution)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting managed_execution: %s", err)
	}

	d.Set("name", stackSet.StackSetName)
	d.Set("permission_model", stackSet.PermissionModel)

	if err := d.Set("parameters", flattenAllParameters(stackSet.Parameters)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting parameters: %s", err)
	}

	d.Set("stack_set_id", stackSet.StackSetId)

	SetTagsOut(ctx, stackSet.Tags)

	d.Set("template_body", stackSet.TemplateBody)

	return diags
}

func resourceStackSetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CloudFormationConn(ctx)

	input := &cloudformation.UpdateStackSetInput{
		OperationId:  aws.String(id.UniqueId()),
		StackSetName: aws.String(d.Id()),
		Tags:         []*cloudformation.Tag{},
		TemplateBody: aws.String(d.Get("template_body").(string)),
	}

	if v, ok := d.GetOk("administration_role_arn"); ok {
		input.AdministrationRoleARN = aws.String(v.(string))
	}

	if v, ok := d.GetOk("capabilities"); ok {
		input.Capabilities = flex.ExpandStringSet(v.(*schema.Set))
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = aws.String(v.(string))
	}

	if v, ok := d.GetOk("execution_role_name"); ok {
		input.ExecutionRoleName = aws.String(v.(string))
	}

	if v, ok := d.GetOk("managed_execution"); ok {
		input.ManagedExecution = expandManagedExecution(v.([]interface{}))
	}

	if v, ok := d.GetOk("operation_preferences"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		input.OperationPreferences = expandOperationPreferences(v.([]interface{})[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("parameters"); ok {
		input.Parameters = expandParameters(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("permission_model"); ok {
		input.PermissionModel = aws.String(v.(string))
	}

	callAs := d.Get("call_as").(string)
	if v, ok := d.GetOk("call_as"); ok {
		input.CallAs = aws.String(v.(string))
	}

	if tags := GetTagsIn(ctx); len(tags) > 0 {
		input.Tags = tags
	}

	if v, ok := d.GetOk("template_url"); ok {
		// ValidationError: Exactly one of TemplateBody or TemplateUrl must be specified
		// TemplateBody is always present when TemplateUrl is used so remove TemplateBody if TemplateUrl is set
		input.TemplateBody = nil
		input.TemplateURL = aws.String(v.(string))
	}

	// When `auto_deployment` is set, ignore `administration_role_arn` and
	// `execution_role_name` fields since it's using the SERVICE_MANAGED
	// permission model
	if v, ok := d.GetOk("auto_deployment"); ok {
		input.AdministrationRoleARN = nil
		input.ExecutionRoleName = nil
		input.AutoDeployment = expandAutoDeployment(v.([]interface{}))
	}

	log.Printf("[DEBUG] Updating CloudFormation StackSet: %s", input)
	output, err := conn.UpdateStackSetWithContext(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "updating CloudFormation StackSet (%s): %s", d.Id(), err)
	}

	if _, err := WaitStackSetOperationSucceeded(ctx, conn, d.Id(), aws.StringValue(output.OperationId), callAs, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for CloudFormation StackSet (%s) update: %s", d.Id(), err)
	}

	return append(diags, resourceStackSetRead(ctx, d, meta)...)
}

func resourceStackSetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CloudFormationConn(ctx)

	input := &cloudformation.DeleteStackSetInput{
		StackSetName: aws.String(d.Id()),
	}

	if v, ok := d.GetOk("call_as"); ok {
		input.CallAs = aws.String(v.(string))
	}

	log.Printf("[DEBUG] Deleting CloudFormation StackSet: %s", d.Id())
	_, err := conn.DeleteStackSetWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, cloudformation.ErrCodeStackSetNotFoundException) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting CloudFormation StackSet (%s): %s", d.Id(), err)
	}

	return diags
}

func expandAutoDeployment(l []interface{}) *cloudformation.AutoDeployment {
	if len(l) == 0 {
		return nil
	}

	m := l[0].(map[string]interface{})

	autoDeployment := &cloudformation.AutoDeployment{
		Enabled:                      aws.Bool(m["enabled"].(bool)),
		RetainStacksOnAccountRemoval: aws.Bool(m["retain_stacks_on_account_removal"].(bool)),
	}

	return autoDeployment
}

func expandManagedExecution(l []interface{}) *cloudformation.ManagedExecution {
	if len(l) == 0 {
		return nil
	}

	m := l[0].(map[string]interface{})

	managedExecution := &cloudformation.ManagedExecution{
		Active: aws.Bool(m["active"].(bool)),
	}

	return managedExecution
}

func flattenStackSetAutoDeploymentResponse(autoDeployment *cloudformation.AutoDeployment) []map[string]interface{} {
	if autoDeployment == nil {
		return []map[string]interface{}{}
	}

	m := map[string]interface{}{
		"enabled":                          aws.BoolValue(autoDeployment.Enabled),
		"retain_stacks_on_account_removal": aws.BoolValue(autoDeployment.RetainStacksOnAccountRemoval),
	}

	return []map[string]interface{}{m}
}

func flattenStackSetManagedExecution(managedExecution *cloudformation.ManagedExecution) []map[string]interface{} {
	if managedExecution == nil {
		return []map[string]interface{}{}
	}

	m := map[string]interface{}{
		"active": aws.BoolValue(managedExecution.Active),
	}

	return []map[string]interface{}{m}
}
