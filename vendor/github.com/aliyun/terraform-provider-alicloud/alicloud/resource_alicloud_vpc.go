package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudVpcCreate,
		Read:   resourceAlicloudVpcRead,
		Update: resourceAlicloudVpcUpdate,
		Delete: resourceAlicloudVpcDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cidr_block": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "172.16.0.0/12",
				ValidateFunc:  validateCIDRNetworkAddress,
				ConflictsWith: []string{"enable_ipv6"},
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"enable_ipv6": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"cidr_block"},
			},
			"ipv6_cidr_block": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"router_table_id": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Attribute router_table_id has been deprecated and replaced with route_table_id.",
			},
			"router_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secondary_cidr_blocks": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"user_cidrs": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ForceNew: true,
			},
			"vpc_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
				ValidateFunc:  validateNormalName,
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'name' has been deprecated from provider version 1.119.0. New field 'vpc_name' instead.",
				ConflictsWith: []string{"vpc_name"},
				ValidateFunc:  validateNormalName,
			},
		},
	}
}

func resourceAlicloudVpcCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	action := "CreateVpc"
	request := make(map[string]interface{})
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("cidr_block"); ok {
		request["CidrBlock"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	if v, ok := d.GetOkExists("enable_ipv6"); ok {
		request["EnableIpv6"] = v
	}

	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("user_cidrs"); ok && v != nil {
		request["UserCidr"] = convertListToCommaSeparate(v.([]interface{}))
	}

	if v, ok := d.GetOk("vpc_name"); ok {
		request["VpcName"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["VpcName"] = v
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("CreateVpc")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"TaskConflict", "Throttling", "UnknownError"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["VpcId"]))
	stateConf := BuildStateConf([]string{"Pending"}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcService.VpcStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudVpcUpdate(d, meta)
}
func resourceAlicloudVpcRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeVpc(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_vpc vpcService.DescribeVpc Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("cidr_block", object["CidrBlock"])
	d.Set("description", object["Description"])
	d.Set("ipv6_cidr_block", object["Ipv6CidrBlock"])
	d.Set("router_id", object["VRouterId"])
	d.Set("secondary_cidr_blocks", object["SecondaryCidrBlocks"].(map[string]interface{})["SecondaryCidrBlock"])
	d.Set("status", object["Status"])
	if v, ok := object["Tags"].(map[string]interface{}); ok {
		d.Set("tags", tagsToMap(v["Tag"]))
	}
	d.Set("user_cidrs", object["UserCidrs"].(map[string]interface{})["UserCidr"])
	d.Set("vpc_name", object["VpcName"])
	d.Set("name", object["VpcName"])

	describeRouteTableListObject, err := vpcService.DescribeRouteTableList(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("route_table_id", describeRouteTableListObject["RouteTableId"])
	d.Set("router_table_id", describeRouteTableListObject["RouteTableId"])
	return nil
}
func resourceAlicloudVpcUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	d.Partial(true)
	if err := vpcService.setInstanceSecondaryCidrBlocks(d); err != nil {
		return WrapError(err)
	}

	if d.HasChange("tags") {
		if err := vpcService.SetResourceTags(d, "vpc"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	update := false
	moveResourceGroupReq := map[string]interface{}{
		"ResourceId": d.Id(),
	}
	moveResourceGroupReq["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	moveResourceGroupReq["NewResourceGroupId"] = d.Get("resource_group_id")
	moveResourceGroupReq["ResourceType"] = "vpc"
	if update {
		action := "MoveResourceGroup"
		conn, err := client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, moveResourceGroupReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, moveResourceGroupReq)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("resource_group_id")
	}
	update = false
	modifyVpcAttributeReq := map[string]interface{}{
		"VpcId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("cidr_block") {
		update = true
		modifyVpcAttributeReq["CidrBlock"] = d.Get("cidr_block")
	}
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		modifyVpcAttributeReq["Description"] = d.Get("description")
	}

	modifyVpcAttributeReq["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("vpc_name") {
		update = true
		modifyVpcAttributeReq["VpcName"] = d.Get("vpc_name")
	}
	if !d.IsNewResource() && d.HasChange("name") {
		update = true
		modifyVpcAttributeReq["VpcName"] = d.Get("name")
	}
	if update {
		if _, ok := d.GetOkExists("enable_ipv6"); ok {
			modifyVpcAttributeReq["EnableIPv6"] = d.Get("enable_ipv6")
		}
		action := "ModifyVpcAttribute"
		conn, err := client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, modifyVpcAttributeReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, modifyVpcAttributeReq)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("cidr_block")
		d.SetPartial("description")
		d.SetPartial("name")
		d.SetPartial("vpc_name")
	}
	d.Partial(false)
	return resourceAlicloudVpcRead(d, meta)
}
func resourceAlicloudVpcDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	action := "DeleteVpc"
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"VpcId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"Forbidden.VpcNotFound", "InvalidVpcID.NotFound", "InvalidVpcId.NotFound"}) {
				return nil
			}
			wait()
			return resource.RetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{"Pending"}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcService.VpcStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
