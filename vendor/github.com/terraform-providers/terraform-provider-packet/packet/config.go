package packet

import (
	"context"
	"crypto/x509"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-sdk/helper/logging"
	"github.com/packethost/packngo"
)

const (
	consumerToken = "aZ9GmqHTPtxevvFq9SK3Pi2yr9YCbRzduCSXF2SNem5sjB91mDq7Th3ZwTtRqMWZ"
)

type Config struct {
	AuthToken string
}

var redirectsErrorRe = regexp.MustCompile(`stopped after \d+ redirects\z`)

func PacketRetryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	if err != nil {
		if v, ok := err.(*url.Error); ok {
			// Don't retry if the error was due to too many redirects.
			if redirectsErrorRe.MatchString(v.Error()) {
				return false, nil
			}

			// Don't retry if the error was due to TLS cert verification failure.
			if _, ok := v.Err.(x509.UnknownAuthorityError); ok {
				return false, nil
			}
		}

		// The error is likely recoverable so retry.
		return true, nil
	}
	return false, nil
}

// Client returns a new client for accessing Packet's API.
func (c *Config) Client() *packngo.Client {
	httpClient := retryablehttp.NewClient()
	httpClient.RetryWaitMin = time.Second
	httpClient.RetryWaitMax = 30 * time.Second
	httpClient.RetryMax = 10
	httpClient.CheckRetry = PacketRetryPolicy
	httpClient.HTTPClient.Transport = logging.NewTransport(
		"Packet",
		httpClient.HTTPClient.Transport)

	return packngo.NewClientWithAuth(consumerToken, c.AuthToken, httpClient)
}
