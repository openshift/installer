package baremetal

import (
	"fmt"

	"github.com/metal3-io/baremetal-operator/pkg/hardware"

	machineapi "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	baremetalhost "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"

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
}

// Hosts returns the HostSettings with details of the hardware being
// used to construct the cluster.
func Hosts(config *types.InstallConfig, machines []machineapi.Machine) (*HostSettings, error) {
	settings := &HostSettings{}

	if config.Platform.BareMetal == nil {
		return nil, fmt.Errorf("no baremetal platform in configuration")
	}

	for i, host := range config.Platform.BareMetal.Hosts {
		// Each host needs a secret to hold the credentials for
		// communicating with the management controller that drives
		// that host.
		secret := corev1.Secret{
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
		settings.Secrets = append(settings.Secrets, secret)

		// Map string 'default' to hardware.DefaultProfileName
		if host.HardwareProfile == "default" {
			host.HardwareProfile = hardware.DefaultProfileName
		}

		newHost := baremetalhost.BareMetalHost{
			TypeMeta: metav1.TypeMeta{
				APIVersion: baremetalhost.SchemeGroupVersion.String(),
				Kind:       "BareMetalHost",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      host.Name,
				Namespace: "openshift-machine-api",
			},
			Spec: baremetalhost.BareMetalHostSpec{
				Online: true,
				BMC: baremetalhost.BMCDetails{
					Address:                        host.BMC.Address,
					CredentialsName:                secret.Name,
					DisableCertificateVerification: host.BMC.DisableCertificateVerification,
				},
				BootMACAddress:  host.BootMACAddress,
				HardwareProfile: host.HardwareProfile,
			},
		}

		// we will always default to UEFI
		if host.BootMode == baremetal.Legacy {
			newHost.Spec.BootMode = baremetalhost.Legacy
		} else {
			newHost.Spec.BootMode = baremetalhost.UEFI
		}

		if i < len(machines) {
			// Setting ExternallyProvisioned to true and adding a
			// ConsumerRef without setting Image associates the host
			// with a machine without triggering provisioning. We only
			// want to do that for control plane hosts. We assume the
			// first known hosts are the control plane and that the
			// hosts are in the same order as the control plane
			// machines.
			newHost.Spec.ExternallyProvisioned = true
			machine := machines[i]
			newHost.Spec.ConsumerRef = &corev1.ObjectReference{
				APIVersion: machine.TypeMeta.APIVersion,
				Kind:       machine.TypeMeta.Kind,
				Namespace:  machine.ObjectMeta.Namespace,
				Name:       machine.ObjectMeta.Name,
			}
			newHost.Spec.Online = true
		}
		settings.Hosts = append(settings.Hosts, newHost)
	}

	return settings, nil
}
