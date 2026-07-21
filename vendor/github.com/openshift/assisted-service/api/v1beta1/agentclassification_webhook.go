/*
Copyright 2020.

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
	"context"
	"fmt"
	"strings"

	"github.com/itchyny/gojq"

	"k8s.io/apimachinery/pkg/apis/meta/v1/validation"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
	ClassificationLabelPrefix = "agentclassification." + Group + "/"
)

// log is for logging in this package.
var agentclassificationlog = logf.Log.WithName("agentclassification-resource")

// agentClassificationWebhook implements webhook.CustomValidator for AgentClassification.
type agentClassificationWebhook struct{}

var _ webhook.CustomValidator = &agentClassificationWebhook{}

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *AgentClassification) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, r).
		WithCustomValidator(&agentClassificationWebhook{}).
		Complete()
}

//+kubebuilder:webhook:path=/validate-agent-install-openshift-io-v1beta1-agentclassification,mutating=false,failurePolicy=fail,sideEffects=None,groups=agent-install.openshift.io,resources=agentclassifications,verbs=create;update,versions=v1beta1,name=vagentclassification.kb.io,admissionReviewVersions=v1

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type
func (w *agentClassificationWebhook) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	r, ok := obj.(*AgentClassification)
	if !ok {
		return nil, fmt.Errorf("object is not an AgentClassification")
	}
	agentclassificationlog.Info("validate create", "name", r.Name)
	f := field.NewPath("spec")
	errs := validation.ValidateLabels(map[string]string{ClassificationLabelPrefix + r.Spec.LabelKey: r.Spec.LabelValue}, f)
	if strings.HasPrefix(r.Spec.LabelValue, "QUERYERROR") {
		errs = append(errs, field.Invalid(f, r.Spec.LabelValue, "label must not start with QUERYERROR as this is reserved"))
	}

	// Validate that we can parse the specified query
	_, err := gojq.Parse(r.Spec.Query)
	if err != nil {
		errs = append(errs, field.Invalid(f, r.Spec.Query, err.Error()))
	}

	if len(errs) > 0 {
		err := fmt.Errorf("Validation failed: %s", errs.ToAggregate().Error())
		return nil, err
	}

	agentclassificationlog.Info("Successful validation")
	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type
func (w *agentClassificationWebhook) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	r, ok := newObj.(*AgentClassification)
	if !ok {
		return nil, fmt.Errorf("new object is not an AgentClassification")
	}
	agentclassificationlog.Info("validate update", "name", r.Name)

	oldAgentClassification, ok := oldObj.(*AgentClassification)
	if !ok {
		return nil, fmt.Errorf("old object is not an AgentClassification")
	}

	// Validate that the label key and value haven't changed
	if (oldAgentClassification.Spec.LabelKey != r.Spec.LabelKey) || (oldAgentClassification.Spec.LabelValue != r.Spec.LabelValue) {
		return nil, fmt.Errorf("Label modified: the specified label may not be modified after creation")
	}

	// If we get here, then all checks passed, so the object is valid.
	agentclassificationlog.Info("Successful validation")
	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type
func (w *agentClassificationWebhook) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}
