package gcp

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	capg "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"

	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/gcp"
)

const (
	testInfraID     = "test-infra-abcde"
	testMachineCIDR = "10.0.0.0/16"

	// publicRegion is any region served by the public GCP health check
	// probers, as opposed to the GCD regions in gcdHealthCheckSourceRanges.
	publicRegion = "us-central1"
	gcdRegion    = "u-germany-northeast1"
)

// publicProberRanges are the health check prober ranges used in every public
// GCP region. gcdProberRanges are the equivalents for gcdRegion.
var (
	publicProberRanges   = []string{"35.191.0.0/16", "130.211.0.0/22"}
	externalLBRanges     = []string{"209.85.152.0/22", "209.85.204.0/22"}
	gcdRegionProberRange = []string{
		"34.3.144.0/23",
		"34.3.151.0/26",
		"34.3.151.64/26",
		"136.124.104.0/22",
		"136.124.108.0/22",
	}
)

func testClusterID() *installconfig.ClusterID {
	return &installconfig.ClusterID{InfraID: testInfraID}
}

// testInstallConfig builds the minimum install config the firewall rule
// helpers read, then applies each option. Subnet names are deliberately left
// empty so that GenerateClusterAssets does not call out to the compute API.
func testInstallConfig(opts ...func(*types.InstallConfig)) *installconfig.InstallConfig {
	ic := &types.InstallConfig{
		Publish: types.ExternalPublishingStrategy,
		Networking: &types.Networking{
			MachineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR(testMachineCIDR)},
			},
		},
		ControlPlane: &types.MachinePool{},
		Compute:      []types.MachinePool{{}},
		Platform: types.Platform{
			GCP: &gcp.Platform{
				ProjectID: "test-project",
				Region:    publicRegion,
			},
		},
	}

	for _, opt := range opts {
		opt(ic)
	}

	return &installconfig.InstallConfig{
		AssetBase: installconfig.AssetBase{Config: ic},
	}
}

func withRegion(region string) func(*types.InstallConfig) {
	return func(ic *types.InstallConfig) { ic.Platform.GCP.Region = region }
}

func withPublish(publish types.PublishingStrategy) func(*types.InstallConfig) {
	return func(ic *types.InstallConfig) { ic.Publish = publish }
}

func withNetworkProjectID(projectID string) func(*types.InstallConfig) {
	return func(ic *types.InstallConfig) { ic.Platform.GCP.NetworkProjectID = projectID }
}

func withFirewallRulesManagement(policy gcp.FirewallRulesManagementPolicy) func(*types.InstallConfig) {
	return func(ic *types.InstallConfig) { ic.Platform.GCP.FirewallRulesManagement = policy }
}

func withAPIServerPublishing(strategy string) func(*types.InstallConfig) {
	return func(ic *types.InstallConfig) {
		ic.OperatorPublishingStrategy = &types.OperatorPublishingStrategy{APIServer: strategy}
	}
}

func ruleNames(rules []capg.FirewallRule) []string {
	names := make([]string, 0, len(rules))
	for _, rule := range rules {
		names = append(names, rule.Name)
	}
	return names
}

// findRule returns the rule with the given name, failing the test if it is
// absent. Callers that expect a rule to be missing should assert on
// ruleNames instead.
func findRule(t *testing.T, rules []capg.FirewallRule, name string) capg.FirewallRule {
	t.Helper()

	for _, rule := range rules {
		if rule.Name == name {
			return rule
		}
	}

	require.Failf(t, "firewall rule not found", "no rule named %q in %v", name, ruleNames(rules))
	return capg.FirewallRule{}
}

// allRules is the set of rules the installer declares on GCPCluster for an
// install that manages its own firewall rules in a public GCP region.
func allRules() []string {
	return []string{
		fmt.Sprintf("%s-bootstrap-in-ssh", testInfraID),
		fmt.Sprintf("%s-control-plane", testInfraID),
		fmt.Sprintf("%s-etcd", testInfraID),
		fmt.Sprintf("%s-health-checks", testInfraID),
		fmt.Sprintf("%s-internal-cluster", testInfraID),
		fmt.Sprintf("%s-api", testInfraID),
		fmt.Sprintf("%s-internal-network", testInfraID),
		fmt.Sprintf("%s-bootstrap-in-konnectivity", testInfraID),
	}
}

// gcpClusterFrom runs GenerateClusterAssets and returns the GCPCluster it
// produced.
func gcpClusterFrom(t *testing.T, installConfig *installconfig.InstallConfig) *capg.GCPCluster {
	t.Helper()

	out, err := GenerateClusterAssets(installConfig, testClusterID())
	require.NoError(t, err)
	require.Len(t, out.Manifests, 1)

	gcpCluster, ok := out.Manifests[0].Object.(*capg.GCPCluster)
	require.True(t, ok, "expected a GCPCluster manifest, got %T", out.Manifests[0].Object)

	return gcpCluster
}

// TestGenerateClusterAssetsFirewallRuleManagement covers which installs hand
// their firewall rules to CAPG. The installer only declares rules when it owns
// them: a shared VPC install leaves the rules to the host project, and an
// install that opts out of rule management leaves them to the user. In both
// cases CAPG's own default rules must be turned off as well, otherwise it
// creates its permissive defaults in a network the cluster does not own.
func TestGenerateClusterAssetsFirewallRuleManagement(t *testing.T) {
	cases := []struct {
		name               string
		opts               []func(*types.InstallConfig)
		expectedManagement capg.RulesManagementPolicy
		expectedRules      []string
	}{
		{
			name:               "managed by default declares every rule",
			expectedManagement: capg.RulesManagementManaged,
			expectedRules:      allRules(),
		},
		{
			name:               "explicitly managed declares every rule",
			opts:               []func(*types.InstallConfig){withFirewallRulesManagement(gcp.ManagedFirewallRules)},
			expectedManagement: capg.RulesManagementManaged,
			expectedRules:      allRules(),
		},
		{
			name:               "unmanaged declares no rules",
			opts:               []func(*types.InstallConfig){withFirewallRulesManagement(gcp.UnmanagedFirewallRules)},
			expectedManagement: capg.RulesManagementUnmanaged,
			expectedRules:      []string{},
		},
		{
			name:               "shared vpc declares no rules",
			opts:               []func(*types.InstallConfig){withNetworkProjectID("host-project")},
			expectedManagement: capg.RulesManagementUnmanaged,
			expectedRules:      []string{},
		},
		{
			name: "shared vpc declares no rules even when managed is requested",
			opts: []func(*types.InstallConfig){
				withNetworkProjectID("host-project"),
				withFirewallRulesManagement(gcp.ManagedFirewallRules),
			},
			expectedManagement: capg.RulesManagementUnmanaged,
			expectedRules:      []string{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			firewall := gcpClusterFrom(t, testInstallConfig(tc.opts...)).Spec.Network.Firewall

			assert.Equal(t, tc.expectedManagement, firewall.DefaultRulesManagement)
			assert.Equal(t, tc.expectedRules, ruleNames(firewall.FirewallRules))
		})
	}
}

// TestGenerateClusterAssetsRuleNames pins the generated rule names. The
// destroy path in pkg/destroy/gcp and removeBootstrapFirewallRules both find
// these rules by name, so renaming one here leaks the rule on teardown.
func TestGenerateClusterAssetsRuleNames(t *testing.T) {
	t.Run("public region", func(t *testing.T) {
		gcpCluster := gcpClusterFrom(t, testInstallConfig())

		assert.Equal(t, allRules(), ruleNames(gcpCluster.Spec.Network.Firewall.FirewallRules))
	})

	t.Run("gcd region adds the service health check rule", func(t *testing.T) {
		gcpCluster := gcpClusterFrom(t, testInstallConfig(withRegion(gcdRegion)))

		assert.ElementsMatch(t,
			append(allRules(), fmt.Sprintf("%s-gcd-service-health-checks", testInfraID)),
			ruleNames(gcpCluster.Spec.Network.Firewall.FirewallRules))
	})
}

// TestHealthCheckSourceRanges covers the source ranges on the API health check
// rule. The public GCP prober ranges apply everywhere, a public API adds the
// ranges the external load balancer probes from, and a GCD region adds the
// sovereign prober ranges that region uses instead.
func TestHealthCheckSourceRanges(t *testing.T) {
	cases := []struct {
		name     string
		opts     []func(*types.InstallConfig)
		expected []string
	}{
		{
			name:     "public region, internal publishing",
			opts:     []func(*types.InstallConfig){withPublish(types.InternalPublishingStrategy)},
			expected: publicProberRanges,
		},
		{
			name:     "public region, external publishing adds load balancer ranges",
			opts:     []func(*types.InstallConfig){withPublish(types.ExternalPublishingStrategy)},
			expected: append(append([]string{}, publicProberRanges...), externalLBRanges...),
		},
		{
			name: "gcd region, internal publishing adds sovereign ranges",
			opts: []func(*types.InstallConfig){
				withRegion(gcdRegion),
				withPublish(types.InternalPublishingStrategy),
			},
			expected: append(append([]string{}, publicProberRanges...), gcdRegionProberRange...),
		},
		{
			name: "gcd region, external publishing adds both",
			opts: []func(*types.InstallConfig){
				withRegion(gcdRegion),
				withPublish(types.ExternalPublishingStrategy),
			},
			expected: append(append(append([]string{}, publicProberRanges...), externalLBRanges...), gcdRegionProberRange...),
		},
		{
			name:     "every gcd region contributes its own ranges",
			opts:     []func(*types.InstallConfig){withRegion("u-france-east1"), withPublish(types.InternalPublishingStrategy)},
			expected: append(append([]string{}, publicProberRanges...), gcdHealthCheckSourceRanges["u-france-east1"]...),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rules := createFirewallRulesForCAPG(testInstallConfig(tc.opts...), testClusterID())
			rule := findRule(t, rules, fmt.Sprintf("%s-health-checks", testInfraID))

			assert.Equal(t, tc.expected, rule.SourceRanges)
		})
	}
}

// TestGCDServiceHealthCheckRule covers the additive rule that exists only in
// GCD regions. cloud-provider-gcp opens the service-type LoadBalancer health
// check ports using the public GCP prober ranges, which never reach a GCD
// cluster, so the installer declares the same ports with the GCD ranges.
func TestGCDServiceHealthCheckRule(t *testing.T) {
	ruleName := fmt.Sprintf("%s-gcd-service-health-checks", testInfraID)

	t.Run("absent in a public region", func(t *testing.T) {
		rules := createFirewallRulesForCAPG(testInstallConfig(), testClusterID())

		assert.NotContains(t, ruleNames(rules), ruleName)
	})

	t.Run("uses the gcd ranges and targets both machine roles", func(t *testing.T) {
		rules := createFirewallRulesForCAPG(testInstallConfig(withRegion(gcdRegion)), testClusterID())
		rule := findRule(t, rules, ruleName)

		assert.Equal(t, gcdRegionProberRange, rule.SourceRanges)
		assert.ElementsMatch(t,
			[]string{fmt.Sprintf("%s-worker", testInfraID), fmt.Sprintf("%s-control-plane", testInfraID)},
			rule.TargetTags)
		assert.Equal(t, []capg.FirewallDescriptor{{
			IPProtocol: capg.FirewallProtocolTCP,
			Ports:      []string{"10256", "30000-32767"},
		}}, rule.Allowed)
	})
}

// TestAPIFirewallRuleSourceRanges covers who may reach the kube-apiserver.
//
// An empty source range is load bearing rather than an omission: CAPG passes
// the empty list through to GCP, which fills it in with 0.0.0.0/0 on a rule
// that names no source tags. Narrowing a public API to the machine CIDR, or
// widening an internal one to the internet, both show up here.
func TestAPIFirewallRuleSourceRanges(t *testing.T) {
	cases := []struct {
		name     string
		opts     []func(*types.InstallConfig)
		expected []string
	}{
		{
			name:     "external publishing is open to the internet",
			opts:     []func(*types.InstallConfig){withPublish(types.ExternalPublishingStrategy)},
			expected: []string{},
		},
		{
			name:     "unset publishing defaults to external",
			opts:     []func(*types.InstallConfig){withPublish("")},
			expected: []string{},
		},
		{
			name:     "internal publishing is limited to the machine CIDR",
			opts:     []func(*types.InstallConfig){withPublish(types.InternalPublishingStrategy)},
			expected: []string{testMachineCIDR},
		},
		{
			name: "internal publishing with an external apiserver is open to the internet",
			opts: []func(*types.InstallConfig){
				withPublish(types.InternalPublishingStrategy),
				withAPIServerPublishing("External"),
			},
			expected: []string{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rules := createFirewallRulesForCAPG(testInstallConfig(tc.opts...), testClusterID())
			rule := findRule(t, rules, fmt.Sprintf("%s-api", testInfraID))

			assert.Equal(t, tc.expected, rule.SourceRanges)
			assert.Equal(t, []string{fmt.Sprintf("%s-control-plane", testInfraID)}, rule.TargetTags)
		})
	}
}

// TestBootstrapFirewallRuleSourceRanges covers who may SSH to the bootstrap
// node. Note that this rule branches on Publish alone, while the API rule
// above branches on PublicAPI(), which also honours
// OperatorPublishingStrategy. The last case pins that divergence: an install
// publishing its apiserver externally from an otherwise internal cluster opens
// the API to the internet but keeps bootstrap SSH on the machine CIDR.
func TestBootstrapFirewallRuleSourceRanges(t *testing.T) {
	cases := []struct {
		name     string
		opts     []func(*types.InstallConfig)
		expected []string
	}{
		{
			name:     "external publishing is open to the internet",
			opts:     []func(*types.InstallConfig){withPublish(types.ExternalPublishingStrategy)},
			expected: []string{"0.0.0.0/0"},
		},
		{
			name:     "internal publishing is limited to the machine CIDR",
			opts:     []func(*types.InstallConfig){withPublish(types.InternalPublishingStrategy)},
			expected: []string{testMachineCIDR},
		},
		{
			name:     "unset publishing is limited to the machine CIDR",
			opts:     []func(*types.InstallConfig){withPublish("")},
			expected: []string{testMachineCIDR},
		},
		{
			name: "internal publishing with an external apiserver stays on the machine CIDR",
			opts: []func(*types.InstallConfig){
				withPublish(types.InternalPublishingStrategy),
				withAPIServerPublishing("External"),
			},
			expected: []string{testMachineCIDR},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rules := createBootstrapFirewallRuleForCAPG(testInstallConfig(tc.opts...), testClusterID())
			require.Len(t, rules, 1)

			assert.Equal(t, fmt.Sprintf("%s-bootstrap-in-ssh", testInfraID), rules[0].Name)
			assert.Equal(t, tc.expected, rules[0].SourceRanges)
			assert.Equal(t, []string{fmt.Sprintf("%s-control-plane", testInfraID)}, rules[0].TargetTags)
		})
	}
}

// TestFirewallRuleTargets covers the tags each rule applies to. A rule that
// targets the wrong role either leaves a port closed that the cluster needs or
// opens one on machines that should not serve it.
func TestFirewallRuleTargets(t *testing.T) {
	masterTag := fmt.Sprintf("%s-control-plane", testInfraID)
	workerTag := fmt.Sprintf("%s-worker", testInfraID)

	cases := []struct {
		rule               string
		expectedTargetTags []string
		expectedSourceTags []string
	}{
		{
			rule:               fmt.Sprintf("%s-control-plane", testInfraID),
			expectedTargetTags: []string{masterTag},
			expectedSourceTags: []string{workerTag, masterTag},
		},
		{
			rule:               fmt.Sprintf("%s-etcd", testInfraID),
			expectedTargetTags: []string{masterTag},
			expectedSourceTags: []string{masterTag},
		},
		{
			rule:               fmt.Sprintf("%s-health-checks", testInfraID),
			expectedTargetTags: []string{masterTag},
			expectedSourceTags: []string{},
		},
		{
			rule:               fmt.Sprintf("%s-internal-cluster", testInfraID),
			expectedTargetTags: []string{workerTag, masterTag},
			expectedSourceTags: []string{workerTag, masterTag},
		},
		{
			rule:               fmt.Sprintf("%s-internal-network", testInfraID),
			expectedTargetTags: []string{workerTag, masterTag},
			expectedSourceTags: []string{},
		},
		{
			rule:               fmt.Sprintf("%s-bootstrap-in-konnectivity", testInfraID),
			expectedTargetTags: []string{masterTag},
			expectedSourceTags: []string{workerTag, masterTag},
		},
	}

	rules := createFirewallRulesForCAPG(testInstallConfig(), testClusterID())

	for _, tc := range cases {
		t.Run(tc.rule, func(t *testing.T) {
			rule := findRule(t, rules, tc.rule)

			assert.ElementsMatch(t, tc.expectedTargetTags, rule.TargetTags)
			assert.ElementsMatch(t, tc.expectedSourceTags, rule.SourceTags)
			assert.Equal(t, capg.FirewallRuleDirectionIngress, rule.Direction)
			assert.Equal(t, resourceDescription, rule.Description)
		})
	}
}

// TestInternalNetworkRuleSourceRanges covers SSH and ICMP reachability, which
// is scoped to the machine network regardless of publishing strategy.
func TestInternalNetworkRuleSourceRanges(t *testing.T) {
	for _, publish := range []types.PublishingStrategy{types.ExternalPublishingStrategy, types.InternalPublishingStrategy} {
		t.Run(string(publish), func(t *testing.T) {
			rules := createFirewallRulesForCAPG(testInstallConfig(withPublish(publish)), testClusterID())
			rule := findRule(t, rules, fmt.Sprintf("%s-internal-network", testInfraID))

			assert.Equal(t, []string{testMachineCIDR}, rule.SourceRanges)
			assert.Equal(t, []capg.FirewallDescriptor{
				{IPProtocol: capg.FirewallProtocolTCP, Ports: []string{"22"}},
				{IPProtocol: capg.FirewallProtocolICMP},
			}, rule.Allowed)
		})
	}
}
