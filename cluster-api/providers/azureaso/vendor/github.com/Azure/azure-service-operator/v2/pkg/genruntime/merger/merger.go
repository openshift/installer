/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package merger

import (
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/configmaps"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/secrets"
)

func MergeObjects(objs []client.Object) ([]client.Object, error) {
	// We have special handling for secrets and configmaps
	var secretSlice []*v1.Secret
	var configMapSlice []*v1.ConfigMap
	var other []client.Object

	var namespace string
	for _, obj := range objs {
		if namespace == "" {
			namespace = obj.GetNamespace()
		}

		switch typedObj := obj.(type) {
		case *v1.Secret:
			if namespace != obj.GetNamespace() {
				return nil, errors.Errorf("cannot merge objects from different namespaces: %s : %s", namespace, obj.GetNamespace())
			}
			secretSlice = append(secretSlice, typedObj)
		case *v1.ConfigMap:
			if namespace != obj.GetNamespace() {
				return nil, errors.Errorf("cannot merge objects from different namespaces: %s : %s", namespace, obj.GetNamespace())
			}
			configMapSlice = append(configMapSlice, typedObj)
		default:
			other = append(other, typedObj)
		}
	}

	mergedSecrets, err := mergeSecrets(namespace, secretSlice)
	if err != nil {
		return nil, errors.Wrap(err, "failed merging secrets")
	}
	mergedConfigMaps, err := mergeConfigMaps(namespace, configMapSlice)
	if err != nil {
		return nil, errors.Wrap(err, "failed merging config maps")
	}

	result := append(other, secrets.SliceToClientObjectSlice(mergedSecrets)...)
	result = append(result, configmaps.SliceToClientObjectSlice(mergedConfigMaps)...)
	return result, nil
}

func mergeSecrets(namespace string, s []*v1.Secret) ([]*v1.Secret, error) {
	collector := secrets.NewCollector(namespace)
	for _, secret := range s {
		for key, value := range secret.StringData {
			collector.AddValue(
				&genruntime.SecretDestination{
					Name: secret.Name,
					Key:  key,
				}, value)
		}
		for key, value := range secret.Data {
			collector.AddBinaryValue(
				&genruntime.SecretDestination{
					Name: secret.Name,
					Key:  key,
				}, value)
		}
	}

	return collector.Values()
}

func mergeConfigMaps(namespace string, c []*v1.ConfigMap) ([]*v1.ConfigMap, error) {
	collector := configmaps.NewCollector(namespace)
	for _, configMap := range c {
		for key, value := range configMap.Data {
			collector.AddValue(
				&genruntime.ConfigMapDestination{
					Name: configMap.Name,
					Key:  key,
				}, value)
		}
		for key, value := range configMap.BinaryData {
			collector.AddBinaryValue(
				&genruntime.ConfigMapDestination{
					Name: configMap.Name,
					Key:  key,
				}, value)
		}
	}

	return collector.Values()
}
