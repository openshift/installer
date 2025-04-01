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

package mutators

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Azure/azure-service-operator/v2/pkg/common/annotations"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ResourcesMutator mutates in-place a slice of ASO resources to be reconciled. These mutations make only the
// changes strictly necessary for CAPZ resources to play nice with Cluster API. Any mutations should be logged
// and mutations that conflict with user-defined values should be rejected by returning Incompatible.
type ResourcesMutator func(context.Context, []*unstructured.Unstructured) error

type mutation struct {
	location string
	val      any
	reason   string
}

func logMutation(log logr.Logger, mutation mutation) {
	log.V(4).Info(fmt.Sprintf("setting %s to %v %s", mutation.location, mutation.val, mutation.reason))
}

// Incompatible describes an error where a piece of user-defined configuration does not match what CAPZ
// requires.
type Incompatible struct {
	mutation
	userVal any
}

func (e Incompatible) Error() string {
	return fmt.Sprintf("incompatible value: value at %s set by user to %v but CAPZ must set it to %v %s. The user-defined value must not be defined, or must match CAPZ's desired value.", e.location, e.userVal, e.val, e.reason)
}

// ApplyMutators applies the given mutators to the given resources.
func ApplyMutators(ctx context.Context, resources []runtime.RawExtension, mutators ...ResourcesMutator) ([]*unstructured.Unstructured, error) {
	us := []*unstructured.Unstructured{}
	for _, resource := range resources {
		u := &unstructured.Unstructured{}
		if err := u.UnmarshalJSON(resource.Raw); err != nil {
			return nil, fmt.Errorf("failed to unmarshal resource JSON: %w", err)
		}
		us = append(us, u)
	}
	for _, mutator := range mutators {
		if err := mutator(ctx, us); err != nil {
			err = fmt.Errorf("failed to run mutator: %w", err)
			if errors.As(err, &Incompatible{}) {
				err = reconcile.TerminalError(err)
			}
			return nil, err
		}
	}
	return us, nil
}

// ToUnstructured converts the given resources to Unstructured.
func ToUnstructured(ctx context.Context, resources []runtime.RawExtension) ([]*unstructured.Unstructured, error) {
	return ApplyMutators(ctx, resources)
}

// Pause sets the "skip" reconcile policy on all resources to facilitate a CAPI pause.
func Pause(ctx context.Context, resources []*unstructured.Unstructured) error {
	_, log, done := tele.StartSpanWithLogger(ctx, "mutators.Pause")
	defer done()

	for i, resource := range resources {
		resourcePath := "spec.resources[" + strconv.Itoa(i) + "]"
		policyPath := []string{"metadata", "annotations", annotations.ReconcilePolicy}
		capiPolicy := string(annotations.ReconcilePolicySkip)
		userPolicy, userDefined := resource.GetAnnotations()[annotations.ReconcilePolicy]

		setPolicy := mutation{
			location: resourcePath + "." + strings.Join(policyPath, "."),
			val:      capiPolicy,
			reason:   "because the CAPZ resource is paused",
		}
		if userDefined && userPolicy != capiPolicy {
			return Incompatible{
				mutation: setPolicy,
				userVal:  userPolicy,
			}
		}

		logMutation(log, setPolicy)
		anns := resource.GetAnnotations()
		if anns == nil {
			anns = make(map[string]string)
		}
		anns[annotations.ReconcilePolicy] = capiPolicy
		resource.SetAnnotations(anns)
	}

	return nil
}
