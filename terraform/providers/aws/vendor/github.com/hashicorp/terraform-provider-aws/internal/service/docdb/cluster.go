package docdb

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/docdb"
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
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"availability_zones": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				ForceNew: true,
				Computed: true,
				Set:      schema.HashString,
			},

			"cluster_identifier": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"cluster_identifier_prefix"},
				ValidateFunc:  validIdentifier,
			},
			"cluster_identifier_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"cluster_identifier"},
				ValidateFunc:  validIdentifierPrefix,
			},

			"cluster_members": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
				Set:      schema.HashString,
			},

			"db_subnet_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"db_cluster_parameter_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"global_cluster_identifier": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validGlobalCusterIdentifier,
			},

			"reader_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"hosted_zone_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"engine": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "docdb",
				ForceNew:     true,
				ValidateFunc: validEngine(),
			},

			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"storage_encrypted": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
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

			"skip_final_snapshot": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"master_username": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},

			"master_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"snapshot_identifier": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      27017,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1150, 65535),
			},

			"apply_immediately": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"vpc_security_group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
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

			"backup_retention_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntAtMost(35),
			},

			"kms_key_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidARN,
			},

			"cluster_resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"enabled_cloudwatch_logs_exports": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"audit",
						"profiler",
					}, false),
				},
			},

			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceClusterImport(
	d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// Neither skip_final_snapshot nor final_snapshot_identifier can be fetched
	// from any API call, so we need to default skip_final_snapshot to true so
	// that final_snapshot_identifier is not required
	d.Set("skip_final_snapshot", true)
	return []*schema.ResourceData{d}, nil
}

func resourceClusterCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).DocDBConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	// Some API calls (e.g. RestoreDBClusterFromSnapshot do not support all
	// parameters to correctly apply all settings in one pass. For missing
	// parameters or unsupported configurations, we may need to call
	// ModifyDBInstance afterwadocdb to prevent Terraform operators from API
	// errors or needing to double apply.
	var requiresModifyDbCluster bool
	modifyDbClusterInput := &docdb.ModifyDBClusterInput{
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

	if _, ok := d.GetOk("snapshot_identifier"); ok {
		opts := docdb.RestoreDBClusterFromSnapshotInput{
			DBClusterIdentifier: aws.String(identifier),
			Engine:              aws.String(d.Get("engine").(string)),
			SnapshotIdentifier:  aws.String(d.Get("snapshot_identifier").(string)),
			DeletionProtection:  aws.Bool(d.Get("deletion_protection").(bool)),
			Tags:                Tags(tags.IgnoreAWS()),
		}

		if attr := d.Get("availability_zones").(*schema.Set); attr.Len() > 0 {
			opts.AvailabilityZones = flex.ExpandStringSet(attr)
		}

		if attr, ok := d.GetOk("backup_retention_period"); ok {
			modifyDbClusterInput.BackupRetentionPeriod = aws.Int64(int64(attr.(int)))
			requiresModifyDbCluster = true
		}

		if attr, ok := d.GetOk("db_subnet_group_name"); ok {
			opts.DBSubnetGroupName = aws.String(attr.(string))
		}

		if attr, ok := d.GetOk("db_cluster_parameter_group_name"); ok {
			modifyDbClusterInput.DBClusterParameterGroupName = aws.String(attr.(string))
			requiresModifyDbCluster = true
		}

		if attr, ok := d.GetOk("enabled_cloudwatch_logs_exports"); ok && len(attr.([]interface{})) > 0 {
			opts.EnableCloudwatchLogsExports = flex.ExpandStringList(attr.([]interface{}))
		}

		if attr, ok := d.GetOk("engine_version"); ok {
			opts.EngineVersion = aws.String(attr.(string))
		}

		if attr, ok := d.GetOk("kms_key_id"); ok {
			opts.KmsKeyId = aws.String(attr.(string))
		}

		if attr, ok := d.GetOk("port"); ok {
			opts.Port = aws.Int64(int64(attr.(int)))
		}

		if attr, ok := d.GetOk("preferred_backup_window"); ok {
			modifyDbClusterInput.PreferredBackupWindow = aws.String(attr.(string))
			requiresModifyDbCluster = true
		}

		if attr, ok := d.GetOk("preferred_maintenance_window"); ok {
			modifyDbClusterInput.PreferredMaintenanceWindow = aws.String(attr.(string))
			requiresModifyDbCluster = true
		}

		if attr := d.Get("vpc_security_group_ids").(*schema.Set); attr.Len() > 0 {
			opts.VpcSecurityGroupIds = flex.ExpandStringSet(attr)
		}

		log.Printf("[DEBUG] DocDB Cluster restore from snapshot configuration: %s", opts)
		err := resource.Retry(propagationTimeout, func() *resource.RetryError {
			_, err := conn.RestoreDBClusterFromSnapshot(&opts)
			if err != nil {
				if tfawserr.ErrMessageContains(err, "InvalidParameterValue", "IAM role ARN value is invalid or does not include the required permissions") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if tfresource.TimedOut(err) {
			_, err = conn.RestoreDBClusterFromSnapshot(&opts)
		}
		if err != nil {
			return fmt.Errorf("Error creating DocDB Cluster: %s", err)
		}
	} else {
		// Secondary DocDB clusters part of a global cluster will not supply the master_password
		if _, ok := d.GetOk("global_cluster_identifier"); !ok {
			if _, ok := d.GetOk("master_password"); !ok {
				return fmt.Errorf(`provider.aws: aws_docdb_cluster: %s: "master_password": required field is not set`, identifier)
			}
		}

		// Secondary DocDB clusters part of a global cluster will not supply the master_username
		if _, ok := d.GetOk("global_cluster_identifier"); !ok {
			if _, ok := d.GetOk("master_username"); !ok {
				return fmt.Errorf(`provider.aws: aws_docdb_cluster: %s: "master_username": required field is not set`, identifier)
			}
		}

		createOpts := &docdb.CreateDBClusterInput{
			DBClusterIdentifier: aws.String(identifier),
			Engine:              aws.String(d.Get("engine").(string)),
			MasterUserPassword:  aws.String(d.Get("master_password").(string)),
			MasterUsername:      aws.String(d.Get("master_username").(string)),
			DeletionProtection:  aws.Bool(d.Get("deletion_protection").(bool)),
			Tags:                Tags(tags.IgnoreAWS()),
		}

		if attr, ok := d.GetOk("global_cluster_identifier"); ok {
			createOpts.GlobalClusterIdentifier = aws.String(attr.(string))
		}

		if attr, ok := d.GetOk("port"); ok {
			createOpts.Port = aws.Int64(int64(attr.(int)))
		}

		if attr, ok := d.GetOk("db_subnet_group_name"); ok {
			createOpts.DBSubnetGroupName = aws.String(attr.(string))
		}

		if attr, ok := d.GetOk("db_cluster_parameter_group_name"); ok {
			createOpts.DBClusterParameterGroupName = aws.String(attr.(string))
		}

		if attr, ok := d.GetOk("engine_version"); ok {
			createOpts.EngineVersion = aws.String(attr.(string))
		}

		if attr := d.Get("vpc_security_group_ids").(*schema.Set); attr.Len() > 0 {
			createOpts.VpcSecurityGroupIds = flex.ExpandStringSet(attr)
		}

		if attr := d.Get("availability_zones").(*schema.Set); attr.Len() > 0 {
			createOpts.AvailabilityZones = flex.ExpandStringSet(attr)
		}

		if v, ok := d.GetOk("backup_retention_period"); ok {
			createOpts.BackupRetentionPeriod = aws.Int64(int64(v.(int)))
		}

		if v, ok := d.GetOk("preferred_backup_window"); ok {
			createOpts.PreferredBackupWindow = aws.String(v.(string))
		}

		if v, ok := d.GetOk("preferred_maintenance_window"); ok {
			createOpts.PreferredMaintenanceWindow = aws.String(v.(string))
		}

		if attr, ok := d.GetOk("kms_key_id"); ok {
			createOpts.KmsKeyId = aws.String(attr.(string))
		}

		if attr, ok := d.GetOk("enabled_cloudwatch_logs_exports"); ok && len(attr.([]interface{})) > 0 {
			createOpts.EnableCloudwatchLogsExports = flex.ExpandStringList(attr.([]interface{}))
		}

		if attr, ok := d.GetOkExists("storage_encrypted"); ok {
			createOpts.StorageEncrypted = aws.Bool(attr.(bool))
		}

		log.Printf("[DEBUG] DocDB Cluster create options: %s", createOpts)
		var resp *docdb.CreateDBClusterOutput
		err := resource.Retry(propagationTimeout, func() *resource.RetryError {
			var err error
			resp, err = conn.CreateDBCluster(createOpts)
			if err != nil {
				if tfawserr.ErrMessageContains(err, "InvalidParameterValue", "IAM role ARN value is invalid or does not include the required permissions") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if tfresource.TimedOut(err) {
			resp, err = conn.CreateDBCluster(createOpts)
		}
		if err != nil {
			return fmt.Errorf("error creating DocDB cluster: %s", err)
		}

		log.Printf("[DEBUG]: DocDB Cluster create response: %s", resp)
	}

	d.SetId(identifier)

	log.Printf("[INFO] DocDB Cluster ID: %s", d.Id())

	log.Println(
		"[INFO] Waiting for DocDB Cluster to be available")

	stateConf := &resource.StateChangeConf{
		Pending:    resourceClusterCreatePendingStates,
		Target:     []string{"available"},
		Refresh:    resourceClusterStateRefreshFunc(conn, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second,
	}

	// Wait, catching any errors
	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for DocDB Cluster state to be \"available\": %s", err)
	}

	if requiresModifyDbCluster {
		modifyDbClusterInput.DBClusterIdentifier = aws.String(d.Id())

		log.Printf("[INFO] DocDB Cluster (%s) configuration requires ModifyDBCluster: %s", d.Id(), modifyDbClusterInput)
		_, err := conn.ModifyDBCluster(modifyDbClusterInput)
		if err != nil {
			return fmt.Errorf("error modifying DocDB Cluster (%s): %s", d.Id(), err)
		}

		log.Printf("[INFO] Waiting for DocDB Cluster (%s) to be available", d.Id())
		err = waitForClusterUpdate(conn, d.Id(), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return fmt.Errorf("error waiting for DocDB Cluster (%s) to be available: %s", d.Id(), err)
		}
	}

	return resourceClusterRead(d, meta)
}

func resourceClusterRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).DocDBConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	input := &docdb.DescribeDBClustersInput{
		DBClusterIdentifier: aws.String(d.Id()),
	}

	log.Printf("[DEBUG] Describing DocDB Cluster: %s", input)
	resp, err := conn.DescribeDBClusters(input)

	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, docdb.ErrCodeDBClusterNotFoundFault) {
		log.Printf("[WARN] DocDB Cluster (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error describing DocDB Cluster (%s): %s", d.Id(), err)
	}

	if resp == nil {
		return fmt.Errorf("Error retrieving DocDB cluster: empty response for: %s", input)
	}

	var dbc *docdb.DBCluster
	for _, c := range resp.DBClusters {
		if aws.StringValue(c.DBClusterIdentifier) == d.Id() {
			dbc = c
			break
		}
	}

	if !d.IsNewResource() && dbc == nil {
		log.Printf("[WARN] DocDB Cluster (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	globalCluster, err := findGlobalClusterByARN(context.TODO(), conn, aws.StringValue(dbc.DBClusterArn))

	// Ignore the following API error for regions/partitions that do not support DocDB Global Clusters:
	// InvalidParameterValue: Access Denied to API Version: APIGlobalDatabases
	if err != nil && !tfawserr.ErrMessageContains(err, "InvalidParameterValue", "Access Denied to API Version: APIGlobalDatabases") {
		return fmt.Errorf("error reading DocDB Global Cluster information for DB Cluster (%s): %w", d.Id(), err)
	}

	if globalCluster != nil {
		d.Set("global_cluster_identifier", globalCluster.GlobalClusterIdentifier)
	} else {
		d.Set("global_cluster_identifier", "")
	}

	if err := d.Set("availability_zones", aws.StringValueSlice(dbc.AvailabilityZones)); err != nil {
		return fmt.Errorf("error setting availability_zones: %s", err)
	}

	d.Set("arn", dbc.DBClusterArn)
	d.Set("backup_retention_period", dbc.BackupRetentionPeriod)
	d.Set("cluster_identifier", dbc.DBClusterIdentifier)

	var cm []string
	for _, m := range dbc.DBClusterMembers {
		cm = append(cm, aws.StringValue(m.DBInstanceIdentifier))
	}
	if err := d.Set("cluster_members", cm); err != nil {
		return fmt.Errorf("error setting cluster_members: %s", err)
	}

	d.Set("cluster_resource_id", dbc.DbClusterResourceId)
	d.Set("db_cluster_parameter_group_name", dbc.DBClusterParameterGroup)
	d.Set("db_subnet_group_name", dbc.DBSubnetGroup)

	if err := d.Set("enabled_cloudwatch_logs_exports", aws.StringValueSlice(dbc.EnabledCloudwatchLogsExports)); err != nil {
		return fmt.Errorf("error setting enabled_cloudwatch_logs_exports: %s", err)
	}

	d.Set("endpoint", dbc.Endpoint)
	d.Set("engine_version", dbc.EngineVersion)
	d.Set("engine", dbc.Engine)
	d.Set("hosted_zone_id", dbc.HostedZoneId)

	d.Set("kms_key_id", dbc.KmsKeyId)
	d.Set("master_username", dbc.MasterUsername)
	d.Set("port", dbc.Port)
	d.Set("preferred_backup_window", dbc.PreferredBackupWindow)
	d.Set("preferred_maintenance_window", dbc.PreferredMaintenanceWindow)
	d.Set("reader_endpoint", dbc.ReaderEndpoint)
	d.Set("storage_encrypted", dbc.StorageEncrypted)
	d.Set("deletion_protection", dbc.DeletionProtection)

	var vpcg []string
	for _, g := range dbc.VpcSecurityGroups {
		vpcg = append(vpcg, aws.StringValue(g.VpcSecurityGroupId))
	}
	if err := d.Set("vpc_security_group_ids", vpcg); err != nil {
		return fmt.Errorf("error setting vpc_security_group_ids: %s", err)
	}

	tags, err := ListTags(conn, d.Get("arn").(string))

	if err != nil {
		return fmt.Errorf("error listing tags for DocumentDB Cluster (%s): %s", d.Get("arn").(string), err)
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

func resourceClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).DocDBConn
	requestUpdate := false

	req := &docdb.ModifyDBClusterInput{
		ApplyImmediately:    aws.Bool(d.Get("apply_immediately").(bool)),
		DBClusterIdentifier: aws.String(d.Id()),
	}

	if d.HasChange("master_password") {
		req.MasterUserPassword = aws.String(d.Get("master_password").(string))
		requestUpdate = true
	}

	if d.HasChange("engine_version") {
		req.EngineVersion = aws.String(d.Get("engine_version").(string))
		requestUpdate = true
	}

	if d.HasChange("vpc_security_group_ids") {
		if attr := d.Get("vpc_security_group_ids").(*schema.Set); attr.Len() > 0 {
			req.VpcSecurityGroupIds = flex.ExpandStringSet(attr)
		} else {
			req.VpcSecurityGroupIds = []*string{}
		}
		requestUpdate = true
	}

	if d.HasChange("preferred_backup_window") {
		req.PreferredBackupWindow = aws.String(d.Get("preferred_backup_window").(string))
		requestUpdate = true
	}

	if d.HasChange("preferred_maintenance_window") {
		req.PreferredMaintenanceWindow = aws.String(d.Get("preferred_maintenance_window").(string))
		requestUpdate = true
	}

	if d.HasChange("backup_retention_period") {
		req.BackupRetentionPeriod = aws.Int64(int64(d.Get("backup_retention_period").(int)))
		requestUpdate = true
	}

	if d.HasChange("db_cluster_parameter_group_name") {
		req.DBClusterParameterGroupName = aws.String(d.Get("db_cluster_parameter_group_name").(string))
		requestUpdate = true
	}

	if d.HasChange("enabled_cloudwatch_logs_exports") {
		req.CloudwatchLogsExportConfiguration = buildCloudWatchLogsExportConfiguration(d)
		requestUpdate = true
	}

	if d.HasChange("deletion_protection") {
		req.DeletionProtection = aws.Bool(d.Get("deletion_protection").(bool))
		requestUpdate = true
	}

	if d.HasChange("global_cluster_identifier") {
		oRaw, nRaw := d.GetChange("global_cluster_identifier")
		o := oRaw.(string)
		n := nRaw.(string)

		if o == "" {
			return errors.New("existing DocDB Clusters cannot be added to an existing DocDB Global Cluster")
		}

		if n != "" {
			return errors.New("existing DocDB Clusters cannot be migrated between existing DocDB Global Clusters")
		}

		input := &docdb.RemoveFromGlobalClusterInput{
			DbClusterIdentifier:     aws.String(d.Get("arn").(string)),
			GlobalClusterIdentifier: aws.String(o),
		}

		log.Printf("[DEBUG] Removing DocDB Cluster from DocDB Global Cluster: %s", input)
		_, err := conn.RemoveFromGlobalCluster(input)

		if err != nil && !tfawserr.ErrCodeEquals(err, docdb.ErrCodeGlobalClusterNotFoundFault) && !tfawserr.ErrMessageContains(err, "InvalidParameterValue", "is not found in global cluster") {
			return fmt.Errorf("error removing DocDB Cluster (%s) from DocDB Global Cluster: %w", d.Id(), err)
		}
	}

	if requestUpdate {
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, err := conn.ModifyDBCluster(req)
			if err != nil {
				if tfawserr.ErrMessageContains(err, "InvalidParameterValue", "IAM role ARN value is invalid or does not include the required permissions") {
					return resource.RetryableError(err)
				}

				if tfawserr.ErrMessageContains(err, docdb.ErrCodeInvalidDBClusterStateFault, "is not currently in the available state") {
					return resource.RetryableError(err)
				}

				if tfawserr.ErrMessageContains(err, docdb.ErrCodeInvalidDBClusterStateFault, "DB cluster is not available for modification") {
					return resource.RetryableError(err)
				}

				return resource.NonRetryableError(err)
			}
			return nil
		})
		if tfresource.TimedOut(err) {
			_, err = conn.ModifyDBCluster(req)
		}
		if err != nil {
			return fmt.Errorf("Failed to modify DocDB Cluster (%s): %s", d.Id(), err)
		}

		log.Printf("[INFO] Waiting for DocDB Cluster (%s) to be available", d.Id())
		err = waitForClusterUpdate(conn, d.Id(), d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("error waiting for DocDB Cluster (%s) to be available: %s", d.Id(), err)
		}
	}

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := UpdateTags(conn, d.Get("arn").(string), o, n); err != nil {
			return fmt.Errorf("error updating DocumentDB Cluster (%s) tags: %s", d.Get("arn").(string), err)
		}

	}

	return resourceClusterRead(d, meta)
}

func resourceClusterDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).DocDBConn
	log.Printf("[DEBUG] Destroying DocDB Cluster (%s)", d.Id())

	// Automatically remove from global cluster to bypass this error on deletion:
	// InvalidDBClusterStateFault: This cluster is a part of a global cluster, please remove it from globalcluster first
	if d.Get("global_cluster_identifier").(string) != "" {
		input := &docdb.RemoveFromGlobalClusterInput{
			DbClusterIdentifier:     aws.String(d.Get("arn").(string)),
			GlobalClusterIdentifier: aws.String(d.Get("global_cluster_identifier").(string)),
		}

		log.Printf("[DEBUG] Removing DocDB Cluster from DocDB Global Cluster: %s", input)
		_, err := conn.RemoveFromGlobalCluster(input)

		if err != nil && !tfawserr.ErrCodeEquals(err, docdb.ErrCodeGlobalClusterNotFoundFault) && !tfawserr.ErrMessageContains(err, "InvalidParameterValue", "is not found in global cluster") {
			return fmt.Errorf("error removing DocDB Cluster (%s) from DocDB Global Cluster: %w", d.Id(), err)
		}
	}

	deleteOpts := docdb.DeleteDBClusterInput{
		DBClusterIdentifier: aws.String(d.Id()),
	}

	skipFinalSnapshot := d.Get("skip_final_snapshot").(bool)
	deleteOpts.SkipFinalSnapshot = aws.Bool(skipFinalSnapshot)

	if !skipFinalSnapshot {
		if name, present := d.GetOk("final_snapshot_identifier"); present {
			deleteOpts.FinalDBSnapshotIdentifier = aws.String(name.(string))
		} else {
			return fmt.Errorf("DocDB Cluster FinalSnapshotIdentifier is required when a final snapshot is required")
		}
	}

	log.Printf("[DEBUG] DocDB Cluster delete options: %s", deleteOpts)

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := conn.DeleteDBCluster(&deleteOpts)
		if err != nil {
			if tfawserr.ErrMessageContains(err, docdb.ErrCodeInvalidDBClusterStateFault, "is not currently in the available state") {
				return resource.RetryableError(err)
			}
			if tfawserr.ErrMessageContains(err, docdb.ErrCodeInvalidDBClusterStateFault, "cluster is a part of a global cluster") {
				return resource.RetryableError(err)
			}
			if tfawserr.ErrCodeEquals(err, docdb.ErrCodeDBClusterNotFoundFault) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if tfresource.TimedOut(err) {
		_, err = conn.DeleteDBCluster(&deleteOpts)
	}
	if err != nil {
		return fmt.Errorf("DocDB Cluster cannot be deleted: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    resourceClusterDeletePendingStates,
		Target:     []string{"destroyed"},
		Refresh:    resourceClusterStateRefreshFunc(conn, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting DocDB Cluster (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceClusterStateRefreshFunc(conn *docdb.DocDB, dbClusterIdentifier string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := conn.DescribeDBClusters(&docdb.DescribeDBClustersInput{
			DBClusterIdentifier: aws.String(dbClusterIdentifier),
		})

		if tfawserr.ErrCodeEquals(err, docdb.ErrCodeDBClusterNotFoundFault) {
			return 42, "destroyed", nil
		}

		if err != nil {
			return nil, "", err
		}

		var dbc *docdb.DBCluster

		for _, c := range resp.DBClusters {
			if aws.StringValue(c.DBClusterIdentifier) == dbClusterIdentifier {
				dbc = c
			}
		}

		if dbc == nil {
			return 42, "destroyed", nil
		}

		if dbc.Status != nil {
			log.Printf("[DEBUG] DB Cluster status (%s): %s", dbClusterIdentifier, *dbc.Status)
		}

		return dbc, aws.StringValue(dbc.Status), nil
	}
}

var resourceClusterCreatePendingStates = []string{
	"creating",
	"backing-up",
	"modifying",
	"preparing-data-migration",
	"migrating",
	"resetting-master-credentials",
}

var resourceClusterDeletePendingStates = []string{
	"available",
	"deleting",
	"backing-up",
	"modifying",
}

var resourceClusterUpdatePendingStates = []string{
	"backing-up",
	"modifying",
	"resetting-master-credentials",
	"upgrading",
}

func waitForClusterUpdate(conn *docdb.DocDB, id string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:    resourceClusterUpdatePendingStates,
		Target:     []string{"available"},
		Refresh:    resourceClusterStateRefreshFunc(conn, id),
		Timeout:    timeout,
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second, // Wait 30 secs before starting
	}
	_, err := stateConf.WaitForState()
	return err
}

func buildCloudWatchLogsExportConfiguration(d *schema.ResourceData) *docdb.CloudwatchLogsExportConfiguration {

	oraw, nraw := d.GetChange("enabled_cloudwatch_logs_exports")
	o := oraw.([]interface{})
	n := nraw.([]interface{})

	create, disable := diffCloudWatchLogsExportConfiguration(o, n)

	return &docdb.CloudwatchLogsExportConfiguration{
		EnableLogTypes:  flex.ExpandStringList(create),
		DisableLogTypes: flex.ExpandStringList(disable),
	}
}

func diffCloudWatchLogsExportConfiguration(old, new []interface{}) ([]interface{}, []interface{}) {
	add := make([]interface{}, 0)
	disable := make([]interface{}, 0)

	for _, n := range new {
		if _, contains := verify.SliceContainsString(old, n.(string)); !contains {
			add = append(add, n)
		}
	}

	for _, o := range old {
		if _, contains := verify.SliceContainsString(new, o.(string)); !contains {
			disable = append(disable, o)
		}
	}

	return add, disable
}
