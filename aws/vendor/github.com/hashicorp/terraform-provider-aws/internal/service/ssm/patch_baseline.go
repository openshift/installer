package ssm

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_ssm_patch_baseline", name="Patch Baseline")
// @Tags(identifierAttribute="id", resourceType="PatchBaseline")
func ResourcePatchBaseline() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourcePatchBaselineCreate,
		ReadWithoutTimeout:   resourcePatchBaselineRead,
		UpdateWithoutTimeout: resourcePatchBaselineUpdate,
		DeleteWithoutTimeout: resourcePatchBaselineDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(3, 128),
					validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9_\-.]{3,128}$`), "must contain only alphanumeric, underscore, hyphen, or period characters"),
				),
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 1024),
			},

			"global_filter": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 4,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(ssm.PatchFilterKey_Values(), false),
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 20,
							MinItems: 1,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringLenBetween(1, 64),
							},
						},
					},
				},
			},

			"approval_rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"approve_after_days": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 100),
						},

						"approve_until_date": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile(`([12]\d{3}-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01]))`), "must be formatted YYYY-MM-DD"),
						},

						"compliance_level": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      ssm.PatchComplianceLevelUnspecified,
							ValidateFunc: validation.StringInSlice(ssm.PatchComplianceLevel_Values(), false),
						},

						"enable_non_security": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},

						"patch_filter": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 10,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(ssm.PatchFilterKey_Values(), false),
									},
									"values": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 20,
										MinItems: 1,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringLenBetween(1, 64),
										},
									},
								},
							},
						},
					},
				},
			},

			"approved_patches": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 50,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringLenBetween(1, 100),
				},
			},

			"rejected_patches": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 50,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringLenBetween(1, 100),
				},
			},

			"operating_system": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      ssm.OperatingSystemWindows,
				ValidateFunc: validation.StringInSlice(ssm.OperatingSystem_Values(), false),
			},

			"approved_patches_compliance_level": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      ssm.PatchComplianceLevelUnspecified,
				ValidateFunc: validation.StringInSlice(ssm.PatchComplianceLevel_Values(), false),
			},
			"rejected_patches_action": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(ssm.PatchAction_Values(), false),
			},
			"approved_patches_enable_non_security": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"source": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 20,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.All(
								validation.StringLenBetween(3, 50),
								validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9_\-.]{3,50}$`), "must contain only alphanumeric, underscore, hyphen, or period characters"),
							),
						},

						"configuration": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 1024),
						},

						"products": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 20,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringLenBetween(1, 128),
							},
						},
					},
				},
			},

			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

const (
	resNamePatchBaseline = "Patch Baseline"
)

func resourcePatchBaselineCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).SSMConn(ctx)

	name := d.Get("name").(string)
	input := &ssm.CreatePatchBaselineInput{
		ApprovedPatchesComplianceLevel: aws.String(d.Get("approved_patches_compliance_level").(string)),
		Name:                           aws.String(name),
		OperatingSystem:                aws.String(d.Get("operating_system").(string)),
		Tags:                           GetTagsIn(ctx),
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = aws.String(v.(string))
	}

	if v, ok := d.GetOk("approved_patches"); ok && v.(*schema.Set).Len() > 0 {
		input.ApprovedPatches = flex.ExpandStringSet(v.(*schema.Set))
	}

	if v, ok := d.GetOk("rejected_patches"); ok && v.(*schema.Set).Len() > 0 {
		input.RejectedPatches = flex.ExpandStringSet(v.(*schema.Set))
	}

	if _, ok := d.GetOk("global_filter"); ok {
		input.GlobalFilters = expandPatchFilterGroup(d)
	}

	if _, ok := d.GetOk("approval_rule"); ok {
		input.ApprovalRules = expandPatchRuleGroup(d)
	}

	if _, ok := d.GetOk("source"); ok {
		input.Sources = expandPatchSource(d)
	}

	if v, ok := d.GetOk("approved_patches_enable_non_security"); ok {
		input.ApprovedPatchesEnableNonSecurity = aws.Bool(v.(bool))
	}

	if v, ok := d.GetOk("rejected_patches_action"); ok {
		input.RejectedPatchesAction = aws.String(v.(string))
	}

	output, err := conn.CreatePatchBaselineWithContext(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating SSM Patch Baseline (%s): %s", name, err)
	}

	d.SetId(aws.StringValue(output.BaselineId))

	return append(diags, resourcePatchBaselineRead(ctx, d, meta)...)
}

func resourcePatchBaselineRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).SSMConn(ctx)

	params := &ssm.GetPatchBaselineInput{
		BaselineId: aws.String(d.Id()),
	}

	resp, err := conn.GetPatchBaselineWithContext(ctx, params)
	if err != nil {
		if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, ssm.ErrCodeDoesNotExistException) {
			log.Printf("[WARN] SSM Patch Baseline (%s) not found, removing from state", d.Id())
			d.SetId("")
			return diags
		}
		return sdkdiag.AppendErrorf(diags, "reading SSM Patch Baseline (%s): %s", d.Id(), err)
	}

	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Region:    meta.(*conns.AWSClient).Region,
		Service:   "ssm",
		AccountID: meta.(*conns.AWSClient).AccountID,
		Resource:  fmt.Sprintf("patchbaseline/%s", strings.TrimPrefix(d.Id(), "/")),
	}
	d.Set("arn", arn.String())
	d.Set("name", resp.Name)
	d.Set("description", resp.Description)
	d.Set("operating_system", resp.OperatingSystem)
	d.Set("approved_patches_compliance_level", resp.ApprovedPatchesComplianceLevel)
	d.Set("approved_patches", flex.FlattenStringList(resp.ApprovedPatches))
	d.Set("rejected_patches", flex.FlattenStringList(resp.RejectedPatches))
	d.Set("rejected_patches_action", resp.RejectedPatchesAction)
	d.Set("approved_patches_enable_non_security", resp.ApprovedPatchesEnableNonSecurity)

	if err := d.Set("global_filter", flattenPatchFilterGroup(resp.GlobalFilters)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting global filters: %s", err)
	}

	if err := d.Set("approval_rule", flattenPatchRuleGroup(resp.ApprovalRules)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting approval rules: %s", err)
	}

	if err := d.Set("source", flattenPatchSource(resp.Sources)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting patch sources: %s", err)
	}

	return diags
}

func resourcePatchBaselineUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).SSMConn(ctx)

	if d.HasChangesExcept("tags", "tags_all") {
		input := &ssm.UpdatePatchBaselineInput{
			BaselineId: aws.String(d.Id()),
		}

		if d.HasChange("name") {
			input.Name = aws.String(d.Get("name").(string))
		}

		if d.HasChange("description") {
			input.Description = aws.String(d.Get("description").(string))
		}

		if d.HasChange("approved_patches") {
			input.ApprovedPatches = flex.ExpandStringSet(d.Get("approved_patches").(*schema.Set))
		}

		if d.HasChange("rejected_patches") {
			input.RejectedPatches = flex.ExpandStringSet(d.Get("rejected_patches").(*schema.Set))
		}

		if d.HasChange("approved_patches_compliance_level") {
			input.ApprovedPatchesComplianceLevel = aws.String(d.Get("approved_patches_compliance_level").(string))
		}

		if d.HasChange("approval_rule") {
			input.ApprovalRules = expandPatchRuleGroup(d)
		}

		if d.HasChange("global_filter") {
			input.GlobalFilters = expandPatchFilterGroup(d)
		}

		if d.HasChange("source") {
			input.Sources = expandPatchSource(d)
		}

		if d.HasChange("approved_patches_enable_non_security") {
			input.ApprovedPatchesEnableNonSecurity = aws.Bool(d.Get("approved_patches_enable_non_security").(bool))
		}

		if d.HasChange("rejected_patches_action") {
			input.RejectedPatchesAction = aws.String(d.Get("rejected_patches_action").(string))
		}

		_, err := conn.UpdatePatchBaselineWithContext(ctx, input)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating SSM Patch Baseline (%s): %s", d.Id(), err)
		}
	}

	return append(diags, resourcePatchBaselineRead(ctx, d, meta)...)
}

func resourcePatchBaselineDelete(ctx context.Context, d *schema.ResourceData, meta any) (diags diag.Diagnostics) {
	conn := meta.(*conns.AWSClient).SSMConn(ctx)

	log.Printf("[INFO] Deleting SSM Patch Baseline: %s", d.Id())

	params := &ssm.DeletePatchBaselineInput{
		BaselineId: aws.String(d.Id()),
	}

	_, err := conn.DeletePatchBaselineWithContext(ctx, params)
	if tfawserr.ErrCodeEquals(err, ssm.ErrCodeResourceInUseException) {
		// Reset the default patch baseline before retrying
		diags = append(diags, defaultPatchBaselineRestoreOSDefault(ctx, meta.(*conns.AWSClient).SSMClient(ctx), types.OperatingSystem(d.Get("operating_system").(string)))...)
		if diags.HasError() {
			return
		}
		_, err = conn.DeletePatchBaselineWithContext(ctx, params)
	}
	if err != nil {
		diags = sdkdiag.AppendErrorf(diags, "deleting SSM Patch Baseline (%s): %s", d.Id(), err)
	}

	return
}

func expandPatchFilterGroup(d *schema.ResourceData) *ssm.PatchFilterGroup {
	var filters []*ssm.PatchFilter

	filterConfig := d.Get("global_filter").([]interface{})

	for _, fConfig := range filterConfig {
		config := fConfig.(map[string]interface{})

		filter := &ssm.PatchFilter{
			Key:    aws.String(config["key"].(string)),
			Values: flex.ExpandStringList(config["values"].([]interface{})),
		}

		filters = append(filters, filter)
	}

	return &ssm.PatchFilterGroup{
		PatchFilters: filters,
	}
}

func flattenPatchFilterGroup(group *ssm.PatchFilterGroup) []map[string]interface{} {
	if len(group.PatchFilters) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(group.PatchFilters))

	for _, filter := range group.PatchFilters {
		f := make(map[string]interface{})
		f["key"] = aws.StringValue(filter.Key)
		f["values"] = flex.FlattenStringList(filter.Values)

		result = append(result, f)
	}

	return result
}

func expandPatchRuleGroup(d *schema.ResourceData) *ssm.PatchRuleGroup {
	var rules []*ssm.PatchRule

	ruleConfig := d.Get("approval_rule").([]interface{})

	for _, rConfig := range ruleConfig {
		rCfg := rConfig.(map[string]interface{})

		var filters []*ssm.PatchFilter
		filterConfig := rCfg["patch_filter"].([]interface{})

		for _, fConfig := range filterConfig {
			fCfg := fConfig.(map[string]interface{})

			filter := &ssm.PatchFilter{
				Key:    aws.String(fCfg["key"].(string)),
				Values: flex.ExpandStringList(fCfg["values"].([]interface{})),
			}

			filters = append(filters, filter)
		}

		filterGroup := &ssm.PatchFilterGroup{
			PatchFilters: filters,
		}

		rule := &ssm.PatchRule{
			PatchFilterGroup:  filterGroup,
			ComplianceLevel:   aws.String(rCfg["compliance_level"].(string)),
			EnableNonSecurity: aws.Bool(rCfg["enable_non_security"].(bool)),
		}

		if v, ok := rCfg["approve_until_date"].(string); ok && v != "" {
			rule.ApproveUntilDate = aws.String(v)
		} else if v, ok := rCfg["approve_after_days"].(int); ok {
			rule.ApproveAfterDays = aws.Int64(int64(v))
		}

		rules = append(rules, rule)
	}

	return &ssm.PatchRuleGroup{
		PatchRules: rules,
	}
}

func flattenPatchRuleGroup(group *ssm.PatchRuleGroup) []map[string]interface{} {
	if len(group.PatchRules) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(group.PatchRules))

	for _, rule := range group.PatchRules {
		r := make(map[string]interface{})
		r["compliance_level"] = aws.StringValue(rule.ComplianceLevel)
		r["enable_non_security"] = aws.BoolValue(rule.EnableNonSecurity)
		r["patch_filter"] = flattenPatchFilterGroup(rule.PatchFilterGroup)

		if rule.ApproveAfterDays != nil {
			r["approve_after_days"] = aws.Int64Value(rule.ApproveAfterDays)
		}

		if rule.ApproveUntilDate != nil {
			r["approve_until_date"] = aws.StringValue(rule.ApproveUntilDate)
		}

		result = append(result, r)
	}

	return result
}

func expandPatchSource(d *schema.ResourceData) []*ssm.PatchSource {
	var sources []*ssm.PatchSource

	sourceConfigs := d.Get("source").([]interface{})

	for _, sConfig := range sourceConfigs {
		config := sConfig.(map[string]interface{})

		source := &ssm.PatchSource{
			Name:          aws.String(config["name"].(string)),
			Configuration: aws.String(config["configuration"].(string)),
			Products:      flex.ExpandStringList(config["products"].([]interface{})),
		}

		sources = append(sources, source)
	}

	return sources
}

func flattenPatchSource(sources []*ssm.PatchSource) []map[string]interface{} {
	if len(sources) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(sources))

	for _, source := range sources {
		s := make(map[string]interface{})
		s["name"] = aws.StringValue(source.Name)
		s["configuration"] = aws.StringValue(source.Configuration)
		s["products"] = flex.FlattenStringList(source.Products)
		result = append(result, s)
	}

	return result
}
