/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package metrics

type Metrics interface {
	RegisterMetrics()
}

func RegisterMetrics(metrics ...Metrics) {
	for _, metric := range metrics {
		metric.RegisterMetrics()
	}
}
