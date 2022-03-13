package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCloudStorageGatewayStocks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCloudStorageGatewayStocksRead,
		Schema: map[string]*schema.Schema{
			"gateway_class": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Basic", "Standard", "Enhanced", "Advanced"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"stocks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"available_gateway_classes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCloudStorageGatewayStocksRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DescribeGatewayStock"
	request := make(map[string]interface{})
	request["GatewayRegionId"] = client.RegionId
	gatewayClass, gatewayClassOk := d.GetOk("gateway_class")
	var objects []map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewHcsSgwClient()
	if err != nil {
		return WrapError(err)
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-11"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cloud_storage_gateway_stocks", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	resp, err := jsonpath.Get("$.Stocks.Stock", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Stocks.Stock", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		stockInfo, _ := convertJsonStringToMap(item["StockInfo"].(string))
		if gatewayClassOk && gatewayClass.(string) != "" && !stockInfo[gatewayClass.(string)].(bool) {
			continue
		}

		stockInfoKeys := make([]string, 0)
		for k, v := range stockInfo {
			if v.(bool) {
				stockInfoKeys = append(stockInfoKeys, k)
			}
		}

		if len(stockInfoKeys) > 0 {
			object := map[string]interface{}{
				"zone_id":                   fmt.Sprint(item["ZoneId"]),
				"available_gateway_classes": stockInfoKeys,
			}

			objects = append(objects, object)
		}
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 16))
	if err := d.Set("stocks", objects); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), objects)
	}

	return nil
}
