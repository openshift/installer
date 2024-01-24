package alicloud

import (
	"fmt"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudLogAlert() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLogAlertCreate,
		Read:   resourceAlicloudLogAlertRead,
		Update: resourceAlicloudLogAlertUpdate,
		Delete: resourceAlicloudLogAlertDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"alert_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"alert_displayname": {
				Type:     schema.TypeString,
				Required: true,
			},
			"alert_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"condition": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dashboard": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mute_until": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(0),
				Default:      time.Now().Unix(),
			},
			"throttling": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notify_threshold": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},

			"query_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"chart_title": {
							Type:     schema.TypeString,
							Required: true,
						},
						"logstore": {
							Type:     schema.TypeString,
							Required: true,
						},
						"query": {
							Type:     schema.TypeString,
							Required: true,
						},
						"start": {
							Type:     schema.TypeString,
							Required: true,
						},
						"end": {
							Type:     schema.TypeString,
							Required: true,
						},
						"time_span_type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "Custom",
						},
					},
				},
			},

			"notification_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								sls.NotificationTypeSMS,
								sls.NotificationTypeDingTalk,
								sls.NotificationTypeEmail,
								sls.NotificationTypeMessageCenter},
								false),
						},
						"content": {
							Type:     schema.TypeString,
							Required: true,
						},
						"service_uri": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"mobile_list": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"email_list": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

			"schedule_interval": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "60s",
			},
			"schedule_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "FixedRate",
			},
		},
	}
}

func resourceAlicloudLogAlertCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	project_name := d.Get("project_name").(string)
	alert_name := d.Get("alert_name").(string)
	alert_displayname := d.Get("alert_displayname").(string)

	alert := &sls.Alert{
		Name:        alert_name,
		DisplayName: alert_displayname,
		Description: d.Get("alert_description").(string),
		State:       "Enabled",
		Schedule: &sls.Schedule{
			Type:     d.Get("schedule_type").(string),
			Interval: d.Get("schedule_interval").(string),
		},
	}
	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			dashboard := d.Get("dashboard").(string)
			err := CreateDashboard(project_name, dashboard, slsClient)
			if err != nil {
				return nil, err
			}
			alert.Configuration = createAlertConfig(d, project_name, dashboard, slsClient)
			return nil, slsClient.CreateAlert(project_name, alert)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_alert", "CreateLogstoreAlert", AliyunLogGoSdkERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", project_name, COLON_SEPARATED, alert_name))
	return resourceAlicloudLogAlertRead(d, meta)
}

func resourceAlicloudLogAlertRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	object, err := logService.DescribeLogAlert(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("project_name", parts[0])
	d.Set("alert_name", parts[1])
	d.Set("alert_displayname", object.DisplayName)
	d.Set("alert_description", object.Description)
	d.Set("condition", object.Configuration.Condition)
	d.Set("dashboard", object.Configuration.Dashboard)
	d.Set("mute_until", object.Configuration.MuteUntil)
	d.Set("throttling", object.Configuration.Throttling)
	d.Set("notify_threshold", object.Configuration.NotifyThreshold)
	d.Set("schedule_interval", object.Schedule.Interval)
	d.Set("schedule_type", object.Schedule.Type)

	var notiList []map[string]interface{}

	for _, v := range object.Configuration.NotificationList {
		mapping := getNotiMap(v)
		notiList = append(notiList, mapping)
	}

	var queryList []map[string]interface{}
	for _, v := range object.Configuration.QueryList {
		mapping := map[string]interface{}{
			"chart_title":    v.ChartTitle,
			"logstore":       v.LogStore,
			"query":          v.Query,
			"start":          v.Start,
			"end":            v.End,
			"time_span_type": v.TimeSpanType,
		}
		queryList = append(queryList, mapping)
	}

	d.Set("notification_list", notiList)
	d.Set("query_list", queryList)

	return nil

}

func resourceAlicloudLogAlertUpdate(d *schema.ResourceData, meta interface{}) error {

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	client := meta.(*connectivity.AliyunClient)
	params := &sls.Alert{
		Name:        parts[1],
		DisplayName: d.Get("alert_displayname").(string),
		Description: d.Get("alert_description").(string),
		State:       "Enabled",
		Schedule: &sls.Schedule{
			Type:     d.Get("schedule_type").(string),
			Interval: d.Get("schedule_interval").(string),
		},
	}

	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			project_name := d.Get("project_name").(string)
			dashboard := d.Get("dashboard").(string)
			err := CreateDashboard(project_name, dashboard, slsClient)
			if err != nil {
				return nil, err
			}
			params.Configuration = createAlertConfig(d, project_name, dashboard, slsClient)
			return nil, slsClient.UpdateAlert(parts[0], params)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateAlert", AliyunLogGoSdkERROR)
	}

	return resourceAlicloudLogAlertRead(d, meta)
}

func resourceAlicloudLogAlertDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	var requestInfo *sls.Client
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return nil, slsClient.DeleteAlert(parts[0], parts[1])
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			addDebug("DeleteAlert", raw, requestInfo, map[string]interface{}{
				"project_name": parts[0],
				"alert":        parts[1],
			})
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_alert", "DeleteAlert", AliyunLogGoSdkERROR)
	}
	return WrapError(logService.WaitForLogstoreAlert(d.Id(), Deleted, DefaultTimeout))
}

func createAlertConfig(d *schema.ResourceData, project, dashboard string, client *sls.Client) *sls.AlertConfiguration {

	noti := []*sls.Notification{}
	if v, ok := d.GetOk("notification_list"); ok {
		for _, e := range v.([]interface{}) {
			noti_map := e.(map[string]interface{})
			content := noti_map["content"].(string)

			email_list := []string{}
			email_list_temp := noti_map["email_list"].(*schema.Set).List()
			for _, v := range email_list_temp {
				new_v := v.(string)
				email_list = append(email_list, new_v)
			}
			mobile_list_temp := noti_map["mobile_list"].(*schema.Set).List()
			mobile_list := []string{}
			if len(mobile_list_temp) > 0 {
				for _, v := range mobile_list_temp {
					new_v := v.(string)
					mobile_list = append(mobile_list, new_v)
				}
			}

			if noti_map["type"].(string) == sls.NotificationTypeEmail {
				email := &sls.Notification{
					Type:      sls.NotificationTypeEmail,
					EmailList: email_list,
					Content:   content,
				}
				noti = append(noti, email)
			}

			if noti_map["type"].(string) == sls.NotificationTypeSMS {
				sms := &sls.Notification{
					Type:       sls.NotificationTypeSMS,
					MobileList: mobile_list,
					Content:    content,
				}
				noti = append(noti, sms)
			}
			if noti_map["type"].(string) == sls.NotificationTypeDingTalk {
				ding := &sls.Notification{
					Type:       sls.NotificationTypeDingTalk,
					ServiceUri: noti_map["service_uri"].(string),
					Content:    content,
				}
				noti = append(noti, ding)
			}
			if noti_map["type"].(string) == sls.NotificationTypeMessageCenter {
				messageCenter := &sls.Notification{
					Type:    sls.NotificationTypeMessageCenter,
					Content: content,
				}
				noti = append(noti, messageCenter)
			}
		}
	}

	queryList := []*sls.AlertQuery{}

	if v, ok := d.GetOk("query_list"); ok {
		for _, e := range v.([]interface{}) {
			query_map := e.(map[string]interface{})
			query := &sls.AlertQuery{
				ChartTitle:   GetCharTitile(project, dashboard, query_map["chart_title"].(string), client),
				LogStore:     query_map["logstore"].(string),
				Query:        query_map["query"].(string),
				Start:        query_map["start"].(string),
				End:          query_map["end"].(string),
				TimeSpanType: query_map["time_span_type"].(string),
			}
			queryList = append(queryList, query)

		}
	}

	config := &sls.AlertConfiguration{
		Condition:        d.Get("condition").(string),
		Dashboard:        d.Get("dashboard").(string),
		QueryList:        queryList,
		MuteUntil:        int64(d.Get("mute_until").(int)),
		NotificationList: noti,
		Throttling:       d.Get("throttling").(string),
		NotifyThreshold:  int32(d.Get("notify_threshold").(int)),
	}
	return config
}

func getNotiMap(v *sls.Notification) map[string]interface{} {
	mapping := make(map[string]interface{})

	mapping["content"] = v.Content
	if v.Type == sls.NotificationTypeSMS {
		mapping["type"] = sls.NotificationTypeSMS
		mapping["mobile_list"] = v.MobileList
	}

	if v.Type == sls.NotificationTypeEmail {
		mapping["type"] = sls.NotificationTypeEmail
		mapping["email_list"] = v.EmailList
	}

	if v.Type == sls.NotificationTypeDingTalk {
		mapping["type"] = sls.NotificationTypeDingTalk
		mapping["service_uri"] = v.ServiceUri
	}

	if v.Type == sls.NotificationTypeMessageCenter {
		mapping["type"] = sls.NotificationTypeMessageCenter
	}
	return mapping

}
