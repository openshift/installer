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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/servicemgmt/v1

// _Amazon Web Services_ specific settings of a cluster.
type AWSBuilder struct {
	fieldSet_       []bool
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
	return &AWSBuilder{
		fieldSet_: make([]bool, 7),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AWSBuilder) Empty() bool {
	if b == nil || len(b.fieldSet_) == 0 {
		return true
	}
	for _, set := range b.fieldSet_ {
		if set {
			return false
		}
	}
	return true
}

// STS sets the value of the 'STS' attribute to the given value.
//
// Contains the necessary attributes to support role-based authentication on AWS.
func (b *AWSBuilder) STS(value *STSBuilder) *AWSBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.sts = value
	if value != nil {
		b.fieldSet_[0] = true
	} else {
		b.fieldSet_[0] = false
	}
	return b
}

// AccessKeyID sets the value of the 'access_key_ID' attribute to the given value.
func (b *AWSBuilder) AccessKeyID(value string) *AWSBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.accessKeyID = value
	b.fieldSet_[1] = true
	return b
}

// AccountID sets the value of the 'account_ID' attribute to the given value.
func (b *AWSBuilder) AccountID(value string) *AWSBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.accountID = value
	b.fieldSet_[2] = true
	return b
}

// PrivateLink sets the value of the 'private_link' attribute to the given value.
func (b *AWSBuilder) PrivateLink(value bool) *AWSBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.privateLink = value
	b.fieldSet_[3] = true
	return b
}

// SecretAccessKey sets the value of the 'secret_access_key' attribute to the given value.
func (b *AWSBuilder) SecretAccessKey(value string) *AWSBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.secretAccessKey = value
	b.fieldSet_[4] = true
	return b
}

// SubnetIDs sets the value of the 'subnet_IDs' attribute to the given values.
func (b *AWSBuilder) SubnetIDs(values ...string) *AWSBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.subnetIDs = make([]string, len(values))
	copy(b.subnetIDs, values)
	b.fieldSet_[5] = true
	return b
}

// Tags sets the value of the 'tags' attribute to the given value.
func (b *AWSBuilder) Tags(value map[string]string) *AWSBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.tags = value
	if value != nil {
		b.fieldSet_[6] = true
	} else {
		b.fieldSet_[6] = false
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AWSBuilder) Copy(object *AWS) *AWSBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
