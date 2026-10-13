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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1

// LogForwarderS3Config represents the values of the 'log_forwarder_S3_config' type.
//
// S3 configuration for log forwarding.
type LogForwarderS3Config struct {
	fieldSet_    []bool
	bucketName   string
	bucketPrefix string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *LogForwarderS3Config) Empty() bool {
	if o == nil || len(o.fieldSet_) == 0 {
		return true
	}
	for _, set := range o.fieldSet_ {
		if set {
			return false
		}
	}
	return true
}

// BucketName returns the value of the 'bucket_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The name of the S3 bucket.
func (o *LogForwarderS3Config) BucketName() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.bucketName
	}
	return ""
}

// GetBucketName returns the value of the 'bucket_name' attribute and
// a flag indicating if the attribute has a value.
//
// The name of the S3 bucket.
func (o *LogForwarderS3Config) GetBucketName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.bucketName
	}
	return
}

// BucketPrefix returns the value of the 'bucket_prefix' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The prefix to use for objects stored in the S3 bucket.
func (o *LogForwarderS3Config) BucketPrefix() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.bucketPrefix
	}
	return ""
}

// GetBucketPrefix returns the value of the 'bucket_prefix' attribute and
// a flag indicating if the attribute has a value.
//
// The prefix to use for objects stored in the S3 bucket.
func (o *LogForwarderS3Config) GetBucketPrefix() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.bucketPrefix
	}
	return
}

// LogForwarderS3ConfigListKind is the name of the type used to represent list of objects of
// type 'log_forwarder_S3_config'.
const LogForwarderS3ConfigListKind = "LogForwarderS3ConfigList"

// LogForwarderS3ConfigListLinkKind is the name of the type used to represent links to list
// of objects of type 'log_forwarder_S3_config'.
const LogForwarderS3ConfigListLinkKind = "LogForwarderS3ConfigListLink"

// LogForwarderS3ConfigNilKind is the name of the type used to nil lists of objects of
// type 'log_forwarder_S3_config'.
const LogForwarderS3ConfigListNilKind = "LogForwarderS3ConfigListNil"

// LogForwarderS3ConfigList is a list of values of the 'log_forwarder_S3_config' type.
type LogForwarderS3ConfigList struct {
	href  string
	link  bool
	items []*LogForwarderS3Config
}

// Len returns the length of the list.
func (l *LogForwarderS3ConfigList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *LogForwarderS3ConfigList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *LogForwarderS3ConfigList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *LogForwarderS3ConfigList) SetItems(items []*LogForwarderS3Config) {
	l.items = items
}

// Items returns the items of the list.
func (l *LogForwarderS3ConfigList) Items() []*LogForwarderS3Config {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *LogForwarderS3ConfigList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *LogForwarderS3ConfigList) Get(i int) *LogForwarderS3Config {
	if l == nil || i < 0 || i >= len(l.items) {
		return nil
	}
	return l.items[i]
}

// Slice returns an slice containing the items of the list. The returned slice is a
// copy of the one used internally, so it can be modified without affecting the
// internal representation.
//
// If you don't need to modify the returned slice consider using the Each or Range
// functions, as they don't need to allocate a new slice.
func (l *LogForwarderS3ConfigList) Slice() []*LogForwarderS3Config {
	var slice []*LogForwarderS3Config
	if l == nil {
		slice = make([]*LogForwarderS3Config, 0)
	} else {
		slice = make([]*LogForwarderS3Config, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *LogForwarderS3ConfigList) Each(f func(item *LogForwarderS3Config) bool) {
	if l == nil {
		return
	}
	for _, item := range l.items {
		if !f(item) {
			break
		}
	}
}

// Range runs the given function for each index and item of the list, in order. If
// the function returns false the iteration stops, otherwise it continues till all
// the elements of the list have been processed.
func (l *LogForwarderS3ConfigList) Range(f func(index int, item *LogForwarderS3Config) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
