package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudMseGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMseGatewayCreate,
		Read:   resourceAlicloudMseGatewayRead,
		Update: resourceAlicloudMseGatewayUpdate,
		Delete: resourceAlicloudMseGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"backup_vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enterprise_security_group": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"gateway_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"internet_slb_spec": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"replica": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 30),
			},
			"slb_spec": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"spec": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"MSE_GTW_16_32_200_c", "MSE_GTW_2_4_200_c", "MSE_GTW_4_8_200_c", "MSE_GTW_8_16_200_c"}, false),
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
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"delete_slb": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"slb_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"associate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slb_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slb_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slb_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gmt_create": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_slb_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_slb_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudMseGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AddGateway"
	request := make(map[string]interface{})
	conn, err := client.NewMseClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("backup_vswitch_id"); ok {
		request["VSwitchId2"] = v
	}
	if v, ok := d.GetOkExists("enterprise_security_group"); ok {
		request["EnterpriseSecurityGroup"] = v
	}
	if v, ok := d.GetOk("gateway_name"); ok {
		request["Name"] = v
	}
	request["Region"] = client.RegionId
	request["Replica"] = d.Get("replica")
	if v, ok := d.GetOk("slb_spec"); ok {
		request["SlbSpec"] = v
	}
	if v, ok := d.GetOk("internet_slb_spec"); ok {
		request["InternetSlbSpec"] = v
	}
	request["Spec"] = d.Get("spec")

	request["VSwitchId"] = d.Get("vswitch_id")
	request["Vpc"] = d.Get("vpc_id")

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-05-31"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_mse_gateway", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Code"]) != "200" || fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	responseData := response["Data"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseData["GatewayUniqueId"]))
	mseService := MseService{client}
	stateConf := BuildStateConf([]string{}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, mseService.MseGatewayStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudMseGatewayRead(d, meta)
}
func resourceAlicloudMseGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mseService := MseService{client}
	object, err := mseService.DescribeMseGateway(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_mse_gateway mseService.DescribeMseGateway Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("backup_vswitch_id", object["Vswitch2"])
	d.Set("gateway_name", object["Name"])
	d.Set("replica", formatInt(object["Replica"]))
	d.Set("spec", object["Spec"])
	d.Set("status", fmt.Sprint(formatInt(object["Status"])))
	d.Set("vswitch_id", object["Vswitch"])
	d.Set("vpc_id", object["Vpc"])

	slbListObject, err := mseService.ListGatewaySlb(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_mse_gateway mseService.ListGatewaySlb Failed!!! %s", err)
			return nil
		}
		return WrapError(err)
	}
	slbList := slbListObject["SlbList"]

	slbMapList := make([]map[string]interface{}, 0)
	if v, ok := slbList.([]interface{}); ok && len(v) > 0 {
		for _, slb := range v {
			slbMap := slb.(map[string]interface{})
			mapping := map[string]interface{}{
				"associate_id":       slbMap["Id"],
				"slb_id":             slbMap["SlbId"],
				"slb_ip":             slbMap["SlbIp"],
				"slb_port":           slbMap["SlbPort"],
				"type":               slbMap["Type"],
				"gmt_create":         slbMap["GmtCreate"],
				"gateway_slb_mode":   slbMap["GatewaySlbMode"],
				"gateway_slb_status": slbMap["GatewaySlbStatus"],
			}
			slbMapList = append(slbMapList, mapping)
		}
	}
	d.Set("slb_list", slbMapList)

	return nil
}
func resourceAlicloudMseGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := map[string]interface{}{
		"GatewayUniqueId": d.Id(),
	}
	var response map[string]interface{}
	if d.HasChange("gateway_name") {
		if v, ok := d.GetOk("gateway_name"); ok {
			request["Name"] = v
		}
	}
	action := "UpdateGatewayName"
	conn, err := client.NewMseClient()
	if err != nil {
		return WrapError(err)
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2019-05-31"), StringPointer("AK"), request, nil, &util.RuntimeOptions{})
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
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	return resourceAlicloudMseGatewayRead(d, meta)
}
func resourceAlicloudMseGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteGateway"
	var response map[string]interface{}
	conn, err := client.NewMseClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"GatewayUniqueId": d.Id(),
	}

	if v, ok := d.GetOkExists("delete_slb"); ok {
		request["DeleteSlb"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2019-05-31"), StringPointer("AK"), request, nil, &util.RuntimeOptions{})
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
	if IsExpectedErrorCodes(fmt.Sprint(response["Code"]), []string{"404"}) {
		return nil
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	mseService := MseService{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 10*time.Second, mseService.MseGatewayStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
