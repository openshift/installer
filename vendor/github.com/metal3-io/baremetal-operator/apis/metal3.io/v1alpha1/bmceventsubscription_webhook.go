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

	"k8s.io/apimachinery/pkg/util/errors"

	"k8s.io/apimachinery/pkg/runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// bmcsubscriptionlog is for logging in this package.
var bmcsubscriptionlog = logf.Log.WithName("webhooks").WithName("BMCEventSubscription")

//+kubebuilder:webhook:verbs=create;update,path=/validate-metal3-io-v1alpha1-bmceventsubscription,mutating=false,failurePolicy=fail,sideEffects=none,admissionReviewVersions=v1;v1beta,groups=metal3.io,resources=bmceventsubscriptions,versions=v1alpha1,name=bmceventsubscription.metal3.io

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (s *BMCEventSubscription) ValidateCreate() (admission.Warnings, error) {
	bmcsubscriptionlog.Info("validate create", "namespace", s.Namespace, "name", s.Name)
	return nil, errors.NewAggregate(s.validateSubscription())
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
//
// We prevent updates to the spec.  All other updates (e.g. status, finalizers) are allowed.
func (s *BMCEventSubscription) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	bmcsubscriptionlog.Info("validate update", "namespace", s.Namespace, "name", s.Name)

	bes, casted := old.(*BMCEventSubscription)
	if !casted {
		bmcsubscriptionlog.Error(fmt.Errorf("old object conversion error for %s/%s", s.Namespace, s.Name), "validate update error")
		return nil, nil
	}

	if s.Spec != bes.Spec {
		return nil, fmt.Errorf("subscriptions cannot be updated, please recreate it")
	}

	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (s *BMCEventSubscription) ValidateDelete() (admission.Warnings, error) {
	return nil, nil
}
