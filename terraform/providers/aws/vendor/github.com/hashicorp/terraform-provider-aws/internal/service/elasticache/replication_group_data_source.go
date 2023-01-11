package elasticache

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
)

func DataSourceReplicationGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceReplicationGroupRead,
		Schema: map[string]*schema.Schema{
			"replication_group_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateReplicationGroupID,
			},
			"replication_group_description": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Use description instead",
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auth_token_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"automatic_failover_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"configuration_endpoint_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"primary_endpoint_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reader_endpoint_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"num_cache_clusters": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"num_node_groups": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"number_cache_clusters": {
				Type:       schema.TypeInt,
				Computed:   true,
				Deprecated: "Use num_cache_clusters instead",
			},
			"member_clusters": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"multi_az_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"node_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"replicas_per_node_group": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"log_delivery_configuration": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"destination": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_format": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"snapshot_window": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_retention_limit": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceReplicationGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ElastiCacheConn

	groupID := d.Get("replication_group_id").(string)

	rg, err := FindReplicationGroupByID(conn, groupID)
	if err != nil {
		return fmt.Errorf("error reading ElastiCache Replication Group (%s): %w", groupID, err)
	}

	d.SetId(aws.StringValue(rg.ReplicationGroupId))
	d.Set("description", rg.Description)
	d.Set("replication_group_description", rg.Description)
	d.Set("arn", rg.ARN)
	d.Set("auth_token_enabled", rg.AuthTokenEnabled)

	if rg.AutomaticFailover != nil {
		switch aws.StringValue(rg.AutomaticFailover) {
		case elasticache.AutomaticFailoverStatusDisabled, elasticache.AutomaticFailoverStatusDisabling:
			d.Set("automatic_failover_enabled", false)
		case elasticache.AutomaticFailoverStatusEnabled, elasticache.AutomaticFailoverStatusEnabling:
			d.Set("automatic_failover_enabled", true)
		}
	}

	if rg.MultiAZ != nil {
		switch strings.ToLower(aws.StringValue(rg.MultiAZ)) {
		case elasticache.MultiAZStatusEnabled:
			d.Set("multi_az_enabled", true)
		case elasticache.MultiAZStatusDisabled:
			d.Set("multi_az_enabled", false)
		default:
			log.Printf("Unknown MultiAZ state %q", aws.StringValue(rg.MultiAZ))
		}
	}

	if rg.ConfigurationEndpoint != nil {
		d.Set("port", rg.ConfigurationEndpoint.Port)
		d.Set("configuration_endpoint_address", rg.ConfigurationEndpoint.Address)
	} else {
		if rg.NodeGroups == nil {
			d.SetId("")
			return fmt.Errorf("ElastiCache Replication Group (%s) doesn't have node groups", aws.StringValue(rg.ReplicationGroupId))
		}
		d.Set("port", rg.NodeGroups[0].PrimaryEndpoint.Port)
		d.Set("primary_endpoint_address", rg.NodeGroups[0].PrimaryEndpoint.Address)
		d.Set("reader_endpoint_address", rg.NodeGroups[0].ReaderEndpoint.Address)
	}

	d.Set("num_cache_clusters", len(rg.MemberClusters))
	d.Set("number_cache_clusters", len(rg.MemberClusters))
	if err := d.Set("member_clusters", flex.FlattenStringList(rg.MemberClusters)); err != nil {
		return fmt.Errorf("error setting member_clusters: %w", err)
	}
	d.Set("node_type", rg.CacheNodeType)
	d.Set("num_node_groups", len(rg.NodeGroups))
	d.Set("replicas_per_node_group", len(rg.NodeGroups[0].NodeGroupMembers)-1)
	d.Set("log_delivery_configuration", flattenLogDeliveryConfigurations(rg.LogDeliveryConfigurations))
	d.Set("snapshot_window", rg.SnapshotWindow)
	d.Set("snapshot_retention_limit", rg.SnapshotRetentionLimit)
	return nil
}
