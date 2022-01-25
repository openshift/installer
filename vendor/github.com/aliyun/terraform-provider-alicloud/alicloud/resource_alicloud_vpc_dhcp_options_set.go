package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudVpcDhcpOptionsSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudVpcDhcpOptionsSetCreate,
		Read:   resourceAlicloudVpcDhcpOptionsSetRead,
		Update: resourceAlicloudVpcDhcpOptionsSetUpdate,
		Delete: resourceAlicloudVpcDhcpOptionsSetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"associate_vpcs": {
				Type:       schema.TypeSet,
				Optional:   true,
				Deprecated: "Field 'associate_vpcs' has been deprecated from provider version 1.153.0 and it will be removed in the future version. Please use the new resource 'alicloud_vpc_dhcp_options_set_attachment' to attach DhcpOptionsSet and Vpc.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"associate_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"dhcp_options_set_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dhcp_options_set_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[a-zA-Z\u4E00-\u9FA5][\u4E00-\u9FA5A-Za-z0-9_-]{2,128}$"), "The name must be 2 to 128 characters in length and can contain letters, Chinese characters, digits, underscores (_), and hyphens (-). It must start with a letter or a Chinese character."),
			},
			"domain_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain_name_servers": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudVpcDhcpOptionsSetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateDhcpOptionsSet"
	request := make(map[string]interface{})
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("dhcp_options_set_description"); ok {
		request["DhcpOptionsSetDescription"] = v
	}
	if v, ok := d.GetOk("dhcp_options_set_name"); ok {
		request["DhcpOptionsSetName"] = v
	}
	if v, ok := d.GetOk("domain_name"); ok {
		request["DomainName"] = v
	}
	if v, ok := d.GetOk("domain_name_servers"); ok {
		request["DomainNameServers"] = v
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["RegionId"] = client.RegionId
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("CreateDhcpOptionsSet")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_dhcp_options_set", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DhcpOptionsSetId"]))
	vpcService := VpcService{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcService.VpcDhcpOptionsSetStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudVpcDhcpOptionsSetUpdate(d, meta)
}
func resourceAlicloudVpcDhcpOptionsSetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeVpcDhcpOptionsSet(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_dhcp_options_set vpcService.DescribeVpcDhcpOptionsSet Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if associateVpcsList, ok := object["AssociateVpcs"]; ok && associateVpcsList != nil {
		associateVpcsMaps := make([]map[string]interface{}, 0)
		for _, associateVpcsListItem := range associateVpcsList.([]interface{}) {
			if associateVpcsListItemMap, ok := associateVpcsListItem.(map[string]interface{}); ok {
				associateVpcsListItemMap["associate_status"] = associateVpcsListItemMap["AssociateStatus"]
				associateVpcsListItemMap["vpc_id"] = associateVpcsListItemMap["VpcId"]
				associateVpcsMaps = append(associateVpcsMaps, associateVpcsListItemMap)
			}
		}
		d.Set("associate_vpcs", associateVpcsMaps)
	}

	d.Set("dhcp_options_set_description", object["DhcpOptionsSetDescription"])
	d.Set("dhcp_options_set_name", object["DhcpOptionsSetName"])
	d.Set("domain_name", object["DhcpOptions"].(map[string]interface{})["DomainName"])
	d.Set("domain_name_servers", object["DhcpOptions"].(map[string]interface{})["DomainNameServers"])
	d.Set("owner_id", fmt.Sprint(formatInt(object["OwnerId"])))
	d.Set("status", object["Status"])
	return nil
}
func resourceAlicloudVpcDhcpOptionsSetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"DhcpOptionsSetId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("dhcp_options_set_description") {
		update = true
		if v, ok := d.GetOk("dhcp_options_set_description"); ok {
			request["DhcpOptionsSetDescription"] = v
		}
	}
	if d.HasChange("dhcp_options_set_name") {
		update = true
		if v, ok := d.GetOk("dhcp_options_set_name"); ok {
			request["DhcpOptionsSetName"] = v
		}
	}
	if d.HasChange("domain_name") {
		update = true
		if v, ok := d.GetOk("domain_name"); ok {
			request["DomainName"] = v
		}
	}
	if d.HasChange("domain_name_servers") {
		update = true
		if v, ok := d.GetOk("domain_name_servers"); ok {
			request["DomainNameServers"] = v
		}
	}
	if update {
		if v, ok := d.GetOkExists("dry_run"); ok {
			request["DryRun"] = v
		}
		action := "UpdateDhcpOptionsSetAttribute"
		conn, err := client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			request["ClientToken"] = buildClientToken("UpdateDhcpOptionsSetAttribute")
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcService.VpcDhcpOptionsSetStateRefreshFunc(d.Id(), []string{"InUse"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("dhcp_options_set_description")
		d.SetPartial("dhcp_options_set_name")
		d.SetPartial("domain_name")
		d.SetPartial("domain_name_servers")
	}
	d.Partial(false)
	if d.HasChange("associate_vpcs") {
		oldAssociateVpcs, newAssociateVpcs := d.GetChange("associate_vpcs")
		oldAssociateVpcsSet := oldAssociateVpcs.(*schema.Set)
		newAssociateVpcsSet := newAssociateVpcs.(*schema.Set)
		removed := oldAssociateVpcsSet.Difference(newAssociateVpcsSet)
		added := newAssociateVpcsSet.Difference(oldAssociateVpcsSet)
		if removed.Len() > 0 {
			action := "DetachDhcpOptionsSetFromVpc"
			detachVpcDhcpOptionsSetRequest := map[string]interface{}{
				"DhcpOptionsSetId": d.Id(),
			}
			detachVpcDhcpOptionsSetRequest["RegionId"] = client.RegionId
			if _, ok := d.GetOkExists("dry_run"); ok {
				detachVpcDhcpOptionsSetRequest["DryRun"] = d.Get("dry_run")
			}
			for _, associateVpcs := range removed.List() {
				associateVpc := associateVpcs.(map[string]interface{})
				detachVpcDhcpOptionsSetRequest["VpcId"] = associateVpc["vpc_id"].(string)
				conn, err := client.NewVpcClient()
				if err != nil {
					return WrapError(err)
				}
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					request["ClientToken"] = buildClientToken("DetachDhcpOptionsSetFromVpc")
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, detachVpcDhcpOptionsSetRequest, &runtime)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, detachVpcDhcpOptionsSetRequest)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}
			d.SetPartial("associate_vpcs")
		}
		if added.Len() > 0 {
			action := "AttachDhcpOptionsSetToVpc"
			attachVpcDhcpOptionsSetRequest := map[string]interface{}{
				"DhcpOptionsSetId": d.Id(),
			}
			attachVpcDhcpOptionsSetRequest["RegionId"] = client.RegionId
			if _, ok := d.GetOkExists("dry_run"); ok {
				attachVpcDhcpOptionsSetRequest["DryRun"] = d.Get("dry_run")
			}
			for _, associateVpcs := range added.List() {
				associateVpc := associateVpcs.(map[string]interface{})
				attachVpcDhcpOptionsSetRequest["VpcId"] = associateVpc["vpc_id"].(string)
				conn, err := client.NewVpcClient()
				if err != nil {
					return WrapError(err)
				}
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					request["ClientToken"] = buildClientToken("AttachDhcpOptionsSetToVpc")
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, attachVpcDhcpOptionsSetRequest, &runtime)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, attachVpcDhcpOptionsSetRequest)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}
			d.SetPartial("associate_vpcs")
		}
	}
	return resourceAlicloudVpcDhcpOptionsSetRead(d, meta)
}
func resourceAlicloudVpcDhcpOptionsSetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDhcpOptionsSet"
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DhcpOptionsSetId": d.Id(),
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["RegionId"] = client.RegionId
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 10*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("DeleteDhcpOptionsSet")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectStatus.DhcpOptionsSet"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDhcpOptionsSetId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
