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

func resourceAlicloudGaEndpointGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudGaEndpointGroupCreate,
		Read:   resourceAlicloudGaEndpointGroupRead,
		Update: resourceAlicloudGaEndpointGroupUpdate,
		Delete: resourceAlicloudGaEndpointGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(6 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"accelerator_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"endpoint_configurations": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_clientip_preservation": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"endpoint": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"Domain", "Ip", "PublicIp", "ECS", "SLB"}, false),
						},
						"weight": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 255),
						},
					},
				},
			},
			"endpoint_group_region": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"endpoint_group_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"default", "virtual"}, false),
				Default:      "default",
			},
			"endpoint_request_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"HTTP", "HTTPS"}, false),
			},
			"health_check_interval_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"health_check_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"health_check_port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"health_check_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"http", "https", "tcp"}, false),
			},
			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"port_overrides": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint_port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"listener_port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"threshold_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3,
			},
			"traffic_percentage": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudGaEndpointGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	action := "CreateEndpointGroup"
	request := make(map[string]interface{})
	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}
	request["AcceleratorId"] = d.Get("accelerator_id")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	EndpointConfigurations := make([]map[string]interface{}, len(d.Get("endpoint_configurations").([]interface{})))
	for i, EndpointConfigurationsValue := range d.Get("endpoint_configurations").([]interface{}) {
		EndpointConfigurationsMap := EndpointConfigurationsValue.(map[string]interface{})
		EndpointConfigurations[i] = make(map[string]interface{})
		EndpointConfigurations[i]["EnableClientIPPreservation"] = EndpointConfigurationsMap["enable_clientip_preservation"]
		EndpointConfigurations[i]["Endpoint"] = EndpointConfigurationsMap["endpoint"]
		EndpointConfigurations[i]["Type"] = EndpointConfigurationsMap["type"]
		EndpointConfigurations[i]["Weight"] = EndpointConfigurationsMap["weight"]
	}
	request["EndpointConfigurations"] = EndpointConfigurations

	request["EndpointGroupRegion"] = d.Get("endpoint_group_region")
	if v, ok := d.GetOk("endpoint_group_type"); ok {
		request["EndpointGroupType"] = v
	}

	if v, ok := d.GetOk("endpoint_request_protocol"); ok {
		request["EndpointRequestProtocol"] = v
	}

	if v, ok := d.GetOk("health_check_interval_seconds"); ok {
		request["HealthCheckIntervalSeconds"] = v
	}

	if v, ok := d.GetOk("health_check_path"); ok {
		request["HealthCheckPath"] = v
	}

	if v, ok := d.GetOk("health_check_port"); ok {
		request["HealthCheckPort"] = v
	}

	if v, ok := d.GetOk("health_check_protocol"); ok {
		request["HealthCheckProtocol"] = v
	}

	request["ListenerId"] = d.Get("listener_id")
	if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	}

	if v, ok := d.GetOk("port_overrides"); ok {
		PortOverrides := make([]map[string]interface{}, len(v.([]interface{})))
		for i, PortOverridesValue := range v.([]interface{}) {
			PortOverridesMap := PortOverridesValue.(map[string]interface{})
			PortOverrides[i] = make(map[string]interface{})
			PortOverrides[i]["EndpointPort"] = PortOverridesMap["endpoint_port"]
			PortOverrides[i]["ListenerPort"] = PortOverridesMap["listener_port"]
		}
		request["PortOverrides"] = PortOverrides

	}

	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("threshold_count"); ok {
		request["ThresholdCount"] = v
	}

	if v, ok := d.GetOk("traffic_percentage"); ok {
		request["TrafficPercentage"] = v
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("CreateEndpointGroup")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"GA_NOT_STEADY", "StateError.Accelerator", "StateError.EndPointGroup"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ga_endpoint_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["EndpointGroupId"]))
	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, gaService.GaEndpointGroupStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudGaEndpointGroupRead(d, meta)
}
func resourceAlicloudGaEndpointGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	object, err := gaService.DescribeGaEndpointGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ga_endpoint_group gaService.DescribeGaEndpointGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("description", object["Description"])

	endpointConfigurations := make([]map[string]interface{}, 0)
	if endpointConfigurationsList, ok := object["EndpointConfigurations"].([]interface{}); ok {
		for _, v := range endpointConfigurationsList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"enable_clientip_preservation": m1["EnableClientIPPreservation"],
					"endpoint":                     m1["Endpoint"],
					"type":                         m1["Type"],
					"weight":                       m1["Weight"],
				}
				endpointConfigurations = append(endpointConfigurations, temp1)

			}
		}
	}
	if err := d.Set("endpoint_configurations", endpointConfigurations); err != nil {
		return WrapError(err)
	}
	d.Set("endpoint_group_region", object["EndpointGroupRegion"])
	d.Set("health_check_interval_seconds", formatInt(object["HealthCheckIntervalSeconds"]))
	d.Set("health_check_path", object["HealthCheckPath"])
	d.Set("health_check_port", formatInt(object["HealthCheckPort"]))
	d.Set("health_check_protocol", object["HealthCheckProtocol"])
	d.Set("listener_id", object["ListenerId"])
	d.Set("name", object["Name"])

	portOverrides := make([]map[string]interface{}, 0)
	if portOverridesList, ok := object["PortOverrides"].([]interface{}); ok {
		for _, v := range portOverridesList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"endpoint_port": m1["EndpointPort"],
					"listener_port": m1["ListenerPort"],
				}
				portOverrides = append(portOverrides, temp1)

			}
		}
	}
	if err := d.Set("port_overrides", portOverrides); err != nil {
		return WrapError(err)
	}
	d.Set("status", object["State"])
	d.Set("threshold_count", formatInt(object["ThresholdCount"]))
	d.Set("traffic_percentage", formatInt(object["TrafficPercentage"]))
	return nil
}
func resourceAlicloudGaEndpointGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"EndpointGroupId": d.Id(),
	}
	if d.HasChange("endpoint_configurations") {
		update = true
	}
	EndpointConfigurations := make([]map[string]interface{}, len(d.Get("endpoint_configurations").([]interface{})))
	for i, EndpointConfigurationsValue := range d.Get("endpoint_configurations").([]interface{}) {
		EndpointConfigurationsMap := EndpointConfigurationsValue.(map[string]interface{})
		EndpointConfigurations[i] = make(map[string]interface{})
		EndpointConfigurations[i]["EnableClientIPPreservation"] = EndpointConfigurationsMap["enable_clientip_preservation"]
		EndpointConfigurations[i]["Endpoint"] = EndpointConfigurationsMap["endpoint"]
		EndpointConfigurations[i]["Type"] = EndpointConfigurationsMap["type"]
		EndpointConfigurations[i]["Weight"] = EndpointConfigurationsMap["weight"]
	}
	request["EndpointConfigurations"] = EndpointConfigurations

	request["EndpointGroupRegion"] = d.Get("endpoint_group_region")
	request["RegionId"] = client.RegionId
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}
	if d.HasChange("health_check_interval_seconds") {
		update = true
		request["HealthCheckIntervalSeconds"] = d.Get("health_check_interval_seconds")
	}
	if d.HasChange("health_check_path") {
		update = true
		request["HealthCheckPath"] = d.Get("health_check_path")
	}
	if d.HasChange("health_check_port") {
		update = true
		request["HealthCheckPort"] = d.Get("health_check_port")
	}
	if d.HasChange("health_check_protocol") {
		update = true
		request["HealthCheckProtocol"] = d.Get("health_check_protocol")
	}
	if d.HasChange("name") {
		update = true
		request["Name"] = d.Get("name")
	}
	if d.HasChange("port_overrides") {
		update = true
	}
	PortOverrides := make([]map[string]interface{}, len(d.Get("port_overrides").([]interface{})))
	for i, PortOverridesValue := range d.Get("port_overrides").([]interface{}) {
		PortOverridesMap := PortOverridesValue.(map[string]interface{})
		PortOverrides[i] = make(map[string]interface{})
		PortOverrides[i]["EndpointPort"] = PortOverridesMap["endpoint_port"]
		PortOverrides[i]["ListenerPort"] = PortOverridesMap["listener_port"]
	}
	request["PortOverrides"] = PortOverrides

	if d.HasChange("threshold_count") {
		update = true
	}
	request["ThresholdCount"] = d.Get("threshold_count")
	if d.HasChange("traffic_percentage") {
		update = true
		request["TrafficPercentage"] = d.Get("traffic_percentage")
	}
	if update {
		if _, ok := d.GetOk("endpoint_request_protocol"); ok {
			request["EndpointRequestProtocol"] = d.Get("endpoint_request_protocol")
		}
		action := "UpdateEndpointGroup"
		conn, err := client.NewGaplusClient()
		if err != nil {
			return WrapError(err)
		}
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			request["ClientToken"] = buildClientToken("UpdateEndpointGroup")
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.EndPointGroup"}) {
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
		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gaService.GaEndpointGroupStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return resourceAlicloudGaEndpointGroupRead(d, meta)
}
func resourceAlicloudGaEndpointGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	action := "DeleteEndpointGroup"
	var response map[string]interface{}
	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"EndpointGroupId": d.Id(),
	}

	request["AcceleratorId"] = d.Get("accelerator_id")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("DeleteEndpointGroup")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.EndPointGroup"}) {
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gaService.GaEndpointGroupStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
