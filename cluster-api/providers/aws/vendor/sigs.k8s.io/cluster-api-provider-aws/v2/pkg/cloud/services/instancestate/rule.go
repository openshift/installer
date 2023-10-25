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
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
)

// Ec2StateChangeNotification defines the EC2 instance's state change notification.
const Ec2StateChangeNotification = "EC2 Instance State-change Notification"

// reconcileRules creates rules and attaches the queue as a target.
func (s Service) reconcileRules() error {
	var ruleNotFound bool
	ruleResp, err := s.EventBridgeClient.DescribeRule(&eventbridge.DescribeRuleInput{
		Name: aws.String(s.getEC2RuleName()),
	})
	if err != nil {
		if resourceNotFoundError(err) {
			ruleNotFound = true
		} else {
			return errors.Wrapf(err, "unable to describe rule %s", s.getEC2RuleName())
		}
	}

	if ruleNotFound {
		err = s.createRule()
		if err != nil {
			return errors.Wrap(err, "unable to create rule")
		}
		// fetch newly created rule
		ruleResp, err = s.EventBridgeClient.DescribeRule(&eventbridge.DescribeRuleInput{
			Name: aws.String(s.getEC2RuleName()),
		})

		if err != nil {
			return errors.Wrapf(err, "unable to describe new rule %s", s.getEC2RuleName())
		}
	}

	queueURLResp, err := s.SQSClient.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(GenerateQueueName(s.scope.Name())),
	})

	if err != nil {
		return errors.Wrap(err, "unable to get queue URL")
	}
	queueAttrs, err := s.SQSClient.GetQueueAttributes(&sqs.GetQueueAttributesInput{
		AttributeNames: aws.StringSlice([]string{sqs.QueueAttributeNameQueueArn, sqs.QueueAttributeNamePolicy}),
		QueueUrl:       queueURLResp.QueueUrl,
	})

	if err != nil {
		return errors.Wrap(err, "unable to get queue attributes")
	}

	targetsResp, err := s.EventBridgeClient.ListTargetsByRule(&eventbridge.ListTargetsByRuleInput{
		Rule: aws.String(s.getEC2RuleName()),
	})
	if err != nil {
		return errors.Wrapf(err, "unable to list targets for rule %s", s.getEC2RuleName())
	}

	targetFound := false
	for _, target := range targetsResp.Targets {
		// check if queue is already added as a target
		if *target.Id == GenerateQueueName(s.scope.Name()) && *target.Arn == *queueAttrs.Attributes[sqs.QueueAttributeNameQueueArn] {
			targetFound = true
		}
	}

	if !targetFound {
		_, err = s.EventBridgeClient.PutTargets(&eventbridge.PutTargetsInput{
			Rule: ruleResp.Name,
			Targets: []*eventbridge.Target{{
				Arn: queueAttrs.Attributes[sqs.QueueAttributeNameQueueArn],
				Id:  aws.String(GenerateQueueName(s.scope.Name())),
			}},
		})

		if err != nil {
			return errors.Wrapf(err, "unable to add SQS target %s to rule %s", GenerateQueueName(s.scope.Name()), s.getEC2RuleName())
		}
	}

	if queueAttrs.Attributes[sqs.QueueAttributeNamePolicy] == nil {
		// add a policy for the rule so the rule is authorized to emit messages to the queue
		err = s.createPolicyForRule(&createPolicyForRuleInput{
			QueueArn: *queueAttrs.Attributes[sqs.QueueAttributeNameQueueArn],
			QueueURL: *queueURLResp.QueueUrl,
			RuleArn:  *ruleResp.Arn,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s Service) createRule() error {
	eventPattern := eventPattern{
		Source:     []string{"aws.ec2"},
		DetailType: []string{Ec2StateChangeNotification},
		EventDetail: &eventDetail{
			States: []infrav1.InstanceState{infrav1.InstanceStateShuttingDown, infrav1.InstanceStateTerminated},
		},
	}
	data, err := json.Marshal(eventPattern)
	if err != nil {
		return err
	}
	// create in disabled state so the rule doesn't pick up all EC2 instances. As machines get created,
	// the rule will get updated to track those machines
	_, err = s.EventBridgeClient.PutRule(&eventbridge.PutRuleInput{
		Name:         aws.String(s.getEC2RuleName()),
		EventPattern: aws.String(string(data)),
		State:        aws.String(eventbridge.RuleStateDisabled),
	})

	return err
}

func (s Service) deleteRules() error {
	_, err := s.EventBridgeClient.RemoveTargets(&eventbridge.RemoveTargetsInput{
		Rule: aws.String(s.getEC2RuleName()),
		Ids:  aws.StringSlice([]string{GenerateQueueName(s.scope.Name())}),
	})
	if err != nil && !resourceNotFoundError(err) {
		return errors.Wrapf(err, "unable to remove target %s for rule %s", GenerateQueueName(s.scope.Name()), s.getEC2RuleName())
	}
	_, err = s.EventBridgeClient.DeleteRule(&eventbridge.DeleteRuleInput{
		Name: aws.String(s.getEC2RuleName()),
	})

	if err != nil && resourceNotFoundError(err) {
		return nil
	}

	return err
}

// AddInstanceToEventPattern will add an instance to an event pattern.
func (s Service) AddInstanceToEventPattern(instanceID string) error {
	ruleResp, err := s.EventBridgeClient.DescribeRule(&eventbridge.DescribeRuleInput{
		Name: aws.String(s.getEC2RuleName()),
	})
	if err != nil {
		return errors.Wrapf(err, "unable to describe rule %s", s.getEC2RuleName())
	}
	e := eventPattern{}
	err = json.Unmarshal([]byte(*ruleResp.EventPattern), &e)
	if err != nil {
		return err
	}
	e.DetailType = []string{Ec2StateChangeNotification}

	for _, r := range e.EventDetail.InstanceIDs {
		if r == instanceID {
			// instance is already tracked by rule
			return nil
		}
	}

	e.EventDetail.InstanceIDs = append(e.EventDetail.InstanceIDs, instanceID)
	eventData, err := json.Marshal(e)
	if err != nil {
		return err
	}
	_, err = s.EventBridgeClient.PutRule(&eventbridge.PutRuleInput{
		Name:         aws.String(s.getEC2RuleName()),
		EventPattern: aws.String(string(eventData)),
		State:        aws.String(eventbridge.RuleStateEnabled),
	})
	return err
}

// RemoveInstanceFromEventPattern attempts a best effort update to the event rule to remove the instance.
// Any errors encountered won't be blocking.
func (s Service) RemoveInstanceFromEventPattern(instanceID string) {
	ruleResp, err := s.EventBridgeClient.DescribeRule(&eventbridge.DescribeRuleInput{
		Name: aws.String(s.getEC2RuleName()),
	})
	if err != nil {
		return
	}
	e := eventPattern{}
	err = json.Unmarshal([]byte(*ruleResp.EventPattern), &e)
	if err != nil {
		return
	}
	e.DetailType = []string{Ec2StateChangeNotification}

	found := false
	for i, r := range e.EventDetail.InstanceIDs {
		if r == instanceID {
			found = true
			e.EventDetail.InstanceIDs = append(e.EventDetail.InstanceIDs[:i], e.EventDetail.InstanceIDs[i+1:]...)
			break
		}
	}

	if found {
		eventData, err := json.Marshal(e)
		if err != nil {
			return
		}
		input := &eventbridge.PutRuleInput{
			Name:         aws.String(s.getEC2RuleName()),
			EventPattern: aws.String(string(eventData)),
			State:        aws.String(eventbridge.RuleStateEnabled),
		}

		if len(e.EventDetail.InstanceIDs) == 0 {
			input.State = aws.String(eventbridge.RuleStateDisabled)
		}
		_, _ = s.EventBridgeClient.PutRule(input)
	}
}

func (s Service) getEC2RuleName() string {
	return fmt.Sprintf("%s-ec2-rule", s.scope.Name())
}

func resourceNotFoundError(err error) bool {
	if aerr, ok := err.(awserr.Error); ok && aerr.Code() == eventbridge.ErrCodeResourceNotFoundException {
		return true
	}
	return false
}

type eventPattern struct {
	Source      []string     `json:"source"`
	DetailType  []string     `json:"detail-type,omitempty"`
	EventDetail *eventDetail `json:"detail,omitempty"`
}

type eventDetail struct {
	InstanceIDs []string                `json:"instance-id,omitempty"`
	States      []infrav1.InstanceState `json:"state,omitempty"`
}
