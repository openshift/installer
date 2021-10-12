package alicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type EcdService struct {
	client *connectivity.AliyunClient
}

func (s *EcdService) DescribeEcdPolicyGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewGwsecdClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribePolicyGroups"
	request := map[string]interface{}{
		"RegionId":      s.client.RegionId,
		"PolicyGroupId": []string{id},
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &runtime)
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
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.DescribePolicyGroups", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.DescribePolicyGroups", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("ECD", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["PolicyGroupId"]) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("ECD", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func parsingResourceType(d *schema.ResourceData, resourceType string) ([]interface{}, []interface{}) {
	oraw, nraw := d.GetChange(resourceType)
	removed := oraw.(*schema.Set).List()
	added := nraw.(*schema.Set).List()
	return removed, added
}

func (s *EcdService) setAuthAccessPolicyRules(d *schema.ResourceData, request map[string]interface{}, resourceType string) error {
	if d.HasChange(resourceType) {
		removed, added := parsingResourceType(d, resourceType)
		if len(removed) > 0 {
			action := "ModifyPolicyGroup"
			conn, err := s.client.NewGwsecdClient()
			if err != nil {
				return WrapError(err)
			}
			req := map[string]interface{}{
				"PolicyGroupId": d.Id(),
			}
			var response map[string]interface{}
			for authorizeAccessPolicyRulesPtr, authorizeAccessPolicyRules := range removed {
				authorizeAccessPolicyRulesArg := authorizeAccessPolicyRules.(map[string]interface{})
				req["RevokeAccessPolicyRule."+fmt.Sprint(authorizeAccessPolicyRulesPtr+1)+".CidrIp"] = authorizeAccessPolicyRulesArg["cidr_ip"]
				req["RevokeAccessPolicyRule."+fmt.Sprint(authorizeAccessPolicyRulesPtr+1)+".Description"] = authorizeAccessPolicyRulesArg["description"]
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, req, &util.RuntimeOptions{})
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
		}
		if len(added) > 0 {
			for authorizeAccessPolicyRulesPtr, authorizeAccessPolicyRules := range added {
				authorizeAccessPolicyRulesArg := authorizeAccessPolicyRules.(map[string]interface{})
				request["AuthorizeAccessPolicyRule."+fmt.Sprint(authorizeAccessPolicyRulesPtr+1)+".CidrIp"] = authorizeAccessPolicyRulesArg["cidr_ip"]
				request["AuthorizeAccessPolicyRule."+fmt.Sprint(authorizeAccessPolicyRulesPtr+1)+".Description"] = authorizeAccessPolicyRulesArg["description"]
			}
		}
	}
	return nil
}

func (s *EcdService) setAuthSecurityPolicyRules(d *schema.ResourceData, request map[string]interface{}, resourceType string) error {
	if d.HasChange(resourceType) {
		removed, added := parsingResourceType(d, resourceType)
		if len(removed) > 0 {
			action := "ModifyPolicyGroup"
			conn, err := s.client.NewGwsecdClient()
			if err != nil {
				return WrapError(err)
			}
			req := map[string]interface{}{
				"PolicyGroupId": d.Id(),
			}
			var response map[string]interface{}
			for authorizeSecurityPolicyRulesPtr, authorizeSecurityPolicyRules := range removed {
				authorizeSecurityPolicyRulesArg := authorizeSecurityPolicyRules.(map[string]interface{})
				req["RevokeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".CidrIp"] = authorizeSecurityPolicyRulesArg["cidr_ip"]
				req["RevokeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".Description"] = authorizeSecurityPolicyRulesArg["description"]
				req["RevokeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".IpProtocol"] = authorizeSecurityPolicyRulesArg["ip_protocol"]
				req["RevokeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".Policy"] = authorizeSecurityPolicyRulesArg["policy"]
				req["RevokeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".PortRange"] = authorizeSecurityPolicyRulesArg["port_range"]
				req["RevokeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".Priority"] = authorizeSecurityPolicyRulesArg["priority"]
				req["RevokeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".Type"] = authorizeSecurityPolicyRulesArg["type"]
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, req, &util.RuntimeOptions{})
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
		}
		if len(added) > 0 {
			for authorizeSecurityPolicyRulesPtr, authorizeSecurityPolicyRules := range added {
				authorizeSecurityPolicyRulesArg := authorizeSecurityPolicyRules.(map[string]interface{})
				request["AuthorizeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".CidrIp"] = authorizeSecurityPolicyRulesArg["cidr_ip"]
				request["AuthorizeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".Description"] = authorizeSecurityPolicyRulesArg["description"]
				request["AuthorizeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".IpProtocol"] = authorizeSecurityPolicyRulesArg["ip_protocol"]
				request["AuthorizeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".Policy"] = authorizeSecurityPolicyRulesArg["policy"]
				request["AuthorizeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".PortRange"] = authorizeSecurityPolicyRulesArg["port_range"]
				request["AuthorizeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".Priority"] = authorizeSecurityPolicyRulesArg["priority"]
				request["AuthorizeSecurityPolicyRule."+fmt.Sprint(authorizeSecurityPolicyRulesPtr+1)+".Type"] = authorizeSecurityPolicyRulesArg["type"]
			}
		}
	}
	return nil
}
