/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package secrets

import (
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/core"
)

// Exporter defines an interface for exporting Secrets based on CEL expressions.
type Exporter interface {
	SecretDestinationExpressions() []*core.DestinationExpression
}
