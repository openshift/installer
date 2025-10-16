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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

// QuotaRules represents the values of the 'quota_rules' type.
type QuotaRules struct {
	bitmap_          uint32
	availabilityZone string
	billingModel     string
	byoc             string
	cloud            string
	cost             int
	name             string
	product          string
	quotaId          string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *QuotaRules) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// AvailabilityZone returns the value of the 'availability_zone' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *QuotaRules) AvailabilityZone() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.availabilityZone
	}
	return ""
}

// GetAvailabilityZone returns the value of the 'availability_zone' attribute and
// a flag indicating if the attribute has a value.
func (o *QuotaRules) GetAvailabilityZone() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.availabilityZone
	}
	return
}

// BillingModel returns the value of the 'billing_model' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *QuotaRules) BillingModel() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.billingModel
	}
	return ""
}

// GetBillingModel returns the value of the 'billing_model' attribute and
// a flag indicating if the attribute has a value.
func (o *QuotaRules) GetBillingModel() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.billingModel
	}
	return
}

// Byoc returns the value of the 'byoc' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *QuotaRules) Byoc() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.byoc
	}
	return ""
}

// GetByoc returns the value of the 'byoc' attribute and
// a flag indicating if the attribute has a value.
func (o *QuotaRules) GetByoc() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.byoc
	}
	return
}

// Cloud returns the value of the 'cloud' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *QuotaRules) Cloud() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.cloud
	}
	return ""
}

// GetCloud returns the value of the 'cloud' attribute and
// a flag indicating if the attribute has a value.
func (o *QuotaRules) GetCloud() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.cloud
	}
	return
}

// Cost returns the value of the 'cost' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *QuotaRules) Cost() int {
	if o != nil && o.bitmap_&16 != 0 {
		return o.cost
	}
	return 0
}

// GetCost returns the value of the 'cost' attribute and
// a flag indicating if the attribute has a value.
func (o *QuotaRules) GetCost() (value int, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.cost
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *QuotaRules) Name() string {
	if o != nil && o.bitmap_&32 != 0 {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
func (o *QuotaRules) GetName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.name
	}
	return
}

// Product returns the value of the 'product' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *QuotaRules) Product() string {
	if o != nil && o.bitmap_&64 != 0 {
		return o.product
	}
	return ""
}

// GetProduct returns the value of the 'product' attribute and
// a flag indicating if the attribute has a value.
func (o *QuotaRules) GetProduct() (value string, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.product
	}
	return
}

// QuotaId returns the value of the 'quota_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *QuotaRules) QuotaId() string {
	if o != nil && o.bitmap_&128 != 0 {
		return o.quotaId
	}
	return ""
}

// GetQuotaId returns the value of the 'quota_id' attribute and
// a flag indicating if the attribute has a value.
func (o *QuotaRules) GetQuotaId() (value string, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.quotaId
	}
	return
}

// QuotaRulesListKind is the name of the type used to represent list of objects of
// type 'quota_rules'.
const QuotaRulesListKind = "QuotaRulesList"

// QuotaRulesListLinkKind is the name of the type used to represent links to list
// of objects of type 'quota_rules'.
const QuotaRulesListLinkKind = "QuotaRulesListLink"

// QuotaRulesNilKind is the name of the type used to nil lists of objects of
// type 'quota_rules'.
const QuotaRulesListNilKind = "QuotaRulesListNil"

// QuotaRulesList is a list of values of the 'quota_rules' type.
type QuotaRulesList struct {
	href  string
	link  bool
	items []*QuotaRules
}

// Len returns the length of the list.
func (l *QuotaRulesList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *QuotaRulesList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *QuotaRulesList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *QuotaRulesList) SetItems(items []*QuotaRules) {
	l.items = items
}

// Items returns the items of the list.
func (l *QuotaRulesList) Items() []*QuotaRules {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *QuotaRulesList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *QuotaRulesList) Get(i int) *QuotaRules {
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
func (l *QuotaRulesList) Slice() []*QuotaRules {
	var slice []*QuotaRules
	if l == nil {
		slice = make([]*QuotaRules, 0)
	} else {
		slice = make([]*QuotaRules, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *QuotaRulesList) Each(f func(item *QuotaRules) bool) {
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
func (l *QuotaRulesList) Range(f func(index int, item *QuotaRules) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
