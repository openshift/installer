package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type SecurityGroup struct {
	Attributes        ecs.DescribeSecurityGroupAttributeResponse
	CreationTime      string
	SecurityGroupType string
	ResourceGroupId   string
	Tags              ecs.TagsInDescribeSecurityGroups
}

func dataSourceAlicloudSecurityGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSecurityGroupsRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Computed values
			"groups": {
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"inner_access": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": tagsSchema(),
					},
				},
			},
		},
	}
}

func dataSourceAlicloudSecurityGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	request := ecs.CreateDescribeSecurityGroupsRequest()
	request.RegionId = client.RegionId
	request.VpcId = d.Get("vpc_id").(string)
	request.PageNumber = requests.NewInteger(1)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.ResourceGroupId = d.Get("resource_group_id").(string)
	var sg []SecurityGroup
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		if r, err := regexp.Compile(v.(string)); err == nil {
			nameRegex = r
		}
	}
	if v, ok := d.GetOk("tags"); ok {
		var tags []ecs.DescribeSecurityGroupsTag

		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, ecs.DescribeSecurityGroupsTag{
				Key:   key,
				Value: value.(string),
			})
		}
		request.Tag = &tags
	}

	// ids
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeSecurityGroups(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "security_groups", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ecs.DescribeSecurityGroupsResponse)
		if len(response.SecurityGroups.SecurityGroup) < 1 {
			break
		}

		for _, item := range response.SecurityGroups.SecurityGroup {
			if nameRegex != nil {
				if !nameRegex.MatchString(item.SecurityGroupName) {
					continue
				}
			}

			if len(idsMap) > 0 {
				if _, ok := idsMap[item.SecurityGroupId]; !ok {
					continue
				}
			}

			attr, err := ecsService.DescribeSecurityGroup(item.SecurityGroupId)
			if err != nil {
				return WrapError(err)
			}

			sg = append(sg,
				SecurityGroup{
					Attributes:        attr,
					CreationTime:      item.CreationTime,
					SecurityGroupType: item.SecurityGroupType,
					ResourceGroupId:   item.ResourceGroupId,
					Tags:              item.Tags,
				},
			)
		}

		if len(response.SecurityGroups.SecurityGroup) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return WrapError(err)
		} else {
			request.PageNumber = page
		}
	}
	return securityGroupsDescription(d, sg, meta)
}

func securityGroupsDescription(d *schema.ResourceData, sg []SecurityGroup, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	var ids []string
	var names []string
	var s []map[string]interface{}

	for _, item := range sg {
		mapping := map[string]interface{}{
			"id":                  item.Attributes.SecurityGroupId,
			"name":                item.Attributes.SecurityGroupName,
			"description":         item.Attributes.Description,
			"vpc_id":              item.Attributes.VpcId,
			"resource_group_id":   item.ResourceGroupId,
			"security_group_type": item.SecurityGroupType,
			"inner_access":        item.Attributes.InnerAccessPolicy == string(GroupInnerAccept),
			"creation_time":       item.CreationTime,
			"tags":                ecsService.tagsToMap(item.Tags.Tag),
		}

		ids = append(ids, string(item.Attributes.SecurityGroupId))
		names = append(names, item.Attributes.SecurityGroupName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("groups", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
