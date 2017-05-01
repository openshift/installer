package server

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptrace"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/coreos/tectonic-installer/installer/server/aws/cloudforms"
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

// EtcdHealth returns the ServiceStatus of the given Node's etcd instance.
func EtcdHealth(client *http.Client, endpoint string) ServiceStatus {
	status := ServiceStatus{
		Instance: endpoint,
		Ready:    false,
	}

	// etcd health endpoint must return {"health": "true"}
	resp, err := client.Get(fmt.Sprintf("http://%s:%d/health", endpoint, 2379))
	if err != nil {
		status.Message = err.Error()
		return status
	}
	defer resp.Body.Close()

	type health struct {
		Health string `json:"health"`
	}
	etcd := new(health)
	err = json.NewDecoder(resp.Body).Decode(etcd)
	if err != nil {
		status.Message = err.Error()
		return status
	}

	if etcd.Health == "true" {
		status.Ready = true
	}
	return status
}

// KubeletHealth returns the ServiceStatus of the Kubelet via read-only port.
func KubeletHealth(client *http.Client, endpoint string) ServiceStatus {
	status := ServiceStatus{
		Instance: endpoint,
		Ready:    false,
	}

	// kubelet read-only port
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/healthz", endpoint, 10255), nil)
	// TODO: (ggreer) check the health of *all* IPs that endpoint resolves to instead of a random IP
	trace := &httptrace.ClientTrace{
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			for _, addr := range dnsInfo.Addrs {
				status.Addrs = append(status.Addrs, addr.String())
			}
			sort.Strings(status.Addrs)
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			status.RemoteAddr, _, _ = net.SplitHostPort(connInfo.Conn.RemoteAddr().String())
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	resp, err := client.Do(req)
	if err != nil {
		status.Message = err.Error()
		return status
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		status.Ready = true
	}
	return status
}

// TectonicConsoleHealth returns the ServiceStatus of the Tectonic Console.
func TectonicConsoleHealth(client *http.Client, endpoint string) ServiceStatus {
	status := ServiceStatus{
		Instance: endpoint,
		Ready:    false,
	}

	// Tectonic Console health endpoint must return {"status": "ok"}
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

// cfResourceStatus gives the status of the Cloud Formation resources
type cfResourceStatus struct {
	LogicalResourceID    string `json:"logicalResourceID"`
	PhysicalResourceID   string `json:"physicalResourceID"`
	ResourceType         string `json:"resourceType"`
	ResourceStatus       string `json:"resourceStatus"`
	ResourceStatusReason string `json:"resourceStatusReason"`
}

// cloudFormationStatus gives the status of the AWS Cloud Formation
type cloudFormationStatus struct {
	Resources []cfResourceStatus `json:"resources"`
	Instance  string             `json:"instance"`
	ID        string             `json:"id"`
	Message   string             `json:"message"`
	Ready     bool               `json:"ready"`
	Error     bool               `json:"error"`
}

func toStr(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return string(*ptr)
}

// CloudFormationHealth returns an indicator of CloudFormation status.
func cloudFormationHealth(cluster *cloudforms.Cluster, sess *session.Session) cloudFormationStatus {
	status := cloudFormationStatus{
		Instance:  cluster.ControllerDomain,
		Resources: []cfResourceStatus{},
		Ready:     false,
		ID:        "",
	}
	cfStatus, err := cluster.Status(sess)
	if err != nil {
		status.Message = err.Error()
		status.Error = true
		return status
	}

	for _, resource := range cfStatus.Resources {
		r := cfResourceStatus{
			LogicalResourceID:    toStr(resource.LogicalResourceId),
			PhysicalResourceID:   toStr(resource.PhysicalResourceId),
			ResourceType:         toStr(resource.ResourceType),
			ResourceStatus:       toStr(resource.ResourceStatus),
			ResourceStatusReason: toStr(resource.ResourceStatusReason),
		}
		status.Resources = append(status.Resources, r)
	}
	status.Message = cfStatus.StatusString
	status.Ready = cfStatus.Ready
	status.Error = cfStatus.Error
	status.ID = cfStatus.ID
	return status
}
