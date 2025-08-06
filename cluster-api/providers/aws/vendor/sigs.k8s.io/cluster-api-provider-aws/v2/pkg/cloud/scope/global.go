/*
Copyright 2020 The Kubernetes Authors.

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

// Package scope provides a global scope for CAPA controllers.
package scope

import (
	awsv2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/pkg/errors"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/throttle"
)

// NewGlobalScope creates a new Scope from the supplied parameters.
func NewGlobalScope(params GlobalScopeParams) (*GlobalScope, error) {
	if params.Region == "" {
		return nil, errors.New("region required to create session")
	}
	if params.ControllerName == "" {
		return nil, errors.New("controller name required to generate global scope")
	}

	ns2, limiters, err := sessionForRegion(params.Region)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create aws V2 session")
	}
	return &GlobalScope{
		session:         *ns2,
		serviceLimiters: limiters,
		controllerName:  params.ControllerName,
	}, nil
}

// GlobalScopeParams defines the parameters acceptable for GlobalScope.
type GlobalScopeParams struct {
	ControllerName string
	Region         string
}

// GlobalScope defines the specs for the GlobalScope.
type GlobalScope struct {
	session         awsv2.Config
	serviceLimiters throttle.ServiceLimiters
	controllerName  string
}

// Session returns the AWS SDK V2 config. Used for creating clients.
func (s *GlobalScope) Session() awsv2.Config {
	return s.session
}

// ServiceLimiter returns the AWS SDK session. Used for creating clients.
func (s *GlobalScope) ServiceLimiter(service string) *throttle.ServiceLimiter {
	if sl, ok := s.serviceLimiters[service]; ok {
		return sl
	}
	return nil
}

// ControllerName returns the name of the controller that
// created the GlobalScope.
func (s *GlobalScope) ControllerName() string {
	return s.controllerName
}
