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
	v1beta1patch "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/patch"
)

// ROSAOCMRoleConfigScopeParams defines the input parameters used to create a new ROSAOCMRoleConfigScope.
type ROSAOCMRoleConfigScopeParams struct {
	Client            client.Client
	ControllerName    string
	Logger            *logger.Logger
	ROSAOCMRoleConfig *expinfrav1.ROSAOCMRoleConfig
}

// ROSAOCMRoleConfigScope defines the basic context for an actuator to operate upon.
type ROSAOCMRoleConfigScope struct {
	logger.Logger
	Client            client.Client
	controllerName    string
	patchHelper       *v1beta1patch.Helper
	ROSAOCMRoleConfig *expinfrav1.ROSAOCMRoleConfig
	serviceLimiters   throttle.ServiceLimiters
	session           aws.Config
	iamClient         *iam.Client
}

// NewROSAOCMRoleConfigScope creates a new ROSAOCMRoleConfigScope from the supplied parameters.
func NewROSAOCMRoleConfigScope(params ROSAOCMRoleConfigScopeParams) (*ROSAOCMRoleConfigScope, error) {
	if params.Logger == nil {
		log := klog.Background()
		params.Logger = logger.NewLogger(log)
	}

	ocmRoleConfigScope := &ROSAOCMRoleConfigScope{
		Logger:            *params.Logger,
		Client:            params.Client,
		controllerName:    params.ControllerName,
		patchHelper:       nil,
		ROSAOCMRoleConfig: params.ROSAOCMRoleConfig,
	}

	session, serviceLimiters, err := sessionForClusterWithRegion(params.Client, ocmRoleConfigScope, "", params.Logger)
	if err != nil {
		return nil, errors.Errorf("failed to create aws V2 session: %v", err)
	}

	iamClient := iam.NewFromConfig(*session)

	patchHelper, err := v1beta1patch.NewHelper(params.ROSAOCMRoleConfig, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}

	ocmRoleConfigScope.patchHelper = patchHelper
	ocmRoleConfigScope.session = *session
	ocmRoleConfigScope.serviceLimiters = serviceLimiters
	ocmRoleConfigScope.iamClient = iamClient

	return ocmRoleConfigScope, nil
}

// IdentityRef returns the AWSIdentityReference object.
func (s *ROSAOCMRoleConfigScope) IdentityRef() *infrav1.AWSIdentityReference {
	return s.ROSAOCMRoleConfig.Spec.IdentityRef
}

// Session returns the AWS SDK V2 session. Used for creating clients.
func (s *ROSAOCMRoleConfigScope) Session() aws.Config {
	return s.session
}

// IAMClient returns the IAM client.
func (s *ROSAOCMRoleConfigScope) IAMClient() *iam.Client {
	return s.iamClient
}

// CredentialsSecret returns the credentials secret for OCM.
func (s *ROSAOCMRoleConfigScope) CredentialsSecret() *corev1.Secret {
	secretRef := s.ROSAOCMRoleConfig.Spec.CredentialsSecretRef
	if secretRef == nil {
		return nil
	}

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretRef.Name,
			Namespace: secretRef.Namespace,
		},
	}
}

// GetClient returns the controller-runtime client.
func (s *ROSAOCMRoleConfigScope) GetClient() client.Client {
	return s.Client
}

// PatchObject persists the ROSAOCMRoleConfig spec and status.
func (s *ROSAOCMRoleConfigScope) PatchObject() error {
	return s.patchHelper.Patch(context.TODO(), s.ROSAOCMRoleConfig)
}

// ServiceLimiter returns the AWS SDK session. Used for creating clients.
func (s *ROSAOCMRoleConfigScope) ServiceLimiter(service string) *throttle.ServiceLimiter {
	if sl, ok := s.serviceLimiters[service]; ok {
		return sl
	}
	return nil
}

// ControllerName returns the name of the controller that created the scope.
func (s *ROSAOCMRoleConfigScope) ControllerName() string {
	return s.controllerName
}

// InfraCluster returns the ROSAOCMRoleConfig object.
// The method is then used in session.go to set proper Conditions for the ROSAOCMRoleConfig object.
func (s *ROSAOCMRoleConfigScope) InfraCluster() cloud.ClusterObject {
	return s.ROSAOCMRoleConfig
}

// InfraClusterName returns the name of the ROSAOCMRoleConfig object.
// The method is then used in session.go to set the key to the AWS session cache.
func (s *ROSAOCMRoleConfigScope) InfraClusterName() string {
	return s.ROSAOCMRoleConfig.Name
}

// Namespace returns empty string since ROSAOCMRoleConfig is cluster-scoped.
func (s *ROSAOCMRoleConfigScope) Namespace() string {
	return ""
}
