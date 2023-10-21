package workbooktemplatesapis

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkbookTemplatesAPIsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewWorkbookTemplatesAPIsClientWithBaseURI(endpoint string) WorkbookTemplatesAPIsClient {
	return WorkbookTemplatesAPIsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
