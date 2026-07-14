package clusterapi

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
	capa "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	k8sClient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/dns"
	"github.com/openshift/installer/pkg/types/network"
)

// getIPsFromDNSName resolves a DNS name to IP addresses with exponential backoff retry.
// It returns a slice of IP address strings or an error if all retries fail.
func getIPsFromDNSName(ctx context.Context, dnsName, lbType string) ([]string, error) {
	resolver := &net.Resolver{}
	var ips []net.IP
	err := wait.ExponentialBackoffWithContext(ctx, wait.Backoff{
		Duration: 1 * time.Second,
		Factor:   2.0,
		Jitter:   0.1,
		Steps:    10, // ~15 minutes total timeout with exponential backoff
	}, func(ctx context.Context) (bool, error) {
		var err error
		ips, err = resolver.LookupIP(ctx, "ip", dnsName)
		if err != nil {
			logrus.Debugf("AWS: DNS lookup for %s DNS name %q failed, retrying: %v", lbType, dnsName, err)
			return false, nil
		}
		if len(ips) == 0 {
			logrus.Debugf("AWS: DNS lookup for %s DNS name %q returned no IPs, retrying", lbType, dnsName)
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to lookup IP for %s DNS name %q after retries: %w", lbType, dnsName, err)
	}

	var ipAddresses []string
	for _, ip := range ips {
		ipAddresses = append(ipAddresses, ip.String())
	}
	logrus.Debugf("AWS: Resolved %s DNS name %q to IPs: %v", lbType, dnsName, ipAddresses)
	return ipAddresses, nil
}

func editIgnitionForCustomDNS(ctx context.Context, in clusterapi.IgnitionInput) (*clusterapi.IgnitionOutput, error) {
	if in.InstallConfig.Config.AWS.UserProvisionedDNS != dns.UserProvisionedDNSEnabled {
		return &clusterapi.IgnitionOutput{
			UpdatedBootstrapIgn: in.BootstrapIgnData,
			UpdatedMasterIgn:    in.MasterIgnData,
			UpdatedWorkerIgn:    in.WorkerIgnData}, nil
	}

	awsCluster := &capa.AWSCluster{}
	key := k8sClient.ObjectKey{
		Name:      in.InfraID,
		Namespace: capiutils.Namespace,
	}
	if err := in.Client.Get(ctx, key, awsCluster); err != nil {
		return nil, fmt.Errorf("failed to get AWSCluster: %w", err)
	}

	var publicIPAddresses, privateIPAddresses []string

	// Get private LB IP addresses from APIServerELB DNS name
	if dnsName := awsCluster.Status.Network.APIServerELB.DNSName; dnsName != "" {
		ips, err := getIPsFromDNSName(ctx, dnsName, fmt.Sprintf("APIServerELB (scheme: %s)", awsCluster.Status.Network.APIServerELB.Scheme))
		if err != nil {
			return nil, err
		}
		privateIPAddresses = ips
	} else {
		return nil, fmt.Errorf("internal API load balancer not found")
	}

	if in.InstallConfig.Config.PublicAPI() {
		// Get public LB IP addresses from SecondaryAPIServerELB DNS name
		if dnsName := awsCluster.Status.Network.SecondaryAPIServerELB.DNSName; dnsName != "" {
			ips, err := getIPsFromDNSName(ctx, dnsName, fmt.Sprintf("SecondaryAPIServerELB (scheme: %s)", awsCluster.Status.Network.SecondaryAPIServerELB.Scheme))
			if err != nil {
				return nil, err
			}
			publicIPAddresses = ips
		} else {
			return nil, fmt.Errorf("public API load balancer not found")
		}
	} else {
		// For private cluster installs, the API LB IP is the same as the API-Int LB IP
		publicIPAddresses = privateIPAddresses
	}

	logrus.Debugf("AWS: Editing Ignition files to start in-cluster DNS when UserProvisionedDNS is enabled")
	return clusterapi.EditIgnitionForCustomDNS(in, awstypes.Name, publicIPAddresses, privateIPAddresses)
}

func editIgnitionForDualStack(ctx context.Context, in clusterapi.IgnitionInput) (*clusterapi.IgnitionOutput, error) {
	ic := in.InstallConfig.Config
	machineCIDRs := capiutils.MachineCIDRsFromInstallConfig(in.InstallConfig)

	// If the machine network entries contain IPv6 CIDRs, the users must have added in manually for BYO subnets.
	// In this case, those CIDRs are already passed to the AWSCluster node port ingress rule spec
	if !ic.AWS.IPFamily.DualStackEnabled() || len(capiutils.GetIPv6CIDRs(machineCIDRs)) > 0 {
		return &clusterapi.IgnitionOutput{
			UpdatedBootstrapIgn: in.BootstrapIgnData,
			UpdatedMasterIgn:    in.MasterIgnData,
			UpdatedWorkerIgn:    in.WorkerIgnData}, nil
	}

	awsCluster := &capa.AWSCluster{}
	key := k8sClient.ObjectKey{
		Name:      in.InfraID,
		Namespace: capiutils.Namespace,
	}
	if err := in.Client.Get(ctx, key, awsCluster); err != nil {
		return nil, fmt.Errorf("failed to get AWSCluster: %w", err)
	}

	vpcSpec := awsCluster.Spec.NetworkSpec.VPC
	if vpcSpec.IPv6 == nil || vpcSpec.IPv6.CidrBlock == "" {
		return nil, fmt.Errorf("dualstack networking is enabled, but VPC does not have IPV6 CIDR")
	}

	machineNetworks := ic.MachineNetwork
	cidr, err := ipnet.ParseCIDR(vpcSpec.IPv6.CidrBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to parse VPC IPv6 CIDR block %q: %w", vpcSpec.IPv6.CidrBlock, err)
	}
	ipv6Entry := []types.MachineNetworkEntry{
		{
			CIDR: *cidr,
		},
	}

	if ic.AWS.IPFamily == network.DualStackIPv6Primary {
		machineNetworks = append(ipv6Entry, machineNetworks...)
	} else {
		machineNetworks = append(machineNetworks, ipv6Entry...)
	}

	return clusterapi.EditIgnitionForDualStack(in, awstypes.Name, machineNetworks)
}
