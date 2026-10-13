/*
Copyright 2021 The Kubernetes Authors.

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

package ec2

import "errors"

var (
	// ErrInstanceNotFoundByID defines an error for when the instance with the provided provider ID is missing.
	ErrInstanceNotFoundByID = errors.New("failed to find instance by id")

	// ErrDescribeInstance defines an error for when AWS SDK returns error when describing instances.
	ErrDescribeInstance = errors.New("failed to describe instance by id")
)
