package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudAlikafkaSaslAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAlikafkaSaslAclCreate,
		Read:   resourceAlicloudAlikafkaSaslAclRead,
		Delete: resourceAlicloudAlikafkaSaslAclDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"username": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"acl_resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Group", "Topic"}, false),
			},
			"acl_resource_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"acl_resource_pattern_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"LITERAL", "PREFIXED"}, false),
			},
			"acl_operation_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Read", "Write"}, false),
			},
			"host": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudAlikafkaSaslAclCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}

	instanceId := d.Get("instance_id").(string)
	regionId := client.RegionId
	username := d.Get("username").(string)
	aclResourceType := d.Get("acl_resource_type").(string)
	aclResourceName := d.Get("acl_resource_name").(string)
	aclResourcePatternType := d.Get("acl_resource_pattern_type").(string)
	aclOperationType := d.Get("acl_operation_type").(string)

	request := alikafka.CreateCreateAclRequest()
	request.InstanceId = instanceId
	request.RegionId = regionId
	request.Username = username
	request.AclResourceType = aclResourceType
	request.AclResourceName = aclResourceName
	request.AclResourcePatternType = aclResourcePatternType
	request.AclOperationType = aclOperationType

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.CreateAcl(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				time.Sleep(2 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alikafka_sasl_acl", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	// Server may have cache, sleep a while.
	time.Sleep(60 * time.Second)
	d.SetId(fmt.Sprintf("%s:%s:%s:%s:%s:%s", instanceId, username, aclResourceType, aclResourceName, aclResourcePatternType, aclOperationType))
	return resourceAlicloudAlikafkaSaslAclRead(d, meta)
}

func resourceAlicloudAlikafkaSaslAclRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}

	parts, err := ParseResourceId(d.Id(), 6)
	if err != nil {
		return WrapError(err)
	}
	object, err := alikafkaService.DescribeAlikafkaSaslAcl(d.Id())
	if err != nil {
		// Handle exceptions
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", parts[0])
	d.Set("username", object.Username)
	d.Set("acl_resource_type", object.AclResourceType)
	d.Set("acl_resource_name", object.AclResourceName)
	d.Set("acl_resource_pattern_type", object.AclResourcePatternType)
	d.Set("acl_operation_type", object.AclOperationType)
	d.Set("host", object.Host)

	return nil
}

func resourceAlicloudAlikafkaSaslAclDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}

	parts, err := ParseResourceId(d.Id(), 6)
	if err != nil {
		return WrapError(err)
	}
	instanceId := parts[0]
	username := parts[1]
	aclResourceType := parts[2]
	aclResourceName := parts[3]
	aclResourcePatternType := parts[4]
	aclOperationType := parts[5]

	request := alikafka.CreateDeleteAclRequest()
	request.RegionId = client.RegionId
	request.InstanceId = instanceId
	request.Username = username
	request.AclResourceType = aclResourceType
	request.AclResourceName = aclResourceName
	request.AclResourcePatternType = aclResourcePatternType
	request.AclOperationType = aclOperationType

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.DeleteAcl(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	// Server may have cache, sleep a while.
	time.Sleep(60 * time.Second)
	return WrapError(alikafkaService.WaitForAlikafkaSaslAcl(d.Id(), Deleted, DefaultTimeoutMedium))
}
