package replicationpolicies

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicationPoliciesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewReplicationPoliciesClientWithBaseURI(endpoint string) ReplicationPoliciesClient {
	return ReplicationPoliciesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
