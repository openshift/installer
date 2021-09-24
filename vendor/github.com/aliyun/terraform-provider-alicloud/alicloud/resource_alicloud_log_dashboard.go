package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sort"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudLogDashboard() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLogDashboardCreate,
		Read:   resourceAlicloudLogDashboardRead,
		Update: resourceAlicloudLogDashboardUpdate,
		Delete: resourceAlicloudLogDashboardDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dashboard_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"char_list": {
				Type:             schema.TypeString,
				DiffSuppressFunc: jsonPolicyDiffSuppress,
				Required:         true,
			},
		},
	}
}

func resourceAlicloudLogDashboardCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var requestInfo *sls.Client

	dashBoard := sls.Dashboard{
		DashboardName: d.Get("dashboard_name").(string),
		DisplayName:   d.Get("display_name").(string),
	}
	jsonErr := json.Unmarshal([]byte(d.Get("char_list").(string)), &dashBoard.ChartList)
	if jsonErr != nil {
		return WrapError(jsonErr)
	}

	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return nil, slsClient.CreateDashboard(d.Get("project_name").(string), dashBoard)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("CreateDashboard", dashBoard, requestInfo, map[string]interface{}{
			"dashBoard": dashBoard,
		})
		d.SetId(fmt.Sprintf("%s%s%s", d.Get("project_name").(string), COLON_SEPARATED, d.Get("dashboard_name").(string)))
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_dashboard", "CreateDashboard", AliyunLogGoSdkERROR)
	}
	return resourceAlicloudLogDashboardRead(d, meta)
}

func resourceAlicloudLogDashboardRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	object, err := logService.DescribeLogDashboard(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("project_name", parts[0])
	d.Set("dashboard_name", object.DashboardName)
	d.Set("display_name", object.DisplayName)
	charlist, err := json.Marshal(object.ChartList)
	if err != nil {
		return WrapError(err)
	}
	d.Set("char_list", string(charlist))
	return nil
}

func resourceAlicloudLogDashboardUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	update := false
	dashboard := sls.Dashboard{}
	dashboard.DisplayName = d.Get("display_name").(string)
	data := d.Get("char_list").(string)
	err := json.Unmarshal([]byte(data), &dashboard.ChartList)
	if err != nil {
		return WrapError(err)
	}

	if d.HasChange("display_name") {
		update = true
	}
	if d.HasChange("char_list") {
		update = true
	}

	if update {
		parts, err := ParseResourceId(d.Id(), 2)
		if err != nil {
			return WrapError(err)
		}
		dashboard.DashboardName = parts[1]
		_, err = client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.UpdateDashboard(parts[0], dashboard)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateDashboard", AliyunLogGoSdkERROR)
		}
	}
	return resourceAlicloudLogDashboardRead(d, meta)
}

func resourceAlicloudLogDashboardDelete(d *schema.ResourceData, meta interface{}) error {
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
			return nil, slsClient.DeleteDashboard(parts[0], parts[1])
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout, "RequestTimeout"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("DeleteDashboard", raw, requestInfo, map[string]interface{}{
			"project_name": parts[0],
			"dashboard":    parts[1],
		})
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"DashboardNotExist", "ProjectNotExist"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteDashboard", AliyunLogGoSdkERROR)
	}
	return WrapError(logService.WaitForLogDashboard(d.Id(), Deleted, DefaultTimeout))
}

func jsonPolicyDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	if old == "" && new == "" {
		return true
	}

	var oldChartList, newChartList []sls.Chart
	if old != "" && new != "" {
		if err := json.Unmarshal([]byte(old), &oldChartList); err != nil {
			log.Printf("[ERROR] Could not unmarshal old chart list %s: %v", old, err)
			return false
		}
		if err := json.Unmarshal([]byte(new), &newChartList); err != nil {
			log.Printf("[ERROR] Could not unmarshal new chart list %s: %v", new, err)
			return false
		}

		return compareChartList(newChartList, oldChartList)
	}

	return false
}

type chartSort []sls.Chart

func (a chartSort) Len() int {
	return len(a)
}

func (a chartSort) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a chartSort) Less(i, j int) bool {
	return a[i].Title < a[j].Title
}

func compareChartList(a, b []sls.Chart) bool {
	if len(a) != len(b) {
		return false
	}
	sort.Sort(chartSort(a))
	sort.Sort(chartSort(b))
	for i, chart := range a {
		if !reflect.DeepEqual(chart, b[i]) {
			return false
		}
	}
	return true
}
