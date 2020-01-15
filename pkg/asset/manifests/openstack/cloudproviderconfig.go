package openstack

import (
	"strconv"
	"strings"

	"github.com/gophercloud/utils/openstack/clientconfig"
)

// CloudProviderConfig generates the cloud provider config for the OpenStack platform.
func CloudProviderConfig(cloud *clientconfig.Cloud) string {
	res := `[Global]
secret-name = openstack-credentials
secret-namespace = kube-system
`
	if cloud.RegionName != "" {
		res += "region = " + cloud.RegionName + "\n"
	}

	if cloud.CACertFile != "" {
		res += "ca-file = /etc/kubernetes/static-pod-resources/configmaps/cloud-config/ca-bundle.pem\n"
	}

	return res
}

// CloudProviderConfigSecret generates the cloud provider config for the OpenStack
// platform, that will be stored in the system secret.
func CloudProviderConfigSecret(cloud *clientconfig.Cloud) ([]byte, error) {
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
	if cloud.CACertFile != "" {
		res.WriteString("ca-file = /etc/kubernetes/static-pod-resources/configmaps/cloud-config/ca-bundle.pem\n")
	}

	return []byte(res.String()), nil
}
