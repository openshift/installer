package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCenFlowlogs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCenFlowlogsRead,
		Schema: map[string]*schema.Schema{
			"cen_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"log_store_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Active", "Inactive"}, false),
				Default:      "Active",
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flowlogs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cen_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"flow_log_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"flow_log_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_store_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCenFlowlogsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := cbn.CreateDescribeFlowlogsRequest()
	if v, ok := d.GetOk("cen_id"); ok {
		request.CenId = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {
		request.Description = v.(string)
	}
	if v, ok := d.GetOk("flow_log_name"); ok {
		request.FlowLogName = v.(string)
	}
	if v, ok := d.GetOk("log_store_name"); ok {
		request.LogStoreName = v.(string)
	}
	if v, ok := d.GetOk("project_name"); ok {
		request.ProjectName = v.(string)
	}
	request.RegionId = client.RegionId
	if v, ok := d.GetOk("status"); ok {
		request.Status = v.(string)
	}
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var objects []cbn.FlowLog
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
	}
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	for {
		request.ClientToken = buildClientToken(request.GetActionName())
		raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeFlowlogs(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cen_flowlogs", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*cbn.DescribeFlowlogsResponse)

		for _, item := range response.FlowLogs.FlowLog {
			if nameRegex != nil {
				if !nameRegex.MatchString(item.FlowLogName) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[item.FlowLogId]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(response.FlowLogs.FlowLog) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}
	ids := make([]string, len(objects))
	names := make([]string, len(objects))
	s := make([]map[string]interface{}, len(objects))

	for i, object := range objects {
		mapping := map[string]interface{}{
			"cen_id":         object.CenId,
			"description":    object.Description,
			"id":             object.FlowLogId,
			"flow_log_id":    object.FlowLogId,
			"flow_log_name":  object.FlowLogName,
			"log_store_name": object.LogStoreName,
			"project_name":   object.ProjectName,
			"status":         object.Status,
		}
		ids[i] = object.FlowLogId
		names[i] = object.FlowLogName
		s[i] = mapping
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("flowlogs", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
