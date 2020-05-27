// +build okd

package ignition

import (
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"path/filepath"

	"github.com/coreos/ignition/config/util"
	igntypes3 "github.com/coreos/ignition/v2/config/v3_1/types"
	mcfgv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/openshift/installer/pkg/asset/openshiftinstall"

	"github.com/openshift/installer/pkg/asset"
	"github.com/vincent-petithory/dataurl"
)

// Config is ignition v3 Config
type Config igntypes3.Config

// Dropin in an abstraction over igntypes2.SystemdDropin
type Dropin struct {
	Name     string
	Contents string
}

// FilesFromAsset creates ignition files for each of the files in the specified
// asset.
func FilesFromAsset(pathPrefix string, username string, mode int, asset asset.WritableAsset) []igntypes3.File {
	var files []igntypes3.File
	for _, f := range asset.Files() {
		files = append(files, FileFromBytes(filepath.Join(pathPrefix, f.Filename), username, mode, f.Data))
	}
	return files
}

// FileFromString creates an ignition-config file with the given contents.
func FileFromString(path string, username string, mode int, contents string) igntypes3.File {
	return FileFromBytes(path, username, mode, []byte(contents))
}

// FileFromBytes creates an ignition-config file with the given contents.
func FileFromBytes(path string, username string, mode int, contents []byte) igntypes3.File {
	contentsString := dataurl.EncodeBytes(contents)
	overwrite := true
	return igntypes3.File{
		Node: igntypes3.Node{
			Path: path,
			User: igntypes3.NodeUser{
				Name: &username,
			},
			Overwrite: &overwrite,
		},
		FileEmbedded1: igntypes3.FileEmbedded1{
			Mode: &mode,
			Contents: igntypes3.Resource{
				Source: &contentsString,
			},
		},
	}
}

// SetAppendToFile sets append flag
func SetAppendToFile(file *igntypes3.File) {
	file.Append = []igntypes3.Resource{file.Contents}
}

// PointerIgnitionConfig generates a config which references the remote config
// served by the machine config server.
func PointerIgnitionConfig(url string, rootCA []byte) *Config {
	rootCAString := dataurl.EncodeBytes(rootCA)
	return &Config{
		Ignition: igntypes3.Ignition{
			Version: igntypes3.MaxVersion.String(),
			Config: igntypes3.IgnitionConfig{
				Merge: []igntypes3.Resource{{
					Source: &url,
				}},
			},
			Security: igntypes3.Security{
				TLS: igntypes3.TLS{
					CertificateAuthorities: []igntypes3.Resource{{
						Source: &rootCAString,
					}},
				},
			},
		},
	}
}

// MarshalOrDie returns a marshalled interface or panics
func MarshalOrDie(input interface{}) []byte {
	bytes, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}
	return bytes
}

// ForAuthorizedKeys creates the MachineConfig to set the authorized key for `core` user.
func ForAuthorizedKeys(key string, role string) *mcfgv1.MachineConfig {
	return &mcfgv1.MachineConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: mcfgv1.SchemeGroupVersion.String(),
			Kind:       "MachineConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("99-%s-ssh", role),
			Labels: map[string]string{
				"machineconfiguration.openshift.io/role": role,
			},
		},
		Spec: mcfgv1.MachineConfigSpec{
			Config: runtime.RawExtension{
				Raw: MarshalOrDie(&igntypes3.Config{
					Ignition: igntypes3.Ignition{
						Version: igntypes3.MaxVersion.String(),
					},
					Passwd: igntypes3.Passwd{
						Users: []igntypes3.PasswdUser{{
							Name: "core", SSHAuthorizedKeys: []igntypes3.SSHAuthorizedKey{igntypes3.SSHAuthorizedKey(key)},
						}},
					},
				}),
			},
		},
	}
}

// ForFIPSEnabled creates the MachineConfig to enable FIPS.
// See also https://github.com/openshift/machine-config-operator/pull/889
func ForFIPSEnabled(role string) *mcfgv1.MachineConfig {
	return &mcfgv1.MachineConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "machineconfiguration.openshift.io/v1",
			Kind:       "MachineConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("99-%s-fips", role),
			Labels: map[string]string{
				"machineconfiguration.openshift.io/role": role,
			},
		},
		Spec: mcfgv1.MachineConfigSpec{
			Config: runtime.RawExtension{
				Raw: MarshalOrDie(&igntypes3.Config{
					Ignition: igntypes3.Ignition{
						Version: igntypes3.MaxVersion.String(),
					},
				}),
			},
			FIPS: true,
		},
	}
}

// ForHyperthreadingDisabled creates the MachineConfig to disable hyperthreading.
// RHCOS ships with pivot.service that uses the `/etc/pivot/kernel-args` to override the kernel arguments for hosts.
func ForHyperthreadingDisabled(role string) *mcfgv1.MachineConfig {
	return &mcfgv1.MachineConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "machineconfiguration.openshift.io/v1",
			Kind:       "MachineConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("99-%s-disable-hyperthreading", role),
			Labels: map[string]string{
				"machineconfiguration.openshift.io/role": role,
			},
		},
		Spec: mcfgv1.MachineConfigSpec{
			Config: runtime.RawExtension{
				Raw: MarshalOrDie(&igntypes3.Config{
					Ignition: igntypes3.Ignition{
						Version: igntypes3.MaxVersion.String(),
					},
					Storage: igntypes3.Storage{
						Files: []igntypes3.File{
							FileFromString("/etc/pivot/kernel-args", "root", 0600, "ADD nosmt"),
						},
					},
				}),
			},
		},
	}
}

// ForMitigationsDisabled creates the MachineConfig to disable mitigatations.
// FCOS uses `/etc/pivot/kernel-args` to override the kernel arguments for hosts during pivot.
func ForMitigationsDisabled(role string) *mcfgv1.MachineConfig {
	return &mcfgv1.MachineConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "machineconfiguration.openshift.io/v1",
			Kind:       "MachineConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("99-%s-disable-mitigations", role),
			Labels: map[string]string{
				"machineconfiguration.openshift.io/role": role,
			},
		},
		Spec: mcfgv1.MachineConfigSpec{
			Config: runtime.RawExtension{
				Raw: MarshalOrDie(&igntypes3.Config{
					Ignition: igntypes3.Ignition{
						Version: igntypes3.MaxVersion.String(),
					},
					Storage: igntypes3.Storage{
						Files: []igntypes3.File{
							FileFromString("/etc/pivot/kernel-args", "root", 0600, "DELETE mitigations=auto,nosmt"),
						},
					},
				}),
			},
		},
	}
}

// InjectInstallInfo adds information about the installer and its invoker as a
// ConfigMap to the provided bootstrap Ignition config.
func InjectInstallInfo(bootstrap []byte) (string, error) {
	config := &igntypes3.Config{}
	if err := json.Unmarshal(bootstrap, &config); err != nil {
		return "", errors.Wrap(err, "failed to unmarshal bootstrap Ignition config")
	}

	cm, err := openshiftinstall.CreateInstallConfigMap("openshift-install")
	if err != nil {
		return "", errors.Wrap(err, "failed to generate openshift-install config")
	}

	config.Storage.Files = append(config.Storage.Files, FileFromString("/opt/openshift/manifests/openshift-install.yaml", "root", 0644, cm))

	ign, err := json.Marshal(config)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal bootstrap Ignition config")
	}

	return string(ign), nil
}

// GenerateMinimalConfig returns a minimal ignition v2 config
func GenerateMinimalConfig() *Config {
	return &Config{
		Ignition: igntypes3.Ignition{
			Version: igntypes3.MaxVersion.String(),
		},
	}
}

// AddSSHKey returns a minimal ignition v2 config
func (c *Config) AddSSHKey(sshKey, bootstrapSSHKeyPair string) {
	c.Passwd.Users = append(
		c.Passwd.Users,
		igntypes3.PasswdUser{
			Name: "core",
			SSHAuthorizedKeys: []igntypes3.SSHAuthorizedKey{
				igntypes3.SSHAuthorizedKey(sshKey),
				igntypes3.SSHAuthorizedKey(bootstrapSSHKeyPair),
			},
		},
	)
}

// AddSystemdUnit appends contents in Ignition config
func (c *Config) AddSystemdUnit(name string, contents string, enabled bool) {
	unit := igntypes3.Unit{
		Name:     name,
		Contents: &contents,
	}
	if enabled {
		unit.Enabled = util.BoolToPtr(true)
	}
	c.Systemd.Units = append(c.Systemd.Units, unit)

}

// AddSystemdDropins appends systemd dropins in the config
func (c *Config) AddSystemdDropins(name string, children []Dropin, enabled bool) {
	dropins := []igntypes3.Dropin{}
	for _, childInfo := range children {

		dropins = append(dropins, igntypes3.Dropin{
			Name:     childInfo.Name,
			Contents: &childInfo.Contents,
		})
	}
	unit := igntypes3.Unit{
		Name:    name,
		Dropins: dropins,
	}
	if enabled {
		unit.Enabled = util.BoolToPtr(true)
	}
	c.Systemd.Units = append(c.Systemd.Units, unit)
}

// ReplaceOrAppend is a function which ensures duplicate files are not added in the file list
func (c *Config) ReplaceOrAppend(file igntypes3.File) {
	for i, f := range c.Storage.Files {
		if f.Node.Path == file.Node.Path {
			c.Storage.Files[i] = file
			return
		}
	}
	c.Storage.Files = append(c.Storage.Files, file)
}

// To allow Ignition to download its config on the bootstrap machine from a location secured by a
// self-signed certificate, we have to provide it a valid custom ca bundle.
// To do so we generate a small ignition config that contains just Security section with the bundle
// and later append it to the main ignition config.
// We can't do it directly in Terraform, because Ignition provider suppors only 2.1 version, but
// Security section was added in 2.2 only.

// GenerateIgnitionShim is used to generate an ignition file that contains a user ca bundle
// in its Security section.
func GenerateIgnitionShim(userCA string, clusterID string, bootstrapConfigURL string, tokenID string) (string, error) {
	fileMode := 420

	// Hostname Config
	contents := fmt.Sprintf("%s-bootstrap", clusterID)

	contentsString := dataurl.EncodeBytes([]byte(contents))

	hostnameConfigFile := igntypes3.File{
		Node: igntypes3.Node{
			Path: "/etc/hostname",
		},
		FileEmbedded1: igntypes3.FileEmbedded1{
			Mode: &fileMode,
			Contents: igntypes3.Resource{
				Source: &contentsString,
			},
		},
	}

	userCAString := dataurl.EncodeBytes([]byte(userCA))

	// Openstack Ca Cert file
	openstackCAFile := igntypes3.File{
		Node: igntypes3.Node{
			Path: "/opt/openshift/tls/cloud-ca-cert.pem",
		},
		FileEmbedded1: igntypes3.FileEmbedded1{
			Mode: &fileMode,
			Contents: igntypes3.Resource{
				Source: &userCAString,
			},
		},
	}

	security := igntypes3.Security{}
	if userCA != "" {
		carefs := []igntypes3.Resource{}
		rest := []byte(userCA)

		for {
			var block *pem.Block
			block, rest = pem.Decode(rest)
			if block == nil {
				return "", fmt.Errorf("unable to parse certificate, please check the cacert section of clouds.yaml")
			}

			blockString := dataurl.EncodeBytes(pem.EncodeToMemory(block))
			carefs = append(carefs, igntypes3.Resource{Source: &blockString})

			if len(rest) == 0 {
				break
			}
		}

		security = igntypes3.Security{
			TLS: igntypes3.TLS{
				CertificateAuthorities: carefs,
			},
		}
	}

	headers := []igntypes3.HTTPHeader{
		{
			Name:  "X-Auth-Token",
			Value: &tokenID,
		},
	}

	ign := igntypes3.Config{
		Ignition: igntypes3.Ignition{
			Version:  igntypes3.MaxVersion.String(),
			Security: security,
			Config: igntypes3.IgnitionConfig{
				Merge: []igntypes3.Resource{
					{
						Source:      &bootstrapConfigURL,
						HTTPHeaders: headers,
					},
				},
			},
		},
		Storage: igntypes3.Storage{
			Files: []igntypes3.File{
				hostnameConfigFile,
				openstackCAFile,
			},
		},
	}

	data, err := json.Marshal(ign)
	if err != nil {
		return "", err
	}

	// Check the size of the base64-rendered ignition shim isn't to big for nova
	// https://docs.openstack.org/nova/latest/user/metadata.html#user-data
	if len(base64.StdEncoding.EncodeToString(data)) > 65535 {
		return "", fmt.Errorf("rendered bootstrap ignition shim exceeds the 64KB limit for nova user data -- try reducing the size of your CA cert bundle")
	}

	return string(data), nil
}
