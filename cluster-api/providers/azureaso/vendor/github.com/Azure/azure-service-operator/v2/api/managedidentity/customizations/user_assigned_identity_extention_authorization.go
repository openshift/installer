/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package customizations

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	v20230131s "github.com/Azure/azure-service-operator/v2/api/managedidentity/v1api20230131/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	. "github.com/Azure/azure-service-operator/v2/internal/logging"
	"github.com/Azure/azure-service-operator/v2/internal/set"
	"github.com/Azure/azure-service-operator/v2/internal/util/to"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/secrets"
)

var _ genruntime.KubernetesSecretExporter = &UserAssignedIdentityExtension{}

func (ext *UserAssignedIdentityExtension) ExportKubernetesSecrets(
	ctx context.Context,
	obj genruntime.MetaObject,
	_ set.Set[string],
	armClient *genericarmclient.GenericClient,
	log logr.Logger,
) (*genruntime.KubernetesSecretExportResult, error) {
	typedObj, ok := obj.(*v20230131s.UserAssignedIdentity)
	if !ok {
		return nil, fmt.Errorf(
			"cannot run on unknown resource type %T, expected *v20230131s.UserAssignedIdentity", obj)
	}

	// Type assert that we are the hub type. This will fail to compile if
	// the hub type has been changed but this extension has not
	var _ conversion.Hub = typedObj

	hasSecrets := secretsSpecified(typedObj)
	if !hasSecrets {
		log.V(Debug).Info("No secrets retrieval to perform as operatorSpec.Secrets is empty")
		return nil, nil
	}

	collector := secrets.NewCollector(typedObj.Namespace)
	if typedObj.Spec.OperatorSpec != nil && typedObj.Spec.OperatorSpec.Secrets != nil {
		collector.AddValue(typedObj.Spec.OperatorSpec.Secrets.ClientId, to.Value(typedObj.Status.ClientId))
		collector.AddValue(typedObj.Spec.OperatorSpec.Secrets.PrincipalId, to.Value(typedObj.Status.PrincipalId))
		collector.AddValue(typedObj.Spec.OperatorSpec.Secrets.TenantId, to.Value(typedObj.Status.TenantId))
	}

	result, err := collector.Values()
	if err != nil {
		return nil, err
	}
	return &genruntime.KubernetesSecretExportResult{
		Objs: secrets.SliceToClientObjectSlice(result),
	}, nil
}

func secretsSpecified(obj *v20230131s.UserAssignedIdentity) bool {
	if obj.Spec.OperatorSpec == nil || obj.Spec.OperatorSpec.Secrets == nil {
		return false
	}

	specSecrets := obj.Spec.OperatorSpec.Secrets
	hasSecrets := false
	if specSecrets.ClientId != nil ||
		specSecrets.PrincipalId != nil ||
		specSecrets.TenantId != nil {
		hasSecrets = true
	}

	return hasSecrets
}
