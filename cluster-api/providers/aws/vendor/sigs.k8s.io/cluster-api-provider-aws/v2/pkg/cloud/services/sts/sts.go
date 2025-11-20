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

// Package sts provides an interface for AWS STS operations using the AWS SDK v2.
package sts

import (
	"context"

	signerv4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// STSClient interface for STS operations using AWS SDK v2.
type STSClient interface {
	GetCallerIdentity(ctx context.Context, params *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error)
	PresignGetCallerIdentity(ctx context.Context, params *sts.GetCallerIdentityInput, optFns ...func(*sts.PresignOptions)) (*signerv4.PresignedHTTPRequest, error)
	AssumeRole(ctx context.Context, params *sts.AssumeRoleInput, optFns ...func(*sts.Options)) (*sts.AssumeRoleOutput, error)
}

// ClientWrapper wraps both the regular STS client and presign client to implement STSClient interface.
type ClientWrapper struct {
	client        *sts.Client
	presignClient *sts.PresignClient
}

// NewClientWrapper creates a new STS client wrapper.
func NewClientWrapper(client *sts.Client) *ClientWrapper {
	return &ClientWrapper{
		client:        client,
		presignClient: sts.NewPresignClient(client),
	}
}

// GetCallerIdentity calls the regular STS GetCallerIdentity operation.
func (c *ClientWrapper) GetCallerIdentity(ctx context.Context, params *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error) {
	return c.client.GetCallerIdentity(ctx, params, optFns...)
}

// PresignGetCallerIdentity creates a presigned URL for the GetCallerIdentity operation.
func (c *ClientWrapper) PresignGetCallerIdentity(ctx context.Context, params *sts.GetCallerIdentityInput, optFns ...func(*sts.PresignOptions)) (*signerv4.PresignedHTTPRequest, error) {
	return c.presignClient.PresignGetCallerIdentity(ctx, params, optFns...)
}

// AssumeRole calls the STS AssumeRole operation.
func (c *ClientWrapper) AssumeRole(ctx context.Context, params *sts.AssumeRoleInput, optFns ...func(*sts.Options)) (*sts.AssumeRoleOutput, error) {
	return c.client.AssumeRole(ctx, params, optFns...)
}

// Ensure our wrapper implements the STSClient interface.
var _ STSClient = (*ClientWrapper)(nil)
