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
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	controlplanev1 "sigs.k8s.io/cluster-api/controlplane/kubeadm/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Wrapper is used to expose methods on the cluster-api objects responsible
// for the creation of Machine objects that need to be anti-affined.
type Wrapper interface {
	client.Object

	// GetTemplatePath is used to fetch the path that contains the infrastructure
	// machine template details in use.
	GetTemplatePath() []string

	// IsControlPlane is used to determine whether the cluster-api object is
	// responsible for control plane VMs.
	IsControlPlane() bool
}

// NewWrapper returns the correct wrapper for the passed in object.
func NewWrapper(obj client.Object) Wrapper {
	if obj.GetObjectKind().GroupVersionKind().Kind == "KubeadmControlPlane" {
		kcp, _ := obj.(*controlplanev1.KubeadmControlPlane)
		return kcpWrapper{kcp}
	}
	md, _ := obj.(*clusterv1.MachineDeployment)
	return mdWrapper{md}
}

type kcpWrapper struct {
	*controlplanev1.KubeadmControlPlane
}

func (w kcpWrapper) GetTemplatePath() []string {
	return []string{"spec", "machineTemplate", "infrastructureRef"}
}

func (w kcpWrapper) IsControlPlane() bool {
	return true
}

type mdWrapper struct {
	*clusterv1.MachineDeployment
}

func (w mdWrapper) GetTemplatePath() []string {
	return []string{"spec", "template", "spec", "infrastructureRef"}
}

func (w mdWrapper) IsControlPlane() bool {
	return false
}
