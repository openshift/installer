/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package customizations

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	mariadb "github.com/Azure/azure-service-operator/v2/api/dbformariadb/v1api20180601/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	. "github.com/Azure/azure-service-operator/v2/internal/logging"
	"github.com/Azure/azure-service-operator/v2/internal/util/to"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/secrets"
)

var _ genruntime.KubernetesExporter = &ServerExtension{}

func (ext *ServerExtension) ExportKubernetesResources(
	ctx context.Context,
	obj genruntime.MetaObject,
	armClient *genericarmclient.GenericClient,
	log logr.Logger) ([]client.Object, error) {

	// This has to be the current storage version. It will need to be updated
	// if the storage version changes.
	typedObj, ok := obj.(*mariadb.Server)
	if !ok {
		return nil, errors.Errorf("cannot run on unknown resource type %T", obj)
	}

	// Type assert that we are the hub type. This will fail to compile if
	// the hub type has been changed but this extension has not
	var _ conversion.Hub = typedObj

	hasSecrets := secretsSpecified(typedObj)
	if !hasSecrets {
		log.V(Debug).Info("No secrets retrieval to perform as operatorSpec is empty")
		return nil, nil
	}

	secretSlice, err := secretsToWrite(typedObj)
	if err != nil {
		return nil, err
	}

	return secrets.SliceToClientObjectSlice(secretSlice), nil
}

func secretsSpecified(obj *mariadb.Server) bool {
	if obj.Spec.OperatorSpec == nil || obj.Spec.OperatorSpec.Secrets == nil {
		return false
	}

	secrets := obj.Spec.OperatorSpec.Secrets
	if secrets.FullyQualifiedDomainName != nil {
		return true
	}

	return false
}

func secretsToWrite(obj *mariadb.Server) ([]*v1.Secret, error) {
	operatorSpecSecrets := obj.Spec.OperatorSpec.Secrets
	if operatorSpecSecrets == nil {
		return nil, errors.Errorf("unexpected nil OperatorSpec")
	}

	collector := secrets.NewCollector(obj.Namespace)
	collector.AddValue(operatorSpecSecrets.FullyQualifiedDomainName, to.Value(obj.Status.FullyQualifiedDomainName))

	return collector.Values()
}
