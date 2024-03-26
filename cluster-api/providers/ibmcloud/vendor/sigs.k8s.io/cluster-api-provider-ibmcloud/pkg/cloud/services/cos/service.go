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
	"net"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/http/httpproxy"

	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/aws/credentials/ibmiam"
	"github.com/IBM/ibm-cos-sdk-go/aws/request"
	cosSession "github.com/IBM/ibm-cos-sdk-go/aws/session"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
)

// iamEndpoint represent the IAM authorisation URL.
const (
	iamEndpoint  = "https://iam.cloud.ibm.com/identity/token"
	cosURLDomain = "cloud-object-storage.appdomain.cloud"
)

// Service holds the IBM Cloud Resource Controller Service specific information.
type Service struct {
	client *s3.S3
}

// ServiceOptions holds the IBM Cloud Resource Controller Service Options specific information.
type ServiceOptions struct {
	*cosSession.Options
}

// GetBucketByName returns a bucket with the given name.
func (s *Service) GetBucketByName(name string) (*s3.HeadBucketOutput, error) {
	input := &s3.HeadBucketInput{
		Bucket: &name,
	}
	return s.client.HeadBucket(input)
}

// CreateBucket creates a new bucket in the COS instance.
func (s *Service) CreateBucket(input *s3.CreateBucketInput) (*s3.CreateBucketOutput, error) {
	return s.client.CreateBucket(input)
}

// CreateBucketWithContext creates a new bucket with an addition ability to pass context.
func (s *Service) CreateBucketWithContext(ctx aws.Context, input *s3.CreateBucketInput, opts ...request.Option) (*s3.CreateBucketOutput, error) {
	return s.client.CreateBucketWithContext(ctx, input, opts...)
}

// PutObject adds an object to a bucket.
func (s *Service) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return s.client.PutObject(input)
}

// GetObjectRequest generates a "aws/request.Request" representing the client's request for the GetObject operation.
func (s *Service) GetObjectRequest(input *s3.GetObjectInput) (*request.Request, *s3.GetObjectOutput) {
	return s.client.GetObjectRequest(input)
}

// ListObjects returns the list of objects in a bucket.
func (s *Service) ListObjects(input *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
	return s.client.ListObjects(input)
}

// DeleteObject deletes a object in a bucket.
func (s *Service) DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	return s.client.DeleteObject(input)
}

// PutPublicAccessBlock creates or modifies the PublicAccessBlock configuration for a bucket.
func (s *Service) PutPublicAccessBlock(input *s3.PutPublicAccessBlockInput) (*s3.PutPublicAccessBlockOutput, error) {
	return s.client.PutPublicAccessBlock(input)
}

// NewService returns a new service for the IBM Cloud Resource Controller api client.
func NewService(options ServiceOptions, apikey, serviceInstance string) (*Service, error) {
	if options.Options == nil {
		options.Options = &cosSession.Options{}
	}
	options.Config.S3ForcePathStyle = aws.Bool(true)
	options.Config.HTTPClient = &http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return httpproxy.FromEnvironment().ProxyFunc()(req.URL)
			},
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	options.Config.Credentials = ibmiam.NewStaticCredentials(aws.NewConfig(), iamEndpoint, apikey, serviceInstance)

	sess, err := cosSession.NewSessionWithOptions(*options.Options)
	if err != nil {
		return nil, err
	}
	return &Service{
		client: s3.New(sess),
	}, nil
}
