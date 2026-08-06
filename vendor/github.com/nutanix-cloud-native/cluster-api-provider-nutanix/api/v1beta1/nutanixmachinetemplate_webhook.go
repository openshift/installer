/*
Copyright 2026 Nutanix

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

	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/nutanix-cloud-native/cluster-api-provider-nutanix/pkg/feature"
)

// Placeholder values for brownfield compatibility.
// These are injected by the defaulting webhook when a NutanixMachineTemplate has
// spec.template.spec.image with type set but the corresponding name/uuid value missing.
// This can happen with legacy objects created before the CEL validation rules were added.
// The placeholders satisfy CEL validation at admission time; the actual image resolution
// is expected to be handled later by ClusterClass variable patches or the controller.
const (
	// ImageNamePlaceholder is the placeholder value set for image.name when type is "name"
	// but name is not provided. It is clearly identifiable as a non-real value.
	ImageNamePlaceholder = "PLACEHOLDER"

	// ImageUUIDPlaceholder is the placeholder UUID set for image.uuid when type is "uuid"
	// but uuid is not provided. It uses the nil UUID format (all zeros) which is a valid
	// UUID format that satisfies the CEL rule requiring 36 characters and a dash.
	ImageUUIDPlaceholder = "00000000-0000-0000-0000-000000000000"
)

// NutanixMachineTemplateDefaulter implements a defaulting webhook for NutanixMachineTemplate.
//
// This is a brownfield compatibility workaround: existing NutanixMachineTemplate objects
// in production may have spec.template.spec.image with type set to "name" or "uuid" but
// the corresponding value field (name or uuid) unset. The CEL validation rules on
// NutanixResourceIdentifier require these fields to be present. This webhook fills in
// placeholder values at admission time so the CEL rules pass, without weakening the
// API contract globally.
//
// The defaulter is idempotent and will not overwrite valid user-provided values.
//
// Each defaulting behavior is individually controlled via a Kubernetes feature gate:
//   - feature.DefaultToPlaceholderImageName controls defaulting for type "name"
//   - feature.DefaultToPlaceholderImageUUID controls defaulting for type "uuid"
//
// Enable via: --feature-gates=DefaultToPlaceholderImageName=true,DefaultToPlaceholderImageUUID=true
//
// +kubebuilder:object:generate=false
type NutanixMachineTemplateDefaulter struct{}

var _ admission.Defaulter[*NutanixMachineTemplate] = &NutanixMachineTemplateDefaulter{}

// SetupWebhookWithManager registers the defaulting webhook for NutanixMachineTemplate.
//
// +kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1beta1-nutanixmachinetemplate,mutating=true,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=nutanixmachinetemplates,verbs=create;update,versions=v1beta1,name=default.nutanixmachinetemplate.infrastructure.cluster.x-k8s.io,admissionReviewVersions=v1
func (d *NutanixMachineTemplateDefaulter) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &NutanixMachineTemplate{}).
		WithDefaulter(d).
		Complete()
}

// Default implements admission.Defaulter.
// It normalizes spec.template.spec.image for brownfield NutanixMachineTemplate objects
// that have the image type set but the corresponding identifier value missing.
// Each defaulting behavior is controlled by its respective feature gate.
func (d *NutanixMachineTemplateDefaulter) Default(_ context.Context, nmt *NutanixMachineTemplate) error {
	defaultNutanixMachineTemplateImage(nmt,
		feature.Gates.Enabled(feature.DefaultToPlaceholderImageName),
		feature.Gates.Enabled(feature.DefaultToPlaceholderImageUUID),
	)
	return nil
}

// defaultNutanixMachineTemplateImage applies placeholder defaults to
// spec.template.spec.image when the identifier value is missing or empty.
// This function is scoped narrowly to the image field only.
// Each defaulting behavior is individually gated by the corresponding feature flag.
//
// For brownfield upgrade: older CRDs allowed image.name or image.uuid to be empty
// (or omitted). After upgrading to a CAPX version with CEL + MinLength=1, existing
// objects with name="" or uuid="" remain in etcd but any update would fail validation.
// Treating empty string as "missing" here ensures the first update after upgrade
// gets defaulted to the placeholder and passes validation.
func defaultNutanixMachineTemplateImage(nmt *NutanixMachineTemplate, defaultName, defaultUUID bool) {
	image := nmt.Spec.Template.Spec.Image
	if image == nil {
		return
	}

	switch image.Type {
	case NutanixIdentifierName:
		nameMissingOrEmpty := image.Name == nil || (image.Name != nil && ptr.Deref(image.Name, "") == "")
		if defaultName && nameMissingOrEmpty {
			placeholder := ImageNamePlaceholder
			image.Name = &placeholder
		}
	case NutanixIdentifierUUID:
		uuidMissingOrEmpty := image.UUID == nil || (image.UUID != nil && ptr.Deref(image.UUID, "") == "")
		if defaultUUID && uuidMissingOrEmpty {
			placeholder := ImageUUIDPlaceholder
			image.UUID = &placeholder
		}
	}
}
