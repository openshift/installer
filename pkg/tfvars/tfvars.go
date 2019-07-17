// Package tfvars generates Terraform variables for launching the cluster.
package tfvars

import (
	"encoding/json"
	"net"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type config struct {
	ClusterID     string `json:"cluster_id,omitempty"`
	ClusterDomain string `json:"cluster_domain,omitempty"`
	BaseDomain    string `json:"base_domain,omitempty"`
	MachineCIDR   string `json:"machine_cidr"`
	Masters       int    `json:"master_count,omitempty"`
	EtcdPivot     string `json:"etcd_pivot,omitempty"`

	IgnitionBootstrap string `json:"ignition_bootstrap,omitempty"`
	IgnitionMaster    string `json:"ignition_master,omitempty"`
}

// TFVars generates terraform.tfvar JSON for launching the cluster.
func TFVars(clusterID string, clusterDomain string, baseDomain string, machineCIDR *net.IPNet, bootstrapIgn string, masterIgn string, masterCount int) ([]byte, error) {
	etcdPivot := "false"
	if pivot, ok := os.LookupEnv("OPENSHIFT_INSTALL_ETCD_ENABLE_PIVOT"); ok && pivot != "" {
		logrus.Warn("Found override for etcd provisioning style. Please be warned, this is still under development")
		etcdPivot = pivot
	}
	config := &config{
		ClusterID:         clusterID,
		ClusterDomain:     strings.TrimSuffix(clusterDomain, "."),
		BaseDomain:        strings.TrimSuffix(baseDomain, "."),
		MachineCIDR:       machineCIDR.String(),
		Masters:           masterCount,
		IgnitionBootstrap: bootstrapIgn,
		IgnitionMaster:    masterIgn,
		EtcdPivot:         etcdPivot,
	}

	return json.MarshalIndent(config, "", "  ")
}
