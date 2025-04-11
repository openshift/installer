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

package controllers

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/scope"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/aksextensions"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/fleetsmembers"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/groups"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/managedclusters"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/privateendpoints"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/resourcehealth"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/subnets"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/virtualnetworks"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/secret"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// azureManagedControlPlaneService contains the services required by the cluster controller.
type azureManagedControlPlaneService struct {
	kubeclient client.Client
	scope      managedclusters.ManagedClusterScope
	services   []azure.ServiceReconciler
}

// newAzureManagedControlPlaneReconciler populates all the services based on input scope.
func newAzureManagedControlPlaneReconciler(scope *scope.ManagedControlPlaneScope) (*azureManagedControlPlaneService, error) {
	resourceHealthSvc, err := resourcehealth.New(scope)
	if err != nil {
		return nil, err
	}
	return &azureManagedControlPlaneService{
		kubeclient: scope.Client,
		scope:      scope,
		services: []azure.ServiceReconciler{
			groups.New(scope),
			virtualnetworks.New(scope),
			subnets.New(scope),
			managedclusters.New(scope),
			privateendpoints.New(scope),
			fleetsmembers.New(scope),
			aksextensions.New(scope),
			resourceHealthSvc,
		},
	}, nil
}

// Reconcile reconciles all the services in a predetermined order.
func (r *azureManagedControlPlaneService) Reconcile(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "controllers.azureManagedControlPlaneService.Reconcile")
	defer done()

	for _, service := range r.services {
		if err := service.Reconcile(ctx); err != nil {
			return errors.Wrapf(err, "failed to reconcile AzureManagedControlPlane service %s", service.Name())
		}
	}

	if err := r.reconcileKubeconfig(ctx); err != nil {
		return errors.Wrap(err, "failed to reconcile kubeconfig secret")
	}

	return nil
}

// Pause pauses all components making up the cluster.
func (r *azureManagedControlPlaneService) Pause(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "controllers.azureManagedControlPlaneService.Pause")
	defer done()

	for _, service := range r.services {
		pauser, ok := service.(azure.Pauser)
		if !ok {
			continue
		}
		if err := pauser.Pause(ctx); err != nil {
			return errors.Wrapf(err, "failed to pause AzureManagedControlPlane service %s", service.Name())
		}
	}

	return nil
}

// Delete reconciles all the services in a predetermined order.
func (r *azureManagedControlPlaneService) Delete(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "controllers.azureManagedControlPlaneService.Delete")
	defer done()

	// Delete services in reverse order of creation.
	for i := len(r.services) - 1; i >= 0; i-- {
		if err := r.services[i].Delete(ctx); err != nil {
			return errors.Wrapf(err, "failed to delete AzureManagedControlPlane service %s", r.services[i].Name())
		}
	}

	return nil
}

func (r *azureManagedControlPlaneService) reconcileKubeconfig(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "controllers.azureManagedControlPlaneService.reconcileKubeconfig")
	defer done()

	kubeConfigs := [][]byte{r.scope.GetAdminKubeconfigData(), r.scope.GetUserKubeconfigData()}

	for i, kubeConfigData := range kubeConfigs {
		if len(kubeConfigData) == 0 {
			continue
		}
		kubeConfigSecret := r.scope.MakeEmptyKubeConfigSecret()
		if i == 1 {
			// 2nd kubeconfig is the user kubeconfig
			kubeConfigSecret.Name = fmt.Sprintf("%s-user", kubeConfigSecret.Name)
		}
		if _, err := controllerutil.CreateOrUpdate(ctx, r.kubeclient, &kubeConfigSecret, func() error {
			kubeConfigSecret.Data = map[string][]byte{
				secret.KubeconfigDataName: kubeConfigData,
			}

			// When upgrading from an older version of CAPI, the kubeconfig secret may not have the required
			// cluster name label. Add it here to avoid kubeconfig issues during upgrades.
			if _, ok := kubeConfigSecret.Labels[clusterv1.ClusterNameLabel]; !ok {
				if kubeConfigSecret.Labels == nil {
					kubeConfigSecret.Labels = make(map[string]string)
				}
				kubeConfigSecret.Labels[clusterv1.ClusterNameLabel] = r.scope.ClusterName()
			}
			return nil
		}); err != nil {
			return errors.Wrap(err, "failed to reconcile kubeconfig secret for cluster")
		}
	}

	// store cluster-info for the cluster with the admin kubeconfig.
	kubeconfigFile, err := clientcmd.Load(kubeConfigs[0])
	if err != nil {
		return errors.Wrap(err, "failed to turn aks credentials into kubeconfig file struct")
	}

	cluster := kubeconfigFile.Contexts[kubeconfigFile.CurrentContext].Cluster
	caData := kubeconfigFile.Clusters[cluster].CertificateAuthorityData
	caSecret := r.scope.MakeClusterCA()
	if _, err := controllerutil.CreateOrUpdate(ctx, r.kubeclient, caSecret, func() error {
		caSecret.Data = map[string][]byte{
			secret.TLSCrtDataName: caData,
			secret.TLSKeyDataName: []byte("foo"),
		}
		return nil
	}); err != nil {
		return errors.Wrapf(err, "failed to reconcile certificate authority data secret for cluster")
	}

	if err := r.scope.StoreClusterInfo(ctx, caData); err != nil {
		return errors.Wrap(err, "failed to construct cluster-info")
	}

	return nil
}
