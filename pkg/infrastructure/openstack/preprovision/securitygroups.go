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

		if err := attributestags.Add(ctx, networkClient, "security-groups", masterGroup.ID, "openshiftClusterID="+infraID).ExtractErr(); err != nil {
			return fmt.Errorf("failed to tag the Control plane security group: %w", err)
		}

		workerGroup, err = groups.Create(ctx, networkClient, groups.CreateOpts{
			Name:        infraID + "-worker",
			Description: description,
		}).Extract()
		if err != nil {
			return fmt.Errorf("failed to create the Compute security group: %w", err)
		}

		if err := attributestags.Add(ctx, networkClient, "security-groups", workerGroup.ID, "openshiftClusterID="+infraID).ExtractErr(); err != nil {
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

	addRules := func(ctx context.Context, ch chan<- rules.CreateOpts, securityGroupID string, s service, ipVersion rules.RuleEtherType, remoteCIDRs []string) {
		for _, proto := range s.Protocols() {
			for _, remoteCIDR := range remoteCIDRs {
				select {
				case ch <- rules.CreateOpts{
					Direction:      rules.DirIngress,
					Description:    description,
					EtherType:      ipVersion,
					SecGroupID:     securityGroupID,
					PortRangeMax:   s.maxPort,
					PortRangeMin:   s.minPort,
					Protocol:       proto,
					RemoteIPPrefix: remoteCIDR,
				}:
				case <-ctx.Done():
					return
				}
			}
		}
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// TODO(mandre) Explicitly enable egress

	r := make(chan rules.CreateOpts)
	go func() {
		defer close(r)

		// In this loop: security groups with a catch-all remote IP
		for ipVersion, anyIP := range map[rules.RuleEtherType][]string{
			rules.EtherType4: {"0.0.0.0/0"},
			rules.EtherType6: func() []string {
				switch len(machineV6CIDRs) {
				case 0:
					return []string{}
				default:
					return []string{"::/0"}
				}
			}(),
		} {
			addRules(ctx, r, masterGroup.ID, serviceAPI, ipVersion, anyIP)
			addRules(ctx, r, masterGroup.ID, serviceICMP, ipVersion, anyIP)
			addRules(ctx, r, workerGroup.ID, serviceHTTP, ipVersion, anyIP)
			addRules(ctx, r, workerGroup.ID, serviceHTTPS, ipVersion, anyIP)
			addRules(ctx, r, workerGroup.ID, serviceICMP, ipVersion, anyIP)
			if mastersSchedulable {
				addRules(ctx, r, masterGroup.ID, serviceHTTP, ipVersion, anyIP)
				addRules(ctx, r, masterGroup.ID, serviceHTTPS, ipVersion, anyIP)
			}
		}

		// In this loop: security groups with the machine CIDR as remote IPs
		for ipVersion, CIDRs := range map[rules.RuleEtherType][]string{
			rules.EtherType4: machineV4CIDRs,
			rules.EtherType6: machineV6CIDRs,
		} {
			// In this loop: rules that equally apply to masters and workers
			for _, groupID := range [...]string{masterGroup.ID, workerGroup.ID} {
				addRules(ctx, r, groupID, serviceESP, ipVersion, CIDRs)
				addRules(ctx, r, groupID, serviceGeneve, ipVersion, CIDRs)
				addRules(ctx, r, groupID, serviceIKE, ipVersion, CIDRs)
				addRules(ctx, r, groupID, serviceInternal, ipVersion, CIDRs)
				addRules(ctx, r, groupID, serviceKubelet, ipVersion, CIDRs)
				addRules(ctx, r, groupID, serviceNodeport, ipVersion, CIDRs)
				addRules(ctx, r, groupID, serviceSSH, ipVersion, CIDRs)
				addRules(ctx, r, groupID, serviceVRRP, ipVersion, CIDRs)
				addRules(ctx, r, groupID, serviceVXLAN, ipVersion, CIDRs)
			}

			addRules(ctx, r, masterGroup.ID, serviceDNS, ipVersion, CIDRs)
			addRules(ctx, r, masterGroup.ID, serviceETCD, ipVersion, CIDRs)
			addRules(ctx, r, masterGroup.ID, serviceKCM, ipVersion, CIDRs)
			addRules(ctx, r, masterGroup.ID, serviceKubeScheduler, ipVersion, CIDRs)
			addRules(ctx, r, masterGroup.ID, serviceMCS, ipVersion, CIDRs)
			addRules(ctx, r, masterGroup.ID, serviceOVNDB, ipVersion, CIDRs)
			addRules(ctx, r, workerGroup.ID, serviceRouter, ipVersion, CIDRs)
			if mastersSchedulable {
				addRules(ctx, r, masterGroup.ID, serviceRouter, ipVersion, CIDRs)
			}
		}

		// IPv4-only rules
		addRules(ctx, r, masterGroup.ID, serviceIKENat, rules.EtherType4, machineV4CIDRs)
		addRules(ctx, r, workerGroup.ID, serviceIKENat, rules.EtherType4, machineV4CIDRs)
	}()

	for ruleCreateOpts := range r {
		if _, err := rules.Create(ctx, networkClient, ruleCreateOpts).Extract(); err != nil {
			return fmt.Errorf("failed to create the security group rule on group %q for %s %s on ports %d-%d: %w", ruleCreateOpts.SecGroupID, ruleCreateOpts.EtherType, ruleCreateOpts.Protocol, ruleCreateOpts.PortRangeMin, ruleCreateOpts.PortRangeMax, err)
		}
	}
	return nil
}
