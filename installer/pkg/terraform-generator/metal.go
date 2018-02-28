package terraformgenerator

import (
	"github.com/coreos/tectonic-installer/installer/pkg/config"
)

// Metal defines all variables for this platform.
type Metal struct {
	CalicoMTU           string `json:"tectonic_metal_calico_mtu,omitempty"`
	ControllerDomain    string `json:"tectonic_metal_controller_domain,omitempty"`
	ControllerDomains   string `json:"tectonic_metal_controller_domains,omitempty"`
	ControllerMACs      string `json:"tectonic_metal_controller_macs,omitempty"`
	ControllerNames     string `json:"tectonic_metal_controller_names,omitempty"`
	IngressDomain       string `json:"tectonic_metal_ingress_domain,omitempty"`
	MatchboxCA          string `json:"tectonic_metal_matchbox_ca,omitempty"`
	MatchboxClientCert  string `json:"tectonic_metal_matchbox_client_cert,omitempty"`
	MatchboxClientKey   string `json:"tectonic_metal_matchbox_client_key,omitempty"`
	MatchboxHTTPURL     string `json:"tectonic_metal_matchbox_http_url,omitempty"`
	MatchboxRPCEndpoint string `json:"tectonic_metal_matchbox_rpc_endpoint,omitempty"`
	SSHAuthorizedKey    string `json:"tectonic_ssh_authorized_key,omitempty"` // TODO(spangenberg): Prefix with metal.
	WorkerDomains       string `json:"tectonic_metal_worker_domains,omitempty"`
	WorkerMACs          string `json:"tectonic_metal_worker_macs,omitempty"`
	WorkerNames         string `json:"tectonic_metal_worker_names,omitempty"`
}

// NewMetal returns the config for Metal.
func NewMetal(cluster config.Cluster) Metal {
	return Metal{
		CalicoMTU:           cluster.Metal.CalicoMTU,
		ControllerDomain:    cluster.Metal.Controller.Domain,
		ControllerDomains:   cluster.Metal.Controller.Domains,
		ControllerMACs:      cluster.Metal.Controller.MACs,
		ControllerNames:     cluster.Metal.Controller.Names,
		IngressDomain:       cluster.Metal.IngressDomain,
		MatchboxCA:          cluster.Metal.Matchbox.CA,
		MatchboxClientCert:  cluster.Metal.Matchbox.Client.Cert,
		MatchboxClientKey:   cluster.Metal.Matchbox.Client.Key,
		MatchboxHTTPURL:     cluster.Metal.Matchbox.HTTPURL,
		MatchboxRPCEndpoint: cluster.Metal.Matchbox.RPCEndpoint,
		SSHAuthorizedKey:    cluster.Metal.SSHAuthorizedKey,
		WorkerDomains:       cluster.Metal.Worker.Domains,
		WorkerMACs:          cluster.Metal.Worker.MACs,
		WorkerNames:         cluster.Metal.Worker.Names,
	}
}
