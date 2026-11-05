/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package customizations

import (
	"strings"

	"github.com/go-logr/logr"

	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/core"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/extensions"
)

var _ extensions.ErrorClassifier = &CassandraClusterExtension{}

// ClassifyError evaluates the provided error, returning details including whether it is fatal or can be retried.
// cloudError is the error returned from ARM.
// apiVersion is the ARM API version used for the request.
// log is a logger that can be used for telemetry.
// next is the default classification implementation to call.
// Returns CloudErrorDetails with classification and an error if classification itself fails.
func (ext *CassandraClusterExtension) ClassifyError(
	cloudError *genericarmclient.CloudError,
	apiVersion string,
	log logr.Logger,
	next extensions.ErrorClassifierFunc,
) (core.CloudErrorDetails, error) {
	details, err := next(cloudError)
	if err != nil {
		return core.CloudErrorDetails{}, err
	}

	// It's weird to retry on BadRequest as it certainly sounds like a fatal error, but this can happen while we're
	// still waiting on a RoleAssignment to take effect.
	if details.Code == "BadRequest" && strings.Contains(details.Message, "service principal") {
		details.Classification = core.ErrorRetryable
	}

	return details, nil
}
