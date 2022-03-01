package alicloud

import (
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/hbase"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudHBaseInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudHBaseInstanceCreate,
		Read:   resourceAlicloudHBaseInstanceRead,
		Update: resourceAlicloudHBaseInstanceUpdate,
		Delete: resourceAlicloudHBaseInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"engine": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"hbase", "hbaseue", "bds"}, false),
				Default:      "hbase",
			},
			"engine_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"master_instance_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"master_instance_quantity": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"core_instance_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"core_instance_quantity": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(1, 200),
				Optional:     true,
				Default:      2,
			},
			"core_disk_type": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ValidateFunc:     validation.StringInSlice([]string{"cloud_ssd", "cloud_essd_pl1", "cloud_efficiency", "local_hdd_pro", "local_ssd_pro", ""}, false),
				DiffSuppressFunc: engineDiffSuppressFunc,
			},
			"core_disk_size": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateFunc:     validation.Any(validation.IntBetween(20, 64000), validation.IntInSlice([]int{0})),
				Default:          400,
				DiffSuppressFunc: engineDiffSuppressFunc,
			},
			"pay_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{string(PostPaid), string(PrePaid)}, false),
				Default:      PostPaid,
			},
			"duration": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				DiffSuppressFunc: payTypePostPaidDiffSuppressFunc,
			},
			"auto_renew": {
				Type:             schema.TypeBool,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				DiffSuppressFunc: payTypePostPaidDiffSuppressFunc,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cold_storage_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.Any(validation.IntBetween(800, 1000000), validation.IntInSlice([]int{0})),
				Default:      0,
			},
			"maintain_start_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"maintain_end_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"immediate_delete_flag": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"tags": tagsSchema(),
			"account": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 128),
			},
			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringLenBetween(0, 128),
			},
			"ip_white": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: whiteIpListDiffSuppressFunc,
			},
			"security_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ui_proxy_conn_addrs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"conn_addr_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"conn_addr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"net_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"zk_conn_addrs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"conn_addr_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"conn_addr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"net_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"slb_conn_addrs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"conn_addr_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"conn_addr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"net_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func checkParams(request *hbase.CreateClusterRequest) error {
	if request.Engine == "bds" && request.VSwitchId == "" {
		return WrapError(Error("bds is not support classic"))
	}
	return nil
}

func buildHBaseCreateRequest(d *schema.ResourceData, meta interface{}) (*hbase.CreateClusterRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	request := hbase.CreateCreateClusterRequest()
	request.ClusterName = Trim(d.Get("name").(string))
	request.RegionId = string(client.Region)
	request.ZoneId = Trim(d.Get("zone_id").(string))
	request.Engine = Trim(d.Get("engine").(string))
	request.EngineVersion = Trim(d.Get("engine_version").(string))
	request.MasterInstanceType = Trim(d.Get("master_instance_type").(string))
	request.CoreInstanceType = Trim(d.Get("core_instance_type").(string))
	request.NodeCount = requests.NewInteger(d.Get("core_instance_quantity").(int))
	request.DiskType = Trim(d.Get("core_disk_type").(string))
	request.DiskSize = requests.NewInteger(d.Get("core_disk_size").(int))
	request.PayType = Trim(d.Get("pay_type").(string))
	if d.Get("duration").(int) > 9 {
		request.PeriodUnit = "year"
		request.Period = requests.NewInteger(d.Get("duration").(int) / 12)
	} else {
		request.PeriodUnit = "month"
		request.Period = requests.NewInteger(d.Get("duration").(int))
	}

	if d.Get("auto_renew").(bool) {
		request.AutoRenewPeriod = requests.NewInteger(1)
	}

	vswitchId := Trim(d.Get("vswitch_id").(string))
	if vswitchId != "" {
		request.VSwitchId = vswitchId
		// check vswitchId in zone
		vsw, err := vpcService.DescribeVSwitch(vswitchId)
		if err != nil {
			return nil, WrapError(err)
		}

		if request.ZoneId == "" {
			request.ZoneId = vsw.ZoneId
		} else if strings.Contains(request.ZoneId, MULTI_IZ_SYMBOL) {
			zonestr := strings.Split(strings.SplitAfter(request.ZoneId, "(")[1], ")")[0]
			if !strings.Contains(zonestr, string([]byte(vsw.ZoneId)[len(vsw.ZoneId)-1])) {
				return nil, WrapError(Error("The specified vswitch %s isn't in the multi zone %s.", vsw.VSwitchId, request.ZoneId))
			}
		} else if request.ZoneId != vsw.ZoneId {
			return nil, WrapError(Error("The specified vswitch %s isn't in the zone %s.", vsw.VSwitchId, request.ZoneId))
		}

		request.VpcId = vsw.VpcId
	}

	if d.Get("cold_storage_size").(int) < 0 {
		return nil, WrapError(Error("cold_storage_size=%s is invalid", d.Get("cold_storage_size")))
	}
	request.ColdStorageSize = requests.NewInteger(d.Get("cold_storage_size").(int))

	request.SecurityIPList = d.Get("ip_white").(string)
	return request, checkParams(request)
}

func resourceAlicloudHBaseInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hBaseService := HBaseService{client}

	request, err := buildHBaseCreateRequest(d, meta)
	if err != nil {
		return WrapError(err)
	}

	raw, err := client.WithHbaseClient(func(client *hbase.Client) (interface{}, error) {
		return client.CreateCluster(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_hbase_instance", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*hbase.CreateClusterResponse)

	d.SetId(response.ClusterId)

	stateConf := BuildStateConf([]string{Hb_LAUNCHING, Hb_CREATING}, []string{Hb_ACTIVATION}, d.Timeout(schema.TimeoutCreate),
		5*time.Minute, hBaseService.HBaseClusterStateRefreshFunc(d.Id(), []string{Hb_CREATE_FAILED}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudHBaseInstanceUpdate(d, meta)
}

func resourceAlicloudHBaseInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbaseService := HBaseService{client}

	instance, err := hbaseService.DescribeHBaseInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", instance["InstanceName"])
	d.Set("zone_id", instance["ZoneId"])
	d.Set("engine", instance["Engine"])
	d.Set("engine_version", instance["MajorVersion"])
	d.Set("master_instance_type", instance["MasterInstanceType"])
	d.Set("master_instance_quantity", instance["MasterNodeCount"])
	d.Set("core_instance_type", instance["CoreInstanceType"])
	d.Set("core_instance_quantity", instance["CoreNodeCount"])
	d.Set("core_disk_size", formatInt(instance["CoreDiskCount"])*formatInt(instance["CoreDiskSize"]))
	d.Set("core_disk_type", instance["CoreDiskType"])
	// Postpaid -> PostPaid
	if instance["PayType"] == string(Postpaid) {
		d.Set("pay_type", string(PostPaid))
	} else if instance["PayType"] == string(Prepaid) {
		d.Set("pay_type", string(PrePaid))
		period, err := computePeriodByUnit(instance["CreatedTimeUTC"], instance["ExpireTimeUTC"], d.Get("duration").(int), "Month")
		if err != nil {
			return WrapError(err)
		}
		d.Set("duration", period)
	}
	// now sdk can not get right value, "auto_renew", "is_cold_storage".
	d.Set("auto_renew", instance["AutoRenewal"])
	d.Set("cold_storage_size", instance["ColdStorageSize"])
	d.Set("vpc_id", instance["VpcId"])
	d.Set("vswitch_id", instance["VswitchId"])
	d.Set("maintain_start_time", instance["MaintainStartTime"])
	d.Set("maintain_end_time", instance["MaintainEndTime"])
	d.Set("deletion_protection", instance["IsDeletionProtection"])
	d.Set("tags", tagsToMap(instance["Tags"]))

	ipWhitelist, err := hbaseService.DescribeIpWhitelist(d.Id())
	if err != nil {
		return WrapError(err)
	}
	ipWhite := LOCAL_HOST_IP
	for _, value := range ipWhitelist.Groups.Group {
		if value.GroupName == "default" {
			ipWhite = strings.Join(value.IpList.Ip, ",")
		}
	}
	d.Set("ip_white", ipWhite)

	securityGroups, err := hbaseService.DescribeSecurityGroups(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("security_groups", securityGroups.SecurityGroupIds.SecurityGroupId)

	conns, err := hbaseService.DescribeClusterConnection(d.Id())
	if err != nil {
		return WrapError(err)
	}
	// UiProxyConnAddrInfo
	var uiConnAddrs []map[string]interface{}
	e := map[string]interface{}{
		"conn_addr_port": conns.UiProxyConnAddrInfo.ConnAddrPort,
		"conn_addr":      conns.UiProxyConnAddrInfo.ConnAddr,
		"net_type":       conns.UiProxyConnAddrInfo.NetType,
	}
	uiConnAddrs = append(uiConnAddrs, e)
	d.Set("ui_proxy_conn_addrs", uiConnAddrs)
	// ZKConnAddrs
	var zkConnAddrs []map[string]interface{}
	for _, v := range conns.ZkConnAddrs.ZkConnAddr {
		e := map[string]interface{}{
			"conn_addr_port": v.ConnAddrPort,
			"conn_addr":      v.ConnAddr,
			"net_type":       v.NetType,
		}
		zkConnAddrs = append(zkConnAddrs, e)
	}
	d.Set("zk_conn_addrs", zkConnAddrs)
	// slbConnAddrs
	var slbConnAddrs []map[string]interface{}
	for _, v := range conns.SlbConnAddrs.SlbConnAddr {
		e := map[string]interface{}{
			"conn_addr_port": v.ConnAddrInfo.ConnAddrPort,
			"conn_addr":      v.ConnAddrInfo.ConnAddr,
			"net_type":       v.ConnAddrInfo.NetType,
		}
		slbConnAddrs = append(slbConnAddrs, e)
	}
	d.Set("slb_conn_addrs", slbConnAddrs)
	return nil
}

func resourceAlicloudHBaseInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hBaseService := HBaseService{client}
	d.Partial(true)
	if d.HasChange("pay_type") {
		object, err := hBaseService.DescribeHBaseInstance(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := strings.ToLower(d.Get("pay_type").(string))
		if strings.ToLower(object["PayType"].(string)) != target {
			request := hbase.CreateConvertInstanceRequest()
			request.ClusterId = d.Id()
			request.PayType = string(Postpaid)
			if target == "prepaid" {
				request.PayType = string(Prepaid)
			}
			if d.Get("duration").(int) > 9 {
				request.PricingCycle = "year"
				request.Duration = requests.NewInteger(d.Get("duration").(int) / 12)
			} else {
				request.PricingCycle = "month"
				request.Duration = requests.NewInteger(d.Get("duration").(int))
			}
			raw, err := client.WithHbaseClient(func(hbaseClient *hbase.Client) (interface{}, error) {
				return hbaseClient.ConvertInstance(request)
			})
			addDebug(request.GetActionName(), raw)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			stateConf := BuildStateConf([]string{}, []string{Hb_ACTIVATION}, d.Timeout(schema.TimeoutUpdate),
				2*time.Minute, hBaseService.HBaseClusterStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
			d.SetPartial("pay_type")
		}
	}

	if d.HasChange("ip_white") {
		request := hbase.CreateModifyIpWhitelistRequest()
		request.ClusterId = d.Id()
		request.IpList = d.Get("ip_white").(string)
		request.GroupName = "default"
		request.IpVersion = "4"
		raw, err := client.WithHbaseClient(func(hbaseClient *hbase.Client) (interface{}, error) {
			return hbaseClient.ModifyIpWhitelist(request)
		})
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("ip_white")
	}

	if d.HasChange("security_groups") {
		request := hbase.CreateModifySecurityGroupsRequest()
		request.ClusterId = d.Id()
		securityGroups := d.Get("security_groups").(*schema.Set).List()
		if securityGroups != nil && len(securityGroups) > 0 {
			request.SecurityGroupIds = strings.Join(expandStringList(securityGroups), ",")
		} else {
			request.SecurityGroupIds = " "
		}
		raw, err := client.WithHbaseClient(func(hbaseClient *hbase.Client) (interface{}, error) {
			return hbaseClient.ModifySecurityGroups(request)
		})
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{Hb_ACTIVATION}, d.Timeout(schema.TimeoutUpdate),
			5*time.Second, hBaseService.HBaseClusterStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("security_groups")
	}

	if d.HasChange("maintain_start_time") || d.HasChange("maintain_end_time") {
		request := hbase.CreateModifyInstanceMaintainTimeRequest()
		request.ClusterId = d.Id()
		request.MaintainStartTime = d.Get("maintain_start_time").(string)
		request.MaintainEndTime = d.Get("maintain_end_time").(string)

		raw, err := client.WithHbaseClient(func(client *hbase.Client) (interface{}, error) {
			return client.ModifyInstanceMaintainTime(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("maintain_start_time")
		d.SetPartial("maintain_end_time")
	}

	if d.HasChange("deletion_protection") {
		if err := hBaseService.ModifyClusterDeletionProtection(d.Id(), d.Get("deletion_protection").(bool)); err != nil {
			return WrapError(err)
		}
		d.SetPartial("deletion_protection")
	}

	if err := hBaseService.setInstanceTags(d); err != nil {
		return WrapError(err)
	}

	if d.HasChange("account") || d.HasChange("password") {
		request := hbase.CreateModifyUIAccountPasswordRequest()
		request.ClusterId = d.Id()
		request.AccountName = d.Get("account").(string)
		request.AccountPassword = d.Get("password").(string)

		raw, err := client.WithHbaseClient(func(client *hbase.Client) (interface{}, error) {
			return client.ModifyUIAccountPassword(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("account")
		d.SetPartial("password")
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlicloudHBaseInstanceRead(d, meta)
	}

	if d.HasChange("name") {
		request := hbase.CreateModifyInstanceNameRequest()
		request.ClusterId = d.Id()
		request.ClusterName = d.Get("name").(string)

		raw, err := client.WithHbaseClient(func(client *hbase.Client) (interface{}, error) {
			return client.ModifyInstanceName(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("name")
	}

	if d.HasChange("core_instance_quantity") {
		request := hbase.CreateResizeNodeCountRequest()
		request.ClusterId = d.Id()
		request.NodeCount = requests.NewInteger(d.Get("core_instance_quantity").(int))

		raw, err := client.WithHbaseClient(func(client *hbase.Client) (interface{}, error) {
			return client.ResizeNodeCount(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		// Cumbersome operation，async call, wait for state change
		// wait instance status is running after modifying
		stateConf := BuildStateConf([]string{Hb_NODE_RESIZING}, []string{Hb_ACTIVATION}, d.Timeout(schema.TimeoutUpdate),
			5*time.Minute, hBaseService.HBaseClusterStateRefreshFunc(d.Id(), []string{Hb_NODE_RESIZING_FAILED}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("core_instance_quantity")
	}

	if d.HasChange("master_instance_type") || d.HasChange("core_instance_type") {
		request := hbase.CreateModifyInstanceTypeRequest()
		request.ClusterId = d.Id()
		typeChange := false
		if d.HasChange("master_instance_type") &&
			d.Get("engine") != "bds" && d.Get("core_instance_quantity").(int) > 1 {
			request.MasterInstanceType = d.Get("master_instance_type").(string)
			typeChange = true
		}

		if d.HasChange("core_instance_type") {
			request.CoreInstanceType = d.Get("core_instance_type").(string)
			typeChange = true
		}

		if typeChange {
			raw, err := client.WithHbaseClient(func(client *hbase.Client) (interface{}, error) {
				return client.ModifyInstanceType(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			// Cumbersome operation，async call, wait for state change
			// wait instance status is running after modifying
			stateConf := BuildStateConf([]string{Hb_LEVEL_MODIFY}, []string{Hb_ACTIVATION}, d.Timeout(schema.TimeoutUpdate),
				5*time.Minute, hBaseService.HBaseClusterStateRefreshFunc(d.Id(), []string{Hb_LEVEL_MODIFY_FAILED}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			d.SetPartial("master_instance_type")
			d.SetPartial("core_instance_type")
		}
	}

	if d.HasChange("core_disk_size") && d.Get("engine") != "bds" {
		request := hbase.CreateResizeDiskSizeRequest()
		request.ClusterId = d.Id()
		request.NodeDiskSize = requests.NewInteger(d.Get("core_disk_size").(int))

		raw, err := client.WithHbaseClient(func(client *hbase.Client) (interface{}, error) {
			return client.ResizeDiskSize(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		// Cumbersome operation，async call, wait for state change
		// wait instance status is running after modifying
		stateConf := BuildStateConf([]string{Hb_DISK_RESIZING}, []string{Hb_ACTIVATION}, d.Timeout(schema.TimeoutUpdate),
			2*time.Minute, hBaseService.HBaseClusterStateRefreshFunc(d.Id(), []string{Hb_DISK_RESIZE_FAILED}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("core_disk_size")
	}

	if d.HasChange("cold_storage_size") {
		request := hbase.CreateResizeColdStorageSizeRequest()
		request.ClusterId = d.Id()
		request.ColdStorageSize = requests.NewInteger(d.Get("cold_storage_size").(int))

		raw, err := client.WithHbaseClient(func(client *hbase.Client) (interface{}, error) {
			return client.ResizeColdStorageSize(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		// Cumbersome operation，async call, wait for state change
		// wait instance status is running after modifying
		stateConf := BuildStateConf([]string{Hb_HBASE_COLD_EXPANDING}, []string{Hb_ACTIVATION}, d.Timeout(schema.TimeoutUpdate),
			10*time.Second, hBaseService.HBaseClusterStateRefreshFunc(d.Id(), []string{Hb_DISK_RESIZE_FAILED}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("cold_storage_size")
	}

	d.Partial(false)
	return resourceAlicloudHBaseInstanceRead(d, meta)
}

func resourceAlicloudHBaseInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbaseService := HBaseService{client}

	request := hbase.CreateDeleteInstanceRequest()
	request.ClusterId = d.Id()
	if v, ok := d.GetOk("immediate_delete_flag"); ok {
		request.ImmediateDeleteFlag = requests.NewBoolean(v.(bool))
	}

	raw, err := client.WithHbaseClient(func(hbaseClient *hbase.Client) (interface{}, error) {
		return hbaseClient.DeleteInstance(request)
	})

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"Instance.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{Hb_DELETING}, []string{}, d.Timeout(schema.TimeoutDelete), 1*time.Minute, hbaseService.HBaseClusterStateRefreshFunc(d.Id(), []string{}))
	_, err = stateConf.WaitForState()
	return WrapError(err)
}
