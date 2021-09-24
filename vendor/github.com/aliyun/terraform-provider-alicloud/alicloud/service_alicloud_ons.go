package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type OnsService struct {
	client *connectivity.AliyunClient
}

func (s *OnsService) DescribeOnsInstance(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewOnsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "OnsInstanceBaseInfo"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-02-14"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"INSTANCE_NOT_FOUNDError", "InvalidDomainName.NoExist"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("OnsInstance", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.InstanceBaseInfo", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.InstanceBaseInfo", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *OnsService) DescribeOnsTopic(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewOnsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "OnsTopicList"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": parts[0],
		"Topic":      parts[1],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-02-14"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"AUTH_RESOURCE_OWNER_ERROR", "INSTANCE_NOT_FOUND"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("OnsTopic", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.Data.PublishInfoDo", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data.PublishInfoDo", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("Ons", id)), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["Topic"].(string) != parts[1] {
			return object, WrapErrorf(Error(GetNotFoundMessage("Ons", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *OnsService) DescribeOnsGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewOnsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "OnsGroupList"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"GroupId":    parts[1],
		"InstanceId": parts[0],
		"GroupType":  "all",
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-02-14"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"AUTH_RESOURCE_OWNER_ERROR", "INSTANCE_NOT_FOUND"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("OnsGroup", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.Data.SubscribeInfoDo", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data.SubscribeInfoDo", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("Ons", id)), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["GroupId"].(string) != parts[1] {
			return object, WrapErrorf(Error(GetNotFoundMessage("Ons", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *OnsService) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	var parts []string
	var err error
	if resourceType == "GROUP" || resourceType == "TOPIC" {
		parts, err = ParseResourceId(d.Id(), 2)
		if err != nil {
			return WrapError(err)
		}
	}
	if d.HasChange("tags") {
		added, removed := parsingTags(d)
		conn, err := s.client.NewOnsClient()
		if err != nil {
			return WrapError(err)
		}

		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action := "UntagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
			}
			if resourceType == "INSTANCE" {
				request["ResourceId.1"] = d.Id()
			}
			if resourceType == "GROUP" || resourceType == "TOPIC" {
				request["InstanceId"] = parts[0]
				request["ResourceId.1"] = parts[1]
			}
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-02-14"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
				if err != nil {
					if IsThrottling(err) {
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
		}
		if len(added) > 0 {
			action := "TagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
			}
			if resourceType == "INSTANCE" {
				request["ResourceId.1"] = d.Id()
			}
			if resourceType == "GROUP" || resourceType == "TOPIC" {
				request["InstanceId"] = parts[0]
				request["ResourceId.1"] = parts[1]
			}
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-02-14"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
				if err != nil {
					if IsThrottling(err) {
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
		}
		d.SetPartial("tags")
	}
	return nil
}

func (s *OnsService) OnsTopicStatus(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewOnsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "OnsTopicStatus"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": parts[0],
		"Topic":      parts[1],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(11*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-02-14"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"Throttling.User"}) {
				wait()
				return resource.RetryableError(err)
			}
			if IsExpectedErrors(err, []string{"AUTH_RESOURCE_OWNER_ERROR", "INSTANCE_NOT_FOUND"}) {
				err = WrapErrorf(Error(GetNotFoundMessage("OnsTopic", id)), NotFoundMsg, ProviderERROR)
				return resource.NonRetryableError(err)
			}
			err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		v, err := jsonpath.Get("$.Data", response)
		if err != nil {
			return resource.NonRetryableError(WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response))
		}
		object = v.(map[string]interface{})
		return nil
	})
	return
}

func (s *OnsService) ListTagResources(id string, resourceType string) (object interface{}, err error) {
	conn, err := s.client.NewOnsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListTagResources"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"ResourceType": resourceType,
		"ResourceId.1": id,
	}
	tags := make([]interface{}, 0)
	var response map[string]interface{}

	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-02-14"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{Throttling}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			v, err := jsonpath.Get("$.TagResources", response)
			if err != nil {
				return resource.NonRetryableError(WrapErrorf(err, FailedGetAttributeMsg, id, "$.TagResources.TagResource", response))
			}
			if v != nil {
				if v != nil {
					tags = append(tags, v.([]interface{})...)
				}
			}
			return nil
		})
		if err != nil {
			err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
			return
		}
		if response["NextToken"] == nil {
			break
		}
		request["NextToken"] = response["NextToken"]
	}

	return tags, nil
}
