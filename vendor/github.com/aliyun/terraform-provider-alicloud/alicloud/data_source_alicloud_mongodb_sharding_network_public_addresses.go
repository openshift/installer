package alicloud

import (
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudMongodbShardingNetworkPublicAddresses() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudMongodbShardingNetworkPublicAddressesRead,
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"role": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Primary", "Secondary"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
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

func dataSourceAlicloudMongodbShardingNetworkPublicAddressesRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "DescribeShardingNetworkAddress"
	request := make(map[string]interface{})
	request["DBInstanceId"] = d.Get("db_instance_id")
	if v, ok := d.GetOk("node_id"); ok {
		request["NodeId"] = v
	}
	role, roleOk := d.GetOk("role")
	var response map[string]interface{}
	var objects []map[string]interface{}
	conn, err := client.NewDdsClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_mongodb_sharding_network_public_addresses", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.NetworkAddresses.NetworkAddress", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.NetworkAddresses.NetworkAddress", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if roleOk && role.(string) != "" && role.(string) != item["Role"].(string) {
			continue
		}
		if item["NetworkType"].(string) != "Public" {
			continue
		}
		objects = append(objects, item)
	}

	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"db_instance_id":  request["DBInstanceId"],
			"expired_time":    object["ExpiredTime"],
			"ip_address":      object["IPAddress"],
			"network_address": object["NetworkAddress"],
			"network_type":    object["NetworkType"],
			"node_type":       object["NodeType"],
			"port":            object["Port"],
			"role":            object["Role"],
			"vpc_id":          object["VPCId"],
			"vswitch_id":      object["VswitchId"],
		}
		if v, ok := d.GetOk("node_id"); ok {
			mapping["node_id"] = v
		} else {
			mapping["node_id"] = object["NodeId"]
		}
		s = append(s, mapping)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 16))
	if err := d.Set("addresses", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
