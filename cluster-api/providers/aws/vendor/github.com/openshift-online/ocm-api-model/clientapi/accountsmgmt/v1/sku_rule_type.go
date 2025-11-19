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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/accountsmgmt/v1

// SkuRuleKind is the name of the type used to represent objects
// of type 'sku_rule'.
const SkuRuleKind = "SkuRule"

// SkuRuleLinkKind is the name of the type used to represent links
// to objects of type 'sku_rule'.
const SkuRuleLinkKind = "SkuRuleLink"

// SkuRuleNilKind is the name of the type used to nil references
// to objects of type 'sku_rule'.
const SkuRuleNilKind = "SkuRuleNil"

// SkuRule represents the values of the 'sku_rule' type.
//
// Identifies sku rule
type SkuRule struct {
	fieldSet_ []bool
	id        string
	href      string
	allowed   int
	quotaId   string
	sku       string
}

// Kind returns the name of the type of the object.
func (o *SkuRule) Kind() string {
	if o == nil {
		return SkuRuleNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return SkuRuleLinkKind
	}
	return SkuRuleKind
}

// Link returns true if this is a link.
func (o *SkuRule) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *SkuRule) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *SkuRule) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *SkuRule) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *SkuRule) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *SkuRule) Empty() bool {
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

// Allowed returns the value of the 'allowed' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Specifies the allowed parameter for calculation
func (o *SkuRule) Allowed() int {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.allowed
	}
	return 0
}

// GetAllowed returns the value of the 'allowed' attribute and
// a flag indicating if the attribute has a value.
//
// Specifies the allowed parameter for calculation
func (o *SkuRule) GetAllowed() (value int, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.allowed
	}
	return
}

// QuotaId returns the value of the 'quota_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Specifies the quota id
func (o *SkuRule) QuotaId() string {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.quotaId
	}
	return ""
}

// GetQuotaId returns the value of the 'quota_id' attribute and
// a flag indicating if the attribute has a value.
//
// Specifies the quota id
func (o *SkuRule) GetQuotaId() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.quotaId
	}
	return
}

// Sku returns the value of the 'sku' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Specifies the sku, such as ""MW00504""
func (o *SkuRule) Sku() string {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.sku
	}
	return ""
}

// GetSku returns the value of the 'sku' attribute and
// a flag indicating if the attribute has a value.
//
// Specifies the sku, such as ""MW00504""
func (o *SkuRule) GetSku() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.sku
	}
	return
}

// SkuRuleListKind is the name of the type used to represent list of objects of
// type 'sku_rule'.
const SkuRuleListKind = "SkuRuleList"

// SkuRuleListLinkKind is the name of the type used to represent links to list
// of objects of type 'sku_rule'.
const SkuRuleListLinkKind = "SkuRuleListLink"

// SkuRuleNilKind is the name of the type used to nil lists of objects of
// type 'sku_rule'.
const SkuRuleListNilKind = "SkuRuleListNil"

// SkuRuleList is a list of values of the 'sku_rule' type.
type SkuRuleList struct {
	href  string
	link  bool
	items []*SkuRule
}

// Kind returns the name of the type of the object.
func (l *SkuRuleList) Kind() string {
	if l == nil {
		return SkuRuleListNilKind
	}
	if l.link {
		return SkuRuleListLinkKind
	}
	return SkuRuleListKind
}

// Link returns true iif this is a link.
func (l *SkuRuleList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *SkuRuleList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *SkuRuleList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *SkuRuleList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *SkuRuleList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *SkuRuleList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *SkuRuleList) SetItems(items []*SkuRule) {
	l.items = items
}

// Items returns the items of the list.
func (l *SkuRuleList) Items() []*SkuRule {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *SkuRuleList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *SkuRuleList) Get(i int) *SkuRule {
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
func (l *SkuRuleList) Slice() []*SkuRule {
	var slice []*SkuRule
	if l == nil {
		slice = make([]*SkuRule, 0)
	} else {
		slice = make([]*SkuRule, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *SkuRuleList) Each(f func(item *SkuRule) bool) {
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
func (l *SkuRuleList) Range(f func(index int, item *SkuRule) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
