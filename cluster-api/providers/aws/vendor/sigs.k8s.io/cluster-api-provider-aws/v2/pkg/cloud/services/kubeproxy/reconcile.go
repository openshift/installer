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

package kubeproxy

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
)

const (
	kubeProxyName      = "kube-proxy"
	kubeProxyNamespace = "kube-system"
)

// ReconcileKubeProxy will reconcile kube-proxy.
func (s *Service) ReconcileKubeProxy(ctx context.Context) error {
	s.scope.Info("Reconciling kube-proxy DaemonSet in cluster", "cluster", klog.KRef(s.scope.Namespace(), s.scope.Name()))

	remoteClient, err := s.scope.RemoteClient()
	if err != nil {
		s.scope.Error(err, "getting client for remote cluster")
		return fmt.Errorf("getting client for remote cluster: %w", err)
	}

	if s.scope.DisableKubeProxy() {
		if err := s.deleteKubeProxy(ctx, remoteClient); err != nil {
			return fmt.Errorf("disabling kube-proxy: %w", err)
		}
	}

	return nil
}

func (s *Service) deleteKubeProxy(ctx context.Context, remoteClient client.Client) error {
	s.scope.Info("Ensuring the kube-proxy DaemonSet in cluster is deleted", "cluster", klog.KRef(s.scope.Namespace(), s.scope.Name()))

	ds := &appsv1.DaemonSet{}
	if err := remoteClient.Get(ctx, types.NamespacedName{Namespace: kubeProxyNamespace, Name: kubeProxyName}, ds); err != nil {
		if apierrors.IsNotFound(err) {
			s.scope.Debug("The kube-proxy DaemonSet is not found, no action")
			return nil
		}
		return fmt.Errorf("getting kube-proxy daemonset: %w", err)
	}

	s.scope.Debug("The kube-proxy DaemonSet found, deleting")
	if err := remoteClient.Delete(ctx, ds, &client.DeleteOptions{}); err != nil {
		if apierrors.IsNotFound(err) {
			s.scope.Debug("The kube-proxy DaemonSet is not found, not deleted")
			return nil
		}
		return fmt.Errorf("deleting kube-proxy DaemonSet: %w", err)
	}
	record.Eventf(s.scope.InfraCluster(), "DeletedKubeProxy", "Kube-proxy has been removed from the cluster. Ensure you enable kube-proxy functionality via another mechanism")

	return nil
}
