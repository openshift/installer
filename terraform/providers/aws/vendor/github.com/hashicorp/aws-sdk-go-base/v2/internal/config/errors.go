package config

import (
	"fmt"
)

// CannotAssumeRoleError occurs when AssumeRole cannot complete.
type CannotAssumeRoleError struct {
	Config *Config
	Err    error
}

func (e CannotAssumeRoleError) Error() string {
	if e.Config == nil || e.Config.AssumeRole == nil {
		return fmt.Sprintf("cannot assume role: %s", e.Err)
	}

	return fmt.Sprintf(`IAM Role (%s) cannot be assumed.

There are a number of possible causes of this - the most common are:
  * The credentials used in order to assume the role are invalid
  * The credentials do not have appropriate permission to assume the role
  * The role ARN is not valid

Error: %s
`, e.Config.AssumeRole.RoleARN, e.Err)
}

func (e CannotAssumeRoleError) Unwrap() error {
	return e.Err
}

func (c *Config) NewCannotAssumeRoleError(err error) CannotAssumeRoleError {
	return CannotAssumeRoleError{Config: c, Err: err}
}

// NoValidCredentialSourcesError occurs when all credential lookup methods have been exhausted without results.
type NoValidCredentialSourcesError struct {
	Config *Config
	Err    error
}

func (e NoValidCredentialSourcesError) Error() string {
	if e.Config == nil {
		return fmt.Sprintf("no valid credential sources found: %s", e.Err)
	}

	return fmt.Sprintf(`no valid credential sources for %[1]s found.

Please see %[2]s
for more information about providing credentials.

Error: %[3]s
`, e.Config.CallerName, e.Config.CallerDocumentationURL, e.Err)
}

func (e NoValidCredentialSourcesError) Unwrap() error {
	return e.Err
}

func (c *Config) NewNoValidCredentialSourcesError(err error) NoValidCredentialSourcesError {
	return NoValidCredentialSourcesError{Config: c, Err: err}
}
