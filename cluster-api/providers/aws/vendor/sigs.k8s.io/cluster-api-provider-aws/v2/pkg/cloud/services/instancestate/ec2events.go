/*
Copyright 2020 The Kubernetes Authors.

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

package instancestate

import "context"

// ReconcileEC2Events will reconcile a Service's EC2 events.
func (s Service) ReconcileEC2Events(ctx context.Context) error {
	if err := s.reconcileSQSQueue(ctx); err != nil {
		return err
	}

	return s.reconcileRules(ctx)
}

// DeleteEC2Events will delete a Service's EC2 events.
func (s Service) DeleteEC2Events(ctx context.Context) error {
	if err := s.deleteRules(ctx); err != nil {
		return err
	}

	return s.deleteSQSQueue(ctx)
}
