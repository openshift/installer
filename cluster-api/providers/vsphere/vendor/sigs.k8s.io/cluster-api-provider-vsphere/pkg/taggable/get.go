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

// Package taggable handles tagging objects in VSphere.
package taggable

import (
	"context"

	"github.com/pkg/errors"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services/govmomi/find"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/session"
)

type taggableContext interface {
	GetSession() *session.Session
}

// GetObjects returns the objects for a given failure domain.
func GetObjects(ctx context.Context, taggableCtx taggableContext, failureDomain *infrav1.VSphereFailureDomain, fdType infrav1.FailureDomainType) (Objects, error) {
	finderFunc := find.ObjectFunc(fdType, failureDomain.Spec.Topology, taggableCtx.GetSession().Finder)
	objRefs, err := finderFunc(ctx)
	if err != nil {
		return nil, err
	}
	if len(objRefs) == 0 {
		return nil, errors.New("unable to find taggable object")
	}

	objects := make(Objects, len(objRefs))
	for i, ref := range objRefs {
		objects[i] = managedObject{
			tagManager: taggableCtx.GetSession().TagManager,
			ref:        ref,
		}
	}
	return objects, nil
}
