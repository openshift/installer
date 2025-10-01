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

package v1 // github.com/openshift-online/ocm-sdk-go/servicelogs/v1

import (
	time "time"
)

// LogEntryKind is the name of the type used to represent objects
// of type 'log_entry'.
const LogEntryKind = "LogEntry"

// LogEntryLinkKind is the name of the type used to represent links
// to objects of type 'log_entry'.
const LogEntryLinkKind = "LogEntryLink"

// LogEntryNilKind is the name of the type used to nil references
// to objects of type 'log_entry'.
const LogEntryNilKind = "LogEntryNil"

// LogEntry represents the values of the 'log_entry' type.
type LogEntry struct {
	bitmap_        uint32
	id             string
	href           string
	clusterID      string
	clusterUUID    string
	createdAt      time.Time
	createdBy      string
	description    string
	docReferences  []string
	eventStreamID  string
	logType        LogType
	serviceName    string
	severity       Severity
	subscriptionID string
	summary        string
	timestamp      time.Time
	username       string
	internalOnly   bool
}

// Kind returns the name of the type of the object.
func (o *LogEntry) Kind() string {
	if o == nil {
		return LogEntryNilKind
	}
	if o.bitmap_&1 != 0 {
		return LogEntryLinkKind
	}
	return LogEntryKind
}

// Link returns true if this is a link.
func (o *LogEntry) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *LogEntry) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *LogEntry) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *LogEntry) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *LogEntry) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *LogEntry) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// ClusterID returns the value of the 'cluster_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Internal cluster ID.
func (o *LogEntry) ClusterID() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.clusterID
	}
	return ""
}

// GetClusterID returns the value of the 'cluster_ID' attribute and
// a flag indicating if the attribute has a value.
//
// Internal cluster ID.
func (o *LogEntry) GetClusterID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.clusterID
	}
	return
}

// ClusterUUID returns the value of the 'cluster_UUID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// External cluster ID.
func (o *LogEntry) ClusterUUID() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.clusterUUID
	}
	return ""
}

// GetClusterUUID returns the value of the 'cluster_UUID' attribute and
// a flag indicating if the attribute has a value.
//
// External cluster ID.
func (o *LogEntry) GetClusterUUID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.clusterUUID
	}
	return
}

// CreatedAt returns the value of the 'created_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The time at which the cluster log was created.
func (o *LogEntry) CreatedAt() time.Time {
	if o != nil && o.bitmap_&32 != 0 {
		return o.createdAt
	}
	return time.Time{}
}

// GetCreatedAt returns the value of the 'created_at' attribute and
// a flag indicating if the attribute has a value.
//
// The time at which the cluster log was created.
func (o *LogEntry) GetCreatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.createdAt
	}
	return
}

// CreatedBy returns the value of the 'created_by' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The name of the user who created the cluster log.
func (o *LogEntry) CreatedBy() string {
	if o != nil && o.bitmap_&64 != 0 {
		return o.createdBy
	}
	return ""
}

// GetCreatedBy returns the value of the 'created_by' attribute and
// a flag indicating if the attribute has a value.
//
// The name of the user who created the cluster log.
func (o *LogEntry) GetCreatedBy() (value string, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.createdBy
	}
	return
}

// Description returns the value of the 'description' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Full description of the log entry content (supports Markdown format as well).
func (o *LogEntry) Description() string {
	if o != nil && o.bitmap_&128 != 0 {
		return o.description
	}
	return ""
}

// GetDescription returns the value of the 'description' attribute and
// a flag indicating if the attribute has a value.
//
// Full description of the log entry content (supports Markdown format as well).
func (o *LogEntry) GetDescription() (value string, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.description
	}
	return
}

// DocReferences returns the value of the 'doc_references' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The list of documentation references (i.e links) contained in the event.
func (o *LogEntry) DocReferences() []string {
	if o != nil && o.bitmap_&256 != 0 {
		return o.docReferences
	}
	return nil
}

// GetDocReferences returns the value of the 'doc_references' attribute and
// a flag indicating if the attribute has a value.
//
// The list of documentation references (i.e links) contained in the event.
func (o *LogEntry) GetDocReferences() (value []string, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.docReferences
	}
	return
}

// EventStreamID returns the value of the 'event_stream_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Log custom event id for a simple search of related cluster logs.
func (o *LogEntry) EventStreamID() string {
	if o != nil && o.bitmap_&512 != 0 {
		return o.eventStreamID
	}
	return ""
}

// GetEventStreamID returns the value of the 'event_stream_ID' attribute and
// a flag indicating if the attribute has a value.
//
// Log custom event id for a simple search of related cluster logs.
func (o *LogEntry) GetEventStreamID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&512 != 0
	if ok {
		value = o.eventStreamID
	}
	return
}

// InternalOnly returns the value of the 'internal_only' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// A flag that indicates whether the log entry should be internal/private only.
func (o *LogEntry) InternalOnly() bool {
	if o != nil && o.bitmap_&1024 != 0 {
		return o.internalOnly
	}
	return false
}

// GetInternalOnly returns the value of the 'internal_only' attribute and
// a flag indicating if the attribute has a value.
//
// A flag that indicates whether the log entry should be internal/private only.
func (o *LogEntry) GetInternalOnly() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&1024 != 0
	if ok {
		value = o.internalOnly
	}
	return
}

// LogType returns the value of the 'log_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Type of the service log entry.
func (o *LogEntry) LogType() LogType {
	if o != nil && o.bitmap_&2048 != 0 {
		return o.logType
	}
	return LogType("")
}

// GetLogType returns the value of the 'log_type' attribute and
// a flag indicating if the attribute has a value.
//
// Type of the service log entry.
func (o *LogEntry) GetLogType() (value LogType, ok bool) {
	ok = o != nil && o.bitmap_&2048 != 0
	if ok {
		value = o.logType
	}
	return
}

// ServiceName returns the value of the 'service_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The name of the service who created the log.
func (o *LogEntry) ServiceName() string {
	if o != nil && o.bitmap_&4096 != 0 {
		return o.serviceName
	}
	return ""
}

// GetServiceName returns the value of the 'service_name' attribute and
// a flag indicating if the attribute has a value.
//
// The name of the service who created the log.
func (o *LogEntry) GetServiceName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4096 != 0
	if ok {
		value = o.serviceName
	}
	return
}

// Severity returns the value of the 'severity' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Log severity for the specific log entry.
func (o *LogEntry) Severity() Severity {
	if o != nil && o.bitmap_&8192 != 0 {
		return o.severity
	}
	return Severity("")
}

// GetSeverity returns the value of the 'severity' attribute and
// a flag indicating if the attribute has a value.
//
// Log severity for the specific log entry.
func (o *LogEntry) GetSeverity() (value Severity, ok bool) {
	ok = o != nil && o.bitmap_&8192 != 0
	if ok {
		value = o.severity
	}
	return
}

// SubscriptionID returns the value of the 'subscription_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The related subscription id of the cluster.
func (o *LogEntry) SubscriptionID() string {
	if o != nil && o.bitmap_&16384 != 0 {
		return o.subscriptionID
	}
	return ""
}

// GetSubscriptionID returns the value of the 'subscription_ID' attribute and
// a flag indicating if the attribute has a value.
//
// The related subscription id of the cluster.
func (o *LogEntry) GetSubscriptionID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16384 != 0
	if ok {
		value = o.subscriptionID
	}
	return
}

// Summary returns the value of the 'summary' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Title of the log entry.
func (o *LogEntry) Summary() string {
	if o != nil && o.bitmap_&32768 != 0 {
		return o.summary
	}
	return ""
}

// GetSummary returns the value of the 'summary' attribute and
// a flag indicating if the attribute has a value.
//
// Title of the log entry.
func (o *LogEntry) GetSummary() (value string, ok bool) {
	ok = o != nil && o.bitmap_&32768 != 0
	if ok {
		value = o.summary
	}
	return
}

// Timestamp returns the value of the 'timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *LogEntry) Timestamp() time.Time {
	if o != nil && o.bitmap_&65536 != 0 {
		return o.timestamp
	}
	return time.Time{}
}

// GetTimestamp returns the value of the 'timestamp' attribute and
// a flag indicating if the attribute has a value.
func (o *LogEntry) GetTimestamp() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&65536 != 0
	if ok {
		value = o.timestamp
	}
	return
}

// Username returns the value of the 'username' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The username that triggered the event (if available).
func (o *LogEntry) Username() string {
	if o != nil && o.bitmap_&131072 != 0 {
		return o.username
	}
	return ""
}

// GetUsername returns the value of the 'username' attribute and
// a flag indicating if the attribute has a value.
//
// The username that triggered the event (if available).
func (o *LogEntry) GetUsername() (value string, ok bool) {
	ok = o != nil && o.bitmap_&131072 != 0
	if ok {
		value = o.username
	}
	return
}

// LogEntryListKind is the name of the type used to represent list of objects of
// type 'log_entry'.
const LogEntryListKind = "LogEntryList"

// LogEntryListLinkKind is the name of the type used to represent links to list
// of objects of type 'log_entry'.
const LogEntryListLinkKind = "LogEntryListLink"

// LogEntryNilKind is the name of the type used to nil lists of objects of
// type 'log_entry'.
const LogEntryListNilKind = "LogEntryListNil"

// LogEntryList is a list of values of the 'log_entry' type.
type LogEntryList struct {
	href  string
	link  bool
	items []*LogEntry
}

// Kind returns the name of the type of the object.
func (l *LogEntryList) Kind() string {
	if l == nil {
		return LogEntryListNilKind
	}
	if l.link {
		return LogEntryListLinkKind
	}
	return LogEntryListKind
}

// Link returns true iif this is a link.
func (l *LogEntryList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *LogEntryList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *LogEntryList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *LogEntryList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *LogEntryList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *LogEntryList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *LogEntryList) SetItems(items []*LogEntry) {
	l.items = items
}

// Items returns the items of the list.
func (l *LogEntryList) Items() []*LogEntry {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *LogEntryList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *LogEntryList) Get(i int) *LogEntry {
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
func (l *LogEntryList) Slice() []*LogEntry {
	var slice []*LogEntry
	if l == nil {
		slice = make([]*LogEntry, 0)
	} else {
		slice = make([]*LogEntry, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *LogEntryList) Each(f func(item *LogEntry) bool) {
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
func (l *LogEntryList) Range(f func(index int, item *LogEntry) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
