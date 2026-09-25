package powervs

import (
	"fmt"

	capibmcloud "sigs.k8s.io/cluster-api-provider-ibmcloud/api/powervs/v1beta3"
)

const (
	controlPlaneSGNameSuffix = "sg-control-plane"
	clusterWideSGNameSuffix  = "sg-cluster-wide"
	kubeAPILBSGNameSuffix    = "sg-kube-api-lb"
)

func buildControlPlaneSecurityGroup(infraID string) capibmcloud.VPCSecurityGroupSource {
	kubeAPILBSGName := fmt.Sprintf("%s-%s", infraID, controlPlaneSGNameSuffix)
	return capibmcloud.VPCSecurityGroupSource{
		Type: capibmcloud.SourceTypeProvision,
		Provision: capibmcloud.VPCSecurityGroupProvision{
			Name: kubeAPILBSGName,
			Rules: []capibmcloud.VPCSecurityGroupRule{
				{
					Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
					Source: capibmcloud.VPCSecurityGroupRulePrototype{
						PortRange: capibmcloud.VPCSecurityGroupPortRange{
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
					Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
					Source: capibmcloud.VPCSecurityGroupRulePrototype{
						PortRange: capibmcloud.VPCSecurityGroupPortRange{
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
					// Konnectivity
					Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
					Source: capibmcloud.VPCSecurityGroupRulePrototype{
						PortRange: capibmcloud.VPCSecurityGroupPortRange{
							MaximumPort: 8091,
							MinimumPort: 8091,
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
					Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
					Source: capibmcloud.VPCSecurityGroupRulePrototype{
						PortRange: capibmcloud.VPCSecurityGroupPortRange{
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
		},
	}
}

func buildKubeAPILBSecurityGroup(infraID string) capibmcloud.VPCSecurityGroupSource {
	kubeAPILBSGName := fmt.Sprintf("%s-%s", infraID, kubeAPILBSGNameSuffix)
	return capibmcloud.VPCSecurityGroupSource{
		Type: capibmcloud.SourceTypeProvision,
		Provision: capibmcloud.VPCSecurityGroupProvision{
			Name: kubeAPILBSGName,
			Rules: []capibmcloud.VPCSecurityGroupRule{
				{
					Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
					Source: capibmcloud.VPCSecurityGroupRulePrototype{
						PortRange: capibmcloud.VPCSecurityGroupPortRange{
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
		},
	}
}

func buildClusterWideSecurityGroup(infraID string) capibmcloud.VPCSecurityGroupSource {
	clusterWideSGName := fmt.Sprintf("%s-%s", infraID, clusterWideSGNameSuffix)
	return capibmcloud.VPCSecurityGroupSource{
		Type: capibmcloud.SourceTypeProvision,
		Provision: capibmcloud.VPCSecurityGroupProvision{
			Name: clusterWideSGName,
			Rules: []capibmcloud.VPCSecurityGroupRule{
				{
					// SSH inbound
					Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
					Source: capibmcloud.VPCSecurityGroupRulePrototype{
						PortRange: capibmcloud.VPCSecurityGroupPortRange{
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
					Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
					Source: capibmcloud.VPCSecurityGroupRulePrototype{
						PortRange: capibmcloud.VPCSecurityGroupPortRange{
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
					Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
					Source: capibmcloud.VPCSecurityGroupRulePrototype{
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
					Destination: capibmcloud.VPCSecurityGroupRulePrototype{
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
		},
	}
}

func getVPCSecurityGroups(infraID string) []capibmcloud.VPCSecurityGroupSource {
	// IBM Power VS will rely on 3 SecurityGroups to manage traffic.
	securityGroups := make([]capibmcloud.VPCSecurityGroupSource, 0, 3)
	securityGroups = append(securityGroups, buildClusterWideSecurityGroup(infraID))
	securityGroups = append(securityGroups, buildControlPlaneSecurityGroup(infraID))
	securityGroups = append(securityGroups, buildKubeAPILBSecurityGroup(infraID))
	return securityGroups
}
