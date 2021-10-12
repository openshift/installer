package alicloud

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudKvstoreInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudKvstoreInstancesRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"architecture_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"SplitRW", "cluster", "standard"}, false),
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
			"edition_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Enterprise", "Community"}, false),
			},
			"engine_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"2.8", "4.0", "5.0", "6.0"}, false),
			},
			"expired": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"global_instance": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"instance_class": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Memcache", "Redis"}, false),
			},
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"CLASSIC", "VPC"}, false),
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PostPaid", "PrePaid"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"search_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Changing", "CleaningUpExpiredData", "Creating", "Flushing", "HASwitching", "Inactive", "MajorVersionUpgrading", "Migrating", "NetworkModifying", "Normal", "Rebooting", "SSLModifying", "Transforming", "ZoneMigrating"}, false),
			},
			"tags": tagsSchema(),
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"architecture_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auto_renew": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"auto_renew_period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"config": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"connection_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"destroy_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expire_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"has_renew_change_order": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"instance_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_release_protection": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_rds": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"maintain_end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"maintain_start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_connections": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"connections": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"package_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"qps": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"replacate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ssl_enable": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"search_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_ip_group_attribute": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_ip_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secondary_zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_ips": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_auth_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_cloud_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_domain": {
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

func dataSourceAlicloudKvstoreInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := r_kvstore.CreateDescribeInstancesRequest()
	request.RegionId = client.RegionId
	if v, ok := d.GetOk("architecture_type"); ok {
		request.ArchitectureType = v.(string)
	}
	if v, ok := d.GetOk("edition_type"); ok {
		request.EditionType = v.(string)
	}
	if v, ok := d.GetOk("engine_version"); ok {
		request.EngineVersion = v.(string)
	}
	if v, ok := d.GetOk("expired"); ok {
		request.Expired = v.(string)
	}
	if v, ok := d.GetOkExists("global_instance"); ok {
		request.GlobalInstance = requests.NewBoolean(v.(bool))
	}
	if v, ok := d.GetOk("instance_class"); ok {
		request.InstanceClass = v.(string)
	}
	if v, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = v.(string)
	}
	if v, ok := d.GetOk("network_type"); ok {
		request.NetworkType = v.(string)
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request.ChargeType = v.(string)
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResourceGroupId = v.(string)
	}
	if v, ok := d.GetOk("search_key"); ok {
		request.SearchKey = v.(string)
	}
	if v, ok := d.GetOk("status"); ok {
		request.InstanceStatus = v.(string)
	}
	if v, ok := d.GetOk("tags"); ok {
		tags := make([]r_kvstore.DescribeInstancesTag, len(v.(map[string]interface{})))
		i := 0
		for key, value := range v.(map[string]interface{}) {
			tags[i] = r_kvstore.DescribeInstancesTag{
				Key:   key,
				Value: value.(string),
			}
			i++
		}
		request.Tag = &tags
	}
	if v, ok := d.GetOk("vswitch_id"); ok {
		request.VSwitchId = v.(string)
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = v.(string)
	}
	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = v.(string)
	}
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var objects []r_kvstore.KVStoreInstance
	var dBInstanceNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		dBInstanceNameRegex = r
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
	var response *r_kvstore.DescribeInstancesResponse
	for {
		raw, err := client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.DescribeInstances(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_kvstore_instances", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ = raw.(*r_kvstore.DescribeInstancesResponse)

		for _, item := range response.Instances.KVStoreInstance {
			if dBInstanceNameRegex != nil {
				if !dBInstanceNameRegex.MatchString(item.InstanceName) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[item.InstanceId]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(response.Instances.KVStoreInstance) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}
	ids := make([]string, 0)
	names := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		m := make(map[string]string)
		err := json.Unmarshal([]byte(object.Config), &m)
		mapping := map[string]interface{}{
			"architecture_type":      object.ArchitectureType,
			"bandwidth":              object.Bandwidth,
			"capacity":               object.Capacity,
			"config":                 m,
			"connection_mode":        object.ConnectionMode,
			"id":                     object.InstanceId,
			"db_instance_id":         object.InstanceId,
			"db_instance_name":       object.InstanceName,
			"name":                   object.InstanceName,
			"destroy_time":           object.DestroyTime,
			"end_time":               object.EndTime,
			"expire_time":            object.EndTime,
			"engine_version":         object.EngineVersion,
			"has_renew_change_order": object.HasRenewChangeOrder,
			"instance_class":         object.InstanceClass,
			"instance_type":          object.InstanceType,
			"is_rds":                 object.IsRds,
			"max_connections":        object.Connections,
			"connections":            object.Connections,
			"network_type":           object.NetworkType,
			"node_type":              object.NodeType,
			"package_type":           object.PackageType,
			"payment_type":           object.ChargeType,
			"charge_type":            object.ChargeType,
			"port":                   object.Port,
			"private_ip":             object.PrivateIp,
			"qps":                    object.QPS,
			"replacate_id":           object.ReplacateId,
			"resource_group_id":      object.ResourceGroupId,
			"search_key":             object.SearchKey,
			"status":                 object.InstanceStatus,
			"vswitch_id":             object.VSwitchId,
			"vpc_cloud_instance_id":  object.VpcCloudInstanceId,
			"vpc_id":                 object.VpcId,
			"zone_id":                object.ZoneId,
			"availability_zone":      object.ZoneId,
			"region_id":              object.RegionId,
			"create_time":            object.CreateTime,
			"user_name":              object.UserName,
			"connection_domain":      object.ConnectionDomain,
		}
		ids = append(ids, object.InstanceId)
		tags := make(map[string]string)
		for _, t := range object.Tags.Tag {
			if !ignoredTags(t.Key, t.Value) {
				tags[t.Key] = t.Value
			}
		}
		mapping["tags"] = tags
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			names = append(names, object.InstanceName)
			s = append(s, mapping)
			continue
		}

		request := r_kvstore.CreateDescribeInstanceAttributeRequest()
		request.RegionId = client.RegionId
		request.InstanceId = object.InstanceId
		raw, err := client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.DescribeInstanceAttribute(request)
		})
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		responseGet, _ := raw.(*r_kvstore.DescribeInstanceAttributeResponse)
		if len(responseGet.Instances.DBInstanceAttribute) > 0 {
			mapping["instance_release_protection"] = responseGet.Instances.DBInstanceAttribute[0].InstanceReleaseProtection
			mapping["maintain_end_time"] = responseGet.Instances.DBInstanceAttribute[0].MaintainEndTime
			mapping["maintain_start_time"] = responseGet.Instances.DBInstanceAttribute[0].MaintainStartTime
			mapping["vpc_auth_mode"] = responseGet.Instances.DBInstanceAttribute[0].VpcAuthMode
			mapping["secondary_zone_id"] = responseGet.Instances.DBInstanceAttribute[0].SecondaryZoneId
		}

		request1 := r_kvstore.CreateDescribeInstanceAutoRenewalAttributeRequest()
		request1.RegionId = client.RegionId
		request1.DBInstanceId = object.InstanceId
		request1.ClientToken = buildClientToken(request1.GetActionName())
		raw1, err := client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.DescribeInstanceAutoRenewalAttribute(request1)
		})
		if err != nil {
			if !IsExpectedErrors(err, []string{"InvalidOrderCharge.NotSupport"}) {
				return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_kvstore_instances", request1.GetActionName(), AlibabaCloudSdkGoERROR)
			}
		}
		addDebug(request1.GetActionName(), raw1, request1.RpcRequest, request1)
		responseGet1, _ := raw1.(*r_kvstore.DescribeInstanceAutoRenewalAttributeResponse)
		if len(responseGet1.Items.Item) > 0 {
			mapping["auto_renew"] = responseGet1.Items.Item[0].AutoRenew
			mapping["auto_renew_period"] = responseGet1.Items.Item[0].Duration
		}

		request2 := r_kvstore.CreateDescribeInstanceSSLRequest()
		request2.RegionId = client.RegionId
		request2.InstanceId = object.InstanceId
		raw2, err := client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.DescribeInstanceSSL(request2)
		})
		addDebug(request2.GetActionName(), raw2, request2.RpcRequest, request2)
		responseGet2, _ := raw2.(*r_kvstore.DescribeInstanceSSLResponse)
		mapping["ssl_enable"] = responseGet2.SSLEnabled

		request3 := r_kvstore.CreateDescribeSecurityGroupConfigurationRequest()
		request3.RegionId = client.RegionId
		request3.InstanceId = object.InstanceId
		raw3, err := client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.DescribeSecurityGroupConfiguration(request3)
		})
		addDebug(request3.GetActionName(), raw3, request3.RpcRequest, request3)
		responseGet3, _ := raw3.(*r_kvstore.DescribeSecurityGroupConfigurationResponse)
		if len(responseGet3.Items.EcsSecurityGroupRelation) > 0 {
			mapping["security_group_id"] = responseGet3.Items.EcsSecurityGroupRelation[0].SecurityGroupId
		}

		request4 := r_kvstore.CreateDescribeSecurityIpsRequest()
		request4.RegionId = client.RegionId
		request4.InstanceId = object.InstanceId
		raw4, err := client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.DescribeSecurityIps(request4)
		})
		addDebug(request4.GetActionName(), raw4, request4.RpcRequest, request4)
		responseGet4, _ := raw4.(*r_kvstore.DescribeSecurityIpsResponse)
		if len(responseGet4.SecurityIpGroups.SecurityIpGroup) > 0 {
			mapping["security_ip_group_attribute"] = responseGet4.SecurityIpGroups.SecurityIpGroup[0].SecurityIpGroupAttribute
			mapping["security_ip_group_name"] = responseGet4.SecurityIpGroups.SecurityIpGroup[0].SecurityIpGroupName
			mapping["security_ips"] = strings.Split(responseGet4.SecurityIpGroups.SecurityIpGroup[0].SecurityIpList, ",")
		}

		names = append(names, object.InstanceName)
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
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
