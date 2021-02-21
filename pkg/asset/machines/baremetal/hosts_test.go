package baremetal

import (
	"fmt"
	"testing"

	machineapi "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/types"
	baremetaltypes "github.com/openshift/installer/pkg/types/baremetal"

	baremetalhost "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
)

func makeConfig(nHosts int) *types.InstallConfig {
	config := &types.InstallConfig{
		Platform: types.Platform{
			BareMetal: &baremetaltypes.Platform{},
		},
	}

	config.Platform.BareMetal.Hosts = make([]*baremetaltypes.Host, nHosts)
	for i := 0; i < nHosts; i++ {
		config.Platform.BareMetal.Hosts[i] = &baremetaltypes.Host{
			Name: fmt.Sprintf("host%d", i),
			BMC: baremetaltypes.BMC{
				Username: fmt.Sprintf("user%d", i),
				Password: fmt.Sprintf("password%d", i),
			},
		}
	}

	return config
}

func makeMachines(n int) []machineapi.Machine {
	results := []machineapi.Machine{}

	for i := 0; i < n; i++ {
		results = append(results, machineapi.Machine{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("machine%d", i),
			},
		})
	}

	return results
}

func TestHosts(t *testing.T) {
	testCases := []struct {
		Scenario    string
		Machines    []machineapi.Machine
		Config      *types.InstallConfig
		Expected    HostSettings
		ExpectError bool
	}{
		{
			Scenario:    "no platform",
			Config:      &types.InstallConfig{},
			ExpectError: true,
		},

		{
			Scenario: "no hosts",
			Config:   makeConfig(0),
			Expected: HostSettings{},
		},

		{
			Scenario: "3 hosts and machines",
			Config:   makeConfig(3),
			Machines: makeMachines(3),
			Expected: HostSettings{
				Hosts: []baremetalhost.BareMetalHost{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "host0",
						},
						Spec: baremetalhost.BareMetalHostSpec{
							Online: true,
							BMC: baremetalhost.BMCDetails{
								CredentialsName: "host0-bmc-secret",
							},
							ConsumerRef: &corev1.ObjectReference{
								Name: "machine0",
							},
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "host1",
						},
						Spec: baremetalhost.BareMetalHostSpec{
							Online: true,
							BMC: baremetalhost.BMCDetails{
								CredentialsName: "host1-bmc-secret",
							},
							ConsumerRef: &corev1.ObjectReference{
								Name: "machine1",
							},
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "host2",
						},
						Spec: baremetalhost.BareMetalHostSpec{
							Online: true,
							BMC: baremetalhost.BMCDetails{
								CredentialsName: "host2-bmc-secret",
							},
							ConsumerRef: &corev1.ObjectReference{
								Name: "machine2",
							},
						},
					},
				},
			},
		},

		{
			Scenario: "4 hosts and 3 machines",
			Config:   makeConfig(4),
			Machines: makeMachines(3),
			Expected: HostSettings{
				Hosts: []baremetalhost.BareMetalHost{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "host0",
						},
						Spec: baremetalhost.BareMetalHostSpec{
							Online: true,
							BMC: baremetalhost.BMCDetails{
								CredentialsName: "host0-bmc-secret",
							},
							ConsumerRef: &corev1.ObjectReference{
								Name: "machine0",
							},
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "host1",
						},
						Spec: baremetalhost.BareMetalHostSpec{
							Online: true,
							BMC: baremetalhost.BMCDetails{
								CredentialsName: "host1-bmc-secret",
							},
							ConsumerRef: &corev1.ObjectReference{
								Name: "machine1",
							},
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "host2",
						},
						Spec: baremetalhost.BareMetalHostSpec{
							Online: true,
							BMC: baremetalhost.BMCDetails{
								CredentialsName: "host2-bmc-secret",
							},
							ConsumerRef: &corev1.ObjectReference{
								Name: "machine2",
							},
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "host3",
						},
						Spec: baremetalhost.BareMetalHostSpec{
							Online: true,
							BMC: baremetalhost.BMCDetails{
								CredentialsName: "host3-bmc-secret",
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			actual, err := Hosts(tc.Config, tc.Machines)

			if tc.ExpectError {
				if err == nil {
					t.Error("expected error but did not get one")
				}
			} else {

				if err != nil {
					t.Error("did not expect error but got one")
					return
				}

				if len(tc.Expected.Hosts) != len(actual.Hosts) {
					t.Errorf("Expected %d hosts but received %d (%v)",
						len(tc.Expected.Hosts),
						len(actual.Hosts),
						actual.Hosts,
					)
					return
				}

				if len(tc.Expected.Hosts) != len(actual.Secrets) {
					t.Errorf("Expected %d secrets to go with hosts, got %d",
						len(tc.Expected.Hosts), len(actual.Secrets))
					return
				}

				// Make sure the first few hosts that correspond to
				// machines have their consumer ref set correctly
				for i := 0; i < len(tc.Machines); i++ {
					if actual.Hosts[i].Spec.ConsumerRef.Name != tc.Expected.Hosts[i].Spec.ConsumerRef.Name {
						t.Errorf("Expected host %d to link to %q but got %q",
							i,
							tc.Expected.Hosts[i].Spec.ConsumerRef.Name,
							actual.Hosts[i].Spec.ConsumerRef.Name,
						)
					}
				}

				// Make sure any extra hosts do not have a consumer set.
				for i := len(actual.Hosts) - 1; i > len(tc.Machines); i-- {
					if actual.Hosts[i].Spec.ConsumerRef != nil {
						t.Errorf("Expected host %d to have no consumer but has %v",
							i, actual.Hosts[i].Spec.ConsumerRef,
						)
					}
				}

				for i, sec := range actual.Secrets {
					expectedName := fmt.Sprintf("host%d-bmc-secret", i)
					if sec.Name != expectedName {
						t.Errorf("Expected secret name %d to be %q but got %q",
							i,
							expectedName,
							sec.Name,
						)
					}
				}

			}
		})
	}
}
