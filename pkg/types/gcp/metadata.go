package gcp

import configv1 "github.com/openshift/api/config/v1"

// Metadata contains GCP metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	Region            string                        `json:"region"`
	ProjectID         string                        `json:"projectID"`
	NetworkProjectID  string                        `json:"networkProjectID,omitempty"`
	PrivateZoneDomain string                        `json:"privateZoneDomain,omitempty"`
	ServiceEndpoints  []configv1.GCPServiceEndpoint `json:"serviceEndpoints,omitempty"`
}
