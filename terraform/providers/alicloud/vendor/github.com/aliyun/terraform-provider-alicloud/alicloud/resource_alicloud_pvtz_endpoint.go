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

func resourceAlicloudPvtzEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPvtzEndpointCreate,
		Read:   resourceAlicloudPvtzEndpointRead,
		Update: resourceAlicloudPvtzEndpointUpdate,
		Delete: resourceAlicloudPvtzEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"endpoint_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ip_configs": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 6,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"cidr_block": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudPvtzEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AddResolverEndpoint"
	request := make(map[string]interface{})
	conn, err := client.NewPvtzClient()
	if err != nil {
		return WrapError(err)
	}
	request["Name"] = d.Get("endpoint_name")
	for k, ipConfig := range d.Get("ip_configs").(*schema.Set).List() {
		ipConfigArg := ipConfig.(map[string]interface{})
		request[fmt.Sprintf("IpConfig.%d.VSwitchId", k+1)] = ipConfigArg["vswitch_id"]
		request[fmt.Sprintf("IpConfig.%d.AzId", k+1)] = ipConfigArg["zone_id"]
		request[fmt.Sprintf("IpConfig.%d.CidrBlock", k+1)] = ipConfigArg["cidr_block"]
		request[fmt.Sprintf("IpConfig.%d.Ip", k+1)] = ipConfigArg["ip"]
	}

	request["Lang"] = "en"
	request["SecurityGroupId"] = d.Get("security_group_id")
	request["VpcId"] = d.Get("vpc_id")
	request["VpcRegionId"] = d.Get("vpc_region_id")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_pvtz_endpoint", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["EndpointId"]))
	pvtzService := PvtzService{client}
	stateConf := BuildStateConf([]string{}, []string{"SUCCESS"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, pvtzService.PvtzEndpointStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudPvtzEndpointRead(d, meta)
}
func resourceAlicloudPvtzEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	pvtzService := PvtzService{client}
	object, err := pvtzService.DescribePvtzEndpoint(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_pvtz_endpoint pvtzService.DescribePvtzEndpoint Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("endpoint_name", object["Name"])
	ipConfigsSli := make([]map[string]interface{}, 0)
	if ipConfigs, ok := object["IpConfigs"].([]interface{}); ok {
		for _, ipConfigArgs := range ipConfigs {
			ipConfigArg := ipConfigArgs.(map[string]interface{})
			ipConfigsMap := make(map[string]interface{})
			ipConfigsMap["zone_id"] = ipConfigArg["AzId"]
			ipConfigsMap["cidr_block"] = ipConfigArg["CidrBlock"]
			ipConfigsMap["ip"] = ipConfigArg["Ip"]
			ipConfigsMap["vswitch_id"] = ipConfigArg["VSwitchId"]
			ipConfigsSli = append(ipConfigsSli, ipConfigsMap)
		}
	}
	d.Set("ip_configs", ipConfigsSli)
	d.Set("security_group_id", object["SecurityGroupId"])
	d.Set("status", object["Status"])
	d.Set("vpc_id", object["VpcId"])
	d.Set("vpc_region_id", object["VpcRegionId"])
	return nil
}
func resourceAlicloudPvtzEndpointUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewPvtzClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	request := map[string]interface{}{
		"EndpointId": d.Id(),
	}

	update := false
	if d.HasChange("endpoint_name") {
		update = true
		if v, ok := d.GetOk("endpoint_name"); ok {
			request["Name"] = v
		}
	}

	if d.HasChange("ip_configs") {
		update = true
		for k, ipConfig := range d.Get("ip_configs").(*schema.Set).List() {
			ipConfigArg := ipConfig.(map[string]interface{})
			request[fmt.Sprintf("IpConfig.%d.VSwitchId", k+1)] = ipConfigArg["vswitch_id"]
			request[fmt.Sprintf("IpConfig.%d.AzId", k+1)] = ipConfigArg["zone_id"]
			request[fmt.Sprintf("IpConfig.%d.CidrBlock", k+1)] = ipConfigArg["cidr_block"]
			request[fmt.Sprintf("IpConfig.%d.Ip", k+1)] = ipConfigArg["ip"]
		}
	}
	request["Lang"] = "en"
	if update {
		action := "UpdateResolverEndpoint"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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

		pvtzService := PvtzService{client}
		stateConf := BuildStateConf([]string{}, []string{"SUCCESS"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, pvtzService.PvtzEndpointStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

	}
	return resourceAlicloudPvtzEndpointRead(d, meta)
}
func resourceAlicloudPvtzEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteResolverEndpoint"
	var response map[string]interface{}
	conn, err := client.NewPvtzClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"EndpointId": d.Id(),
	}

	request["Lang"] = "en"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"ResolverEndpoint.NotExists"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
