/*
Copyright 2021 The Kubernetes Authors.

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

package taggable

import (
	"context"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vapi/tags"
)

type Object interface {
	HasTag(ctx context.Context, tagName string) (bool, error)

	AttachTag(ctx context.Context, tagName string) error
}

type Objects []Object

type managedObject struct {
	tagManager *tags.Manager

	ref object.Reference
}

func (m managedObject) HasTag(ctx context.Context, tagName string) (bool, error) {
	attachedTags, err := m.tagManager.GetAttachedTags(ctx, m.ref)
	if err != nil {
		return false, err
	}
	for _, tag := range attachedTags {
		if tag.Name == tagName {
			return true, nil
		}
	}
	return false, nil
}

func (m managedObject) AttachTag(ctx context.Context, tagName string) error {
	return m.tagManager.AttachTag(ctx, tagName, m.ref)
}

func (m managedObject) String() string {
	return m.ref.Reference().String()
}
