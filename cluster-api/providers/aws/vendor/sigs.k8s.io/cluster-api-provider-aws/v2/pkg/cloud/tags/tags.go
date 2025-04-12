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

// Package tags provides a way to tag cloud resources.
package tags

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/pkg/errors"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
)

const (
	// AwsInternalTagPrefix is the prefix for AWS internal tags, which are reserved for internal AWS use.
	AwsInternalTagPrefix = "aws:"
)

var (
	// ErrBuildParamsRequired defines an error for when no build params are supplied.
	ErrBuildParamsRequired = errors.New("no build params supplied")

	// ErrApplyFuncRequired defines an error for when tags are not supplied.
	ErrApplyFuncRequired = errors.New("no tags apply function supplied")
)

// BuilderOption represents an option when creating a tags builder.
type BuilderOption func(*Builder)

// Builder is the interface for a tags builder.
type Builder struct {
	params    *infrav1.BuildParams
	applyFunc func(params *infrav1.BuildParams) error
}

// New creates a new TagsBuilder with the specified build parameters
// and with optional configuration.
func New(params *infrav1.BuildParams, opts ...BuilderOption) *Builder {
	builder := &Builder{
		params: params,
	}

	for _, opt := range opts {
		opt(builder)
	}

	return builder
}

// Apply tags a resource with tags including the cluster tag.
func (b *Builder) Apply() error {
	if b.params == nil {
		return ErrBuildParamsRequired
	}
	if b.applyFunc == nil {
		return ErrApplyFuncRequired
	}

	if err := b.applyFunc(b.params); err != nil {
		return fmt.Errorf("failed applying tags: %w", err)
	}
	return nil
}

// Ensure applies the tags if the current tags differ from the params.
func (b *Builder) Ensure(current infrav1.Tags) error {
	if b.params == nil {
		return ErrBuildParamsRequired
	}
	if diff := computeDiff(current, *b.params); len(diff) > 0 {
		return b.Apply()
	}
	return nil
}

// WithEC2 is used to denote that the tags builder will be using EC2.
func WithEC2(ec2client ec2iface.EC2API) BuilderOption {
	return func(b *Builder) {
		b.applyFunc = func(params *infrav1.BuildParams) error {
			tags := infrav1.Build(*params)
			awsTags := make([]*ec2.Tag, 0, len(tags))

			// For testing, we need sorted keys
			sortedKeys := make([]string, 0, len(tags))
			for k := range tags {
				// We want to filter out the tag keys that start with `aws:` as they are reserved for internal AWS use.
				if !strings.HasPrefix(k, AwsInternalTagPrefix) {
					sortedKeys = append(sortedKeys, k)
				}
			}
			sort.Strings(sortedKeys)

			for _, key := range sortedKeys {
				tag := &ec2.Tag{
					Key:   aws.String(key),
					Value: aws.String(tags[key]),
				}
				awsTags = append(awsTags, tag)
			}

			createTagsInput := &ec2.CreateTagsInput{
				Resources: aws.StringSlice([]string{params.ResourceID}),
				Tags:      awsTags,
			}

			_, err := ec2client.CreateTagsWithContext(context.TODO(), createTagsInput)
			return errors.Wrapf(err, "failed to tag resource %q in cluster %q", params.ResourceID, params.ClusterName)
		}
	}
}

// WithEKS is used to specify that the tags builder will be targeting EKS.
func WithEKS(eksclient eksiface.EKSAPI) BuilderOption {
	return func(b *Builder) {
		b.applyFunc = func(params *infrav1.BuildParams) error {
			tags := infrav1.Build(*params)
			eksTags := make(map[string]*string, len(tags))
			for k, v := range tags {
				// We want to filter out the tag keys that start with `aws:` as they are reserved for internal AWS use.
				if !strings.HasPrefix(k, AwsInternalTagPrefix) {
					eksTags[k] = aws.String(v)
				}
			}

			tagResourcesInput := &eks.TagResourceInput{
				ResourceArn: aws.String(params.ResourceID),
				Tags:        eksTags,
			}

			_, err := eksclient.TagResource(tagResourcesInput)
			if err != nil {
				return errors.Wrapf(err, "failed to tag eks cluster %q in cluster %q", params.ResourceID, params.ClusterName)
			}

			return nil
		}
	}
}

func computeDiff(current infrav1.Tags, buildParams infrav1.BuildParams) infrav1.Tags {
	want := infrav1.Build(buildParams)

	// Some tags could be external set by some external entities
	// and that means even if there is no change in cluster
	// managed tags, tags would be updated as "current" and
	// "want" would be different due to external tags.
	// This fix makes sure that tags are updated only if
	// there is a change in cluster managed tags.
	return want.Difference(current)
}

// BuildParamsToTagSpecification builds a TagSpecification for the specified resource type.
func BuildParamsToTagSpecification(ec2ResourceType string, params infrav1.BuildParams) *ec2.TagSpecification {
	tags := infrav1.Build(params)

	tagSpec := &ec2.TagSpecification{ResourceType: aws.String(ec2ResourceType)}

	// For testing, we need sorted keys
	sortedKeys := make([]string, 0, len(tags))
	for k := range tags {
		sortedKeys = append(sortedKeys, k)
	}

	sort.Strings(sortedKeys)

	for _, key := range sortedKeys {
		tagSpec.Tags = append(tagSpec.Tags, &ec2.Tag{
			Key:   aws.String(key),
			Value: aws.String(tags[key]),
		})
	}

	return tagSpec
}
