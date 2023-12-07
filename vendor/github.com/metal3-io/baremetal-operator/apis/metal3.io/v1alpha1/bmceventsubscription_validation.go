/*

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

package v1alpha1

import (
	"errors"
	"fmt"
	"net/url"
)

// validateSubscription validates BMCEventSubscription resource for creation
func (s *BMCEventSubscription) validateSubscription() []error {
	var errs []error

	if s.Spec.HostName == "" {
		errs = append(errs, errors.New("hostName cannot be empty"))
	}

	if s.Spec.Destination == "" {
		errs = append(errs, errors.New("destination cannot be empty"))
	} else {
		destinationUrl, err := url.ParseRequestURI(s.Spec.Destination)

		if err != nil {
			errs = append(errs, fmt.Errorf("destination is invalid: %w", err))
		} else {
			if destinationUrl.Path == "" {
				errs = append(errs, errors.New("hostname-only destination must have a trailing slash"))
			}
		}
	}

	return errs
}
