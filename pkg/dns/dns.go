package dns

import (
	"errors"

	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/openstack"
)

//ConfigProvider is an interface that provides means to fetch the DNS settings
type ConfigProvider interface {
	GetBaseDomain() (string, error)
	GetPublicZone(name string) string
}

//NewDNSConfig is a factory method to return the platform specific implementation of dnsConfig
func NewDNSConfig(platform string) (ConfigProvider, error) {
	switch platform {
	case azure.Name:
		return azure.NewDNSConfig()
	case libvirt.Name, none.Name, openstack.Name, aws.Name:
		return nil, nil //not using the common interface yet
	case "fake":
		return &MockConfigProvider{}, nil
	default:
		return nil, errors.New("fail")
	}
}

//MockConfigProvider allows faking the dns settings
type MockConfigProvider struct {
}

//GetBaseDomain returns the fake base domain
func (*MockConfigProvider) GetBaseDomain() (string, error) {
	return "cloudapp.azure.com", nil
}

//GetPublicZone return the fake public zone
func (*MockConfigProvider) GetPublicZone(name string) string {
	return ""
}
