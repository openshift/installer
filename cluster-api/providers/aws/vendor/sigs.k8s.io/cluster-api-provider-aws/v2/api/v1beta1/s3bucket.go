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

	"sigs.k8s.io/cluster-api-provider-aws/v2/feature"
)

// Validate validates S3Bucket fields.
func (b *S3Bucket) Validate() []*field.Error {
	var errs field.ErrorList

	if b == nil {
		return errs
	}

	if b.Name == "" {
		errs = append(errs, field.Required(field.NewPath("spec", "s3Bucket", "name"), "can't be empty"))
	}

	// Feature gate is not enabled but ignition is enabled then send a forbidden error.
	if !feature.Gates.Enabled(feature.BootstrapFormatIgnition) {
		errs = append(errs, field.Forbidden(field.NewPath("spec", "s3Bucket"),
			"can be set only if the BootstrapFormatIgnition feature gate is enabled"))
	}

	if b.ControlPlaneIAMInstanceProfile == "" {
		errs = append(errs,
			field.Required(field.NewPath("spec", "s3Bucket", "controlPlaneIAMInstanceProfiles"), "can't be empty"))
	}

	if len(b.NodesIAMInstanceProfiles) == 0 {
		errs = append(errs,
			field.Required(field.NewPath("spec", "s3Bucket", "nodesIAMInstanceProfiles"), "can't be empty"))
	}

	for i, iamInstanceProfile := range b.NodesIAMInstanceProfiles {
		if iamInstanceProfile == "" {
			errs = append(errs,
				field.Required(field.NewPath("spec", "s3Bucket", fmt.Sprintf("nodesIAMInstanceProfiles[%d]", i)), "can't be empty"))
		}
	}

	if b.Name != "" {
		errs = append(errs, validateS3BucketName(b.Name)...)
	}

	return errs
}

// Validation rules taken from https://docs.aws.amazon.com/AmazonS3/latest/userguide/bucketnamingrules.html.
func validateS3BucketName(name string) []*field.Error {
	var errs field.ErrorList

	path := field.NewPath("spec", "s3Bucket", "name")

	if net.ParseIP(name) != nil {
		errs = append(errs, field.Invalid(path, name, "must not be formatted as an IP address (for example, 192.168.5.4)"))
	}

	return errs
}
