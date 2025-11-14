package powervs

import (
	"fmt"

	"k8s.io/utils/ptr"
	capibmcloud "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
)

const (
	controlPlaneSGNameSuffix = "sg-control-plane"
	clusterWideSGNameSuffix  = "sg-cluster-wide"
	kubeAPILBSGNameSuffix    = "sg-kube-api-lb"
)

func buildControlPlaneSecurityGroup(infraID string) capibmcloud.VPCSecurityGroup {
	kubeAPILBSGNamePtr := ptr.To(fmt.Sprintf("%s-%s", infraID, controlPlaneSGNameSuffix))
	return capibmcloud.VPCSecurityGroup{
		Name: kubeAPILBSGNamePtr,
		Rules: []*capibmcloud.VPCSecurityGroupRule{
			{
				Action:    capibmcloud.VPCSecurityGroupRuleActionAllow,
				Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
				Source: &capibmcloud.VPCSecurityGroupRulePrototype{
					PortRange: &capibmcloud.VPCSecurityGroupPortRange{
						MaximumPort: 10258,
						MinimumPort: 10258,
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
							RemoteType: capibmcloud.VPCSecurityGroupRuleRemoteTypeAny,
						},
					},
				},
			},
			{
				Action:    capibmcloud.VPCSecurityGroupRuleActionAllow,
				Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
				Source: &capibmcloud.VPCSecurityGroupRulePrototype{
					PortRange: &capibmcloud.VPCSecurityGroupPortRange{
						MaximumPort: 443,
						MinimumPort: 443,
					},
					Protocol: capibmcloud.VPCSecurityGroupRuleProtocolTCP,
					Remotes: []capibmcloud.VPCSecurityGroupRuleRemote{
						{
							RemoteType: capibmcloud.VPCSecurityGroupRuleRemoteTypeAny,
						},
					},
				},
			},
		},
	}
}

func buildKubeAPILBSecurityGroup(infraID string) capibmcloud.VPCSecurityGroup {
	kubeAPILBSGNamePtr := ptr.To(fmt.Sprintf("%s-%s", infraID, kubeAPILBSGNameSuffix))
	return capibmcloud.VPCSecurityGroup{
		Name: kubeAPILBSGNamePtr,
		Rules: []*capibmcloud.VPCSecurityGroupRule{
			{
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
		},
	}
}

func buildClusterWideSecurityGroup(infraID string) capibmcloud.VPCSecurityGroup {
	kubeAPILBSGNamePtr := ptr.To(fmt.Sprintf("%s-%s", infraID, clusterWideSGNameSuffix))
	return capibmcloud.VPCSecurityGroup{
		Name: kubeAPILBSGNamePtr,
		Rules: []*capibmcloud.VPCSecurityGroupRule{
			{
				// SSH inbound
				Action:    capibmcloud.VPCSecurityGroupRuleActionAllow,
				Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
				Source: &capibmcloud.VPCSecurityGroupRulePrototype{
					PortRange: &capibmcloud.VPCSecurityGroupPortRange{
						MaximumPort: 22,
						MinimumPort: 22,
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
				Action:    capibmcloud.VPCSecurityGroupRuleActionAllow,
				Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
				Source: &capibmcloud.VPCSecurityGroupRulePrototype{
					PortRange: &capibmcloud.VPCSecurityGroupPortRange{
						MaximumPort: 5000,
						MinimumPort: 5000,
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
				// ping
				Action:    capibmcloud.VPCSecurityGroupRuleActionAllow,
				Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
				Source: &capibmcloud.VPCSecurityGroupRulePrototype{
					Protocol: capibmcloud.VPCSecurityGroupRuleProtocolIcmp,
					Remotes: []capibmcloud.VPCSecurityGroupRuleRemote{
						{
							RemoteType: capibmcloud.VPCSecurityGroupRuleRemoteTypeAny,
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

func getVPCSecurityGroups(infraID string) []capibmcloud.VPCSecurityGroup {
	// IBM Power VS will rely on 3 SecurityGroups to manage traffic.
	securityGroups := make([]capibmcloud.VPCSecurityGroup, 0, 3)
	securityGroups = append(securityGroups, buildClusterWideSecurityGroup(infraID))
	securityGroups = append(securityGroups, buildControlPlaneSecurityGroup(infraID))
	securityGroups = append(securityGroups, buildKubeAPILBSecurityGroup(infraID))
	return securityGroups
}
