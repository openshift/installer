package alicloud

import (
	"encoding/json"
	"regexp"
	"strconv"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const dbConnectionPrefixWithSuffixRegex = "^([a-zA-Z0-9\\-_]+)" + dbConnectionSuffixRegex + "$"

var dbConnectionPrefixWithSuffixRegexp = regexp.MustCompile(dbConnectionPrefixWithSuffixRegex)

func resourceAlicloudDBReadWriteSplittingConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDBReadWriteSplittingConnectionCreate,
		Read:   resourceAlicloudDBReadWriteSplittingConnectionRead,
		Update: resourceAlicloudDBReadWriteSplittingConnectionUpdate,
		Delete: resourceAlicloudDBReadWriteSplittingConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"connection_prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 31),
			},
			"distribution_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Standard", "Custom"}, false),
			},
			"weight": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"max_delay_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudDBReadWriteSplittingConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	action := "AllocateReadWriteSplittingConnection"
	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"DBInstanceId": Trim(d.Get("instance_id").(string)),
		"MaxDelayTime": strconv.Itoa(d.Get("max_delay_time").(int)),
		"SourceIp":     client.SourceIp,
	}
	prefix, ok := d.GetOk("connection_prefix")
	if ok && prefix.(string) != "" {
		request["ConnectionStringPrefix"] = prefix
	}

	port, ok := d.GetOk("port")
	if ok {
		request["Port"] = strconv.Itoa(port.(int))
	}

	request["DistributionType"] = d.Get("distribution_type")

	if weight, ok := d.GetOk("weight"); ok && weight != nil && len(weight.(map[string]interface{})) > 0 {
		if serial, err := json.Marshal(weight); err != nil {
			return WrapError(err)
		} else {
			request["Weight"] = string(serial)
		}
	}
	runtime := util.RuntimeOptions{}
	conn, err := client.NewRdsClient()
	if err != nil {
		return WrapError(err)
	}
	if err := resource.Retry(60*time.Minute, func() *resource.RetryError {
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, DBReadInstanceNotReadyStatus) || NeedRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(request["DBInstanceId"].(string))

	// wait read write splitting connection ready after creation
	// for it may take up to 10 hours to create a readonly instance
	if err := rdsService.WaitForDBReadWriteSplitting(request["DBInstanceId"].(string), "", 60*60*10); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudDBReadWriteSplittingConnectionUpdate(d, meta)
}

func resourceAlicloudDBReadWriteSplittingConnectionRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	proxy, proxyErr := rdsService.DescribeDBProxy(d.Id())
	if proxyErr != nil {
		if NotFoundError(proxyErr) {
			d.SetId("")
			return nil
		}
		return WrapError(proxyErr)
	}
	if proxy["DBProxyInstanceType"] == "2" {
		return resourceAlicloudDBProxyEndpointRead(d, rdsService, proxy["DBProxyInstanceName"].(string))
	}

	err := rdsService.WaitForDBReadWriteSplitting(d.Id(), "", DefaultLongTimeout)
	if err != nil {
		return WrapError(err)
	}

	object, err := rdsService.DescribeDBReadWriteSplittingConnection(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", d.Id())
	d.Set("connection_string", object["ConnectionString"])
	d.Set("distribution_type", object["DistributionType"])
	if port, err := strconv.Atoi(object["Port"].(string)); err == nil {
		d.Set("port", port)
	}
	if mdt, err := strconv.Atoi(object["MaxDelayTime"].(string)); err == nil {
		d.Set("max_delay_time", mdt)
	}
	if w, ok := d.GetOk("weight"); ok {
		documented := w.(map[string]interface{})
		dBInstanceWeights := object["DBInstanceWeights"].(map[string]interface{})["DBInstanceWeight"].([]interface{})
		for _, config := range dBInstanceWeights {
			config := config.(map[string]interface{})
			if config["Availability"] != "Available" {
				delete(documented, config["DBInstanceId"].(string))
				continue
			}
			if config["Weight"] != "0" {
				if _, ok := documented[config["DBInstanceId"].(string)]; ok {
					documented[config["DBInstanceId"].(string)] = config["Weight"]
				}
			}
		}
		d.Set("weight", documented)
	}
	submatch := dbConnectionPrefixWithSuffixRegexp.FindStringSubmatch(object["ConnectionString"].(string))
	if len(submatch) > 1 {
		d.Set("connection_prefix", submatch[1])
	}

	return nil
}

func resourceAlicloudDBReadWriteSplittingConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	action := "ModifyReadWriteSplittingConnection"
	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"DBInstanceId": d.Id(),
		"SourceIp":     client.SourceIp,
	}
	update := false

	if d.HasChange("max_delay_time") {
		request["MaxDelayTime"] = strconv.Itoa(d.Get("max_delay_time").(int))
		update = true
	}

	if !update && d.IsNewResource() {
		return resourceAlicloudDBReadWriteSplittingConnectionRead(d, meta)
	}

	if d.HasChange("weight") {
		if weight, ok := d.GetOk("weight"); ok && weight != nil && len(weight.(map[string]interface{})) > 0 {
			if serial, err := json.Marshal(weight); err != nil {
				return err
			} else {
				request["Weight"] = string(serial)
			}
		}
		update = true
	}

	if d.HasChange("distribution_type") {
		request["DistributionType"] = d.Get("distribution_type")
		update = true
	}

	if update {
		// wait instance running before modifying
		if err := rdsService.WaitForDBInstance(request["DBInstanceId"].(string), Running, 60*60); err != nil {
			return WrapError(err)
		}
		conn, err := client.NewRdsClient()
		if err != nil {
			return WrapError(err)
		}
		runtime := util.RuntimeOptions{}
		if err := resource.Retry(30*time.Minute, func() *resource.RetryError {
			response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, OperationDeniedDBStatus) || IsExpectedErrors(err, DBReadInstanceNotReadyStatus) || NeedRetry(err) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		// wait instance running after modifying
		if err := rdsService.WaitForDBInstance(request["DBInstanceId"].(string), Running, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
	}

	return resourceAlicloudDBReadWriteSplittingConnectionRead(d, meta)
}

func resourceAlicloudDBReadWriteSplittingConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	action := "ReleaseReadWriteSplittingConnection"
	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"DBInstanceId": d.Id(),
		"SourceIp":     client.SourceIp,
	}
	conn, err := client.NewRdsClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	if err := resource.Retry(30*time.Minute, func() *resource.RetryError {
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, OperationDeniedDBStatus) || NeedRetry(err) {
				return resource.RetryableError(err)
			}
			if NotFoundError(err) || IsExpectedErrors(err, []string{"InvalidRwSplitNetType.NotFound"}) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return WrapError(rdsService.WaitForDBReadWriteSplitting(d.Id(), Deleted, DefaultLongTimeout))
}

func resourceAlicloudDBProxyEndpointRead(d *schema.ResourceData, rdsService RdsService, endPointName string) error {
	endpointInfo, endpointError := rdsService.DescribeDBProxyEndpoint(d.Id(), endPointName)
	if endpointError != nil {
		if NotFoundError(endpointError) {
			d.SetId("")
			return nil
		}
		return WrapError(endpointError)
	}
	d.Set("instance_id", d.Id())
	d.Set("connection_string", endpointInfo["DBProxyConnectString"])
	d.Set("distribution_type", endpointInfo["ReadOnlyInstanceDistributionType"])
	if port, err := strconv.Atoi(endpointInfo["DBProxyConnectStringPort"].(string)); err == nil {
		d.Set("port", port)
	}

	if mdt, err := strconv.Atoi(endpointInfo["ReadOnlyInstanceMaxDelayTime"].(string)); err == nil {
		d.Set("max_delay_time", mdt)
	}
	submatch := dbConnectionPrefixWithSuffixRegexp.FindStringSubmatch(endpointInfo["DBProxyConnectString"].(string))
	if len(submatch) > 1 {
		d.Set("connection_prefix", submatch[1])
	}

	var documented map[string]interface{}
	if w, ok := d.GetOk("weight"); ok {
		documented = w.(map[string]interface{})
	} else {
		documented = make(map[string]interface{})
	}
	var weight []map[string]interface{}
	rawData := []byte(endpointInfo["ReadOnlyInstanceWeight"].(string))
	parseErr := json.Unmarshal(rawData, &weight)
	if parseErr != nil {
		return WrapError(parseErr)
	}
	for _, configNode := range weight {
		var dbInstanceId string
		if instanceId, ok := configNode["DBInstanceId"]; ok {
			dbInstanceId = instanceId.(string)
		}
		if _, ok := configNode["Availability"]; ok && configNode["Availability"] != "Available" {
			delete(documented, dbInstanceId)
			continue
		}
		if _, ok := configNode["Weight"]; ok && configNode["Weight"] != "0" {
			if _, ok := documented[dbInstanceId]; ok {
				documented[dbInstanceId] = configNode["Weight"]
			}
		}
	}
	d.Set("weight", documented)
	return nil
}
