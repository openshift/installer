/*
Copyright 2024 The Kubernetes Authors.

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
	"errors"
	"fmt"

	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/conditions"
	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infrav1alpha "sigs.k8s.io/cluster-api-provider-azure/api/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-azure/pkg/mutators"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

// ResourceReconciler reconciles a set of arbitrary ASO resources.
type ResourceReconciler struct {
	client.Client
	resources []*unstructured.Unstructured
	owner     resourceStatusObject
	watcher   watcher
}

type watcher interface {
	Watch(log logr.Logger, obj client.Object, handler handler.EventHandler, p ...predicate.Predicate) error
}

type resourceStatusObject interface {
	client.Object
	GetResourceStatuses() []infrav1alpha.ResourceStatus
	SetResourceStatuses([]infrav1alpha.ResourceStatus)
}

// Reconcile creates or updates the specified resources.
func (r *ResourceReconciler) Reconcile(ctx context.Context) error {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "controllers.ResourceReconciler.Reconcile")
	defer done()
	log.V(4).Info("reconciling resources")

	var newResourceStatuses []infrav1alpha.ResourceStatus

	for _, spec := range r.resources {
		gvk := spec.GroupVersionKind()
		spec.SetNamespace(r.owner.GetNamespace())

		log := log.WithValues("resource", klog.KObj(spec), "resourceVersion", gvk.GroupVersion(), "resourceKind", gvk.Kind)

		if err := controllerutil.SetControllerReference(r.owner, spec, r.Scheme()); err != nil {
			return fmt.Errorf("failed to set owner reference: %w", err)
		}

		if err := r.watcher.Watch(log, spec, handler.EnqueueRequestForOwner(r.Client.Scheme(), r.Client.RESTMapper(), r.owner)); err != nil {
			return fmt.Errorf("failed to watch resource: %w", err)
		}

		log.V(4).Info("applying resource")
		err := r.Patch(ctx, spec, client.Apply, client.FieldOwner("capz-manager"), client.ForceOwnership)
		if err != nil {
			return fmt.Errorf("failed to apply resource: %w", err)
		}

		ready, err := readyStatus(ctx, spec)
		if err != nil {
			return fmt.Errorf("failed to get ready status: %w", err)
		}
		newResourceStatuses = append(newResourceStatuses, infrav1alpha.ResourceStatus{
			Resource: infrav1alpha.StatusResource{
				Group:   gvk.Group,
				Version: gvk.Version,
				Kind:    gvk.Kind,
				Name:    spec.GetName(),
			},
			Ready: ready,
		})
	}

	for _, oldStatus := range r.owner.GetResourceStatuses() {
		needsDelete := true
		for _, newStatus := range newResourceStatuses {
			if oldStatus.Resource.Group == newStatus.Resource.Group &&
				oldStatus.Resource.Kind == newStatus.Resource.Kind &&
				oldStatus.Resource.Name == newStatus.Resource.Name {
				needsDelete = false
				break
			}
		}

		if needsDelete {
			updatedStatus, err := r.deleteResource(ctx, oldStatus.Resource)
			if err != nil {
				return err
			}
			if updatedStatus != nil {
				newResourceStatuses = append(newResourceStatuses, *updatedStatus)
			}
		}
	}

	r.owner.SetResourceStatuses(newResourceStatuses)

	return nil
}

// Pause pauses reconciliation of the specified resources.
func (r *ResourceReconciler) Pause(ctx context.Context) error {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "controllers.ResourceReconciler.Pause")
	defer done()
	log.V(4).Info("pausing resources")

	err := mutators.Pause(ctx, r.resources)
	if err != nil {
		if errors.As(err, &mutators.Incompatible{}) {
			err = reconcile.TerminalError(err)
		}
		return err
	}

	for _, spec := range r.resources {
		gvk := spec.GroupVersionKind()
		spec.SetNamespace(r.owner.GetNamespace())

		log := log.WithValues("resource", klog.KObj(spec), "resourceVersion", gvk.GroupVersion(), "resourceKind", gvk.Kind)

		log.V(4).Info("pausing resource")
		err := r.Patch(ctx, spec, client.Apply, client.FieldOwner("capz-manager"))
		if err != nil {
			return fmt.Errorf("failed to patch resource: %w", err)
		}
	}

	return nil
}

// Delete deletes the specified resources.
func (r *ResourceReconciler) Delete(ctx context.Context) error {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "controllers.ResourceReconciler.Delete")
	defer done()
	log.V(4).Info("deleting resources")

	var newResourceStatuses []infrav1alpha.ResourceStatus

	for _, spec := range r.owner.GetResourceStatuses() {
		newStatus, err := r.deleteResource(ctx, spec.Resource)
		if err != nil {
			return err
		}
		if newStatus != nil {
			newResourceStatuses = append(newResourceStatuses, *newStatus)
		}
	}

	r.owner.SetResourceStatuses(newResourceStatuses)

	return nil
}

func (r *ResourceReconciler) deleteResource(ctx context.Context, resource infrav1alpha.StatusResource) (*infrav1alpha.ResourceStatus, error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "controllers.ResourceReconciler.deleteResource")
	defer done()

	spec := &unstructured.Unstructured{}
	spec.SetGroupVersionKind(schema.GroupVersionKind{Group: resource.Group, Version: resource.Version, Kind: resource.Kind})
	spec.SetNamespace(r.owner.GetNamespace())
	spec.SetName(resource.Name)

	log = log.WithValues("resource", klog.KObj(spec), "resourceVersion", spec.GroupVersionKind().GroupVersion(), "resourceKind", spec.GetKind())

	log.V(4).Info("deleting resource")
	err := r.Client.Delete(ctx, spec)
	if apierrors.IsNotFound(err) {
		log.V(4).Info("resource has been deleted")
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to delete resource: %w", err)
	}

	err = r.Client.Get(ctx, client.ObjectKeyFromObject(spec), spec)
	if apierrors.IsNotFound(err) {
		log.V(4).Info("resource has been deleted")
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get resource: %w", err)
	}
	ready, err := readyStatus(ctx, spec)
	if err != nil {
		return nil, fmt.Errorf("failed to get ready status: %w", err)
	}

	return &infrav1alpha.ResourceStatus{
		Resource: resource,
		Ready:    ready,
	}, nil
}

func readyStatus(ctx context.Context, u *unstructured.Unstructured) (bool, error) {
	_, log, done := tele.StartSpanWithLogger(ctx, "controllers.ResourceReconciler.readyStatus")
	defer done()

	statusConditions, found, err := unstructured.NestedSlice(u.Object, "status", "conditions")
	if err != nil {
		return false, err
	}
	if !found {
		return false, nil
	}

	for _, el := range statusConditions {
		condition, ok := el.(map[string]interface{})
		if !ok {
			continue
		}
		condType, found, err := unstructured.NestedString(condition, "type")
		if !found || err != nil || condType != conditions.ConditionTypeReady {
			continue
		}

		observedGen, _, err := unstructured.NestedInt64(condition, "observedGeneration")
		if err != nil {
			return false, err
		}
		if observedGen < u.GetGeneration() {
			log.V(4).Info("waiting for ASO to reconcile the resource")
			return false, nil
		}

		readyStatus, _, err := unstructured.NestedString(condition, "status")
		if err != nil {
			return false, err
		}
		return readyStatus == string(metav1.ConditionTrue), nil
	}

	// no ready condition is set
	return false, nil
}
