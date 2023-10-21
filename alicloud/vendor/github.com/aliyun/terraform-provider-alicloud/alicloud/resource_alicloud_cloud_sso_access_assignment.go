package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudCloudSsoAccessAssignment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCloudSsoAccessAssignmentCreate,
		Read:   resourceAlicloudCloudSsoAccessAssignmentRead,
		Update: resourceAlicloudCloudSsoAccessAssignmentUpdate,
		Delete: resourceAlicloudCloudSsoAccessAssignmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"access_configuration_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"deprovision_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"DeprovisionForLastAccessAssignmentOnAccount", "None"}, false),
			},
			"directory_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"principal_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"principal_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Group", "User"}, false),
			},
			"target_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"RD-Account"}, false),
			},
		},
	}
}

func resourceAlicloudCloudSsoAccessAssignmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateAccessAssignment"
	request := make(map[string]interface{})
	conn, err := client.NewCloudssoClient()
	if err != nil {
		return WrapError(err)
	}

	request["AccessConfigurationId"] = d.Get("access_configuration_id")
	request["DirectoryId"] = d.Get("directory_id")
	request["PrincipalId"] = d.Get("principal_id")
	request["PrincipalType"] = d.Get("principal_type")
	request["TargetId"] = d.Get("target_id")
	request["TargetType"] = d.Get("target_type")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationConflict.Task"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_sso_access_assignment", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(request["DirectoryId"], ":", request["AccessConfigurationId"], ":", request["TargetType"], ":", request["TargetId"], ":", request["PrincipalType"], ":", request["PrincipalId"]))

	v, err := jsonpath.Get("$.Task", response)
	if err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	response = v.(map[string]interface{})
	cloudssoService := CloudssoService{client}
	stateConf := BuildStateConf([]string{}, []string{"Success"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cloudssoService.CloudssoServiceAccessAssignmentStateRefreshFunc(fmt.Sprint(request["DirectoryId"]), fmt.Sprint(response["TaskId"]), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudCloudSsoAccessAssignmentRead(d, meta)
}
func resourceAlicloudCloudSsoAccessAssignmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudssoService := CloudssoService{client}
	_, err := cloudssoService.DescribeCloudSsoAccessAssignment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_sso_access_assignment cloudssoService.DescribeCloudSsoAccessAssignment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 6)
	if err != nil {
		return WrapError(err)
	}
	d.Set("directory_id", parts[0])
	d.Set("access_configuration_id", parts[1])
	d.Set("target_type", parts[2])
	d.Set("target_id", parts[3])
	d.Set("principal_type", parts[4])
	d.Set("principal_id", parts[5])

	return nil
}
func resourceAlicloudCloudSsoAccessAssignmentUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAlicloudCloudSsoAccessAssignmentRead(d, meta)
}
func resourceAlicloudCloudSsoAccessAssignmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 6)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteAccessAssignment"
	var response map[string]interface{}
	conn, err := client.NewCloudssoClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DirectoryId":           parts[0],
		"AccessConfigurationId": parts[1],
		"TargetType":            parts[2],
		"TargetId":              parts[3],
		"PrincipalType":         parts[4],
		"PrincipalId":           parts[5],
	}

	if v, ok := d.GetOk("deprovision_strategy"); ok {
		request["DeprovisionStrategy"] = v
	} else {
		request["DeprovisionStrategy"] = "DeprovisionForLastAccessAssignmentOnAccount"
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationConflict.Task", "DeletionConflict.AccessConfigurationProvisioning.AccessAssignment"}) || NeedRetry(err) {
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

	v, err := jsonpath.Get("$.Task", response)
	if err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	response = v.(map[string]interface{})
	cloudssoService := CloudssoService{client}
	stateConf := BuildStateConf([]string{}, []string{"Success"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cloudssoService.CloudssoServiceAccessAssignmentStateRefreshFunc(fmt.Sprint(request["DirectoryId"]), fmt.Sprint(response["TaskId"]), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
