package alicloud

import (
	"strconv"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudSecurityGroupRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSecurityGroupRulesRead,

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"nic_type": {
				Type:     schema.TypeString,
				Optional: true,
				// must be one of GroupRuleInternet, GroupRuleIntranet
				ValidateFunc: validation.StringInSlice([]string{"internet", "intranet"}, false),
			},
			"direction": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ingress", "egress"}, false),
			},
			"ip_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				// must be one of Tcp, Udp, Icmp, Gre, All
				ValidateFunc: validation.StringInSlice([]string{
					string(Tcp),
					string(Udp),
					string(Icmp),
					string(Gre),
					string(All),
				}, false),
			},
			"policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"accept", "drop"}, false),
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_cidr_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_group_owner_account": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dest_cidr_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dest_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dest_group_owner_account": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nic_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"direction": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_desc": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAlicloudSecurityGroupRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	req := ecs.CreateDescribeSecurityGroupAttributeRequest()
	req.RegionId = client.RegionId
	req.SecurityGroupId = d.Get("group_id").(string)
	req.NicType = d.Get("nic_type").(string)
	req.Direction = d.Get("direction").(string)
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeSecurityGroupAttribute(req)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "security_group_rules", req.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(req.GetActionName(), raw, req.RpcRequest, req)
	attr, _ := raw.(*ecs.DescribeSecurityGroupAttributeResponse)
	var rules []map[string]interface{}

	if attr != nil {
		for _, item := range attr.Permissions.Permission {
			if v, ok := d.GetOk("ip_protocol"); ok && strings.ToLower(string(item.IpProtocol)) != v.(string) {
				continue
			}

			if v, ok := d.GetOk("policy"); ok && strings.ToLower(string(item.Policy)) != v.(string) {
				continue
			}

			mapping := map[string]interface{}{
				"ip_protocol":                strings.ToLower(string(item.IpProtocol)),
				"port_range":                 item.PortRange,
				"source_cidr_ip":             item.SourceCidrIp,
				"source_group_id":            item.SourceGroupId,
				"source_group_owner_account": item.SourceGroupOwnerAccount,
				"dest_cidr_ip":               item.DestCidrIp,
				"dest_group_id":              item.DestGroupId,
				"dest_group_owner_account":   item.DestGroupOwnerAccount,
				"policy":                     strings.ToLower(string(item.Policy)),
				"nic_type":                   item.NicType,
				"direction":                  item.Direction,
				"description":                item.Description,
			}

			pri, err := strconv.Atoi(item.Priority)
			if err != nil {
				return WrapError(err)
			}
			mapping["priority"] = pri
			rules = append(rules, mapping)
		}

		if err := d.Set("group_name", attr.SecurityGroupName); err != nil {
			return WrapError(err)
		}

		if err := d.Set("group_desc", attr.Description); err != nil {
			return WrapError(err)
		}
	}

	d.SetId(d.Get("group_id").(string))

	if err := d.Set("rules", rules); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), rules)
	}
	return nil
}
