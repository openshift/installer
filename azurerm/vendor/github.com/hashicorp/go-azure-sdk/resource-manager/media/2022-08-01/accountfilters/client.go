package accountfilters

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountFiltersClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAccountFiltersClientWithBaseURI(endpoint string) AccountFiltersClient {
	return AccountFiltersClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
