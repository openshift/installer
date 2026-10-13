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

package v1beta1

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/pkg/errors"
)

// Errors for validation of Amazon EKS nodes that are registered with the control plane.
var (
	ErrRoleARNRequired  = errors.New("rolearn is required")
	ErrUserARNRequired  = errors.New("userarn is required")
	ErrUserNameRequired = errors.New("username is required")
	ErrGroupsRequired   = errors.New("groups are required")
	ErrIsNotARN         = errors.New("supplied value is not a ARN")
	ErrIsNotRoleARN     = errors.New("supplied ARN is not a role ARN")
	ErrIsNotUserARN     = errors.New("supplied ARN is not a user ARN")
)

// Validate will return nil is there are no errors with the role mapping.
func (r *RoleMapping) Validate() []error {
	errs := []error{}

	if strings.TrimSpace(r.RoleARN) == "" {
		errs = append(errs, ErrRoleARNRequired)
	}
	if strings.TrimSpace(r.UserName) == "" {
		errs = append(errs, ErrUserNameRequired)
	}
	if len(r.Groups) == 0 {
		errs = append(errs, ErrGroupsRequired)
	}

	if !arn.IsARN(r.RoleARN) {
		errs = append(errs, ErrIsNotARN)
	} else {
		parsedARN, err := arn.Parse(r.RoleARN)
		if err != nil {
			errs = append(errs, err)
		} else if !strings.Contains(parsedARN.Resource, "role/") {
			errs = append(errs, ErrIsNotRoleARN)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

// Validate will return nil is there are no errors with the user mapping.
func (u *UserMapping) Validate() []error {
	errs := []error{}

	if strings.TrimSpace(u.UserARN) == "" {
		errs = append(errs, ErrUserARNRequired)
	}
	if strings.TrimSpace(u.UserName) == "" {
		errs = append(errs, ErrUserNameRequired)
	}
	if len(u.Groups) == 0 {
		errs = append(errs, ErrGroupsRequired)
	}

	if !arn.IsARN(u.UserARN) {
		errs = append(errs, ErrIsNotARN)
	} else {
		parsedARN, err := arn.Parse(u.UserARN)
		if err != nil {
			errs = append(errs, err)
		} else if !strings.Contains(parsedARN.Resource, "user/") {
			errs = append(errs, ErrIsNotUserARN)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}
