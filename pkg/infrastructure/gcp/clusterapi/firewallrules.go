package clusterapi

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/compute/v1"

	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	"github.com/openshift/installer/pkg/types"
)

func getControlPlanePorts() []*compute.FirewallAllowed {
	return []*compute.FirewallAllowed{
		{
			IPProtocol: "tcp",
			Ports: []string{
				"22623", // Ignition
			},
		},
		{
			IPProtocol: "tcp",
			Ports: []string{
				"10257", // Kube manager
			},
		},
		{
			IPProtocol: "tcp",
			Ports: []string{
				"10259", // Kube scheduler
			},
		},
	}
}

func getInternalClusterPorts() []*compute.FirewallAllowed {
	return []*compute.FirewallAllowed{
		{
			IPProtocol: "tcp",
			Ports: []string{
				"30000-32767", // k8s NodePorts
			},
		},
		{
			IPProtocol: "udp",
			Ports: []string{
				"30000-32767", // k8s NodePorts
			},
		},
		{
			IPProtocol: "tcp",
			Ports: []string{
				"9000-9999", // host-level services
			},
		},
		{
			IPProtocol: "udp",
			Ports: []string{
				"9000-9999", // host-level services
			},
		},
		{
			IPProtocol: "udp",
			Ports: []string{
				"4789", "6081", // VXLAN and GENEVE
			},
		},
		{
			IPProtocol: "udp",
			Ports: []string{
				"500", "4500", // IKE and IKE(NAT-T)
			},
		},
		{
			IPProtocol: "tcp",
			Ports: []string{
				"10250", // kubelet secure
			},
		},
		{
			IPProtocol: "esp",
		},
	}
}

func getAPIPorts() []*compute.FirewallAllowed {
	return []*compute.FirewallAllowed{
		{
			IPProtocol: "tcp",
			Ports: []string{
				"6443", // kube-apiserver
			},
		},
	}
}

func getInternalNetworkPorts() []*compute.FirewallAllowed {
	return []*compute.FirewallAllowed{
		{
			IPProtocol: "tcp",
			Ports: []string{
				"22", // SSH
			},
		},
		{
			IPProtocol: "icmp",
		},
	}
}

func getIngressPorts() []*compute.FirewallAllowed {
	return []*compute.FirewallAllowed{
		{
			IPProtocol: "tcp",
			Ports: []string{
				"6080", // ingress, *.apps
			},
		},
	}
}

// addFirewallRule creates the firewall rule and adds it the compute's firewalls.
func addFirewallRule(ctx context.Context, name, network, projectID string, ports []*compute.FirewallAllowed, srcTags, targetTags, srcRanges []string) error {
	service, err := NewComputeService()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Minute*3)
	defer cancel()

	firewallRule := &compute.Firewall{
		Name:        name,
		Description: resourceDescription,
		Direction:   "INGRESS",
		Network:     network,
		Allowed:     ports,
		SourceTags:  srcTags,
		TargetTags:  targetTags,
	}
	if len(srcTags) > 0 {
		firewallRule.SourceTags = srcTags
	}
	if len(srcRanges) > 0 {
		firewallRule.SourceRanges = srcRanges
	}

	op, err := service.Firewalls.Insert(projectID, firewallRule).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to create %s firewall rule: %w", name, err)
	}

	if err := WaitForOperationGlobal(ctx, projectID, op); err != nil {
		return fmt.Errorf("failed to wait for inserting %s firewall rule: %w", name, err)
	}

	return nil
}

// createFirewallRules creates the rules needed between tthe worker and master nodes.
func createFirewallRules(ctx context.Context, in clusterapi.InfraReadyInput, network string) error {
	projectID := in.InstallConfig.Config.Platform.GCP.ProjectID
	workerTag := fmt.Sprintf("%s-worker", in.InfraID)
	masterTag := fmt.Sprintf("%s-control-plane", in.InfraID)

	// control-plane rules are needed for worker<->master communication for worker provisioning
	firewallName := fmt.Sprintf("%s-control-plane", in.InfraID)
	srcTags := []string{workerTag, masterTag}
	targetTags := []string{masterTag}
	srcRanges := []string{}
	if err := addFirewallRule(ctx, firewallName, network, projectID, getControlPlanePorts(), srcTags, targetTags, srcRanges); err != nil {
		return err
	}

	// internal-cluster rules are needed for worker<->master communication for k8s nodes
	firewallName = fmt.Sprintf("%s-internal-cluster", in.InfraID)
	srcTags = []string{workerTag, masterTag}
	targetTags = []string{workerTag, masterTag}
	srcRanges = []string{}
	if err := addFirewallRule(ctx, firewallName, network, projectID, getInternalClusterPorts(), srcTags, targetTags, srcRanges); err != nil {
		return err
	}

	// api rules are needed to access the kube-apiserver on master nodes
	firewallName = fmt.Sprintf("%s-api", in.InfraID)
	srcTags = []string{}
	targetTags = []string{masterTag}
	srcRanges = []string{}
	if err := addFirewallRule(ctx, firewallName, network, projectID, getAPIPorts(), srcTags, targetTags, srcRanges); err != nil {
		return err
	}

	// internal-network rules are used to access ssh and icmp over the machine network
	firewallName = fmt.Sprintf("%s-internal-network", in.InfraID)
	srcTags = []string{}
	targetTags = []string{workerTag, masterTag}
	machineCIDR := in.InstallConfig.Config.Networking.MachineNetwork[0].CIDR.String()
	srcRanges = []string{machineCIDR}
	if err := addFirewallRule(ctx, firewallName, network, projectID, getInternalNetworkPorts(), srcTags, targetTags, srcRanges); err != nil {
		return err
	}

	// health-check rules are used to allow GCP health probes to reach the control plane.
	// https://cloud.google.com/load-balancing/docs/health-checks#fw-rule
	firewallName = fmt.Sprintf("%s-health-checks", in.InfraID)
	srcTags = []string{}
	targetTags = []string{masterTag}
	srcRanges = []string{"35.191.0.0/16", "130.211.0.0/22"}
	if in.InstallConfig.Config.Publish == types.ExternalPublishingStrategy {
		srcRanges = append(srcRanges, "209.85.152.0/22", "209.85.204.0/22")
	}
	err := addFirewallRule(ctx, firewallName, network, projectID, getIngressPorts(), srcTags, targetTags, srcRanges)

	return err
}
