package ovirt

import (
	"bytes"
	"encoding/json"
)

// CloudProviderConfig generates the cloud provider config for the oVirt platform.
func CloudProviderConfig(storageDomainID string, clusterID string, networkName string) (string, error) {
	c := config{
		StorageDomainID: storageDomainID,
		ClusterID:       clusterID,
		NetworkName:     networkName,
	}

	buff := &bytes.Buffer{}
	encoder := json.NewEncoder(buff)
	encoder.SetIndent("", "\t")
	if err := encoder.Encode(c); err != nil {
		return "", err
	}
	return buff.String(), nil
}
