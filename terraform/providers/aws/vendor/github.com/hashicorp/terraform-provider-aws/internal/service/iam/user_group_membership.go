package iam

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func ResourceUserGroupMembership() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserGroupMembershipCreate,
		Read:   resourceUserGroupMembershipRead,
		Update: resourceUserGroupMembershipUpdate,
		Delete: resourceUserGroupMembershipDelete,
		Importer: &schema.ResourceImporter{
			State: resourceUserGroupMembershipImport,
		},

		Schema: map[string]*schema.Schema{
			"user": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"groups": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceUserGroupMembershipCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn

	user := d.Get("user").(string)
	groupList := flex.ExpandStringSet(d.Get("groups").(*schema.Set))

	if err := addUserToGroups(conn, user, groupList); err != nil {
		return err
	}

	//lintignore:R015 // Allow legacy unstable ID usage in managed resource
	d.SetId(resource.UniqueId())

	return resourceUserGroupMembershipRead(d, meta)
}

func resourceUserGroupMembershipRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn

	user := d.Get("user").(string)
	groups := d.Get("groups").(*schema.Set)

	input := &iam.ListGroupsForUserInput{
		UserName: aws.String(user),
	}

	var gl []string

	err := resource.Retry(propagationTimeout, func() *resource.RetryError {
		err := conn.ListGroupsForUserPages(input, func(page *iam.ListGroupsForUserOutput, lastPage bool) bool {
			if page == nil {
				return !lastPage
			}

			for _, group := range page.Groups {
				if groups.Contains(aws.StringValue(group.GroupName)) {
					gl = append(gl, aws.StringValue(group.GroupName))
				}
			}

			return !lastPage
		})

		if d.IsNewResource() && tfawserr.ErrCodeEquals(err, iam.ErrCodeNoSuchEntityException) {
			return resource.RetryableError(err)
		}

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if tfresource.TimedOut(err) {
		err = conn.ListGroupsForUserPages(input, func(page *iam.ListGroupsForUserOutput, lastPage bool) bool {
			if page == nil {
				return !lastPage
			}

			for _, group := range page.Groups {
				if groups.Contains(aws.StringValue(group.GroupName)) {
					gl = append(gl, aws.StringValue(group.GroupName))
				}
			}

			return !lastPage
		})
	}

	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, iam.ErrCodeNoSuchEntityException) {
		log.Printf("[WARN] IAM User Group Membership (%s) not found, removing from state", user)
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading IAM User Group Membership (%s): %w", user, err)
	}

	if err := d.Set("groups", gl); err != nil {
		return fmt.Errorf("Error setting group list from IAM (%s), error: %s", user, err)
	}

	return nil
}

func resourceUserGroupMembershipUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn

	if d.HasChange("groups") {
		user := d.Get("user").(string)

		o, n := d.GetChange("groups")
		if o == nil {
			o = new(schema.Set)
		}
		if n == nil {
			n = new(schema.Set)
		}

		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := flex.ExpandStringSet(os.Difference(ns))
		add := flex.ExpandStringSet(ns.Difference(os))

		if err := removeUserFromGroups(conn, user, remove); err != nil {
			return err
		}

		if err := addUserToGroups(conn, user, add); err != nil {
			return err
		}
	}

	return resourceUserGroupMembershipRead(d, meta)
}

func resourceUserGroupMembershipDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).IAMConn
	user := d.Get("user").(string)
	groups := flex.ExpandStringSet(d.Get("groups").(*schema.Set))

	err := removeUserFromGroups(conn, user, groups)
	return err
}

func removeUserFromGroups(conn *iam.IAM, user string, groups []*string) error {
	for _, group := range groups {
		_, err := conn.RemoveUserFromGroup(&iam.RemoveUserFromGroupInput{
			UserName:  &user,
			GroupName: group,
		})
		if err != nil {
			if tfawserr.ErrCodeEquals(err, iam.ErrCodeNoSuchEntityException) {
				continue
			}
			return err
		}
	}

	return nil
}

func addUserToGroups(conn *iam.IAM, user string, groups []*string) error {
	for _, group := range groups {
		_, err := conn.AddUserToGroup(&iam.AddUserToGroupInput{
			UserName:  &user,
			GroupName: group,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceUserGroupMembershipImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	idParts := strings.Split(d.Id(), "/")
	if len(idParts) < 2 {
		return nil, fmt.Errorf("unexpected format of ID (%q), expected <user-name>/<group-name1>/...", d.Id())
	}

	userName := idParts[0]
	groupList := idParts[1:]

	d.Set("user", userName)
	d.Set("groups", groupList)

	//lintignore:R015 // Allow legacy unstable ID usage in managed resource
	d.SetId(resource.UniqueId())

	return []*schema.ResourceData{d}, nil
}
