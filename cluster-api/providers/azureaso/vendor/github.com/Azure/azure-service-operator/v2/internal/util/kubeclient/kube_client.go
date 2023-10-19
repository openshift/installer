/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package kubeclient

import (
	"context"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Client interface {
	client.Client

	// Additional helpers

	GetObject(ctx context.Context, namespacedName types.NamespacedName, gvk schema.GroupVersionKind) (client.Object, error)
	GetObjectOrDefault(ctx context.Context, namespacedName types.NamespacedName, gvk schema.GroupVersionKind) (client.Object, error)
	CommitObject(ctx context.Context, obj client.Object) error
}

type clientHelper struct {
	client client.Client
}

var _ Client = &clientHelper{}

func NewClient(client client.Client) Client {
	return &clientHelper{
		client: client,
	}
}

func (c *clientHelper) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	return c.client.Get(ctx, key, obj, opts...)
}

func (c *clientHelper) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	return c.client.List(ctx, list, opts...)
}

func (c *clientHelper) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	return c.client.Create(ctx, obj, opts...)
}

func (c *clientHelper) Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error {
	return c.client.Delete(ctx, obj, opts...)
}

func (c *clientHelper) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	return c.client.Update(ctx, obj, opts...)
}

func (c *clientHelper) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
	return c.client.Patch(ctx, obj, patch, opts...)
}

func (c *clientHelper) DeleteAllOf(ctx context.Context, obj client.Object, opts ...client.DeleteAllOfOption) error {
	return c.client.DeleteAllOf(ctx, obj, opts...)
}

func (c *clientHelper) Status() client.StatusWriter {
	return c.client.Status()
}

func (c *clientHelper) Scheme() *runtime.Scheme {
	return c.client.Scheme()
}

func (c *clientHelper) RESTMapper() meta.RESTMapper {
	return c.client.RESTMapper()
}

func (c *clientHelper) SubResource(subResource string) client.SubResourceClient {
	return c.client.SubResource(subResource)
}

func (c *clientHelper) GroupVersionKindFor(obj runtime.Object) (schema.GroupVersionKind, error) {
	return c.client.GroupVersionKindFor(obj)
}

func (c *clientHelper) IsObjectNamespaced(obj runtime.Object) (bool, error) {
	return c.client.IsObjectNamespaced(obj)
}

func (c *clientHelper) GetObject(ctx context.Context, namespacedName types.NamespacedName, gvk schema.GroupVersionKind) (client.Object, error) {
	obj, err := c.Scheme().New(gvk)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create object from gvk %s with", gvk)
	}

	clientObj, ok := obj.(client.Object)
	if !ok {
		return nil, errors.Errorf("gvk %s doesn't implement client.Object", gvk)
	}

	if err := c.Get(ctx, namespacedName, clientObj); err != nil {
		return nil, err
	}

	// Ensure GVK is populated
	clientObj.GetObjectKind().SetGroupVersionKind(gvk)

	return clientObj, nil
}

func (c *clientHelper) GetObjectOrDefault(ctx context.Context, namespacedName types.NamespacedName, gvk schema.GroupVersionKind) (client.Object, error) {
	result, err := c.GetObject(ctx, namespacedName, gvk)
	if apierrors.IsNotFound(err) {
		return nil, nil
	}

	return result, err
}

// CommitObject persists the contents of obj to etcd by using the Kubernetes client.
// Note that after this method has been called, obj contains the result of the update
// from APIServer (including an updated resourceVersion). Both Spec and Status are written
func (c *clientHelper) CommitObject(ctx context.Context, obj client.Object) error {
	// Order of updates (spec first or status first) matters here.
	// If the status is updated first: clients that are waiting on status
	// Condition Ready == true might see that quickly enough, and make a spec
	// update fast enough, to conflict with the second write (that of the spec).
	// This will trigger extra requests to Azure and fail our recording tests but is
	// otherwise harmless in an actual deployment.
	// We update the spec first to avoid the above problem.

	// We must clone here because the result of this update could contain
	// fields such as status.location that may not be set but are not omitempty.
	// This will cause the contents we have in Status.Location to be overwritten.
	clone := obj.DeepCopyObject().(client.Object)

	err := c.Update(ctx, clone)
	if err != nil {
		return errors.Wrapf(err, "updating %s/%s resource", obj.GetNamespace(), obj.GetName())
	}

	obj.SetResourceVersion(clone.GetResourceVersion())

	// Note that subsequent calls to GET can (if using a cached client) can miss the updates we've just done.
	// See: https://github.com/kubernetes-sigs/controller-runtime/issues/1464.
	err = c.Status().Update(ctx, obj)
	if err != nil {
		return errors.Wrapf(err, "updating %s/%s resource status", obj.GetNamespace(), obj.GetName())
	}

	return nil
}
