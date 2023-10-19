/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package arm

import (
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

const (
	PollerResumeTokenAnnotation = "serviceoperator.azure.com/poller-resume-token"
	PollerResumeIDAnnotation    = "serviceoperator.azure.com/poller-resume-id"
)

// GetPollerResumeToken returns a poller ID and the poller token
func GetPollerResumeToken(obj genruntime.MetaObject) (string, string, bool) {
	token, hasResumeToken := obj.GetAnnotations()[PollerResumeTokenAnnotation]
	id, hasResumeID := obj.GetAnnotations()[PollerResumeIDAnnotation]

	return id, token, hasResumeToken && hasResumeID
}

func SetPollerResumeToken(obj genruntime.MetaObject, id string, token string) {
	genruntime.AddAnnotation(obj, PollerResumeTokenAnnotation, token)
	genruntime.AddAnnotation(obj, PollerResumeIDAnnotation, id)
}

// ClearPollerResumeToken clears the poller resume token and ID annotations
func ClearPollerResumeToken(obj genruntime.MetaObject) {
	genruntime.RemoveAnnotation(obj, PollerResumeTokenAnnotation)
	genruntime.RemoveAnnotation(obj, PollerResumeIDAnnotation)
}
