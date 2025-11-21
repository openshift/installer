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

package context

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
	v1beta2conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions/v1beta2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// MachineContext is the context used in VSphereMachine reconciliation.
type MachineContext interface {
	String() string
	Patch(ctx context.Context) error
	GetVSphereMachine() VSphereMachine
	GetReady() bool
	GetObjectMeta() metav1.ObjectMeta
	GetCluster() *clusterv1.Cluster
	GetMachine() *clusterv1.Machine
	SetBaseMachineContext(base *BaseMachineContext)
}

// VSphereMachine is a common interface used for VSphereMachines across VMOperator and non-VMOperator modes.
type VSphereMachine interface {
	client.Object
	v1beta1conditions.Setter
	v1beta2conditions.Setter
}
