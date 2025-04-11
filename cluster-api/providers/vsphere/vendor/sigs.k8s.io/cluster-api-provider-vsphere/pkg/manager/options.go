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
	"os"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
	ctrlmgr "sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/yaml"

	capvcontext "sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
)

// AddToManagerFunc is a function that can be optionally specified with
// the manager's Options in order to explicitly decide what controllers and
// webhooks to add to the manager.
type AddToManagerFunc func(context.Context, *capvcontext.ControllerManagerContext, ctrlmgr.Manager) error

// Options describes the options used to create a new CAPV manager.
type Options struct {
	ctrlmgr.Options

	// PodNamespace is the namespace in which the pod running the controller
	// maintains a leader election lock.
	//
	// Defaults to the eponymous constant in this package.
	PodNamespace string

	// PodName is the name of the pod running the controller manager.
	//
	// Defaults to the eponymous constant in this package.
	PodName string

	// Username is the username for the account used to access remote vSphere
	// endpoints.
	Username string

	// Password is the password for the account used to access remote vSphere
	// endpoints.
	Password string

	// CredentialsFile is the file that contains credentials of CAPV
	CredentialsFile string

	KubeConfig *rest.Config

	// AddToManager is a function that can be optionally specified with
	// the manager's Options in order to explicitly decide what controllers
	// and webhooks to add to the manager.
	AddToManager AddToManagerFunc

	// NetworkProvider is the network provider used by Supervisor based clusters.
	// If not set, it will default to a DummyNetworkProvider which is intended for testing purposes.
	// VIM based clusters and managers will not need to set this flag.
	NetworkProvider string

	// WatchFilterValue is used to filter incoming objects by label.
	//
	// Defaults to the empty string and by that not filter anything.
	WatchFilterValue string
}

func (o *Options) defaults() {
	if o.Logger.GetSink() == nil {
		o.Logger = ctrllog.Log
	}

	if o.PodName == "" {
		o.PodName = DefaultPodName
	}

	if o.KubeConfig == nil {
		o.KubeConfig = config.GetConfigOrDie()
	}

	if o.Scheme == nil {
		o.Scheme = runtime.NewScheme()
	}

	if o.Username == "" || o.Password == "" {
		o.readAndSetCredentials()
	}

	if ns, ok := os.LookupEnv("POD_NAMESPACE"); ok {
		o.PodNamespace = ns
	} else if data, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace"); err == nil {
		if ns := strings.TrimSpace(string(data)); ns != "" {
			o.PodNamespace = ns
		}
	} else {
		o.PodNamespace = DefaultPodNamespace
	}
}

func (o *Options) getCredentials() map[string]string {
	file, err := os.ReadFile(o.CredentialsFile)
	if err != nil {
		o.Logger.Error(err, "error opening credentials file")
		return map[string]string{}
	}

	credentials := map[string]string{}
	if err := yaml.Unmarshal(file, &credentials); err != nil {
		o.Logger.Error(err, "error unmarshaling credentials to yaml")
		return map[string]string{}
	}

	return credentials
}

func (o *Options) readAndSetCredentials() {
	credentials := o.getCredentials()
	o.Username = credentials["username"]
	o.Password = credentials["password"]
}
