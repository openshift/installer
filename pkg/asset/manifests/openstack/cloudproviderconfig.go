package openstack

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/openshift/installer/pkg/asset/installconfig/openstack"
	"github.com/openshift/installer/pkg/types"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
)

func getCaBundle(cloud *clientconfig.Cloud, installConfig types.InstallConfig) (string, error) {
	// install-config's additionalTrustBundle property overrides
	// clouds.yaml's cacert
	if caCert := installConfig.AdditionalTrustBundle; caCert != "" {
		return caCert, nil
	}

	// Get the ca-cert-bundle key if there is a value for cacert in clouds.yaml
	if caPath := cloud.CACertFile; caPath != "" {
		caCert, err := ioutil.ReadFile(caPath)
		if err != nil {
			return "", errors.Wrap(err, "failed to read clouds.yaml ca-cert from disk")
		}
		return string(caCert), nil
	}

	return "", nil
}

// GenerateCloudProviderConfig adds the cloud provider config for the OpenStack
// platform in the provided configmap.
func GenerateCloudProviderConfig(cm *corev1.ConfigMap, cloudProviderConfigDataKey string, installConfig types.InstallConfig) error {
	cloud, err := openstack.GetSession(installConfig.Platform.OpenStack.Cloud)
	if err != nil {
		return errors.Wrap(err, "failed to get cloud config for openstack")
	}

	cm.Data[cloudProviderConfigDataKey] = `[Global]
secret-name = openstack-credentials
secret-namespace = kube-system
`
	if regionName := cloud.CloudConfig.RegionName; regionName != "" {
		cm.Data[cloudProviderConfigDataKey] += "region = " + regionName + "\n"
	}

	caBundle, err := getCaBundle(cloud.CloudConfig, installConfig)
	if err != nil {
		return errors.Wrap(err, "failed to get the additional trust-bundle")
	}
	if caBundle != "" {
		cm.Data[cloudProviderConfigDataKey] += "ca-file = /etc/kubernetes/static-pod-resources/configmaps/cloud-config/ca-bundle.pem\n"
		cm.Data["ca-bundle.pem"] = caBundle
	}

	return nil
}

// CloudProviderConfigSecret generates the cloud provider config for the OpenStack
// platform, that will be stored in the system secret.
func CloudProviderConfigSecret(cloud *clientconfig.Cloud, installConfig types.InstallConfig) ([]byte, error) {
	domainID := cloud.AuthInfo.DomainID
	if domainID == "" {
		domainID = cloud.AuthInfo.UserDomainID
	}

	domainName := cloud.AuthInfo.DomainName
	if domainName == "" {
		domainName = cloud.AuthInfo.UserDomainName
	}

	// We have to generate this config manually without "go-ini" library, because its
	// output data is incompatible with "gcfg".
	// For instance, if there is a string with a # character, then "go-ini" wraps it in bacticks,
	// like `aaa#bbb`, but gcfg doesn't recognize it and  parses the data as `aaa, skipping
	// everything after the #.
	// For more information: https://bugzilla.redhat.com/show_bug.cgi?id=1771358
	var res strings.Builder
	res.WriteString("[Global]\n")
	if cloud.AuthInfo.AuthURL != "" {
		res.WriteString("auth-url = " + strconv.Quote(cloud.AuthInfo.AuthURL) + "\n")
	}
	if cloud.AuthInfo.Username != "" {
		res.WriteString("username = " + strconv.Quote(cloud.AuthInfo.Username) + "\n")
	}
	if cloud.AuthInfo.Password != "" {
		res.WriteString("password = " + strconv.Quote(cloud.AuthInfo.Password) + "\n")
	}
	if cloud.AuthInfo.ProjectID != "" {
		res.WriteString("tenant-id = " + strconv.Quote(cloud.AuthInfo.ProjectID) + "\n")
	}
	if cloud.AuthInfo.ProjectName != "" {
		res.WriteString("tenant-name = " + strconv.Quote(cloud.AuthInfo.ProjectName) + "\n")
	}
	if domainID != "" {
		res.WriteString("domain-id = " + strconv.Quote(domainID) + "\n")
	}
	if domainName != "" {
		res.WriteString("domain-name = " + strconv.Quote(domainName) + "\n")
	}
	if cloud.RegionName != "" {
		res.WriteString("region = " + strconv.Quote(cloud.RegionName) + "\n")
	}
	caBundle, err := getCaBundle(cloud, installConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the additional trust-bundle")
	}
	if caBundle != "" {
		res.WriteString("ca-file = /etc/kubernetes/static-pod-resources/configmaps/cloud-config/ca-bundle.pem\n")
	}

	return []byte(res.String()), nil
}
