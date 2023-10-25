/*
Copyright 2019 The Kubernetes Authors.

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

package util

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	kcfg "sigs.k8s.io/cluster-api/util/kubeconfig"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewKubeClient returns a new client for the target cluster using the KubeConfig
// secret stored in the management cluster.
func NewKubeClient(
	ctx context.Context,
	controllerClient client.Client,
	cluster *clusterv1.Cluster) (kubernetes.Interface, error) {
	clusterKey := client.ObjectKey{Namespace: cluster.Namespace, Name: cluster.Name}
	kubeconfig, err := kcfg.FromSecret(ctx, controllerClient, clusterKey)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to retrieve kubeconfig secret for Cluster %q in namespace %q",
			cluster.Name, cluster.Namespace)
	}

	restConfig, err := clientcmd.RESTConfigFromKubeConfig(kubeconfig)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create client configuration for Cluster %q in namespace %q",
			cluster.Name, cluster.Namespace)
	}
	// sets the timeout, otherwise this will default to 0 (i.e. no timeout) which might cause tests to hang
	restConfig.Timeout = 10 * time.Second

	return kubernetes.NewForConfig(restConfig)
}
