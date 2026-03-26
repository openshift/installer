package scope

import (
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/api/compute/v1"

	infrav1 "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
)

// createFirewallRules
func createFirewallRules(clusterName, networkLink string, policy infrav1.RulesManagementPolicy, userSpecifiedRules []infrav1.FirewallRule) []*compute.Firewall {
	firewallRules := []*compute.Firewall{}

	// Only when the user explicitly states that it is unmanaged, the rules should be skipped.
	// When the policy is Managed or missing/empty the rules should be included.
	if policy != infrav1.RulesManagementUnmanaged {
		firewallRules = append(firewallRules, []*compute.Firewall{
			{
				Name:    fmt.Sprintf("allow-%s-healthchecks", clusterName),
				Network: networkLink,
				Allowed: []*compute.FirewallAllowed{
					{
						IPProtocol: "TCP",
						Ports: []string{
							strconv.FormatInt(6443, 10),
						},
					},
				},
				Direction: "INGRESS",
				SourceRanges: []string{
					"35.191.0.0/16",
					"130.211.0.0/22",
				},
				TargetTags: []string{
					clusterName + "-control-plane",
				},
			},
			{
				Name:    fmt.Sprintf("allow-%s-cluster", clusterName),
				Network: networkLink,
				Allowed: []*compute.FirewallAllowed{
					{
						IPProtocol: "all",
					},
				},
				Direction: "INGRESS",
				SourceTags: []string{
					clusterName + "-control-plane",
					clusterName + "-node",
				},
				TargetTags: []string{
					clusterName + "-control-plane",
					clusterName + "-node",
				},
			},
		}...)
	}

	// Add user defined firewall rules.
	for _, rule := range userSpecifiedRules {
		allowed := []*compute.FirewallAllowed{}
		for _, a := range rule.Allowed {
			allowed = append(allowed, &compute.FirewallAllowed{
				IPProtocol: strings.ToLower(string(a.IPProtocol)),
				Ports:      a.Ports,
			})
		}

		denied := []*compute.FirewallDenied{}
		for _, d := range rule.Denied {
			denied = append(denied, &compute.FirewallDenied{
				IPProtocol: strings.ToLower(string(d.IPProtocol)),
				Ports:      d.Ports,
			})
		}

		direction := strings.ToUpper(string(rule.Direction))
		name := fmt.Sprintf("%s-%s", clusterName, strings.ToLower(direction))
		if rule.Name != "" {
			name = rule.Name
			if !strings.HasPrefix(name, clusterName) {
				name = fmt.Sprintf("%s-%s", clusterName, name)
			}
		}
		name = name[:min(len(name), 63)]
		name = strings.TrimSuffix(name, "-")

		description := rule.Description
		if description == "" {
			description = "Created by Cluster API GCP Provider"
		}

		firewallRules = append(firewallRules, &compute.Firewall{
			Name:         name,
			Description:  description,
			Network:      networkLink,
			Allowed:      allowed,
			Denied:       denied,
			Direction:    direction,
			Priority:     int64(rule.Priority),
			Disabled:     false,
			SourceRanges: rule.SourceRanges,
			TargetTags:   rule.TargetTags,
			SourceTags:   rule.SourceTags,
		})
	}

	return firewallRules
}
