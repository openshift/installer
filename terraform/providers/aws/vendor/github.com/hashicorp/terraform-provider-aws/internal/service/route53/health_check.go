package route53

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceHealthCheck() *schema.Resource {
	return &schema.Resource{
		Create: resourceHealthCheckCreate,
		Read:   resourceHealthCheckRead,
		Update: resourceHealthCheckUpdate,
		Delete: resourceHealthCheckDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				StateFunc: func(val interface{}) string {
					return strings.ToUpper(val.(string))
				},
				ValidateFunc: validation.StringInSlice(route53.HealthCheckType_Values(), true),
			},
			"failure_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 10),
			},
			"request_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true, // todo this should be updateable but the awslabs route53 service doesnt have the ability
				ValidateFunc: validation.IntInSlice([]int{10, 30}),
			},
			"ip_address": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPAddress,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return net.ParseIP(old).Equal(net.ParseIP(new))
				},
			},
			"fqdn": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IsPortNumber,
			},

			"invert_healthcheck": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"resource_path": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
			},

			"search_string": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
			},

			"measure_latency": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"child_healthchecks": {
				Type:     schema.TypeSet,
				MaxItems: 256,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringLenBetween(0, 64),
				},
				Optional: true,
			},
			"child_health_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtMost(256),
			},

			"cloudwatch_alarm_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"cloudwatch_alarm_region": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"insufficient_data_health_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(route53.InsufficientDataHealthStatus_Values(), true),
			},
			"reference_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				// The max length of the reference name is 64 characters for the API.
				// Terraform appends a 37-character unique ID to the provided
				// reference_name. This limits the length of the resource argument to 27.
				//
				// Example generated suffix: -terraform-20190122200019880700000001
				ValidateFunc: validation.StringLenBetween(0, (64 - resource.UniqueIDSuffixLength - 11)),
			},
			"enable_sni": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"regions": {
				Type:     schema.TypeSet,
				MinItems: 3,
				MaxItems: 64,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						endpoints.UsWest1RegionID,
						endpoints.UsWest2RegionID,
						endpoints.UsEast1RegionID,
						endpoints.EuWest1RegionID,
						endpoints.SaEast1RegionID,
						endpoints.ApSoutheast1RegionID,
						endpoints.ApSoutheast2RegionID,
						endpoints.ApNortheast1RegionID,
					}, true),
				},
				Optional: true,
			},

			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"routing_control_arn": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: verify.ValidARN,
			},

			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceHealthCheckCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).Route53Conn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	healthCheckType := d.Get("type").(string)
	healthConfig := &route53.HealthCheckConfig{
		Type: aws.String(healthCheckType),
	}

	if v, ok := d.GetOk("request_interval"); ok {
		healthConfig.RequestInterval = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("failure_threshold"); ok {
		healthConfig.FailureThreshold = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("fqdn"); ok {
		healthConfig.FullyQualifiedDomainName = aws.String(v.(string))
	}

	if v, ok := d.GetOk("search_string"); ok {
		healthConfig.SearchString = aws.String(v.(string))
	}

	if v, ok := d.GetOk("ip_address"); ok {
		healthConfig.IPAddress = aws.String(v.(string))
	}

	if v, ok := d.GetOk("port"); ok {
		healthConfig.Port = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("resource_path"); ok {
		healthConfig.ResourcePath = aws.String(v.(string))
	}

	if v, ok := d.GetOk("invert_healthcheck"); ok {
		healthConfig.Inverted = aws.Bool(v.(bool))
	}

	if v, ok := d.GetOk("enable_sni"); ok {
		healthConfig.EnableSNI = aws.Bool(v.(bool))
	}

	switch healthCheckType {
	case route53.HealthCheckTypeCalculated:
		if v, ok := d.GetOk("child_healthchecks"); ok {
			healthConfig.ChildHealthChecks = flex.ExpandStringSet(v.(*schema.Set))
		}

		if v, ok := d.GetOk("child_health_threshold"); ok {
			healthConfig.HealthThreshold = aws.Int64(int64(v.(int)))
		}
	case route53.HealthCheckTypeCloudwatchMetric:
		cloudwatchAlarmIdentifier := &route53.AlarmIdentifier{}

		if v, ok := d.GetOk("cloudwatch_alarm_name"); ok {
			cloudwatchAlarmIdentifier.Name = aws.String(v.(string))
		}

		if v, ok := d.GetOk("cloudwatch_alarm_region"); ok {
			cloudwatchAlarmIdentifier.Region = aws.String(v.(string))
		}

		healthConfig.AlarmIdentifier = cloudwatchAlarmIdentifier

		if v, ok := d.GetOk("insufficient_data_health_status"); ok {
			healthConfig.InsufficientDataHealthStatus = aws.String(v.(string))
		}
	case route53.HealthCheckTypeRecoveryControl:
		if v, ok := d.GetOk("routing_control_arn"); ok {
			healthConfig.RoutingControlArn = aws.String(v.(string))
		}
		fallthrough
	default:
		if v, ok := d.GetOk("measure_latency"); ok {
			healthConfig.MeasureLatency = aws.Bool(v.(bool))
		}
	}

	if v, ok := d.GetOk("regions"); ok {
		healthConfig.Regions = flex.ExpandStringSet(v.(*schema.Set))
	}

	callerRef := resource.UniqueId()
	if v, ok := d.GetOk("reference_name"); ok {
		callerRef = fmt.Sprintf("%s-%s", v.(string), callerRef)
	}

	if v, ok := d.GetOk("disabled"); ok {
		healthConfig.Disabled = aws.Bool(v.(bool))
	}

	input := &route53.CreateHealthCheckInput{
		CallerReference:   aws.String(callerRef),
		HealthCheckConfig: healthConfig,
	}

	resp, err := conn.CreateHealthCheck(input)

	if err != nil {
		return fmt.Errorf("creating Route53 Health Check (%s): %w", d.Id(), err)
	}

	d.SetId(aws.StringValue(resp.HealthCheck.Id))

	if err := UpdateTags(conn, d.Id(), route53.TagResourceTypeHealthcheck, nil, tags); err != nil {
		return fmt.Errorf("setting Route53 Health Check (%s) tags: %w", d.Id(), err)
	}

	return resourceHealthCheckRead(d, meta)
}

func resourceHealthCheckRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).Route53Conn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	output, err := FindHealthCheckByID(conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Route53 Health Check (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("reading Route53 Health Check (%s): %w", d.Id(), err)
	}

	healthCheckConfig := output.HealthCheckConfig
	d.Set("type", healthCheckConfig.Type)
	d.Set("failure_threshold", healthCheckConfig.FailureThreshold)
	d.Set("request_interval", healthCheckConfig.RequestInterval)
	d.Set("fqdn", healthCheckConfig.FullyQualifiedDomainName)
	d.Set("search_string", healthCheckConfig.SearchString)
	d.Set("ip_address", healthCheckConfig.IPAddress)
	d.Set("port", healthCheckConfig.Port)
	d.Set("resource_path", healthCheckConfig.ResourcePath)
	d.Set("measure_latency", healthCheckConfig.MeasureLatency)
	d.Set("invert_healthcheck", healthCheckConfig.Inverted)
	d.Set("disabled", healthCheckConfig.Disabled)
	d.Set("routing_control_arn", healthCheckConfig.RoutingControlArn)

	if err := d.Set("child_healthchecks", flex.FlattenStringList(healthCheckConfig.ChildHealthChecks)); err != nil {
		return fmt.Errorf("setting child_healthchecks: %w", err)
	}

	d.Set("child_health_threshold", healthCheckConfig.HealthThreshold)
	d.Set("insufficient_data_health_status", healthCheckConfig.InsufficientDataHealthStatus)
	d.Set("enable_sni", healthCheckConfig.EnableSNI)

	d.Set("regions", flex.FlattenStringList(healthCheckConfig.Regions))

	if healthCheckConfig.AlarmIdentifier != nil {
		d.Set("cloudwatch_alarm_name", healthCheckConfig.AlarmIdentifier.Name)
		d.Set("cloudwatch_alarm_region", healthCheckConfig.AlarmIdentifier.Region)
	}

	tags, err := ListTags(conn, d.Id(), route53.TagResourceTypeHealthcheck)

	if err != nil {
		return fmt.Errorf("listing tags for Route53 Health Check (%s): %w", d.Id(), err)
	}

	tags = tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("setting tags_all: %w", err)
	}

	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Service:   "route53",
		Resource:  fmt.Sprintf("healthcheck/%s", d.Id()),
	}.String()
	d.Set("arn", arn)

	return nil
}

func resourceHealthCheckUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).Route53Conn

	if d.HasChangesExcept("tags", "tags_all") {
		updateHealthCheck := &route53.UpdateHealthCheckInput{
			HealthCheckId: aws.String(d.Id()),
		}

		if d.HasChange("failure_threshold") {
			updateHealthCheck.FailureThreshold = aws.Int64(int64(d.Get("failure_threshold").(int)))
		}

		if d.HasChange("fqdn") {
			updateHealthCheck.FullyQualifiedDomainName = aws.String(d.Get("fqdn").(string))
		}

		if d.HasChange("port") {
			updateHealthCheck.Port = aws.Int64(int64(d.Get("port").(int)))
		}

		if d.HasChange("resource_path") {
			updateHealthCheck.ResourcePath = aws.String(d.Get("resource_path").(string))
		}

		if d.HasChange("invert_healthcheck") {
			updateHealthCheck.Inverted = aws.Bool(d.Get("invert_healthcheck").(bool))
		}

		if d.HasChange("child_healthchecks") {
			updateHealthCheck.ChildHealthChecks = flex.ExpandStringSet(d.Get("child_healthchecks").(*schema.Set))
		}

		if d.HasChange("child_health_threshold") {
			updateHealthCheck.HealthThreshold = aws.Int64(int64(d.Get("child_health_threshold").(int)))
		}

		if d.HasChange("search_string") {
			updateHealthCheck.SearchString = aws.String(d.Get("search_string").(string))
		}

		if d.HasChanges("cloudwatch_alarm_name", "cloudwatch_alarm_region") {
			cloudwatchAlarm := &route53.AlarmIdentifier{
				Name:   aws.String(d.Get("cloudwatch_alarm_name").(string)),
				Region: aws.String(d.Get("cloudwatch_alarm_region").(string)),
			}

			updateHealthCheck.AlarmIdentifier = cloudwatchAlarm
		}

		if d.HasChange("insufficient_data_health_status") {
			updateHealthCheck.InsufficientDataHealthStatus = aws.String(d.Get("insufficient_data_health_status").(string))
		}

		if d.HasChange("enable_sni") {
			updateHealthCheck.EnableSNI = aws.Bool(d.Get("enable_sni").(bool))
		}

		if d.HasChange("regions") {
			updateHealthCheck.Regions = flex.ExpandStringSet(d.Get("regions").(*schema.Set))
		}

		if d.HasChange("disabled") {
			updateHealthCheck.Disabled = aws.Bool(d.Get("disabled").(bool))
		}

		if d.HasChange("ip_address") {
			updateHealthCheck.IPAddress = aws.String(d.Get("ip_address").(string))
		}

		_, err := conn.UpdateHealthCheck(updateHealthCheck)

		if err != nil {
			return fmt.Errorf("updating Route53 Health Check (%s): %w", d.Id(), err)
		}
	}

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := UpdateTags(conn, d.Id(), route53.TagResourceTypeHealthcheck, o, n); err != nil {
			return fmt.Errorf("updating Route53 Health Check (%s) tags: %w", d.Id(), err)
		}
	}

	return resourceHealthCheckRead(d, meta)
}

func resourceHealthCheckDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).Route53Conn

	log.Printf("[DEBUG] Deleting Route53 Health Check: %s", d.Id())
	_, err := conn.DeleteHealthCheck(&route53.DeleteHealthCheckInput{HealthCheckId: aws.String(d.Id())})

	if tfawserr.ErrCodeEquals(err, route53.ErrCodeNoSuchHealthCheck) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("deleting Route53 Health Check (%s): %w", d.Id(), err)
	}

	return nil
}
