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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/webrca/v1

import (
	time "time"
)

// ProductKind is the name of the type used to represent objects
// of type 'product'.
const ProductKind = "Product"

// ProductLinkKind is the name of the type used to represent links
// to objects of type 'product'.
const ProductLinkKind = "ProductLink"

// ProductNilKind is the name of the type used to nil references
// to objects of type 'product'.
const ProductNilKind = "ProductNil"

// Product represents the values of the 'product' type.
//
// Definition of a Web RCA product.
type Product struct {
	fieldSet_   []bool
	id          string
	href        string
	createdAt   time.Time
	deletedAt   time.Time
	productId   string
	productName string
	updatedAt   time.Time
}

// Kind returns the name of the type of the object.
func (o *Product) Kind() string {
	if o == nil {
		return ProductNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return ProductLinkKind
	}
	return ProductKind
}

// Link returns true if this is a link.
func (o *Product) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *Product) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *Product) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *Product) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *Product) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *Product) Empty() bool {
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

// CreatedAt returns the value of the 'created_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Object creation timestamp.
func (o *Product) CreatedAt() time.Time {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.createdAt
	}
	return time.Time{}
}

// GetCreatedAt returns the value of the 'created_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object creation timestamp.
func (o *Product) GetCreatedAt() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.createdAt
	}
	return
}

// DeletedAt returns the value of the 'deleted_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Object deletion timestamp.
func (o *Product) DeletedAt() time.Time {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.deletedAt
	}
	return time.Time{}
}

// GetDeletedAt returns the value of the 'deleted_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object deletion timestamp.
func (o *Product) GetDeletedAt() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.deletedAt
	}
	return
}

// ProductId returns the value of the 'product_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The product ID from status board
func (o *Product) ProductId() string {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.productId
	}
	return ""
}

// GetProductId returns the value of the 'product_id' attribute and
// a flag indicating if the attribute has a value.
//
// The product ID from status board
func (o *Product) GetProductId() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.productId
	}
	return
}

// ProductName returns the value of the 'product_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The name of the product from status-board.
func (o *Product) ProductName() string {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.productName
	}
	return ""
}

// GetProductName returns the value of the 'product_name' attribute and
// a flag indicating if the attribute has a value.
//
// The name of the product from status-board.
func (o *Product) GetProductName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.productName
	}
	return
}

// UpdatedAt returns the value of the 'updated_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Object modification timestamp.
func (o *Product) UpdatedAt() time.Time {
	if o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7] {
		return o.updatedAt
	}
	return time.Time{}
}

// GetUpdatedAt returns the value of the 'updated_at' attribute and
// a flag indicating if the attribute has a value.
//
// Object modification timestamp.
func (o *Product) GetUpdatedAt() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7]
	if ok {
		value = o.updatedAt
	}
	return
}

// ProductListKind is the name of the type used to represent list of objects of
// type 'product'.
const ProductListKind = "ProductList"

// ProductListLinkKind is the name of the type used to represent links to list
// of objects of type 'product'.
const ProductListLinkKind = "ProductListLink"

// ProductNilKind is the name of the type used to nil lists of objects of
// type 'product'.
const ProductListNilKind = "ProductListNil"

// ProductList is a list of values of the 'product' type.
type ProductList struct {
	href  string
	link  bool
	items []*Product
}

// Kind returns the name of the type of the object.
func (l *ProductList) Kind() string {
	if l == nil {
		return ProductListNilKind
	}
	if l.link {
		return ProductListLinkKind
	}
	return ProductListKind
}

// Link returns true iif this is a link.
func (l *ProductList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *ProductList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *ProductList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *ProductList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ProductList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ProductList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ProductList) SetItems(items []*Product) {
	l.items = items
}

// Items returns the items of the list.
func (l *ProductList) Items() []*Product {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ProductList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ProductList) Get(i int) *Product {
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
func (l *ProductList) Slice() []*Product {
	var slice []*Product
	if l == nil {
		slice = make([]*Product, 0)
	} else {
		slice = make([]*Product, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ProductList) Each(f func(item *Product) bool) {
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
func (l *ProductList) Range(f func(index int, item *Product) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
