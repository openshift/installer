/*
Copyright 2025 The Kubernetes Authors.

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

// Package capiutils contains utility functions for working with Cluster API resources.
// These have mostly been inlined as part of the CAPI 1.10 -> 1.11 upgrade,
// and should be removed when we switch to reading CAPI v1beta2 objects.
//
// Deprecated: This package is deprecated and is going to be removed when support for v1beta1 will be dropped.
package capiutils

import (
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"sigs.k8s.io/cluster-api/util/predicates"

	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
)

// ClusterUpdateInfraReady returns a predicate that returns true for an update event when a cluster has Status.InfrastructureReady changed from false to true
// it also returns true if the resource provided is not a Cluster to allow for use with controller-runtime NewControllerManagedBy.
func ClusterUpdateInfraReady(scheme *runtime.Scheme, logger logr.Logger) predicate.Funcs {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			log := logger.WithValues("predicate", "ClusterUpdateInfraReady", "eventType", "update")
			if gvk, err := apiutil.GVKForObject(e.ObjectOld, scheme); err == nil {
				log = log.WithValues(gvk.Kind, klog.KObj(e.ObjectOld))
			}

			oldCluster, ok := e.ObjectOld.(*clusterv1.Cluster)
			if !ok {
				log.V(4).Info("Expected Cluster", "type", fmt.Sprintf("%T", e.ObjectOld))
				return false
			}

			newCluster := e.ObjectNew.(*clusterv1.Cluster)

			if !oldCluster.Status.InfrastructureReady && newCluster.Status.InfrastructureReady {
				log.V(6).Info("Cluster infrastructure became ready, allowing further processing")
				return true
			}

			log.V(4).Info("Cluster infrastructure did not become ready, blocking further processing")
			return false
		},
		CreateFunc:  func(event.CreateEvent) bool { return false },
		DeleteFunc:  func(event.DeleteEvent) bool { return false },
		GenericFunc: func(event.GenericEvent) bool { return false },
	}
}

// ClusterPausedTransitions returns a predicate that returns true for an update event when a cluster has Spec.Paused changed.
func ClusterPausedTransitions(scheme *runtime.Scheme, logger logr.Logger) predicate.Funcs {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			log := logger.WithValues("predicate", "ClusterPausedTransitions", "eventType", "update")
			if gvk, err := apiutil.GVKForObject(e.ObjectOld, scheme); err == nil {
				log = log.WithValues(gvk.Kind, klog.KObj(e.ObjectOld))
			}

			oldCluster, ok := e.ObjectOld.(*clusterv1.Cluster)
			if !ok {
				log.V(4).Info("Expected Cluster", "type", fmt.Sprintf("%T", e.ObjectOld))
				return false
			}

			newCluster := e.ObjectNew.(*clusterv1.Cluster)

			if oldCluster.Spec.Paused && !newCluster.Spec.Paused {
				log.V(6).Info("Cluster unpausing, allowing further processing")
				return true
			}

			if !oldCluster.Spec.Paused && newCluster.Spec.Paused {
				log.V(6).Info("Cluster pausing, allowing further processing")
				return true
			}

			// This predicate always work in "or" with Paused predicates
			// so the logs are adjusted to not provide false negatives/verbosity at V<=5.
			log.V(6).Info("Cluster paused state was not changed, blocking further processing")
			return false
		},
		CreateFunc:  func(event.CreateEvent) bool { return false },
		DeleteFunc:  func(event.DeleteEvent) bool { return false },
		GenericFunc: func(event.GenericEvent) bool { return false },
	}
}

// ClusterPausedTransitionsOrInfrastructureReady returns a Predicate that returns true on Cluster Update events where
// either Cluster.Spec.Paused transitions or Cluster.Status.InfrastructureReady transitions to true.
// This implements a common requirement for some cluster-api and provider controllers (such as Machine Infrastructure
// controllers) to resume reconciliation when the Cluster gets paused or unpaused and when the infrastructure becomes ready.
// Example use:
//
//	err := controller.Watch(
//	    source.Kind(cache, &clusterv1.Cluster{}),
//	    handler.EnqueueRequestsFromMapFunc(clusterToMachines)
//	    predicates.ClusterPausedTransitionsOrInfrastructureReady(mgr.GetScheme(), r.Log),
//	)
func ClusterPausedTransitionsOrInfrastructureReady(scheme *runtime.Scheme, logger logr.Logger) predicate.Funcs {
	log := logger.WithValues("predicate", "ClusterPausedTransitionsOrInfrastructureReady")

	return predicates.Any(scheme, log, ClusterPausedTransitions(scheme, log), ClusterUpdateInfraReady(scheme, log))
}
