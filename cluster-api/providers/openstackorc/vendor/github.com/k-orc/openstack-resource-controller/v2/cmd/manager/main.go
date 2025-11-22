/*
Copyright 2024.

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

package main

import (
	"flag"
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/flavor"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/floatingip"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/image"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/network"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/port"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/project"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/router"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/routerinterface"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/securitygroup"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/server"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/servergroup"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/subnet"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/volume"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/volumetype"
	internalmanager "github.com/k-orc/openstack-resource-controller/v2/internal/manager"
	"github.com/k-orc/openstack-resource-controller/v2/internal/scheme"
	"github.com/k-orc/openstack-resource-controller/v2/internal/scope"
	// +kubebuilder:scaffold:imports
)

var (
	defaultCACertsPath string
	namespaceList      []string
)

func main() {
	setupLog := ctrl.Log.WithName("setup")

	orcOpts := internalmanager.Options{}
	flag.StringVar(&orcOpts.MetricsAddr, "metrics-bind-address", "0", "The address the metrics endpoint binds to. "+
		"Use :8443 for HTTPS or :8080 for HTTP, or leave as 0 to disable the metrics service.")
	flag.StringVar(&orcOpts.ProbeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&orcOpts.EnableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&orcOpts.SecureMetrics, "metrics-secure", true,
		"If set, the metrics endpoint is served securely via HTTPS. Use --metrics-secure=false to use HTTP instead.")
	flag.BoolVar(&orcOpts.EnableHTTP2, "enable-http2", false,
		"If set, HTTP/2 will be enabled for the metrics and webhook servers")
	flag.IntVar(&orcOpts.ScopeCacheMaxSize, "scope-cache-max-size", 10,
		"The maximum credentials count the operator should keep in cache. "+
			"Setting this value to 0 means no cache.")
	flag.StringVar(&defaultCACertsPath, "default-ca-certs", "",
		"The path to a PEM-encoded CA Certificate file to supply as default for OpenStack API requests.")
	flag.Func("namespace", "A namespace that the controller watches to reconcile ORC objects. "+
		"Can be specified multiple times.", func(ns string) error {
		namespaceList = append(namespaceList, ns)
		return nil
	})

	zapOpts := zap.Options{
		Development: true,
	}
	zapOpts.BindFlags(flag.CommandLine)
	flag.Parse()

	log := zap.New(zap.UseFlagOptions(&zapOpts))
	ctrl.SetLogger(log)

	// Setup the context that's going to be used in controllers and for the manager.
	ctx := ctrl.SetupSignalHandler()

	var caCerts []byte
	if defaultCACertsPath != "" {
		var err error
		caCerts, err = os.ReadFile(defaultCACertsPath)
		if err != nil {
			setupLog.Error(err, "unable to read provided ca certificates file")
			os.Exit(1)
		}
	}
	scopeFactory := scope.NewFactory(orcOpts.ScopeCacheMaxSize, caCerts)

	controllers := []interfaces.Controller{
		image.New(scopeFactory),
		network.New(scopeFactory),
		subnet.New(scopeFactory),
		router.New(scopeFactory),
		routerinterface.New(scopeFactory),
		port.New(scopeFactory),
		floatingip.New(scopeFactory),
		flavor.New(scopeFactory),
		securitygroup.New(scopeFactory),
		server.New(scopeFactory),
		servergroup.New(scopeFactory),
		project.New(scopeFactory),
		volume.New(scopeFactory),
		volumetype.New(scopeFactory),
	}

	restConfig := ctrl.GetConfigOrDie()
	orcOpts.WatchNamespaces = namespaceList
	err := internalmanager.Run(ctx, &orcOpts, restConfig, scheme.New(), setupLog, log, controllers)
	if err != nil {
		setupLog.Error(err, "Error starting manager")
		os.Exit(1)
	}
}
