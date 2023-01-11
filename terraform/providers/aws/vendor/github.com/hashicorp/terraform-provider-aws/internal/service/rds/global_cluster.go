package rds

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

const (
	GlobalClusterRemovalTimeout = 30 * time.Minute
	globalClusterCreateTimeout  = 30 * time.Minute
	globalClusterUpdateTimeout  = 90 * time.Minute
)

func ResourceGlobalCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceGlobalClusterCreate,
		Read:   resourceGlobalClusterRead,
		Update: resourceGlobalClusterUpdate,
		Delete: resourceGlobalClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(globalClusterCreateTimeout),
			Update: schema.DefaultTimeout(globalClusterUpdateTimeout),
			Delete: schema.DefaultTimeout(GlobalClusterRemovalTimeout),
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"database_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"engine": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"source_db_cluster_identifier"},
				ValidateFunc: validation.StringInSlice([]string{
					"aurora",
					"aurora-mysql",
					"aurora-postgresql",
				}, false),
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
			"force_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"global_cluster_identifier": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"global_cluster_members": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_cluster_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_writer": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"global_cluster_resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_db_cluster_identifier": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"engine"},
				RequiredWith:  []string{"force_destroy"},
			},
			"storage_encrypted": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceGlobalClusterCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).RDSConn

	input := &rds.CreateGlobalClusterInput{
		GlobalClusterIdentifier: aws.String(d.Get("global_cluster_identifier").(string)),
	}

	if v, ok := d.GetOk("database_name"); ok {
		input.DatabaseName = aws.String(v.(string))
	}

	if v, ok := d.GetOk("deletion_protection"); ok {
		input.DeletionProtection = aws.Bool(v.(bool))
	}

	if v, ok := d.GetOk("engine"); ok {
		input.Engine = aws.String(v.(string))
	}

	if v, ok := d.GetOk("engine_version"); ok {
		input.EngineVersion = aws.String(v.(string))
	}

	if v, ok := d.GetOk("source_db_cluster_identifier"); ok {
		input.SourceDBClusterIdentifier = aws.String(v.(string))
	}

	if v, ok := d.GetOk("storage_encrypted"); ok {
		input.StorageEncrypted = aws.Bool(v.(bool))
	}

	// Prevent the following error and keep the previous default,
	// since we cannot have Engine default after adding SourceDBClusterIdentifier:
	// InvalidParameterValue: When creating standalone global cluster, value for engineName should be specified
	if input.Engine == nil && input.SourceDBClusterIdentifier == nil {
		input.Engine = aws.String("aurora")
	}

	output, err := conn.CreateGlobalCluster(input)
	if err != nil {
		return fmt.Errorf("error creating RDS Global Cluster: %s", err)
	}

	d.SetId(aws.StringValue(output.GlobalCluster.GlobalClusterIdentifier))

	if err := waitForGlobalClusterCreation(conn, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
		return fmt.Errorf("error waiting for RDS Global Cluster (%s) availability: %s", d.Id(), err)
	}

	return resourceGlobalClusterRead(d, meta)
}

func resourceGlobalClusterRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).RDSConn

	globalCluster, err := DescribeGlobalCluster(conn, d.Id())

	if tfawserr.ErrCodeEquals(err, rds.ErrCodeGlobalClusterNotFoundFault) {
		log.Printf("[WARN] RDS Global Cluster (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading RDS Global Cluster: %s", err)
	}

	if globalCluster == nil {
		log.Printf("[WARN] RDS Global Cluster (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if aws.StringValue(globalCluster.Status) == "deleting" || aws.StringValue(globalCluster.Status) == "deleted" {
		log.Printf("[WARN] RDS Global Cluster (%s) in deleted state (%s), removing from state", d.Id(), aws.StringValue(globalCluster.Status))
		d.SetId("")
		return nil
	}

	d.Set("arn", globalCluster.GlobalClusterArn)
	d.Set("database_name", globalCluster.DatabaseName)
	d.Set("deletion_protection", globalCluster.DeletionProtection)
	d.Set("engine", globalCluster.Engine)
	d.Set("global_cluster_identifier", globalCluster.GlobalClusterIdentifier)
	if err := d.Set("global_cluster_members", flattenGlobalClusterMembers(globalCluster.GlobalClusterMembers)); err != nil {
		return fmt.Errorf("error setting global_cluster_members: %w", err)
	}
	d.Set("global_cluster_resource_id", globalCluster.GlobalClusterResourceId)
	d.Set("storage_encrypted", globalCluster.StorageEncrypted)

	oldEngineVersion := d.Get("engine_version").(string)
	newEngineVersion := aws.StringValue(globalCluster.EngineVersion)

	// For example a configured engine_version of "5.6.10a" and a returned engine_version of "5.6.global_10a".
	if oldParts, newParts := strings.Split(oldEngineVersion, "."), strings.Split(newEngineVersion, "."); len(oldParts) == 3 &&
		len(oldParts) == len(newParts) &&
		oldParts[0] == newParts[0] &&
		oldParts[1] == newParts[1] &&
		strings.HasSuffix(newParts[2], oldParts[2]) {
		d.Set("engine_version", oldEngineVersion)
		d.Set("engine_version_actual", newEngineVersion)
	} else {
		d.Set("engine_version", newEngineVersion)
		d.Set("engine_version_actual", newEngineVersion)
	}

	return nil
}

func resourceGlobalClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).RDSConn

	input := &rds.ModifyGlobalClusterInput{
		DeletionProtection:      aws.Bool(d.Get("deletion_protection").(bool)),
		GlobalClusterIdentifier: aws.String(d.Id()),
	}

	if d.HasChange("engine_version") {
		if err := globalClusterUpgradeEngineVersion(d, meta, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return err
		}
	}

	log.Printf("[DEBUG] Updating RDS Global Cluster (%s): %s", d.Id(), input)
	_, err := conn.ModifyGlobalCluster(input)

	if tfawserr.ErrCodeEquals(err, rds.ErrCodeGlobalClusterNotFoundFault) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting RDS Global Cluster: %s", err)
	}

	if err := waitForGlobalClusterUpdate(conn, d.Id(), d.Timeout(schema.TimeoutUpdate)); err != nil {
		return fmt.Errorf("error waiting for RDS Global Cluster (%s) update: %s", d.Id(), err)
	}

	return resourceGlobalClusterRead(d, meta)
}

func resourceGlobalClusterDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).RDSConn

	if d.Get("force_destroy").(bool) {
		for _, globalClusterMemberRaw := range d.Get("global_cluster_members").(*schema.Set).List() {
			globalClusterMember, ok := globalClusterMemberRaw.(map[string]interface{})

			if !ok {
				continue
			}

			dbClusterArn, ok := globalClusterMember["db_cluster_arn"].(string)

			if !ok {
				continue
			}

			input := &rds.RemoveFromGlobalClusterInput{
				DbClusterIdentifier:     aws.String(dbClusterArn),
				GlobalClusterIdentifier: aws.String(d.Id()),
			}

			_, err := conn.RemoveFromGlobalCluster(input)

			if tfawserr.ErrMessageContains(err, "InvalidParameterValue", "is not found in global cluster") {
				continue
			}

			if err != nil {
				return fmt.Errorf("error removing RDS DB Cluster (%s) from Global Cluster (%s): %w", dbClusterArn, d.Id(), err)
			}

			if err := waitForGlobalClusterRemoval(conn, dbClusterArn, d.Timeout(schema.TimeoutDelete)); err != nil {
				return fmt.Errorf("error waiting for RDS DB Cluster (%s) removal from RDS Global Cluster (%s): %w", dbClusterArn, d.Id(), err)
			}
		}
	}

	input := &rds.DeleteGlobalClusterInput{
		GlobalClusterIdentifier: aws.String(d.Id()),
	}

	log.Printf("[DEBUG] Deleting RDS Global Cluster (%s): %s", d.Id(), input)

	// Allow for eventual consistency
	// InvalidGlobalClusterStateFault: Global Cluster arn:aws:rds::123456789012:global-cluster:tf-acc-test-5618525093076697001-0 is not empty
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := conn.DeleteGlobalCluster(input)

		if tfawserr.ErrMessageContains(err, rds.ErrCodeInvalidGlobalClusterStateFault, "is not empty") {
			return resource.RetryableError(err)
		}

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if tfresource.TimedOut(err) {
		_, err = conn.DeleteGlobalCluster(input)
	}

	if tfawserr.ErrCodeEquals(err, rds.ErrCodeGlobalClusterNotFoundFault) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting RDS Global Cluster: %s", err)
	}

	if err := WaitForGlobalClusterDeletion(conn, d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return fmt.Errorf("error waiting for RDS Global Cluster (%s) deletion: %s", d.Id(), err)
	}

	return nil
}

func flattenGlobalClusterMembers(apiObjects []*rds.GlobalClusterMember) []interface{} {
	if len(apiObjects) == 0 {
		return nil
	}

	var tfList []interface{}

	for _, apiObject := range apiObjects {
		tfMap := map[string]interface{}{
			"db_cluster_arn": aws.StringValue(apiObject.DBClusterArn),
			"is_writer":      aws.BoolValue(apiObject.IsWriter),
		}

		tfList = append(tfList, tfMap)
	}

	return tfList
}

func DescribeGlobalCluster(conn *rds.RDS, globalClusterID string) (*rds.GlobalCluster, error) {
	var globalCluster *rds.GlobalCluster

	input := &rds.DescribeGlobalClustersInput{
		GlobalClusterIdentifier: aws.String(globalClusterID),
	}

	log.Printf("[DEBUG] Reading RDS Global Cluster (%s): %s", globalClusterID, input)
	err := conn.DescribeGlobalClustersPages(input, func(page *rds.DescribeGlobalClustersOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, gc := range page.GlobalClusters {
			if gc == nil {
				continue
			}

			if aws.StringValue(gc.GlobalClusterIdentifier) == globalClusterID {
				globalCluster = gc
				return false
			}
		}

		return !lastPage
	})

	return globalCluster, err
}

func DescribeGlobalClusterFromClusterARN(conn *rds.RDS, dbClusterARN string) (*rds.GlobalCluster, error) {
	var globalCluster *rds.GlobalCluster

	input := &rds.DescribeGlobalClustersInput{
		Filters: []*rds.Filter{
			{
				Name:   aws.String("db-cluster-id"),
				Values: []*string{aws.String(dbClusterARN)},
			},
		},
	}

	log.Printf("[DEBUG] Reading RDS Global Clusters: %s", input)
	err := conn.DescribeGlobalClustersPages(input, func(page *rds.DescribeGlobalClustersOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, gc := range page.GlobalClusters {
			if gc == nil {
				continue
			}

			for _, globalClusterMember := range gc.GlobalClusterMembers {
				if aws.StringValue(globalClusterMember.DBClusterArn) == dbClusterARN {
					globalCluster = gc
					return false
				}
			}
		}

		return !lastPage
	})

	return globalCluster, err
}

func globalClusterRefreshFunc(conn *rds.RDS, globalClusterID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		globalCluster, err := DescribeGlobalCluster(conn, globalClusterID)

		if tfawserr.ErrCodeEquals(err, rds.ErrCodeGlobalClusterNotFoundFault) {
			return nil, "deleted", nil
		}

		if err != nil {
			return nil, "", fmt.Errorf("error reading RDS Global Cluster (%s): %s", globalClusterID, err)
		}

		if globalCluster == nil {
			return nil, "deleted", nil
		}

		return globalCluster, aws.StringValue(globalCluster.Status), nil
	}
}

func waitForGlobalClusterCreation(conn *rds.RDS, globalClusterID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"creating"},
		Target:  []string{"available"},
		Refresh: globalClusterRefreshFunc(conn, globalClusterID),
		Timeout: timeout,
	}

	log.Printf("[DEBUG] Waiting for RDS Global Cluster (%s) availability", globalClusterID)
	_, err := stateConf.WaitForState()

	return err
}

func waitForGlobalClusterUpdate(conn *rds.RDS, globalClusterID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"modifying", "upgrading"},
		Target:  []string{"available"},
		Refresh: globalClusterRefreshFunc(conn, globalClusterID),
		Timeout: timeout,
		Delay:   30 * time.Second,
	}

	log.Printf("[DEBUG] Waiting for RDS Global Cluster (%s) availability", globalClusterID)
	_, err := stateConf.WaitForState()

	return err
}

func WaitForGlobalClusterDeletion(conn *rds.RDS, globalClusterID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			"available",
			"deleting",
		},
		Target:         []string{"deleted"},
		Refresh:        globalClusterRefreshFunc(conn, globalClusterID),
		Timeout:        timeout,
		NotFoundChecks: 1,
	}

	log.Printf("[DEBUG] Waiting for RDS Global Cluster (%s) deletion", globalClusterID)
	_, err := stateConf.WaitForState()

	if tfresource.NotFound(err) {
		return nil
	}

	return err
}

func waitForGlobalClusterRemoval(conn *rds.RDS, dbClusterIdentifier string, timeout time.Duration) error {
	var globalCluster *rds.GlobalCluster
	stillExistsErr := fmt.Errorf("RDS DB Cluster still exists in RDS Global Cluster")

	err := resource.Retry(timeout, func() *resource.RetryError {
		var err error

		globalCluster, err = DescribeGlobalClusterFromClusterARN(conn, dbClusterIdentifier)

		if err != nil {
			return resource.NonRetryableError(err)
		}

		if globalCluster != nil {
			return resource.RetryableError(stillExistsErr)
		}

		return nil
	})

	if tfresource.TimedOut(err) {
		_, err = DescribeGlobalClusterFromClusterARN(conn, dbClusterIdentifier)
	}

	if err != nil {
		return err
	}

	if globalCluster != nil {
		return stillExistsErr
	}

	return nil
}

func globalClusterUpgradeMajorEngineVersion(meta interface{}, clusterID string, engineVersion string, timeout time.Duration) error {
	conn := meta.(*conns.AWSClient).RDSConn

	input := &rds.ModifyGlobalClusterInput{
		GlobalClusterIdentifier: aws.String(clusterID),
	}

	input.AllowMajorVersionUpgrade = aws.Bool(true)
	input.EngineVersion = aws.String(engineVersion)

	err := resource.Retry(timeout, func() *resource.RetryError {
		_, err := conn.ModifyGlobalCluster(input)

		if err != nil {
			if tfawserr.ErrCodeEquals(err, rds.ErrCodeGlobalClusterNotFoundFault) {
				return resource.NonRetryableError(err)
			}

			if tfawserr.ErrMessageContains(err, "InvalidParameterValue", "only supports Major Version Upgrades") {
				return resource.NonRetryableError(err)
			}

			return resource.RetryableError(err)
		}

		return nil
	})

	if tfresource.TimedOut(err) {
		_, err = conn.ModifyGlobalCluster(input)
	}

	if err != nil {
		return fmt.Errorf("while upgrading major version of RDS Global Cluster (%s): %w", clusterID, err)
	}

	globalCluster, err := DescribeGlobalCluster(conn, clusterID)

	if err != nil {
		return fmt.Errorf("while upgrading major version of RDS Global Cluster (%s): %w", clusterID, err)
	}

	for _, clusterMember := range globalCluster.GlobalClusterMembers {
		arnID := aws.StringValue(clusterMember.DBClusterArn)

		if arnID == "" {
			continue
		}

		dbi, clusterRegion, err := ClusterIDRegionFromARN(arnID)

		if err != nil {
			return fmt.Errorf("while upgrading RDS Global Cluster Cluster minor engine version: %w", err)
		}

		if dbi == "" {
			continue
		}

		useConn := conn // clusters may not all be in the same region

		if clusterRegion != meta.(*conns.AWSClient).Region {
			useConn = rds.New(meta.(*conns.AWSClient).Session, aws.NewConfig().WithRegion(clusterRegion))
		}

		if err := WaitForClusterUpdate(useConn, dbi, timeout); err != nil {
			return fmt.Errorf("failed to update engine_version, waiting for RDS Global Cluster (%s) to update: %s", dbi, err)
		}
	}

	return err
}

func ClusterIDRegionFromARN(arnID string) (string, string, error) {
	parsedARN, err := arn.Parse(arnID)

	if err != nil {
		return "", "", fmt.Errorf("could not parse ARN (%s): %w", arnID, err)
	}

	dbi := ""

	if parsedARN.Resource != "" {
		parts := strings.Split(parsedARN.Resource, ":")

		if len(parts) < 2 {
			return "", "", fmt.Errorf("could not get DB Cluster ID from parsing ARN (%s): %w", arnID, err)
		}

		if parsedARN.Service != endpoints.RdsServiceID || parts[0] != "cluster" {
			return "", "", fmt.Errorf("wrong ARN (%s) for a DB Cluster", arnID)
		}

		dbi = parts[1]
	}

	return dbi, parsedARN.Region, nil
}

func globalClusterUpgradeMinorEngineVersion(meta interface{}, clusterMembers *schema.Set, clusterID, engineVersion string, timeout time.Duration) error {
	conn := meta.(*conns.AWSClient).RDSConn

	log.Printf("[INFO] Performing RDS Global Cluster (%s) minor version (%s) upgrade", clusterID, engineVersion)

	for _, clusterMemberRaw := range clusterMembers.List() {
		clusterMember := clusterMemberRaw.(map[string]interface{})

		// DBClusterIdentifier supposedly can be either ARN or ID, and both used to work,
		// but as of now, only ID works
		if clusterMemberArn, ok := clusterMember["db_cluster_arn"]; !ok || clusterMemberArn.(string) == "" {
			continue
		}

		arnID := clusterMember["db_cluster_arn"].(string)

		dbi, clusterRegion, err := ClusterIDRegionFromARN(arnID)

		if err != nil {
			return fmt.Errorf("while upgrading RDS Global Cluster Cluster minor engine version: %w", err)
		}

		if dbi == "" {
			continue
		}

		useConn := conn

		if clusterRegion != meta.(*conns.AWSClient).Region {
			useConn = rds.New(meta.(*conns.AWSClient).Session, aws.NewConfig().WithRegion(clusterRegion))
		}

		modInput := &rds.ModifyDBClusterInput{
			ApplyImmediately:    aws.Bool(true),
			DBClusterIdentifier: aws.String(dbi),
			EngineVersion:       aws.String(engineVersion),
		}

		log.Printf("[INFO] Performing RDS Global Cluster (%s) Cluster (%s) minor version (%s) upgrade", clusterID, dbi, engineVersion)

		err = resource.Retry(timeout, func() *resource.RetryError {
			_, err := useConn.ModifyDBCluster(modInput)

			if err != nil {
				if tfawserr.ErrMessageContains(err, "InvalidParameterValue", "IAM role ARN value is invalid or does not include the required permissions") {
					return resource.RetryableError(err)
				}

				if tfawserr.ErrMessageContains(err, rds.ErrCodeInvalidDBClusterStateFault, "Cannot modify engine version without a primary instance in DB cluster") {
					return resource.NonRetryableError(err)
				}

				if tfawserr.ErrCodeEquals(err, rds.ErrCodeInvalidDBClusterStateFault) {
					return resource.RetryableError(err)
				}

				return resource.NonRetryableError(err)
			}
			return nil
		})

		if tfresource.TimedOut(err) {
			_, err := useConn.ModifyDBCluster(modInput)

			if err != nil {
				return err
			}
		}

		if err != nil {
			return fmt.Errorf("failed to update engine_version on RDS Global Cluster Cluster (%s): %s", dbi, err)
		}

		log.Printf("[INFO] Waiting for RDS Global Cluster (%s) Cluster (%s) minor version (%s) upgrade", clusterID, dbi, engineVersion)
		if err := WaitForClusterUpdate(useConn, dbi, timeout); err != nil {
			return fmt.Errorf("failed to update engine_version, waiting for RDS Global Cluster Cluster (%s) to update: %s", dbi, err)
		}
	}

	globalCluster, err := DescribeGlobalCluster(conn, clusterID)

	if tfawserr.ErrCodeEquals(err, rds.ErrCodeGlobalClusterNotFoundFault) {
		return fmt.Errorf("after upgrading engine_version, could not find RDS Global Cluster (%s): %s", clusterID, err)
	}

	if err != nil {
		return fmt.Errorf("after minor engine_version upgrade to RDS Global Cluster (%s): %s", clusterID, err)
	}

	if globalCluster == nil {
		return fmt.Errorf("after minor engine_version upgrade to RDS Global Cluster (%s): empty response", clusterID)
	}

	if aws.StringValue(globalCluster.EngineVersion) != engineVersion {
		log.Printf("[DEBUG] RDS Global Cluster (%s) upgrade did not take effect, trying again", clusterID)
		return globalClusterUpgradeMinorEngineVersion(meta, clusterMembers, clusterID, engineVersion, timeout)
	}

	return nil
}

func globalClusterUpgradeEngineVersion(d *schema.ResourceData, meta interface{}, timeout time.Duration) error {
	log.Printf("[DEBUG] Upgrading RDS Global Cluster (%s) engine version: %s", d.Id(), d.Get("engine_version"))

	err := globalClusterUpgradeMajorEngineVersion(meta, d.Id(), d.Get("engine_version").(string), timeout)

	if tfawserr.ErrMessageContains(err, "InvalidParameterValue", "only supports Major Version Upgrades") {
		err = globalClusterUpgradeMinorEngineVersion(meta, d.Get("global_cluster_members").(*schema.Set), d.Id(), d.Get("engine_version").(string), timeout)

		if err != nil {
			return fmt.Errorf("while upgrading minor version of RDS Global Cluster (%s): %w", d.Id(), err)
		}

		return nil
	}

	if err != nil {
		return fmt.Errorf("while upgrading major version of RDS Global Cluster (%s): %w", d.Id(), err)
	}

	return nil
}

var resourceClusterUpdatePendingStates = []string{
	"backing-up",
	"configuring-iam-database-auth",
	"modifying",
	"renaming",
	"resetting-master-credentials",
	"upgrading",
}

func WaitForClusterUpdate(conn *rds.RDS, id string, timeout time.Duration) error {
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

func resourceClusterStateRefreshFunc(conn *rds.RDS, dbClusterIdentifier string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := conn.DescribeDBClusters(&rds.DescribeDBClustersInput{
			DBClusterIdentifier: aws.String(dbClusterIdentifier),
		})

		if tfawserr.ErrCodeEquals(err, rds.ErrCodeDBClusterNotFoundFault) {
			return 42, "destroyed", nil
		}

		if err != nil {
			return nil, "", err
		}

		var dbc *rds.DBCluster

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
