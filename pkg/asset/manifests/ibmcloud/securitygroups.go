package ibmcloud

import (
	"fmt"

	"k8s.io/utils/ptr"
	capibmcloud "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"

	ibmcloudic "github.com/openshift/installer/pkg/asset/installconfig/ibmcloud"
	"github.com/openshift/installer/pkg/types"
)

const (
	clusterWideSGNameSuffix  = "sg-cluster-wide"
	openshiftNetSGNameSuffix = "sg-openshift-net"
	kubeAPILBSGNameSuffix    = "sg-kube-api-lb"
	controlPlaneSGNameSuffix = "sg-control-plane"
	cpInternalSGNameSuffix   = "sg-cp-internal"
)

func buildClusterWideSecurityGroup(infraID string, allSubnets []capibmcloud.Subnet) capibmcloud.VPCSecurityGroup {
	clusterWideSGNamePtr := ptr.To(fmt.Sprintf("%s-%s", infraID, clusterWideSGNameSuffix))

	// Build set of Remotes for Security Group Rules
	// - cluster-wide SSH rule (for CP and Compute subnets)
	clusterWideSSHRemotes := make([]capibmcloud.VPCSecurityGroupRuleRemote, len(allSubnets))
	for index, subnet := range allSubnets {
		clusterWideSSHRemotes[index] = capibmcloud.VPCSecurityGroupRuleRemote{
			RemoteType:     capibmcloud.VPCSecurityGroupRuleRemoteTypeCIDR,
			CIDRSubnetName: subnet.Name,
		}
	}

	return capibmcloud.VPCSecurityGroup{
		Name: clusterWideSGNamePtr,
		Rules: []*capibmcloud.VPCSecurityGroupRule{
			{
				// SSH inbound cluster-wide
				Action:    capibmcloud.VPCSecurityGroupRuleActionAllow,
				Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
				Source: &capibmcloud.VPCSecurityGroupRulePrototype{
					PortRange: &capibmcloud.VPCSecurityGroupPortRange{
						MaximumPort: 22,
						MinimumPort: 22,
					},
					Protocol: capibmcloud.VPCSecurityGroupRuleProtocolTCP,
					Remotes:  clusterWideSSHRemotes,
				},
			},
			{
				// ICMP inbound cluster-wide
				Action:    capibmcloud.VPCSecurityGroupRuleActionAllow,
				Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
				Source: &capibmcloud.VPCSecurityGroupRulePrototype{
					Protocol: capibmcloud.VPCSecurityGroupRuleProtocolIcmp,
					Remotes: []capibmcloud.VPCSecurityGroupRuleRemote{
						{
							RemoteType:        capibmcloud.VPCSecurityGroupRuleRemoteTypeSG,
							SecurityGroupName: clusterWideSGNamePtr,
						},
					},
				},
			},
			{
				// VXLAN and Geneve - port 4789
				Action:    capibmcloud.VPCSecurityGroupRuleActionAllow,
				Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
				Source: &capibmcloud.VPCSecurityGroupRulePrototype{
					PortRange: &capibmcloud.VPCSecurityGroupPortRange{
						MaximumPort: 4789,
						MinimumPort: 4789,
					},
					Protocol: capibmcloud.VPCSecurityGroupRuleProtocolUDP,
					Remotes: []capibmcloud.VPCSecurityGroupRuleRemote{
						{
							RemoteType:        capibmcloud.VPCSecurityGroupRuleRemoteTypeSG,
							SecurityGroupName: clusterWideSGNamePtr,
						},
					},
				},
			},
			{
				// VXLAN and Geneve - port 6081
				Action:    capibmcloud.VPCSecurityGroupRuleActionAllow,
				Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
				Source: &capibmcloud.VPCSecurityGroupRulePrototype{
					PortRange: &capibmcloud.VPCSecurityGroupPortRange{
						MaximumPort: 6081,
						MinimumPort: 6081,
					},
					Protocol: capibmcloud.VPCSecurityGroupRuleProtocolUDP,
					Remotes: []capibmcloud.VPCSecurityGroupRuleRemote{
						{
							RemoteType:        capibmcloud.VPCSecurityGroupRuleRemoteTypeSG,
							SecurityGroupName: clusterWideSGNamePtr,
						},
					},
				},
			},
			{
				// Outbound for cluster-wide
				Action: capibmcloud.VPCSecurityGroupRuleActionAllow,
				Destination: &capibmcloud.VPCSecurityGroupRulePrototype{
					Protocol: capibmcloud.VPCSecurityGroupRuleProtocolAll,
					Remotes: []capibmcloud.VPCSecurityGroupRuleRemote{
						{
							RemoteType: capibmcloud.VPCSecurityGroupRuleRemoteTypeAny,
						},
					},
				},
				Direction: capibmcloud.VPCSecurityGroupRuleDirectionOutbound,
			},
		},
	}
}

func buildOpenshiftNetSecurityGroup(infraID string, allSubnets []capibmcloud.Subnet) capibmcloud.VPCSecurityGroup {
	openshiftNetSGNamePtr := ptr.To(fmt.Sprintf("%s-%s", infraID, openshiftNetSGNameSuffix))

	// Build sets of Remotes for Security Group Rules
	// - openshift-net TCP rule for Node Ports (for CP and Compute subnets)
	openshiftNetworkNodePortTCPRemotes := make([]capibmcloud.VPCSecurityGroupRuleRemote, len(allSubnets))
	// - openshift-net UDP rule for Node Ports (for CP and Compute subnets)
	openshiftNetworkNodePortUDPRemotes := make([]capibmcloud.VPCSecurityGroupRuleRemote, len(allSubnets))
	for index, subnet := range allSubnets {
		openshiftNetworkNodePortTCPRemotes[index] = capibmcloud.VPCSecurityGroupRuleRemote{
			RemoteType:     capibmcloud.VPCSecurityGroupRuleRemoteTypeCIDR,
			CIDRSubnetName: subnet.Name,
		}
		openshiftNetworkNodePortUDPRemotes[index] = capibmcloud.VPCSecurityGroupRuleRemote{
			RemoteType:     capibmcloud.VPCSecurityGroupRuleRemoteTypeCIDR,
			CIDRSubnetName: subnet.Name,
		}
	}

	return capibmcloud.VPCSecurityGroup{
		Name: openshiftNetSGNamePtr,
		Rules: []*capibmcloud.VPCSecurityGroupRule{
			{
				// Host level services - TCP
				Action:    capibmcloud.VPCSecurityGroupRuleActionAllow,
				Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
				Source: &capibmcloud.VPCSecurityGroupRulePrototype{
					PortRange: &capibmcloud.VPCSecurityGroupPortRange{
						MaximumPort: 9999,
						MinimumPort: 9000,
					},
					Protocol: capibmcloud.VPCSecurityGroupRuleProtocolTCP,
					Remotes: []capibmcloud.VPCSecurityGroupRuleRemote{
						{
							RemoteType:        capibmcloud.VPCSecurityGroupRuleRemoteTypeSG,
							SecurityGroupName: openshiftNetSGNamePtr,
						},
					},
				},
			},
			{
				// Host level services - UDP
				Action:    capibmcloud.VPCSecurityGroupRuleActionAllow,
				Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
				Source: &capibmcloud.VPCSecurityGroupRulePrototype{
					PortRange: &capibmcloud.VPCSecurityGroupPortRange{
						MaximumPort: 9999,
						MinimumPort: 9000,
					},
					Protocol: capibmcloud.VPCSecurityGroupRuleProtocolUDP,
					Remotes: []capibmcloud.VPCSecurityGroupRuleRemote{
						{
							RemoteType:        capibmcloud.VPCSecurityGroupRuleRemoteTypeSG,
							SecurityGroupName: openshiftNetSGNamePtr,
						},
					},
				},
			},
			{
				// Kubernetes default ports
				Action:    capibmcloud.VPCSecurityGroupRuleActionAllow,
				Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
				Source: &capibmcloud.VPCSecurityGroupRulePrototype{
					PortRange: &capibmcloud.VPCSecurityGroupPortRange{
						MaximumPort: 10250,
						MinimumPort: 10250,
					},
					Protocol: capibmcloud.VPCSecurityGroupRuleProtocolTCP,
					Remotes: []capibmcloud.VPCSecurityGroupRuleRemote{
						{
							RemoteType:        capibmcloud.VPCSecurityGroupRuleRemoteTypeSG,
							SecurityGroupName: openshiftNetSGNamePtr,
						},
					},
				},
			},
			{
				// IPsec IKE - port 500
				Action:    capibmcloud.VPCSecurityGroupRuleActionAllow,
				Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
				Source: &capibmcloud.VPCSecurityGroupRulePrototype{
					PortRange: &capibmcloud.VPCSecurityGroupPortRange{
						MaximumPort: 500,
						MinimumPort: 500,
					},
					Protocol: capibmcloud.VPCSecurityGroupRuleProtocolUDP,
					Remotes: []capibmcloud.VPCSecurityGroupRuleRemote{
						{
							RemoteType:        capibmcloud.VPCSecurityGroupRuleRemoteTypeSG,
							SecurityGroupName: openshiftNetSGNamePtr,
						},
					},
				},
			},
			{
				// IPsec IKE NAT-T - port 4500
				Action:    capibmcloud.VPCSecurityGroupRuleActionAllow,
				Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
				Source: &capibmcloud.VPCSecurityGroupRulePrototype{
					PortRange: &capibmcloud.VPCSecurityGroupPortRange{
						MaximumPort: 4500,
						MinimumPort: 4500,
					},
					Protocol: capibmcloud.VPCSecurityGroupRuleProtocolUDP,
					Remotes: []capibmcloud.VPCSecurityGroupRuleRemote{
						{
							RemoteType:        capibmcloud.VPCSecurityGroupRuleRemoteTypeSG,
							SecurityGroupName: openshiftNetSGNamePtr,
						},
					},
				},
			},
			{
				// Kubernetes node ports - TCP
				// Allows access to node ports from within VPC subnets to accommodate CCM LBs
				Action:    capibmcloud.VPCSecurityGroupRuleActionAllow,
				Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
				Source: &capibmcloud.VPCSecurityGroupRulePrototype{
					PortRange: &capibmcloud.VPCSecurityGroupPortRange{
						MaximumPort: 32767,
						MinimumPort: 30000,
					},
					Protocol: capibmcloud.VPCSecurityGroupRuleProtocolTCP,
					Remotes:  openshiftNetworkNodePortTCPRemotes,
				},
			},
			{
				// Kubernetes node ports - UDP
				// Allows access to node ports from within VPC subnets to accommodate CCM LBs
				Action:    capibmcloud.VPCSecurityGroupRuleActionAllow,
				Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
				Source: &capibmcloud.VPCSecurityGroupRulePrototype{
					PortRange: &capibmcloud.VPCSecurityGroupPortRange{
						MaximumPort: 32767,
						MinimumPort: 30000,
					},
					Protocol: capibmcloud.VPCSecurityGroupRuleProtocolUDP,
					Remotes:  openshiftNetworkNodePortUDPRemotes,
				},
			},
		},
	}
}

func buildKubeAPILBSecurityGroup(infraID string) capibmcloud.VPCSecurityGroup {
	kubeAPILBSGNamePtr := ptr.To(fmt.Sprintf("%s-%s", infraID, kubeAPILBSGNameSuffix))
	controlPlaneSGNamePtr := ptr.To(fmt.Sprintf("%s-%s", infraID, controlPlaneSGNameSuffix))
	clusterWideSGNamePtr := ptr.To(fmt.Sprintf("%s-%s", infraID, clusterWideSGNameSuffix))

	return capibmcloud.VPCSecurityGroup{
		Name: kubeAPILBSGNamePtr,
		Rules: []*capibmcloud.VPCSecurityGroupRule{
			{
				// Kubernetes API LB - inbound
				Action:    capibmcloud.VPCSecurityGroupRuleActionAllow,
				Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
				Source: &capibmcloud.VPCSecurityGroupRulePrototype{
					PortRange: &capibmcloud.VPCSecurityGroupPortRange{
						MaximumPort: 6443,
						MinimumPort: 6443,
					},
					Protocol: capibmcloud.VPCSecurityGroupRuleProtocolTCP,
					Remotes: []capibmcloud.VPCSecurityGroupRuleRemote{
						{
							RemoteType: capibmcloud.VPCSecurityGroupRuleRemoteTypeAny,
						},
					},
				},
			},
			{
				// Kubernetes API LB - outbound
				Action: capibmcloud.VPCSecurityGroupRuleActionAllow,
				Destination: &capibmcloud.VPCSecurityGroupRulePrototype{
					PortRange: &capibmcloud.VPCSecurityGroupPortRange{
						MaximumPort: 6443,
						MinimumPort: 6443,
					},
					Protocol: capibmcloud.VPCSecurityGroupRuleProtocolTCP,
					Remotes: []capibmcloud.VPCSecurityGroupRuleRemote{
						{
							RemoteType:        capibmcloud.VPCSecurityGroupRuleRemoteTypeSG,
							SecurityGroupName: controlPlaneSGNamePtr,
						},
					},
				},
				Direction: capibmcloud.VPCSecurityGroupRuleDirectionOutbound,
			},
			{
				// Machine Config Server LB - inbound
				Action:    capibmcloud.VPCSecurityGroupRuleActionAllow,
				Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
				Source: &capibmcloud.VPCSecurityGroupRulePrototype{
					PortRange: &capibmcloud.VPCSecurityGroupPortRange{
						MaximumPort: 22623,
						MinimumPort: 22623,
					},
					Protocol: capibmcloud.VPCSecurityGroupRuleProtocolTCP,
					Remotes: []capibmcloud.VPCSecurityGroupRuleRemote{
						{
							RemoteType:        capibmcloud.VPCSecurityGroupRuleRemoteTypeSG,
							SecurityGroupName: clusterWideSGNamePtr,
						},
					},
				},
			},
			{
				// Machine Config Server LB - outbound
				Action: capibmcloud.VPCSecurityGroupRuleActionAllow,
				Destination: &capibmcloud.VPCSecurityGroupRulePrototype{
					PortRange: &capibmcloud.VPCSecurityGroupPortRange{
						MaximumPort: 22623,
						MinimumPort: 22623,
					},
					Protocol: capibmcloud.VPCSecurityGroupRuleProtocolTCP,
					Remotes: []capibmcloud.VPCSecurityGroupRuleRemote{
						{
							RemoteType:        capibmcloud.VPCSecurityGroupRuleRemoteTypeSG,
							SecurityGroupName: controlPlaneSGNamePtr,
						},
					},
				},
				Direction: capibmcloud.VPCSecurityGroupRuleDirectionOutbound,
			},
		},
	}
}

func buildControlPlaneSecurityGroup(infraID string) capibmcloud.VPCSecurityGroup {
	controlPlaneSGNamePtr := ptr.To(fmt.Sprintf("%s-%s", infraID, controlPlaneSGNameSuffix))
	clusterWideSGNamePtr := ptr.To(fmt.Sprintf("%s-%s", infraID, clusterWideSGNameSuffix))
	kubeAPILBSGNamePtr := ptr.To(fmt.Sprintf("%s-%s", infraID, kubeAPILBSGNameSuffix))

	return capibmcloud.VPCSecurityGroup{
		Name: controlPlaneSGNamePtr,
		Rules: []*capibmcloud.VPCSecurityGroupRule{
			{
				// Kubernetes API - inbound via cluster
				Action:    capibmcloud.VPCSecurityGroupRuleActionAllow,
				Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
				Source: &capibmcloud.VPCSecurityGroupRulePrototype{
					PortRange: &capibmcloud.VPCSecurityGroupPortRange{
						MaximumPort: 6443,
						MinimumPort: 6443,
					},
					Protocol: capibmcloud.VPCSecurityGroupRuleProtocolTCP,
					Remotes: []capibmcloud.VPCSecurityGroupRuleRemote{
						{
							RemoteType:        capibmcloud.VPCSecurityGroupRuleRemoteTypeSG,
							SecurityGroupName: clusterWideSGNamePtr,
						},
					},
				},
			},
			{
				// Kubernetes API - inbound via LB
				Action:    capibmcloud.VPCSecurityGroupRuleActionAllow,
				Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
				Source: &capibmcloud.VPCSecurityGroupRulePrototype{
					PortRange: &capibmcloud.VPCSecurityGroupPortRange{
						MaximumPort: 6443,
						MinimumPort: 6443,
					},
					Protocol: capibmcloud.VPCSecurityGroupRuleProtocolTCP,
					Remotes: []capibmcloud.VPCSecurityGroupRuleRemote{
						{
							RemoteType:        capibmcloud.VPCSecurityGroupRuleRemoteTypeSG,
							SecurityGroupName: kubeAPILBSGNamePtr,
						},
					},
				},
			},
			{
				// Machine Config Server - inbound via LB
				Action:    capibmcloud.VPCSecurityGroupRuleActionAllow,
				Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
				Source: &capibmcloud.VPCSecurityGroupRulePrototype{
					PortRange: &capibmcloud.VPCSecurityGroupPortRange{
						MaximumPort: 22623,
						MinimumPort: 22623,
					},
					Protocol: capibmcloud.VPCSecurityGroupRuleProtocolTCP,
					Remotes: []capibmcloud.VPCSecurityGroupRuleRemote{
						{
							RemoteType:        capibmcloud.VPCSecurityGroupRuleRemoteTypeSG,
							SecurityGroupName: kubeAPILBSGNamePtr,
						},
					},
				},
			},
			{
				// Kubernetes default ports
				Action:    capibmcloud.VPCSecurityGroupRuleActionAllow,
				Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
				Source: &capibmcloud.VPCSecurityGroupRulePrototype{
					PortRange: &capibmcloud.VPCSecurityGroupPortRange{
						MaximumPort: 10259,
						MinimumPort: 10257,
					},
					Protocol: capibmcloud.VPCSecurityGroupRuleProtocolTCP,
					Remotes: []capibmcloud.VPCSecurityGroupRuleRemote{
						{
							RemoteType:        capibmcloud.VPCSecurityGroupRuleRemoteTypeSG,
							SecurityGroupName: clusterWideSGNamePtr,
						},
					},
				},
			},
		},
	}
}

func buildCPInternalSecurityGroup(infraID string) capibmcloud.VPCSecurityGroup {
	cpInternalSGNamePtr := ptr.To(fmt.Sprintf("%s-%s", infraID, cpInternalSGNameSuffix))

	return capibmcloud.VPCSecurityGroup{
		Name: cpInternalSGNamePtr,
		Rules: []*capibmcloud.VPCSecurityGroupRule{
			{
				// etcd internal traffic
				Action:    capibmcloud.VPCSecurityGroupRuleActionAllow,
				Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
				Source: &capibmcloud.VPCSecurityGroupRulePrototype{
					PortRange: &capibmcloud.VPCSecurityGroupPortRange{
						MaximumPort: 2380,
						MinimumPort: 2379,
					},
					Protocol: capibmcloud.VPCSecurityGroupRuleProtocolTCP,
					Remotes: []capibmcloud.VPCSecurityGroupRuleRemote{
						{
							RemoteType:        capibmcloud.VPCSecurityGroupRuleRemoteTypeSG,
							SecurityGroupName: cpInternalSGNamePtr,
						},
					},
				},
			},
		},
	}
}

func buildBootstrapSecurityGroup(infraID string, allSubnets []capibmcloud.Subnet, publishStrategy types.PublishingStrategy) capibmcloud.VPCSecurityGroup {
	bootstrapSGNamePtr := ptr.To(fmt.Sprintf("%s-%s", infraID, ibmcloudic.BootstrapSGNameSuffix))

	remotes := make([]capibmcloud.VPCSecurityGroupRuleRemote, 0, len(allSubnets))

	// Based on the Cluster's PublishingStrategy, SSH access to the bootstrap node is restricted to within the Cluster subnets' CIDR's, or simply publicly available.
	switch publishStrategy {
	case types.ExternalPublishingStrategy:
		remotes = append(remotes, capibmcloud.VPCSecurityGroupRuleRemote{
			RemoteType: capibmcloud.VPCSecurityGroupRuleRemoteTypeAny,
		})
	case types.InternalPublishingStrategy:
		for _, subnet := range allSubnets {
			remotes = append(remotes, capibmcloud.VPCSecurityGroupRuleRemote{
				RemoteType:     capibmcloud.VPCSecurityGroupRuleRemoteTypeCIDR,
				CIDRSubnetName: subnet.Name,
			})
		}
	}
	return capibmcloud.VPCSecurityGroup{
		Name: bootstrapSGNamePtr,
		Rules: []*capibmcloud.VPCSecurityGroupRule{
			{
				// SSH traffic
				Action:    capibmcloud.VPCSecurityGroupRuleActionAllow,
				Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
				Source: &capibmcloud.VPCSecurityGroupRulePrototype{
					PortRange: &capibmcloud.VPCSecurityGroupPortRange{
						MaximumPort: 22,
						MinimumPort: 22,
					},
					Protocol: capibmcloud.VPCSecurityGroupRuleProtocolTCP,
					Remotes:  remotes,
				},
			},
		},
	}
}

func getVPCSecurityGroups(infraID string, allSubnets []capibmcloud.Subnet, publishStrategy types.PublishingStrategy) []capibmcloud.VPCSecurityGroup {
	// IBM Cloud currently relies on 5 SecurityGroups to manage traffic and 1 SecurityGroup for bootstrapping.
	securityGroups := make([]capibmcloud.VPCSecurityGroup, 0, 6)
	// Generate the Cluster's primary SG's.
	securityGroups = append(securityGroups, buildClusterWideSecurityGroup(infraID, allSubnets))
	securityGroups = append(securityGroups, buildOpenshiftNetSecurityGroup(infraID, allSubnets))
	securityGroups = append(securityGroups, buildKubeAPILBSecurityGroup(infraID))
	securityGroups = append(securityGroups, buildControlPlaneSecurityGroup(infraID))
	securityGroups = append(securityGroups, buildCPInternalSecurityGroup(infraID))

	// Generate the bootstrap SG.
	securityGroups = append(securityGroups, buildBootstrapSecurityGroup(infraID, allSubnets, publishStrategy))

	return securityGroups
}
