package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudPrivatelinkVpcEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPrivatelinkVpcEndpointCreate,
		Read:   resourceAlicloudPrivatelinkVpcEndpointRead,
		Update: resourceAlicloudPrivatelinkVpcEndpointUpdate,
		Delete: resourceAlicloudPrivatelinkVpcEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(6 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bandwidth": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"connection_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"endpoint_business_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"endpoint_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"service_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"service_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_endpoint_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudPrivatelinkVpcEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	privatelinkService := PrivatelinkService{client}
	var response map[string]interface{}
	action := "CreateVpcEndpoint"
	request := make(map[string]interface{})
	conn, err := client.NewPrivatelinkClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	if v, ok := d.GetOk("endpoint_description"); ok {
		request["EndpointDescription"] = v
	}

	request["RegionId"] = client.RegionId
	request["SecurityGroupId"] = d.Get("security_group_ids").(*schema.Set).List()
	if v, ok := d.GetOk("service_id"); ok {
		request["ServiceId"] = v
	}

	if v, ok := d.GetOk("service_name"); ok {
		request["ServiceName"] = v
	}

	if v, ok := d.GetOk("vpc_endpoint_name"); ok {
		request["EndpointName"] = v
	}

	request["VpcId"] = d.Get("vpc_id")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_privatelink_vpc_endpoint", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["EndpointId"]))
	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 60*time.Second, privatelinkService.PrivatelinkVpcEndpointStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudPrivatelinkVpcEndpointRead(d, meta)
}
func resourceAlicloudPrivatelinkVpcEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	privatelinkService := PrivatelinkService{client}
	object, err := privatelinkService.DescribePrivatelinkVpcEndpoint(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_privatelink_vpc_endpoint privatelinkService.DescribePrivatelinkVpcEndpoint Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("bandwidth", formatInt(object["Bandwidth"]))
	d.Set("connection_status", object["ConnectionStatus"])
	d.Set("endpoint_business_status", object["EndpointBusinessStatus"])
	d.Set("endpoint_description", object["EndpointDescription"])
	d.Set("endpoint_domain", object["EndpointDomain"])
	d.Set("service_id", object["ServiceId"])
	d.Set("service_name", object["ServiceName"])
	d.Set("status", object["EndpointStatus"])
	d.Set("vpc_endpoint_name", object["EndpointName"])
	d.Set("vpc_id", object["VpcId"])

	listVpcEndpointSecurityGroupsObject, err := privatelinkService.ListVpcEndpointSecurityGroups(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("security_group_ids", convertSecurityGroupIdToStringList(listVpcEndpointSecurityGroupsObject["SecurityGroups"]))
	return nil
}
func resourceAlicloudPrivatelinkVpcEndpointUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	privatelinkService := PrivatelinkService{client}
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"EndpointId": d.Id(),
	}
	request["RegionId"] = client.RegionId

	if d.HasChange("security_group_ids") && !d.IsNewResource() {
		if _, ok := d.GetOkExists("dry_run"); ok {
			request["DryRun"] = d.Get("dry_run")
		}

		conn, err := client.NewPrivatelinkClient()
		if err != nil {
			return WrapError(err)
		}

		o, n := d.GetChange("security_group_ids")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		rl := expandStringList(os.Difference(ns).List())
		al := expandStringList(ns.Difference(os).List())

		if len(al) > 0 {
			for _, a := range al {
				action := "AttachSecurityGroupToVpcEndpoint"
				request["SecurityGroupId"] = a
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					addDebug(action, response, request)
					return nil
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
				stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, privatelinkService.PrivatelinkVpcEndpointStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
		}

		if len(rl) > 0 {
			for _, r := range rl {
				action := "DetachSecurityGroupFromVpcEndpoint"
				request["SecurityGroupId"] = r
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					addDebug(action, response, request)
					return nil
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
				stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, privatelinkService.PrivatelinkVpcEndpointStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
		}
		d.SetPartial("security_group_ids")
	}

	update = false
	updateVpcEndpointAttributeReq := map[string]interface{}{
		"EndpointId": d.Id(),
	}
	updateVpcEndpointAttributeReq["RegionId"] = client.RegionId
	if d.HasChange("endpoint_description") {
		update = true
		updateVpcEndpointAttributeReq["EndpointDescription"] = d.Get("endpoint_description")
	}
	if d.HasChange("vpc_endpoint_name") {
		update = true
		updateVpcEndpointAttributeReq["EndpointName"] = d.Get("vpc_endpoint_name")
	}
	if update {
		if _, ok := d.GetOkExists("dry_run"); ok {
			updateVpcEndpointAttributeReq["DryRun"] = d.Get("dry_run")
		}
		action := "UpdateVpcEndpointAttribute"
		conn, err := client.NewPrivatelinkClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 10*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), nil, updateVpcEndpointAttributeReq, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"EndpointConnectionOperationDenied", "EndpointLocked", "EndpointOperationDenied"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("endpoint_description")
		d.SetPartial("vpc_endpoint_name")
	}
	d.Partial(false)
	return resourceAlicloudPrivatelinkVpcEndpointRead(d, meta)
}
func resourceAlicloudPrivatelinkVpcEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteVpcEndpoint"
	var response map[string]interface{}
	conn, err := client.NewPrivatelinkClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"EndpointId": d.Id(),
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"EndpointOperationDenied"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"EndpointNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
func convertSecurityGroupIdToStringList(src interface{}) (result []interface{}) {
	if err, ok := src.([]interface{}); !ok {
		panic(err)
	}
	for _, v := range src.([]interface{}) {
		vv := v.(map[string]interface{})
		result = append(result, vv["SecurityGroupId"].(string))
	}
	return
}
