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

package v1 // github.com/openshift-online/ocm-sdk-go/servicemgmt/v1

// AWSBuilder contains the data and logic needed to build 'AWS' objects.
//
// _Amazon Web Services_ specific settings of a cluster.
type AWSBuilder struct {
	bitmap_         uint32
	sts             *STSBuilder
	accessKeyID     string
	accountID       string
	secretAccessKey string
	subnetIDs       []string
	tags            map[string]string
	privateLink     bool
}

// NewAWS creates a new builder of 'AWS' objects.
func NewAWS() *AWSBuilder {
	return &AWSBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AWSBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// STS sets the value of the 'STS' attribute to the given value.
//
// Contains the necessary attributes to support role-based authentication on AWS.
func (b *AWSBuilder) STS(value *STSBuilder) *AWSBuilder {
	b.sts = value
	if value != nil {
		b.bitmap_ |= 1
	} else {
		b.bitmap_ &^= 1
	}
	return b
}

// AccessKeyID sets the value of the 'access_key_ID' attribute to the given value.
func (b *AWSBuilder) AccessKeyID(value string) *AWSBuilder {
	b.accessKeyID = value
	b.bitmap_ |= 2
	return b
}

// AccountID sets the value of the 'account_ID' attribute to the given value.
func (b *AWSBuilder) AccountID(value string) *AWSBuilder {
	b.accountID = value
	b.bitmap_ |= 4
	return b
}

// PrivateLink sets the value of the 'private_link' attribute to the given value.
func (b *AWSBuilder) PrivateLink(value bool) *AWSBuilder {
	b.privateLink = value
	b.bitmap_ |= 8
	return b
}

// SecretAccessKey sets the value of the 'secret_access_key' attribute to the given value.
func (b *AWSBuilder) SecretAccessKey(value string) *AWSBuilder {
	b.secretAccessKey = value
	b.bitmap_ |= 16
	return b
}

// SubnetIDs sets the value of the 'subnet_IDs' attribute to the given values.
func (b *AWSBuilder) SubnetIDs(values ...string) *AWSBuilder {
	b.subnetIDs = make([]string, len(values))
	copy(b.subnetIDs, values)
	b.bitmap_ |= 32
	return b
}

// Tags sets the value of the 'tags' attribute to the given value.
func (b *AWSBuilder) Tags(value map[string]string) *AWSBuilder {
	b.tags = value
	if value != nil {
		b.bitmap_ |= 64
	} else {
		b.bitmap_ &^= 64
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AWSBuilder) Copy(object *AWS) *AWSBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	if object.sts != nil {
		b.sts = NewSTS().Copy(object.sts)
	} else {
		b.sts = nil
	}
	b.accessKeyID = object.accessKeyID
	b.accountID = object.accountID
	b.privateLink = object.privateLink
	b.secretAccessKey = object.secretAccessKey
	if object.subnetIDs != nil {
		b.subnetIDs = make([]string, len(object.subnetIDs))
		copy(b.subnetIDs, object.subnetIDs)
	} else {
		b.subnetIDs = nil
	}
	if len(object.tags) > 0 {
		b.tags = map[string]string{}
		for k, v := range object.tags {
			b.tags[k] = v
		}
	} else {
		b.tags = nil
	}
	return b
}

// Build creates a 'AWS' object using the configuration stored in the builder.
func (b *AWSBuilder) Build() (object *AWS, err error) {
	object = new(AWS)
	object.bitmap_ = b.bitmap_
	if b.sts != nil {
		object.sts, err = b.sts.Build()
		if err != nil {
			return
		}
	}
	object.accessKeyID = b.accessKeyID
	object.accountID = b.accountID
	object.privateLink = b.privateLink
	object.secretAccessKey = b.secretAccessKey
	if b.subnetIDs != nil {
		object.subnetIDs = make([]string, len(b.subnetIDs))
		copy(object.subnetIDs, b.subnetIDs)
	}
	if b.tags != nil {
		object.tags = make(map[string]string)
		for k, v := range b.tags {
			object.tags[k] = v
		}
	}
	return
}
