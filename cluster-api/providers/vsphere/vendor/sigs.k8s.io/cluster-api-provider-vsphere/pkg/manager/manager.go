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

package manager

import (
	"context"
	"fmt"

	pkgerrors "github.com/pkg/errors"
	netopv1 "github.com/vmware-tanzu/net-operator-api/api/v1alpha1"
	nsxvpcv1 "github.com/vmware-tanzu/nsx-operator/pkg/apis/vpc/v1alpha1"
	vmoprv1alpha2 "github.com/vmware-tanzu/vm-operator/api/v1alpha2"
	vmoprv1alpha5 "github.com/vmware-tanzu/vm-operator/api/v1alpha5"
	vmoprv1alpha6 "github.com/vmware-tanzu/vm-operator/api/v1alpha6"
	ncpv1 "github.com/vmware-tanzu/vm-operator/external/ncp/api/v1alpha1"
	topologyv1 "github.com/vmware-tanzu/vm-operator/external/tanzu-topology/api/v1alpha1"
	"gopkg.in/fsnotify.v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	bootstrapv1 "sigs.k8s.io/cluster-api/api/bootstrap/kubeadm/v1beta2"
	controlplanev1 "sigs.k8s.io/cluster-api/api/controlplane/kubeadm/v1beta2"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	ipamv1 "sigs.k8s.io/cluster-api/api/ipam/v1beta2"
	ctrl "sigs.k8s.io/controller-runtime"

	infrav1beta1 "sigs.k8s.io/cluster-api-provider-vsphere/api/govmomi/v1beta1"
	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/api/govmomi/v1beta2"
	vmwarev1beta1 "sigs.k8s.io/cluster-api-provider-vsphere/api/supervisor/v1beta1"
	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/api/supervisor/v1beta2"
	capvcontext "sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	vmoprvhub "sigs.k8s.io/cluster-api-provider-vsphere/pkg/conversion/api/vmoperator/hub"
	conversionclient "sigs.k8s.io/cluster-api-provider-vsphere/pkg/conversion/client"
)

// Manager is a CAPV controller manager.
type Manager interface {
	ctrl.Manager

	// GetControllerManagerContext returns the controller manager's context.
	GetControllerManagerContext() *capvcontext.ControllerManagerContext
}

// New returns a new CAPV controller manager.
func New(ctx context.Context, opts Options) (Manager, error) {
	// Ensure the default options are set.
	opts.defaults()

	utilruntime.Must(apiextensionsv1.AddToScheme(opts.Scheme))
	utilruntime.Must(clientgoscheme.AddToScheme(opts.Scheme))
	utilruntime.Must(clusterv1.AddToScheme(opts.Scheme))
	utilruntime.Must(infrav1beta1.AddToScheme(opts.Scheme))
	utilruntime.Must(infrav1.AddToScheme(opts.Scheme))
	utilruntime.Must(controlplanev1.AddToScheme(opts.Scheme))
	utilruntime.Must(bootstrapv1.AddToScheme(opts.Scheme))
	utilruntime.Must(vmwarev1beta1.AddToScheme(opts.Scheme))
	utilruntime.Must(vmwarev1.AddToScheme(opts.Scheme))
	utilruntime.Must(vmoprvhub.AddToScheme(opts.Scheme))
	utilruntime.Must(vmoprv1alpha2.AddToScheme(opts.Scheme))
	utilruntime.Must(vmoprv1alpha5.AddToScheme(opts.Scheme))
	utilruntime.Must(vmoprv1alpha6.AddToScheme(opts.Scheme))
	utilruntime.Must(ncpv1.AddToScheme(opts.Scheme))
	utilruntime.Must(netopv1.AddToScheme(opts.Scheme))
	utilruntime.Must(nsxvpcv1.AddToScheme(opts.Scheme))
	utilruntime.Must(topologyv1.AddToScheme(opts.Scheme))
	utilruntime.Must(ipamv1.AddToScheme(opts.Scheme))

	// Build the controller manager.
	mgr, err := ctrl.NewManager(opts.KubeConfig, opts.Options)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "unable to create manager")
	}

	cc := mgr.GetClient()
	if opts.Converter != nil {
		cc, err = conversionclient.NewWithConverter(mgr.GetClient(), opts.Converter)
		if err != nil {
			return nil, pkgerrors.Wrap(err, "failed to create a conversion client")
		}
	}

	// Build the controller manager context.
	controllerManagerContext := &capvcontext.ControllerManagerContext{
		WatchNamespaces:         opts.Cache.DefaultNamespaces,
		Namespace:               opts.PodNamespace,
		Name:                    opts.PodName,
		LeaderElectionID:        opts.LeaderElectionID,
		LeaderElectionNamespace: opts.LeaderElectionNamespace,
		// NOTE: use a client that can handle conversions from API versions that exist in the supervisor
		// and the internal hub version used in the reconciler.
		Client:           cc,
		Logger:           opts.Logger,
		Scheme:           opts.Scheme,
		Username:         opts.Username,
		Password:         opts.Password,
		NetworkProvider:  opts.NetworkProvider,
		WatchFilterValue: opts.WatchFilterValue,
	}

	// Add the requested items to the manager.
	if err := opts.AddToManager(ctx, controllerManagerContext, mgr); err != nil {
		return nil, pkgerrors.Wrap(err, "failed to add resources to the manager")
	}

	return &manager{
		Manager:              mgr,
		controllerManagerCtx: controllerManagerContext,
	}, nil
}

type manager struct {
	ctrl.Manager
	controllerManagerCtx *capvcontext.ControllerManagerContext
}

func (m *manager) GetControllerManagerContext() *capvcontext.ControllerManagerContext {
	return m.controllerManagerCtx
}

// UpdateCredentials reads and updates credentials from the credentials file.
func UpdateCredentials(opts *Options) {
	opts.readAndSetCredentials()
}

// InitializeWatch adds a filesystem watcher for the capv credentials file.
// In case of any update to the credentials file, the new credentials are passed to the capv manager context.
func InitializeWatch(controllerManagerContext *capvcontext.ControllerManagerContext, managerOpts *Options) (watch *fsnotify.Watcher, err error) {
	capvCredentialsFile := managerOpts.CredentialsFile
	updateEventCh := make(chan bool)
	watch, err = fsnotify.NewWatcher()
	if err != nil {
		return nil, pkgerrors.Wrap(err, fmt.Sprintf("failed to create new Watcher for %s", capvCredentialsFile))
	}
	if err = watch.Add(capvCredentialsFile); err != nil {
		return nil, pkgerrors.Wrap(err, "received error on CAPV credential watcher")
	}
	go func() {
		for {
			select {
			case err := <-watch.Errors:
				controllerManagerContext.Logger.Error(err, "Received error on CAPV credential watcher")
			case event := <-watch.Events:
				controllerManagerContext.Logger.Info(fmt.Sprintf("Received event %v on the credential file %s", event, capvCredentialsFile))
				updateEventCh <- true
			}
		}
	}()

	go func() {
		for range updateEventCh {
			UpdateCredentials(managerOpts)
		}
	}()

	return watch, err
}
