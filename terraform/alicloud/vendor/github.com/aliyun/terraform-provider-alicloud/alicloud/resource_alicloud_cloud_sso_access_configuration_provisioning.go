package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudCloudSsoAccessConfigurationProvisioning() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCloudSsoAccessConfigurationProvisioningCreate,
		Read:   resourceAlicloudCloudSsoAccessConfigurationProvisioningRead,
		Update: resourceAlicloudCloudSsoAccessConfigurationProvisioningUpdate,
		Delete: resourceAlicloudCloudSsoAccessConfigurationProvisioningDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"access_configuration_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"directory_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"status": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Provisioned"}, false),
			},
		},
	}
}

func resourceAlicloudCloudSsoAccessConfigurationProvisioningCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudssoService := CloudssoService{client}
	err := cloudssoService.CloudssoServicAccessConfigurationProvisioning(fmt.Sprint(d.Get("directory_id")), fmt.Sprint(d.Get("access_configuration_id")), fmt.Sprint(d.Get("target_type")), fmt.Sprint(d.Get("target_id")))
	if err != nil {
		return WrapError(err)
	}
	d.SetId(fmt.Sprint(d.Get("directory_id"), ":", d.Get("access_configuration_id"), ":", d.Get("target_type"), ":", d.Get("target_id")))
	return resourceAlicloudCloudSsoAccessConfigurationProvisioningRead(d, meta)
}
func resourceAlicloudCloudSsoAccessConfigurationProvisioningRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudssoService := CloudssoService{client}
	object, err := cloudssoService.DescribeCloudSsoAccessConfigurationProvisioning(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_sso_access_configuration_provisioning cloudssoService.DescribeCloudSsoAccessConfigurationProvisioning Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}

	d.Set("access_configuration_id", parts[1])
	d.Set("directory_id", parts[0])
	d.Set("target_id", parts[3])
	d.Set("target_type", parts[2])
	d.Set("status", object["Status"])
	return nil
}
func resourceAlicloudCloudSsoAccessConfigurationProvisioningUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudssoService := CloudssoService{client}

	if d.HasChange("status") {
		object, err := cloudssoService.DescribeCloudSsoAccessConfigurationProvisioning(d.Id())
		if err != nil && !NotFoundError(err) {
			return WrapError(err)
		}

		parts, err := ParseResourceId(d.Id(), 4)
		if err != nil {
			return WrapError(err)
		}

		if v, ok := object["Status"]; !ok || fmt.Sprint(v) != "Provisioned" {
			cloudssoService := CloudssoService{client}
			err := cloudssoService.CloudssoServicAccessConfigurationProvisioning(fmt.Sprint(parts[0]), fmt.Sprint(parts[1]), fmt.Sprint(parts[2]), fmt.Sprint(parts[3]))
			if err != nil {
				return WrapError(err)
			}
		}
	}

	return resourceAlicloudCloudSsoAccessConfigurationProvisioningRead(d, meta)
}
func resourceAlicloudCloudSsoAccessConfigurationProvisioningDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}
	action := "DeprovisionAccessConfiguration"
	var response map[string]interface{}
	conn, err := client.NewCloudssoClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"AccessConfigurationId": parts[1],
		"DirectoryId":           parts[0],
		"TargetId":              parts[3],
		"TargetType":            parts[2],
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
		if IsExpectedErrors(err, []string{"EntityNotExists.AccessConfigurationProvisioning"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Tasks", response)
	if err != nil || len(v.([]interface{})) < 1 {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	response = v.([]interface{})[0].(map[string]interface{})
	cloudssoService := CloudssoService{client}
	_, err = cloudssoService.GetTaskStatus(fmt.Sprint(request["DirectoryId"]), fmt.Sprint(response["TaskId"]))
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}
	stateConf := BuildStateConf([]string{}, []string{"Success"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cloudssoService.CloudssoServiceAccessConfigurationProvisioningStateRefreshFunc(fmt.Sprint(request["DirectoryId"]), fmt.Sprint(response["TaskId"]), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
