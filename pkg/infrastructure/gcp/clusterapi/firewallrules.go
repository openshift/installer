package clusterapi

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"

	configv1 "github.com/openshift/api/config/v1"
	gcpconfig "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	"github.com/openshift/installer/pkg/types"
)

const (
	// gcpFirewallPermission is the role/permission to create or skip the creation of
	// firewall rules for GCP during a xpn installation.
	gcpFirewallPermission = "compute.firewalls.create"
)

func getEtcdPorts() []*compute.FirewallAllowed {
	return []*compute.FirewallAllowed{
		{
			IPProtocol: "tcp",
			Ports: []string{
				"2379-2380",
			},
		},
	}
}

func getHealthChecksPorts() []*compute.FirewallAllowed {
	return []*compute.FirewallAllowed{
		{
			IPProtocol: "tcp",
			Ports: []string{
				"6080",
				"6443",
				"22624",
			},
		},
	}
}

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

func getBootstrapSSHPorts() []*compute.FirewallAllowed {
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

// addFirewallRule creates the firewall rule and adds it the compute's firewalls.
func addFirewallRule(ctx context.Context, svc *compute.Service, name, network, projectID string, ports []*compute.FirewallAllowed, srcTags, targetTags, srcRanges []string) error {
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

	op, err := svc.Firewalls.Insert(projectID, firewallRule).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to create %s firewall rule: %w", name, err)
	}

	if err := WaitForOperationGlobal(ctx, svc, projectID, op); err != nil {
		return fmt.Errorf("failed to wait for inserting %s firewall rule: %w", name, err)
	}

	return nil
}

// deleteFirewallRule deletes the firewall rule identified by name.
func deleteFirewallRule(ctx context.Context, svc *compute.Service, name, projectID string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*3)
	defer cancel()

	op, err := svc.Firewalls.Delete(projectID, name).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to delete %s firewall rule: %w", name, err)
	}

	if err := WaitForOperationGlobal(ctx, svc, projectID, op); err != nil {
		return fmt.Errorf("failed to wait for delete %s firewall rule: %w", name, err)
	}

	return nil
}

func hasFirewallPermission(ctx context.Context, projectID string, endpoints []configv1.GCPServiceEndpoint) (bool, error) {
	client, err := gcpconfig.NewClient(ctx, endpoints)
	if err != nil {
		return false, fmt.Errorf("failed to create client during firewall permission check: %w", err)
	}

	permissions, err := client.GetProjectPermissions(ctx, projectID, []string{
		gcpFirewallPermission,
	})
	if err != nil {
		return false, fmt.Errorf("failed to find project permissions during firewall permission check: %w", err)
	}

	hasPermission := permissions.Has(gcpFirewallPermission)
	if !hasPermission {
		logrus.Warnf("failed to find permission %s, skipping firewall rule creation", gcpFirewallPermission)
	}

	return hasPermission, nil
}

// createFirewallRules creates the rules needed between the worker and master nodes.
func createFirewallRules(ctx context.Context, in clusterapi.InfraReadyInput, network string) error {
	projectID := in.InstallConfig.Config.Platform.GCP.ProjectID
	if in.InstallConfig.Config.GCP.NetworkProjectID != "" {
		projectID = in.InstallConfig.Config.GCP.NetworkProjectID

		createFwRules, err := hasFirewallPermission(ctx, projectID, []configv1.GCPServiceEndpoint{})
		if err != nil {
			return fmt.Errorf("failed to create cluster firewall rules: %w", err)
		}
		if !createFwRules {
			return nil
		}
	}

	svc, err := gcpconfig.GetComputeService(ctx, []configv1.GCPServiceEndpoint{})
	if err != nil {
		return fmt.Errorf("failed to get copmute service for firewall rule creation: %w", err)
	}

	workerTag := fmt.Sprintf("%s-worker", in.InfraID)
	masterTag := fmt.Sprintf("%s-control-plane", in.InfraID)

	// control-plane rules are needed for worker<->master communication for worker provisioning
	firewallName := fmt.Sprintf("%s-control-plane", in.InfraID)
	srcTags := []string{workerTag, masterTag}
	targetTags := []string{masterTag}
	srcRanges := []string{}
	if err := addFirewallRule(ctx, svc, firewallName, network, projectID, getControlPlanePorts(), srcTags, targetTags, srcRanges); err != nil {
		return err
	}

	// etcd are needed for master communication for etcd nodes
	firewallName = fmt.Sprintf("%s-etcd", in.InfraID)
	srcTags = []string{masterTag}
	targetTags = []string{masterTag}
	srcRanges = []string{}
	if err := addFirewallRule(ctx, svc, firewallName, network, projectID, getEtcdPorts(), srcTags, targetTags, srcRanges); err != nil {
		return err
	}

	// Add a single firewall rule to allow the Google Cloud Engine health checks to access all of the services.
	// This rule enables the ingress load balancers to determine the health status of their instances.
	firewallName = fmt.Sprintf("%s-health-checks", in.InfraID)
	srcTags = []string{}
	targetTags = []string{masterTag}
	srcRanges = []string{"35.191.0.0/16", "130.211.0.0/22"}
	if in.InstallConfig.Config.Publish == types.ExternalPublishingStrategy {
		// public installs require additional google ip addresses for health checks
		srcRanges = append(srcRanges, []string{"209.85.152.0/22", "209.85.204.0/22"}...)
	}
	if err := addFirewallRule(ctx, svc, firewallName, network, projectID, getHealthChecksPorts(), srcTags, targetTags, srcRanges); err != nil {
		return err
	}

	// internal-cluster rules are needed for worker<->master communication for k8s nodes
	firewallName = fmt.Sprintf("%s-internal-cluster", in.InfraID)
	srcTags = []string{workerTag, masterTag}
	targetTags = []string{workerTag, masterTag}
	srcRanges = []string{}
	if err := addFirewallRule(ctx, svc, firewallName, network, projectID, getInternalClusterPorts(), srcTags, targetTags, srcRanges); err != nil {
		return err
	}

	machineCIDR := in.InstallConfig.Config.Networking.MachineNetwork[0].CIDR.String()
	// api rules are needed to access the kube-apiserver on master nodes
	firewallName = fmt.Sprintf("%s-api", in.InfraID)
	srcTags = []string{}
	targetTags = []string{masterTag}
	srcRanges = []string{}
	if !in.InstallConfig.Config.PublicAPI() {
		// For Internal, limit the source to the machineCIDR
		srcRanges = append(srcRanges, machineCIDR)
	}
	if err := addFirewallRule(ctx, svc, firewallName, network, projectID, getAPIPorts(), srcTags, targetTags, srcRanges); err != nil {
		return err
	}

	// internal-network rules are used to access ssh and icmp over the machine network
	firewallName = fmt.Sprintf("%s-internal-network", in.InfraID)
	srcTags = []string{}
	targetTags = []string{workerTag, masterTag}
	srcRanges = []string{machineCIDR}
	err = addFirewallRule(ctx, svc, firewallName, network, projectID, getInternalNetworkPorts(), srcTags, targetTags, srcRanges)

	return err
}

// createBootstrapFirewallRules creates the rules needed for the bootstrap node.
func createBootstrapFirewallRules(ctx context.Context, in clusterapi.InfraReadyInput, network string) error {
	projectID := in.InstallConfig.Config.Platform.GCP.ProjectID
	if in.InstallConfig.Config.Platform.GCP.NetworkProjectID != "" {
		projectID = in.InstallConfig.Config.Platform.GCP.NetworkProjectID

		createFwRules, err := hasFirewallPermission(ctx, projectID, []configv1.GCPServiceEndpoint{})
		if err != nil {
			return fmt.Errorf("failed to create bootstrap firewall rules: %w", err)
		}
		if !createFwRules {
			return nil
		}
	}

	svc, err := gcpconfig.GetComputeService(ctx, []configv1.GCPServiceEndpoint{})
	if err != nil {
		return fmt.Errorf("failed to get compute service for bootstrap firewall rule creation: %w", err)
	}

	firewallName := fmt.Sprintf("%s-bootstrap-in-ssh", in.InfraID)
	srcTags := []string{}
	bootstrapTag := fmt.Sprintf("%s-control-plane", in.InfraID)
	targetTags := []string{bootstrapTag}
	var srcRanges []string
	if in.InstallConfig.Config.Publish == types.ExternalPublishingStrategy {
		srcRanges = []string{"0.0.0.0/0"}
	} else {
		machineCIDR := in.InstallConfig.Config.Networking.MachineNetwork[0].CIDR.String()
		srcRanges = []string{machineCIDR}
	}
	return addFirewallRule(ctx, svc, firewallName, network, projectID, getBootstrapSSHPorts(), srcTags, targetTags, srcRanges)
}

// removeBootstrapFirewallRules removes the rules created for the bootstrap node.
func removeBootstrapFirewallRules(ctx context.Context, infraID, projectID string, endpoints []configv1.GCPServiceEndpoint) error {
	svc, err := gcpconfig.GetComputeService(ctx, endpoints)
	if err != nil {
		return err
	}

	firewallName := fmt.Sprintf("%s-bootstrap-in-ssh", infraID)
	return deleteFirewallRule(ctx, svc, firewallName, projectID)
}

// removeCAPGFirewallRules removes the overly permissive firewall rules created by cluster-api-provider-gcp.
func removeCAPGFirewallRules(ctx context.Context, infraID, projectID string, endpoints []configv1.GCPServiceEndpoint) error {
	svc, err := gcpconfig.GetComputeService(ctx, endpoints)
	if err != nil {
		return err
	}

	firewallName := fmt.Sprintf("allow-%s-cluster", infraID)
	if err := deleteFirewallRule(ctx, svc, firewallName, projectID); err != nil {
		return err
	}

	firewallName = fmt.Sprintf("allow-%s-healthchecks", infraID)
	return deleteFirewallRule(ctx, svc, firewallName, projectID)
}
