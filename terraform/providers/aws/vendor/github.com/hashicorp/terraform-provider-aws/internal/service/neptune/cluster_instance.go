package neptune

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/neptune"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_neptune_cluster_instance", name="Cluster")
// @Tags(identifierAttribute="arn")
func ResourceClusterInstance() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceClusterInstanceCreate,
		ReadWithoutTimeout:   resourceClusterInstanceRead,
		UpdateWithoutTimeout: resourceClusterInstanceUpdate,
		DeleteWithoutTimeout: resourceClusterInstanceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"apply_immediately": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_minor_version_upgrade": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"cluster_identifier": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dbi_resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "neptune",
				ValidateFunc: validEngine(),
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"identifier": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"identifier_prefix"},
				ValidateFunc:  validIdentifier,
			},
			"identifier_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"identifier"},
				ValidateFunc:  validIdentifierPrefix,
			},
			"instance_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"kms_key_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"neptune_parameter_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default.neptune1",
			},
			"neptune_subnet_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  DefaultPort,
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
			"promotion_tier": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"publicly_accessible": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"storage_encrypted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
			"writer": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceClusterInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).NeptuneConn(ctx)

	var instanceID string
	if v, ok := d.GetOk("identifier"); ok {
		instanceID = v.(string)
	} else if v, ok := d.GetOk("identifier_prefix"); ok {
		instanceID = id.PrefixedUniqueId(v.(string))
	} else {
		instanceID = id.PrefixedUniqueId("tf-")
	}

	input := &neptune.CreateDBInstanceInput{
		AutoMinorVersionUpgrade: aws.Bool(d.Get("auto_minor_version_upgrade").(bool)),
		DBClusterIdentifier:     aws.String(d.Get("cluster_identifier").(string)),
		DBInstanceClass:         aws.String(d.Get("instance_class").(string)),
		DBInstanceIdentifier:    aws.String(instanceID),
		Engine:                  aws.String(d.Get("engine").(string)),
		PromotionTier:           aws.Int64(int64(d.Get("promotion_tier").(int))),
		PubliclyAccessible:      aws.Bool(d.Get("publicly_accessible").(bool)),
		Tags:                    GetTagsIn(ctx),
	}

	if v, ok := d.GetOk("availability_zone"); ok {
		input.AvailabilityZone = aws.String(v.(string))
	}

	if v, ok := d.GetOk("engine_version"); ok {
		input.EngineVersion = aws.String(v.(string))
	}

	if v, ok := d.GetOk("neptune_parameter_group_name"); ok {
		input.DBParameterGroupName = aws.String(v.(string))
	}

	if v, ok := d.GetOk("neptune_subnet_group_name"); ok {
		input.DBSubnetGroupName = aws.String(v.(string))
	}

	if v, ok := d.GetOk("preferred_backup_window"); ok {
		input.PreferredBackupWindow = aws.String(v.(string))
	}

	if v, ok := d.GetOk("preferred_maintenance_window"); ok {
		input.PreferredMaintenanceWindow = aws.String(v.(string))
	}

	outputRaw, err := tfresource.RetryWhenAWSErrMessageContains(ctx, propagationTimeout, func() (interface{}, error) {
		return conn.CreateDBInstanceWithContext(ctx, input)
	}, "InvalidParameterValue", "IAM role ARN value is invalid or does not include the required permissions")

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating Neptune Cluster Instance (%s): %s", instanceID, err)
	}

	d.SetId(aws.StringValue(outputRaw.(*neptune.CreateDBInstanceOutput).DBInstance.DBInstanceIdentifier))

	if _, err := waitClusterInstanceAvailable(ctx, conn, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for Neptune Cluster Instance (%s) create: %s", d.Id(), err)
	}

	return append(diags, resourceClusterInstanceRead(ctx, d, meta)...)
}

func resourceClusterInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).NeptuneConn(ctx)

	db, err := FindClusterInstanceByID(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Neptune Cluster Instance (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Neptune Cluster Instance (%s): %s", d.Id(), err)
	}

	clusterID := aws.StringValue(db.DBClusterIdentifier)
	d.Set("arn", db.DBInstanceArn)
	d.Set("auto_minor_version_upgrade", db.AutoMinorVersionUpgrade)
	d.Set("availability_zone", db.AvailabilityZone)
	d.Set("cluster_identifier", clusterID)
	d.Set("dbi_resource_id", db.DbiResourceId)
	d.Set("engine_version", db.EngineVersion)
	d.Set("engine", db.Engine)
	d.Set("identifier", db.DBInstanceIdentifier)
	d.Set("instance_class", db.DBInstanceClass)
	d.Set("kms_key_arn", db.KmsKeyId)
	if len(db.DBParameterGroups) > 0 {
		d.Set("neptune_parameter_group_name", db.DBParameterGroups[0].DBParameterGroupName)
	}
	if db.DBSubnetGroup != nil {
		d.Set("neptune_subnet_group_name", db.DBSubnetGroup.DBSubnetGroupName)
	}
	d.Set("preferred_backup_window", db.PreferredBackupWindow)
	d.Set("preferred_maintenance_window", db.PreferredMaintenanceWindow)
	d.Set("promotion_tier", db.PromotionTier)
	d.Set("publicly_accessible", db.PubliclyAccessible)
	d.Set("storage_encrypted", db.StorageEncrypted)

	if db.Endpoint != nil {
		address := aws.StringValue(db.Endpoint.Address)
		port := int(aws.Int64Value(db.Endpoint.Port))

		d.Set("address", address)
		d.Set("endpoint", fmt.Sprintf("%s:%d", address, port))
		d.Set("port", port)
	}

	m, err := FindClusterMemberByInstanceByTwoPartKey(ctx, conn, clusterID, d.Id())

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Neptune Cluster (%s) member (%s): %s", clusterID, d.Id(), err)
	}

	d.Set("writer", m.IsClusterWriter)

	return diags
}

func resourceClusterInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).NeptuneConn(ctx)

	if d.HasChangesExcept("tags", "tags_all") {
		input := &neptune.ModifyDBInstanceInput{
			ApplyImmediately:     aws.Bool(d.Get("apply_immediately").(bool)),
			DBInstanceIdentifier: aws.String(d.Id()),
		}

		if d.HasChange("auto_minor_version_upgrade") {
			input.AutoMinorVersionUpgrade = aws.Bool(d.Get("auto_minor_version_upgrade").(bool))
		}

		if d.HasChange("instance_class") {
			input.DBInstanceClass = aws.String(d.Get("instance_class").(string))
		}

		if d.HasChange("neptune_parameter_group_name") {
			input.DBParameterGroupName = aws.String(d.Get("neptune_parameter_group_name").(string))
		}

		if d.HasChange("preferred_backup_window") {
			input.PreferredBackupWindow = aws.String(d.Get("preferred_backup_window").(string))
		}

		if d.HasChange("preferred_maintenance_window") {
			input.PreferredMaintenanceWindow = aws.String(d.Get("preferred_maintenance_window").(string))
		}

		if d.HasChange("promotion_tier") {
			input.PromotionTier = aws.Int64(int64(d.Get("promotion_tier").(int)))
		}

		_, err := tfresource.RetryWhenAWSErrMessageContains(ctx, propagationTimeout, func() (interface{}, error) {
			return conn.ModifyDBInstanceWithContext(ctx, input)
		}, "InvalidParameterValue", "IAM role ARN value is invalid or does not include the required permissions")

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "modifying Neptune Cluster Instance (%s): %s", d.Id(), err)
		}

		if _, err := waitClusterInstanceAvailable(ctx, conn, d.Id(), d.Timeout(schema.TimeoutUpdate)); err != nil {
			return sdkdiag.AppendErrorf(diags, "waiting for Neptune Cluster Instance (%s) update: %s", d.Id(), err)
		}
	}

	return append(diags, resourceClusterInstanceRead(ctx, d, meta)...)
}

func resourceClusterInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).NeptuneConn(ctx)

	log.Printf("[DEBUG] Deleting Neptune Cluster Instance: %s", d.Id())
	_, err := conn.DeleteDBInstanceWithContext(ctx, &neptune.DeleteDBInstanceInput{
		DBInstanceIdentifier: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, neptune.ErrCodeDBInstanceNotFoundFault) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting Neptune Cluster Instance (%s): %s", d.Id(), err)
	}

	if _, err := waitClusterInstanceDeleted(ctx, conn, d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for Neptune Cluster Instance (%s) delete: %s", d.Id(), err)
	}

	return nil
}

func FindClusterInstanceByID(ctx context.Context, conn *neptune.Neptune, id string) (*neptune.DBInstance, error) {
	input := &neptune.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(id),
	}

	output, err := conn.DescribeDBInstancesWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, neptune.ErrCodeDBInstanceNotFoundFault) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || len(output.DBInstances) == 0 || output.DBInstances[0] == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	dbInstance := output.DBInstances[0]

	// Eventual consistency check.
	if aws.StringValue(dbInstance.DBInstanceIdentifier) != id {
		return nil, &retry.NotFoundError{
			LastRequest: input,
		}
	}

	return dbInstance, nil
}

func FindClusterMemberByInstanceByTwoPartKey(ctx context.Context, conn *neptune.Neptune, clusterID, instanceID string) (*neptune.DBClusterMember, error) {
	dbc, err := FindClusterByID(ctx, conn, clusterID)

	if err != nil {
		return nil, err
	}

	for _, v := range dbc.DBClusterMembers {
		if v == nil {
			continue
		}

		if aws.StringValue(v.DBInstanceIdentifier) == instanceID {
			return v, nil
		}
	}

	return nil, &retry.NotFoundError{}
}

func statusClusterInstance(ctx context.Context, conn *neptune.Neptune, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindClusterInstanceByID(ctx, conn, id)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.DBInstanceStatus), nil
	}
}

func waitClusterInstanceAvailable(ctx context.Context, conn *neptune.Neptune, id string, timeout time.Duration) (*neptune.DBInstance, error) { //nolint:unparam
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			"backing-up",
			"configuring-enhanced-monitoring",
			"configuring-iam-database-auth",
			"configuring-log-exports",
			"creating",
			"maintenance",
			"modifying",
			"rebooting",
			"renaming",
			"resetting-master-credentials",
			"starting",
			"storage-optimization",
			"upgrading",
		},
		Target:     []string{"available"},
		Refresh:    statusClusterInstance(ctx, conn, id),
		Timeout:    timeout,
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*neptune.DBInstance); ok {
		return output, err
	}

	return nil, err
}

func waitClusterInstanceDeleted(ctx context.Context, conn *neptune.Neptune, id string, timeout time.Duration) (*neptune.DBInstance, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			"modifying",
			"deleting",
		},
		Target:     []string{},
		Refresh:    statusClusterInstance(ctx, conn, id),
		Timeout:    timeout,
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*neptune.DBInstance); ok {
		return output, err
	}

	return nil, err
}
