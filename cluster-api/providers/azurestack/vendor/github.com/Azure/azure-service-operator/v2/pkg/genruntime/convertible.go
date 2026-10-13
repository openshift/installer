/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package genruntime

import (
	"reflect"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

func findStorageGVK(scheme *runtime.Scheme, gk schema.GroupKind) (schema.GroupVersionKind, error) {
	interfaceType := reflect.TypeOf((*conversion.Hub)(nil)).Elem()

	for gvk, rt := range scheme.AllKnownTypes() {
		if gvk.Group != gk.Group || gvk.Kind != gk.Kind {
			continue
		}

		if reflect.PointerTo(rt).Implements(interfaceType) {
			return gvk, nil
		}
	}

	return schema.GroupVersionKind{}, errors.Errorf("couldn't find hub type for group kind %s", gk.String())
}

// ObjAsOriginalVersion returns the obj as the original API version used to create it.
func ObjAsOriginalVersion(obj ARMMetaObject, scheme *runtime.Scheme) (ARMMetaObject, error) {
	return ObjAsVersion(obj, scheme, GetOriginalGVK(obj))
}

// ObjAsVersion returns the object as the specified version, or an error if it cannot be converted to the
// requested version.
func ObjAsVersion(obj ARMMetaObject, scheme *runtime.Scheme, gvk schema.GroupVersionKind) (ARMMetaObject, error) {
	objGVK := obj.GetObjectKind().GroupVersionKind()

	versionedResource, err := NewEmptyVersionedResourceFromGVK(scheme, gvk)
	if err != nil {
		return nil, errors.Wrap(err, "getting empty versioned resource")
	}
	if objGVK == gvk {
		// No conversion needed, resource GVK is the same as what we want
		return obj, nil
	}

	// if the obj isn't the storage version (doesn't implement conversion.Hub) then we find the version that is the storage
	// version, construct an empty one and convert obj to the storage version.
	hub, ok := obj.(conversion.Hub)
	if !ok {
		// Note that if this function is called by a control-loop, the expectation is that the control-loop
		// is operating on the storage/hub version already, so it shouldn't be necessarily to take this code path.
		gk := schema.GroupKind{Group: gvk.Group, Kind: gvk.Kind}
		var storageGVK schema.GroupVersionKind
		storageGVK, err = findStorageGVK(scheme, gk)
		if err != nil {
			return nil, errors.Wrapf(err, "couldn't find storage GVK for %s", gk)
		}

		var storageObj ARMMetaObject
		storageObj, err = NewEmptyVersionedResourceFromGVK(scheme, storageGVK)
		if err != nil {
			return nil, errors.Wrap(err, "getting empty hub versioned resource")
		}

		hub, ok = storageObj.(conversion.Hub)
		if !ok {
			return nil, errors.Errorf("storage object %T with GVK %s is not a Hub object", storageObj, storageGVK)
		}

		if convertible, ok := obj.(conversion.Convertible); ok {
			err = convertible.ConvertTo(hub)
			if err != nil {
				return nil, errors.Wrapf(err, "couldn't convert %s to hub", objGVK)
			}
		} else {
			// This is unexpected/a bug
			return nil, errors.Errorf("obj %T was not convertible", obj)
		}
	}

	// Special case, if the destination is the hub we can just return here because we've already got it
	if _, ok := versionedResource.(conversion.Hub); ok {
		return hub.(ARMMetaObject), nil
	}

	// Convert from the storage version to the destination version.
	if convertible, ok := versionedResource.(conversion.Convertible); ok {
		err = convertible.ConvertFrom(hub)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to convert resource to expected version. have: %s, want: %s", objGVK, gvk)
		}
	} else {
		return nil, errors.Errorf("obj %T was not convertible", versionedResource)
	}

	obj = versionedResource

	return obj, nil
}
