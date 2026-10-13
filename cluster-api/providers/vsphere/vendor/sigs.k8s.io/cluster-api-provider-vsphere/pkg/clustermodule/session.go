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
	"context"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	capvcontext "sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/identity"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/session"
)

func (s *service) fetchSessionForObject(ctx context.Context, clusterCtx *capvcontext.ClusterContext, template *infrav1.VSphereMachineTemplate) (*session.Session, error) {
	params := s.newParams(*clusterCtx)
	// Datacenter is necessary since we use the finder.
	params = params.WithDatacenter(template.Spec.Template.Spec.Datacenter)

	return s.fetchSession(ctx, clusterCtx, params)
}

func (s *service) newParams(clusterCtx capvcontext.ClusterContext) *session.Params {
	return session.NewParams().
		WithServer(clusterCtx.VSphereCluster.Spec.Server).
		WithThumbprint(clusterCtx.VSphereCluster.Spec.Thumbprint)
}

func (s *service) fetchSession(ctx context.Context, clusterCtx *capvcontext.ClusterContext, params *session.Params) (*session.Session, error) {
	if clusterCtx.VSphereCluster.Spec.IdentityRef != nil {
		creds, err := identity.GetCredentials(ctx, s.Client, clusterCtx.VSphereCluster, s.ControllerManagerContext.Namespace)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get credentials from IdentityRef")
		}

		params = params.WithUserInfo(creds.Username, creds.Password)
		return session.GetOrCreate(ctx, params)
	}

	params = params.WithUserInfo(s.ControllerManagerContext.Username, s.ControllerManagerContext.Password)
	return session.GetOrCreate(ctx, params)
}

func (s *service) fetchTemplateRef(ctx context.Context, input Wrapper) (*corev1.ObjectReference, error) {
	obj := new(unstructured.Unstructured)
	obj.SetAPIVersion(input.GetObjectKind().GroupVersionKind().GroupVersion().String())
	obj.SetKind(input.GetObjectKind().GroupVersionKind().Kind)
	obj.SetName(input.GetName())
	key := client.ObjectKey{Name: obj.GetName(), Namespace: input.GetNamespace()}
	if err := s.Client.Get(ctx, key, obj); err != nil {
		return nil, errors.Wrapf(err, "failed to get %s external object %q/%q", obj.GetKind(), key.Namespace, key.Name)
	}

	objRef := corev1.ObjectReference{}
	if err := util.UnstructuredUnmarshalField(obj, &objRef, input.GetTemplatePath()...); err != nil && err != util.ErrUnstructuredFieldNotFound {
		return nil, err
	}
	return &objRef, nil
}

func (s *service) fetchMachineTemplate(ctx context.Context, input Wrapper, templateName string) (*infrav1.VSphereMachineTemplate, error) {
	template := &infrav1.VSphereMachineTemplate{}
	if err := s.Client.Get(ctx, client.ObjectKey{
		Name:      templateName,
		Namespace: input.GetNamespace(),
	}, template); err != nil {
		return nil, err
	}
	return template, nil
}
