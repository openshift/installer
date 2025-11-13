package powervs

import (
	"fmt"

	"k8s.io/utils/ptr"
	capibmcloud "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"

	"github.com/openshift/installer/pkg/types"
)

const (
	bootstrapSGNameSuffix    = "sg-bootstrap"
	clusterWideSGNameSuffix  = "sg-clusterwide"
	controlPlaneSGNameSuffix = "sg-control-plane"
	cpInternalSGNameSuffix   = "sg-cp-internal"
	kubeAPILBSGNameSuffix    = "sg-kube-api-lb"
)

func buildBootstrapSecurityGroup(infraID string) capibmcloud.VPCSecurityGroup {
	bootstrapSGNamePtr := ptr.To(fmt.Sprintf("%s-%s", infraID, bootstrapSGNameSuffix))
	return capibmcloud.VPCSecurityGroup{
		Name: bootstrapSGNamePtr,
		Rules: []*capibmcloud.VPCSecurityGroupRule{
			{
				// SSH inbound bootstrap
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
		},
	}
}

func buildClusterWideSecurityGroup(infraID string) capibmcloud.VPCSecurityGroup {
	clusterWideSGNamePtr := ptr.To(fmt.Sprintf("%s-%s", infraID, clusterWideSGNameSuffix))
	return capibmcloud.VPCSecurityGroup{
		Name: clusterWideSGNamePtr,
	}
}

func buildControlPlaneSecurityGroup(infraID string) capibmcloud.VPCSecurityGroup {
	controlPlaneSGNamePtr := ptr.To(fmt.Sprintf("%s-%s", infraID, controlPlaneSGNameSuffix))
	return capibmcloud.VPCSecurityGroup{
		Name: controlPlaneSGNamePtr,
	}
}

func buildCPInternalSecurityGroup(infraID string) capibmcloud.VPCSecurityGroup {
	cpInternalSGNamePtr := ptr.To(fmt.Sprintf("%s-%s", infraID, cpInternalSGNameSuffix))
	return capibmcloud.VPCSecurityGroup{
		Name: cpInternalSGNamePtr,
	}
}

func buildKubeAPILBSecurityGroup(infraID string) capibmcloud.VPCSecurityGroup {
	kubeAPILBSGNamePtr := ptr.To(fmt.Sprintf("%s-%s", infraID, kubeAPILBSGNameSuffix))
	return capibmcloud.VPCSecurityGroup{
		Name: kubeAPILBSGNamePtr,
		Rules: []*capibmcloud.VPCSecurityGroupRule{
			{
				// SSH inbound bootstrap
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

func getVPCSecurityGroups(infraID string, publishStrategy types.PublishingStrategy) []capibmcloud.VPCSecurityGroup {
	// IBM Power VS will rely on 6 SecurityGroups to manage traffic and 1 SecurityGroup for bootstrapping.
	securityGroups := make([]capibmcloud.VPCSecurityGroup, 0, 6)
	securityGroups = append(securityGroups, buildKubeAPILBSecurityGroup(infraID))
	securityGroups = append(securityGroups, buildBootstrapSecurityGroup(infraID))
	securityGroups = append(securityGroups, buildClusterWideSecurityGroup(infraID))
	securityGroups = append(securityGroups, buildControlPlaneSecurityGroup(infraID))
	securityGroups = append(securityGroups, buildCPInternalSecurityGroup(infraID))
	return securityGroups
}
