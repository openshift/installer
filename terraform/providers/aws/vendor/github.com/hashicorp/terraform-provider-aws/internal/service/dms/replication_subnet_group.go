package dms

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	dms "github.com/aws/aws-sdk-go/service/databasemigrationservice"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceReplicationSubnetGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceReplicationSubnetGroupCreate,
		Read:   resourceReplicationSubnetGroupRead,
		Update: resourceReplicationSubnetGroupUpdate,
		Delete: resourceReplicationSubnetGroupDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"replication_subnet_group_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"replication_subnet_group_description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"replication_subnet_group_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validReplicationSubnetGroupID,
			},
			"subnet_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Required: true,
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceReplicationSubnetGroupCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).DMSConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	request := &dms.CreateReplicationSubnetGroupInput{
		ReplicationSubnetGroupIdentifier:  aws.String(d.Get("replication_subnet_group_id").(string)),
		ReplicationSubnetGroupDescription: aws.String(d.Get("replication_subnet_group_description").(string)),
		SubnetIds:                         flex.ExpandStringSet(d.Get("subnet_ids").(*schema.Set)),
		Tags:                              Tags(tags.IgnoreAWS()),
	}

	log.Println("[DEBUG] DMS create replication subnet group:", request)

	err := resource.Retry(propagationTimeout, func() *resource.RetryError {
		_, err := conn.CreateReplicationSubnetGroup(request)

		if tfawserr.ErrCodeEquals(err, dms.ErrCodeAccessDeniedFault) {
			return resource.RetryableError(err)
		}

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if tfresource.TimedOut(err) {
		_, err = conn.CreateReplicationSubnetGroup(request)

		if err != nil {
			return err
		}
	}

	d.SetId(d.Get("replication_subnet_group_id").(string))
	return resourceReplicationSubnetGroupRead(d, meta)
}

func resourceReplicationSubnetGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).DMSConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	response, err := conn.DescribeReplicationSubnetGroups(&dms.DescribeReplicationSubnetGroupsInput{
		Filters: []*dms.Filter{
			{
				Name:   aws.String("replication-subnet-group-id"),
				Values: []*string{aws.String(d.Id())}, // Must use d.Id() to work with import.
			},
		},
	})
	if err != nil {
		return err
	}
	if len(response.ReplicationSubnetGroups) == 0 {
		d.SetId("")
		return nil
	}

	// The AWS API for DMS subnet groups does not return the ARN which is required to
	// retrieve tags. This ARN can be built.
	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Service:   "dms",
		Region:    meta.(*conns.AWSClient).Region,
		AccountID: meta.(*conns.AWSClient).AccountID,
		Resource:  fmt.Sprintf("subgrp:%s", d.Id()),
	}.String()
	d.Set("replication_subnet_group_arn", arn)

	err = resourceReplicationSubnetGroupSetState(d, response.ReplicationSubnetGroups[0])
	if err != nil {
		return err
	}

	tags, err := ListTags(conn, arn)

	if err != nil {
		return fmt.Errorf("error listing tags for DMS Replication Subnet Group (%s): %s", arn, err)
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

func resourceReplicationSubnetGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).DMSConn

	// Updates to subnet groups are only valid when sending SubnetIds even if there are no
	// changes to SubnetIds.
	request := &dms.ModifyReplicationSubnetGroupInput{
		ReplicationSubnetGroupIdentifier: aws.String(d.Get("replication_subnet_group_id").(string)),
		SubnetIds:                        flex.ExpandStringSet(d.Get("subnet_ids").(*schema.Set)),
	}

	if d.HasChange("replication_subnet_group_description") {
		request.ReplicationSubnetGroupDescription = aws.String(d.Get("replication_subnet_group_description").(string))
	}

	if d.HasChange("tags_all") {
		arn := d.Get("replication_subnet_group_arn").(string)
		o, n := d.GetChange("tags_all")

		if err := UpdateTags(conn, arn, o, n); err != nil {
			return fmt.Errorf("error updating DMS Replication Subnet Group (%s) tags: %s", arn, err)
		}
	}

	log.Println("[DEBUG] DMS update replication subnet group:", request)

	_, err := conn.ModifyReplicationSubnetGroup(request)
	if err != nil {
		return err
	}

	return resourceReplicationSubnetGroupRead(d, meta)
}

func resourceReplicationSubnetGroupDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).DMSConn

	request := &dms.DeleteReplicationSubnetGroupInput{
		ReplicationSubnetGroupIdentifier: aws.String(d.Get("replication_subnet_group_id").(string)),
	}

	log.Printf("[DEBUG] DMS delete replication subnet group: %#v", request)

	_, err := conn.DeleteReplicationSubnetGroup(request)
	return err
}

func resourceReplicationSubnetGroupSetState(d *schema.ResourceData, group *dms.ReplicationSubnetGroup) error {
	d.SetId(aws.StringValue(group.ReplicationSubnetGroupIdentifier))

	subnet_ids := []string{}
	for _, subnet := range group.Subnets {
		subnet_ids = append(subnet_ids, aws.StringValue(subnet.SubnetIdentifier))
	}

	d.Set("replication_subnet_group_description", group.ReplicationSubnetGroupDescription)
	d.Set("replication_subnet_group_id", group.ReplicationSubnetGroupIdentifier)
	d.Set("subnet_ids", subnet_ids)
	d.Set("vpc_id", group.VpcId)

	return nil
}
