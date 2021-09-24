package alicloud

import (
	"regexp"
	"sort"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCrEEInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCrEEInstancesRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			// Computed values
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
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
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"specification": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace_quota": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace_usage": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repo_quota": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repo_usage": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_endpoints": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"public_endpoints": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"authorization_token": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"temp_username": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCrEEInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := &CrService{client}
	pageNo := 1
	pageSize := 50

	var instances []cr_ee.InstancesItem
	for {
		resp, err := crService.ListCrEEInstances(pageNo, pageSize)
		if err != nil {
			return WrapError(err)
		}
		instances = append(instances, resp.Instances...)
		if len(resp.Instances) < pageSize {
			break
		}
		pageNo++
	}

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	var idsMap map[string]string
	if v, ok := d.GetOk("ids"); ok {
		idsMap = make(map[string]string)
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var targetInstances []cr_ee.InstancesItem
	for _, instance := range instances {
		if nameRegex != nil && !nameRegex.MatchString(instance.InstanceName) {
			continue
		}

		if idsMap != nil && idsMap[instance.InstanceId] == "" {
			continue
		}

		targetInstances = append(targetInstances, instance)
	}

	instances = targetInstances

	sort.SliceStable(instances, func(i, j int) bool {
		return instances[i].CreateTime < instances[j].CreateTime
	})

	var (
		ids          []string
		names        []string
		instanceMaps []map[string]interface{}
	)

	for _, instance := range instances {
		usageResp, err := crService.GetCrEEInstanceUsage(instance.InstanceId)
		if err != nil {
			return WrapError(err)
		}
		endpointResp, err := crService.ListCrEEInstanceEndpoint(instance.InstanceId)
		if err != nil {
			return WrapError(err)
		}

		var (
			publicDomains []string
			vpcDomains    []string
		)
		for _, endpoint := range endpointResp.Endpoints {
			if !endpoint.Enable {
				continue
			}
			if endpoint.EndpointType == "internet" {
				for _, d := range endpoint.Domains {
					publicDomains = append(publicDomains, d.Domain)
				}
			} else if endpoint.EndpointType == "vpc" {
				for _, d := range endpoint.Domains {
					vpcDomains = append(vpcDomains, d.Domain)
				}
			}
		}

		mapping := make(map[string]interface{})
		mapping["id"] = instance.InstanceId
		mapping["name"] = instance.InstanceName
		mapping["region"] = instance.RegionId
		mapping["specification"] = instance.InstanceSpecification
		mapping["namespace_quota"] = usageResp.NamespaceQuota
		mapping["namespace_usage"] = usageResp.NamespaceUsage
		mapping["repo_quota"] = usageResp.RepoQuota
		mapping["repo_usage"] = usageResp.RepoUsage
		mapping["vpc_endpoints"] = vpcDomains
		mapping["public_endpoints"] = publicDomains

		ids = append(ids, instance.InstanceId)
		names = append(names, instance.InstanceName)

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			instanceMaps = append(instanceMaps, mapping)
			continue
		}

		response := &cr_ee.GetAuthorizationTokenResponse{}
		request := cr_ee.CreateGetAuthorizationTokenRequest()
		request.InstanceId = instance.InstanceId
		action := request.GetActionName()

		raw, err := client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
			return creeClient.GetAuthorizationToken(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "data_alicloud_cr_ee_instances", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, raw, request.RpcRequest, request)

		response, _ = raw.(*cr_ee.GetAuthorizationTokenResponse)
		if !response.GetAuthorizationTokenIsSuccess {
			return WrapErrorf(err, DefaultErrorMsg, "data_alicloud_cr_ee_instances", action, AlibabaCloudSdkGoERROR)
		}
		mapping["authorization_token"] = response.AuthorizationToken
		mapping["temp_username"] = response.TempUsername

		instanceMaps = append(instanceMaps, mapping)
	}

	d.SetId(dataResourceIdHash(names))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("instances", instanceMaps); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok {
		writeToFile(output.(string), instanceMaps)
	}

	return nil
}
