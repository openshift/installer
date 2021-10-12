package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cassandra"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudCassandraDataCenter() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCassandraDataCenterCreate,
		Read:   resourceAlicloudCassandraDataCenterRead,
		Update: resourceAlicloudCassandraDataCenterUpdate,
		Delete: resourceAlicloudCassandraDataCenterDelete,
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
				ValidateFunc: validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 12}),
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data_center_id": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
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
			"node_count": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"pay_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
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
		},
	}
}

func resourceAlicloudCassandraDataCenterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cassandraService := CassandraService{client}

	request := cassandra.CreateCreateDataCenterRequest()
	if v, ok := d.GetOkExists("auto_renew"); ok {
		request.AutoRenew = requests.NewBoolean(v.(bool))
	}
	if v, ok := d.GetOk("auto_renew_period"); ok {
		request.AutoRenewPeriod = requests.NewInteger(v.(int))
	}
	request.ClusterId = d.Get("cluster_id").(string)
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
	request.NodeCount = requests.NewInteger(d.Get("node_count").(int))
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
		return cassandraClient.CreateDataCenter(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cassandra_data_center", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*cassandra.CreateDataCenterResponse)
	d.SetId(fmt.Sprintf("%v:%v", response.DataCenterId, request.ClusterId))

	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 8*time.Minute, cassandraService.CassandraDataCenterStateRefreshFunc(d.Id(), []string{"CreateFailed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudCassandraDataCenterUpdate(d, meta)
}
func resourceAlicloudCassandraDataCenterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cassandraService := CassandraService{client}
	object, err := cassandraService.DescribeCassandraDataCenter(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cassandra_data_center cassandraService.DescribeCassandraDataCenter Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("id", fmt.Sprintf("%v:%v", object.DataCenterId, object.ClusterId))
	d.Set("cluster_id", object.ClusterId)
	d.Set("data_center_id", object.DataCenterId)
	d.Set("data_center_name", object.DataCenterName)
	d.Set("disk_size", object.DiskSize)
	d.Set("disk_type", object.DiskType)
	d.Set("instance_type", object.InstanceType)
	d.Set("node_count", object.NodeCount)
	d.Set("pay_type", object.PayType)
	d.Set("status", object.Status)
	d.Set("vswitch_id", object.VswitchId)
	d.Set("zone_id", object.ZoneId)
	d.Set("auto_renew", object.AutoRenewal)
	if d.Get("auto_renew").(bool) {
		d.Set("auto_renew_period", object.AutoRenewPeriod)
	} else {
		if v, ok := d.GetOkExists("auto_renew_period"); ok {
			d.Set("auto_renew_period", v.(int))
		}
	}

	d.Set("enable_public", false)
	d.Set("public_points", nil)
	endPoints, err := cassandraService.DescribeCassandraEndpoints(d.Id())
	if err != nil {
		return WrapError(err)
	}
	for _, contactPoint := range endPoints.ContactPoints.ContactPoint {
		if object.ZoneId == contactPoint.DataCenterId && len(contactPoint.PublicAddresses.PublicAddress) > 0 {
			d.Set("enable_public", true)
			d.Set("public_points", contactPoint.PublicAddresses.PublicAddress)
			break
		}
	}
	return nil
}
func resourceAlicloudCassandraDataCenterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cassandraService := CassandraService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Partial(true)
	if (d.IsNewResource() && d.Get("enable_public") != nil && d.Get("enable_public").(bool)) || (!d.IsNewResource() && d.HasChange("enable_public")) {
		if d.Get("enable_public").(bool) {
			request := cassandra.CreateAllocatePublicContactPointsRequest()
			request.ClusterId = parts[1]
			request.DataCenterId = parts[0]
			raw, err := client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
				return cassandraClient.AllocatePublicContactPoints(request)
			})
			addDebug(request.GetActionName(), raw)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cassandraService.CassandraDataCenterStateRefreshFunc(d.Id(), []string{"AllocatePubConnFailed"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		} else {
			request := cassandra.CreateReleasePublicContactPointsRequest()
			request.ClusterId = parts[1]
			request.DataCenterId = parts[0]
			raw, err := client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
				return cassandraClient.ReleasePublicContactPoints(request)
			})
			addDebug(request.GetActionName(), raw)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cassandraService.CassandraDataCenterStateRefreshFunc(d.Id(), []string{"ReleasePubConnFailed"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
		d.SetPartial("enable_public")
	}
	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlicloudCassandraDataCenterRead(d, meta)
	}

	if d.HasChange("data_center_name") {
		request := cassandra.CreateModifyDataCenterRequest()
		request.ClusterId = parts[1]
		request.DataCenterId = parts[0]
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
	if d.HasChange("instance_type") {
		request := cassandra.CreateModifyInstanceTypeRequest()
		request.ClusterId = parts[1]
		request.DataCenterId = parts[0]
		request.InstanceType = d.Get("instance_type").(string)
		raw, err := client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
			return cassandraClient.ModifyInstanceType(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cassandraService.CassandraDataCenterStateRefreshFunc(d.Id(), []string{"ModifyInstanceTypeFailed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("instance_type")
	}
	if d.HasChange("disk_size") {
		request := cassandra.CreateResizeDiskSizeRequest()
		request.ClusterId = parts[1]
		request.DataCenterId = parts[0]
		request.DiskSize = requests.NewInteger(d.Get("disk_size").(int))
		raw, err := client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
			return cassandraClient.ResizeDiskSize(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cassandraService.CassandraDataCenterStateRefreshFunc(d.Id(), []string{"ResizeDiskSizeFailed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("disk_size")
	}
	if d.HasChange("node_count") {
		request := cassandra.CreateResizeNodeCountRequest()
		request.ClusterId = parts[1]
		request.DataCenterId = parts[0]
		request.NodeCount = requests.NewInteger(d.Get("node_count").(int))
		raw, err := client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
			return cassandraClient.ResizeNodeCount(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cassandraService.CassandraDataCenterStateRefreshFunc(d.Id(), []string{"ResizeNodeCountFailed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("node_count")
	}
	d.Partial(false)
	return resourceAlicloudCassandraDataCenterRead(d, meta)
}
func resourceAlicloudCassandraDataCenterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := cassandra.CreateDeleteDataCenterRequest()
	request.ClusterId = parts[1]
	request.DataCenterId = parts[0]
	raw, err := client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
		return cassandraClient.DeleteDataCenter(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		if IsExpectedErrors(err, []string{"DataCenter.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	cassandraService := CassandraService{client}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cassandraService.CassandraClusterStateRefreshFunc(parts[1], []string{"DeleteOtherDataCenterFailed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, parts[1])
	}
	return nil
}
