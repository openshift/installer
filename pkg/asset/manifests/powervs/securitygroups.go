package powervs

import (
	"fmt"

	capibm "sigs.k8s.io/cluster-api-provider-ibmcloud/api/powervs/v1beta3"
)

const (
	controlPlaneSGNameSuffix = "sg-control-plane"
	clusterWideSGNameSuffix  = "sg-cluster-wide"
	kubeAPILBSGNameSuffix    = "sg-kube-api-lb"
)

func sgInboundTCP(minPort, maxPort int64) capibm.VPCSecurityGroupRule {
	return capibm.VPCSecurityGroupRule{
		Direction: capibm.VPCSecurityGroupRuleDirectionInbound,
		Source: capibm.VPCSecurityGroupRulePrototype{
			Protocol: capibm.VPCSecurityGroupRuleProtocolTCP,
			PortRange: capibm.VPCSecurityGroupPortRange{
				MinimumPort: minPort,
				MaximumPort: maxPort,
			},
			Remotes: []capibm.VPCSecurityGroupRuleRemote{
				{RemoteType: capibm.VPCSecurityGroupRuleRemoteTypeAny},
			},
		},
	}
}

func buildControlPlaneSecurityGroup(infraID string) capibm.VPCSecurityGroupSource {
	return capibm.VPCSecurityGroupSource{
		Type: capibm.SourceTypeProvision,
		Provision: capibm.VPCSecurityGroupProvision{
			Name: fmt.Sprintf("%s-%s", infraID, controlPlaneSGNameSuffix),
			Rules: []capibm.VPCSecurityGroupRule{
				sgInboundTCP(10258, 10258),
				sgInboundTCP(22623, 22623),
				// Konnectivity
				sgInboundTCP(8091, 8091),
				sgInboundTCP(443, 443),
			},
		},
	}
}

func buildKubeAPILBSecurityGroup(infraID string) capibm.VPCSecurityGroupSource {
	return capibm.VPCSecurityGroupSource{
		Type: capibm.SourceTypeProvision,
		Provision: capibm.VPCSecurityGroupProvision{
			Name: fmt.Sprintf("%s-%s", infraID, kubeAPILBSGNameSuffix),
			Rules: []capibm.VPCSecurityGroupRule{
				sgInboundTCP(6443, 6443),
			},
		},
	}
}

func buildClusterWideSecurityGroup(infraID string) capibm.VPCSecurityGroupSource {
	return capibm.VPCSecurityGroupSource{
		Type: capibm.SourceTypeProvision,
		Provision: capibm.VPCSecurityGroupProvision{
			Name: fmt.Sprintf("%s-%s", infraID, clusterWideSGNameSuffix),
			Rules: []capibm.VPCSecurityGroupRule{
				// SSH inbound
				sgInboundTCP(22, 22),
				// Image registry
				sgInboundTCP(5000, 5000),
				// ping (ICMP inbound)
				{
					Direction: capibm.VPCSecurityGroupRuleDirectionInbound,
					Source: capibm.VPCSecurityGroupRulePrototype{
						Protocol: capibm.VPCSecurityGroupRuleProtocolIcmp,
						Remotes: []capibm.VPCSecurityGroupRuleRemote{
							{RemoteType: capibm.VPCSecurityGroupRuleRemoteTypeAny},
						},
					},
				},
				// Outbound — all protocols
				{
					Direction: capibm.VPCSecurityGroupRuleDirectionOutbound,
					Destination: capibm.VPCSecurityGroupRulePrototype{
						Protocol: capibm.VPCSecurityGroupRuleProtocolAll,
						Remotes: []capibm.VPCSecurityGroupRuleRemote{
							{RemoteType: capibm.VPCSecurityGroupRuleRemoteTypeAny},
						},
					},
				},
			},
		},
	}
}

func getVPCSecurityGroups(infraID string) []capibm.VPCSecurityGroupSource {
	return []capibm.VPCSecurityGroupSource{
		buildClusterWideSecurityGroup(infraID),
		buildControlPlaneSecurityGroup(infraID),
		buildKubeAPILBSecurityGroup(infraID),
	}
}
