package gcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	capg "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"

	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/gcp"
)

// validTestInstallConfig creates a valid base install config for testing.
func validTestInstallConfig() *installconfig.InstallConfig {
	return installconfig.MakeAsset(&types.InstallConfig{
		BaseDomain: "example.com",
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-cluster",
		},
		Networking: &types.Networking{
			MachineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR("10.0.0.0/16")},
			},
		},
		Platform: types.Platform{
			GCP: &gcp.Platform{
				ProjectID: "test-project",
				Region:    "us-east1",
			},
		},
		Publish: types.ExternalPublishingStrategy,
	})
}

func TestCreateFirewallRulesForCAPG(t *testing.T) {
	tests := []struct {
		name               string
		installConfig      *installconfig.InstallConfig
		expectedRulesCount int
		validateRules      func(*testing.T, []capg.FirewallRule, string, string)
	}{
		{
			name:               "External publishing strategy",
			installConfig:      validTestInstallConfig(),
			expectedRulesCount: 6,
			validateRules: func(t *testing.T, rules []capg.FirewallRule, infraID string, machineCIDR string) {
				ruleMap := make(map[string]capg.FirewallRule)
				for _, rule := range rules {
					ruleMap[rule.Name] = rule
				}

				masterTag := infraID + "-control-plane"
				workerTag := infraID + "-worker"

				// Verify control-plane rule
				cpRule, found := ruleMap[infraID+"-control-plane"]
				if !assert.True(t, found, "control-plane rule should exist") {
					return
				}
				assert.Equal(t, []string{masterTag}, cpRule.TargetTags)
				assert.Equal(t, []string{workerTag, masterTag}, cpRule.SourceTags)
				assert.Empty(t, cpRule.SourceRanges)
				assert.Equal(t, capg.FirewallRuleDirectionIngress, cpRule.Direction)
				assert.Equal(t, resourceDescription, cpRule.Description)
				assert.Equal(t, 1000, cpRule.Priority)
				assert.Len(t, cpRule.Allowed, 3)
				assert.Contains(t, cpRule.Allowed, capg.FirewallDescriptor{
					IPProtocol: capg.FirewallProtocolTCP,
					Ports:      []string{"22623"},
				})
				assert.Contains(t, cpRule.Allowed, capg.FirewallDescriptor{
					IPProtocol: capg.FirewallProtocolTCP,
					Ports:      []string{"10257"},
				})
				assert.Contains(t, cpRule.Allowed, capg.FirewallDescriptor{
					IPProtocol: capg.FirewallProtocolTCP,
					Ports:      []string{"10259"},
				})

				// Verify etcd rule
				etcdRule, found := ruleMap[infraID+"-etcd"]
				if !assert.True(t, found, "etcd rule should exist") {
					return
				}
				assert.Equal(t, []string{masterTag}, etcdRule.TargetTags)
				assert.Equal(t, []string{masterTag}, etcdRule.SourceTags)
				assert.Empty(t, etcdRule.SourceRanges)
				assert.Len(t, etcdRule.Allowed, 1)
				assert.Contains(t, etcdRule.Allowed, capg.FirewallDescriptor{
					IPProtocol: capg.FirewallProtocolTCP,
					Ports:      []string{"2379-2380"},
				})

				// Verify health-checks rule
				healthRule, found := ruleMap[infraID+"-health-checks"]
				if !assert.True(t, found, "health-checks rule should exist") {
					return
				}
				assert.Equal(t, []string{masterTag}, healthRule.TargetTags)
				assert.Empty(t, healthRule.SourceTags)
				// For external publishing, should have base health check ranges only
				assert.Equal(t, []string{"35.191.0.0/16", "130.211.0.0/22"}, healthRule.SourceRanges)
				assert.Len(t, healthRule.Allowed, 1)
				assert.Contains(t, healthRule.Allowed, capg.FirewallDescriptor{
					IPProtocol: capg.FirewallProtocolTCP,
					Ports:      []string{"6080", "6443", "22624"},
				})

				// Verify internal-cluster rule
				internalRule, found := ruleMap[infraID+"-internal-cluster"]
				if !assert.True(t, found, "internal-cluster rule should exist") {
					return
				}
				assert.Equal(t, []string{workerTag, masterTag}, internalRule.TargetTags)
				assert.Equal(t, []string{workerTag, masterTag}, internalRule.SourceTags)
				assert.Empty(t, internalRule.SourceRanges)
				assert.Len(t, internalRule.Allowed, 8)
				// Verify key protocols and ports
				assert.Contains(t, internalRule.Allowed, capg.FirewallDescriptor{
					IPProtocol: capg.FirewallProtocolTCP,
					Ports:      []string{"30000-32767"},
				})
				assert.Contains(t, internalRule.Allowed, capg.FirewallDescriptor{
					IPProtocol: capg.FirewallProtocolUDP,
					Ports:      []string{"30000-32767"},
				})
				assert.Contains(t, internalRule.Allowed, capg.FirewallDescriptor{
					IPProtocol: capg.FirewallProtocolTCP,
					Ports:      []string{"9000-9999"},
				})
				assert.Contains(t, internalRule.Allowed, capg.FirewallDescriptor{
					IPProtocol: capg.FirewallProtocolUDP,
					Ports:      []string{"9000-9999"},
				})
				assert.Contains(t, internalRule.Allowed, capg.FirewallDescriptor{
					IPProtocol: capg.FirewallProtocolUDP,
					Ports:      []string{"4789", "6081"},
				})
				assert.Contains(t, internalRule.Allowed, capg.FirewallDescriptor{
					IPProtocol: capg.FirewallProtocolUDP,
					Ports:      []string{"500", "4500"},
				})
				assert.Contains(t, internalRule.Allowed, capg.FirewallDescriptor{
					IPProtocol: capg.FirewallProtocolTCP,
					Ports:      []string{"10250"},
				})
				assert.Contains(t, internalRule.Allowed, capg.FirewallDescriptor{
					IPProtocol: capg.FirewallProtocolESP,
				})

				// Verify api rule - for external publishing, source ranges should be empty
				apiRule, found := ruleMap[infraID+"-api"]
				if !assert.True(t, found, "api rule should exist") {
					return
				}
				assert.Equal(t, []string{masterTag}, apiRule.TargetTags)
				assert.Empty(t, apiRule.SourceRanges, "API rule should allow access from anywhere for external publishing")
				assert.Len(t, apiRule.Allowed, 1)
				assert.Contains(t, apiRule.Allowed, capg.FirewallDescriptor{
					IPProtocol: capg.FirewallProtocolTCP,
					Ports:      []string{"6443"},
				})

				// Verify internal-network rule
				networkRule, found := ruleMap[infraID+"-internal-network"]
				if !assert.True(t, found, "internal-network rule should exist") {
					return
				}
				assert.Equal(t, []string{workerTag, masterTag}, networkRule.TargetTags)
				assert.Empty(t, networkRule.SourceTags)
				assert.Equal(t, []string{machineCIDR}, networkRule.SourceRanges)
				assert.Len(t, networkRule.Allowed, 2)
				assert.Contains(t, networkRule.Allowed, capg.FirewallDescriptor{
					IPProtocol: capg.FirewallProtocolTCP,
					Ports:      []string{"22"},
				})
				assert.Contains(t, networkRule.Allowed, capg.FirewallDescriptor{
					IPProtocol: capg.FirewallProtocolICMP,
				})
			},
		},
		{
			name: "Internal publishing strategy",
			installConfig: func() *installconfig.InstallConfig {
				ic := validTestInstallConfig()
				ic.Config.Publish = types.InternalPublishingStrategy
				return ic
			}(),
			expectedRulesCount: 6,
			validateRules: func(t *testing.T, rules []capg.FirewallRule, infraID string, machineCIDR string) {
				ruleMap := make(map[string]capg.FirewallRule)
				for _, rule := range rules {
					ruleMap[rule.Name] = rule
				}

				// For internal publishing, health checks should include additional ranges
				healthRule, found := ruleMap[infraID+"-health-checks"]
				if !assert.True(t, found, "health-checks rule should exist") {
					return
				}
				assert.Equal(t, []string{
					"35.191.0.0/16",
					"130.211.0.0/22",
					"209.85.152.0/22",
					"209.85.204.0/22",
				}, healthRule.SourceRanges, "Internal publishing should include additional health check ranges")

				// For internal publishing, API should be limited to machine CIDR
				apiRule, found := ruleMap[infraID+"-api"]
				if !assert.True(t, found, "api rule should exist") {
					return
				}
				assert.Equal(t, []string{machineCIDR}, apiRule.SourceRanges, "API rule should be restricted to machine CIDR for internal publishing")
			},
		},
		{
			name: "Custom infraID",
			installConfig: func() *installconfig.InstallConfig {
				ic := validTestInstallConfig()
				return ic
			}(),
			expectedRulesCount: 6,
			validateRules: func(t *testing.T, rules []capg.FirewallRule, infraID string, machineCIDR string) {
				// Verify all rule names contain the infraID
				for _, rule := range rules {
					assert.Contains(t, rule.Name, infraID, "Rule name should contain infraID")
				}

				// Verify all tags contain the infraID
				for _, rule := range rules {
					for _, tag := range rule.TargetTags {
						assert.Contains(t, tag, infraID, "Target tag should contain infraID")
					}
					for _, tag := range rule.SourceTags {
						assert.Contains(t, tag, infraID, "Source tag should contain infraID")
					}
				}
			},
		},
		{
			name: "Different machine CIDR",
			installConfig: func() *installconfig.InstallConfig {
				ic := validTestInstallConfig()
				ic.Config.Networking.MachineNetwork = []types.MachineNetworkEntry{
					{CIDR: *ipnet.MustParseCIDR("192.168.0.0/16")},
				}
				return ic
			}(),
			expectedRulesCount: 6,
			validateRules: func(t *testing.T, rules []capg.FirewallRule, infraID string, machineCIDR string) {
				ruleMap := make(map[string]capg.FirewallRule)
				for _, rule := range rules {
					ruleMap[rule.Name] = rule
				}

				// Verify internal-network rule uses the custom machine CIDR
				networkRule, found := ruleMap[infraID+"-internal-network"]
				if !assert.True(t, found) {
					return
				}
				assert.Equal(t, []string{"192.168.0.0/16"}, networkRule.SourceRanges)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clusterID := &installconfig.ClusterID{
				InfraID: "test-infra",
			}

			rules := createFirewallRulesForCAPG(tt.installConfig, clusterID)

			assert.Len(t, rules, tt.expectedRulesCount, "unexpected number of firewall rules")

			// All rules should have the resource description and ingress direction
			for _, rule := range rules {
				assert.Equal(t, resourceDescription, rule.Description)
				assert.Equal(t, capg.FirewallRuleDirectionIngress, rule.Direction)
			}

			machineCIDR := tt.installConfig.Config.Networking.MachineNetwork[0].CIDR.String()
			if tt.validateRules != nil {
				tt.validateRules(t, rules, clusterID.InfraID, machineCIDR)
			}
		})
	}
}

func TestCreateFirewallRulesForCAPG_AllRulesPresent(t *testing.T) {
	// Test that all expected rules are created
	installConfig := validTestInstallConfig()
	clusterID := &installconfig.ClusterID{
		InfraID: "test-cluster-infra",
	}

	rules := createFirewallRulesForCAPG(installConfig, clusterID)

	expectedRuleNames := []string{
		"test-cluster-infra-control-plane",
		"test-cluster-infra-etcd",
		"test-cluster-infra-health-checks",
		"test-cluster-infra-internal-cluster",
		"test-cluster-infra-api",
		"test-cluster-infra-internal-network",
	}

	actualRuleNames := make([]string, len(rules))
	for i, rule := range rules {
		actualRuleNames[i] = rule.Name
	}

	assert.ElementsMatch(t, expectedRuleNames, actualRuleNames, "All expected firewall rules should be present")
}

func TestCreateBootstrapFirewallRuleForCAPG(t *testing.T) {
	tests := []struct {
		name                 string
		installConfig        *installconfig.InstallConfig
		expectedRuleName     string
		expectedSourceRanges []string
		expectedTargetTags   []string
	}{
		{
			name:                 "External publishing strategy",
			installConfig:        validTestInstallConfig(),
			expectedRuleName:     "test-infra-bootstrap-in-ssh",
			expectedSourceRanges: []string{"0.0.0.0/0"},
			expectedTargetTags:   []string{"test-infra-control-plane"},
		},
		{
			name: "Internal publishing strategy",
			installConfig: func() *installconfig.InstallConfig {
				ic := validTestInstallConfig()
				ic.Config.Publish = types.InternalPublishingStrategy
				return ic
			}(),
			expectedRuleName:     "test-infra-bootstrap-in-ssh",
			expectedSourceRanges: []string{"10.0.0.0/16"},
			expectedTargetTags:   []string{"test-infra-control-plane"},
		},
		{
			name: "Custom machine CIDR with internal publishing",
			installConfig: func() *installconfig.InstallConfig {
				ic := validTestInstallConfig()
				ic.Config.Publish = types.InternalPublishingStrategy
				ic.Config.Networking.MachineNetwork = []types.MachineNetworkEntry{
					{CIDR: *ipnet.MustParseCIDR("172.16.0.0/12")},
				}
				return ic
			}(),
			expectedRuleName:     "test-infra-bootstrap-in-ssh",
			expectedSourceRanges: []string{"172.16.0.0/12"},
			expectedTargetTags:   []string{"test-infra-control-plane"},
		},
		{
			name: "Different infraID",
			installConfig: func() *installconfig.InstallConfig {
				ic := validTestInstallConfig()
				return ic
			}(),
			expectedRuleName:     "custom-id-bootstrap-in-ssh",
			expectedSourceRanges: []string{"0.0.0.0/0"},
			expectedTargetTags:   []string{"custom-id-control-plane"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var infraID string
			if tt.name == "Different infraID" {
				infraID = "custom-id"
			} else {
				infraID = "test-infra"
			}

			clusterID := &installconfig.ClusterID{
				InfraID: infraID,
			}

			rules := createBootstrapFirewallRuleForCAPG(tt.installConfig, clusterID)

			if !assert.Len(t, rules, 1, "should return exactly one rule") {
				return
			}

			rule := rules[0]
			assert.Equal(t, tt.expectedRuleName, rule.Name)
			assert.Equal(t, tt.expectedTargetTags, rule.TargetTags)
			assert.Empty(t, rule.SourceTags)
			assert.Equal(t, tt.expectedSourceRanges, rule.SourceRanges)
			assert.Equal(t, resourceDescription, rule.Description)
			assert.Equal(t, capg.FirewallRuleDirectionIngress, rule.Direction)
			assert.Equal(t, 1000, rule.Priority)

			// Check allowed protocols and ports
			assert.Len(t, rule.Allowed, 2, "should have two allowed descriptors")
			assert.Contains(t, rule.Allowed, capg.FirewallDescriptor{
				IPProtocol: capg.FirewallProtocolTCP,
				Ports:      []string{"22"},
			})
			assert.Contains(t, rule.Allowed, capg.FirewallDescriptor{
				IPProtocol: capg.FirewallProtocolICMP,
			})
		})
	}
}

func TestCreateBootstrapFirewallRuleForCAPG_BootstrapTagFormat(t *testing.T) {
	// Verify that the bootstrap tag follows the expected format
	installConfig := validTestInstallConfig()
	clusterID := &installconfig.ClusterID{
		InfraID: "my-test-cluster",
	}

	rules := createBootstrapFirewallRuleForCAPG(installConfig, clusterID)

	if !assert.Len(t, rules, 1) {
		return
	}
	rule := rules[0]

	// Bootstrap uses control-plane tag, not a separate bootstrap tag
	assert.Equal(t, []string{"my-test-cluster-control-plane"}, rule.TargetTags,
		"Bootstrap should target control-plane tag")
	assert.Equal(t, "my-test-cluster-bootstrap-in-ssh", rule.Name,
		"Bootstrap rule name should include 'bootstrap-in-ssh'")
}

func TestFirewallRules_Priorities(t *testing.T) {
	// Verify firewall rules have priorities set correctly
	installConfig := validTestInstallConfig()
	clusterID := &installconfig.ClusterID{
		InfraID: "test-infra",
	}

	allRules := append(
		createFirewallRulesForCAPG(installConfig, clusterID),
		createBootstrapFirewallRuleForCAPG(installConfig, clusterID)...,
	)

	// Note: Most rules have priority 1000, but the API rule has priority 0 (default)
	// This is because the API rule doesn't explicitly set Priority in the code
	expectedPriorities := map[string]int{
		"test-infra-control-plane":    1000,
		"test-infra-etcd":             1000,
		"test-infra-health-checks":    1000,
		"test-infra-internal-cluster": 1000,
		"test-infra-api":              0, // Not explicitly set in code
		"test-infra-internal-network": 1000,
		"test-infra-bootstrap-in-ssh": 1000,
	}

	for _, rule := range allRules {
		expectedPriority, ok := expectedPriorities[rule.Name]
		if assert.True(t, ok, "Rule %s should have known priority", rule.Name) {
			assert.Equal(t, expectedPriority, rule.Priority, "Rule %s should have priority %d", rule.Name, expectedPriority)
		}
	}
}

func TestFirewallRules_NoSourceRangesAndSourceTagsMixed(t *testing.T) {
	// Verify that rules either use source ranges OR source tags, not both
	installConfig := validTestInstallConfig()
	clusterID := &installconfig.ClusterID{
		InfraID: "test-infra",
	}

	allRules := append(
		createFirewallRulesForCAPG(installConfig, clusterID),
		createBootstrapFirewallRuleForCAPG(installConfig, clusterID)...,
	)

	for _, rule := range allRules {
		// It's valid to have both empty, but if one is set, verify the logic
		if len(rule.SourceRanges) > 0 {
			// Rules with source ranges typically don't use source tags
			t.Logf("Rule %s uses source ranges: %v", rule.Name, rule.SourceRanges)
		}
		if len(rule.SourceTags) > 0 {
			// Rules with source tags typically don't use source ranges
			assert.Empty(t, rule.SourceRanges, "Rule %s has both source tags and ranges", rule.Name)
		}
	}
}

func TestFirewallRules_InternalClusterProtocols(t *testing.T) {
	// Detailed test for the internal-cluster rule which has the most complex allowed list
	installConfig := validTestInstallConfig()
	clusterID := &installconfig.ClusterID{
		InfraID: "test-infra",
	}

	rules := createFirewallRulesForCAPG(installConfig, clusterID)

	var internalClusterRule *capg.FirewallRule
	for i := range rules {
		if rules[i].Name == "test-infra-internal-cluster" {
			internalClusterRule = &rules[i]
			break
		}
	}

	if !assert.NotNil(t, internalClusterRule, "internal-cluster rule should exist") {
		return
	}

	// Verify all 8 expected protocol/port combinations
	expectedAllowed := []capg.FirewallDescriptor{
		{IPProtocol: capg.FirewallProtocolTCP, Ports: []string{"30000-32767"}},
		{IPProtocol: capg.FirewallProtocolUDP, Ports: []string{"30000-32767"}},
		{IPProtocol: capg.FirewallProtocolTCP, Ports: []string{"9000-9999"}},
		{IPProtocol: capg.FirewallProtocolUDP, Ports: []string{"9000-9999"}},
		{IPProtocol: capg.FirewallProtocolUDP, Ports: []string{"4789", "6081"}},
		{IPProtocol: capg.FirewallProtocolUDP, Ports: []string{"500", "4500"}},
		{IPProtocol: capg.FirewallProtocolTCP, Ports: []string{"10250"}},
		{IPProtocol: capg.FirewallProtocolESP},
	}

	// The order might differ, so check each expected descriptor is present
	for _, expected := range expectedAllowed {
		assert.Contains(t, internalClusterRule.Allowed, expected,
			"internal-cluster rule should contain %+v", expected)
	}

	assert.Len(t, internalClusterRule.Allowed, 8, "internal-cluster rule should have exactly 8 allowed descriptors")
}
