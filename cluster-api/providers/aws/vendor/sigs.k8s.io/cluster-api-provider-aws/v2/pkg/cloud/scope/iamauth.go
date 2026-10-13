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

package scope

import (
	"sigs.k8s.io/controller-runtime/pkg/client"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
)

// IAMAuthScope is the interface for the scope to be used with iamauth reconciling service.
type IAMAuthScope interface {
	cloud.ClusterScoper

	// RemoteClient returns the Kubernetes client for connecting to the workload cluster.
	RemoteClient() (client.Client, error)
	// IAMAuthConfig returns the IAM authenticator config
	IAMAuthConfig() *ekscontrolplanev1.IAMAuthenticatorConfig
}
