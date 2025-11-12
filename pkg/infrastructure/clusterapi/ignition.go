package clusterapi

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strings"

	ignutil "github.com/coreos/ignition/v2/config/util"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/cmd/openshift-install/command"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/lbconfig"
	"github.com/openshift/installer/pkg/asset/machines"
	"github.com/openshift/installer/pkg/asset/tls"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

const (
	infrastructureFilepath = "/opt/openshift/manifests/cluster-infrastructure-02-config.yml"
	mcsCertKeyFilepath     = "/opt/openshift/manifests/machine-config-server-tls-secret.yaml"
	mcsKeyFile             = "/opt/openshift/tls/machine-config-server.key"
	mcsCertFile            = "/opt/openshift/tls/machine-config-server.crt"
	masterUserDataFile     = "/opt/openshift/openshift/99_openshift-cluster-api_master-user-data-secret.yaml"
	workerUserDataFile     = "/opt/openshift/openshift/99_openshift-cluster-api_worker-user-data-secret.yaml"

	// header is the string that precedes the encoded data in the ignition data.
	// The data must be replaced before decoding the string, and the string must be
	// prepended to the encoded data.
	header = "data:text/plain;charset=utf-8;base64,"

	masterRole = "master"
	workerRole = "worker"
)

// EditIgnition attempts to edit the contents of the bootstrap ignition when the user has selected
// a custom DNS configuration. Find the public and private load balancer addresses and fill in the
// infrastructure file within the ignition struct.
func EditIgnition(in IgnitionInput, platform string, publicIPAddresses, privateIPAddresses []string) (*IgnitionOutput, error) {
	ignData := &igntypes.Config{}
	err := json.Unmarshal(in.BootstrapIgnData, ignData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal bootstrap ignition: %w", err)
	}

	err = addLoadBalancersToInfra(platform, ignData, publicIPAddresses, privateIPAddresses)
	if err != nil {
		return nil, fmt.Errorf("failed to add load balancers to ignition config: %w", err)
	}
	logrus.Debugf("Successfully added API and API-Int Load Balancer IPs to infrastructure manifest")

	publicIPsSingleStr := strings.Join(publicIPAddresses, ",")
	privateIPsSingleStr := strings.Join(privateIPAddresses, ",")

	lbConfig, err := lbconfig.GenerateLBConfigOverride(privateIPsSingleStr, publicIPsSingleStr)
	if err != nil {
		return nil, err
	}
	if err := asset.NewDefaultFileWriter(lbConfig).PersistToFile(command.RootOpts.Dir); err != nil {
		return nil, fmt.Errorf("failed to save %s to state file: %w", lbConfig.Name(), err)
	}
	logrus.Debugf("Successfully added API and API-Int Load Balancer IPs to lbConfig configmap")

	err = updateMCSCertKey(in, platform, ignData, privateIPAddresses)
	if err != nil {
		return nil, fmt.Errorf("failed to update MCS Cert and Key in the bootstrap ignition: %w", err)
	}
	logrus.Debugf("Successfully updated bootstrap machine-config-server certs")

	editedMasterIgn, err := updatePointerIgnition(in, privateIPAddresses, masterRole)
	if err != nil {
		return nil, fmt.Errorf("failed to update %s Pointer Ignition: %w", masterRole, err)
	}
	logrus.Debugf("Successfully updated %s pointer ignition with API LB IP", masterRole)

	editedWorkerIgn, err := updatePointerIgnition(in, privateIPAddresses, workerRole)
	if err != nil {
		return nil, fmt.Errorf("failed to update %s Pointer Ignition: %w", workerRole, err)
	}
	logrus.Debugf("Successfully updated %s pointer ignition with API LB IP", workerRole)

	err = updateUserDataSecret(in, masterRole, ignData, editedMasterIgn)
	if err != nil {
		return nil, fmt.Errorf("failed to update %s user data secret: %w", masterRole, err)
	}
	logrus.Debugf("Successfully updated %s user data secret", masterRole)

	err = updateUserDataSecret(in, workerRole, ignData, editedWorkerIgn)
	if err != nil {
		return nil, fmt.Errorf("failed to update %s user data secret: %w", workerRole, err)
	}
	logrus.Debugf("Successfully updated %s user data secret", workerRole)

	editedIgnBytes, err := json.Marshal(ignData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert ignition data to json: %w", err)
	}
	logrus.Debugf("Successfully updated bootstrap ignition with updated manifests")

	ignOutput := &IgnitionOutput{
		UpdatedBootstrapIgn: editedIgnBytes,
		UpdatedMasterIgn:    editedMasterIgn,
		UpdatedWorkerIgn:    editedWorkerIgn,
	}
	return ignOutput, nil
}

// AddLoadBalancersToInfra will load the public and private load balancer information into
// the infrastructure CR. This will occur after the data has already been inserted into the
// ignition file.
func addLoadBalancersToInfra(platform string, config *igntypes.Config, publicLBs []string, privateLBs []string) error {
	for i, fileData := range config.Storage.Files {
		// update the contents of this file
		if fileData.Path == infrastructureFilepath {
			contents := strings.Split(*config.Storage.Files[i].Contents.Source, ",")
			rawDecodedText, err := base64.StdEncoding.DecodeString(contents[1])
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
				if infra.Status.PlatformStatus.AWS.CloudLoadBalancerConfig.DNSType == configv1.ClusterHostedDNSType {
					infra.Status.PlatformStatus.AWS.CloudLoadBalancerConfig.ClusterHosted = &cloudLBInfo
				}
			case azuretypes.Name:
				if infra.Status.PlatformStatus.Azure.CloudLoadBalancerConfig.DNSType == configv1.ClusterHostedDNSType {
					infra.Status.PlatformStatus.Azure.CloudLoadBalancerConfig.ClusterHosted = &cloudLBInfo
				}
			default:
				return fmt.Errorf("invalid platform %s", platform)
			}

			// convert the infrastructure back to an encoded string
			infraContents, err := yaml.Marshal(infra)
			if err != nil {
				return fmt.Errorf("failed to marshal infrastructure: %w", err)
			}

			encoded := fmt.Sprintf("%s%s", header, base64.StdEncoding.EncodeToString(infraContents))
			// replace the contents with the edited information
			config.Storage.Files[i].Contents.Source = &encoded

			break
		}
	}

	return nil
}

func updatePointerIgnition(in IgnitionInput, privateLBs []string, role string) ([]byte, error) {
	ignData := &igntypes.Config{}
	switch role {
	case masterRole:
		if len(privateLBs) == 0 {
			return in.MasterIgnData, nil
		}
		err := json.Unmarshal(in.MasterIgnData, ignData)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal %s ignition: %w", role, err)
		}
	case workerRole:
		if len(privateLBs) == 0 {
			return in.WorkerIgnData, nil
		}
		err := json.Unmarshal(in.WorkerIgnData, ignData)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal %s ignition: %w", role, err)
		}
	default:
		return nil, fmt.Errorf("unrecognized role %s for updating pointer ignition", role)
	}
	ignitionHost := net.JoinHostPort(privateLBs[0], "22623")
	ignData.Ignition.Config.Merge[0].Source = ignutil.StrToPtr(func() *url.URL {
		return &url.URL{
			Scheme: "https",
			Host:   ignitionHost,
			Path:   fmt.Sprintf("/config/%s", role),
		}
	}().String())
	editedIgnBytes, err := json.Marshal(ignData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert %s pointer ignition data to json: %w", role, err)
	}
	return editedIgnBytes, nil
}

func updateMCSCertKey(in IgnitionInput, platform string, config *igntypes.Config, privateLBs []string) error {
	if len(privateLBs) > 0 {
		keyRaw, certRaw, err := tls.RegenerateMCSCertKey(in.InstallConfig, in.RootCA, privateLBs)
		if err != nil {
			return fmt.Errorf("failed to regenerate MCS Cert and Key: %w", err)
		}

		for i, fileData := range config.Storage.Files {
			switch fileData.Path {
			case mcsCertKeyFilepath:
				contents := strings.Split(*config.Storage.Files[i].Contents.Source, ",")
				rawDecodedText, err := base64.StdEncoding.DecodeString(contents[1])
				if err != nil {
					return fmt.Errorf("failed to decode contents of ignition file %s: %w", mcsCertKeyFilepath, err)
				}
				mcsSecret := &corev1.Secret{}
				if err := yaml.Unmarshal(rawDecodedText, mcsSecret); err != nil {
					return fmt.Errorf("failed to unmarshal MCSCertKey within ignition: %w", err)
				}
				mcsSecret.Data[corev1.TLSCertKey] = certRaw
				mcsSecret.Data[corev1.TLSPrivateKeyKey] = keyRaw
				// convert the mcsSecret back to an encoded string
				mcsSecretContents, err := yaml.Marshal(mcsSecret)
				if err != nil {
					return fmt.Errorf("failed to marshal MCS Secret: %w", err)
				}
				encoded := fmt.Sprintf("%s%s", header, base64.StdEncoding.EncodeToString(mcsSecretContents))
				// replace the contents with the edited information
				config.Storage.Files[i].Contents.Source = &encoded

				logrus.Debugf("Updated MCSCertKey file %s with new cert and key", mcsCertKeyFilepath)
			case mcsKeyFile:
				encoded := fmt.Sprintf("%s%s", header, base64.StdEncoding.EncodeToString(keyRaw))
				// replace the contents with the edited information
				config.Storage.Files[i].Contents.Source = &encoded
				logrus.Debugf("Updated MCSKey file %s with new key", mcsKeyFile)
			case mcsCertFile:
				encoded := fmt.Sprintf("%s%s", header, base64.StdEncoding.EncodeToString(certRaw))
				// replace the contents with the edited information
				config.Storage.Files[i].Contents.Source = &encoded
				logrus.Debugf("Updated MCSCert file %s with new cert", mcsCertFile)
			}
		}
	}
	return nil
}

func updateUserDataSecret(in IgnitionInput, role string, config *igntypes.Config, updatedPointerIgnition []byte) error {
	userDataFile := ""
	name := ""
	switch role {
	case masterRole:
		userDataFile = masterUserDataFile
		name = "master-user-data"
	case workerRole:
		userDataFile = workerUserDataFile
		name = "worker-user-data"
	default:
		return fmt.Errorf("user data secret cannot be updated for unrecognized role %s", role)
	}

	for i, fileData := range config.Storage.Files {
		if fileData.Path == userDataFile {
			userDataSecret, err := machines.UserDataSecret(name, updatedPointerIgnition)
			if err != nil {
				return fmt.Errorf("failed to generate updated userData Secret for %s: %w", role, err)
			}
			encoded := fmt.Sprintf("%s%s", header, base64.StdEncoding.EncodeToString(userDataSecret))
			// replace the contents with the edited information
			config.Storage.Files[i].Contents.Source = &encoded
			return nil
		}
	}
	return nil
}
