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

func resourceAlicloudEcpKeyPair() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcpKeyPairCreate,
		Read:   resourceAlicloudEcpKeyPairRead,
		Update: resourceAlicloudEcpKeyPairUpdate,
		Delete: resourceAlicloudEcpKeyPairDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"key_pair_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"public_key_body": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAlicloudEcpKeyPairCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "ImportKeyPair"
	request := make(map[string]interface{})
	conn, err := client.NewCloudphoneClient()
	if err != nil {
		return WrapError(err)
	}
	request["KeyPairName"] = d.Get("key_pair_name")
	request["PublicKeyBody"] = d.Get("public_key_body")
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-12-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecp_key_pair", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["KeyPairName"]))

	return resourceAlicloudEcpKeyPairRead(d, meta)
}
func resourceAlicloudEcpKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudphoneService := CloudphoneService{client}
	_, err := cloudphoneService.DescribeEcpKeyPair(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecp_key_pair cloudphoneService.DescribeEcpKeyPair Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("key_pair_name", d.Id())
	return nil
}
func resourceAlicloudEcpKeyPairUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAlicloudEcpKeyPairRead(d, meta)
}
func resourceAlicloudEcpKeyPairDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteKeyPairs"
	var response map[string]interface{}
	conn, err := client.NewCloudphoneClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"KeyPairName": []string{d.Id()},
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-12-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"KeyPair.WithInstance"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"KeyPairsNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
