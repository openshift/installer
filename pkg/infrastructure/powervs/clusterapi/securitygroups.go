package clusterapi

import (
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"k8s.io/utils/ptr"
)

type Rules interface {
	GetRules(SGIDCollection) []*vpcv1.SecurityGroupRulePrototype
}
type ClusterWide struct{}
type ControlPlane struct{}
type CPInternal struct{}
type KubeAPILB struct{}
type OpenShiftNet struct{}

func GetSGRules(r Rules, sgIDs SGIDCollection) []*vpcv1.SecurityGroupRulePrototype {
	return r.GetRules(sgIDs)
}

func (ClusterWide) GetRules(sgIDs SGIDCollection) []*vpcv1.SecurityGroupRulePrototype {
	return []*vpcv1.SecurityGroupRulePrototype{
		{
			Direction: ptr.To("inbound"),
			Protocol:  ptr.To("icmp"),
			Remote: &vpcv1.SecurityGroupRuleRemotePrototypeSecurityGroupIdentitySecurityGroupIdentityByID{
				ID: &sgIDs.ClusterWideSGID,
			},
		},
		{
			Direction: ptr.To("inbound"),
			Protocol:  ptr.To("udp"),
			Remote: &vpcv1.SecurityGroupRuleRemotePrototypeSecurityGroupIdentitySecurityGroupIdentityByID{
				ID: &sgIDs.ClusterWideSGID,
			},
			PortMin: ptr.To(int64(4789)),
			PortMax: ptr.To(int64(4789)),
		},
		{
			Direction: ptr.To("inbound"),
			Protocol:  ptr.To("udp"),
			Remote: &vpcv1.SecurityGroupRuleRemotePrototypeSecurityGroupIdentitySecurityGroupIdentityByID{
				ID: &sgIDs.ClusterWideSGID,
			},
			PortMin: ptr.To(int64(6081)),
			PortMax: ptr.To(int64(6081)),
		},
	}
}

func (ControlPlane) GetRules(sgIDs SGIDCollection) []*vpcv1.SecurityGroupRulePrototype {
	return []*vpcv1.SecurityGroupRulePrototype{
		{
			Direction: ptr.To("inbound"),
			Protocol:  ptr.To("tcp"),
			Remote: &vpcv1.SecurityGroupRuleRemotePrototypeSecurityGroupIdentitySecurityGroupIdentityByID{
				ID: &sgIDs.ClusterWideSGID,
			},
			PortMin: ptr.To(int64(5443)),
			PortMax: ptr.To(int64(5443)),
		},
		{
			Direction: ptr.To("inbound"),
			Protocol:  ptr.To("tcp"),
			Remote: &vpcv1.SecurityGroupRuleRemotePrototypeSecurityGroupIdentitySecurityGroupIdentityByID{
				ID: &sgIDs.ClusterWideSGID,
			},
			PortMin: ptr.To(int64(10257)),
			PortMax: ptr.To(int64(10259)),
		},
		{
			Direction: ptr.To("inbound"),
			Protocol:  ptr.To("tcp"),
			Remote: &vpcv1.SecurityGroupRuleRemotePrototypeSecurityGroupIdentitySecurityGroupIdentityByID{
				ID: &sgIDs.KubeAPILBSGID,
			},
			PortMin: ptr.To(int64(6443)),
			PortMax: ptr.To(int64(6443)),
		},
		{
			Direction: ptr.To("inbound"),
			Protocol:  ptr.To("tcp"),
			Remote: &vpcv1.SecurityGroupRuleRemotePrototypeSecurityGroupIdentitySecurityGroupIdentityByID{
				ID: &sgIDs.KubeAPILBSGID,
			},
			PortMin: ptr.To(int64(22623)),
			PortMax: ptr.To(int64(22623)),
		},
	}
}

func (CPInternal) GetRules(sgIDs SGIDCollection) []*vpcv1.SecurityGroupRulePrototype {
	return []*vpcv1.SecurityGroupRulePrototype{
		{
			Direction: ptr.To("inbound"),
			Protocol:  ptr.To("tcp"),
			Remote: &vpcv1.SecurityGroupRuleRemotePrototypeSecurityGroupIdentitySecurityGroupIdentityByID{
				ID: &sgIDs.CPInternalSGID,
			},
			PortMin: ptr.To(int64(2379)),
			PortMax: ptr.To(int64(2380)),
		},
	}
}

func (KubeAPILB) GetRules(sgIDs SGIDCollection) []*vpcv1.SecurityGroupRulePrototype {
	return []*vpcv1.SecurityGroupRulePrototype{
		{
			Direction: ptr.To("inbound"),
			Protocol:  ptr.To("tcp"),
			Remote: &vpcv1.SecurityGroupRuleRemotePrototypeSecurityGroupIdentitySecurityGroupIdentityByID{
				ID: &sgIDs.ClusterWideSGID,
			},
			PortMin: ptr.To(int64(22623)),
			PortMax: ptr.To(int64(22623)),
		},
	}
}

func (OpenShiftNet) GetRules(sgIDs SGIDCollection) []*vpcv1.SecurityGroupRulePrototype {
	return []*vpcv1.SecurityGroupRulePrototype{
		{
			Direction: ptr.To("inbound"),
			Protocol:  ptr.To("tcp"),
			Remote: &vpcv1.SecurityGroupRuleRemotePrototypeSecurityGroupIdentitySecurityGroupIdentityByID{
				ID: &sgIDs.ClusterWideSGID,
			},
		},
		{
			Direction: ptr.To("inbound"),
			Protocol:  ptr.To("udp"),
			Remote: &vpcv1.SecurityGroupRuleRemotePrototypeSecurityGroupIdentitySecurityGroupIdentityByID{
				ID: &sgIDs.ClusterWideSGID,
			},
			PortMin: ptr.To(int64(4789)),
			PortMax: ptr.To(int64(4789)),
		},
		{
			Direction: ptr.To("inbound"),
			Protocol:  ptr.To("udp"),
			Remote: &vpcv1.SecurityGroupRuleRemotePrototypeSecurityGroupIdentitySecurityGroupIdentityByID{
				ID: &sgIDs.ClusterWideSGID,
			},
			PortMin: ptr.To(int64(6081)),
			PortMax: ptr.To(int64(6081)),
		},
	}
}
