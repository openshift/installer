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
)

// KubernetesExporter defines a resource which can create other resources in Kubernetes.
type KubernetesExporter interface {
	// ExportKubernetesResources provides a list of Kubernetes resource for the operator to create once the resource which
	// implements this interface is successfully provisioned. This method is invoked once a resource has been
	// successfully created in Azure, but before the Ready condition has been marked successful.
	ExportKubernetesResources(
		ctx context.Context,
		obj MetaObject,
		armClient *genericarmclient.GenericClient,
		log logr.Logger) ([]client.Object, error)
}
