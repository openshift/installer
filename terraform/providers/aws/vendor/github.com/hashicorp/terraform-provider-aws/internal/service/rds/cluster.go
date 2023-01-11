package rds

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
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

const (
	clusterScalingConfiguration_DefaultMinCapacity = 1
	clusterScalingConfiguration_DefaultMaxCapacity = 16
	clusterTimeoutDelete                           = 2 * time.Minute
)

func ResourceCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceClusterCreate,
		Read:   resourceClusterRead,
		Update: resourceClusterUpdate,
		Delete: resourceClusterDelete,

		Importer: &schema.ResourceImporter{
			State: resourceClusterImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Update: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(120 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"allocated_storage": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"allow_major_version_upgrade": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			// apply_immediately is used to determine when the update modifications take place.
			// See http://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Overview.DBInstance.Modifying.html
			"apply_immediately": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zones": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"backup_retention_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntAtMost(35),
			},
			"backtrack_window": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 259200),
			},
			"cluster_identifier": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validIdentifier,
				ConflictsWith: []string{"cluster_identifier_prefix"},
			},
			"cluster_identifier_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validIdentifierPrefix,
				ConflictsWith: []string{"cluster_identifier"},
			},
			"cluster_members": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"cluster_resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"copy_tags_to_snapshot": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"database_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"db_cluster_instance_class": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_cluster_parameter_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"db_instance_parameter_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_subnet_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_global_write_forwarding": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"enabled_cloudwatch_logs_exports": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(ClusterExportableLogType_Values(), false),
				},
			},
			"enable_http_endpoint": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      ClusterEngineAurora,
				ValidateFunc: validClusterEngine(),
			},
			"engine_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      EngineModeProvisioned,
				ValidateFunc: validation.StringInSlice(EngineMode_Values(), false),
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"engine_version_actual": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"final_snapshot_identifier": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, es []error) {
					value := v.(string)
					if !regexp.MustCompile(`^[0-9A-Za-z-]+$`).MatchString(value) {
						es = append(es, fmt.Errorf(
							"only alphanumeric characters and hyphens allowed in %q", k))
					}
					if regexp.MustCompile(`--`).MatchString(value) {
						es = append(es, fmt.Errorf("%q cannot contain two consecutive hyphens", k))
					}
					if regexp.MustCompile(`-$`).MatchString(value) {
						es = append(es, fmt.Errorf("%q cannot end in a hyphen", k))
					}
					return
				},
			},
			"global_cluster_identifier": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hosted_zone_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"iam_database_authentication_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"iam_roles": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"iops": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"kms_key_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidARN,
			},
			"master_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"master_username": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(NetworkType_Values(), false),
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"preferred_backup_window": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: verify.ValidOnceADayWindowFormat,
			},
			"preferred_maintenance_window": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				StateFunc: func(val interface{}) string {
					if val == nil {
						return ""
					}
					return strings.ToLower(val.(string))
				},
				ValidateFunc: verify.ValidOnceAWeekWindowFormat,
			},
			"reader_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"replication_source_identifier": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"restore_to_point_in_time": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"restore_to_time": {
							Type:          schema.TypeString,
							Optional:      true,
							ForceNew:      true,
							ValidateFunc:  verify.ValidUTCTimestamp,
							ConflictsWith: []string{"restore_to_point_in_time.0.use_latest_restorable_time"},
						},
						"restore_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice(RestoreType_Values(), false),
						},
						"source_cluster_identifier": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.Any(
								verify.ValidARN,
								validIdentifier,
							),
						},
						"use_latest_restorable_time": {
							Type:          schema.TypeBool,
							Optional:      true,
							ForceNew:      true,
							ConflictsWith: []string{"restore_to_point_in_time.0.restore_to_time"},
						},
					},
				},
				ConflictsWith: []string{
					"s3_import",
					"snapshot_identifier",
				},
			},
			"s3_import": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"bucket_prefix": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"ingestion_role": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"source_engine": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"source_engine_version": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
				ConflictsWith: []string{
					"snapshot_identifier",
					"restore_to_point_in_time",
				},
			},
			"scaling_configuration": {
				Type:             schema.TypeList,
				Optional:         true,
				MaxItems:         1,
				DiffSuppressFunc: verify.SuppressMissingOptionalConfigurationBlock,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_pause": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"max_capacity": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  clusterScalingConfiguration_DefaultMaxCapacity,
						},
						"min_capacity": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  clusterScalingConfiguration_DefaultMinCapacity,
						},
						"seconds_until_auto_pause": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      300,
							ValidateFunc: validation.IntBetween(300, 86400),
						},
						"timeout_action": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      TimeoutActionRollbackCapacityChange,
							ValidateFunc: validation.StringInSlice(TimeoutAction_Values(), false),
						},
					},
				},
			},
			"serverlessv2_scaling_configuration": {
				Type:             schema.TypeList,
				Optional:         true,
				MaxItems:         1,
				DiffSuppressFunc: verify.SuppressMissingOptionalConfigurationBlock,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_capacity": {
							Type:         schema.TypeFloat,
							Required:     true,
							ValidateFunc: validation.FloatBetween(0.5, 128),
						},
						"min_capacity": {
							Type:         schema.TypeFloat,
							Required:     true,
							ValidateFunc: validation.FloatBetween(0.5, 128),
						},
					},
				},
			},
			"skip_final_snapshot": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"snapshot_identifier": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"storage_encrypted": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"storage_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
			"vpc_security_group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceClusterCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).RDSConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	// Some API calls (e.g. RestoreDBClusterFromSnapshot do not support all
	// parameters to correctly apply all settings in one pass. For missing
	// parameters or unsupported configurations, we may need to call
	// ModifyDBInstance afterwards to prevent Terraform operators from API
	// errors or needing to double apply.
	var requiresModifyDbCluster bool
	modifyDbClusterInput := &rds.ModifyDBClusterInput{
		ApplyImmediately: aws.Bool(true),
	}

	var identifier string
	if v, ok := d.GetOk("cluster_identifier"); ok {
		identifier = v.(string)
	} else if v, ok := d.GetOk("cluster_identifier_prefix"); ok {
		identifier = resource.PrefixedUniqueId(v.(string))
	} else {
		identifier = resource.PrefixedUniqueId("tf-")
	}

	if v, ok := d.GetOk("snapshot_identifier"); ok {
		input := &rds.RestoreDBClusterFromSnapshotInput{
			CopyTagsToSnapshot:  aws.Bool(d.Get("copy_tags_to_snapshot").(bool)),
			DBClusterIdentifier: aws.String(identifier),
			DeletionProtection:  aws.Bool(d.Get("deletion_protection").(bool)),
			Engine:              aws.String(d.Get("engine").(string)),
			EngineMode:          aws.String(d.Get("engine_mode").(string)),
			SnapshotIdentifier:  aws.String(v.(string)),
			Tags:                Tags(tags.IgnoreAWS()),
		}

		if v, ok := d.GetOk("availability_zones"); ok && v.(*schema.Set).Len() > 0 {
			input.AvailabilityZones = flex.ExpandStringSet(v.(*schema.Set))
		}

		if v, ok := d.GetOk("backtrack_window"); ok {
			input.BacktrackWindow = aws.Int64(int64(v.(int)))
		}

		if v, ok := d.GetOk("backup_retention_period"); ok {
			modifyDbClusterInput.BackupRetentionPeriod = aws.Int64(int64(v.(int)))
			requiresModifyDbCluster = true
		}

		if v, ok := d.GetOk("database_name"); ok {
			input.DatabaseName = aws.String(v.(string))
		}

		if v, ok := d.GetOk("db_cluster_parameter_group_name"); ok {
			input.DBClusterParameterGroupName = aws.String(v.(string))
		}

		if v, ok := d.GetOk("db_subnet_group_name"); ok {
			input.DBSubnetGroupName = aws.String(v.(string))
		}

		if v, ok := d.GetOk("enabled_cloudwatch_logs_exports"); ok && v.(*schema.Set).Len() > 0 {
			input.EnableCloudwatchLogsExports = flex.ExpandStringSet(v.(*schema.Set))
		}

		if v, ok := d.GetOk("engine_version"); ok {
			input.EngineVersion = aws.String(v.(string))
		}

		if v, ok := d.GetOk("kms_key_id"); ok {
			input.KmsKeyId = aws.String(v.(string))
		}

		if v, ok := d.GetOk("master_password"); ok {
			modifyDbClusterInput.MasterUserPassword = aws.String(v.(string))
			requiresModifyDbCluster = true
		}

		if v, ok := d.GetOk("network_type"); ok {
			input.NetworkType = aws.String(v.(string))
		}

		if v, ok := d.GetOk("option_group_name"); ok {
			input.OptionGroupName = aws.String(v.(string))
		}

		if v, ok := d.GetOk("port"); ok {
			input.Port = aws.Int64(int64(v.(int)))
		}

		if v, ok := d.GetOk("preferred_backup_window"); ok {
			modifyDbClusterInput.PreferredBackupWindow = aws.String(v.(string))
			requiresModifyDbCluster = true
		}

		if v, ok := d.GetOk("preferred_maintenance_window"); ok {
			modifyDbClusterInput.PreferredMaintenanceWindow = aws.String(v.(string))
			requiresModifyDbCluster = true
		}

		if v, ok := d.GetOk("scaling_configuration"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
			input.ScalingConfiguration = expandScalingConfiguration(v.([]interface{})[0].(map[string]interface{}))
		}

		if v, ok := d.GetOk("serverlessv2_scaling_configuration"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
			modifyDbClusterInput.ServerlessV2ScalingConfiguration = expandServerlessV2ScalingConfiguration(v.([]interface{})[0].(map[string]interface{}))
			requiresModifyDbCluster = true
		}

		if v, ok := d.GetOk("vpc_security_group_ids"); ok && v.(*schema.Set).Len() > 0 {
			input.VpcSecurityGroupIds = flex.ExpandStringSet(v.(*schema.Set))
		}

		log.Printf("[DEBUG] Creating RDS Cluster: %s", input)
		_, err := tfresource.RetryWhenAWSErrMessageContains(propagationTimeout,
			func() (interface{}, error) {
				return conn.RestoreDBClusterFromSnapshot(input)
			},
			errCodeInvalidParameterValue, "IAM role ARN value is invalid or does not include the required permissions")

		if err != nil {
			return fmt.Errorf("creating RDS Cluster (restore from snapshot) (%s): %w", identifier, err)
		}
	} else if v, ok := d.GetOk("s3_import"); ok {
		if _, ok := d.GetOk("master_password"); !ok {
			return fmt.Errorf(`provider.aws: aws_db_instance: %s: "master_password": required field is not set`, d.Get("name").(string))
		}
		if _, ok := d.GetOk("master_username"); !ok {
			return fmt.Errorf(`provider.aws: aws_db_instance: %s: "master_username": required field is not set`, d.Get("name").(string))
		}

		tfMap := v.([]interface{})[0].(map[string]interface{})
		input := &rds.RestoreDBClusterFromS3Input{
			CopyTagsToSnapshot:  aws.Bool(d.Get("copy_tags_to_snapshot").(bool)),
			DBClusterIdentifier: aws.String(identifier),
			DeletionProtection:  aws.Bool(d.Get("deletion_protection").(bool)),
			Engine:              aws.String(d.Get("engine").(string)),
			MasterUsername:      aws.String(d.Get("master_username").(string)),
			MasterUserPassword:  aws.String(d.Get("master_password").(string)),
			S3BucketName:        aws.String(tfMap["bucket_name"].(string)),
			S3IngestionRoleArn:  aws.String(tfMap["ingestion_role"].(string)),
			S3Prefix:            aws.String(tfMap["bucket_prefix"].(string)),
			SourceEngine:        aws.String(tfMap["source_engine"].(string)),
			SourceEngineVersion: aws.String(tfMap["source_engine_version"].(string)),
			Tags:                Tags(tags.IgnoreAWS()),
		}

		if v, ok := d.GetOk("availability_zones"); ok && v.(*schema.Set).Len() > 0 {
			input.AvailabilityZones = flex.ExpandStringSet(v.(*schema.Set))
		}

		if v, ok := d.GetOk("backtrack_window"); ok {
			input.BacktrackWindow = aws.Int64(int64(v.(int)))
		}

		if v, ok := d.GetOk("backup_retention_period"); ok {
			input.BackupRetentionPeriod = aws.Int64(int64(v.(int)))
		}

		if v := d.Get("database_name"); v.(string) != "" {
			input.DatabaseName = aws.String(v.(string))
		}

		if v, ok := d.GetOk("db_cluster_parameter_group_name"); ok {
			input.DBClusterParameterGroupName = aws.String(v.(string))
		}

		if v, ok := d.GetOk("db_subnet_group_name"); ok {
			input.DBSubnetGroupName = aws.String(v.(string))
		}

		if v, ok := d.GetOk("enabled_cloudwatch_logs_exports"); ok && v.(*schema.Set).Len() > 0 {
			input.EnableCloudwatchLogsExports = flex.ExpandStringSet(v.(*schema.Set))
		}

		if v, ok := d.GetOk("engine_version"); ok {
			input.EngineVersion = aws.String(v.(string))
		}

		if v, ok := d.GetOk("iam_database_authentication_enabled"); ok {
			input.EnableIAMDatabaseAuthentication = aws.Bool(v.(bool))
		}

		if v, ok := d.GetOk("kms_key_id"); ok {
			input.KmsKeyId = aws.String(v.(string))
		}

		if v, ok := d.GetOk("network_type"); ok {
			input.NetworkType = aws.String(v.(string))
		}

		if v, ok := d.GetOk("port"); ok {
			input.Port = aws.Int64(int64(v.(int)))
		}

		if v, ok := d.GetOk("preferred_backup_window"); ok {
			input.PreferredBackupWindow = aws.String(v.(string))
		}

		if v, ok := d.GetOk("preferred_maintenance_window"); ok {
			input.PreferredMaintenanceWindow = aws.String(v.(string))
		}

		if v, ok := d.GetOkExists("storage_encrypted"); ok {
			input.StorageEncrypted = aws.Bool(v.(bool))
		}

		if v, ok := d.GetOk("vpc_security_group_ids"); ok && v.(*schema.Set).Len() > 0 {
			input.VpcSecurityGroupIds = flex.ExpandStringSet(v.(*schema.Set))
		}

		log.Printf("[DEBUG] Creating RDS Cluster: %s", input)
		_, err := tfresource.RetryWhen(propagationTimeout,
			func() (interface{}, error) {
				return conn.RestoreDBClusterFromS3(input)
			},
			func(err error) (bool, error) {
				// InvalidParameterValue: Files from the specified Amazon S3 bucket cannot be downloaded.
				// Make sure that you have created an AWS Identity and Access Management (IAM) role that lets Amazon RDS access Amazon S3 for you.
				if tfawserr.ErrMessageContains(err, errCodeInvalidParameterValue, "Files from the specified Amazon S3 bucket cannot be downloaded") {
					return true, err
				}

				if tfawserr.ErrMessageContains(err, errCodeInvalidParameterValue, "S3_SNAPSHOT_INGESTION") {
					return true, err
				}

				if tfawserr.ErrMessageContains(err, errCodeInvalidParameterValue, "S3 bucket cannot be found") {
					return true, err
				}

				return false, err
			},
		)

		if err != nil {
			return fmt.Errorf("creating RDS Cluster (restore from S3) (%s): %w", identifier, err)
		}
	} else if v, ok := d.GetOk("restore_to_point_in_time"); ok {
		tfMap := v.([]interface{})[0].(map[string]interface{})
		input := &rds.RestoreDBClusterToPointInTimeInput{
			DBClusterIdentifier:       aws.String(identifier),
			DeletionProtection:        aws.Bool(d.Get("deletion_protection").(bool)),
			SourceDBClusterIdentifier: aws.String(tfMap["source_cluster_identifier"].(string)),
			Tags:                      Tags(tags.IgnoreAWS()),
		}

		if v, ok := tfMap["restore_to_time"].(string); ok && v != "" {
			v, _ := time.Parse(time.RFC3339, v)
			input.RestoreToTime = aws.Time(v)
		}

		if v, ok := tfMap["use_latest_restorable_time"].(bool); ok && v {
			input.UseLatestRestorableTime = aws.Bool(v)
		}

		if input.RestoreToTime == nil && input.UseLatestRestorableTime == nil {
			return fmt.Errorf(`provider.aws: aws_rds_cluster: %s: Either "restore_to_time" or "use_latest_restorable_time" must be set`, d.Get("database_name").(string))
		}

		if v, ok := d.GetOk("backtrack_window"); ok {
			input.BacktrackWindow = aws.Int64(int64(v.(int)))
		}

		if v, ok := d.GetOk("backup_retention_period"); ok {
			modifyDbClusterInput.BackupRetentionPeriod = aws.Int64(int64(v.(int)))
			requiresModifyDbCluster = true
		}

		if v, ok := d.GetOk("db_cluster_parameter_group_name"); ok {
			input.DBClusterParameterGroupName = aws.String(v.(string))
		}

		if v, ok := d.GetOk("db_subnet_group_name"); ok {
			input.DBSubnetGroupName = aws.String(v.(string))
		}

		if v, ok := d.GetOk("enabled_cloudwatch_logs_exports"); ok && v.(*schema.Set).Len() > 0 {
			input.EnableCloudwatchLogsExports = flex.ExpandStringSet(v.(*schema.Set))
		}

		if v, ok := d.GetOk("iam_database_authentication_enabled"); ok {
			input.EnableIAMDatabaseAuthentication = aws.Bool(v.(bool))
		}

		if v, ok := d.GetOk("kms_key_id"); ok {
			input.KmsKeyId = aws.String(v.(string))
		}

		if v, ok := d.GetOk("master_password"); ok {
			modifyDbClusterInput.MasterUserPassword = aws.String(v.(string))
			requiresModifyDbCluster = true
		}

		if v, ok := d.GetOk("network_type"); ok {
			input.NetworkType = aws.String(v.(string))
		}

		if v, ok := d.GetOk("option_group_name"); ok {
			input.OptionGroupName = aws.String(v.(string))
		}

		if v, ok := d.GetOk("port"); ok {
			input.Port = aws.Int64(int64(v.(int)))
		}

		if v, ok := d.GetOk("preferred_backup_window"); ok {
			modifyDbClusterInput.PreferredBackupWindow = aws.String(v.(string))
			requiresModifyDbCluster = true
		}

		if v, ok := d.GetOk("preferred_maintenance_window"); ok {
			modifyDbClusterInput.PreferredMaintenanceWindow = aws.String(v.(string))
			requiresModifyDbCluster = true
		}

		if v, ok := tfMap["restore_type"].(string); ok {
			input.RestoreType = aws.String(v)
		}

		if v, ok := d.GetOk("scaling_configuration"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
			modifyDbClusterInput.ScalingConfiguration = expandScalingConfiguration(v.([]interface{})[0].(map[string]interface{}))
			requiresModifyDbCluster = true
		}

		if v, ok := d.GetOk("serverlessv2_scaling_configuration"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
			modifyDbClusterInput.ServerlessV2ScalingConfiguration = expandServerlessV2ScalingConfiguration(v.([]interface{})[0].(map[string]interface{}))
			requiresModifyDbCluster = true
		}

		if v, ok := d.GetOk("vpc_security_group_ids"); ok && v.(*schema.Set).Len() > 0 {
			input.VpcSecurityGroupIds = flex.ExpandStringSet(v.(*schema.Set))
		}

		log.Printf("[DEBUG] Creating RDS Cluster: %s", input)
		_, err := conn.RestoreDBClusterToPointInTime(input)

		if err != nil {
			return fmt.Errorf("creating RDS Cluster (restore to point-in-time) (%s): %w", identifier, err)
		}
	} else {
		input := &rds.CreateDBClusterInput{
			CopyTagsToSnapshot:  aws.Bool(d.Get("copy_tags_to_snapshot").(bool)),
			DBClusterIdentifier: aws.String(identifier),
			DeletionProtection:  aws.Bool(d.Get("deletion_protection").(bool)),
			Engine:              aws.String(d.Get("engine").(string)),
			EngineMode:          aws.String(d.Get("engine_mode").(string)),
			Tags:                Tags(tags.IgnoreAWS()),
		}

		if v, ok := d.GetOkExists("allocated_storage"); ok {
			input.AllocatedStorage = aws.Int64(int64(v.(int)))
		}

		if v, ok := d.GetOk("availability_zones"); ok && v.(*schema.Set).Len() > 0 {
			input.AvailabilityZones = flex.ExpandStringSet(v.(*schema.Set))
		}

		if v, ok := d.GetOk("backtrack_window"); ok {
			input.BacktrackWindow = aws.Int64(int64(v.(int)))
		}

		if v, ok := d.GetOk("backup_retention_period"); ok {
			input.BackupRetentionPeriod = aws.Int64(int64(v.(int)))
		}

		if v := d.Get("database_name"); v.(string) != "" {
			input.DatabaseName = aws.String(v.(string))
		}

		if v, ok := d.GetOk("db_cluster_instance_class"); ok {
			input.DBClusterInstanceClass = aws.String(v.(string))
		}

		if v, ok := d.GetOk("db_cluster_parameter_group_name"); ok {
			input.DBClusterParameterGroupName = aws.String(v.(string))
		}

		if v, ok := d.GetOk("db_subnet_group_name"); ok {
			input.DBSubnetGroupName = aws.String(v.(string))
		}

		if v, ok := d.GetOk("enable_global_write_forwarding"); ok {
			input.EnableGlobalWriteForwarding = aws.Bool(v.(bool))
		}

		if v, ok := d.GetOk("enable_http_endpoint"); ok {
			input.EnableHttpEndpoint = aws.Bool(v.(bool))
		}

		if v, ok := d.GetOk("enabled_cloudwatch_logs_exports"); ok && v.(*schema.Set).Len() > 0 {
			input.EnableCloudwatchLogsExports = flex.ExpandStringSet(v.(*schema.Set))
		}

		if v, ok := d.GetOk("engine_version"); ok {
			input.EngineVersion = aws.String(v.(string))
		}

		if v, ok := d.GetOk("global_cluster_identifier"); ok {
			input.GlobalClusterIdentifier = aws.String(v.(string))
		}

		if v, ok := d.GetOk("iam_database_authentication_enabled"); ok {
			input.EnableIAMDatabaseAuthentication = aws.Bool(v.(bool))
		}

		if v, ok := d.GetOkExists("iops"); ok {
			input.Iops = aws.Int64(int64(v.(int)))
		}

		if v, ok := d.GetOk("kms_key_id"); ok {
			input.KmsKeyId = aws.String(v.(string))
		}

		// Note: Username and password credentials are required and valid
		// unless the cluster is a read-replica. This also applies to clusters
		// within a global cluster. Providing a password and/or username for
		// a replica will result in an InvalidParameterValue error.
		if v, ok := d.GetOk("master_password"); ok {
			input.MasterUserPassword = aws.String(v.(string))
		}

		if v, ok := d.GetOk("master_username"); ok {
			input.MasterUsername = aws.String(v.(string))
		}

		if v, ok := d.GetOk("network_type"); ok {
			input.NetworkType = aws.String(v.(string))
		}

		if v, ok := d.GetOk("port"); ok {
			input.Port = aws.Int64(int64(v.(int)))
		}

		if v, ok := d.GetOk("preferred_backup_window"); ok {
			input.PreferredBackupWindow = aws.String(v.(string))
		}

		if v, ok := d.GetOk("preferred_maintenance_window"); ok {
			input.PreferredMaintenanceWindow = aws.String(v.(string))
		}

		if v, ok := d.GetOk("replication_source_identifier"); ok && input.GlobalClusterIdentifier == nil {
			input.ReplicationSourceIdentifier = aws.String(v.(string))
		}

		if v, ok := d.GetOk("scaling_configuration"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
			input.ScalingConfiguration = expandScalingConfiguration(v.([]interface{})[0].(map[string]interface{}))
		}

		if v, ok := d.GetOk("serverlessv2_scaling_configuration"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
			input.ServerlessV2ScalingConfiguration = expandServerlessV2ScalingConfiguration(v.([]interface{})[0].(map[string]interface{}))
		}

		if v, ok := d.GetOk("source_region"); ok {
			input.SourceRegion = aws.String(v.(string))
		}

		if v, ok := d.GetOkExists("storage_encrypted"); ok {
			input.StorageEncrypted = aws.Bool(v.(bool))
		}

		if v, ok := d.GetOkExists("storage_type"); ok {
			input.StorageType = aws.String(v.(string))
		}

		if v, ok := d.GetOk("vpc_security_group_ids"); ok && v.(*schema.Set).Len() > 0 {
			input.VpcSecurityGroupIds = flex.ExpandStringSet(v.(*schema.Set))
		}

		log.Printf("[DEBUG] Creating RDS Cluster: %s", input)
		_, err := tfresource.RetryWhenAWSErrMessageContains(propagationTimeout,
			func() (interface{}, error) {
				return conn.CreateDBCluster(input)
			},
			errCodeInvalidParameterValue, "IAM role ARN value is invalid or does not include the required permissions")

		if err != nil {
			return fmt.Errorf("creating RDS Cluster (%s): %w", identifier, err)
		}
	}

	d.SetId(identifier)

	if _, err := waitDBClusterCreated(conn, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
		return fmt.Errorf("waiting for RDS Cluster (%s) create: %w", d.Id(), err)
	}

	if v, ok := d.GetOk("iam_roles"); ok && v.(*schema.Set).Len() > 0 {
		for _, v := range v.(*schema.Set).List() {
			if err := addIAMRoleToCluster(conn, d.Id(), v.(string)); err != nil {
				return err
			}
		}
	}

	if requiresModifyDbCluster {
		modifyDbClusterInput.DBClusterIdentifier = aws.String(d.Id())

		log.Printf("[INFO] Modifying RDS Cluster: %s", modifyDbClusterInput)
		_, err := conn.ModifyDBCluster(modifyDbClusterInput)

		if err != nil {
			return fmt.Errorf("updating RDS Cluster (%s): %w", d.Id(), err)
		}

		if _, err := waitDBClusterUpdated(conn, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
			return fmt.Errorf("waiting for RDS Cluster (%s) update: %w", d.Id(), err)
		}
	}

	return resourceClusterRead(d, meta)
}

func resourceClusterRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).RDSConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	dbc, err := FindDBClusterByID(conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] RDS Cluster (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("reading RDS Cluster (%s): %w", d.Id(), err)
	}

	d.Set("allocated_storage", dbc.AllocatedStorage)
	clusterARN := aws.StringValue(dbc.DBClusterArn)
	d.Set("arn", clusterARN)
	d.Set("availability_zones", aws.StringValueSlice(dbc.AvailabilityZones))
	d.Set("backtrack_window", dbc.BacktrackWindow)
	d.Set("backup_retention_period", dbc.BackupRetentionPeriod)
	d.Set("cluster_identifier", dbc.DBClusterIdentifier)
	var clusterMembers []string
	for _, v := range dbc.DBClusterMembers {
		clusterMembers = append(clusterMembers, aws.StringValue(v.DBInstanceIdentifier))
	}
	d.Set("cluster_members", clusterMembers)
	d.Set("cluster_resource_id", dbc.DbClusterResourceId)
	d.Set("copy_tags_to_snapshot", dbc.CopyTagsToSnapshot)
	// Only set the DatabaseName if it is not nil. There is a known API bug where
	// RDS accepts a DatabaseName but does not return it, causing a perpetual
	// diff.
	//	See https://github.com/hashicorp/terraform/issues/4671 for backstory
	if dbc.DatabaseName != nil {
		d.Set("database_name", dbc.DatabaseName)
	}
	d.Set("db_cluster_instance_class", dbc.DBClusterInstanceClass)
	d.Set("db_cluster_parameter_group_name", dbc.DBClusterParameterGroup)
	d.Set("db_subnet_group_name", dbc.DBSubnetGroup)
	d.Set("deletion_protection", dbc.DeletionProtection)
	d.Set("enabled_cloudwatch_logs_exports", aws.StringValueSlice(dbc.EnabledCloudwatchLogsExports))
	d.Set("enable_http_endpoint", dbc.HttpEndpointEnabled)
	d.Set("endpoint", dbc.Endpoint)
	d.Set("engine", dbc.Engine)
	d.Set("engine_mode", dbc.EngineMode)
	clusterSetResourceDataEngineVersionFromCluster(d, dbc)
	d.Set("hosted_zone_id", dbc.HostedZoneId)
	d.Set("iam_database_authentication_enabled", dbc.IAMDatabaseAuthenticationEnabled)
	var iamRoleARNs []string
	for _, v := range dbc.AssociatedRoles {
		iamRoleARNs = append(iamRoleARNs, aws.StringValue(v.RoleArn))
	}
	d.Set("iam_roles", iamRoleARNs)
	d.Set("iops", dbc.Iops)
	d.Set("kms_key_id", dbc.KmsKeyId)
	d.Set("master_username", dbc.MasterUsername)
	d.Set("network_type", dbc.NetworkType)
	d.Set("port", dbc.Port)
	d.Set("preferred_backup_window", dbc.PreferredBackupWindow)
	d.Set("preferred_maintenance_window", dbc.PreferredMaintenanceWindow)
	d.Set("reader_endpoint", dbc.ReaderEndpoint)
	d.Set("replication_source_identifier", dbc.ReplicationSourceIdentifier)
	if dbc.ScalingConfigurationInfo != nil {
		if err := d.Set("scaling_configuration", []interface{}{flattenScalingConfigurationInfo(dbc.ScalingConfigurationInfo)}); err != nil {
			return fmt.Errorf("setting scaling_configuration: %w", err)
		}
	} else {
		d.Set("scaling_configuration", nil)
	}
	if dbc.ServerlessV2ScalingConfiguration != nil {
		if err := d.Set("serverlessv2_scaling_configuration", []interface{}{flattenServerlessV2ScalingConfigurationInfo(dbc.ServerlessV2ScalingConfiguration)}); err != nil {
			return fmt.Errorf("setting serverlessv2_scaling_configuration: %w", err)
		}
	} else {
		d.Set("serverlessv2_scaling_configuration", nil)
	}
	d.Set("storage_encrypted", dbc.StorageEncrypted)
	d.Set("storage_type", dbc.StorageType)
	var securityGroupIDs []string
	for _, v := range dbc.VpcSecurityGroups {
		securityGroupIDs = append(securityGroupIDs, aws.StringValue(v.VpcSecurityGroupId))
	}
	d.Set("vpc_security_group_ids", securityGroupIDs)

	tags, err := ListTags(conn, clusterARN)

	if err != nil {
		return fmt.Errorf("listing tags for RDS Cluster (%s): %w", d.Id(), err)
	}

	tags = tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("setting tags_all: %w", err)
	}

	// Fetch and save Global Cluster if engine mode global
	d.Set("global_cluster_identifier", "")

	if aws.StringValue(dbc.EngineMode) == EngineModeGlobal || aws.StringValue(dbc.EngineMode) == EngineModeProvisioned {
		globalCluster, err := FindGlobalClusterByDBClusterARN(conn, aws.StringValue(dbc.DBClusterArn))

		if err == nil {
			d.Set("global_cluster_identifier", globalCluster.GlobalClusterIdentifier)
		} else if tfresource.NotFound(err) || tfawserr.ErrMessageContains(err, errCodeInvalidParameterValue, "Access Denied to API Version: APIGlobalDatabases") {
			// Ignore the following API error for regions/partitions that do not support RDS Global Clusters:
			// InvalidParameterValue: Access Denied to API Version: APIGlobalDatabases
		} else {
			return fmt.Errorf("reading RDS Global Cluster for RDS Cluster (%s): %w", d.Id(), err)
		}
	}

	return nil
}

func resourceClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).RDSConn

	if d.HasChangesExcept(
		"allow_major_version_upgrade",
		"final_snapshot_identifier",
		"global_cluster_identifier",
		"iam_roles",
		"replication_source_identifier",
		"skip_final_snapshot",
		"tags", "tags_all") {
		input := &rds.ModifyDBClusterInput{
			ApplyImmediately:    aws.Bool(d.Get("apply_immediately").(bool)),
			DBClusterIdentifier: aws.String(d.Id()),
		}

		if d.HasChange("allocated_storage") {
			input.AllocatedStorage = aws.Int64(int64(d.Get("allocated_storage").(int)))
		}

		if v, ok := d.GetOk("allow_major_version_upgrade"); ok {
			input.AllowMajorVersionUpgrade = aws.Bool(v.(bool))
		}

		if d.HasChange("backtrack_window") {
			input.BacktrackWindow = aws.Int64(int64(d.Get("backtrack_window").(int)))
		}

		if d.HasChange("backup_retention_period") {
			input.BackupRetentionPeriod = aws.Int64(int64(d.Get("backup_retention_period").(int)))
		}

		if d.HasChange("copy_tags_to_snapshot") {
			input.CopyTagsToSnapshot = aws.Bool(d.Get("copy_tags_to_snapshot").(bool))
		}

		if d.HasChange("db_cluster_instance_class") {
			input.EngineVersion = aws.String(d.Get("db_cluster_instance_class").(string))
		}

		if d.HasChange("db_cluster_parameter_group_name") {
			input.DBClusterParameterGroupName = aws.String(d.Get("db_cluster_parameter_group_name").(string))
		}

		if d.HasChange("db_instance_parameter_group_name") {
			input.DBInstanceParameterGroupName = aws.String(d.Get("db_instance_parameter_group_name").(string))
		}

		if d.HasChange("deletion_protection") {
			input.DeletionProtection = aws.Bool(d.Get("deletion_protection").(bool))
		}

		if d.HasChange("enable_global_write_forwarding") {
			input.EnableGlobalWriteForwarding = aws.Bool(d.Get("enable_global_write_forwarding").(bool))
		}

		if d.HasChange("enable_http_endpoint") {
			input.EnableHttpEndpoint = aws.Bool(d.Get("enable_http_endpoint").(bool))
		}

		if d.HasChange("enabled_cloudwatch_logs_exports") {
			oraw, nraw := d.GetChange("enabled_cloudwatch_logs_exports")
			o := oraw.(*schema.Set)
			n := nraw.(*schema.Set)

			input.CloudwatchLogsExportConfiguration = &rds.CloudwatchLogsExportConfiguration{
				DisableLogTypes: flex.ExpandStringSet(o.Difference(n)),
				EnableLogTypes:  flex.ExpandStringSet(n.Difference(o)),
			}
		}

		if d.HasChange("engine_version") {
			input.EngineVersion = aws.String(d.Get("engine_version").(string))
		}

		if d.HasChange("iam_database_authentication_enabled") {
			input.EnableIAMDatabaseAuthentication = aws.Bool(d.Get("iam_database_authentication_enabled").(bool))
		}

		if d.HasChange("iops") {
			input.Iops = aws.Int64(int64(d.Get("iops").(int)))
		}

		if d.HasChange("master_password") {
			input.MasterUserPassword = aws.String(d.Get("master_password").(string))
		}

		if d.HasChange("network_type") {
			input.NetworkType = aws.String(d.Get("network_type").(string))
		}

		if d.HasChange("port") {
			input.Port = aws.Int64(int64(d.Get("port").(int)))
		}

		if d.HasChange("preferred_backup_window") {
			input.PreferredBackupWindow = aws.String(d.Get("preferred_backup_window").(string))
		}

		if d.HasChange("preferred_maintenance_window") {
			input.PreferredMaintenanceWindow = aws.String(d.Get("preferred_maintenance_window").(string))
		}

		if d.HasChange("scaling_configuration") {
			if v, ok := d.GetOk("scaling_configuration"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
				input.ScalingConfiguration = expandScalingConfiguration(v.([]interface{})[0].(map[string]interface{}))
			}
		}

		if d.HasChange("serverlessv2_scaling_configuration") {
			if v, ok := d.GetOk("serverlessv2_scaling_configuration"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
				input.ServerlessV2ScalingConfiguration = expandServerlessV2ScalingConfiguration(v.([]interface{})[0].(map[string]interface{}))
			}
		}

		if d.HasChange("storage_type") {
			input.StorageType = aws.String(d.Get("storage_type").(string))
		}

		if d.HasChange("vpc_security_group_ids") {
			if v, ok := d.GetOk("vpc_security_group_ids"); ok && v.(*schema.Set).Len() > 0 {
				input.VpcSecurityGroupIds = flex.ExpandStringSet(v.(*schema.Set))
			} else {
				input.VpcSecurityGroupIds = aws.StringSlice(nil)
			}
		}

		log.Printf("[DEBUG] Modifying RDS Cluster: %s", input)
		_, err := tfresource.RetryWhen(5*time.Minute,
			func() (interface{}, error) {
				return conn.ModifyDBCluster(input)
			},
			func(err error) (bool, error) {
				if tfawserr.ErrMessageContains(err, errCodeInvalidParameterValue, "IAM role ARN value is invalid or does not include the required permissions") {
					return true, err
				}

				if tfawserr.ErrCodeEquals(err, rds.ErrCodeInvalidDBClusterStateFault) {
					return true, err
				}

				return false, err
			},
		)

		if err != nil {
			return fmt.Errorf("updating RDS Cluster (%s): %w", d.Id(), err)
		}

		if _, err := waitDBClusterUpdated(conn, d.Id(), d.Timeout(schema.TimeoutUpdate)); err != nil {
			return fmt.Errorf("waiting for RDS Cluster (%s) update: %w", d.Id(), err)
		}
	}

	if d.HasChange("global_cluster_identifier") {
		oRaw, nRaw := d.GetChange("global_cluster_identifier")
		o := oRaw.(string)
		n := nRaw.(string)

		if o == "" {
			return errors.New("Existing RDS Clusters cannot be added to an existing RDS Global Cluster")
		}

		if n != "" {
			return errors.New("Existing RDS Clusters cannot be migrated between existing RDS Global Clusters")
		}

		input := &rds.RemoveFromGlobalClusterInput{
			DbClusterIdentifier:     aws.String(d.Get("arn").(string)),
			GlobalClusterIdentifier: aws.String(o),
		}

		log.Printf("[DEBUG] Removing RDS Cluster from RDS Global Cluster: %s", input)
		_, err := conn.RemoveFromGlobalCluster(input)

		if err != nil && !tfawserr.ErrCodeEquals(err, rds.ErrCodeGlobalClusterNotFoundFault) && !tfawserr.ErrMessageContains(err, "InvalidParameterValue", "is not found in global cluster") {
			return fmt.Errorf("removing RDS Cluster (%s) from RDS Global Cluster: %w", d.Id(), err)
		}
	}

	if d.HasChange("iam_roles") {
		oraw, nraw := d.GetChange("iam_roles")
		if oraw == nil {
			oraw = new(schema.Set)
		}
		if nraw == nil {
			nraw = new(schema.Set)
		}
		os := oraw.(*schema.Set)
		ns := nraw.(*schema.Set)

		for _, v := range ns.Difference(os).List() {
			if err := addIAMRoleToCluster(conn, d.Id(), v.(string)); err != nil {
				return err
			}
		}

		for _, v := range os.Difference(ns).List() {
			if err := removeIAMRoleFromCluster(conn, d.Id(), v.(string)); err != nil {
				return err
			}
		}
	}

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := UpdateTags(conn, d.Get("arn").(string), o, n); err != nil {
			return fmt.Errorf("updating RDS Cluster (%s) tags: %w", d.Get("arn").(string), err)
		}
	}

	return resourceClusterRead(d, meta)
}

func resourceClusterDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).RDSConn

	// Automatically remove from global cluster to bypass this error on deletion:
	// InvalidDBClusterStateFault: This cluster is a part of a global cluster, please remove it from globalcluster first
	if d.Get("global_cluster_identifier").(string) != "" {
		globalClusterID := d.Get("global_cluster_identifier").(string)
		input := &rds.RemoveFromGlobalClusterInput{
			DbClusterIdentifier:     aws.String(d.Get("arn").(string)),
			GlobalClusterIdentifier: aws.String(globalClusterID),
		}

		log.Printf("[DEBUG] Removing RDS Cluster from RDS Global Cluster: %s", input)
		_, err := conn.RemoveFromGlobalCluster(input)

		if err != nil && !tfawserr.ErrCodeEquals(err, rds.ErrCodeGlobalClusterNotFoundFault) && !tfawserr.ErrMessageContains(err, "InvalidParameterValue", "is not found in global cluster") {
			return fmt.Errorf("removing RDS Cluster (%s) from RDS Global Cluster (%s): %w", d.Id(), globalClusterID, err)
		}
	}

	skipFinalSnapshot := d.Get("skip_final_snapshot").(bool)
	input := &rds.DeleteDBClusterInput{
		DBClusterIdentifier: aws.String(d.Id()),
		SkipFinalSnapshot:   aws.Bool(skipFinalSnapshot),
	}

	if !skipFinalSnapshot {
		if v, ok := d.GetOk("final_snapshot_identifier"); ok {
			input.FinalDBSnapshotIdentifier = aws.String(v.(string))
		} else {
			return fmt.Errorf("RDS Cluster final_snapshot_identifier is required when skip_final_snapshot is false")
		}
	}

	log.Printf("[DEBUG] Deleting RDS Cluster: %s", d.Id())
	_, err := tfresource.RetryWhen(clusterTimeoutDelete,
		func() (interface{}, error) {
			return conn.DeleteDBCluster(input)
		},
		func(err error) (bool, error) {
			if tfawserr.ErrMessageContains(err, "InvalidParameterCombination", "disable deletion pro") {
				if v, ok := d.GetOk("deletion_protection"); (!ok || !v.(bool)) && d.Get("apply_immediately").(bool) {
					_, err := tfresource.RetryWhen(d.Timeout(schema.TimeoutDelete),
						func() (interface{}, error) {
							return conn.ModifyDBCluster(&rds.ModifyDBClusterInput{
								ApplyImmediately:    aws.Bool(true),
								DBClusterIdentifier: aws.String(d.Id()),
								DeletionProtection:  aws.Bool(false),
							})
						},
						func(err error) (bool, error) {
							if tfawserr.ErrMessageContains(err, errCodeInvalidParameterValue, "IAM role ARN value is invalid or does not include the required permissions") {
								return true, err
							}

							if tfawserr.ErrCodeEquals(err, rds.ErrCodeInvalidDBClusterStateFault) {
								return true, err
							}

							return false, err
						},
					)

					if err != nil {
						return false, fmt.Errorf("modifying RDS Cluster (%s) DeletionProtection=false: %w", d.Id(), err)
					}

					if _, err := waitDBClusterUpdated(conn, d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
						return false, fmt.Errorf("waiting for RDS Cluster (%s) update: %w", d.Id(), err)
					}
				}

				return true, err
			}

			if tfawserr.ErrMessageContains(err, rds.ErrCodeInvalidDBClusterStateFault, "is not currently in the available state") {
				return true, err
			}

			if tfawserr.ErrMessageContains(err, rds.ErrCodeInvalidDBClusterStateFault, "cluster is a part of a global cluster") {
				return true, err
			}

			return false, err
		},
	)

	if tfawserr.ErrCodeEquals(err, rds.ErrCodeDBClusterNotFoundFault) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("deleting RDS Cluster (%s): %w", d.Id(), err)
	}

	if _, err := waitDBClusterDeleted(conn, d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return fmt.Errorf("waiting for RDS Cluster (%s) delete: %w", d.Id(), err)
	}

	return nil
}

func resourceClusterImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// Neither skip_final_snapshot nor final_snapshot_identifier can be fetched
	// from any API call, so we need to default skip_final_snapshot to true so
	// that final_snapshot_identifier is not required
	d.Set("skip_final_snapshot", true)
	return []*schema.ResourceData{d}, nil
}

func addIAMRoleToCluster(conn *rds.RDS, clusterID, roleARN string) error {
	input := &rds.AddRoleToDBClusterInput{
		DBClusterIdentifier: aws.String(clusterID),
		RoleArn:             aws.String(roleARN),
	}

	_, err := conn.AddRoleToDBCluster(input)

	if err != nil {
		return fmt.Errorf("adding IAM Role (%s) to RDS Cluster (%s): %w", roleARN, clusterID, err)
	}

	return nil
}

func removeIAMRoleFromCluster(conn *rds.RDS, clusterID, roleARN string) error {
	input := &rds.RemoveRoleFromDBClusterInput{
		DBClusterIdentifier: aws.String(clusterID),
		RoleArn:             aws.String(roleARN),
	}

	_, err := conn.RemoveRoleFromDBCluster(input)

	if err != nil {
		return fmt.Errorf("removing IAM Role (%s) from RDS Cluster (%s): %w", roleARN, clusterID, err)
	}

	return err
}

func clusterSetResourceDataEngineVersionFromCluster(d *schema.ResourceData, c *rds.DBCluster) {
	oldVersion := d.Get("engine_version").(string)
	newVersion := aws.StringValue(c.EngineVersion)
	compareActualEngineVersion(d, oldVersion, newVersion)
}
