package baremetal

import (
	"fmt"
	"net"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	utilsnet "k8s.io/utils/net"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/manifests"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
)

// TemplateData holds data specific to templates used for the baremetal platform.
type TemplateData struct {
	// ProvisioningInterfaceMAC holds the interface's MAC address that the bootstrap node will use to host the ProvisioningIP below.
	// When the provisioning network is disabled, this is the external baremetal network MAC address.
	ProvisioningInterfaceMAC string

	// ProvisioningIP holds the IP the bootstrap node will use to service Ironic, TFTP, etc.
	ProvisioningIP string

	// ProvisioningIPv6 determines if we are using IPv6 or not.
	ProvisioningIPv6 bool

	// ProvisioningCIDR has the integer CIDR notation, e.g. 255.255.255.0 should be "24"
	ProvisioningCIDR int

	// ProvisioningDNSMasq determines if we start the dnsmasq service on the bootstrap node.
	ProvisioningDNSMasq bool

	// ProvisioningDHCPRange has the DHCP range, if DHCP is not external. Otherwise it
	// should be blank.
	ProvisioningDHCPRange string

	// ProvisioningDHCPAllowList contains a space-separated list of all of the control plane's boot
	// MAC addresses. Requests to bootstrap DHCP from other hosts will be ignored.
	ProvisioningDHCPAllowList string

	// IronicUsername contains the username for authentication to Ironic
	IronicUsername string

	// IronicUsername contains the password for authentication to Ironic
	IronicPassword string

	// BaremetalEndpointOverride contains the url for the baremetal endpoint
	BaremetalEndpointOverride string

	// BaremetalIntrospectionEndpointOverride contains the url for the baremetal introspection endpoint
	BaremetalIntrospectionEndpointOverride string

	// ClusterOSImage contains 4 URLs to download RHCOS live iso, kernel, rootfs and initramfs
	ClusterOSImage string

	// API VIP for use by ironic during bootstrap.
	APIVIPs []string

	// Hosts is the information needed to create the objects in Ironic.
	Hosts []*baremetal.Host

	// How many of the Hosts are control plane machines?
	ControlPlaneReplicas int64

	// ProvisioningNetwork displays the type of provisioning network being used
	ProvisioningNetwork string

	// ExternalStaticIP is the static IP of the bootstrap node
	ExternalStaticIP string

	// ExternalStaticIP is the static gateway of the bootstrap node
	ExternalStaticGateway string

	// ExternalStaticDNS is the static DNS of the bootstrap node
	ExternalStaticDNS string

	ExternalSubnetCIDR int

	ExternalMACAddress string

	// ExternalURLv6 is a callback URL for the node if the node and the BMC use different network families
	ExternalURLv6 string

	// DisableIronicVirtualMediaTLS enables or disables TLS in ironic virtual media deployments
	DisableIronicVirtualMediaTLS bool

	// AdditionalNTPServers holds a list of additional NTP servers to be used for provisioning
	AdditionalNTPServers []string
}

func externalURLs(apiVIPs []string, protocol string) (externalURLv4 string, externalURLv6 string) {
	if len(apiVIPs) > 1 {
		// IPv6 BMCs may not be able to reach IPv4 servers, use the right callback URL for them.
		// Warning: when backporting to 4.12 or earlier, change the port to 80!
		port := "6180"
		if protocol == "https" {
			port = "6183"
		}
		externalURL := fmt.Sprintf("%s://%s/", protocol, net.JoinHostPort(apiVIPs[1], port))
		if utilsnet.IsIPv6String(apiVIPs[1]) {
			externalURLv6 = externalURL
		}
		if utilsnet.IsIPv4String(apiVIPs[1]) {
			externalURLv4 = externalURL
		}
	}

	return
}

// GetTemplateData returns platform-specific data for bootstrap templates.
func GetTemplateData(config *baremetal.Platform, networks []types.MachineNetworkEntry, controlPlaneReplicaCount int64, ironicUsername, ironicPassword string, dependencies asset.Parents) *TemplateData {
	var templateData TemplateData

	templateData.Hosts = config.Hosts
	templateData.ControlPlaneReplicas = controlPlaneReplicaCount

	templateData.ProvisioningIP = config.BootstrapProvisioningIP
	templateData.ProvisioningNetwork = string(config.ProvisioningNetwork)
	templateData.ExternalStaticIP = config.BootstrapExternalStaticIP
	templateData.ExternalStaticGateway = config.BootstrapExternalStaticGateway
	templateData.ExternalStaticDNS = config.BootstrapExternalStaticDNS
	templateData.ExternalMACAddress = config.ExternalMACAddress

	if len(config.AdditionalNTPServers) > 0 {
		templateData.AdditionalNTPServers = config.AdditionalNTPServers
	}

	// If the user has manually set disableVirtualMediaTLS to False in the Provisioning CR, then enable TLS in ironic.
	// The default value is 'false'.
	templateData.DisableIronicVirtualMediaTLS = false
	protocol := "https"

	type provisioningCRTemplate struct {
		Spec struct {
			DisableVirtualMediaTLS *bool `yaml:"disableVirtualMediaTLS"`
		} `yaml:"spec"`
	}
	provisioningCR := &provisioningCRTemplate{}
	openshiftManifests := &manifests.Openshift{}
	dependencies.Get(openshiftManifests)
	var provisioningCRBytes []byte
	for _, file := range openshiftManifests.Files() {
		if strings.Contains(file.Filename, "99_baremetal-provisioning-config.yaml") {
			provisioningCRBytes = file.Data
			break
		}
	}
	if provisioningCRBytes != nil {
		err := yaml.Unmarshal(provisioningCRBytes, provisioningCR)
		if err != nil {
			logrus.Errorf("Error in unmarshalling Provisioning CR while generating TLS certificate for ironic virtual media: %s", err)
		}
	} else {
		logrus.Errorf("No Provisioning CR data found while generating TLS certificate for ironic virtual media")
	}
	if provisioningCR.Spec.DisableVirtualMediaTLS != nil && *provisioningCR.Spec.DisableVirtualMediaTLS {
		templateData.DisableIronicVirtualMediaTLS = *provisioningCR.Spec.DisableVirtualMediaTLS
		logrus.Debugf("TLS is disabled for ironic virtual media")
		protocol = "http"
	}
	_, externalURLv6 := externalURLs(config.APIVIPs, protocol)
	templateData.ExternalURLv6 = externalURLv6

	if len(config.APIVIPs) > 0 {
		templateData.APIVIPs = config.APIVIPs
		templateData.BaremetalEndpointOverride = fmt.Sprintf("https://%s/v1", net.JoinHostPort(config.APIVIPs[0], "6385"))
		templateData.BaremetalIntrospectionEndpointOverride = fmt.Sprintf("http://%s/v1", net.JoinHostPort(config.APIVIPs[0], "5050"))
	}

	if config.BootstrapExternalStaticIP != "" {
		for _, network := range networks {
			cidr, _ := network.CIDR.Mask.Size()
			templateData.ExternalSubnetCIDR = cidr
			break
		}
	}

	if config.ProvisioningNetwork != baremetal.DisabledProvisioningNetwork {
		cidr, _ := config.ProvisioningNetworkCIDR.Mask.Size()
		templateData.ProvisioningCIDR = cidr
		templateData.ProvisioningIPv6 = config.ProvisioningNetworkCIDR.IP.To4() == nil
		templateData.ProvisioningInterfaceMAC = config.ProvisioningMACAddress
		templateData.ProvisioningDNSMasq = true
	}

	switch config.ProvisioningNetwork {
	case baremetal.ManagedProvisioningNetwork:
		cidr, _ := config.ProvisioningNetworkCIDR.Mask.Size()

		// When provisioning network is managed, we set a DHCP range including
		// netmask for dnsmasq.
		templateData.ProvisioningDHCPRange = fmt.Sprintf("%s,%d", config.ProvisioningDHCPRange, cidr)

		var dhcpAllowList []string
		for _, host := range config.Hosts {
			if host.IsMaster() {
				dhcpAllowList = append(dhcpAllowList, host.BootMACAddress)
			}
		}
		templateData.ProvisioningDHCPAllowList = strings.Join(dhcpAllowList, " ")
	case baremetal.DisabledProvisioningNetwork:
		templateData.ProvisioningInterfaceMAC = config.ExternalMACAddress
		templateData.ProvisioningDNSMasq = false

		if templateData.ProvisioningIP != "" {
			for _, network := range networks {
				if network.CIDR.Contains(net.ParseIP(templateData.ProvisioningIP)) {
					templateData.ProvisioningIPv6 = network.CIDR.IP.To4() == nil

					cidr, _ := network.CIDR.Mask.Size()
					templateData.ProvisioningCIDR = cidr
				}
			}
		}
	}

	templateData.IronicUsername = ironicUsername
	templateData.IronicPassword = ironicPassword
	templateData.ClusterOSImage = config.ClusterOSImage

	return &templateData
}
