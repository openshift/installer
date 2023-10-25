/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package scope

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"

	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"

	"k8s.io/klog/v2/klogr"

	"sigs.k8s.io/controller-runtime/pkg/client"

	capiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/patch"

	infrav1beta2 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/powervs"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/resourcecontroller"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/endpoints"
)

const (
	// DEBUGLEVEL indicates the debug level of the logs.
	DEBUGLEVEL = 5
)

// PowerVSClusterScopeParams defines the input parameters used to create a new PowerVSClusterScope.
type PowerVSClusterScopeParams struct {
	Client            client.Client
	Logger            logr.Logger
	Cluster           *capiv1beta1.Cluster
	IBMPowerVSCluster *infrav1beta2.IBMPowerVSCluster
	ServiceEndpoint   []endpoints.ServiceEndpoint
}

// PowerVSClusterScope defines a scope defined around a Power VS Cluster.
type PowerVSClusterScope struct {
	logr.Logger
	Client      client.Client
	patchHelper *patch.Helper

	IBMPowerVSClient  powervs.PowerVS
	Cluster           *capiv1beta1.Cluster
	IBMPowerVSCluster *infrav1beta2.IBMPowerVSCluster
	ServiceEndpoint   []endpoints.ServiceEndpoint
}

// NewPowerVSClusterScope creates a new PowerVSClusterScope from the supplied parameters.
func NewPowerVSClusterScope(params PowerVSClusterScopeParams) (scope *PowerVSClusterScope, err error) {
	scope = &PowerVSClusterScope{}

	if params.Client == nil {
		err = errors.New("failed to generate new scope from nil Client")
		return nil, err
	}
	scope.Client = params.Client

	if params.Cluster == nil {
		err = errors.New("failed to generate new scope from nil Cluster")
		return nil, err
	}
	scope.Cluster = params.Cluster

	if params.IBMPowerVSCluster == nil {
		err = errors.New("failed to generate new scope from nil IBMPowerVSCluster")
		return nil, err
	}
	scope.IBMPowerVSCluster = params.IBMPowerVSCluster

	if params.Logger == (logr.Logger{}) {
		params.Logger = klogr.New()
	}
	scope.Logger = params.Logger

	helper, err := patch.NewHelper(params.IBMPowerVSCluster, params.Client)
	if err != nil {
		err = errors.Wrap(err, "failed to init patch helper")
		return nil, err
	}
	scope.patchHelper = helper

	spec := params.IBMPowerVSCluster.Spec

	rc, err := resourcecontroller.NewService(resourcecontroller.ServiceOptions{})
	if err != nil {
		return nil, err
	}

	// Fetch the resource controller endpoint.
	if rcEndpoint := endpoints.FetchRCEndpoint(params.ServiceEndpoint); rcEndpoint != "" {
		if err := rc.SetServiceURL(rcEndpoint); err != nil {
			return nil, errors.Wrap(err, "failed to set resource controller endpoint")
		}
		scope.Logger.V(3).Info("Overriding the default resource controller endpoint")
	}

	res, _, err := rc.GetResourceInstance(
		&resourcecontrollerv2.GetResourceInstanceOptions{
			ID: core.StringPtr(spec.ServiceInstanceID),
		})
	if err != nil {
		err = errors.Wrap(err, "failed to get resource instance")
		return nil, err
	}

	options := powervs.ServiceOptions{
		IBMPIOptions: &ibmpisession.IBMPIOptions{
			Debug: params.Logger.V(DEBUGLEVEL).Enabled(),
			Zone:  *res.RegionID,
		},
		CloudInstanceID: spec.ServiceInstanceID,
	}

	// Fetch the service endpoint.
	if svcEndpoint := endpoints.FetchPVSEndpoint(endpoints.CostructRegionFromZone(*res.RegionID), params.ServiceEndpoint); svcEndpoint != "" {
		options.IBMPIOptions.URL = svcEndpoint
		scope.Logger.V(3).Info("Overriding the default powervs service endpoint")
	}

	c, err := powervs.NewService(options)
	if err != nil {
		err = fmt.Errorf("failed to create NewIBMPowerVSClient")
		return nil, err
	}
	scope.IBMPowerVSClient = c

	return scope, nil
}

// PatchObject persists the cluster configuration and status.
func (s *PowerVSClusterScope) PatchObject() error {
	return s.patchHelper.Patch(context.TODO(), s.IBMPowerVSCluster)
}

// Close closes the current scope persisting the cluster configuration and status.
func (s *PowerVSClusterScope) Close() error {
	return s.PatchObject()
}
