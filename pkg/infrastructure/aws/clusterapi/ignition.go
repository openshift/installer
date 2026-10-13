package clusterapi

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/sirupsen/logrus"
	capa "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	k8sClient "sigs.k8s.io/controller-runtime/pkg/client"

	awsconfig "github.com/openshift/installer/pkg/asset/installconfig/aws"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/dns"
	"github.com/openshift/installer/pkg/types/network"
)

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

	// There is no direct access to load balancer IP addresses, so the security groups
	// are used here to find the network interfaces that correspond to the load balancers.
	securityGroupIDs := make([]string, 0, len(awsCluster.Status.Network.SecurityGroups))
	for _, securityGroup := range awsCluster.Status.Network.SecurityGroups {
		securityGroupIDs = append(securityGroupIDs, securityGroup.ID)
	}
	nicInput := ec2.DescribeNetworkInterfacesInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("group-id"),
				Values: securityGroupIDs,
			},
		},
	}

	platformAWS := in.InstallConfig.Config.AWS
	client, err := awsconfig.NewEC2Client(ctx, awsconfig.EndpointOptions{
		Region:    platformAWS.Region,
		Endpoints: platformAWS.ServiceEndpoints,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create ec2 client: %w", err)
	}

	nicOutput, err := client.DescribeNetworkInterfaces(ctx, &nicInput)
	if err != nil {
		return nil, fmt.Errorf("failed to describe network interfaces: %w", err)
	}

	// The only network interfaces existing at this stage are those from the load balancers.
	// If this stage is executed after control plane nodes are provisioned, there may be
	// other network interfaces available.
	publicIPAddresses := []string{}
	privateIPAddresses := []string{}
	for _, nic := range nicOutput.NetworkInterfaces {
		if nic.Association != nil && nic.Association.PublicIp != nil {
			logrus.Debugf("found public IP address %s associated with %s", *nic.Association.PublicIp, *nic.Description)
			publicIPAddresses = append(publicIPAddresses, *nic.Association.PublicIp)
		} else if nic.PrivateIpAddress != nil {
			logrus.Debugf("found private IP address %s associated with %s", *nic.PrivateIpAddress, *nic.Description)
			privateIPAddresses = append(privateIPAddresses, *nic.PrivateIpAddress)
		}
	}
	if !in.InstallConfig.Config.PublicAPI() {
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
