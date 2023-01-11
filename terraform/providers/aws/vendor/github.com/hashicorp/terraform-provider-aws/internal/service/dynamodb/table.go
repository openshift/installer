package dynamodb

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

const (
	provisionedThroughputMinValue = 1
	ResNameTable                  = "Table"
)

func ResourceTable() *schema.Resource {
	//lintignore:R011
	return &schema.Resource{
		Create: resourceTableCreate,
		Read:   resourceTableRead,
		Update: resourceTableUpdate,
		Delete: resourceTableDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(createTableTimeout),
			Delete: schema.DefaultTimeout(deleteTableTimeout),
			Update: schema.DefaultTimeout(updateTableTimeoutTotal),
		},

		CustomizeDiff: customdiff.All(
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return validStreamSpec(diff)
			},
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return validateTableAttributes(diff)
			},
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				if diff.Id() != "" && diff.HasChange("server_side_encryption") {
					o, n := diff.GetChange("server_side_encryption")
					if isTableOptionDisabled(o) && isTableOptionDisabled(n) {
						return diff.Clear("server_side_encryption")
					}
				}
				return nil
			},
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				if diff.Id() != "" && diff.HasChange("point_in_time_recovery") {
					o, n := diff.GetChange("point_in_time_recovery")
					if isTableOptionDisabled(o) && isTableOptionDisabled(n) {
						return diff.Clear("point_in_time_recovery")
					}
				}
				return nil
			},
			func(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
				if v := diff.Get("restore_source_name"); v != "" {
					return nil
				}

				var errs *multierror.Error
				if err := validateProvisionedThroughputField(diff, "read_capacity"); err != nil {
					errs = multierror.Append(errs, err)
				}
				if err := validateProvisionedThroughputField(diff, "write_capacity"); err != nil {
					errs = multierror.Append(errs, err)
				}
				return errs.ErrorOrNil()
			},
			customdiff.ForceNewIfChange("restore_source_name", func(_ context.Context, old, new, meta interface{}) bool {
				// If they differ force new unless new is cleared
				// https://github.com/hashicorp/terraform-provider-aws/issues/25214
				return old.(string) != new.(string) && new.(string) != ""
			}),
			verify.SetTagsDiff,
		),

		SchemaVersion: 1,
		MigrateState:  resourceTableMigrateState,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attribute": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								dynamodb.ScalarAttributeTypeB,
								dynamodb.ScalarAttributeTypeN,
								dynamodb.ScalarAttributeTypeS,
							}, false),
						},
					},
				},
				Set: func(v interface{}) int {
					var buf bytes.Buffer
					m := v.(map[string]interface{})
					buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
					return create.StringHashcode(buf.String())
				},
			},
			"billing_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      dynamodb.BillingModeProvisioned,
				ValidateFunc: validation.StringInSlice(dynamodb.BillingMode_Values(), false),
			},
			"global_secondary_index": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hash_key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"non_key_attributes": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"projection_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(dynamodb.ProjectionType_Values(), false),
						},
						"range_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"read_capacity": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"write_capacity": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"hash_key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"local_secondary_index": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"non_key_attributes": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"projection_type": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice(dynamodb.ProjectionType_Values(), false),
						},
						"range_key": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
				Set: func(v interface{}) int {
					var buf bytes.Buffer
					m := v.(map[string]interface{})
					buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
					return create.StringHashcode(buf.String())
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"point_in_time_recovery": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},
			"range_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"read_capacity": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"replica": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kms_key_arn": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: verify.ValidARN,
						},
						"point_in_time_recovery": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"propagate_tags": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"region_name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"restore_date_time": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidUTCTimestamp,
			},
			"restore_source_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"restore_to_latest_time": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"server_side_encryption": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"kms_key_arn": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: verify.ValidARN,
						},
					},
				},
			},
			"stream_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"stream_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"stream_label": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"stream_view_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				StateFunc: func(v interface{}) string {
					value := v.(string)
					return strings.ToUpper(value)
				},
				ValidateFunc: validation.StringInSlice([]string{
					"",
					dynamodb.StreamViewTypeNewImage,
					dynamodb.StreamViewTypeOldImage,
					dynamodb.StreamViewTypeNewAndOldImages,
					dynamodb.StreamViewTypeKeysOnly,
				}, false),
			},
			"table_class": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(dynamodb.TableClass_Values(), false),
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
			"ttl": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attribute_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
				DiffSuppressFunc: verify.SuppressMissingOptionalConfigurationBlock,
			},
			"write_capacity": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceTableCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).DynamoDBConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	keySchemaMap := map[string]interface{}{
		"hash_key": d.Get("hash_key").(string),
	}
	if v, ok := d.GetOk("range_key"); ok {
		keySchemaMap["range_key"] = v.(string)
	}

	log.Printf("[DEBUG] Creating DynamoDB table with key schema: %#v", keySchemaMap)

	if _, ok := d.GetOk("restore_source_name"); ok {
		input := &dynamodb.RestoreTableToPointInTimeInput{
			TargetTableName: aws.String(d.Get("name").(string)),
			SourceTableName: aws.String(d.Get("restore_source_name").(string)),
		}

		if v, ok := d.GetOk("restore_date_time"); ok {
			t, _ := time.Parse(time.RFC3339, v.(string))
			input.RestoreDateTime = aws.Time(t)
		}

		if attr, ok := d.GetOk("restore_to_latest_time"); ok {
			input.UseLatestRestorableTime = aws.Bool(attr.(bool))
		}

		if v, ok := d.GetOk("local_secondary_index"); ok {
			lsiSet := v.(*schema.Set)
			input.LocalSecondaryIndexOverride = expandLocalSecondaryIndexes(lsiSet.List(), keySchemaMap)
		}

		billingModeOverride := d.Get("billing_mode").(string)

		if _, ok := d.GetOk("write_capacity"); ok {
			if _, ok := d.GetOk("read_capacity"); ok {
				capacityMap := map[string]interface{}{
					"write_capacity": d.Get("write_capacity"),
					"read_capacity":  d.Get("read_capacity"),
				}
				input.ProvisionedThroughputOverride = expandProvisionedThroughput(capacityMap, billingModeOverride)
			}
		}

		if v, ok := d.GetOk("local_secondary_index"); ok {
			lsiSet := v.(*schema.Set)
			input.LocalSecondaryIndexOverride = expandLocalSecondaryIndexes(lsiSet.List(), keySchemaMap)
		}

		if v, ok := d.GetOk("global_secondary_index"); ok {
			globalSecondaryIndexes := []*dynamodb.GlobalSecondaryIndex{}
			gsiSet := v.(*schema.Set)

			for _, gsiObject := range gsiSet.List() {
				gsi := gsiObject.(map[string]interface{})
				if err := validateGSIProvisionedThroughput(gsi, billingModeOverride); err != nil {
					return create.Error(names.DynamoDB, create.ErrActionCreating, ResNameTable, d.Get("name").(string), err)
				}

				gsiObject := expandGlobalSecondaryIndex(gsi, billingModeOverride)
				globalSecondaryIndexes = append(globalSecondaryIndexes, gsiObject)
			}
			input.GlobalSecondaryIndexOverride = globalSecondaryIndexes
		}

		if v, ok := d.GetOk("server_side_encryption"); ok {
			input.SSESpecificationOverride = expandEncryptAtRestOptions(v.([]interface{}))
		}

		var output *dynamodb.RestoreTableToPointInTimeOutput
		err := resource.Retry(createTableTimeout, func() *resource.RetryError {
			var err error
			output, err = conn.RestoreTableToPointInTime(input)
			if err != nil {
				if tfawserr.ErrCodeEquals(err, "ThrottlingException") {
					return resource.RetryableError(err)
				}
				if tfawserr.ErrMessageContains(err, dynamodb.ErrCodeLimitExceededException, "can be created, updated, or deleted simultaneously") {
					return resource.RetryableError(err)
				}
				if tfawserr.ErrMessageContains(err, dynamodb.ErrCodeLimitExceededException, "indexed tables that can be created simultaneously") {
					return resource.RetryableError(err)
				}

				return resource.NonRetryableError(err)
			}
			return nil
		})

		if tfresource.TimedOut(err) {
			output, err = conn.RestoreTableToPointInTime(input)
		}

		if err != nil {
			return create.Error(names.DynamoDB, create.ErrActionCreating, ResNameTable, d.Get("name").(string), err)
		}

		if output == nil || output.TableDescription == nil {
			return errors.New("error creating DynamoDB Table: empty response")
		}

	} else {
		input := &dynamodb.CreateTableInput{
			TableName:   aws.String(d.Get("name").(string)),
			BillingMode: aws.String(d.Get("billing_mode").(string)),
			KeySchema:   expandKeySchema(keySchemaMap),
		}

		if len(tags) > 0 {
			input.Tags = Tags(tags.IgnoreAWS())
		}

		billingMode := d.Get("billing_mode").(string)

		capacityMap := map[string]interface{}{
			"write_capacity": d.Get("write_capacity"),
			"read_capacity":  d.Get("read_capacity"),
		}

		input.ProvisionedThroughput = expandProvisionedThroughput(capacityMap, billingMode)

		if v, ok := d.GetOk("attribute"); ok {
			aSet := v.(*schema.Set)
			input.AttributeDefinitions = expandAttributes(aSet.List())
		}

		if v, ok := d.GetOk("local_secondary_index"); ok {
			lsiSet := v.(*schema.Set)
			input.LocalSecondaryIndexes = expandLocalSecondaryIndexes(lsiSet.List(), keySchemaMap)
		}

		if v, ok := d.GetOk("global_secondary_index"); ok {
			globalSecondaryIndexes := []*dynamodb.GlobalSecondaryIndex{}
			gsiSet := v.(*schema.Set)

			for _, gsiObject := range gsiSet.List() {
				gsi := gsiObject.(map[string]interface{})
				if err := validateGSIProvisionedThroughput(gsi, billingMode); err != nil {
					return create.Error(names.DynamoDB, create.ErrActionCreating, ResNameTable, d.Get("name").(string), err)
				}

				gsiObject := expandGlobalSecondaryIndex(gsi, billingMode)
				globalSecondaryIndexes = append(globalSecondaryIndexes, gsiObject)
			}
			input.GlobalSecondaryIndexes = globalSecondaryIndexes
		}

		if v, ok := d.GetOk("stream_enabled"); ok {
			input.StreamSpecification = &dynamodb.StreamSpecification{
				StreamEnabled:  aws.Bool(v.(bool)),
				StreamViewType: aws.String(d.Get("stream_view_type").(string)),
			}
		}

		if v, ok := d.GetOk("server_side_encryption"); ok {
			input.SSESpecification = expandEncryptAtRestOptions(v.([]interface{}))
		}

		if v, ok := d.GetOk("table_class"); ok {
			input.TableClass = aws.String(v.(string))
		}

		var output *dynamodb.CreateTableOutput
		err := resource.Retry(createTableTimeout, func() *resource.RetryError {
			var err error
			output, err = conn.CreateTable(input)
			if err != nil {
				if tfawserr.ErrCodeEquals(err, "ThrottlingException", "") {
					return resource.RetryableError(err)
				}
				if tfawserr.ErrMessageContains(err, dynamodb.ErrCodeLimitExceededException, "can be created, updated, or deleted simultaneously") {
					return resource.RetryableError(err)
				}
				if tfawserr.ErrMessageContains(err, dynamodb.ErrCodeLimitExceededException, "indexed tables that can be created simultaneously") {
					return resource.RetryableError(err)
				}

				return resource.NonRetryableError(err)
			}
			return nil
		})

		if tfresource.TimedOut(err) {
			output, err = conn.CreateTable(input)
		}

		if err != nil {
			return create.Error(names.DynamoDB, create.ErrActionCreating, ResNameTable, d.Get("name").(string), err)
		}

		if output == nil || output.TableDescription == nil {
			return errors.New("error creating DynamoDB Table: empty response")
		}
	}

	d.SetId(d.Get("name").(string))

	var output *dynamodb.TableDescription
	var err error
	if output, err = waitTableActive(conn, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
		return create.Error(names.DynamoDB, create.ErrActionWaitingForCreation, ResNameTable, d.Id(), err)
	}

	if v, ok := d.GetOk("global_secondary_index"); ok {
		gsiSet := v.(*schema.Set)

		for _, gsiObject := range gsiSet.List() {
			gsi := gsiObject.(map[string]interface{})

			if _, err := waitGSIActive(conn, d.Id(), gsi["name"].(string), d.Timeout(schema.TimeoutUpdate)); err != nil {
				return create.Error(names.DynamoDB, create.ErrActionWaitingForCreation, ResNameTable, d.Id(), fmt.Errorf("GSI (%s): %w", gsi["name"].(string), err))
			}
		}
	}

	if d.Get("ttl.0.enabled").(bool) {
		if err := updateTimeToLive(conn, d.Id(), d.Get("ttl").([]interface{}), d.Timeout(schema.TimeoutCreate)); err != nil {
			return create.Error(names.DynamoDB, create.ErrActionCreating, ResNameTable, d.Id(), fmt.Errorf("enabling TTL: %w", err))
		}
	}

	if d.Get("point_in_time_recovery.0.enabled").(bool) {
		if err := updatePITR(conn, d.Id(), true, aws.StringValue(conn.Config.Region), meta.(*conns.AWSClient).TerraformVersion, d.Timeout(schema.TimeoutCreate)); err != nil {
			return create.Error(names.DynamoDB, create.ErrActionCreating, ResNameTable, d.Id(), fmt.Errorf("enabling point in time recovery: %w", err))
		}
	}

	if v := d.Get("replica").(*schema.Set); v.Len() > 0 {
		if err := createReplicas(conn, d.Id(), v.List(), meta.(*conns.AWSClient).TerraformVersion, true, d.Timeout(schema.TimeoutCreate)); err != nil {
			return create.Error(names.DynamoDB, create.ErrActionCreating, ResNameTable, d.Id(), fmt.Errorf("replicas: %w", err))
		}

		if err := updateReplicaTags(conn, aws.StringValue(output.TableArn), v.List(), tags, meta.(*conns.AWSClient).TerraformVersion); err != nil {
			return create.Error(names.DynamoDB, create.ErrActionCreating, ResNameTable, d.Id(), fmt.Errorf("replica tags: %w", err))
		}
	}

	return resourceTableRead(d, meta)
}

func resourceTableRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).DynamoDBConn

	result, err := conn.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(d.Id()),
	})

	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, dynamodb.ErrCodeResourceNotFoundException) {
		log.Printf("[WARN] Dynamodb Table (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return create.Error(names.DynamoDB, create.ErrActionReading, ResNameTable, d.Id(), err)
	}

	if result == nil || result.Table == nil {
		if d.IsNewResource() {
			return create.Error(names.DynamoDB, create.ErrActionReading, ResNameTable, d.Id(), errors.New("empty output after creation"))
		}
		create.LogNotFoundRemoveState(names.DynamoDB, create.ErrActionReading, ResNameTable, d.Id())
		d.SetId("")
		return nil
	}

	table := result.Table

	d.Set("arn", table.TableArn)
	d.Set("name", table.TableName)

	if table.BillingModeSummary != nil {
		d.Set("billing_mode", table.BillingModeSummary.BillingMode)
	} else {
		d.Set("billing_mode", dynamodb.BillingModeProvisioned)
	}

	if table.ProvisionedThroughput != nil {
		d.Set("write_capacity", table.ProvisionedThroughput.WriteCapacityUnits)
		d.Set("read_capacity", table.ProvisionedThroughput.ReadCapacityUnits)
	}

	if err := d.Set("attribute", flattenTableAttributeDefinitions(table.AttributeDefinitions)); err != nil {
		return create.SettingError(names.DynamoDB, ResNameTable, d.Id(), "attribute", err)
	}

	for _, attribute := range table.KeySchema {
		if aws.StringValue(attribute.KeyType) == dynamodb.KeyTypeHash {
			d.Set("hash_key", attribute.AttributeName)
		}

		if aws.StringValue(attribute.KeyType) == dynamodb.KeyTypeRange {
			d.Set("range_key", attribute.AttributeName)
		}
	}

	if err := d.Set("local_secondary_index", flattenTableLocalSecondaryIndex(table.LocalSecondaryIndexes)); err != nil {
		return create.SettingError(names.DynamoDB, ResNameTable, d.Id(), "local_secondary_index", err)
	}

	if err := d.Set("global_secondary_index", flattenTableGlobalSecondaryIndex(table.GlobalSecondaryIndexes)); err != nil {
		return create.SettingError(names.DynamoDB, ResNameTable, d.Id(), "global_secondary_index", err)
	}

	if table.StreamSpecification != nil {
		d.Set("stream_view_type", table.StreamSpecification.StreamViewType)
		d.Set("stream_enabled", table.StreamSpecification.StreamEnabled)
	} else {
		d.Set("stream_view_type", "")
		d.Set("stream_enabled", false)
	}

	d.Set("stream_arn", table.LatestStreamArn)
	d.Set("stream_label", table.LatestStreamLabel)

	if err := d.Set("server_side_encryption", flattenTableServerSideEncryption(table.SSEDescription)); err != nil {
		return create.SettingError(names.DynamoDB, ResNameTable, d.Id(), "server_side_encryption", err)
	}

	replicas := flattenReplicaDescriptions(table.Replicas)

	if replicas, err = addReplicaPITRs(conn, d.Id(), meta.(*conns.AWSClient).TerraformVersion, replicas); err != nil {
		return create.Error(names.DynamoDB, create.ErrActionReading, ResNameTable, d.Id(), err)
	}

	replicas = addReplicaTagPropagates(d.Get("replica").(*schema.Set), replicas)

	if err := d.Set("replica", replicas); err != nil {
		return create.SettingError(names.DynamoDB, ResNameTable, d.Id(), "replica", err)
	}

	if table.TableClassSummary != nil {
		d.Set("table_class", table.TableClassSummary.TableClass)
	} else {
		d.Set("table_class", nil)
	}

	pitrOut, err := conn.DescribeContinuousBackups(&dynamodb.DescribeContinuousBackupsInput{
		TableName: aws.String(d.Id()),
	})
	// When a Table is `ARCHIVED`, DescribeContinuousBackups returns `TableNotFoundException`
	if err != nil && !tfawserr.ErrCodeEquals(err, "UnknownOperationException", dynamodb.ErrCodeTableNotFoundException) {
		return create.Error(names.DynamoDB, create.ErrActionReading, ResNameTable, d.Id(), fmt.Errorf("continuous backups: %w", err))
	}

	if err := d.Set("point_in_time_recovery", flattenPITR(pitrOut)); err != nil {
		return create.SettingError(names.DynamoDB, ResNameTable, d.Id(), "point_in_time_recovery", err)
	}

	ttlOut, err := conn.DescribeTimeToLive(&dynamodb.DescribeTimeToLiveInput{
		TableName: aws.String(d.Id()),
	})

	if err != nil {
		return create.Error(names.DynamoDB, create.ErrActionReading, ResNameTable, d.Id(), fmt.Errorf("TTL: %w", err))
	}

	if err := d.Set("ttl", flattenTTL(ttlOut)); err != nil {
		return create.SettingError(names.DynamoDB, ResNameTable, d.Id(), "ttl", err)
	}

	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	tags, err := ListTags(conn, d.Get("arn").(string))
	// When a Table is `ARCHIVED`, ListTags returns `ResourceNotFoundException`
	if err != nil && !(tfawserr.ErrMessageContains(err, "UnknownOperationException", "Tagging is not currently supported in DynamoDB Local.") || tfresource.NotFound(err)) {
		return create.Error(names.DynamoDB, create.ErrActionReading, ResNameTable, d.Id(), fmt.Errorf("tags: %w", err))
	}

	tags = tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return create.SettingError(names.DynamoDB, ResNameTable, d.Id(), "tags", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return create.SettingError(names.DynamoDB, ResNameTable, d.Id(), "tags_all", err)
	}

	return nil
}

func resourceTableUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).DynamoDBConn
	o, n := d.GetChange("billing_mode")
	billingMode := n.(string)
	oldBillingMode := o.(string)

	// Global Secondary Index operations must occur in multiple phases
	// to prevent various error scenarios. If there are no detected required
	// updates in the Terraform configuration, later validation or API errors
	// will signal the problems.
	var gsiUpdates []*dynamodb.GlobalSecondaryIndexUpdate

	if d.HasChange("global_secondary_index") {
		var err error
		o, n := d.GetChange("global_secondary_index")
		gsiUpdates, err = UpdateDiffGSI(o.(*schema.Set).List(), n.(*schema.Set).List(), billingMode)

		if err != nil {
			return create.Error(names.DynamoDB, create.ErrActionUpdating, ResNameTable, d.Id(), fmt.Errorf("computing GSI difference: %w", err))
		}

		log.Printf("[DEBUG] Computed DynamoDB Table (%s) Global Secondary Index updates: %s", d.Id(), gsiUpdates)
	}

	// Phase 1 of Global Secondary Index Operations: Delete Only
	//  * Delete indexes first to prevent error when simultaneously updating
	//    BillingMode to PROVISIONED, which requires updating index
	//    ProvisionedThroughput first, but we have no definition
	//  * Only 1 online index can be deleted simultaneously per table
	for _, gsiUpdate := range gsiUpdates {
		if gsiUpdate.Delete == nil {
			continue
		}

		idxName := aws.StringValue(gsiUpdate.Delete.IndexName)
		input := &dynamodb.UpdateTableInput{
			GlobalSecondaryIndexUpdates: []*dynamodb.GlobalSecondaryIndexUpdate{gsiUpdate},
			TableName:                   aws.String(d.Id()),
		}

		if _, err := conn.UpdateTable(input); err != nil {
			return create.Error(names.DynamoDB, create.ErrActionDeleting, ResNameTable, d.Id(), fmt.Errorf("GSI (%s): %w", idxName, err))
		}

		if _, err := waitGSIDeleted(conn, d.Id(), idxName, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return create.Error(names.DynamoDB, create.ErrActionWaitingForDeletion, ResNameTable, d.Id(), fmt.Errorf("GSI (%s): %w", idxName, err))
		}
	}

	hasTableUpdate := false
	input := &dynamodb.UpdateTableInput{
		TableName: aws.String(d.Id()),
	}

	if d.HasChanges("billing_mode", "read_capacity", "write_capacity") {
		hasTableUpdate = true

		capacityMap := map[string]interface{}{
			"write_capacity": d.Get("write_capacity"),
			"read_capacity":  d.Get("read_capacity"),
		}

		input.BillingMode = aws.String(billingMode)
		input.ProvisionedThroughput = expandProvisionedThroughputUpdate(d.Id(), capacityMap, billingMode, oldBillingMode)
	}

	if d.HasChanges("stream_enabled", "stream_view_type") {
		hasTableUpdate = true

		input.StreamSpecification = &dynamodb.StreamSpecification{
			StreamEnabled: aws.Bool(d.Get("stream_enabled").(bool)),
		}
		if d.Get("stream_enabled").(bool) {
			input.StreamSpecification.StreamViewType = aws.String(d.Get("stream_view_type").(string))
		}
	}

	// Phase 2 of Global Secondary Index Operations: Update Only
	// Cannot create or delete index while updating table ProvisionedThroughput
	// Must skip all index updates when switching BillingMode from PROVISIONED to PAY_PER_REQUEST
	// Must update all indexes when switching BillingMode from PAY_PER_REQUEST to PROVISIONED
	if billingMode == dynamodb.BillingModeProvisioned {
		for _, gsiUpdate := range gsiUpdates {
			if gsiUpdate.Update == nil {
				continue
			}

			hasTableUpdate = true
			input.GlobalSecondaryIndexUpdates = append(input.GlobalSecondaryIndexUpdates, gsiUpdate)
		}
	}

	if d.HasChange("table_class") {
		hasTableUpdate = true
		input.TableClass = aws.String(d.Get("table_class").(string))
	}

	if hasTableUpdate {
		log.Printf("[DEBUG] Updating DynamoDB Table: %s", input)
		_, err := conn.UpdateTable(input)

		if err != nil {
			return create.Error(names.DynamoDB, create.ErrActionUpdating, ResNameTable, d.Id(), err)
		}

		if _, err := waitTableActive(conn, d.Id(), d.Timeout(schema.TimeoutUpdate)); err != nil {
			return create.Error(names.DynamoDB, create.ErrActionWaitingForUpdate, ResNameTable, d.Id(), err)
		}

		for _, gsiUpdate := range gsiUpdates {
			if gsiUpdate.Update == nil {
				continue
			}

			idxName := aws.StringValue(gsiUpdate.Update.IndexName)

			if _, err := waitGSIActive(conn, d.Id(), idxName, d.Timeout(schema.TimeoutUpdate)); err != nil {
				return create.Error(names.DynamoDB, create.ErrActionWaitingForUpdate, ResNameTable, d.Id(), fmt.Errorf("GSI (%s): %w", idxName, err))
			}
		}
	}

	// Phase 3 of Global Secondary Index Operations: Create Only
	// Only 1 online index can be created simultaneously per table
	for _, gsiUpdate := range gsiUpdates {
		if gsiUpdate.Create == nil {
			continue
		}

		idxName := aws.StringValue(gsiUpdate.Create.IndexName)
		input := &dynamodb.UpdateTableInput{
			AttributeDefinitions:        expandAttributes(d.Get("attribute").(*schema.Set).List()),
			GlobalSecondaryIndexUpdates: []*dynamodb.GlobalSecondaryIndexUpdate{gsiUpdate},
			TableName:                   aws.String(d.Id()),
		}

		if _, err := conn.UpdateTable(input); err != nil {
			return create.Error(names.DynamoDB, create.ErrActionUpdating, ResNameTable, d.Id(), fmt.Errorf("creating GSI (%s): %w", idxName, err))
		}

		if _, err := waitGSIActive(conn, d.Id(), idxName, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return create.Error(names.DynamoDB, create.ErrActionUpdating, ResNameTable, d.Id(), fmt.Errorf("%s GSI (%s): %w", create.ErrActionWaitingForCreation, idxName, err))
		}
	}

	if d.HasChange("server_side_encryption") {
		// "ValidationException: One or more parameter values were invalid: Server-Side Encryption modification must be the only operation in the request".
		_, err := conn.UpdateTable(&dynamodb.UpdateTableInput{
			TableName:        aws.String(d.Id()),
			SSESpecification: expandEncryptAtRestOptions(d.Get("server_side_encryption").([]interface{})),
		})
		if err != nil {
			return create.Error(names.DynamoDB, create.ErrActionUpdating, ResNameTable, d.Id(), fmt.Errorf("SSE: %w", err))
		}

		if _, err := waitSSEUpdated(conn, d.Id(), d.Timeout(schema.TimeoutUpdate)); err != nil {
			return create.Error(names.DynamoDB, create.ErrActionWaitingForUpdate, ResNameTable, d.Id(), err)
		}
	}

	if d.HasChange("ttl") {
		if err := updateTimeToLive(conn, d.Id(), d.Get("ttl").([]interface{}), d.Timeout(schema.TimeoutUpdate)); err != nil {
			return create.Error(names.DynamoDB, create.ErrActionUpdating, ResNameTable, d.Id(), err)
		}
	}

	replicaTagsChange := false
	if d.HasChange("replica") {
		replicaTagsChange = true

		if err := updateReplica(d, conn, meta.(*conns.AWSClient).TerraformVersion); err != nil {
			return create.Error(names.DynamoDB, create.ErrActionUpdating, ResNameTable, d.Id(), err)
		}
	}

	if d.HasChange("tags_all") {
		replicaTagsChange = true

		o, n := d.GetChange("tags_all")
		if err := UpdateTags(conn, d.Get("arn").(string), o, n); err != nil {
			return create.Error(names.DynamoDB, create.ErrActionUpdating, ResNameTable, d.Id(), err)
		}
	}

	if replicaTagsChange {
		if v, ok := d.Get("replica").(*schema.Set); ok && v.Len() > 0 {
			if err := updateReplicaTags(conn, d.Get("arn").(string), v.List(), d.Get("tags_all"), meta.(*conns.AWSClient).TerraformVersion); err != nil {
				return create.Error(names.DynamoDB, create.ErrActionUpdating, ResNameTable, d.Id(), err)
			}
		}
	}

	if d.HasChange("point_in_time_recovery") {
		if err := updatePITR(conn, d.Id(), d.Get("point_in_time_recovery.0.enabled").(bool), aws.StringValue(conn.Config.Region), meta.(*conns.AWSClient).TerraformVersion, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return create.Error(names.DynamoDB, create.ErrActionUpdating, ResNameTable, d.Id(), err)
		}
	}

	return resourceTableRead(d, meta)
}

func resourceTableDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).DynamoDBConn

	log.Printf("[DEBUG] DynamoDB delete table: %s", d.Id())

	if replicas := d.Get("replica").(*schema.Set).List(); len(replicas) > 0 {
		if err := deleteReplicas(conn, d.Id(), replicas, d.Timeout(schema.TimeoutDelete)); err != nil {
			// ValidationException: Replica specified in the Replica Update or Replica Delete action of the request was not found.
			if !tfawserr.ErrMessageContains(err, "ValidationException", "request was not found") {
				return create.Error(names.DynamoDB, create.ErrActionDeleting, ResNameTable, d.Id(), err)
			}
		}
	}

	err := deleteTable(conn, d.Id())
	if err != nil {
		if tfawserr.ErrMessageContains(err, dynamodb.ErrCodeResourceNotFoundException, "Requested resource not found: Table: ") {
			return nil
		}
		return create.Error(names.DynamoDB, create.ErrActionDeleting, ResNameTable, d.Id(), err)
	}

	if _, err := waitTableDeleted(conn, d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return create.Error(names.DynamoDB, create.ErrActionWaitingForDeletion, ResNameTable, d.Id(), err)
	}

	return nil
}

// custom diff

func isTableOptionDisabled(v interface{}) bool {
	options := v.([]interface{})
	if len(options) == 0 {
		return true
	}
	e := options[0].(map[string]interface{})["enabled"]
	return !e.(bool)
}

// CRUD helpers

func createReplicas(conn *dynamodb.DynamoDB, tableName string, tfList []interface{}, tfVersion string, create bool, timeout time.Duration) error {
	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]interface{})

		if !ok {
			continue
		}

		var replicaInput = &dynamodb.CreateReplicationGroupMemberAction{}

		if v, ok := tfMap["region_name"].(string); ok && v != "" {
			replicaInput.RegionName = aws.String(v)
		}

		if v, ok := tfMap["kms_key_arn"].(string); ok && v != "" {
			replicaInput.KMSMasterKeyId = aws.String(v)
		}

		input := &dynamodb.UpdateTableInput{
			TableName: aws.String(tableName),
			ReplicaUpdates: []*dynamodb.ReplicationGroupUpdate{
				{
					Create: replicaInput,
				},
			},
		}

		if !create {
			var replicaInput = &dynamodb.UpdateReplicationGroupMemberAction{}

			if v, ok := tfMap["region_name"].(string); ok && v != "" {
				replicaInput.RegionName = aws.String(v)
			}

			if v, ok := tfMap["kms_key_arn"].(string); ok && v != "" {
				replicaInput.KMSMasterKeyId = aws.String(v)
			}

			input = &dynamodb.UpdateTableInput{
				TableName: aws.String(tableName),
				ReplicaUpdates: []*dynamodb.ReplicationGroupUpdate{
					{
						Update: replicaInput,
					},
				},
			}
		}

		err := resource.Retry(maxDuration(replicaUpdateTimeout, timeout), func() *resource.RetryError {
			_, err := conn.UpdateTable(input)
			if err != nil {
				if tfawserr.ErrCodeEquals(err, "ThrottlingException") {
					return resource.RetryableError(err)
				}
				if tfawserr.ErrMessageContains(err, dynamodb.ErrCodeLimitExceededException, "can be created, updated, or deleted simultaneously") {
					return resource.RetryableError(err)
				}
				if tfawserr.ErrCodeEquals(err, dynamodb.ErrCodeResourceInUseException) {
					return resource.RetryableError(err)
				}

				return resource.NonRetryableError(err)
			}
			return nil
		})

		if tfresource.TimedOut(err) {
			_, err = conn.UpdateTable(input)
		}

		if create && tfawserr.ErrMessageContains(err, "ValidationException", "already exist") {
			return createReplicas(conn, tableName, tfList, tfVersion, false, timeout)
		}

		if err != nil && !tfawserr.ErrMessageContains(err, "ValidationException", "no actions specified") {
			return fmt.Errorf("creating replica (%s): %w", tfMap["region_name"].(string), err)
		}

		if err := waitReplicaActive(conn, tableName, tfMap["region_name"].(string), timeout); err != nil {
			return fmt.Errorf("waiting for replica (%s) creation: %w", tfMap["region_name"].(string), err)
		}

		// pitr
		if err = updatePITR(conn, tableName, tfMap["point_in_time_recovery"].(bool), tfMap["region_name"].(string), tfVersion, timeout); err != nil {
			return fmt.Errorf("updating replica (%s) point in time recovery: %w", tfMap["region_name"].(string), err)
		}
	}

	return nil
}

func updateReplicaTags(conn *dynamodb.DynamoDB, rn string, replicas []interface{}, newTags interface{}, terraformVersion string) error {
	for _, tfMapRaw := range replicas {
		tfMap, ok := tfMapRaw.(map[string]interface{})

		if !ok {
			continue
		}

		region, ok := tfMap["region_name"].(string)

		if !ok || region == "" {
			continue
		}

		if v, ok := tfMap["propagate_tags"].(bool); ok && v {
			session, err := conns.NewSessionForRegion(&conn.Config, region, terraformVersion)
			if err != nil {
				return fmt.Errorf("updating replica (%s) tags: %w", region, err)
			}

			conn = dynamodb.New(session)

			repARN, err := ARNForNewRegion(rn, region)
			if err != nil {
				return fmt.Errorf("per region ARN for replica (%s): %w", region, err)
			}

			oldTags, err := ListTags(conn, repARN)
			if err != nil {
				return fmt.Errorf("listing tags (%s): %w", repARN, err)
			}

			if err := UpdateTags(conn, repARN, oldTags, newTags); err != nil {
				return fmt.Errorf("updating tags: %w", err)
			}
		}
	}

	return nil
}

func updateTimeToLive(conn *dynamodb.DynamoDB, tableName string, ttlList []interface{}, timeout time.Duration) error {
	ttlMap := ttlList[0].(map[string]interface{})

	input := &dynamodb.UpdateTimeToLiveInput{
		TableName: aws.String(tableName),
		TimeToLiveSpecification: &dynamodb.TimeToLiveSpecification{
			AttributeName: aws.String(ttlMap["attribute_name"].(string)),
			Enabled:       aws.Bool(ttlMap["enabled"].(bool)),
		},
	}

	log.Printf("[DEBUG] Updating DynamoDB Table (%s) Time To Live: %s", tableName, input)
	if _, err := conn.UpdateTimeToLive(input); err != nil {
		return fmt.Errorf("updating Time To Live: %w", err)
	}

	log.Printf("[DEBUG] Waiting for DynamoDB Table (%s) Time to Live update to complete", tableName)

	if _, err := waitTTLUpdated(conn, tableName, ttlMap["enabled"].(bool), timeout); err != nil {
		return fmt.Errorf("waiting for Time To Live update: %w", err)
	}

	return nil
}

func updatePITR(conn *dynamodb.DynamoDB, tableName string, enabled bool, region string, tfVersion string, timeout time.Duration) error {
	// pitr must be modified from region where the main/replica resides
	log.Printf("[DEBUG] Updating DynamoDB point in time recovery status to %v (%s)", enabled, region)
	input := &dynamodb.UpdateContinuousBackupsInput{
		TableName: aws.String(tableName),
		PointInTimeRecoverySpecification: &dynamodb.PointInTimeRecoverySpecification{
			PointInTimeRecoveryEnabled: aws.Bool(enabled),
		},
	}

	if aws.StringValue(conn.Config.Region) != region {
		session, err := conns.NewSessionForRegion(&conn.Config, region, tfVersion)
		if err != nil {
			return fmt.Errorf("new session for region (%s): %w", region, err)
		}

		conn = dynamodb.New(session)
	}

	err := resource.Retry(updateTableContinuousBackupsTimeout, func() *resource.RetryError {
		_, err := conn.UpdateContinuousBackups(input)
		if err != nil {
			// Backups are still being enabled for this newly created table
			if tfawserr.ErrMessageContains(err, dynamodb.ErrCodeContinuousBackupsUnavailableException, "Backups are being enabled") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if tfresource.TimedOut(err) {
		_, err = conn.UpdateContinuousBackups(input)
	}

	if err != nil {
		return fmt.Errorf("updating PITR: %w", err)
	}

	if _, err := waitPITRUpdated(conn, tableName, enabled, timeout); err != nil {
		return fmt.Errorf("waiting for PITR update: %w", err)
	}

	return nil
}

func updateReplica(d *schema.ResourceData, conn *dynamodb.DynamoDB, tfVersion string) error {
	oRaw, nRaw := d.GetChange("replica")
	o := oRaw.(*schema.Set)
	n := nRaw.(*schema.Set)

	removed := o.Difference(n).List()
	added := n.Difference(o).List()

	// For true updates, don't remove and add, just update (i.e., keep in added
	// but remove from removed)
	for _, a := range added {
		for j, r := range removed {
			ma := a.(map[string]interface{})
			mr := r.(map[string]interface{})
			if ma["region_name"].(string) == mr["region_name"].(string) {
				removed = append(removed[:j], removed[j+1:]...)
				continue
			}
		}
	}

	if len(added) > 0 {
		if err := createReplicas(conn, d.Id(), added, tfVersion, true, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return fmt.Errorf("updating replicas, while creating: %w", err)
		}
	}

	if len(removed) > 0 {
		if err := deleteReplicas(conn, d.Id(), removed, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return fmt.Errorf("updating replicas, while deleting: %w", err)
		}
	}

	return nil
}

func UpdateDiffGSI(oldGsi, newGsi []interface{}, billingMode string) (ops []*dynamodb.GlobalSecondaryIndexUpdate, e error) {
	// Transform slices into maps
	oldGsis := make(map[string]interface{})
	for _, gsidata := range oldGsi {
		m := gsidata.(map[string]interface{})
		oldGsis[m["name"].(string)] = m
	}
	newGsis := make(map[string]interface{})
	for _, gsidata := range newGsi {
		m := gsidata.(map[string]interface{})
		// validate throughput input early, to avoid unnecessary processing
		if e = validateGSIProvisionedThroughput(m, billingMode); e != nil {
			return
		}
		newGsis[m["name"].(string)] = m
	}

	for _, data := range newGsi {
		newMap := data.(map[string]interface{})
		newName := newMap["name"].(string)

		if _, exists := oldGsis[newName]; !exists {
			m := data.(map[string]interface{})
			idxName := m["name"].(string)

			ops = append(ops, &dynamodb.GlobalSecondaryIndexUpdate{
				Create: &dynamodb.CreateGlobalSecondaryIndexAction{
					IndexName:             aws.String(idxName),
					KeySchema:             expandKeySchema(m),
					ProvisionedThroughput: expandProvisionedThroughput(m, billingMode),
					Projection:            expandProjection(m),
				},
			})
		}
	}

	for _, data := range oldGsi {
		oldMap := data.(map[string]interface{})
		oldName := oldMap["name"].(string)

		newData, exists := newGsis[oldName]
		if exists {
			newMap := newData.(map[string]interface{})
			idxName := newMap["name"].(string)

			oldWriteCapacity, oldReadCapacity := oldMap["write_capacity"].(int), oldMap["read_capacity"].(int)
			newWriteCapacity, newReadCapacity := newMap["write_capacity"].(int), newMap["read_capacity"].(int)
			capacityChanged := (oldWriteCapacity != newWriteCapacity || oldReadCapacity != newReadCapacity)

			// pluck non_key_attributes from oldAttributes and newAttributes as reflect.DeepEquals will compare
			// ordinal of elements in its equality (which we actually don't care about)
			nonKeyAttributesChanged := checkIfNonKeyAttributesChanged(oldMap, newMap)

			oldAttributes, err := stripCapacityAttributes(oldMap)
			if err != nil {
				return ops, err
			}
			oldAttributes, err = stripNonKeyAttributes(oldAttributes)
			if err != nil {
				return ops, err
			}
			newAttributes, err := stripCapacityAttributes(newMap)
			if err != nil {
				return ops, err
			}
			newAttributes, err = stripNonKeyAttributes(newAttributes)
			if err != nil {
				return ops, err
			}
			otherAttributesChanged := nonKeyAttributesChanged || !reflect.DeepEqual(oldAttributes, newAttributes)

			if capacityChanged && !otherAttributesChanged {
				update := &dynamodb.GlobalSecondaryIndexUpdate{
					Update: &dynamodb.UpdateGlobalSecondaryIndexAction{
						IndexName:             aws.String(idxName),
						ProvisionedThroughput: expandProvisionedThroughput(newMap, billingMode),
					},
				}
				ops = append(ops, update)
			} else if otherAttributesChanged {
				// Other attributes cannot be updated
				ops = append(ops, &dynamodb.GlobalSecondaryIndexUpdate{
					Delete: &dynamodb.DeleteGlobalSecondaryIndexAction{
						IndexName: aws.String(idxName),
					},
				})

				ops = append(ops, &dynamodb.GlobalSecondaryIndexUpdate{
					Create: &dynamodb.CreateGlobalSecondaryIndexAction{
						IndexName:             aws.String(idxName),
						KeySchema:             expandKeySchema(newMap),
						ProvisionedThroughput: expandProvisionedThroughput(newMap, billingMode),
						Projection:            expandProjection(newMap),
					},
				})
			}
		} else {
			idxName := oldName
			ops = append(ops, &dynamodb.GlobalSecondaryIndexUpdate{
				Delete: &dynamodb.DeleteGlobalSecondaryIndexAction{
					IndexName: aws.String(idxName),
				},
			})
		}
	}
	return ops, nil
}

func deleteTable(conn *dynamodb.DynamoDB, tableName string) error {
	input := &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	}

	err := resource.Retry(deleteTableTimeout, func() *resource.RetryError {
		_, err := conn.DeleteTable(input)
		if err != nil {
			// Subscriber limit exceeded: Only 10 tables can be created, updated, or deleted simultaneously
			if tfawserr.ErrMessageContains(err, dynamodb.ErrCodeLimitExceededException, "simultaneously") {
				return resource.RetryableError(err)
			}
			// This handles multiple scenarios in the DynamoDB API:
			// 1. Updating a table immediately before deletion may return:
			//    ResourceInUseException: Attempt to change a resource which is still in use: Table is being updated:
			// 2. Removing a table from a DynamoDB global table may return:
			//    ResourceInUseException: Attempt to change a resource which is still in use: Table is being deleted:
			if tfawserr.ErrCodeEquals(err, dynamodb.ErrCodeResourceInUseException) {
				return resource.RetryableError(err)
			}
			if tfawserr.ErrMessageContains(err, dynamodb.ErrCodeResourceNotFoundException, "Requested resource not found: Table: ") {
				return resource.NonRetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if tfresource.TimedOut(err) {
		_, err = conn.DeleteTable(input)
	}

	return err
}

func deleteReplicas(conn *dynamodb.DynamoDB, tableName string, tfList []interface{}, timeout time.Duration) error {
	var g multierror.Group

	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]interface{})

		if !ok {
			continue
		}

		var regionName string

		if v, ok := tfMap["region_name"].(string); ok {
			regionName = v
		}

		if regionName == "" {
			continue
		}

		g.Go(func() error {
			input := &dynamodb.UpdateTableInput{
				TableName: aws.String(tableName),
				ReplicaUpdates: []*dynamodb.ReplicationGroupUpdate{
					{
						Delete: &dynamodb.DeleteReplicationGroupMemberAction{
							RegionName: aws.String(regionName),
						},
					},
				},
			}

			err := resource.Retry(updateTableTimeout, func() *resource.RetryError {
				_, err := conn.UpdateTable(input)
				if err != nil {
					if tfawserr.ErrCodeEquals(err, "ThrottlingException") {
						return resource.RetryableError(err)
					}
					if tfawserr.ErrMessageContains(err, dynamodb.ErrCodeLimitExceededException, "can be created, updated, or deleted simultaneously") {
						return resource.RetryableError(err)
					}
					if tfawserr.ErrCodeEquals(err, dynamodb.ErrCodeResourceInUseException) {
						return resource.RetryableError(err)
					}

					return resource.NonRetryableError(err)
				}
				return nil
			})

			if tfresource.TimedOut(err) {
				_, err = conn.UpdateTable(input)
			}

			if err != nil {
				return fmt.Errorf("deleting replica (%s): %w", regionName, err)
			}

			if err := waitReplicaDeleted(conn, tableName, regionName, timeout); err != nil {
				return fmt.Errorf("waiting for replica (%s) deletion: %w", regionName, err)
			}

			return nil
		})
	}

	return g.Wait().ErrorOrNil()
}

func replicaPITR(conn *dynamodb.DynamoDB, tableName string, region string, tfVersion string) (bool, error) {
	// At a future time, replicas should probably have a separate resource because,
	// to manage them, you need connections from the different regions. However, they
	// have to be created from the starting/main region...
	session, err := conns.NewSessionForRegion(&conn.Config, region, tfVersion)
	if err != nil {
		return false, fmt.Errorf("new session for replica (%s) PITR: %w", region, err)
	}

	conn = dynamodb.New(session)

	pitrOut, err := conn.DescribeContinuousBackups(&dynamodb.DescribeContinuousBackupsInput{
		TableName: aws.String(tableName),
	})
	// When a Table is `ARCHIVED`, DescribeContinuousBackups returns `TableNotFoundException`
	if err != nil && !tfawserr.ErrCodeEquals(err, "UnknownOperationException", dynamodb.ErrCodeTableNotFoundException) {
		return false, fmt.Errorf("describing Continuous Backups: %w", err)
	}

	if pitrOut == nil {
		return false, nil
	}

	enabled := false

	if pitrOut.ContinuousBackupsDescription != nil {
		pitr := pitrOut.ContinuousBackupsDescription.PointInTimeRecoveryDescription
		if pitr != nil {
			enabled = (aws.StringValue(pitr.PointInTimeRecoveryStatus) == dynamodb.PointInTimeRecoveryStatusEnabled)
		}
	}

	return enabled, nil
}

func addReplicaPITRs(conn *dynamodb.DynamoDB, tableName string, tfVersion string, replicas []interface{}) ([]interface{}, error) {
	// This non-standard approach is needed because PITR info for a replica
	// must come from a region-specific connection. A future table_replica
	// resource may improve this.

	for i, replicaRaw := range replicas {
		replica := replicaRaw.(map[string]interface{})

		var enabled bool
		var err error
		if enabled, err = replicaPITR(conn, tableName, replica["region_name"].(string), tfVersion); err != nil {
			return nil, err
		}
		replica["point_in_time_recovery"] = enabled
		replicas[i] = replica
	}

	return replicas, nil
}

func addReplicaTagPropagates(configReplicas *schema.Set, replicas []interface{}) []interface{} {
	if configReplicas.Len() == 0 {
		return replicas
	}

	l := configReplicas.List()

	for i, replicaRaw := range replicas {
		replica := replicaRaw.(map[string]interface{})

		prop := false

		for _, configReplicaRaw := range l {
			configReplica := configReplicaRaw.(map[string]interface{})

			if v, ok := configReplica["region_name"].(string); ok && v != replica["region_name"].(string) {
				continue
			}

			if v, ok := configReplica["propagate_tags"].(bool); ok && v {
				prop = true
				break
			}
		}
		replica["propagate_tags"] = prop
		replicas[i] = replica
	}

	return replicas
}

// flatteners, expanders

func flattenTableAttributeDefinitions(definitions []*dynamodb.AttributeDefinition) []interface{} {
	if len(definitions) == 0 {
		return []interface{}{}
	}

	var attributes []interface{}

	for _, d := range definitions {
		if d == nil {
			continue
		}

		m := map[string]string{
			"name": aws.StringValue(d.AttributeName),
			"type": aws.StringValue(d.AttributeType),
		}

		attributes = append(attributes, m)
	}

	return attributes
}

func flattenTableLocalSecondaryIndex(lsi []*dynamodb.LocalSecondaryIndexDescription) []interface{} {
	if len(lsi) == 0 {
		return []interface{}{}
	}

	var output []interface{}

	for _, l := range lsi {
		if l == nil {
			continue
		}

		m := map[string]interface{}{
			"name": aws.StringValue(l.IndexName),
		}

		if l.Projection != nil {
			m["projection_type"] = aws.StringValue(l.Projection.ProjectionType)
			m["non_key_attributes"] = aws.StringValueSlice(l.Projection.NonKeyAttributes)
		}

		for _, attribute := range l.KeySchema {
			if attribute == nil {
				continue
			}
			if aws.StringValue(attribute.KeyType) == dynamodb.KeyTypeRange {
				m["range_key"] = aws.StringValue(attribute.AttributeName)
			}
		}

		output = append(output, m)
	}

	return output
}

func flattenTableGlobalSecondaryIndex(gsi []*dynamodb.GlobalSecondaryIndexDescription) []interface{} {
	if len(gsi) == 0 {
		return []interface{}{}
	}

	var output []interface{}

	for _, g := range gsi {
		if g == nil {
			continue
		}

		gsi := make(map[string]interface{})

		if g.ProvisionedThroughput != nil {
			gsi["write_capacity"] = aws.Int64Value(g.ProvisionedThroughput.WriteCapacityUnits)
			gsi["read_capacity"] = aws.Int64Value(g.ProvisionedThroughput.ReadCapacityUnits)
			gsi["name"] = aws.StringValue(g.IndexName)
		}

		for _, attribute := range g.KeySchema {
			if attribute == nil {
				continue
			}

			if aws.StringValue(attribute.KeyType) == dynamodb.KeyTypeHash {
				gsi["hash_key"] = aws.StringValue(attribute.AttributeName)
			}

			if aws.StringValue(attribute.KeyType) == dynamodb.KeyTypeRange {
				gsi["range_key"] = aws.StringValue(attribute.AttributeName)
			}
		}

		if g.Projection != nil {
			gsi["projection_type"] = aws.StringValue(g.Projection.ProjectionType)
			gsi["non_key_attributes"] = aws.StringValueSlice(g.Projection.NonKeyAttributes)
		}

		output = append(output, gsi)
	}

	return output
}

func flattenTableServerSideEncryption(description *dynamodb.SSEDescription) []interface{} {
	if description == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"enabled":     aws.StringValue(description.Status) == dynamodb.SSEStatusEnabled,
		"kms_key_arn": aws.StringValue(description.KMSMasterKeyArn),
	}

	return []interface{}{m}
}

func expandAttributes(cfg []interface{}) []*dynamodb.AttributeDefinition {
	attributes := make([]*dynamodb.AttributeDefinition, len(cfg))
	for i, attribute := range cfg {
		attr := attribute.(map[string]interface{})
		attributes[i] = &dynamodb.AttributeDefinition{
			AttributeName: aws.String(attr["name"].(string)),
			AttributeType: aws.String(attr["type"].(string)),
		}
	}
	return attributes
}

func flattenReplicaDescription(apiObject *dynamodb.ReplicaDescription) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if apiObject.KMSMasterKeyId != nil {
		tfMap["kms_key_arn"] = aws.StringValue(apiObject.KMSMasterKeyId)
	}

	if apiObject.RegionName != nil {
		tfMap["region_name"] = aws.StringValue(apiObject.RegionName)
	}

	return tfMap
}

func flattenReplicaDescriptions(apiObjects []*dynamodb.ReplicaDescription) []interface{} {
	if len(apiObjects) == 0 {
		return nil
	}

	var tfList []interface{}

	for _, apiObject := range apiObjects {
		if apiObject == nil {
			continue
		}

		tfList = append(tfList, flattenReplicaDescription(apiObject))
	}

	return tfList
}

func flattenTTL(ttlOutput *dynamodb.DescribeTimeToLiveOutput) []interface{} {
	m := map[string]interface{}{
		"enabled": false,
	}

	if ttlOutput == nil || ttlOutput.TimeToLiveDescription == nil {
		return []interface{}{m}
	}

	ttlDesc := ttlOutput.TimeToLiveDescription

	m["attribute_name"] = aws.StringValue(ttlDesc.AttributeName)
	m["enabled"] = (aws.StringValue(ttlDesc.TimeToLiveStatus) == dynamodb.TimeToLiveStatusEnabled)

	return []interface{}{m}
}

func flattenPITR(pitrDesc *dynamodb.DescribeContinuousBackupsOutput) []interface{} {
	m := map[string]interface{}{
		"enabled": false,
	}

	if pitrDesc == nil {
		return []interface{}{m}
	}

	if pitrDesc.ContinuousBackupsDescription != nil {
		pitr := pitrDesc.ContinuousBackupsDescription.PointInTimeRecoveryDescription
		if pitr != nil {
			m["enabled"] = (aws.StringValue(pitr.PointInTimeRecoveryStatus) == dynamodb.PointInTimeRecoveryStatusEnabled)
		}
	}

	return []interface{}{m}
}

// TODO: Get rid of keySchemaM - the user should just explicitly define
// this in the config, we shouldn't magically be setting it like this.
// Removal will however require config change, hence BC. :/
func expandLocalSecondaryIndexes(cfg []interface{}, keySchemaM map[string]interface{}) []*dynamodb.LocalSecondaryIndex {
	indexes := make([]*dynamodb.LocalSecondaryIndex, len(cfg))
	for i, lsi := range cfg {
		m := lsi.(map[string]interface{})
		idxName := m["name"].(string)

		// TODO: See https://github.com/hashicorp/terraform-provider-aws/issues/3176
		if _, ok := m["hash_key"]; !ok {
			m["hash_key"] = keySchemaM["hash_key"]
		}

		indexes[i] = &dynamodb.LocalSecondaryIndex{
			IndexName:  aws.String(idxName),
			KeySchema:  expandKeySchema(m),
			Projection: expandProjection(m),
		}
	}
	return indexes
}

func expandGlobalSecondaryIndex(data map[string]interface{}, billingMode string) *dynamodb.GlobalSecondaryIndex {
	return &dynamodb.GlobalSecondaryIndex{
		IndexName:             aws.String(data["name"].(string)),
		KeySchema:             expandKeySchema(data),
		Projection:            expandProjection(data),
		ProvisionedThroughput: expandProvisionedThroughput(data, billingMode),
	}
}

func expandProvisionedThroughput(data map[string]interface{}, billingMode string) *dynamodb.ProvisionedThroughput {
	return expandProvisionedThroughputUpdate("", data, billingMode, "")
}

func expandProvisionedThroughputUpdate(id string, data map[string]interface{}, billingMode, oldBillingMode string) *dynamodb.ProvisionedThroughput {
	if billingMode == dynamodb.BillingModePayPerRequest {
		return nil
	}

	return &dynamodb.ProvisionedThroughput{
		ReadCapacityUnits:  aws.Int64(expandProvisionedThroughputField(id, data, "read_capacity", billingMode, oldBillingMode)),
		WriteCapacityUnits: aws.Int64(expandProvisionedThroughputField(id, data, "write_capacity", billingMode, oldBillingMode)),
	}
}

func expandProvisionedThroughputField(id string, data map[string]interface{}, key, billingMode, oldBillingMode string) int64 {
	v := data[key].(int)
	if v == 0 && billingMode == dynamodb.BillingModeProvisioned && oldBillingMode == dynamodb.BillingModePayPerRequest {
		log.Printf("[WARN] Overriding %[1]s on DynamoDB Table (%[2]s) to %[3]d. Switching from billing mode %[4]q to %[5]q without value for %[1]s. Assuming changes are being ignored.",
			key, id, provisionedThroughputMinValue, oldBillingMode, billingMode)
		v = provisionedThroughputMinValue
	}
	return int64(v)
}

func expandProjection(data map[string]interface{}) *dynamodb.Projection {
	projection := &dynamodb.Projection{
		ProjectionType: aws.String(data["projection_type"].(string)),
	}

	if v, ok := data["non_key_attributes"].([]interface{}); ok && len(v) > 0 {
		projection.NonKeyAttributes = flex.ExpandStringList(v)
	}

	if v, ok := data["non_key_attributes"].(*schema.Set); ok && v.Len() > 0 {
		projection.NonKeyAttributes = flex.ExpandStringSet(v)
	}

	return projection
}

func expandKeySchema(data map[string]interface{}) []*dynamodb.KeySchemaElement {
	keySchema := []*dynamodb.KeySchemaElement{}

	if v, ok := data["hash_key"]; ok && v != nil && v != "" {
		keySchema = append(keySchema, &dynamodb.KeySchemaElement{
			AttributeName: aws.String(v.(string)),
			KeyType:       aws.String(dynamodb.KeyTypeHash),
		})
	}

	if v, ok := data["range_key"]; ok && v != nil && v != "" {
		keySchema = append(keySchema, &dynamodb.KeySchemaElement{
			AttributeName: aws.String(v.(string)),
			KeyType:       aws.String(dynamodb.KeyTypeRange),
		})
	}

	return keySchema
}

func expandEncryptAtRestOptions(vOptions []interface{}) *dynamodb.SSESpecification {
	options := &dynamodb.SSESpecification{}

	enabled := false
	if len(vOptions) > 0 {
		mOptions := vOptions[0].(map[string]interface{})

		enabled = mOptions["enabled"].(bool)
		if enabled {
			if vKmsKeyArn, ok := mOptions["kms_key_arn"].(string); ok && vKmsKeyArn != "" {
				options.KMSMasterKeyId = aws.String(vKmsKeyArn)
				options.SSEType = aws.String(dynamodb.SSETypeKms)
			}
		}
	}
	options.Enabled = aws.Bool(enabled)

	return options
}

// validators

func validateTableAttributes(d *schema.ResourceDiff) error {
	// Collect all indexed attributes
	indexedAttributes := map[string]bool{}

	if v, ok := d.GetOk("hash_key"); ok {
		indexedAttributes[v.(string)] = true
	}
	if v, ok := d.GetOk("range_key"); ok {
		indexedAttributes[v.(string)] = true
	}
	if v, ok := d.GetOk("local_secondary_index"); ok {
		indexes := v.(*schema.Set).List()
		for _, idx := range indexes {
			index := idx.(map[string]interface{})
			rangeKey := index["range_key"].(string)
			indexedAttributes[rangeKey] = true
		}
	}
	if v, ok := d.GetOk("global_secondary_index"); ok {
		indexes := v.(*schema.Set).List()
		for _, idx := range indexes {
			index := idx.(map[string]interface{})

			hashKey := index["hash_key"].(string)
			indexedAttributes[hashKey] = true

			if rk, ok := index["range_key"].(string); ok && rk != "" {
				indexedAttributes[rk] = true
			}
		}
	}

	// Check if all indexed attributes have an attribute definition
	attributes := d.Get("attribute").(*schema.Set).List()
	unindexedAttributes := []string{}
	for _, attr := range attributes {
		attribute := attr.(map[string]interface{})
		attrName := attribute["name"].(string)

		if _, ok := indexedAttributes[attrName]; !ok {
			unindexedAttributes = append(unindexedAttributes, attrName)
		} else {
			delete(indexedAttributes, attrName)
		}
	}

	var err *multierror.Error

	if len(unindexedAttributes) > 0 {
		err = multierror.Append(err, fmt.Errorf("all attributes must be indexed. Unused attributes: %q", unindexedAttributes))
	}

	if len(indexedAttributes) > 0 {
		missingIndexes := []string{}
		for index := range indexedAttributes {
			missingIndexes = append(missingIndexes, index)
		}

		err = multierror.Append(err, fmt.Errorf("all indexes must match a defined attribute. Unmatched indexes: %q", missingIndexes))
	}

	return err.ErrorOrNil()
}

func validateGSIProvisionedThroughput(data map[string]interface{}, billingMode string) error {
	// if billing mode is PAY_PER_REQUEST, don't need to validate the throughput settings
	if billingMode == dynamodb.BillingModePayPerRequest {
		return nil
	}

	writeCapacity, writeCapacitySet := data["write_capacity"].(int)
	readCapacity, readCapacitySet := data["read_capacity"].(int)

	if !writeCapacitySet || !readCapacitySet {
		return fmt.Errorf("read and write capacity should be set when billing mode is %s", dynamodb.BillingModeProvisioned)
	}

	if writeCapacity < 1 {
		return fmt.Errorf("write capacity must be > 0 when billing mode is %s", dynamodb.BillingModeProvisioned)
	}

	if readCapacity < 1 {
		return fmt.Errorf("read capacity must be > 0 when billing mode is %s", dynamodb.BillingModeProvisioned)
	}

	return nil
}

func validateProvisionedThroughputField(diff *schema.ResourceDiff, key string) error {
	oldBillingMode, billingMode := diff.GetChange("billing_mode")
	v := diff.Get(key).(int)
	if billingMode == dynamodb.BillingModeProvisioned {
		if v < provisionedThroughputMinValue {
			// Assuming the field is ignored, likely due to autoscaling
			if oldBillingMode == dynamodb.BillingModePayPerRequest {
				return nil
			}
			return fmt.Errorf("%s must be at least 1 when billing_mode is %q", key, billingMode)
		}
	} else if billingMode == dynamodb.BillingModePayPerRequest && oldBillingMode != dynamodb.BillingModeProvisioned {
		if v != 0 {
			return fmt.Errorf("%s can not be set when billing_mode is %q", key, dynamodb.BillingModePayPerRequest)
		}
	}
	return nil
}
