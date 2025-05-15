/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package customizations

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicebus/armservicebus"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	servicebus "github.com/Azure/azure-service-operator/v2/api/servicebus/v1api20211101/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	. "github.com/Azure/azure-service-operator/v2/internal/logging"
	"github.com/Azure/azure-service-operator/v2/internal/set"
	"github.com/Azure/azure-service-operator/v2/internal/util/to"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/secrets"
)

var _ genruntime.KubernetesSecretExporter = &NamespacesAuthorizationRuleExtension{}

func (ext *NamespacesAuthorizationRuleExtension) ExportKubernetesSecrets(
	ctx context.Context,
	obj genruntime.MetaObject,
	additionalSecrets set.Set[string],
	armClient *genericarmclient.GenericClient,
	log logr.Logger,
) (*genruntime.KubernetesSecretExportResult, error) {
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

	primarySecrets := authorizationRuleSecretsSpecified(rule)
	requestedSecrets := set.Union(primarySecrets, additionalSecrets)
	if len(requestedSecrets) == 0 {
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

	resolvedSecrets := map[string]string{}
	if to.Value(response.PrimaryKey) != "" {
		resolvedSecrets[primaryKey] = to.Value(response.PrimaryKey)
	}
	if to.Value(response.PrimaryConnectionString) != "" {
		resolvedSecrets[primaryConnectionString] = to.Value(response.PrimaryConnectionString)
	}
	if to.Value(response.SecondaryKey) != "" {
		resolvedSecrets[secondaryKey] = to.Value(response.SecondaryKey)
	}
	if to.Value(response.SecondaryConnectionString) != "" {
		resolvedSecrets[secondaryConnectionString] = to.Value(response.SecondaryConnectionString)
	}

	return &genruntime.KubernetesSecretExportResult{
		Objs:       secrets.SliceToClientObjectSlice(ruleSecrets),
		RawSecrets: secrets.SelectSecrets(additionalSecrets, resolvedSecrets),
	}, nil
}

func authorizationRuleSecretsSpecified(rule *servicebus.NamespacesAuthorizationRule) set.Set[string] {
	if rule.Spec.OperatorSpec == nil ||
		rule.Spec.OperatorSpec.Secrets == nil {
		return nil
	}

	secrets := rule.Spec.OperatorSpec.Secrets

	result := make(set.Set[string])
	if secrets.PrimaryKey != nil {
		result.Add(primaryKey)
	}
	if secrets.PrimaryConnectionString != nil {
		result.Add(primaryConnectionString)
	}
	if secrets.SecondaryKey != nil {
		result.Add(secondaryKey)
	}
	if secrets.SecondaryConnectionString != nil {
		result.Add(secondaryConnectionString)
	}

	return result
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
