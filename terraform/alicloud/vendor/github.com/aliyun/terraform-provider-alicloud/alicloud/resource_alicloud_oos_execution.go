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

func resourceAlicloudOosExecution() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudOosExecutionCreate,
		Read:   resourceAlicloudOosExecutionRead,
		Delete: resourceAlicloudOosExecutionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(11 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"counters": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"end_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"executed_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_parent": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"loop_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Automatic", "Debug"}, false),
				Default:      "Automatic",
			},
			"outputs": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parameters": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "{}",
				ValidateFunc: validation.ValidateJsonString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"parent_execution_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ram_role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"safety_check": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"start_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_content": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"template_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"template_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"update_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudOosExecutionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	oosService := OosService{client}
	var response map[string]interface{}
	action := "StartExecution"
	request := make(map[string]interface{})
	conn, err := client.NewOosClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("loop_mode"); ok {
		request["LoopMode"] = v
	}

	if v, ok := d.GetOk("mode"); ok {
		request["Mode"] = v
	}

	if v, ok := d.GetOk("parameters"); ok {
		request["Parameters"] = v
	}

	if v, ok := d.GetOk("parent_execution_id"); ok {
		request["ParentExecutionId"] = v
	}

	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("safety_check"); ok {
		request["SafetyCheck"] = v
	}

	if v, ok := d.GetOk("template_content"); ok {
		request["TemplateContent"] = v
	}

	request["TemplateName"] = d.Get("template_name")
	if v, ok := d.GetOk("template_version"); ok {
		request["TemplateVersion"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-06-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oos_execution", action, AlibabaCloudSdkGoERROR)
	}
	responseExecution := response["Execution"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseExecution["ExecutionId"]))
	stateConf := BuildStateConf([]string{}, []string{"Success"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, oosService.OosExecutionStateRefreshFunc(d.Id(), []string{"Failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudOosExecutionRead(d, meta)
}
func resourceAlicloudOosExecutionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	oosService := OosService{client}
	object, err := oosService.DescribeOosExecution(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oos_execution oosService.DescribeOosExecution Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("counters", object["Counters"])
	d.Set("create_date", object["CreateDate"])
	d.Set("end_date", object["EndDate"])
	d.Set("executed_by", object["ExecutedBy"])
	d.Set("is_parent", object["IsParent"])
	d.Set("mode", object["Mode"])
	d.Set("outputs", object["Outputs"])
	d.Set("parameters", object["Parameters"])
	d.Set("parent_execution_id", object["ParentExecutionId"])
	d.Set("ram_role", object["RamRole"])
	d.Set("start_date", object["StartDate"])
	d.Set("status", object["Status"])
	d.Set("status_message", object["StatusMessage"])
	d.Set("template_id", object["TemplateId"])
	d.Set("template_name", object["TemplateName"])
	d.Set("template_version", object["TemplateVersion"])
	d.Set("update_date", object["UpdateDate"])
	return nil
}
func resourceAlicloudOosExecutionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteExecutions"
	var response map[string]interface{}
	conn, err := client.NewOosClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ExecutionIds": convertListToJsonString(convertListStringToListInterface([]string{d.Id()})),
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-06-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return nil
}
