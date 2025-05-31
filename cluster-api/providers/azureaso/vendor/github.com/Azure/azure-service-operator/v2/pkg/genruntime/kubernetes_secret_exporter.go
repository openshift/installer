/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package genruntime

import (
	"context"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/set"
)

type KubernetesSecretExportResult struct {
	// Objs is the set of objects (secrets) to export.
	// Only secrets defined on the operatorSpec.secrets are included here. Secrets referenced via a "secret expression"
	// in operatorSpec.secretExpressions are returned in RawSecrets for later use.
	Objs []client.Object

	// RawSecrets contains the raw secret values from Azure.
	// The keys are the "names" of the secrets as defined on operatorSpec.secrets (JSON-cased), and the
	// values are the actual secrets. So for example ManagedCluster has "adminCredentials" and "userCredentials".
	// This will ONLY contain secrets that were requested via additionalSecrets, NOT secrets requested via
	// self.spec.operatorSpec.secrets.
	RawSecrets map[string]string
}

// KubernetesSecretExporter defines a resource which can create retrieve secrets from Azure and export them to
// Kubernetes secrets.
type KubernetesSecretExporter interface {
	// ExportKubernetesSecrets provides a list of Kubernetes resource for the operator to create once the resource which
	// implements this interface is successfully provisioned. This method is invoked once a resource has been
	// successfully created in Azure, but before the Ready condition has been marked successful.
	ExportKubernetesSecrets(
		ctx context.Context,
		obj MetaObject,
		additionalSecrets set.Set[string], // This exists to avoid making multiple calls to the secrets API - instead we capture all the secrets we need and then get them
		armClient *genericarmclient.GenericClient,
		log logr.Logger,
	) (*KubernetesSecretExportResult, error)
}
