package docdb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/docdb"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

// @SDKResource("aws_docdb_global_cluster")
func ResourceGlobalCluster() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceGlobalClusterCreate,
		ReadWithoutTimeout:   resourceGlobalClusterRead,
		UpdateWithoutTimeout: resourceGlobalClusterUpdate,
		DeleteWithoutTimeout: resourceGlobalClusterDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		//Timeouts will scale per number of resources in the cluster. Timeouts implemented on each resource action.
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(GlobalClusterCreateTimeout),
			Update: schema.DefaultTimeout(GlobalClusterUpdateTimeout),
			Delete: schema.DefaultTimeout(GlobalClusterDeleteTimeout),
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
				AtLeastOneOf:  []string{"engine", "source_db_cluster_identifier"},
				ConflictsWith: []string{"source_db_cluster_identifier"},
				ValidateFunc:  validEngine(),
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"global_cluster_identifier": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validGlobalCusterIdentifier,
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
				AtLeastOneOf:  []string{"engine", "source_db_cluster_identifier"},
				ConflictsWith: []string{"engine"},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
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

func resourceGlobalClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).DocDBConn(ctx)

	input := &docdb.CreateGlobalClusterInput{
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

	output, err := conn.CreateGlobalClusterWithContext(ctx, input)
	if err != nil {
		return diag.Errorf("creating DocumentDB Global Cluster: %s", err)
	}

	d.SetId(aws.StringValue(output.GlobalCluster.GlobalClusterIdentifier))

	if err := waitForGlobalClusterCreation(ctx, conn, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("waiting for DocumentDB Global Cluster (%s) availability: %s", d.Id(), err)
	}

	return resourceGlobalClusterRead(ctx, d, meta)
}

func resourceGlobalClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).DocDBConn(ctx)

	globalCluster, err := FindGlobalClusterById(ctx, conn, d.Id())

	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, docdb.ErrCodeGlobalClusterNotFoundFault) {
		log.Printf("[WARN] DocumentDB Global Cluster (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.Errorf("reading DocumentDB Global Cluster: %s", err)
	}

	if !d.IsNewResource() && globalCluster == nil {
		log.Printf("[WARN] DocumentDB Global Cluster (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if !d.IsNewResource() && (aws.StringValue(globalCluster.Status) == GlobalClusterStatusDeleting || aws.StringValue(globalCluster.Status) == GlobalClusterStatusDeleted) {
		log.Printf("[WARN] DocumentDB Global Cluster (%s) in deleted state (%s), removing from state", d.Id(), aws.StringValue(globalCluster.Status))
		d.SetId("")
		return nil
	}

	d.Set("arn", globalCluster.GlobalClusterArn)
	d.Set("database_name", globalCluster.DatabaseName)
	d.Set("deletion_protection", globalCluster.DeletionProtection)
	d.Set("engine", globalCluster.Engine)
	d.Set("engine_version", globalCluster.EngineVersion)
	d.Set("global_cluster_identifier", globalCluster.GlobalClusterIdentifier)

	if err := d.Set("global_cluster_members", flattenGlobalClusterMembers(globalCluster.GlobalClusterMembers)); err != nil {
		return diag.Errorf("setting global_cluster_members: %s", err)
	}

	d.Set("global_cluster_resource_id", globalCluster.GlobalClusterResourceId)
	d.Set("storage_encrypted", globalCluster.StorageEncrypted)

	return nil
}

func resourceGlobalClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).DocDBConn(ctx)

	input := &docdb.ModifyGlobalClusterInput{
		DeletionProtection:      aws.Bool(d.Get("deletion_protection").(bool)),
		GlobalClusterIdentifier: aws.String(d.Id()),
	}

	log.Printf("[DEBUG] Updating DocumentDB Global Cluster (%s): %s", d.Id(), input)

	if d.HasChange("engine_version") {
		if err := resourceGlobalClusterUpgradeEngineVersion(ctx, d, conn); err != nil {
			return diag.FromErr(err)
		}
	}

	_, err := conn.ModifyGlobalClusterWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, docdb.ErrCodeGlobalClusterNotFoundFault) {
		return nil
	}

	if err != nil {
		return diag.Errorf("updating DocumentDB Global Cluster: %s", err)
	}

	if err := waitForGlobalClusterUpdate(ctx, conn, d.Id(), d.Timeout(schema.TimeoutUpdate)); err != nil {
		return diag.Errorf("waiting for DocumentDB Global Cluster (%s) update: %s", d.Id(), err)
	}

	return resourceGlobalClusterRead(ctx, d, meta)
}

func resourceGlobalClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).DocDBConn(ctx)

	for _, globalClusterMemberRaw := range d.Get("global_cluster_members").(*schema.Set).List() {
		globalClusterMember, ok := globalClusterMemberRaw.(map[string]interface{})
		if !ok {
			continue
		}

		dbClusterArn, ok := globalClusterMember["db_cluster_arn"].(string)
		if !ok {
			continue
		}

		input := &docdb.RemoveFromGlobalClusterInput{
			DbClusterIdentifier:     aws.String(dbClusterArn),
			GlobalClusterIdentifier: aws.String(d.Id()),
		}

		_, err := conn.RemoveFromGlobalClusterWithContext(ctx, input)
		if tfawserr.ErrMessageContains(err, "InvalidParameterValue", "is not found in global cluster") {
			continue
		}
		if err != nil {
			return diag.Errorf("removing DocumentDB Cluster (%s) from Global Cluster (%s): %s", dbClusterArn, d.Id(), err)
		}

		if err := waitForGlobalClusterRemoval(ctx, conn, dbClusterArn, d.Timeout(schema.TimeoutDelete)); err != nil {
			return diag.Errorf("waiting for DocumentDB Cluster (%s) removal from DocumentDB Global Cluster (%s): %s", dbClusterArn, d.Id(), err)
		}
	}

	input := &docdb.DeleteGlobalClusterInput{
		GlobalClusterIdentifier: aws.String(d.Id()),
	}

	// Allow for eventual consistency
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		_, err := conn.DeleteGlobalClusterWithContext(ctx, input)

		if tfawserr.ErrMessageContains(err, docdb.ErrCodeInvalidGlobalClusterStateFault, "is not empty") {
			return retry.RetryableError(err)
		}

		if err != nil {
			return retry.NonRetryableError(err)
		}

		return nil
	})

	if tfresource.TimedOut(err) {
		_, err = conn.DeleteGlobalClusterWithContext(ctx, input)
	}

	if tfawserr.ErrCodeEquals(err, docdb.ErrCodeGlobalClusterNotFoundFault) {
		return nil
	}

	if err != nil {
		return diag.Errorf("deleting DocumentDB Global Cluster: %s", err)
	}

	if err := WaitForGlobalClusterDeletion(ctx, conn, d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("waiting for DocumentDB Global Cluster (%s) deletion: %s", d.Id(), err)
	}

	return nil
}

func flattenGlobalClusterMembers(apiObjects []*docdb.GlobalClusterMember) []interface{} {
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

// Updating major versions is not supported by documentDB
// To support minor version upgrades, we will upgrade all cluster members
func resourceGlobalClusterUpgradeEngineVersion(ctx context.Context, d *schema.ResourceData, conn *docdb.DocDB) error {
	log.Printf("[DEBUG] Upgrading DocumentDB Global Cluster (%s) engine version: %s", d.Id(), d.Get("engine_version"))
	err := resourceGlobalClusterUpgradeMinorEngineVersion(ctx, d.Get("global_cluster_members").(*schema.Set), d.Get("engine_version").(string), conn, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}
	globalCluster, err := FindGlobalClusterById(ctx, conn, d.Id())
	if err != nil {
		return err
	}
	for _, clusterMember := range globalCluster.GlobalClusterMembers {
		err := waitForClusterUpdate(ctx, conn, findGlobalClusterIDByARN(ctx, conn, aws.StringValue(clusterMember.DBClusterArn)), d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceGlobalClusterUpgradeMinorEngineVersion(ctx context.Context, clusterMembers *schema.Set, engineVersion string, conn *docdb.DocDB, timeout time.Duration) error {
	for _, clusterMemberRaw := range clusterMembers.List() {
		clusterMember := clusterMemberRaw.(map[string]interface{})
		if clusterMemberArn, ok := clusterMember["db_cluster_arn"]; ok && clusterMemberArn.(string) != "" {
			modInput := &docdb.ModifyDBClusterInput{
				ApplyImmediately:    aws.Bool(true),
				DBClusterIdentifier: aws.String(clusterMemberArn.(string)),
				EngineVersion:       aws.String(engineVersion),
			}
			err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
				_, err := conn.ModifyDBClusterWithContext(ctx, modInput)
				if err != nil {
					if tfawserr.ErrMessageContains(err, "InvalidParameterValue", "IAM role ARN value is invalid or does not include the required permissions") {
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			if tfresource.TimedOut(err) {
				_, err := conn.ModifyDBClusterWithContext(ctx, modInput)
				if err != nil {
					return err
				}
			}
			if err != nil {
				return fmt.Errorf("failed to update engine_version on global cluster member (%s): %w", clusterMemberArn, err)
			}
		}
	}
	return nil
}
