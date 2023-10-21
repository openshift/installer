package analysisservices

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AnalysisServicesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAnalysisServicesClientWithBaseURI(endpoint string) AnalysisServicesClient {
	return AnalysisServicesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
