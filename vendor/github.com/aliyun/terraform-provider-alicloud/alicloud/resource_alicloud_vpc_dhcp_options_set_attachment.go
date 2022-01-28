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

func resourceAlicloudVpcDhcpOptionsSetAttachement() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudVpcDhcpOptionsAttachmentCreate,
		Read:   resourceAlicloudVpcDhcpOptionsSetAttachmentRead,
		Update: resourceAlicloudVpcDhcpOptionsSetAttachmentUpdate,
		Delete: resourceAlicloudVpcDhcpOptionsSetAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"dhcp_options_set_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudVpcDhcpOptionsAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AttachDhcpOptionsSetToVpc"
	request := make(map[string]interface{})
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}

	request["RegionId"] = client.RegionId
	request["DhcpOptionsSetId"] = d.Get("dhcp_options_set_id")
	request["VpcId"] = d.Get("vpc_id")
	request["ClientToken"] = buildClientToken(action)
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_dhcp_options_set_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["VpcId"], ":", request["DhcpOptionsSetId"]))

	vpcService := VpcService{client}
	stateConf := BuildStateConf([]string{"Pending"}, []string{"InUse"}, d.Timeout(schema.TimeoutCreate), 3*time.Second, vpcService.DescribeVpcDhcpOptionsSetAttachmentStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudVpcDhcpOptionsSetAttachmentRead(d, meta)
}
func resourceAlicloudVpcDhcpOptionsSetAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	object, err := vpcService.DescribeVpcDhcpOptionsSetAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_dhcp_options_set_attachment vpcService.DescribeVpcDhcpOptionsSetAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if associateVpcsList, ok := object["AssociateVpcs"]; ok {
		for _, associateVpcsListItem := range associateVpcsList.([]interface{}) {
			if associateVpcsListItem != nil {
				associateVpcsListItemMap, ok := associateVpcsListItem.(map[string]interface{})
				if ok && associateVpcsListItemMap["VpcId"] == parts[0] {
					d.Set("vpc_id", associateVpcsListItemMap["VpcId"])
					d.Set("dhcp_options_set_id", parts[1])
					d.Set("status", associateVpcsListItemMap["AssociateStatus"])
					break
				}
			}
		}
	}
	return nil
}
func resourceAlicloudVpcDhcpOptionsSetAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAlicloudVpcDhcpOptionsSetAttachmentRead(d, meta)
}
func resourceAlicloudVpcDhcpOptionsSetAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DetachDhcpOptionsSetFromVpc"
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DhcpOptionsSetId": parts[1],
	}
	request["RegionId"] = client.RegionId
	request["VpcId"] = d.Get("vpc_id")
	request["ClientToken"] = buildClientToken(action)
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 10*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectStatus.DhcpOptionsSet"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDhcpOptionsSetId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	vpcService := VpcService{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcService.VpcDhcpOptionsSetStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
