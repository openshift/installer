/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package identity

import "github.com/Azure/azure-sdk-for-go/sdk/azidentity"

type TokenCredentialProvider interface {
	NewClientSecretCredential(tenantID string, clientID string, clientSecret string, options *azidentity.ClientSecretCredentialOptions) (*azidentity.ClientSecretCredential, error)
	NewClientCertificateCredential(tenantID, clientID string, clientCertificate, password []byte) (*azidentity.ClientCertificateCredential, error)
	NewManagedIdentityCredential(options *azidentity.ManagedIdentityCredentialOptions) (*azidentity.ManagedIdentityCredential, error)
	NewWorkloadIdentityCredential(options *azidentity.WorkloadIdentityCredentialOptions) (*azidentity.WorkloadIdentityCredential, error)
}

var _ TokenCredentialProvider = &tokenCredentialProvider{}

type tokenCredentialProvider struct{}

func (t *tokenCredentialProvider) NewClientSecretCredential(tenantID string, clientID string, clientSecret string, options *azidentity.ClientSecretCredentialOptions) (*azidentity.ClientSecretCredential, error) {
	return azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, options)
}

func (t *tokenCredentialProvider) NewClientCertificateCredential(tenantID, clientID string, clientCertificate, password []byte) (*azidentity.ClientCertificateCredential, error) {
	return NewClientCertificateCredential(tenantID, clientID, clientCertificate, password)
}

func (t *tokenCredentialProvider) NewManagedIdentityCredential(options *azidentity.ManagedIdentityCredentialOptions) (*azidentity.ManagedIdentityCredential, error) {
	return azidentity.NewManagedIdentityCredential(options)
}

func (t *tokenCredentialProvider) NewWorkloadIdentityCredential(options *azidentity.WorkloadIdentityCredentialOptions) (*azidentity.WorkloadIdentityCredential, error) {
	return azidentity.NewWorkloadIdentityCredential(options)
}

func DefaultTokenCredentialProvider() TokenCredentialProvider {
	return &tokenCredentialProvider{}
}
