/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package extensions

import (
	"context"

	. "github.com/Azure/azure-service-operator/v2/internal/logging"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/set"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

// This file contains helpers for exporting Kubernetes resources via extensions

type KubernetesSecretExportFunc = func(obj genruntime.MetaObject, additionalSecrets set.Set[string]) (*genruntime.KubernetesSecretExportResult, error)

// CreateKubernetesSecretExporter creates a function to create Kubernetes secrets. If the resource
// in question has not been configured with the genruntime.KubernetesSecretExporter interface, the returned function
// is a no-op.
func CreateKubernetesSecretExporter(
	ctx context.Context,
	host genruntime.ResourceExtension,
	armClient *genericarmclient.GenericClient,
	log logr.Logger,
) KubernetesSecretExportFunc {
	impl, ok := host.(genruntime.KubernetesSecretExporter)
	if !ok {
		return func(obj genruntime.MetaObject, additionalSecrets set.Set[string]) (*genruntime.KubernetesSecretExportResult, error) {
			return nil, nil
		}
	}

	return func(obj genruntime.MetaObject, additionalSecrets set.Set[string]) (*genruntime.KubernetesSecretExportResult, error) {
		log.V(Info).Info("Getting Kubernetes secrets for export")
		result, err := impl.ExportKubernetesSecrets(ctx, obj, additionalSecrets, armClient, log)
		if err != nil {
			return result, err
		}

		var objs []client.Object
		var rawSecrets map[string]string
		if result != nil {
			objs = result.Objs
			rawSecrets = result.RawSecrets
		}

		log.V(Info).Info(
			"Successfully retrieved Kubernetes secrets for export",
			"ResourcesToWrite", len(objs), "RawSecrets", len(rawSecrets))

		return result, nil
	}
}
