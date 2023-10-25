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

package eks

import "github.com/pkg/errors"

var (
	// ErrClusterExists is an error if a EKS cluster already exists with
	// the same name in the spec but that isn't owned by the CAPI cluster.
	ErrClusterExists = errors.New("an EKS cluster already exists with same name but isn't owned by cluster")
	// ErrUnknownTokenMethod defines an error if a unsupported token generation method is supplied.
	ErrUnknownTokenMethod = errors.New("unknown token method")
	// ErrClusterRoleNameMissing if no role name is specified.
	ErrClusterRoleNameMissing = errors.New("a cluster role name must be specified")
	// ErrClusterRoleNotFound is an error if the specified role couldn't be founbd in AWS.
	ErrClusterRoleNotFound = errors.New("the specified cluster role couldn't be found")
	// ErrNodegroupRoleNotFound is an error if the specified role couldn't be founbd in AWS.
	ErrNodegroupRoleNotFound = errors.New("the specified nodegroup role couldn't be found")
	// ErrFargateRoleNotFound is an error if the specified role couldn't be founbd in AWS.
	ErrFargateRoleNotFound = errors.New("the specified fargate role couldn't be found")
	// ErrCannotUseAdditionalRoles is an error if the spec contains additional role and the
	// EKSAllowAddRoles feature flag isn't enabled.
	ErrCannotUseAdditionalRoles = errors.New("additional rules cannot be added as this has been disabled")
	// ErrNoSecurityGroup is an error when no security group is found for an EKS cluster.
	ErrNoSecurityGroup = errors.New("no security group for EKS cluster")
)
