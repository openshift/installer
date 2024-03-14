package aws

import capa "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"

func getDefaultNetworkCNIIngressRules() capa.CNIIngressRules {

	return capa.CNIIngressRules{
		{
			Description: "ICMP",
			Protocol:    capa.SecurityGroupProtocolICMP,
			FromPort:    -1,
			ToPort:      -1,
		},
		{
			Description: "Port 22 (TCP)",
			Protocol:    capa.SecurityGroupProtocolTCP,
			FromPort:    22,
			ToPort:      22,
		},
		{
			Description: "Port 4789 (UDP) for VXLAN",
			Protocol:    capa.SecurityGroupProtocolUDP,
			FromPort:    4789,
			ToPort:      4789,
		},
		{
			Description: "Port 6081 (UDP) for geneve",
			Protocol:    capa.SecurityGroupProtocolUDP,
			FromPort:    6081,
			ToPort:      6081,
		},
		{
			Description: "Port 500 (UDP) for IKE",
			Protocol:    capa.SecurityGroupProtocolUDP,
			FromPort:    500,
			ToPort:      500,
		},
		{
			Description: "Port 4500 (UDP) for IKE NAT",
			Protocol:    capa.SecurityGroupProtocolUDP,
			FromPort:    4500,
			ToPort:      4500,
		},
		{
			Description: "ESP",
			Protocol:    capa.SecurityGroupProtocolESP,
			FromPort:    -1,
			ToPort:      -1,
		},
		{
			Description: "Port 6441-6442 (TCP) for ovndb",
			Protocol:    capa.SecurityGroupProtocolTCP,
			FromPort:    6441,
			ToPort:      6442,
		},
		{
			Description: "Port 9000-9999 for node ports (TCP)",
			Protocol:    capa.SecurityGroupProtocolTCP,
			FromPort:    9000,
			ToPort:      9999,
		},
		{
			Description: "Port 9000-9999 for node ports (UDP)",
			Protocol:    capa.SecurityGroupProtocolUDP,
			FromPort:    9000,
			ToPort:      9999,
		},
		{
			Description: "Service node ports (TCP)",
			Protocol:    capa.SecurityGroupProtocolTCP,
			FromPort:    30000,
			ToPort:      32767,
		},
		{
			Description: "Service node ports (UDP)",
			Protocol:    capa.SecurityGroupProtocolUDP,
			FromPort:    30000,
			ToPort:      32767,
		},
		{
			Description: "Machine Config Server (MCS)",
			Protocol:    capa.SecurityGroupProtocolTCP,
			FromPort:    22623,
			ToPort:      22623,
		},
	}
}

func getDefaultNetworkAdditionalControlPlaneIngressRules() []capa.IngressRule {

	return []capa.IngressRule{
		{
			Description: "MCS traffic from cluster network",
			Protocol:    capa.SecurityGroupProtocolTCP,
			FromPort:    22623,
			ToPort:      22623,
			//SourceSecurityGroupRoles: []capa.SecurityGroupRole{"node", "controlplane"},
			CidrBlocks: []string{"10.0.0.0/16"}, //TODO(padillon): figure out security group rules
		},
		{
			Description:              "controller-manager",
			Protocol:                 capa.SecurityGroupProtocolTCP,
			FromPort:                 10257,
			ToPort:                   10257,
			SourceSecurityGroupRoles: []capa.SecurityGroupRole{"controlplane", "node"},
		},
		{
			Description:              "kube-scheduler",
			Protocol:                 capa.SecurityGroupProtocolTCP,
			FromPort:                 10259,
			ToPort:                   10259,
			SourceSecurityGroupRoles: []capa.SecurityGroupRole{"controlplane", "node"},
		},
		{
			Description: "SSH everyone",
			Protocol:    capa.SecurityGroupProtocolTCP,
			FromPort:    22,
			ToPort:      22,
			CidrBlocks:  []string{"0.0.0.0/0"},
		},
		{
			Description: "public api", //TESTING
			Protocol:    capa.SecurityGroupProtocolTCP,
			FromPort:    6443,
			ToPort:      6443,
			CidrBlocks:  []string{"0.0.0.0/0"},
		},
	}
}
