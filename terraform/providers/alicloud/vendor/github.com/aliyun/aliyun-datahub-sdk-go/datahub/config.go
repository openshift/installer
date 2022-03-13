package datahub

import (
    "fmt"
    "math"
    "net/http"
    "runtime"
    "time"
)

type Config struct {
    UserAgent            string
    CompressorType       CompressorType
    EnableBinary         bool
    EnableSchemaRegistry bool
    HttpClient           *http.Client
}

func NewDefaultConfig() *Config {
    return &Config{
        UserAgent:            DefaultUserAgent(),
        CompressorType:       NOCOMPRESS,
        EnableBinary:         true,
        EnableSchemaRegistry: false,
        HttpClient:           DefaultHttpClient(),
    }
}

// DefaultHttpClient returns a default HTTP client with sensible values.
func DefaultHttpClient() *http.Client {
    return &http.Client{
        Transport: &http.Transport{
            DialContext:           TraceDialContext(10 * time.Second),
            Proxy:                 http.ProxyFromEnvironment,
            MaxIdleConns:          math.MaxInt32,
            MaxIdleConnsPerHost:   math.MaxInt32,
            MaxConnsPerHost:       math.MaxInt32,
            IdleConnTimeout:       30 * time.Second,
            TLSHandshakeTimeout:   10 * time.Second,
            ExpectContinueTimeout: 1 * time.Second,
            ResponseHeaderTimeout: 100 * time.Second,
        },
    }
}

// DefaultUserAgent returns a default user agent
func DefaultUserAgent() string {
    return fmt.Sprintf("godatahub/%s golang/%s %s", DATAHUB_SDK_VERSION, runtime.Version(), runtime.GOOS)
}
