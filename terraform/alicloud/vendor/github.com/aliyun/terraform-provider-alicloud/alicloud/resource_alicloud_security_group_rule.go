package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliyunSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSecurityGroupRuleCreate,
		Read:   resourceAliyunSecurityGroupRuleRead,
		Update: resourceAliyunSecurityGroupRuleUpdate,
		Delete: resourceAliyunSecurityGroupRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ingress", "egress"}, false),
				Description:  "Type of rule, ingress (inbound) or egress (outbound).",
			},

			"ip_protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"tcp", "udp", "icmp", "gre", "all"}, false),
			},

			"nic_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"internet", "intranet"}, false),
			},

			"policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      GroupRulePolicyAccept,
				ValidateFunc: validation.StringInSlice([]string{"accept", "drop"}, false),
			},

			"port_range": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          AllPortRange,
				DiffSuppressFunc: ecsSecurityGroupRulePortRangeDiffSuppressFunc,
			},

			"priority": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 100),
			},

			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"cidr_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"source_security_group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"cidr_ip"},
			},

			"source_group_owner_account": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliyunSecurityGroupRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	direction := d.Get("type").(string)
	sgId := d.Get("security_group_id").(string)
	ptl := d.Get("ip_protocol").(string)
	port := d.Get("port_range").(string)
	if port == "" {
		return WrapError(fmt.Errorf("'port_range': required field is not set or invalid."))
	}
	nicType := d.Get("nic_type").(string)
	policy := d.Get("policy").(string)
	priority := d.Get("priority").(int)

	if _, ok := d.GetOk("cidr_ip"); !ok {
		if _, ok := d.GetOk("source_security_group_id"); !ok {
			return WrapError(fmt.Errorf("Either 'cidr_ip' or 'source_security_group_id' must be specified."))
		}
	}

	request, err := buildAliyunSGRuleRequest(d, meta)
	if err != nil {
		return WrapError(err)
	}

	if direction == string(DirectionIngress) {
		request.ApiName = "AuthorizeSecurityGroup"
		_, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
	} else {
		request.ApiName = "AuthorizeSecurityGroupEgress"
		_, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
	}
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_security_group_rule", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	var cidr_ip string
	if ip, ok := d.GetOk("cidr_ip"); ok {
		cidr_ip = ip.(string)
	} else {
		cidr_ip = d.Get("source_security_group_id").(string)
	}
	d.SetId(sgId + ":" + direction + ":" + ptl + ":" + port + ":" + nicType + ":" + cidr_ip + ":" + policy + ":" + strconv.Itoa(priority))

	return resourceAliyunSecurityGroupRuleRead(d, meta)
}

func resourceAliyunSecurityGroupRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	parts := strings.Split(d.Id(), ":")
	policy := parseSecurityRuleId(d, meta, 6)
	strPriority := parseSecurityRuleId(d, meta, 7)
	var priority int
	if policy == "" || strPriority == "" {
		policy = d.Get("policy").(string)
		priority = d.Get("priority").(int)
		d.SetId(d.Id() + ":" + policy + ":" + strconv.Itoa(priority))
	} else {
		prior, err := strconv.Atoi(strPriority)
		if err != nil {
			return WrapError(err)
		}
		priority = prior
	}
	sgId := parts[0]
	direction := parts[1]

	// wait the rule exist
	var object ecs.Permission
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err := resource.Retry(10*time.Minute, func() *resource.RetryError {
		obj, err := ecsService.DescribeSecurityGroupRule(d.Id())
		if err != nil && d.IsNewResource() {
			wait()
			return resource.RetryableError(err)
		} else {
			object = obj
			return resource.NonRetryableError(err)
		}
	})
	if err != nil {
		if NotFoundError(err) && !d.IsNewResource() {
			log.Printf("[DEBUG] Resource alicloud_security_group_rule ecsService.DescribeSecurityGroupRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("type", object.Direction)
	d.Set("ip_protocol", strings.ToLower(string(object.IpProtocol)))
	d.Set("nic_type", object.NicType)
	d.Set("policy", strings.ToLower(string(object.Policy)))
	d.Set("port_range", object.PortRange)
	d.Set("description", object.Description)
	if pri, err := strconv.Atoi(object.Priority); err != nil {
		return WrapError(err)
	} else {
		d.Set("priority", pri)
	}
	d.Set("security_group_id", sgId)
	//support source and desc by type
	if direction == string(DirectionIngress) {
		d.Set("cidr_ip", object.SourceCidrIp)
		d.Set("source_security_group_id", object.SourceGroupId)
		d.Set("source_group_owner_account", object.SourceGroupOwnerAccount)
	} else {
		d.Set("cidr_ip", object.DestCidrIp)
		d.Set("source_security_group_id", object.DestGroupId)
		d.Set("source_group_owner_account", object.DestGroupOwnerAccount)
	}
	return nil
}

func resourceAliyunSecurityGroupRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	policy := parseSecurityRuleId(d, meta, 6)
	strPriority := parseSecurityRuleId(d, meta, 7)
	var priority int
	if policy == "" || strPriority == "" {
		policy = d.Get("policy").(string)
		priority = d.Get("priority").(int)
		d.SetId(d.Id() + ":" + policy + ":" + strconv.Itoa(priority))
	} else {
		prior, err := strconv.Atoi(strPriority)
		if err != nil {
			return WrapError(err)
		}
		priority = prior
	}

	request, err := buildAliyunSGRuleRequest(d, meta)
	if err != nil {
		return WrapError(err)
	}

	direction := d.Get("type").(string)

	if direction == string(DirectionIngress) {
		request.ApiName = "ModifySecurityGroupRule"
		_, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
	} else {
		request.ApiName = "ModifySecurityGroupEgressRule"
		_, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
	}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.Headers, request)
	return resourceAliyunSecurityGroupRuleRead(d, meta)
}

func deleteSecurityGroupRule(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ruleType := d.Get("type").(string)
	request, err := buildAliyunSGRuleRequest(d, meta)
	if err != nil {
		return WrapError(err)
	}

	if ruleType == string(DirectionIngress) {
		request.ApiName = "RevokeSecurityGroup"
		_, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
	} else {
		request.ApiName = "RevokeSecurityGroupEgress"
		_, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
	}

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}

func resourceAliyunSecurityGroupRuleDelete(d *schema.ResourceData, meta interface{}) error {
	policy := parseSecurityRuleId(d, meta, 6)
	strPriority := parseSecurityRuleId(d, meta, 7)
	var priority int
	if policy == "" || strPriority == "" {
		policy = d.Get("policy").(string)
		priority = d.Get("priority").(int)
		d.SetId(d.Id() + ":" + policy + ":" + strconv.Itoa(priority))
	} else {
		prior, err := strconv.Atoi(strPriority)
		if err != nil {
			return WrapError(err)
		}
		priority = prior
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		err := deleteSecurityGroupRule(d, meta)
		if err != nil {
			if NotFoundError(err) || IsExpectedErrors(err, []string{"InvalidSecurityGroupId.NotFound"}) {
				return nil
			}
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapError(err)
	}
	return nil
}

func buildAliyunSGRuleRequest(d *schema.ResourceData, meta interface{}) (*requests.CommonRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	// Get product code from the built request
	ruleReq := ecs.CreateModifySecurityGroupRuleRequest()
	request, err := client.NewCommonRequest(ruleReq.GetProduct(), ruleReq.GetLocationServiceCode(), strings.ToUpper(string(Https)), connectivity.ApiVersion20140526)
	if err != nil {
		return request, WrapError(err)
	}

	direction := d.Get("type").(string)

	port_range := d.Get("port_range").(string)
	request.QueryParams["PortRange"] = port_range

	if v, ok := d.GetOk("ip_protocol"); ok {
		request.QueryParams["IpProtocol"] = v.(string)
		if v.(string) == string(Tcp) || v.(string) == string(Udp) {
			if port_range == AllPortRange {
				return nil, fmt.Errorf("'tcp' and 'udp' can support port range: [1, 65535]. Please correct it and try again.")
			}
		} else if port_range != AllPortRange {
			return nil, fmt.Errorf("'icmp', 'gre' and 'all' only support port range '-1/-1'. Please correct it and try again.")
		}
	}

	if v, ok := d.GetOk("policy"); ok {
		request.QueryParams["Policy"] = v.(string)
	}

	if v, ok := d.GetOk("priority"); ok {
		request.QueryParams["Priority"] = strconv.Itoa(v.(int))
	}

	if v, ok := d.GetOk("cidr_ip"); ok {
		if direction == string(DirectionIngress) {
			request.QueryParams["SourceCidrIp"] = v.(string)
		} else {
			request.QueryParams["DestCidrIp"] = v.(string)
		}
	}

	var targetGroupId string
	if v, ok := d.GetOk("source_security_group_id"); ok {
		targetGroupId = v.(string)
		if direction == string(DirectionIngress) {
			request.QueryParams["SourceGroupId"] = targetGroupId
		} else {
			request.QueryParams["DestGroupId"] = targetGroupId
		}
	}

	if v, ok := d.GetOk("source_group_owner_account"); ok {
		if direction == string(DirectionIngress) {
			request.QueryParams["SourceGroupOwnerAccount"] = v.(string)
		} else {
			request.QueryParams["DestGroupOwnerAccount"] = v.(string)
		}
	}

	sgId := d.Get("security_group_id").(string)

	group, err := ecsService.DescribeSecurityGroup(sgId)
	if err != nil {
		return nil, WrapError(err)
	}

	if v, ok := d.GetOk("nic_type"); ok {
		if group.VpcId != "" || targetGroupId != "" {
			if GroupRuleNicType(v.(string)) != GroupRuleIntranet {
				return nil, fmt.Errorf("When security group in the vpc or authorizing permission for source/destination security group, " +
					"the nic_type must be 'intranet'.")
			}
		}
		request.QueryParams["NicType"] = v.(string)
	}

	request.QueryParams["SecurityGroupId"] = sgId

	description := d.Get("description").(string)
	request.QueryParams["Description"] = description

	return request, nil
}

func parseSecurityRuleId(d *schema.ResourceData, meta interface{}, index int) (result string) {
	parts := strings.Split(d.Id(), ":")
	defer func() {
		if e := recover(); e != nil {
			log.Printf("Panicing %s\r\n", e)
			result = ""
		}
	}()
	return parts[index]
}
