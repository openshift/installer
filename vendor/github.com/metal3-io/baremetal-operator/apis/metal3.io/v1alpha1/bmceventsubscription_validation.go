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
	"fmt"
	"net/url"

	"github.com/pkg/errors"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// bmclog is for logging in this package.
var bmclog = logf.Log.WithName("bmceventsubscription-validation")

// validateSubscription validates BMCEventSubscription resource for creation
func (s *BMCEventSubscription) validateSubscription() []error {
	bmclog.Info("validate create", "name", s.Name)
	var errs []error

	if s.Spec.HostName == "" {
		errs = append(errs, fmt.Errorf("HostName cannot be empty"))
	}

	if s.Spec.Destination == "" {
		errs = append(errs, fmt.Errorf("Destination cannot be empty"))
	} else {
		destinationUrl, err := url.ParseRequestURI(s.Spec.Destination)

		if err != nil {
			errs = append(errs, errors.Wrap(err, "Destination is an invalid URL"))
		} else {
			if destinationUrl.Path == "" {
				errs = append(errs, fmt.Errorf("Hostname-only destination must have a trailing slash"))
			}
		}
	}

	return errs
}
