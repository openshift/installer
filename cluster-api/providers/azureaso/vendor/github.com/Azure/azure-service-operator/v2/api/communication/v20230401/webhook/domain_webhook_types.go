/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package webhook

import (
	"context"
	"strings"

	"github.com/rotisserie/eris"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	v20230401 "github.com/Azure/azure-service-operator/v2/api/communication/v20230401"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

var _ genruntime.Validator[*v20230401.Domain] = &Domain{}

// CreateValidations returns validation functions for Domain creation.
func (domain *Domain) CreateValidations() []func(ctx context.Context, obj *v20230401.Domain) (admission.Warnings, error) {
	return []func(ctx context.Context, obj *v20230401.Domain) (admission.Warnings, error){
		domain.validateAzureManagedDomainName,
	}
}

// UpdateValidations returns validation functions for Domain updates.
func (domain *Domain) UpdateValidations() []func(ctx context.Context, oldObj *v20230401.Domain, newObj *v20230401.Domain) (admission.Warnings, error) {
	return []func(ctx context.Context, oldObj *v20230401.Domain, newObj *v20230401.Domain) (admission.Warnings, error){
		func(ctx context.Context, oldObj *v20230401.Domain, newObj *v20230401.Domain) (admission.Warnings, error) {
			return domain.validateAzureManagedDomainName(ctx, newObj)
		},
	}
}

// DeleteValidations returns validation functions for Domain deletion.
func (domain *Domain) DeleteValidations() []func(ctx context.Context, obj *v20230401.Domain) (admission.Warnings, error) {
	return nil
}

// validateAzureManagedDomainName validates that when DomainManagement is AzureManaged,
// the AzureName must be "AzureManagedDomain" as required by the Azure API.
func (domain *Domain) validateAzureManagedDomainName(_ context.Context, obj *v20230401.Domain) (admission.Warnings, error) {
	if obj.Spec.DomainManagement == nil {
		return nil, nil
	}

	if *obj.Spec.DomainManagement == v20230401.DomainManagement_AzureManaged && !strings.EqualFold(obj.Spec.AzureName, "AzureManagedDomain") {
		return nil, eris.Errorf(
			"when spec.domainManagement is %q, spec.azureName must be %q, but got %q",
			v20230401.DomainManagement_AzureManaged,
			"AzureManagedDomain",
			obj.Spec.AzureName)
	}

	return nil, nil
}
