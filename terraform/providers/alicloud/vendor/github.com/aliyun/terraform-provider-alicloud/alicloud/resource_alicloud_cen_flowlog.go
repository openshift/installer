package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudCenFlowlog() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenFlowlogCreate,
		Read:   resourceAlicloudCenFlowlogRead,
		Update: resourceAlicloudCenFlowlogUpdate,
		Delete: resourceAlicloudCenFlowlogDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cen_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flow_log_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"log_store_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Active", "Inactive"}, false),
				Default:      "Active",
			},
		},
	}
}

func resourceAlicloudCenFlowlogCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := cbn.CreateCreateFlowlogRequest()
	request.CenId = d.Get("cen_id").(string)
	if v, ok := d.GetOk("description"); ok {
		request.Description = v.(string)
	}
	if v, ok := d.GetOk("flow_log_name"); ok {
		request.FlowLogName = v.(string)
	}
	request.LogStoreName = d.Get("log_store_name").(string)
	request.ProjectName = d.Get("project_name").(string)
	request.RegionId = client.RegionId
	request.ClientToken = buildClientToken(request.GetActionName())
	raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
		return cbnClient.CreateFlowlog(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_flowlog", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*cbn.CreateFlowlogResponse)
	d.SetId(response.FlowLogId)

	return resourceAlicloudCenFlowlogUpdate(d, meta)
}
func resourceAlicloudCenFlowlogRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	object, err := cbnService.DescribeCenFlowlog(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("cen_id", object.CenId)
	d.Set("description", object.Description)
	d.Set("flow_log_name", object.FlowLogName)
	d.Set("log_store_name", object.LogStoreName)
	d.Set("project_name", object.ProjectName)
	d.Set("status", object.Status)
	return nil
}
func resourceAlicloudCenFlowlogUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	d.Partial(true)

	update := false
	request := cbn.CreateModifyFlowLogAttributeRequest()
	request.FlowLogId = d.Id()
	request.CenId = d.Get("cen_id").(string)
	request.RegionId = client.RegionId
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request.Description = d.Get("description").(string)
	}
	if !d.IsNewResource() && d.HasChange("flow_log_name") {
		update = true
		request.FlowLogName = d.Get("flow_log_name").(string)
	}
	if update {
		err := resource.Retry(30*time.Second, func() *resource.RetryError {
			raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
				return cbnClient.ModifyFlowLogAttribute(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"LOCK_ERROR"}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("description")
		d.SetPartial("flow_log_name")
	}
	if d.HasChange("status") {
		object, err := cbnService.DescribeCenFlowlog(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := d.Get("status").(string)
		if object.Status != target {
			if target == "Active" {
				request := cbn.CreateActiveFlowLogRequest()
				request.FlowLogId = d.Id()
				request.CenId = d.Get("cen_id").(string)
				request.RegionId = client.RegionId
				err := resource.Retry(30*time.Second, func() *resource.RetryError {
					raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
						return cbnClient.ActiveFlowLog(request)
					})
					if err != nil {
						if IsExpectedErrors(err, []string{"LOCK_ERROR"}) {
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					addDebug(request.GetActionName(), raw)
					return nil
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
			}
			if target == "Inactive" {
				request := cbn.CreateDeactiveFlowLogRequest()
				request.FlowLogId = d.Id()
				request.CenId = d.Get("cen_id").(string)
				request.RegionId = client.RegionId
				err := resource.Retry(30*time.Second, func() *resource.RetryError {
					raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
						return cbnClient.DeactiveFlowLog(request)
					})
					if err != nil {
						if IsExpectedErrors(err, []string{"LOCK_ERROR"}) {
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					addDebug(request.GetActionName(), raw)
					return nil
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
			}
			d.SetPartial("status")
		}
	}
	d.Partial(false)
	return resourceAlicloudCenFlowlogRead(d, meta)
}
func resourceAlicloudCenFlowlogDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := cbn.CreateDeleteFlowlogRequest()
	request.FlowLogId = d.Id()
	request.CenId = d.Get("cen_id").(string)
	request.RegionId = client.RegionId
	raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
		return cbnClient.DeleteFlowlog(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidFlowlogId.NotFound", "ProjectOrLogstoreNotExist", "SourceProjectNotExist"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}
