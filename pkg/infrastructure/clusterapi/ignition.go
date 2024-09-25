package clusterapi

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	configv1 "github.com/openshift/api/config/v1"
	"sigs.k8s.io/yaml"

	"github.com/openshift/installer/cmd/openshift-install/command"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/lbconfig"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

const (
	infrastructureFilepath = "/opt/openshift/manifests/cluster-infrastructure-02-config.yml"

	// replaceable is the string that precedes the encoded data in the ignition data.
	// The data must be replaced before decoding the string, and the string must be
	// prepended to the encoded data.
	replaceable = "data:text/plain;charset=utf-8;base64,"
)

// EditIgnition attempts to edit the contents of the bootstrap ignition when the user has selected
// a custom DNS configuration. Find the public and private load balancer addresses and fill in the
// infrastructure file within the ignition struct.
func EditIgnition(ctx context.Context, in IgnitionInput, platform string, publicIPAddresses, privateIPAddresses []string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*2)
	defer cancel()

	ignData := &igntypes.Config{}
	err := json.Unmarshal(in.BootstrapIgnData, ignData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal bootstrap ignition: %w", err)
	}

	err = AddLoadBalancersToInfra(platform, ignData, publicIPAddresses, privateIPAddresses)
	if err != nil {
		return nil, fmt.Errorf("failed to add load balancers to ignition config: %w", err)
	}

	publicIPsSingleStr := strings.Join(publicIPAddresses, ",")
	privateIPsSingleStr := strings.Join(privateIPAddresses, ",")

	lbConfig, err := lbconfig.GenerateLBConfigOverride(privateIPsSingleStr, publicIPsSingleStr)
	if err != nil {
		return nil, err
	}
	if err := asset.NewDefaultFileWriter(lbConfig).PersistToFile(command.RootOpts.Dir); err != nil {
		return nil, fmt.Errorf("failed to save %s to state file: %w", lbConfig.Name(), err)
	}

	editedIgnBytes, err := json.Marshal(ignData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert ignition data to json: %w", err)
	}

	return editedIgnBytes, nil
}

// AddLoadBalancersToInfra will load the public and private load balancer information into
// the infrastructure CR. This will occur after the data has already been inserted into the
// ignition file.
func AddLoadBalancersToInfra(platform string, config *igntypes.Config, publicLBs []string, privateLBs []string) error {
	for i, fileData := range config.Storage.Files {
		// update the contents of this file
		if fileData.Path == infrastructureFilepath {
			contents := config.Storage.Files[i].Contents.Source
			replaced := strings.Replace(*contents, replaceable, "", 1)

			rawDecodedText, err := base64.StdEncoding.DecodeString(replaced)
			if err != nil {
				return fmt.Errorf("failed to decode contents of ignition file: %w", err)
			}

			infra := &configv1.Infrastructure{}
			if err := yaml.Unmarshal(rawDecodedText, infra); err != nil {
				return fmt.Errorf("failed to unmarshal infrastructure: %w", err)
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
			case gcptypes.Name:
				if infra.Status.PlatformStatus.GCP.CloudLoadBalancerConfig.DNSType == configv1.ClusterHostedDNSType {
					infra.Status.PlatformStatus.GCP.CloudLoadBalancerConfig.ClusterHosted = &cloudLBInfo
				}
			case awstypes.Name:
				return fmt.Errorf("%s platform is incomplete", platform)
			default:
				return fmt.Errorf("invalid platform %s", platform)
			}

			// convert the infrastructure back to an encoded string
			infraContents, err := yaml.Marshal(infra)
			if err != nil {
				return fmt.Errorf("failed to marshal infrastructure: %w", err)
			}

			encoded := fmt.Sprintf("%s%s", replaceable, base64.StdEncoding.EncodeToString(infraContents))
			// replace the contents with the edited information
			config.Storage.Files[i].Contents.Source = &encoded

			break
		}
	}

	return nil
}
