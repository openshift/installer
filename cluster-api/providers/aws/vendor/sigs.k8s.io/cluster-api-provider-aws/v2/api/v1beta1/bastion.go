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
	"fmt"
	"net"

	"k8s.io/apimachinery/pkg/util/validation/field"
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
