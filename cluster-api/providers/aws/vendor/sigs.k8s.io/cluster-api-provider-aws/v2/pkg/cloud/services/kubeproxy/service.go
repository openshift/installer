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

// Package kubeproxy provides a way to interact with the kube-proxy service.
package kubeproxy

import (
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
)

// Service defines the spec for a service.
type Service struct {
	scope scope.KubeProxyScope
}

// NewService will create a new service.
func NewService(kubeproxyScope scope.KubeProxyScope) *Service {
	return &Service{
		scope: kubeproxyScope,
	}
}
