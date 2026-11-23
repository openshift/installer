/*
Copyright The Kubernetes Authors.

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

package v1beta1

import (
	"context"
	"time"

	"github.com/pkg/errors"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
)

// GetOwnerAzureClusterNameAndNamespace returns the owner azure cluster's name and namespace for the given cluster name and namespace.
func GetOwnerAzureClusterNameAndNamespace(cli client.Client, clusterName string, namespace string, maxAttempts int) (azureClusterName string, azureClusterNamespace string, err error) {
	ctx := context.Background()

	ownerCluster := &clusterv1.Cluster{}
	key := client.ObjectKey{
		Namespace: namespace,
		Name:      clusterName,
	}

	for i := 1; ; i++ {
		if err := cli.Get(ctx, key, ownerCluster); err != nil {
			if i > maxAttempts {
				return "", "", errors.Wrapf(err, "failed to find owner cluster for %s/%s", namespace, clusterName)
			}
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	return ownerCluster.Spec.InfrastructureRef.Name, ownerCluster.Namespace, nil
}

// GetSubscriptionID returns the subscription ID for the AzureCluster given the cluster name and namespace.
func GetSubscriptionID(cli client.Client, ownerAzureClusterName string, ownerAzureClusterNamespace string, maxAttempts int) (string, error) {
	ctx := context.Background()

	ownerAzureCluster := &infrav1.AzureCluster{}
	key := client.ObjectKey{
		Namespace: ownerAzureClusterNamespace,
		Name:      ownerAzureClusterName,
	}
	for i := 1; ; i++ {
		if err := cli.Get(ctx, key, ownerAzureCluster); err != nil {
			if i >= maxAttempts {
				return "", errors.Wrapf(err, "failed to find AzureCluster for owner cluster %s/%s", ownerAzureClusterNamespace, ownerAzureClusterName)
			}
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	return ownerAzureCluster.Spec.SubscriptionID, nil
}
