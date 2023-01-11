package ssm

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourcePatchBaseline() *schema.Resource {
	return &schema.Resource{
		Create: resourcePatchBaselineCreate,
		Read:   resourcePatchBaselineRead,
		Update: resourcePatchBaselineUpdate,
		Delete: resourcePatchBaselineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourcePatchBaselineCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SSMConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	params := &ssm.CreatePatchBaselineInput{
		Name:                           aws.String(d.Get("name").(string)),
		ApprovedPatchesComplianceLevel: aws.String(d.Get("approved_patches_compliance_level").(string)),
		OperatingSystem:                aws.String(d.Get("operating_system").(string)),
	}

	if len(tags) > 0 {
		params.Tags = Tags(tags.IgnoreAWS())
	}

	if v, ok := d.GetOk("description"); ok {
		params.Description = aws.String(v.(string))
	}

	if v, ok := d.GetOk("approved_patches"); ok && v.(*schema.Set).Len() > 0 {
		params.ApprovedPatches = flex.ExpandStringSet(v.(*schema.Set))
	}

	if v, ok := d.GetOk("rejected_patches"); ok && v.(*schema.Set).Len() > 0 {
		params.RejectedPatches = flex.ExpandStringSet(v.(*schema.Set))
	}

	if _, ok := d.GetOk("global_filter"); ok {
		params.GlobalFilters = expandPatchFilterGroup(d)
	}

	if _, ok := d.GetOk("approval_rule"); ok {
		params.ApprovalRules = expandPatchRuleGroup(d)
	}

	if _, ok := d.GetOk("source"); ok {
		params.Sources = expandPatchSource(d)
	}

	if v, ok := d.GetOk("approved_patches_enable_non_security"); ok {
		params.ApprovedPatchesEnableNonSecurity = aws.Bool(v.(bool))
	}

	if v, ok := d.GetOk("rejected_patches_action"); ok {
		params.RejectedPatchesAction = aws.String(v.(string))
	}

	resp, err := conn.CreatePatchBaseline(params)

	if err != nil {
		return err
	}

	d.SetId(aws.StringValue(resp.BaselineId))
	return resourcePatchBaselineRead(d, meta)
}

func resourcePatchBaselineUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SSMConn

	params := &ssm.UpdatePatchBaselineInput{
		BaselineId: aws.String(d.Id()),
	}

	if d.HasChange("name") {
		params.Name = aws.String(d.Get("name").(string))
	}

	if d.HasChange("description") {
		params.Description = aws.String(d.Get("description").(string))
	}

	if d.HasChange("approved_patches") {
		params.ApprovedPatches = flex.ExpandStringSet(d.Get("approved_patches").(*schema.Set))
	}

	if d.HasChange("rejected_patches") {
		params.RejectedPatches = flex.ExpandStringSet(d.Get("rejected_patches").(*schema.Set))
	}

	if d.HasChange("approved_patches_compliance_level") {
		params.ApprovedPatchesComplianceLevel = aws.String(d.Get("approved_patches_compliance_level").(string))
	}

	if d.HasChange("approval_rule") {
		params.ApprovalRules = expandPatchRuleGroup(d)
	}

	if d.HasChange("global_filter") {
		params.GlobalFilters = expandPatchFilterGroup(d)
	}

	if d.HasChange("source") {
		params.Sources = expandPatchSource(d)
	}

	if d.HasChange("approved_patches_enable_non_security") {
		params.ApprovedPatchesEnableNonSecurity = aws.Bool(d.Get("approved_patches_enable_non_security").(bool))
	}

	if d.HasChange("rejected_patches_action") {
		params.RejectedPatchesAction = aws.String(d.Get("rejected_patches_action").(string))
	}

	if d.HasChangesExcept("tags", "tags_all") {
		_, err := conn.UpdatePatchBaseline(params)
		if err != nil {
			return fmt.Errorf("error updating SSM Patch Baseline (%s): %w", d.Id(), err)
		}
	}

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := UpdateTags(conn, d.Id(), ssm.ResourceTypeForTaggingPatchBaseline, o, n); err != nil {
			return fmt.Errorf("error updating SSM Patch Baseline (%s) tags: %s", d.Id(), err)
		}
	}

	return resourcePatchBaselineRead(d, meta)
}
func resourcePatchBaselineRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SSMConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	params := &ssm.GetPatchBaselineInput{
		BaselineId: aws.String(d.Id()),
	}

	resp, err := conn.GetPatchBaseline(params)
	if err != nil {
		if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, ssm.ErrCodeDoesNotExistException) {
			log.Printf("[WARN] SSM Patch Baseline (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", resp.Name)
	d.Set("description", resp.Description)
	d.Set("operating_system", resp.OperatingSystem)
	d.Set("approved_patches_compliance_level", resp.ApprovedPatchesComplianceLevel)
	d.Set("approved_patches", flex.FlattenStringList(resp.ApprovedPatches))
	d.Set("rejected_patches", flex.FlattenStringList(resp.RejectedPatches))
	d.Set("rejected_patches_action", resp.RejectedPatchesAction)
	d.Set("approved_patches_enable_non_security", resp.ApprovedPatchesEnableNonSecurity)

	if err := d.Set("global_filter", flattenPatchFilterGroup(resp.GlobalFilters)); err != nil {
		return fmt.Errorf("Error setting global filters error: %#v", err)
	}

	if err := d.Set("approval_rule", flattenPatchRuleGroup(resp.ApprovalRules)); err != nil {
		return fmt.Errorf("Error setting approval rules error: %#v", err)
	}

	if err := d.Set("source", flattenPatchSource(resp.Sources)); err != nil {
		return fmt.Errorf("Error setting patch sources error: %#v", err)
	}

	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Region:    meta.(*conns.AWSClient).Region,
		Service:   "ssm",
		AccountID: meta.(*conns.AWSClient).AccountID,
		Resource:  fmt.Sprintf("patchbaseline/%s", strings.TrimPrefix(d.Id(), "/")),
	}
	d.Set("arn", arn.String())

	tags, err := ListTags(conn, d.Id(), ssm.ResourceTypeForTaggingPatchBaseline)

	if err != nil {
		return fmt.Errorf("error listing tags for SSM Patch Baseline (%s): %s", d.Id(), err)
	}

	tags = tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("error setting tags_all: %w", err)
	}

	return nil
}

func resourcePatchBaselineDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SSMConn

	log.Printf("[INFO] Deleting SSM Patch Baseline: %s", d.Id())

	params := &ssm.DeletePatchBaselineInput{
		BaselineId: aws.String(d.Id()),
	}

	_, err := conn.DeletePatchBaseline(params)
	if err != nil {
		return fmt.Errorf("error deleting SSM Patch Baseline (%s): %s", d.Id(), err)
	}

	return nil
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
