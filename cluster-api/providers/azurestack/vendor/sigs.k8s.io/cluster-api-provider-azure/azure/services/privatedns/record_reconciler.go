/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package privatedns

import (
	"context"

	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

func (s *Service) reconcileRecords(ctx context.Context, records []azure.ResourceSpecGetter) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "privatedns.Service.reconcileRecords")
	defer done()

	var resErr error

	// We go through the list of links to reconcile each one, independently of the result of the previous one.
	// If multiple errors occur, we return the most pressing one.
	// Order of precedence (highest -> lowest) is: error that is not an operationNotDoneError (i.e. error creating) -> operationNotDoneError (i.e. creating in progress) -> no error (i.e. created)
	for _, recordSpec := range records {
		if _, err := s.recordReconciler.CreateOrUpdateResource(ctx, recordSpec, serviceName); err != nil {
			if !azure.IsOperationNotDoneError(err) || resErr == nil {
				resErr = err
			}
		}
	}

	return resErr
}
