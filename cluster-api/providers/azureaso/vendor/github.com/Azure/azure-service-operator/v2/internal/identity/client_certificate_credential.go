/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package identity

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/rotisserie/eris"
)

func NewClientCertificateCredential(
	tenantID string,
	clientID string,
	clientCertificate,
	password []byte,
	options *azidentity.ClientCertificateCredentialOptions,
) (*azidentity.ClientCertificateCredential, error) {
	certs, key, err := azidentity.ParseCertificates(clientCertificate, password)
	if err != nil {
		return nil, eris.Wrapf(err, "failed to parse certificate for '%s'", clientID)
	}

	cred, err := azidentity.NewClientCertificateCredential(
		tenantID,
		clientID,
		certs,
		key,
		options)
	if err != nil {
		return nil, err
	}
	return cred, nil
}
