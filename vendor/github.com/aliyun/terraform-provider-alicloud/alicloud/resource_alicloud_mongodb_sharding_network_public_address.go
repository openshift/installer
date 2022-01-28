package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudMongodbShardingNetworkPublicAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMongodbShardingNetworkPublicAddressCreate,
		Read:   resourceAlicloudMongodbShardingNetworkPublicAddressRead,
		Delete: resourceAlicloudMongodbShardingNetworkPublicAddressDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network_address": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expired_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudMongodbShardingNetworkPublicAddressCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AllocatePublicNetworkAddress"
	request := make(map[string]interface{})
	conn, err := client.NewDdsClient()
	if err != nil {
		return WrapError(err)
	}
	request["DBInstanceId"] = d.Get("db_instance_id")
	request["NodeId"] = d.Get("node_id")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_mongodb_sharding_network_public_address", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(request["DBInstanceId"], ":", request["NodeId"]))

	MongoDBService := MongoDBService{client}
	nodeType, err := MongoDBService.DescribeShardingNodeType(d.Id())
	if err != nil {
		return WrapError(err)
	}

	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, MongoDBService.MongodbShardingNetworkPublicAddressStateRefreshFunc(d.Id(), nodeType, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudMongodbShardingNetworkPublicAddressRead(d, meta)
}
func resourceAlicloudMongodbShardingNetworkPublicAddressRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	MongoDBService := MongoDBService{client}
	object, err := MongoDBService.DescribeMongodbShardingNetworkPublicAddress(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_mongodb_sharding_network_public_address MongoDBService.DescribeMongodbShardingNetworkPublicAddress Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("db_instance_id", parts[0])
	d.Set("node_id", parts[1])
	s := make([]map[string]interface{}, 0)
	for _, item := range object["NetworkAddress"].([]map[string]interface{}) {
		mapping := map[string]interface{}{
			"expired_time":    item["ExpiredTime"],
			"ip_address":      item["IPAddress"],
			"network_address": item["NetworkAddress"],
			"network_type":    item["NetworkType"],
			"node_type":       item["NodeType"],
			"port":            item["Port"],
			"role":            item["Role"],
			"vpc_id":          item["VPCId"],
			"vswitch_id":      item["VswitchId"],
			"node_id":         parts[1],
		}
		s = append(s, mapping)
	}
	d.Set("network_address", s)
	return nil
}
func resourceAlicloudMongodbShardingNetworkPublicAddressDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "ReleasePublicNetworkAddress"
	var response map[string]interface{}
	conn, err := client.NewDdsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DBInstanceId": parts[0],
		"NodeId":       parts[1],
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	MongoDBService := MongoDBService{client}
	nodeType, err := MongoDBService.DescribeShardingNodeType(d.Id())
	if err != nil {
		return WrapError(err)
	}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, MongoDBService.MongodbShardingNetworkPublicAddressStateRefreshFunc(d.Id(), nodeType, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
