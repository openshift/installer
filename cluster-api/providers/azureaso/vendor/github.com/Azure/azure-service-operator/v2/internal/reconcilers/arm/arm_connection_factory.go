/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package arm

import (
	"context"

	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

type ARMConnectionFactory func(context.Context, genruntime.ARMMetaObject) (Connection, error)
