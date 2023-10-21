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

func resourceAlicloudEciVirtualNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEciVirtualNodeCreate,
		Read:   resourceAlicloudEciVirtualNodeRead,
		Update: resourceAlicloudEciVirtualNodeUpdate,
		Delete: resourceAlicloudEciVirtualNodeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"eip_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enable_public_network": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"kube_config": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"taints": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"effect": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"NoSchedule", "NoExecute", "PreferNoSchedule"}, false),
						},
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"virtual_node_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[\u4e00-\u9fa5a-zA-Z][\u4e00-\u9fa5a-zA-Z0-9-_:]{1,127}$"), "The length of the name is limited to `2` to `128` characters. It can contain uppercase and lowercase letters, Chinese characters, numbers, half-width colon (:), underscores (_), or hyphens (-), and must start with  letters."),
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudEciVirtualNodeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateVirtualNode"
	request := make(map[string]interface{})
	conn, err := client.NewEciClient()
	if err != nil {
		return WrapError(err)
	}
	request["RegionId"] = client.RegionId
	request["VSwitchId"] = d.Get("vswitch_id")
	request["KubeConfig"] = d.Get("kube_config")
	request["SecurityGroupId"] = d.Get("security_group_id")
	if v, ok := d.GetOk("eip_instance_id"); ok {
		request["EipInstanceId"] = v
	}
	if v, ok := d.GetOkExists("enable_public_network"); ok {
		request["EnablePublicNetwork"] = v
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		index := 1
		for key, value := range v.(map[string]interface{}) {
			request[fmt.Sprintf("Tag.%d.Key", index)] = key
			request[fmt.Sprintf("Tag.%d.Value", index)] = value.(string)
			index = index + 1
		}
	}
	if m, ok := d.GetOk("taints"); ok {
		for k, v := range m.(*schema.Set).List() {
			taint := v.(map[string]interface{})
			request[fmt.Sprintf("Taint.%d.Key", k+1)] = taint["key"].(string)
			request[fmt.Sprintf("Taint.%d.Value", k+1)] = taint["value"].(string)
			request[fmt.Sprintf("Taint.%d.Effect", k+1)] = taint["effect"].(string)
		}
	}
	if v, ok := d.GetOk("virtual_node_name"); ok {
		request["VirtualNodeName"] = v
	}
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}
	request["ClientToken"] = buildClientToken("CreateVirtualNode")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-08-08"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_eci_virtual_node", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["VirtualNodeId"]))
	eciService := EciService{client}
	stateConf := BuildStateConf([]string{}, []string{"Ready"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, eciService.EciVirtualNodeStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudEciVirtualNodeRead(d, meta)
}
func resourceAlicloudEciVirtualNodeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eciService := EciService{client}
	object, err := eciService.DescribeEciVirtualNode(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_eci_virtual_node eciService.DescribeEciVirtualNode Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("security_group_id", object["SecurityGroupId"])
	d.Set("status", object["Status"])
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("virtual_node_name", object["VirtualNodeName"])
	d.Set("zone_id", object["ZoneId"])
	d.Set("tags", tagsToMap(object["Tags"]))
	return nil
}
func resourceAlicloudEciVirtualNodeUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAlicloudEciVirtualNodeRead(d, meta)
}
func resourceAlicloudEciVirtualNodeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteVirtualNode"
	var response map[string]interface{}
	conn, err := client.NewEciClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"VirtualNodeId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("DeleteVirtualNode")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-08-08"), StringPointer("AK"), nil, request, &runtime)
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
	return nil
}
