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

package context

import (
	"sync"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

// ControllerManagerContext is the context of the controller that owns the
// controllers.
type ControllerManagerContext struct {
	// Namespace is the namespace in which the resource is located responsible
	// for running the controller manager.
	Namespace string

	// Name is the name of the controller manager.
	Name string

	// LeaderElectionID is the information used to identify the object
	// responsible for synchronizing leader election.
	LeaderElectionID string

	// LeaderElectionNamespace is the namespace in which the LeaderElection
	// object is located.
	LeaderElectionNamespace string

	// WatchNamespaces are the namespaces the controllers watches for changes. If
	// no value is specified then all namespaces are watched.
	WatchNamespaces map[string]cache.Config

	// Client is the controller manager's client.
	Client client.Client

	// Logger is the controller manager's logger.
	Logger logr.Logger

	// Scheme is the controller manager's API scheme.
	Scheme *runtime.Scheme

	// Username is the username for the account used to access remote vSphere
	// endpoints.
	Username string

	// Password is the password for the account used to access remote vSphere
	// endpoints.
	Password string

	// NetworkProvider is the network provider used by Supervisor based clusters
	NetworkProvider string

	// WatchFilterValue is used to filter incoming objects by label.
	WatchFilterValue string

	genericEventCache sync.Map
}

// String returns ControllerManagerName.
func (c *ControllerManagerContext) String() string {
	return c.Name
}

// GetGenericEventChannelFor returns a generic event channel for a resource
// specified by the provided GroupVersionKind.
func (c *ControllerManagerContext) GetGenericEventChannelFor(gvk schema.GroupVersionKind) chan event.GenericEvent {
	if val, ok := c.genericEventCache.Load(gvk); ok {
		return val.(chan event.GenericEvent)
	}
	val, _ := c.genericEventCache.LoadOrStore(gvk, make(chan event.GenericEvent))
	return val.(chan event.GenericEvent)
}
