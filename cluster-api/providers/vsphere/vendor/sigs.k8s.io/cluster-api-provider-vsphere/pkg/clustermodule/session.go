/*
Copyright 2022 The Kubernetes Authors.

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

package clustermodule

import (
	goctx "context"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/identity"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/session"
)

func fetchSessionForObject(ctx *context.ClusterContext, template *infrav1.VSphereMachineTemplate) (*session.Session, error) {
	params := newParams(*ctx)
	// Datacenter is necessary since we use the finder.
	params = params.WithDatacenter(template.Spec.Template.Spec.Datacenter)

	return fetchSession(ctx, params)
}

func newParams(ctx context.ClusterContext) *session.Params {
	return session.NewParams().
		WithServer(ctx.VSphereCluster.Spec.Server).
		WithThumbprint(ctx.VSphereCluster.Spec.Thumbprint).
		WithFeatures(session.Feature{
			EnableKeepAlive:   ctx.EnableKeepAlive,
			KeepAliveDuration: ctx.KeepAliveDuration,
		})
}

func fetchSession(ctx *context.ClusterContext, params *session.Params) (*session.Session, error) {
	if ctx.VSphereCluster.Spec.IdentityRef != nil {
		creds, err := identity.GetCredentials(ctx, ctx.Client, ctx.VSphereCluster, ctx.Namespace)
		if err != nil {
			return nil, err
		}

		params = params.WithUserInfo(creds.Username, creds.Password)
		return session.GetOrCreate(ctx, params)
	}

	params = params.WithUserInfo(ctx.Username, ctx.Password)
	return session.GetOrCreate(ctx, params)
}

func fetchTemplateRef(ctx goctx.Context, c client.Client, input Wrapper) (*corev1.ObjectReference, error) {
	obj := new(unstructured.Unstructured)
	obj.SetAPIVersion(input.GetObjectKind().GroupVersionKind().GroupVersion().String())
	obj.SetKind(input.GetObjectKind().GroupVersionKind().Kind)
	obj.SetName(input.GetName())
	key := client.ObjectKey{Name: obj.GetName(), Namespace: input.GetNamespace()}
	if err := c.Get(ctx, key, obj); err != nil {
		return nil, errors.Wrapf(err, "failed to retrieve %s external object %q/%q", obj.GetKind(), key.Namespace, key.Name)
	}

	objRef := corev1.ObjectReference{}
	if err := util.UnstructuredUnmarshalField(obj, &objRef, input.GetTemplatePath()...); err != nil && err != util.ErrUnstructuredFieldNotFound {
		return nil, err
	}
	return &objRef, nil
}

func fetchMachineTemplate(ctx *context.ClusterContext, input Wrapper, templateName string) (*infrav1.VSphereMachineTemplate, error) {
	template := &infrav1.VSphereMachineTemplate{}
	if err := ctx.Client.Get(ctx, client.ObjectKey{
		Name:      templateName,
		Namespace: input.GetNamespace(),
	}, template); err != nil {
		return nil, err
	}
	return template, nil
}
