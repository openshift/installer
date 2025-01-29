package preprovision

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/rules"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset/installconfig"
	openstackdefaults "github.com/openshift/installer/pkg/types/openstack/defaults"
)

const description = "Created By OpenShift Installer"

type protocol uint8

const (
	_ protocol = 1 << iota
	tcp
	udp
	icmp
	esp
	vrrp
)

func (p protocol) Protocols() []rules.RuleProtocol {
	protocols := make([]rules.RuleProtocol, 0, 2)
	if p&tcp == tcp {
		protocols = append(protocols, rules.ProtocolTCP)
	}
	if p&udp == udp {
		protocols = append(protocols, rules.ProtocolUDP)
	}
	if p&icmp == icmp {
		protocols = append(protocols, rules.ProtocolICMP)
	}
	if p&esp == esp {
		protocols = append(protocols, rules.ProtocolESP)
	}
	if p&vrrp == vrrp {
		protocols = append(protocols, rules.RuleProtocol("112"))
	}
	return protocols
}

// SecurityGroups creates the master and worker security groups with the security group rules.
func SecurityGroups(ctx context.Context, installConfig *installconfig.InstallConfig, infraID string, mastersSchedulable bool) error {
	networkClient, err := openstackdefaults.NewServiceClient(ctx, "network", openstackdefaults.DefaultClientOpts(installConfig.Config.Platform.OpenStack.Cloud))
	if err != nil {
		return err
	}

	logrus.Debugf("Creating the security groups")
	var masterGroup, workerGroup *groups.SecGroup
	{
		masterGroup, err = groups.Create(ctx, networkClient, groups.CreateOpts{
			Name:        infraID + "-master",
			Description: description,
		}).Extract()
		if err != nil {
			return fmt.Errorf("failed to create the Control plane security group: %w", err)
		}

		// The Neutron call to add a tag
		// (https://docs.openstack.org/api-ref/network/v2/#add-a-tag)
		// doesn't accept all special characters. Here we use the
		// "replace-all-tags" call instead, because it accepts a more
		// robust JSON body.
		//
		// see: https://bugzilla.redhat.com/show_bug.cgi?id=2299208
		if _, err := attributestags.ReplaceAll(ctx, networkClient, "security-groups", masterGroup.ID, attributestags.ReplaceAllOpts{
			Tags: []string{"openshiftClusterID=" + infraID},
		}).Extract(); err != nil {
			return fmt.Errorf("failed to tag the Control plane security group: %w", err)
		}

		workerGroup, err = groups.Create(ctx, networkClient, groups.CreateOpts{
			Name:        infraID + "-worker",
			Description: description,
		}).Extract()
		if err != nil {
			return fmt.Errorf("failed to create the Compute security group: %w", err)
		}

		// See comment above
		if _, err := attributestags.ReplaceAll(ctx, networkClient, "security-groups", workerGroup.ID, attributestags.ReplaceAllOpts{
			Tags: []string{"openshiftClusterID=" + infraID},
		}).Extract(); err != nil {
			return fmt.Errorf("failed to tag the Compute security group: %w", err)
		}
	}

	logrus.Debugf("Creating the security group rules")
	var (
		machineV4CIDRs = make([]string, 0, len(installConfig.Config.Networking.MachineNetwork))
		machineV6CIDRs = make([]string, 0, len(installConfig.Config.Networking.MachineNetwork))
	)
	for _, network := range installConfig.Config.Networking.MachineNetwork {
		if network.CIDR.IPNet.IP.To4() != nil {
			machineV4CIDRs = append(machineV4CIDRs, network.CIDR.IPNet.String())
		} else {
			machineV6CIDRs = append(machineV6CIDRs, network.CIDR.IPNet.String())
		}
	}

	type service struct {
		protocol
		minPort int
		maxPort int
	}

	var (
		serviceAPI           = service{tcp | udp, 6443, 6443}
		serviceDNS           = service{tcp | udp, 53, 53}
		serviceESP           = service{protocol: esp}
		serviceETCD          = service{tcp, 2379, 2380}
		serviceGeneve        = service{udp, 6081, 6081}
		serviceHTTP          = service{tcp, 80, 80}
		serviceHTTPS         = service{tcp, 443, 443}
		serviceICMP          = service{protocol: icmp}
		serviceIKE           = service{udp, 500, 500}
		serviceIKENat        = service{udp, 4500, 4500}
		serviceInternal      = service{tcp | udp, 9000, 9999}
		serviceKCM           = service{tcp, 10257, 10257}
		serviceKubeScheduler = service{tcp, 10259, 10259}
		serviceKubelet       = service{tcp, 10250, 10250}
		serviceMCS           = service{tcp, 22623, 22623}
		serviceNodeport      = service{tcp | udp, 30000, 32767}
		serviceOVNDB         = service{tcp, 6641, 6642}
		serviceRouter        = service{tcp, 1936, 1936}
		serviceSSH           = service{tcp, 22, 22}
		serviceVRRP          = service{protocol: vrrp}
		serviceVXLAN         = service{udp, 4789, 4789}
	)

	buildRules := func(groupID string, s service, ipVersion rules.RuleEtherType, remoteCIDRs []string) (r []rules.CreateOpts) {
		for _, proto := range s.Protocols() {
			for _, remoteCIDR := range remoteCIDRs {
				r = append(r, rules.CreateOpts{
					Direction:      rules.DirIngress,
					Description:    description,
					EtherType:      ipVersion,
					SecGroupID:     groupID,
					PortRangeMax:   s.maxPort,
					PortRangeMin:   s.minPort,
					Protocol:       proto,
					RemoteIPPrefix: remoteCIDR,
				})
			}
		}
		return r
	}

	var masterRules []rules.CreateOpts
	var workerRules []rules.CreateOpts
	addMasterRules := func(s service, ipVersion rules.RuleEtherType, remoteCIDRs []string) {
		masterRules = append(masterRules, buildRules(masterGroup.ID, s, ipVersion, remoteCIDRs)...)
	}
	addWorkerRules := func(s service, ipVersion rules.RuleEtherType, remoteCIDRs []string) {
		workerRules = append(workerRules, buildRules(workerGroup.ID, s, ipVersion, remoteCIDRs)...)
	}

	// TODO(mandre) Explicitly enable egress

	// In this loop: security groups with a catch-all remote IP
	for ipVersion, anyIP := range map[rules.RuleEtherType][]string{
		rules.EtherType4: func() []string {
			switch len(machineV4CIDRs) {
			case 0:
				return []string{}
			default:
				return []string{"0.0.0.0/0"}
			}
		}(),
		rules.EtherType6: func() []string {
			switch len(machineV6CIDRs) {
			case 0:
				return []string{}
			default:
				return []string{"::/0"}
			}
		}(),
	} {
		addMasterRules(serviceAPI, ipVersion, anyIP)
		addMasterRules(serviceICMP, ipVersion, anyIP)
		addWorkerRules(serviceHTTP, ipVersion, anyIP)
		addWorkerRules(serviceHTTPS, ipVersion, anyIP)
		addWorkerRules(serviceICMP, ipVersion, anyIP)
		if mastersSchedulable {
			addMasterRules(serviceHTTP, ipVersion, anyIP)
			addMasterRules(serviceHTTPS, ipVersion, anyIP)
		}
	}

	// In this loop: security groups with the machine CIDR as remote IPs
	for ipVersion, CIDRs := range map[rules.RuleEtherType][]string{
		rules.EtherType4: machineV4CIDRs,
		rules.EtherType6: machineV6CIDRs,
	} {
		// In this loop: rules that equally apply to masters and workers
		for _, addRules := range [...]func(service, rules.RuleEtherType, []string){addMasterRules, addWorkerRules} {
			addRules(serviceESP, ipVersion, CIDRs)
			addRules(serviceGeneve, ipVersion, CIDRs)
			addRules(serviceIKE, ipVersion, CIDRs)
			addRules(serviceInternal, ipVersion, CIDRs)
			addRules(serviceKubelet, ipVersion, CIDRs)
			addRules(serviceNodeport, ipVersion, CIDRs)
			addRules(serviceSSH, ipVersion, CIDRs)
			addRules(serviceVRRP, ipVersion, CIDRs)
			addRules(serviceVXLAN, ipVersion, CIDRs)
		}

		addMasterRules(serviceDNS, ipVersion, CIDRs)
		addMasterRules(serviceETCD, ipVersion, CIDRs)
		addMasterRules(serviceKCM, ipVersion, CIDRs)
		addMasterRules(serviceKubeScheduler, ipVersion, CIDRs)
		addMasterRules(serviceMCS, ipVersion, CIDRs)
		addMasterRules(serviceOVNDB, ipVersion, CIDRs)
		addWorkerRules(serviceRouter, ipVersion, CIDRs)
		if mastersSchedulable {
			addMasterRules(serviceRouter, ipVersion, CIDRs)
		}
	}

	// IPv4-only rules
	addMasterRules(serviceIKENat, rules.EtherType4, machineV4CIDRs)
	addWorkerRules(serviceIKENat, rules.EtherType4, machineV4CIDRs)

	if _, err := rules.CreateBulk(ctx, networkClient, masterRules).Extract(); err != nil {
		return fmt.Errorf("failed to add the rules to the Control plane security group: %w", err)
	}
	if _, err := rules.CreateBulk(ctx, networkClient, workerRules).Extract(); err != nil {
		return fmt.Errorf("failed to add the rules to the Compute security group: %w", err)
	}

	return nil
}
