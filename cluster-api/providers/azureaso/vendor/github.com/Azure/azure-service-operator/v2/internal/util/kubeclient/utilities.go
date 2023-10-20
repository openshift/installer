/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package kubeclient

import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func IgnoreNotFound(err error) error {
	return client.IgnoreNotFound(err)
}

func IgnoreNotFoundAndConflict(err error) error {
	if IsNotFoundOrConflict(err) {
		return nil
	}
	return err
}

func IsNotFoundOrConflict(err error) bool {
	return apierrors.IsConflict(err) || apierrors.IsNotFound(err)
}
