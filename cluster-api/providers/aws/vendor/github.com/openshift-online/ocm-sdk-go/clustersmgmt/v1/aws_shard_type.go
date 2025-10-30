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

// AWSShard represents the values of the 'AWS_shard' type.
//
// Config for AWS provision shards
type AWSShard struct {
	bitmap_           uint32
	ecrRepositoryURLs []string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AWSShard) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// ECRRepositoryURLs returns the value of the 'ECR_repository_URLs' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// ECR repository URLs of the provision shard
func (o *AWSShard) ECRRepositoryURLs() []string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.ecrRepositoryURLs
	}
	return nil
}

// GetECRRepositoryURLs returns the value of the 'ECR_repository_URLs' attribute and
// a flag indicating if the attribute has a value.
//
// ECR repository URLs of the provision shard
func (o *AWSShard) GetECRRepositoryURLs() (value []string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.ecrRepositoryURLs
	}
	return
}

// AWSShardListKind is the name of the type used to represent list of objects of
// type 'AWS_shard'.
const AWSShardListKind = "AWSShardList"

// AWSShardListLinkKind is the name of the type used to represent links to list
// of objects of type 'AWS_shard'.
const AWSShardListLinkKind = "AWSShardListLink"

// AWSShardNilKind is the name of the type used to nil lists of objects of
// type 'AWS_shard'.
const AWSShardListNilKind = "AWSShardListNil"

// AWSShardList is a list of values of the 'AWS_shard' type.
type AWSShardList struct {
	href  string
	link  bool
	items []*AWSShard
}

// Len returns the length of the list.
func (l *AWSShardList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AWSShardList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AWSShardList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AWSShardList) SetItems(items []*AWSShard) {
	l.items = items
}

// Items returns the items of the list.
func (l *AWSShardList) Items() []*AWSShard {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AWSShardList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AWSShardList) Get(i int) *AWSShard {
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
func (l *AWSShardList) Slice() []*AWSShard {
	var slice []*AWSShard
	if l == nil {
		slice = make([]*AWSShard, 0)
	} else {
		slice = make([]*AWSShard, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AWSShardList) Each(f func(item *AWSShard) bool) {
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
func (l *AWSShardList) Range(f func(index int, item *AWSShard) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
