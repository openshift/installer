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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/jobqueue/v1

import (
	time "time"
)

// JobKind is the name of the type used to represent objects
// of type 'job'.
const JobKind = "Job"

// JobLinkKind is the name of the type used to represent links
// to objects of type 'job'.
const JobLinkKind = "JobLink"

// JobNilKind is the name of the type used to nil references
// to objects of type 'job'.
const JobNilKind = "JobNil"

// Job represents the values of the 'job' type.
//
// This struct is a job in a Job Queue.
type Job struct {
	fieldSet_   []bool
	id          string
	href        string
	abandonedAt time.Time
	arguments   string
	attempts    int
	createdAt   time.Time
	receiptId   string
	updatedAt   time.Time
}

// Kind returns the name of the type of the object.
func (o *Job) Kind() string {
	if o == nil {
		return JobNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return JobLinkKind
	}
	return JobKind
}

// Link returns true if this is a link.
func (o *Job) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *Job) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *Job) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *Job) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *Job) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *Job) Empty() bool {
	if o == nil || len(o.fieldSet_) == 0 {
		return true
	}

	// Check all fields except the link flag (index 0)
	for i := 1; i < len(o.fieldSet_); i++ {
		if o.fieldSet_[i] {
			return false
		}
	}
	return true
}

// AbandonedAt returns the value of the 'abandoned_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// DLQ sent timestamp
func (o *Job) AbandonedAt() time.Time {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.abandonedAt
	}
	return time.Time{}
}

// GetAbandonedAt returns the value of the 'abandoned_at' attribute and
// a flag indicating if the attribute has a value.
//
// DLQ sent timestamp
func (o *Job) GetAbandonedAt() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.abandonedAt
	}
	return
}

// Arguments returns the value of the 'arguments' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Arguments to run Job with.
func (o *Job) Arguments() string {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.arguments
	}
	return ""
}

// GetArguments returns the value of the 'arguments' attribute and
// a flag indicating if the attribute has a value.
//
// Arguments to run Job with.
func (o *Job) GetArguments() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.arguments
	}
	return
}

// Attempts returns the value of the 'attempts' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Number of retries.
func (o *Job) Attempts() int {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.attempts
	}
	return 0
}

// GetAttempts returns the value of the 'attempts' attribute and
// a flag indicating if the attribute has a value.
//
// Number of retries.
func (o *Job) GetAttempts() (value int, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.attempts
	}
	return
}

// CreatedAt returns the value of the 'created_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Job) CreatedAt() time.Time {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.createdAt
	}
	return time.Time{}
}

// GetCreatedAt returns the value of the 'created_at' attribute and
// a flag indicating if the attribute has a value.
func (o *Job) GetCreatedAt() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.createdAt
	}
	return
}

// ReceiptId returns the value of the 'receipt_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Each time a specific job is pop'd, the receiptId will change, while the ID stays the same.
func (o *Job) ReceiptId() string {
	if o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7] {
		return o.receiptId
	}
	return ""
}

// GetReceiptId returns the value of the 'receipt_id' attribute and
// a flag indicating if the attribute has a value.
//
// Each time a specific job is pop'd, the receiptId will change, while the ID stays the same.
func (o *Job) GetReceiptId() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7]
	if ok {
		value = o.receiptId
	}
	return
}

// UpdatedAt returns the value of the 'updated_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Job) UpdatedAt() time.Time {
	if o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8] {
		return o.updatedAt
	}
	return time.Time{}
}

// GetUpdatedAt returns the value of the 'updated_at' attribute and
// a flag indicating if the attribute has a value.
func (o *Job) GetUpdatedAt() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8]
	if ok {
		value = o.updatedAt
	}
	return
}

// JobListKind is the name of the type used to represent list of objects of
// type 'job'.
const JobListKind = "JobList"

// JobListLinkKind is the name of the type used to represent links to list
// of objects of type 'job'.
const JobListLinkKind = "JobListLink"

// JobNilKind is the name of the type used to nil lists of objects of
// type 'job'.
const JobListNilKind = "JobListNil"

// JobList is a list of values of the 'job' type.
type JobList struct {
	href  string
	link  bool
	items []*Job
}

// Kind returns the name of the type of the object.
func (l *JobList) Kind() string {
	if l == nil {
		return JobListNilKind
	}
	if l.link {
		return JobListLinkKind
	}
	return JobListKind
}

// Link returns true iif this is a link.
func (l *JobList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *JobList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *JobList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *JobList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *JobList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *JobList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *JobList) SetItems(items []*Job) {
	l.items = items
}

// Items returns the items of the list.
func (l *JobList) Items() []*Job {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *JobList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *JobList) Get(i int) *Job {
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
func (l *JobList) Slice() []*Job {
	var slice []*Job
	if l == nil {
		slice = make([]*Job, 0)
	} else {
		slice = make([]*Job, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *JobList) Each(f func(item *Job) bool) {
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
func (l *JobList) Range(f func(index int, item *Job) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
