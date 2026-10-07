package scope

import (
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/api/compute/v1"
	infrav1 "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	firewallutil "sigs.k8s.io/cluster-api-provider-gcp/util/firewall"
)

// createFirewallRules
func createFirewallRules(clusterName, networkLink string, policy infrav1.RulesManagementPolicy, userSpecifiedRules []infrav1.FirewallRule) ([]*compute.Firewall, error) {
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
						IPProtocol: protocolTCP,
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

	// Rules are normally named by the defaulting webhook, but it leaves the spec of a
	// ClusterClass owned cluster alone (see SkipRuleNameDefaulting), so those rules
	// arrive here unnamed and are named below instead. The name is derived from the
	// rule, so it matches the one the webhook would have written and is the same on
	// every reconcile.
	takenNames := firewallutil.TakenRuleNames(userSpecifiedRules)

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
		name := rule.Name
		if name == "" {
			generated, err := firewallutil.GenerateRuleName(clusterName, rule, takenNames)
			if err != nil {
				return nil, err
			}
			name = generated
			takenNames.Insert(name)
		} else {
			name = firewallutil.QualifyRuleName(clusterName, name)
		}

		description := rule.Description
		if description == "" {
			description = "Created by Cluster API GCP Provider"
		}

		firewallRules = append(firewallRules, &compute.Firewall{
			Name:              name,
			Description:       description,
			Network:           networkLink,
			Allowed:           allowed,
			Denied:            denied,
			Direction:         direction,
			Priority:          int64(rule.Priority),
			Disabled:          false,
			SourceRanges:      rule.SourceRanges,
			DestinationRanges: rule.DestinationRanges,
			TargetTags:        rule.TargetTags,
			SourceTags:        rule.SourceTags,
		})
	}

	return firewallRules, nil
}
