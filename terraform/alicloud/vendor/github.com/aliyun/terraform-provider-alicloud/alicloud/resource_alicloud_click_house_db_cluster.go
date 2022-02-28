package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudClickHouseDbCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudClickHouseDbClusterCreate,
		Read:   resourceAlicloudClickHouseDbClusterRead,
		Update: resourceAlicloudClickHouseDbClusterUpdate,
		Delete: resourceAlicloudClickHouseDbClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"category": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Basic", "HighAvailability"}, false),
				ForceNew:     true,
			},
			"db_cluster_access_white_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_cluster_ip_array_attribute": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"db_cluster_ip_array_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"security_ip_list": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"db_cluster_class": {
				Type:         schema.TypeString,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"S4-NEW", "S8", "S16", "S32", "S64", "S104", "C4-NEW", "C8", "C16", "C32", "C64", "C104"}, false),
				Required:     true,
			},
			"db_cluster_network_type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"vpc"}, false),
				Required:     true,
				ForceNew:     true,
			},
			"db_cluster_version": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"19.15.2.2", "20.3.10.75", "20.8.7.15"}, false),
			},
			"db_node_storage": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"db_node_group_count": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(1, 48),
				Required:     true,
				ForceNew:     true,
			},

			"encryption_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"encryption_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
				ForceNew:     true,
			},
			"period": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Month", "Year"}, false),
			},
			"storage_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"cloud_essd", "cloud_efficiency", "cloud_essd_pl2", "cloud_essd_pl3"}, false),
			},
			"used_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"db_cluster_description": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Creating", "Deleting", "Restarting", "Preparing", "Running"}, false),
			},
			"maintain_time": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudClickHouseDbClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	clickhouseService := ClickhouseService{client}
	var response map[string]interface{}
	action := "CreateDBInstance"
	request := make(map[string]interface{})
	conn, err := client.NewClickhouseClient()
	if err != nil {
		return WrapError(err)
	}
	request["DBClusterCategory"] = d.Get("category")
	request["DBClusterClass"] = d.Get("db_cluster_class")
	if v, ok := d.GetOk("db_cluster_description"); ok {
		request["DBClusterDescription"] = v
	}
	request["DBClusterNetworkType"] = d.Get("db_cluster_network_type")
	request["DBClusterVersion"] = d.Get("db_cluster_version")
	request["DBNodeGroupCount"] = d.Get("db_node_group_count")
	request["DBNodeStorage"] = d.Get("db_node_storage")
	if v, ok := d.GetOk("encryption_key"); ok {
		request["EncryptionKey"] = v
	}
	if v, ok := d.GetOk("encryption_type"); ok {
		request["EncryptionType"] = v
	}
	request["PayType"] = convertClickHouseDbClusterPaymentTypeRequest(d.Get("payment_type").(string))

	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	request["RegionId"] = client.RegionId
	request["DbNodeStorageType"] = d.Get("storage_type")
	if v, ok := d.GetOk("used_time"); ok {
		request["UsedTime"] = v
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
	request["ClientToken"] = buildClientToken("CreateDBInstance")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-11"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_click_house_db_cluster", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DBClusterId"]))
	stateConf := BuildStateConf([]string{"Creating", "Deleting", "Restarting", "Preparing"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, clickhouseService.ClickHouseDbClusterStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudClickHouseDbClusterUpdate(d, meta)
}
func resourceAlicloudClickHouseDbClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	clickhouseService := ClickhouseService{client}
	object, err := clickhouseService.DescribeClickHouseDbCluster(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_click_house_db_cluster clickhouseService.DescribeClickHouseDbCluster Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("category", object["Category"])
	d.Set("db_cluster_description", object["DBClusterDescription"])
	d.Set("db_cluster_network_type", object["DBClusterNetworkType"])
	d.Set("db_node_storage", fmt.Sprint(formatInt(object["DBNodeStorage"])))
	d.Set("encryption_key", object["EncryptionKey"])
	d.Set("encryption_type", object["EncryptionType"])
	d.Set("maintain_time", object["MaintainTime"])
	d.Set("status", object["DBClusterStatus"])
	d.Set("payment_type", convertClickHouseDbClusterPaymentTypeResponse(object["PayType"].(string)))
	d.Set("storage_type", convertClickHouseDbClusterStorageTypeResponse(object["StorageType"].(string)))
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("zone_id", object["ZoneId"])
	describeDBClusterAccessWhiteListObject, err := clickhouseService.DescribeDBClusterAccessWhiteList(d.Id())
	if err != nil {
		return WrapError(err)
	}
	if dBClusterAccessWhiteListMap, ok := describeDBClusterAccessWhiteListObject["DBClusterAccessWhiteList"].(map[string]interface{}); ok && dBClusterAccessWhiteListMap != nil {
		if iPArrayList, ok := dBClusterAccessWhiteListMap["IPArray"]; ok && iPArrayList != nil {
			dBClusterAccessWhiteListMaps := make([]map[string]interface{}, 0)
			for _, iPArrayListItem := range iPArrayList.([]interface{}) {
				if v, ok := iPArrayListItem.(map[string]interface{}); ok {
					if v["DBClusterIPArrayName"].(string) == "default" {
						continue
					}
					iPArrayListItemMap := make(map[string]interface{})
					iPArrayListItemMap["db_cluster_ip_array_attribute"] = v["DBClusterIPArrayAttribute"]
					iPArrayListItemMap["db_cluster_ip_array_name"] = v["DBClusterIPArrayName"]
					iPArrayListItemMap["security_ip_list"] = v["SecurityIPList"]
					dBClusterAccessWhiteListMaps = append(dBClusterAccessWhiteListMaps, iPArrayListItemMap)
				}
			}
			d.Set("db_cluster_access_white_list", dBClusterAccessWhiteListMaps)
		}
	}

	return nil
}
func resourceAlicloudClickHouseDbClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"DBClusterId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("db_cluster_description") {
		update = true
		if v, ok := d.GetOk("db_cluster_description"); ok {
			request["DBClusterDescription"] = v
		}
	}
	if update {
		action := "ModifyDBClusterDescription"
		conn, err := client.NewClickhouseClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-11"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectDBInstanceState"}) || NeedRetry(err) {
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
		d.SetPartial("db_cluster_description")
	}
	update = false
	modifyDBClusterMaintainTimeReq := map[string]interface{}{
		"DBClusterId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("maintain_time") {
		update = true
	}
	if v, ok := d.GetOk("maintain_time"); ok {
		modifyDBClusterMaintainTimeReq["MaintainTime"] = v
	}
	if update {
		action := "ModifyDBClusterMaintainTime"
		conn, err := client.NewClickhouseClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-11"), StringPointer("AK"), nil, modifyDBClusterMaintainTimeReq, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectDBInstanceState"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyDBClusterMaintainTimeReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("maintain_time")
	}
	if d.HasChange("db_cluster_access_white_list") {

		oraw, nraw := d.GetChange("db_cluster_access_white_list")
		remove := oraw.(*schema.Set).Difference(nraw.(*schema.Set)).List()
		create := nraw.(*schema.Set).Difference(oraw.(*schema.Set)).List()
		if len(remove) > 0 {
			removeWhiteListReq := map[string]interface{}{
				"DBClusterId": d.Id(),
				"ModifyMode":  "Delete",
			}

			for _, whiteList := range remove {
				whiteListArg := whiteList.(map[string]interface{})
				removeWhiteListReq["DBClusterIPArrayAttribute"] = whiteListArg["db_cluster_ip_array_attribute"]
				removeWhiteListReq["DBClusterIPArrayName"] = whiteListArg["db_cluster_ip_array_name"]
				removeWhiteListReq["SecurityIps"] = whiteListArg["security_ip_list"]

				action := "ModifyDBClusterAccessWhiteList"
				conn, err := client.NewClickhouseClient()
				if err != nil {
					return WrapError(err)
				}
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-11"), StringPointer("AK"), nil, removeWhiteListReq, &util.RuntimeOptions{})
					if err != nil {
						if IsExpectedErrors(err, []string{"IncorrectDBInstanceState"}) || NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, removeWhiteListReq)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}
		}
		if len(create) > 0 {
			createWhiteListReq := map[string]interface{}{
				"DBClusterId": d.Id(),
				"ModifyMode":  "Append",
			}

			for _, whiteList := range create {
				whiteListArg := whiteList.(map[string]interface{})
				createWhiteListReq["DBClusterIPArrayAttribute"] = whiteListArg["db_cluster_ip_array_attribute"]
				createWhiteListReq["DBClusterIPArrayName"] = whiteListArg["db_cluster_ip_array_name"]
				createWhiteListReq["SecurityIps"] = whiteListArg["security_ip_list"]

				action := "ModifyDBClusterAccessWhiteList"
				conn, err := client.NewClickhouseClient()
				if err != nil {
					return WrapError(err)
				}
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-11"), StringPointer("AK"), nil, createWhiteListReq, &util.RuntimeOptions{})
					if err != nil {
						if IsExpectedErrors(err, []string{"IncorrectDBInstanceState"}) || NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, createWhiteListReq)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}
		}
		d.SetPartial("db_cluster_access_white_list")
	}
	if d.HasChange("status") {
		clickhouseService := ClickhouseService{client}
		object, err := clickhouseService.DescribeClickHouseDbCluster(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := d.Get("status").(string)
		if object["DBClusterStatus"].(string) != target {
			if target == "Running" {
				request := map[string]interface{}{
					"DBClusterId": d.Id(),
				}
				request["RegionId"] = client.RegionId
				action := "RestartInstance"
				conn, err := client.NewClickhouseClient()
				if err != nil {
					return WrapError(err)
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-11"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
					if err != nil {
						if IsExpectedErrors(err, []string{"IncorrectDBInstanceState"}) || NeedRetry(err) {
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
			}
			d.SetPartial("status")
		}
		stateConf := BuildStateConf([]string{"RESTARTING"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, clickhouseService.ClickHouseDbClusterStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	d.Partial(false)
	return resourceAlicloudClickHouseDbClusterRead(d, meta)
}
func resourceAlicloudClickHouseDbClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDBCluster"
	var response map[string]interface{}
	conn, err := client.NewClickhouseClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DBClusterId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-11"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectDBInstanceState"}) || NeedRetry(err) {
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

func convertClickHouseDbClusterPaymentTypeRequest(source string) string {
	switch source {
	case "PayAsYouGo":
		return "Postpaid"
	case "Subscription":
		return "Prepay"
	}
	return source
}
func convertClickHouseDbClusterPaymentTypeResponse(source string) string {
	switch source {
	case "Postpaid":
		return "PayAsYouGo"
	case "Prepay":
		return "Subscription"
	}
	return source
}

func convertClickHouseDbClusterStorageTypeResponse(source string) string {
	switch source {
	case "CloudESSD":
		return "cloud_essd"
	case "CloudEfficiency":
		return "cloud_efficiency"
	case "CloudESSD_PL2":
		return "cloud_essd_pl2"
	case "CloudESSD_PL3":
		return "cloud_essd_pl3"

	}
	return source
}
