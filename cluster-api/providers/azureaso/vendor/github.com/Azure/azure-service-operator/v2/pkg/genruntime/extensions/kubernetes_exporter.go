/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package extensions

import (
	"context"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/Azure/azure-service-operator/v2/internal/logging"

	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

// This file contains helpers for exporting Kubernetes resources via extensions

type KubernetesExportFunc = func(obj genruntime.MetaObject) ([]client.Object, error)

// CreateKubernetesExporter creates a function to create Kubernetes resources. If the resource
// in question has not been configured with the genruntime.KubernetesExporter interface, the returned function
// is a no-op.
func CreateKubernetesExporter(
	ctx context.Context,
	host genruntime.ResourceExtension,
	armClient *genericarmclient.GenericClient,
	log logr.Logger) KubernetesExportFunc {

	impl, ok := host.(genruntime.KubernetesExporter)
	if !ok {
		return func(obj genruntime.MetaObject) ([]client.Object, error) {
			return nil, nil
		}
	}

	return func(obj genruntime.MetaObject) ([]client.Object, error) {
		log.V(Info).Info("Getting Kubernetes resources for export")
		resources, err := impl.ExportKubernetesResources(ctx, obj, armClient, log)
		if err != nil {
			return resources, err
		}

		log.V(Info).Info(
			"Successfully retrieved Kubernetes resources for export",
			"ResourcesToWrite", len(resources))

		return resources, nil
	}
}
