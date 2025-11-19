/*
 Copyright 2025 The Kubernetes Authors.

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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/throttle"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api/util/patch"
)

// RosaRoleConfigScopeParams defines the input parameters used to create a new RosaRoleConfigScope.
type RosaRoleConfigScopeParams struct {
	Client         client.Client
	ControllerName string
	Logger         *logger.Logger
	RosaRoleConfig *expinfrav1.ROSARoleConfig
}

// RosaRoleConfigScope defines the basic context for an actuator to operate upon.
type RosaRoleConfigScope struct {
	logger.Logger
	Client          client.Client
	controllerName  string
	patchHelper     *patch.Helper
	RosaRoleConfig  *expinfrav1.ROSARoleConfig
	serviceLimiters throttle.ServiceLimiters
	session         aws.Config
	iamClient       *iam.Client
}

// NewRosaRoleConfigScope creates a new RosaRoleConfigScope from the supplied parameters.
func NewRosaRoleConfigScope(params RosaRoleConfigScopeParams) (*RosaRoleConfigScope, error) {
	if params.Logger == nil {
		log := klog.Background()
		params.Logger = logger.NewLogger(log)
	}

	RosaRoleConfigScope := &RosaRoleConfigScope{
		Logger:         *params.Logger,
		Client:         params.Client,
		controllerName: params.ControllerName,
		patchHelper:    nil,
		RosaRoleConfig: params.RosaRoleConfig,
	}

	session, serviceLimiters, err := sessionForClusterWithRegion(params.Client, RosaRoleConfigScope, "", params.Logger)

	if err != nil {
		return nil, errors.Errorf("failed to create aws V2 session: %v", err)
	}

	iamClient := iam.NewFromConfig(*session)

	patchHelper, err := patch.NewHelper(params.RosaRoleConfig, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}

	RosaRoleConfigScope.patchHelper = patchHelper
	RosaRoleConfigScope.session = *session
	RosaRoleConfigScope.serviceLimiters = serviceLimiters
	RosaRoleConfigScope.iamClient = iamClient

	return RosaRoleConfigScope, nil
}

// IdentityRef returns the AWSIdentityReference object.
func (s *RosaRoleConfigScope) IdentityRef() *infrav1.AWSIdentityReference {
	return s.RosaRoleConfig.Spec.IdentityRef
}

// Session returns the AWS SDK V2 session. Used for creating clients.
func (s *RosaRoleConfigScope) Session() aws.Config {
	return s.session
}

// ServiceLimiter returns the AWS SDK session (used for creating clients).
func (s *RosaRoleConfigScope) ServiceLimiter(service string) *throttle.ServiceLimiter {
	if sl, ok := s.serviceLimiters[service]; ok {
		return sl
	}
	return nil
}

// ControllerName returns the name of the controller.
func (s *RosaRoleConfigScope) ControllerName() string {
	return s.controllerName
}

// InfraCluster returns the RosaRoleConfig object.
// The method is then used in session.go to set proper Conditions for the RosaRoleConfig object.
func (s *RosaRoleConfigScope) InfraCluster() cloud.ClusterObject {
	return s.RosaRoleConfig
}

// InfraClusterName returns the name of the RosaRoleConfig object.
// The method is then used in session.go to set the key to the AWS session cache.
func (s *RosaRoleConfigScope) InfraClusterName() string {
	return s.RosaRoleConfig.Name
}

// Namespace returns the namespace of the RosaRoleConfig object.
// The method is then used in session.go to set the key to the AWS session cache.
func (s *RosaRoleConfigScope) Namespace() string {
	return s.RosaRoleConfig.Namespace
}

// GetClient Returns RosaRoleConfigScope client.
func (s *RosaRoleConfigScope) GetClient() client.Client {
	return s.Client
}

// PatchObject persists the RosaRoleConfig configuration and status.
func (s *RosaRoleConfigScope) PatchObject() error {
	return s.patchHelper.Patch(
		context.Background(),
		s.RosaRoleConfig)
}

// CredentialsSecret returns the CredentialsSecret object.
func (s *RosaRoleConfigScope) CredentialsSecret() *corev1.Secret {
	secretRef := s.RosaRoleConfig.Spec.CredentialsSecretRef
	if secretRef == nil {
		return nil
	}

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.RosaRoleConfig.Spec.CredentialsSecretRef.Name,
			Namespace: s.RosaRoleConfig.Namespace,
		},
	}
}

// IAMClient returns the IAM client.
func (s *RosaRoleConfigScope) IAMClient() *iam.Client {
	return s.iamClient
}
