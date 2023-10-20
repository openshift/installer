/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package reconcilers

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
)

// TODO: It's not clear to me that this file holds any value...
type ActionFunc = func(ctx context.Context) (ctrl.Result, error)
