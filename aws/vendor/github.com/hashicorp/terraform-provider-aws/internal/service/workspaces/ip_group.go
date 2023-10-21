package workspaces

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/workspaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_workspaces_ip_group", name="IP Group")
// @Tags(identifierAttribute="id")
func ResourceIPGroup() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceIPGroupCreate,
		ReadWithoutTimeout:   resourceIPGroupRead,
		UpdateWithoutTimeout: resourceIPGroupUpdate,
		DeleteWithoutTimeout: resourceIPGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"rules": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsCIDR,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
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

func resourceIPGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).WorkSpacesConn(ctx)

	rules := d.Get("rules").(*schema.Set).List()
	resp, err := conn.CreateIpGroupWithContext(ctx, &workspaces.CreateIpGroupInput{
		GroupName: aws.String(d.Get("name").(string)),
		GroupDesc: aws.String(d.Get("description").(string)),
		UserRules: expandIPGroupRules(rules),
		Tags:      GetTagsIn(ctx),
	})

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating WorkSpaces IP Group: %s", err)
	}

	d.SetId(aws.StringValue(resp.GroupId))

	return append(diags, resourceIPGroupRead(ctx, d, meta)...)
}

func resourceIPGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).WorkSpacesConn(ctx)

	resp, err := conn.DescribeIpGroupsWithContext(ctx, &workspaces.DescribeIpGroupsInput{
		GroupIds: []*string{aws.String(d.Id())},
	})
	if err != nil {
		if len(resp.Result) == 0 {
			log.Printf("[WARN] WorkSpaces IP Group (%s) not found, removing from state", d.Id())
			d.SetId("")
			return diags
		}

		return sdkdiag.AppendErrorf(diags, "reading WorkSpaces IP Group (%s): %s", d.Id(), err)
	}

	ipGroups := resp.Result

	if len(ipGroups) == 0 {
		log.Printf("[WARN] WorkSpaces Ip Group (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	ipGroup := ipGroups[0]

	d.Set("name", ipGroup.GroupName)
	d.Set("description", ipGroup.GroupDesc)
	d.Set("rules", flattenIPGroupRules(ipGroup.UserRules))

	return diags
}

func resourceIPGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).WorkSpacesConn(ctx)

	if d.HasChange("rules") {
		rules := d.Get("rules").(*schema.Set).List()

		_, err := conn.UpdateRulesOfIpGroupWithContext(ctx, &workspaces.UpdateRulesOfIpGroupInput{
			GroupId:   aws.String(d.Id()),
			UserRules: expandIPGroupRules(rules),
		})
		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating WorkSpaces IP Group (%s): %s", d.Id(), err)
		}
	}

	return append(diags, resourceIPGroupRead(ctx, d, meta)...)
}

func resourceIPGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).WorkSpacesConn(ctx)

	var found bool
	log.Printf("[DEBUG] Finding directories associated with WorkSpaces IP Group (%s)", d.Id())
	err := conn.DescribeWorkspaceDirectoriesPagesWithContext(ctx, nil, func(page *workspaces.DescribeWorkspaceDirectoriesOutput, lastPage bool) bool {
		for _, dir := range page.Directories {
			for _, ipg := range dir.IpGroupIds {
				groupID := aws.StringValue(ipg)
				if groupID == d.Id() {
					found = true
					log.Printf("[DEBUG] WorkSpaces IP Group (%s) associated with WorkSpaces Directory (%s), disassociating", groupID, aws.StringValue(dir.DirectoryId))
					_, err := conn.DisassociateIpGroupsWithContext(ctx, &workspaces.DisassociateIpGroupsInput{
						DirectoryId: dir.DirectoryId,
						GroupIds:    aws.StringSlice([]string{d.Id()}),
					})
					if err != nil {
						diags = sdkdiag.AppendErrorf(diags, "disassociating WorkSpaces IP Group (%s) from WorkSpaces Directory (%s): %s", d.Id(), aws.StringValue(dir.DirectoryId), err)
						continue
					}
					log.Printf("[INFO] WorkSpaces IP Group (%s) disassociated from WorkSpaces Directory (%s)", d.Id(), aws.StringValue(dir.DirectoryId))
				}
			}
		}
		return !lastPage
	})
	if err != nil {
		diags = sdkdiag.AppendErrorf(diags, "describing WorkSpaces Directories: %s", err)
	}
	if diags.HasError() {
		return diags
	}

	if !found {
		log.Printf("[DEBUG] WorkSpaces IP Group (%s) not associated with any WorkSpaces Directories", d.Id())
	}

	log.Printf("[DEBUG] Deleting WorkSpaces IP Group (%s)", d.Id())
	_, err = conn.DeleteIpGroupWithContext(ctx, &workspaces.DeleteIpGroupInput{
		GroupId: aws.String(d.Id()),
	})
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting WorkSpaces IP Group (%s): %s", d.Id(), err)
	}
	log.Printf("[INFO] WorkSpaces IP Group (%s) deleted", d.Id())

	return diags
}

func expandIPGroupRules(rules []interface{}) []*workspaces.IpRuleItem {
	var result []*workspaces.IpRuleItem
	for _, rule := range rules {
		r := rule.(map[string]interface{})

		result = append(result, &workspaces.IpRuleItem{
			IpRule:   aws.String(r["source"].(string)),
			RuleDesc: aws.String(r["description"].(string)),
		})
	}
	return result
}

func flattenIPGroupRules(rules []*workspaces.IpRuleItem) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(rules))
	for _, rule := range rules {
		r := map[string]interface{}{}

		if v := rule.IpRule; v != nil {
			r["source"] = aws.StringValue(v)
		}

		if v := rule.RuleDesc; v != nil {
			r["description"] = aws.StringValue(rule.RuleDesc)
		}

		result = append(result, r)
	}
	return result
}
