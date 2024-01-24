package alicloud

import (
	"errors"
	"strconv"
	"time"

	"github.com/denverdino/aliyungo/common"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudAlikafkaInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAlikafkaInstanceCreate,
		Read:   resourceAlicloudAlikafkaInstanceRead,
		Update: resourceAlicloudAlikafkaInstanceUpdate,
		Delete: resourceAlicloudAlikafkaInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(3, 64),
			},
			"topic_quota": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"disk_type": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"disk_size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"deploy_type": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"io_max": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"paid_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
				Default:      PostPaid,
			},
			"spec_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "normal",
			},
			"eip_max": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"security_group": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"config": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.ValidateJsonString,
				DiffSuppressFunc: alikafkaInstanceConfigDiffSuppressFunc,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_point": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAlicloudAlikafkaInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}
	vpcService := VpcService{client}

	regionId := client.RegionId
	topicQuota := d.Get("topic_quota").(int)
	diskType := d.Get("disk_type").(int)
	diskSize := d.Get("disk_size").(int)
	deployType := d.Get("deploy_type").(int)
	ioMax := d.Get("io_max").(int)
	vswitchId := d.Get("vswitch_id").(string)
	paidType := d.Get("paid_type").(string)
	specType := d.Get("spec_type").(string)

	// Get vswitch info by vswitchId
	vsw, err := vpcService.DescribeVSwitch(vswitchId)
	if err != nil {
		return WrapError(err)
	}

	// 1. Create order
	createOrderReq := alikafka.CreateCreatePostPayOrderRequest()
	createOrderReq.RegionId = regionId
	createOrderReq.TopicQuota = requests.NewInteger(topicQuota)
	createOrderReq.DiskType = strconv.Itoa(diskType)
	createOrderReq.DiskSize = requests.NewInteger(diskSize)
	createOrderReq.DeployType = requests.NewInteger(deployType)
	createOrderReq.IoMax = requests.NewInteger(ioMax)
	createOrderReq.PaidType = requests.NewInteger(1)
	createOrderReq.SpecType = specType
	if paidType == string(PrePaid) {
		createOrderReq.PaidType = requests.NewInteger(0)
	}
	if v, ok := d.GetOk("eip_max"); ok {
		createOrderReq.EipMax = requests.NewInteger(v.(int))
	}

	var createOrderResp *alikafka.CreatePostPayOrderResponse
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.CreatePostPayOrder(createOrderReq)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL", "ONS_SYSTEM_ERROR"}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(createOrderReq.GetActionName(), raw, createOrderReq.RpcRequest, createOrderReq)
		v, _ := raw.(*alikafka.CreatePostPayOrderResponse)
		createOrderResp = v
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alikafka_instance", createOrderReq.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	alikafkaInstanceVO, err := alikafkaService.DescribeAlikafkaInstanceByOrderId(createOrderResp.OrderId, 60)

	if err != nil {
		return WrapError(err)
	}

	instanceId := alikafkaInstanceVO.InstanceId
	d.SetId(instanceId)

	// 3. Start instance
	startInstanceReq := alikafka.CreateStartInstanceRequest()
	startInstanceReq.RegionId = regionId
	startInstanceReq.InstanceId = instanceId
	startInstanceReq.VpcId = vsw.VpcId
	startInstanceReq.VSwitchId = vswitchId
	startInstanceReq.ZoneId = vsw.ZoneId
	if _, ok := d.GetOk("eip_max"); ok {
		startInstanceReq.IsEipInner = requests.NewBoolean(true)
		startInstanceReq.DeployModule = "eip"
	}
	if v, ok := d.GetOk("name"); ok {
		startInstanceReq.Name = v.(string)
	}
	if v, ok := d.GetOk("security_group"); ok {
		startInstanceReq.SecurityGroup = v.(string)
	}
	if v, ok := d.GetOk("service_version"); ok {
		startInstanceReq.ServiceVersion = v.(string)
	}
	if v, ok := d.GetOk("config"); ok {
		startInstanceReq.Config = v.(string)
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.StartInstance(startInstanceReq)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(startInstanceReq.GetActionName(), raw, startInstanceReq.RpcRequest, startInstanceReq)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alikafka_instance", startInstanceReq.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	// 3. wait until running
	err = alikafkaService.WaitForAlikafkaInstance(d.Id(), Running, DefaultLongTimeout)

	if err != nil {
		return WrapError(err)
	}

	return resourceAlicloudAlikafkaInstanceUpdate(d, meta)
}

func resourceAlicloudAlikafkaInstanceRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}

	object, err := alikafkaService.DescribeAlikafkaInstance(d.Id())
	if err != nil {
		// Handle exceptions
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object.Name)
	d.Set("topic_quota", object.TopicNumLimit)
	d.Set("disk_type", object.DiskType)
	d.Set("disk_size", object.DiskSize)
	d.Set("deploy_type", object.DeployType)
	d.Set("io_max", object.IoMax)
	d.Set("eip_max", object.EipMax)
	d.Set("vpc_id", object.VpcId)
	d.Set("vswitch_id", object.VSwitchId)
	d.Set("zone_id", object.ZoneId)
	d.Set("paid_type", PostPaid)
	d.Set("spec_type", object.SpecType)
	d.Set("security_group", object.SecurityGroup)
	d.Set("end_point", object.EndPoint)
	// object.UpgradeServiceDetailInfo.UpgradeServiceDetailInfoVO[0].Current2OpenSourceVersion can guaranteed not to be null
	d.Set("service_version", object.UpgradeServiceDetailInfo.Current2OpenSourceVersion)
	d.Set("config", object.AllConfig)
	if object.PaidType == 0 {
		d.Set("paid_type", PrePaid)
	}

	tags, err := alikafkaService.DescribeTags(d.Id(), nil, TagResourceInstance)
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", alikafkaService.tagsToMap(tags))

	return nil
}

func resourceAlicloudAlikafkaInstanceUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}
	d.Partial(true)
	if err := alikafkaService.setInstanceTags(d, TagResourceInstance); err != nil {
		return WrapError(err)
	}
	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlicloudAlikafkaInstanceRead(d, meta)
	}
	// Process change instance name.
	if d.HasChange("name") {
		var name string
		if v, ok := d.GetOk("name"); ok {
			name = v.(string)
		}
		modifyInstanceNameReq := alikafka.CreateModifyInstanceNameRequest()
		modifyInstanceNameReq.RegionId = client.RegionId
		modifyInstanceNameReq.InstanceId = d.Id()
		modifyInstanceNameReq.InstanceName = name

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
				return alikafkaClient.ModifyInstanceName(modifyInstanceNameReq)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
					time.Sleep(10 * time.Second)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(modifyInstanceNameReq.GetActionName(), raw, modifyInstanceNameReq.RpcRequest, modifyInstanceNameReq)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifyInstanceNameReq.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("name")
	}

	// Process paid type change, note only support change from post to pre pay.
	if d.HasChange("paid_type") {
		o, n := d.GetChange("paid_type")
		oldPaidType := o.(string)
		newPaidType := n.(string)
		oldPaidTypeInt := 1
		newPaidTypeInt := 1
		if oldPaidType == string(PrePaid) {
			oldPaidTypeInt = 0
		}
		if newPaidType == string(PrePaid) {
			newPaidTypeInt = 0
		}
		if oldPaidTypeInt == 1 && newPaidTypeInt == 0 {

			convertPostPayOrderReq := alikafka.CreateConvertPostPayOrderRequest()
			convertPostPayOrderReq.InstanceId = d.Id()
			convertPostPayOrderReq.RegionId = client.RegionId
			err := resource.Retry(5*time.Minute, func() *resource.RetryError {
				raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
					return alikafkaClient.ConvertPostPayOrder(convertPostPayOrderReq)
				})
				if err != nil {
					if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
						time.Sleep(10 * time.Second)
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(convertPostPayOrderReq.GetActionName(), raw, convertPostPayOrderReq.RpcRequest, convertPostPayOrderReq)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), convertPostPayOrderReq.GetActionName(), AlibabaCloudSdkGoERROR)
			}

			// Make sure convert success
			object, err := alikafkaService.DescribeAlikafkaInstance(d.Id())
			if err != nil {
				return WrapError(err)
			}

			err = alikafkaService.WaitForAlikafkaInstanceUpdated(d.Id(), object.TopicNumLimit,
				object.DiskSize, object.IoMax, object.EipMax, newPaidTypeInt, object.SpecType, DefaultTimeoutMedium)
			if err != nil {
				return WrapError(err)
			}
		} else {
			return WrapError(errors.New("paid type only support change from post pay to pre pay"))
		}

		d.SetPartial("paid_type")
	}

	attributeUpdate := false

	upgradeReq := alikafka.CreateUpgradePostPayOrderRequest()
	upgradeReq.RegionId = client.RegionId
	upgradeReq.InstanceId = d.Id()
	upgradeReq.TopicQuota = requests.NewInteger(d.Get("topic_quota").(int))
	upgradeReq.DiskSize = requests.NewInteger(d.Get("disk_size").(int))
	upgradeReq.IoMax = requests.NewInteger(d.Get("io_max").(int))
	upgradeReq.SpecType = d.Get("spec_type").(string)

	if d.HasChange("topic_quota") || d.HasChange("disk_size") || d.HasChange("io_max") || d.HasChange("spec_type") {
		attributeUpdate = true
	}
	eipMax := 0
	if v, ok := d.GetOk("eip_max"); ok {
		eipMax = v.(int)
	}
	if d.HasChange("eip_max") {
		if v, ok := d.GetOk("eip_max"); ok {
			eipMax = v.(int)
		}
		upgradeReq.EipMax = requests.NewInteger(eipMax)
		attributeUpdate = true
	}

	if attributeUpdate {
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
				return alikafkaClient.UpgradePostPayOrder(upgradeReq)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
					time.Sleep(10 * time.Second)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(upgradeReq.GetActionName(), raw, upgradeReq.RpcRequest, upgradeReq)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), upgradeReq.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("topic_quota")
		d.SetPartial("disk_size")
		d.SetPartial("io_max")
		d.SetPartial("spec_type")
		d.SetPartial("eip_max")
	}

	paidType := 1
	if d.Get("paid_type").(string) == string(PrePaid) {
		paidType = 0
	}
	err := alikafkaService.WaitForAlikafkaInstanceUpdated(d.Id(), d.Get("topic_quota").(int), d.Get("disk_size").(int),
		d.Get("io_max").(int), eipMax, paidType, d.Get("spec_type").(string), DefaultTimeoutMedium)

	if err != nil {
		return WrapError(err)
	}

	err = alikafkaService.WaitForAlikafkaInstance(d.Id(), Running, 6000)

	if err != nil {
		return WrapError(err)
	}

	if d.HasChange("service_version") {
		var serviceVersion string
		if v, ok := d.GetOk("service_version"); ok {
			serviceVersion = v.(string)
		}
		upgradeInstanceVersionReq := alikafka.CreateUpgradeInstanceVersionRequest()
		upgradeInstanceVersionReq.RegionId = client.RegionId
		upgradeInstanceVersionReq.InstanceId = d.Id()
		upgradeInstanceVersionReq.TargetVersion = serviceVersion

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
				return alikafkaClient.UpgradeInstanceVersion(upgradeInstanceVersionReq)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
					time.Sleep(10 * time.Second)
					return resource.RetryableError(err)
				}
				// means no need to update version
				if IsExpectedErrors(err, []string{"ONS_INIT_ENV_ERROR"}) {
					return nil
				}
				return resource.NonRetryableError(err)
			}
			addDebug(upgradeInstanceVersionReq.GetActionName(), raw, upgradeInstanceVersionReq.RpcRequest, upgradeInstanceVersionReq)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), upgradeInstanceVersionReq.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		// wait for upgrade task be invoke
		time.Sleep(60 * time.Second)
		// upgrade service may be last a long time
		err = alikafkaService.WaitForAlikafkaInstance(d.Id(), Running, 10000)
		if err != nil {
			return WrapError(err)
		}
		d.SetPartial("service_version")
	}

	if d.HasChange("config") {
		var config string
		if v, ok := d.GetOk("config"); ok {
			config = v.(string)
		}
		upgradeInstanceConfigReq := alikafka.CreateUpdateInstanceConfigRequest()
		upgradeInstanceConfigReq.RegionId = client.RegionId
		upgradeInstanceConfigReq.InstanceId = d.Id()
		upgradeInstanceConfigReq.Config = config

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
				return alikafkaClient.UpdateInstanceConfig(upgradeInstanceConfigReq)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
					time.Sleep(10 * time.Second)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(upgradeInstanceConfigReq.GetActionName(), raw, upgradeInstanceConfigReq.RpcRequest, upgradeInstanceConfigReq)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), upgradeInstanceConfigReq.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		// wait for upgrade task be invoke
		time.Sleep(60 * time.Second)
		err = alikafkaService.WaitForAlikafkaInstance(d.Id(), Running, 6000)
		if err != nil {
			return WrapError(err)
		}
		d.SetPartial("config")
	}

	d.Partial(false)
	return resourceAlicloudAlikafkaInstanceRead(d, meta)
}

func resourceAlicloudAlikafkaInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}

	// Pre paid instance can not be release.
	if d.Get("paid_type").(string) == string(PrePaid) {
		return nil
	}

	request := alikafka.CreateReleaseInstanceRequest()
	request.InstanceId = d.Id()
	request.RegionId = client.RegionId
	request.ReleaseIgnoreTime = requests.NewBoolean(true)
	request.ForceDeleteInstance = requests.NewBoolean(true)

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.ReleaseInstance(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return WrapError(alikafkaService.WaitForAllAlikafkaNodeRelease(d.Id(), "released", DefaultTimeoutMedium))
}
