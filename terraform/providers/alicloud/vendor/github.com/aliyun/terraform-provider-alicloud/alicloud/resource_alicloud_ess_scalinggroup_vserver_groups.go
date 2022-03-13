package alicloud

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEssScalingGroupVserverGroups() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEssVserverGroupsCreate,
		Read:   resourceAliyunEssVserverGroupsRead,
		Update: resourceAliyunEssVserverGroupsUpdate,
		Delete: resourceAliyunEssVserverGroupsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"vserver_groups": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"loadbalancer_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"vserver_attributes": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vserver_group_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"weight": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
							Set: func(v interface{}) int {
								var buf bytes.Buffer
								m := v.(map[string]interface{})
								if v, ok := m["vserver_group_id"]; ok {
									buf.WriteString(fmt.Sprintf("%s-", v.(string)))
								}
								if v, ok := m["port"]; ok {
									buf.WriteString(fmt.Sprintf("%d-", v.(int)))
								}
								return hashcode.String(buf.String())
							},
						},
					},
				},
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				DiffSuppressFunc: func(k, old string, new string, d *schema.ResourceData) bool {
					return old == "" && new == "true" && d.Id() != ""
				},
			},
		},
	}
}

func resourceAliyunEssVserverGroupsCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(d.Get("scaling_group_id").(string))
	return resourceAliyunEssVserverGroupsUpdate(d, meta)
}

func resourceAliyunEssVserverGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	object, err := essService.DescribeEssScalingGroup(d.Id())
	if err != nil {
		return WrapError(err)
	}
	err = d.Set("scaling_group_id", object.ScalingGroupId)
	if err != nil {
		return WrapError(err)
	}
	err = d.Set("vserver_groups", essService.flattenVserverGroupList(object.VServerGroups.VServerGroup))
	if err != nil {
		return WrapError(err)
	}
	return nil
}

func resourceAliyunEssVserverGroupsUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	object, err := essService.DescribeEssScalingGroup(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Partial(true)
	vserverGroupsMapFromScalingGroup := vserverGroupMapFromScalingGroup(object.VServerGroups.VServerGroup)
	vserverGroupsMapFromConfig := vserverGroupMapFromConfig(d.Get("vserver_groups").(*schema.Set))
	err = vserverGroupMapPreCheckLb(vserverGroupsMapFromConfig, meta.(*connectivity.AliyunClient))
	if err != nil {
		return WrapError(err)
	}
	attachMap, detachMap := attachOrDetachVserverGroupMap(vserverGroupsMapFromConfig, vserverGroupsMapFromScalingGroup)
	v, ok := d.GetOkExists("force")
	force := true
	if ok {
		force = v.(bool)
	}
	err = detachVserverGroups(d, client, detachMap, force)
	if err != nil {
		return WrapError(err)
	}
	d.SetPartial("vserver_groups")

	err = attachVserverGroups(d, client, attachMap, force)
	if err != nil {
		return WrapError(err)
	}
	d.SetPartial("vserver_groups")
	d.Partial(false)
	return resourceAliyunEssVserverGroupsRead(d, meta)
}

func resourceAliyunEssVserverGroupsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vserverGroupsFromConfig := vserverGroupMapFromConfig(d.Get("vserver_groups").(*schema.Set))
	_, detachMap := attachOrDetachVserverGroupMap(make(map[string]string, 0), vserverGroupsFromConfig)
	v, ok := d.GetOkExists("force")
	force := true
	if ok {
		force = v.(bool)
	}
	err := detachVserverGroups(d, client, detachMap, force)
	if err != nil {
		return WrapError(err)
	}
	return nil
}

func vserverGroupMapFromScalingGroup(vServerGroups []ess.VServerGroup) map[string]string {
	vserverGroupMap := make(map[string]string)
	if vServerGroups != nil && len(vServerGroups) > 0 {
		for _, v := range vServerGroups {
			vserverGroupAttributes := v.VServerGroupAttributes.VServerGroupAttribute
			for _, a := range vserverGroupAttributes {
				key := fmt.Sprintf("%s_%s_%d_%d", v.LoadBalancerId, a.VServerGroupId, a.Port, a.Weight)
				vserverGroupMap[key] = key
			}
		}
	}
	return vserverGroupMap
}

func vserverGroupMapFromConfig(vserverGroups *schema.Set) map[string]string {
	vserverGroupMap := make(map[string]string)
	vserverGroupList := vserverGroups.List()
	if len(vserverGroupList) > 0 {
		for _, v := range vserverGroupList {
			vserverGroup := v.(map[string]interface{})
			loadBalancerId := vserverGroup["loadbalancer_id"].(string)
			attrs := vserverGroup["vserver_attributes"].(*schema.Set).List()
			for _, e := range attrs {
				vserverAttribute := e.(map[string]interface{})
				vserverGroupId := vserverAttribute["vserver_group_id"].(string)
				port := vserverAttribute["port"].(int)
				weight := vserverAttribute["weight"].(int)
				key := fmt.Sprintf("%s_%s_%d_%d", loadBalancerId, vserverGroupId, port, weight)
				vserverGroupMap[key] = key
			}
		}
	}
	return vserverGroupMap
}

func attachOrDetachVserverGroupMap(newMap map[string]string, oldMap map[string]string) (map[string]string, map[string]string) {
	attachMap := make(map[string]string)
	detachMap := make(map[string]string)
	for k, v := range newMap {
		if _, ok := oldMap[k]; !ok {
			attachMap[k] = v
		}
	}
	for k, v := range oldMap {
		if _, ok := newMap[k]; !ok {
			detachMap[k] = v
		}
	}
	return attachMap, detachMap
}

func buildEssVserverGroupListMap(vserverGroupMap map[string]string) map[string][]string {
	vserverGroupRequestMap := make(map[string][]string, 0)
	for _, v := range vserverGroupMap {
		attrs := strings.Split(v, "_")
		loadbalancerId := attrs[0]
		if _, ok := vserverGroupRequestMap[loadbalancerId]; !ok {
			vserverGroupAttributes := make([]string, 0)
			vserverGroupAttributes = append(vserverGroupAttributes, v)
			vserverGroupRequestMap[loadbalancerId] = vserverGroupAttributes
		} else {
			vserverGroupAttributes := vserverGroupRequestMap[loadbalancerId]
			vserverGroupAttributes = append(vserverGroupAttributes, v)
			vserverGroupRequestMap[loadbalancerId] = vserverGroupAttributes
		}
	}
	return vserverGroupRequestMap
}

func attachVserverGroups(d *schema.ResourceData, client *connectivity.AliyunClient, attachMap map[string]string, force bool) error {
	if len(attachMap) > 0 {
		vserverGroupListMap := buildEssVserverGroupListMap(attachMap)
		attachScalingGroupVserverGroups := make([]ess.AttachVServerGroupsVServerGroup, 0)
		for k, v := range vserverGroupListMap {
			vserverAttributes := make([]ess.AttachVServerGroupsVServerGroupAttribute, 0)
			for _, e := range v {
				attrs := strings.Split(e, "_")
				vserverAttribute := ess.AttachVServerGroupsVServerGroupAttribute{
					VServerGroupId: attrs[1],
					Port:           attrs[2],
					Weight:         attrs[3],
				}
				vserverAttributes = append(vserverAttributes, vserverAttribute)
			}
			vserverGroup := ess.AttachVServerGroupsVServerGroup{
				LoadBalancerId:        k,
				VServerGroupAttribute: &vserverAttributes,
			}
			attachScalingGroupVserverGroups = append(attachScalingGroupVserverGroups, vserverGroup)
		}
		request := ess.CreateAttachVServerGroupsRequest()
		request.RegionId = client.RegionId
		request.ScalingGroupId = d.Id()
		request.ForceAttach = requests.NewBoolean(force)
		request.VServerGroup = &attachScalingGroupVserverGroups
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.AttachVServerGroups(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return nil
}

func detachVserverGroups(d *schema.ResourceData, client *connectivity.AliyunClient, detachMap map[string]string, force bool) error {
	if len(detachMap) > 0 {
		vserverGroupListMap := buildEssVserverGroupListMap(detachMap)
		detachScalingGroupVserverGroups := make([]ess.DetachVServerGroupsVServerGroup, 0)
		for k, v := range vserverGroupListMap {
			vserverAttributes := make([]ess.DetachVServerGroupsVServerGroupAttribute, 0)
			for _, e := range v {
				attrs := strings.Split(e, "_")
				vserverAttribute := ess.DetachVServerGroupsVServerGroupAttribute{
					VServerGroupId: attrs[1],
					Port:           attrs[2],
				}
				vserverAttributes = append(vserverAttributes, vserverAttribute)
			}
			vserverGroup := ess.DetachVServerGroupsVServerGroup{
				LoadBalancerId:        k,
				VServerGroupAttribute: &vserverAttributes,
			}
			detachScalingGroupVserverGroups = append(detachScalingGroupVserverGroups, vserverGroup)
		}
		request := ess.CreateDetachVServerGroupsRequest()
		request.RegionId = client.RegionId
		request.ScalingGroupId = d.Id()
		request.ForceDetach = requests.NewBoolean(force)
		request.VServerGroup = &detachScalingGroupVserverGroups
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DetachVServerGroups(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return nil
}

func vserverGroupMapPreCheckLb(vserverGroupsMapFromConfig map[string]string, client *connectivity.AliyunClient) error {
	slbService := SlbService{client}
	for _, v := range vserverGroupsMapFromConfig {
		attrs := strings.Split(v, "_")
		loadbalancerID := attrs[0]
		if err := slbService.WaitForSlb(loadbalancerID, Active, DefaultTimeout); err != nil {
			return WrapError(err)
		}
	}
	return nil
}
