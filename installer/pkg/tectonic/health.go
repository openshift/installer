package tectonic

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// StatusChecker is a client for checking the status of cluster components.
type StatusChecker interface {
	// Status returns the status of cluster components and services.
	Status() ([]byte, error)
}

// ServiceStatus gives the status of an instance of a service.
type ServiceStatus struct {
	Instance   string   `json:"instance"`
	Message    string   `json:"message"`
	Ready      bool     `json:"ready"`
	RemoteAddr string   `json:"remoteAddr"`
	Addrs      []string `json:"addrs"`
}

// defaultStatusClient respects proxies, sets reasonable timeouts, and allows
// checking the status of services running with self-signed certificates.
var defaultStatusClient = &http.Client{
	Timeout: time.Duration(10 * time.Second),
	Transport: &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}

// ConsoleHealth returns the ServiceStatus of the Tectonic Console.
func ConsoleHealth(client *http.Client, endpoint string) ServiceStatus {
	status := ServiceStatus{
		Instance: endpoint,
		Ready:    false,
	}

	// Tectonic Console health endpoint must return {"status": "ok"}
	if client == nil {
		client = defaultStatusClient
	}
	resp, err := client.Get(fmt.Sprintf("https://%s/health", endpoint))
	if err != nil {
		status.Message = err.Error()
		return status
	}
	defer resp.Body.Close()

	type health struct {
		Status string `json:"status"`
	}
	h := new(health)
	err = json.NewDecoder(resp.Body).Decode(h)
	if err != nil {
		status.Message = err.Error()
		return status
	}

	if h.Status == "ok" {
		status.Ready = true
	}
	return status
}
