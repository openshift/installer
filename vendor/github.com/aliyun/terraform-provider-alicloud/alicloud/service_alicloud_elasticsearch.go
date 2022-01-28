package alicloud

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/elasticsearch"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ElasticsearchService struct {
	client *connectivity.AliyunClient
}

func (s *ElasticsearchService) DescribeElasticsearchInstance(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewElasticsearchClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeInstance"
	response, err = conn.DoRequestWithAction(StringPointer(action), StringPointer("2017-06-13"), nil, StringPointer("GET"), StringPointer("AK"),
		String(fmt.Sprintf("/openapi/instances/%s", id)), nil, nil, nil, &util.RuntimeOptions{})

	addDebug(action, response, nil)
	if err != nil {
		if IsExpectedErrors(err, []string{"InstanceNotFound"}) {
			return object, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.body.Result", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.body.Result", response)
	}
	object = v.(map[string]interface{})
	if (object["instanceId"].(string)) != id {
		return object, WrapErrorf(Error(GetNotFoundMessage("Elasticsearch Instance", id)), NotFoundWithResponse, response)
	}

	return object, WrapError(err)
}

func (s *ElasticsearchService) ElasticsearchStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeElasticsearchInstance(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["status"].(string) == failState {
				return object, object["status"].(string), WrapError(Error(FailedToReachTargetStatus, object["status"].(string)))
			}
		}

		return object, object["status"].(string), nil
	}
}

func (s *ElasticsearchService) ElasticsearchRetryFunc(wait func(), errorCodeList []string, do func(*elasticsearch.Client) (interface{}, error)) (interface{}, error) {
	var raw interface{}
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithElasticsearchClient(do)

		if err != nil {
			if IsExpectedErrors(err, errorCodeList) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	return raw, WrapError(err)
}

func (s *ElasticsearchService) TriggerNetwork(d *schema.ResourceData, content map[string]interface{}, meta interface{}) error {
	var response map[string]interface{}
	conn, err := s.client.NewElasticsearchClient()
	if err != nil {
		return WrapError(err)
	}
	action := "TriggerNetwork"
	requestQuery := map[string]*string{
		"clientToken": StringPointer(buildClientToken(action)),
	}

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := conn.DoRequestWithAction(StringPointer(action), StringPointer("2017-06-13"), nil, StringPointer("POST"), StringPointer("AK"),
			String(fmt.Sprintf("/openapi/instances/%s/actions/network-trigger", d.Id())), requestQuery, nil, content, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction", "InternalServerError"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, content)
		return nil
	})

	addDebug(action, response, content)
	if err != nil {
		if IsExpectedErrors(err, []string{"RepetitionOperationError"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, s.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func (s *ElasticsearchService) ModifyWhiteIps(d *schema.ResourceData, content map[string]interface{}, meta interface{}) error {
	var response map[string]interface{}
	conn, err := s.client.NewElasticsearchClient()
	if err != nil {
		return WrapError(err)
	}
	action := "ModifyWhiteIps"
	requestQuery := map[string]*string{
		"clientToken": StringPointer(buildClientToken(action)),
	}

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := conn.DoRequestWithAction(StringPointer(action), StringPointer("2017-06-13"), nil, StringPointer("POST"), StringPointer("AK"),
			String(fmt.Sprintf("/openapi/instances/%s/actions/modify-white-ips", d.Id())), requestQuery, nil, content, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction", "InternalServerError"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, nil)
		return nil
	})

	addDebug(action, response, nil)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, s.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func (s *ElasticsearchService) DescribeElasticsearchTags(id string) (tags map[string]string, err error) {
	resourceIds, err := json.Marshal([]string{id})
	if err != nil {
		tmp := make(map[string]string)
		return tmp, WrapError(err)
	}

	request := elasticsearch.CreateListTagResourcesRequest()
	request.RegionId = s.client.RegionId
	request.ResourceIds = string(resourceIds)
	request.ResourceType = strings.ToUpper(string(TagResourceInstance))

	raw, err := s.client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
		return elasticsearchClient.ListTagResources(request)
	})

	addDebug(request.GetActionName(), raw, request.RoaRequest, request)
	if err != nil {
		tmp := make(map[string]string)
		return tmp, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*elasticsearch.ListTagResourcesResponse)
	return s.tagsToMap(response.TagResources.TagResource), nil
}

func (s *ElasticsearchService) tagsToMap(tagSet []elasticsearch.TagResourceItem) (tags map[string]string) {
	result := make(map[string]string)
	for _, t := range tagSet {
		if !elasticsearchTagIgnored(t.TagKey, t.TagValue) {
			result[t.TagKey] = t.TagValue
		}
	}

	return result
}

func (s *ElasticsearchService) diffElasticsearchTags(oldTags, newTags map[string]interface{}) (remove []string, add []map[string]string) {
	for k, _ := range oldTags {
		remove = append(remove, k)
	}
	for k, v := range newTags {
		tag := map[string]string{
			"key":   k,
			"value": v.(string),
		}

		add = append(add, tag)
	}
	return
}

func (s *ElasticsearchService) getActionType(actionType bool) string {
	if actionType == true {
		return string(OPEN)
	} else {
		return string(CLOSE)
	}
}

func updateDescription(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "UpdateDescription"

	content := make(map[string]interface{})
	content["description"] = d.Get("description").(string)
	requestQuery := map[string]*string{
		"clientToken": StringPointer(buildClientToken(action)),
	}
	elasticsearchClient, err := client.NewElasticsearchClient()
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := elasticsearchClient.DoRequestWithAction(StringPointer(action), StringPointer("2017-06-13"), nil, StringPointer("POST"), StringPointer("AK"),
			String(fmt.Sprintf("/openapi/instances/%s/description", d.Id())), requestQuery, nil, content, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"GetCustomerLabelFail"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, nil)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

func updateInstanceTags(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}

	oraw, nraw := d.GetChange("tags")
	o := oraw.(map[string]interface{})
	n := nraw.(map[string]interface{})
	remove, add := elasticsearchService.diffElasticsearchTags(o, n)

	// 对系统 Tag 进行过滤
	removeTagKeys := make([]string, 0)
	for _, v := range remove {
		if !elasticsearchTagIgnored(v, "") {
			removeTagKeys = append(removeTagKeys, v)
		}
	}
	if len(removeTagKeys) > 0 {
		tagKeys, err := json.Marshal(removeTagKeys)
		if err != nil {
			return WrapError(err)
		}

		resourceIds, err := json.Marshal([]string{d.Id()})
		if err != nil {
			return WrapError(err)
		}
		request := elasticsearch.CreateUntagResourcesRequest()
		request.RegionId = client.RegionId
		request.TagKeys = string(tagKeys)
		request.ResourceType = strings.ToUpper(string(TagResourceInstance))
		request.ResourceIds = string(resourceIds)
		request.SetContentType("application/json")

		raw, err := client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
			return elasticsearchClient.UntagResources(request)
		})

		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}

	if len(add) > 0 {
		content := make(map[string]interface{})
		content["ResourceIds"] = []string{d.Id()}
		content["ResourceType"] = strings.ToUpper(string(TagResourceInstance))
		content["Tags"] = add
		data, err := json.Marshal(content)
		if err != nil {
			return WrapError(err)
		}

		request := elasticsearch.CreateTagResourcesRequest()
		request.RegionId = client.RegionId
		request.SetContent(data)
		request.SetContentType("application/json")
		raw, err := client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
			return elasticsearchClient.TagResources(request)
		})

		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}

	return nil
}

func updateInstanceChargeType(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchClient, err := client.NewElasticsearchClient()
	action := "UpdateInstanceChargeType"

	content := make(map[string]interface{})
	content["paymentType"] = strings.ToLower(d.Get("instance_charge_type").(string))
	if d.Get("instance_charge_type").(string) == string(PrePaid) {
		paymentInfo := make(map[string]interface{})
		if d.Get("period").(int) >= 12 {
			paymentInfo["duration"] = d.Get("period").(int) / 12
			paymentInfo["pricingCycle"] = string(Year)
		} else {
			paymentInfo["duration"] = d.Get("period").(int)
			paymentInfo["pricingCycle"] = string(Month)
		}

		content["paymentInfo"] = paymentInfo
	}
	requestQuery := map[string]*string{
		"clientToken": StringPointer(buildClientToken(action)),
	}
	response, err := elasticsearchClient.DoRequestWithAction(StringPointer(action), StringPointer("2017-06-13"), nil, StringPointer("POST"), StringPointer("AK"),
		String(fmt.Sprintf("/openapi/instances/%s/actions/convert-pay-type", d.Id())), requestQuery, nil, content, &util.RuntimeOptions{})

	time.Sleep(10 * time.Second)

	addDebug(action, response, content)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), response, AlibabaCloudSdkGoERROR)
	}
	return nil
}

func renewInstance(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchClient, err := client.NewElasticsearchClient()
	action := "RenewInstance"

	content := make(map[string]interface{})
	if d.Get("period").(int) >= 12 {
		content["duration"] = d.Get("period").(int) / 12
		content["pricingCycle"] = string(Year)
	} else {
		content["duration"] = d.Get("period").(int)
		content["pricingCycle"] = string(Month)
	}
	requestQuery := map[string]*string{
		"clientToken": StringPointer(buildClientToken(action)),
	}
	response, err := elasticsearchClient.DoRequestWithAction(StringPointer(action), StringPointer("2017-06-13"), nil, StringPointer("POST"), StringPointer("AK"),
		String(fmt.Sprintf("/openapi/instances/%s/actions/renew", d.Id())), requestQuery, nil, content, &util.RuntimeOptions{})

	time.Sleep(10 * time.Second)

	addDebug(action, response, content)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), response, AlibabaCloudSdkGoERROR)
	}
	return nil
}

func updateDataNodeAmount(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}
	conn, err := client.NewElasticsearchClient()
	if err != nil {
		return WrapError(err)
	}
	action := "UpdateInstance"

	var response map[string]interface{}
	requestQuery := map[string]*string{
		"clientToken": StringPointer(buildClientToken(action)),
	}
	content := make(map[string]interface{})
	content["nodeAmount"] = d.Get("data_node_amount").(int)

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := conn.DoRequestWithAction(StringPointer(action), StringPointer("2017-06-13"), nil, StringPointer("PUT"), StringPointer("AK"),
			String(fmt.Sprintf("/openapi/instances/%s", d.Id())), requestQuery, nil, content, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, content)
		return nil
	})

	addDebug(action, response, content)
	if err != nil && !IsExpectedErrors(err, []string{"MustChangeOneResource", "CssCheckUpdowngradeError"}) {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func updateDataNodeSpec(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}
	conn, err := client.NewElasticsearchClient()
	if err != nil {
		return WrapError(err)
	}
	action := "UpdateInstance"

	var response map[string]interface{}
	content := make(map[string]interface{})
	spec := make(map[string]interface{})
	spec["spec"] = d.Get("data_node_spec")
	spec["disk"] = d.Get("data_node_disk_size")
	spec["diskType"] = d.Get("data_node_disk_type")
	content["nodeSpec"] = spec
	requestQuery := map[string]*string{
		"clientToken": StringPointer(buildClientToken(action)),
	}

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := conn.DoRequestWithAction(StringPointer(action), StringPointer("2017-06-13"), nil, StringPointer("PUT"), StringPointer("AK"),
			String(fmt.Sprintf("/openapi/instances/%s", d.Id())), requestQuery, nil, content, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, content)
		return nil
	})

	addDebug(action, response, content)
	if err != nil && !IsExpectedErrors(err, []string{"MustChangeOneResource", "CssCheckUpdowngradeError"}) {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func updateMasterNode(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}
	conn, err := client.NewElasticsearchClient()
	if err != nil {
		return WrapError(err)
	}
	action := "UpdateInstance"

	var response map[string]interface{}
	content := make(map[string]interface{})
	if d.Get("master_node_spec") != nil {
		master := make(map[string]interface{})
		master["spec"] = d.Get("master_node_spec").(string)
		master["amount"] = "3"
		master["diskType"] = "cloud_ssd"
		master["disk"] = "20"
		content["masterConfiguration"] = master
		content["advancedDedicateMaster"] = true
	} else {
		content["advancedDedicateMaster"] = false
	}
	requestQuery := map[string]*string{
		"clientToken": StringPointer(buildClientToken(action)),
	}

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := conn.DoRequestWithAction(StringPointer(action), StringPointer("2017-06-13"), nil, StringPointer("PUT"), StringPointer("AK"),
			String(fmt.Sprintf("/openapi/instances/%s", d.Id())), requestQuery, nil, content, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, content)
		return nil
	})

	if err != nil && !IsExpectedErrors(err, []string{"MustChangeOneResource", "CssCheckUpdowngradeError"}) {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, content)

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func updatePassword(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}
	conn, err := client.NewElasticsearchClient()
	if err != nil {
		return WrapError(err)
	}
	action := "UpdateAdminPassword"

	var response map[string]interface{}
	content := make(map[string]interface{})
	password := d.Get("password").(string)
	kmsPassword := d.Get("kms_encrypted_password").(string)
	if password == "" && kmsPassword == "" {
		return WrapError(Error("One of the 'password' and 'kms_encrypted_password' should be set."))
	}
	if password != "" {
		d.SetPartial("password")
		content["esAdminPassword"] = password
	} else {
		kmsService := KmsService{meta.(*connectivity.AliyunClient)}
		decryptResp, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		content["esAdminPassword"] = decryptResp
		d.SetPartial("kms_encrypted_password")
		d.SetPartial("kms_encryption_context")
	}
	requestQuery := map[string]*string{
		"clientToken": StringPointer(buildClientToken(action)),
	}

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := conn.DoRequestWithAction(StringPointer(action), StringPointer("2017-06-13"), nil, StringPointer("POST"), StringPointer("AK"),
			String(fmt.Sprintf("/openapi/instances/%s/admin-pwd", d.Id())), requestQuery, nil, content, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, content)
		return nil
	})

	addDebug(action, response, content)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func getChargeType(paymentType string) string {
	if strings.ToLower(paymentType) == strings.ToLower(string(PostPaid)) {
		return string(PostPaid)
	} else {
		return string(PrePaid)
	}
}

func filterWhitelist(destIPs []string, localIPs *schema.Set) []string {
	var whitelist []string
	if destIPs != nil {
		for _, ip := range destIPs {
			if (ip == "::1" || ip == "::/0" || ip == "127.0.0.1" || ip == "0.0.0.0/0") && !localIPs.Contains(ip) {
				continue
			}
			whitelist = append(whitelist, ip)
		}
	}
	return whitelist
}

func updateClientNode(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}
	conn, err := client.NewElasticsearchClient()
	if err != nil {
		return WrapError(err)
	}
	action := "UpdateInstance"

	var response map[string]interface{}
	content := make(map[string]interface{})
	content["isHaveClientNode"] = true

	spec := make(map[string]interface{})
	spec["spec"] = d.Get("client_node_spec")
	if d.Get("client_node_amount") == nil {
		spec["amount"] = "2"
	} else {
		spec["amount"] = d.Get("client_node_amount")
	}
	spec["disk"] = "20"
	spec["diskType"] = "cloud_efficiency"
	content["clientNodeConfiguration"] = spec
	requestQuery := map[string]*string{
		"clientToken": StringPointer(buildClientToken(action)),
	}

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := conn.DoRequestWithAction(StringPointer(action), StringPointer("2017-06-13"), nil, StringPointer("PUT"), StringPointer("AK"),
			String(fmt.Sprintf("/openapi/instances/%s", d.Id())), requestQuery, nil, content, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, content)
		return nil
	})

	addDebug(action, response, content)
	if err != nil && !IsExpectedErrors(err, []string{"MustChangeOneResource", "CssCheckUpdowngradeError"}) {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func openHttps(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}
	conn, err := client.NewElasticsearchClient()
	if err != nil {
		return WrapError(err)
	}
	action := "OpenHttps"

	var response map[string]interface{}
	requestQuery := map[string]*string{
		"clientToken": StringPointer(buildClientToken(action)),
	}

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := conn.DoRequestWithAction(StringPointer(action), StringPointer("2017-06-13"), nil, StringPointer("POST"), StringPointer("AK"),
			String(fmt.Sprintf("/openapi/instances/%s/actions/open-https", d.Id())), requestQuery, nil, nil, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, nil)
		return nil
	})

	addDebug(action, response, nil)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func closeHttps(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}
	conn, err := client.NewElasticsearchClient()
	if err != nil {
		return WrapError(err)
	}
	action := "CloseHttps"

	var response map[string]interface{}
	requestQuery := map[string]*string{
		"clientToken": StringPointer(buildClientToken(action)),
	}

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := conn.DoRequestWithAction(StringPointer(action), StringPointer("2017-06-13"), nil, StringPointer("POST"), StringPointer("AK"),
			String(fmt.Sprintf("/openapi/instances/%s/actions/close-https", d.Id())), requestQuery, nil, nil, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, nil)
		return nil
	})

	addDebug(action, response, nil)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
