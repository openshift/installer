package awsbase

import (
	"fmt"
	"net/http"
	"net/url"

	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/hashicorp/aws-sdk-go-base/v2/internal/config"
)

func defaultHttpClient(c *config.Config) (*awshttp.BuildableClient, error) {
	var err error

	httpClient := awshttp.NewBuildableClient().
		WithTransportOptions(func(tr *http.Transport) {
			if c.Insecure {
				tlsConfig := tr.TLSClientConfig
				tlsConfig.InsecureSkipVerify = true
			}
			if c.HTTPProxy != "" {
				var proxyUrl *url.URL
				proxyUrl, parseErr := url.Parse(c.HTTPProxy)
				if parseErr != nil {
					err = fmt.Errorf("error parsing HTTP proxy URL: %w", parseErr)
				}
				tr.Proxy = http.ProxyURL(proxyUrl)
			}
		})

	return httpClient, err
}
