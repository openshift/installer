/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package entra

import (
	"context"

	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

type EntraConnectionFactory func(context.Context, genruntime.EntraMetaObject) (Connection, error)
