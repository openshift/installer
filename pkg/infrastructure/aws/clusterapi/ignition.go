package clusterapi

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sirupsen/logrus"
	capa "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	k8sClient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	awstypes "github.com/openshift/installer/pkg/types/aws"
)

func editIgnition(ctx context.Context, in clusterapi.IgnitionInput) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*2)
	defer cancel()

	// TODO: only execute here if custom DNS enabled
	//if in.InstallConfig.Config.AWS.UserProvisionedDNS != aws.UserProvisionedDNSEnabled {
	//	return nil, nil
	//}

	awsCluster := &capa.AWSCluster{}
	key := k8sClient.ObjectKey{
		Name:      in.InfraID,
		Namespace: capiutils.Namespace,
	}
	if err := in.Client.Get(ctx, key, awsCluster); err != nil {
		return nil, fmt.Errorf("failed to get AWSCluster: %w", err)
	}

	awsSession, err := in.InstallConfig.AWS.Session(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get aws session: %w", err)
	}

	securityGroupIDs := []*string{}
	for _, securityGroup := range awsCluster.Status.Network.SecurityGroups {
		logrus.Warnf("Found Security Group ID: %s", securityGroup.ID)
		securityGroupIDs = append(securityGroupIDs, aws.String(securityGroup.ID))
	}
	nicInput := ec2.DescribeNetworkInterfacesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("group-id"),
				Values: securityGroupIDs,
			},
		},
	}
	nicOutput, err := ec2.New(awsSession).DescribeNetworkInterfacesWithContext(ctx, &nicInput)
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
			logrus.Debugf("found public IP address associated with %s", *nic.Description)
			publicIPAddresses = append(publicIPAddresses, *nic.Association.PublicIp)
		} else if nic.PrivateIpAddress != nil {
			logrus.Debugf("found private IP address associated with %s", *nic.Description)
			privateIPAddresses = append(privateIPAddresses, *nic.PrivateIpAddress)
		}
	}

	// TODO: this wont do anything, just shows us that the data is available
	logrus.Warnf("public load balancer addresses: %+v", publicIPAddresses)
	logrus.Warnf("private load balancer addresses: %+v", privateIPAddresses)

	return clusterapi.EditIgnition(ctx, in, awstypes.Name, publicIPAddresses, privateIPAddresses)
}
