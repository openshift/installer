package alicloud

import (
	cs "github.com/alibabacloud-go/cs-20151215/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCSKubernetesPermissions() *schema.Resource {
	return &schema.Resource{
		Read: dataAlicloudCSKubernetesPermissionsRead,

		Schema: map[string]*schema.Schema{
			"uid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"permissions": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"resource_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"cluster", "namespace", "console"}, false),
						},
						"role_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"role_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"is_owner": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"is_ram_role": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func dataAlicloudCSKubernetesPermissionsRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "InitializeClient", err)
	}

	// Query existing permissions, DescribeUserPermission
	uid := d.Get("uid").(string)
	perms, _err := describeUserPermissions(client, uid)
	if _err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "DescribeUserPermission", err)
	}

	_ = d.Set("permissions", flattenPermissionsConfig(perms))
	_ = d.Set("uid", uid)

	d.SetId(tea.ToString(hashcode.String(uid)))
	return nil
}

func describeUserPermissions(client *cs.Client, uid string) ([]*cs.DescribeUserPermissionResponseBody, error) {
	resp, err := client.DescribeUserPermission(tea.String(uid))
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func flattenPermissionsConfig(permissions []*cs.DescribeUserPermissionResponseBody) (m []map[string]interface{}) {
	if permissions == nil {
		return []map[string]interface{}{}
	}
	for _, permission := range permissions {
		m = append(m, map[string]interface{}{
			"resource_id":   permission.ResourceId,
			"resource_type": permission.ResourceType,
			"role_name":     permission.RoleName,
			"role_type":     permission.RoleType,
			"is_owner":      convertToBool(permission.IsOwner),
			"is_ram_role":   convertToBool(permission.IsRamRole),
		})
	}

	return m
}

func convertToBool(i *int64) bool {
	in := tea.Int64Value(i)
	var b bool
	if in == 0 {
		b = false
	}
	if in == 1 {
		b = true
	}
	return b
}
