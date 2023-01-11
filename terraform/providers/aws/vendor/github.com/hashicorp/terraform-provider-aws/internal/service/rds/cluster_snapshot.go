package rds

import (
	"fmt"
	"log"
	"regexp"
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

const clusterSnapshotCreateTimeout = 2 * time.Minute

func ResourceClusterSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceClusterSnapshotCreate,
		Read:   resourceClusterSnapshotRead,
		Delete: resourceClusterSnapshotDelete,
		Update: resourcedbClusterSnapshotUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"db_cluster_snapshot_identifier": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 63),
					validation.StringMatch(regexp.MustCompile(`^[0-9a-z-]+$`), "must contain only lowercase alphanumeric characters and hyphens"),
					validation.StringMatch(regexp.MustCompile(`^[a-z]`), "must begin with a lowercase letter"),
					validation.StringDoesNotMatch(regexp.MustCompile(`--`), "cannot contain two consecutive hyphens"),
					validation.StringDoesNotMatch(regexp.MustCompile(`-$`), "cannot end with a hyphen"),
				),
			},
			"db_cluster_identifier": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"allocated_storage": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"availability_zones": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"db_cluster_snapshot_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_encrypted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"license_model": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"source_db_cluster_snapshot_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceClusterSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).RDSConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	params := &rds.CreateDBClusterSnapshotInput{
		DBClusterIdentifier:         aws.String(d.Get("db_cluster_identifier").(string)),
		DBClusterSnapshotIdentifier: aws.String(d.Get("db_cluster_snapshot_identifier").(string)),
		Tags:                        Tags(tags.IgnoreAWS()),
	}

	err := resource.Retry(clusterSnapshotCreateTimeout, func() *resource.RetryError {
		_, err := conn.CreateDBClusterSnapshot(params)
		if err != nil {
			if tfawserr.ErrCodeEquals(err, rds.ErrCodeInvalidDBClusterStateFault) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if tfresource.TimedOut(err) {
		_, err = conn.CreateDBClusterSnapshot(params)
	}
	if err != nil {
		return fmt.Errorf("error creating RDS DB Cluster Snapshot: %w", err)
	}
	d.SetId(d.Get("db_cluster_snapshot_identifier").(string))

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating"},
		Target:     []string{"available"},
		Refresh:    resourceClusterSnapshotStateRefreshFunc(d.Id(), conn),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 10 * time.Second,
		Delay:      5 * time.Second,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for RDS DB Cluster Snapshot %q to create: %s", d.Id(), err)
	}

	return resourceClusterSnapshotRead(d, meta)
}

func resourceClusterSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).RDSConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	params := &rds.DescribeDBClusterSnapshotsInput{
		DBClusterSnapshotIdentifier: aws.String(d.Id()),
	}
	resp, err := conn.DescribeDBClusterSnapshots(params)
	if err != nil {
		if tfawserr.ErrCodeEquals(err, rds.ErrCodeDBClusterSnapshotNotFoundFault) {
			log.Printf("[WARN] RDS DB Cluster Snapshot %q not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error reading RDS DB Cluster Snapshot %q: %s", d.Id(), err)
	}

	if resp == nil || len(resp.DBClusterSnapshots) == 0 || resp.DBClusterSnapshots[0] == nil || aws.StringValue(resp.DBClusterSnapshots[0].DBClusterSnapshotIdentifier) != d.Id() {
		log.Printf("[WARN] RDS DB Cluster Snapshot %q not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	snapshot := resp.DBClusterSnapshots[0]

	d.Set("allocated_storage", snapshot.AllocatedStorage)
	if err := d.Set("availability_zones", flex.FlattenStringList(snapshot.AvailabilityZones)); err != nil {
		return fmt.Errorf("error setting availability_zones: %s", err)
	}
	d.Set("db_cluster_identifier", snapshot.DBClusterIdentifier)
	d.Set("db_cluster_snapshot_arn", snapshot.DBClusterSnapshotArn)
	d.Set("db_cluster_snapshot_identifier", snapshot.DBClusterSnapshotIdentifier)
	d.Set("engine_version", snapshot.EngineVersion)
	d.Set("engine", snapshot.Engine)
	d.Set("kms_key_id", snapshot.KmsKeyId)
	d.Set("license_model", snapshot.LicenseModel)
	d.Set("port", snapshot.Port)
	d.Set("snapshot_type", snapshot.SnapshotType)
	d.Set("source_db_cluster_snapshot_arn", snapshot.SourceDBClusterSnapshotArn)
	d.Set("status", snapshot.Status)
	d.Set("storage_encrypted", snapshot.StorageEncrypted)
	d.Set("vpc_id", snapshot.VpcId)

	tags, err := ListTags(conn, d.Get("db_cluster_snapshot_arn").(string))

	if err != nil {
		return fmt.Errorf("error listing tags for RDS DB Cluster Snapshot (%s): %s", d.Get("db_cluster_snapshot_arn").(string), err)
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

func resourcedbClusterSnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).RDSConn

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := UpdateTags(conn, d.Get("db_cluster_snapshot_arn").(string), o, n); err != nil {
			return fmt.Errorf("error updating RDS DB Cluster Snapshot (%s) tags: %s", d.Get("db_cluster_snapshot_arn").(string), err)
		}
	}

	return nil
}

func resourceClusterSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).RDSConn

	params := &rds.DeleteDBClusterSnapshotInput{
		DBClusterSnapshotIdentifier: aws.String(d.Id()),
	}
	_, err := conn.DeleteDBClusterSnapshot(params)
	if err != nil {
		if tfawserr.ErrCodeEquals(err, rds.ErrCodeDBClusterSnapshotNotFoundFault) {
			return nil
		}
		return fmt.Errorf("error deleting RDS DB Cluster Snapshot %q: %s", d.Id(), err)
	}

	return nil
}

func resourceClusterSnapshotStateRefreshFunc(dbClusterSnapshotIdentifier string, conn *rds.RDS) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		opts := &rds.DescribeDBClusterSnapshotsInput{
			DBClusterSnapshotIdentifier: aws.String(dbClusterSnapshotIdentifier),
		}

		log.Printf("[DEBUG] DB Cluster Snapshot describe configuration: %#v", opts)

		resp, err := conn.DescribeDBClusterSnapshots(opts)
		if err != nil {
			if tfawserr.ErrCodeEquals(err, rds.ErrCodeDBClusterSnapshotNotFoundFault) {
				return nil, "", nil
			}
			return nil, "", fmt.Errorf("Error retrieving DB Cluster Snapshots: %s", err)
		}

		if resp == nil || len(resp.DBClusterSnapshots) == 0 || resp.DBClusterSnapshots[0] == nil {
			return nil, "", fmt.Errorf("No snapshots returned for %s", dbClusterSnapshotIdentifier)
		}

		snapshot := resp.DBClusterSnapshots[0]

		return resp, aws.StringValue(snapshot.Status), nil
	}
}
