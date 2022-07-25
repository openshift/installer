package awsv1shim

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"

	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/hashicorp/aws-sdk-go-base/v2/internal/config"
	"github.com/hashicorp/go-cleanhttp"
)

func defaultHttpClient(c *config.Config) (*http.Client, error) {
	httpClient := cleanhttp.DefaultPooledClient()
	transport := httpClient.Transport.(*http.Transport)

	transport.MaxIdleConnsPerHost = awshttp.DefaultHTTPTransportMaxIdleConnsPerHost

	tlsConfig := transport.TLSClientConfig
	if tlsConfig == nil {
		tlsConfig = &tls.Config{}
		transport.TLSClientConfig = tlsConfig
	}
	tlsConfig.MinVersion = tls.VersionTLS12

	if c.Insecure {
		tlsConfig.InsecureSkipVerify = true
	}

	if c.HTTPProxy != "" {
		proxyUrl, err := url.Parse(c.HTTPProxy)
		if err != nil {
			return nil, fmt.Errorf("error parsing HTTP proxy URL: %w", err)
		}

		transport.Proxy = http.ProxyURL(proxyUrl)
	}

	return httpClient, nil
}
