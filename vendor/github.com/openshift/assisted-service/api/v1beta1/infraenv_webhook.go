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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var infraenvlog = logf.Log.WithName("infraenv-resource")

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *InfraEnv) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/validate-agent-install-openshift-io-v1beta1-infraenv,mutating=false,failurePolicy=fail,sideEffects=None,groups=agent-install.openshift.io,resources=infraenvs,verbs=create;update,versions=v1beta1,name=vinfraenv.kb.io,admissionReviewVersions=v1

var _ webhook.CustomValidator = &InfraEnv{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type
func (r *InfraEnv) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	infraEnv, ok := obj.(*InfraEnv)
	if !ok {
		return nil, fmt.Errorf("object is not an InfraEnv")
	}
	infraenvlog.Info("validate create", "name", infraEnv.Name)
	if infraEnv.Spec.ClusterRef != nil && infraEnv.Spec.OSImageVersion != "" {
		err := fmt.Errorf("Failed validation: Either Spec.ClusterRef or Spec.OSImageVersion should be specified (not both).")
		infraenvlog.Info(err.Error())
		return nil, err
	}
	infraenvlog.Info("Successful validation")
	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type
func (r *InfraEnv) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	infraEnv, ok := newObj.(*InfraEnv)
	if !ok {
		return nil, fmt.Errorf("new object is not an InfraEnv")
	}
	infraenvlog.Info("validate update", "name", infraEnv.Name)
	oldInfraEnv, ok := oldObj.(*InfraEnv)
	if !ok {
		return nil, fmt.Errorf("old object is not an InfraEnv")
	}
	if !areClusterRefsEqual(oldInfraEnv.Spec.ClusterRef, infraEnv.Spec.ClusterRef) {
		err := fmt.Errorf("Failed validation: Attempted to change Spec.ClusterRef which is immutable after InfraEnv creation.")
		return nil, err
	}
	infraenvlog.Info("Successful validation")
	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type
func (r *InfraEnv) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	return nil, nil
}
