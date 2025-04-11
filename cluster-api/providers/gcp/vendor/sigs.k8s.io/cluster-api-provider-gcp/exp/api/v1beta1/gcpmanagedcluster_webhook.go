/*
Copyright 2022 The Kubernetes Authors.

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
	"github.com/google/go-cmp/cmp"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var gcpmanagedclusterlog = logf.Log.WithName("gcpmanagedcluster-resource")

func (r *GCPManagedCluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1beta1-gcpmanagedcluster,mutating=true,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=gcpmanagedclusters,verbs=create;update,versions=v1beta1,name=mgcpmanagedcluster.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &GCPManagedCluster{}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (r *GCPManagedCluster) Default() {
	gcpmanagedclusterlog.Info("default", "name", r.Name)
}

//+kubebuilder:webhook:path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-gcpmanagedcluster,mutating=false,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=gcpmanagedclusters,verbs=create;update,versions=v1beta1,name=vgcpmanagedcluster.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &GCPManagedCluster{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *GCPManagedCluster) ValidateCreate() (admission.Warnings, error) {
	gcpmanagedclusterlog.Info("validate create", "name", r.Name)

	return r.validate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *GCPManagedCluster) ValidateUpdate(oldRaw runtime.Object) (admission.Warnings, error) {
	gcpmanagedclusterlog.Info("validate update", "name", r.Name)
	var allErrs field.ErrorList
	old := oldRaw.(*GCPManagedCluster)

	if !cmp.Equal(r.Spec.Project, old.Spec.Project) {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "Project"),
				r.Spec.Project, "field is immutable"),
		)
	}

	if !cmp.Equal(r.Spec.Region, old.Spec.Region) {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "Region"),
				r.Spec.Region, "field is immutable"),
		)
	}

	if !cmp.Equal(r.Spec.CredentialsRef, old.Spec.CredentialsRef) {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "CredentialsRef"),
				r.Spec.CredentialsRef, "field is immutable"),
		)
	}

	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(GroupVersion.WithKind("GCPManagedCluster").GroupKind(), r.Name, allErrs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (r *GCPManagedCluster) ValidateDelete() (admission.Warnings, error) {
	gcpmanagedclusterlog.Info("validate delete", "name", r.Name)

	return nil, nil
}

func (r *GCPManagedCluster) validate() (admission.Warnings, error) {
	validators := []func() error{
		r.validateCustomSubnet,
	}

	var errs []error
	for _, validator := range validators {
		if err := validator(); err != nil {
			errs = append(errs, err)
		}
	}

	return nil, kerrors.NewAggregate(errs)
}

func (r *GCPManagedCluster) validateCustomSubnet() error {
	gcpmanagedclusterlog.Info("validate custom subnet", "name", r.Name)
	if r.Spec.Network.AutoCreateSubnetworks == nil || *r.Spec.Network.AutoCreateSubnetworks {
		return nil
	}
	isSubnetExistInClusterRegion := false
	for _, subnet := range r.Spec.Network.Subnets {
		if subnet.Region == r.Spec.Region {
			isSubnetExistInClusterRegion = true
		}
	}

	if !isSubnetExistInClusterRegion {
		return field.Required(field.NewPath("spec", "network", "subnet"), "at least one given subnets region should be same as spec.network.region when spec.network.autoCreateSubnetworks is false")
	}
	return nil
}
