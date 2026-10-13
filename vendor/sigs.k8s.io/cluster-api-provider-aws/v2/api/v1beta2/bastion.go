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

package v1beta2

import (
	"fmt"
	"net"
	"regexp"

	"k8s.io/apimachinery/pkg/util/validation/field"
)

var (
	sshKeyValidNameRegex = regexp.MustCompile(`^[[:graph:]]+([[:print:]]*[[:graph:]]+)*$`)
)

// Validate will validate the bastion fields.
func (b *Bastion) Validate() []*field.Error {
	var errs field.ErrorList

	if b.DisableIngressRules && len(b.AllowedCIDRBlocks) > 0 {
		errs = append(errs,
			field.Forbidden(field.NewPath("spec", "bastion", "allowedCIDRBlocks"), "cannot be set if spec.bastion.disableIngressRules is true"),
		)
		return errs
	}

	for i, cidr := range b.AllowedCIDRBlocks {
		if _, _, err := net.ParseCIDR(cidr); err != nil {
			errs = append(errs,
				field.Invalid(field.NewPath("spec", "bastion", fmt.Sprintf("allowedCIDRBlocks[%d]", i)), cidr, "must be a valid CIDR block"),
			)
		}
	}
	return errs
}

func validateSSHKeyName(sshKeyName *string) field.ErrorList {
	var allErrs field.ErrorList
	switch {
	case sshKeyName == nil:
	// nil is accepted
	case sshKeyName != nil && *sshKeyName == "":
	// empty string is accepted
	case sshKeyName != nil && !sshKeyValidNameRegex.MatchString(*sshKeyName):
		allErrs = append(allErrs, field.Invalid(field.NewPath("sshKeyName"), sshKeyName, "Name is invalid. Must be specified in ASCII and must not start or end in whitespace"))
	}
	return allErrs
}
