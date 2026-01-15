// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package customizations

import (
	"context"
	"fmt"

	. "github.com/Azure/azure-service-operator/v2/internal/logging"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/go-logr/logr"
	"github.com/rotisserie/eris"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	network "github.com/Azure/azure-service-operator/v2/api/network/v1api20240301/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/resolver"
	"github.com/Azure/azure-service-operator/v2/internal/util/to"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/configmaps"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/extensions"
)

var _ extensions.PostReconciliationChecker = &PrivateEndpointExtension{}

func (extension *PrivateEndpointExtension) PostReconcileCheck(
	_ context.Context,
	obj genruntime.MetaObject,
	_ genruntime.MetaObject,
	_ *resolver.Resolver,
	_ *genericarmclient.GenericClient,
	_ logr.Logger,
	_ extensions.PostReconcileCheckFunc,
) (extensions.PostReconcileCheckResult, error) {
	endpoint, ok := obj.(*network.PrivateEndpoint)
	if !ok {
		return extensions.PostReconcileCheckResult{},
			eris.Errorf("cannot run on unknown resource type %T, expected *network.PrivateEndpoint", obj)
	}

	// Type assert that we are the hub type. This will fail to compile if
	// the hub type has been changed but this extension has not
	var _ conversion.Hub = endpoint

	var reqApprovals []string
	// We want to check `ManualPrivateLinkServiceConnections` as these are the ones which are not auto-approved.
	if connections := endpoint.Status.ManualPrivateLinkServiceConnections; connections != nil {
		for _, connection := range connections {
			if *connection.PrivateLinkServiceConnectionState.Status != "Approved" {
				reqApprovals = append(reqApprovals, *connection.Id)
			}
		}
	}

	if len(reqApprovals) > 0 {
		// Returns 'conditions.NewReadyConditionImpactingError' error
		return extensions.PostReconcileCheckResultFailure(
			fmt.Sprintf(
				"Private connection(s) '%q' to the PrivateEndpoint requires approval",
				reqApprovals)), nil
	}

	return extensions.PostReconcileCheckResultSuccess(), nil
}

var _ genruntime.KubernetesConfigExporter = &PrivateEndpointExtension{}

func (extension *PrivateEndpointExtension) ExportKubernetesConfigMaps(
	ctx context.Context,
	obj genruntime.MetaObject,
	armClient *genericarmclient.GenericClient,
	log logr.Logger,
) ([]client.Object, error) {
	// This has to be the current hub storage version. It will need to be updated
	// if the hub storage version changes.
	endpoint, ok := obj.(*network.PrivateEndpoint)
	if !ok {
		return nil, eris.Errorf("cannot run on unknown resource type %T, expected *network.PrivateEndpoint", obj)
	}

	// Type assert that we are the hub type. This will fail to compile if
	// the hub type has been changed but this extension has not
	var _ conversion.Hub = endpoint

	hasIPConfiguration := hasConfigMaps(endpoint)
	if !hasIPConfiguration {
		log.V(Debug).Info("no configmap retrieval to perform as operatorSpec is empty")
		return nil, nil
	}

	if len(endpoint.Status.NetworkInterfaces) == 0 {
		log.V(Debug).Info("no configmap retrieval to perform as there are no NetworkInterfaces attached")
		return nil, nil
	}

	if endpoint.Status.NetworkInterfaces[0].Id == nil {
		log.V(Debug).Info("no configmap retrieval to perform, failed to fetch the attached NetworkInterfaces")
		return nil, nil
	}

	nicID, err := arm.ParseResourceID(*endpoint.Status.NetworkInterfaces[0].Id)
	if err != nil {
		return nil, err
	}

	// The default primary ip configuration for PrivateEndpoint is on NetworkInterfaceController. Hence, we fetch it from there.
	var interfacesClient *armnetwork.InterfacesClient
	interfacesClient, err = armnetwork.NewInterfacesClient(nicID.SubscriptionID, armClient.Creds(), armClient.ClientOptions())
	if err != nil {
		return nil, eris.Wrapf(err, "failed to create new NetworkInterfacesClient")
	}

	var resp armnetwork.InterfacesClientGetResponse
	resp, err = interfacesClient.Get(ctx, nicID.ResourceGroupName, nicID.Name, nil)
	if err != nil {
		return nil, eris.Wrapf(err, "failed getting NetworkInterfaceController")
	}

	configsByName := configByName(log, resp.Interface)
	configs, err := configMapToWrite(endpoint, configsByName)
	if err != nil {
		return nil, err
	}

	return configmaps.SliceToClientObjectSlice(configs), nil
}

func configByName(log logr.Logger, nic armnetwork.Interface) map[string]string {
	result := make(map[string]string)

	if nic.Properties != nil && nic.Properties.IPConfigurations != nil {
		for _, ipConfiguration := range nic.Properties.IPConfigurations {
			if ipConfiguration.Properties == nil || ipConfiguration.Properties.PrivateIPAddress == nil {
				log.V(Debug).Info("skipping IPConfiguration properties nil for IPConfiguration")
				continue
			}

			if !to.Value(ipConfiguration.Properties.Primary) {
				// This ipConfiguration is not primary
				continue
			}

			result["primaryNICPrivateIPAddress"] = *nic.Properties.IPConfigurations[0].Properties.PrivateIPAddress
			break
		}
	}

	return result
}

func configMapToWrite(obj *network.PrivateEndpoint, configs map[string]string) ([]*v1.ConfigMap, error) {
	operatorSpecConfigs := obj.Spec.OperatorSpec.ConfigMaps
	if operatorSpecConfigs == nil {
		return nil, eris.Errorf("unexpected nil operatorspec")
	}

	collector := configmaps.NewCollector(obj.Namespace)

	primaryNICPrivateIPAddress, ok := configs["primaryNICPrivateIPAddress"]
	if ok {
		collector.AddValue(operatorSpecConfigs.PrimaryNicPrivateIpAddress, primaryNICPrivateIPAddress)
	}

	return collector.Values()
}

func hasConfigMaps(endpoint *network.PrivateEndpoint) bool {
	if endpoint.Spec.OperatorSpec == nil || endpoint.Spec.OperatorSpec.ConfigMaps == nil {
		return false
	}

	hasIPConfiguration := false
	configMaps := endpoint.Spec.OperatorSpec.ConfigMaps

	if configMaps != nil && configMaps.PrimaryNicPrivateIpAddress != nil {
		hasIPConfiguration = true
	}

	return hasIPConfiguration
}
