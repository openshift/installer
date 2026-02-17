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

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *AgentClassification) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/validate-agent-install-openshift-io-v1beta1-agentclassification,mutating=false,failurePolicy=fail,sideEffects=None,groups=agent-install.openshift.io,resources=agentclassifications,verbs=create;update,versions=v1beta1,name=vagentclassification.kb.io,admissionReviewVersions=v1

var _ webhook.CustomValidator = &AgentClassification{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type
func (r *AgentClassification) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	agentClassification, ok := obj.(*AgentClassification)
	if !ok {
		return nil, fmt.Errorf("object is not an AgentClassification")
	}
	agentclassificationlog.Info("validate create", "name", agentClassification.Name)
	f := field.NewPath("spec")
	errs := validation.ValidateLabels(map[string]string{ClassificationLabelPrefix + agentClassification.Spec.LabelKey: agentClassification.Spec.LabelValue}, f)
	if strings.HasPrefix(agentClassification.Spec.LabelValue, "QUERYERROR") {
		errs = append(errs, field.Invalid(f, agentClassification.Spec.LabelValue, "label must not start with QUERYERROR as this is reserved"))
	}

	// Validate that we can parse the specified query
	_, err := gojq.Parse(agentClassification.Spec.Query)
	if err != nil {
		errs = append(errs, field.Invalid(f, agentClassification.Spec.Query, err.Error()))
	}

	if len(errs) > 0 {
		err := fmt.Errorf("Validation failed: %s", errs.ToAggregate().Error())
		return nil, err
	}

	agentclassificationlog.Info("Successful validation")
	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type
func (r *AgentClassification) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	agentClassification, ok := newObj.(*AgentClassification)
	if !ok {
		return nil, fmt.Errorf("new object is not an AgentClassification")
	}
	agentclassificationlog.Info("validate update", "name", agentClassification.Name)

	oldAgentClassification, ok := oldObj.(*AgentClassification)
	if !ok {
		return nil, fmt.Errorf("old object is not an AgentClassification")
	}

	// Validate that the label key and value haven't changed
	if (oldAgentClassification.Spec.LabelKey != agentClassification.Spec.LabelKey) || (oldAgentClassification.Spec.LabelValue != agentClassification.Spec.LabelValue) {
		return nil, fmt.Errorf("Label modified: the specified label may not be modified after creation")
	}

	// If we get here, then all checks passed, so the object is valid.
	agentclassificationlog.Info("Successful validation")
	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type
func (r *AgentClassification) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	return nil, nil
}
