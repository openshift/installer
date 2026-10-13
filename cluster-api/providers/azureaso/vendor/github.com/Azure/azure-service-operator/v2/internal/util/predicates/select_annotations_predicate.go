/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package predicates

import (
	"reflect"

	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type HasAnnotationChanged func(old *string, new *string) bool

func HasBasicAnnotationChanged(old *string, new *string) bool {
	return !reflect.DeepEqual(old, new)
}

// MakeSelectAnnotationChangedPredicate creates a selectAnnotationChangedPredicate watching for
// changes to select annotations.
// annotations is a map of annotations to HasAnnotationChanged handlers which define if the annotation has been
// changed in a way we care about.
func MakeSelectAnnotationChangedPredicate(annotations map[string]HasAnnotationChanged) predicate.Predicate {
	return selectAnnotationChangedPredicate{
		annotations: annotations,
	}
}

type selectAnnotationChangedPredicate struct {
	predicate.Funcs

	annotations map[string]HasAnnotationChanged
}

var _ predicate.Predicate = selectAnnotationChangedPredicate{}

// Update implements UpdateEvent filter for annotation changes.
func (p selectAnnotationChangedPredicate) Update(e event.UpdateEvent) bool {
	if e.ObjectOld == nil {
		return false
	}
	if e.ObjectNew == nil {
		return false
	}

	newAnnotations := e.ObjectNew.GetAnnotations()
	oldAnnotations := e.ObjectOld.GetAnnotations()

	for k, f := range p.annotations {
		oldAnnotation := valueOrNil(oldAnnotations, k)
		newAnnotation := valueOrNil(newAnnotations, k)

		changed := f(oldAnnotation, newAnnotation)
		if changed {
			return true
		}
	}

	return false
}

func valueOrNil(annotations map[string]string, key string) *string {
	val, ok := annotations[key]
	if !ok {
		return nil
	}
	return &val
}
