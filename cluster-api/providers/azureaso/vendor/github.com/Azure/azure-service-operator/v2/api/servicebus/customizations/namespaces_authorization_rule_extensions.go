/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package customizations

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicebus/armservicebus"
	v1 "k8s.io/api/core/v1"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	servicebus "github.com/Azure/azure-service-operator/v2/api/servicebus/v1api20211101/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	. "github.com/Azure/azure-service-operator/v2/internal/logging"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/secrets"
)

// Ensure that NamespacesAuthorizationRuleExtension implements the KubernetesExporter interface
var _ genruntime.KubernetesExporter = &NamespacesAuthorizationRuleExtension{}

// ExportKubernetesResources implements genruntime.KubernetesExporter
func (*NamespacesAuthorizationRuleExtension) ExportKubernetesResources(
	ctx context.Context,
	obj genruntime.MetaObject,
	armClient *genericarmclient.GenericClient,
	log logr.Logger,
) ([]client.Object, error) {
	// Make sure we're working with the current hub version of the resource
	// This will need to be updated if the hub version changes
	rule, ok := obj.(*servicebus.NamespacesAuthorizationRule)
	if !ok {
		return nil, errors.Errorf(
			"cannot run on unknown resource type %T, expected *servicebus.NamespacesAuthorizationRule",
			obj)
	}

	// Type assert that we are the hub type. This will fail to compile if
	// the hub type has been changed but this extension has not
	var _ conversion.Hub = rule

	hasSecrets := authorizationRuleSecretsSpecified(rule)
	if !hasSecrets {
		log.V(Debug).Info("No secrets retrieval to perform as operatorSpec is empty")
		return nil, nil
	}

	id, err := genruntime.GetAndParseResourceID(rule)
	if err != nil {
		return nil, err
	}

	namespaceID := id.Parent
	subscription := id.SubscriptionID

	// Using armClient.ClientOptions() here ensures we share the same HTTP connection, so this is not opening a new
	// connection each time through
	clientFactory, err := armservicebus.NewClientFactory(subscription, armClient.Creds(), armClient.ClientOptions())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create ARM servicebus client factory")
	}

	client := clientFactory.NewNamespacesClient()
	options := armservicebus.NamespacesClientListKeysOptions{}
	response, err := client.ListKeys(ctx, id.ResourceGroupName, namespaceID.Name, rule.Name, &options)
	if err != nil {
		return nil, errors.Wrapf(
			err,
			"failed to retrieve keys for authorization rule %q",
			rule.Name)
	}

	ruleSecrets, err := authorizationRuleSecretsToWrite(rule, response)
	if err != nil {
		return nil, errors.Wrapf(
			err,
			"failed to create secrets for authorization rule %q",
			rule.Name)
	}

	return secrets.SliceToClientObjectSlice(ruleSecrets), nil
}

func authorizationRuleSecretsSpecified(rule *servicebus.NamespacesAuthorizationRule) bool {
	if rule.Spec.OperatorSpec == nil ||
		rule.Spec.OperatorSpec.Secrets == nil {
		return false
	}

	secrets := rule.Spec.OperatorSpec.Secrets

	return secrets.PrimaryKey != nil ||
		secrets.PrimaryConnectionString != nil ||
		secrets.SecondaryKey != nil ||
		secrets.SecondaryConnectionString != nil
}

func authorizationRuleSecretsToWrite(
	rule *servicebus.NamespacesAuthorizationRule,
	response armservicebus.NamespacesClientListKeysResponse,
) ([]*v1.Secret, error) {
	if rule.Spec.OperatorSpec == nil ||
		rule.Spec.OperatorSpec.Secrets == nil {
		return nil, errors.Errorf(
			"authorization rule %q has no secrets specified",
			rule.Name)
	}

	specSecrets := rule.Spec.OperatorSpec.Secrets

	collector := secrets.NewCollector(rule.Namespace)
	collector.AddValue(specSecrets.PrimaryKey, *response.PrimaryKey)
	collector.AddValue(specSecrets.PrimaryConnectionString, *response.PrimaryConnectionString)
	collector.AddValue(specSecrets.SecondaryKey, *response.SecondaryKey)
	collector.AddValue(specSecrets.SecondaryConnectionString, *response.SecondaryConnectionString)

	return collector.Values()
}
