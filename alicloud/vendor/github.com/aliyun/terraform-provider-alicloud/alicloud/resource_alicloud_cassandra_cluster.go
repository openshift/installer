package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cassandra"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudCassandraCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCassandraClusterCreate,
		Read:   resourceAlicloudCassandraClusterRead,
		Update: resourceAlicloudCassandraClusterUpdate,
		Delete: resourceAlicloudCassandraClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"auto_renew_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 3, 6, 12}),
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_center_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disk_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"disk_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"maintain_end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"maintain_start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"major_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"node_count": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"pay_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"period_unit": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enable_public": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"public_points": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ip_white": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"security_groups": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAlicloudCassandraClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cassandraService := CassandraService{client}

	request := cassandra.CreateCreateClusterRequest()
	if v, ok := d.GetOkExists("auto_renew"); ok {
		request.AutoRenew = requests.NewBoolean(v.(bool))
	}
	if v, ok := d.GetOk("auto_renew_period"); ok {
		request.AutoRenewPeriod = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("cluster_name"); ok {
		request.ClusterName = v.(string)
	}
	if v, ok := d.GetOk("data_center_name"); ok {
		request.DataCenterName = v.(string)
	}
	if v, ok := d.GetOk("disk_size"); ok {
		request.DiskSize = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("disk_type"); ok {
		request.DiskType = v.(string)
	}
	request.InstanceType = d.Get("instance_type").(string)
	request.MajorVersion = d.Get("major_version").(string)
	request.NodeCount = requests.NewInteger(d.Get("node_count").(int))
	if v, ok := d.GetOk("password"); ok {
		request.Password = v.(string)
	}
	request.PayType = d.Get("pay_type").(string)
	if v, ok := d.GetOk("period"); ok {
		request.Period = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("period_unit"); ok {
		request.PeriodUnit = v.(string)
	}
	request.RegionId = client.RegionId
	request.VswitchId = d.Get("vswitch_id").(string)
	request.ZoneId = d.Get("zone_id").(string)
	vswitchId := Trim(d.Get("vswitch_id").(string))
	if vswitchId != "" {
		vpcService := VpcService{client}
		request.VswitchId = vswitchId
		// check vswitchId in zone
		vsw, err := vpcService.DescribeVSwitch(vswitchId)
		if err != nil {
			return WrapError(err)
		}

		if request.ZoneId == "" {
			request.ZoneId = vsw.ZoneId
		}
		request.VpcId = vsw.VpcId
	}
	raw, err := client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
		return cassandraClient.CreateCluster(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cassandra_cluster", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*cassandra.CreateClusterResponse)
	d.SetId(fmt.Sprintf("%v", response.ClusterId))
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 8*time.Minute, cassandraService.CassandraClusterStateRefreshFunc(d.Id(), []string{"CreateFailed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudCassandraClusterUpdate(d, meta)
}
func resourceAlicloudCassandraClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cassandraService := CassandraService{client}
	object, err := cassandraService.DescribeCassandraCluster(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cassandra_cluster cassandraService.DescribeCassandraCluster Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("id", object.ClusterId)
	d.Set("cluster_name", object.ClusterName)
	d.Set("maintain_end_time", object.MaintainEndTime)
	d.Set("maintain_start_time", object.MaintainStartTime)
	d.Set("major_version", object.MajorVersion)
	d.Set("pay_type", object.PayType)
	d.Set("status", object.Status)
	d.Set("tags", cassandraService.tagsToMap(object.Tags.Tag))
	firstDc, err := cassandraService.DescribeClusterDataCenter(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("zone_id", firstDc.ZoneId)
	d.Set("vswitch_id", firstDc.VswitchId)
	d.Set("data_center_name", firstDc.DataCenterName)
	d.Set("data_center_id", firstDc.DataCenterId)
	d.Set("disk_type", firstDc.DiskType)
	d.Set("disk_size", firstDc.DiskSize)
	d.Set("instance_type", firstDc.InstanceType)
	d.Set("node_count", firstDc.NodeCount)
	d.Set("auto_renew", firstDc.AutoRenewal)
	if d.Get("auto_renew").(bool) {
		d.Set("auto_renew_period", firstDc.AutoRenewPeriod)
	} else {
		if v, ok := d.GetOkExists("auto_renew_period"); ok {
			d.Set("auto_renew_period", v.(int))
		}
	}

	d.Set("enable_public", false)
	d.Set("public_points", nil)
	d.Set("security_groups", nil)
	ipWhitelist, err := cassandraService.DescribeIpWhitelist(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("ip_white", strings.Join(ipWhitelist.IpList.IP, ","))
	securityGroups, err := cassandraService.DescribeSecurityGroups(d.Id())
	if err != nil {
		return WrapError(err)
	}
	if len(securityGroups.SecurityGroupIds.SecurityGroupId) > 0 {
		d.Set("security_groups", securityGroups.SecurityGroupIds.SecurityGroupId)
	}
	endPoints, err := cassandraService.DescribeCassandraEndpoints(d.Id())
	if err != nil {
		return WrapError(err)
	}
	for _, contactPoint := range endPoints.ContactPoints.ContactPoint {
		if d.Get("zone_id").(string) == contactPoint.DataCenterId && len(contactPoint.PublicAddresses.PublicAddress) > 0 {
			d.Set("enable_public", true)
			d.Set("public_points", contactPoint.PublicAddresses.PublicAddress)
			break
		}
	}
	return nil
}
func resourceAlicloudCassandraClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cassandraService := CassandraService{client}
	d.Partial(true)

	if err := cassandraService.setInstanceTags(d); err != nil {
		return WrapError(err)
	}
	if d.HasChange("ip_white") {
		modifyIpWhitelistGroupReq := cassandra.CreateModifyIpWhitelistGroupRequest()
		modifyIpWhitelistGroupReq.ClusterId = d.Id()
		modifyIpWhitelistGroupReq.IpList = d.Get("ip_white").(string)
		modifyIpWhitelistGroupReq.GroupName = "default"
		modifyIpWhitelistGroupReq.IpVersion = "4"
		raw, err := client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
			return cassandraClient.ModifyIpWhitelistGroup(modifyIpWhitelistGroupReq)
		})
		addDebug(modifyIpWhitelistGroupReq.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifyIpWhitelistGroupReq.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("ip_white")
	}
	if d.HasChange("security_groups") {
		modifySecurityGroupsReq := cassandra.CreateModifySecurityGroupsRequest()
		modifySecurityGroupsReq.ClusterId = d.Id()
		securityGroups := d.Get("security_groups").([]interface{})
		if securityGroups != nil && len(securityGroups) > 0 {
			security_groups := strings.Join(expandStringList(securityGroups), ",")
			modifySecurityGroupsReq.SecurityGroupIds = security_groups
		} else {
			modifySecurityGroupsReq.SecurityGroupIds = " "
		}
		raw, err := client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
			return cassandraClient.ModifySecurityGroups(modifySecurityGroupsReq)
		})
		addDebug(modifySecurityGroupsReq.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifySecurityGroupsReq.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cassandraService.CassandraClusterStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("security_groups")
	}
	if d.HasChange("enable_public") {
		if d.Get("enable_public").(bool) {
			request := cassandra.CreateAllocatePublicContactPointsRequest()
			request.ClusterId = d.Id()
			request.DataCenterId = d.Get("zone_id").(string)
			raw, err := client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
				return cassandraClient.AllocatePublicContactPoints(request)
			})
			addDebug(request.GetActionName(), raw)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cassandraService.CassandraClusterStateRefreshFunc(d.Id(), []string{"AllocatePubConnFailed"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		} else {
			request := cassandra.CreateReleasePublicContactPointsRequest()
			request.ClusterId = d.Id()
			request.DataCenterId = d.Get("zone_id").(string)
			raw, err := client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
				return cassandraClient.ReleasePublicContactPoints(request)
			})
			addDebug(request.GetActionName(), raw)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cassandraService.CassandraClusterStateRefreshFunc(d.Id(), []string{"ReleasePubConnFailed"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
		d.SetPartial("enable_public")
	}
	if (d.HasChange("maintain_end_time") || d.HasChange("maintain_start_time")) && d.Get("maintain_end_time").(string) != "" && d.Get("maintain_start_time").(string) != "" {
		modifyInstanceMaintainTimeReq := cassandra.CreateModifyInstanceMaintainTimeRequest()
		modifyInstanceMaintainTimeReq.ClusterId = d.Id()
		modifyInstanceMaintainTimeReq.MaintainEndTime = d.Get("maintain_end_time").(string)
		modifyInstanceMaintainTimeReq.MaintainStartTime = d.Get("maintain_start_time").(string)
		raw, err := client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
			return cassandraClient.ModifyInstanceMaintainTime(modifyInstanceMaintainTimeReq)
		})
		addDebug(modifyInstanceMaintainTimeReq.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifyInstanceMaintainTimeReq.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("maintain_end_time")
		d.SetPartial("maintain_start_time")
	}
	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlicloudCassandraClusterRead(d, meta)
	}
	if d.HasChange("cluster_name") {
		request := cassandra.CreateModifyClusterRequest()
		request.ClusterId = d.Id()
		request.ClusterName = d.Get("cluster_name").(string)
		raw, err := client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
			return cassandraClient.ModifyCluster(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("cluster_name")
	}
	if d.HasChange("data_center_name") {
		request := cassandra.CreateModifyDataCenterRequest()
		request.ClusterId = d.Id()
		request.DataCenterId = d.Get("zone_id").(string)
		request.DataCenterName = d.Get("data_center_name").(string)
		raw, err := client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
			return cassandraClient.ModifyDataCenter(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("data_center_name")
	}
	update := false
	modifyInstanceTypeReq := cassandra.CreateModifyInstanceTypeRequest()
	modifyInstanceTypeReq.ClusterId = d.Id()
	if d.HasChange("instance_type") {
		update = true
	}
	modifyInstanceTypeReq.InstanceType = d.Get("instance_type").(string)
	modifyInstanceTypeReq.DataCenterId = d.Get("zone_id").(string)
	if update {
		raw, err := client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
			return cassandraClient.ModifyInstanceType(modifyInstanceTypeReq)
		})
		addDebug(modifyInstanceTypeReq.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifyInstanceTypeReq.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cassandraService.CassandraClusterStateRefreshFunc(d.Id(), []string{"ModifyInstanceTypeFailed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("instance_type")
	}
	update = false
	resizeDiskSizeReq := cassandra.CreateResizeDiskSizeRequest()
	resizeDiskSizeReq.ClusterId = d.Id()
	if d.HasChange("disk_size") {
		update = true
	}
	resizeDiskSizeReq.DiskSize = requests.NewInteger(d.Get("disk_size").(int))
	resizeDiskSizeReq.DataCenterId = d.Get("zone_id").(string)
	if update {
		raw, err := client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
			return cassandraClient.ResizeDiskSize(resizeDiskSizeReq)
		})
		addDebug(resizeDiskSizeReq.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), resizeDiskSizeReq.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cassandraService.CassandraClusterStateRefreshFunc(d.Id(), []string{"ResizeDiskSizeFailed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("disk_size")
	}
	update = false
	resizeNodeCountReq := cassandra.CreateResizeNodeCountRequest()
	resizeNodeCountReq.ClusterId = d.Id()
	if d.HasChange("node_count") {
		update = true
	}
	resizeNodeCountReq.NodeCount = requests.NewInteger(d.Get("node_count").(int))
	resizeNodeCountReq.DataCenterId = d.Get("zone_id").(string)
	if update {
		raw, err := client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
			return cassandraClient.ResizeNodeCount(resizeNodeCountReq)
		})
		addDebug(resizeNodeCountReq.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), resizeNodeCountReq.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cassandraService.CassandraClusterStateRefreshFunc(d.Id(), []string{"ResizeNodeCountFailed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("node_count")
	}
	if d.HasChange("password") {
		request := cassandra.CreateModifyAccountPasswordRequest()
		request.ClusterId = d.Id()
		request.Account = "cassandra"
		request.Password = d.Get("password").(string)
		raw, err := client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
			return cassandraClient.ModifyAccountPassword(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("password")
	}
	d.Partial(false)
	return resourceAlicloudCassandraClusterRead(d, meta)
}
func resourceAlicloudCassandraClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := cassandra.CreateDeleteClusterRequest()
	request.ClusterId = d.Id()
	raw, err := client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
		return cassandraClient.DeleteCluster(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		if IsExpectedErrors(err, []string{"Cluster.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	cassandraService := CassandraService{client}
	stateConf := BuildStateConf([]string{"Deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cassandraService.CassandraClusterStateRefreshFunc(d.Id(), []string{"DeleteFailed"}))
	_, err = stateConf.WaitForState()
	return WrapError(err)
}
