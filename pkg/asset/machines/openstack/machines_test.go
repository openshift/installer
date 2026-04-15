package openstack

import (
	"fmt"
	"net"
	"strings"
	"testing"

	machinev1 "github.com/openshift/api/machine/v1"
	"k8s.io/utils/ptr"
	capo "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	"sigs.k8s.io/yaml"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
)

func mpWithZones(zones ...string) func(*openstack.MachinePool) {
	return func(mpool *openstack.MachinePool) {
		mpool.Zones = zones
	}
}

func mpWithRootVolumeZones(zones ...string) func(*openstack.MachinePool) {
	return func(mpool *openstack.MachinePool) {
		if mpool.RootVolume != nil {
			mpool.RootVolume.Zones = zones
		} else {
			mpool.RootVolume = &openstack.RootVolume{Zones: zones}
		}
	}
}

func mpWithRootVolumeTypes(types ...string) func(*openstack.MachinePool) {
	return func(mpool *openstack.MachinePool) {
		if mpool.RootVolume != nil {
			mpool.RootVolume.Types = types
		} else {
			mpool.RootVolume = &openstack.RootVolume{Types: types}
		}
	}
}

func generateMachinePool(options ...func(*openstack.MachinePool)) openstack.MachinePool {
	mpool := openstack.MachinePool{}
	for _, apply := range options {
		apply(&mpool)
	}
	return mpool
}

func TestFailureDomains(t *testing.T) {
	type checkFunc func([]machinev1.OpenStackFailureDomain, error) error
	check := func(fns ...checkFunc) []checkFunc { return fns }

	hasNFailureDomains := func(want int) checkFunc {
		return func(fds []machinev1.OpenStackFailureDomain, _ error) error {
			if have := len(fds); want != have {
				return fmt.Errorf("expected %d failure domains, got %d", want, have)
			}
			return nil
		}
	}

	hasComputeZones := func(wantZones ...string) checkFunc {
		return func(fds []machinev1.OpenStackFailureDomain, _ error) error {
			haveZones := make([]string, len(fds))
			for i := range fds {
				haveZones[i] = fds[i].AvailabilityZone
			}

			if wantLen, haveLen := len(wantZones), len(haveZones); wantLen != haveLen {
				return fmt.Errorf("expected compute zones %v (len %d), got %v (len %d)", wantZones, wantLen, haveZones, haveLen)
			}

			for i := range fds {
				if want, have := wantZones[i], haveZones[i]; want != have {
					return fmt.Errorf("expected compute zones %v, got %v", wantZones, haveZones)
				}
			}

			return nil
		}
	}

	hasNilRootVolume := func(fds []machinev1.OpenStackFailureDomain, _ error) error {
		for i := range fds {
			if fds[i].RootVolume != nil {
				return fmt.Errorf("failure domain %d has unexpectedly non-nil RootVolume", i)
			}
		}
		return nil
	}

	hasRootVolumeZones := func(wantZones ...string) checkFunc {
		return func(fds []machinev1.OpenStackFailureDomain, _ error) error {
			haveZones := make([]string, len(fds))
			for i := range fds {
				if fds[i].RootVolume == nil {
					return fmt.Errorf("failure domain %d has unexpectedly nil RootVolume", i)
				}
				haveZones[i] = fds[i].RootVolume.AvailabilityZone
			}

			if wantLen, haveLen := len(wantZones), len(haveZones); wantLen != haveLen {
				return fmt.Errorf("expected root volume zones %v, got %v", wantZones, haveZones)
			}

			for i := range fds {
				if want, have := wantZones[i], haveZones[i]; want != have {
					return fmt.Errorf("expected root volume zones %v, got %v", wantZones, haveZones)
				}
			}

			return nil
		}
	}

	hasRootVolumeTypes := func(wantTypes ...string) checkFunc {
		return func(fds []machinev1.OpenStackFailureDomain, _ error) error {
			haveTypes := make([]string, len(fds))
			for i := range fds {
				if fds[i].RootVolume == nil {
					return fmt.Errorf("failure domain %d has unexpectedly nil RootVolume", i)
				}
				haveTypes[i] = fds[i].RootVolume.VolumeType
			}

			if wantLen, haveLen := len(wantTypes), len(haveTypes); wantLen != haveLen {
				return fmt.Errorf("expected root volume types %v, got %v", wantTypes, haveTypes)
			}

			for i := range fds {
				if want, have := wantTypes[i], haveTypes[i]; want != have {
					return fmt.Errorf("expected root volume types %v, got %v", wantTypes, haveTypes)
				}
			}

			return nil
		}
	}

	doesNotPanic := func(_ []machinev1.OpenStackFailureDomain, have error) error {
		if have != nil {
			return fmt.Errorf("unexpected panic: %w", have)
		}
		return nil
	}

	panicsWith := func(want string) checkFunc {
		return func(_ []machinev1.OpenStackFailureDomain, have error) error {
			if have == nil {
				return fmt.Errorf("unexpectedly, didn't panic")
			}
			if have := fmt.Sprintf("%v", have); !strings.Contains(have, want) {
				return fmt.Errorf("expected panic with %q, got %q", want, have)
			}
			return nil
		}
	}

	for _, tc := range [...]struct {
		name   string
		mpool  openstack.MachinePool
		checks []checkFunc
	}{
		{
			"no_zones",
			generateMachinePool(),
			check(
				hasNFailureDomains(1),
				hasComputeZones(""),
				hasNilRootVolume,
				doesNotPanic,
			),
		},
		{
			"one_compute_zone",
			generateMachinePool(
				mpWithZones("one"),
			),
			check(
				hasNFailureDomains(1),
				hasComputeZones("one"),
				hasNilRootVolume,
				doesNotPanic,
			),
		},
		{
			"three_compute_zones",
			generateMachinePool(
				mpWithZones("one", "two", "three"),
			),
			check(
				hasNFailureDomains(3),
				hasComputeZones("one", "two", "three"),
				hasNilRootVolume,
				doesNotPanic,
			),
		},
		{
			"three_compute_zones_one_root_volume_zone",
			generateMachinePool(
				mpWithZones("one", "two", "three"),
				mpWithRootVolumeZones("volume_one"),
				mpWithRootVolumeTypes("type-1"),
			),
			check(
				hasNFailureDomains(3),
				hasComputeZones("one", "two", "three"),
				hasRootVolumeZones("volume_one", "volume_one", "volume_one"),
				hasRootVolumeTypes("type-1", "type-1", "type-1"),
				doesNotPanic,
			),
		},
		{
			"one_compute_zone_three_root_volume_zones",
			generateMachinePool(
				mpWithZones("one"),
				mpWithRootVolumeZones("volume_one", "volume_two", "volume_three"),
				mpWithRootVolumeTypes("type-1"),
			),
			check(
				hasNFailureDomains(3),
				hasComputeZones("one", "one", "one"),
				hasRootVolumeZones("volume_one", "volume_two", "volume_three"),
				hasRootVolumeTypes("type-1", "type-1", "type-1"),
				doesNotPanic,
			),
		},
		{
			"three_compute_zone_two_root_volume_zones_panics",
			generateMachinePool(
				mpWithZones("one", "two", "three"),
				mpWithRootVolumeZones("volume_one", "volume_two"),
			),
			check(
				// We have to check for a partial result here, because the mapping
				// of compute zones to root volume zones is handled in a map therefore
				// the order is not deterministic.
				panicsWith("availability zones should have equal length"),
			),
		},
		{
			"three_compute_zones_three_root_volume_types",
			generateMachinePool(
				mpWithZones("one", "two", "three"),
				mpWithRootVolumeZones("volume_one", "volume_two", "volume_three"),
				mpWithRootVolumeTypes("type-1", "type-2", "type-3"),
			),
			check(
				hasNFailureDomains(3),
				hasComputeZones("one", "two", "three"),
				hasRootVolumeZones("volume_one", "volume_two", "volume_three"),
				hasRootVolumeTypes("type-1", "type-2", "type-3"),
				doesNotPanic,
			),
		},
		{
			"three_root_volume_types",
			generateMachinePool(
				mpWithRootVolumeTypes("type-1", "type-2", "type-3"),
			),
			check(
				hasNFailureDomains(3),
				hasComputeZones("", "", ""),
				hasRootVolumeZones("", "", ""),
				hasRootVolumeTypes("type-1", "type-2", "type-3"),
				doesNotPanic,
			),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			failureDomains, recoveredPanic := func() (fds []machinev1.OpenStackFailureDomain, recoveredPanic error) {
				defer func() {
					if r := recover(); r != nil {
						recoveredPanic = fmt.Errorf("%v", r)
					}
				}()

				fds = failureDomainsFromSpec(tc.mpool)
				return
			}()

			for _, check := range tc.checks {
				if err := check(failureDomains, recoveredPanic); err != nil {
					t.Error(err)
				}
			}
		})
	}
}

// installConfigYAML is the structure used to parse the relevant subset of the
// install-config YAML for the integration tests below.
type installConfigYAML struct {
	Platform struct {
		OpenStack struct {
			Cloud           string   `yaml:"cloud"`
			BootstrapFlavor string   `yaml:"bootstrapFlavor,omitempty"`
			APIVIPs         []string `yaml:"apiVIPs,omitempty"`
			IngressVIPs     []string `yaml:"ingressVIPs,omitempty"`
		} `yaml:"openstack"`
	} `yaml:"platform"`
	ControlPlane struct {
		Platform struct {
			OpenStack struct {
				Type string `yaml:"type"`
			} `yaml:"openstack"`
		} `yaml:"platform"`
		Replicas *int64 `yaml:"replicas,omitempty"`
	} `yaml:"controlPlane"`
}

// parseInstallConfigFromYAML parses a subset install-config YAML into the Go
// types used by GenerateMachines, simulating the YAML-to-struct path that the
// real installer follows.
func parseInstallConfigFromYAML(t *testing.T, raw string) (*types.InstallConfig, *types.MachinePool) {
	t.Helper()

	var cfg installConfigYAML
	if err := yaml.Unmarshal([]byte(raw), &cfg); err != nil {
		t.Fatalf("failed to unmarshal install-config YAML: %v", err)
	}

	ic := &types.InstallConfig{
		Networking: &types.Networking{
			MachineNetwork: []types.MachineNetworkEntry{
				{
					CIDR: ipnet.IPNet{
						IPNet: net.IPNet{
							IP:   net.ParseIP("10.0.0.0"),
							Mask: net.CIDRMask(16, 32),
						},
					},
				},
			},
		},
		Platform: types.Platform{
			OpenStack: &openstack.Platform{
				Cloud:           cfg.Platform.OpenStack.Cloud,
				BootstrapFlavor: cfg.Platform.OpenStack.BootstrapFlavor,
				APIVIPs:         cfg.Platform.OpenStack.APIVIPs,
				IngressVIPs:     cfg.Platform.OpenStack.IngressVIPs,
			},
		},
	}

	replicas := int64(1)
	if cfg.ControlPlane.Replicas != nil {
		replicas = *cfg.ControlPlane.Replicas
	}

	pool := &types.MachinePool{
		Name:     "master",
		Replicas: ptr.To(replicas),
		Platform: types.MachinePoolPlatform{
			OpenStack: &openstack.MachinePool{
				FlavorName: cfg.ControlPlane.Platform.OpenStack.Type,
			},
		},
	}

	return ic, pool
}

// TestBootstrapFlavorIntegration is an end-to-end integration test that
// simulates the complete flow from install-config.yaml parsing through machine
// manifest generation. It verifies that:
//
//   - A bootstrap OpenStackMachine receives the explicitly configured
//     bootstrapFlavor from the install-config.
//   - Control plane OpenStackMachines are not affected and continue to use the
//     control plane pool flavor.
//   - The serialized YAML representation of each manifest contains the expected
//     flavor string.
func TestBootstrapFlavorIntegration(t *testing.T) {
	const (
		clusterID = "integration-test-cluster"
		osImage   = "rhcos-4.14"
	)

	tests := []struct {
		name                string
		installConfigYAML   string
		wantBootstrapFlavor string
		wantCPFlavor        string
	}{
		{
			name: "custom bootstrap flavor propagates to bootstrap manifest",
			installConfigYAML: `
platform:
  openstack:
    cloud: mycloud
    bootstrapFlavor: m1.xlarge
    apiVIPs:
      - 10.0.0.5
    ingressVIPs:
      - 10.0.0.7
controlPlane:
  replicas: 3
  platform:
    openstack:
      type: m1.large
`,
			wantBootstrapFlavor: "m1.xlarge",
			wantCPFlavor:        "m1.large",
		},
		{
			name: "bootstrap flavor with spaces in name is preserved",
			installConfigYAML: `
platform:
  openstack:
    cloud: mycloud
    bootstrapFlavor: "my bootstrap flavor"
    apiVIPs:
      - 10.0.0.5
    ingressVIPs:
      - 10.0.0.7
controlPlane:
  replicas: 1
  platform:
    openstack:
      type: m1.large
`,
			wantBootstrapFlavor: "my bootstrap flavor",
			wantCPFlavor:        "m1.large",
		},
		{
			name: "bootstrap flavor with mixed case is preserved",
			installConfigYAML: `
platform:
  openstack:
    cloud: mycloud
    bootstrapFlavor: Bootstrap-Flavor-XLarge
    apiVIPs:
      - 10.0.0.5
    ingressVIPs:
      - 10.0.0.7
controlPlane:
  replicas: 1
  platform:
    openstack:
      type: m1.large
`,
			wantBootstrapFlavor: "Bootstrap-Flavor-XLarge",
			wantCPFlavor:        "m1.large",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ic, pool := parseInstallConfigFromYAML(t, tt.installConfigYAML)

			// --- Bootstrap machine ---
			bootstrapFiles, err := GenerateMachines(clusterID, ic, pool, osImage, bootstrapRole)
			if err != nil {
				t.Fatalf("GenerateMachines(bootstrap) unexpected error: %v", err)
			}
			if len(bootstrapFiles) == 0 {
				t.Fatal("GenerateMachines(bootstrap) returned no files")
			}

			bootstrapMachine, ok := bootstrapFiles[0].Object.(*capo.OpenStackMachine)
			if !ok {
				t.Fatalf("bootstrap file[0]: expected *capo.OpenStackMachine, got %T", bootstrapFiles[0].Object)
			}
			if bootstrapMachine.Spec.Flavor == nil {
				t.Fatal("bootstrap OpenStackMachine.Spec.Flavor is nil")
			}
			if got := *bootstrapMachine.Spec.Flavor; got != tt.wantBootstrapFlavor {
				t.Errorf("bootstrap machine flavor = %q, want %q", got, tt.wantBootstrapFlavor)
			}

			// Verify the YAML serialization of the bootstrap manifest contains
			// the expected flavor string — this tests the "generated YAML structure
			// is correct" acceptance criterion.
			bootstrapYAML, err := yaml.Marshal(bootstrapMachine)
			if err != nil {
				t.Fatalf("failed to marshal bootstrap OpenStackMachine to YAML: %v", err)
			}
			bootstrapYAMLStr := string(bootstrapYAML)
			if !strings.Contains(bootstrapYAMLStr, tt.wantBootstrapFlavor) {
				t.Errorf("bootstrap machine YAML does not contain flavor %q:\n%s", tt.wantBootstrapFlavor, bootstrapYAMLStr)
			}

			// --- Control plane machines ---
			cpFiles, err := GenerateMachines(clusterID, ic, pool, osImage, masterRole)
			if err != nil {
				t.Fatalf("GenerateMachines(master) unexpected error: %v", err)
			}

			replicas := *pool.Replicas
			expectedCPFiles := int(replicas) * 2 // OpenStackMachine + CAPI Machine per replica
			if len(cpFiles) != expectedCPFiles {
				t.Fatalf("expected %d CP files (replicas=%d), got %d", expectedCPFiles, replicas, len(cpFiles))
			}

			for i := int64(0); i < replicas; i++ {
				fileIdx := i * 2
				cpMachine, ok := cpFiles[fileIdx].Object.(*capo.OpenStackMachine)
				if !ok {
					t.Errorf("CP file[%d]: expected *capo.OpenStackMachine, got %T", fileIdx, cpFiles[fileIdx].Object)
					continue
				}
				if cpMachine.Spec.Flavor == nil {
					t.Errorf("CP file[%d]: OpenStackMachine.Spec.Flavor is nil", fileIdx)
					continue
				}
				if got := *cpMachine.Spec.Flavor; got != tt.wantCPFlavor {
					t.Errorf("CP file[%d] flavor = %q, want %q", fileIdx, got, tt.wantCPFlavor)
				}

				// Verify YAML serialization for each control plane machine.
				cpYAML, err := yaml.Marshal(cpMachine)
				if err != nil {
					t.Fatalf("failed to marshal CP OpenStackMachine[%d] to YAML: %v", i, err)
				}
				if !strings.Contains(string(cpYAML), tt.wantCPFlavor) {
					t.Errorf("CP machine[%d] YAML does not contain flavor %q:\n%s", i, tt.wantCPFlavor, string(cpYAML))
				}
			}
		})
	}
}

// TestBootstrapFlavorBackwardCompatibility verifies that when bootstrapFlavor
// is not specified in the install-config, the bootstrap machine inherits the
// control plane flavor — preserving backward-compatible behavior for existing
// clusters that do not set bootstrapFlavor.
func TestBootstrapFlavorBackwardCompatibility(t *testing.T) {
	const (
		clusterID = "compat-test-cluster"
		osImage   = "rhcos-4.14"
	)

	tests := []struct {
		name                string
		installConfigYAML   string
		wantBootstrapFlavor string
		wantCPFlavor        string
	}{
		{
			name: "no bootstrapFlavor: bootstrap inherits control plane flavor",
			installConfigYAML: `
platform:
  openstack:
    cloud: mycloud
    apiVIPs:
      - 10.0.0.5
    ingressVIPs:
      - 10.0.0.7
controlPlane:
  replicas: 3
  platform:
    openstack:
      type: m1.large
`,
			wantBootstrapFlavor: "m1.large",
			wantCPFlavor:        "m1.large",
		},
		{
			name: "empty bootstrapFlavor string: bootstrap inherits control plane flavor",
			installConfigYAML: `
platform:
  openstack:
    cloud: mycloud
    bootstrapFlavor: ""
    apiVIPs:
      - 10.0.0.5
    ingressVIPs:
      - 10.0.0.7
controlPlane:
  replicas: 1
  platform:
    openstack:
      type: m1.xlarge
`,
			wantBootstrapFlavor: "m1.xlarge",
			wantCPFlavor:        "m1.xlarge",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ic, pool := parseInstallConfigFromYAML(t, tt.installConfigYAML)

			// Bootstrap machine must use the control plane flavor as fallback.
			bootstrapFiles, err := GenerateMachines(clusterID, ic, pool, osImage, bootstrapRole)
			if err != nil {
				t.Fatalf("GenerateMachines(bootstrap) unexpected error: %v", err)
			}
			if len(bootstrapFiles) == 0 {
				t.Fatal("GenerateMachines(bootstrap) returned no files")
			}

			bootstrapMachine, ok := bootstrapFiles[0].Object.(*capo.OpenStackMachine)
			if !ok {
				t.Fatalf("bootstrap file[0]: expected *capo.OpenStackMachine, got %T", bootstrapFiles[0].Object)
			}
			if bootstrapMachine.Spec.Flavor == nil {
				t.Fatal("bootstrap OpenStackMachine.Spec.Flavor is nil")
			}
			if got := *bootstrapMachine.Spec.Flavor; got != tt.wantBootstrapFlavor {
				t.Errorf("bootstrap machine flavor = %q, want %q (should inherit CP flavor when bootstrapFlavor unset)", got, tt.wantBootstrapFlavor)
			}

			// Verify YAML serialization.
			bootstrapYAML, err := yaml.Marshal(bootstrapMachine)
			if err != nil {
				t.Fatalf("failed to marshal bootstrap OpenStackMachine to YAML: %v", err)
			}
			if !strings.Contains(string(bootstrapYAML), tt.wantBootstrapFlavor) {
				t.Errorf("bootstrap machine YAML does not contain expected flavor %q:\n%s", tt.wantBootstrapFlavor, string(bootstrapYAML))
			}

			// Control plane machines must use the pool flavor.
			cpFiles, err := GenerateMachines(clusterID, ic, pool, osImage, masterRole)
			if err != nil {
				t.Fatalf("GenerateMachines(master) unexpected error: %v", err)
			}

			replicas := *pool.Replicas
			for i := int64(0); i < replicas; i++ {
				fileIdx := i * 2
				cpMachine, ok := cpFiles[fileIdx].Object.(*capo.OpenStackMachine)
				if !ok {
					t.Errorf("CP file[%d]: expected *capo.OpenStackMachine, got %T", fileIdx, cpFiles[fileIdx].Object)
					continue
				}
				if cpMachine.Spec.Flavor == nil {
					t.Errorf("CP file[%d]: OpenStackMachine.Spec.Flavor is nil", fileIdx)
					continue
				}
				if got := *cpMachine.Spec.Flavor; got != tt.wantCPFlavor {
					t.Errorf("CP machine[%d] flavor = %q, want %q", i, got, tt.wantCPFlavor)
				}
			}
		})
	}
}

func TestPruneFailureDomains(t *testing.T) {
}
