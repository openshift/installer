package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCddcDedicatedHostGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCddcDedicatedHostGroupCreate,
		Read:   resourceAlicloudCddcDedicatedHostGroupRead,
		Update: resourceAlicloudCddcDedicatedHostGroupUpdate,
		Delete: resourceAlicloudCddcDedicatedHostGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"allocation_policy": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Evenly", "Intensively"}, false),
			},
			"cpu_allocation_ratio": {
				Type:         schema.TypeInt,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.IntBetween(100, 300),
			},
			"dedicated_host_group_desc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disk_allocation_ratio": {
				Type:         schema.TypeInt,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.IntBetween(100, 300),
			},
			"engine": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Redis", "SQLServer", "MySQL", "PostgreSQL", "MongoDB"}, false),
			},
			"host_replace_policy": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"mem_allocation_ratio": {
				Type:         schema.TypeInt,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 100),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudCddcDedicatedHostGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateDedicatedHostGroup"
	request := make(map[string]interface{})
	conn, err := client.NewCddcClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("allocation_policy"); ok {
		request["AllocationPolicy"] = v
	}
	if v, ok := d.GetOk("cpu_allocation_ratio"); ok {
		request["CpuAllocationRatio"] = v
	}
	if v, ok := d.GetOk("dedicated_host_group_desc"); ok {
		request["DedicatedHostGroupDesc"] = v
	}

	request["Engine"] = d.Get("engine")
	if v, ok := d.GetOk("disk_allocation_ratio"); ok {
		if d.Get("engine").(string) == "SQLServer" && v.(int) > 100 {
			return WrapError(fmt.Errorf("disk_allocation_ratio needs to be less than 100 under the SQLServer"))
		}
		request["DiskAllocationRatio"] = v
	}
	if v, ok := d.GetOk("host_replace_policy"); ok {
		request["HostReplacePolicy"] = v
	}
	if v, ok := d.GetOk("mem_allocation_ratio"); ok {
		request["MemAllocationRatio"] = v
	}
	request["RegionId"] = client.RegionId
	request["VPCId"] = d.Get("vpc_id")
	request["ClientToken"] = buildClientToken("CreateDedicatedHostGroup")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-20"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cddc_dedicated_host_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DedicatedHostGroupId"]))

	return resourceAlicloudCddcDedicatedHostGroupRead(d, meta)
}
func resourceAlicloudCddcDedicatedHostGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cddcService := CddcService{client}
	object, err := cddcService.DescribeCddcDedicatedHostGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cddc_dedicated_host_group cddcService.DescribeCddcDedicatedHostGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("allocation_policy", object["AllocationPolicy"])
	if v, ok := object["CpuAllocationRatio"]; ok && fmt.Sprint(v) != "0" {
		d.Set("cpu_allocation_ratio", formatInt(v))
	}
	d.Set("dedicated_host_group_desc", object["DedicatedHostGroupDesc"])
	if v, ok := object["DiskAllocationRatio"]; ok && fmt.Sprint(v) != "0" {
		d.Set("disk_allocation_ratio", formatInt(v))
	}

	d.Set("engine", switchEngine(object["Engine"].(string)))
	d.Set("host_replace_policy", object["HostReplacePolicy"])
	if v, ok := object["MemAllocationRatio"]; ok && fmt.Sprint(v) != "0" {
		d.Set("mem_allocation_ratio", formatInt(v))
	}
	d.Set("vpc_id", object["VPCId"])
	return nil
}
func resourceAlicloudCddcDedicatedHostGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"DedicatedHostGroupId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("allocation_policy") {
		update = true
		if v, ok := d.GetOk("allocation_policy"); ok {
			request["AllocationPolicy"] = v
		}
	}
	if d.HasChange("cpu_allocation_ratio") {
		update = true
		if v, ok := d.GetOk("cpu_allocation_ratio"); ok {
			request["CpuAllocationRatio"] = v
		}
	}
	if d.HasChange("dedicated_host_group_desc") {
		update = true
		if v, ok := d.GetOk("dedicated_host_group_desc"); ok {
			request["DedicatedHostGroupDesc"] = v
		}
	}
	if d.HasChange("disk_allocation_ratio") {
		update = true
		if v, ok := d.GetOk("disk_allocation_ratio"); ok {
			request["DiskAllocationRatio"] = v
		}
	}
	if d.HasChange("host_replace_policy") {
		update = true
		if v, ok := d.GetOk("host_replace_policy"); ok {
			request["HostReplacePolicy"] = v
		}
	}
	if d.HasChange("mem_allocation_ratio") {
		update = true
		if v, ok := d.GetOk("mem_allocation_ratio"); ok {
			request["MemAllocationRatio"] = v
		}
	}
	if update {
		action := "ModifyDedicatedHostGroupAttribute"
		conn, err := client.NewCddcClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	}
	return resourceAlicloudCddcDedicatedHostGroupRead(d, meta)
}
func resourceAlicloudCddcDedicatedHostGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDedicatedHostGroup"
	var response map[string]interface{}
	conn, err := client.NewCddcClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DedicatedHostGroupId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return nil
}

func switchEngine(engine string) string {
	switch engine {
	case "mysql":
		engine = "MySQL"
	case "redis":
		engine = "Redis"
	case "mssql":
		engine = "SQLServer"
	case "pgsql":
		engine = "PostgreSQL"
	case "mongodb":
		engine = "MongoDB"
	}
	return engine
}
