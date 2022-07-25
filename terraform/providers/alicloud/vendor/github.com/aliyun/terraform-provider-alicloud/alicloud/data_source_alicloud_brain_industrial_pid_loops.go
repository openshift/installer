package alicloud

import (
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudBrainIndustrialPidLoops() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudBrainIndustrialPidLoopsRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"pid_loop_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"pid_project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"0", "1"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"loops": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pid_loop_configuration": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pid_loop_dcs_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pid_loop_desc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pid_loop_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pid_loop_is_crucial": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"pid_loop_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pid_loop_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pid_project_id": {
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
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAlicloudBrainIndustrialPidLoopsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListPidLoops"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("pid_loop_name"); ok {
		request["PidLoopName"] = v
	}
	request["PidProjectId"] = d.Get("pid_project_id")
	request["PageSize"] = PageSizeLarge
	request["CurrentPage"] = 1
	var objects []map[string]interface{}
	var pidLoopNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		pidLoopNameRegex = r
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
	status, statusOk := d.GetOk("status")
	var response map[string]interface{}
	conn, err := client.NewAistudioClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-20"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_brain_industrial_pid_loops", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.PidLoopList", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.PidLoopList", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if pidLoopNameRegex != nil {
				if !pidLoopNameRegex.MatchString(fmt.Sprint(item["PidLoopName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["PidLoopId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["ReleaseStatus"].(string) {
				continue
			}
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["CurrentPage"] = request["CurrentPage"].(int) + 1
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"pid_loop_dcs_type": object["PidLoopDcsType"],
			"id":                fmt.Sprint(object["PidLoopId"]),
			"pid_loop_id":       fmt.Sprint(object["PidLoopId"]),
			"pid_loop_name":     object["PidLoopName"],
			"pid_loop_type":     object["PidLoopType"],
			"status":            object["ReleaseStatus"],
		}
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(object["PidLoopId"]))
			names = append(names, object["PidLoopName"])
			s = append(s, mapping)
			continue
		}

		brain_industrialService := Brain_industrialService{client}
		id := fmt.Sprint(object["PidLoopId"])
		getResp, err := brain_industrialService.DescribeBrainIndustrialPidLoop(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["pid_loop_configuration"] = getResp["PidLoopConfiguration"]
		mapping["pid_loop_desc"] = getResp["PidLoopDesc"]
		mapping["pid_loop_is_crucial"] = getResp["PidLoopIsCrucial"]
		mapping["pid_project_id"] = getResp["PidProjectId"]
		ids = append(ids, fmt.Sprint(object["PidLoopId"]))
		names = append(names, object["PidLoopName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("loops", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
