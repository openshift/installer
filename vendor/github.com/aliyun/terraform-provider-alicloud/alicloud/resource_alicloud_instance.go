package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/denverdino/aliyungo/common"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"encoding/base64"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliyunInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunInstanceCreate,
		Read:   resourceAliyunInstanceRead,
		Update: resourceAliyunInstanceUpdate,
		Delete: resourceAliyunInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"image_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^ecs\..*`), "prefix must be 'ecs.'"),
			},

			"credit_specification": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(CreditSpecificationStandard),
					string(CreditSpecificationUnlimited),
				}, false),
			},

			"security_groups": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},

			"allocate_public_ip": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "Field 'allocate_public_ip' has been deprecated from provider version 1.6.1. Setting 'internet_max_bandwidth_out' larger than 0 will allocate public ip for instance.",
			},

			"instance_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ECS-Instance",
				ValidateFunc: validation.StringLenBetween(2, 128),
			},

			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},

			"internet_charge_type": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringInSlice([]string{"PayByBandwidth", "PayByTraffic"}, false),
				Default:          PayByTraffic,
				DiffSuppressFunc: ecsInternetDiffSuppressFunc,
			},
			"internet_max_bandwidth_in": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: ecsInternetDiffSuppressFunc,
				Deprecated:       "The attribute is invalid and no any affect for the instance. So it has been deprecated from version v1.121.2.",
			},
			"internet_max_bandwidth_out": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"kms_encrypted_password": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: kmsDiffSuppressFunc,
			},
			"kms_encryption_context": {
				Type:     schema.TypeMap,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("kms_encrypted_password") == ""
				},
				Elem: schema.TypeString,
			},
			"io_optimized": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Attribute io_optimized has been deprecated on instance resource. All the launched alicloud instances will be IO optimized. Suggest to remove it from your template.",
			},
			"is_outdated": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"system_disk_category": {
				Type:         schema.TypeString,
				Default:      DiskCloudEfficiency,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"all", "cloud", "ephemeral_ssd", "cloud_essd", "cloud_efficiency", "cloud_ssd", "local_disk"}, false),
			},
			"system_disk_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"system_disk_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"system_disk_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  40,
			},
			"system_disk_performance_level": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: ecsSystemDiskPerformanceLevelSuppressFunc,
				ValidateFunc:     validation.StringInSlice([]string{"PL0", "PL1", "PL2", "PL3"}, false),
			},
			"system_disk_auto_snapshot_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"data_disks": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				MaxItems: 16,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringLenBetween(2, 128),
						},
						"size": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"category": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"all", "cloud", "ephemeral_ssd", "cloud_essd", "cloud_efficiency", "cloud_ssd", "local_disk"}, false),
							Default:      DiskCloudEfficiency,
							ForceNew:     true,
						},
						"encrypted": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
							ForceNew: true,
						},
						"kms_key_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"auto_snapshot_policy_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"delete_with_instance": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
							Default:  true,
						},
						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringLenBetween(2, 256),
						},
						"performance_level": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"PL0", "PL1", "PL2", "PL3"}, false),
						},
					},
				},
			},

			//subnet_id and vswitch_id both exists, cause compatible old version, and aws habit.
			"subnet_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true, //add this schema cause subnet_id not used enter parameter, will different, so will be ForceNew
				ConflictsWith: []string{"vswitch_id"},
			},

			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"private_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
				Default:      PostPaid,
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateFunc: validation.Any(
					validation.IntBetween(1, 9),
					validation.IntInSlice([]int{12, 24, 36, 48, 60})),
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"period_unit": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          Month,
				ValidateFunc:     validation.StringInSlice([]string{"Week", "Month"}, false),
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"renewal_status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  RenewNormal,
				ValidateFunc: validation.StringInSlice([]string{
					string(RenewAutoRenewal),
					string(RenewNormal),
					string(RenewNotRenewal)}, false),
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"auto_renew_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 6, 12}),
				DiffSuppressFunc: ecsNotAutoRenewDiffSuppressFunc,
			},
			"include_data_disks": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          true,
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Running", "Stopped"}, false),
				Default:      "Running",
			},

			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"role_name": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				DiffSuppressFunc: vpcTypeResourceDiffSuppressFunc,
			},

			"key_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"spot_strategy": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          NoSpot,
				ValidateFunc:     validation.StringInSlice([]string{"NoSpot", "SpotAsPriceGo", "SpotWithPriceLimit"}, false),
				DiffSuppressFunc: ecsSpotStrategyDiffSuppressFunc,
			},

			"spot_price_limit": {
				Type:             schema.TypeFloat,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: ecsSpotPriceLimitDiffSuppressFunc,
			},

			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"force_delete": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          false,
				Description:      descriptions["A behavior mark used to delete 'PrePaid' ECS instance forcibly."],
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},

			"security_enhancement_strategy": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(ActiveSecurityEnhancementStrategy),
					string(DeactiveSecurityEnhancementStrategy),
				}, false),
			},

			"tags":        tagsSchemaWithIgnore(),
			"volume_tags": tagsSchemaComputed(),

			"auto_release_time": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					diff := d.Get("instance_charge_type").(string) == "PrePaid"
					if diff {
						return diff
					}
					if old != "" && new != "" && strings.HasPrefix(new, strings.Trim(old, "Z")) {
						diff = true
					}
					return diff
				},
			},
		},
	}
}

func resourceAliyunInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	// Ensure instance_type is valid
	//zoneId, validZones, requestId, err := ecsService.DescribeAvailableResources(d, meta, InstanceTypeResource)
	//if err != nil {
	//	return WrapError(err)
	//}
	//if err := ecsService.InstanceTypeValidation(d.Get("instance_type").(string), zoneId, validZones); err != nil {
	//	return WrapError(Error("%s. RequestId: %s", err, requestId))
	//}

	request, err := buildAliyunInstanceArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}
	request.IoOptimized = "optimized"
	if d.Get("is_outdated").(bool) == true {
		request.IoOptimized = "none"
	}
	wait := incrementalWait(1*time.Second, 1*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.RunInstances(request)
		})
		if err != nil {
			if IsThrottling(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ecs.RunInstancesResponse)
		d.SetId(response.InstanceIdSets.InstanceIdSet[0])
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_instance", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"Pending", "Starting", "Stopped"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, ecsService.InstanceStateRefreshFunc(d.Id(), []string{"Stopping"}))

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliyunInstanceUpdate(d, meta)
}

func resourceAliyunInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	instance, err := ecsService.DescribeInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_instance ecsService.DescribeInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	var disk ecs.Disk
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		disk, err = ecsService.DescribeInstanceSystemDisk(d.Id(), instance.ResourceGroupId)
		if err != nil {
			if NotFoundError(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapError(err)
	}
	d.Set("system_disk_category", disk.Category)
	d.Set("system_disk_size", disk.Size)
	d.Set("system_disk_auto_snapshot_policy_id", disk.AutoSnapshotPolicyId)
	d.Set("volume_tags", ecsService.tagsToMap(disk.Tags.Tag))
	d.Set("system_disk_performance_level", disk.PerformanceLevel)
	d.Set("instance_name", instance.InstanceName)
	d.Set("resource_group_id", instance.ResourceGroupId)
	d.Set("description", instance.Description)
	d.Set("status", instance.Status)
	d.Set("availability_zone", instance.ZoneId)
	d.Set("host_name", instance.HostName)
	d.Set("image_id", instance.ImageId)
	d.Set("instance_type", instance.InstanceType)
	d.Set("password", d.Get("password").(string))
	d.Set("internet_max_bandwidth_out", instance.InternetMaxBandwidthOut)
	d.Set("internet_max_bandwidth_in", instance.InternetMaxBandwidthIn)
	d.Set("instance_charge_type", instance.InstanceChargeType)
	d.Set("key_name", instance.KeyPairName)
	d.Set("spot_strategy", instance.SpotStrategy)
	d.Set("spot_price_limit", instance.SpotPriceLimit)
	d.Set("internet_charge_type", instance.InternetChargeType)
	d.Set("deletion_protection", instance.DeletionProtection)
	d.Set("credit_specification", instance.CreditSpecification)
	d.Set("auto_release_time", instance.AutoReleaseTime)
	d.Set("tags", ecsService.tagsToMap(instance.Tags.Tag))

	if len(instance.PublicIpAddress.IpAddress) > 0 {
		d.Set("public_ip", instance.PublicIpAddress.IpAddress[0])
	} else {
		d.Set("public_ip", "")
	}
	d.Set("subnet_id", instance.VpcAttributes.VSwitchId)
	d.Set("vswitch_id", instance.VpcAttributes.VSwitchId)

	if len(instance.VpcAttributes.PrivateIpAddress.IpAddress) > 0 {
		d.Set("private_ip", instance.VpcAttributes.PrivateIpAddress.IpAddress[0])
	} else {
		d.Set("private_ip", strings.Join(instance.InnerIpAddress.IpAddress, ","))
	}

	sgs := make([]string, 0, len(instance.SecurityGroupIds.SecurityGroupId))
	for _, sg := range instance.SecurityGroupIds.SecurityGroupId {
		sgs = append(sgs, sg)
	}
	if err := d.Set("security_groups", sgs); err != nil {
		return WrapError(err)
	}

	if !d.IsNewResource() || d.HasChange("user_data") {
		dataRequest := ecs.CreateDescribeUserDataRequest()
		dataRequest.RegionId = client.RegionId
		dataRequest.InstanceId = d.Id()
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeUserData(dataRequest)
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), dataRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(dataRequest.GetActionName(), raw, dataRequest.RpcRequest, dataRequest)
		response, _ := raw.(*ecs.DescribeUserDataResponse)
		d.Set("user_data", userDataHashSum(response.UserData))
	}

	if len(instance.VpcAttributes.VSwitchId) > 0 && (!d.IsNewResource() || d.HasChange("role_name")) {
		request := ecs.CreateDescribeInstanceRamRoleRequest()
		request.RegionId = client.RegionId
		request.InstanceIds = convertListToJsonString([]interface{}{d.Id()})
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeInstanceRamRole(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ecs.DescribeInstanceRamRoleResponse)
		if len(response.InstanceRamRoleSets.InstanceRamRoleSet) >= 1 {
			d.Set("role_name", response.InstanceRamRoleSets.InstanceRamRoleSet[0].RamRoleName)
		}
	}

	if instance.InstanceChargeType == string(PrePaid) {
		request := ecs.CreateDescribeInstanceAutoRenewAttributeRequest()
		request.RegionId = client.RegionId
		request.InstanceId = d.Id()
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeInstanceAutoRenewAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ecs.DescribeInstanceAutoRenewAttributeResponse)
		periodUnit := d.Get("period_unit").(string)
		if periodUnit == "" {
			periodUnit = "Month"
		}
		if len(response.InstanceRenewAttributes.InstanceRenewAttribute) > 0 {
			renew := response.InstanceRenewAttributes.InstanceRenewAttribute[0]
			d.Set("renewal_status", renew.RenewalStatus)
			d.Set("auto_renew_period", renew.Duration)
			if renew.RenewalStatus == "AutoRenewal" {
				periodUnit = renew.PeriodUnit
			}
			if periodUnit == "Year" {
				periodUnit = "Month"
				d.Set("auto_renew_period", renew.Duration*12)
			}
		}
		//period, err := computePeriodByUnit(instance.CreationTime, instance.ExpiredTime, d.Get("period").(int), periodUnit)
		//if err != nil {
		//	return WrapError(err)
		//}
		//thisPeriod := d.Get("period").(int)
		//if thisPeriod != 0 && thisPeriod != period {
		//	d.Set("period", thisPeriod)
		//} else {
		//	d.Set("period", period)
		//}
		d.Set("period_unit", periodUnit)
	}

	return nil
}

func resourceAliyunInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	d.Partial(true)

	if !d.IsNewResource() {
		if err := setTags(client, TagResourceInstance, d); err != nil {
			return WrapError(err)
		} else {
			d.SetPartial("tags")
		}
	}
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		action := "JoinResourceGroup"
		request := map[string]interface{}{
			"ResourceType":    "instance",
			"ResourceId":      d.Id(),
			"RegionId":        client.RegionId,
			"ResourceGroupId": d.Get("resource_group_id"),
		}
		conn, err := client.NewEcsClient()
		if err != nil {
			return WrapError(err)
		}
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		d.SetPartial("resource_group_id")
	}

	if err := setVolumeTags(client, TagResourceDisk, d); err != nil {
		return WrapError(err)
	} else {
		d.SetPartial("volume_tags")
	}

	if d.HasChange("security_groups") {
		if !d.IsNewResource() || d.Get("vswitch_id").(string) == "" {
			o, n := d.GetChange("security_groups")
			os := o.(*schema.Set)
			ns := n.(*schema.Set)

			rl := expandStringList(os.Difference(ns).List())
			al := expandStringList(ns.Difference(os).List())

			if len(al) > 0 {
				err := ecsService.JoinSecurityGroups(d.Id(), al)
				if err != nil {
					return WrapError(err)
				}
			}
			if len(rl) > 0 {
				err := ecsService.LeaveSecurityGroups(d.Id(), rl)
				if err != nil {
					return WrapError(err)
				}
			}

			d.SetPartial("security_groups")
		}
	}

	if !d.IsNewResource() && d.HasChange("system_disk_size") {
		diskReq := ecs.CreateDescribeDisksRequest()
		diskReq.InstanceId = d.Id()
		diskReq.DiskType = "system"
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeDisks(diskReq)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), diskReq.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(diskReq.GetActionName(), raw, diskReq.RpcRequest, diskReq)
		resp := raw.(*ecs.DescribeDisksResponse)

		instance, errDesc := ecsService.DescribeInstance(d.Id())
		if errDesc != nil {
			return WrapError(errDesc)
		}

		request := ecs.CreateResizeDiskRequest()
		request.NewSize = requests.NewInteger(d.Get("system_disk_size").(int))
		if instance.Status == string(Stopped) {
			request.Type = "offline"
		} else {
			request.Type = "online"
		}
		request.DiskId = resp.Disks.Disk[0].DiskId
		_, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ResizeDisk(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("system_disk_size")
	}

	run := false
	imageUpdate, err := modifyInstanceImage(d, meta, run)
	if err != nil {
		return WrapError(err)
	}

	vpcUpdate, err := modifyVpcAttribute(d, meta, run)
	if err != nil {
		return WrapError(err)
	}

	passwordUpdate, err := modifyInstanceAttribute(d, meta)
	if err != nil {
		return WrapError(err)
	}

	if d.HasChange("auto_release_time") {
		request := ecs.CreateModifyInstanceAutoReleaseTimeRequest()
		request.InstanceId = d.Id()
		request.RegionId = client.RegionId
		request.AutoReleaseTime = d.Get("auto_release_time").(string)
		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifyInstanceAutoReleaseTime(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("auto_release_time")
	}

	typeUpdate, err := modifyInstanceType(d, meta, run)
	if err != nil {
		return WrapError(err)
	}
	target := d.Get("status").(string)
	statusUpdate := d.HasChange("status")
	if d.IsNewResource() && target == string(Running) {
		statusUpdate = false
	}
	if imageUpdate || vpcUpdate || passwordUpdate || typeUpdate || statusUpdate {
		run = true
		instance, errDesc := ecsService.DescribeInstance(d.Id())
		if errDesc != nil {
			return WrapError(errDesc)
		}
		if (statusUpdate && target == string(Stopped)) || instance.Status == string(Running) {
			stopRequest := ecs.CreateStopInstanceRequest()
			stopRequest.RegionId = client.RegionId
			stopRequest.InstanceId = d.Id()
			stopRequest.ForceStop = requests.NewBoolean(false)
			err := resource.Retry(5*time.Minute, func() *resource.RetryError {
				raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
					return ecsClient.StopInstance(stopRequest)
				})
				if err != nil {
					if IsExpectedErrors(err, []string{"IncorrectInstanceStatus"}) {
						time.Sleep(time.Second)
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(stopRequest.GetActionName(), raw)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), stopRequest.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			stateConf := BuildStateConf([]string{"Pending", "Running", "Stopping"}, []string{"Stopped"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ecsService.InstanceStateRefreshFunc(d.Id(), []string{}))

			if _, err = stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
		if _, err := modifyInstanceImage(d, meta, run); err != nil {
			return WrapError(err)
		}

		if _, err := modifyVpcAttribute(d, meta, run); err != nil {
			return WrapError(err)
		}

		if _, err := modifyInstanceType(d, meta, run); err != nil {
			return WrapError(err)
		}

		if target == string(Running) {
			startRequest := ecs.CreateStartInstanceRequest()
			startRequest.InstanceId = d.Id()

			err := resource.Retry(5*time.Minute, func() *resource.RetryError {
				raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
					return ecsClient.StartInstance(startRequest)
				})
				if err != nil {
					if IsExpectedErrors(err, []string{"IncorrectInstanceStatus"}) {
						time.Sleep(time.Second)
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(startRequest.GetActionName(), raw)
				return nil
			})

			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), startRequest.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			// Start instance sometimes costs more than 8 minutes when os type is centos.
			stateConf := &resource.StateChangeConf{
				Pending:    []string{"Pending", "Starting", "Stopped"},
				Target:     []string{"Running"},
				Refresh:    ecsService.InstanceStateRefreshFunc(d.Id(), []string{}),
				Timeout:    d.Timeout(schema.TimeoutUpdate),
				Delay:      5 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			if _, err = stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
		if d.HasChange("status") {
			d.SetPartial("status")
		}
	}

	if err := modifyInstanceNetworkSpec(d, meta); err != nil {
		return WrapError(err)
	}

	if d.HasChange("force_delete") {
		d.SetPartial("force_delete")
	}

	if err := modifyInstanceChargeType(d, meta, false); err != nil {
		return WrapError(err)
	}

	// Only PrePaid instance can support modifying renewal attribute
	if d.Get("instance_charge_type").(string) == string(PrePaid) &&
		(d.HasChange("renewal_status") || d.HasChange("auto_renew_period")) {
		status := d.Get("renewal_status").(string)
		request := ecs.CreateModifyInstanceAutoRenewAttributeRequest()
		request.InstanceId = d.Id()
		request.RenewalStatus = status

		if status == string(RenewAutoRenewal) {
			request.PeriodUnit = d.Get("period_unit").(string)
			request.Duration = requests.NewInteger(d.Get("auto_renew_period").(int))
		}

		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifyInstanceAutoRenewAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("renewal_status")
		d.SetPartial("auto_renew_period")
	}

	d.Partial(false)
	return resourceAliyunInstanceRead(d, meta)
}

func resourceAliyunInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	if d.Get("instance_charge_type").(string) == string(PrePaid) {
		force := d.Get("force_delete").(bool)
		if !force {
			return WrapError(Error("Please convert 'PrePaid' instance to 'PostPaid' or set 'force_delete' as true before deleting 'PrePaid' instance."))
		} else if err := modifyInstanceChargeType(d, meta, force); err != nil {
			return WrapError(err)
		}
	}
	stopRequest := ecs.CreateStopInstanceRequest()
	stopRequest.InstanceId = d.Id()
	stopRequest.ForceStop = requests.NewBoolean(true)

	deleteRequest := ecs.CreateDeleteInstanceRequest()
	deleteRequest.InstanceId = d.Id()
	deleteRequest.Force = requests.NewBoolean(true)

	wait := incrementalWait(1*time.Second, 1*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteInstance(deleteRequest)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectInstanceStatus", "DependencyViolation.RouteEntry", "IncorrectInstanceStatus.Initializing"}) {
				return resource.RetryableError(err)
			}
			if IsExpectedErrors(err, []string{Throttling, "LastTokenProcessing"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(deleteRequest.GetActionName(), raw)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, EcsNotFound) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), deleteRequest.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"Pending", "Running", "Stopped", "Stopping"}, []string{}, d.Timeout(schema.TimeoutDelete), 10*time.Second, ecsService.InstanceStateRefreshFunc(d.Id(), []string{}))

	if _, err = stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func buildAliyunInstanceArgs(d *schema.ResourceData, meta interface{}) (*ecs.RunInstancesRequest, error) {
	client := meta.(*connectivity.AliyunClient)

	request := ecs.CreateRunInstancesRequest()
	request.RegionId = client.RegionId
	request.InstanceType = d.Get("instance_type").(string)

	imageID := d.Get("image_id").(string)

	request.ImageId = imageID

	systemDiskCategory := DiskCategory(d.Get("system_disk_category").(string))

	if v, ok := d.GetOk("availability_zone"); ok && v.(string) != "" {
		request.ZoneId = v.(string)
	}

	DiskName := d.Get("system_disk_name").(string)

	Description := d.Get("system_disk_description").(string)

	request.SystemDiskDiskName = DiskName

	request.SystemDiskDescription = Description

	request.SystemDiskPerformanceLevel = d.Get("system_disk_performance_level").(string)

	request.SystemDiskCategory = string(systemDiskCategory)
	request.SystemDiskSize = strconv.Itoa(d.Get("system_disk_size").(int))

	if v, ok := d.GetOk("system_disk_auto_snapshot_policy_id"); ok && v.(string) != "" {
		request.SystemDiskAutoSnapshotPolicyId = v.(string)
	}

	if v, ok := d.GetOk("security_groups"); ok {
		// At present, the classic network instance does not support multi sg in runInstances
		sgs := expandStringList(v.(*schema.Set).List())
		if d.Get("vswitch_id").(string) == "" && len(sgs) > 0 {
			request.SecurityGroupId = sgs[0]
		} else {
			request.SecurityGroupIds = &sgs
		}
	}

	if v := d.Get("instance_name").(string); v != "" {
		request.InstanceName = v
	}

	if v := d.Get("credit_specification").(string); v != "" {
		request.CreditSpecification = v
	}

	if v := d.Get("resource_group_id").(string); v != "" {
		request.ResourceGroupId = v
	}

	if v := d.Get("description").(string); v != "" {
		request.Description = v
	}

	if v := d.Get("internet_charge_type").(string); v != "" {
		request.InternetChargeType = v
	}

	request.InternetMaxBandwidthOut = requests.NewInteger(d.Get("internet_max_bandwidth_out").(int))

	if v, ok := d.GetOk("internet_max_bandwidth_in"); ok {
		request.InternetMaxBandwidthIn = requests.NewInteger(v.(int))
	}

	if v := d.Get("host_name").(string); v != "" {
		request.HostName = v
	}

	if v := d.Get("password").(string); v != "" {
		request.Password = v
	}

	if v := d.Get("kms_encrypted_password").(string); v != "" {
		kmsService := KmsService{client}
		decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
		if err != nil {
			return request, WrapError(err)
		}
		request.Password = decryptResp
	}

	vswitchValue := d.Get("subnet_id").(string)
	if vswitchValue == "" {
		vswitchValue = d.Get("vswitch_id").(string)
	}
	if vswitchValue != "" {
		request.VSwitchId = vswitchValue
		if v, ok := d.GetOk("private_ip"); ok && v.(string) != "" {
			request.PrivateIpAddress = v.(string)
		}
	}

	if v := d.Get("instance_charge_type").(string); v != "" {
		request.InstanceChargeType = v
	}

	if request.InstanceChargeType == string(PrePaid) {
		request.Period = requests.NewInteger(d.Get("period").(int))
		request.PeriodUnit = d.Get("period_unit").(string)
	} else {
		if v := d.Get("spot_strategy").(string); v != "" {
			request.SpotStrategy = v
		}
		if v := d.Get("spot_price_limit").(float64); v > 0 {
			request.SpotPriceLimit = requests.NewFloat(v)
		}
	}

	if v := d.Get("user_data").(string); v != "" {
		_, base64DecodeError := base64.StdEncoding.DecodeString(v)
		if base64DecodeError == nil {
			request.UserData = v
		} else {
			request.UserData = base64.StdEncoding.EncodeToString([]byte(v))
		}
	}

	if v := d.Get("role_name").(string); v != "" {
		request.RamRoleName = v
	}

	if v := d.Get("key_name").(string); v != "" {
		request.KeyPairName = v
	}

	if v, ok := d.GetOk("security_enhancement_strategy"); ok {
		request.SecurityEnhancementStrategy = v.(string)
	}
	if v, ok := d.GetOk("auto_release_time"); ok && v.(string) != "" {
		request.AutoReleaseTime = v.(string)
	}
	request.DryRun = requests.NewBoolean(d.Get("dry_run").(bool))
	request.DeletionProtection = requests.NewBoolean(d.Get("deletion_protection").(bool))

	if v, ok := d.GetOk("tags"); ok && len(v.(map[string]interface{})) > 0 {
		tags := make([]ecs.RunInstancesTag, 0)
		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, ecs.RunInstancesTag{
				Key:   key,
				Value: value.(string),
			})
		}
		request.Tag = &tags
	}
	request.ClientToken = buildClientToken(request.GetActionName())

	if v, ok := d.GetOk("data_disks"); ok {
		disks := v.([]interface{})
		var dataDiskRequests []ecs.RunInstancesDataDisk
		for i := range disks {
			disk := disks[i].(map[string]interface{})

			dataDiskRequest := ecs.RunInstancesDataDisk{
				Category:           disk["category"].(string),
				DeleteWithInstance: strconv.FormatBool(disk["delete_with_instance"].(bool)),
				Encrypted:          strconv.FormatBool(disk["encrypted"].(bool)),
			}

			if kmsKeyId, ok := disk["kms_key_id"]; ok {
				dataDiskRequest.KMSKeyId = kmsKeyId.(string)
			}
			if name, ok := disk["name"]; ok {
				dataDiskRequest.DiskName = name.(string)
			}
			if snapshotId, ok := disk["snapshot_id"]; ok {
				dataDiskRequest.SnapshotId = snapshotId.(string)
			}
			if description, ok := disk["description"]; ok {
				dataDiskRequest.Description = description.(string)
			}
			if autoSnapshotPolicyId, ok := disk["auto_snapshot_policy_id"]; ok {
				dataDiskRequest.AutoSnapshotPolicyId = autoSnapshotPolicyId.(string)
			}
			dataDiskRequest.Size = fmt.Sprintf("%d", disk["size"].(int))
			dataDiskRequest.Category = disk["category"].(string)
			if dataDiskRequest.Category == string(DiskEphemeralSSD) {
				dataDiskRequest.DeleteWithInstance = ""
			}
			if performanceLevel, ok := disk["performance_level"]; ok && dataDiskRequest.Category == string(DiskCloudESSD) {
				dataDiskRequest.PerformanceLevel = performanceLevel.(string)
			}

			dataDiskRequests = append(dataDiskRequests, dataDiskRequest)
		}
		request.DataDisk = &dataDiskRequests
	}
	return request, nil
}

func modifyInstanceChargeType(d *schema.ResourceData, meta interface{}, forceDelete bool) error {
	if d.IsNewResource() {
		d.Partial(false)
		return nil
	}

	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	chargeType := d.Get("instance_charge_type").(string)
	if d.HasChange("instance_charge_type") || forceDelete {
		if forceDelete {
			chargeType = string(PostPaid)
		}
		request := ecs.CreateModifyInstanceChargeTypeRequest()
		request.InstanceIds = convertListToJsonString(append(make([]interface{}, 0, 1), d.Id()))
		request.IncludeDataDisks = requests.NewBoolean(d.Get("include_data_disks").(bool))
		request.AutoPay = requests.NewBoolean(true)
		request.DryRun = requests.NewBoolean(d.Get("dry_run").(bool))
		request.ClientToken = fmt.Sprintf("terraform-modify-instance-charge-type-%s", d.Id())
		if chargeType == string(PrePaid) {
			request.Period = requests.NewInteger(d.Get("period").(int))
			request.PeriodUnit = d.Get("period_unit").(string)
		}
		request.InstanceChargeType = chargeType
		if err := resource.Retry(6*time.Minute, func() *resource.RetryError {
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ModifyInstanceChargeType(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{Throttling}) {
					time.Sleep(10 * time.Second)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		// Wait for instance charge type has been changed
		if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			if instance, err := ecsService.DescribeInstance(d.Id()); err != nil {
				return resource.NonRetryableError(err)
			} else if instance.InstanceChargeType == chargeType {
				return nil
			}
			return resource.RetryableError(Error("Waitting for instance %s to be %s timeout.", d.Id(), chargeType))
		}); err != nil {
			return WrapError(err)
		}

		d.SetPartial("instance_charge_type")
		return nil
	}

	return nil
}

func modifyInstanceImage(d *schema.ResourceData, meta interface{}, run bool) (bool, error) {
	if d.IsNewResource() {
		d.Partial(false)
		return false, nil
	}
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	update := false
	if d.HasChange("image_id") {
		update = true
		if !run {
			return update, nil
		}
		instance, err := ecsService.DescribeInstance(d.Id())
		if err != nil {
			return update, WrapError(err)
		}
		keyPairName := instance.KeyPairName
		request := ecs.CreateReplaceSystemDiskRequest()
		request.InstanceId = d.Id()
		request.ImageId = d.Get("image_id").(string)
		request.SystemDiskSize = requests.NewInteger(d.Get("system_disk_size").(int))
		request.ClientToken = buildClientToken(request.GetActionName())
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ReplaceSystemDisk(request)
		})
		if err != nil {
			return update, WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		// Ensure instance's image has been replaced successfully.
		timeout := DefaultTimeoutMedium
		for {
			instance, errDesc := ecsService.DescribeInstance(d.Id())
			if errDesc != nil {
				return update, WrapError(errDesc)
			}
			var disk ecs.Disk
			err := resource.Retry(2*time.Minute, func() *resource.RetryError {
				disk, err = ecsService.DescribeInstanceSystemDisk(d.Id(), instance.ResourceGroupId)
				if err != nil {
					if NotFoundError(err) {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			if err != nil {
				return update, WrapError(err)
			}

			if instance.ImageId == d.Get("image_id") && disk.Size == d.Get("system_disk_size").(int) {
				break
			}
			time.Sleep(DefaultIntervalShort * time.Second)

			timeout = timeout - DefaultIntervalShort
			if timeout <= 0 {
				return update, WrapError(GetTimeErrorFromString(fmt.Sprintf("Replacing instance %s system disk timeout.", d.Id())))
			}
		}

		d.SetPartial("system_disk_size")
		d.SetPartial("image_id")

		// After updating image, it need to re-attach key pair
		if keyPairName != "" {
			if err := ecsService.AttachKeyPair(keyPairName, []interface{}{d.Id()}); err != nil {
				return update, WrapError(err)
			}
		}
	}
	return update, nil
}

func modifyInstanceAttribute(d *schema.ResourceData, meta interface{}) (bool, error) {
	if d.IsNewResource() {
		d.Partial(false)
		return false, nil
	}

	update := false
	reboot := false
	request := ecs.CreateModifyInstanceAttributeRequest()
	request.InstanceId = d.Id()

	if d.HasChange("instance_name") {
		d.SetPartial("instance_name")
		request.InstanceName = d.Get("instance_name").(string)
		update = true
	}

	if d.HasChange("description") {
		d.SetPartial("description")
		request.Description = d.Get("description").(string)
		update = true
	}

	if d.HasChange("user_data") {
		d.SetPartial("user_data")
		if v, ok := d.GetOk("user_data"); ok && v.(string) != "" {
			_, base64DecodeError := base64.StdEncoding.DecodeString(v.(string))
			if base64DecodeError == nil {
				request.UserData = v.(string)
			} else {
				request.UserData = base64.StdEncoding.EncodeToString([]byte(v.(string)))
			}
		}
		update = true
		reboot = true
	}

	if d.HasChange("host_name") {
		d.SetPartial("host_name")
		request.HostName = d.Get("host_name").(string)
		update = true
		reboot = true
	}

	if d.HasChange("password") || d.HasChange("kms_encrypted_password") {
		if v := d.Get("password").(string); v != "" {
			d.SetPartial("password")
			request.Password = v
			update = true
			reboot = true
		}
		if v := d.Get("kms_encrypted_password").(string); v != "" {
			kmsService := KmsService{meta.(*connectivity.AliyunClient)}
			decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return reboot, WrapError(err)
			}
			request.Password = decryptResp
			d.SetPartial("kms_encrypted_password")
			d.SetPartial("kms_encryption_context")
			update = true
			reboot = true
		}
	}

	if d.HasChange("deletion_protection") {
		d.SetPartial("deletion_protection")
		request.DeletionProtection = requests.NewBoolean(d.Get("deletion_protection").(bool))
		update = true
	}

	if d.HasChange("credit_specification") {
		d.SetPartial("credit_specification")
		request.CreditSpecification = d.Get("credit_specification").(string)
		update = true
	}

	client := meta.(*connectivity.AliyunClient)
	if update {
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ModifyInstanceAttribute(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"InvalidChargeType.ValueNotSupported"}) {
					time.Sleep(time.Minute)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})
		if err != nil {
			return reboot, WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return reboot, nil
}

func modifyVpcAttribute(d *schema.ResourceData, meta interface{}, run bool) (bool, error) {
	if d.IsNewResource() {
		d.Partial(false)
		return false, nil
	}

	update := false
	request := ecs.CreateModifyInstanceVpcAttributeRequest()
	request.InstanceId = d.Id()
	request.VSwitchId = d.Get("vswitch_id").(string)

	if d.HasChange("vswitch_id") {
		update = true
		if d.Get("vswitch_id").(string) == "" {
			return update, WrapError(Error("Field 'vswitch_id' is required when modifying the instance VPC attribute."))
		}
		d.SetPartial("vswitch_id")
	}

	if d.HasChange("subnet_id") {
		update = true
		if d.Get("subnet_id").(string) == "" {
			return update, WrapError(Error("Field 'subnet_id' is required when modifying the instance VPC attribute."))
		}
		request.VSwitchId = d.Get("subnet_id").(string)
		d.SetPartial("subnet_id")
	}

	if request.VSwitchId != "" && d.HasChange("private_ip") {
		request.PrivateIpAddress = d.Get("private_ip").(string)
		update = true
		d.SetPartial("private_ip")
	}

	if !run {
		return update, nil
	}

	if update {
		client := meta.(*connectivity.AliyunClient)
		err := resource.Retry(1*time.Minute, func() *resource.RetryError {
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ModifyInstanceVpcAttribute(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"OperationConflict"}) {
					time.Sleep(1 * time.Second)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})

		if err != nil {
			return update, WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		ecsService := EcsService{client}
		if err := ecsService.WaitForVpcAttributesChanged(d.Id(), request.VSwitchId, request.PrivateIpAddress); err != nil {
			return update, WrapError(err)
		}
	}
	return update, nil
}

func modifyInstanceType(d *schema.ResourceData, meta interface{}, run bool) (bool, error) {
	if d.IsNewResource() {
		d.Partial(false)
		return false, nil
	}
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	update := false
	if d.HasChange("instance_type") {
		update = true
		if !run {
			return update, nil
		}
		// Ensure instance_type is valid
		//zoneId, validZones, err := ecsService.DescribeAvailableResources(d, meta, InstanceTypeResource)
		//if err != nil {
		//	return update, WrapError(err)
		//}
		//if err = ecsService.InstanceTypeValidation(d.Get("instance_type").(string), zoneId, validZones); err != nil {
		//	return update, WrapError(err)
		//}

		// There should use the old instance charge type to decide API method because of instance_charge_type will be updated at last step
		oldCharge, _ := d.GetChange("instance_charge_type")
		if oldCharge.(string) == string(PrePaid) {
			request := ecs.CreateModifyPrepayInstanceSpecRequest()
			request.InstanceId = d.Id()
			request.InstanceType = d.Get("instance_type").(string)

			err := resource.Retry(6*time.Minute, func() *resource.RetryError {
				raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
					return ecsClient.ModifyPrepayInstanceSpec(request)
				})
				if err != nil {
					if IsExpectedErrors(err, []string{Throttling}) {
						time.Sleep(5 * time.Second)
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				return nil
			})
			if err != nil {
				return update, WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
		} else {
			//An instance that was successfully modified once cannot be modified again within 5 minutes.
			request := ecs.CreateModifyInstanceSpecRequest()
			request.InstanceId = d.Id()
			request.InstanceType = d.Get("instance_type").(string)
			request.ClientToken = buildClientToken(request.GetActionName())

			err := resource.Retry(6*time.Minute, func() *resource.RetryError {
				args := *request
				raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
					return ecsClient.ModifyInstanceSpec(&args)
				})
				if err != nil {
					if IsExpectedErrors(err, []string{Throttling}) {
						time.Sleep(10 * time.Second)
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				return nil
			})
			if err != nil {
				return update, WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
		}

		// Ensure instance's type has been replaced successfully.
		timeout := DefaultTimeoutMedium
		for {
			instance, err := ecsService.DescribeInstance(d.Id())

			if err != nil {
				return update, WrapError(err)
			}

			if instance.InstanceType == d.Get("instance_type").(string) {
				break
			}

			timeout = timeout - DefaultIntervalShort
			if timeout <= 0 {
				return update, WrapErrorf(err, WaitTimeoutMsg, d.Id(), GetFunc(1), timeout, instance.InstanceType, d.Get("instance_type"), ProviderERROR)
			}

			time.Sleep(DefaultIntervalShort * time.Second)
		}
		d.SetPartial("instance_type")
	}
	return update, nil
}

func modifyInstanceNetworkSpec(d *schema.ResourceData, meta interface{}) error {
	if d.IsNewResource() {
		d.Partial(false)
		return nil
	}

	allocate := false
	update := false
	request := ecs.CreateModifyInstanceNetworkSpecRequest()
	request.InstanceId = d.Id()
	request.ClientToken = buildClientToken(request.GetActionName())

	if d.HasChange("internet_charge_type") {
		request.NetworkChargeType = d.Get("internet_charge_type").(string)
		update = true
		d.SetPartial("internet_charge_type")
	}

	if d.HasChange("internet_max_bandwidth_out") {
		o, n := d.GetChange("internet_max_bandwidth_out")
		if o.(int) <= 0 && n.(int) > 0 {
			allocate = true
		}
		request.InternetMaxBandwidthOut = requests.NewInteger(n.(int))
		update = true
		d.SetPartial("internet_max_bandwidth_out")
	}

	if d.HasChange("internet_max_bandwidth_in") {
		request.InternetMaxBandwidthIn = requests.NewInteger(d.Get("internet_max_bandwidth_in").(int))
		update = true
		d.SetPartial("internet_max_bandwidth_in")
	}

	//An instance that was successfully modified once cannot be modified again within 5 minutes.
	wait := incrementalWait(2*time.Second, 2*time.Second)
	client := meta.(*connectivity.AliyunClient)
	if update {
		if err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ModifyInstanceNetworkSpec(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{Throttling, "LastOrderProcessing", "LastRequestProcessing", "LastTokenProcessing"}) {
					wait()
					return resource.RetryableError(err)
				}
				if IsExpectedErrors(err, []string{"InternalError"}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		ecsService := EcsService{client: client}

		deadline := time.Now().Add(DefaultTimeout * time.Second)
		for {
			instance, err := ecsService.DescribeInstance(d.Id())
			if err != nil {
				return WrapError(err)
			}

			if instance.InternetMaxBandwidthOut == d.Get("internet_max_bandwidth_out").(int) &&
				instance.InternetChargeType == d.Get("internet_charge_type").(string) {
				break
			}

			if time.Now().After(deadline) {
				return WrapError(Error(`wait for internet update timeout! expect internet_charge_type value %s, get %s
					expect internet_max_bandwidth_out value %d, get %d,`,
					d.Get("internet_charge_type").(string), instance.InternetChargeType, d.Get("internet_max_bandwidth_out").(int),
					instance.InternetMaxBandwidthOut))
			}
			time.Sleep(1 * time.Second)
		}

		if allocate {
			request := ecs.CreateAllocatePublicIpAddressRequest()
			request.InstanceId = d.Id()
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.AllocatePublicIpAddress(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}
	}
	return nil
}
