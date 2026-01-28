// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package customizations

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/rotisserie/eris"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	network "github.com/Azure/azure-service-operator/v2/api/network/v1api20240301/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/resolver"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/extensions"
)

var _ extensions.PostReconciliationChecker = &PrivateEndpointExtension{}

func (extension *VirtualNetworksSubnetExtension) PostReconcileCheck(
	_ context.Context,
	obj genruntime.MetaObject,
	_ genruntime.MetaObject,
	_ *resolver.Resolver,
	_ *genericarmclient.GenericClient,
	_ logr.Logger,
	_ extensions.PostReconcileCheckFunc,
) (extensions.PostReconcileCheckResult, error) {
	subnet, ok := obj.(*network.VirtualNetworksSubnet)
	if !ok {
		return extensions.PostReconcileCheckResult{},
			eris.Errorf("cannot run on unknown resource type %T, expected *network.VirtualNetworksSubnet", obj)
	}

	// Type assert that we are the hub type. This will fail to compile if
	// the hub type has been changed but this extension has not
	var _ conversion.Hub = subnet

	// Subnets can have a HUGE number of ipConfigurations in some modes. So many that it can cause Kubernetes to be unable
	// to fit the resource. We have to omit them after some point to avoid blowing out the resource size and causing
	// kube-apiserver to reject us. See https://github.com/Azure/azure-service-operator/issues/4428.

	// This limit was chosen based on a 300 character long IPConfiguration ID,
	// which would be 300 bytes in UTF-8. 2000*300 = ~.6mb, which is about around 1/3rd the max allowed size of a
	// Kubernetes object.
	if len(subnet.Status.IpConfigurations) > 2000 {
		subnet.Status.IpConfigurations = nil
	}

	return extensions.PostReconcileCheckResultSuccess(), nil
}
