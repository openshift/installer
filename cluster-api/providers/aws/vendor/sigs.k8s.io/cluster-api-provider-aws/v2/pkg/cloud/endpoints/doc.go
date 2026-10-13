/*
Copyright 2025 The Kubernetes Authors.

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

/*
This file downloads the partition.json from AWS SDK GO V2 github and
generates Partition and data for use in the project.
*/

// Package endpoints provides AWS Service Endpoints, Region, Partition related methods.
package endpoints

//go:generate /usr/bin/env bash -c "grep 'github.com/aws/aws-sdk-go-v2 ' ../../../go.mod | cut -d' ' -f2 > .awsSDKversion && curl -ssLO https://raw.githubusercontent.com/aws/aws-sdk-go-v2/refs/tags/$(cat .awsSDKversion)/internal/endpoints/awsrulesfn/partitions.json && rm .awsSDKversion"
//go:generate go run codegen.go -model partitions.json -output partitions.go
//go:generate /usr/bin/env bash -c "cat ../../../hack/boilerplate/boilerplate.generatego.txt partitions.go > _partitions.go && mv _partitions.go partitions.go && rm partitions.json"
