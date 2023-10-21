package replicationrecoveryplans

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicationRecoveryPlansClient struct {
	Client  autorest.Client
	baseUri string
}

func NewReplicationRecoveryPlansClientWithBaseURI(endpoint string) ReplicationRecoveryPlansClient {
	return ReplicationRecoveryPlansClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
