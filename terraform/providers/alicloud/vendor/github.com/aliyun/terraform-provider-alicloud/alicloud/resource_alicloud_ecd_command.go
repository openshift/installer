package alicloud

import (
	"encoding/base64"
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudEcdCommand() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcdCommandCreate,
		Read:   resourceAlicloudEcdCommandRead,
		Update: resourceAlicloudEcdCommandUpdate,
		Delete: resourceAlicloudEcdCommandDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"command_content": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"command_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"RunBatScript", "RunPowerShellScript"}, false),
			},
			"content_encoding": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Base64", "PlainText"}, false),
			},
			"desktop_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"timeout": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudEcdCommandCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecdService := EcdService{client}
	var response map[string]interface{}
	action := "RunCommand"
	request := make(map[string]interface{})
	conn, err := client.NewGwsecdClient()
	if err != nil {
		return WrapError(err)
	}
	request["CommandContent"] = d.Get("command_content")
	request["Type"] = d.Get("command_type")
	if v, ok := d.GetOk("content_encoding"); ok {
		request["ContentEncoding"] = v
	}
	request["DesktopId.1"] = d.Get("desktop_id")
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("timeout"); ok {
		request["Timeout"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecd_command", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["InvokeId"]))
	stateConf := BuildStateConf([]string{}, []string{"Success"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, ecdService.EcdCommandStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudEcdCommandRead(d, meta)
}
func resourceAlicloudEcdCommandRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecdService := EcdService{client}
	object, err := ecdService.DescribeEcdCommand(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecd_command ecdService.DescribeEcdCommand Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	commandContent, _ := base64.StdEncoding.DecodeString(object["CommandContent"].(string))
	d.Set("command_content", string(commandContent))
	d.Set("command_type", object["CommandType"])
	d.Set("status", object["InvocationStatus"])
	temp1 := map[string]interface{}{}
	if invokeDesktopsList, ok := object["InvokeDesktops"].([]interface{}); ok {
		for _, v := range invokeDesktopsList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 = map[string]interface{}{
					"desktop_id": m1["DesktopId"],
				}
			}
		}
	}
	d.Set("desktop_id", temp1["desktop_id"])
	return nil
}
func resourceAlicloudEcdCommandUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return nil
}
func resourceAlicloudEcdCommandDelete(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not delete operation."))
	return nil

}
