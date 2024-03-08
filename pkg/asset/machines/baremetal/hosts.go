package baremetal

import (
	"fmt"

	baremetalhost "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	hardware "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1/profile"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	machineapi "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
)

// HostSettings hold the information needed to build the manifests to
// register hosts with the cluster.
type HostSettings struct {
	// Hosts are the BareMetalHost objects for the hardware making up
	// the cluster.
	Hosts []baremetalhost.BareMetalHost
	// Secrets holds the credential information for communicating with
	// the management controllers on the hosts.
	Secrets []corev1.Secret
	// NetworkConfigSecrets holds the networking configuration defined
	// on the host.
	NetworkConfigSecrets []corev1.Secret
}

func createNetworkConfigSecret(host *baremetal.Host) (*corev1.Secret, error) {

	yamlNetworkConfig, err := yaml.Marshal(host.NetworkConfig)
	if err != nil {
		return nil, errors.Wrap(err, "Error while creating network config secret")
	}

	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-network-config-secret", host.Name),
			Namespace: "openshift-machine-api",
		},
		Data: map[string][]byte{"nmstate": yamlNetworkConfig},
	}, nil
}

func createSecret(host *baremetal.Host) (*corev1.Secret, baremetalhost.BMCDetails) {
	bmc := baremetalhost.BMCDetails{}
	if host.BMC.Username != "" && host.BMC.Password != "" {
		// Each host needs a secret to hold the credentials for
		// communicating with the management controller that drives
		// that host.
		secret := &corev1.Secret{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Secret",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-bmc-secret", host.Name),
				Namespace: "openshift-machine-api",
			},
			Data: map[string][]byte{
				"username": []byte(host.BMC.Username),
				"password": []byte(host.BMC.Password),
			},
		}
		bmc.Address = host.BMC.Address
		bmc.CredentialsName = secret.Name
		bmc.DisableCertificateVerification = host.BMC.DisableCertificateVerification

		return secret, bmc
	}
	return nil, bmc
}

func createBaremetalHost(host *baremetal.Host, bmc baremetalhost.BMCDetails) baremetalhost.BareMetalHost {

	// Map string 'default' to hardware.DefaultProfileName
	if host.HardwareProfile == "default" {
		host.HardwareProfile = hardware.DefaultProfileName
	}

	newHost := baremetalhost.BareMetalHost{
		TypeMeta: metav1.TypeMeta{
			APIVersion: baremetalhost.GroupVersion.String(),
			Kind:       "BareMetalHost",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      host.Name,
			Namespace: "openshift-machine-api",
		},
		Spec: baremetalhost.BareMetalHostSpec{
			Online:          true,
			BMC:             bmc,
			BootMACAddress:  host.BootMACAddress,
			HardwareProfile: host.HardwareProfile,
			BootMode:        baremetalhost.BootMode(host.BootMode),
			RootDeviceHints: host.RootDeviceHints.MakeCRDHints(),
		},
	}

	return newHost
}

// Hosts returns the HostSettings with details of the hardware being
// used to construct the cluster.
func Hosts(config *types.InstallConfig, machines []machineapi.Machine, userDataSecret string) (*HostSettings, error) {
	settings := &HostSettings{}

	if config.Platform.BareMetal == nil {
		return nil, fmt.Errorf("no baremetal platform in configuration")
	}

	numRequiredMasters := len(machines)
	numMasters := 0
	for _, host := range config.Platform.BareMetal.Hosts {

		secret, bmc := createSecret(host)
		if secret != nil {
			settings.Secrets = append(settings.Secrets, *secret)
		}
		newHost := createBaremetalHost(host, bmc)

		if host.NetworkConfig != nil {
			networkConfigSecret, err := createNetworkConfigSecret(host)
			if err != nil {
				return nil, err
			}
			settings.NetworkConfigSecrets = append(settings.NetworkConfigSecrets, *networkConfigSecret)
			newHost.Spec.PreprovisioningNetworkDataName = networkConfigSecret.Name
		}

		if !host.IsWorker() && numMasters < numRequiredMasters {
			// Setting CustomDeploy early ensures that the
			// corresponding Ironic node gets correctly configured
			// from the beginning.
			newHost.Spec.CustomDeploy = &baremetalhost.CustomDeploy{
				Method: "install_coreos",
			}

			newHost.ObjectMeta.Labels = map[string]string{
				"installer.openshift.io/role": "control-plane",
			}

			// Link the new host to the currently available machine
			machine := machines[numMasters]
			newHost.Spec.ConsumerRef = &corev1.ObjectReference{
				APIVersion: machine.TypeMeta.APIVersion,
				Kind:       machine.TypeMeta.Kind,
				Namespace:  machine.ObjectMeta.Namespace,
				Name:       machine.ObjectMeta.Name,
			}
			newHost.Spec.Online = true

			// userDataSecret carries a reference to the master ignition file
			newHost.Spec.UserData = &corev1.SecretReference{Name: userDataSecret}
			numMasters++
		} else {
			// Pause workers until the real control plane is up.
			newHost.ObjectMeta.Annotations = map[string]string{
				"baremetalhost.metal3.io/paused": "",
			}
		}

		settings.Hosts = append(settings.Hosts, newHost)
	}

	return settings, nil
}
