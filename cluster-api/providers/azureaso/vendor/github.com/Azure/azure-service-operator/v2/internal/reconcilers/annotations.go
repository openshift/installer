/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package reconcilers

// Annotation labels, used to store metadata about the state of the resource.
const (
	PollerResumeTokenAnnotation = "serviceoperator.azure.com/poller-resume-token"
	PollerResumeIDAnnotation    = "serviceoperator.azure.com/poller-resume-id"
	LatestReconciledGeneration  = "serviceoperator.azure.com/latest-reconciled-generation"
)
