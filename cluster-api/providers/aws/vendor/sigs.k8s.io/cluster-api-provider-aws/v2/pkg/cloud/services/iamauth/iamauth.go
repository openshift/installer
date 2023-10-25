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

package iamauth

import (
	crclient "sigs.k8s.io/controller-runtime/pkg/client"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
)

const (
	// EC2NodeUserName is the username required for EC2 nodes.
	EC2NodeUserName = "system:node:{{EC2PrivateDNSName}}"
)

var (
	// NodeGroups is the groups that are required for a node.
	NodeGroups = []string{"system:bootstrappers", "system:nodes"}
)

// AuthenticatorBackend is the interface that represents an aws-iam-authenticator backend.
type AuthenticatorBackend interface {
	// MapRole is used to map a role ARN to a user and set of groups
	MapRole(mapping ekscontrolplanev1.RoleMapping) error
	// MapUser is used to map a user ARN to a user and set of groups
	MapUser(mapping ekscontrolplanev1.UserMapping) error
}

// BackendType is a type that represents the different aws-iam-authenticator backends.
type BackendType string

var (
	// BackendTypeConfigMap is the Kubernetes config map backend.
	BackendTypeConfigMap = BackendType("config-map")
	// BackendTypeCRD is the CRD based backend.
	BackendTypeCRD = BackendType("crd")
)

// NewBackend will create a new authenticate backend for a given type. Only use BackendTypeConfigMap
// with EKS.
func NewBackend(backendType BackendType, client crclient.Client) (AuthenticatorBackend, error) {
	if client == nil {
		return nil, ErrClientRequired
	}

	switch backendType {
	case BackendTypeConfigMap:
		return &configMapBackend{client: client}, nil
	case BackendTypeCRD:
		return &crdBackend{client: client}, nil
	default:
		return nil, ErrInvalidBackendType
	}
}
