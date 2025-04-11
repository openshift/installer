/*
Copyright 2023 The Kubernetes Authors.

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

// Run go generate to regenerate this mock.
//
//go:generate ../hack/tools/bin/mockgen -destination nodelister_mock.go -package controllers -source azuremanagedmachinepool_reconciler.go Nodelister
//go:generate /usr/bin/env bash -c "cat ../hack/boilerplate/boilerplate.generatego.txt nodelister_mock.go > _nodelister_mock.go && mv _nodelister_mock.go nodelister_mock.go"
package controllers
