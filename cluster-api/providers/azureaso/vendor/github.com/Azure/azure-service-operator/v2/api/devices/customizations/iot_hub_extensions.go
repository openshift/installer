/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package customizations

import (
	"context"

	. "github.com/Azure/azure-service-operator/v2/internal/logging"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/iothub/armiothub"
	"github.com/go-logr/logr"
	"github.com/rotisserie/eris"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	devices "github.com/Azure/azure-service-operator/v2/api/devices/v1api20210702/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/set"
	"github.com/Azure/azure-service-operator/v2/internal/util/to"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/secrets"
)

const (
	iotHubOwnerPrimaryKey         = "iotHubOwnerPrimaryKey"
	iotHubOwnerSecondaryKey       = "iotHubOwnerSecondaryKey"
	servicePrimaryKey             = "servicePrimaryKey"
	serviceSecondaryKey           = "serviceSecondaryKey"
	registryReadWritePrimaryKey   = "registryReadWritePrimaryKey"
	registryReadWriteSecondaryKey = "registryReadWriteSecondaryKey"
	registryReadPrimaryKey        = "registryReadPrimaryKey"
	registryReadSecondaryKey      = "registryReadSecondaryKey"
	devicePrimaryKey              = "devicePrimaryKey"
	deviceSecondaryKey            = "deviceSecondaryKey"
)

var _ genruntime.KubernetesSecretExporter = &IotHubExtension{}

func (ext *IotHubExtension) ExportKubernetesSecrets(
	ctx context.Context,
	obj genruntime.MetaObject,
	additionalSecrets set.Set[string],
	armClient *genericarmclient.GenericClient,
	log logr.Logger,
) (*genruntime.KubernetesSecretExportResult, error) {
	// This has to be the current hub devices version. It will need to be updated
	// if the hub devices version changes.
	typedObj, ok := obj.(*devices.IotHub)
	if !ok {
		return nil, eris.Errorf("cannot run on unknown resource type %T, expected *devices.IotHub", obj)
	}

	// Type assert that we are the hub type. This will fail to compile if
	// the hub type has been changed but this extension has not
	var _ conversion.Hub = typedObj

	primarySecrets := secretsSpecified(typedObj)
	requestedSecrets := set.Union(primarySecrets, additionalSecrets)
	if len(requestedSecrets) == 0 {
		log.V(Debug).Info("No secrets retrieval to perform as operatorSpec is empty")
		return nil, nil
	}

	id, err := genruntime.GetAndParseResourceID(typedObj)
	if err != nil {
		return nil, err
	}

	keys := make(map[string]armiothub.SharedAccessSignatureAuthorizationRule)
	// Only bother calling ListKeys if there are secrets to retrieve
	if len(requestedSecrets) > 0 {
		subscription := id.SubscriptionID
		// Using armClient.ClientOptions() here ensures we share the same HTTP connection, so this is not opening a new
		// connection each time through
		var resClient *armiothub.ResourceClient
		resClient, err = armiothub.NewResourceClient(subscription, armClient.Creds(), armClient.ClientOptions())
		if err != nil {
			return nil, eris.Wrapf(err, "failed to create new DevicesClient")
		}

		var pager *runtime.Pager[armiothub.ResourceClientListKeysResponse]
		var resp armiothub.ResourceClientListKeysResponse
		pager = resClient.NewListKeysPager(id.ResourceGroupName, typedObj.AzureName(), nil)
		for pager.More() {
			resp, err = pager.NextPage(ctx)
			if err != nil {
				return nil, eris.Wrapf(err, "failed listing keys")
			}
			addSecretsToMap(resp.Value, keys)
		}
	}

	secretSlice, err := secretsToWrite(typedObj, keys)
	if err != nil {
		return nil, err
	}

	resolvedSecrets := makeResolvedSecretsMap(keys)

	return &genruntime.KubernetesSecretExportResult{
		Objs:       secrets.SliceToClientObjectSlice(secretSlice),
		RawSecrets: secrets.SelectSecrets(additionalSecrets, resolvedSecrets),
	}, nil
}

func secretsSpecified(obj *devices.IotHub) set.Set[string] {
	if obj.Spec.OperatorSpec == nil || obj.Spec.OperatorSpec.Secrets == nil {
		return nil
	}

	secrets := obj.Spec.OperatorSpec.Secrets

	result := make(set.Set[string])
	if secrets.IotHubOwnerPrimaryKey != nil {
		result.Add(iotHubOwnerPrimaryKey)
	}
	if secrets.IotHubOwnerSecondaryKey != nil {
		result.Add(iotHubOwnerSecondaryKey)
	}
	if secrets.ServicePrimaryKey != nil {
		result.Add(servicePrimaryKey)
	}
	if secrets.ServiceSecondaryKey != nil {
		result.Add(serviceSecondaryKey)
	}
	if secrets.RegistryReadWritePrimaryKey != nil {
		result.Add(registryReadWritePrimaryKey)
	}
	if secrets.RegistryReadWriteSecondaryKey != nil {
		result.Add(registryReadWriteSecondaryKey)
	}
	if secrets.RegistryReadPrimaryKey != nil {
		result.Add(registryReadPrimaryKey)
	}
	if secrets.RegistryReadSecondaryKey != nil {
		result.Add(registryReadSecondaryKey)
	}
	if secrets.DevicePrimaryKey != nil {
		result.Add(devicePrimaryKey)
	}
	if secrets.DeviceSecondaryKey != nil {
		result.Add(deviceSecondaryKey)
	}

	return result
}

func addSecretsToMap(keys []*armiothub.SharedAccessSignatureAuthorizationRule, result map[string]armiothub.SharedAccessSignatureAuthorizationRule) {
	for _, key := range keys {
		if key == nil || key.KeyName == nil {
			continue
		}
		result[*key.KeyName] = *key
	}
}

func secretsToWrite(obj *devices.IotHub, keys map[string]armiothub.SharedAccessSignatureAuthorizationRule) ([]*v1.Secret, error) {
	operatorSpecSecrets := obj.Spec.OperatorSpec.Secrets
	if operatorSpecSecrets == nil {
		return nil, nil
	}

	// Documentation for keys : https://learn.microsoft.com/en-us/rest/api/iothub/iot-hub-resource/list-keys?tabs=HTTP#sharedaccesssignatureauthorizationrule
	collector := secrets.NewCollector(obj.Namespace)
	iothubOwner, ok := keys["iothubowner"]
	if ok {
		collector.AddValue(operatorSpecSecrets.IotHubOwnerPrimaryKey, to.Value(iothubOwner.PrimaryKey))
		collector.AddValue(operatorSpecSecrets.IotHubOwnerSecondaryKey, to.Value(iothubOwner.SecondaryKey))
	}

	service, ok := keys["service"]
	if ok {
		collector.AddValue(operatorSpecSecrets.ServicePrimaryKey, to.Value(service.PrimaryKey))
		collector.AddValue(operatorSpecSecrets.ServiceSecondaryKey, to.Value(service.SecondaryKey))
	}

	device, ok := keys["device"]
	if ok {
		collector.AddValue(operatorSpecSecrets.DevicePrimaryKey, to.Value(device.PrimaryKey))
		collector.AddValue(operatorSpecSecrets.DeviceSecondaryKey, to.Value(device.SecondaryKey))
	}

	registryRead, ok := keys["registryRead"]
	if ok {
		collector.AddValue(operatorSpecSecrets.RegistryReadPrimaryKey, to.Value(registryRead.PrimaryKey))
		collector.AddValue(operatorSpecSecrets.RegistryReadSecondaryKey, to.Value(registryRead.SecondaryKey))
	}

	registryReadWrite, ok := keys["registryReadWrite"]
	if ok {
		collector.AddValue(operatorSpecSecrets.RegistryReadWritePrimaryKey, to.Value(registryReadWrite.PrimaryKey))
		collector.AddValue(operatorSpecSecrets.RegistryReadWriteSecondaryKey, to.Value(registryReadWrite.SecondaryKey))
	}

	return collector.Values()
}

func makeResolvedSecretsMap(keys map[string]armiothub.SharedAccessSignatureAuthorizationRule) map[string]string {
	result := make(map[string]string)

	iothubOwner, ok := keys["iothubowner"]
	if ok {
		if to.Value(iothubOwner.PrimaryKey) != "" {
			result[iotHubOwnerPrimaryKey] = to.Value(iothubOwner.PrimaryKey)
		}
		if to.Value(iothubOwner.SecondaryKey) != "" {
			result[iotHubOwnerSecondaryKey] = to.Value(iothubOwner.SecondaryKey)
		}
	}

	service, ok := keys["service"]
	if ok {
		if to.Value(service.PrimaryKey) != "" {
			result[servicePrimaryKey] = to.Value(service.PrimaryKey)
		}
		if to.Value(service.SecondaryKey) != "" {
			result[serviceSecondaryKey] = to.Value(service.SecondaryKey)
		}
	}

	device, ok := keys["device"]
	if ok {
		if to.Value(device.PrimaryKey) != "" {
			result[devicePrimaryKey] = to.Value(device.PrimaryKey)
		}
		if to.Value(device.SecondaryKey) != "" {
			result[deviceSecondaryKey] = to.Value(device.SecondaryKey)
		}
	}

	registryRead, ok := keys["registryRead"]
	if ok {
		if to.Value(registryRead.PrimaryKey) != "" {
			result[registryReadPrimaryKey] = to.Value(registryRead.PrimaryKey)
		}
		if to.Value(registryRead.SecondaryKey) != "" {
			result[registryReadSecondaryKey] = to.Value(registryRead.SecondaryKey)
		}
	}

	registryReadWrite, ok := keys["registryReadWrite"]
	if ok {
		if to.Value(registryReadWrite.PrimaryKey) != "" {
			result[registryReadWritePrimaryKey] = to.Value(registryReadWrite.PrimaryKey)
		}
		if to.Value(registryReadWrite.SecondaryKey) != "" {
			result[registryReadWriteSecondaryKey] = to.Value(registryReadWrite.SecondaryKey)
		}
	}

	return result
}
