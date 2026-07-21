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

	"github.com/openshift/assisted-service/models"
	"github.com/thoas/go-funk"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var agentlog = logf.Log.WithName("agent-resource")

// agentWebhook implements webhook.CustomValidator for Agent.
type agentWebhook struct{}

var _ webhook.CustomValidator = &agentWebhook{}

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *Agent) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, r).
		WithCustomValidator(&agentWebhook{}).
		Complete()
}

//+kubebuilder:webhook:path=/validate-agent-install-openshift-io-v1beta1-agent,mutating=false,failurePolicy=fail,sideEffects=None,groups=agent-install.openshift.io,resources=agents,verbs=create;update,versions=v1beta1,name=vagent.kb.io,admissionReviewVersions=v1

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type
func (w *agentWebhook) ValidateCreate(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type
func (w *agentWebhook) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	r, ok := newObj.(*Agent)
	if !ok {
		return nil, fmt.Errorf("new object is not an Agent")
	}
	agentlog.Info("validate update", "name", r.Name)
	oldObject, ok := oldObj.(*Agent)
	if !ok {
		return nil, fmt.Errorf("old object is not an Agent")
	}
	if !areClusterRefsEqual(oldObject.Spec.ClusterDeploymentName, r.Spec.ClusterDeploymentName) {
		installingStatuses := []string{
			models.HostStatusPreparingForInstallation,
			models.HostStatusPreparingFailed,
			models.HostStatusPreparingSuccessful,
			models.HostStatusInstalling,
			models.HostStatusInstallingInProgress,
			models.HostStatusInstallingPendingUserAction,
		}
		if funk.ContainsString(installingStatuses, r.Status.DebugInfo.State) {
			err := fmt.Errorf("Failed validation: Attempted to change Spec.ClusterDeploymentName which is immutable during Agent installation.")
			agentlog.Info(err.Error())
			return nil, err
		}
	}
	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type
func (w *agentWebhook) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

func areClusterRefsEqual(clusterRef1 *ClusterReference, clusterRef2 *ClusterReference) bool {
	if clusterRef1 == nil && clusterRef2 == nil {
		return true
	} else if clusterRef1 != nil && clusterRef2 != nil {
		return (clusterRef1.Name == clusterRef2.Name && clusterRef1.Namespace == clusterRef2.Namespace)
	} else {
		return false
	}
}
