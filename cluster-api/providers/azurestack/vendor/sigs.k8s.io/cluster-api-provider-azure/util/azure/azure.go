/*
Copyright 2022 The Kubernetes Authors.

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

package azure

import (
	"context"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/pkg/errors"
	expv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// AzureSystemNodeLabelPrefix is a standard node label prefix for Azure features, e.g., kubernetes.azure.com/scalesetpriority.
const AzureSystemNodeLabelPrefix = "kubernetes.azure.com"

const (
	// ProviderIDPrefix will be appended to the beginning of Azure resource IDs to form the Kubernetes Provider ID.
	// NOTE: this format matches the 2 slashes format used in cloud-provider and cluster-autoscaler.
	ProviderIDPrefix = "azure://"
)

// IsAzureSystemNodeLabelKey is a helper function that determines whether a node label key is an Azure "system" label.
func IsAzureSystemNodeLabelKey(labelKey string) bool {
	return strings.HasPrefix(labelKey, AzureSystemNodeLabelPrefix)
}

// FindParentMachinePool finds the parent MachinePool for the AzureMachinePool.
func FindParentMachinePool(ampName string, cli client.Client) (*expv1.MachinePool, error) {
	ctx := context.Background()
	machinePoolList := &expv1.MachinePoolList{}
	if err := cli.List(ctx, machinePoolList); err != nil {
		return nil, errors.Wrapf(err, "failed to list MachinePools for %s", ampName)
	}
	for _, mp := range machinePoolList.Items {
		if mp.Spec.Template.Spec.InfrastructureRef.Name == ampName {
			return &mp, nil
		}
	}
	return nil, errors.Errorf("failed to get MachinePool for %s", ampName)
}

// FindParentMachinePoolWithRetry finds the parent MachinePool for the AzureMachinePool with retry.
func FindParentMachinePoolWithRetry(ampName string, cli client.Client, maxAttempts int) (*expv1.MachinePool, error) {
	for i := 1; ; i++ {
		p, err := FindParentMachinePool(ampName, cli)
		if err != nil {
			if i >= maxAttempts {
				return nil, errors.Wrap(err, "failed to find parent MachinePool")
			}
			time.Sleep(1 * time.Second)
			continue
		}
		return p, nil
	}
}

// ParseResourceID parses a string to an *arm.ResourceID, first removing any "azure://" prefix.
func ParseResourceID(id string) (*arm.ResourceID, error) {
	return arm.ParseResourceID(strings.TrimPrefix(id, ProviderIDPrefix))
}
