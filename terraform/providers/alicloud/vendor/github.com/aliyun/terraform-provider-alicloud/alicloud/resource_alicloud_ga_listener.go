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

func resourceAlicloudGaListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudGaListenerCreate,
		Read:   resourceAlicloudGaListenerRead,
		Update: resourceAlicloudGaListenerUpdate,
		Delete: resourceAlicloudGaListenerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(6 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"accelerator_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"certificates": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"client_affinity": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"NONE", "SOURCE_IP"}, false),
				Default:      "NONE",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if d.Get("protocol") == "UDP" && new == "SOURCE_IP" {
						return true
					}
					return false
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"port_ranges": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"to_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP", "HTTP", "HTTPS"}, false),
				Default:      "TCP",
			},
			"proxy_protocol": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudGaListenerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	action := "CreateListener"
	request := make(map[string]interface{})
	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}
	request["AcceleratorId"] = d.Get("accelerator_id")
	if v, ok := d.GetOk("certificates"); ok {
		Certificates := make([]map[string]interface{}, len(v.([]interface{})))
		for i, CertificatesValue := range v.([]interface{}) {
			CertificatesMap := CertificatesValue.(map[string]interface{})
			Certificates[i] = make(map[string]interface{})
			Certificates[i]["Id"] = CertificatesMap["id"]
		}
		request["Certificates"] = Certificates

	}

	if v, ok := d.GetOk("client_affinity"); ok {
		request["ClientAffinity"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	}

	PortRanges := make([]map[string]interface{}, len(d.Get("port_ranges").([]interface{})))
	for i, PortRangesValue := range d.Get("port_ranges").([]interface{}) {
		PortRangesMap := PortRangesValue.(map[string]interface{})
		PortRanges[i] = make(map[string]interface{})
		PortRanges[i]["FromPort"] = PortRangesMap["from_port"]
		PortRanges[i]["ToPort"] = PortRangesMap["to_port"]
	}
	request["PortRanges"] = PortRanges

	if v, ok := d.GetOk("protocol"); ok {
		request["Protocol"] = v
	}

	if v, ok := d.GetOkExists("proxy_protocol"); ok {
		request["ProxyProtocol"] = v
	}

	request["RegionId"] = client.RegionId
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("CreateListener")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"StateError.Accelerator"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ga_listener", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ListenerId"]))
	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, gaService.GaListenerStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudGaListenerRead(d, meta)
}
func resourceAlicloudGaListenerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	object, err := gaService.DescribeGaListener(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ga_listener gaService.DescribeGaListener Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	certificates := make([]map[string]interface{}, 0)
	if certificatesList, ok := object["Certificates"].([]interface{}); ok {
		for _, v := range certificatesList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"id": m1["Id"],
				}
				certificates = append(certificates, temp1)

			}
		}
	}
	if err := d.Set("certificates", certificates); err != nil {
		return WrapError(err)
	}
	d.Set("client_affinity", object["ClientAffinity"])
	d.Set("description", object["Description"])
	d.Set("name", object["Name"])
	if val, ok := d.GetOk("proxy_protocol"); ok {
		d.Set("proxy_protocol", val)
	}

	portRanges := make([]map[string]interface{}, 0)
	if portRangesList, ok := object["PortRanges"].([]interface{}); ok {
		for _, v := range portRangesList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"from_port": m1["FromPort"],
					"to_port":   m1["ToPort"],
				}
				portRanges = append(portRanges, temp1)

			}
		}
	}
	if err := d.Set("port_ranges", portRanges); err != nil {
		return WrapError(err)
	}
	d.Set("protocol", object["Protocol"])
	d.Set("status", object["State"])
	return nil
}
func resourceAlicloudGaListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"ListenerId": d.Id(),
	}
	if d.HasChange("certificates") {
		update = true
		Certificates := make([]map[string]interface{}, len(d.Get("certificates").([]interface{})))
		for i, CertificatesValue := range d.Get("certificates").([]interface{}) {
			CertificatesMap := CertificatesValue.(map[string]interface{})
			Certificates[i] = make(map[string]interface{})
			Certificates[i]["Id"] = CertificatesMap["id"]
		}
		request["Certificates"] = Certificates

	}
	if d.HasChange("client_affinity") {
		update = true
		request["ClientAffinity"] = d.Get("client_affinity")
	}
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}
	if d.HasChange("name") {
		update = true
		request["Name"] = d.Get("name")
	}
	if d.HasChange("port_ranges") {
		update = true
		PortRanges := make([]map[string]interface{}, len(d.Get("port_ranges").([]interface{})))
		for i, PortRangesValue := range d.Get("port_ranges").([]interface{}) {
			PortRangesMap := PortRangesValue.(map[string]interface{})
			PortRanges[i] = make(map[string]interface{})
			PortRanges[i]["FromPort"] = PortRangesMap["from_port"]
			PortRanges[i]["ToPort"] = PortRangesMap["to_port"]
		}
		request["PortRanges"] = PortRanges

	}
	if d.HasChange("protocol") {
		update = true
		request["Protocol"] = d.Get("protocol")
	}
	request["RegionId"] = client.RegionId
	if update {
		if _, ok := d.GetOkExists("proxy_protocol"); ok {
			request["ProxyProtocol"] = d.Get("proxy_protocol")
		}
		action := "UpdateListener"
		conn, err := client.NewGaplusClient()
		if err != nil {
			return WrapError(err)
		}
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			request["ClientToken"] = buildClientToken("UpdateListener")
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"StateError.Accelerator"}) {
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
		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gaService.GaListenerStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return resourceAlicloudGaListenerRead(d, meta)
}
func resourceAlicloudGaListenerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	action := "DeleteListener"
	var response map[string]interface{}
	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ListenerId": d.Id(),
	}

	request["AcceleratorId"] = d.Get("accelerator_id")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("DeleteListener")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"StateError.Accelerator"}) {
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gaService.GaListenerStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
