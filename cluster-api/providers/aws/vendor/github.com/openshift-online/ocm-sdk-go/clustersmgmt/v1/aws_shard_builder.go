/*
Copyright (c) 2020 Red Hat, Inc.

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

// IMPORTANT: This file has been generated automatically, refrain from modifying it manually as all
// your changes will be lost when the file is generated again.

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

// AWSShardBuilder contains the data and logic needed to build 'AWS_shard' objects.
//
// Config for AWS provision shards
type AWSShardBuilder struct {
	bitmap_           uint32
	ecrRepositoryURLs []string
}

// NewAWSShard creates a new builder of 'AWS_shard' objects.
func NewAWSShard() *AWSShardBuilder {
	return &AWSShardBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AWSShardBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// ECRRepositoryURLs sets the value of the 'ECR_repository_URLs' attribute to the given values.
func (b *AWSShardBuilder) ECRRepositoryURLs(values ...string) *AWSShardBuilder {
	b.ecrRepositoryURLs = make([]string, len(values))
	copy(b.ecrRepositoryURLs, values)
	b.bitmap_ |= 1
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AWSShardBuilder) Copy(object *AWSShard) *AWSShardBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	if object.ecrRepositoryURLs != nil {
		b.ecrRepositoryURLs = make([]string, len(object.ecrRepositoryURLs))
		copy(b.ecrRepositoryURLs, object.ecrRepositoryURLs)
	} else {
		b.ecrRepositoryURLs = nil
	}
	return b
}

// Build creates a 'AWS_shard' object using the configuration stored in the builder.
func (b *AWSShardBuilder) Build() (object *AWSShard, err error) {
	object = new(AWSShard)
	object.bitmap_ = b.bitmap_
	if b.ecrRepositoryURLs != nil {
		object.ecrRepositoryURLs = make([]string, len(b.ecrRepositoryURLs))
		copy(object.ecrRepositoryURLs, b.ecrRepositoryURLs)
	}
	return
}
