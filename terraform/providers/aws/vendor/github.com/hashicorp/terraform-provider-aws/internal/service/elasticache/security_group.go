package elasticache

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func ResourceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityGroupCreate,
		Read:   resourceSecurityGroupRead,
		Delete: resourceSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "Managed by Terraform",
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_names": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},

		DeprecationMessage: `With the retirement of EC2-Classic the aws_elasticache_security_group resource has been deprecated and will be removed in a future version.`,
	}
}

func resourceSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	return errors.New(`with the retirement of EC2-Classic no new ElastiCache Security Groups can be created`)
}

func resourceSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ElastiCacheConn
	req := &elasticache.DescribeCacheSecurityGroupsInput{
		CacheSecurityGroupName: aws.String(d.Id()),
	}

	res, err := conn.DescribeCacheSecurityGroups(req)
	if err != nil {
		return err
	}
	if len(res.CacheSecurityGroups) == 0 {
		return fmt.Errorf("Error missing %v", d.Id())
	}

	var group *elasticache.CacheSecurityGroup
	for _, g := range res.CacheSecurityGroups {
		log.Printf("[DEBUG] CacheSecurityGroupName: %v, id: %v", g.CacheSecurityGroupName, d.Id())
		if aws.StringValue(g.CacheSecurityGroupName) == d.Id() {
			group = g
		}
	}
	if group == nil {
		return fmt.Errorf("Error retrieving cache security group: %v", res)
	}

	d.Set("name", group.CacheSecurityGroupName)
	d.Set("description", group.Description)

	sgNames := make([]string, 0, len(group.EC2SecurityGroups))
	for _, sg := range group.EC2SecurityGroups {
		sgNames = append(sgNames, *sg.EC2SecurityGroupName)
	}
	d.Set("security_group_names", sgNames)

	return nil
}

func resourceSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ElastiCacheConn

	log.Printf("[DEBUG] Cache security group delete: %s", d.Id())

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := conn.DeleteCacheSecurityGroup(&elasticache.DeleteCacheSecurityGroupInput{
			CacheSecurityGroupName: aws.String(d.Id()),
		})
		if err != nil {
			apierr, ok := err.(awserr.Error)
			if !ok {
				return resource.RetryableError(err)
			}
			log.Printf("[DEBUG] APIError.Code: %v", apierr.Code())
			switch apierr.Code() {
			case "InvalidCacheSecurityGroupState":
				return resource.RetryableError(err)
			case "DependencyViolation":
				// If it is a dependency violation, we want to retry
				return resource.RetryableError(err)
			default:
				return resource.NonRetryableError(err)
			}
		}
		return nil
	})

	if tfresource.TimedOut(err) {
		_, err = conn.DeleteCacheSecurityGroup(&elasticache.DeleteCacheSecurityGroupInput{
			CacheSecurityGroupName: aws.String(d.Id()),
		})
	}

	return err
}
