/*
Copyright 2026 The Kubernetes Authors.

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

package client

import (
	"context"

	pkgerrors "github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// MergeFromWithOptions creates a Patch that patches using the merge-patch strategy with the given object as base.
// When required, the generated patch performs conversion for both/one of the original or the target object.
// Note: obj must be an hub type.
func MergeFromWithOptions(ctx context.Context, c client.Client, obj client.Object, opts ...client.MergeFromOption) (client.Patch, error) {
	cc, ok := c.(*conversionClient)
	if !ok {
		return nil, pkgerrors.Errorf("client must be created using sigs.k8s.io/cluster-api-provider-vsphere/pkg/conversion/client.NewWithConverter")
	}

	_, err := cc.converter.SpokeGroupVersionKindFor(obj)
	if err != nil {
		return nil, err
	}

	options := &client.MergeFromOptions{}
	for _, opt := range opts {
		opt.ApplyToMergeFrom(options)
	}

	return &conversionMergePatch{
		conversionCtx: ctx,
		from:          obj,
		client:        cc,
		options:       opts,
	}, nil
}

// MergeFrom creates a Patch that patches using the merge-patch strategy with the given object as base.
// When required, the generated patch performs conversion for both/one of the original or the target object.
// Note: obj must be an hub type.
func MergeFrom(ctx context.Context, c client.Client, obj client.Object) (client.Patch, error) {
	return MergeFromWithOptions(ctx, c, obj)
}

type conversionMergePatch struct {
	conversionCtx context.Context //nolint:containedctx
	client        *conversionClient

	from    client.Object
	options []client.MergeFromOption
}

// conversionMergePatch must implement client.Patch.
var _ client.Patch = &conversionMergePatch{}

// Type is the PatchType of the patch.
func (p *conversionMergePatch) Type() types.PatchType {
	return types.MergePatchType
}

// Data is the raw data representing the patch.
// Note: obj can be either an object to be converted or an object already converted.
func (p *conversionMergePatch) Data(obj client.Object) ([]byte, error) {
	fromObj, err := p.client.newSpokeObjectFor(p.from)
	if err != nil {
		return nil, err
	}
	if err := p.client.converter.Convert(p.conversionCtx, p.from, fromObj); err != nil {
		return nil, pkgerrors.Wrapf(err, "failed to convert original %s to target version while computing patch data", klog.KObj(p.from))
	}

	toObj := obj
	if p.client.converter.IsHub(obj) {
		toObj, err = p.client.newSpokeObjectFor(obj)
		if err != nil {
			return nil, err
		}
		if err := p.client.converter.Convert(p.conversionCtx, obj, toObj); err != nil {
			return nil, pkgerrors.Wrapf(err, "failed to convert modified %s to target version while computing patch data", klog.KObj(obj))
		}
	}

	gvkFrom, err := p.client.GroupVersionKindFor(fromObj)
	if err != nil {
		return nil, err
	}
	gvkTo, err := p.client.GroupVersionKindFor(toObj)
	if err != nil {
		return nil, err
	}
	if gvkFrom != gvkTo {
		return nil, pkgerrors.Errorf("cannot generate patch data between %s and %s", gvkFrom, gvkTo)
	}

	return client.MergeFromWithOptions(fromObj, p.options...).Data(toObj)
}
