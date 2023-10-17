/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package genericarmclient

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

var defaultHttpClient *http.Client

func init() {
	// If httpClient is not provided, we use a HTTPClient with default Transport + settings
	// to establish multiple TCP ARMClient connections to avoid throttling.
	// TODO: Use https://github.com/Azure/go-armbalancer here once its prod ready.
	httpTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second, // default transport value
			KeepAlive: 30 * time.Second, // default transport value
		}).DialContext,
		ForceAttemptHTTP2:     false,            // default is true; since HTTP/2 multiplexes a single TCP connection. we'd want to use HTTP/1, which would use multiple TCP connections.
		MaxIdleConns:          100,              // default transport value
		MaxIdleConnsPerHost:   10,               // default is 2, so we want to increase the number to use establish more connections.
		IdleConnTimeout:       90 * time.Second, // default transport value
		TLSHandshakeTimeout:   10 * time.Second, // default transport value
		ExpectContinueTimeout: 1 * time.Second,  // default transport value
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12, // default tls version
		},
	}

	defaultHttpClient = &http.Client{
		Transport: httpTransport,
	}

}
