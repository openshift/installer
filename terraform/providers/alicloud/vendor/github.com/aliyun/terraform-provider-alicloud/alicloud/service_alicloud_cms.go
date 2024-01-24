package alicloud

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type CmsService struct {
	client *connectivity.AliyunClient
}

type IspCities []map[string]string

func (s *CmsService) BuildCmsCommonRequest(region string) *requests.CommonRequest {
	request := requests.NewCommonRequest()
	return request
}

func (s *CmsService) BuildCmsAlarmRequest(id string) *requests.CommonRequest {

	request := s.BuildCmsCommonRequest(s.client.RegionId)
	request.QueryParams["Id"] = id

	return request
}

func (s *CmsService) DescribeAlarm(id string) (alarm cms.AlarmInDescribeMetricRuleList, err error) {
	request := cms.CreateDescribeMetricRuleListRequest()
	request.RuleIds = id

	wait := incrementalWait(3*time.Second, 5*time.Second)
	var response *cms.DescribeMetricRuleListResponse
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DescribeMetricRuleList(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*cms.DescribeMetricRuleListResponse)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InternalError", "ResourceNotFound"}) {
			return alarm, WrapErrorf(Error(GetNotFoundMessage("Alarm Rule", id)), NotFoundWithResponse, response)
		}
		return alarm, err
	}

	if len(response.Alarms.Alarm) < 1 {
		return alarm, GetNotFoundErrorFromString(GetNotFoundMessage("Alarm Rule", id))
	}

	return response.Alarms.Alarm[0], nil
}

func (s *CmsService) WaitForCmsAlarm(id string, enabled bool, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		alarm, err := s.DescribeAlarm(id)
		if err != nil {
			return err
		}

		if alarm.EnableState == enabled {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("Alarm", strconv.FormatBool(enabled)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *CmsService) BuildJsonWebhook(webhook string) string {
	if webhook != "" {
		return fmt.Sprintf("{\"method\":\"post\",\"url\":\"%s\"}", webhook)
	}
	return ""
}

func (s *CmsService) ExtractWebhookFromJson(webhookJson string) (string, error) {
	byt := []byte(webhookJson)
	var dat map[string]interface{}
	if err := json.Unmarshal(byt, &dat); err != nil {
		return "", err
	}
	return dat["url"].(string), nil
}

func (s *CmsService) DescribeSiteMonitor(id, keyword string) (siteMonitor cms.SiteMonitor, err error) {
	listRequest := cms.CreateDescribeSiteMonitorListRequest()
	listRequest.Keyword = keyword
	listRequest.TaskId = id
	var raw interface{}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DescribeSiteMonitorList(listRequest)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser, "ExceedingQuota"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return siteMonitor, WrapError(err)
	}
	list := raw.(*cms.DescribeSiteMonitorListResponse)
	if len(list.SiteMonitors.SiteMonitor) < 1 {
		return siteMonitor, GetNotFoundErrorFromString(GetNotFoundMessage("Site Monitor", id))

	}
	for _, v := range list.SiteMonitors.SiteMonitor {
		if v.TaskName == keyword || v.TaskId == id {
			return v, nil
		}
	}
	return siteMonitor, GetNotFoundErrorFromString(GetNotFoundMessage("Site Monitor", id))
}

func (s *CmsService) GetIspCities(id string) (ispCities IspCities, err error) {
	request := cms.CreateDescribeSiteMonitorAttributeRequest()
	request.TaskId = id

	var raw interface{}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DescribeSiteMonitorAttribute(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser, "ExceedingQuota"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return nil, WrapError(err)
	}

	response := raw.(*cms.DescribeSiteMonitorAttributeResponse)
	ispCity := response.SiteMonitors.IspCities.IspCity

	var list []map[string]string
	for _, element := range ispCity {
		list = append(list, map[string]string{"city": element.City, "isp": element.Isp})
	}

	return list, nil
}

func (s *CmsService) DescribeCmsAlarmContact(id string) (object cms.Contact, err error) {
	request := cms.CreateDescribeContactListRequest()
	request.RegionId = s.client.RegionId

	request.ContactName = id

	raw, err := s.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.DescribeContactList(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ContactNotExists", "ResourceNotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("CmsAlarmContact", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*cms.DescribeContactListResponse)
	if response.Code != "200" {
		err = Error("DescribeContactList failed for " + response.Message)
		return
	}

	if len(response.Contacts.Contact) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("CmsAlarmContact", id)), NotFoundMsg, ProviderERROR, response.RequestId)
		return
	}
	return response.Contacts.Contact[0], nil
}

func (s *CmsService) DescribeCmsAlarmContactGroup(id string) (object cms.ContactGroup, err error) {
	request := cms.CreateDescribeContactGroupListRequest()
	request.RegionId = s.client.RegionId

	request.PageNumber = requests.NewInteger(1)
	request.PageSize = requests.NewInteger(20)
	for {

		var raw interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err = s.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
				return cmsClient.DescribeContactGroupList(request)
			})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			if IsExpectedErrors(err, []string{"ContactGroupNotExists", "ResourceNotFound"}) {
				err = WrapErrorf(Error(GetNotFoundMessage("CmsAlarmContactGroup", id)), NotFoundMsg, ProviderERROR)
				return object, err
			}
			err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
			return object, err
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*cms.DescribeContactGroupListResponse)
		if response.Code != "200" {
			err = Error("DescribeContactGroupList failed for " + response.Message)
			return object, err
		}

		if len(response.ContactGroupList.ContactGroup) < 1 {
			err = WrapErrorf(Error(GetNotFoundMessage("CmsAlarmContactGroup", id)), NotFoundMsg, ProviderERROR, response.RequestId)
			return object, err
		}
		for _, object := range response.ContactGroupList.ContactGroup {
			if object.Name == id {
				return object, nil
			}
		}
		if len(response.ContactGroupList.ContactGroup) < PageSizeMedium {
			break
		}
		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return object, WrapError(err)
		} else {
			request.PageNumber = page
		}
	}
	err = WrapErrorf(Error(GetNotFoundMessage("CmsAlarmContactGroup", id)), NotFoundMsg, ProviderERROR)
	return
}

func (s *CmsService) DescribeCmsGroupMetricRule(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCmsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeMetricRuleList"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"RuleIds":  id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(6*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"ExceedingQuota", "Throttling.User"}) {
				wait()
				return resource.RetryableError(err)
			}
			err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if IsExpectedErrors(err, []string{"GroupMetricRuleNotExists", "ResourceNotFound", "ResourceNotFoundError"}) {
		err = WrapErrorf(Error(GetNotFoundMessage("CmsGroupMetricRule", id)), NotFoundMsg, ProviderERROR)
		return object, err
	}
	if fmt.Sprintf(`%v`, response["Code"]) != "200" {
		err = Error("DescribeMetricRuleList failed for " + response["Message"].(string))
		return object, err
	}
	v, err := jsonpath.Get("$.Alarms.Alarm", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Alarms.Alarm", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("CloudMonitorService", id)), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["RuleId"].(string) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudMonitorService", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *CmsService) SetResourceTags(d *schema.ResourceData, resourceType string) error {

	if d.HasChange("tags") {
		added, removed := parsingTags(d)
		conn, err := s.client.NewCmsClient()
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
			action := "RemoveTags"
			request := map[string]interface{}{
				"RegionId":   s.client.RegionId,
				"GroupIds.1": d.Id(),
			}
			oraw, _ := d.GetChange("tags")
			removedTags := oraw.(map[string]interface{})
			count := 1
			for _, key := range removedTagKeys {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = removedTags[key]
			}
			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			action := "AddTags"
			request := map[string]interface{}{
				"RegionId":   s.client.RegionId,
				"GroupIds.1": d.Id(),
			}
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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

func (s *CmsService) DescribeCmsMonitorGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCmsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeMonitorGroups"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"GroupId":  id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(action, response, request)
	if IsExpectedErrorCodes(fmt.Sprintf("%v", response["Code"]), []string{"GroupsNotExists", "ResourceNotFound"}) {
		err = WrapErrorf(Error(GetNotFoundMessage("CmsMonitorGroup", id)), NotFoundMsg, ProviderERROR)
		return object, err
	}
	if fmt.Sprintf(`%v`, response["Code"]) != "200" {
		err = Error("DescribeMonitorGroups failed for " + response["Message"].(string))
		return object, err
	}
	v, err := jsonpath.Get("$.Resources.Resource", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Resources.Resource", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("CloudMonitorService", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(formatInt(v.([]interface{})[0].(map[string]interface{})["GroupId"])) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudMonitorService", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *CmsService) DescribeCmsMonitorGroupInstances(id string) (object []map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCmsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeMonitorGroupInstances"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"GroupId":  id,
	}
	request["PageSize"] = PageSizeMedium
	request["PageNumber"] = 1
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &runtime)
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
			if IsExpectedErrors(err, []string{"ResourceNotFound", "ResourceNotFoundError"}) {
				err = WrapErrorf(Error(GetNotFoundMessage("CmsMonitorGroupInstances", id)), NotFoundMsg, ProviderERROR)
				return object, err
			}
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		v, err := jsonpath.Get("$.Resources.Resource", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Resources.Resource", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("CmsMonitorGroupInstances", id)), NotFoundWithResponse, response)
		}

		for _, v := range v.([]interface{}) {
			object = append(object, v.(map[string]interface{}))
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	return object, nil
}
func (s *CmsService) DescribeCmsMetricRuleTemplate(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCmsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeMetricRuleTemplateAttribute"
	request := map[string]interface{}{
		"TemplateId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"ResourceNotFound"}) {
			return nil, WrapErrorf(Error(GetNotFoundMessage("CloudMonitorService:MetricRuleTemplate", id)), NotFoundMsg, ProviderERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	v, err := jsonpath.Get("$.Resource", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Resource", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *CmsService) DescribeCmsDynamicTagGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCmsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeDynamicTagRuleList"
	request := map[string]interface{}{
		"TagRegionId": s.client.RegionId,
		"PageSize":    PageSizeLarge,
		"PageNumber":  1,
	}
	idExist := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &runtime)
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
		if fmt.Sprint(response["Success"]) == "false" {
			return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		v, err := jsonpath.Get("$.TagGroupList.TagGroup", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.TagGroupList.TagGroup", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudMonitorService", id)), NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["DynamicTagRuleId"]) == id {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("CloudMonitorService", id)), NotFoundWithResponse, response)
	}
	return
}
