package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudMongodbShardingNetworkPrivateAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMongodbShardingNetworkPrivateAddressCreate,
		Read:   resourceAlicloudMongodbShardingNetworkPrivateAddressRead,
		Update: resourceAlicloudMongodbShardingNetworkPrivateAddressUpdate,
		Delete: resourceAlicloudMongodbShardingNetworkPrivateAddressDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"account_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"account_password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[\w!@#$%^&*()_+=]{6,32}$`), "The account password must be 6 to 32 characters in length, and can contain letters, digits, and special characters（!@#$%^&*()_+-=)."),
			},
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceAlicloudMongodbShardingNetworkPrivateAddressCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AllocateNodePrivateNetworkAddress"
	request := make(map[string]interface{})
	conn, err := client.NewDdsClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("account_name"); ok {
		request["AccountName"] = v
	}
	if v, ok := d.GetOk("account_password"); ok {
		request["AccountPassword"] = v
	}
	request["DBInstanceId"] = d.Get("db_instance_id")
	request["NodeId"] = d.Get("node_id")
	request["RegionId"] = client.RegionId
	request["ZoneId"] = d.Get("zone_id")
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_mongodb_sharding_network_private_address", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["DBInstanceId"], ":", request["NodeId"]))
	ddsService := MongoDBService{client}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, ddsService.RdsMongodbDBInstanceStateRefreshFunc(fmt.Sprint(request["DBInstanceId"]), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudMongodbShardingNetworkPrivateAddressRead(d, meta)
}
func resourceAlicloudMongodbShardingNetworkPrivateAddressRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	MongoDBService := MongoDBService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	object, err := MongoDBService.DescribeMongodbShardingNetworkPrivateAddress(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_mongodb_sharding_network_private_address MongoDBService.DescribeMongodbShardingNetworkPrivateAddress Failed!!! %s", err)
			d.SetId("")
			return nil
		}
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
			"node_id":         item["NodeId"],
		}
		s = append(s, mapping)
	}
	d.Set("network_address", s)
	return nil
}
func resourceAlicloudMongodbShardingNetworkPrivateAddressUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAlicloudMongodbShardingNetworkPrivateAddressRead(d, meta)
}
func resourceAlicloudMongodbShardingNetworkPrivateAddressDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "ReleaseNodePrivateNetworkAddress"
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
	ddsService := MongoDBService{client}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, ddsService.RdsMongodbDBInstanceStateRefreshFunc(fmt.Sprint(parts[0]), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}
	return nil
}
