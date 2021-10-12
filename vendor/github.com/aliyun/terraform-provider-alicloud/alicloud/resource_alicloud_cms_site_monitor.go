package alicloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudCmsSiteMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCmsSiteMonitorCreate,
		Read:   resourceAlicloudCmsSiteMonitorRead,
		Update: resourceAlicloudCmsSiteMonitorUpdate,
		Delete: resourceAlicloudCmsSiteMonitorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"task_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"task_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{SiteMonitorHTTP, SiteMonitorDNS, SiteMonitorFTP, SiteMonitorPOP3,
					SiteMonitorPing, SiteMonitorSMTP, SiteMonitorTCP, SiteMonitorUDP}, false),
			},
			"alert_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntInSlice([]int{1, 5, 15}),
			},
			"options_json": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"isp_cities": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"city": {
							Type:     schema.TypeString,
							Required: true,
						},
						"isp": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"task_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCmsSiteMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	taskName := d.Get("task_name").(string)
	request := cms.CreateCreateSiteMonitorRequest()
	request.Address = d.Get("address").(string)
	request.TaskName = taskName
	request.TaskType = d.Get("task_type").(string)
	request.Interval = strconv.Itoa(d.Get("interval").(int))
	request.OptionsJson = d.Get("options_json").(string)
	alertIds := d.Get("alert_ids").([]interface{})
	alertId := getAlertId(alertIds)
	if alertId != "" {
		request.AlertIds = alertId
	}

	if isp_cities, ok := d.GetOk("isp_cities"); ok {
		var a []map[string]interface{}
		for _, element := range isp_cities.(*schema.Set).List() {
			isp_city := element.(map[string]interface{})
			a = append(a, isp_city)
		}
		b, err := json.Marshal(a)
		if err != nil {
			return WrapError(err)
		}
		request.IspCities = bytes.NewBuffer(b).String()
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.CreateSiteMonitor(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser, "ExceedingQuota"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), request.RpcRequest, request)
		resp, _ := raw.(*cms.CreateSiteMonitorResponse)
		d.SetId(resp.CreateResultList.CreateResultListItem[0].TaskId)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_site_monitor", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return resourceAlicloudCmsSiteMonitorRead(d, meta)
}

func resourceAlicloudCmsSiteMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}

	siteMonitor, err := cmsService.DescribeSiteMonitor(d.Id(), "")
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("address", siteMonitor.Address)
	d.Set("task_name", siteMonitor.TaskName)
	d.Set("task_type", siteMonitor.TaskType)
	d.Set("task_state", siteMonitor.TaskState)
	d.Set("interval", siteMonitor.Interval)
	d.Set("options_json", siteMonitor.OptionsJson)
	d.Set("create_time", siteMonitor.CreateTime)
	d.Set("update_time", siteMonitor.UpdateTime)

	ispCities, err := cmsService.GetIspCities(d.Id())
	var list []map[string]interface{}

	for _, e := range ispCities {
		list = append(list, map[string]interface{}{"city": e["city"], "isp": e["isp"]})
	}

	if err = d.Set("isp_cities", list); err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceAlicloudCmsSiteMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := cms.CreateModifySiteMonitorRequest()
	request.TaskId = d.Id()
	request.Address = d.Get("address").(string)
	request.Interval = strconv.Itoa(d.Get("interval").(int))
	request.OptionsJson = d.Get("options_json").(string)
	request.TaskName = d.Get("task_name").(string)
	alertIds := d.Get("alert_ids").([]interface{})
	alertId := getAlertId(alertIds)
	if alertId != "" {
		request.AlertIds = alertId
	}

	if isp_cities, ok := d.GetOk("isp_cities"); ok {
		var a []map[string]interface{}
		for _, element := range isp_cities.(*schema.Set).List() {
			isp_city := element.(map[string]interface{})
			a = append(a, isp_city)
		}
		b, err := json.Marshal(a)
		if err != nil {
			return WrapError(err)
		}
		request.IspCities = bytes.NewBuffer(b).String()
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.ModifySiteMonitor(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser, "ExceedingQuota"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapError(err)
	}

	return resourceAlicloudCmsSiteMonitorRead(d, meta)
}

func resourceAlicloudCmsSiteMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}
	request := cms.CreateDeleteSiteMonitorsRequest()

	request.TaskIds = d.Id()
	request.IsDeleteAlarms = "false"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DeleteSiteMonitors(request)
		})

		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser, "ExceedingQuota"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(fmt.Errorf("Deleting site monitor got an error: %#v", err))
		}

		_, err = cmsService.DescribeSiteMonitor(d.Id(), "")
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("DescribeSiteMonitor got an error: %#v", err))
		}

		return resource.RetryableError(fmt.Errorf("Deleting site monitor got an error: %#v", err))

	})
}

func getAlertId(alertIds []interface{}) string {
	if alertIds != nil && len(alertIds) > 0 {
		alertId := strings.Join(expandStringList(alertIds), ",")
		return alertId
	}
	return ""
}
