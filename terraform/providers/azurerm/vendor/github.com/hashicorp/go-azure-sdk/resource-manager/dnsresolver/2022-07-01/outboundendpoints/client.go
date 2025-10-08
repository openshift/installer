package outboundendpoints

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OutboundEndpointsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewOutboundEndpointsClientWithBaseURI(endpoint string) OutboundEndpointsClient {
	return OutboundEndpointsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
