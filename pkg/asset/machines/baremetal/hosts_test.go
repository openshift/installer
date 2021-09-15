package baremetal

import (
	"testing"

	machineapi "github.com/openshift/api/machine/v1beta1"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/types"
	baremetaltypes "github.com/openshift/installer/pkg/types/baremetal"

	baremetalhost "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
)

func TestHosts(t *testing.T) {
	testCases := []struct {
		Scenario        string
		Machines        []machineapi.Machine
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
				secrets(secret("master-0-bmc-secret").data("usr0", "pwd0")).
				hosts(host("master-0").consumerRef("machine-0").annotation("baremetalhost.metal3.io/paused", "").externallyProvisioned()).build(),
		},
		{
			Scenario: "default-norole",
			Machines: machines(machine("machine-0")),
			Config:   configHosts(hostType("master-0").bmc("usr0", "pwd0")),

			ExpectedSetting: settings().
				secrets(secret("master-0-bmc-secret").data("usr0", "pwd0")).
				hosts(host("master-0").consumerRef("machine-0").annotation("baremetalhost.metal3.io/paused", "").externallyProvisioned()).build(),
		},
		{
			Scenario: "3-hosts-3-machines",
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
					secret("master-0-bmc-secret").data("usr0", "pwd0"),
					secret("master-1-bmc-secret").data("usr1", "pwd1"),
					secret("master-2-bmc-secret").data("usr2", "pwd2")).
				hosts(
					host("master-0").consumerRef("machine-0").annotation("baremetalhost.metal3.io/paused", "").externallyProvisioned(),
					host("master-1").consumerRef("machine-1").annotation("baremetalhost.metal3.io/paused", "").externallyProvisioned(),
					host("master-2").consumerRef("machine-2").annotation("baremetalhost.metal3.io/paused", "").externallyProvisioned()).build(),
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
				hostType("master-3").bmc("usr3", "pwd3").role("worker")),

			ExpectedSetting: settings().
				secrets(
					secret("master-0-bmc-secret").data("usr0", "pwd0"),
					secret("master-1-bmc-secret").data("usr1", "pwd1"),
					secret("master-2-bmc-secret").data("usr2", "pwd2"),
					secret("master-3-bmc-secret").data("usr3", "pwd3")).
				hosts(
					host("master-0").consumerRef("machine-0").annotation("baremetalhost.metal3.io/paused", "").externallyProvisioned(),
					host("master-1").consumerRef("machine-1").annotation("baremetalhost.metal3.io/paused", "").externallyProvisioned(),
					host("master-2").consumerRef("machine-2").annotation("baremetalhost.metal3.io/paused", "").externallyProvisioned(),
					host("master-3")).build(),
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
				hostType("master-3").bmc("usr3", "pwd3")),

			ExpectedSetting: settings().
				secrets(
					secret("master-0-bmc-secret").data("usr0", "pwd0"),
					secret("master-1-bmc-secret").data("usr1", "pwd1"),
					secret("master-2-bmc-secret").data("usr2", "pwd2"),
					secret("master-3-bmc-secret").data("usr3", "pwd3")).
				hosts(
					host("master-0").consumerRef("machine-0").annotation("baremetalhost.metal3.io/paused", "").externallyProvisioned(),
					host("master-1").consumerRef("machine-1").annotation("baremetalhost.metal3.io/paused", "").externallyProvisioned(),
					host("master-2").consumerRef("machine-2").annotation("baremetalhost.metal3.io/paused", "").externallyProvisioned(),
					host("master-3")).build(),
		},
		{
			Scenario: "5-hosts-3-machines-mixed",
			Machines: machines(
				machine("machine-0"),
				machine("machine-1"),
				machine("machine-2")),
			Config: configHosts(
				hostType("master-0").bmc("usr0", "pwd0").role("master"),
				hostType("worker-0").bmc("wrk0", "pwd0").role("worker"),
				hostType("master-1").bmc("usr1", "pwd1").role("master"),
				hostType("worker-1").bmc("wrk1", "pwd1").role("worker"),
				hostType("master-2").bmc("usr2", "pwd2").role("master")),

			ExpectedSetting: settings().
				secrets(
					secret("master-0-bmc-secret").data("usr0", "pwd0"),
					secret("worker-0-bmc-secret").data("wrk0", "pwd0"),
					secret("master-1-bmc-secret").data("usr1", "pwd1"),
					secret("worker-1-bmc-secret").data("wrk1", "pwd1"),
					secret("master-2-bmc-secret").data("usr2", "pwd2")).
				hosts(
					host("master-0").consumerRef("machine-0").annotation("baremetalhost.metal3.io/paused", "").externallyProvisioned(),
					host("worker-0"),
					host("master-1").consumerRef("machine-1").annotation("baremetalhost.metal3.io/paused", "").externallyProvisioned(),
					host("worker-1"),
					host("master-2").consumerRef("machine-2").annotation("baremetalhost.metal3.io/paused", "").externallyProvisioned()).build(),
		},
		{
			Scenario: "5-hosts-3-machines-mixed-norole",
			Machines: machines(
				machine("machine-0"),
				machine("machine-1"),
				machine("machine-2")),
			Config: configHosts(
				hostType("master-0").bmc("usr0", "pwd0"),
				hostType("worker-0").bmc("wrk0", "pwd0").role("worker"),
				hostType("master-1").bmc("usr1", "pwd1").role("master"),
				hostType("worker-1").bmc("wrk1", "pwd1").role("worker"),
				hostType("master-2").bmc("usr2", "pwd2")),

			ExpectedSetting: settings().
				secrets(
					secret("master-0-bmc-secret").data("usr0", "pwd0"),
					secret("worker-0-bmc-secret").data("wrk0", "pwd0"),
					secret("master-1-bmc-secret").data("usr1", "pwd1"),
					secret("worker-1-bmc-secret").data("wrk1", "pwd1"),
					secret("master-2-bmc-secret").data("usr2", "pwd2")).
				hosts(
					host("master-0").consumerRef("machine-0").annotation("baremetalhost.metal3.io/paused", "").externallyProvisioned(),
					host("worker-0"),
					host("master-1").consumerRef("machine-1").annotation("baremetalhost.metal3.io/paused", "").externallyProvisioned(),
					host("worker-1"),
					host("master-2").consumerRef("machine-2").annotation("baremetalhost.metal3.io/paused", "").externallyProvisioned()).build(),
		},
		{
			Scenario: "3-hosts-mixed",
			Machines: machines(
				machine("machine-0")),
			Config: configHosts(
				hostType("server-0").bmc("usr0", "pwd0"),
				hostType("server-1").bmc("usr1", "pwd1"),
				hostType("server-2").bmc("usr2", "pwd2").role("master")),

			ExpectedSetting: settings().
				secrets(
					secret("server-0-bmc-secret").data("usr0", "pwd0"),
					secret("server-1-bmc-secret").data("usr1", "pwd1"),
					secret("server-2-bmc-secret").data("usr2", "pwd2")).
				hosts(
					host("server-0").consumerRef("machine-0").annotation("baremetalhost.metal3.io/paused", "").externallyProvisioned(),
					host("server-1"),
					host("server-2")).build(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			settings, err := Hosts(tc.Config, tc.Machines)

			if tc.ExpectedError != "" {
				assert.EqualError(t, err, tc.ExpectedError)
			}

			if tc.ExpectedSetting != nil {
				assert.Equal(t, tc.ExpectedSetting, settings)
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

func (hb *hostBuilder) annotation(key, value string) *hostBuilder {
	if hb.Annotations == nil {
		hb.Annotations = map[string]string{}
	}

	hb.Annotations[key] = value
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

func (sb *secretBuilder) data(user, password string) *secretBuilder {
	sb.Data = map[string][]byte{
		"username": []byte(user),
		"password": []byte(password),
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
