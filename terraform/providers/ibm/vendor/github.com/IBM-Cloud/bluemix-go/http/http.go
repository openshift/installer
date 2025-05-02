package http

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"time"

	"github.com/IBM-Cloud/bluemix-go"
)

//NewHTTPClient ...
func NewHTTPClient(config *bluemix.Config) *http.Client {
	return &http.Client{
		Transport: makeTransport(config),
		Timeout:   config.HTTPTimeout,
	}
}

func makeTransport(config *bluemix.Config) http.RoundTripper {
	proxyFunc := http.ProxyFromEnvironment
	if config.HTTPClient != nil && config.HTTPClient.Transport != nil {
		if t, ok := config.HTTPClient.Transport.(*http.Transport); ok {
			proxyFunc = t.Proxy
		}
	}
	return NewTraceLoggingTransport(&http.Transport{
		Proxy: proxyFunc,
		Dial: (&net.Dialer{
			Timeout:   50 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 20 * time.Second,
		DisableCompression:  true,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.SSLDisable,
		},
	})
}

//UserAgent ...
func UserAgent() string {
	return fmt.Sprintf("Bluemix-go SDK %s / %s ", bluemix.Version, runtime.GOOS)
}
