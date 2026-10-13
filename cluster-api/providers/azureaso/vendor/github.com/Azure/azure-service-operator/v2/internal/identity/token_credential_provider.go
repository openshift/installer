/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package identity

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/msi-dataplane/pkg/dataplane"
)

type TokenCredentialProvider interface {
	NewClientSecretCredential(
		tenantID string,
		clientID string,
		clientSecret string,
		options *azidentity.ClientSecretCredentialOptions,
	) (*azidentity.ClientSecretCredential, error)

	NewClientCertificateCredential(
		tenantID string,
		clientID string,
		clientCertificate,
		password []byte,
		options *azidentity.ClientCertificateCredentialOptions,
	) (*azidentity.ClientCertificateCredential, error)

	NewManagedIdentityCredential(options *azidentity.ManagedIdentityCredentialOptions) (*azidentity.ManagedIdentityCredential, error)

	NewWorkloadIdentityCredential(options *azidentity.WorkloadIdentityCredentialOptions) (*azidentity.WorkloadIdentityCredential, error)
	NewUserAssignedIdentityCredentials(ctx context.Context, credentialPath string, opts ...dataplane.Option) (azcore.TokenCredential, error)
}

var _ TokenCredentialProvider = &tokenCredentialProvider{}

type tokenCredentialProvider struct{}

func (t *tokenCredentialProvider) NewClientSecretCredential(
	tenantID string,
	clientID string,
	clientSecret string,
	options *azidentity.ClientSecretCredentialOptions,
) (*azidentity.ClientSecretCredential, error) {
	return azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, options)
}

func (t *tokenCredentialProvider) NewClientCertificateCredential(
	tenantID string,
	clientID string,
	clientCertificate []byte,
	password []byte,
	options *azidentity.ClientCertificateCredentialOptions,
) (*azidentity.ClientCertificateCredential, error) {
	return NewClientCertificateCredential(tenantID, clientID, clientCertificate, password, options)
}

func (t *tokenCredentialProvider) NewManagedIdentityCredential(options *azidentity.ManagedIdentityCredentialOptions) (*azidentity.ManagedIdentityCredential, error) {
	return azidentity.NewManagedIdentityCredential(options)
}

func (t *tokenCredentialProvider) NewWorkloadIdentityCredential(options *azidentity.WorkloadIdentityCredentialOptions) (*azidentity.WorkloadIdentityCredential, error) {
	return azidentity.NewWorkloadIdentityCredential(options)
}

func (t *tokenCredentialProvider) NewUserAssignedIdentityCredentials(ctx context.Context, credentialPath string, opts ...dataplane.Option) (azcore.TokenCredential, error) {
	return dataplane.NewUserAssignedIdentityCredential(ctx, credentialPath, opts...)
}

func DefaultTokenCredentialProvider() TokenCredentialProvider {
	return &tokenCredentialProvider{}
}
