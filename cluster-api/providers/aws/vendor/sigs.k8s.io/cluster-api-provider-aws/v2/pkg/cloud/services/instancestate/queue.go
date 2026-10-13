/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package instancestate

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/pkg/errors"

	iamv1 "sigs.k8s.io/cluster-api-provider-aws/v2/iam/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
)

func (s *Service) reconcileSQSQueue(ctx context.Context) error {
	attrs := make(map[string]string)
	attrs[string(sqstypes.QueueAttributeNameReceiveMessageWaitTimeSeconds)] = "20"

	_, err := s.SQSClient.CreateQueue(ctx, &sqs.CreateQueueInput{
		QueueName:  aws.String(GenerateQueueName(s.scope.Name())),
		Attributes: attrs,
	})

	smithyErr := awserrors.ParseSmithyError(err)
	if smithyErr != nil {
		if smithyErr.ErrorCode() == (&sqstypes.QueueNameExists{}).ErrorCode() {
			return nil
		}
	}
	return errors.Wrap(err, "unable to create new queue")
}

func (s *Service) deleteSQSQueue(ctx context.Context) error {
	resp, err := s.SQSClient.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(GenerateQueueName(s.scope.Name()))})
	if err != nil {
		if queueNotFoundError(err) {
			return nil
		}
		return errors.Wrap(err, "unable to get queue URL")
	}
	_, err = s.SQSClient.DeleteQueue(ctx, &sqs.DeleteQueueInput{QueueUrl: resp.QueueUrl})
	if err != nil && queueNotFoundError(err) {
		return nil
	}

	return errors.Wrap(err, "unable to delete queue")
}

func (s *Service) createPolicyForRule(ctx context.Context, input *createPolicyForRuleInput) error {
	attrs := make(map[string]string)
	policy := iamv1.PolicyDocument{
		Version: iamv1.CurrentVersion,
		ID:      input.QueueArn,
		Statement: iamv1.Statements{
			iamv1.StatementEntry{
				Sid:       fmt.Sprintf("CAPAEvents_%s_%s", s.getEC2RuleName(), GenerateQueueName(s.scope.Name())),
				Effect:    iamv1.EffectAllow,
				Principal: iamv1.Principals{iamv1.PrincipalService: iamv1.PrincipalID{"events.amazonaws.com"}},
				Action:    iamv1.Actions{"sqs:SendMessage"},
				Resource:  iamv1.Resources{input.QueueArn},
				Condition: iamv1.Conditions{
					"ArnEquals": map[string]string{"aws:SourceArn": input.RuleArn},
				},
			},
		},
	}
	policyData, err := json.Marshal(policy)
	if err != nil {
		return errors.Wrap(err, "unable to JSON marshal policy")
	}
	attrs[string(sqstypes.QueueAttributeNamePolicy)] = string(policyData)

	_, err = s.SQSClient.SetQueueAttributes(ctx, &sqs.SetQueueAttributesInput{
		QueueUrl:   aws.String(input.QueueURL),
		Attributes: attrs,
	})

	return errors.Wrap(err, "unable to update queue attributes")
}

// GenerateQueueName will generate a queue name.
func GenerateQueueName(clusterName string) string {
	adjusted := strings.ReplaceAll(clusterName, ".", "-")
	return fmt.Sprintf("%s-queue", adjusted)
}

func queueNotFoundError(err error) bool {
	smithyErr := awserrors.ParseSmithyError(err)
	if smithyErr == nil {
		return false
	}
	return smithyErr.ErrorCode() == (&sqstypes.QueueDoesNotExist{}).ErrorCode()
}

type createPolicyForRuleInput struct {
	QueueArn string
	QueueURL string
	RuleArn  string
}
