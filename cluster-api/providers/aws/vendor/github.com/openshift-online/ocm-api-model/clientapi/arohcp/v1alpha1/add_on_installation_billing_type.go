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

package v1alpha1 // github.com/openshift-online/ocm-api-model/clientapi/arohcp/v1alpha1

// AddOnInstallationBillingKind is the name of the type used to represent objects
// of type 'add_on_installation_billing'.
const AddOnInstallationBillingKind = "AddOnInstallationBilling"

// AddOnInstallationBillingLinkKind is the name of the type used to represent links
// to objects of type 'add_on_installation_billing'.
const AddOnInstallationBillingLinkKind = "AddOnInstallationBillingLink"

// AddOnInstallationBillingNilKind is the name of the type used to nil references
// to objects of type 'add_on_installation_billing'.
const AddOnInstallationBillingNilKind = "AddOnInstallationBillingNil"

// AddOnInstallationBilling represents the values of the 'add_on_installation_billing' type.
//
// Representation of an add-on installation billing.
type AddOnInstallationBilling struct {
	fieldSet_                 []bool
	id                        string
	href                      string
	billingMarketplaceAccount string
	billingModel              BillingModel
}

// Kind returns the name of the type of the object.
func (o *AddOnInstallationBilling) Kind() string {
	if o == nil {
		return AddOnInstallationBillingNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return AddOnInstallationBillingLinkKind
	}
	return AddOnInstallationBillingKind
}

// Link returns true if this is a link.
func (o *AddOnInstallationBilling) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *AddOnInstallationBilling) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *AddOnInstallationBilling) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *AddOnInstallationBilling) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *AddOnInstallationBilling) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AddOnInstallationBilling) Empty() bool {
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

// BillingMarketplaceAccount returns the value of the 'billing_marketplace_account' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Account ID for billing market place
func (o *AddOnInstallationBilling) BillingMarketplaceAccount() string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.billingMarketplaceAccount
	}
	return ""
}

// GetBillingMarketplaceAccount returns the value of the 'billing_marketplace_account' attribute and
// a flag indicating if the attribute has a value.
//
// Account ID for billing market place
func (o *AddOnInstallationBilling) GetBillingMarketplaceAccount() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.billingMarketplaceAccount
	}
	return
}

// BillingModel returns the value of the 'billing_model' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Billing Model for addon resources
func (o *AddOnInstallationBilling) BillingModel() BillingModel {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.billingModel
	}
	return BillingModel("")
}

// GetBillingModel returns the value of the 'billing_model' attribute and
// a flag indicating if the attribute has a value.
//
// Billing Model for addon resources
func (o *AddOnInstallationBilling) GetBillingModel() (value BillingModel, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.billingModel
	}
	return
}

// AddOnInstallationBillingListKind is the name of the type used to represent list of objects of
// type 'add_on_installation_billing'.
const AddOnInstallationBillingListKind = "AddOnInstallationBillingList"

// AddOnInstallationBillingListLinkKind is the name of the type used to represent links to list
// of objects of type 'add_on_installation_billing'.
const AddOnInstallationBillingListLinkKind = "AddOnInstallationBillingListLink"

// AddOnInstallationBillingNilKind is the name of the type used to nil lists of objects of
// type 'add_on_installation_billing'.
const AddOnInstallationBillingListNilKind = "AddOnInstallationBillingListNil"

// AddOnInstallationBillingList is a list of values of the 'add_on_installation_billing' type.
type AddOnInstallationBillingList struct {
	href  string
	link  bool
	items []*AddOnInstallationBilling
}

// Kind returns the name of the type of the object.
func (l *AddOnInstallationBillingList) Kind() string {
	if l == nil {
		return AddOnInstallationBillingListNilKind
	}
	if l.link {
		return AddOnInstallationBillingListLinkKind
	}
	return AddOnInstallationBillingListKind
}

// Link returns true iif this is a link.
func (l *AddOnInstallationBillingList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *AddOnInstallationBillingList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *AddOnInstallationBillingList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *AddOnInstallationBillingList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AddOnInstallationBillingList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AddOnInstallationBillingList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AddOnInstallationBillingList) SetItems(items []*AddOnInstallationBilling) {
	l.items = items
}

// Items returns the items of the list.
func (l *AddOnInstallationBillingList) Items() []*AddOnInstallationBilling {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AddOnInstallationBillingList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AddOnInstallationBillingList) Get(i int) *AddOnInstallationBilling {
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
func (l *AddOnInstallationBillingList) Slice() []*AddOnInstallationBilling {
	var slice []*AddOnInstallationBilling
	if l == nil {
		slice = make([]*AddOnInstallationBilling, 0)
	} else {
		slice = make([]*AddOnInstallationBilling, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AddOnInstallationBillingList) Each(f func(item *AddOnInstallationBilling) bool) {
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
func (l *AddOnInstallationBillingList) Range(f func(index int, item *AddOnInstallationBilling) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
