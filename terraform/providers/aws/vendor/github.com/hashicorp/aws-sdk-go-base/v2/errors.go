// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package awsbase

import (
	"errors"

	"github.com/hashicorp/aws-sdk-go-base/v2/internal/config"
)

// CannotAssumeRoleError occurs when AssumeRole cannot complete.
type CannotAssumeRoleError = config.CannotAssumeRoleError

// IsCannotAssumeRoleError returns true if the error contains the CannotAssumeRoleError type.
func IsCannotAssumeRoleError(err error) bool {
	var e CannotAssumeRoleError
	return errors.As(err, &e)
}

// NoValidCredentialSourcesError occurs when all credential lookup methods have been exhausted without results.
type NoValidCredentialSourcesError = config.NoValidCredentialSourcesError

// IsNoValidCredentialSourcesError returns true if the error contains the NoValidCredentialSourcesError type.
func IsNoValidCredentialSourcesError(err error) bool {
	var e NoValidCredentialSourcesError
	return errors.As(err, &e)
}
