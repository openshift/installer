/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package customizations

import (
	"regexp"

	"github.com/go-logr/logr"

	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/core"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/extensions"
)

var _ extensions.ErrorClassifier = &DiskExtension{}

// diskAccessNotFoundRegex matches the error returned by CRP when the DiskAccess doesn't exist yet.
var diskAccessNotFoundRegex = regexp.MustCompile("DiskAccess.*not found")

// diskAccessFailedStateRegex matches the error returned by CRP when a DiskAccess is not in successful state yet.
// Note that even though this says failed, this is returned even if the resource is in a transitioning state (and
// will succeed eventually)
var diskAccessFailedStateRegex = regexp.MustCompile("DiskAccess.*is in failed state.")

// ClassifyError evaluates the provided error, returning whether it is fatal or can be retried.
func (e *DiskExtension) ClassifyError(
	cloudError *genericarmclient.CloudError,
	apiVersion string,
	log logr.Logger,
	next extensions.ErrorClassifierFunc,
) (core.CloudErrorDetails, error) {
	details, err := next(cloudError)
	if err != nil {
		return core.CloudErrorDetails{}, err
	}

	// If the DiskAccess doesn't exist yet, we retry as it may be being created
	if shouldRetry(details) {
		details.Classification = core.ErrorRetryable
	}

	return details, nil
}

func shouldRetry(details core.CloudErrorDetails) bool {
	if details.Code == "BadRequest" && diskAccessNotFoundRegex.MatchString(details.Message) {
		return true
	}

	if details.Code == "Conflict" && diskAccessFailedStateRegex.MatchString(details.Message) {
		return true
	}

	return false
}
