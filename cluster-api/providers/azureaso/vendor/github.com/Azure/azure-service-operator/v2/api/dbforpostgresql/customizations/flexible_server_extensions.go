/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package customizations

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	postgresql "github.com/Azure/azure-service-operator/v2/api/dbforpostgresql/v1api20221201/storage"
	. "github.com/Azure/azure-service-operator/v2/internal/logging"

	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/resolver"
	"github.com/Azure/azure-service-operator/v2/internal/set"
	"github.com/Azure/azure-service-operator/v2/internal/util/to"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/extensions"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/secrets"
)

var _ genruntime.KubernetesExporter = &FlexibleServerExtension{}

func (ext *FlexibleServerExtension) ExportKubernetesResources(
	_ context.Context,
	obj genruntime.MetaObject,
	_ *genericarmclient.GenericClient,
	log logr.Logger) ([]client.Object, error) {

	// This has to be the current hub storage version. It will need to be updated
	// if the hub storage version changes.
	typedObj, ok := obj.(*postgresql.FlexibleServer)
	if !ok {
		return nil, errors.Errorf("cannot run on unknown resource type %T, expected *postgresql.FlexibleServer", obj)
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

func secretsSpecified(obj *postgresql.FlexibleServer) bool {
	if obj.Spec.OperatorSpec == nil || obj.Spec.OperatorSpec.Secrets == nil {
		return false
	}

	operatorSecrets := obj.Spec.OperatorSpec.Secrets
	if operatorSecrets.FullyQualifiedDomainName != nil {
		return true
	}

	return false
}

func secretsToWrite(obj *postgresql.FlexibleServer) ([]*v1.Secret, error) {
	operatorSpecSecrets := obj.Spec.OperatorSpec.Secrets
	if operatorSpecSecrets == nil {
		return nil, errors.Errorf("unexpected nil operatorspec")
	}

	collector := secrets.NewCollector(obj.Namespace)
	collector.AddValue(operatorSpecSecrets.FullyQualifiedDomainName, to.Value(obj.Status.FullyQualifiedDomainName))

	return collector.Values()
}

var _ extensions.PreReconciliationChecker = &FlexibleServerExtension{}

// If the provisioningState of a flexible server is not in this set, it will reject any attempt to PUT the resource
// out of hand; so there's no point in even trying. This is true even if the PUT we're doing will have no effect on
// the state of the server.
// These are all listed lowercase, so we can do a case-insensitive match.
var nonBlockingFlexibleServerStates = set.Make(
	"succeeded",
	"failed",
	"canceled",
	"ready",
)

func (ext *FlexibleServerExtension) PreReconcileCheck(
	_ context.Context,
	obj genruntime.MetaObject,
	_ genruntime.MetaObject,
	_ *resolver.Resolver,
	_ *genericarmclient.GenericClient,
	_ logr.Logger,
	_ extensions.PreReconcileCheckFunc,
) (extensions.PreReconcileCheckResult, error) {
	// This has to be the current hub storage version. It will need to be updated
	// if the hub storage version changes.
	server, ok := obj.(*postgresql.FlexibleServer)
	if !ok {
		return extensions.PreReconcileCheckResult{},
			errors.Errorf("cannot run on unknown resource type %T, expected *postgresql.FlexibleServer", obj)
	}

	// Type assert that we are the hub type. This will fail to compile if
	// the hub type has been changed but this extension has not been updated
	var _ conversion.Hub = server

	state := server.Status.State
	if state != nil && flexibleServerStateBlocksReconciliation(*state) {
		return extensions.BlockReconcile(
			fmt.Sprintf(
				"Flexible Server is in provisioning state %q",
				*state)), nil
	}

	return extensions.ProceedWithReconcile(), nil
}

func flexibleServerStateBlocksReconciliation(state string) bool {
	return !nonBlockingFlexibleServerStates.Contains(strings.ToLower(state))
}
