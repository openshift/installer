package alicloud

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudAlikafkaSaslAcls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAlikafkaSaslAclsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"username": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
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
				ValidateFunc: validation.ValidateRegexp,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"acls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"acl_resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"acl_resource_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"acl_resource_pattern_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"acl_operation_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudAlikafkaSaslAclsRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}

	request := alikafka.CreateDescribeAclsRequest()
	request.InstanceId = d.Get("instance_id").(string)
	request.RegionId = client.RegionId
	request.Username = d.Get("username").(string)
	request.AclResourceType = d.Get("acl_resource_type").(string)
	request.AclResourceName = d.Get("acl_resource_name").(string)

	raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
		return alikafkaClient.DescribeAcls(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"BIZ_SUBSCRIPTION_NOT_FOUND", "BIZ_TOPIC_NOT_FOUND"}) {
			var emptyValue []alikafka.KafkaAclVO
			return alikafkaSaslAclsDecriptionAttributes(d, emptyValue, meta)
		}
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_alikafka_sasl_acls", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*alikafka.DescribeAclsResponse)

	return alikafkaSaslAclsDecriptionAttributes(d, response.KafkaAclList.KafkaAclVO, meta)
}

func alikafkaSaslAclsDecriptionAttributes(d *schema.ResourceData, kafkaAclsInfo []alikafka.KafkaAclVO, meta interface{}) error {

	var names []string
	var s []map[string]interface{}

	for _, item := range kafkaAclsInfo {
		mapping := map[string]interface{}{
			"username":                  item.Username,
			"acl_resource_type":         item.AclResourceType,
			"acl_resource_name":         item.AclResourceName,
			"acl_resource_pattern_type": item.AclResourcePatternType,
			"host":                      item.Host,
			"acl_operation_type":        item.AclOperationType,
		}

		name := fmt.Sprintf("%s:%s:%s:%s:%s", item.Username, item.AclResourceType, item.AclResourceName, item.AclResourcePatternType, item.AclOperationType)
		names = append(names, name)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(names))

	if err := d.Set("acls", s); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
