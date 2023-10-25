/*
Copyright 2023 The Kubernetes Authors.

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

package scope

import (
	"context"

	amazoncni "github.com/aws/amazon-vpc-cni-k8s/pkg/apis/crd/v1alpha1"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	rosacontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/patch"
)

func init() {
	_ = amazoncni.AddToScheme(scheme)
	_ = appsv1.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)
	_ = rbacv1.AddToScheme(scheme)
}

type ROSAControlPlaneScopeParams struct {
	Client         client.Client
	Logger         *logger.Logger
	Cluster        *clusterv1.Cluster
	ControlPlane   *rosacontrolplanev1.ROSAControlPlane
	ControllerName string
}

func NewROSAControlPlaneScope(params ROSAControlPlaneScopeParams) (*ROSAControlPlaneScope, error) {
	if params.Cluster == nil {
		return nil, errors.New("failed to generate new scope from nil Cluster")
	}
	if params.ControlPlane == nil {
		return nil, errors.New("failed to generate new scope from nil AWSManagedControlPlane")
	}
	if params.Logger == nil {
		log := klog.Background()
		params.Logger = logger.NewLogger(log)
	}

	managedScope := &ROSAControlPlaneScope{
		Logger:       *params.Logger,
		Client:       params.Client,
		Cluster:      params.Cluster,
		ControlPlane: params.ControlPlane,
		patchHelper:  nil,
	}

	helper, err := patch.NewHelper(params.ControlPlane, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}

	managedScope.patchHelper = helper
	return managedScope, nil
}

// ROSAControlPlaneScope defines the basic context for an actuator to operate upon.
type ROSAControlPlaneScope struct {
	logger.Logger
	Client      client.Client
	patchHelper *patch.Helper

	Cluster      *clusterv1.Cluster
	ControlPlane *rosacontrolplanev1.ROSAControlPlane
}

// Name returns the CAPI cluster name.
func (s *ROSAControlPlaneScope) Name() string {
	return s.Cluster.Name
}

// InfraClusterName returns the AWS cluster name.

func (s *ROSAControlPlaneScope) InfraClusterName() string {
	return s.ControlPlane.Name
}

// Namespace returns the cluster namespace.
func (s *ROSAControlPlaneScope) Namespace() string {
	return s.Cluster.Namespace
}

// PatchObject persists the control plane configuration and status.
func (s *ROSAControlPlaneScope) PatchObject() error {
	return s.patchHelper.Patch(
		context.TODO(),
		s.ControlPlane,
		patch.WithOwnedConditions{Conditions: []clusterv1.ConditionType{}})
}

// Close closes the current scope persisting the control plane configuration and status.
func (s *ROSAControlPlaneScope) Close() error {
	return s.PatchObject()
}
