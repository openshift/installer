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

var _ genruntime.Validator[*v20230401.SenderUsername] = &SenderUsername{}

// CreateValidations returns validation functions for SenderUsername creation.
func (username *SenderUsername) CreateValidations() []func(ctx context.Context, obj *v20230401.SenderUsername) (admission.Warnings, error) {
	return []func(ctx context.Context, obj *v20230401.SenderUsername) (admission.Warnings, error){
		username.validateAzureNameMatchesUsername,
	}
}

// UpdateValidations returns validation functions for SenderUsername updates.
func (username *SenderUsername) UpdateValidations() []func(ctx context.Context, oldObj *v20230401.SenderUsername, newObj *v20230401.SenderUsername) (admission.Warnings, error) {
	return []func(ctx context.Context, oldObj *v20230401.SenderUsername, newObj *v20230401.SenderUsername) (admission.Warnings, error){
		func(ctx context.Context, oldObj *v20230401.SenderUsername, newObj *v20230401.SenderUsername) (admission.Warnings, error) {
			return username.validateAzureNameMatchesUsername(ctx, newObj)
		},
	}
}

// DeleteValidations returns validation functions for SenderUsername deletion.
func (username *SenderUsername) DeleteValidations() []func(ctx context.Context, obj *v20230401.SenderUsername) (admission.Warnings, error) {
	return nil
}

// validateAzureNameMatchesUsername validates that spec.azureName matches spec.username,
// as required by the Azure API (the resource name in the URL must match the username field).
func (username *SenderUsername) validateAzureNameMatchesUsername(_ context.Context, obj *v20230401.SenderUsername) (admission.Warnings, error) {
	if obj.Spec.Username == nil {
		return nil, nil
	}

	if obj.Spec.AzureName != "" && !strings.EqualFold(obj.Spec.AzureName, *obj.Spec.Username) {
		return nil, eris.Errorf(
			"spec.azureName (%q) must match spec.username (%q) (case-insensitive) — Azure requires the resource name to match the username",
			obj.Spec.AzureName,
			*obj.Spec.Username)
	}

	return nil, nil
}
