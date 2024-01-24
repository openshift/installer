package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudRamGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRamGroupsRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"user_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 64),
			},
			"policy_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 128),
			},
			"policy_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				// must be ram.System, ram.Custom
				ValidateFunc: validation.StringInSlice([]string{"System", "Custom"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
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
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comments": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudRamGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}
	allGroups := []interface{}{}

	allGroupsMap := make(map[string]interface{})
	userFilterGroupsMap := make(map[string]interface{})
	policyFilterGroupsMap := make(map[string]interface{})

	dataMap := []map[string]interface{}{}

	userName, userNameOk := d.GetOk("user_name")
	policyName, policyNameOk := d.GetOk("policy_name")
	policyType, policyTypeOk := d.GetOk("policy_type")
	nameRegex, nameRegexOk := d.GetOk("name_regex")

	if policyTypeOk && !policyNameOk {
		return WrapError(Error("You must set 'policy_name' at one time when you set 'policy_type'."))
	}

	// groups filtered by name_regex
	request := ram.CreateListGroupsRequest()
	request.RegionId = client.RegionId
	request.MaxItems = requests.NewInteger(1000)
	for {
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListGroups(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_groups", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ram.ListGroupsResponse)
		for _, v := range response.Groups.Group {
			if nameRegexOk {
				r, err := regexp.Compile(nameRegex.(string))
				if err != nil {
					return WrapError(err)
				}
				if !r.MatchString(v.GroupName) {
					continue
				}
			}
			allGroupsMap[v.GroupName] = v
		}
		if !response.IsTruncated {
			break
		}
		request.Marker = response.Marker
	}

	// groups for user
	if userNameOk {
		request := ram.CreateListGroupsForUserRequest()
		request.UserName = userName.(string)
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListGroupsForUser(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_groups", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ram.ListGroupsForUserResponse)
		for _, v := range response.Groups.Group {
			userFilterGroupsMap[v.GroupName] = v
		}
		dataMap = append(dataMap, userFilterGroupsMap)
	}

	// groups which attach with this policy
	if policyNameOk {
		pType := "System"
		if policyTypeOk {
			pType = policyType.(string)
		}
		request := ram.CreateListEntitiesForPolicyRequest()
		request.PolicyType = pType
		request.PolicyName = policyName.(string)
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListEntitiesForPolicy(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_groups", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ram.ListEntitiesForPolicyResponse)
		for _, v := range response.Groups.Group {
			policyFilterGroupsMap[v.GroupName] = v
		}
		dataMap = append(dataMap, policyFilterGroupsMap)
	}

	// GetIntersection of each map
	allGroups = ramService.GetIntersection(dataMap, allGroupsMap)

	return ramGroupsDescriptionAttributes(d, allGroups)
}

func ramGroupsDescriptionAttributes(d *schema.ResourceData, groups []interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for _, v := range groups {
		group := v.(ram.GroupInListGroups)
		mapping := map[string]interface{}{
			"name":     group.GroupName,
			"comments": group.Comments,
		}
		ids = append(ids, v.(ram.GroupInListGroups).GroupName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("groups", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", ids); err != nil {
		return WrapError(err)
	}
	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
