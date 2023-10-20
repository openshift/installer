/*
Copyright 2018 The Kubernetes Authors.

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

package converters

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/ssm"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
)

// TagsToMap converts a []*ec2.Tag into a infrav1.Tags.
func TagsToMap(src []*ec2.Tag) infrav1.Tags {
	tags := make(infrav1.Tags, len(src))

	for _, t := range src {
		tags[*t.Key] = *t.Value
	}

	return tags
}

// MapPtrToMap converts a [string]*string into a infrav1.Tags.
func MapPtrToMap(src map[string]*string) infrav1.Tags {
	tags := make(infrav1.Tags, len(src))

	for k, v := range src {
		tags[k] = *v
	}

	return tags
}

// MapToTags converts a infrav1.Tags to a []*ec2.Tag.
func MapToTags(src infrav1.Tags) []*ec2.Tag {
	tags := make([]*ec2.Tag, 0, len(src))

	for k, v := range src {
		tag := &ec2.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		tags = append(tags, tag)
	}

	return tags
}

// ELBTagsToMap converts a []*elb.Tag into a infrav1.Tags.
func ELBTagsToMap(src []*elb.Tag) infrav1.Tags {
	tags := make(infrav1.Tags, len(src))

	for _, t := range src {
		tags[*t.Key] = *t.Value
	}

	return tags
}

// V2TagsToMap converts a []*elbv2.Tag into a infrav1.Tags.
func V2TagsToMap(src []*elbv2.Tag) infrav1.Tags {
	tags := make(infrav1.Tags, len(src))

	for _, t := range src {
		tags[*t.Key] = *t.Value
	}

	return tags
}

// MapToELBTags converts a infrav1.Tags to a []*elb.Tag.
func MapToELBTags(src infrav1.Tags) []*elb.Tag {
	tags := make([]*elb.Tag, 0, len(src))

	for k, v := range src {
		tag := &elb.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		tags = append(tags, tag)
	}

	return tags
}

// MapToV2Tags converts a infrav1.Tags to a []*elbv2.Tag.
func MapToV2Tags(src infrav1.Tags) []*elbv2.Tag {
	tags := make([]*elbv2.Tag, 0, len(src))

	for k, v := range src {
		tag := &elbv2.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		tags = append(tags, tag)
	}

	return tags
}

// MapToSecretsManagerTags converts a infrav1.Tags to a []*secretsmanager.Tag.
func MapToSecretsManagerTags(src infrav1.Tags) []*secretsmanager.Tag {
	tags := make([]*secretsmanager.Tag, 0, len(src))

	for k, v := range src {
		tag := &secretsmanager.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		tags = append(tags, tag)
	}

	return tags
}

// MapToSSMTags converts a infrav1.Tags to a []*ssm.Tag.
func MapToSSMTags(src infrav1.Tags) []*ssm.Tag {
	tags := make([]*ssm.Tag, 0, len(src))

	for k, v := range src {
		tag := &ssm.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		tags = append(tags, tag)
	}

	return tags
}

// MapToIAMTags converts a infrav1.Tags to a []*iam.Tag.
func MapToIAMTags(src infrav1.Tags) []*iam.Tag {
	tags := make([]*iam.Tag, 0, len(src))

	for k, v := range src {
		tag := &iam.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		tags = append(tags, tag)
	}

	return tags
}

// ASGTagsToMap converts a []*autoscaling.TagDescription into a infrav1.Tags.
func ASGTagsToMap(src []*autoscaling.TagDescription) infrav1.Tags {
	tags := make(infrav1.Tags, len(src))

	for _, t := range src {
		tags[*t.Key] = *t.Value
	}

	return tags
}
