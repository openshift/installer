/*
 *
<<<<<<<< HEAD:cluster-api/providers/ibmcloud/vendor/google.golang.org/grpc/grpclog/internal/grpclog.go
 * Copyright 2024 gRPC authors.
========
 * Copyright 2025 gRPC authors.
>>>>>>>> 8a4761186a (1. Update the base go.mod file):vendor/google.golang.org/grpc/encoding/internal/internal.go
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

<<<<<<<< HEAD:cluster-api/providers/ibmcloud/vendor/google.golang.org/grpc/grpclog/internal/grpclog.go
// Package internal contains functionality internal to the grpclog package.
package internal

// LoggerV2Impl is the logger used for the non-depth log functions.
var LoggerV2Impl LoggerV2

// DepthLoggerV2Impl is the logger used for the depth log functions.
var DepthLoggerV2Impl DepthLoggerV2
========
// Package internal contains code internal to the encoding package.
package internal

// RegisterCompressorForTesting registers a compressor in the global compressor
// registry. It returns a cleanup function that should be called at the end
// of the test to unregister the compressor.
//
// This prevents compressors registered in one test from appearing in the
// encoding headers of subsequent tests.
var RegisterCompressorForTesting any // func RegisterCompressor(c Compressor) func()
>>>>>>>> 8a4761186a (1. Update the base go.mod file):vendor/google.golang.org/grpc/encoding/internal/internal.go
