package alicloud

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudGpdbElasticInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudGpdbElasticInstanceCreate,
		Read:   resourceAlicloudGpdbElasticInstanceRead,
		Update: resourceAlicloudGpdbElasticInstanceUpdate,
		Delete: resourceAlicloudGpdbElasticInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(50 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"engine": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"gpdb"}, false),
			},
			"engine_version": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"6.0"}, false),
			},
			"seg_storage_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"cloud_essd", "cloud_efficiency"}, false),
			},
			"seg_node_num": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 4 || v > 256 || v%4 != 0 {
						errs = append(errs, fmt.Errorf("%q must be between 0 and 256 inclusive, and multiple of 4, got: %d", key, v))
					}
					return
				},
			},
			"storage_size": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 50 || v > 4000 || v%50 != 0 {
						errs = append(errs, fmt.Errorf("%q must be between 50 and 4000 inclusive, and multiple of 50, got: %d", key, v))
					}
					return
				},
			},
			"instance_spec": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"2C16G", "4C32G", "16C128G"}, false),
			},
			"db_instance_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"VPC"}, false),
				Default:      "VPC",
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
				Default:      "PayAsYouGo",
			},
			"payment_duration_unit": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringInSlice([]string{"Month", "Year"}, false),
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"payment_duration": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateFunc:     validation.IntBetween(1, 12),
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
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
			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_ip_list": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudGpdbElasticInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbService := GpdbService{client}
	var response map[string]interface{}
	action := "CreateECSDBInstance"
	request := make(map[string]interface{})
	conn, err := client.NewGpdbClient()
	if err != nil {
		return WrapError(err)
	}

	request["Engine"] = d.Get("engine")
	request["EngineVersion"] = d.Get("engine_version")
	request["SegStorageType"] = d.Get("seg_storage_type")
	request["SegNodeNum"] = d.Get("seg_node_num")
	request["StorageSize"] = d.Get("storage_size")
	request["InstanceSpec"] = d.Get("instance_spec")
	request["PayType"] = convertGpdbInstancePaymentTypeRequest(d.Get("payment_type").(string))
	if request["PayType"].(string) == "Prepaid" {
		request["Period"] = d.Get("payment_duration_unit")
		paymentDuration := d.Get("payment_duration").(int)
		request["UsedTime"] = strconv.Itoa(paymentDuration)
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}

	vswitchId := Trim(d.Get("vswitch_id").(string))
	if vswitchId != "" {
		vpcService := VpcService{client}
		vsw, err := vpcService.DescribeVSwitchWithTeadsl(vswitchId)
		if err != nil {
			return WrapError(err)
		}
		request["VPCId"] = vsw["VpcId"]
		request["VSwitchId"] = vswitchId
		if v, ok := request["ZoneId"].(string); !ok || v == "" {
			request["ZoneId"] = vsw["ZoneId"]
		}
	}
	request["ClientToken"] = buildClientToken("CreateECSDBInstance")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_gpdb_elastic_instance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DBInstanceId"]))
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 10*time.Minute, gpdbService.GpdbElasticInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudGpdbElasticInstanceUpdate(d, meta)
}

func resourceAlicloudGpdbElasticInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbService := GpdbService{client}
	instance, err := gpdbService.DescribeGpdbElasticInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("engine", instance["Engine"])
	d.Set("engine_version", instance["EngineVersion"])
	d.Set("seg_storage_type", instance["StorageType"])
	d.Set("seg_node_num", instance["SegNodeNum"])
	d.Set("storage_size", instance["StorageSize"])
	d.Set("payment_type", convertGpdbInstancePaymentTypeResponse(instance["PayType"].(string)))
	d.Set("instance_spec", convertDBInstanceClassToInstanceSpec(instance["DBInstanceClass"].(string)))
	d.Set("status", instance["DBInstanceStatus"])
	d.Set("db_instance_description", instance["DBInstanceDescription"])
	d.Set("instance_network_type", instance["InstanceNetworkType"])
	d.Set("vswitch_id", instance["VSwitchId"])
	d.Set("zone_id", instance["ZoneId"])
	d.Set("connection_string", instance["ConnectionString"])
	securityIps, err := gpdbService.DescribeGpdbSecurityIps(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("security_ip_list", securityIps)
	return nil
}

func resourceAlicloudGpdbElasticInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbService := GpdbService{client}
	d.Partial(true)
	conn, err := client.NewGpdbClient()
	if err != nil {
		return WrapError(err)
	}
	if d.HasChange("db_instance_description") {
		action := "ModifyDBInstanceDescription"
		request := map[string]interface{}{
			"RegionId":              client.RegionId,
			"DBInstanceId":          d.Id(),
			"DBInstanceDescription": d.Get("db_instance_description"),
			"SourceIp":              client.SourceIp,
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("db_instance_description")
	}
	if d.HasChange("security_ip_list") {
		ipList := expandStringList(d.Get("security_ip_list").([]interface{}))
		ipStr := strings.Join(ipList[:], COMMA_SEPARATED)
		if ipStr == "" {
			ipStr = LOCAL_HOST_IP
		}
		if err := gpdbService.ModifyGpdbSecurityIps(d.Id(), ipStr); err != nil {
			return WrapError(err)
		}
		d.SetPartial("security_ip_list")
	}

	d.Partial(false)
	return resourceAlicloudGpdbElasticInstanceRead(d, meta)
}

func resourceAlicloudGpdbElasticInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDBInstance"
	var response map[string]interface{}
	conn, err := client.NewGpdbClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"OperationDenied.DBInstancePayType"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

func convertGpdbInstancePaymentTypeRequest(source string) string {
	switch source {
	case "PayAsYouGo":
		return "Postpaid"
	case "Subscription":
		return "Prepaid"
	}
	return source
}
func convertGpdbInstancePaymentTypeResponse(source string) string {
	switch source {
	case "Postpaid":
		return "PayAsYouGo"
	case "Prepaid":
		return "Subscription"
	}
	return source
}

func convertDBInstanceClassToInstanceSpec(instanceClass string) string {
	splitClass := strings.Split(instanceClass, ".")
	return strings.ToUpper(splitClass[len(splitClass)-1])
}
