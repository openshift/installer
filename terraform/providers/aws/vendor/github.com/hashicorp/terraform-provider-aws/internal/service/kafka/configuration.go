package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kafka"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
)

func ResourceConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceConfigurationCreate,
		Read:   resourceConfigurationRead,
		Update: resourceConfigurationUpdate,
		Delete: resourceConfigurationDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		CustomizeDiff: customdiff.Sequence(
			customdiff.ComputedIf("latest_revision", func(_ context.Context, diff *schema.ResourceDiff, meta interface{}) bool {
				return diff.HasChange("server_properties")
			}),
		),

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"kafka_versions": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"latest_revision": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"server_properties": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).KafkaConn

	input := &kafka.CreateConfigurationInput{
		Name:             aws.String(d.Get("name").(string)),
		ServerProperties: []byte(d.Get("server_properties").(string)),
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = aws.String(v.(string))
	}

	if v, ok := d.GetOk("kafka_versions"); ok && v.(*schema.Set).Len() > 0 {
		input.KafkaVersions = flex.ExpandStringSet(v.(*schema.Set))
	}

	output, err := conn.CreateConfiguration(input)

	if err != nil {
		return fmt.Errorf("error creating MSK Configuration: %s", err)
	}

	d.SetId(aws.StringValue(output.Arn))

	return resourceConfigurationRead(d, meta)
}

func resourceConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).KafkaConn

	configurationInput := &kafka.DescribeConfigurationInput{
		Arn: aws.String(d.Id()),
	}

	configurationOutput, err := conn.DescribeConfiguration(configurationInput)

	if tfawserr.ErrMessageContains(err, kafka.ErrCodeBadRequestException, "Configuration ARN does not exist") {
		log.Printf("[WARN] MSK Configuration (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error describing MSK Configuration (%s): %s", d.Id(), err)
	}

	if configurationOutput == nil {
		return fmt.Errorf("error describing MSK Configuration (%s): missing result", d.Id())
	}

	if configurationOutput.LatestRevision == nil {
		return fmt.Errorf("error describing MSK Configuration (%s): missing latest revision", d.Id())
	}

	revision := configurationOutput.LatestRevision.Revision
	revisionInput := &kafka.DescribeConfigurationRevisionInput{
		Arn:      aws.String(d.Id()),
		Revision: revision,
	}

	revisionOutput, err := conn.DescribeConfigurationRevision(revisionInput)

	if err != nil {
		return fmt.Errorf("error describing MSK Configuration (%s) Revision (%d): %s", d.Id(), aws.Int64Value(revision), err)
	}

	if revisionOutput == nil {
		return fmt.Errorf("error describing MSK Configuration (%s) Revision (%d): missing result", d.Id(), aws.Int64Value(revision))
	}

	d.Set("arn", configurationOutput.Arn)
	d.Set("description", revisionOutput.Description)

	if err := d.Set("kafka_versions", aws.StringValueSlice(configurationOutput.KafkaVersions)); err != nil {
		return fmt.Errorf("error setting kafka_versions: %s", err)
	}

	d.Set("latest_revision", revision)
	d.Set("name", configurationOutput.Name)
	d.Set("server_properties", string(revisionOutput.ServerProperties))

	return nil
}

func resourceConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).KafkaConn

	input := &kafka.UpdateConfigurationInput{
		Arn:              aws.String(d.Id()),
		ServerProperties: []byte(d.Get("server_properties").(string)),
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = aws.String(v.(string))
	}

	_, err := conn.UpdateConfiguration(input)

	if err != nil {
		return fmt.Errorf("error updating MSK Configuration (%s): %w", d.Id(), err)
	}

	return resourceConfigurationRead(d, meta)
}

func resourceConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).KafkaConn

	input := &kafka.DeleteConfigurationInput{
		Arn: aws.String(d.Id()),
	}

	_, err := conn.DeleteConfiguration(input)

	if err != nil {
		return fmt.Errorf("error deleting MSK Configuration (%s): %w", d.Id(), err)
	}

	if _, err := waitConfigurationDeleted(conn, d.Id()); err != nil {
		return fmt.Errorf("error waiting for MSK Configuration (%s): %w", d.Id(), err)
	}

	return nil
}
