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

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

// ReleaseImagesListBuilder contains the data and logic needed to build
// 'release_images' objects.
type ReleaseImagesListBuilder struct {
	items []*ReleaseImagesBuilder
}

// NewReleaseImagesList creates a new builder of 'release_images' objects.
func NewReleaseImagesList() *ReleaseImagesListBuilder {
	return new(ReleaseImagesListBuilder)
}

// Items sets the items of the list.
func (b *ReleaseImagesListBuilder) Items(values ...*ReleaseImagesBuilder) *ReleaseImagesListBuilder {
	b.items = make([]*ReleaseImagesBuilder, len(values))
	copy(b.items, values)
	return b
}

// Empty returns true if the list is empty.
func (b *ReleaseImagesListBuilder) Empty() bool {
	return b == nil || len(b.items) == 0
}

// Copy copies the items of the given list into this builder, discarding any previous items.
func (b *ReleaseImagesListBuilder) Copy(list *ReleaseImagesList) *ReleaseImagesListBuilder {
	if list == nil || list.items == nil {
		b.items = nil
	} else {
		b.items = make([]*ReleaseImagesBuilder, len(list.items))
		for i, v := range list.items {
			b.items[i] = NewReleaseImages().Copy(v)
		}
	}
	return b
}

// Build creates a list of 'release_images' objects using the
// configuration stored in the builder.
func (b *ReleaseImagesListBuilder) Build() (list *ReleaseImagesList, err error) {
	items := make([]*ReleaseImages, len(b.items))
	for i, item := range b.items {
		items[i], err = item.Build()
		if err != nil {
			return
		}
	}
	list = new(ReleaseImagesList)
	list.items = items
	return
}
