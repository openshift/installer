package alicloud

import (
	"strings"
	"time"

	ali_mns "github.com/aliyun/aliyun-mns-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

type MnsService struct {
	client *connectivity.AliyunClient
}

func (s *MnsService) SubscriptionNotExistFunc(err error) bool {
	return strings.Contains(err.Error(), "SubscriptionNotExist")
}
func (s *MnsService) TopicNotExistFunc(err error) bool {
	return strings.Contains(err.Error(), "TopicNotExist")
}

func (s *MnsService) QueueNotExistFunc(err error) bool {
	return strings.Contains(err.Error(), "QueueNotExist")
}

func (s *MnsService) DescribeMnsQueue(id string) (response ali_mns.QueueAttribute, err error) {
	raw, err := s.client.WithMnsQueueManager(func(queueManager ali_mns.AliQueueManager) (interface{}, error) {
		return queueManager.GetQueueAttributes(id)
	})
	if err != nil {
		if s.QueueNotExistFunc(err) {
			return response, WrapErrorf(err, NotFoundMsg, AliMnsERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, "GetQueueAttributes", AliMnsERROR)
	}
	addDebug("GetQueueAttributes", raw)
	response, _ = raw.(ali_mns.QueueAttribute)
	if response.QueueName == "" {
		return response, WrapErrorf(Error(GetNotFoundMessage("MnsQueue", id)), NotFoundMsg, ProviderERROR)
	}
	return
}

func (s *MnsService) WaitForMnsQueue(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeMnsQueue(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.QueueName == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.QueueName, id, ProviderERROR)
		}
	}
}

func (s *MnsService) DescribeMnsTopic(id string) (*ali_mns.TopicAttribute, error) {
	response := &ali_mns.TopicAttribute{}
	raw, err := s.client.WithMnsTopicManager(func(topicManager ali_mns.AliTopicManager) (interface{}, error) {
		return topicManager.GetTopicAttributes(id)
	})
	if err != nil {
		if s.TopicNotExistFunc(err) {
			return response, WrapErrorf(err, NotFoundMsg, AliMnsERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, "GetTopicAttributes", AliMnsERROR)
	}
	addDebug("GetTopicAttributes", raw)
	resp, _ := raw.(ali_mns.TopicAttribute)
	if resp.TopicName == "" {
		return response, WrapErrorf(Error(GetNotFoundMessage("MnsTopic", id)), NotFoundMsg, ProviderERROR)
	}
	return &resp, nil
}

func (s *MnsService) WaitForMnsTopic(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeMnsTopic(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.TopicName == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.TopicName, id, ProviderERROR)
		}
	}
}

func (s *MnsService) DescribeMnsTopicSubscription(id string) (*ali_mns.SubscriptionAttribute, error) {
	response := &ali_mns.SubscriptionAttribute{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return response, WrapError(err)
	}
	topicName, name := parts[0], parts[1]

	raw, err := s.client.WithMnsSubscriptionManagerByTopicName(topicName, func(subscriptionManager ali_mns.AliMNSTopic) (interface{}, error) {
		return subscriptionManager.GetSubscriptionAttributes(name)
	})
	if err != nil {
		if s.TopicNotExistFunc(err) || s.SubscriptionNotExistFunc(err) {
			return response, WrapErrorf(err, NotFoundMsg, AliMnsERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, "GetSubscriptionAttributes", AliMnsERROR)
	}
	addDebug("GetSubscriptionAttributes", raw)
	resp, _ := raw.(ali_mns.SubscriptionAttribute)
	response = &resp
	if response.SubscriptionName == "" {
		return response, WrapErrorf(Error(GetNotFoundMessage("MnsTopicSubscription", id)), NotFoundMsg, ProviderERROR)
	}
	return response, nil
}

func (s *MnsService) WaitForMnsTopicSubscription(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeMnsTopicSubscription(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.TopicName+":"+object.SubscriptionName == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.SubscriptionName, id, ProviderERROR)
		}
	}
}
