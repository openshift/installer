package gcp

import (
	"bytes"
	"fmt"

	"github.com/pkg/errors"
	ini "gopkg.in/ini.v1"
)

// https://github.com/kubernetes/kubernetes/blob/368ee4bb8ee7a0c18431cd87ee49f0c890aa53e5/staging/src/k8s.io/legacy-cloud-providers/gce/gce.go#L188
type config struct {
	Global global `ini:"global"`
}

type global struct {
	ProjectID string `ini:"project-id"`

	Regional  bool `ini:"regional"`
	Multizone bool `ini:"multizone"`

	NodeTags []string `ini:"node-tags"`

	SubnetworkName string `ini:"subnetwork-name"`
}

// CloudProviderConfig generates the cloud provider config for the GCP platform.
func CloudProviderConfig(infraID, projectID string) (string, error) {
	file := ini.Empty()
	config := &config{
		Global: global{
			ProjectID: projectID,

			// To make sure k8s cloud provider is looking for instances in all zones.
			Regional:  true,
			Multizone: true,

			// To make sure k8s cloud provide has tags for firewal for load balancer.
			NodeTags: []string{fmt.Sprintf("%s-worker", infraID)},

			// Used for internal load balancers
			SubnetworkName: fmt.Sprintf("%s-worker-subnet", infraID),
		},
	}
	if err := file.ReflectFrom(config); err != nil {
		return "", errors.Wrap(err, "failed to reflect from config")
	}
	buf := &bytes.Buffer{}
	if _, err := file.WriteTo(buf); err != nil {
		return "", errors.Wrap(err, "failed to write out cloud provider config")
	}
	return buf.String(), nil
}
