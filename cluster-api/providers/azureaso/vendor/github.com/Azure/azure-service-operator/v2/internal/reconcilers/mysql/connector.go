/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package mysql

import (
	"context"

	"github.com/pkg/errors"
	ctrlconversion "sigs.k8s.io/controller-runtime/pkg/conversion"

	asomysql "github.com/Azure/azure-service-operator/v2/api/dbformysql/v1"
	dbformysql "github.com/Azure/azure-service-operator/v2/api/dbformysql/v1api20230630/storage"
	"github.com/Azure/azure-service-operator/v2/internal/resolver"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/conditions"
)

type Connector interface {
	CreateOrUpdate(ctx context.Context) error
	Delete(ctx context.Context) error
	Exists(ctx context.Context) (bool, error)
}

func getServerFQDN(ctx context.Context, resourceResolver *resolver.Resolver, user *asomysql.User) (string, error) {
	// Get the owner - at this point it must exist
	ownerDetails, err := resourceResolver.ResolveOwner(ctx, user)
	if err != nil {
		return "", err
	}

	// Note that this is not actually possible for this type because we don't allow ARMID references for these owners,
	// but protecting against it here anyway.
	if !ownerDetails.FoundKubernetesOwner() {
		return "", errors.Errorf("user owner must exist in Kubernetes for user %s", user.Name)
	}

	flexibleServer, ok := ownerDetails.Owner.(*dbformysql.FlexibleServer)
	if !ok {
		return "", errors.Errorf("owner was not type FlexibleServer, instead: %T", ownerDetails)
	}

	// Assertion to ensure that this is still the storage type
	// If this doesn't compile, update the version being imported to the new Hub version
	var _ ctrlconversion.Hub = &dbformysql.FlexibleServer{}

	if flexibleServer.Status.FullyQualifiedDomainName == nil {
		// This possibly means that the server hasn't finished deploying yet
		err = errors.Errorf("owning Flexibleserver %q '.status.fullyQualifiedDomainName' not set. Has the server been provisioned successfully?", flexibleServer.Name)
		return "", conditions.NewReadyConditionImpactingError(err, conditions.ConditionSeverityWarning, conditions.ReasonWaitingForOwner)
	}
	serverFQDN := *flexibleServer.Status.FullyQualifiedDomainName

	return serverFQDN, nil
}
