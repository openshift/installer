package dms

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKDataSource("aws_dms_replication_task")
func DataSourceReplicationTask() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceReplicationTaskRead,

		Schema: map[string]*schema.Schema{
			"cdc_start_position": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cdc_start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"migration_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"replication_instance_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"replication_task_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"replication_task_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"replication_task_settings": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_endpoint_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"start_replication_task": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"table_mappings": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_endpoint_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tftags.TagsSchemaComputed(),
		},
	}
}

const (
	DSNameReplicationTask = "Replication Task Data Source"
)

func dataSourceReplicationTaskRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).DMSConn(ctx)
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	taskID := d.Get("replication_task_id").(string)

	task, err := FindReplicationTaskByID(ctx, conn, taskID)
	if err != nil {
		return create.DiagError(names.DMS, create.ErrActionReading, DSNameReplicationTask, d.Id(), err)
	}

	d.SetId(aws.StringValue(task.ReplicationTaskIdentifier))
	d.Set("cdc_start_position", task.CdcStartPosition)
	d.Set("migration_type", task.MigrationType)
	d.Set("replication_instance_arn", task.ReplicationInstanceArn)
	d.Set("replication_task_arn", task.ReplicationTaskArn)
	d.Set("replication_task_id", task.ReplicationTaskIdentifier)
	d.Set("source_endpoint_arn", task.SourceEndpointArn)
	d.Set("status", task.Status)
	d.Set("table_mappings", task.TableMappings)
	d.Set("target_endpoint_arn", task.TargetEndpointArn)

	settings, err := replicationTaskRemoveReadOnlySettings(aws.StringValue(task.ReplicationTaskSettings))
	if err != nil {
		return create.DiagError(names.DMS, create.ErrActionReading, DSNameReplicationTask, d.Id(), err)
	}

	d.Set("replication_task_settings", settings)

	tags, err := ListTags(ctx, conn, aws.StringValue(task.ReplicationTaskArn))
	if err != nil {
		return create.DiagError(names.DMS, create.ErrActionReading, DSNameReplicationTask, d.Id(), err)
	}

	tags = tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return create.DiagError(names.DMS, create.ErrActionSetting, DSNameReplicationTask, d.Id(), err)
	}

	return nil
}
