/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package networking

var defaultRules = []resolvedSecurityGroupRuleSpec{
	{
		Direction:      "egress",
		Description:    "Full open",
		EtherType:      "IPv4",
		PortRangeMin:   0,
		PortRangeMax:   0,
		Protocol:       "",
		RemoteIPPrefix: "",
	},
	{
		Direction:      "egress",
		Description:    "Full open",
		EtherType:      "IPv6",
		PortRangeMin:   0,
		PortRangeMax:   0,
		Protocol:       "",
		RemoteIPPrefix: "",
	},
}

// Permit traffic for etcd, kubelet.
func getSGControlPlaneCommon(remoteGroupIDSelf, secWorkerGroupID string) []resolvedSecurityGroupRuleSpec {
	return []resolvedSecurityGroupRuleSpec{
		{
			Description:   "Etcd",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  2379,
			PortRangeMax:  2380,
			Protocol:      "tcp",
			RemoteGroupID: remoteGroupIDSelf,
		},
		{
			// kubeadm says this is needed
			Description:   "Kubelet API",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  10250,
			PortRangeMax:  10250,
			Protocol:      "tcp",
			RemoteGroupID: remoteGroupIDSelf,
		},
		{
			// This is needed to support metrics-server deployments
			Description:   "Kubelet API",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  10250,
			PortRangeMax:  10250,
			Protocol:      "tcp",
			RemoteGroupID: secWorkerGroupID,
		},
	}
}

// Permit traffic for kubelet.
func getSGWorkerCommon(remoteGroupIDSelf, secControlPlaneGroupID string) []resolvedSecurityGroupRuleSpec {
	return []resolvedSecurityGroupRuleSpec{
		{
			// This is needed to support metrics-server deployments
			Description:   "Kubelet API",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  10250,
			PortRangeMax:  10250,
			Protocol:      "tcp",
			RemoteGroupID: remoteGroupIDSelf,
		},
		{
			Description:   "Kubelet API",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  10250,
			PortRangeMax:  10250,
			Protocol:      "tcp",
			RemoteGroupID: secControlPlaneGroupID,
		},
	}
}

// Permit traffic for ssh control plane.
func getSGControlPlaneSSH(secBastionGroupID string) []resolvedSecurityGroupRuleSpec {
	return []resolvedSecurityGroupRuleSpec{
		{
			Description:   "SSH",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  22,
			PortRangeMax:  22,
			Protocol:      "tcp",
			RemoteGroupID: secBastionGroupID,
		},
	}
}

// Permit traffic for ssh worker.
func getSGWorkerSSH(secBastionGroupID string) []resolvedSecurityGroupRuleSpec {
	return []resolvedSecurityGroupRuleSpec{
		{
			Description:   "SSH",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  22,
			PortRangeMax:  22,
			Protocol:      "tcp",
			RemoteGroupID: secBastionGroupID,
		},
	}
}

// Allow all traffic, including from outside the cluster, to access the API.
func getSGControlPlaneHTTPS() []resolvedSecurityGroupRuleSpec {
	return []resolvedSecurityGroupRuleSpec{
		{
			Description:  "Kubernetes API",
			Direction:    "ingress",
			EtherType:    "IPv4",
			PortRangeMin: 6443,
			PortRangeMax: 6443,
			Protocol:     "tcp",
		},
	}
}

// Allow all traffic, including from outside the cluster, to access node port services.
func getSGWorkerNodePort(secWorkerGroupID string, secControlPlaneGroupID string) []resolvedSecurityGroupRuleSpec {
	return []resolvedSecurityGroupRuleSpec{
		{
			Description:   "Node Port Services",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  30000,
			PortRangeMax:  32767,
			Protocol:      "tcp",
			RemoteGroupID: secWorkerGroupID,
		},
		{
			Description:   "Node Port Services",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  30000,
			PortRangeMax:  32767,
			Protocol:      "udp",
			RemoteGroupID: secWorkerGroupID,
		},
		{
			Description:   "Node Port Services",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  30000,
			PortRangeMax:  32767,
			Protocol:      "tcp",
			RemoteGroupID: secControlPlaneGroupID,
		},
		{
			Description:   "Node Port Services",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  30000,
			PortRangeMax:  32767,
			Protocol:      "udp",
			RemoteGroupID: secControlPlaneGroupID,
		},
	}
}

// Allow all traffic from a specific CIDR to access node port services.
func getSGWorkerNodePortCIDR(cidr string) []resolvedSecurityGroupRuleSpec {
	return []resolvedSecurityGroupRuleSpec{
		{
			Description:    "Node Port Services",
			Direction:      "ingress",
			EtherType:      "IPv4",
			PortRangeMin:   30000,
			PortRangeMax:   32767,
			Protocol:       "tcp",
			RemoteIPPrefix: cidr,
		},
		{
			Description:    "Node Port Services",
			Direction:      "ingress",
			EtherType:      "IPv4",
			PortRangeMin:   30000,
			PortRangeMax:   32767,
			Protocol:       "udp",
			RemoteIPPrefix: cidr,
		},
	}
}

// Permit all ingress from the cluster security groups.
func getSGControlPlaneAllowAll(remoteGroupIDSelf, secWorkerGroupID string) []resolvedSecurityGroupRuleSpec {
	return []resolvedSecurityGroupRuleSpec{
		{
			Description:   "In-cluster Ingress",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  0,
			PortRangeMax:  0,
			Protocol:      "",
			RemoteGroupID: remoteGroupIDSelf,
		},
		{
			Description:   "In-cluster Ingress",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  0,
			PortRangeMax:  0,
			Protocol:      "",
			RemoteGroupID: secWorkerGroupID,
		},
	}
}

// Permit all ingress from the cluster security groups.
func getSGWorkerAllowAll(remoteGroupIDSelf, secControlPlaneGroupID string) []resolvedSecurityGroupRuleSpec {
	return []resolvedSecurityGroupRuleSpec{
		{
			Description:   "In-cluster Ingress",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  0,
			PortRangeMax:  0,
			Protocol:      "",
			RemoteGroupID: remoteGroupIDSelf,
		},
		{
			Description:   "In-cluster Ingress",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  0,
			PortRangeMax:  0,
			Protocol:      "",
			RemoteGroupID: secControlPlaneGroupID,
		},
	}
}

// Permit ports that defined in openStackCluster.Spec.APIServerLoadBalancer.AdditionalPorts.
func getSGControlPlaneAdditionalPorts(ports []int) []resolvedSecurityGroupRuleSpec {
	controlPlaneRules := []resolvedSecurityGroupRuleSpec{}

	r := []resolvedSecurityGroupRuleSpec{
		{
			Description: "Additional ports",
			Direction:   "ingress",
			EtherType:   "IPv4",
			Protocol:    "tcp",
		},
	}
	for i, p := range ports {
		r[i].PortRangeMin = p
		r[i].PortRangeMax = p
		controlPlaneRules = append(controlPlaneRules, r...)
	}
	return controlPlaneRules
}

func getSGControlPlaneGeneral(remoteGroupIDSelf, secWorkerGroupID string) []resolvedSecurityGroupRuleSpec {
	controlPlaneRules := []resolvedSecurityGroupRuleSpec{}
	controlPlaneRules = append(controlPlaneRules, getSGControlPlaneCommon(remoteGroupIDSelf, secWorkerGroupID)...)
	return controlPlaneRules
}

func getSGWorkerGeneral(remoteGroupIDSelf, secControlPlaneGroupID string) []resolvedSecurityGroupRuleSpec {
	workerRules := []resolvedSecurityGroupRuleSpec{}
	workerRules = append(workerRules, getSGWorkerCommon(remoteGroupIDSelf, secControlPlaneGroupID)...)
	return workerRules
}
