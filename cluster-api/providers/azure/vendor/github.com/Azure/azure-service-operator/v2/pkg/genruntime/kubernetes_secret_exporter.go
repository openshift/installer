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

// KubernetesSecretExporter defines a resource which can retrieve secrets from Azure and export them to
// Kubernetes secrets. This extension is invoked after a resource has been successfully created or updated
// in Azure, giving resources the ability to make sensitive data available in Kubernetes.
// Implement this extension when:
// - The resource generates secrets in Azure (keys, connection strings, passwords)
// - Secret values are available in the resource status
// - Additional ARM API calls are needed to retrieve secrets
type KubernetesSecretExporter interface {
	// ExportKubernetesSecrets retrieves secrets from Azure and returns Kubernetes Secret objects to be created.
	// This method is invoked once a resource has been successfully created or updated in Azure,
	// but before the Ready condition has been marked successful.
	// ctx is the current operation context.
	// obj is the resource that owns the secrets.
	// additionalSecrets is a set of secret names to retrieve (for secret expressions). This exists to avoid
	// making multiple calls to the secrets API - instead we capture all the secrets we need and then get them.
	// armClient allows making ARM API calls to retrieve secrets.
	// log is a logger for the current operation.
	// Returns a KubernetesSecretExportResult containing Secret objects to create and raw secret values.
	ExportKubernetesSecrets(
		ctx context.Context,
		obj MetaObject,
		additionalSecrets set.Set[string], // This exists to avoid making multiple calls to the secrets API - instead we capture all the secrets we need and then get them
		armClient *genericarmclient.GenericClient,
		log logr.Logger,
	) (*KubernetesSecretExportResult, error)
}
