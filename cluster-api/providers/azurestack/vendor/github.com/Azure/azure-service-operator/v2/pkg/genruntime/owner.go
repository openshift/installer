/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package genruntime

import (
	"context"
	"reflect"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/Azure/azure-service-operator/v2/internal/ownerutil"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/core"
)

// CheckTargetOwnedByObj raises an error if the target object is not owned by obj.
func CheckTargetOwnedByObj(obj client.Object, target client.Object) error {
	ownerRefs := target.GetOwnerReferences()
	owned := false
	for _, ref := range ownerRefs {
		if ref.UID == obj.GetUID() {
			owned = true
			break
		}
	}

	if !owned {
		return core.NewNotOwnedError(
			target.GetNamespace(),
			target.GetName(),
			target.GetObjectKind().GroupVersionKind(),
			obj.GetName(),
			obj.GetObjectKind().GroupVersionKind())
	}

	return nil
}

func safeNewObj(obj client.Object) (result client.Object, err error) {
	defer func() {
		if oops := recover(); oops != nil {
			err = errors.Errorf("failed to create new %T.", obj)
		}
	}()
	result = reflect.New(reflect.TypeOf(obj).Elem()).Interface().(client.Object)

	// No return needed here as
	return result, err
}

// ApplyObjAndEnsureOwner applies the object (similar to kubectl apply). If the object does not exist
// it is created. If it exists, it is updated.
func ApplyObjAndEnsureOwner(ctx context.Context, c client.Client, owner client.Object, obj client.Object) (controllerutil.OperationResult, error) {
	updatedObj, err := safeNewObj(obj)
	if err != nil {
		return controllerutil.OperationResultNone, err
	}
	updatedObj.SetNamespace(obj.GetNamespace())
	updatedObj.SetName(obj.GetName())

	objProps, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return controllerutil.OperationResultNone, errors.Wrapf(err, "failed to convert obj to unstructured")
	}

	result, err := controllerutil.CreateOrUpdate(ctx, c, updatedObj, func() error {
		// If the secret exists but isn't owned by our resource then it must have been created
		// by the user. We want to avoid overwriting or otherwise modifying secrets of theirs.
		if updatedObj.GetResourceVersion() != "" {
			if err = CheckTargetOwnedByObj(owner, updatedObj); err != nil {
				return err
			}
		}
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(objProps, updatedObj)
		if err != nil {
			return errors.Wrap(err, "failed to convert unstructured to obj")
		}

		ownerRef := ownerutil.MakeOwnerReference(owner)
		updatedObj.SetOwnerReferences(ownerutil.EnsureOwnerRef(updatedObj.GetOwnerReferences(), ownerRef))
		return nil
	})
	if err != nil {
		return controllerutil.OperationResultNone, err
	}

	return result, nil
}

// ApplyObjsAndEnsureOwner applies the specified collection of objects (similar to kubectl apply). If the objects do not exist
// they are created. If they exist, they are updated. An attempt is made to apply each object before returning an error.
func ApplyObjsAndEnsureOwner(ctx context.Context, client client.Client, owner client.Object, objs []client.Object) ([]controllerutil.OperationResult, error) {
	var errs []error
	results := make([]controllerutil.OperationResult, 0, len(objs))

	for _, obj := range objs {
		result, err := ApplyObjAndEnsureOwner(ctx, client, owner, obj)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		results = append(results, result)
	}

	err := kerrors.NewAggregate(errs)
	if err != nil {
		return nil, err
	}

	return results, nil
}
