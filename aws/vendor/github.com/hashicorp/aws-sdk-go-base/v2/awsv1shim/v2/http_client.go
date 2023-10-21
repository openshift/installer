// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package awsv1shim

import (
	"net/http"

	"github.com/hashicorp/aws-sdk-go-base/v2/internal/config"
	"github.com/hashicorp/go-cleanhttp"
)

func defaultHttpClient(c *config.Config) (*http.Client, error) {
	opts, err := c.HTTPTransportOptions()
	if err != nil {
		return nil, err
	}

	httpClient := cleanhttp.DefaultPooledClient()
	opts(httpClient.Transport.(*http.Transport))

	return httpClient, nil
}
