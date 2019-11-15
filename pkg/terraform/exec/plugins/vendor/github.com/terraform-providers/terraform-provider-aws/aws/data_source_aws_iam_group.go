package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAwsIAMGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAwsIAMGroupRead,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAwsIAMGroupRead(d *schema.ResourceData, meta interface{}) error {
	iamconn := meta.(*AWSClient).iamconn

	groupName := d.Get("group_name").(string)

	req := &iam.GetGroupInput{
		GroupName: aws.String(groupName),
	}

	log.Printf("[DEBUG] Reading IAM Group: %s", req)
	resp, err := iamconn.GetGroup(req)
	if err != nil {
		return fmt.Errorf("Error getting group: %s", err)
	}
	if resp == nil {
		return fmt.Errorf("no IAM group found")
	}

	group := resp.Group

	d.SetId(*group.GroupId)
	d.Set("arn", group.Arn)
	d.Set("path", group.Path)
	d.Set("group_id", group.GroupId)
	if err := d.Set("users", dataSourceUsersRead(resp.Users)); err != nil {
		return fmt.Errorf("error setting users: %s", err)
	}

	return nil
}

func dataSourceUsersRead(iamUsers []*iam.User) []map[string]interface{} {
	users := make([]map[string]interface{}, 0, len(iamUsers))
	for _, i := range iamUsers {
		u := make(map[string]interface{})
		u["arn"] = aws.StringValue(i.Arn)
		u["user_id"] = aws.StringValue(i.UserId)
		u["user_name"] = aws.StringValue(i.UserName)
		u["path"] = aws.StringValue(i.Path)
		users = append(users, u)
	}
	return users
}
