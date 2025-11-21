/*
 Copyright The Kubernetes Authors.

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

	awsv2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/throttle"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/patch"
)

// ROSANetworkScopeParams defines the input parameters used to create a new ROSANetworkScope.
type ROSANetworkScopeParams struct {
	Client         client.Client
	ControllerName string
	Logger         *logger.Logger
	ROSANetwork    *expinfrav1.ROSANetwork
}

// ROSANetworkScope defines the basic context for an actuator to operate upon.
type ROSANetworkScope struct {
	logger.Logger
	Client          client.Client
	controllerName  string
	patchHelper     *patch.Helper
	ROSANetwork     *expinfrav1.ROSANetwork
	serviceLimiters throttle.ServiceLimiters
	session         awsv2.Config
}

// NewROSANetworkScope creates a new NewROSANetworkScope from the supplied parameters.
func NewROSANetworkScope(params ROSANetworkScopeParams) (*ROSANetworkScope, error) {
	if params.Logger == nil {
		log := klog.Background()
		params.Logger = logger.NewLogger(log)
	}

	rosaNetworkScope := &ROSANetworkScope{
		Logger:         *params.Logger,
		Client:         params.Client,
		controllerName: params.ControllerName,
		patchHelper:    nil,
		ROSANetwork:    params.ROSANetwork,
	}

	session, serviceLimiters, err := sessionForClusterWithRegion(params.Client, rosaNetworkScope, params.ROSANetwork.Spec.Region, params.Logger)
	if err != nil {
		return nil, errors.Errorf("failed to create aws V2 session: %v", err)
	}

	patchHelper, err := patch.NewHelper(params.ROSANetwork, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}

	rosaNetworkScope.patchHelper = patchHelper
	rosaNetworkScope.session = *session
	rosaNetworkScope.serviceLimiters = serviceLimiters

	return rosaNetworkScope, nil
}

// Session returns the AWS SDK V2 Config. Used for creating clients.
func (s *ROSANetworkScope) Session() awsv2.Config {
	return s.session
}

// IdentityRef returns the AWSIdentityReference object.
func (s *ROSANetworkScope) IdentityRef() *infrav1.AWSIdentityReference {
	return s.ROSANetwork.Spec.IdentityRef
}

// ServiceLimiter returns the AWS SDK session (used for creating clients).
func (s *ROSANetworkScope) ServiceLimiter(service string) *throttle.ServiceLimiter {
	if sl, ok := s.serviceLimiters[service]; ok {
		return sl
	}
	return nil
}

// ControllerName returns the name of the controller.
func (s *ROSANetworkScope) ControllerName() string {
	return s.controllerName
}

// InfraCluster returns the ROSANetwork object.
// The method is then used in session.go to set proper Conditions for the ROSANetwork object.
func (s *ROSANetworkScope) InfraCluster() cloud.ClusterObject {
	return s.ROSANetwork
}

// InfraClusterName returns the name of the ROSANetwork object.
// The method is then used in session.go to set the key to the AWS session cache.
func (s *ROSANetworkScope) InfraClusterName() string {
	return s.ROSANetwork.Name
}

// Namespace returns the namespace of the ROSANetwork object.
// The method is then used in session.go to set the key to the AWS session cache.
func (s *ROSANetworkScope) Namespace() string {
	return s.ROSANetwork.Namespace
}

// PatchObject persists the rosanetwork configuration and status.
func (s *ROSANetworkScope) PatchObject() error {
	return s.patchHelper.Patch(
		context.TODO(),
		s.ROSANetwork,
		patch.WithOwnedConditions{Conditions: []clusterv1.ConditionType{
			expinfrav1.ROSANetworkReadyCondition,
		}})
}
