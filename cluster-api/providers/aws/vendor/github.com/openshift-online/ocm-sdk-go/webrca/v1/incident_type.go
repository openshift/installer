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

package v1 // github.com/openshift-online/ocm-sdk-go/webrca/v1

import (
	time "time"
)

// IncidentKind is the name of the type used to represent objects
// of type 'incident'.
const IncidentKind = "Incident"

// IncidentLinkKind is the name of the type used to represent links
// to objects of type 'incident'.
const IncidentLinkKind = "IncidentLink"

// IncidentNilKind is the name of the type used to nil references
// to objects of type 'incident'.
const IncidentNilKind = "IncidentNil"

// Incident represents the values of the 'incident' type.
//
// Definition of a Web RCA incident.
type Incident struct {
	bitmap_              uint32
	id                   string
	href                 string
	createdAt            time.Time
	creatorId            string
	deletedAt            time.Time
	description          string
	externalCoordination []string
	incidentId           string
	incidentType         string
	lastUpdated          time.Time
	primaryTeam          string
	severity             string
	status               string
	summary              string
	updatedAt            time.Time
	workedAt             time.Time
}

// Kind returns the name of the type of the object.
func (o *Incident) Kind() string {
	if o == nil {
		return IncidentNilKind
	}
	if o.bitmap_&1 != 0 {
		return IncidentLinkKind
	}
	return IncidentKind
}

// Link returns true if this is a link.
func (o *Incident) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *Incident) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *Incident) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *Incident) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *Incident) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *Incident) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// CreatedAt returns the value of the 'created_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Object creation timestamp.
func (o *Incident) CreatedAt() time.Time {
	if o != nil && o.bitmap_&8 != 0 {
		return o.createdAt
	}
	return time.Time{}
}

// GetCreatedAt returns the value of the 'created_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object creation timestamp.
func (o *Incident) GetCreatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.createdAt
	}
	return
}

// CreatorId returns the value of the 'creator_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Incident) CreatorId() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.creatorId
	}
	return ""
}

// GetCreatorId returns the value of the 'creator_id' attribute and
// a flag indicating if the attribute has a value.
func (o *Incident) GetCreatorId() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.creatorId
	}
	return
}

// DeletedAt returns the value of the 'deleted_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Object deletion timestamp.
func (o *Incident) DeletedAt() time.Time {
	if o != nil && o.bitmap_&32 != 0 {
		return o.deletedAt
	}
	return time.Time{}
}

// GetDeletedAt returns the value of the 'deleted_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object deletion timestamp.
func (o *Incident) GetDeletedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.deletedAt
	}
	return
}

// Description returns the value of the 'description' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Incident) Description() string {
	if o != nil && o.bitmap_&64 != 0 {
		return o.description
	}
	return ""
}

// GetDescription returns the value of the 'description' attribute and
// a flag indicating if the attribute has a value.
func (o *Incident) GetDescription() (value string, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.description
	}
	return
}

// ExternalCoordination returns the value of the 'external_coordination' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Incident) ExternalCoordination() []string {
	if o != nil && o.bitmap_&128 != 0 {
		return o.externalCoordination
	}
	return nil
}

// GetExternalCoordination returns the value of the 'external_coordination' attribute and
// a flag indicating if the attribute has a value.
func (o *Incident) GetExternalCoordination() (value []string, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.externalCoordination
	}
	return
}

// IncidentId returns the value of the 'incident_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Incident) IncidentId() string {
	if o != nil && o.bitmap_&256 != 0 {
		return o.incidentId
	}
	return ""
}

// GetIncidentId returns the value of the 'incident_id' attribute and
// a flag indicating if the attribute has a value.
func (o *Incident) GetIncidentId() (value string, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.incidentId
	}
	return
}

// IncidentType returns the value of the 'incident_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Incident) IncidentType() string {
	if o != nil && o.bitmap_&512 != 0 {
		return o.incidentType
	}
	return ""
}

// GetIncidentType returns the value of the 'incident_type' attribute and
// a flag indicating if the attribute has a value.
func (o *Incident) GetIncidentType() (value string, ok bool) {
	ok = o != nil && o.bitmap_&512 != 0
	if ok {
		value = o.incidentType
	}
	return
}

// LastUpdated returns the value of the 'last_updated' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Incident) LastUpdated() time.Time {
	if o != nil && o.bitmap_&1024 != 0 {
		return o.lastUpdated
	}
	return time.Time{}
}

// GetLastUpdated returns the value of the 'last_updated' attribute and
// a flag indicating if the attribute has a value.
func (o *Incident) GetLastUpdated() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&1024 != 0
	if ok {
		value = o.lastUpdated
	}
	return
}

// PrimaryTeam returns the value of the 'primary_team' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Incident) PrimaryTeam() string {
	if o != nil && o.bitmap_&2048 != 0 {
		return o.primaryTeam
	}
	return ""
}

// GetPrimaryTeam returns the value of the 'primary_team' attribute and
// a flag indicating if the attribute has a value.
func (o *Incident) GetPrimaryTeam() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2048 != 0
	if ok {
		value = o.primaryTeam
	}
	return
}

// Severity returns the value of the 'severity' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Incident) Severity() string {
	if o != nil && o.bitmap_&4096 != 0 {
		return o.severity
	}
	return ""
}

// GetSeverity returns the value of the 'severity' attribute and
// a flag indicating if the attribute has a value.
func (o *Incident) GetSeverity() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4096 != 0
	if ok {
		value = o.severity
	}
	return
}

// Status returns the value of the 'status' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Incident) Status() string {
	if o != nil && o.bitmap_&8192 != 0 {
		return o.status
	}
	return ""
}

// GetStatus returns the value of the 'status' attribute and
// a flag indicating if the attribute has a value.
func (o *Incident) GetStatus() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8192 != 0
	if ok {
		value = o.status
	}
	return
}

// Summary returns the value of the 'summary' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Incident) Summary() string {
	if o != nil && o.bitmap_&16384 != 0 {
		return o.summary
	}
	return ""
}

// GetSummary returns the value of the 'summary' attribute and
// a flag indicating if the attribute has a value.
func (o *Incident) GetSummary() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16384 != 0
	if ok {
		value = o.summary
	}
	return
}

// UpdatedAt returns the value of the 'updated_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Object modification timestamp.
func (o *Incident) UpdatedAt() time.Time {
	if o != nil && o.bitmap_&32768 != 0 {
		return o.updatedAt
	}
	return time.Time{}
}

// GetUpdatedAt returns the value of the 'updated_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object modification timestamp.
func (o *Incident) GetUpdatedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&32768 != 0
	if ok {
		value = o.updatedAt
	}
	return
}

// WorkedAt returns the value of the 'worked_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *Incident) WorkedAt() time.Time {
	if o != nil && o.bitmap_&65536 != 0 {
		return o.workedAt
	}
	return time.Time{}
}

// GetWorkedAt returns the value of the 'worked_at' attribute and
// a flag indicating if the attribute has a value.
func (o *Incident) GetWorkedAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&65536 != 0
	if ok {
		value = o.workedAt
	}
	return
}

// IncidentListKind is the name of the type used to represent list of objects of
// type 'incident'.
const IncidentListKind = "IncidentList"

// IncidentListLinkKind is the name of the type used to represent links to list
// of objects of type 'incident'.
const IncidentListLinkKind = "IncidentListLink"

// IncidentNilKind is the name of the type used to nil lists of objects of
// type 'incident'.
const IncidentListNilKind = "IncidentListNil"

// IncidentList is a list of values of the 'incident' type.
type IncidentList struct {
	href  string
	link  bool
	items []*Incident
}

// Kind returns the name of the type of the object.
func (l *IncidentList) Kind() string {
	if l == nil {
		return IncidentListNilKind
	}
	if l.link {
		return IncidentListLinkKind
	}
	return IncidentListKind
}

// Link returns true iif this is a link.
func (l *IncidentList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *IncidentList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *IncidentList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *IncidentList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *IncidentList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *IncidentList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *IncidentList) SetItems(items []*Incident) {
	l.items = items
}

// Items returns the items of the list.
func (l *IncidentList) Items() []*Incident {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *IncidentList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *IncidentList) Get(i int) *Incident {
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
func (l *IncidentList) Slice() []*Incident {
	var slice []*Incident
	if l == nil {
		slice = make([]*Incident, 0)
	} else {
		slice = make([]*Incident, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *IncidentList) Each(f func(item *Incident) bool) {
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
func (l *IncidentList) Range(f func(index int, item *Incident) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
