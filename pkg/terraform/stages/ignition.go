package stages

import (
	"encoding/base64"
	"fmt"
	"strings"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"sigs.k8s.io/yaml"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types/gcp"
)

const (
	// replaceable is the string that precedes the encoded data in the ignition data.
	// The data must be replaced before decoding the string, and the string must be
	// prepended to the encoded data.
	replaceable = "data:text/plain;charset=utf-8;base64,"
)

// AddLoadBalancersToInfra will load the public and private load balancer information into
// the infrastructure CR. This will occur after the data has already been inserted into the
// ignition file.
func AddLoadBalancersToInfra(platform string, config *igntypes.Config, publicLBs []string, privateLBs []string) error {
	index := -1
	for i, fileData := range config.Storage.Files {
		// update the contents of this file
		if fileData.Path == "/opt/openshift/manifests/cluster-infrastructure-02-config.yml" {
			index = i
			break
		}
	}

	if index >= 0 {
		contents := config.Storage.Files[index].Contents.Source
		replaced := strings.Replace(*contents, replaceable, "", 1)

		rawDecodedText, err := base64.StdEncoding.DecodeString(replaced)
		if err != nil {
			return err
		}

		infra := &configv1.Infrastructure{}
		if err := yaml.Unmarshal(rawDecodedText, infra); err != nil {
			return err
		}

		// convert the list of strings to a list of IPs
		apiIntLbs := []configv1.IP{}
		for _, ip := range privateLBs {
			apiIntLbs = append(apiIntLbs, configv1.IP(ip))
		}
		apiLbs := []configv1.IP{}
		for _, ip := range publicLBs {
			apiLbs = append(apiLbs, configv1.IP(ip))
		}
		cloudLBInfo := configv1.CloudLoadBalancerIPs{
			APIIntLoadBalancerIPs: apiIntLbs,
			APILoadBalancerIPs:    apiLbs,
		}

		switch platform {
		case gcp.Name:
			if infra.Status.PlatformStatus.GCP.CloudLoadBalancerConfig.DNSType == configv1.ClusterHostedDNSType {
				infra.Status.PlatformStatus.GCP.CloudLoadBalancerConfig.ClusterHosted = &cloudLBInfo
			}
		default:
			return fmt.Errorf("failed to set load balancer info for platform %s", platform)
		}

		// convert the infrastructure back to an encoded string
		infraContents, err := yaml.Marshal(infra)
		if err != nil {
			return err
		}

		encoded := fmt.Sprintf("%s%s", replaceable, base64.StdEncoding.EncodeToString(infraContents))
		config.Storage.Files[index].Contents.Source = &encoded
	}

	return nil
}
