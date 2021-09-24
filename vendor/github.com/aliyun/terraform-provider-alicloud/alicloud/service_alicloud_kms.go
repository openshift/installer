package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type KmsService struct {
	client *connectivity.AliyunClient
}

func (s *KmsService) DescribeKmsKey(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewKmsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeKey"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"KeyId":    id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"Forbidden.AliasNotFound", "Forbidden.KeyNotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("KMS:Key", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.KeyMetadata", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.KeyMetadata", response)
	}
	object = v.(map[string]interface{})
	if object["KeyState"] == "PendingDeletion" {
		log.Printf("[WARN] Removing KmsKey  %s because it's already gone", id)
		return object, WrapErrorf(Error(GetNotFoundMessage("KmsKey", id)), NotFoundMsg, ProviderERROR)
	}
	return object, nil
}

func (s *KmsService) Decrypt(ciphertextBlob string, encryptionContext map[string]interface{}) (plaintext string, err error) {
	context, err := json.Marshal(encryptionContext)
	if err != nil {
		return plaintext, WrapError(err)
	}

	var response map[string]interface{}
	conn, err := s.client.NewKmsClient()
	if err != nil {
		return plaintext, WrapError(err)
	}
	action := "Decrypt"
	request := map[string]interface{}{
		"RegionId":          s.client.RegionId,
		"CiphertextBlob":    ciphertextBlob,
		"EncryptionContext": string(context[:]),
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return plaintext, WrapErrorf(err, DefaultErrorMsg, context, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Plaintext", response)
	if err != nil {
		return plaintext, WrapErrorf(err, FailedGetAttributeMsg, context, "$.Plaintext", response)
	}

	return fmt.Sprint(v), err
}

func (s *KmsService) DescribeKmsSecret(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewKmsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeSecret"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"SecretName": id,
		"FetchTags":  "true",
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"Forbidden.ResourceNotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("KMS:Secret", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *KmsService) GetSecretValue(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewKmsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetSecretValue"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"SecretName": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), nil, request, &runtime)
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
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *KmsService) DescribeKmsAlias(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewKmsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListAliases"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"PageNumber": 1,
		"PageSize":   20,
	}
	idExist := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), nil, request, &runtime)
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
		v, err := jsonpath.Get("$.Aliases.Alias", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Aliases.Alias", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("KMS", id)), NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if v.(map[string]interface{})["AliasName"].(string) == id {
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
		return object, WrapErrorf(Error(GetNotFoundMessage("KMS", id)), NotFoundWithResponse, response)
	}
	return
}

func (s *KmsService) DescribeKmsKeyVersion(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewKmsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeKeyVersion"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"KeyId":        parts[0],
		"KeyVersionId": parts[1],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), nil, request, &runtime)
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
	v, err := jsonpath.Get("$.KeyVersion", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.KeyVersion", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *KmsService) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	oldItems, newItems := d.GetChange("tags")
	added := make([]JsonTag, 0)
	for key, value := range newItems.(map[string]interface{}) {
		added = append(added, JsonTag{
			TagKey:   key,
			TagValue: value.(string),
		})
	}
	removed := make([]string, 0)
	for key, _ := range oldItems.(map[string]interface{}) {
		removed = append(removed, key)
	}
	// 对系统 Tag 进行过滤
	removedTagKeys := make([]string, 0)
	for _, v := range removed {
		if !ignoredTags(v, "") {
			removedTagKeys = append(removedTagKeys, v)
		}
	}
	if len(removedTagKeys) > 0 {
		var response map[string]interface{}
		conn, err := s.client.NewKmsClient()
		if err != nil {
			return WrapError(err)
		}
		action := "UntagResource"
		request := map[string]interface{}{
			"RegionId": s.client.RegionId,
		}
		if resourceType == "key" {
			request["KeyId"] = d.Id()
		}
		if resourceType == "secret" {
			request["SecretName"] = d.Id()
		}
		remove, err := json.Marshal(removed)
		if err != nil {
			return WrapError(err)
		}
		request["TagKeys"] = string(remove)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		var response map[string]interface{}
		conn, err := s.client.NewKmsClient()
		if err != nil {
			return WrapError(err)
		}
		action := "TagResource"
		request := map[string]interface{}{
			"RegionId": s.client.RegionId,
		}
		if resourceType == "key" {
			request["KeyId"] = d.Id()
		}
		if resourceType == "secret" {
			request["SecretName"] = d.Id()
		}
		add, err := json.Marshal(added)
		if err != nil {
			return WrapError(err)
		}
		request["Tags"] = string(add)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return nil
}
