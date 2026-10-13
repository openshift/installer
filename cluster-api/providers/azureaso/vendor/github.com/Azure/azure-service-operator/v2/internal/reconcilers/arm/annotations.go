/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package arm

import (
	"strconv"

	"github.com/Azure/azure-service-operator/v2/internal/reconcilers"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

// GetPollerResumeToken returns a poller ID and the poller token
func GetPollerResumeToken(obj genruntime.MetaObject) (string, string, bool) {
	token, hasResumeToken := obj.GetAnnotations()[reconcilers.PollerResumeTokenAnnotation]
	id, hasResumeID := obj.GetAnnotations()[reconcilers.PollerResumeIDAnnotation]

	return id, token, hasResumeToken && hasResumeID
}

func SetPollerResumeToken(obj genruntime.MetaObject, id string, token string) {
	genruntime.AddAnnotation(obj, reconcilers.PollerResumeTokenAnnotation, token)
	genruntime.AddAnnotation(obj, reconcilers.PollerResumeIDAnnotation, id)
}

// ClearPollerResumeToken clears the poller resume token and ID annotations
func ClearPollerResumeToken(obj genruntime.MetaObject) {
	genruntime.RemoveAnnotation(obj, reconcilers.PollerResumeTokenAnnotation)
	genruntime.RemoveAnnotation(obj, reconcilers.PollerResumeIDAnnotation)
}

func SetLatestReconciledGeneration(obj genruntime.MetaObject) {
	genruntime.AddAnnotation(obj, reconcilers.LatestReconciledGeneration, strconv.FormatInt(obj.GetGeneration(), 10))
}

func GetLatestReconciledGeneration(obj genruntime.MetaObject) (int64, bool) {
	val, hasGeneration := obj.GetAnnotations()[reconcilers.LatestReconciledGeneration]
	gen, err := strconv.Atoi(val)
	if err != nil {
		return 0, false
	}
	return int64(gen), hasGeneration
}
