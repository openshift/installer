/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package genruntime

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// ConvertibleStatus is implemented by status types to allow conversion among the different versions of a given status
//
// Why do we need both directions of conversion? See ConvertibleSpec for details.
type ConvertibleStatus interface {
	// ConvertStatusTo will populate the passed Status by copying over all available information from this one
	ConvertStatusTo(destination ConvertibleStatus) error

	// ConvertStatusFrom will populate this status by copying over all available information from the passed one
	ConvertStatusFrom(source ConvertibleStatus) error
}

// GetVersionedStatus returns a versioned status for the provided resource; the original API version used when the
// resource was first created is used to identify the version to return
func GetVersionedStatus(metaObject ARMMetaObject, scheme *runtime.Scheme) (ConvertibleStatus, error) {
	rsrc, err := NewEmptyVersionedResource(metaObject, scheme)
	if err != nil {
		return nil, errors.Wrap(err, "getting versioned status")
	}

	if rsrc.GetObjectKind().GroupVersionKind() == metaObject.GetObjectKind().GroupVersionKind() {
		// No conversion needed
		return metaObject.GetStatus(), nil
	}

	// Get a blank status and populate it
	status := rsrc.GetStatus()
	err = status.ConvertStatusFrom(metaObject.GetStatus())
	if err != nil {
		return nil, errors.Wrapf(err, "failed conversion of status")
	}

	return status, nil
}

// NewEmptyVersionedStatus returns a blank versioned status for the provided resource; the original API version used
// when the resource was first created is used to identify the version to return
func NewEmptyVersionedStatus(metaObject ARMMetaObject, scheme *runtime.Scheme) (ConvertibleStatus, error) {
	return NewEmptyVersionedStatusFromGVK(metaObject, scheme, GetOriginalGVK(metaObject))
}

// NewEmptyVersionedStatusFromGVK returns a blank versioned status for the provided resource and GVK
func NewEmptyVersionedStatusFromGVK(metaObject ARMMetaObject, scheme *runtime.Scheme, gvk schema.GroupVersionKind) (ConvertibleStatus, error) {
	rsrc, err := NewEmptyVersionedResourceFromGVK(scheme, gvk)
	if err != nil {
		return nil, errors.Wrap(err, "creating new empty versioned status")
	}

	if rsrc.GetObjectKind().GroupVersionKind() == metaObject.GetObjectKind().GroupVersionKind() {
		// No conversion needed, return an empty status from metaObject
		return metaObject.NewEmptyStatus(), nil
	}

	// Return the versioned one
	return rsrc.NewEmptyStatus(), nil
}

// NewEmptyARMStatus returns an empty ARM status object ready for deserialization from ARM; the original API version
// used when the resource was first created is used to create the appropriate version
func NewEmptyARMStatus(metaObject ARMMetaObject, scheme *runtime.Scheme) (ARMResourceStatus, error) {
	status, err := GetVersionedStatus(metaObject, scheme)
	if err != nil {
		return nil, errors.Wrap(err, "creating ARM status")
	}

	converter, ok := status.(FromARMConverter)
	if !ok {
		return nil, errors.Errorf("expected %T to implement genruntime.FromARMConverter", status)
	}

	return converter.NewEmptyARMValue(), nil
}
