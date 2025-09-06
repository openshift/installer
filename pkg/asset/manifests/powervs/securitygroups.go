package powervs

import (
	"fmt"

	"k8s.io/utils/ptr"
	capibmcloud "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"

	"github.com/openshift/installer/pkg/types"
)

const (
	kubeAPILBSGNameSuffix = "sg-kube-api-lb"
)

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
	return securityGroups
}
