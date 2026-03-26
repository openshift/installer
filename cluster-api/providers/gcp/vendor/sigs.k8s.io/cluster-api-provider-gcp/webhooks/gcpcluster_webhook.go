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

package webhooks

import (
	"context"
	"fmt"
	"reflect"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	infrav1 "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// clusterlog is for logging in this package.
var clusterlog = logf.Log.WithName("gcpcluster-resource")

// SetupWebhookWithManager sets up and registers the webhook with the manager.
func (c *GCPCluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&infrav1.GCPCluster{}).
		WithValidator(c).
		WithDefaulter(c).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-gcpcluster,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=gcpclusters,versions=v1beta1,name=validation.gcpcluster.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta1-gcpcluster,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=gcpclusters,versions=v1beta1,name=default.gcpcluster.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1

// GCPCluster implements a validating and defaulting webhook for GCPCluster.
type GCPCluster struct{}

var (
	_ webhook.CustomValidator = &GCPCluster{}
	_ webhook.CustomDefaulter = &GCPCluster{}
)

// Default implements webhook.CustomDefaulter so a webhook will be registered for the type.
func (*GCPCluster) Default(_ context.Context, _ runtime.Object) error {
	return nil
}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type.
func (*GCPCluster) ValidateCreate(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type.
func (*GCPCluster) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	c, ok := newObj.(*infrav1.GCPCluster)
	if !ok {
		return nil, fmt.Errorf("expected an GCPCluster object but got %T", c)
	}

	clusterlog.Info("validate update", "name", c.Name)
	var allErrs field.ErrorList
	old := oldObj.(*infrav1.GCPCluster)

	if !reflect.DeepEqual(c.Spec.Project, old.Spec.Project) {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "Project"),
				c.Spec.Project, "field is immutable"),
		)
	}

	if !reflect.DeepEqual(c.Spec.Region, old.Spec.Region) {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "Region"),
				c.Spec.Region, "field is immutable"),
		)
	}

	if !reflect.DeepEqual(c.Spec.CredentialsRef, old.Spec.CredentialsRef) {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "CredentialsRef"),
				c.Spec.CredentialsRef, "field is immutable"),
		)
	}

	if !reflect.DeepEqual(c.Spec.LoadBalancer, old.Spec.LoadBalancer) {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "LoadBalancer"),
				c.Spec.LoadBalancer, "field is immutable"),
		)
	}

	if c.Spec.Network.Mtu < int64(1300) {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "Network", "Mtu"),
				c.Spec.Network.Mtu, "field cannot be lesser than 1300"),
		)
	}

	if c.Spec.Network.Mtu > int64(8896) {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "Network", "Mtu"),
				c.Spec.Network.Mtu, "field cannot be greater than 8896"),
		)
	}

	for i, firewallRule := range c.Spec.Network.Firewall.FirewallRules {
		for j, allowRule := range firewallRule.Allowed {
			if allowRule.IPProtocol != infrav1.FirewallProtocolTCP && allowRule.IPProtocol != infrav1.FirewallProtocolUDP &&
				len(allowRule.Ports) > 0 {
				allErrs = append(allErrs,
					field.Invalid(field.NewPath("spec", "Network", "Firewall",
						fmt.Sprintf("FirewallRules[%d]", i), fmt.Sprintf("Allowed[%d]", j),
					),
						allowRule.Ports,
						"field should not exist unless IPProtocol is TCP or UDP"),
				)
			}
		}
		for j, denyRule := range firewallRule.Denied {
			if denyRule.IPProtocol != infrav1.FirewallProtocolTCP && denyRule.IPProtocol != infrav1.FirewallProtocolUDP &&
				len(denyRule.Ports) > 0 {
				allErrs = append(allErrs,
					field.Invalid(field.NewPath("spec", "Network", "Firewall",
						fmt.Sprintf("FirewallRules[%d]", i), fmt.Sprintf("Denied[%d]", j),
					),
						denyRule.Ports,
						"field should not exist unless IPProtocol is TCP or UDP"),
				)
			}
		}
	}

	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(infrav1.GroupVersion.WithKind("GCPCluster").GroupKind(), c.Name, allErrs)
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type.
func (*GCPCluster) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}
