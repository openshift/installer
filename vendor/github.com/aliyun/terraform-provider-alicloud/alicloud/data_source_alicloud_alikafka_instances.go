package alicloud

import (
	"regexp"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudAlikafkaInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAlikafkaInstancesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"deploy_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"io_max": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"eip_max": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"disk_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"disk_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"topic_quota": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"paid_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spec_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_point": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"config": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudAlikafkaInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}

	request := alikafka.CreateGetInstanceListRequest()
	request.RegionId = client.RegionId

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
		return alikafkaClient.GetInstanceList(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_alikafka_instances", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*alikafka.GetInstanceListResponse)

	var filteredInstances []alikafka.InstanceVO
	nameRegex, ok := d.GetOk("name_regex")
	if (ok && nameRegex.(string) != "") || (len(idsMap) > 0) {
		var r *regexp.Regexp
		if nameRegex != "" {
			r, err = regexp.Compile(nameRegex.(string))
			if err != nil {
				return WrapError(err)
			}
		}
		for _, instance := range response.InstanceList.InstanceVO {
			if r != nil && !r.MatchString(instance.Name) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[instance.InstanceId]; !ok {
					continue
				}
			}

			filteredInstances = append(filteredInstances, instance)
		}
	} else {
		filteredInstances = response.InstanceList.InstanceVO
	}
	return alikafkaInstancesDecriptionAttributes(d, filteredInstances, meta)
}

func alikafkaInstancesDecriptionAttributes(d *schema.ResourceData, instancesInfo []alikafka.InstanceVO, meta interface{}) error {
	var ids []string
	var names []string
	var s []map[string]interface{}

	for _, item := range instancesInfo {

		paidType := PostPaid
		if item.PaidType == 0 {
			paidType = PrePaid
		}
		mapping := map[string]interface{}{
			"id":              item.InstanceId,
			"name":            item.Name,
			"create_time":     time.Unix(int64(item.CreateTime)/1000, 0).Format("2006-01-02 03:04:05"),
			"service_status":  item.ServiceStatus,
			"deploy_type":     item.DeployType,
			"vpc_id":          item.VpcId,
			"vswitch_id":      item.VSwitchId,
			"io_max":          item.IoMax,
			"eip_max":         item.EipMax,
			"disk_type":       item.DiskType,
			"disk_size":       item.DiskSize,
			"topic_quota":     item.TopicNumLimit,
			"paid_type":       paidType,
			"spec_type":       item.SpecType,
			"zone_id":         item.ZoneId,
			"end_point":       item.EndPoint,
			"security_group":  item.SecurityGroup,
			"service_version": item.UpgradeServiceDetailInfo.Current2OpenSourceVersion,
			"config":          item.AllConfig,
		}

		ids = append(ids, item.InstanceId)
		names = append(names, item.Name)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("instances", s); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil

}
