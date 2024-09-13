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

package cos

import (
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/aws/request"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
)

//go:generate ../../../../hack/tools/bin/mockgen -source=./cos.go -destination=./mock/cos_generated.go -package=mock
//go:generate /usr/bin/env bash -c "cat ../../../../hack/boilerplate/boilerplate.generatego.txt ./mock/cos_generated.go > ./mock/_cos_generated.go && mv ./mock/_cos_generated.go ./mock/cos_generated.go"

// Cos interface defines a method that a IBMCLOUD service object should implement in order to
// use the cos package for listing resource instances.
type Cos interface {
	GetBucketByName(name string) (*s3.HeadBucketOutput, error)
	CreateBucket(input *s3.CreateBucketInput) (*s3.CreateBucketOutput, error)
	CreateBucketWithContext(ctx aws.Context, input *s3.CreateBucketInput, opts ...request.Option) (*s3.CreateBucketOutput, error)
	PutObject(*s3.PutObjectInput) (*s3.PutObjectOutput, error)
	GetObjectRequest(*s3.GetObjectInput) (*request.Request, *s3.GetObjectOutput)
	ListObjects(input *s3.ListObjectsInput) (*s3.ListObjectsOutput, error)
	DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error)
	PutPublicAccessBlock(input *s3.PutPublicAccessBlockInput) (*s3.PutPublicAccessBlockOutput, error)
}
