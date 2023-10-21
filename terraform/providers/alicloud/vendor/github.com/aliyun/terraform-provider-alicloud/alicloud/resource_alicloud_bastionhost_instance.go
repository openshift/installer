package alicloud

import (
	"fmt"
	"strconv"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/yundun_bastionhost"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudBastionhostInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudBastionhostInstanceCreate,
		Read:   resourceAlicloudBastionhostInstanceRead,
		Update: resourceAlicloudBastionhostInstanceUpdate,
		Delete: resourceAlicloudBastionhostInstanceDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"license_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"period": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntInSlice([]int{1, 3, 6, 12, 24, 36}),
				Optional:     true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": tagsSchema(),

			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_public_access": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudBastionhostInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateInstance"
	request := make(map[string]interface{})
	parameterMapList := make([]map[string]interface{}, 0)
	conn, err := client.NewBssopenapiClient()
	if err != nil {
		return WrapError(err)
	}
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "NetworkType",
		"Value": "vpc",
	})
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "LicenseCode",
		"Value": d.Get("license_code").(string),
	})
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "PlanCode",
		"Value": "cloudbastion",
	})
	request["SubscriptionType"] = "Subscription"
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	request["ProductCode"] = "bastionhost"
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "RegionId",
		"Value": client.RegionId,
	})
	request["Parameter"] = parameterMapList
	request["ClientToken"] = buildClientToken("CreateInstance")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if IsExpectedErrors(err, []string{"NotApplicable"}) {
				conn.Endpoint = String(connectivity.BssOpenAPIEndpointInternational)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_bastionhost_instance", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Code"]) != "Success" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	responseData := response["Data"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseData["InstanceId"]))

	bastionhostService := YundunBastionhostService{client}

	// check RAM policy
	if err := bastionhostService.ProcessRolePolicy(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	// wait for order complete
	stateConf := BuildStateConf([]string{}, []string{"PENDING"}, d.Timeout(schema.TimeoutCreate), 20*time.Second, bastionhostService.BastionhostInstanceRefreshFunc(d.Id(), []string{"UPGRADING", "UPGRADE_FAILED", "CREATE_FAILED"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	rawSecurityGroupIds := d.Get("security_group_ids").([]interface{})
	securityGroupIds := make([]string, len(rawSecurityGroupIds))
	for index, rawSecurityGroupId := range rawSecurityGroupIds {
		securityGroupIds[index] = rawSecurityGroupId.(string)
	}
	// start instance
	if err := bastionhostService.StartBastionhostInstance(d.Id(), d.Get("vswitch_id").(string), securityGroupIds); err != nil {
		return WrapError(err)
	}
	// wait for pending
	stateConf = BuildStateConf([]string{"PENDING", "CREATING"}, []string{"RUNNING"}, d.Timeout(schema.TimeoutCreate), 600*time.Second, bastionhostService.BastionhostInstanceRefreshFunc(d.Id(), []string{"UPGRADING", "UPGRADE_FAILED", "CREATE_FAILED"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudBastionhostInstanceUpdate(d, meta)
}

func resourceAlicloudBastionhostInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	BastionhostService := YundunBastionhostService{client}
	instance, err := BastionhostService.DescribeBastionhostInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("description", instance["Description"])
	d.Set("license_code", instance["LicenseCode"])
	d.Set("vswitch_id", instance["VswitchId"])
	d.Set("security_group_ids", instance["AuthorizedSecurityGroups"])
	d.Set("enable_public_access", instance["PublicNetworkAccess"])
	tags, err := BastionhostService.DescribeTags(d.Id(), nil, TagResourceInstance)
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", BastionhostService.tagsToMap(tags))
	return nil
}

func resourceAlicloudBastionhostInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	bastionhostService := YundunBastionhostService{client}

	d.Partial(true)
	if d.HasChange("tags") {
		if err := bastionhostService.setInstanceTags(d, TagResourceInstance); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}

	if d.HasChange("description") {
		if err := bastionhostService.UpdateBastionhostInstanceDescription(d.Id(), d.Get("description").(string)); err != nil {
			return WrapError(err)
		}
		d.SetPartial("description")
	}

	if d.HasChange("resource_group_id") {
		if err := bastionhostService.UpdateResourceGroup(d.Id(), d.Get("resource_group_id").(string)); err != nil {
			return WrapError(err)
		}
		d.SetPartial("resource_group_id")
	}

	if !d.IsNewResource() && d.HasChange("license_code") {
		params := map[string]string{
			"LicenseCode": "license_code",
		}
		if err := bastionhostService.UpdateInstanceSpec(params, d, meta); err != nil {
			return WrapError(err)
		}
		stateConf := BuildStateConf([]string{"UPGRADING"}, []string{"PENDING", "RUNNING"}, d.Timeout(schema.TimeoutUpdate), 20*time.Second, bastionhostService.BastionhostInstanceRefreshFunc(d.Id(), []string{"CREATING", "UPGRADE_FAILED", "CREATE_FAILED"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("license_code")
	}

	if !d.IsNewResource() && d.HasChange("security_group_ids") {
		securityGroupIds := d.Get("security_group_ids").([]interface{})
		sgs := make([]string, 0, len(securityGroupIds))
		for _, rawSecurityGroupId := range securityGroupIds {
			sgs = append(sgs, rawSecurityGroupId.(string))
		}
		if err := bastionhostService.UpdateBastionhostSecurityGroups(d.Id(), sgs); err != nil {
			return WrapError(err)
		}
		stateConf := BuildStateConf([]string{"UPGRADING"}, []string{"RUNNING"}, d.Timeout(schema.TimeoutUpdate), 20*time.Second, bastionhostService.BastionhostInstanceRefreshFunc(d.Id(), []string{"CREATING", "UPGRADE_FAILED", "CREATE_FAILED"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("security_group_ids")
	}

	if d.HasChange("enable_public_access") {
		client := meta.(*connectivity.AliyunClient)
		BastionhostService := YundunBastionhostService{client}
		instance, err := BastionhostService.DescribeBastionhostInstance(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := strconv.FormatBool(d.Get("enable_public_access").(bool))
		if strconv.FormatBool(instance["PublicNetworkAccess"].(bool)) != target {
			if target == "false" {
				err := BastionhostService.DisableInstancePublicAccess(d.Id())
				if err != nil {
					return WrapError(err)
				}
			} else {
				err := BastionhostService.EnableInstancePublicAccess(d.Id())
				if err != nil {
					return WrapError(err)
				}
			}
		}
		d.SetPartial("enable_public_access")
	}

	d.Partial(false)
	// wait for order complete
	return resourceAlicloudBastionhostInstanceRead(d, meta)
}

func resourceAlicloudBastionhostInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	bastionhostService := YundunBastionhostService{client}
	request := yundun_bastionhost.CreateRefundInstanceRequest()
	request.InstanceId = d.Id()

	raw, err := bastionhostService.client.WithBastionhostClient(func(BastionhostClient *yundun_bastionhost.Client) (interface{}, error) {
		return BastionhostClient.RefundInstance(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	// Wait for the release procedure of cloud resource dependencies. Instance can not be fetched through api as soon as release has
	// been invoked, however the resources have not been fully destroyed yet. Therefore, a certain amount time of waiting
	// is quite necessary (conservative estimation cloud be less then 3 minutes)
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, bastionhostService.BastionhostInstanceRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
