/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package configmaps

import (
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/core"
)

// Exporter defines an interface for exporting ConfigMaps based on CEL expressions.
type Exporter interface {
	ConfigMapDestinationExpressions() []*core.DestinationExpression
}
