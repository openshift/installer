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

// LogForwarderKind is the name of the type used to represent objects
// of type 'log_forwarder'.
const LogForwarderKind = "LogForwarder"

// LogForwarderLinkKind is the name of the type used to represent links
// to objects of type 'log_forwarder'.
const LogForwarderLinkKind = "LogForwarderLink"

// LogForwarderNilKind is the name of the type used to nil references
// to objects of type 'log_forwarder'.
const LogForwarderNilKind = "LogForwarderNil"

// LogForwarder represents the values of the 'log_forwarder' type.
//
// Representation of a log forwarder configuration for a cluster.
type LogForwarder struct {
	fieldSet_    []bool
	id           string
	href         string
	s3           *LogForwarderS3Config
	applications []string
	cloudwatch   *LogForwarderCloudWatchConfig
	clusterID    string
	groups       []*LogForwarderGroup
	status       *LogForwarderStatus
}

// Kind returns the name of the type of the object.
func (o *LogForwarder) Kind() string {
	if o == nil {
		return LogForwarderNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return LogForwarderLinkKind
	}
	return LogForwarderKind
}

// Link returns true if this is a link.
func (o *LogForwarder) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *LogForwarder) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *LogForwarder) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *LogForwarder) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *LogForwarder) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *LogForwarder) Empty() bool {
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

// S3 returns the value of the 'S3' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// S3 configuration for log forwarding destination.
func (o *LogForwarder) S3() *LogForwarderS3Config {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.s3
	}
	return nil
}

// GetS3 returns the value of the 'S3' attribute and
// a flag indicating if the attribute has a value.
//
// S3 configuration for log forwarding destination.
func (o *LogForwarder) GetS3() (value *LogForwarderS3Config, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.s3
	}
	return
}

// Applications returns the value of the 'applications' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of additional applications to forward logs for.
func (o *LogForwarder) Applications() []string {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.applications
	}
	return nil
}

// GetApplications returns the value of the 'applications' attribute and
// a flag indicating if the attribute has a value.
//
// List of additional applications to forward logs for.
func (o *LogForwarder) GetApplications() (value []string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.applications
	}
	return
}

// Cloudwatch returns the value of the 'cloudwatch' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// CloudWatch configuration for log forwarding destination.
func (o *LogForwarder) Cloudwatch() *LogForwarderCloudWatchConfig {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.cloudwatch
	}
	return nil
}

// GetCloudwatch returns the value of the 'cloudwatch' attribute and
// a flag indicating if the attribute has a value.
//
// CloudWatch configuration for log forwarding destination.
func (o *LogForwarder) GetCloudwatch() (value *LogForwarderCloudWatchConfig, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.cloudwatch
	}
	return
}

// ClusterID returns the value of the 'cluster_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Identifier of the cluster.
func (o *LogForwarder) ClusterID() string {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.clusterID
	}
	return ""
}

// GetClusterID returns the value of the 'cluster_ID' attribute and
// a flag indicating if the attribute has a value.
//
// Identifier of the cluster.
func (o *LogForwarder) GetClusterID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.clusterID
	}
	return
}

// Groups returns the value of the 'groups' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of log forwarder groups.
func (o *LogForwarder) Groups() []*LogForwarderGroup {
	if o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7] {
		return o.groups
	}
	return nil
}

// GetGroups returns the value of the 'groups' attribute and
// a flag indicating if the attribute has a value.
//
// List of log forwarder groups.
func (o *LogForwarder) GetGroups() (value []*LogForwarderGroup, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7]
	if ok {
		value = o.groups
	}
	return
}

// Status returns the value of the 'status' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Status of the log forwarder.
func (o *LogForwarder) Status() *LogForwarderStatus {
	if o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8] {
		return o.status
	}
	return nil
}

// GetStatus returns the value of the 'status' attribute and
// a flag indicating if the attribute has a value.
//
// Status of the log forwarder.
func (o *LogForwarder) GetStatus() (value *LogForwarderStatus, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8]
	if ok {
		value = o.status
	}
	return
}

// LogForwarderListKind is the name of the type used to represent list of objects of
// type 'log_forwarder'.
const LogForwarderListKind = "LogForwarderList"

// LogForwarderListLinkKind is the name of the type used to represent links to list
// of objects of type 'log_forwarder'.
const LogForwarderListLinkKind = "LogForwarderListLink"

// LogForwarderNilKind is the name of the type used to nil lists of objects of
// type 'log_forwarder'.
const LogForwarderListNilKind = "LogForwarderListNil"

// LogForwarderList is a list of values of the 'log_forwarder' type.
type LogForwarderList struct {
	href  string
	link  bool
	items []*LogForwarder
}

// Kind returns the name of the type of the object.
func (l *LogForwarderList) Kind() string {
	if l == nil {
		return LogForwarderListNilKind
	}
	if l.link {
		return LogForwarderListLinkKind
	}
	return LogForwarderListKind
}

// Link returns true iif this is a link.
func (l *LogForwarderList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *LogForwarderList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *LogForwarderList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *LogForwarderList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *LogForwarderList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *LogForwarderList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *LogForwarderList) SetItems(items []*LogForwarder) {
	l.items = items
}

// Items returns the items of the list.
func (l *LogForwarderList) Items() []*LogForwarder {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *LogForwarderList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *LogForwarderList) Get(i int) *LogForwarder {
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
func (l *LogForwarderList) Slice() []*LogForwarder {
	var slice []*LogForwarder
	if l == nil {
		slice = make([]*LogForwarder, 0)
	} else {
		slice = make([]*LogForwarder, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *LogForwarderList) Each(f func(item *LogForwarder) bool) {
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
func (l *LogForwarderList) Range(f func(index int, item *LogForwarder) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
