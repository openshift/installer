package clusterapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	configv2 "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/sirupsen/logrus"
	capa "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	k8sClient "sigs.k8s.io/controller-runtime/pkg/client"

	machinev1 "github.com/openshift/api/machineconfiguration/v1"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/dns"
)

func editIgnition(ctx context.Context, in clusterapi.IgnitionInput) (*clusterapi.IgnitionOutput, error) {
	// To be revised: The IgnitionOutput and IgnitionInput can be chained in a series of editIgnition calls.
	// To update output with latest ignition output, use (*IgnitionOutput).Set(*IgnitionOutput).
	// To update input with latest ignition output (to use as input for another editIgnition call), use (*IgnitionInput).SetIgnFromOutput(*IgnitionOutput).
	ignOutput := &clusterapi.IgnitionOutput{
		UpdatedBootstrapIgn: in.BootstrapIgnData,
		UpdatedMasterIgn:    in.MasterIgnData,
		UpdatedWorkerIgn:    in.WorkerIgnData,
	}

	// Dualstack mode
	if in.InstallConfig.Config.IsDualStackInfra() {
		// If there is no IPv6 machine network in the install-config, that means AWS will assign a IPv6 GUA range during
		// infrastructure provisioning. Thus, after VPC is ready, we need to populate the install config with the VPC IPv6 CIDR.
		if len(capiutils.MachineCIDRsFromInstallConfig(in.InstallConfig).IPv6Nets()) == 0 {
			dualStackIgnOut, err := editIgnitionForDualStack(ctx, in)
			if err != nil {
				return ignOutput, err
			}
			ignOutput.Set(dualStackIgnOut)
			in.SetIgnFromOutput(ignOutput)
		}
	}

	// Custom DNS
	if in.InstallConfig.Config.AWS.UserProvisionedDNS == dns.UserProvisionedDNSEnabled {
		customDNSIgnOut, err := editIgnitionForCustomDNS(ctx, in)
		if err != nil {
			return ignOutput, err
		}
		ignOutput.Set(customDNSIgnOut)
		in.SetIgnFromOutput(ignOutput)
	}

	return ignOutput, nil
}

func editIgnitionForCustomDNS(ctx context.Context, in clusterapi.IgnitionInput) (*clusterapi.IgnitionOutput, error) {
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

	cfg, err := configv2.LoadDefaultConfig(ctx, configv2.WithRegion(in.InstallConfig.Config.AWS.Region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := ec2.NewFromConfig(cfg, func(options *ec2.Options) {
		options.Region = in.InstallConfig.Config.AWS.Region
		for _, endpoint := range in.InstallConfig.Config.AWS.ServiceEndpoints {
			if strings.EqualFold(endpoint.Name, "ec2") {
				options.BaseEndpoint = aws.String(endpoint.URL)
			}
		}
	})

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
	awsCluster := &capa.AWSCluster{}
	key := k8sClient.ObjectKey{
		Name:      in.InfraID,
		Namespace: capiutils.Namespace,
	}
	if err := in.Client.Get(ctx, key, awsCluster); err != nil {
		return nil, fmt.Errorf("failed to get AWSCluster: %w", err)
	}

	machineNetworks := in.InstallConfig.Config.MachineNetwork
	// Assumption: When IPv6 is enabled, the AWSCluster IPv6 block should be non-nil
	// And the VPC IPv6 CIDR is valid (i.e. set by CAPA)
	vpcIPv6CIDR := awsCluster.Spec.NetworkSpec.VPC.IPv6.CidrBlock

	ic := in.InstallConfig.Config
	ipv6Entry := []types.MachineNetworkEntry{
		{
			CIDR: *ipnet.MustParseCIDR(vpcIPv6CIDR),
		},
	}
	if ic.InfraStack() == machinev1.IPFamiliesDualStackIPv6Primary {
		machineNetworks = append(ipv6Entry, machineNetworks...)
	} else {
		machineNetworks = append(machineNetworks, ipv6Entry...)
	}

	logrus.Debugf("AWS: Editing Ignition files to add VPC IPv6 CIDR to the machine")
	return clusterapi.EditIgnitionForDualStack(in, awstypes.Name, machineNetworks)
}
