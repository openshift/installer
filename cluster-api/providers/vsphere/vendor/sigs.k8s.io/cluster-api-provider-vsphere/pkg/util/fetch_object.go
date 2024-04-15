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

package util

import (
	"context"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	controlplanev1 "sigs.k8s.io/cluster-api/controlplane/kubeadm/api/v1beta1"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// FetchObjectInput contains data to fetch objects.
type FetchObjectInput struct {
	ctrlclient.Client
	Object ctrlclient.Object
}

// FetchControlPlaneOwnerObject returns the ControlPlane owner for a Machine.
func FetchControlPlaneOwnerObject(ctx context.Context, input FetchObjectInput) (ctrlclient.Object, error) {
	gvk := controlplanev1.GroupVersion
	kcp := &controlplanev1.KubeadmControlPlane{}
	if err := fetchOwnerOfKindInto(ctx, input.Client, gvk, "KubeadmControlPlane", input.Object, kcp); err != nil {
		return nil, err
	}
	return kcp, nil
}

// FetchMachineDeploymentOwnerObject returns the MachineDeployment owner for a Machine.
func FetchMachineDeploymentOwnerObject(ctx context.Context, input FetchObjectInput) (ctrlclient.Object, error) {
	gvk := clusterv1.GroupVersion

	ms := &clusterv1.MachineSet{}
	if err := fetchOwnerOfKindInto(ctx, input.Client, gvk, "MachineSet", input.Object, ms); err != nil {
		return nil, err
	}

	md := &clusterv1.MachineDeployment{}
	if err := fetchOwnerOfKindInto(ctx, input.Client, gvk, "MachineDeployment", ms, md); err != nil {
		return nil, err
	}
	return md, nil
}

func fetchOwnerOfKindInto(ctx context.Context, c ctrlclient.Client, gvk schema.GroupVersion, kind string, fromObject ctrlclient.Object, intoObj ctrlclient.Object) error {
	ref, err := findOwnerRefWithKind(fromObject.GetOwnerReferences(), gvk, kind)
	if err != nil {
		return err
	}

	return c.Get(ctx, ctrlclient.ObjectKey{
		Namespace: fromObject.GetNamespace(),
		Name:      ref.Name,
	}, intoObj)
}

func findOwnerRefWithKind(ownerRefs []metav1.OwnerReference, gvk schema.GroupVersion, kind string) (*metav1.OwnerReference, error) {
	for _, ref := range ownerRefs {
		gv, err := schema.ParseGroupVersion(ref.APIVersion)
		if err != nil {
			return nil, err
		}
		if gv.Group == gvk.Group && ref.Kind == kind {
			return &ref, nil
		}
	}
	return nil, errors.Errorf("unable to find owner reference with APIVersion %s and Kind %s", gvk, kind)
}
