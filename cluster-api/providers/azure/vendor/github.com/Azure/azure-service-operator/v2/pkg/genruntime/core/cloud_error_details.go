/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package core

import (
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/retry"
)

type CloudErrorDetails struct {
	// Classification specifies if the error is fatal or transient
	Classification ErrorClassification
	// Retry defines the speed at which the error should be retried. If this is not set,
	// the default is retry.Slow.
	Retry retry.Classification
	// Code is the error code
	Code string
	// Message is the error message
	Message string
}
