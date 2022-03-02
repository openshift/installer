package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudDtsSynchronizationJobs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDtsSynchronizationJobsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Downgrade", "Failed", "Finished", "InitializeFailed", "Initializing", "Locked", "Modifying", "NotConfigured", "NotStarted", "PreCheckPass", "PrecheckFailed", "Prechecking", "Retrying", "Suspending", "Synchronizing", "Upgrade"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"jobs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"checkpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_initialization": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"data_synchronization": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"db_list": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"destination_endpoint_data_base_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"destination_endpoint_engine_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"destination_endpoint_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"destination_endpoint_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"destination_endpoint_instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"destination_endpoint_oracle_sid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"destination_endpoint_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"destination_endpoint_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"destination_endpoint_user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dts_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dts_job_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dts_job_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expire_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_endpoint_database_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_endpoint_engine_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_endpoint_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_endpoint_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_endpoint_instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_endpoint_oracle_sid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_endpoint_owner_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_endpoint_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_endpoint_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_endpoint_role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_endpoint_user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"structure_initialization": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"synchronization_direction": {
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

func dataSourceAlicloudDtsSynchronizationJobsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeDtsJobs"
	request := make(map[string]interface{})
	request["JobType"] = "SYNC"
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var jobNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		jobNameRegex = r
	}
	status, statusOk := d.GetOk("status")

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	var response map[string]interface{}
	conn, err := client.NewDtsClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_dts_synchronization_jobs", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.DtsJobList", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.DtsJobList", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			// see comment on function DescribeDtsSynchronizationJob
			if item["Status"] == "synchronizing" || item["Status"] == "Initializing" {
				item["Status"] = "Synchronizing"
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["DtsJobId"])]; !ok {
					continue
				}
			}
			if jobNameRegex != nil && !jobNameRegex.MatchString(fmt.Sprint(item["DtsJobName"])) {
				continue
			}
			if statusOk && status.(string) != "" && status.(string) != item["Status"].(string) {
				continue
			}
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		migrationModeObj := object["MigrationMode"].(map[string]interface{})
		destinationEndpointObj := object["DestinationEndpoint"].(map[string]interface{})
		sourceEndpointObj := object["SourceEndpoint"].(map[string]interface{})
		mapping := map[string]interface{}{
			"checkpoint":                          object["Checkpoint"],
			"create_time":                         object["CreateTime"],
			"data_initialization":                 migrationModeObj["DataInitialization"],
			"data_synchronization":                migrationModeObj["DataSynchronization"],
			"structure_initialization":            migrationModeObj["StructureInitialization"],
			"db_list":                             object["DbObject"],
			"destination_endpoint_data_base_name": destinationEndpointObj["DatabaseName"],
			"destination_endpoint_engine_name":    destinationEndpointObj["EngineName"],
			"destination_endpoint_ip":             destinationEndpointObj["Ip"],
			"destination_endpoint_instance_id":    destinationEndpointObj["InstanceID"],
			"destination_endpoint_instance_type":  destinationEndpointObj["InstanceType"],
			"destination_endpoint_oracle_sid":     destinationEndpointObj["OracleSID"],
			"destination_endpoint_port":           destinationEndpointObj["Port"],
			"destination_endpoint_region":         destinationEndpointObj["Region"],
			"destination_endpoint_user_name":      destinationEndpointObj["UserName"],
			"dts_instance_id":                     object["DtsInstanceID"],
			"id":                                  fmt.Sprint(object["DtsJobId"]),
			"dts_job_id":                          fmt.Sprint(object["DtsJobId"]),
			"dts_job_name":                        object["DtsJobName"],
			"expire_time":                         object["ExpireTime"],
			"source_endpoint_database_name":       object["ReverseJob"].(map[string]interface{})["SourceEndpoint"].(map[string]interface{})["DatabaseName"],
			"source_endpoint_engine_name":         sourceEndpointObj["EngineName"],
			"source_endpoint_ip":                  sourceEndpointObj["Ip"],
			"source_endpoint_instance_id":         sourceEndpointObj["InstanceID"],
			"source_endpoint_instance_type":       sourceEndpointObj["InstanceType"],
			"source_endpoint_oracle_sid":          sourceEndpointObj["OracleSID"],
			"source_endpoint_port":                sourceEndpointObj["Port"],
			"source_endpoint_region":              sourceEndpointObj["Region"],
			"source_endpoint_user_name":           sourceEndpointObj["UserName"],
			"status":                              object["Status"],
			"synchronization_direction":           object["DtsJobDirection"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["DtsJobId"])
		dtsService := DtsService{client}
		getResp, err := dtsService.DescribeDtsSynchronizationJob(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["source_endpoint_owner_id"] = getResp["SourceEndpoint"].(map[string]interface{})["AliyunUid"]
		mapping["source_endpoint_role"] = getResp["SourceEndpoint"].(map[string]interface{})["RoleName"]
		if err != nil {
			return WrapError(err)
		}

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("jobs", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
