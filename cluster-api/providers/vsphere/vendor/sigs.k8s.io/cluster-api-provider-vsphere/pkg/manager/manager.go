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
	goctx "context"
	"fmt"
	"os"

	"github.com/pkg/errors"
	netopv1 "github.com/vmware-tanzu/net-operator-api/api/v1alpha1"
	vmoprv1 "github.com/vmware-tanzu/vm-operator/api/v1alpha1"
	ncpv1 "github.com/vmware-tanzu/vm-operator/external/ncp/api/v1alpha1"
	topologyv1 "github.com/vmware-tanzu/vm-operator/external/tanzu-topology/api/v1alpha1"
	"gopkg.in/fsnotify.v1"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	bootstrapv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1beta1"
	controlplanev1 "sigs.k8s.io/cluster-api/controlplane/kubeadm/api/v1beta1"
	ipamv1 "sigs.k8s.io/cluster-api/exp/ipam/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"

	infrav1a3 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1alpha3"
	infrav1a4 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1alpha4"
	infrav1b1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	vmwarev1b1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/vmware/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/record"
)

// Manager is a CAPV controller manager.
type Manager interface {
	ctrl.Manager

	// GetContext returns the controller manager's context.
	GetContext() *context.ControllerManagerContext
}

// New returns a new CAPV controller manager.
func New(opts Options) (Manager, error) {
	// Ensure the default options are set.
	opts.defaults()

	_ = clientgoscheme.AddToScheme(opts.Scheme)
	_ = clusterv1.AddToScheme(opts.Scheme)
	_ = infrav1a3.AddToScheme(opts.Scheme)
	_ = infrav1a4.AddToScheme(opts.Scheme)
	_ = infrav1b1.AddToScheme(opts.Scheme)
	_ = controlplanev1.AddToScheme(opts.Scheme)
	_ = bootstrapv1.AddToScheme(opts.Scheme)
	_ = vmwarev1b1.AddToScheme(opts.Scheme)
	_ = vmoprv1.AddToScheme(opts.Scheme)
	_ = ncpv1.AddToScheme(opts.Scheme)
	_ = netopv1.AddToScheme(opts.Scheme)
	_ = topologyv1.AddToScheme(opts.Scheme)
	_ = ipamv1.AddToScheme(opts.Scheme)
	// +kubebuilder:scaffold:scheme

	podName, err := os.Hostname()
	if err != nil {
		podName = DefaultPodName
	}

	// Build the controller manager.
	mgr, err := ctrl.NewManager(opts.KubeConfig, opts.Options)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create manager")
	}

	// Build the controller manager context.
	controllerManagerContext := &context.ControllerManagerContext{
		Context:                 goctx.Background(),
		WatchNamespaces:         opts.Cache.Namespaces,
		Namespace:               opts.PodNamespace,
		Name:                    opts.PodName,
		LeaderElectionID:        opts.LeaderElectionID,
		LeaderElectionNamespace: opts.LeaderElectionNamespace,
		Client:                  mgr.GetClient(),
		Logger:                  opts.Logger.WithName(opts.PodName),
		Recorder:                record.New(mgr.GetEventRecorderFor(fmt.Sprintf("%s/%s", opts.PodNamespace, podName))),
		Scheme:                  opts.Scheme,
		Username:                opts.Username,
		Password:                opts.Password,
		EnableKeepAlive:         opts.EnableKeepAlive,
		KeepAliveDuration:       opts.KeepAliveDuration,
		NetworkProvider:         opts.NetworkProvider,
		WatchFilterValue:        opts.WatchFilterValue,
	}

	// Add the requested items to the manager.
	if err := opts.AddToManager(controllerManagerContext, mgr); err != nil {
		return nil, errors.Wrap(err, "failed to add resources to the manager")
	}

	// +kubebuilder:scaffold:builder

	return &manager{
		Manager: mgr,
		ctx:     controllerManagerContext,
	}, nil
}

type manager struct {
	ctrl.Manager
	ctx *context.ControllerManagerContext
}

func (m *manager) GetContext() *context.ControllerManagerContext {
	return m.ctx
}

func UpdateCredentials(opts *Options) {
	opts.readAndSetCredentials()
}

// InitializeWatch adds a filesystem watcher for the capv credentials file
// In case of any update to the credentials file, the new credentials are passed to the capv manager context.
func InitializeWatch(ctx *context.ControllerManagerContext, managerOpts *Options) (watch *fsnotify.Watcher, err error) {
	capvCredentialsFile := managerOpts.CredentialsFile
	updateEventCh := make(chan bool)
	watch, err = fsnotify.NewWatcher()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to create new Watcher for %s", capvCredentialsFile))
	}
	if err = watch.Add(capvCredentialsFile); err != nil {
		return nil, errors.Wrap(err, "received error on CAPV credential watcher")
	}
	go func() {
		for {
			select {
			case err := <-watch.Errors:
				ctx.Logger.Error(err, "received error on CAPV credential watcher")
			case event := <-watch.Events:
				ctx.Logger.Info(fmt.Sprintf("received event %v on the credential file %s", event, capvCredentialsFile))
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
