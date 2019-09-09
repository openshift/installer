package openstack

import (
	"bytes"
	"strconv"

	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/pkg/errors"
	ini "gopkg.in/ini.v1"
)

type config struct {
	Global global
}

type global struct {
	AuthURL    string `ini:"auth-url,omitempty"`
	Username   string `ini:"username,omitempty"`
	UserID     string `ini:"user-id,omitempty"`
	Password   string `ini:"password,omitempty"`
	TenantID   string `ini:"tenant-id,omitempty"`
	TenantName string `ini:"tenant-name,omitempty"`
	DomainID   string `ini:"domain-id,omitempty"`
	DomainName string `ini:"domain-name,omitempty"`
	Region     string `ini:"region,omitempty"`
	CAFile     string `ini:"ca-file,omitempty"`
}

// CloudProviderConfig generates the cloud provider config for the OpenStack platform.
func CloudProviderConfig(cloud *clientconfig.Cloud) string {
	res := `[Global]
secret-name = openstack-credentials
secret-namespace = kube-system
kubeconfig-path = /var/lib/kubelet/kubeconfig
`
	if cloud.RegionName != "" {
		res += "region = " + cloud.RegionName + "\n"
	}

	return res
}

// CloudProviderConfigSecret generates the cloud provider config for the OpenStack
// platform, that will be stored in the system secret.
func CloudProviderConfigSecret(cloud *clientconfig.Cloud) ([]byte, error) {
	file := ini.Empty()

	domainID := cloud.AuthInfo.DomainID
	if domainID == "" {
		domainID = cloud.AuthInfo.UserDomainID
	}

	domainName := cloud.AuthInfo.DomainName
	if domainName == "" {
		domainName = cloud.AuthInfo.UserDomainName
	}

	config := &config{
		Global: global{
			AuthURL:    cloud.AuthInfo.AuthURL,
			Username:   cloud.AuthInfo.Username,
			UserID:     cloud.AuthInfo.UserID,
			Password:   strconv.Quote(cloud.AuthInfo.Password),
			TenantID:   cloud.AuthInfo.ProjectID,
			TenantName: cloud.AuthInfo.ProjectName,
			DomainID:   domainID,
			DomainName: domainName,
			Region:     cloud.RegionName,
			CAFile:     cloud.CACertFile,
		},
	}
	if err := file.ReflectFrom(config); err != nil {
		return nil, errors.Wrap(err, "failed to reflect from config")
	}

	buf := &bytes.Buffer{}
	if _, err := file.WriteTo(buf); err != nil {
		return nil, errors.Wrap(err, "failed to write out cloud provider config")
	}
	return buf.Bytes(), nil
}
