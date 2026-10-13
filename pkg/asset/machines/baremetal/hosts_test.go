package baremetal

import (
	"fmt"
	"testing"

	baremetalhost "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	machineapi "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/types"
	baremetaltypes "github.com/openshift/installer/pkg/types/baremetal"
)

func TestHosts(t *testing.T) {

	nmstate := `interfaces:
- name: eth0
  type: ethernet
routes:
  config:
  - destination: 198.51.100.0/24
    metric: 150
    next-hop-address: 192.0.2.1
    next-hop-interface: eth1
    table-id: 254
`

	testCases := []struct {
		Scenario        string
		Machines        []machineapi.Machine
		ArbiterMachines []machineapi.Machine
		Config          *types.InstallConfig
		ExpectedSecrets []corev1.Secret
		ExpectedHosts   []baremetalhost.BareMetalHost
		ExpectedError   string
		ExpectedSetting *HostSettings
	}{
		{
			Scenario: "no-platform",
			Config: &types.InstallConfig{
				Platform: types.Platform{
					BareMetal: nil,
				},
			},

			ExpectedError: "no baremetal platform in configuration",
		},
		{
			Scenario: "no-hosts",
			Config:   config().build(),

			ExpectedSetting: settings().build(),
		},
		{
			Scenario: "default",
			Machines: machines(machine("machine-0")),
			Config:   configHosts(hostType("master-0").bmc("usr0", "pwd0").role("master")),

			ExpectedSetting: settings().
				secrets(secret("master-0-bmc-secret").creds("usr0", "pwd0")).
				hosts(host("master-0").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-0").customDeploy().architecture("x86_64")).build(),
		},
		{
			Scenario: "default-norole",
			Machines: machines(machine("machine-0")),
			Config:   configHosts(hostType("master-0").bmc("usr0", "pwd0")),

			ExpectedSetting: settings().
				secrets(secret("master-0-bmc-secret").creds("usr0", "pwd0")).
				hosts(host("master-0").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-0").customDeploy().architecture("x86_64")).build(),
		},
		{
			Scenario: "network-config",
			Machines: machines(machine("machine-0")),
			Config: configHosts(
				hostType("master-0").
					bmc("usr0", "pwd0").
					networkConfig(nmstate)),

			ExpectedSetting: settings().
				secrets(secret("master-0-bmc-secret").creds("usr0", "pwd0")).
				networkConfigSecrets(secret("master-0-network-config-secret").nmstate(nmstate)).
				hosts(
					host("master-0").
						label("installer.openshift.io/role", "control-plane").
						consumerRef("machine-0").
						userDataRef("user-data-secret").
						preprovisioningNetworkDataName("master-0-network-config-secret").
						customDeploy().
						architecture("x86_64")).build(),
		},
		{
			Scenario: "3-hosts-3-machines-norole-all",
			Machines: machines(
				machine("machine-0"),
				machine("machine-1"),
				machine("machine-2")),
			Config: configHosts(
				hostType("master-0").bmc("usr0", "pwd0"),
				hostType("master-1").bmc("usr1", "pwd1"),
				hostType("master-2").bmc("usr2", "pwd2")),

			ExpectedSetting: settings().
				secrets(
					secret("master-0-bmc-secret").creds("usr0", "pwd0"),
					secret("master-1-bmc-secret").creds("usr1", "pwd1"),
					secret("master-2-bmc-secret").creds("usr2", "pwd2")).
				hosts(
					host("master-0").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-0").customDeploy().architecture("x86_64"),
					host("master-1").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-1").customDeploy().architecture("x86_64"),
					host("master-2").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-2").customDeploy().architecture("x86_64")).build(),
		},
		{
			Scenario: "4-hosts-3-machines",
			Machines: machines(
				machine("machine-0"),
				machine("machine-1"),
				machine("machine-2")),
			Config: configHosts(
				hostType("master-0").bmc("usr0", "pwd0").role("master"),
				hostType("master-1").bmc("usr1", "pwd1").role("master"),
				hostType("master-2").bmc("usr2", "pwd2").role("master"),
				hostType("worker-0").bmc("usr3", "pwd3").role("worker")),

			ExpectedSetting: settings().
				secrets(
					secret("master-0-bmc-secret").creds("usr0", "pwd0"),
					secret("master-1-bmc-secret").creds("usr1", "pwd1"),
					secret("master-2-bmc-secret").creds("usr2", "pwd2"),
					secret("worker-0-bmc-secret").creds("usr3", "pwd3")).
				hosts(
					host("master-0").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-0").customDeploy().architecture("x86_64"),
					host("master-1").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-1").customDeploy().architecture("x86_64"),
					host("master-2").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-2").customDeploy().architecture("x86_64"),
					host("worker-0").annotation("baremetalhost.metal3.io/paused", "").architecture("x86_64"),
				).build(),
		},
		{
			Scenario: "4-hosts-3-machines-norole",
			Machines: machines(
				machine("machine-0"),
				machine("machine-1"),
				machine("machine-2")),
			Config: configHosts(
				hostType("master-0").bmc("usr0", "pwd0"),
				hostType("master-1").bmc("usr1", "pwd1"),
				hostType("master-2").bmc("usr2", "pwd2"),
				hostType("worker-0").bmc("wrk0", "pwd0")),

			ExpectedSetting: settings().
				secrets(
					secret("master-0-bmc-secret").creds("usr0", "pwd0"),
					secret("master-1-bmc-secret").creds("usr1", "pwd1"),
					secret("master-2-bmc-secret").creds("usr2", "pwd2"),
					secret("worker-0-bmc-secret").creds("wrk0", "pwd0")).
				hosts(
					host("master-0").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-0").customDeploy().architecture("x86_64"),
					host("master-1").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-1").customDeploy().architecture("x86_64"),
					host("master-2").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-2").customDeploy().architecture("x86_64"),
					host("worker-0").annotation("baremetalhost.metal3.io/paused", "").architecture("x86_64"),
				).build(),
		},
		{
			Scenario: "5-hosts-3-machines",
			Machines: machines(
				machine("machine-0"),
				machine("machine-1"),
				machine("machine-2")),
			Config: configHosts(
				hostType("master-0").bmc("usr0", "pwd0").role("master"),
				hostType("master-1").bmc("usr1", "pwd1").role("master"),
				hostType("master-2").bmc("usr2", "pwd2").role("master"),
				hostType("worker-0").bmc("wrk0", "pwd0").role("worker"),
				hostType("worker-1").bmc("wrk1", "pwd1").role("worker")),

			ExpectedSetting: settings().
				secrets(
					secret("master-0-bmc-secret").creds("usr0", "pwd0"),
					secret("master-1-bmc-secret").creds("usr1", "pwd1"),
					secret("master-2-bmc-secret").creds("usr2", "pwd2"),
					secret("worker-0-bmc-secret").creds("wrk0", "pwd0"),
					secret("worker-1-bmc-secret").creds("wrk1", "pwd1")).
				hosts(
					host("master-0").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-0").customDeploy().architecture("x86_64"),
					host("master-1").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-1").customDeploy().architecture("x86_64"),
					host("master-2").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-2").customDeploy().architecture("x86_64"),
					host("worker-0").annotation("baremetalhost.metal3.io/paused", "").architecture("x86_64"),
					host("worker-1").annotation("baremetalhost.metal3.io/paused", "").architecture("x86_64"),
				).build(),
		},
		{
			Scenario: "5-hosts-3-machines-mixed",
			Machines: machines(
				machine("machine-0"),
				machine("machine-1"),
				machine("machine-2")),
			Config: configHosts(
				hostType("master-1").bmc("usr1", "pwd1").role("master"),
				hostType("worker-0").bmc("wrk0", "pwd0").role("worker"),
				hostType("worker-1").bmc("wrk1", "pwd1").role("worker"),
				hostType("master-0").bmc("usr0", "pwd0"),
				hostType("master-2").bmc("usr2", "pwd2")),

			ExpectedSetting: settings().
				secrets(
					secret("master-1-bmc-secret").creds("usr1", "pwd1"),
					secret("worker-0-bmc-secret").creds("wrk0", "pwd0"),
					secret("worker-1-bmc-secret").creds("wrk1", "pwd1"),
					secret("master-0-bmc-secret").creds("usr0", "pwd0"),
					secret("master-2-bmc-secret").creds("usr2", "pwd2")).
				hosts(
					host("master-1").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-0").customDeploy().architecture("x86_64"),
					host("worker-0").annotation("baremetalhost.metal3.io/paused", "").architecture("x86_64"),
					host("worker-1").annotation("baremetalhost.metal3.io/paused", "").architecture("x86_64"),
					host("master-0").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-1").customDeploy().architecture("x86_64"),
					host("master-2").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-2").customDeploy().architecture("x86_64")).build(),
		},
		{
			Scenario: "4-hosts-3-machines-norole-master",
			Machines: machines(
				machine("machine-0"),
				machine("machine-1"),
				machine("machine-2")),
			Config: configHosts(
				hostType("worker-0").bmc("wrk0", "pwd0").role("worker"),
				hostType("master-0").bmc("usr0", "pwd0"),
				hostType("master-1").bmc("usr1", "pwd1"),
				hostType("master-2").bmc("usr2", "pwd2")),

			ExpectedSetting: settings().
				secrets(
					secret("worker-0-bmc-secret").creds("wrk0", "pwd0"),
					secret("master-0-bmc-secret").creds("usr0", "pwd0"),
					secret("master-1-bmc-secret").creds("usr1", "pwd1"),
					secret("master-2-bmc-secret").creds("usr2", "pwd2")).
				hosts(
					host("worker-0").annotation("baremetalhost.metal3.io/paused", "").architecture("x86_64"),
					host("master-0").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-0").customDeploy().architecture("x86_64"),
					host("master-1").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-1").customDeploy().architecture("x86_64"),
					host("master-2").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-2").customDeploy().architecture("x86_64")).build(),
		},
		{
			Scenario: "4-hosts-3-machines-norole-worker",
			Machines: machines(
				machine("machine-0"),
				machine("machine-1"),
				machine("machine-2")),
			Config: configHosts(
				hostType("master-0").bmc("usr0", "pwd0").role("master"),
				hostType("master-1").bmc("usr1", "pwd1").role("master"),
				hostType("master-2").bmc("usr2", "pwd2").role("master"),
				hostType("worker-0").bmc("wrk0", "pwd0")),

			ExpectedSetting: settings().
				secrets(
					secret("master-0-bmc-secret").creds("usr0", "pwd0"),
					secret("master-1-bmc-secret").creds("usr1", "pwd1"),
					secret("master-2-bmc-secret").creds("usr2", "pwd2"),
					secret("worker-0-bmc-secret").creds("wrk0", "pwd0")).
				hosts(
					host("master-0").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-0").customDeploy().architecture("x86_64"),
					host("master-1").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-1").customDeploy().architecture("x86_64"),
					host("master-2").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-2").customDeploy().architecture("x86_64"),
					host("worker-0").annotation("baremetalhost.metal3.io/paused", "").architecture("x86_64")).build(),
		},
		{
			Scenario: "3-hosts-2-masters-no-arbiter-render",
			Machines: machines(
				machine("machine-0"),
				machine("machine-1"),
				machine("machine-2")),
			Config: configHosts(
				hostType("master-0").bmc("usr0", "pwd0").role("master"),
				hostType("master-1").bmc("usr1", "pwd1").role("master"),
				hostType("arbiter-0").bmc("usr2", "pwd2").role("arbiter")),

			ExpectedSetting: settings().
				secrets(
					secret("master-0-bmc-secret").creds("usr0", "pwd0"),
					secret("master-1-bmc-secret").creds("usr1", "pwd1"),
				).
				hosts(
					host("master-0").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-0").customDeploy().architecture("x86_64"),
					host("master-1").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-1").customDeploy().architecture("x86_64")).build(),
		},
		{
			Scenario: "3-hosts-2-masters-1-arbiter",
			Machines: machines(
				machine("machine-0"),
				machine("machine-1")),
			ArbiterMachines: machines(
				machine("machine-2")),
			Config: config().withArbiter(1).hosts(
				hostType("master-0").bmc("usr0", "pwd0").role("master"),
				hostType("master-1").bmc("usr1", "pwd1").role("master"),
				hostType("arbiter-0").bmc("usr2", "pwd2").role("arbiter")).
				build(),

			ExpectedSetting: settings().
				secrets(
					secret("master-0-bmc-secret").creds("usr0", "pwd0"),
					secret("master-1-bmc-secret").creds("usr1", "pwd1"),
					secret("arbiter-0-bmc-secret").creds("usr2", "pwd2"),
				).
				hosts(
					host("master-0").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-0").customDeploy().architecture("x86_64"),
					host("master-1").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-1").customDeploy().architecture("x86_64"),
					host("arbiter-0").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret-arbiter").consumerRef("machine-2").customDeploy().architecture("x86_64"),
				).build(),
		},
		{
			Scenario: "3-hosts-2-masters-1-arbiter-norole",
			Machines: machines(
				machine("machine-0"),
				machine("machine-1")),
			ArbiterMachines: machines(
				machine("machine-2")),
			Config: config().withArbiter(1).withControlPlane(2).hosts(
				hostType("master-0").bmc("usr0", "pwd0"),
				hostType("master-1").bmc("usr1", "pwd1"),
				hostType("arbiter-0").bmc("usr2", "pwd2")).
				build(),

			ExpectedSetting: settings().
				secrets(
					secret("master-0-bmc-secret").creds("usr0", "pwd0"),
					secret("master-1-bmc-secret").creds("usr1", "pwd1"),
					secret("arbiter-0-bmc-secret").creds("usr2", "pwd2"),
				).
				hosts(
					host("master-0").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-0").customDeploy().architecture("x86_64"),
					host("master-1").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-1").customDeploy().architecture("x86_64"),
					host("arbiter-0").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret-arbiter").consumerRef("machine-2").customDeploy().architecture("x86_64"),
				).build(),
		},
		{
			Scenario: "3-hosts-2-masters-1-arbiter-norole-arbiter",
			Machines: machines(
				machine("machine-0"),
				machine("machine-1")),
			ArbiterMachines: machines(
				machine("machine-2")),
			Config: config().withArbiter(1).withControlPlane(2).hosts(
				hostType("master-0").bmc("usr0", "pwd0").role("master"),
				hostType("master-1").bmc("usr1", "pwd1").role("master"),
				hostType("arbiter-0").bmc("usr2", "pwd2")).
				build(),

			ExpectedSetting: settings().
				secrets(
					secret("master-0-bmc-secret").creds("usr0", "pwd0"),
					secret("master-1-bmc-secret").creds("usr1", "pwd1"),
					secret("arbiter-0-bmc-secret").creds("usr2", "pwd2"),
				).
				hosts(
					host("master-0").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-0").customDeploy().architecture("x86_64"),
					host("master-1").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-1").customDeploy().architecture("x86_64"),
					host("arbiter-0").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret-arbiter").consumerRef("machine-2").customDeploy().architecture("x86_64"),
				).build(),
		},
		{
			Scenario: "3-hosts-2-masters-1-arbiter-norole-master",
			Machines: machines(
				machine("machine-0"),
				machine("machine-1")),
			ArbiterMachines: machines(
				machine("machine-2")),
			Config: config().withArbiter(1).withControlPlane(2).hosts(
				hostType("master-0").bmc("usr0", "pwd0"),
				hostType("master-1").bmc("usr1", "pwd1"),
				hostType("arbiter-0").bmc("usr2", "pwd2").role("arbiter")).
				build(),

			ExpectedSetting: settings().
				secrets(
					secret("master-0-bmc-secret").creds("usr0", "pwd0"),
					secret("master-1-bmc-secret").creds("usr1", "pwd1"),
					secret("arbiter-0-bmc-secret").creds("usr2", "pwd2"),
				).
				hosts(
					host("master-0").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-0").customDeploy().architecture("x86_64"),
					host("master-1").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-1").customDeploy().architecture("x86_64"),
					host("arbiter-0").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret-arbiter").consumerRef("machine-2").customDeploy().architecture("x86_64"),
				).build(),
		},
		{
			Scenario: "5-hosts-3-masters-2-arbiter-mixed",
			Machines: machines(
				machine("machine-0"),
				machine("machine-1"),
				machine("machine-2")),
			ArbiterMachines: machines(
				machine("machine-3"),
				machine("machine-4")),
			Config: config().withArbiter(2).withControlPlane(3).hosts(
				hostType("master-0").bmc("usr0", "pwd0").role("master"),
				hostType("arbiter-0").bmc("usr3", "pwd3").role("arbiter"),
				hostType("master-1").bmc("usr1", "pwd1"),
				hostType("master-2").bmc("usr2", "pwd2"),
				hostType("arbiter-1").bmc("usr4", "pwd4")).
				build(),

			ExpectedSetting: settings().
				secrets(
					secret("master-0-bmc-secret").creds("usr0", "pwd0"),
					secret("master-1-bmc-secret").creds("usr1", "pwd1"),
					secret("master-2-bmc-secret").creds("usr2", "pwd2"),
					secret("arbiter-0-bmc-secret").creds("usr3", "pwd3"),
					secret("arbiter-1-bmc-secret").creds("usr4", "pwd4"),
				).
				hosts(
					host("master-0").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-0").customDeploy().architecture("x86_64"),
					host("master-1").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-1").customDeploy().architecture("x86_64"),
					host("master-2").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-2").customDeploy().architecture("x86_64"),
					host("arbiter-0").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret-arbiter").consumerRef("machine-3").customDeploy().architecture("x86_64"),
					host("arbiter-1").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret-arbiter").consumerRef("machine-4").customDeploy().architecture("x86_64"),
				).build(),
		},
		{
			Scenario: "4-hosts-2-masters-1-arbiter-1-worker-norole",
			Machines: machines(
				machine("machine-0"),
				machine("machine-1")),
			ArbiterMachines: machines(
				machine("machine-2")),
			Config: config().withArbiter(1).withControlPlane(2).hosts(
				hostType("master-0").bmc("usr0", "pwd0"),
				hostType("master-1").bmc("usr1", "pwd1"),
				hostType("arbiter-0").bmc("usr2", "pwd2"),
				hostType("worker-0").bmc("usr3", "pwd3")).
				build(),

			ExpectedSetting: settings().
				secrets(
					secret("master-0-bmc-secret").creds("usr0", "pwd0"),
					secret("master-1-bmc-secret").creds("usr1", "pwd1"),
					secret("worker-0-bmc-secret").creds("usr3", "pwd3"),
					secret("arbiter-0-bmc-secret").creds("usr2", "pwd2"),
				).
				hosts(
					host("master-0").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-0").customDeploy().architecture("x86_64"),
					host("master-1").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret").consumerRef("machine-1").customDeploy().architecture("x86_64"),
					host("worker-0").annotation("baremetalhost.metal3.io/paused", "").architecture("x86_64"),
					host("arbiter-0").label("installer.openshift.io/role", "control-plane").userDataRef("user-data-secret-arbiter").consumerRef("machine-2").customDeploy().architecture("x86_64"),
				).build(),
		},
		{
			Scenario: "architecture-amd64-control-plane",
			Machines: machines(machine("machine-0")),
			Config: config().
				withControlPlaneArchitecture(types.ArchitectureAMD64).
				hosts(hostType("master-0").bmc("usr0", "pwd0").role("master")).
				build(),

			ExpectedSetting: settings().
				secrets(secret("master-0-bmc-secret").creds("usr0", "pwd0")).
				hosts(
					host("master-0").
						label("installer.openshift.io/role", "control-plane").
						userDataRef("user-data-secret").
						consumerRef("machine-0").
						customDeploy().
						architecture("x86_64"),
				).build(),
		},
		{
			Scenario: "architecture-arm64-control-plane",
			Machines: machines(machine("machine-0")),
			Config: config().
				withControlPlaneArchitecture(types.ArchitectureARM64).
				hosts(hostType("master-0").bmc("usr0", "pwd0").role("master")).
				build(),

			ExpectedSetting: settings().
				secrets(secret("master-0-bmc-secret").creds("usr0", "pwd0")).
				hosts(
					host("master-0").
						label("installer.openshift.io/role", "control-plane").
						userDataRef("user-data-secret").
						consumerRef("machine-0").
						customDeploy().
						architecture("aarch64"),
				).build(),
		},
		{
			Scenario: "architecture-amd64-worker",
			Machines: machines(machine("machine-0")),
			Config: config().
				withControlPlaneArchitecture(types.ArchitectureAMD64).
				withComputeArchitecture(types.ArchitectureAMD64).
				hosts(
					hostType("master-0").bmc("usr0", "pwd0").role("master"),
					hostType("worker-0").bmc("usr1", "pwd1").role("worker"),
				).
				build(),

			ExpectedSetting: settings().
				secrets(
					secret("master-0-bmc-secret").creds("usr0", "pwd0"),
					secret("worker-0-bmc-secret").creds("usr1", "pwd1"),
				).
				hosts(
					host("master-0").
						label("installer.openshift.io/role", "control-plane").
						userDataRef("user-data-secret").
						consumerRef("machine-0").
						customDeploy().
						architecture("x86_64"),
					host("worker-0").
						annotation("baremetalhost.metal3.io/paused", "").
						architecture("x86_64"),
				).build(),
		},
		{
			Scenario: "architecture-arm64-worker",
			Machines: machines(machine("machine-0")),
			Config: config().
				withControlPlaneArchitecture(types.ArchitectureARM64).
				withComputeArchitecture(types.ArchitectureARM64).
				hosts(
					hostType("master-0").bmc("usr0", "pwd0").role("master"),
					hostType("worker-0").bmc("usr1", "pwd1").role("worker"),
				).
				build(),

			ExpectedSetting: settings().
				secrets(
					secret("master-0-bmc-secret").creds("usr0", "pwd0"),
					secret("worker-0-bmc-secret").creds("usr1", "pwd1"),
				).
				hosts(
					host("master-0").
						label("installer.openshift.io/role", "control-plane").
						userDataRef("user-data-secret").
						consumerRef("machine-0").
						customDeploy().
						architecture("aarch64"),
					host("worker-0").
						annotation("baremetalhost.metal3.io/paused", "").
						architecture("aarch64"),
				).build(),
		},
		{
			Scenario: "architecture-s390x-control-plane",
			Machines: machines(machine("machine-0")),
			Config: config().
				withControlPlaneArchitecture(types.ArchitectureS390X).
				hosts(hostType("master-0").bmc("usr0", "pwd0").role("master")).
				build(),

			ExpectedSetting: settings().
				secrets(secret("master-0-bmc-secret").creds("usr0", "pwd0")).
				hosts(
					host("master-0").
						label("installer.openshift.io/role", "control-plane").
						userDataRef("user-data-secret").
						consumerRef("machine-0").
						customDeploy().
						architecture("s390x"),
				).build(),
		},
		{
			Scenario: "architecture-ppc64le-control-plane",
			Machines: machines(machine("machine-0")),
			Config: config().
				withControlPlaneArchitecture(types.ArchitecturePPC64LE).
				hosts(hostType("master-0").bmc("usr0", "pwd0").role("master")).
				build(),

			ExpectedSetting: settings().
				secrets(secret("master-0-bmc-secret").creds("usr0", "pwd0")).
				hosts(
					host("master-0").
						label("installer.openshift.io/role", "control-plane").
						userDataRef("user-data-secret").
						consumerRef("machine-0").
						customDeploy().
						architecture("ppc64le"),
				).build(),
		},
		{
			Scenario: "architecture-amd64-arbiter",
			Machines: machines(
				machine("machine-0"),
				machine("machine-1")),
			ArbiterMachines: machines(
				machine("machine-2")),
			Config: config().
				withArbiter(1).
				withArbiterArchitecture(types.ArchitectureAMD64).
				withControlPlane(2).
				withControlPlaneArchitecture(types.ArchitectureAMD64).
				hosts(
					hostType("master-0").bmc("usr0", "pwd0").role("master"),
					hostType("master-1").bmc("usr1", "pwd1").role("master"),
					hostType("arbiter-0").bmc("usr2", "pwd2").role("arbiter"),
				).
				build(),

			ExpectedSetting: settings().
				secrets(
					secret("master-0-bmc-secret").creds("usr0", "pwd0"),
					secret("master-1-bmc-secret").creds("usr1", "pwd1"),
					secret("arbiter-0-bmc-secret").creds("usr2", "pwd2"),
				).
				hosts(
					host("master-0").
						label("installer.openshift.io/role", "control-plane").
						userDataRef("user-data-secret").
						consumerRef("machine-0").
						customDeploy().
						architecture("x86_64"),
					host("master-1").
						label("installer.openshift.io/role", "control-plane").
						userDataRef("user-data-secret").
						consumerRef("machine-1").
						customDeploy().
						architecture("x86_64"),
					host("arbiter-0").
						label("installer.openshift.io/role", "control-plane").
						userDataRef("user-data-secret-arbiter").
						consumerRef("machine-2").
						customDeploy().
						architecture("x86_64"),
				).build(),
		},
		{
			Scenario: "architecture-arm64-arbiter",
			Machines: machines(
				machine("machine-0"),
				machine("machine-1")),
			ArbiterMachines: machines(
				machine("machine-2")),
			Config: config().
				withArbiter(1).
				withArbiterArchitecture(types.ArchitectureARM64).
				withControlPlane(2).
				withControlPlaneArchitecture(types.ArchitectureAMD64).
				hosts(
					hostType("master-0").bmc("usr0", "pwd0").role("master"),
					hostType("master-1").bmc("usr1", "pwd1").role("master"),
					hostType("arbiter-0").bmc("usr2", "pwd2").role("arbiter"),
				).
				build(),

			ExpectedSetting: settings().
				secrets(
					secret("master-0-bmc-secret").creds("usr0", "pwd0"),
					secret("master-1-bmc-secret").creds("usr1", "pwd1"),
					secret("arbiter-0-bmc-secret").creds("usr2", "pwd2"),
				).
				hosts(
					host("master-0").
						label("installer.openshift.io/role", "control-plane").
						userDataRef("user-data-secret").
						consumerRef("machine-0").
						customDeploy().
						architecture("x86_64"),
					host("master-1").
						label("installer.openshift.io/role", "control-plane").
						userDataRef("user-data-secret").
						consumerRef("machine-1").
						customDeploy().
						architecture("x86_64"),
					host("arbiter-0").
						label("installer.openshift.io/role", "control-plane").
						userDataRef("user-data-secret-arbiter").
						consumerRef("machine-2").
						customDeploy().
						architecture("aarch64"),
				).build(),
		},
		{
			Scenario: "architecture-s390x-arbiter",
			Machines: machines(
				machine("machine-0"),
				machine("machine-1")),
			ArbiterMachines: machines(
				machine("machine-2")),
			Config: config().
				withArbiter(1).
				withArbiterArchitecture(types.ArchitectureS390X).
				withControlPlane(2).
				withControlPlaneArchitecture(types.ArchitectureAMD64).
				hosts(
					hostType("master-0").bmc("usr0", "pwd0").role("master"),
					hostType("master-1").bmc("usr1", "pwd1").role("master"),
					hostType("arbiter-0").bmc("usr2", "pwd2").role("arbiter"),
				).
				build(),

			ExpectedSetting: settings().
				secrets(
					secret("master-0-bmc-secret").creds("usr0", "pwd0"),
					secret("master-1-bmc-secret").creds("usr1", "pwd1"),
					secret("arbiter-0-bmc-secret").creds("usr2", "pwd2"),
				).
				hosts(
					host("master-0").
						label("installer.openshift.io/role", "control-plane").
						userDataRef("user-data-secret").
						consumerRef("machine-0").
						customDeploy().
						architecture("x86_64"),
					host("master-1").
						label("installer.openshift.io/role", "control-plane").
						userDataRef("user-data-secret").
						consumerRef("machine-1").
						customDeploy().
						architecture("x86_64"),
					host("arbiter-0").
						label("installer.openshift.io/role", "control-plane").
						userDataRef("user-data-secret-arbiter").
						consumerRef("machine-2").
						customDeploy().
						architecture("s390x"),
				).build(),
		},
		{
			Scenario: "architecture-ppc64le-arbiter",
			Machines: machines(
				machine("machine-0"),
				machine("machine-1")),
			ArbiterMachines: machines(
				machine("machine-2")),
			Config: config().
				withArbiter(1).
				withArbiterArchitecture(types.ArchitecturePPC64LE).
				withControlPlane(2).
				withControlPlaneArchitecture(types.ArchitectureAMD64).
				hosts(
					hostType("master-0").bmc("usr0", "pwd0").role("master"),
					hostType("master-1").bmc("usr1", "pwd1").role("master"),
					hostType("arbiter-0").bmc("usr2", "pwd2").role("arbiter"),
				).
				build(),

			ExpectedSetting: settings().
				secrets(
					secret("master-0-bmc-secret").creds("usr0", "pwd0"),
					secret("master-1-bmc-secret").creds("usr1", "pwd1"),
					secret("arbiter-0-bmc-secret").creds("usr2", "pwd2"),
				).
				hosts(
					host("master-0").
						label("installer.openshift.io/role", "control-plane").
						userDataRef("user-data-secret").
						consumerRef("machine-0").
						customDeploy().
						architecture("x86_64"),
					host("master-1").
						label("installer.openshift.io/role", "control-plane").
						userDataRef("user-data-secret").
						consumerRef("machine-1").
						customDeploy().
						architecture("x86_64"),
					host("arbiter-0").
						label("installer.openshift.io/role", "control-plane").
						userDataRef("user-data-secret-arbiter").
						consumerRef("machine-2").
						customDeploy().
						architecture("ppc64le"),
				).build(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			settings, err := Hosts(tc.Config, tc.Machines, "user-data-secret")

			arbiterSettings, arbiterErr := ArbiterHosts(tc.Config, tc.ArbiterMachines, "user-data-secret-arbiter")
			if tc.ExpectedError == "" {
				assert.Nil(t, arbiterErr)
			}
			if settings != nil && arbiterSettings != nil {
				settings.Hosts = append(settings.Hosts, arbiterSettings.Hosts...)
				settings.Secrets = append(settings.Secrets, arbiterSettings.Secrets...)
			}

			if tc.ExpectedError != "" {
				assert.EqualError(t, err, tc.ExpectedError)
			}

			if tc.ExpectedSetting != nil {
				for i, h := range tc.ExpectedSetting.Hosts {
					assert.Equal(t, h, settings.Hosts[i], fmt.Sprintf("%s and %s are not equal", h.Name, settings.Hosts[i].Name))
				}

				for i, s := range tc.ExpectedSetting.Secrets {
					assert.Equal(t, s, settings.Secrets[i], s.Name, fmt.Sprintf("%s and %s are not equal", s.Name, settings.Secrets[i].Name))
				}

				for i, s := range tc.ExpectedSetting.NetworkConfigSecrets {
					assert.Equal(t, s, settings.NetworkConfigSecrets[i], s.Name, fmt.Sprintf("%s and %s are not equal", s.Name, settings.NetworkConfigSecrets[i].Name))
				}
			}
		})
	}
}

func configHosts(builders ...*hostTypeBuilder) *types.InstallConfig {
	return config().hosts(builders...).build()
}

type installConfigBuilder struct {
	types.InstallConfig
}

func (ib *installConfigBuilder) build() *types.InstallConfig {
	return &ib.InstallConfig
}

func config() *installConfigBuilder {
	return &installConfigBuilder{
		types.InstallConfig{
			Platform: types.Platform{
				BareMetal: &baremetaltypes.Platform{},
			},
			ControlPlane: &types.MachinePool{
				Architecture: types.ArchitectureAMD64,
			},
			Compute: []types.MachinePool{
				{
					Architecture: types.ArchitectureAMD64,
				},
			},
		},
	}
}

func (ib *installConfigBuilder) hosts(builders ...*hostTypeBuilder) *installConfigBuilder {
	ib.Platform.BareMetal.Hosts = []*baremetaltypes.Host{}

	for _, hb := range builders {
		ib.Platform.BareMetal.Hosts = append(ib.Platform.BareMetal.Hosts, hb.build())
	}
	return ib
}

func (ib *installConfigBuilder) withArbiter(count int) *installConfigBuilder {
	replicas := int64(count)
	ib.Arbiter = &types.MachinePool{
		Replicas:     &replicas,
		Architecture: types.ArchitectureAMD64,
	}
	return ib
}

func (ib *installConfigBuilder) withControlPlane(count int) *installConfigBuilder {
	replicas := int64(count)
	if ib.ControlPlane == nil {
		ib.ControlPlane = &types.MachinePool{}
	}
	ib.ControlPlane.Replicas = &replicas
	if ib.ControlPlane.Architecture == "" {
		ib.ControlPlane.Architecture = types.ArchitectureAMD64
	}
	return ib
}

func (ib *installConfigBuilder) withControlPlaneArchitecture(arch types.Architecture) *installConfigBuilder {
	if ib.ControlPlane == nil {
		ib.ControlPlane = &types.MachinePool{}
	}
	ib.ControlPlane.Architecture = arch
	return ib
}

func (ib *installConfigBuilder) withComputeArchitecture(arch types.Architecture) *installConfigBuilder {
	if len(ib.Compute) == 0 {
		ib.Compute = []types.MachinePool{{}}
	}
	ib.Compute[0].Architecture = arch
	return ib
}

func (ib *installConfigBuilder) withArbiterArchitecture(arch types.Architecture) *installConfigBuilder {
	if ib.Arbiter == nil {
		ib.Arbiter = &types.MachinePool{}
	}
	ib.Arbiter.Architecture = arch
	return ib
}

type hostTypeBuilder struct {
	baremetaltypes.Host
}

func (htb *hostTypeBuilder) build() *baremetaltypes.Host {
	return &htb.Host
}

func hostType(name string) *hostTypeBuilder {
	return &hostTypeBuilder{
		Host: baremetaltypes.Host{
			Name:           name,
			BootMACAddress: "c0:ff:ee:ca:fe:00",
			BootMode:       baremetaltypes.UEFI,
			RootDeviceHints: &baremetaltypes.RootDeviceHints{
				DeviceName: "userd_devicename",
			},
		},
	}
}

func (htb *hostTypeBuilder) role(role string) *hostTypeBuilder {
	htb.Role = role
	return htb
}

func (htb *hostTypeBuilder) bmc(user, password string) *hostTypeBuilder {
	htb.BMC = baremetaltypes.BMC{
		Username: user,
		Password: password,
	}
	return htb
}

func (htb *hostTypeBuilder) networkConfig(config string) *hostTypeBuilder {
	yaml.Unmarshal([]byte(config), &htb.NetworkConfig)
	return htb
}

type machineBuilder struct {
	machineapi.Machine
}

func (mb *machineBuilder) build() *machineapi.Machine {
	return &mb.Machine
}

func machine(name string) *machineBuilder {
	return &machineBuilder{
		machineapi.Machine{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Machine",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: "namespace",
			},
		},
	}
}

func machines(builders ...*machineBuilder) []machineapi.Machine {
	m := []machineapi.Machine{}

	for _, mb := range builders {
		m = append(m, *mb.build())
	}
	return m
}

type hostBuilder struct {
	baremetalhost.BareMetalHost
}

func host(name string) *hostBuilder {
	return &hostBuilder{
		baremetalhost.BareMetalHost{

			TypeMeta: metav1.TypeMeta{
				APIVersion: "metal3.io/v1alpha1",
				Kind:       "BareMetalHost",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: "openshift-machine-api",
			},
			Spec: baremetalhost.BareMetalHostSpec{
				BMC: baremetalhost.BMCDetails{
					CredentialsName: name + "-bmc-secret",
				},
				RootDeviceHints: &baremetalhost.RootDeviceHints{
					DeviceName: "userd_devicename",
				},
				BootMode:       "UEFI",
				BootMACAddress: "c0:ff:ee:ca:fe:00",
				Online:         true,
			},
		},
	}
}

func (hb *hostBuilder) build() *baremetalhost.BareMetalHost {
	return &hb.BareMetalHost
}

func (hb *hostBuilder) externallyProvisioned() *hostBuilder {
	hb.Spec.ExternallyProvisioned = true
	return hb
}

func (hb *hostBuilder) customDeploy() *hostBuilder {
	hb.Spec.CustomDeploy = &baremetalhost.CustomDeploy{
		Method: "install_coreos",
	}
	return hb
}

func (hb *hostBuilder) annotation(key, value string) *hostBuilder {
	if hb.Annotations == nil {
		hb.Annotations = map[string]string{}
	}

	hb.Annotations[key] = value
	return hb
}

func (hb *hostBuilder) label(key, value string) *hostBuilder {
	if hb.Labels == nil {
		hb.Labels = map[string]string{}
	}

	hb.Labels[key] = value
	return hb
}

func (hb *hostBuilder) consumerRef(name string) *hostBuilder {
	hb.Spec.ConsumerRef = &corev1.ObjectReference{
		APIVersion: "v1",
		Kind:       "Machine",
		Namespace:  "namespace",
		Name:       name,
	}
	return hb
}

func (hb *hostBuilder) userDataRef(name string) *hostBuilder {
	hb.Spec.UserData = &corev1.SecretReference{Name: name}
	return hb
}

func (hb *hostBuilder) preprovisioningNetworkDataName(secretName string) *hostBuilder {
	hb.Spec.PreprovisioningNetworkDataName = secretName
	return hb
}

func (hb *hostBuilder) architecture(arch string) *hostBuilder {
	hb.Spec.Architecture = arch
	return hb
}

type secretBuilder struct {
	corev1.Secret
}

func secret(name string) *secretBuilder {
	return &secretBuilder{
		corev1.Secret{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Secret",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: "openshift-machine-api",
			},
		},
	}
}

func (sb *secretBuilder) creds(user, password string) *secretBuilder {
	sb.Data = map[string][]byte{
		"username": []byte(user),
		"password": []byte(password),
	}
	return sb
}

func (sb *secretBuilder) nmstate(config string) *secretBuilder {
	sb.Data = map[string][]byte{
		"nmstate": []byte(config),
	}
	return sb
}

func (sb *secretBuilder) build() *corev1.Secret {
	return &sb.Secret
}

type hostSettingsBuilder struct {
	HostSettings
}

func (hsb *hostSettingsBuilder) secrets(builders ...*secretBuilder) *hostSettingsBuilder {
	hsb.Secrets = []corev1.Secret{}
	for _, sb := range builders {
		hsb.Secrets = append(hsb.Secrets, *sb.build())
	}
	return hsb
}

func (hsb *hostSettingsBuilder) networkConfigSecrets(builders ...*secretBuilder) *hostSettingsBuilder {
	hsb.Secrets = []corev1.Secret{}
	for _, sb := range builders {
		hsb.NetworkConfigSecrets = append(hsb.NetworkConfigSecrets, *sb.build())
	}
	return hsb
}

func (hsb *hostSettingsBuilder) hosts(builders ...*hostBuilder) *hostSettingsBuilder {
	hsb.Hosts = []baremetalhost.BareMetalHost{}
	for _, hb := range builders {
		hsb.Hosts = append(hsb.Hosts, *hb.build())
	}
	return hsb
}

func (hsb *hostSettingsBuilder) build() *HostSettings {
	return &hsb.HostSettings
}

func settings() *hostSettingsBuilder {
	return &hostSettingsBuilder{
		HostSettings: HostSettings{},
	}
}
