/*
Copyright 2025 The ORC Authors.

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

package predicates

import (
	"fmt"
	"slices"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/logging"
)

// newServerResourceChanged is a generic predicate factory that filters Server events based on
// changes to a list of resource IDs extracted by the provided getIDs function.
func newServerResourceChanged(
	log logr.Logger,
	functionName string,
	getIDs func(*orcv1alpha1.Server) []string,
) predicate.Predicate {
	getServer := func(obj client.Object, event string) *orcv1alpha1.Server {
		server, ok := obj.(*orcv1alpha1.Server)
		if !ok {
			log.Info(fmt.Sprintf("%s got unexpected object type", functionName),
				"got", fmt.Sprintf("%T", obj),
				"expected", fmt.Sprintf("%T", &orcv1alpha1.Server{}),
				"event", event)
			return nil
		}
		return server
	}

	log = log.WithValues("watchKind", fmt.Sprintf("%T", &orcv1alpha1.Server{}))

	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			log := log.WithValues("name", e.Object.GetName(), "namespace", e.Object.GetNamespace())
			log.V(logging.Debug).Info("Observed create")

			server := getServer(e.Object, "create")
			if server == nil {
				return false
			}

			// Trigger reconciliation if server is created with resources
			ids := getIDs(server)
			return len(ids) > 0
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			log := log.WithValues("name", e.ObjectOld.GetName(), "namespace", e.ObjectOld.GetNamespace())
			log.V(logging.Debug).Info("Observed update")

			oldServer := getServer(e.ObjectOld, "update")
			newServer := getServer(e.ObjectNew, "update")

			if oldServer == nil || newServer == nil {
				return false
			}

			oldIDs := getIDs(oldServer)
			newIDs := getIDs(newServer)

			// Trigger reconciliation if the resource IDs changed
			return !slices.Equal(oldIDs, newIDs)
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			log := log.WithValues("name", e.Object.GetName(), "namespace", e.Object.GetNamespace())
			log.V(logging.Debug).Info("Observed delete")

			server := getServer(e.Object, "delete")
			if server == nil {
				return false
			}

			// Trigger reconciliation if server is deleted with resources
			ids := getIDs(server)
			return len(ids) > 0
		},
	}
}

// NewServerVolumesChanged returns a predicate that filters Server events to only those
// where the status.resource.volumes field changed.
func NewServerVolumesChanged(log logr.Logger) predicate.Predicate {
	getVolumeIDs := func(server *orcv1alpha1.Server) []string {
		if server.Status.Resource == nil {
			return nil
		}
		volumeIDs := make([]string, 0, len(server.Status.Resource.Volumes))
		for i := range server.Status.Resource.Volumes {
			volumeID := server.Status.Resource.Volumes[i].ID
			if volumeID != "" {
				volumeIDs = append(volumeIDs, volumeID)
			}
		}
		return volumeIDs
	}

	return newServerResourceChanged(log, "NewServerVolumesChanged", getVolumeIDs)
}

// NewServerInterfacesChanged returns a predicate that filters Server events to only those
// where the status.resource.interfaces field changed.
func NewServerInterfacesChanged(log logr.Logger) predicate.Predicate {
	getPortIDs := func(server *orcv1alpha1.Server) []string {
		if server.Status.Resource == nil {
			return nil
		}
		portIDs := make([]string, 0, len(server.Status.Resource.Interfaces))
		for i := range server.Status.Resource.Interfaces {
			portID := server.Status.Resource.Interfaces[i].PortID
			if portID != "" {
				portIDs = append(portIDs, portID)
			}
		}
		return portIDs
	}

	return newServerResourceChanged(log, "NewServerInterfacesChanged", getPortIDs)
}
