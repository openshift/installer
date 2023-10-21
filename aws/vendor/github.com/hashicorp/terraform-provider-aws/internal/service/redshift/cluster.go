package redshift

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/redshift"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_redshift_cluster", name="Cluster")
// @Tags(identifierAttribute="arn")
func ResourceCluster() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceClusterCreate,
		ReadWithoutTimeout:   resourceClusterRead,
		UpdateWithoutTimeout: resourceClusterUpdate,
		DeleteWithoutTimeout: resourceClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceClusterImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(75 * time.Minute),
			Update: schema.DefaultTimeout(75 * time.Minute),
			Delete: schema.DefaultTimeout(40 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"allow_version_upgrade": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"apply_immediately": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"aqua_configuration_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(redshift.AquaConfigurationStatus_Values(), false),
				Deprecated:   "This parameter is no longer supported by the AWS API. It will be removed in the next major version of the provider.",
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					return true
				},
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"automated_snapshot_retention_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntAtMost(35),
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"availability_zone_relocation_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"cluster_identifier": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[0-9a-z-]+$`), "must contain only lowercase alphanumeric characters and hyphens"),
					validation.StringMatch(regexp.MustCompile(`(?i)^[a-z]`), "first character must be a letter"),
					validation.StringDoesNotMatch(regexp.MustCompile(`--`), "cannot contain two consecutive hyphens"),
					validation.StringDoesNotMatch(regexp.MustCompile(`-$`), "cannot end with a hyphen"),
				),
			},
			"cluster_namespace_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"cluster_parameter_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_public_key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_revision_number": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_subnet_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"cluster_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "1.0",
			},
			"database_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 64),
					validation.StringMatch(regexp.MustCompile(`^[0-9a-z_$]+$`), "must contain only lowercase alphanumeric characters, underscores, and dollar signs"),
					validation.StringMatch(regexp.MustCompile(`(?i)^[a-z_]`), "first character must be a letter or underscore"),
				),
			},
			"default_iam_role_arn": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: verify.ValidARN,
			},
			"dns_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"elastic_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encrypted": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enhanced_vpc_routing": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"final_snapshot_identifier": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 255),
					validation.StringMatch(regexp.MustCompile(`^[0-9A-Za-z-]+$`), "must only contain alphanumeric characters and hyphens"),
					validation.StringDoesNotMatch(regexp.MustCompile(`--`), "cannot contain two consecutive hyphens"),
					validation.StringDoesNotMatch(regexp.MustCompile(`-$`), "cannot end in a hyphen"),
				),
			},
			"iam_roles": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: verify.ValidARN,
				},
			},
			"kms_key_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: verify.ValidARN,
			},
			"logging": {
				Type:             schema.TypeList,
				MaxItems:         1,
				Optional:         true,
				DiffSuppressFunc: verify.SuppressMissingOptionalConfigurationBlock,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"enable": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"log_destination_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(redshift.LogDestinationType_Values(), false),
						},
						"log_exports": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"s3_key_prefix": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"maintenance_track_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "current",
			},
			"manual_snapshot_retention_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      -1,
				ValidateFunc: validation.IntBetween(-1, 3653),
			},
			"master_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(8, 64),
					validation.StringMatch(regexp.MustCompile(`^.*[a-z].*`), "must contain at least one lowercase letter"),
					validation.StringMatch(regexp.MustCompile(`^.*[A-Z].*`), "must contain at least one uppercase letter"),
					validation.StringMatch(regexp.MustCompile(`^.*[0-9].*`), "must contain at least one number"),
					validation.StringMatch(regexp.MustCompile(`^[^\@\/'" ]*$`), "cannot contain [/@\"' ]"),
				),
			},
			"master_username": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 128),
					validation.StringMatch(regexp.MustCompile(`^\w+$`), "must contain only alphanumeric characters"),
					validation.StringMatch(regexp.MustCompile(`(?i)^[a-z_]`), "first character must be a letter"),
				),
			},
			"node_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"number_of_nodes": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"owner_account": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: verify.ValidAccountID,
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5439,
				ValidateFunc: validation.IntBetween(1115, 65535),
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
			"publicly_accessible": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"skip_final_snapshot": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"snapshot_cluster_identifier": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"snapshot_copy": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination_region": {
							Type:     schema.TypeString,
							Required: true,
						},
						"grant_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"retention_period": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  7,
						},
					},
				},
			},
			"snapshot_identifier": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
			"vpc_security_group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},

		CustomizeDiff: customdiff.All(
			verify.SetTagsDiff,
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				if diff.Id() == "" {
					return nil
				}
				if diff.Get("availability_zone_relocation_enabled").(bool) {
					return nil
				}
				o, n := diff.GetChange("availability_zone")
				if o.(string) != n.(string) {
					return fmt.Errorf("cannot change `availability_zone` if `availability_zone_relocation_enabled` is not true")
				}
				return nil
			},
		),
	}
}

func resourceClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).RedshiftConn(ctx)

	clusterID := d.Get("cluster_identifier").(string)
	backupInput := &redshift.RestoreFromClusterSnapshotInput{
		AllowVersionUpgrade:              aws.Bool(d.Get("allow_version_upgrade").(bool)),
		AutomatedSnapshotRetentionPeriod: aws.Int64(int64(d.Get("automated_snapshot_retention_period").(int))),
		ClusterIdentifier:                aws.String(clusterID),
		Port:                             aws.Int64(int64(d.Get("port").(int))),
		NodeType:                         aws.String(d.Get("node_type").(string)),
		PubliclyAccessible:               aws.Bool(d.Get("publicly_accessible").(bool)),
	}

	input := &redshift.CreateClusterInput{
		AllowVersionUpgrade:              aws.Bool(d.Get("allow_version_upgrade").(bool)),
		AutomatedSnapshotRetentionPeriod: aws.Int64(int64(d.Get("automated_snapshot_retention_period").(int))),
		ClusterIdentifier:                aws.String(clusterID),
		ClusterVersion:                   aws.String(d.Get("cluster_version").(string)),
		DBName:                           aws.String(d.Get("database_name").(string)),
		MasterUsername:                   aws.String(d.Get("master_username").(string)),
		MasterUserPassword:               aws.String(d.Get("master_password").(string)),
		NodeType:                         aws.String(d.Get("node_type").(string)),
		Port:                             aws.Int64(int64(d.Get("port").(int))),
		PubliclyAccessible:               aws.Bool(d.Get("publicly_accessible").(bool)),
		Tags:                             GetTagsIn(ctx),
	}

	if v, ok := d.GetOk("aqua_configuration_status"); ok {
		backupInput.AquaConfigurationStatus = aws.String(v.(string))
		input.AquaConfigurationStatus = aws.String(v.(string))
	}

	if v, ok := d.GetOk("availability_zone"); ok {
		backupInput.AvailabilityZone = aws.String(v.(string))
		input.AvailabilityZone = aws.String(v.(string))
	}

	if v, ok := d.GetOk("availability_zone_relocation_enabled"); ok {
		backupInput.AvailabilityZoneRelocation = aws.Bool(v.(bool))
		input.AvailabilityZoneRelocation = aws.Bool(v.(bool))
	}

	if v, ok := d.GetOk("cluster_parameter_group_name"); ok {
		backupInput.ClusterParameterGroupName = aws.String(v.(string))
		input.ClusterParameterGroupName = aws.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_subnet_group_name"); ok {
		backupInput.ClusterSubnetGroupName = aws.String(v.(string))
		input.ClusterSubnetGroupName = aws.String(v.(string))
	}

	if v, ok := d.GetOk("default_iam_role_arn"); ok {
		backupInput.DefaultIamRoleArn = aws.String(v.(string))
		input.DefaultIamRoleArn = aws.String(v.(string))
	}

	if v, ok := d.GetOk("elastic_ip"); ok {
		backupInput.ElasticIp = aws.String(v.(string))
		input.ElasticIp = aws.String(v.(string))
	}

	if v, ok := d.GetOk("enhanced_vpc_routing"); ok {
		backupInput.EnhancedVpcRouting = aws.Bool(v.(bool))
		input.EnhancedVpcRouting = aws.Bool(v.(bool))
	}

	if v, ok := d.GetOk("iam_roles"); ok {
		backupInput.IamRoles = flex.ExpandStringSet(v.(*schema.Set))
		input.IamRoles = flex.ExpandStringSet(v.(*schema.Set))
	}

	if v, ok := d.GetOk("kms_key_id"); ok {
		backupInput.KmsKeyId = aws.String(v.(string))
		input.KmsKeyId = aws.String(v.(string))
	}

	if v, ok := d.GetOk("maintenance_track_name"); ok {
		backupInput.MaintenanceTrackName = aws.String(v.(string))
		input.MaintenanceTrackName = aws.String(v.(string))
	}

	if v, ok := d.GetOk("manual_snapshot_retention_period"); ok {
		backupInput.ManualSnapshotRetentionPeriod = aws.Int64(int64(v.(int)))
		input.ManualSnapshotRetentionPeriod = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("number_of_nodes"); ok {
		backupInput.NumberOfNodes = aws.Int64(int64(v.(int)))
		// NumberOfNodes set below for CreateCluster.
	}

	if v, ok := d.GetOk("preferred_maintenance_window"); ok {
		backupInput.PreferredMaintenanceWindow = aws.String(v.(string))
		input.PreferredMaintenanceWindow = aws.String(v.(string))
	}

	if v := d.Get("vpc_security_group_ids").(*schema.Set); v.Len() > 0 {
		backupInput.VpcSecurityGroupIds = flex.ExpandStringSet(v)
		input.VpcSecurityGroupIds = flex.ExpandStringSet(v)
	}

	if v, ok := d.GetOk("snapshot_identifier"); ok {
		backupInput.SnapshotIdentifier = aws.String(v.(string))

		if v, ok := d.GetOk("owner_account"); ok {
			backupInput.OwnerAccount = aws.String(v.(string))
		}

		if v, ok := d.GetOk("snapshot_cluster_identifier"); ok {
			backupInput.SnapshotClusterIdentifier = aws.String(v.(string))
		}

		log.Printf("[DEBUG] Restoring Redshift Cluster: %s", backupInput)
		output, err := conn.RestoreFromClusterSnapshotWithContext(ctx, backupInput)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "restoring Redshift Cluster (%s) from snapshot: %s", clusterID, err)
		}

		d.SetId(aws.StringValue(output.Cluster.ClusterIdentifier))
	} else {
		if _, ok := d.GetOk("master_password"); !ok {
			return sdkdiag.AppendErrorf(diags, `provider.aws: aws_redshift_cluster: %s: "master_password": required field is not set`, d.Get("cluster_identifier").(string))
		}

		if _, ok := d.GetOk("master_username"); !ok {
			return sdkdiag.AppendErrorf(diags, `provider.aws: aws_redshift_cluster: %s: "master_username": required field is not set`, d.Get("cluster_identifier").(string))
		}

		if v, ok := d.GetOk("encrypted"); ok {
			input.Encrypted = aws.Bool(v.(bool))
		}

		if v := d.Get("number_of_nodes").(int); v > 1 {
			input.ClusterType = aws.String(clusterTypeMultiNode)
			input.NumberOfNodes = aws.Int64(int64(d.Get("number_of_nodes").(int)))
		} else {
			input.ClusterType = aws.String(clusterTypeSingleNode)
		}

		output, err := conn.CreateClusterWithContext(ctx, input)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "creating Redshift Cluster (%s): %s", clusterID, err)
		}

		d.SetId(aws.StringValue(output.Cluster.ClusterIdentifier))
	}

	if _, err := waitClusterCreated(ctx, conn, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
		return sdkdiag.AppendErrorf(diags, "creating Redshift Cluster (%s): waiting for completion: %s", d.Id(), err)
	}

	if _, err := waitClusterRelocationStatusResolved(ctx, conn, d.Id()); err != nil {
		return sdkdiag.AppendErrorf(diags, "creating Redshift Cluster (%s): waiting for relocation: %s", d.Id(), err)
	}

	if v, ok := d.GetOk("snapshot_copy"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		if err := enableSnapshotCopy(ctx, conn, d.Id(), v.([]interface{})[0].(map[string]interface{})); err != nil {
			return sdkdiag.AppendErrorf(diags, "creating Redshift Cluster (%s): %s", d.Id(), err)
		}
	}

	if v, ok := d.GetOk("logging"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		tfMap := v.([]interface{})[0].(map[string]interface{})

		if v, ok := tfMap["enable"].(bool); ok && v {
			if err := enableLogging(ctx, conn, d.Id(), tfMap); err != nil {
				return sdkdiag.AppendErrorf(diags, "creating Redshift Cluster (%s): %s", d.Id(), err)
			}
		}
	}

	return append(diags, resourceClusterRead(ctx, d, meta)...)
}

func resourceClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).RedshiftConn(ctx)

	rsc, err := FindClusterByID(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Redshift Cluster (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Redshift Cluster (%s): %s", d.Id(), err)
	}

	loggingStatus, err := conn.DescribeLoggingStatusWithContext(ctx, &redshift.DescribeLoggingStatusInput{
		ClusterIdentifier: aws.String(d.Id()),
	})

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Redshift Cluster (%s) logging status: %s", d.Id(), err)
	}

	d.Set("allow_version_upgrade", rsc.AllowVersionUpgrade)
	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Service:   redshift.ServiceName,
		Region:    meta.(*conns.AWSClient).Region,
		AccountID: meta.(*conns.AWSClient).AccountID,
		Resource:  fmt.Sprintf("cluster:%s", d.Id()),
	}.String()
	d.Set("arn", arn)
	if rsc.AquaConfiguration != nil {
		d.Set("aqua_configuration_status", rsc.AquaConfiguration.AquaConfigurationStatus)
	}
	d.Set("automated_snapshot_retention_period", rsc.AutomatedSnapshotRetentionPeriod)
	d.Set("availability_zone", rsc.AvailabilityZone)
	azr, err := clusterAvailabilityZoneRelocationStatus(rsc)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Redshift Cluster (%s): %s", d.Id(), err)
	}
	d.Set("availability_zone_relocation_enabled", azr)
	d.Set("cluster_identifier", rsc.ClusterIdentifier)
	d.Set("cluster_namespace_arn", rsc.ClusterNamespaceArn)
	if err := d.Set("cluster_nodes", flattenClusterNodes(rsc.ClusterNodes)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting cluster_nodes: %s", err)
	}
	d.Set("cluster_parameter_group_name", rsc.ClusterParameterGroups[0].ParameterGroupName)
	d.Set("cluster_public_key", rsc.ClusterPublicKey)
	d.Set("cluster_revision_number", rsc.ClusterRevisionNumber)
	d.Set("cluster_subnet_group_name", rsc.ClusterSubnetGroupName)
	if len(rsc.ClusterNodes) > 1 {
		d.Set("cluster_type", clusterTypeMultiNode)
	} else {
		d.Set("cluster_type", clusterTypeSingleNode)
	}
	d.Set("cluster_version", rsc.ClusterVersion)
	d.Set("database_name", rsc.DBName)
	d.Set("default_iam_role_arn", rsc.DefaultIamRoleArn)
	d.Set("encrypted", rsc.Encrypted)
	d.Set("enhanced_vpc_routing", rsc.EnhancedVpcRouting)
	d.Set("kms_key_id", rsc.KmsKeyId)
	if err := d.Set("logging", flattenLogging(loggingStatus)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting logging: %s", err)
	}
	d.Set("maintenance_track_name", rsc.MaintenanceTrackName)
	d.Set("manual_snapshot_retention_period", rsc.ManualSnapshotRetentionPeriod)
	d.Set("master_username", rsc.MasterUsername)
	d.Set("node_type", rsc.NodeType)
	d.Set("number_of_nodes", rsc.NumberOfNodes)
	d.Set("preferred_maintenance_window", rsc.PreferredMaintenanceWindow)
	d.Set("publicly_accessible", rsc.PubliclyAccessible)
	if err := d.Set("snapshot_copy", flattenSnapshotCopy(rsc.ClusterSnapshotCopyStatus)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting snapshot_copy: %s", err)
	}

	d.Set("dns_name", nil)
	d.Set("endpoint", nil)
	d.Set("port", nil)
	if endpoint := rsc.Endpoint; endpoint != nil {
		if address := aws.StringValue(endpoint.Address); address != "" {
			d.Set("dns_name", address)
			if port := aws.Int64Value(endpoint.Port); port != 0 {
				d.Set("endpoint", fmt.Sprintf("%s:%d", address, port))
				d.Set("port", port)
			} else {
				d.Set("endpoint", address)
			}
		}
	}

	var apiList []*string

	for _, iamRole := range rsc.IamRoles {
		apiList = append(apiList, iamRole.IamRoleArn)
	}
	d.Set("iam_roles", aws.StringValueSlice(apiList))

	apiList = nil

	for _, vpcSecurityGroup := range rsc.VpcSecurityGroups {
		apiList = append(apiList, vpcSecurityGroup.VpcSecurityGroupId)
	}
	d.Set("vpc_security_group_ids", aws.StringValueSlice(apiList))

	SetTagsOut(ctx, rsc.Tags)

	return diags
}

func resourceClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).RedshiftConn(ctx)

	if d.HasChangesExcept("aqua_configuration_status", "availability_zone", "iam_roles", "logging", "snapshot_copy", "tags", "tags_all") {
		input := &redshift.ModifyClusterInput{
			ClusterIdentifier: aws.String(d.Id()),
		}

		if d.HasChange("allow_version_upgrade") {
			input.AllowVersionUpgrade = aws.Bool(d.Get("allow_version_upgrade").(bool))
		}

		if d.HasChange("automated_snapshot_retention_period") {
			input.AutomatedSnapshotRetentionPeriod = aws.Int64(int64(d.Get("automated_snapshot_retention_period").(int)))
		}

		if d.HasChange("availability_zone_relocation_enabled") {
			input.AvailabilityZoneRelocation = aws.Bool(d.Get("availability_zone_relocation_enabled").(bool))
		}

		if d.HasChange("cluster_parameter_group_name") {
			input.ClusterParameterGroupName = aws.String(d.Get("cluster_parameter_group_name").(string))
		}

		if d.HasChange("maintenance_track_name") {
			input.MaintenanceTrackName = aws.String(d.Get("maintenance_track_name").(string))
		}

		if d.HasChange("manual_snapshot_retention_period") {
			input.ManualSnapshotRetentionPeriod = aws.Int64(int64(d.Get("manual_snapshot_retention_period").(int)))
		}

		// If the cluster type, node type, or number of nodes changed, then the AWS API expects all three
		// items to be sent over.
		if d.HasChanges("cluster_type", "node_type", "number_of_nodes") {
			input.NodeType = aws.String(d.Get("node_type").(string))

			if v := d.Get("number_of_nodes").(int); v > 1 {
				input.ClusterType = aws.String(clusterTypeMultiNode)
				input.NumberOfNodes = aws.Int64(int64(d.Get("number_of_nodes").(int)))
			} else {
				input.ClusterType = aws.String(clusterTypeSingleNode)
			}
		}

		if d.HasChange("cluster_version") {
			input.ClusterVersion = aws.String(d.Get("cluster_version").(string))
		}

		if d.HasChange("encrypted") {
			input.Encrypted = aws.Bool(d.Get("encrypted").(bool))
		}

		if d.HasChange("enhanced_vpc_routing") {
			input.EnhancedVpcRouting = aws.Bool(d.Get("enhanced_vpc_routing").(bool))
		}

		if d.Get("encrypted").(bool) && d.HasChange("kms_key_id") {
			input.KmsKeyId = aws.String(d.Get("kms_key_id").(string))
		}

		if d.HasChange("master_password") {
			input.MasterUserPassword = aws.String(d.Get("master_password").(string))
		}

		if d.HasChange("preferred_maintenance_window") {
			input.PreferredMaintenanceWindow = aws.String(d.Get("preferred_maintenance_window").(string))
		}

		if d.HasChange("publicly_accessible") {
			input.PubliclyAccessible = aws.Bool(d.Get("publicly_accessible").(bool))
		}

		if d.HasChange("vpc_security_group_ids") {
			input.VpcSecurityGroupIds = flex.ExpandStringSet(d.Get("vpc_security_group_ids").(*schema.Set))
		}

		_, err := conn.ModifyClusterWithContext(ctx, input)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "modifying Redshift Cluster (%s): %s", d.Id(), err)
		}

		if _, err := waitClusterUpdated(ctx, conn, d.Id(), d.Timeout(schema.TimeoutUpdate)); err != nil {
			return sdkdiag.AppendErrorf(diags, "waiting for Redshift Cluster (%s) update: %s", d.Id(), err)
		}

		if _, err := waitClusterRelocationStatusResolved(ctx, conn, d.Id()); err != nil {
			return sdkdiag.AppendErrorf(diags, "waiting for Redshift Cluster (%s) Availability Zone Relocation Status resolution: %s", d.Id(), err)
		}
	}

	if d.HasChanges("iam_roles", "default_iam_role_arn") {
		o, n := d.GetChange("iam_roles")
		if o == nil {
			o = new(schema.Set)
		}
		if n == nil {
			n = new(schema.Set)
		}

		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		add := ns.Difference(os)
		del := os.Difference(ns)

		input := &redshift.ModifyClusterIamRolesInput{
			AddIamRoles:       flex.ExpandStringSet(add),
			ClusterIdentifier: aws.String(d.Id()),
			RemoveIamRoles:    flex.ExpandStringSet(del),
			DefaultIamRoleArn: aws.String(d.Get("default_iam_role_arn").(string)),
		}

		log.Printf("[DEBUG] Modifying Redshift Cluster IAM Roles: %s", input)
		_, err := conn.ModifyClusterIamRolesWithContext(ctx, input)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "modifying Redshift Cluster (%s) IAM roles: %s", d.Id(), err)
		}

		if _, err := waitClusterUpdated(ctx, conn, d.Id(), d.Timeout(schema.TimeoutUpdate)); err != nil {
			return sdkdiag.AppendErrorf(diags, "waiting for Redshift Cluster (%s) update: %s", d.Id(), err)
		}
	}

	if d.HasChange("aqua_configuration_status") {
		input := &redshift.ModifyAquaConfigurationInput{
			AquaConfigurationStatus: aws.String(d.Get("aqua_configuration_status").(string)),
			ClusterIdentifier:       aws.String(d.Id()),
		}

		log.Printf("[DEBUG] Modifying Redshift Cluster Aqua Configuration: %s", input)
		_, err := conn.ModifyAquaConfigurationWithContext(ctx, input)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "modifying Redshift Cluster (%s) Aqua Configuration: %s", d.Id(), err)
		}

		if d.Get("apply_immediately").(bool) {
			rebootInput := &redshift.RebootClusterInput{
				ClusterIdentifier: aws.String(d.Id()),
			}

			_, err := tfresource.RetryWhenAWSErrCodeEquals(ctx, clusterInvalidClusterStateFaultTimeout,
				func() (interface{}, error) {
					return conn.RebootClusterWithContext(ctx, rebootInput)
				},
				redshift.ErrCodeInvalidClusterStateFault,
			)

			if err != nil {
				return sdkdiag.AppendErrorf(diags, "rebooting Redshift Cluster (%s): %s", d.Id(), err)
			}

			if _, err := waitClusterRebooted(ctx, conn, d.Id(), d.Timeout(schema.TimeoutUpdate)); err != nil {
				return sdkdiag.AppendErrorf(diags, "waiting for Redshift Cluster (%s) Rebooted: %s", d.Id(), err)
			}

			if _, err := waitClusterAquaApplied(ctx, conn, d.Id(), d.Timeout(schema.TimeoutUpdate)); err != nil {
				return sdkdiag.AppendErrorf(diags, "waiting for Redshift Cluster (%s) Aqua Configuration update: %s", d.Id(), err)
			}
		}
	}

	// Availability Zone cannot be changed at the same time as other settings
	if d.HasChange("availability_zone") {
		input := &redshift.ModifyClusterInput{
			AvailabilityZone:  aws.String(d.Get("availability_zone").(string)),
			ClusterIdentifier: aws.String(d.Id()),
		}

		log.Printf("[DEBUG] Relocating Redshift Cluster: %s", input)
		_, err := conn.ModifyClusterWithContext(ctx, input)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "relocating Redshift Cluster (%s): %s", d.Id(), err)
		}

		if _, err := waitClusterUpdated(ctx, conn, d.Id(), d.Timeout(schema.TimeoutUpdate)); err != nil {
			return sdkdiag.AppendErrorf(diags, "waiting for Redshift Cluster (%s) update: %s", d.Id(), err)
		}
	}

	if d.HasChange("snapshot_copy") {
		if v, ok := d.GetOk("snapshot_copy"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
			if err := enableSnapshotCopy(ctx, conn, d.Id(), v.([]interface{})[0].(map[string]interface{})); err != nil {
				return sdkdiag.AppendErrorf(diags, "updating Redshift Cluster (%s): %s", d.Id(), err)
			}
		} else {
			if err := disableSnapshotCopy(ctx, conn, d.Id()); err != nil {
				return sdkdiag.AppendErrorf(diags, "updating Redshift Cluster (%s): %s", d.Id(), err)
			}
		}
	}

	if d.HasChange("logging") {
		if v, ok := d.GetOk("logging"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
			tfMap := v.([]interface{})[0].(map[string]interface{})

			if v, ok := tfMap["enable"].(bool); ok && v {
				if err := enableLogging(ctx, conn, d.Id(), tfMap); err != nil {
					return sdkdiag.AppendErrorf(diags, "updating Redshift Cluster (%s): %s", d.Id(), err)
				}
			} else {
				if err := disableLogging(ctx, conn, d.Id()); err != nil {
					return sdkdiag.AppendErrorf(diags, "updating Redshift Cluster (%s): %s", d.Id(), err)
				}
			}
		}
	}

	return append(diags, resourceClusterRead(ctx, d, meta)...)
}

func resourceClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).RedshiftConn(ctx)

	skipFinalSnapshot := d.Get("skip_final_snapshot").(bool)
	input := &redshift.DeleteClusterInput{
		ClusterIdentifier:        aws.String(d.Id()),
		SkipFinalClusterSnapshot: aws.Bool(skipFinalSnapshot),
	}

	if !skipFinalSnapshot {
		if v, ok := d.GetOk("final_snapshot_identifier"); ok {
			input.FinalClusterSnapshotIdentifier = aws.String(v.(string))
		} else {
			return sdkdiag.AppendErrorf(diags, "Redshift Cluster Instance FinalSnapshotIdentifier is required when a final snapshot is required")
		}
	}

	log.Printf("[DEBUG] Deleting Redshift Cluster: %s", d.Id())
	_, err := tfresource.RetryWhenAWSErrCodeEquals(ctx, clusterInvalidClusterStateFaultTimeout,
		func() (interface{}, error) {
			return conn.DeleteClusterWithContext(ctx, input)
		},
		redshift.ErrCodeInvalidClusterStateFault,
	)

	if tfawserr.ErrCodeEquals(err, redshift.ErrCodeClusterNotFoundFault) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting Redshift Cluster (%s): %s", d.Id(), err)
	}

	if _, err := waitClusterDeleted(ctx, conn, d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for Redshift Cluster (%s) delete: %s", d.Id(), err)
	}

	return diags
}

func resourceClusterImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// Neither skip_final_snapshot nor final_snapshot_identifier can be fetched
	// from any API call, so we need to default skip_final_snapshot to true so
	// that final_snapshot_identifier is not required.
	d.Set("skip_final_snapshot", true)

	return []*schema.ResourceData{d}, nil
}

func enableLogging(ctx context.Context, conn *redshift.Redshift, clusterID string, tfMap map[string]interface{}) error {
	input := &redshift.EnableLoggingInput{
		ClusterIdentifier: aws.String(clusterID),
	}

	if v, ok := tfMap["bucket_name"].(string); ok && v != "" {
		input.BucketName = aws.String(v)
	}

	if v, ok := tfMap["log_destination_type"].(string); ok && v != "" {
		input.LogDestinationType = aws.String(v)
	}

	if v, ok := tfMap["log_exports"].(*schema.Set); ok && v.Len() > 0 {
		input.LogExports = flex.ExpandStringSet(v)
	}

	if v, ok := tfMap["s3_key_prefix"].(string); ok && v != "" {
		input.S3KeyPrefix = aws.String(v)
	}

	_, err := tfresource.RetryWhenAWSErrCodeEquals(ctx, clusterInvalidClusterStateFaultTimeout,
		func() (interface{}, error) {
			return conn.EnableLoggingWithContext(ctx, input)
		},
		redshift.ErrCodeInvalidClusterStateFault,
	)

	if err != nil {
		return fmt.Errorf("enabling logging: %w", err)
	}

	return nil
}

func disableLogging(ctx context.Context, conn *redshift.Redshift, clusterID string) error {
	input := &redshift.DisableLoggingInput{
		ClusterIdentifier: aws.String(clusterID),
	}

	_, err := tfresource.RetryWhenAWSErrCodeEquals(ctx, clusterInvalidClusterStateFaultTimeout,
		func() (interface{}, error) {
			return conn.DisableLoggingWithContext(ctx, input)
		},
		redshift.ErrCodeInvalidClusterStateFault,
	)

	if err != nil {
		return fmt.Errorf("disabling logging: %w", err)
	}

	return nil
}

func enableSnapshotCopy(ctx context.Context, conn *redshift.Redshift, clusterID string, tfMap map[string]interface{}) error {
	input := &redshift.EnableSnapshotCopyInput{
		ClusterIdentifier: aws.String(clusterID),
		DestinationRegion: aws.String(tfMap["destination_region"].(string)),
	}

	if v, ok := tfMap["retention_period"]; ok {
		input.RetentionPeriod = aws.Int64(int64(v.(int)))
	}

	if v, ok := tfMap["grant_name"]; ok {
		input.SnapshotCopyGrantName = aws.String(v.(string))
	}

	_, err := conn.EnableSnapshotCopyWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("enabling snapshot copy: %w", err)
	}

	return nil
}

func disableSnapshotCopy(ctx context.Context, conn *redshift.Redshift, clusterID string) error {
	input := &redshift.DisableSnapshotCopyInput{
		ClusterIdentifier: aws.String(clusterID),
	}

	_, err := conn.DisableSnapshotCopyWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("disabling snapshot copy: %w", err)
	}

	return nil
}

func flattenClusterNode(apiObject *redshift.ClusterNode) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.NodeRole; v != nil {
		tfMap["node_role"] = aws.StringValue(v)
	}

	if v := apiObject.PrivateIPAddress; v != nil {
		tfMap["private_ip_address"] = aws.StringValue(v)
	}

	if v := apiObject.PublicIPAddress; v != nil {
		tfMap["public_ip_address"] = aws.StringValue(v)
	}

	return tfMap
}

func flattenClusterNodes(apiObjects []*redshift.ClusterNode) []interface{} {
	if len(apiObjects) == 0 {
		return nil
	}

	var tfList []interface{}

	for _, apiObject := range apiObjects {
		if apiObject == nil {
			continue
		}

		tfList = append(tfList, flattenClusterNode(apiObject))
	}

	return tfList
}

func clusterAvailabilityZoneRelocationStatus(cluster *redshift.Cluster) (bool, error) {
	// AvailabilityZoneRelocation is not returned by the API, and AvailabilityZoneRelocationStatus is not implemented as Const at this time.
	switch availabilityZoneRelocationStatus := aws.StringValue(cluster.AvailabilityZoneRelocationStatus); availabilityZoneRelocationStatus {
	case "enabled":
		return true, nil
	case "disabled":
		return false, nil
	default:
		return false, fmt.Errorf("unexpected AvailabilityZoneRelocationStatus value %q returned by API", availabilityZoneRelocationStatus)
	}
}
