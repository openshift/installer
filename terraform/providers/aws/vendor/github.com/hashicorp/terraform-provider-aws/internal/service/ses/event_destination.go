package ses

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceEventDestination() *schema.Resource {
	return &schema.Resource{
		Create: resourceEventDestinationCreate,
		Read:   resourceEventDestinationRead,
		Delete: resourceEventDestinationDelete,
		Importer: &schema.ResourceImporter{
			State: resourceEventDestinationImport,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloudwatch_destination": {
				Type:          schema.TypeSet,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"kinesis_destination", "sns_destination"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default_value": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.All(
								validation.StringLenBetween(1, 256),
								validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z_\-\.@]+$`), "must contain only alphanumeric, underscore, hyphen, period, and at signs characters"),
							),
						},
						"dimension_name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.All(
								validation.StringLenBetween(1, 256),
								validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z_:-]+$`), "must contain only alphanumeric, underscore, and hyphen characters"),
							),
						},
						"value_source": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(ses.DimensionValueSource_Values(), false),
						},
					},
				},
			},
			"configuration_set_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"kinesis_destination": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				MaxItems:      1,
				ConflictsWith: []string{"cloudwatch_destination", "sns_destination"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_arn": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: verify.ValidARN,
						},
						"stream_arn": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: verify.ValidARN,
						},
					},
				},
			},
			"matching_types": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Set:      schema.HashString,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(ses.EventType_Values(), false),
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 64),
					validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z_-]+$`), "must contain only alphanumeric, underscore, and hyphen characters"),
				),
			},
			"sns_destination": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"cloudwatch_destination", "kinesis_destination"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic_arn": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: verify.ValidARN,
						},
					},
				},
			},
		},
	}
}

func resourceEventDestinationCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SESConn

	configurationSetName := d.Get("configuration_set_name").(string)
	eventDestinationName := d.Get("name").(string)
	enabled := d.Get("enabled").(bool)
	matchingEventTypes := d.Get("matching_types").(*schema.Set)

	createOpts := &ses.CreateConfigurationSetEventDestinationInput{
		ConfigurationSetName: aws.String(configurationSetName),
		EventDestination: &ses.EventDestination{
			Name:               aws.String(eventDestinationName),
			Enabled:            aws.Bool(enabled),
			MatchingEventTypes: flex.ExpandStringSet(matchingEventTypes),
		},
	}

	if v, ok := d.GetOk("cloudwatch_destination"); ok {
		destination := v.(*schema.Set).List()
		createOpts.EventDestination.CloudWatchDestination = &ses.CloudWatchDestination{
			DimensionConfigurations: generateCloudWatchDestination(destination),
		}
		log.Printf("[DEBUG] Creating cloudwatch destination: %#v", destination)
	}

	if v, ok := d.GetOk("kinesis_destination"); ok {
		destination := v.([]interface{})

		kinesis := destination[0].(map[string]interface{})
		createOpts.EventDestination.KinesisFirehoseDestination = &ses.KinesisFirehoseDestination{
			DeliveryStreamARN: aws.String(kinesis["stream_arn"].(string)),
			IAMRoleARN:        aws.String(kinesis["role_arn"].(string)),
		}
		log.Printf("[DEBUG] Creating kinesis destination: %#v", kinesis)
	}

	if v, ok := d.GetOk("sns_destination"); ok {
		destination := v.([]interface{})
		sns := destination[0].(map[string]interface{})
		createOpts.EventDestination.SNSDestination = &ses.SNSDestination{
			TopicARN: aws.String(sns["topic_arn"].(string)),
		}
		log.Printf("[DEBUG] Creating sns destination: %#v", sns)
	}

	_, err := conn.CreateConfigurationSetEventDestination(createOpts)
	if err != nil {
		return fmt.Errorf("Error creating SES configuration set event destination: %s", err)
	}

	d.SetId(eventDestinationName)

	log.Printf("[WARN] SES DONE")
	return resourceEventDestinationRead(d, meta)
}

func resourceEventDestinationRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SESConn

	configurationSetName := d.Get("configuration_set_name").(string)
	input := &ses.DescribeConfigurationSetInput{
		ConfigurationSetAttributeNames: aws.StringSlice([]string{ses.ConfigurationSetAttributeEventDestinations}),
		ConfigurationSetName:           aws.String(configurationSetName),
	}

	output, err := conn.DescribeConfigurationSet(input)
	if tfawserr.ErrCodeEquals(err, ses.ErrCodeConfigurationSetDoesNotExistException) {
		log.Printf("[WARN] SES Configuration Set (%s) not found, removing from state", configurationSetName)
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("reading SES Configuration Set Event Destination (%s): %w", d.Id(), err)
	}

	var thisEventDestination *ses.EventDestination
	for _, eventDestination := range output.EventDestinations {
		if aws.StringValue(eventDestination.Name) == d.Id() {
			thisEventDestination = eventDestination
			break
		}
	}
	if thisEventDestination == nil {
		log.Printf("[WARN] SES Configuration Set Event Destination (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("configuration_set_name", output.ConfigurationSet.Name)
	d.Set("enabled", thisEventDestination.Enabled)
	d.Set("name", thisEventDestination.Name)
	if err := d.Set("cloudwatch_destination", flattenCloudWatchDestination(thisEventDestination.CloudWatchDestination)); err != nil {
		return fmt.Errorf("setting cloudwatch_destination: %w", err)
	}
	if err := d.Set("kinesis_destination", flattenKinesisFirehoseDestination(thisEventDestination.KinesisFirehoseDestination)); err != nil {
		return fmt.Errorf("setting kinesis_destination: %w", err)
	}
	if err := d.Set("matching_types", flex.FlattenStringSet(thisEventDestination.MatchingEventTypes)); err != nil {
		return fmt.Errorf("setting matching_types: %w", err)
	}
	if err := d.Set("sns_destination", flattenSNSDestination(thisEventDestination.SNSDestination)); err != nil {
		return fmt.Errorf("setting sns_destination: %w", err)
	}

	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Service:   "ses",
		Region:    meta.(*conns.AWSClient).Region,
		AccountID: meta.(*conns.AWSClient).AccountID,
		Resource:  fmt.Sprintf("configuration-set/%s:event-destination/%s", configurationSetName, d.Id()),
	}.String()
	d.Set("arn", arn)

	return nil
}

func resourceEventDestinationDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SESConn

	log.Printf("[DEBUG] SES Delete Configuration Set Destination: %s", d.Id())
	_, err := conn.DeleteConfigurationSetEventDestination(&ses.DeleteConfigurationSetEventDestinationInput{
		ConfigurationSetName: aws.String(d.Get("configuration_set_name").(string)),
		EventDestinationName: aws.String(d.Id()),
	})

	return err
}

func resourceEventDestinationImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return []*schema.ResourceData{}, fmt.Errorf("wrong format of import ID (%s), use: 'configuration-set-name/event-destination-name'", d.Id())
	}

	configurationSetName := parts[0]
	eventDestinationName := parts[1]
	log.Printf("[DEBUG] Importing SES event destination %s from configuration set %s", eventDestinationName, configurationSetName)

	d.SetId(eventDestinationName)
	d.Set("configuration_set_name", configurationSetName)

	return []*schema.ResourceData{d}, nil
}

func generateCloudWatchDestination(v []interface{}) []*ses.CloudWatchDimensionConfiguration {

	b := make([]*ses.CloudWatchDimensionConfiguration, len(v))

	for i, vI := range v {
		cloudwatch := vI.(map[string]interface{})
		b[i] = &ses.CloudWatchDimensionConfiguration{
			DefaultDimensionValue: aws.String(cloudwatch["default_value"].(string)),
			DimensionName:         aws.String(cloudwatch["dimension_name"].(string)),
			DimensionValueSource:  aws.String(cloudwatch["value_source"].(string)),
		}
	}

	return b
}

func flattenCloudWatchDestination(destination *ses.CloudWatchDestination) []interface{} {
	if destination == nil {
		return []interface{}{}
	}

	vDimensionConfigurations := []interface{}{}

	for _, dimensionConfiguration := range destination.DimensionConfigurations {
		mDimensionConfiguration := map[string]interface{}{
			"default_value":  aws.StringValue(dimensionConfiguration.DefaultDimensionValue),
			"dimension_name": aws.StringValue(dimensionConfiguration.DimensionName),
			"value_source":   aws.StringValue(dimensionConfiguration.DimensionValueSource),
		}

		vDimensionConfigurations = append(vDimensionConfigurations, mDimensionConfiguration)
	}

	return vDimensionConfigurations
}

func flattenKinesisFirehoseDestination(destination *ses.KinesisFirehoseDestination) []interface{} {
	if destination == nil {
		return []interface{}{}
	}

	mDestination := map[string]interface{}{
		"role_arn":   aws.StringValue(destination.IAMRoleARN),
		"stream_arn": aws.StringValue(destination.DeliveryStreamARN),
	}

	return []interface{}{mDestination}
}

func flattenSNSDestination(destination *ses.SNSDestination) []interface{} {
	if destination == nil {
		return []interface{}{}
	}

	mDestination := map[string]interface{}{
		"topic_arn": aws.StringValue(destination.TopicARN),
	}

	return []interface{}{mDestination}
}
