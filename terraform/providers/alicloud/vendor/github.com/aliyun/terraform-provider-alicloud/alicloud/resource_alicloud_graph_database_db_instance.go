package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudGraphDatabaseDbInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudGraphDatabaseDbInstanceCreate,
		Read:   resourceAlicloudGraphDatabaseDbInstanceRead,
		Update: resourceAlicloudGraphDatabaseDbInstanceUpdate,
		Delete: resourceAlicloudGraphDatabaseDbInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"db_instance_ip_array": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_instance_ip_array_attribute": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"db_instance_ip_array_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"security_ips": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"db_instance_category": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"HA"}, false),
			},
			"db_instance_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_instance_network_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"vpc"}, false),
			},
			"db_instance_storage_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"cloud_essd", "cloud_ssd"}, false),
			},
			"db_node_class": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"gdb.r.xlarge", "gdb.r.2xlarge", "gdb.r.4xlarge", "gdb.r.8xlarge", "gdb.r.16xlarge"}, false),
			},
			"db_node_storage": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(20, 32000),
			},
			"db_version": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"1.0", "1.0-OpenCypher"}, false),
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudGraphDatabaseDbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateDBInstance"
	request := make(map[string]interface{})
	conn, err := client.NewGdsClient()
	if err != nil {
		return WrapError(err)
	}
	request["DBInstanceCategory"] = strings.ToLower(d.Get("db_instance_category").(string))
	if v, ok := d.GetOk("db_instance_description"); ok {
		request["DBInstanceDescription"] = v
	}
	request["DBInstanceNetworkType"] = d.Get("db_instance_network_type")
	request["DBNodeStorageType"] = d.Get("db_instance_storage_type")
	request["DBInstanceClass"] = d.Get("db_node_class")
	request["DBNodeStorage"] = d.Get("db_node_storage")
	request["DBInstanceVersion"] = d.Get("db_version")
	request["PayType"] = convertGraphDatabaseDbInstancePaymentTypeRequest(d.Get("payment_type").(string))
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("CreateDBInstance")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-03"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_graph_database_db_instance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DBInstanceId"]))
	gdbService := GdbService{client}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, gdbService.GraphDatabaseDbInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudGraphDatabaseDbInstanceUpdate(d, meta)
}
func resourceAlicloudGraphDatabaseDbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gdbService := GdbService{client}
	object, err := gdbService.DescribeGraphDatabaseDbInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_graph_database_db_instance gdbService.DescribeGraphDatabaseDbInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("db_instance_category", object["Category"])
	d.Set("db_instance_description", object["DBInstanceDescription"])
	d.Set("db_instance_network_type", object["DBInstanceNetworkType"])
	d.Set("db_instance_storage_type", object["DBInstanceStorageType"])
	d.Set("db_node_class", object["DBNodeClass"])
	d.Set("db_node_storage", formatInt(object["DBNodeStorage"]))
	d.Set("db_version", object["DBVersion"])
	d.Set("payment_type", convertGraphDatabaseDbInstancePaymentTypeResponse(object["PayType"]))
	d.Set("status", object["DBInstanceStatus"])
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("zone_id", object["ZoneId"])
	if DBInstanceIPArray, ok := object["DBInstanceIPArray"]; ok {
		DBInstanceIPArrayAry, ok := DBInstanceIPArray.([]interface{})
		if ok && len(DBInstanceIPArrayAry) > 0 {
			DBInstanceIPArraySli := make([]map[string]interface{}, 0)
			for _, DBInstanceIPArrayArg := range DBInstanceIPArrayAry {
				DBInstanceIPArrayMap := make(map[string]interface{})
				DBInstanceIPArrayMap["security_ips"] = DBInstanceIPArrayArg.(map[string]interface{})["SecurityIps"]
				DBInstanceIPArrayMap["db_instance_ip_array_name"] = DBInstanceIPArrayArg.(map[string]interface{})["DBInstanceIPArrayName"]
				if v, ok := DBInstanceIPArrayArg.(map[string]interface{})["DBInstanceIPArrayAttribute"]; ok {
					DBInstanceIPArrayMap["db_instance_ip_array_attribute"] = v
				}
				DBInstanceIPArraySli = append(DBInstanceIPArraySli, DBInstanceIPArrayMap)
			}
			d.Set("db_instance_ip_array", DBInstanceIPArraySli)
		}
	}
	return nil
}
func resourceAlicloudGraphDatabaseDbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gdbService := GdbService{client}
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("db_instance_description") {
		update = true
	}
	if v, ok := d.GetOk("db_instance_description"); ok {
		request["DBInstanceDescription"] = v
	}
	if update {
		action := "ModifyDBInstanceDescription"
		conn, err := client.NewGdsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gdbService.GraphDatabaseDbInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("db_instance_description")
	}
	update = false
	modifyDBInstanceAccessWhiteListReq := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	if d.HasChange("db_instance_ip_array") {
		oraw, nraw := d.GetChange("db_instance_ip_array")
		o := oraw.(*schema.Set)
		n := nraw.(*schema.Set)
		remove := o.Difference(n).List()
		create := n.Difference(o).List()

		if len(remove) > 0 {
			for _, dBInstanceIPArray := range remove {
				dBInstanceIPArrayArg := dBInstanceIPArray.(map[string]interface{})

				action := "ModifyDBInstanceAccessWhiteList"
				conn, err := client.NewGdsClient()
				if err != nil {
					return WrapError(err)
				}
				if v, ok := dBInstanceIPArrayArg["db_instance_ip_array_name"]; !ok || v.(string) == "default" {
					continue
				}
				modifyDBInstanceAccessWhiteListReq["DBInstanceIPArrayName"] = dBInstanceIPArrayArg["db_instance_ip_array_name"]
				modifyDBInstanceAccessWhiteListReq["SecurityIps"] = "Empty"
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-03"), StringPointer("AK"), nil, modifyDBInstanceAccessWhiteListReq, &util.RuntimeOptions{})
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, modifyDBInstanceAccessWhiteListReq)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}
		}

		if len(create) > 0 {
			for _, dBInstanceIPArray := range create {
				dBInstanceIPArrayArg := dBInstanceIPArray.(map[string]interface{})

				action := "ModifyDBInstanceAccessWhiteList"
				conn, err := client.NewGdsClient()
				if err != nil {
					return WrapError(err)
				}
				modifyDBInstanceAccessWhiteListReq["DBInstanceIPArrayAttribute"] = dBInstanceIPArrayArg["db_instance_ip_array_attribute"]
				modifyDBInstanceAccessWhiteListReq["DBInstanceIPArrayName"] = dBInstanceIPArrayArg["db_instance_ip_array_name"]
				modifyDBInstanceAccessWhiteListReq["SecurityIps"] = dBInstanceIPArrayArg["security_ips"]
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-03"), StringPointer("AK"), nil, modifyDBInstanceAccessWhiteListReq, &util.RuntimeOptions{})
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, modifyDBInstanceAccessWhiteListReq)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}
		}

		d.SetPartial("db_instance_ip_array")
	}

	modifyDBInstanceSpecReq := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}

	modifyDBInstanceSpecReq["DBInstanceClass"] = d.Get("db_node_class")
	if !d.IsNewResource() && d.HasChange("db_node_class") {
		update = true
	}
	modifyDBInstanceSpecReq["DBInstanceStorage"] = d.Get("db_node_storage")
	if !d.IsNewResource() && d.HasChange("db_node_storage") {
		update = true
	}
	modifyDBInstanceSpecReq["RegionId"] = client.RegionId
	if update {
		modifyDBInstanceSpecReq["DBInstanceStorageType"] = d.Get("db_instance_storage_type")
		action := "ModifyDBInstanceSpec"
		conn, err := client.NewGdsClient()
		if err != nil {
			return WrapError(err)
		}
		modifyDBInstanceSpecReq["ClientToken"] = buildClientToken("ModifyDBInstanceSpec")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-03"), StringPointer("AK"), nil, modifyDBInstanceSpecReq, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyDBInstanceSpecReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gdbService.GraphDatabaseDbInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("db_instance_storage_type")
		d.SetPartial("db_node_class")
		d.SetPartial("db_node_storage")
	}
	d.Partial(false)
	return resourceAlicloudGraphDatabaseDbInstanceRead(d, meta)
}
func resourceAlicloudGraphDatabaseDbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gdbService := GdbService{client}
	action := "DeleteDBInstance"
	var response map[string]interface{}
	conn, err := client.NewGdsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"InvalidDBInstance.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gdbService.GraphDatabaseDbInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
func convertGraphDatabaseDbInstancePaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "Postpaid"
	}
	return source
}
func convertGraphDatabaseDbInstancePaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "Postpaid":
		return "PayAsYouGo"
	}
	return source
}
