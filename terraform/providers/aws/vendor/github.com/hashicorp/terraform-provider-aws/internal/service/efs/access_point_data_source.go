package efs

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/efs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
)

func DataSourceAccessPoint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAccessPointRead,

		Schema: map[string]*schema.Schema{
			"access_point_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"file_system_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"file_system_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"posix_user": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"uid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"secondary_gids": {
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeInt},
							Set:      schema.HashInt,
							Computed: true,
						},
					},
				},
			},
			"root_directory": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"owner_gid": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"owner_uid": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"permissions": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"tags": tftags.TagsSchema(),
		},
	}
}

func dataSourceAccessPointRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EFSConn
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	resp, err := conn.DescribeAccessPoints(&efs.DescribeAccessPointsInput{
		AccessPointId: aws.String(d.Get("access_point_id").(string)),
	})
	if err != nil {
		return fmt.Errorf("Error reading EFS access point %s: %w", d.Id(), err)
	}
	if len(resp.AccessPoints) != 1 {
		return fmt.Errorf("Search returned %d results, please revise so only one is returned", len(resp.AccessPoints))
	}

	ap := resp.AccessPoints[0]

	log.Printf("[DEBUG] Found EFS access point: %#v", ap)

	d.SetId(aws.StringValue(ap.AccessPointId))

	fsARN := arn.ARN{
		AccountID: meta.(*conns.AWSClient).AccountID,
		Partition: meta.(*conns.AWSClient).Partition,
		Region:    meta.(*conns.AWSClient).Region,
		Resource:  fmt.Sprintf("file-system/%s", aws.StringValue(ap.FileSystemId)),
		Service:   "elasticfilesystem",
	}.String()

	d.Set("file_system_arn", fsARN)
	d.Set("file_system_id", ap.FileSystemId)
	d.Set("arn", ap.AccessPointArn)
	d.Set("owner_id", ap.OwnerId)

	if err := d.Set("posix_user", flattenAccessPointPOSIXUser(ap.PosixUser)); err != nil {
		return fmt.Errorf("error setting posix user: %w", err)
	}

	if err := d.Set("root_directory", flattenAccessPointRootDirectory(ap.RootDirectory)); err != nil {
		return fmt.Errorf("error setting root directory: %w", err)
	}

	if err := d.Set("tags", KeyValueTags(ap.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	return nil
}
