// Package tfvars generates Terraform variables for launching the cluster.
package tfvars

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// Config contains the cluster data for terraform.
type Config struct {
	ClusterID          string   `json:"cluster_id,omitempty"`
	ClusterDomain      string   `json:"cluster_domain,omitempty"`
	BaseDomain         string   `json:"base_domain,omitempty"`
	Masters            int      `json:"master_count,omitempty"`
	MastersSchedulable bool     `json:"masters_schedulable,omitempty"`
	MachineV4CIDRs     []string `json:"machine_v4_cidrs"`
	MachineV6CIDRs     []string `json:"machine_v6_cidrs"`

	UseIPv4 bool `json:"use_ipv4"`
	UseIPv6 bool `json:"use_ipv6"`

	IgnitionBootstrap     string `json:"ignition_bootstrap,omitempty"`
	IgnitionBootstrapFile string `json:"ignition_bootstrap_file,omitempty"`
	IgnitionMaster        string `json:"ignition_master,omitempty"`
}

// TFVars generates terraform.tfvar JSON for launching the cluster.
func TFVars(clusterID string, clusterDomain string, baseDomain string, machineV4CIDRs []string, machineV6CIDRs []string, useIPv4, useIPv6 bool, bootstrapIgn string, bootstrapIgnSize int64, masterIgn string, masterCount int, mastersSchedulable bool) ([]byte, error) {
	f, err := os.CreateTemp("", "openshift-install-bootstrap-*.ign")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create tmp file for bootstrap ignition")
	}
	defer f.Close()

	// In azure, if the storage account is encrypted, page blob is used instead of block blob due to lack of support.
	// Page blobs file size must be a multiple of 512, hence a little padding is needed to push the file.
	// Finding the nearest size divisible by 512 and adding that padding to the file.
	// Since the file is json type, padding at the end results in json parsing error at bootstrap ignition.
	// Adding the paddding to just before the last } in the json file to bypass the parsing error.
	padding := bootstrapIgnSize - int64(len(bootstrapIgn))
	bootstrapIgn = bootstrapIgn[0:len(bootstrapIgn)-1] + strings.Repeat(" ", int(padding)) + string(bootstrapIgn[len(bootstrapIgn)-1])

	if _, err := f.WriteString(bootstrapIgn); err != nil {
		return nil, errors.Wrap(err, "failed to write bootstrap ignition")
	}

	config := &Config{
		ClusterID:             clusterID,
		ClusterDomain:         strings.TrimSuffix(clusterDomain, "."),
		BaseDomain:            strings.TrimSuffix(baseDomain, "."),
		MachineV4CIDRs:        machineV4CIDRs,
		MachineV6CIDRs:        machineV6CIDRs,
		UseIPv4:               useIPv4,
		UseIPv6:               useIPv6,
		Masters:               masterCount,
		MastersSchedulable:    mastersSchedulable,
		IgnitionBootstrap:     bootstrapIgn,
		IgnitionBootstrapFile: f.Name(),
		IgnitionMaster:        masterIgn,
	}

	return json.MarshalIndent(config, "", "  ")
}
