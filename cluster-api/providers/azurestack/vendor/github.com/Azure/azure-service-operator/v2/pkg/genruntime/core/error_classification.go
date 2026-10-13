/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package core

type ErrorClassification string

const (
	ErrorRetryable = ErrorClassification("retryable")
	ErrorFatal     = ErrorClassification("fatal")
)
