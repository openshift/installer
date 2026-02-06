/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package retry

type Classification string

const (
	// None means that this classification is not expected to ever retry. It should only be set for fatal errors
	None     = Classification("None")
	Fast     = Classification("RetryFast")
	Slow     = Classification("RetrySlow")
	VerySlow = Classification("RetryVerySlow")
)
