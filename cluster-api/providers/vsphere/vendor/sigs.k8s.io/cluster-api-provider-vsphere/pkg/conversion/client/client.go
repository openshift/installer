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
	"strings"

	pkgerrors "github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/conversion"
)

// NewWithConverter return a client that can convert objects before write and after read.
func NewWithConverter(c client.Client, converter *conversion.Converter) (client.Client, error) {
	if err := checkConverterAndSchemeAreConsistent(converter, c.Scheme()); err != nil {
		return nil, err
	}

	return &conversionClient{
		internalClient: c,
		converter:      converter,
	}, nil
}

func checkConverterAndSchemeAreConsistent(converter *conversion.Converter, scheme *runtime.Scheme) error {
	allErrs := []error{}
	for gvk := range converter.AllKnownHubTypes() {
		if strings.HasSuffix(gvk.Kind, "List") {
			continue
		}

		if !scheme.Recognizes(gvk) {
			allErrs = append(allErrs, pkgerrors.Errorf("converter is configured to handle %s but the client scheme is not aware of this type", gvk))
		}

		for _, spokeGV := range scheme.PrioritizedVersionsForGroup(gvk.Group) {
			if spokeGV == gvk.GroupVersion() {
				continue
			}
			spokeGVK := spokeGV.WithKind(gvk.Kind)

			if !scheme.Recognizes(spokeGVK) {
				continue
			}
			if !converter.Recognizes(spokeGVK) {
				allErrs = append(allErrs, pkgerrors.Errorf("converter is configured to handle %s but it is not configured for handling conversions from/to %s which is registered in the client scheme", gvk, spokeGV.Version))
			}
		}
	}
	return kerrors.NewAggregate(allErrs)
}

type conversionClient struct {
	converter      *conversion.Converter
	internalClient client.Client
}

// conversionClient must implement client.Client.
var _ client.Client = &conversionClient{}

// Get retrieves an obj for the given object key from the Kubernetes Cluster.
func (c *conversionClient) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	if !c.converter.IsHub(obj) {
		return c.internalClient.Get(ctx, key, obj, opts...)
	}

	spokeObj, err := c.newSpokeObjectFor(obj)
	if err != nil {
		return err
	}

	if err := c.internalClient.Get(ctx, key, spokeObj, opts...); err != nil {
		return err
	}
	if err := c.converter.Convert(ctx, spokeObj, obj); err != nil {
		return pkgerrors.Wrapf(err, "failed to convert %s from spoke version while getting it", klog.KObj(spokeObj))
	}
	return nil
}

// List retrieves list of objects for a given namespace and list options.
func (c *conversionClient) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	if !c.converter.IsHub(list) {
		return c.internalClient.List(ctx, list, opts...)
	}

	spokeList, err := c.newSpokeObjectListFor(list)
	if err != nil {
		return err
	}

	if err := c.internalClient.List(ctx, spokeList, opts...); err != nil {
		return err
	}

	spokeListItems, err := meta.ExtractList(spokeList)
	if err != nil {
		return err
	}

	listItems := []runtime.Object{}
	for _, spokeRuntimeItem := range spokeListItems {
		spokeItem, ok := spokeRuntimeItem.(client.Object)
		if !ok {
			return pkgerrors.Errorf("%T does not implement client.Object", spokeRuntimeItem)
		}

		listItem, err := c.newObjectListItemFor(list)
		if err != nil {
			return err
		}

		if err := c.converter.Convert(ctx, spokeItem, listItem); err != nil {
			return pkgerrors.Wrapf(err, "failed to convert %s from spoke version while listing it", klog.KObj(listItem))
		}
		listItems = append(listItems, listItem)
	}

	return meta.SetList(list, listItems)
}

// Apply applies the given apply configuration to the Kubernetes cluster.
func (c *conversionClient) Apply(ctx context.Context, obj runtime.ApplyConfiguration, opts ...client.ApplyOption) error {
	if cObj, ok := obj.(client.Object); ok && c.converter.IsHub(cObj) {
		return pkgerrors.New("conversionClient only implements methods used in CAPV")
	}

	return c.internalClient.Apply(ctx, obj, opts...)
}

// Create creates the object obj in the Kubernetes cluster.
func (c *conversionClient) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	if !c.converter.IsHub(obj) {
		return c.internalClient.Create(ctx, obj, opts...)
	}

	spokeObj, err := c.newSpokeObjectFor(obj)
	if err != nil {
		return err
	}
	if err := c.converter.Convert(ctx, obj, spokeObj); err != nil {
		return pkgerrors.Wrapf(err, "failed to convert %s to spoke version while creating it", klog.KObj(obj))
	}

	if err := c.internalClient.Create(ctx, spokeObj, opts...); err != nil {
		return err
	}
	if err := c.converter.Convert(ctx, spokeObj, obj); err != nil {
		return pkgerrors.Wrapf(err, "failed to convert %s from spoke version while creating it", klog.KObj(spokeObj))
	}
	return nil
}

// Delete deletes the given obj from Kubernetes cluster.
func (c *conversionClient) Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error {
	if !c.converter.IsHub(obj) {
		return c.internalClient.Delete(ctx, obj, opts...)
	}

	spokeObj, err := c.newSpokeObjectFor(obj)
	if err != nil {
		return err
	}
	if err := c.converter.Convert(ctx, obj, spokeObj); err != nil {
		return pkgerrors.Wrapf(err, "failed to convert %s to spoke version while deleting it", klog.KObj(obj))
	}

	if err := c.internalClient.Delete(ctx, spokeObj, opts...); err != nil {
		return err
	}

	if err := c.converter.Convert(ctx, spokeObj, obj); err != nil {
		return pkgerrors.Wrapf(err, "failed to convert %s from spoke version while deleting it", klog.KObj(spokeObj))
	}
	return nil
}

// Update updates the given obj in the Kubernetes cluster.
func (c *conversionClient) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	if !c.converter.IsHub(obj) {
		return c.internalClient.Update(ctx, obj, opts...)
	}
	return pkgerrors.New("update must not be used for hub types. Use patch instead")
}

// Patch patches the given obj in the Kubernetes cluster.
func (c *conversionClient) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
	if !c.converter.IsHub(obj) {
		return c.internalClient.Patch(ctx, obj, patch, opts...)
	}

	if _, ok := patch.(*conversionMergePatch); !ok {
		return pkgerrors.Errorf("patch must be created using sigs.k8s.io/cluster-api-provider-vsphere/pkg/conversion/client.MergeFrom or MergeFromWithOptions")
	}

	spokeObj, err := c.newSpokeObjectFor(obj)
	if err != nil {
		return err
	}
	if err := c.converter.Convert(ctx, obj, spokeObj); err != nil {
		return pkgerrors.Wrapf(err, "failed to convert %s to spoke version while patching it", klog.KObj(obj))
	}

	if err := c.internalClient.Patch(ctx, spokeObj, patch, opts...); err != nil {
		return err
	}
	if err := c.converter.Convert(ctx, spokeObj, obj); err != nil {
		return pkgerrors.Wrapf(err, "failed to convert %s from spoke version while patching it", klog.KObj(spokeObj))
	}
	return nil
}

// DeleteAllOf deletes all objects of the given type matching the given options.
func (c *conversionClient) DeleteAllOf(ctx context.Context, obj client.Object, opts ...client.DeleteAllOfOption) error {
	if !c.converter.IsHub(obj) {
		return c.internalClient.DeleteAllOf(ctx, obj, opts...)
	}

	return pkgerrors.New("conversionClient only implements methods used in CAPV")
}

func (c *conversionClient) GroupVersionKindFor(obj runtime.Object) (schema.GroupVersionKind, error) {
	return c.internalClient.GroupVersionKindFor(obj)
}

func (c *conversionClient) IsObjectNamespaced(obj runtime.Object) (bool, error) {
	return c.internalClient.IsObjectNamespaced(obj)
}

func (c *conversionClient) Scheme() *runtime.Scheme {
	return c.internalClient.Scheme()
}

func (c *conversionClient) RESTMapper() meta.RESTMapper {
	return c.internalClient.RESTMapper()
}

func (c *conversionClient) Status() client.SubResourceWriter {
	return c.SubResource("status")
}

func (c *conversionClient) SubResource(subResource string) client.SubResourceClient {
	return &conversionSubResourceClient{conversionClient: c, subResource: subResource}
}

type conversionSubResourceClient struct {
	*conversionClient
	subResource string
}

// conversionClient must implement client.Client.
var _ client.SubResourceClient = &conversionSubResourceClient{}

func (c conversionSubResourceClient) Get(_ context.Context, _ client.Object, _ client.Object, _ ...client.SubResourceGetOption) error {
	return pkgerrors.New("conversionClient only implements methods used in CAPV")
}

func (c conversionSubResourceClient) Create(_ context.Context, _ client.Object, _ client.Object, _ ...client.SubResourceCreateOption) error {
	return pkgerrors.New("conversionClient only implements methods used in CAPV")
}

func (c conversionSubResourceClient) Update(_ context.Context, _ client.Object, _ ...client.SubResourceUpdateOption) error {
	return pkgerrors.New("update must not be used for hub types. Use patch instead")
}

func (c conversionSubResourceClient) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.SubResourcePatchOption) error {
	if !c.converter.IsHub(obj) {
		return c.internalClient.Status().Patch(ctx, obj, patch, opts...)
	}

	spokeObj, err := c.newSpokeObjectFor(obj)
	if err != nil {
		return err
	}
	if err := c.converter.Convert(ctx, obj, spokeObj); err != nil {
		return pkgerrors.Wrapf(err, "failed to convert %s to spoke version while patching status on it", klog.KObj(obj))
	}

	if err := c.internalClient.Status().Patch(ctx, spokeObj, patch, opts...); err != nil {
		return err
	}

	if err := c.converter.Convert(ctx, spokeObj, obj); err != nil {
		return pkgerrors.Wrapf(err, "failed to convert %s from spoke version while patching status on it", klog.KObj(spokeObj))
	}
	return nil
}

func (c conversionSubResourceClient) Apply(ctx context.Context, obj runtime.ApplyConfiguration, opts ...client.SubResourceApplyOption) error {
	if cObj, ok := obj.(client.Object); ok && c.converter.IsHub(cObj) {
		return pkgerrors.New("conversionClient only implements methods used in CAPV")
	}

	return c.internalClient.Status().Apply(ctx, obj, opts...)
}

func (c *conversionClient) newSpokeObjectFor(obj client.Object) (client.Object, error) {
	gvk, err := c.converter.SpokeGroupVersionKindFor(obj)
	if err != nil {
		return nil, err
	}

	o, err := c.internalClient.Scheme().New(gvk)
	if err != nil {
		return nil, err
	}

	spokeObj, ok := o.(client.Object)
	if !ok {
		return nil, pkgerrors.Errorf("object for %s does not implement sigs.k8s.io/controller-runtime/pkg/client.Object", gvk)
	}
	return spokeObj, nil
}

func (c *conversionClient) newSpokeObjectListFor(list client.ObjectList) (client.ObjectList, error) {
	gvk, err := c.converter.SpokeGroupVersionKindFor(list)
	if err != nil {
		return nil, err
	}

	o, err := c.internalClient.Scheme().New(gvk)
	if err != nil {
		return nil, err
	}

	spokeList, ok := o.(client.ObjectList)
	if !ok {
		return nil, pkgerrors.Errorf("object for %s does not implement sigs.k8s.io/controller-runtime/pkg/client.ObjectList", gvk)
	}
	return spokeList, nil
}

func (c *conversionClient) newObjectListItemFor(list client.ObjectList) (client.Object, error) {
	gvkList, err := c.internalClient.GroupVersionKindFor(list)
	if err != nil {
		return nil, err
	}

	if !strings.HasSuffix(gvkList.Kind, "List") {
		return nil, pkgerrors.Errorf("object %s does not have a kind with the List suffix", gvkList)
	}

	gvk := schema.GroupVersionKind{
		Group:   gvkList.Group,
		Version: gvkList.Version,
		Kind:    strings.TrimSuffix(gvkList.Kind, "List"),
	}

	o, err := c.internalClient.Scheme().New(gvk)
	if err != nil {
		return nil, err
	}

	obj, ok := o.(client.Object)
	if !ok {
		return nil, pkgerrors.Errorf("object for %s does not implement sigs.k8s.io/controller-runtime/pkg/client.Object", gvk)
	}
	return obj, nil
}
