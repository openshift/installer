package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudRamRoles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRamRolesRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"policy_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 128),
			},
			"policy_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"System", "Custom"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
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
			"roles": {
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
						"arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"assume_role_policy_document": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"document": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudRamRolesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	allRoles := []interface{}{}

	allRolesMap := make(map[string]interface{})

	dataMap := []interface{}{}

	policyName, policyNameOk := d.GetOk("policy_name")
	policyType, policyTypeOk := d.GetOk("policy_type")
	nameRegex, nameRegexOk := d.GetOk("name_regex")
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	// all roles

	request := ram.CreateListRolesRequest()
	request.RegionId = client.RegionId
	request.MaxItems = requests.NewInteger(1000)
	for {
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListRoles(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_roles", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		response, _ := raw.(*ram.ListRolesResponse)
		for _, v := range response.Roles.Role {
			if nameRegexOk {
				r, err := regexp.Compile(nameRegex.(string))
				if err != nil {
					return WrapError(err)
				}
				if !r.MatchString(v.RoleName) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[v.RoleId]; !ok {
					continue
				}
			}
			allRolesMap[v.RoleName] = v
			allRoles = append(allRoles, v)
		}
		if !response.IsTruncated {
			break
		}
		request.Marker = response.Marker
	}

	// roles which attach with this policy
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_roles", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ram.ListEntitiesForPolicyResponse)
		for _, v := range response.Roles.Role {
			role, ok := allRolesMap[v.RoleName]
			if ok {
				dataMap = append(dataMap, role)
			}
		}
		allRoles = dataMap
	}
	return ramRolesDescriptionAttributes(d, meta, allRoles)
}

func ramRolesDescriptionAttributes(d *schema.ResourceData, meta interface{}, roles []interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, v := range roles {
		role := v.(ram.RoleInListRoles)
		request := ram.CreateGetRoleRequest()
		request.RoleName = role.RoleName
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.GetRole(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"EntityNotExist"}) {
				continue
			}
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_roles", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ram.GetRoleResponse)
		mapping := map[string]interface{}{
			"id":                          role.RoleId,
			"name":                        role.RoleName,
			"arn":                         role.Arn,
			"description":                 response.Role.Description,
			"create_date":                 role.CreateDate,
			"update_date":                 role.UpdateDate,
			"assume_role_policy_document": response.Role.AssumeRolePolicyDocument,
			"document":                    response.Role.AssumeRolePolicyDocument,
		}
		ids = append(ids, role.RoleId)
		names = append(names, role.RoleName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("roles", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
