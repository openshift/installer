package elasticache

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/experimental/nullable"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

const (
	defaultRedisPort     = "6379"
	defaultMemcachedPort = "11211"
)

const (
	cacheClusterCreatedTimeout = 40 * time.Minute
)

func ResourceCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceClusterCreate,
		Read:   resourceClusterRead,
		Update: resourceClusterUpdate,
		Delete: resourceClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
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
				Type:         nullable.TypeNullableBool,
				Optional:     true,
				Default:      "true",
				ValidateFunc: nullable.ValidateTypeStringNullableBool,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"az_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					elasticache.AZModeCrossAz,
					elasticache.AZModeSingleAz,
				}, false),
			},
			"cache_nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"cluster_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				StateFunc: func(val interface{}) string {
					// ElastiCache normalizes cluster ids to lowercase,
					// so we have to do this too or else we can end up
					// with non-converging diffs.
					return strings.ToLower(val.(string))
				},
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 50),
					validation.StringMatch(regexp.MustCompile(`^[0-9a-z-]+$`), "must contain only lowercase alphanumeric characters and hyphens"),
					validation.StringMatch(regexp.MustCompile(`^[a-z]`), "must begin with a lowercase letter"),
					validation.StringDoesNotMatch(regexp.MustCompile(`--`), "cannot contain two consecutive hyphens"),
					validation.StringDoesNotMatch(regexp.MustCompile(`-$`), "cannot end with a hyphen"),
				),
			},
			"configuration_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"engine", "replication_group_id"},
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(engine_Values(), false),
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
			"log_delivery_configuration": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 2,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(elasticache.DestinationType_Values(), false),
						},
						"destination": {
							Type:     schema.TypeString,
							Required: true,
						},
						"log_format": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(elasticache.LogFormat_Values(), false),
						},
						"log_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(elasticache.LogType_Values(), false),
						},
					},
				},
			},
			"maintenance_window": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				StateFunc: func(val interface{}) string {
					// ElastiCache always changes the maintenance
					// to lowercase
					return strings.ToLower(val.(string))
				},
				ValidateFunc: verify.ValidOnceAWeekWindowFormat,
			},
			"node_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"notification_topic_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"num_cache_nodes": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"parameter_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// Suppress default memcached/redis ports when not defined
					if !d.IsNewResource() && new == "0" && (old == defaultRedisPort || old == defaultMemcachedPort) {
						return true
					}
					return false
				},
			},
			"preferred_availability_zones": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"replication_group_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"replication_group_id", "engine"},
				ForceNew:     true,
				ValidateFunc: validateReplicationGroupID,
				ConflictsWith: []string{
					"az_mode",
					"engine_version",
					"maintenance_window",
					"node_type",
					"notification_topic_arn",
					"num_cache_nodes",
					"parameter_group_name",
					"port",
					"security_group_ids",
					"security_group_names",
					"snapshot_arns",
					"snapshot_name",
					"snapshot_retention_limit",
					"snapshot_window",
					"subnet_group_name",
				},
			},
			"security_group_names": {
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				ForceNew:   true,
				Elem:       &schema.Schema{Type: schema.TypeString},
				Set:        schema.HashString,
				Deprecated: `With the retirement of EC2-Classic the security_group_names attribute has been deprecated and will be removed in a future version.`,
			},
			"security_group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"snapshot_arns": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.All(
						verify.ValidARN,
						validation.StringDoesNotContainAny(","),
					),
				},
			},
			"snapshot_retention_limit": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtMost(35),
			},
			"snapshot_window": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: verify.ValidOnceADayWindowFormat,
			},
			"snapshot_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"subnet_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"final_snapshot_identifier": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
		},

		CustomizeDiff: customdiff.Sequence(
			CustomizeDiffValidateClusterAZMode,
			CustomizeDiffValidateClusterEngineVersion,
			customizeDiffEngineVersionForceNewOnDowngrade,
			CustomizeDiffValidateClusterNumCacheNodes,
			CustomizeDiffClusterMemcachedNodeType,
			CustomizeDiffValidateClusterMemcachedSnapshotIdentifier,
			verify.SetTagsDiff,
		),
	}
}

func resourceClusterCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ElastiCacheConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	if v, ok := d.GetOk("security_group_names"); ok && v.(*schema.Set).Len() > 0 {
		return errors.New(`with the retirement of EC2-Classic no new ElastiCache Clusters can be created referencing ElastiCache Security Groups`)
	}

	req := &elasticache.CreateCacheClusterInput{}

	if v, ok := d.GetOk("replication_group_id"); ok {
		req.ReplicationGroupId = aws.String(v.(string))
	} else {
		req.SecurityGroupIds = flex.ExpandStringSet(d.Get("security_group_ids").(*schema.Set))

		if len(tags) > 0 {
			req.Tags = Tags(tags.IgnoreAWS())
		}
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		req.CacheClusterId = aws.String(v.(string))
	}

	if v, ok := d.GetOk("node_type"); ok {
		req.CacheNodeType = aws.String(v.(string))
	}

	if v, ok := d.GetOk("num_cache_nodes"); ok {
		req.NumCacheNodes = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("engine"); ok {
		req.Engine = aws.String(v.(string))
	}

	version := d.Get("engine_version").(string)
	if version != "" {
		req.EngineVersion = aws.String(version)
	}

	if v, ok := d.GetOk("auto_minor_version_upgrade"); ok {
		if v, null, _ := nullable.Bool(v.(string)).Value(); !null {
			req.AutoMinorVersionUpgrade = aws.Bool(v)
		}
	}

	if v, ok := d.GetOk("port"); ok {
		req.Port = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("subnet_group_name"); ok {
		req.CacheSubnetGroupName = aws.String(v.(string))
	}

	// parameter groups are optional and can be defaulted by AWS
	if v, ok := d.GetOk("parameter_group_name"); ok {
		req.CacheParameterGroupName = aws.String(v.(string))
	}

	if v, ok := d.GetOk("snapshot_retention_limit"); ok {
		req.SnapshotRetentionLimit = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("snapshot_window"); ok {
		req.SnapshotWindow = aws.String(v.(string))
	}

	if v, ok := d.GetOk("log_delivery_configuration"); ok {
		req.LogDeliveryConfigurations = []*elasticache.LogDeliveryConfigurationRequest{}
		v := v.(*schema.Set).List()
		for _, v := range v {
			logDeliveryConfigurationRequest := expandLogDeliveryConfigurations(v.(map[string]interface{}))
			req.LogDeliveryConfigurations = append(req.LogDeliveryConfigurations, &logDeliveryConfigurationRequest)
		}
	}

	if v, ok := d.GetOk("maintenance_window"); ok {
		req.PreferredMaintenanceWindow = aws.String(v.(string))
	}

	if v, ok := d.GetOk("notification_topic_arn"); ok {
		req.NotificationTopicArn = aws.String(v.(string))
	}

	snaps := d.Get("snapshot_arns").([]interface{})
	if len(snaps) > 0 {
		req.SnapshotArns = flex.ExpandStringList(snaps)
		log.Printf("[DEBUG] Restoring Redis cluster from S3 snapshot: %#v", snaps)
	}

	if v, ok := d.GetOk("snapshot_name"); ok {
		req.SnapshotName = aws.String(v.(string))
	}

	if v, ok := d.GetOk("az_mode"); ok {
		req.AZMode = aws.String(v.(string))
	}

	if v, ok := d.GetOk("availability_zone"); ok {
		req.PreferredAvailabilityZone = aws.String(v.(string))
	}

	if v, ok := d.GetOk("preferred_availability_zones"); ok && len(v.([]interface{})) > 0 {
		req.PreferredAvailabilityZones = flex.ExpandStringList(v.([]interface{}))
	}

	id, arn, err := createCacheCluster(conn, req)
	if err != nil {
		return fmt.Errorf("error creating ElastiCache Cache Cluster: %w", err)
	}

	d.SetId(id)

	_, err = waitCacheClusterAvailable(conn, d.Id(), cacheClusterCreatedTimeout)
	if err != nil {
		return fmt.Errorf("error waiting for ElastiCache Cache Cluster (%s) to be created: %w", d.Id(), err)
	}

	// Only post-create tagging supported in some partitions
	if req.Tags == nil && len(tags) > 0 {
		err := UpdateTags(conn, arn, nil, tags)

		if err != nil {
			if v, ok := d.GetOk("tags"); (ok && len(v.(map[string]interface{})) > 0) || !verify.ErrorISOUnsupported(conn.PartitionID, err) {
				// explicitly setting tags or not an iso-unsupported error
				return fmt.Errorf("failed adding tags after create for ElastiCache Cache Cluster (%s): %w", d.Id(), err)
			}

			log.Printf("[WARN] failed adding tags after create for ElastiCache Cache Cluster (%s): %s", d.Id(), err)
		}
	}

	return resourceClusterRead(d, meta)
}

func resourceClusterRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ElastiCacheConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	c, err := FindCacheClusterWithNodeInfoByID(conn, d.Id())
	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] ElastiCache Cache Cluster (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("error reading ElastiCache Cache Cluster (%s): %w", d.Id(), err)
	}

	d.Set("cluster_id", c.CacheClusterId)

	if err := setFromCacheCluster(d, c); err != nil {
		return err
	}

	d.Set("log_delivery_configuration", flattenLogDeliveryConfigurations(c.LogDeliveryConfigurations))
	d.Set("snapshot_window", c.SnapshotWindow)
	d.Set("snapshot_retention_limit", c.SnapshotRetentionLimit)

	d.Set("num_cache_nodes", c.NumCacheNodes)

	if c.ConfigurationEndpoint != nil {
		d.Set("port", c.ConfigurationEndpoint.Port)
		d.Set("configuration_endpoint", aws.String(fmt.Sprintf("%s:%d", aws.StringValue(c.ConfigurationEndpoint.Address), aws.Int64Value(c.ConfigurationEndpoint.Port))))
		d.Set("cluster_address", c.ConfigurationEndpoint.Address)
	} else if len(c.CacheNodes) > 0 {
		d.Set("port", c.CacheNodes[0].Endpoint.Port)
	}

	if c.ReplicationGroupId != nil {
		d.Set("replication_group_id", c.ReplicationGroupId)
	}

	if c.NotificationConfiguration != nil {
		if aws.StringValue(c.NotificationConfiguration.TopicStatus) == "active" {
			d.Set("notification_topic_arn", c.NotificationConfiguration.TopicArn)
		}
	}
	d.Set("availability_zone", c.PreferredAvailabilityZone)
	if aws.StringValue(c.PreferredAvailabilityZone) == "Multiple" {
		d.Set("az_mode", "cross-az")
	} else {
		d.Set("az_mode", "single-az")
	}

	if err := setCacheNodeData(d, c); err != nil {
		return err
	}

	d.Set("arn", c.ARN)

	tags, err := ListTags(conn, aws.StringValue(c.ARN))

	if err != nil && !verify.ErrorISOUnsupported(conn.PartitionID, err) {
		return fmt.Errorf("error listing tags for ElastiCache Cache Cluster (%s): %w", d.Id(), err)
	}

	if err != nil {
		log.Printf("[WARN] error listing tags for Elasticache Cache Cluster (%s): %s", d.Id(), err)
	}

	if tags != nil {
		tags = tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

		//lintignore:AWSR002
		if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
			return fmt.Errorf("error setting tags: %w", err)
		}

		if err := d.Set("tags_all", tags.Map()); err != nil {
			return fmt.Errorf("error setting tags_all: %w", err)
		}
	}

	return nil
}

func setFromCacheCluster(d *schema.ResourceData, c *elasticache.CacheCluster) error {
	d.Set("node_type", c.CacheNodeType)

	d.Set("engine", c.Engine)
	if aws.StringValue(c.Engine) == engineRedis {
		if err := setEngineVersionRedis(d, c.EngineVersion); err != nil {
			return err
		}
	} else {
		setEngineVersionMemcached(d, c.EngineVersion)
	}
	d.Set("auto_minor_version_upgrade", strconv.FormatBool(aws.BoolValue(c.AutoMinorVersionUpgrade)))

	d.Set("subnet_group_name", c.CacheSubnetGroupName)
	if err := d.Set("security_group_names", flattenSecurityGroupNames(c.CacheSecurityGroups)); err != nil {
		return fmt.Errorf("error setting security_group_names: %w", err)
	}
	if err := d.Set("security_group_ids", flattenSecurityGroupIDs(c.SecurityGroups)); err != nil {
		return fmt.Errorf("error setting security_group_ids: %w", err)
	}

	if c.CacheParameterGroup != nil {
		d.Set("parameter_group_name", c.CacheParameterGroup.CacheParameterGroupName)
	}

	d.Set("maintenance_window", c.PreferredMaintenanceWindow)

	return nil
}

func resourceClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ElastiCacheConn

	req := &elasticache.ModifyCacheClusterInput{
		CacheClusterId:   aws.String(d.Id()),
		ApplyImmediately: aws.Bool(d.Get("apply_immediately").(bool)),
	}

	requestUpdate := false
	if d.HasChange("security_group_ids") {
		if attr := d.Get("security_group_ids").(*schema.Set); attr.Len() > 0 {
			req.SecurityGroupIds = flex.ExpandStringSet(attr)
			requestUpdate = true
		}
	}

	if d.HasChange("parameter_group_name") {
		req.CacheParameterGroupName = aws.String(d.Get("parameter_group_name").(string))
		requestUpdate = true
	}

	if d.HasChange("log_delivery_configuration") {

		oldLogDeliveryConfig, newLogDeliveryConfig := d.GetChange("log_delivery_configuration")

		req.LogDeliveryConfigurations = []*elasticache.LogDeliveryConfigurationRequest{}
		logTypesToSubmit := make(map[string]bool)

		currentLogDeliveryConfig := newLogDeliveryConfig.(*schema.Set).List()
		for _, current := range currentLogDeliveryConfig {
			logDeliveryConfigurationRequest := expandLogDeliveryConfigurations(current.(map[string]interface{}))
			logTypesToSubmit[*logDeliveryConfigurationRequest.LogType] = true
			req.LogDeliveryConfigurations = append(req.LogDeliveryConfigurations, &logDeliveryConfigurationRequest)
		}

		previousLogDeliveryConfig := oldLogDeliveryConfig.(*schema.Set).List()
		for _, previous := range previousLogDeliveryConfig {
			logDeliveryConfigurationRequest := expandEmptyLogDeliveryConfigurations(previous.(map[string]interface{}))
			// if something was removed, send an empty request
			if !logTypesToSubmit[*logDeliveryConfigurationRequest.LogType] {
				req.LogDeliveryConfigurations = append(req.LogDeliveryConfigurations, &logDeliveryConfigurationRequest)
			}
		}

		requestUpdate = true
	}

	if d.HasChange("maintenance_window") {
		req.PreferredMaintenanceWindow = aws.String(d.Get("maintenance_window").(string))
		requestUpdate = true
	}

	if d.HasChange("notification_topic_arn") {
		v := d.Get("notification_topic_arn").(string)
		req.NotificationTopicArn = aws.String(v)
		if v == "" {
			inactive := "inactive"
			req.NotificationTopicStatus = &inactive
		}
		requestUpdate = true
	}

	if d.HasChange("engine_version") {
		req.EngineVersion = aws.String(d.Get("engine_version").(string))
		requestUpdate = true
	}

	if d.HasChange("auto_minor_version_upgrade") {
		v := d.Get("auto_minor_version_upgrade")
		if v, null, _ := nullable.Bool(v.(string)).Value(); !null {
			req.AutoMinorVersionUpgrade = aws.Bool(v)
		}
		requestUpdate = true
	}

	if d.HasChange("snapshot_window") {
		req.SnapshotWindow = aws.String(d.Get("snapshot_window").(string))
		requestUpdate = true
	}

	if d.HasChange("node_type") {
		req.CacheNodeType = aws.String(d.Get("node_type").(string))
		requestUpdate = true
	}

	if d.HasChange("snapshot_retention_limit") {
		req.SnapshotRetentionLimit = aws.Int64(int64(d.Get("snapshot_retention_limit").(int)))
		requestUpdate = true
	}

	if d.HasChange("az_mode") {
		req.AZMode = aws.String(d.Get("az_mode").(string))
		requestUpdate = true
	}

	if d.HasChange("num_cache_nodes") {
		oraw, nraw := d.GetChange("num_cache_nodes")
		o := oraw.(int)
		n := nraw.(int)
		if n < o {
			log.Printf("[INFO] Cluster %s is marked for Decreasing cache nodes from %d to %d", d.Id(), o, n)
			nodesToRemove := getCacheNodesToRemove(o, o-n)
			req.CacheNodeIdsToRemove = nodesToRemove
		} else {
			log.Printf("[INFO] Cluster %s is marked for increasing cache nodes from %d to %d", d.Id(), o, n)
			// SDK documentation for NewAvailabilityZones states:
			// The list of Availability Zones where the new Memcached cache nodes are created.
			//
			// This parameter is only valid when NumCacheNodes in the request is greater
			// than the sum of the number of active cache nodes and the number of cache
			// nodes pending creation (which may be zero). The number of Availability Zones
			// supplied in this list must match the cache nodes being added in this request.
			if v, ok := d.GetOk("preferred_availability_zones"); ok && len(v.([]interface{})) > 0 {
				// Here we check the list length to prevent a potential panic :)
				if len(v.([]interface{})) != n {
					return fmt.Errorf("length of preferred_availability_zones (%d) must match num_cache_nodes (%d)", len(v.([]interface{})), n)
				}
				req.NewAvailabilityZones = flex.ExpandStringList(v.([]interface{})[o:])
			}
		}

		req.NumCacheNodes = aws.Int64(int64(d.Get("num_cache_nodes").(int)))
		requestUpdate = true

	}

	if requestUpdate {
		log.Printf("[DEBUG] Modifying ElastiCache Cluster (%s), opts:\n%s", d.Id(), req)
		_, err := conn.ModifyCacheCluster(req)
		if err != nil {
			return fmt.Errorf("Error updating ElastiCache cluster (%s), error: %w", d.Id(), err)
		}

		_, err = waitCacheClusterAvailable(conn, d.Id(), CacheClusterUpdatedTimeout)
		if err != nil {
			return fmt.Errorf("error waiting for ElastiCache Cache Cluster (%s) to update: %w", d.Id(), err)
		}
	}

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		err := UpdateTags(conn, d.Get("arn").(string), o, n)

		// ISO partitions may not support tagging, giving error
		if err != nil {
			if v, ok := d.GetOk("tags"); (ok && len(v.(map[string]interface{})) > 0) || !verify.ErrorISOUnsupported(conn.PartitionID, err) {
				// explicitly setting tags or not an iso-unsupported error
				return fmt.Errorf("failed updating ElastiCache Cache Cluster (%s) tags: %w", d.Get("arn").(string), err)
			}

			// no non-default tags and iso-unsupported error
			log.Printf("[WARN] failed updating tags for ElastiCache Cache Cluster (%s): %s", d.Get("arn").(string), err)
		}
	}

	return resourceClusterRead(d, meta)
}

func getCacheNodesToRemove(oldNumberOfNodes int, cacheNodesToRemove int) []*string {
	nodesIdsToRemove := []*string{}
	for i := oldNumberOfNodes; i > oldNumberOfNodes-cacheNodesToRemove && i > 0; i-- {
		s := fmt.Sprintf("%04d", i)
		nodesIdsToRemove = append(nodesIdsToRemove, &s)
	}

	return nodesIdsToRemove
}

func setCacheNodeData(d *schema.ResourceData, c *elasticache.CacheCluster) error {
	sortedCacheNodes := make([]*elasticache.CacheNode, len(c.CacheNodes))
	copy(sortedCacheNodes, c.CacheNodes)
	sort.Sort(byCacheNodeId(sortedCacheNodes))

	cacheNodeData := make([]map[string]interface{}, 0, len(sortedCacheNodes))

	for _, node := range sortedCacheNodes {
		if node.CacheNodeId == nil || node.Endpoint == nil || node.Endpoint.Address == nil || node.Endpoint.Port == nil || node.CustomerAvailabilityZone == nil {
			return fmt.Errorf("Unexpected nil pointer in: %s", node)
		}
		cacheNodeData = append(cacheNodeData, map[string]interface{}{
			"id":                aws.StringValue(node.CacheNodeId),
			"address":           aws.StringValue(node.Endpoint.Address),
			"port":              aws.Int64Value(node.Endpoint.Port),
			"availability_zone": aws.StringValue(node.CustomerAvailabilityZone),
		})
	}

	return d.Set("cache_nodes", cacheNodeData)
}

type byCacheNodeId []*elasticache.CacheNode

func (b byCacheNodeId) Len() int      { return len(b) }
func (b byCacheNodeId) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b byCacheNodeId) Less(i, j int) bool {
	return b[i].CacheNodeId != nil && b[j].CacheNodeId != nil &&
		aws.StringValue(b[i].CacheNodeId) < aws.StringValue(b[j].CacheNodeId)
}

func resourceClusterDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ElastiCacheConn

	var finalSnapshotID = d.Get("final_snapshot_identifier").(string)
	err := DeleteCacheCluster(conn, d.Id(), finalSnapshotID)
	if err != nil {
		if tfawserr.ErrCodeEquals(err, elasticache.ErrCodeCacheClusterNotFoundFault) {
			return nil
		}
		return fmt.Errorf("error deleting ElastiCache Cache Cluster (%s): %w", d.Id(), err)
	}
	_, err = WaitCacheClusterDeleted(conn, d.Id(), CacheClusterDeletedTimeout)
	if err != nil {
		return fmt.Errorf("error waiting for ElastiCache Cache Cluster (%s) to be deleted: %w", d.Id(), err)
	}

	return nil
}

func createCacheCluster(conn *elasticache.ElastiCache, input *elasticache.CreateCacheClusterInput) (string, string, error) {
	log.Printf("[DEBUG] Creating ElastiCache Cache Cluster: %s", input)
	output, err := conn.CreateCacheCluster(input)

	// Some partitions may not support tag-on-create
	if input.Tags != nil && verify.ErrorISOUnsupported(conn.PartitionID, err) {
		log.Printf("[WARN] failed creating ElastiCache Cache Cluster with tags: %s. Trying create without tags.", err)

		input.Tags = nil
		output, err = conn.CreateCacheCluster(input)
	}

	if err != nil {
		return "", "", err
	}

	if output == nil || output.CacheCluster == nil {
		return "", "", errors.New("missing cluster ID after creation")
	}
	// Elasticache always retains the id in lower case, so we have to
	// mimic that or else we won't be able to refresh a resource whose
	// name contained uppercase characters.
	return strings.ToLower(aws.StringValue(output.CacheCluster.CacheClusterId)), aws.StringValue(output.CacheCluster.ARN), nil
}

func DeleteCacheCluster(conn *elasticache.ElastiCache, cacheClusterID string, finalSnapshotID string) error {
	input := &elasticache.DeleteCacheClusterInput{
		CacheClusterId: aws.String(cacheClusterID),
	}
	if finalSnapshotID != "" {
		input.FinalSnapshotIdentifier = aws.String(finalSnapshotID)
	}

	log.Printf("[DEBUG] Deleting ElastiCache Cache Cluster: %s", input)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := conn.DeleteCacheCluster(input)
		if err != nil {
			if tfawserr.ErrMessageContains(err, elasticache.ErrCodeInvalidCacheClusterStateFault, "serving as primary") {
				return resource.NonRetryableError(err)
			}
			if tfawserr.ErrMessageContains(err, elasticache.ErrCodeInvalidCacheClusterStateFault, "only member of a replication group") {
				return resource.NonRetryableError(err)
			}
			// The cluster may be just snapshotting, so we retry until it's ready for deletion
			if tfawserr.ErrCodeEquals(err, elasticache.ErrCodeInvalidCacheClusterStateFault) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if tfresource.TimedOut(err) {
		_, err = conn.DeleteCacheCluster(input)
	}

	return err
}
