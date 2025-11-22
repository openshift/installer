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

package cloud

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud"
	corev1 "k8s.io/api/core/v1"
	infrav1 "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
)

// Cloud alias for cloud.Cloud interface.
type Cloud = cloud.Cloud

// Reconciler is a generic interface used by components offering a type of service.
type Reconciler interface {
	Reconcile(ctx context.Context) error
	Delete(ctx context.Context) error
}

// ReconcilerWithResult is a generic interface used by components offering a type of service.
type ReconcilerWithResult interface {
	Reconcile(ctx context.Context) (ctrl.Result, error)
	Delete(ctx context.Context) (ctrl.Result, error)
}

// Client is an interface which can get cloud client.
type Client interface {
	Cloud() Cloud
	NetworkCloud() Cloud
}

// ClusterGetter is an interface which can get cluster information.
type ClusterGetter interface {
	Client
	Project() string
	Region() string
	Name() string
	Namespace() string
	NetworkName() string
	NetworkProject() string
	IsSharedVpc() bool
	SkipFirewallRuleCreation() bool
	Network() *infrav1.Network
	AdditionalLabels() infrav1.Labels
	FailureDomains() clusterv1.FailureDomains
	ControlPlaneEndpoint() clusterv1.APIEndpoint
	ResourceManagerTags() infrav1.ResourceManagerTags
	LoadBalancer() infrav1.LoadBalancerSpec
}

// ClusterSetter is an interface which can set cluster information.
type ClusterSetter interface {
	SetControlPlaneEndpoint(endpoint clusterv1.APIEndpoint)
}

// Cluster is an interface which can get and set cluster information.
type Cluster interface {
	ClusterGetter
	ClusterSetter
}

// MachineGetter is an interface which can get machine information.
type MachineGetter interface {
	Client
	Name() string
	Namespace() string
	Zone() string
	Project() string
	Role() string
	IsControlPlane() bool
	ControlPlaneGroupName() string
	GetInstanceID() *string
	GetProviderID() string
	GetBootstrapData() (string, error)
	GetInstanceStatus() *infrav1.InstanceStatus
}

// MachineSetter is an interface which can set machine information.
type MachineSetter interface {
	SetProviderID()
	SetInstanceStatus(v infrav1.InstanceStatus)
	SetFailureMessage(v error)
	SetFailureReason(v string)
	SetAnnotation(key, value string)
	SetAddresses(addressList []corev1.NodeAddress)
}

// Machine is an interface which can get and set machine information.
type Machine interface {
	MachineGetter
	MachineSetter
}
