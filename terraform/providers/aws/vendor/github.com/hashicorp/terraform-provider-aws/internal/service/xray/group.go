package xray

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/xray"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Read:   resourceGroupRead,
		Update: resourceGroupUpdate,
		Delete: resourceGroupDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		CustomizeDiff: verify.SetTagsDiff,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"filter_expression": {
				Type:     schema.TypeString,
				Required: true,
			},
			"insights_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"insights_enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"notifications_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
		},
	}
}

func resourceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).XRayConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	name := d.Get("group_name").(string)
	input := &xray.CreateGroupInput{
		GroupName:        aws.String(name),
		FilterExpression: aws.String(d.Get("filter_expression").(string)),
		Tags:             Tags(tags.IgnoreAWS()),
	}

	if v, ok := d.GetOk("insights_configuration"); ok {
		input.InsightsConfiguration = expandInsightsConfig(v.([]interface{}))
	}

	output, err := conn.CreateGroup(input)

	if err != nil {
		return fmt.Errorf("error creating XRay Group (%s): %w", name, err)
	}

	d.SetId(aws.StringValue(output.Group.GroupARN))

	return resourceGroupRead(d, meta)
}

func resourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).XRayConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	input := &xray.GetGroupInput{
		GroupARN: aws.String(d.Id()),
	}

	group, err := conn.GetGroup(input)

	if tfawserr.ErrMessageContains(err, xray.ErrCodeInvalidRequestException, "Group not found") {
		log.Printf("[WARN] XRay Group (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading XRay Group (%s): %w", d.Id(), err)
	}

	arn := aws.StringValue(group.Group.GroupARN)
	d.Set("arn", arn)
	d.Set("group_name", group.Group.GroupName)
	d.Set("filter_expression", group.Group.FilterExpression)
	if err := d.Set("insights_configuration", flattenInsightsConfig(group.Group.InsightsConfiguration)); err != nil {
		return fmt.Errorf("error setting insights_configuration: %w", err)
	}

	tags, err := ListTags(conn, arn)

	if err != nil {
		return fmt.Errorf("error listing tags for Xray Group (%q): %s", d.Id(), err)
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

func resourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).XRayConn

	if d.HasChangesExcept("tags", "tags_all") {
		input := &xray.UpdateGroupInput{GroupARN: aws.String(d.Id())}

		if v, ok := d.GetOk("filter_expression"); ok {
			input.FilterExpression = aws.String(v.(string))
		}

		if v, ok := d.GetOk("insights_configuration"); ok {
			input.InsightsConfiguration = expandInsightsConfig(v.([]interface{}))
		}

		_, err := conn.UpdateGroup(input)

		if err != nil {
			return fmt.Errorf("error updating XRay Group (%s): %w", d.Id(), err)
		}
	}

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := UpdateTags(conn, d.Get("arn").(string), o, n); err != nil {
			return fmt.Errorf("error updating tags: %w", err)
		}
	}

	return resourceGroupRead(d, meta)
}

func resourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).XRayConn

	log.Printf("[INFO] Deleting XRay Group: %s", d.Id())
	_, err := conn.DeleteGroup(&xray.DeleteGroupInput{
		GroupARN: aws.String(d.Id()),
	})

	if err != nil {
		return fmt.Errorf("error deleting XRay Group (%s): %w", d.Id(), err)
	}

	return nil
}

func expandInsightsConfig(l []interface{}) *xray.InsightsConfiguration {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})
	config := xray.InsightsConfiguration{}

	if v, ok := m["insights_enabled"]; ok {
		config.InsightsEnabled = aws.Bool(v.(bool))
	}
	if v, ok := m["notifications_enabled"]; ok {
		config.NotificationsEnabled = aws.Bool(v.(bool))
	}

	return &config
}

func flattenInsightsConfig(config *xray.InsightsConfiguration) []interface{} {
	if config == nil {
		return nil
	}

	m := map[string]interface{}{}

	if config.InsightsEnabled != nil {
		m["insights_enabled"] = config.InsightsEnabled
	}
	if config.NotificationsEnabled != nil {
		m["notifications_enabled"] = config.NotificationsEnabled
	}

	return []interface{}{m}
}
