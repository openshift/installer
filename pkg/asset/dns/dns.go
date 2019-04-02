package dns

import (
	"errors"

	azureasset "github.com/openshift/installer/pkg/asset/installconfig/azure"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/openstack"
)

var (
	//DNSConfigOverride allows to override the DNS Config Provider implementation during tests
	DNSConfigOverride ConfigProvider
)

//ConfigProvider is an interface that provides means to fetch the DNS settings
type ConfigProvider interface {
	GetBaseDomain() (string, error)
	GetPublicZone(name string) (string, error)
}

//NewConfig is a factory method to return the platform specific implementation of dnsConfig
func NewConfig(platform string) (ConfigProvider, error) {
	if DNSConfigOverride != nil {
		return DNSConfigOverride, nil
	}
	switch platform {
	case azure.Name:
		return azureasset.NewDNSConfig()
	case libvirt.Name, none.Name, openstack.Name, aws.Name:
		return nil, nil //not using the common interface yet
	default:
		return nil, errors.New("fail")
	}
}

//MockConfigProvider allows faking the dns settings
type MockConfigProvider struct {
	BaseDomain string
	PublicZone string
}

//GetBaseDomain returns the fake base domain
func (p *MockConfigProvider) GetBaseDomain() (string, error) {
	return p.BaseDomain, nil
}

//GetPublicZone return the fake public zone
func (p *MockConfigProvider) GetPublicZone(name string) (string, error) {
	return p.PublicZone, nil
}
